package project

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestGetProjectByName(t *testing.T) {
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("NOT_IN_GROUP")})
	tests := map[string]struct {
		api      func() admin.ProjectsApi
		name     string
		expected *Project
		err      error
	}{
		"should fail to retrieve project from atlas": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(admin.GetProjectByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetProjectByNameExecute(mock.AnythingOfType("admin.GetProjectByNameApiRequest")).
					Return(nil, &http.Response{}, errors.New("fail to retrieve project from atlas"))

				return sdk
			},
			name: "my-project",
			err:  errors.New("fail to retrieve project from atlas"),
		},
		"should return nil when project was not found": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(admin.GetProjectByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetProjectByNameExecute(mock.AnythingOfType("admin.GetProjectByNameApiRequest")).
					Return(nil, &http.Response{}, notFoundErr)

				return sdk
			},
			name: "my-project",
		},
		"should return project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(admin.GetProjectByNameApiRequest{ApiService: sdk})
				sdk.EXPECT().GetProjectByNameExecute(mock.AnythingOfType("admin.GetProjectByNameApiRequest")).
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
			require.Equal(t, tt.err, err)
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
				sdk.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateProjectApiRequest{ApiService: sdk})
				sdk.EXPECT().CreateProjectExecute(mock.AnythingOfType("admin.CreateProjectApiRequest")).
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
				sdk.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateProjectApiRequest{ApiService: sdk})
				sdk.EXPECT().CreateProjectExecute(mock.AnythingOfType("admin.CreateProjectApiRequest")).
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
	notFoundErr.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("GROUP_NOT_FOUND")})
	tests := map[string]struct {
		api     func() admin.ProjectsApi
		project *Project
		err     error
	}{
		"should fail to delete project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().DeleteProject(context.Background(), "my-project-id").
					Return(admin.DeleteProjectApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteProjectExecute(mock.AnythingOfType("admin.DeleteProjectApiRequest")).
					Return(nil, &http.Response{}, errors.New("fail to delete project"))

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
				sdk.EXPECT().DeleteProject(context.Background(), "my-project-id").
					Return(admin.DeleteProjectApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteProjectExecute(mock.AnythingOfType("admin.DeleteProjectApiRequest")).
					Return(nil, &http.Response{}, notFoundErr)

				return sdk
			},
			project: &Project{
				ID: "my-project-id",
			},
		},
		"should delete project": {
			api: func() admin.ProjectsApi {
				sdk := mockadmin.NewProjectsApi(t)
				sdk.EXPECT().DeleteProject(context.Background(), "my-project-id").
					Return(admin.DeleteProjectApiRequest{ApiService: sdk})
				sdk.EXPECT().DeleteProjectExecute(mock.AnythingOfType("admin.DeleteProjectApiRequest")).
					Return(nil, &http.Response{}, nil)

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
