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

	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

// AtlasDeploymentReconciler reconciles an AtlasDeployment object
type AtlasDeploymentReconciler struct {
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
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

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlassearchindexconfigs,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlassearchindexconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlassearchindexconfigs,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlassearchindexconfigs/status,verbs=get;update;patch

// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasDeploymentReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdeployment", req.NamespacedName)

	deployment := &akov2.AtlasDeployment{}
	result := customresource.PrepareResource(context, r.Client, req, deployment, log)
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

	conditions := akov2.InitCondition(deployment, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, context)
	log.Infow("-> Starting AtlasDeployment reconciliation", "spec", deployment.Spec, "status", deployment.Status)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, deployment)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, deployment, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("deployment validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	project := &akov2.AtlasProject{}
	if result := r.readProjectResource(context, deployment, project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	if err := validate.DeploymentSpec(&deployment.Spec, r.AtlasProvider.IsCloudGov(), project.Spec.RegionUsageRestrictions); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(api.ValidationSucceeded)

	if !r.AtlasProvider.IsResourceSupported(deployment) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasDeployment is not supported by Atlas for government").
			WithoutRetry()
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	atlasClient, orgID, err := r.AtlasProvider.Client(workflowCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.OrgID = orgID
	workflowCtx.Client = atlasClient

	atlasSdkClient, _, err := r.AtlasProvider.SdkClient(workflowCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SdkClient = atlasSdkClient

	// Allow users to specify M0/M2/M5 deployments without providing TENANT for Normal and Serverless deployments
	r.verifyNonTenantCase(deployment)

	// convertedDeployment is either serverless or advanced, deployment must be kept unchanged
	// convertedDeployment is always a separate copy, to avoid changes on it to go back to k8s
	convertedDeployment := deployment.DeepCopy()

	deletionRequest, result := r.handleDeletion(workflowCtx, log, prevResult, project, deployment)
	if deletionRequest {
		return result.ReconcileResult(), nil
	}

	if err := uniqueKey(&convertedDeployment.Spec); err != nil {
		log.Errorw("failed to validate tags", "error", err)
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}

	handleDeployment := r.selectDeploymentHandler(convertedDeployment)
	if result, _ := handleDeployment(workflowCtx, project, convertedDeployment, req); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return r.registerConfigAndReturn(workflowCtx, log, deployment, result), nil
	}

	if !convertedDeployment.IsServerless() {
		if result := r.handleAdvancedOptions(workflowCtx, project, convertedDeployment); !result.IsOk() {
			workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
			return r.registerConfigAndReturn(workflowCtx, log, deployment, result), nil
		}
	}

	return r.registerConfigAndReturn(workflowCtx, log, deployment, workflow.OK()), nil
}

func (r *AtlasDeploymentReconciler) registerConfigAndReturn(
	workflowCtx *workflow.Context,
	log *zap.SugaredLogger,
	deployment *akov2.AtlasDeployment, // this must be the original non converted deployment
	result workflow.Result) ctrl.Result {
	if result.IsOk() || result.IsInProgress() {
		err := customresource.ApplyLastConfigApplied(workflowCtx.Context, deployment, r.Client)
		if err != nil {
			alternateResult := workflow.Terminate(workflow.Internal, err.Error())
			workflowCtx.SetConditionFromResult(api.DeploymentReadyType, alternateResult)
			log.Error(result.GetMessage())

			return result.ReconcileResult()
		}
	}
	return result.ReconcileResult()
}

func (r *AtlasDeploymentReconciler) verifyNonTenantCase(deployment *akov2.AtlasDeployment) {
	var pSettings *akov2.ServerlessProviderSettingsSpec
	var deploymentType string

	if deployment.Spec.ServerlessSpec != nil {
		if deployment.Spec.ServerlessSpec.ProviderSettings == nil {
			return
		}
		pSettings = deployment.Spec.ServerlessSpec.ProviderSettings
		deploymentType = "SERVERLESS"
	}

	modifyProviderSettings(pSettings, deploymentType)
}

func (r *AtlasDeploymentReconciler) handleDeletion(
	workflowCtx *workflow.Context,
	log *zap.SugaredLogger,
	prevResult workflow.Result,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment, // this must be the original non converted deployment
) (bool, workflow.Result) {
	if deployment.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
			err := r.Client.Get(workflowCtx.Context, kube.ObjectKeyFromObject(deployment), deployment)
			if err != nil {
				return true, workflow.Terminate(workflow.Internal, err.Error())
			}
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			if err = r.Client.Update(workflowCtx.Context, deployment); err != nil {
				return true, workflow.Terminate(workflow.Internal, err.Error())
			}
		}

		return false, workflow.OK()
	}

	if !customresource.HaveFinalizer(deployment, customresource.FinalizerLabel) {
		return true, prevResult
	}

	if err := r.cleanupBindings(workflowCtx.Context, deployment); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		log.Errorw("failed to cleanup deployment bindings (backups)", "error", err)
		return true, result
	}

	switch {
	case customresource.IsResourcePolicyKeepOrDefault(deployment, r.ObjectDeletionProtection):
		log.Info("Not removing Atlas deployment from Atlas as per configuration")
	case customresource.IsResourcePolicyKeep(deployment):
		log.Infof("Not removing Atlas deployment from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
	case isTerminationProtectionEnabled(deployment):
		msg := fmt.Sprintf("Termination protection for %s deployment enabled. Deployment in Atlas won't be removed", deployment.GetName())
		log.Info(msg)
		r.EventRecorder.Event(deployment, "Warning", "AtlasDeploymentTermination", msg)
	default:
		if err := r.deleteDeploymentFromAtlas(workflowCtx, log, project, deployment); err != nil {
			log.Errorf("failed to remove deployment from Atlas: %s", err)
			result := workflow.Terminate(workflow.Internal, err.Error())
			workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
			return true, result
		}
	}

	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, deployment, customresource.UnsetFinalizer); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		log.Errorw("failed to remove finalizer", "error", err)
		return true, result
	}

	return true, prevResult
}

func isTerminationProtectionEnabled(deployment *akov2.AtlasDeployment) bool {
	return (deployment.Spec.DeploymentSpec != nil &&
		deployment.Spec.DeploymentSpec.TerminationProtectionEnabled) || (deployment.Spec.ServerlessSpec != nil &&
		deployment.Spec.ServerlessSpec.TerminationProtectionEnabled)
}

func (r *AtlasDeploymentReconciler) cleanupBindings(context context.Context, deployment *akov2.AtlasDeployment) error {
	r.Log.Debug("Cleaning up deployment bindings (backup)")

	return r.garbageCollectBackupResource(context, deployment.GetDeploymentName())
}

func modifyProviderSettings(pSettings *akov2.ServerlessProviderSettingsSpec, deploymentType string) {
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

func (r *AtlasDeploymentReconciler) selectDeploymentHandler(deployment *akov2.AtlasDeployment) deploymentHandlerFunc {
	if deployment.IsServerless() {
		return r.handleServerlessInstance
	}
	return r.handleAdvancedDeployment
}

// handleAdvancedDeployment ensures the state of the deployment using the Advanced Deployment API
func (r *AtlasDeploymentReconciler) handleAdvancedDeployment(
	workflowCtx *workflow.Context,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
	req reconcile.Request) (workflow.Result, error) {
	c, result := r.ensureAdvancedDeploymentState(workflowCtx, project, deployment)
	if c != nil && c.StateName != "" {
		workflowCtx.EnsureStatusOption(status.AtlasDeploymentStateNameOption(c.StateName))
	}

	if !result.IsOk() {
		return result, nil
	}

	replicaSetStatus := make([]status.ReplicaSet, 0, len(deployment.Spec.DeploymentSpec.ReplicationSpecs))
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
		workflowCtx, project.ID(),
		deployment,
		backupEnabled,
	); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result, nil
	}

	if csResult := r.ensureConnectionSecrets(workflowCtx, project, c.Name, c.ConnectionStrings, deployment); !csResult.IsOk() {
		return csResult, nil
	}

	workflowCtx.
		SetConditionTrue(api.DeploymentReadyType).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(c.MongoDBVersion)).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(c.ConnectionStrings))

	workflowCtx.SetConditionTrue(api.ReadyType)
	return result, nil
}

// handleServerlessInstance ensures the state of the serverless instance using the serverless API
func (r *AtlasDeploymentReconciler) handleServerlessInstance(
	workflowCtx *workflow.Context,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
	req reconcile.Request) (workflow.Result, error) {
	c, result := r.ensureServerlessInstanceState(workflowCtx, project, deployment)
	return r.ensureConnectionSecretsAndSetStatusOptions(workflowCtx, project, deployment, result, c)
}

// ensureConnectionSecretsAndSetStatusOptions creates the relevant connection secrets and sets
// status options to the given context. This function can be used for regular deployments and serverless instances
func (r *AtlasDeploymentReconciler) ensureConnectionSecretsAndSetStatusOptions(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
	result workflow.Result,
	d *mongodbatlas.Cluster) (workflow.Result, error) {
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
		SetConditionTrue(api.DeploymentReadyType).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(d.MongoDBVersion)).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(d.ConnectionStrings)).
		EnsureStatusOption(status.AtlasDeploymentMongoURIUpdatedOption(d.MongoURIUpdated))

	ctx.SetConditionTrue(api.ReadyType)
	return result, nil
}

func (r *AtlasDeploymentReconciler) handleAdvancedOptions(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment) workflow.Result {
	if deployment.Spec.ProcessArgs == nil {
		return workflow.OK()
	}

	deploymentName := deployment.GetDeploymentName()
	context := context.Background()
	atlasArgs, _, err := ctx.Client.Clusters.GetProcessArgs(context, project.Status.ID, deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.DeploymentAdvancedOptionsReady, "cannot get process args")
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

func (r *AtlasDeploymentReconciler) readProjectResource(ctx context.Context, deployment *akov2.AtlasDeployment, project *akov2.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, deployment.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDeployment").
		For(&akov2.AtlasDeployment{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&akov2.AtlasBackupSchedule{},
			handler.EnqueueRequestsFromMapFunc(r.findDeploymentsForBackupSchedule),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Watches(
			&akov2.AtlasBackupPolicy{},
			handler.EnqueueRequestsFromMapFunc(r.findDeploymentsForBackupPolicy),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Watches(
			&akov2.AtlasSearchIndexConfig{},
			handler.EnqueueRequestsFromMapFunc(r.findDeploymentsForSearchIndexConfig),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Complete(r)
}

func NewAtlasDeploymentReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasDeploymentReconciler {
	return &AtlasDeploymentReconciler{
		Scheme:                   mgr.GetScheme(),
		Client:                   mgr.GetClient(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasDeployment"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasDeploymentReconciler) findDeploymentsForBackupPolicy(ctx context.Context, obj client.Object) []reconcile.Request {
	backupPolicy, ok := obj.(*akov2.AtlasBackupPolicy)
	if !ok {
		r.Log.Warnf("watching AtlasBackupPolicy but got %T", obj)
		return nil
	}

	backupSchedules := &akov2.AtlasBackupScheduleList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasBackupScheduleByBackupPolicyIndex,
			client.ObjectKeyFromObject(backupPolicy).String(),
		),
	}
	err := r.Client.List(ctx, backupSchedules, listOps)
	if err != nil {
		r.Log.Errorf("failed to list Atlas backup schedules: %e", err)
		return []reconcile.Request{}
	}

	deploymentMap := make(map[string]struct{}, len(backupSchedules.Items))
	deployments := make([]reconcile.Request, 0, len(backupSchedules.Items))
	for i := range backupSchedules.Items {
		deploymentKeys := r.findDeploymentsForBackupSchedule(ctx, &backupSchedules.Items[i])
		for j := range deploymentKeys {
			key := deploymentKeys[j].String()
			if _, found := deploymentMap[key]; !found {
				deployments = append(deployments, deploymentKeys[j])
				deploymentMap[key] = struct{}{}
			}
		}
	}

	return deployments
}

func (r *AtlasDeploymentReconciler) findDeploymentsForBackupSchedule(ctx context.Context, obj client.Object) []reconcile.Request {
	backupSchedule, ok := obj.(*akov2.AtlasBackupSchedule)
	if !ok {
		r.Log.Warnf("watching AtlasBackupSchedule but got %T", obj)
		return nil
	}

	deployments := &akov2.AtlasDeploymentList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDeploymentByBackupScheduleIndex,
			client.ObjectKeyFromObject(backupSchedule).String(),
		),
	}
	err := r.Client.List(ctx, deployments, listOps)
	if err != nil {
		r.Log.Errorf("failed to list Atlas deployments: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(deployments.Items))
	for i := range deployments.Items {
		item := deployments.Items[i]
		requests = append(
			requests,
			reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      item.Name,
					Namespace: item.Namespace,
				},
			},
		)
	}

	return requests
}

func (r *AtlasDeploymentReconciler) findDeploymentsForSearchIndexConfig(ctx context.Context, obj client.Object) []reconcile.Request {
	searchIndexConfig, ok := obj.(*akov2.AtlasSearchIndexConfig)
	if !ok {
		r.Log.Warnf("watching AtlasSearchIndexConfig but got %T", obj)
		return nil
	}

	deployments := &akov2.AtlasDeploymentList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDeploymentBySearchIndexIndex,
			client.ObjectKeyFromObject(searchIndexConfig).String(),
		),
	}
	err := r.Client.List(ctx, deployments, listOps)
	if err != nil {
		r.Log.Errorf("failed to list Atlas search index configs: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(deployments.Items))
	for i := range deployments.Items {
		item := deployments.Items[i]
		requests = append(
			requests,
			reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      item.Name,
					Namespace: item.Namespace,
				},
			},
		)
	}

	return requests
}

// Delete implements a handler for the Delete event.
func (r *AtlasDeploymentReconciler) deleteConnectionStrings(
	context context.Context,
	log *zap.SugaredLogger,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
) error {
	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(context, r.Client, "", project.ID(), deployment.GetDeploymentName())
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
	log *zap.SugaredLogger,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
) error {
	log.Infow("-> Starting AtlasDeployment deletion", "spec", deployment.Spec)

	err := r.deleteConnectionStrings(workflowCtx.Context, log, project, deployment)
	if err != nil {
		return err
	}

	atlasClient := workflowCtx.Client
	if deployment.IsServerless() {
		_, err = atlasClient.ServerlessInstances.Delete(workflowCtx.Context, project.Status.ID, deployment.GetDeploymentName())
	} else {
		_, err = atlasClient.AdvancedClusters.Delete(workflowCtx.Context, project.Status.ID, deployment.GetDeploymentName(), nil)
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

func (r *AtlasDeploymentReconciler) removeDeletionFinalizer(context context.Context, deployment *akov2.AtlasDeployment) error {
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

type deploymentHandlerFunc func(workflowCtx *workflow.Context, project *akov2.AtlasProject, deployment *akov2.AtlasDeployment, req reconcile.Request) (workflow.Result, error)

// Parse through tags and verify that all keys are unique. Return error otherwise.
func uniqueKey(deploymentSpec *akov2.AtlasDeploymentSpec) error {
	store := make(map[string]string)
	var arrTags []*akov2.TagSpec

	if deploymentSpec.DeploymentSpec != nil {
		arrTags = deploymentSpec.DeploymentSpec.Tags
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
