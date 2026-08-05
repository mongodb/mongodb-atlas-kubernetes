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

package project

import (
	"testing"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestNewProjectTags(t *testing.T) {
	tests := map[string]struct {
		tags     []*akov2.TagSpec
		expected []*akov2.TagSpec
	}{
		"should keep tags unset when the spec has no tags": {
			tags:     nil,
			expected: nil,
		},
		"should keep an empty tag list distinct from unset": {
			tags:     []*akov2.TagSpec{},
			expected: []*akov2.TagSpec{},
		},
		"should map tags from the spec": {
			tags:     []*akov2.TagSpec{{Key: "test", Value: "AKO"}},
			expected: []*akov2.TagSpec{{Key: "test", Value: "AKO"}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			atlasProject := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
					Tags: tt.tags,
				},
			}

			assert.Equal(t, tt.expected, NewProject(atlasProject, "my-org-id").Tags)
		})
	}
}
