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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func TestGetObservedGeneration(t *testing.T) {
	type args struct {
		obj        client.Object
		prevStatus StatusObject
		nextState  state.ResourceState
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "No previous state, returns current generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 3}},
				prevStatus: newDummyObject(metav1.ObjectMeta{}, nil),
				nextState:  state.StateInitial,
			},
			want: 3,
		},
		{
			name: "Switch from Creating to Created, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 2}},
				prevStatus: prevStatusObject(state.StateCreating, 7),
				nextState:  state.StateCreated,
			},
			want: 7,
		},
		{
			name: "Switch from Updating to Updated, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 3}},
				prevStatus: prevStatusObject(state.StateUpdating, 9),
				nextState:  state.StateUpdated,
			},
			want: 9,
		},
		{
			name: "Switch from Deleting to Deleted, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 1}},
				prevStatus: prevStatusObject(state.StateDeleting, 2),
				nextState:  state.StateDeleted,
			},
			want: 2,
		},
		{
			name: "Polling from Creating, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 2}},
				prevStatus: prevStatusObject(state.StateCreating, 7),
				nextState:  state.StateCreating,
			},
			want: 7,
		},
		{
			name: "Polling Updating, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 3}},
				prevStatus: prevStatusObject(state.StateUpdating, 9),
				nextState:  state.StateUpdating,
			},
			want: 9,
		},
		{
			name: "Polling Deleting, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 1}},
				prevStatus: prevStatusObject(state.StateDeleting, 2),
				nextState:  state.StateDeleting,
			},
			want: 2,
		},
		{
			name: "Start Deleting, uses observed generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 1}},
				prevStatus: prevStatusObject(state.StateDeletionRequested, 2),
				nextState:  state.StateDeleting,
			},
			want: 2,
		},
		{
			name: "Irrelevant state change, returns obj generation",
			args: args{
				obj:        &v1.Pod{ObjectMeta: metav1.ObjectMeta{Generation: 8}},
				prevStatus: prevStatusObject(state.StateInitial, 4),
				nextState:  state.StateInitial,
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getObservedGeneration(tt.args.obj, tt.args.prevStatus.GetConditions(), tt.args.nextState)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewReadyCondition(t *testing.T) {
	tests := []struct {
		name       string
		nextState  state.ResourceState
		wantCond   metav1.ConditionStatus
		wantReason string
		wantMsg    string
	}{
		{
			name:       "Initial - Pending",
			nextState:  state.StateInitial,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is in initial state.",
		},
		{
			name:       "ImportRequested - Pending",
			nextState:  state.StateImportRequested,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is being imported.",
		},
		{
			name:       "Creating - Pending",
			nextState:  state.StateCreating,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is pending.",
		},
		{
			name:       "Updating - Pending",
			nextState:  state.StateUpdating,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is pending.",
		},
		{
			name:       "Deleting - Pending",
			nextState:  state.StateDeleting,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is pending.",
		},
		{
			name:       "DeletionRequested - Pending",
			nextState:  state.StateDeletionRequested,
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonPending,
			wantMsg:    "Resource is pending.",
		},
		{
			name:       "Imported - Settled",
			nextState:  state.StateImported,
			wantCond:   metav1.ConditionTrue,
			wantReason: ReadyReasonSettled,
			wantMsg:    "Resource is imported.",
		},
		{
			name:       "Created - Settled",
			nextState:  state.StateCreated,
			wantCond:   metav1.ConditionTrue,
			wantReason: ReadyReasonSettled,
			wantMsg:    "Resource is settled.",
		},
		{
			name:       "Updated - Settled",
			nextState:  state.StateUpdated,
			wantCond:   metav1.ConditionTrue,
			wantReason: ReadyReasonSettled,
			wantMsg:    "Resource is settled.",
		},
		{
			name:       "Unknown state - Error",
			nextState:  "nonexistent",
			wantCond:   metav1.ConditionFalse,
			wantReason: ReadyReasonError,
			wantMsg:    "unknown state: nonexistent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Result{NextState: tt.nextState}
			cond := NewReadyCondition(result)
			assert.Equal(t, tt.wantCond, cond.Status)
			assert.Equal(t, tt.wantReason, cond.Reason)
			assert.Equal(t, state.ReadyCondition, cond.Type)
			assert.Equal(t, tt.wantMsg, cond.Message)
		})
	}
}

func TestReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	addKnownTestTypes(scheme)

	baseObj := &dummyObject{
		TypeMeta: metav1.TypeMeta{
			Kind:       "dummyObject",
			APIVersion: "test.dummy.example.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "mypod",
			Namespace:  "default",
			Generation: 1,
		},
	}
	objKey := types.NamespacedName{Name: "mypod", Namespace: "default"}

	tests := []struct {
		name         string
		existingObj  client.Object
		interceptors *interceptor.Funcs
		handleState  func(context.Context, *dummyObject) (Result, error)
		wantErr      string
		wantResult   reconcile.Result
	}{
		{
			name:        "get object error",
			existingObj: baseObj,
			interceptors: &interceptor.Funcs{
				Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("simulated get error")
				},
			},
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{NextState: "Initial"}, nil
			},
			wantErr: "unable to get object: simulated get error",
		},
		{
			name:        "object removed is fine",
			existingObj: baseObj,
			interceptors: &interceptor.Funcs{
				Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return apierrors.NewNotFound(schema.GroupResource{}, key.Name)
				},
			},
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{NextState: "Initial"}, nil
			},
			wantResult: reconcile.Result{},
		},
		{
			name:        "failed to set finalizer",
			existingObj: baseObj,
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{NextState: "Initial"}, nil
			},
			interceptors: &interceptor.Funcs{
				Patch: func(ctx context.Context, c client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("simulated patch error")
				},
			},
			wantErr: "failed to manage finalizers: simulated patch error",
		},
		{
			name:        "check state",
			existingObj: baseObj,
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{NextState: "Initial"}, nil
			},
		},
		{
			name:        "should skip reconcile request",
			existingObj: baseObj.WithSkipReconciliationAnnotation(),
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{}}, nil
			},
		},
		{
			name:        "should skip reconcile request on deleted object",
			existingObj: baseObj.WithSkipReconciliationAnnotation().WithDeletedStaze(),
			handleState: func(ctx context.Context, do *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{}}, nil
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme)
			if tc.existingObj != nil {
				builder = builder.WithObjects(tc.existingObj)
				builder = builder.WithStatusSubresource(tc.existingObj)
			}
			if tc.interceptors != nil {
				builder = builder.WithInterceptorFuncs(*tc.interceptors)
			}
			c := builder.Build()
			dummyReconciler := &dummyPodReconciler{handleState: tc.handleState}
			r := &Reconciler[dummyObject]{
				cluster:    &fakeCluster{cli: c},
				reconciler: dummyReconciler,
			}

			req := ctrl.Request{NamespacedName: objKey}
			result, err := r.Reconcile(context.Background(), req)
			assertErrContains(t, tc.wantErr, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestReconcileState(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)
	addKnownTestTypes(scheme)

	ctx := context.Background()

	tests := []struct {
		name       string
		initialObj *dummyObject
		handleFn   func(context.Context, *dummyObject) (Result, error)
		modify     func(t *dummyObject)
		wantResult Result
		wantErr    string
	}{
		{
			name: "simulate error",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateCreated, metav1.ConditionTrue, 1)},
			),
			modify: func(t *dummyObject) {
				// Simulate a state that should cause an error in ReconcileState
				t.Status.Conditions[0].Reason = ""
			},
			wantErr: "unsupported state \"\"",
		},
		{
			name:       "initial state",
			initialObj: newDummyObject(metav1.ObjectMeta{Namespace: "default", Name: "myobj"}, nil),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
		{
			name: "creating",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateCreating, metav1.ConditionFalse, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: true}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: true}},
		},
		{
			name: "created",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateCreated, metav1.ConditionTrue, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
		{
			name: "updating",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateUpdating, metav1.ConditionFalse, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: true}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: true}},
		},
		{
			name: "Updated",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateUpdated, metav1.ConditionTrue, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
		{
			name: "delete request",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateDeletionRequested, metav1.ConditionFalse, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
		{
			name: "deleting",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateDeleting, metav1.ConditionFalse, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: true}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: true}},
		},
		{
			name: "import request",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateImportRequested, metav1.ConditionFalse, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
		{
			name: "importing",
			initialObj: newDummyObject(
				metav1.ObjectMeta{Namespace: "default", Name: "myobj"},
				[]metav1.Condition{newStateCondition(state.StateImported, metav1.ConditionTrue, 1)},
			),
			handleFn: func(context.Context, *dummyObject) (Result, error) {
				return Result{Result: reconcile.Result{Requeue: false}}, nil
			},
			wantResult: Result{Result: reconcile.Result{Requeue: false}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := fake.NewClientBuilder().WithScheme(scheme)
			if tc.initialObj != nil {
				builder = builder.WithObjects(tc.initialObj)
			}
			c := builder.Build()
			dummyReconciler := &dummyPodReconciler{handleState: tc.handleFn}
			r := &Reconciler[dummyObject]{
				cluster:    &fakeCluster{cli: c},
				reconciler: dummyReconciler,
			}
			obj := tc.initialObj.DeepCopy()
			if tc.modify != nil {
				tc.modify(obj)
			}

			gotResult, err := r.ReconcileState(ctx, obj)
			if tc.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult.Requeue, gotResult.Requeue)
			}
		})
	}
}

func addKnownTestTypes(sch *runtime.Scheme) {
	sch.AddKnownTypes(
		schema.GroupVersion{Group: "test.dummy.example.com", Version: "v1"},
		&dummyObject{},
	)
}

type dummyObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            DummyStatus `json:"status,omitempty"`
}

type DummyStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func newDummyObject(objMeta metav1.ObjectMeta, conditions []metav1.Condition) *dummyObject {
	return &dummyObject{
		ObjectMeta: objMeta,
		Status:     DummyStatus{Conditions: conditions},
	}
}

func (*dummyObject) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }

func (do *dummyObject) GetConditions() []metav1.Condition {
	return do.Status.Conditions
}

func (do *dummyObject) DeepCopyObject() runtime.Object {
	return do.DeepCopy()
}

func (do *dummyObject) DeepCopy() *dummyObject {
	if do == nil {
		return nil
	}
	conditions := make([]metav1.Condition, 0, len(do.Status.Conditions))
	for _, condition := range do.Status.Conditions {
		conditions = append(conditions, *condition.DeepCopy())
	}
	return &dummyObject{
		TypeMeta:   do.TypeMeta,
		ObjectMeta: *do.ObjectMeta.DeepCopy(),
		Status:     DummyStatus{Conditions: conditions},
	}
}

func (do *dummyObject) WithSkipReconciliationAnnotation() *dummyObject {
	copyOfDo := do.DeepCopy()

	if copyOfDo.Annotations == nil {
		copyOfDo.Annotations = make(map[string]string)
	}
	copyOfDo.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip

	return copyOfDo
}

func (do *dummyObject) WithDeletedStaze() *dummyObject {
	copyOfDo := do.DeepCopy()

	if len(copyOfDo.Status.Conditions) == 0 {
		copyOfDo.Status.Conditions = []metav1.Condition{}
	}

	copyOfDo.Status.Conditions = append(copyOfDo.Status.Conditions, metav1.Condition{
		Type:   "State",
		Status: metav1.ConditionTrue,
		Reason: "Deleted",
	})

	return copyOfDo
}

func prevStatusObject(state state.ResourceState, observedGen int64) StatusObject {
	return newDummyObject(metav1.ObjectMeta{}, []metav1.Condition{
		newStateCondition(state, metav1.ConditionTrue, observedGen),
	})
}

func newStateCondition(reason state.ResourceState, status metav1.ConditionStatus, observedGen int64) metav1.Condition {
	return metav1.Condition{
		Type:               "State",
		Status:             status,
		Reason:             string(reason),
		ObservedGeneration: observedGen,
	}
}

// Dummy reconciler implementing StateReconciler[do *dummyObject]
type dummyPodReconciler struct {
	handleState func(context.Context, *dummyObject) (Result, error)
}

func (d *dummyPodReconciler) SetupWithManager(_ ctrl.Manager, _ reconcile.Reconciler, _ controller.Options) error {
	return nil
}
func (d *dummyPodReconciler) For() (client.Object, builder.Predicates) {
	return nil, builder.Predicates{}
}
func (d *dummyPodReconciler) HandleInitial(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleImportRequested(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleImported(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleCreating(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleCreated(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleUpdating(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleUpdated(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleDeletionRequested(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}
func (d *dummyPodReconciler) HandleDeleting(ctx context.Context, do *dummyObject) (Result, error) {
	return d.handleState(ctx, do)
}

type fakeCluster struct {
	cluster.Cluster
	cli client.Client
}

func (f *fakeCluster) GetClient() client.Client   { return f.cli }
func (f *fakeCluster) GetScheme() *runtime.Scheme { return f.cli.Scheme() }
