package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

func TestReleaseVersion(t *testing.T) {
	testCases := []struct {
		title    string
		version  string
		expected bool
	}{
		{
			"default is not release",
			version.DefaultVersion,
			false,
		},
		{
			"empty is not release",
			"",
			false,
		},
		{
			"dirty is not release",
			"1.8.0-30-g81233c6-dirty",
			false,
		},
		{
			"semver IS release",
			"1.8.0",
			true,
		},
		{
			"semver certified IS release",
			"1.8.0-certified",
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			assert.Equal(t, tc.expected, version.IsRelease(tc.version))
		})
	}
}
