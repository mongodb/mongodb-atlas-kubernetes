package conditions

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// MatchCondition returns the GomegaMatcher that checks if the 'actual' status.Condition matches the 'expected' one.
func MatchCondition(expected api.Condition) types.GomegaMatcher {
	return &conditionMatcher{ExpectedCondition: expected}
}

// MatchConditions is a convenience method that allows to create the range of matchers simplifying testing
func MatchConditions(expected ...api.Condition) []types.GomegaMatcher {
	result := make([]types.GomegaMatcher, len(expected))
	for i, c := range expected {
		result[i] = MatchCondition(c)
	}
	return result
}

type conditionMatcher struct {
	ExpectedCondition api.Condition
}

func (m *conditionMatcher) Match(actual interface{}) (success bool, err error) {
	var c api.Condition
	var ok bool
	if c, ok = actual.(api.Condition); !ok {
		panic("Expected Condition")
	}
	if m.ExpectedCondition.Reason != "" && c.Reason != m.ExpectedCondition.Reason {
		return false, nil
	}
	if m.ExpectedCondition.Status != "" && c.Status != m.ExpectedCondition.Status {
		return false, nil
	}
	if m.ExpectedCondition.Type != "" && c.Type != m.ExpectedCondition.Type {
		return false, nil
	}
	if m.ExpectedCondition.Message != "" {
		gomega.Expect(c.Message).To(gomega.MatchRegexp(m.ExpectedCondition.Message))
	}

	return true, nil
}

func (m *conditionMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match", m.ExpectedCondition)
}

func (m *conditionMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match", m.ExpectedCondition)
}

func FindConditionByType(conditions []api.Condition, conditionType api.ConditionType) (api.Condition, bool) {
	for _, c := range conditions {
		if c.Type == conditionType {
			return c, true
		}
	}
	return api.Condition{}, false
}
