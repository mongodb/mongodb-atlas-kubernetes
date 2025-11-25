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
	"sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlrtbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/crapi"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/finalizer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type Result struct {
	reconcile.Result
	NextState state.ResourceState
	StateMsg  string
}

type VersionedHandlerFunc[C any, T any] func(client client.Client, atlasClient *C, translator *crapi.Request, deletionProtection bool) StateHandler[T]

type StateHandler[T any] interface {
	SetupWithManager(ctrl.Manager, reconcile.Reconciler, controller.Options) error
	For() (client.Object, builder.Predicates)
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
	supportReapply  bool
}

type ReconcilerOptionFn[T any] func(*Reconciler[T])

func WithCluster[T any](c cluster.Cluster) ReconcilerOptionFn[T] {
	return func(r *Reconciler[T]) {
		r.cluster = c
	}
}

func WithReapplySupport[T any](supportReapply bool) ReconcilerOptionFn[T] {
	return func(r *Reconciler[T]) {
		r.supportReapply = supportReapply
	}
}

type UnstructuredStateReconciler = StateHandler[unstructured.Unstructured]

type ControllerSetupBuilder = ctrlrtbuilder.TypedBuilder[reconcile.Request]

func NewStateReconciler[T any](target StateHandler[T], options ...ReconcilerOptionFn[T]) *Reconciler[T] {
	r := &Reconciler[T]{
		reconciler: target,
	}
	for _, opt := range options {
		opt(r)
	}
	return r
}

func (r *Reconciler[T]) SetupWithManager(mgr ctrl.Manager, defaultOptions controller.Options) error {
	r.cluster = mgr
	return r.reconciler.SetupWithManager(mgr, r, defaultOptions)
}

func (r *Reconciler[T]) For() (client.Object, builder.Predicates) {
	return r.reconciler.For()
}

func (r *Reconciler[T]) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	logger := log.FromContext(ctx).WithName("state")
	logger.Info("reconcile started", "req", req)

	t := new(T)
	obj := any(t).(StatusObject)
	clientObj := any(t).(client.Object)
	if u, ok := clientObj.(*unstructured.Unstructured); ok {
		u.SetGroupVersionKind(r.unstructuredGVK)
	}

	err := r.cluster.GetClient().Get(ctx, req.NamespacedName, clientObj)
	if apierrors.IsNotFound(err) {
		// object is already gone, nothing to do.
		return reconcile.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("unable to get object: %w", err)
	}

	currentStatus := newStatusObject(obj)
	currentState := state.GetState(currentStatus.Status.Conditions)

	if customresource.ReconciliationShouldBeSkipped(clientObj) {
		logger.Info(fmt.Sprintf("Skipping reconciliation by annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip))
		if currentState == state.StateDeleted {
			if err := finalizer.UnsetFinalizers(ctx, r.cluster.GetClient(), clientObj, "mongodb.com/finalizer"); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to unset finalizer: %w", err)
			}
		}

		return ctrl.Result{}, nil
	}

	logger.Info("reconcile started", "currentState", currentState)
	if err := finalizer.EnsureFinalizers(ctx, r.cluster.GetClient(), clientObj, "mongodb.com/finalizer"); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to manage finalizers: %w", err)
	}

	result, reconcileErr := r.ReconcileState(ctx, t)
	stateStatus := true
	if reconcileErr != nil {
		// error message will be displayed in Ready state.
		stateStatus = false
	}

	newStatus := newStatusObject(obj)
	observedGeneration := getObservedGeneration(clientObj, currentStatus.Status.Conditions, result.NextState)
	newStatusConditions := newStatus.Status.Conditions
	state.EnsureState(&newStatusConditions, observedGeneration, result.NextState, result.StateMsg, stateStatus)

	logger.Info("reconcile finished", "nextState", result.NextState)

	if result.NextState == state.StateDeleted {
		if err := finalizer.UnsetFinalizers(ctx, r.cluster.GetClient(), clientObj, "mongodb.com/finalizer"); err != nil {
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
	newStatus.Status.Conditions = newStatusConditions
	if err := patchStatus(ctx, r.cluster.GetClient(), clientObj, newStatus); err != nil {
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
	statusObj := newStatusObject(any(t).(StatusObject))
	currentState := state.GetState(statusObj.Status.Conditions)

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

	if r.supportReapply {
		err := r.reconcileReapply(ctx, obj, result, err)
		if err != nil {
			return Result{}, fmt.Errorf("failed to reconcile reapply: %w", err)
		}
	}

	return result, err
}

func (r *Reconciler[T]) reconcileReapply(ctx context.Context, obj client.Object, result Result, err error) error {
	isReapplyState := result.NextState == state.StateImported ||
		result.NextState == state.StateCreated ||
		result.NextState == state.StateUpdated

	if isReapplyState && result.RequeueAfter == 0 && err == nil {
		requeueAfter, err := PatchReapplyTimestamp(ctx, r.cluster.GetClient(), obj)
		if err != nil {
			return fmt.Errorf("failed to patch reapply timestamp: %w", err)
		}

		result.RequeueAfter = requeueAfter
	}
	return nil
}

func getObservedGeneration(obj client.Object, prevStatusConditions []metav1.Condition, nextState state.ResourceState) int64 {
	observedGeneration := obj.GetGeneration()
	prevState := state.GetState(prevStatusConditions)

	if prevCondition := meta.FindStatusCondition(prevStatusConditions, state.StateCondition); prevCondition != nil {
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
