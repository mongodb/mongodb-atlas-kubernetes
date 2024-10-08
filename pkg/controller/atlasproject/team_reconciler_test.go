package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestTeamManagedByAtlas(t *testing.T) {
	var r AtlasProjectReconciler
	t.Run("should return error when passing wrong resource", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			OrgID:     "orgID",
			SdkClient: &admin.APIClient{},
			Context:   context.Background(),
		}
		checker := r.teamsManagedByAtlas(workflowCtx)
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
		checker := r.teamsManagedByAtlas(workflowCtx)
		result, err := checker(&akov2.AtlasTeam{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when resource was not found in Atlas", func(t *testing.T) {
		atlasClient := admin.APIClient{}
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
		teamService := func() teams.TeamsService {
			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound})
			return service
		}
		r = AtlasProjectReconciler{teamsService: teamService()}
		checker := r.teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return error when failed to fetch the team from Atlas", func(t *testing.T) {
		atlasClient := admin.APIClient{}
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
		teamService := func() teams.TeamsService {
			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(nil, errors.New("unavailable"))
			return service
		}
		r = AtlasProjectReconciler{teamsService: teamService()}
		checker := r.teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.EqualError(t, err, "unavailable")
		assert.False(t, result)
	})

	t.Run("should return false when resource are equal", func(t *testing.T) {
		atlasClient := admin.APIClient{}
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
		teamService := func() teams.TeamsService {
			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(&teams.AssignedTeam{
					TeamID:   "team-id-1",
					TeamName: "My Team",
					Roles:    nil,
				}, nil)
			service.EXPECT().GetTeamUsers(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return([]teams.TeamUser{
					{
						Username: "user1@mongodb.com",
					},
					{
						Username: "user2@mongodb.com",
					},
				}, nil)
			return service
		}
		r = AtlasProjectReconciler{teamsService: teamService()}
		checker := r.teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when resource are different", func(t *testing.T) {
		atlasClient := admin.APIClient{}
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
		teamService := func() teams.TeamsService {
			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(&teams.AssignedTeam{
					TeamID:   "team-id-1",
					TeamName: "My Team",
					Roles:    nil,
				}, nil)
			service.EXPECT().GetTeamUsers(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return([]teams.TeamUser{
					{
						Username: "user1@mongodb.com",
					},
					{
						Username: "user2@mongodb.com",
					},
				}, nil)
			return service
		}
		r = AtlasProjectReconciler{teamsService: teamService()}
		checker := r.teamsManagedByAtlas(workflowCtx)
		result, err := checker(team)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}
