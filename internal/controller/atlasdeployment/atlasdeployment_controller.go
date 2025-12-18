// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

// AtlasDeploymentReconciler reconciles an AtlasDeployment object
type AtlasDeploymentReconciler struct {
	reconciler.AtlasReconciler
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	independentSyncPeriod       time.Duration
	maxConcurrentReconciles     int
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

func (r *AtlasDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdeployment", req.NamespacedName)

	atlasDeployment := &akov2.AtlasDeployment{}
	result := customresource.PrepareResource(ctx, r.Client, req, atlasDeployment, log)
	if !result.IsOk() {
		return result.ReconcileResult()
	}

	if shouldSkip := customresource.ReconciliationShouldBeSkipped(atlasDeployment); shouldSkip {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDeployment reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", atlasDeployment.Spec)
		if !atlasDeployment.GetDeletionTimestamp().IsZero() {
			err := r.removeDeletionFinalizer(ctx, atlasDeployment)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err)
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult()
			}
		}
		return workflow.OK().ReconcileResult()
	}

	conditions := akov2.InitCondition(atlasDeployment, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx, atlasDeployment)
	log.Infow("-> Starting AtlasDeployment reconciliation", "spec", atlasDeployment.Spec, "status", atlasDeployment.Status)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasDeployment)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, atlasDeployment, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("deployment validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult()
	}

	if !r.AtlasProvider.IsResourceSupported(atlasDeployment) {
		result = workflow.Terminate(workflow.AtlasGovUnsupported, errors.New("the AtlasDeployment is not supported by Atlas for government")).
			WithoutRetry()
		workflowCtx.SetConditionFromResult(api.DeploymentReadyType, result)
		return result.ReconcileResult()
	}

	connectionConfig, err := r.ResolveConnectionConfig(workflowCtx.Context, atlasDeployment)
	if err != nil {
		return r.terminate(workflowCtx, workflow.AtlasAPIAccessNotConfigured, err)
	}
	sdkClientSet, err := r.AtlasProvider.SdkClientSet(workflowCtx.Context, connectionConfig.Credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, workflow.AtlasAPIAccessNotConfigured, err)
	}
	workflowCtx.SdkClientSet = sdkClientSet
	projectService := project.NewProjectAPIService(sdkClientSet.SdkClient20250312011.ProjectsApi)
	deploymentService := deployment.NewAtlasDeployments(sdkClientSet.SdkClient20250312011.ClustersApi, sdkClientSet.SdkClient20250312011.GlobalClustersApi, sdkClientSet.SdkClient20250312011.FlexClustersApi, r.AtlasProvider.IsCloudGov())
	atlasProject, err := r.ResolveProject(workflowCtx.Context, sdkClientSet.SdkClient20250312011, atlasDeployment)
	if err != nil {
		return r.terminate(workflowCtx, workflow.AtlasAPIAccessNotConfigured, err)
	}

	if err := validate.AtlasDeployment(atlasDeployment); err != nil {
		result = workflow.Terminate(workflow.Internal, err)
		workflowCtx.SetConditionFromResult(api.ValidationSucceeded, result)
		return result.ReconcileResult()
	}
	workflowCtx.SetConditionTrue(api.ValidationSucceeded)

	deploymentInAKO := deployment.NewDeployment(atlasProject.ID, atlasDeployment)

	if ok, notificationReason, notificationMsg := deploymentInAKO.Notifications(); ok {
		// emit Log and event
		r.Log.Log(zapcore.WarnLevel, notificationMsg)
		r.EventRecorder.Event(deploymentInAKO.GetCustomResource(), corev1.EventTypeWarning, notificationReason, notificationMsg)
	}
	deploymentInAtlas, err := deploymentService.GetDeployment(workflowCtx.Context, atlasProject.ID, atlasDeployment)

	if err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	existsInAtlas := deploymentInAtlas != nil
	if !atlasDeployment.GetDeletionTimestamp().IsZero() {
		if existsInAtlas {
			return r.delete(workflowCtx, deploymentService, deploymentInAKO, deploymentInAtlas)
		}
		return r.unmanage(workflowCtx, deploymentInAKO)
	}

	switch {
	case atlasDeployment.IsServerless(), atlasDeployment.IsFlex():
		return r.handleFlexInstance(workflowCtx, projectService, deploymentService, deploymentInAKO, deploymentInAtlas)

	case atlasDeployment.IsAdvancedDeployment():
		return r.handleAdvancedDeployment(workflowCtx, projectService, deploymentService, deploymentInAKO, deploymentInAtlas)
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasDeploymentReconciler) delete(
	ctx *workflow.Context,
	deploymentService deployment.AtlasDeploymentsService,
	deploymentInAKO deployment.Deployment, // this must be the original non converted deployment
	deploymentInAtlas deployment.Deployment, // this must be the original non converted deployment
) (ctrl.Result, error) {
	if err := r.cleanupBindings(ctx.Context, deploymentInAKO); err != nil {
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("failed to cleanup deployment bindings (backups): %w", err))
	}

	switch {
	case customresource.IsResourcePolicyKeepOrDefault(deploymentInAKO.GetCustomResource(), r.ObjectDeletionProtection):
		ctx.Log.Info("Not removing Atlas deployment from Atlas as per configuration")
	case customresource.IsResourcePolicyKeep(deploymentInAKO.GetCustomResource()):
		ctx.Log.Infof("Not removing Atlas deployment from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
	case isTerminationProtectionEnabled(deploymentInAKO.GetCustomResource()):
		msg := fmt.Sprintf("Termination protection for %s deployment enabled. Deployment in Atlas won't be removed", deploymentInAKO.GetName())
		ctx.Log.Info(msg)
		r.EventRecorder.Event(deploymentInAKO.GetCustomResource(), "Warning", "AtlasDeploymentTermination", msg)
	default:
		if err := r.deleteDeploymentFromAtlas(ctx, deploymentService, deploymentInAKO, deploymentInAtlas); err != nil {
			return r.terminate(ctx, workflow.Internal, fmt.Errorf("failed to remove deployment from Atlas: %w", err))
		}
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, deploymentInAKO.GetCustomResource(), customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("failed to remove finalizer: %w", err))
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasDeploymentReconciler) cleanupBindings(context context.Context, deployment deployment.Deployment) error {
	r.Log.Debug("Cleaning up deployment bindings (backup)")

	return r.garbageCollectBackupResource(context, deployment.GetName())
}

func isTerminationProtectionEnabled(deployment *akov2.AtlasDeployment) bool {
	return (deployment.Spec.DeploymentSpec != nil &&
		deployment.Spec.DeploymentSpec.TerminationProtectionEnabled) || (deployment.Spec.ServerlessSpec != nil &&
		deployment.Spec.ServerlessSpec.TerminationProtectionEnabled)
}

func (r *AtlasDeploymentReconciler) deleteDeploymentFromAtlas(
	ctx *workflow.Context,
	deploymentService deployment.AtlasDeploymentsService,
	deploymentInAKO deployment.Deployment,
	deploymentInAtlas deployment.Deployment,
) error {
	ctx.Log.Infow("-> Starting AtlasDeployment deletion", "spec", deploymentInAKO)

	err := r.deleteConnectionStrings(ctx, deploymentInAKO)
	if err != nil {
		return err
	}

	err = deploymentService.DeleteDeployment(ctx.Context, deploymentInAtlas)
	if err != nil {
		ctx.Log.Errorw("Cannot delete Atlas deployment", "error", err)
		return err
	}

	return nil
}

func (r *AtlasDeploymentReconciler) deleteConnectionStrings(ctx *workflow.Context, deployment deployment.Deployment) error {
	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := secretservice.ListByDeploymentName(ctx.Context, r.Client, "", deployment.GetProjectID(), deployment.GetName())
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(ctx.Context, &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			ctx.Log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
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

type transitionFn func(reason workflow.ConditionReason) (ctrl.Result, error)

func (r *AtlasDeploymentReconciler) transitionFromLegacy(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, projectID string, atlasDeployment *akov2.AtlasDeployment, err error) transitionFn {
	return func(reason workflow.ConditionReason) (ctrl.Result, error) {
		if err != nil {
			return r.terminate(ctx, reason, err)
		}

		deploymentInAtlas, err := deploymentService.GetDeployment(ctx.Context, projectID, atlasDeployment)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.inProgress(ctx, atlasDeployment, deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
	}
}

func (r *AtlasDeploymentReconciler) transitionFromResult(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, projectID string, atlasDeployment *akov2.AtlasDeployment, result workflow.DeprecatedResult) transitionFn {
	if result.IsInProgress() {
		return func(reason workflow.ConditionReason) (ctrl.Result, error) {
			deploymentInAtlas, err := deploymentService.GetDeployment(ctx.Context, projectID, atlasDeployment)
			if err != nil {
				return r.terminate(ctx, workflow.Internal, err)
			}

			return r.inProgress(ctx, atlasDeployment, deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
		}
	}

	if !result.IsOk() {
		return func(reason workflow.ConditionReason) (ctrl.Result, error) {
			return r.terminate(ctx, reason, errors.New(result.GetMessage()))
		}
	}

	return nil
}

func (r *AtlasDeploymentReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err)
	ctx.SetConditionFromResult(api.DeploymentReadyType, terminated)

	return terminated.ReconcileResult()
}

func (r *AtlasDeploymentReconciler) inProgress(ctx *workflow.Context, atlasDeployment *akov2.AtlasDeployment, deploymentInAtlas deployment.Deployment, reason workflow.ConditionReason, msg string) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasDeployment, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	result := workflow.InProgress(reason, msg)
	ctx.SetConditionFromResult(api.DeploymentReadyType, result).
		EnsureStatusOption(status.AtlasDeploymentStateNameOption(deploymentInAtlas.GetState())).
		EnsureStatusOption(status.AtlasDeploymentReplicaSet(deploymentInAtlas.GetReplicaSet())).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(deploymentInAtlas.GetMongoDBVersion()))

	return result.ReconcileResult()
}

func (r *AtlasDeploymentReconciler) ready(ctx *workflow.Context, deploymentInAKO, deploymentInAtlas deployment.Deployment) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, deploymentInAKO.GetCustomResource(), customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	_, _, notificationMsg := deploymentInAKO.Notifications()

	ctx.SetConditionTrue(api.DeploymentReadyType).
		SetConditionTrueMsg(api.ReadyType, notificationMsg).
		EnsureStatusOption(status.AtlasDeploymentStateNameOption(deploymentInAtlas.GetState())).
		EnsureStatusOption(status.AtlasDeploymentReplicaSet(deploymentInAtlas.GetReplicaSet())).
		EnsureStatusOption(status.AtlasDeploymentMongoDBVersionOption(deploymentInAtlas.GetMongoDBVersion())).
		EnsureStatusOption(status.AtlasDeploymentConnectionStringsOption(deploymentInAtlas.GetConnection()))

	if deploymentInAKO.GetCustomResource().Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult()
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasDeploymentReconciler) unmanage(ctx *workflow.Context, atlasDeployment deployment.Deployment) (ctrl.Result, error) {
	err := r.removeDeletionFinalizer(ctx.Context, atlasDeployment.GetCustomResource())
	if err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasDeploymentReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasDeployment{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasDeploymentReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDeployment").
		For(r.For()).
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
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.deploymentsForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:             ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation:      pointer.MakePtr(skipNameValidation),
			MaxConcurrentReconciles: r.maxConcurrentReconciles}).
		Complete(r)
}

func NewAtlasDeploymentReconciler(c cluster.Cluster, predicates []predicate.Predicate, atlasProvider atlas.Provider, deletionProtection bool, independentSyncPeriod time.Duration, logger *zap.Logger, globalSecretref client.ObjectKey, maxConcurrentReconciles int) *AtlasDeploymentReconciler {
	suggaredLogger := logger.Named("controllers").Named("AtlasDeployment").Sugar()

	return &AtlasDeploymentReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             suggaredLogger,
			GlobalSecretRef: globalSecretref,
			AtlasProvider:   atlasProvider,
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasDeployment"),
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
		maxConcurrentReconciles:  maxConcurrentReconciles,
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

func (r *AtlasDeploymentReconciler) findDeploymentsForBackupSchedule(ctx context.Context, obj client.Object) []reconcile.Request { //nolint:dupl
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

func (r *AtlasDeploymentReconciler) findDeploymentsForSearchIndexConfig(ctx context.Context, obj client.Object) []reconcile.Request { //nolint:dupl
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

func (r *AtlasDeploymentReconciler) deploymentsForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasDeploymentCredentialsIndex,
		func() *akov2.AtlasDeploymentList { return &akov2.AtlasDeploymentList{} },
		indexer.DeploymentRequests,
		r.Client,
		r.Log,
	)
}

func (r *AtlasDeploymentReconciler) handleDeleted() (ctrl.Result, error) {
	return workflow.OK().ReconcileResult()
}
