package workflow

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

// Status is a mutable container containing the status of some particular reconciliation. It is expected to be updated
// by a controller and at any time it reflects the status of the reconciled resource. Its state should fully match the
// state of the resource so the information could be used to update the status field of the Custom Resource.
type Status struct {
	options    []api.Option
	conditions []api.Condition
}

func NewStatus(conditions []api.Condition) Status {
	return Status{
		conditions: conditions,
	}
}

func (s *Status) EnsureCondition(condition api.Condition) {
	s.conditions = api.EnsureConditionExists(condition, s.conditions)
}

func (s *Status) GetCondition(conditionType api.ConditionType) (condition api.Condition, found bool) {
	for _, condition := range s.conditions {
		if condition.Type == conditionType {
			return condition, true
		}
	}

	return condition, false
}

func (s *Status) RemoveCondition(conditionType api.ConditionType) {
	s.conditions = api.RemoveConditionIfExists(conditionType, s.conditions)
}

func (s *Status) EnsureOption(option api.Option) {
	// Condition not found - appending (the Option of the same type may be appended more than once)
	// Important! This will work only if the function behind the Option always makes the same updates. If there's a
	// conditional logic and different information is updated this means that we may need some logic to replace the
	// option instead of adding (e.g. some "name" inside the Option)
	s.options = append(s.options, option)
}
