package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

type teamsClient struct {
	GetFunc func() (*mongodbatlas.Team, *mongodbatlas.Response, error)
}

func (c *teamsClient) List(_ context.Context, _ string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.Team, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) Get(_ context.Context, _ string, _ string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	return c.GetFunc()
}

func (c *teamsClient) GetOneTeamByName(_ context.Context, _ string, _ string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) GetTeamUsersAssigned(_ context.Context, _ string, _ string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) Create(_ context.Context, _ string, _ *mongodbatlas.Team) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) Rename(_ context.Context, _ string, _ string, _ string) (*mongodbatlas.Team, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) UpdateTeamRoles(_ context.Context, _ string, _ string, _ *mongodbatlas.TeamUpdateRoles) ([]mongodbatlas.TeamRoles, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) AddUsersToTeam(_ context.Context, _ string, _ string, _ []string) ([]mongodbatlas.AtlasUser, *mongodbatlas.Response, error) {
	return nil, nil, nil
}

func (c *teamsClient) RemoveUserToTeam(_ context.Context, _ string, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func (c *teamsClient) RemoveTeamFromOrganization(_ context.Context, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func (c *teamsClient) RemoveTeamFromProject(_ context.Context, _ string, _ string) (*mongodbatlas.Response, error) {
	return nil, nil
}

func TestTeamManagedByAtlas(t *testing.T) {
	t.Run("should return error when passing wrong resource", func(t *testing.T) {
		checker := teamsManagedByAtlas(context.TODO(), mongodbatlas.Client{}, "orgID")
		result, err := checker(&v1.AtlasProject{})
		assert.EqualError(t, err, "failed to match resource type as AtlasTeams")
		assert.False(t, result)
	})

	t.Run("should return false when resource has no Atlas Team ID", func(t *testing.T) {
		checker := teamsManagedByAtlas(context.TODO(), mongodbatlas.Client{}, "orgID")
		result, err := checker(&v1.AtlasTeam{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when resource was not found in Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &teamsClient{
				GetFunc: func() (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound}
				},
			},
		}
		team := &v1.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		checker := teamsManagedByAtlas(context.TODO(), atlasClient, "orgID")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return error when failed to fetch the team from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &teamsClient{
				GetFunc: func() (*mongodbatlas.Team, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, errors.New("unavailable")
				},
			},
		}
		team := &v1.AtlasTeam{
			Status: status.TeamStatus{
				ID: "team-id-1",
			},
		}
		checker := teamsManagedByAtlas(context.TODO(), atlasClient, "orgID")
		result, err := checker(team)
		assert.EqualError(t, err, "unavailable")
		assert.False(t, result)
	})

	t.Run("should return false when resource are equal", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &teamsClient{
				GetFunc: func() (*mongodbatlas.Team, *mongodbatlas.Response, error) {
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
		checker := teamsManagedByAtlas(context.TODO(), atlasClient, "orgID-1")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true when resource are different", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Teams: &teamsClient{
				GetFunc: func() (*mongodbatlas.Team, *mongodbatlas.Response, error) {
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
		checker := teamsManagedByAtlas(context.TODO(), atlasClient, "orgID-1")
		result, err := checker(team)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}
