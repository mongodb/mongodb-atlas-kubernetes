// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
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
