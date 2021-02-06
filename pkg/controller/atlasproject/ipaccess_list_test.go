package atlasproject

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestValidateSingleIPAccessList(t *testing.T) {
	testCases := []struct {
		in                 mdbv1.ProjectIPAccessList
		errorExpectedRegex string
	}{
		// Date
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "incorrect", IPAddress: "192.158.0.0"}, errorExpectedRegex: "cannot parse"},
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2020/01/02T15:04:05-0700", IPAddress: "192.158.0.0"}, errorExpectedRegex: "cannot parse"},
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2020-01-02T15:04:05-07000", IPAddress: "192.158.0.0"}, errorExpectedRegex: `extra text: "0"`},
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2020-11-02T20:04:05-0700", IPAddress: "192.158.0.0"}},
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2020-11-02T20:04:05+03", IPAddress: "192.158.0.0"}},
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2011-01-02T15:04:05", IPAddress: "192.158.0.0"}},

		// Only one other field is allowed (and required)
		{in: mdbv1.ProjectIPAccessList{DeleteAfterDate: "2011-01-02T15:04:05", IPAddress: "192.158.0.0", CIDRBlock: "203.0.113.0/24"}, errorExpectedRegex: "only one of the "},
		{in: mdbv1.ProjectIPAccessList{IPAddress: "192.158.0.0", AwsSecurityGroup: "sg-0026348ec11780bd1"}, errorExpectedRegex: "only one of the "},
		{in: mdbv1.ProjectIPAccessList{CIDRBlock: "203.0.113.0/24", AwsSecurityGroup: "sg-0026348ec11780bd1"}, errorExpectedRegex: "only one of the "},
		{in: mdbv1.ProjectIPAccessList{}, errorExpectedRegex: "only one of the "},
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
		dateBefore := time.Now().Add(time.Hour * -1).Format("2006-01-02T15:04:05")
		dateAfter := time.Now().Add(time.Hour * 5).Format("2006-01-02T15:04:05")
		ipAccessExpired := mdbv1.ProjectIPAccessList{DeleteAfterDate: dateBefore}
		ipAccessActive := mdbv1.ProjectIPAccessList{DeleteAfterDate: dateAfter}
		active, expired := filterActiveIPAccessLists([]mdbv1.ProjectIPAccessList{ipAccessActive, ipAccessExpired})
		assert.Equal(t, []mdbv1.ProjectIPAccessList{ipAccessActive}, active)
		assert.Equal(t, []mdbv1.ProjectIPAccessList{ipAccessExpired}, expired)
	})
	t.Run("Two active", func(t *testing.T) {
		dateAfter1 := time.Now().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		dateAfter2 := time.Now().Add(time.Hour * 5).Format("2006-01-02T15:04:05")
		ipAccessActive1 := mdbv1.ProjectIPAccessList{DeleteAfterDate: dateAfter1}
		ipAccessActive2 := mdbv1.ProjectIPAccessList{DeleteAfterDate: dateAfter2}
		active, expired := filterActiveIPAccessLists([]mdbv1.ProjectIPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Equal(t, active, []mdbv1.ProjectIPAccessList{ipAccessActive2, ipAccessActive1})
		assert.Empty(t, expired)
	})
}
