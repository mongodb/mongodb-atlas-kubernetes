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

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasDeploymentReconciler reconciles an AtlasDeployment object
type AtlasDeploymentReconciler struct {
	watch.ResourceWatcher
	Client           client.Client
	Log              *zap.SugaredLogger
	Scheme           *runtime.Scheme
	AtlasDomain      string
	GlobalAPISecret  client.ObjectKey
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
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

	if shouldSkip := customresource.ReconciliationShouldBeSkipped(deployment); shouldSkip {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDeployment reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", deployment.Spec)
		if !deployment.GetDeletionTimestamp().IsZero() {
			err := r.removeDeletionFinalizer(context, deployment)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	ctx := customresource.MarkReconciliationStarted(r.Client, deployment, log)
	log.Infow("-> Starting AtlasDeployment reconciliation", "spec", deployment.Spec, "status", deployment.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, deployment)

	resourceVersionIsValid := customresource.ValidateResourceVersion(ctx, deployment, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("deployment validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if err := validate.DeploymentSpec(deployment.Spec); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	ctx.SetConditionTrue(status.ValidationSucceeded)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(context, deployment, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Client = atlasClient

	// Allow users to specify M0/M2/M5 deployments without providing TENANT for Normal and Serverless deployments
	r.verifyNonTenantCase(deployment)

	if deployment.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
			err = r.Client.Get(context, kube.ObjectKeyFromObject(deployment), deployment)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				return result.ReconcileResult(), nil
			}
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			if err = r.Client.Update(context, deployment); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to add finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
	}

	if !deployment.GetDeletionTimestamp().IsZero() {
		if customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
			if customresource.ResourceShouldBeLeftInAtlas(deployment) {
				log.Infof("Not removing Atlas Deployment from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
			} else {
				if err = r.deleteDeploymentFromAtlas(context, deployment, project, log); err != nil {
					log.Errorf("Failed to remove deployment from Atlas: %s", err)
					result = workflow.Terminate(workflow.Internal, err.Error())
					ctx.SetConditionFromResult(status.DeploymentReadyType, result)
					return result.ReconcileResult(), nil
				}
			}
			err = r.removeDeletionFinalizer(context, deployment)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		} else {
			return result.ReconcileResult(), nil
		}
	}

	handleDeployment := r.selectDeploymentHandler(deployment)
	if result, _ := handleDeployment(ctx, project, deployment, req); !result.IsOk() {
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	if !deployment.IsServerless() {
		if result := r.handleAdvancedOptions(ctx, project, deployment); !result.IsOk() {
			ctx.SetConditionFromResult(status.DeploymentReadyType, result)
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
	if deployment.IsAdvancedDeployment() {
		return r.handleAdvancedDeployment
	}
	if deployment.IsServerless() {
		return r.handleServerlessInstance
	}
	return r.handleRegularDeployment
}

// handleAdvancedDeployment ensures the state of the deployment using the Advanced Deployment API
func (r *AtlasDeploymentReconciler) handleAdvancedDeployment(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error) {
	c, result := r.ensureAdvancedDeploymentState(ctx, project, deployment)
	if c != nil && c.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasDeploymentStateNameOption(c.StateName))
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

	ctx.EnsureStatusOption(status.AtlasDeploymentReplicaSet(replicaSetStatus))

	backupEnabled := false
	if c.BackupEnabled != nil {
		backupEnabled = *c.BackupEnabled
	}

	if err := r.ensureBackupScheduleAndPolicy(
		context.Background(),
		ctx, project.ID(),
		deployment,
		backupEnabled,
		req.NamespacedName,
	); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(ctx, project, c.Name, c.ConnectionStrings, deployment); !csResult.IsOk() {
		return csResult, nil
	}

	ctx.
		SetConditionTrue(status.DeploymentReadyType).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(c.ConnectionStrings))

	ctx.SetConditionTrue(status.ReadyType)
	return result, nil
}

// handleServerlessInstance ensures the state of the serverless instance using the serverless API
func (r *AtlasDeploymentReconciler) handleServerlessInstance(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error) {
	c, result := ensureServerlessInstanceState(ctx, project, deployment.Spec.ServerlessSpec)
	return r.ensureConnectionSecretsAndSetStatusOptions(ctx, project, deployment, result, c)
}

// handleRegularDeployment ensures the state of the deployment using the Regular Deployment API
func (r *AtlasDeploymentReconciler) handleRegularDeployment(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error) {
	atlasDeployment, result := ensureDeploymentState(ctx, project, deployment)
	if atlasDeployment != nil && atlasDeployment.StateName != "" {
		ctx.EnsureStatusOption(status.AtlasDeploymentStateNameOption(atlasDeployment.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	replicaSetStatus := make([]status.ReplicaSet, 0, len(deployment.Spec.DeploymentSpec.ReplicationSpecs))
	for _, replicaSet := range atlasDeployment.ReplicationSpecs {
		replicaSetStatus = append(
			replicaSetStatus,
			status.ReplicaSet{
				ID:       replicaSet.ID,
				ZoneName: replicaSet.ZoneName,
			},
		)
	}

	ctx.EnsureStatusOption(status.AtlasDeploymentReplicaSet(replicaSetStatus))

	backupEnabled := false
	providerBackupEnabled := false
	if atlasDeployment.ProviderBackupEnabled != nil {
		providerBackupEnabled = *atlasDeployment.ProviderBackupEnabled
	}
	if atlasDeployment.BackupEnabled != nil {
		backupEnabled = *atlasDeployment.BackupEnabled
	}

	if err := r.ensureBackupScheduleAndPolicy(
		context.Background(),
		ctx,
		project.ID(),
		deployment,
		providerBackupEnabled || backupEnabled,
		req.NamespacedName,
	); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.DeploymentReadyType, result)
		return result, nil
	}

	return r.ensureConnectionSecretsAndSetStatusOptions(ctx, project, deployment, result, atlasDeployment)
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

		workflow.InProgress(workflow.DeploymentAdvancedOptionsReady, "deployment Advanced Configuration Options are being updated")
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
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDeployment{}}, &watch.EventHandlerWithDelete{Controller: r}, r.GlobalPredicates...)
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
func (r *AtlasDeploymentReconciler) Delete(e event.DeleteEvent) error {
	deployment, ok := e.Object.(*mdbv1.AtlasDeployment)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasDeployment{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasdeployment", kube.ObjectKeyFromObject(deployment))

	log.Infow("-> Starting AtlasDeployment deletion", "spec", deployment.Spec)

	context := context.Background()
	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(context, deployment, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	log = log.With("projectID", project.Status.ID, "deploymentName", deployment.GetDeploymentName())

	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(r.Client, deployment.Namespace, project.ID(), deployment.GetDeploymentName())
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

func (r *AtlasDeploymentReconciler) deleteDeploymentFromAtlas(ctx context.Context, deployment *mdbv1.AtlasDeployment, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	deleteDeploymentFunc := atlasClient.Clusters.Delete
	if deployment.Spec.AdvancedDeploymentSpec != nil {
		deleteDeploymentFunc = atlasClient.AdvancedClusters.Delete
	}
	if deployment.IsServerless() {
		deleteDeploymentFunc = atlasClient.ServerlessInstances.Delete
	}

	_, err = deleteDeploymentFunc(ctx, project.Status.ID, deployment.GetDeploymentName())

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

type deploymentHandlerFunc func(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, req reconcile.Request) (workflow.Result, error)
