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

package teams

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	"go.mongodb.org/atlas-sdk/v20250312012/mockadmin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	testProjectID = "project-id"
	testOrgID     = "org-id"
)

var (
	testTeamName = "team-name"
	testTeamID1  = "team1-id"
	testTeamID2  = "team2-id"
	testUserID   = "user-id"
)

func TestTeamsAPI_ListProjectTeams(t *testing.T) {
	ctx := context.Background()
	projectID := testProjectID

	tests := []struct {
		title         string
		mock          func(mockTeamAPI *mockadmin.TeamsApi)
		expectedTeams []AssignedTeam
		expectedErr   error
	}{
		{
			title: "Should return empty when Atlas is also empty",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().ListGroupTeams(ctx, projectID).
					Return(admin.ListGroupTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListGroupTeamsExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{}, &http.Response{}, nil)
			},
			expectedErr:   nil,
			expectedTeams: []AssignedTeam{},
		},
		{
			title: "Should return populated team when team is present on Atlas",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().ListGroupTeams(ctx, projectID).
					Return(admin.ListGroupTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListGroupTeamsExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{
						Results: &[]admin.TeamRole{
							{
								RoleNames: []string{"role1", "role2"},
								TeamId:    testTeamID1,
							},
							{
								RoleNames: []string{"role3", "role4"},
								TeamId:    testTeamID2,
							},
						},
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
			expectedTeams: []AssignedTeam{
				{
					Roles:  []string{"role1", "role2"},
					TeamID: testTeamID1,
				},
				{
					Roles:  []string{"role3", "role4"},
					TeamID: testTeamID2,
				},
			},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().ListGroupTeams(ctx, projectID).
					Return(admin.ListGroupTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().ListGroupTeamsExecute(mock.Anything).
					Return(nil, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:   fmt.Errorf("failed to get project team list from Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			teams, err := ts.ListProjectTeams(ctx, projectID)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeams, teams)
		})
	}
}

func TestTeamsAPI_GetTeamByName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title         string
		mock          func(mockTeamAPI *mockadmin.TeamsApi)
		expectedTeams *AssignedTeam
		expectedErr   error
	}{
		{
			title: "Should return team when team is present on Atlas",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().GetTeamByName(ctx, testOrgID, testTeamName).
					Return(admin.GetTeamByNameApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().GetTeamByNameExecute(mock.Anything).
					Return(&admin.TeamResponse{Id: &testTeamID1, Name: &testTeamName}, &http.Response{}, nil)
			},
			expectedErr:   nil,
			expectedTeams: &AssignedTeam{TeamID: testTeamID1, TeamName: testTeamName},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().GetTeamByName(ctx, testOrgID, testTeamName).
					Return(admin.GetTeamByNameApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().GetTeamByNameExecute(mock.Anything).
					Return(&admin.TeamResponse{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:   fmt.Errorf("failed to get team by name from Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeams: nil,
		},
		{
			title: "Should return empty team and no error when 404 http error occurs",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().GetTeamByName(ctx, testOrgID, testTeamName).
					Return(admin.GetTeamByNameApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().GetTeamByNameExecute(mock.Anything).
					Return(&admin.TeamResponse{}, &http.Response{StatusCode: http.StatusNotFound}, admin.GenericOpenAPIError{})
			},
			expectedErr:   nil,
			expectedTeams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			teams, err := ts.GetTeamByName(ctx, testOrgID, testTeamName)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeams, teams)
		})
	}
}

func TestTeamsAPI_GetTeamByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title         string
		mock          func(mockTeamAPI *mockadmin.TeamsApi)
		expectedTeams *AssignedTeam
		expectedErr   error
	}{
		{
			title: "Should return team when team is present on Atlas",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().GetOrgTeam(ctx, testOrgID, testTeamName).
					Return(admin.GetOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().GetOrgTeamExecute(mock.Anything).
					Return(&admin.TeamResponse{Id: &testTeamID1, Name: &testTeamName}, &http.Response{}, nil)
			},
			expectedErr:   nil,
			expectedTeams: &AssignedTeam{TeamID: testTeamID1, TeamName: testTeamName},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().GetOrgTeam(ctx, testOrgID, testTeamName).
					Return(admin.GetOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().GetOrgTeamExecute(mock.Anything).
					Return(&admin.TeamResponse{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:   fmt.Errorf("failed to get team by ID from Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			teams, err := ts.GetTeamByID(ctx, testOrgID, testTeamName)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeams, teams)
		})
	}
}

func TestTeamsAPI_Assign(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title       string
		mock        func(mockTeamAPI *mockadmin.TeamsApi)
		expectedErr error
	}{
		{
			title: "Should assign team to project",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().AddGroupTeams(ctx, mock.Anything, mock.Anything).
					Return(admin.AddGroupTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().AddGroupTeamsExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{
						Results: &[]admin.TeamRole{
							{
								RoleNames: []string{"role1", "role2"},
								TeamId:    testTeamID1,
							},
						},
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().AddGroupTeams(ctx, mock.Anything, mock.Anything).
					Return(admin.AddGroupTeamsApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().AddGroupTeamsExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr: admin.GenericOpenAPIError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			err := ts.Assign(ctx, &[]AssignedTeam{
				{
					Roles:    []string{"role1", "role2"},
					TeamID:   testTeamID1,
					TeamName: testTeamName,
				},
			}, testProjectID)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTeamsAPI_Unassign(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title       string
		mock        func(mockTeamAPI *mockadmin.TeamsApi)
		expectedErr error
	}{
		{
			title: "Should assign team to project",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RemoveGroupTeam(ctx, mock.Anything, mock.Anything).
					Return(admin.RemoveGroupTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RemoveGroupTeamExecute(mock.Anything).
					Return(&http.Response{}, nil)
			},
			expectedErr: nil,
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RemoveGroupTeam(ctx, mock.Anything, mock.Anything).
					Return(admin.RemoveGroupTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RemoveGroupTeamExecute(mock.Anything).
					Return(&http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr: admin.GenericOpenAPIError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			err := ts.Unassign(ctx, mock.Anything, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTeamsAPI_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title        string
		mock         func(mockTeamAPI *mockadmin.TeamsApi)
		expectedTeam *Team
		expectedErr  error
	}{
		{
			title: "Should create team",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().CreateOrgTeam(ctx, mock.Anything, mock.Anything).
					Return(admin.CreateOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().CreateOrgTeamExecute(mock.Anything).
					Return(&admin.Team{
						Id:        &testTeamID1,
						Name:      testTeamName,
						Usernames: []string{"user@name.com"},
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
			expectedTeam: &Team{
				TeamID:   testTeamID1,
				TeamName: testTeamName,
			},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().CreateOrgTeam(ctx, mock.Anything, mock.Anything).
					Return(admin.CreateOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().CreateOrgTeamExecute(mock.Anything).
					Return(&admin.Team{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:  fmt.Errorf("failed to create team on Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeam: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			team, err := ts.Create(ctx, &Team{}, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeam, team)
		})
	}
}

func TestTeamsAPI_GetTeamUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title         string
		mock          func(mockUsersAPI *mockadmin.MongoDBCloudUsersApi)
		expectedTeams []TeamUser
		expectedErr   error
	}{
		{
			title: "Should return team when team is present on Atlas",
			mock: func(mockUsersAPI *mockadmin.MongoDBCloudUsersApi) {
				mockUsersAPI.EXPECT().ListTeamUsers(ctx, mock.Anything, mock.Anything).
					Return(admin.ListTeamUsersApiRequest{ApiService: mockUsersAPI})
				mockUsersAPI.EXPECT().ListTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedOrgUser{
						Results: &[]admin.OrgUserResponse{
							{
								Username: "user1",
								Id:       testUserID,
							},
						},
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
			expectedTeams: []TeamUser{
				{
					Username: "user1",
					UserID:   testUserID,
				},
			},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockUsersAPI *mockadmin.MongoDBCloudUsersApi) {
				mockUsersAPI.EXPECT().ListTeamUsers(ctx, mock.Anything, mock.Anything).
					Return(admin.ListTeamUsersApiRequest{ApiService: mockUsersAPI})
				mockUsersAPI.EXPECT().ListTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedOrgUser{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:   fmt.Errorf("failed to get team users from Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeams: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockUsersAPI := mockadmin.NewMongoDBCloudUsersApi(t)
			tt.mock(mockUsersAPI)
			ts := &TeamsAPI{
				teamUsersAPI: mockUsersAPI,
			}
			teams, err := ts.GetTeamUsers(ctx, mock.Anything, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeams, teams)
		})
	}
}

func TestTeamsAPI_UpdateRoles(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title       string
		mock        func(mockTeamAPI *mockadmin.TeamsApi)
		newRoles    []akov2.TeamRole
		expectedErr error
	}{
		{
			title:       "should not make API calls when newRole is nil",
			mock:        func(mockTeamAPI *mockadmin.TeamsApi) {},
			newRoles:    nil,
			expectedErr: nil,
		},
		{
			title: "Should successfully update team roles",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().UpdateGroupTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.UpdateGroupTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().UpdateGroupTeamExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{
						Results: &[]admin.TeamRole{
							{
								RoleNames: []string{"role1", "role2"},
								TeamId:    testTeamID1,
							},
						},
					}, &http.Response{}, nil)
			},
			newRoles:    []akov2.TeamRole{"role1", "role2"},
			expectedErr: nil,
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().UpdateGroupTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.UpdateGroupTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().UpdateGroupTeamExecute(mock.Anything).
					Return(&admin.PaginatedTeamRole{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			newRoles:    []akov2.TeamRole{},
			expectedErr: admin.GenericOpenAPIError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			err := ts.UpdateRoles(ctx, &AssignedTeam{}, mock.Anything, tt.newRoles)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTeamsAPI_AddUsers(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title       string
		mock        func(mockTeamAPI *mockadmin.TeamsApi)
		expectedErr error
	}{
		{
			title: "Should successfully add user to team",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().AddTeamUsers(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.AddTeamUsersApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().AddTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedApiAppUser{
						Results: &[]admin.CloudAppUser{
							{
								Username: "user1",
								Id:       &testUserID,
							},
						},
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().AddTeamUsers(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.AddTeamUsersApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().AddTeamUsersExecute(mock.Anything).
					Return(&admin.PaginatedApiAppUser{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr: admin.GenericOpenAPIError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			err := ts.AddUsers(ctx, &[]TeamUser{
				{
					Username: "user@name",
					UserID:   testUserID,
				},
			}, mock.Anything, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTeamsAPI_RemoveUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title       string
		mock        func(mockTeamAPI *mockadmin.TeamsApi)
		expectedErr error
	}{
		{
			title: "Should successfully remove user from team",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RemoveUserFromTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.RemoveUserFromTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RemoveUserFromTeamExecute(mock.Anything).
					Return(&http.Response{}, nil)
			},
			expectedErr: nil,
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RemoveUserFromTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.RemoveUserFromTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RemoveUserFromTeamExecute(mock.Anything).
					Return(&http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr: admin.GenericOpenAPIError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			err := ts.RemoveUser(ctx, mock.Anything, mock.Anything, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestTeamsAPI_Rename(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		title        string
		mock         func(mockTeamAPI *mockadmin.TeamsApi)
		expectedTeam *AssignedTeam
		expectedErr  error
	}{
		{
			title: "Should successfully rename team",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RenameOrgTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.RenameOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RenameOrgTeamExecute(mock.Anything).
					Return(&admin.TeamResponse{
						Id:   &testTeamID1,
						Name: &testTeamName,
					}, &http.Response{}, nil)
			},
			expectedErr: nil,
			expectedTeam: &AssignedTeam{
				TeamID:   testTeamID1,
				TeamName: testTeamName,
			},
		},
		{
			title: "Should return error when request fails",
			mock: func(mockTeamAPI *mockadmin.TeamsApi) {
				mockTeamAPI.EXPECT().RenameOrgTeam(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(admin.RenameOrgTeamApiRequest{ApiService: mockTeamAPI})
				mockTeamAPI.EXPECT().RenameOrgTeamExecute(mock.Anything).
					Return(&admin.TeamResponse{}, &http.Response{}, admin.GenericOpenAPIError{})
			},
			expectedErr:  fmt.Errorf("failed to rename team on Atlas: %w", admin.GenericOpenAPIError{}),
			expectedTeam: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			mockTeamAPI := mockadmin.NewTeamsApi(t)
			tt.mock(mockTeamAPI)
			ts := &TeamsAPI{
				teamsAPI: mockTeamAPI,
			}
			team, err := ts.RenameTeam(ctx, &AssignedTeam{}, mock.Anything, mock.Anything)
			require.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedTeam, team)
		})
	}
}
