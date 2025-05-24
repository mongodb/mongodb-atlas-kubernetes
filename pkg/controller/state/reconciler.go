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

package state

import (
	"context"
	"fmt"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/finalizer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type Result struct {
	reconcile.Result
	NextState state.ResourceState
	StateMsg  string
}

type StateHandler[T any] interface {
	SetupWithManager(ctrl.Manager, reconcile.Reconciler) error
	HandleInitial(context.Context, *T) (Result, error)
	HandleImportRequested(context.Context, *T) (Result, error)
	HandleImported(context.Context, *T) (Result, error)
	HandleCreating(context.Context, *T) (Result, error)
	HandleCreated(context.Context, *T) (Result, error)
	HandleUpdating(context.Context, *T) (Result, error)
	HandleUpdated(context.Context, *T) (Result, error)
	HandleDeletionRequested(context.Context, *T) (Result, error)
	HandleDeleting(context.Context, *T) (Result, error)
	// Deleted, not handled as it is a terminal state
}

const (
	ReadyReasonError   = "Error"
	ReadyReasonPending = "Pending"
	ReadyReasonSettled = "Settled"
)

type Reconciler[T any] struct {
	cluster         cluster.Cluster
	reconciler      StateHandler[T]
	unstructuredGVK schema.GroupVersionKind
}

type UnstructuredStateReconciler = StateHandler[unstructured.Unstructured]

func NewStateReconciler[T any](target StateHandler[T]) *Reconciler[T] {
	return &Reconciler[T]{
		reconciler: target,
	}
}

func NewUnstructuredStateReconciler(target UnstructuredStateReconciler, gvk schema.GroupVersionKind) *Reconciler[unstructured.Unstructured] {
	return &Reconciler[unstructured.Unstructured]{
		reconciler:      target,
		unstructuredGVK: gvk,
	}
}

func (r *Reconciler[T]) SetupWithManager(mgr ctrl.Manager) error {
	r.cluster = mgr
	return r.reconciler.SetupWithManager(mgr, r)
}

func (r *Reconciler[T]) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx).WithName("state")
	logger.Info("reconcile started", "req", req)

	t := new(T)
	obj := any(t).(client.Object)
	if u, ok := obj.(*unstructured.Unstructured); ok {
		u.SetGroupVersionKind(r.unstructuredGVK)
	}

	err := r.cluster.GetClient().Get(ctx, req.NamespacedName, obj)
	if apierrors.IsNotFound(err) {
		// object is already gone, nothing to do.
		return reconcile.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to get object: %w", err)
	}

	currentStatus, err := newStatusObject(obj)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get status object: %w", err)
	}
	currentState := state.GetState(currentStatus.GetConditions())

	logger.Info("reconcile started", "currentState", currentState)
	if err := finalizer.EnsureFinalizers(ctx, r.cluster.GetClient(), obj, "mongodb.com/finalizer"); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to manage finalizers: %w", err)
	}

	result, reconcileErr := r.ReconcileState(ctx, t)
	stateStatus := true
	if reconcileErr != nil {
		// error message will be displayed in Ready state.
		stateStatus = false
	}
	newStatus, err := newStatusObject(obj)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get status object: %w", err)
	}
	observedGeneration := getObservedGeneration(obj, currentStatus, result.NextState)
	newStatusConditions := newStatus.GetConditions()
	state.EnsureState(&newStatusConditions, observedGeneration, result.NextState, result.StateMsg, stateStatus)

	logger.Info("reconcile finished", "nextState", result.NextState)

	if result.NextState == state.StateDeleted {
		if err := finalizer.UnsetFinalizers(ctx, r.cluster.GetClient(), obj, "mongodb.com/finalizer"); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to unset finalizer: %w", err)
		}

		return result.Result, reconcileErr
	}

	ready := NewReadyCondition(result)
	ready.ObservedGeneration = observedGeneration

	if reconcileErr != nil {
		ready.Status = metav1.ConditionFalse
		ready.Reason = ReadyReasonError
		ready.Message = reconcileErr.Error()
	}

	meta.SetStatusCondition(&newStatusConditions, ready)

	if err := patchStatus(ctx, r.cluster.GetClient(), obj, newStatus); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to patch status: %w", err)
	}

	return result.Result, reconcileErr
}

func NewReadyCondition(result Result) metav1.Condition {
	var (
		readyReason, msg string
		cond             metav1.ConditionStatus
	)

	switch result.NextState {
	case state.StateInitial:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is in initial state."

	case state.StateImportRequested:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is being imported."

	case state.StateCreating:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is pending."

	case state.StateUpdating:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is pending."

	case state.StateDeleting:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is pending."

	case state.StateDeletionRequested:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonPending
		msg = "Resource is pending."

	case state.StateImported:
		cond = metav1.ConditionTrue
		readyReason = ReadyReasonSettled
		msg = "Resource is imported."

	case state.StateCreated:
		cond = metav1.ConditionTrue
		readyReason = ReadyReasonSettled
		msg = "Resource is settled."

	case state.StateUpdated:
		cond = metav1.ConditionTrue
		readyReason = ReadyReasonSettled
		msg = "Resource is settled."

	default:
		cond = metav1.ConditionFalse
		readyReason = ReadyReasonError
		msg = fmt.Sprintf("unknown state: %s", result.NextState)
	}

	return metav1.Condition{
		Type:               state.ReadyCondition,
		Status:             cond,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Reason:             readyReason,
		Message:            msg,
	}
}

func (r *Reconciler[T]) ReconcileState(ctx context.Context, t *T) (Result, error) {
	obj := any(t).(client.Object)

	var (
		result = Result{
			Result:    reconcile.Result{},
			NextState: state.StateInitial,
		}

		err error
	)
	statusObj, err := newStatusObject(obj)
	if err != nil {
		return Result{}, fmt.Errorf("failed to get status object: %w", err)
	}
	currentState := state.GetState(statusObj.GetConditions())

	if currentState == state.StateInitial {
		for key := range obj.GetAnnotations() {
			if strings.HasPrefix(key, "mongodb.com/external-") {
				currentState = state.StateImportRequested
			}
		}
	}

	if !obj.GetDeletionTimestamp().IsZero() && currentState != state.StateDeleting {
		currentState = state.StateDeletionRequested
	}

	switch currentState {
	case state.StateInitial:
		result, err = r.reconciler.HandleInitial(ctx, t)
	case state.StateImportRequested:
		result, err = r.reconciler.HandleImportRequested(ctx, t)
	case state.StateImported:
		result, err = r.reconciler.HandleImported(ctx, t)
	case state.StateCreating:
		result, err = r.reconciler.HandleCreating(ctx, t)
	case state.StateCreated:
		result, err = r.reconciler.HandleCreated(ctx, t)
	case state.StateUpdating:
		result, err = r.reconciler.HandleUpdating(ctx, t)
	case state.StateUpdated:
		result, err = r.reconciler.HandleUpdated(ctx, t)
	case state.StateDeletionRequested:
		result, err = r.reconciler.HandleDeletionRequested(ctx, t)
	case state.StateDeleting:
		result, err = r.reconciler.HandleDeleting(ctx, t)
	default:
		return Result{}, fmt.Errorf("unsupported state %q", currentState)
	}

	if result.NextState == "" {
		result.NextState = state.StateInitial
	}

	isReapplyState := result.NextState == state.StateImported ||
		result.NextState == state.StateCreated ||
		result.NextState == state.StateUpdated

	if isReapplyState && result.RequeueAfter == 0 && err == nil {
		requeueAfter, err := PatchReapplyTimestamp(ctx, r.cluster.GetClient(), obj)
		if err != nil {
			return Result{}, err
		}

		result.RequeueAfter = requeueAfter
	}

	return result, err
}

func getObservedGeneration(obj client.Object, prevStatus StatusObject, nextState state.ResourceState) int64 {
	observedGeneration := obj.GetGeneration()
	prevState := state.GetState(prevStatus.GetConditions())

	if prevCondition := meta.FindStatusCondition(prevStatus.GetConditions(), state.StateCondition); prevCondition != nil {
		from := prevState
		to := nextState

		// don't change observed generation if we are:
		// - creating/updating/deleting
		// - just finished creating/updating/deleting
		observedGeneration = prevCondition.ObservedGeneration
		switch {
		case from == state.StateUpdating && to == state.StateUpdating: // polling update
		case from == state.StateUpdating && to == state.StateUpdated: // finished updating

		case from == state.StateCreating && to == state.StateCreating: // polling creation
		case from == state.StateCreating && to == state.StateCreated: // finished creating

		case from == state.StateDeletionRequested && to == state.StateDeleting: // started deletion
		case from == state.StateDeleting && to == state.StateDeleting: // polling deletion
		case from == state.StateDeleting && to == state.StateDeleted: // finshed deletion
		default:
			observedGeneration = obj.GetGeneration()
		}
	}

	return observedGeneration
}
