package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	atlas_mock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestTeamManagedByAtlas(t *testing.T) {
	t.Run("should return error when passing wrong resource", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID")
		result, err := checker(&v1.AtlasProject{})
		assert.EqualError(t, err, "failed to match resource type as AtlasTeams")
		assert.False(t, result)
	})

	t.Run("should return false when resource has no Atlas Team ID", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  mongodbatlas.Client{},
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID")
		result, err := checker(&v1.AtlasTeam{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when resource was not found in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &atlas_mock.TeamsClientMock{
				GetFunc: func(orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound}
				},
			},
		}
		team := &v1.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return error when failed to fetch the team from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &atlas_mock.TeamsClientMock{
				GetFunc: func(orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("unavailable")
				},
			},
		}
		team := &v1.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID")
		result, err := checker(team)
		assert.EqualError(t, err, "unavailable")
		assert.False(t, result)
	})

	t.Run("should return false when resource are equal", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &atlas_mock.TeamsClientMock{
				GetFunc: func(orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return &mongodbatlas.Team{
						ID:        "team-id-1",
						Name:      "My Team",
						Usernames: []string{"user1@mongodb.com", "user2@mongodb.com"},
					}, &mongodbatlas.Response{}, nil
				},
			},
		}
		team := &v1.AtlasTeam{
			Spec: v1.TeamSpec{
				Name:      "My Team",
				Usernames: []v1.TeamUser{"user1@mongodb.com", "user2@mongodb.com"},
			},
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID-1")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when resource are different", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &atlas_mock.TeamsClientMock{
				GetFunc: func(orgID string, teamID string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return &mongodbatlas.Team{
						ID:        "team-id-1",
						Name:      "My Team",
						Usernames: []string{"user1@mongodb.com", "user2@mongodb.com"},
					}, &mongodbatlas.Response{}, nil
				},
			},
		}
		team := &v1.AtlasTeam{
			Spec: v1.TeamSpec{
				Name:      "My Team",
				Usernames: []v1.TeamUser{"user1@mongodb.com"},
			},
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		workflowCtx := &workflow.Context{
			Client:  atlasClient,
			Context: context.TODO(),
		}
		checker := teamsManagedByAtlas(workflowCtx, "orgID-1")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}
