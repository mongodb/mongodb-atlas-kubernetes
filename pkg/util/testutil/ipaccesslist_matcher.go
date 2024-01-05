// TODO: move away from pkg, this code is only usable from tests
package testutil

import (
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/timeutil"
)

// MatchIPAccessList returns the GomegaMatcher that checks if the 'actual' mongodbatlas.ProjectIPAccessList matches
// the 'expected' mdbv1.ProjectIPAccessList  one.
// Note, that we cannot compare them by all the fields as Atlas tends to set default fields after IPAccessList creation
// so we need to compare only the fields that the Operator has set
func MatchIPAccessList(expected project.IPAccessList) types.GomegaMatcher {
	return &ipAccessListMatcher{ExpectedIPAccessList: expected}
}

func BuildMatchersFromExpected(ipLists []project.IPAccessList) []types.GomegaMatcher {
	result := make([]types.GomegaMatcher, len(ipLists))
	for i, list := range ipLists {
		result[i] = MatchIPAccessList(list)
	}
	return result
}

type ipAccessListMatcher struct {
	ExpectedIPAccessList project.IPAccessList
}

func (m *ipAccessListMatcher) Match(actual interface{}) (success bool, err error) {
	var c mongodbatlas.ProjectIPAccessList
	var ok bool
	if c, ok = actual.(mongodbatlas.ProjectIPAccessList); !ok {
		panic("Expected mongodbatlas.ProjectIPAccessList")
	}
	if m.ExpectedIPAccessList.CIDRBlock != "" && c.CIDRBlock != m.ExpectedIPAccessList.CIDRBlock {
		return false, nil
	}
	if m.ExpectedIPAccessList.AwsSecurityGroup != "" && c.AwsSecurityGroup != m.ExpectedIPAccessList.AwsSecurityGroup {
		return false, nil
	}
	if m.ExpectedIPAccessList.IPAddress != "" && c.IPAddress != m.ExpectedIPAccessList.IPAddress {
		return false, nil
	}
	if m.ExpectedIPAccessList.Comment != "" && c.Comment != m.ExpectedIPAccessList.Comment {
		return false, nil
	}
	if m.ExpectedIPAccessList.DeleteAfterDate != "" {
		expected, err := timeutil.ParseISO8601(m.ExpectedIPAccessList.DeleteAfterDate)
		if err != nil {
			return false, err
		}
		fromAtlas, err := timeutil.ParseISO8601(c.DeleteAfterDate)
		if err != nil {
			return false, err
		}
		return expected.Unix() == fromAtlas.Unix(), nil
	}
	return true, nil
}

func (m *ipAccessListMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match", m.ExpectedIPAccessList)
}

func (m *ipAccessListMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match", m.ExpectedIPAccessList)
}
