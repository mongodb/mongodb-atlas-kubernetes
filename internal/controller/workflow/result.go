// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workflow

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
)

const (
	DefaultRetry = time.Second * 10
)

// Note: DeprecatedResult is a legacy type that is used to return the result of the reconciliation logic.
type DeprecatedResult struct {
	terminated   bool
	requeueAfter time.Duration
	message      string
	reason       ConditionReason
	// warning indicates if the reconciliation hasn't ended the expected way. Most of all this may happens in case of
	// an error
	warning bool
	deleted bool
	err     error
}

// OK indicates that the reconciliation logic can proceed further
func OK() DeprecatedResult {
	return DeprecatedResult{
		terminated:   false,
		requeueAfter: -1,
	}
}

func Requeue(period time.Duration) DeprecatedResult {
	return DeprecatedResult{
		terminated:   false,
		requeueAfter: period,
	}
}

// Terminate indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued).
// This is not an expected termination of the reconciliation process so 'warning' flag is set to 'true'.
// 'reason' and 'message' indicate the error state and are supposed to be reflected in the `conditions` for the
// reconciled Custom Resource.
func Terminate(reason ConditionReason, err error) DeprecatedResult {
	dryrun.AddTerminationError(err) // TODO: factor this in favor of controller-runtime error handling

	return DeprecatedResult{
		terminated:   true,
		requeueAfter: DefaultRetry,
		reason:       reason,
		message:      err.Error(),
		warning:      true,
		err:          err,
	}
}

// InProgress indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued).
// This is an expected termination of the reconciliation process so 'warning' flag is set to 'false'.
// 'reason' and 'message' indicate the in-progress state and are supposed to be reflected in the 'conditions' for the reconciled Custom Resource.
func InProgress(reason ConditionReason, message string) DeprecatedResult {
	return DeprecatedResult{
		terminated:   true,
		requeueAfter: DefaultRetry,
		reason:       reason,
		message:      message,
		warning:      false,
	}
}

func Deleted() DeprecatedResult {
	return DeprecatedResult{
		terminated:   false,
		requeueAfter: -1,
		deleted:      true,
	}
}

func (r DeprecatedResult) IsDeleted() bool {
	return r.deleted
}

// TerminateSilently indicates that the reconciliation logic cannot proceed and needs to be finished (and possibly requeued)
// The status of the reconciled Custom Resource is not supposed to be updated.
func TerminateSilently(err error) DeprecatedResult {
	dryrun.AddTerminationError(err)

	return DeprecatedResult{terminated: true, requeueAfter: DefaultRetry, err: err}
}

func (r DeprecatedResult) WithRetry(retry time.Duration) DeprecatedResult {
	r.requeueAfter = retry
	return r
}

// WithoutRetry indicates that no retry must happen after the reconciliation is over. This should usually be used
// in cases when retry won't fix the situation like when the spec is incorrect and requires the user to update it.
func (r DeprecatedResult) WithoutRetry() DeprecatedResult {
	r.err = nil
	r.requeueAfter = -1
	return r
}

func (r DeprecatedResult) WithMessage(message string) DeprecatedResult {
	r.message = message
	return r
}

func (r DeprecatedResult) IsOk() bool {
	return !r.terminated
}

func (r DeprecatedResult) IsWarning() bool {
	return r.warning
}

func (r DeprecatedResult) IsInProgress() bool {
	return r.terminated && !r.warning
}

func (r DeprecatedResult) GetMessage() string {
	return r.message
}

func (r DeprecatedResult) GetError() error { return r.err }

func (r DeprecatedResult) CloneWithoutError() DeprecatedResult {
	return DeprecatedResult{
		terminated:   r.terminated,
		requeueAfter: r.requeueAfter,
		message:      r.message,
		reason:       r.reason,
		warning:      r.warning,
		deleted:      r.deleted,
		err:          nil,
	}
}

func (r DeprecatedResult) ReconcileResult() (reconcile.Result, error) {
	if r.requeueAfter < 0 {
		return reconcile.Result{}, nil
	}
	if r.err != nil {
		return reconcile.Result{}, r.err
	}
	return reconcile.Result{RequeueAfter: r.requeueAfter}, nil
}
