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

package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func TestTagsInSync(t *testing.T) {
	tests := map[string]struct {
		desired  []*akov2.TagSpec
		actual   []*akov2.TagSpec
		expected bool
	}{
		"both empty are in sync": {
			desired:  []*akov2.TagSpec{},
			actual:   []*akov2.TagSpec{},
			expected: true,
		},
		"nil and empty are in sync": {
			desired:  []*akov2.TagSpec{},
			actual:   nil,
			expected: true,
		},
		"identical tags are in sync": {
			desired:  []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			actual:   []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expected: true,
		},
		"tags differing only in order are in sync": {
			desired:  []*akov2.TagSpec{{Key: "env", Value: "dev"}, {Key: "team", Value: "ako"}},
			actual:   []*akov2.TagSpec{{Key: "team", Value: "ako"}, {Key: "env", Value: "dev"}},
			expected: true,
		},
		"a changed value is out of sync": {
			desired:  []*akov2.TagSpec{{Key: "env", Value: "prod"}},
			actual:   []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expected: false,
		},
		"a renamed key is out of sync": {
			desired:  []*akov2.TagSpec{{Key: "environment", Value: "dev"}},
			actual:   []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expected: false,
		},
		"an added tag is out of sync": {
			desired:  []*akov2.TagSpec{{Key: "env", Value: "dev"}, {Key: "team", Value: "ako"}},
			actual:   []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expected: false,
		},
		"a removed tag is out of sync": {
			desired:  []*akov2.TagSpec{},
			actual:   []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expected: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tagsInSync(tt.desired, tt.actual))
		})
	}
}

func TestSyncProjectTags(t *testing.T) {
	tests := map[string]struct {
		specTags       []*akov2.TagSpec
		atlasTags      []*akov2.TagSpec
		updateErr      error
		expectedUpdate []*akov2.TagSpec
		expectedErr    string
	}{
		"should not update when the spec does not manage tags": {
			specTags:  nil,
			atlasTags: []*akov2.TagSpec{{Key: "env", Value: "dev"}},
		},
		"should not update when tags are already in sync": {
			specTags:  []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			atlasTags: []*akov2.TagSpec{{Key: "env", Value: "dev"}},
		},
		"should not update when tags differ only in order": {
			specTags:  []*akov2.TagSpec{{Key: "env", Value: "dev"}, {Key: "team", Value: "ako"}},
			atlasTags: []*akov2.TagSpec{{Key: "team", Value: "ako"}, {Key: "env", Value: "dev"}},
		},
		"should update when a tag value changed": {
			specTags:       []*akov2.TagSpec{{Key: "env", Value: "prod"}},
			atlasTags:      []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expectedUpdate: []*akov2.TagSpec{{Key: "env", Value: "prod"}},
		},
		"should clear tags when the spec has an empty list": {
			specTags:       []*akov2.TagSpec{},
			atlasTags:      []*akov2.TagSpec{{Key: "env", Value: "dev"}},
			expectedUpdate: []*akov2.TagSpec{},
		},
		"should return the error when the update fails": {
			specTags:       []*akov2.TagSpec{{Key: "env", Value: "prod"}},
			atlasTags:      nil,
			updateErr:      errors.New("failed to update project"),
			expectedUpdate: []*akov2.TagSpec{{Key: "env", Value: "prod"}},
			expectedErr:    "failed to update project",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			service := translation.NewProjectServiceMock(t)
			if tt.expectedUpdate != nil {
				service.EXPECT().UpdateProject(context.Background(), &project.Project{
					OrgID: "my-org-id",
					ID:    "projectID",
					Name:  "my-project",
					Tags:  tt.expectedUpdate,
				}).Return(tt.updateErr).Once()
			}

			atlasProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{Name: "my-project", Namespace: "default"},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
					Tags: tt.specTags,
				},
			}
			projectInAtlas := &project.Project{ID: "projectID", Name: "my-project", Tags: tt.atlasTags}
			ctx := &workflow.Context{Context: context.Background()}

			err := (&AtlasProjectReconciler{}).syncProjectTags(ctx, "my-org-id", atlasProject, projectInAtlas, service)
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
