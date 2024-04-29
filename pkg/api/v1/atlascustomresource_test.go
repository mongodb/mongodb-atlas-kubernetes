package v1

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

type resource struct {
	metav1.Object
	runtimeObject runtime.Object
	status.Reader
	status.Writer

	conditions         []status.Condition
	observedGeneration int64
}

// this must be implemented due to naming conflict of metav1.Object and runtime.Object interfaces
func (r *resource) GetObjectKind() schema.ObjectKind {
	return r.runtimeObject.GetObjectKind()
}

// this must be implemented due to naming conflict of metav1.Object and runtime.Object interfaces
func (r *resource) DeepCopyObject() runtime.Object {
	return r.runtimeObject.DeepCopyObject()
}

func (r *resource) GetStatus() status.Status {
	return r
}

func (r *resource) GetConditions() []status.Condition {
	return r.conditions
}

func (r *resource) GetObservedGeneration() int64 {
	return r.observedGeneration
}

func TestInitCondition(t *testing.T) {
	for _, tc := range []struct {
		name             string
		resource         AtlasCustomResource
		defaultCondition status.Condition
		want             []status.Condition
	}{
		{
			name: "keep condition",
			resource: &resource{
				conditions: []status.Condition{
					{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "untouched"},
				},
			},
			defaultCondition: status.Condition{Type: status.ReadyType, Status: corev1.ConditionFalse, Message: "default"},
			want: []status.Condition{
				{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "untouched"},
			},
		},
		{
			name: "set condition",
			resource: &resource{
				conditions: []status.Condition{
					{Type: status.ValidationSucceeded, Status: corev1.ConditionTrue, Message: "untouched"},
				},
			},
			defaultCondition: status.Condition{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			want: []status.Condition{
				{Type: status.ValidationSucceeded, Status: corev1.ConditionTrue, Message: "untouched"},
				{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			},
		},
		{
			name: "set condition on nil list",
			resource: &resource{
				conditions: nil,
			},
			defaultCondition: status.Condition{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			want: []status.Condition{
				{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			},
		},
		{
			name: "set condition on empty list",
			resource: &resource{
				conditions: []status.Condition{},
			},
			defaultCondition: status.Condition{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			want: []status.Condition{
				{Type: status.ReadyType, Status: corev1.ConditionTrue, Message: "default"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := InitCondition(tc.resource, tc.defaultCondition)
			// ignore LastTransitionTime
			for i := range got {
				got[i].LastTransitionTime = metav1.Time{}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got conditions %+v, want %+v", got, tc.want)
			}
		})
	}
}
