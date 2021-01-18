package workflow

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

// Status is a mutable container containing the status of some particular reconciliation. It is expected to be updated
// by a controller and at any time it reflects the status of the reconciled resource. Its state should fully match the
// state of the resource so the information could be used to update the status field of the Custom Resource.
type Status struct {
	options    []status.Option
	conditions []status.Condition
}

func NewStatus(conditions []status.Condition) Status {
	return Status{
		conditions: conditions,
	}
}

func (s *Status) EnsureCondition(condition status.Condition) {
	s.conditions = status.EnsureConditionExists(condition, s.conditions)
}

func (s *Status) EnsureOption(option status.Option) {
	// Condition not found - appending (the Option of the same type may be appended more than once)
	// Important! This will work only if the function behind the Option always makes the same updates. If there's a
	// conditional logic and different information is updated this means that we may need some logic to replace the
	// option instead of adding (e.g. some "name" inside the Option)
	s.options = append(s.options, option)
}
