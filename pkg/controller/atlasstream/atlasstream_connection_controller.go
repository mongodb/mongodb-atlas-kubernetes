package atlasstream

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
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

	conditions := akov2.InitCondition(akoStreamConnection, status.FalseCondition(status.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, akoStreamConnection)

	isValid := customresource.ValidateResourceVersion(workflowCtx, akoStreamConnection, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(akoStreamConnection) {
		return r.unsupport(workflowCtx)
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

func (r *AtlasStreamsConnectionReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasStreamConnection").
		For(&akov2.AtlasStreamConnection{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&akov2.AtlasStreamInstance{},
			handler.EnqueueRequestsFromMapFunc(r.findStreamConnectionsForStreamInstances),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Complete(r)
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
			result := workflow.Terminate(workflow.Internal, err.Error())
			log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasStreamsConnectionReconciler) invalidate(invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasStreamConnection is invalid: %v", invalid)
	return invalid.ReconcileResult(), nil
}

func (r *AtlasStreamsConnectionReconciler) unsupport(ctx *workflow.Context) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, "the AtlasStreamConnection is not supported by Atlas for government").
		WithoutRetry()
	ctx.SetConditionFromResult(status.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

func (r *AtlasStreamsConnectionReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err.Error())
	ctx.SetConditionFromResult(status.ReadyType, terminated)
	return terminated.ReconcileResult(), nil
}

func (r *AtlasStreamsConnectionReconciler) ready(ctx *workflow.Context) (ctrl.Result, error) {
	result := workflow.OK()
	ctx.SetConditionFromResult(status.ReadyType, result)
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
