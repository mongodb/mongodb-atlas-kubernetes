package state

import (
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StateCondition = "State"
	ReadyCondition = "Ready"
)

type ResourceState string

const (
	StateInitial ResourceState = "Initial"

	StateImportRequested ResourceState = "Importing"
	StateImported        ResourceState = "Imported"

	StateCreating ResourceState = "Creating"
	StateCreated  ResourceState = "Created"

	StateUpdating ResourceState = "Updating"
	StateUpdated  ResourceState = "Updated"

	StateDeletionRequested ResourceState = "DeletionRequested"
	StateDeleting          ResourceState = "Deleting"

	// Note: StateDeleted this is a terminal state.
	// Finalizers will be unset and no state handler will be invoked.
	StateDeleted ResourceState = "Deleted"
)

func GetState(conditions []metav1.Condition) ResourceState {
	c := meta.FindStatusCondition(conditions, StateCondition)
	if c == nil {
		return StateInitial
	}
	return ResourceState(c.Reason)
}

func EnsureState(conditions *[]metav1.Condition, observedGeneration int64, state ResourceState, msg string, status bool) {
	s := metav1.ConditionFalse
	if status {
		s = metav1.ConditionTrue
	}

	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:               "State",
		Status:             s,
		ObservedGeneration: observedGeneration,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Reason:             string(state),
		Message:            msg,
	})
}
