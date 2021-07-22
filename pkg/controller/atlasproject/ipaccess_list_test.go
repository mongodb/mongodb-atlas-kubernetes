package atlasproject

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

func TestValidateSingleIPAccessList(t *testing.T) {
	testCases := []struct {
		in                 project.IPAccessList
		errorExpectedRegex string
	}{
		// Date
		{in: project.IPAccessList{DeleteAfterDate: "incorrect", IPAddress: "192.158.0.0"}, errorExpectedRegex: "cannot parse"},
		{in: project.IPAccessList{DeleteAfterDate: "2020/01/02T15:04:05-0700", IPAddress: "192.158.0.0"}, errorExpectedRegex: "cannot parse"},
		{in: project.IPAccessList{DeleteAfterDate: "2020-01-02T15:04:05-07000", IPAddress: "192.158.0.0"}, errorExpectedRegex: "cannot parse"},
		{in: project.IPAccessList{DeleteAfterDate: "2020-11-02T20:04:05-0700", IPAddress: "192.158.0.0"}},
		{in: project.IPAccessList{DeleteAfterDate: "2020-11-02T20:04:05+03", IPAddress: "192.158.0.0"}},
		{in: project.IPAccessList{DeleteAfterDate: "2011-01-02T15:04:05", IPAddress: "192.158.0.0"}},

		// Only one other field is allowed (and required)
		{in: project.IPAccessList{DeleteAfterDate: "2011-01-02T15:04:05", IPAddress: "192.158.0.0", CIDRBlock: "203.0.113.0/24"}, errorExpectedRegex: "only one of the "},
		{in: project.IPAccessList{IPAddress: "192.158.0.0", AwsSecurityGroup: "sg-0026348ec11780bd1"}, errorExpectedRegex: "only one of the "},
		{in: project.IPAccessList{CIDRBlock: "203.0.113.0/24", AwsSecurityGroup: "sg-0026348ec11780bd1"}, errorExpectedRegex: "only one of the "},
		{in: project.IPAccessList{CIDRBlock: "203.0.113.0/24", AwsSecurityGroup: "sg-0026348ec11780bd1", IPAddress: "192.158.0.0"}, errorExpectedRegex: "only one of the "},
		{in: project.IPAccessList{}, errorExpectedRegex: "only one of the "},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			err := validateSingleIPAccessList(testCase.in)
			if testCase.errorExpectedRegex != "" {
				assert.Error(t, err)
				assert.Regexp(t, testCase.errorExpectedRegex, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFilterActiveIPAccessLists(t *testing.T) {
	t.Run("One expired, one active", func(t *testing.T) {
		dateBefore := time.Now().UTC().Add(time.Hour * -1).Format("2006-01-02T15:04:05.999Z")
		dateAfter := time.Now().UTC().Add(time.Hour * 5).Format("2006-01-02T15:04:05.999Z")
		ipAccessExpired := project.IPAccessList{DeleteAfterDate: dateBefore}
		ipAccessActive := project.IPAccessList{DeleteAfterDate: dateAfter}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive, ipAccessExpired})
		assert.Equal(t, []project.IPAccessList{ipAccessActive}, active)
		assert.Equal(t, []project.IPAccessList{ipAccessExpired}, expired)
	})
	t.Run("Two active", func(t *testing.T) {
		dateAfter1 := time.Now().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		dateAfter2 := time.Now().Add(time.Hour * 5).Format("2006-01-02T15:04:05")
		ipAccessActive1 := project.IPAccessList{DeleteAfterDate: dateAfter1}
		ipAccessActive2 := project.IPAccessList{DeleteAfterDate: dateAfter2}
		active, expired := filterActiveIPAccessLists([]project.IPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Equal(t, active, []project.IPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Empty(t, expired)
	})
}
