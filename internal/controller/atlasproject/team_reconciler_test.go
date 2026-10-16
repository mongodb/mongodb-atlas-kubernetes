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
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
	"go.mongodb.org/atlas-sdk/v20250312026/mockadmin"
	"go.uber.org/zap/zaptest"

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
				SdkClient20250312: &admin.APIClient{},
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
				SdkClient20250312: &admin.APIClient{},
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
				SdkClient20250312: &atlasClient,
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
				SdkClient20250312: &atlasClient,
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
				SdkClient20250312: &atlasClient,
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
				SdkClient20250312: &atlasClient,
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

func TestEnsureTeamUsersAreInSync(t *testing.T) {
	const (
		orgID    = "org-id"
		teamID   = "team-id"
		existing = "existing-user@example.com"
		added    = "new-user@example.com"
	)

	team := func(usernames ...string) *akov2.AtlasTeam {
		spec := make([]akov2.TeamUser, 0, len(usernames))
		for _, username := range usernames {
			spec = append(spec, akov2.TeamUser(username))
		}
		return &akov2.AtlasTeam{Spec: akov2.TeamSpec{Name: "My team", Usernames: spec}}
	}

	// Resolves the username and hands back an identifier, the way Atlas does.
	resolvesUsername := func(username, userID string) func(*mockadmin.MongoDBCloudUsersAPI) {
		return func(usersAPI *mockadmin.MongoDBCloudUsersAPI) {
			usersAPI.EXPECT().GetUserByName(mock.Anything, username).
				Return(admin.GetUserByNameApiRequest{ApiService: usersAPI})
			usersAPI.EXPECT().GetUserByNameExecute(mock.Anything).
				Return(&admin.CloudAppUser{Id: new(userID), Username: username}, nil, nil)
		}
	}

	tests := map[string]struct {
		team           *akov2.AtlasTeam
		usersAPI       func(usersAPI *mockadmin.MongoDBCloudUsersAPI)
		teamsService   func(s *translation.TeamsServiceMock)
		expectedResult workflow.DeprecatedResult
	}{
		"adds a user that is not yet a team member": {
			team:     team(existing, added),
			usersAPI: resolvesUsername(added, "new-user-id"),
			teamsService: func(s *translation.TeamsServiceMock) {
				s.EXPECT().GetTeamUsers(mock.Anything, orgID, teamID).
					Return([]teams.TeamUser{{Username: existing, UserID: "existing-user-id"}}, nil)
				s.EXPECT().AddUsers(mock.Anything, &[]teams.TeamUser{{UserID: "new-user-id"}}, orgID, teamID).
					Return(nil)
			},
			expectedResult: workflow.OK(),
		},
		"does not look up users that are already team members": {
			team: team(existing),
			teamsService: func(s *translation.TeamsServiceMock) {
				s.EXPECT().GetTeamUsers(mock.Anything, orgID, teamID).
					Return([]teams.TeamUser{{Username: existing, UserID: "existing-user-id"}}, nil)
			},
			expectedResult: workflow.OK(),
		},
		"removes a team member that left the spec": {
			team: team(existing),
			teamsService: func(s *translation.TeamsServiceMock) {
				s.EXPECT().GetTeamUsers(mock.Anything, orgID, teamID).
					Return([]teams.TeamUser{
						{Username: existing, UserID: "existing-user-id"},
						{Username: "gone@example.com", UserID: "gone-user-id"},
					}, nil)
				s.EXPECT().RemoveUser(mock.Anything, orgID, teamID, "gone-user-id").Return(nil)
			},
			expectedResult: workflow.OK(),
		},
		"terminates when the username cannot be resolved": {
			team: team(added),
			usersAPI: func(usersAPI *mockadmin.MongoDBCloudUsersAPI) {
				usersAPI.EXPECT().GetUserByName(mock.Anything, added).
					Return(admin.GetUserByNameApiRequest{ApiService: usersAPI})
				usersAPI.EXPECT().GetUserByNameExecute(mock.Anything).
					Return(nil, nil, admin.GenericOpenAPIError{})
			},
			teamsService: func(s *translation.TeamsServiceMock) {
				s.EXPECT().GetTeamUsers(mock.Anything, orgID, teamID).Return([]teams.TeamUser{}, nil)
			},
			expectedResult: workflow.Terminate(workflow.TeamUsersNotReady, admin.GenericOpenAPIError{}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			usersAPI := mockadmin.NewMongoDBCloudUsersAPI(t)
			if tc.usersAPI != nil {
				tc.usersAPI(usersAPI)
			}
			teamsService := translation.NewTeamsServiceMock(t)
			tc.teamsService(teamsService)

			r := &AtlasProjectReconciler{}
			workflowCtx := &workflow.Context{
				OrgID:   orgID,
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312: &admin.APIClient{MongoDBCloudUsersAPI: usersAPI},
				},
			}

			result := r.ensureTeamUsersAreInSync(workflowCtx, teamsService, teamID, tc.team)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestEnsureTeamUsersAreInSyncResolvesUsernamesByName(t *testing.T) {
	const (
		orgID    = "org-id"
		teamID   = "team-id"
		username = "new-user@example.com"
	)

	usersAPI := mockadmin.NewMongoDBCloudUsersAPI(t)
	usersAPI.EXPECT().GetUserByName(mock.Anything, username).
		Return(admin.GetUserByNameApiRequest{ApiService: usersAPI})
	usersAPI.EXPECT().GetUserByNameExecute(mock.Anything).
		Return(&admin.CloudAppUser{Id: new("new-user-id"), Username: username}, nil, nil)

	teamsService := translation.NewTeamsServiceMock(t)
	teamsService.EXPECT().GetTeamUsers(mock.Anything, orgID, teamID).Return([]teams.TeamUser{}, nil)
	teamsService.EXPECT().AddUsers(mock.Anything, &[]teams.TeamUser{{UserID: "new-user-id"}}, orgID, teamID).
		Return(nil)

	r := &AtlasProjectReconciler{}
	workflowCtx := &workflow.Context{
		OrgID:   orgID,
		Context: context.Background(),
		Log:     zaptest.NewLogger(t).Sugar(),
		SdkClientSet: &atlas.ClientSet{
			SdkClient20250312: &admin.APIClient{MongoDBCloudUsersAPI: usersAPI},
		},
	}

	result := r.ensureTeamUsersAreInSync(workflowCtx, teamsService, teamID,
		&akov2.AtlasTeam{Spec: akov2.TeamSpec{
			Name:      "My team",
			Usernames: []akov2.TeamUser{username},
		}})

	assert.Equal(t, workflow.OK(), result)
	usersAPI.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
}
