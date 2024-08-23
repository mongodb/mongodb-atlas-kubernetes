package atlasbackupcompliancepolicy

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
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

type AtlasBackupCompliancePolicyReconciler struct {
	Client                      client.Client
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	Log                         *zap.SugaredLogger
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupcompliancepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasbackupcompliancepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasbackupcompliancepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects,verbs=get;list
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects,verbs=get;list

func (r *AtlasBackupCompliancePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasbackupcompliancepolicy", req.NamespacedName)
	log.Infow("-> Starting AtlasBackupCompliancePolicy reonciliation")

	bcp := &akov2.AtlasBackupCompliancePolicy{}
	result := customresource.PrepareResource(ctx, r.Client, req, bcp, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(bcp) {
		return r.skip(ctx, log, bcp), nil
	}

	conditions := akov2.InitCondition(bcp, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, bcp)

	isValid := customresource.ValidateResourceVersion(workflowCtx, bcp, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(bcp) {
		return r.unsupport(workflowCtx)
	}

	return r.ensureAtlasBackupCompliancePolicy(workflowCtx, bcp)
}

func (r *AtlasBackupCompliancePolicyReconciler) ensureAtlasBackupCompliancePolicy(workflowCtx *workflow.Context, bcp *akov2.AtlasBackupCompliancePolicy) (ctrl.Result, error) {
	projects := &akov2.AtlasProjectList{}
	listOpts := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasProjectByBackupCompliancePolicyIndex,
			client.ObjectKeyFromObject(bcp).String(),
		),
	}
	err := r.Client.List(workflowCtx.Context, projects, listOpts)
	if err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	if len(projects.Items) > 0 {
		return r.lock(workflowCtx, bcp)
	}

	return r.release(workflowCtx, bcp)
}

func (r *AtlasBackupCompliancePolicyReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasBackupCompliancePolicy").
		For(&akov2.AtlasBackupCompliancePolicy{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.findBCPForProjects),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func NewAtlasBackupCompliancePolicyReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasBackupCompliancePolicyReconciler {
	return &AtlasBackupCompliancePolicyReconciler{
		Scheme:                   mgr.GetScheme(),
		Client:                   mgr.GetClient(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasBackupCompliancePolicy"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasBackupCompliancePolicy").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasBackupCompliancePolicyReconciler) findBCPForProjects(_ context.Context, obj client.Object) []reconcile.Request {
	project, ok := obj.(*akov2.AtlasProject)
	if !ok {
		r.Log.Warnf("watching AtlasProject but got %T", obj)
		return nil
	}

	if project.Spec.BackupCompliancePolicyRef == nil {
		return nil
	}

	return []reconcile.Request{
		{
			NamespacedName: *project.Spec.BackupCompliancePolicyRef.GetObject(project.Namespace),
		},
	}
}

func (r *AtlasBackupCompliancePolicyReconciler) skip(ctx context.Context, log *zap.SugaredLogger, bcp *akov2.AtlasBackupCompliancePolicy) ctrl.Result {
	log.Infow(fmt.Sprintf("-> Skipping AtlasBackupCompliancePolicy reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", bcp.Spec)
	if !bcp.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, bcp, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err.Error())
			log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasBackupCompliancePolicyReconciler) invalidate(invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasBackupCompliancePolicy is invalid: %v", invalid)
	return invalid.ReconcileResult(), nil
}

func (r *AtlasBackupCompliancePolicyReconciler) unsupport(ctx *workflow.Context) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, "the AtlasBackupCompliancePolicy is not supported by Atlas for government").
		WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

func (r *AtlasBackupCompliancePolicyReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err.Error())
	ctx.SetConditionFromResult(api.ReadyType, terminated)
	return terminated.ReconcileResult(), nil
}

func (r *AtlasBackupCompliancePolicyReconciler) ready(ctx *workflow.Context) (ctrl.Result, error) {
	result := workflow.OK()
	ctx.SetConditionFromResult(api.ReadyType, result)
	return result.ReconcileResult(), nil
}

func (r *AtlasBackupCompliancePolicyReconciler) lock(ctx *workflow.Context, bcp *akov2.AtlasBackupCompliancePolicy) (ctrl.Result, error) {
	if customresource.HaveFinalizer(bcp, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, bcp, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	return r.ready(ctx)
}

func (r *AtlasBackupCompliancePolicyReconciler) release(ctx *workflow.Context, bcp *akov2.AtlasBackupCompliancePolicy) (ctrl.Result, error) {
	if !customresource.HaveFinalizer(bcp, customresource.FinalizerLabel) {
		return r.ready(ctx)
	}

	if err := customresource.ManageFinalizer(ctx.Context, r.Client, bcp, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return r.ready(ctx)
}
