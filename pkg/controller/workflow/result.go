package workflow

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	DefaultRetry   = time.Second * 10
	DefaultTimeout = time.Minute * 20
)

type Result struct {
	terminated   bool
	requeueAfter time.Duration
	message      string
	reason       ConditionReason
	// warning indicates if the reconciliation hasn't ended the expected way. Most of all this may happens in case of
	// an error
	warning bool
}

// OK indicates that the reconciliation logic can proceed further
func OK() Result {
	return Result{
		terminated:   false,
		requeueAfter: -1,
	}
}

// Terminate indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued).
// This is not an expected termination of the reconciliation process so 'warning' flag is set to 'true'.
// 'reason' and 'message' indicate the error state and are supposed to be reflected in the `conditions` for the
// reconciled Custom Resource.
func Terminate(reason ConditionReason, message string) Result {
	return Result{
		terminated:   true,
		requeueAfter: DefaultRetry,
		reason:       reason,
		message:      message,
		warning:      true,
	}
}

// InProgress indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued).
// This is an expected termination of the reconciliation process so 'warning' flag is set to 'false'.
// 'reason' and 'message' indicate the in-progress state and are supposed to be reflected in the 'conditions' for the reconciled Custom Resource.
func InProgress(reason ConditionReason, message string) Result {
	return Result{
		terminated:   true,
		requeueAfter: DefaultRetry,
		reason:       reason,
		message:      message,
		warning:      false,
	}
}

// TerminateSilently indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued)
// The status of the reconciled Custom Resource is not supposed to be updated.
func TerminateSilently() Result {
	return Result{terminated: true, requeueAfter: DefaultRetry}
}

func (r Result) WithRetry(retry time.Duration) Result {
	r.requeueAfter = retry
	return r
}

// WithoutRetry indicates that no retry must happen after the reconciliation is over. This should usually be used
// in cases when retry won't fix the situation like when the spec is incorrect and requires the user to update it.
func (r Result) WithoutRetry() Result {
	r.requeueAfter = -1
	return r
}

func (r Result) IsOk() bool {
	return !r.terminated
}

func (r Result) ReconcileResult() reconcile.Result {
	if r.requeueAfter < 0 {
		return reconcile.Result{}
	}
	return reconcile.Result{RequeueAfter: r.requeueAfter}
}
