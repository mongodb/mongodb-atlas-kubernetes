package atlasstream

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

type AtlasStreamsConnectionReconciler struct {
	Client                      client.Client
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	Log                         *zap.SugaredLogger
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreamconnections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreamconnections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreamconnections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreamconnections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreaminstances,verbs=get;list
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreaminstances,verbs=get;list

func (r *AtlasStreamsConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasstreamconnection", req.NamespacedName)
	log.Infow("-> Starting AtlasStreamConnection reconciliation")

	akoStreamConnection := akov2.AtlasStreamConnection{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoStreamConnection, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	return r.ensureAtlasStreamConnection(ctx, log, &akoStreamConnection)
}

func (r *AtlasStreamsConnectionReconciler) ensureAtlasStreamConnection(ctx context.Context, log *zap.SugaredLogger, akoStreamConnection *akov2.AtlasStreamConnection) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(akoStreamConnection) {
		return r.skip(ctx, log, akoStreamConnection), nil
	}

	conditions := api.InitCondition(akoStreamConnection, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx, akoStreamConnection)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, akoStreamConnection)

	isValid := customresource.ValidateResourceVersion(workflowCtx, akoStreamConnection, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid), nil
	}

	if !r.AtlasProvider.IsResourceSupported(akoStreamConnection) {
		return r.unsupport(workflowCtx), nil
	}

	streamInstances := &akov2.AtlasStreamInstanceList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasStreamInstanceByConnectionIndex,
			client.ObjectKeyFromObject(akoStreamConnection).String(),
		),
	}
	err := r.Client.List(ctx, streamInstances, listOps)
	if err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	if len(streamInstances.Items) > 0 {
		return r.lock(workflowCtx, akoStreamConnection)
	}

	return r.release(workflowCtx, akoStreamConnection)
}

func (r *AtlasStreamsConnectionReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasStreamConnection{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasStreamsConnectionReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasStreamConnection").
		For(r.For()).
		Watches(
			&akov2.AtlasStreamInstance{},
			handler.EnqueueRequestsFromMapFunc(r.findStreamConnectionsForStreamInstances),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func NewAtlasStreamsConnectionReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasStreamsConnectionReconciler {
	return &AtlasStreamsConnectionReconciler{
		Scheme:                   c.GetScheme(),
		Client:                   c.GetClient(),
		EventRecorder:            c.GetEventRecorderFor("AtlasStreamsConnection"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasStreamsConnection").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasStreamsConnectionReconciler) findStreamConnectionsForStreamInstances(_ context.Context, obj client.Object) []reconcile.Request {
	streamInstance, ok := obj.(*akov2.AtlasStreamInstance)
	if !ok {
		r.Log.Warnf("watching AtlasStreamInstance but got %T", obj)
		return nil
	}

	requests := make([]reconcile.Request, 0, len(streamInstance.Spec.ConnectionRegistry))
	for i := range streamInstance.Spec.ConnectionRegistry {
		item := streamInstance.Spec.ConnectionRegistry[i]
		requests = append(
			requests,
			reconcile.Request{
				NamespacedName: *item.GetObject(streamInstance.Namespace),
			},
		)
	}

	return requests
}

func (r *AtlasStreamsConnectionReconciler) skip(ctx context.Context, log *zap.SugaredLogger, streamConnection *akov2.AtlasStreamConnection) ctrl.Result {
	log.Infow(fmt.Sprintf("-> Skipping AtlasStreamConnection reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", streamConnection.Spec)
	if !streamConnection.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, streamConnection, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasStreamsConnectionReconciler) invalidate(invalid workflow.Result) ctrl.Result {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasStreamConnection is invalid: %v", invalid)
	return invalid.ReconcileResult()
}

func (r *AtlasStreamsConnectionReconciler) unsupport(ctx *workflow.Context) ctrl.Result {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, errors.New("the AtlasStreamConnection is not supported by Atlas for government")).
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult()
}

func (r *AtlasStreamsConnectionReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err)
	ctx.SetConditionFromResult(api.ReadyType, terminated)
	return terminated.ReconcileResult(), nil
}

func (r *AtlasStreamsConnectionReconciler) ready(ctx *workflow.Context) (ctrl.Result, error) {
	result := workflow.OK()
	ctx.SetConditionFromResult(api.ReadyType, result)
	return result.ReconcileResult(), nil
}

func (r *AtlasStreamsConnectionReconciler) lock(ctx *workflow.Context, streamConnection *akov2.AtlasStreamConnection) (ctrl.Result, error) {
	if customresource.HaveFinalizer(streamConnection, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, streamConnection, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	return r.ready(ctx)
}

func (r *AtlasStreamsConnectionReconciler) release(ctx *workflow.Context, streamConnection *akov2.AtlasStreamConnection) (ctrl.Result, error) {
	if !customresource.HaveFinalizer(streamConnection, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, streamConnection, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return r.ready(ctx)
}
