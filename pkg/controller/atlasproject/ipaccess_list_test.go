package atlasproject

import (
	"testing"

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
