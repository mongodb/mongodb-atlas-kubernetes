/*
Copyright 2020 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/handler"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasDeploymentReconciler reconciles an AtlasDeployment object
type AtlasDeploymentReconciler struct {
	watch.ResourceWatcher
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	AtlasDomain                 string
	GlobalAPISecret             client.ObjectKey
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupschedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupschedules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackupschedules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackupschedules/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackuppolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackuppolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackuppolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackuppolicies/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasDeploymentReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdeployment", req.NamespacedName)

	deployment := &mdbv1.AtlasDeployment{}
	result := customresource.PrepareResource(r.Client, req, deployment, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	prevResult := result

	if shouldSkip := customresource.ReconciliationShouldBeSkipped(deployment); shouldSkip {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDeployment reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", deployment.Spec)
		if !deployment.GetDeletionTimestamp().IsZero() {
			err := r.removeDeletionFinalizer(context, deployment)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	workflowCtx := customresource.MarkReconciliationStarted(r.Client, deployment, log)
	log.Infow("-> Starting AtlasDeployment reconciliation", "spec", deployment.Spec, "status", deployment.Status)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, deployment)
		r.EnsureMultiplesResourcesAreWatched(req.NamespacedName, log, workflowCtx.ListResourcesToWatch()...)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, deployment, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("deployment validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if err := validate.DeploymentSpec(deployment.Spec); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(status.ValidationSucceeded)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(context, deployment, project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.Client = atlasClient

	// Allow users to specify M0/M2/M5 deployments without providing TENANT for Normal and Serverless deployments
	r.verifyNonTenantCase(deployment)

	if result := r.checkDeploymentIsManaged(workflowCtx, context, log, project, deployment); !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	deletionRequest, result := r.handleDeletion(workflowCtx, context, log, prevResult, project, deployment)
	if deletionRequest {
		return result.ReconcileResult(), nil
	}

	err = customresource.ApplyLastConfigApplied(context, deployment, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	if deployment.IsLegacyDeployment() {
		if err := ConvertLegacyDeployment(&deployment.Spec); err != nil {
			result = workflow.Terminate(workflow.Internal, err.Error())
			log.Errorw("failed to convert legacy deployment", "error", err)
			return result.ReconcileResult(), nil
		}
		deployment.Spec.DeploymentSpec = nil
	}

	if err := uniqueKey(&deployment.Spec); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		log.Errorw("failed to validate tags", "error", err)
		return result.ReconcileResult(), nil
	}

	handleDeployment := r.selectDeploymentHandler(deployment)
	if result, _ := handleDeployment(workflowCtx, project, deployment, req); !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	if !deployment.IsServerless() {
		if result := r.handleAdvancedOptions(workflowCtx, project, deployment); !result.IsOk() {
			workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
			return result.ReconcileResult(), nil
		}
	}
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasDeploymentReconciler) verifyNonTenantCase(deployment *mdbv1.AtlasDeployment) {
	var pSettings *mdbv1.ProviderSettingsSpec
	var deploymentType string
	if deployment.Spec.DeploymentSpec != nil {
		if deployment.Spec.DeploymentSpec.ProviderSettings == nil {
			return
		}
		pSettings = deployment.Spec.DeploymentSpec.ProviderSettings
		deploymentType = "TENANT"
	}

	if deployment.Spec.ServerlessSpec != nil {
		if deployment.Spec.ServerlessSpec.ProviderSettings == nil {
			return
		}
		pSettings = deployment.Spec.ServerlessSpec.ProviderSettings
		deploymentType = "SERVERLESS"
	}

	modifyProviderSettings(pSettings, deploymentType)
}

func (r *AtlasDeploymentReconciler) checkDeploymentIsManaged(
	workflowCtx *workflow.Context,
	context context.Context,
	log *zap.SugaredLogger,
	project *mdbv1.AtlasProject,
	deployment *mdbv1.AtlasDeployment,
) workflow.Result {
	dply := deployment
	if deployment.IsLegacyDeployment() {
		dply = deployment.DeepCopy()
		if err := ConvertLegacyDeployment(&dply.Spec); err != nil {
			result := workflow.Terminate(workflow.Internal, err.Error())
			log.Errorw("failed to temporary convert legacy deployment", "error", err)
			return result
		}
		dply.Spec.DeploymentSpec = nil
	}

	owner, err := customresource.IsOwner(
		dply,
		r.ObjectDeletionProtection,
		customresource.IsResourceManagedByOperator,
		managedByAtlas(context, workflowCtx.Client, project.ID(), log),
	)

	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result
	}

	if !owner {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Deployment due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result
	}

	return workflow.OK()
}

func (r *AtlasDeploymentReconciler) handleDeletion(
	workflowCtx *workflow.Context,
	context context.Context,
	log *zap.SugaredLogger,
	prevResult workflow.Result,
	project *mdbv1.AtlasProject,
	deployment *mdbv1.AtlasDeployment,
) (bool, workflow.Result) {
	if deployment.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
			err := r.Client.Get(context, kube.ObjectKeyFromObject(deployment), deployment)
			if err != nil {
				return true, workflow.Terminate(workflow.Internal, err.Error())
			}
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			if err = r.Client.Update(context, deployment); err != nil {
				return true, workflow.Terminate(workflow.Internal, err.Error())
			}
		}
	}

	if !deployment.GetDeletionTimestamp().IsZero() {
		if customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
			isProtected := customresource.IsResourceProtected(deployment, r.ObjectDeletionProtection)
			if isProtected {
				log.Info("Not removing Atlas deployment from Atlas as per configuration")
			} else {
				if customresource.ResourceShouldBeLeftInAtlas(deployment) {
					log.Infof("Not removing Atlas Deployment from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
				} else {
					if err := r.deleteDeploymentFromAtlas(workflowCtx, context, log, project, deployment); err != nil {
						log.Errorf("failed to remove deployment from Atlas: %s", err)
						result := workflow.Terminate(workflow.Internal, err.Error())
						workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
						return true, result
					}
				}
			}
			err := customresource.ManageFinalizer(context, r.Client, deployment, customresource.UnsetFinalizer)
			if err != nil {
				result := workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return true, result
			}
		}
		return true, prevResult
	}
	return false, workflow.OK()
}

func modifyProviderSettings(pSettings *mdbv1.ProviderSettingsSpec, deploymentType string) {
	if pSettings == nil || string(pSettings.ProviderName) == deploymentType {
		return
	}

	switch strings.ToUpper(deploymentType) {
	case "TENANT":
		switch pSettings.InstanceSizeName {
		case "M0", "M2", "M5":
			pSettings.BackingProviderName = string(pSettings.ProviderName)
			pSettings.ProviderName = provider.ProviderName(deploymentType)
		}
	case "SERVERLESS":
		pSettings.BackingProviderName = string(pSettings.ProviderName)
		pSettings.ProviderName = provider.ProviderName(deploymentType)
	}
}

func (r *AtlasDeploymentReconciler) selectDeploymentHandler(deployment *mdbv1.AtlasDeployment) deploymentHandlerFunc {
	if deployment.IsServerless() {
		return r.handleServerlessInstance
	}
	return r.handleAdvancedDeployment
}

// handleAdvancedDeployment ensures the state of the deployment using the Advanced Deployment API
func (r *AtlasDeploymentReconciler) handleAdvancedDeployment(workflowCtx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error) {
	c, result := r.ensureAdvancedDeploymentState(workflowCtx, project, deployment)
	if c != nil && c.StateName != "" {
		workflowCtx.EnsureStatusOption(status.AtlasDeploymentStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	replicaSetStatus := make([]status.ReplicaSet, 0, len(deployment.Spec.AdvancedDeploymentSpec.ReplicationSpecs))
	for _, replicaSet := range c.ReplicationSpecs {
		replicaSetStatus = append(
			replicaSetStatus,
			status.ReplicaSet{
				ID:       replicaSet.ID,
				ZoneName: replicaSet.ZoneName,
			},
		)
	}

	workflowCtx.EnsureStatusOption(status.AtlasDeploymentReplicaSet(replicaSetStatus))

	backupEnabled := false
	if c.BackupEnabled != nil {
		backupEnabled = *c.BackupEnabled
	}

	if err := r.ensureBackupScheduleAndPolicy(
		context.Background(),
		workflowCtx, project.ID(),
		deployment,
		backupEnabled,
	); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(workflowCtx, project, c.Name, c.ConnectionStrings, deployment); !csResult.IsOk() {
		return csResult, nil
	}

	workflowCtx.
		SetConditionTrue(status.DeploymentReadyType).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(c.ConnectionStrings))

	workflowCtx.SetConditionTrue(status.ReadyType)
	return result, nil
}

// handleServerlessInstance ensures the state of the serverless instance using the serverless API
func (r *AtlasDeploymentReconciler) handleServerlessInstance(workflowCtx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error) {
	c, result := ensureServerlessInstanceState(workflowCtx, project, deployment.Spec.ServerlessSpec)
	return r.ensureConnectionSecretsAndSetStatusOptions(workflowCtx, project, deployment, result, c)
}

// ensureConnectionSecretsAndSetStatusOptions creates the relevant connection secrets and sets
// status options to the given context. This function can be used for regular deployments and serverless instances
func (r *AtlasDeploymentReconciler) ensureConnectionSecretsAndSetStatusOptions(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, result workflow.Result, d *mongodbatlas.Cluster) (workflow.Result, error) {
	if d != nil && d.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasDeploymentStateNameOption(d.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(ctx, project, d.Name, d.ConnectionStrings, deployment); !csResult.IsOk() {
		return csResult, nil
	}

	ctx.
		SetConditionTrue(status.DeploymentReadyType).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(d.MongoDBVersion)).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(d.ConnectionStrings)).
		EnsureStatusOption(status.AtlasDeploymentMongoURIUpdatedOption(d.MongoURIUpdated))

	ctx.SetConditionTrue(status.ReadyType)
	return result, nil
}

func (r *AtlasDeploymentReconciler) handleAdvancedOptions(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment) workflow.Result {
	deploymentName := deployment.GetDeploymentName()
	context := context.Background()
	atlasArgs, _, err := ctx.Client.Clusters.GetProcessArgs(context, project.Status.ID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.DeploymentAdvancedOptionsReady, "cannot get process args")
	}

	if deployment.Spec.ProcessArgs == nil {
		return workflow.OK()
	}

	if !deployment.Spec.ProcessArgs.IsEqual(atlasArgs) {
		options, err := deployment.Spec.ProcessArgs.ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.DeploymentAdvancedOptionsReady, "cannot convert process args to atlas")
		}

		args, resp, err := ctx.Client.Clusters.UpdateProcessArgs(context, project.Status.ID, deploymentName, options)
		ctx.Log.Debugw("ProcessArgs Update", "args", args, "resp", resp.Body, "err", err)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentAdvancedOptionsReady, "cannot update process args")
		}

		// TODO(helderjs): Revisit the advanced options configuration to check if this condition should exist or not
		// workflow.InProgress(workflow.DeploymentAdvancedOptionsReady, "deployment Advanced Configuration Options are being updated")
	}

	return workflow.OK()
}

func (r *AtlasDeploymentReconciler) readProjectResource(ctx context.Context, deployment *mdbv1.AtlasDeployment, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, deployment.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("AtlasDeployment", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AtlasDeployment & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDeployment{}}, &handler.EnqueueRequestForObject{}, r.GlobalPredicates...)
	if err != nil {
		return err
	}

	// Watch for Backup schedules
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupSchedule{}}, watch.NewBackupScheduleHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	// Watch for Backup policies
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupPolicy{}}, watch.NewBackupPolicyHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	return nil
}

// Delete implements a handler for the Delete event.
func (r *AtlasDeploymentReconciler) deleteConnectionStrings(
	context context.Context,
	log *zap.SugaredLogger,
	project *mdbv1.AtlasProject,
	deployment *mdbv1.AtlasDeployment,
) error {
	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(r.Client, "", project.ID(), deployment.GetDeploymentName())
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(context, &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
	}

	return nil
}

func (r *AtlasDeploymentReconciler) deleteDeploymentFromAtlas(
	workflowCtx *workflow.Context,
	context context.Context,
	log *zap.SugaredLogger,
	project *mdbv1.AtlasProject,
	deployment *mdbv1.AtlasDeployment,
) error {
	log.Infow("-> Starting AtlasDeployment deletion", "spec", deployment.Spec)

	err := r.deleteConnectionStrings(context, log, project, deployment)
	if err != nil {
		return err
	}

	atlasClient := workflowCtx.Client
	if deployment.IsServerless() {
		_, err = atlasClient.ServerlessInstances.Delete(context, project.Status.ID, deployment.GetDeploymentName())
	} else {
		_, err = atlasClient.AdvancedClusters.Delete(context, project.Status.ID, deployment.GetDeploymentName(), nil)
	}

	var apiError *mongodbatlas.ErrorResponse
	if errors.As(err, &apiError) && apiError.ErrorCode == atlas.ClusterNotFound {
		log.Info("Deployment doesn't exist or is already deleted")
		return nil
	}

	if err != nil {
		log.Errorw("Cannot delete Atlas deployment", "error", err)
		return err
	}

	return nil
}

func (r *AtlasDeploymentReconciler) removeDeletionFinalizer(context context.Context, deployment *mdbv1.AtlasDeployment) error {
	err := r.Client.Get(context, kube.ObjectKeyFromObject(deployment), deployment)
	if err != nil {
		return fmt.Errorf("cannot get AtlasDeployment while adding finalizer: %w", err)
	}

	customresource.UnsetFinalizer(deployment, customresource.FinalizerLabel)
	if err = r.Client.Update(context, deployment); err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", deployment.GetDeploymentName(), err)
	}
	return nil
}

type deploymentHandlerFunc func(workflowCtx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error)

type atlasClusterType int

const (
	Unset atlasClusterType = iota
	Advanced
	Serverless
)

type atlasTypedCluster struct {
	clusterType atlasClusterType
	serverless  *mongodbatlas.Cluster
	advanced    *mongodbatlas.AdvancedCluster
}

func managedByAtlas(ctx context.Context, atlasClient mongodbatlas.Client, projectID string, log *zap.SugaredLogger) customresource.AtlasChecker {
	return func(resource mdbv1.AtlasCustomResource) (bool, error) {
		deployment, ok := resource.(*mdbv1.AtlasDeployment)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasDeployment")
		}

		typedAtlasCluster, err := findTypedAtlasCluster(ctx, atlasClient, projectID, deployment.GetDeploymentName())
		if typedAtlasCluster == nil || err != nil {
			return false, err
		}

		isSame, err := deploymentMatchesSpec(log, typedAtlasCluster, deployment)
		if err != nil {
			return true, err
		}

		return !isSame, nil
	}
}

func findTypedAtlasCluster(ctx context.Context, atlasClient mongodbatlas.Client, projectID, deploymentName string) (*atlasTypedCluster, error) {
	advancedCluster, _, err := atlasClient.AdvancedClusters.Get(ctx, projectID, deploymentName)
	if err == nil {
		return &atlasTypedCluster{clusterType: Advanced, advanced: advancedCluster}, nil
	}
	var apiError *mongodbatlas.ErrorResponse
	if errors.As(err, &apiError) &&
		apiError.ErrorCode != atlas.ClusterNotFound &&
		apiError.ErrorCode != atlas.ServerlessInstanceFromClusterAPI {
		return nil, err
	}
	// if not found, maybe it is a serverless instead
	serverless, _, err := atlasClient.ServerlessInstances.Get(ctx, projectID, deploymentName)
	if err == nil {
		return &atlasTypedCluster{clusterType: Serverless, serverless: serverless}, nil
	}
	if errors.As(err, &apiError) && apiError.ErrorCode == atlas.ServerlessInstanceNotFound {
		return nil, nil
	}
	return nil, err
}

func deploymentMatchesSpec(log *zap.SugaredLogger, atlasSpec *atlasTypedCluster, deployment *mdbv1.AtlasDeployment) (bool, error) {
	if deployment.IsServerless() {
		if atlasSpec.clusterType != Serverless {
			return false, nil
		}
		return serverlessDeploymentMatchesSpec(log, atlasSpec.serverless, deployment.Spec.ServerlessSpec)
	}
	if atlasSpec.clusterType != Advanced {
		return false, nil
	}
	return advancedDeploymentMatchesSpec(log, atlasSpec.advanced, deployment.Spec.AdvancedDeploymentSpec)
}

func serverlessDeploymentMatchesSpec(log *zap.SugaredLogger, atlasSpec *mongodbatlas.Cluster, operatorSpec *mdbv1.ServerlessSpec) (bool, error) {
	clusterMerged := mongodbatlas.Cluster{}
	if err := compat.JSONCopy(&clusterMerged, atlasSpec); err != nil {
		return false, err
	}

	if err := compat.JSONCopy(&clusterMerged, operatorSpec); err != nil {
		return false, err
	}

	d := cmp.Diff(atlasSpec, &clusterMerged, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Serverless deployment differs from spec: %s", d)
	}

	return d == "", nil
}

func advancedDeploymentMatchesSpec(log *zap.SugaredLogger, atlasSpec *mongodbatlas.AdvancedCluster, operatorSpec *mdbv1.AdvancedDeploymentSpec) (bool, error) {
	clusterMerged := mongodbatlas.AdvancedCluster{}
	if err := compat.JSONCopy(&clusterMerged, atlasSpec); err != nil {
		return false, err
	}

	if err := compat.JSONCopy(&clusterMerged, operatorSpec); err != nil {
		return false, err
	}

	d := cmp.Diff(atlasSpec, &clusterMerged, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Advanced deployment differs from spec: %s", d)
	}

	return d == "", nil
}

// Parse through tags and verfiy that all keys are unique. Return error otherwise.
func uniqueKey(deploymentSpec *mdbv1.AtlasDeploymentSpec) error {
	store := make(map[string]string)
	var arrTags []*mdbv1.TagSpec

	if deploymentSpec.AdvancedDeploymentSpec != nil {
		arrTags = deploymentSpec.AdvancedDeploymentSpec.Tags
	} else {
		arrTags = deploymentSpec.ServerlessSpec.Tags
	}
	for _, currTag := range arrTags {
		if store[currTag.Key] == "" {
			store[currTag.Key] = currTag.Value
		} else {
			return errors.New("duplicate keys found in tags, this is forbidden")
		}
	}
	return nil
}
