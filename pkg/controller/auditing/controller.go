package auditing

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer/auditing"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

var (
	// ErrorNotFound a resources was not found when expected
	ErrorNotFound = errors.New("not Found")

	// ErrorSkipped when a resource needs to be skipped
	ErrorSkipped = errors.New("skipped")
)

type AuditingStatus string

// AtlasAuditingReconciler reconciles an AtlasAuditing object
type AtlasAuditingReconciler struct {
	Client        client.Client
	Log           *zap.SugaredLogger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
	AuditService  auditing.Service
}

func (r *AtlasAuditingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: make sure this is done when constructing AtlasAuditingReconciler
	// log := r.Log.With("atlasdatabaseuser", req.NamespacedName)

	auditing := &v1alpha1.AtlasAuditing{}
	result := customresource.PrepareResource(ctx, r.Client, req, auditing, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult(), fmt.Errorf("%w %s/%s, will not reconcile",
			ErrorNotFound, req.Namespace, req.Name)
	}

	if customresource.ReconciliationShouldBeSkipped(auditing) {
		r.Log.Infow(fmt.Sprintf("-> Skipping AtlasAuditing reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", auditing.Spec)
		return workflow.OK().ReconcileResult(), fmt.Errorf("%w %s/%s has skip annotation, will not reconcile",
			ErrorSkipped, auditing.Namespace, auditing.Name)
	}

	conditions := akov2.InitCondition(auditing, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx)

	if err := validate.Auditing(auditing); err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.ValidationSucceeded, result)
		return result.ReconcileResult(), err
	}

	resultAuditing, err := r.evaluateState(ctx, auditing)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.AuditingReadyType, result)
		return result.ReconcileResult(), err
	}
	if err := customresource.ApplyLastConfigApplied(ctx, resultAuditing, r.Client); err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	workflowCtx.SetConditionTrue(api.AuditingReadyType)
	workflowCtx.SetConditionTrue(api.ReadyType)

	return result.ReconcileResult(), nil
}

// evaluateState detects and updates the auditing config state machine
//
// State machine summary:
//
// UNKNOWN --(has projects)--> LOCKED (set finalizer)
// UNKNOWN --(has NO projects)--> RELEASED
// IN LOCKED: evaluate project auditing pairs
// RELEASED --(delete)--> DELETED (Delete succeeds when no finalizer is set or was released)
func (r *AtlasAuditingReconciler) evaluateState(ctx context.Context, auditing *v1alpha1.AtlasAuditing) (*v1alpha1.AtlasAuditing, error) {
	// UNKNOWN status...
	if hasProjectRefs(auditing) {
		resultAuditing := r.lock(auditing) // LOCKED
		return r.reconcileProjectsAuditing(ctx, resultAuditing)
	}
	resultAuditing := r.release(auditing) // RELEASED
	return resultAuditing, nil
}

func (r *AtlasAuditingReconciler) lock(auditing *v1alpha1.AtlasAuditing) *v1alpha1.AtlasAuditing {
	resultAuditing := auditing.DeepCopy()
	customresource.SetFinalizer(resultAuditing, customresource.FinalizerLabel)
	return resultAuditing
}

func (r *AtlasAuditingReconciler) release(auditing *v1alpha1.AtlasAuditing) *v1alpha1.AtlasAuditing {
	resultAuditing := auditing.DeepCopy()
	customresource.UnsetFinalizer(resultAuditing, customresource.FinalizerLabel)
	return resultAuditing
}

func hasProjectRefs(auditing *v1alpha1.AtlasAuditing) bool {
	return len(projectIDs(auditing)) > 0
}

func projectIDs(auditing *v1alpha1.AtlasAuditing) []string {
	if auditing.Spec.Type == v1alpha1.Standalone {
		return auditing.Spec.ProjectIDs
	}
	return linkedHasProjectIDs(auditing)
}

func linkedHasProjectIDs(auditing *v1alpha1.AtlasAuditing) []string {
	panic("unimplemented")
}
