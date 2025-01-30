package atlassearchindexconfig

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlassearchindexconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlassearchindexconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlassearchindexconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlassearchindexconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdeployments,verbs=get;list
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdeployments,verbs=get;list
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

type AtlasSearchIndexConfigReconciler struct {
	Client                      client.Client
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	Log                         *zap.SugaredLogger
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

func (r *AtlasSearchIndexConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("AtlasSearchIndexConfig", req.NamespacedName)
	log.Infow("-> Starting AtlasSearchIndexConfig reconciliation")

	atlasSearchIndexConfig := &akov2.AtlasSearchIndexConfig{}
	result := customresource.PrepareResource(ctx, r.Client, req, atlasSearchIndexConfig, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(atlasSearchIndexConfig) {
		return r.skip(ctx, log, atlasSearchIndexConfig), nil
	}

	conditions := api.InitCondition(atlasSearchIndexConfig, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx, atlasSearchIndexConfig)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasSearchIndexConfig)

	isValid := customresource.ValidateResourceVersion(workflowCtx, atlasSearchIndexConfig, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid), nil
	}

	if !r.AtlasProvider.IsResourceSupported(atlasSearchIndexConfig) {
		return r.unsupport(workflowCtx), nil
	}

	deployments := &akov2.AtlasDeploymentList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDeploymentBySearchIndexIndex,
			client.ObjectKeyFromObject(atlasSearchIndexConfig).String(),
		),
	}
	err := r.Client.List(ctx, deployments, listOps)
	if err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err), nil
	}

	if len(deployments.Items) > 0 {
		// set finalizer
		return r.lock(workflowCtx, atlasSearchIndexConfig), nil
	}

	// unset finalizer
	return r.release(workflowCtx, atlasSearchIndexConfig), nil
}

func (r *AtlasSearchIndexConfigReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasSearchIndexConfig{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasSearchIndexConfigReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasSearchIndexConfig").
		For(r.For()).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.findReferencesInAtlasDeployments),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func NewAtlasSearchIndexConfigReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasSearchIndexConfigReconciler {
	return &AtlasSearchIndexConfigReconciler{
		Scheme:                   c.GetScheme(),
		Client:                   c.GetClient(),
		EventRecorder:            c.GetEventRecorderFor("AtlasSearchIndexConfig"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasSearchIndexConfig").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasSearchIndexConfigReconciler) findReferencesInAtlasDeployments(ctx context.Context, obj client.Object) []reconcile.Request {
	deployment, ok := obj.(*akov2.AtlasDeployment)
	if !ok {
		r.Log.Warnf("watching AtlasDeployment but got %T", obj)
		return nil
	}

	if deployment.Spec.DeploymentSpec == nil {
		return nil
	}

	requests := []reconcile.Request{}
	for i := range deployment.Spec.DeploymentSpec.SearchIndexes {
		idx := &deployment.Spec.DeploymentSpec.SearchIndexes[i]
		if idx.Search == nil {
			continue
		}
		requests = append(requests, reconcile.Request{
			NamespacedName: *idx.Search.SearchConfigurationRef.GetObject(deployment.GetNamespace())})
	}
	return requests
}

func (r *AtlasSearchIndexConfigReconciler) skip(ctx context.Context, log *zap.SugaredLogger, searchIndexConfig *akov2.AtlasSearchIndexConfig) ctrl.Result {
	log.Infow(fmt.Sprintf("-> Skipping AtlasSearchIndexConfig reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", searchIndexConfig.Spec)
	if !searchIndexConfig.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, searchIndexConfig, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasSearchIndexConfigReconciler) invalidate(invalid workflow.Result) ctrl.Result {
	r.Log.Debugf("AtlasSearchIndexConfig is invalid: %v", invalid)
	return invalid.ReconcileResult()
}

// In case it is not going to be supported
func (r *AtlasSearchIndexConfigReconciler) unsupport(ctx *workflow.Context) ctrl.Result {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, errors.New("the AtlasSearchIndexConfig is not supported by Atlas for government")).
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult()
}

func (r *AtlasSearchIndexConfigReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) ctrl.Result {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err)
	ctx.SetConditionFromResult(api.ReadyType, terminated)
	return terminated.ReconcileResult()
}

func (r *AtlasSearchIndexConfigReconciler) ready(ctx *workflow.Context) ctrl.Result {
	result := workflow.OK()
	ctx.SetConditionFromResult(api.ReadyType, result)
	return result.ReconcileResult()
}

func (r *AtlasSearchIndexConfigReconciler) lock(ctx *workflow.Context, searchIndexConfig *akov2.AtlasSearchIndexConfig) ctrl.Result {
	if customresource.HaveFinalizer(searchIndexConfig, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, searchIndexConfig, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	return r.ready(ctx)
}

func (r *AtlasSearchIndexConfigReconciler) release(ctx *workflow.Context, searchIndexConfig *akov2.AtlasSearchIndexConfig) ctrl.Result {
	if !customresource.HaveFinalizer(searchIndexConfig, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, searchIndexConfig, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return r.ready(ctx)
}
