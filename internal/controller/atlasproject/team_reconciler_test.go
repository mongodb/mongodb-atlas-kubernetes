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
	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
)

func TestTeamManagedByAtlas(t *testing.T) {
	var r AtlasProjectReconciler
	t.Run("should return error when passing wrong resource", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			OrgID: "orgID",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{},
			},
			Context: context.Background(),
		}
		checker := r.teamsManagedByAtlas(workflowCtx, translation.NewTeamsServiceMock(t))
		result, err := checker(&akov2.AtlasProject{})
		assert.EqualError(t, err, "failed to match resource type as AtlasTeams")
		assert.False(t, result)
	})

	t.Run("should return false when resource has no Atlas Team ID", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			OrgID: "orgID",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{},
			},
			Context: context.Background(),
		}
		checker := r.teamsManagedByAtlas(workflowCtx, translation.NewTeamsServiceMock(t))
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
			OrgID: "orgID",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &atlasClient,
			},
			Context: context.Background(),
		}
		teamService := func() teams.TeamsService {
			notFound := &admin.GenericOpenAPIError{}
			notFound.SetModel(admin.ApiError{ErrorCode: atlas.ResourceNotFound})

			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(nil, notFound)
			return service
		}
		r = AtlasProjectReconciler{}
		checker := r.teamsManagedByAtlas(workflowCtx, teamService())
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
			OrgID: "orgID",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &atlasClient,
			},
			Context: context.Background(),
		}
		teamService := func() teams.TeamsService {
			service := translation.NewTeamsServiceMock(t)
			service.EXPECT().GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, "team-id-1").
				Return(nil, errors.New("unavailable"))
			return service
		}
		r = AtlasProjectReconciler{}
		checker := r.teamsManagedByAtlas(workflowCtx, teamService())
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
			OrgID: "orgID-1",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &atlasClient,
			},
			Context: context.Background(),
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
		r = AtlasProjectReconciler{}
		checker := r.teamsManagedByAtlas(workflowCtx, teamService())
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
			OrgID: "orgID-1",
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &atlasClient,
			},
			Context: context.Background(),
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
		r = AtlasProjectReconciler{}
		checker := r.teamsManagedByAtlas(workflowCtx, teamService())
		result, err := checker(team)
		assert.NoError(t, err)
		assert.True(t, result)
	})
}
