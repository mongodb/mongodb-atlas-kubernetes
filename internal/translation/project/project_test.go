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
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestGetProjectByName(t *testing.T) {
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: "NOT_IN_GROUP"})
	tests := map[string]struct {
		api      func() admin.ProjectsApi
		name     string
		expected *Project
		err      string
	}{
		"should fail to retrieve project from atlas": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetGroupByName(context.Background(), "my-project").
					Return(admin.GetGroupByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetGroupByNameExecute(mock.AnythingOfType("admin.GetGroupByNameApiRequest")).
					Return(nil, &http.Response{}, errors.New("fail to retrieve project from atlas"))

				return sdk
			},
			name: "my-project",
			err:  "fail to retrieve project from atlas",
		},
		"should return error when project was not found": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetGroupByName(context.Background(), "my-project").
					Return(admin.GetGroupByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetGroupByNameExecute(mock.AnythingOfType("admin.GetGroupByNameApiRequest")).
					Return(nil, &http.Response{}, notFoundErr)

				return sdk
			},
			name: "my-project",
			err:  "not found\n",
		},
		"should return project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetGroupByName(context.Background(), "my-project").
					Return(admin.GetGroupByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetGroupByNameExecute(mock.AnythingOfType("admin.GetGroupByNameApiRequest")).
					Return(
						&admin.Group{
							OrgId:                     "my-org-id",
							Id:                        pointer.MakePtr("my-project-id"),
							Name:                      "my-project",
							ClusterCount:              0,
							RegionUsageRestrictions:   pointer.MakePtr("NONE"),
							WithDefaultAlertsSettings: pointer.MakePtr(true),
							Tags: &[]admin.ResourceTag{
								{
									Key:   "test",
									Value: "AKO",
								},
							},
						},
						&http.Response{},
						nil,
					)

				return sdk
			},
			name: "my-project",
			expected: &Project{
				OrgID:                     "my-org-id",
				ID:                        "my-project-id",
				Name:                      "my-project",
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				Tags: []*akov2.TagSpec{
					{
						Key:   "test",
						Value: "AKO",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			service := &ProjectAPI{
				projectAPI: tt.api(),
			}
			p, err := service.GetProjectByName(context.Background(), tt.name)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tt.err, gotErr)
			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestCreateProject(t *testing.T) {
	tests := map[string]struct {
		api      func() admin.ProjectsApi
		project  *Project
		expected *Project
		err      error
	}{
		"should fail to create project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().CreateGroup(context.Background(), mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateGroupApiRequest{ApiService: sdk})
				sdk.EXPECT().CreateGroupExecute(mock.AnythingOfType("admin.CreateGroupApiRequest")).
					Return(nil, &http.Response{}, errors.New("fail to create project"))

				return sdk
			},
			project: &Project{
				Name: "my-project",
			},
			err: errors.New("fail to create project"),
		},
		"should create project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().CreateGroup(context.Background(), mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateGroupApiRequest{ApiService: sdk})
				sdk.EXPECT().CreateGroupExecute(mock.AnythingOfType("admin.CreateGroupApiRequest")).
					Return(
						&admin.Group{
							OrgId:                     "my-org-id",
							Id:                        pointer.MakePtr("my-project-id"),
							Name:                      "my-project",
							ClusterCount:              0,
							RegionUsageRestrictions:   pointer.MakePtr("NONE"),
							WithDefaultAlertsSettings: pointer.MakePtr(true),
							Tags: &[]admin.ResourceTag{
								{
									Key:   "test",
									Value: "AKO",
								},
							},
						},
						&http.Response{},
						nil,
					)

				return sdk
			},
			project: &Project{
				Name:                      "my-project",
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				Tags: []*akov2.TagSpec{
					{
						Key:   "test",
						Value: "AKO",
					},
				},
			},
			expected: &Project{
				OrgID:                     "my-org-id",
				ID:                        "my-project-id",
				Name:                      "my-project",
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				Tags: []*akov2.TagSpec{
					{
						Key:   "test",
						Value: "AKO",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			service := &ProjectAPI{
				projectAPI: tt.api(),
			}
			err := service.CreateProject(context.Background(), tt.project)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: "GROUP_NOT_FOUND"})
	tests := map[string]struct {
		api     func() admin.ProjectsApi
		project *Project
		err     error
	}{
		"should fail to delete project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().DeleteGroup(context.Background(), "my-project-id").
					Return(admin.DeleteGroupApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteGroupExecute(mock.AnythingOfType("admin.DeleteGroupApiRequest")).
					Return(&http.Response{}, errors.New("fail to delete project"))

				return sdk
			},
			project: &Project{
				ID: "my-project-id",
			},
			err: errors.New("fail to delete project"),
		},
		"should succeed when project doesn't exist": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().DeleteGroup(context.Background(), "my-project-id").
					Return(admin.DeleteGroupApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteGroupExecute(mock.AnythingOfType("admin.DeleteGroupApiRequest")).
					Return(&http.Response{}, notFoundErr)

				return sdk
			},
			project: &Project{
				ID: "my-project-id",
			},
		},
		"should delete project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().DeleteGroup(context.Background(), "my-project-id").
					Return(admin.DeleteGroupApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteGroupExecute(mock.AnythingOfType("admin.DeleteGroupApiRequest")).
					Return(&http.Response{}, nil)

				return sdk
			},
			project: &Project{
				OrgID:                     "my-org-id",
				ID:                        "my-project-id",
				Name:                      "my-project",
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				Tags: []*akov2.TagSpec{
					{
						Key:   "test",
						Value: "AKO",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			service := &ProjectAPI{
				projectAPI: tt.api(),
			}
			err := service.DeleteProject(context.Background(), tt.project)
			require.Equal(t, tt.err, err)
		})
	}
}
