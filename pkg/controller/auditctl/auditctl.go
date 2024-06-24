package auditctl

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
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
	AuditService  audit.Service
}

func (r *AtlasAuditingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: make sure this is done when constructing AtlasAuditingReconciler
	// log := r.Log.With("atlasdatabaseuser", req.NamespacedName)

	auditing := &v1alpha1.AtlasAuditing{}
	result := customresource.PrepareResource(ctx, r.Client, req, auditing, r.Log)
	if !result.IsOk() {
		return ctrl.Result{}, fmt.Errorf("%w %s/%s, will not reconcile", ErrorNotFound, req.Namespace, req.Name)
	}

	if customresource.ReconciliationShouldBeSkipped(auditing) {
		msg := fmt.Sprintf("-> Skipping AtlasAuditing reconciliation as annotation %s=%s",
			customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip)
		r.Log.Infow(msg, "spec", auditing.Spec)
		return ctrl.Result{}, reconcile.TerminalError(fmt.Errorf("%w: %s", ErrorSkipped, msg))
	}

	if err := validate.Auditing(auditing); err != nil {
		return ctrl.Result{}, err
	}

	resultAuditing, err := r.reconcile(ctx, auditing)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := customresource.ApplyLastConfigApplied(ctx, resultAuditing, r.Client); err != nil {
		return ctrl.Result{}, err
	}
	return result.ReconcileResult(), nil
}

// reconcile detects and updates the auditing config state machine
//
// State machine summary:
//
// UNKNOWN --(has projects)--> LOCKED (set finalizer)
// UNKNOWN --(has NO projects)--> RELEASED
// IN LOCKED: evaluate project auditing pairs
// RELEASED --(delete)--> DELETED (Delete succeeds when no finalizer is set or was released)
func (r *AtlasAuditingReconciler) reconcile(ctx context.Context, auditing *v1alpha1.AtlasAuditing) (*v1alpha1.AtlasAuditing, error) {
	switch hasProjectRefs(auditing) { // Evaluate UNKNOWN status...
	case true: // LOCKED
		resultAuditing := r.lock(auditing)
		return r.reconcileProjectsAuditing(ctx, resultAuditing)
	default: // RELEASED
		resultAuditing := r.release(auditing)
		return resultAuditing, nil
	}
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
