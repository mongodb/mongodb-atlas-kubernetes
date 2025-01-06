package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//+k8s:deepcopy-gen=false

// AtlasCustomResource is the interface common for all Atlas entities
type AtlasCustomResource interface {
	metav1.Object
	runtime.Object
	Reader
	Writer
}

// InitCondition initializes the underlying type of the given condition to the given default value
// if the underlying condition type is unset.
func InitCondition(resource AtlasCustomResource, defaultCondition Condition) []Condition {
	conditions := resource.GetStatus().GetConditions()
	if !HasConditionType(defaultCondition.Type, conditions) {
		return EnsureConditionExists(defaultCondition, conditions)
	}
	return conditions
}
