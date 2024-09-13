package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestTeamManagedByAtlas(t *testing.T) {
	t.Run("should return error when passing wrong resource", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			OrgID:     "orgID",
			SdkClient: &admin.APIClient{},
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(&akov2.AtlasProject{})
		assert.EqualError(t, err, "failed to match resource type as AtlasTeams")
		assert.False(t, result)
	})

	t.Run("should return false when resource has no Atlas Team ID", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			OrgID:     "orgID",
			SdkClient: &admin.APIClient{},
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(&akov2.AtlasTeam{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when resource was not found in Atlas", func(t *testing.T) {
		atlasClient := admin.APIClient{
			TeamsApi: func() *mockadmin.TeamsApi {
				TeamsApi := mockadmin.NewTeamsApi(t)
				TeamsApi.EXPECT().GetTeamById(context.Background(), "orgID", "team-id-1").
					Return(admin.GetTeamByIdApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().GetTeamByIdExecute(mock.Anything).
					Return(nil, &http.Response{}, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound})
				return TeamsApi
			}(),
		}
		team := &akov2.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			OrgID:     "orgID",
			SdkClient: &atlasClient,
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return error when failed to fetch the team from Atlas", func(t *testing.T) {
		atlasClient := admin.APIClient{
			TeamsApi: func() *mockadmin.TeamsApi {
				TeamsApi := mockadmin.NewTeamsApi(t)
				TeamsApi.EXPECT().GetTeamById(context.Background(), "orgID", "team-id-1").
					Return(admin.GetTeamByIdApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().GetTeamByIdExecute(mock.Anything).
					Return(nil, &http.Response{}, errors.New("unavailable"))
				return TeamsApi
			}(),
		}
		team := &akov2.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			OrgID:     "orgID",
			SdkClient: &atlasClient,
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.EqualError(t, err, "unavailable")
		assert.False(t, result)
	})

	t.Run("should return false when resource are equal", func(t *testing.T) {
		atlasClient := admin.APIClient{
			TeamsApi: func() *mockadmin.TeamsApi {
				TeamsApi := mockadmin.NewTeamsApi(t)
				TeamsApi.EXPECT().GetTeamById(context.Background(), "orgID-1", "team-id-1").
					Return(admin.GetTeamByIdApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().GetTeamByIdExecute(mock.Anything).
					Return(&admin.TeamResponse{
						Id:    func(s string) *string { return &s }("team-id-1"),
						Links: nil,
						Name:  func(s string) *string { return &s }("My Team"),
					}, &http.Response{}, nil)
				TeamsApi.EXPECT().ListTeamUsers(context.Background(), "orgID-1", "My Team").
					Return(admin.ListTeamUsersApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().ListTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedApiAppUser{
						Links: nil,
						Results: &[]admin.CloudAppUser{
							{
								Username: "user1@mongodb.com",
							},
							{
								Username: "user2@mongodb.com",
							},
						},
						TotalCount: nil,
					}, &http.Response{}, nil)
				return TeamsApi
			}(),
		}
		team := &akov2.AtlasTeam{
			Spec: akov2.TeamSpec{
				Name:      "My Team",
				Usernames: []akov2.TeamUser{"user1@mongodb.com", "user2@mongodb.com"},
			},
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			OrgID:     "orgID-1",
			SdkClient: &atlasClient,
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when resource are different", func(t *testing.T) {
		atlasClient := admin.APIClient{
			TeamsApi: func() *mockadmin.TeamsApi {
				TeamsApi := mockadmin.NewTeamsApi(t)
				TeamsApi.EXPECT().GetTeamById(context.Background(), "orgID-1", "team-id-1").
					Return(admin.GetTeamByIdApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().GetTeamByIdExecute(mock.Anything).
					Return(&admin.TeamResponse{
						Id:    func(s string) *string { return &s }("team-id-1"),
						Links: nil,
						Name:  func(s string) *string { return &s }("My Team"),
					}, &http.Response{}, nil)
				TeamsApi.EXPECT().ListTeamUsers(context.Background(), "orgID-1", "My Team").
					Return(admin.ListTeamUsersApiRequest{ApiService: TeamsApi})
				TeamsApi.EXPECT().ListTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedApiAppUser{
						Links: nil,
						Results: &[]admin.CloudAppUser{
							{
								Username: "user1@mongodb.com",
							},
							{
								Username: "user2@mongodb.com",
							},
						},
						TotalCount: nil,
					}, &http.Response{}, nil)
				return TeamsApi
			}(),
		}
		team := &akov2.AtlasTeam{
			Spec: akov2.TeamSpec{
				Name:      "My Team",
				Usernames: []akov2.TeamUser{"user1@mongodb.com"},
			},
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			OrgID:     "orgID-1",
			SdkClient: &atlasClient,
			Context:   context.Background(),
		}
		checker := teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}
