package workflow

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const defaultRetry = time.Second * 10

type Result struct {
	terminated   bool
	requeueAfter time.Duration
	message      string
	reason       ConditionReason
}

func OK() Result {
	return Result{
		terminated:   false,
		requeueAfter: -1,
	}
}

func Terminate(reason ConditionReason, message string) Result {
	return Result{
		terminated:   true,
		requeueAfter: defaultRetry,
		reason:       reason,
		message:      message,
	}
}

func InProgress(reason ConditionReason, message string) Result {
	return Result{
		terminated:   false,
		requeueAfter: defaultRetry,
		reason:       reason,
		message:      message,
	}
}

func (r *Result) WithRetry(retry time.Duration) *Result {
	r.requeueAfter = retry
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
