package atlasstream

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
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

func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasstreaminstance", req.NamespacedName)

	streamInstance := akov2.AtlasStreamInstance{}
	result := customresource.PrepareResource(ctx, r.Client, req, &streamInstance, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(&streamInstance) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasStreamInstance reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", streamInstance.Spec)
		if !streamInstance.GetDeletionTimestamp().IsZero() {
			if err := customresource.ManageFinalizer(ctx, r.Client, &streamInstance, customresource.UnsetFinalizer); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to remove finalizer", "error", err)

				return result.ReconcileResult(), nil
			}
		}

		return workflow.OK().ReconcileResult(), nil
	}

	workflowCtx := customresource.MarkReconciliationStarted(r.Client, &streamInstance, log, ctx)
	log.Infow("-> Starting AtlasStreamInstance reconciliation")

	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, &streamInstance)

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, &streamInstance, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("AtlasStreamInstance validation result: %v", resourceVersionIsValid)

		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if !r.AtlasProvider.IsResourceSupported(&streamInstance) {
		result = workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasStreamInstance is not supported by Atlas for government").
			WithoutRetry()
		setCondition(workflowCtx, status.ReadyType, result)

		return result.ReconcileResult(), nil
	}

	project := akov2.AtlasProject{}
	if result = r.readProjectResource(ctx, &streamInstance, &project); !result.IsOk() {
		setCondition(workflowCtx, status.ReadyType, result)

		return result.ReconcileResult(), nil
	}

	atlasClient, orgID, err := r.AtlasProvider.SdkClient(workflowCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		setCondition(workflowCtx, status.ReadyType, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.SdkClient = atlasClient
	workflowCtx.OrgID = orgID

	result = r.ensureStreamInstance(workflowCtx, &project, &streamInstance)
	setCondition(workflowCtx, status.ReadyType, result)

	return result.ReconcileResult(), nil
}

func (r *InstanceReconciler) readProjectResource(ctx context.Context, instance *akov2.AtlasStreamInstance, project *akov2.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, instance.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	return workflow.OK()
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
