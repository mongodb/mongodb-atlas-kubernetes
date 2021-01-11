package workflow

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const DefaultRetry = time.Second * 10

type ConditionReason string

type Result struct {
	done         bool
	requeueAfter time.Duration
	message      string
	reason       ConditionReason
}

// OK indicates that the reconciliation logic can proceed further
func OK() Result {
	return Result{
		done:         false,
		requeueAfter: -1,
	}
}

// Terminate indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued)
// 'reason' and 'message' indicate the error state and are supposed to be reflected in the `conditions` for the
// reconciled Custom Resource
func Terminate(reason ConditionReason, message string) Result {
	return Result{done: true, requeueAfter: DefaultRetry, reason: reason, message: message}
}

// TerminateSilently indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued)
// The status of the reconciled Custom Resource is not supposed to be updated.
func TerminateSilently() Result {
	return Result{done: true, requeueAfter: DefaultRetry}
}

func (r Result) WithRetry(retry time.Duration) Result {
	r.requeueAfter = retry
	return r
}

func (r Result) WithoutRetry() Result {
	r.requeueAfter = -1
	return r
}

func (r Result) IsOk() bool {
	return !r.done
}

func (r Result) ReconcileResult() reconcile.Result {
	if r.requeueAfter < 0 {
		return reconcile.Result{}
	}
	return reconcile.Result{RequeueAfter: r.requeueAfter}
}
