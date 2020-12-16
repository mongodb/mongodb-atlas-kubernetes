package workflow

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const defaultRetry = time.Second * 10

type ConditionReason string

type Result struct {
	done         bool
	requeueAfter time.Duration
	message      string
	reason       ConditionReason
}

func OK() Result {
	return Result{
		done:         false,
		requeueAfter: -1,
	}
}

func Terminate(reason ConditionReason, message string) Result {
	return Result{done: true, requeueAfter: defaultRetry, reason: reason, message: message}
}

func (r *Result) WithRetry(retryInSeconds time.Duration) *Result {
	r.requeueAfter = retryInSeconds
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
