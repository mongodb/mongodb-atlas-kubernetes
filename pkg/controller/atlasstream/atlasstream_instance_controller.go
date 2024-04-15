package atlasstream

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type InstanceReconciler struct {
	watch.ResourceWatcher

	Client                      client.Client
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	Log                         *zap.SugaredLogger
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreaminstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreaminstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreaminstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreaminstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// https://dreampuf.github.io/GraphvizOnline/#digraph%20G%20%7B%0A%20%20%20%20subgraph%20cluster_pending%20%7B%0A%20%20%20%20%20%20%20%20skipped%3B%0A%20%20%20%20%20%20%20%20invalid%3B%0A%20%20%20%20%20%20%20%20unsupported%3B%0A%20%20%20%20%20%20%20%20terminated%3B%0A%20%20%20%20%20%20%20%20label%20%3D%20%22pending%22%3B%0A%20%20%20%20%7D%0A%0A%20%20%20%20deleted%20%5Blabel%3D%22deleted%5Cnfinalizer%20unset%22%5D%0A%0A%20%20%20%20pending%20-%3E%20pending%20%5Blabel%3D%22skip%5Cninvalidate%5Cnunsupport%5Cnterminate%22%5D%0A%20%20%20%20pending%20-%3E%20ready%20%5Blabel%3D%22create%22%5D%0A%20%20%20%20pending%20-%3E%20deleted%20%5Blabel%3D%22delete%22%5D%0A%20%20%20%20ready%20-%3E%20ready%20%5Blabel%3D%22update%22%5D%0A%20%20%20%20ready%20-%3E%20deleted%20%5Blabel%3D%22delete%22%5D%0A%20%20%20%20ready%20-%3E%20pending%20%5Blabel%3D%22terminate%22%5D%0A%7D%0A

func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasstreaminstance", req.NamespacedName)
	log.Infow("-> Starting AtlasStreamInstance reconciliation")

	akoStreamInstance := akov2.AtlasStreamInstance{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoStreamInstance, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	// starting from this point we have an entity at hand and can dispatch
	if customresource.ReconciliationShouldBeSkipped(&akoStreamInstance) {
		return r.skip(ctx, log, &akoStreamInstance), nil
	}

	workflowCtx := customresource.MarkReconciliationStarted(r.Client, &akoStreamInstance, log, ctx)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, &akoStreamInstance)

	// is it "invalid" state?
	isValid := customresource.ValidateResourceVersion(workflowCtx, &akoStreamInstance, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(&akoStreamInstance) {
		return r.unsupport(workflowCtx)
	}

	project := akov2.AtlasProject{}
	if err := r.Client.Get(ctx, akoStreamInstance.AtlasProjectObjectKey(), &project); err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	atlasClient, orgID, err := r.AtlasProvider.SdkClient(workflowCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		return r.terminate(workflowCtx, workflow.AtlasAPIAccessNotConfigured, err)
	}
	workflowCtx.SdkClient = atlasClient
	workflowCtx.OrgID = orgID

	return r.handlePendingOrReady(workflowCtx, &project, &akoStreamInstance)
}

func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasStreamInstance").
		For(&akov2.AtlasStreamInstance{}, builder.WithPredicates(r.GlobalPredicates...)).
		Complete(r)
}

func setCondition(ctx *workflow.Context, condition status.ConditionType, result workflow.Result) {
	ctx.SetConditionFromResult(condition, result)
	logIfWarning(ctx, result)
}

func logIfWarning(ctx *workflow.Context, result workflow.Result) {
	if result.IsWarning() {
		ctx.Log.Warnw(result.GetMessage())
	}
}
