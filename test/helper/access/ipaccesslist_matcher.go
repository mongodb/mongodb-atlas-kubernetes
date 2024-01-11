package access

import (
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go.mongodb.org/atlas-sdk/v20231115003/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

// MatchIPAccessList returns the GomegaMatcher that checks if the 'actual' mongodbatlas.ProjectIPAccessList matches
// the 'expected' mdbv1.ProjectIPAccessList  one.
// Note, that we cannot compare them by all the fields as Atlas tends to set default fields after IPAccessList creation
// so we need to compare only the fields that the Operator has set
func MatchIPAccessList(expected project.IPAccessList) types.GomegaMatcher {
	return &ipAccessListMatcher{ExpectedIPAccessList: expected}
}

type ipAccessListMatcher struct {
	ExpectedIPAccessList project.IPAccessList
}

func (m *ipAccessListMatcher) Match(actual interface{}) (success bool, err error) {
	var c admin.NetworkPermissionEntry
	var ok bool
	if c, ok = actual.(admin.NetworkPermissionEntry); !ok {
		panic("Expected mongodbatlas.ProjectIPAccessList")
	}
	if m.ExpectedIPAccessList.CIDRBlock != "" && c.GetCidrBlock() != m.ExpectedIPAccessList.CIDRBlock {
		return false, nil
	}
	if m.ExpectedIPAccessList.AwsSecurityGroup != "" && c.GetAwsSecurityGroup() != m.ExpectedIPAccessList.AwsSecurityGroup {
		return false, nil
	}
	if m.ExpectedIPAccessList.IPAddress != "" && c.GetIpAddress() != m.ExpectedIPAccessList.IPAddress {
		return false, nil
	}
	if m.ExpectedIPAccessList.Comment != "" && c.GetComment() != m.ExpectedIPAccessList.Comment {
		return false, nil
	}
	if m.ExpectedIPAccessList.DeleteAfterDate != "" {
		expected, err := timeutil.ParseISO8601(m.ExpectedIPAccessList.DeleteAfterDate)
		if err != nil {
			return false, err
		}
		return expected.Unix() == c.GetDeleteAfterDate().Unix(), nil
	}
	return true, nil
}

func (m *ipAccessListMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match", m.ExpectedIPAccessList)
}

func (m *ipAccessListMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match", m.ExpectedIPAccessList)
}
