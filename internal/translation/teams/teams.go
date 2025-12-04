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

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

type TeamsService interface {
	TeamProjectsService
	TeamRolesService
	TeamUsersService
}

type TeamProjectsService interface { // manages Team's associations to Projects
	ListProjectTeams(ctx context.Context, projectID string) ([]AssignedTeam, error)
	Create(ctx context.Context, at *Team, orgID string) (*Team, error)
	Assign(ctx context.Context, at *[]AssignedTeam, projectID string) error
	Unassign(ctx context.Context, projectID, teamID string) error
}

type TeamRolesService interface { // manages Team's Roles
	GetTeamByName(ctx context.Context, orgID, teamName string) (*AssignedTeam, error)
	GetTeamByID(ctx context.Context, orgID, teamID string) (*AssignedTeam, error)
	RenameTeam(ctx context.Context, at *AssignedTeam, orgID, newName string) (*AssignedTeam, error)
	UpdateRoles(ctx context.Context, at *AssignedTeam, projectID string, newRoles []akov2.TeamRole) error
}

type TeamUsersService interface { // manages Team's Members (Users)
	GetTeamUsers(ctx context.Context, orgID, teamID string) ([]TeamUser, error)
	AddUsers(ctx context.Context, usersToAdd *[]TeamUser, orgID, teamID string) error
	RemoveUser(ctx context.Context, orgID, teamID, userID string) error
}

type TeamsAPI struct {
	teamsAPI     admin.TeamsApi
	teamUsersAPI admin.MongoDBCloudUsersApi
}

func NewTeamsAPIService(teamAPI admin.TeamsApi, userAPI admin.MongoDBCloudUsersApi) *TeamsAPI {
	return &TeamsAPI{
		teamsAPI:     teamAPI,
		teamUsersAPI: userAPI,
	}
}

func (tm *TeamsAPI) ListProjectTeams(ctx context.Context, projectID string) ([]AssignedTeam, error) {
	atlasAssignedTeams, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.TeamRole], *http.Response, error) {
		return tm.teamsAPI.ListGroupTeams(ctx, projectID).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get project team list from Atlas: %w", err)
	}
	return TeamRolesFromAtlas(atlasAssignedTeams), err
}

func (tm *TeamsAPI) GetTeamByName(ctx context.Context, orgID, teamName string) (*AssignedTeam, error) {
	atlasTeam, resp, err := tm.teamsAPI.GetTeamByName(ctx, orgID, teamName).Execute()
	if err != nil {
		if httputil.StatusCode(resp) == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by name from Atlas: %w", err)
	}
	return AssignedTeamFromAtlas(atlasTeam), err
}

func (tm *TeamsAPI) GetTeamByID(ctx context.Context, orgID, teamID string) (*AssignedTeam, error) {
	atlasTeam, _, err := tm.teamsAPI.GetOrgTeam(ctx, orgID, teamID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get team by ID from Atlas: %w", err)
	}
	return AssignedTeamFromAtlas(atlasTeam), err
}

func (tm *TeamsAPI) Assign(ctx context.Context, at *[]AssignedTeam, projectID string) error {
	desiredRoles := TeamRolesToAtlas(*at)
	_, _, err := tm.teamsAPI.AddGroupTeams(ctx, projectID, &desiredRoles).Execute()
	return err
}

func (tm *TeamsAPI) Unassign(ctx context.Context, projectID, teamID string) error {
	_, err := tm.teamsAPI.RemoveGroupTeam(ctx, projectID, teamID).Execute()
	return err
}

func (tm *TeamsAPI) Create(ctx context.Context, at *Team, orgID string) (*Team, error) {
	desiredTeam := TeamToAtlas(at)
	atlasTeam, _, err := tm.teamsAPI.CreateOrgTeam(ctx, orgID, desiredTeam).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create team on Atlas: %w", err)
	}

	teamResponse := &admin.TeamResponse{}
	teamResponse.SetId(atlasTeam.GetId())
	teamResponse.SetName(atlasTeam.GetName())
	return TeamFromAtlas(teamResponse), err
}

func (tm *TeamsAPI) RenameTeam(ctx context.Context, at *AssignedTeam, orgID, newName string) (*AssignedTeam, error) {
	teamUpdate := &admin.TeamUpdate{Name: newName}
	atlasTeam, _, err := tm.teamsAPI.RenameOrgTeam(ctx, orgID, at.TeamID, teamUpdate).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to rename team on Atlas: %w", err)
	}
	return AssignedTeamFromAtlas(atlasTeam), err
}

func (tm *TeamsAPI) UpdateRoles(ctx context.Context, at *AssignedTeam, projectID string, newRoles []akov2.TeamRole) error {
	if newRoles == nil {
		return nil
	}
	roles := make([]string, 0, len(newRoles))
	for _, role := range newRoles {
		roles = append(roles, string(role))
	}

	_, _, err := tm.teamsAPI.UpdateGroupTeam(ctx, projectID, at.TeamID, &admin.TeamRole{RoleNames: &roles}).Execute()
	return err
}

func (tm *TeamsAPI) GetTeamUsers(ctx context.Context, orgID, teamID string) ([]TeamUser, error) {
	atlasUsers, _, err := tm.teamUsersAPI.ListTeamUsers(ctx, orgID, teamID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get team users from Atlas: %w", err)
	}

	return UsersFromAtlas(atlasUsers), err
}

func (tm *TeamsAPI) AddUsers(ctx context.Context, usersToAdd *[]TeamUser, orgID, teamID string) error {
	_, _, err := tm.teamsAPI.AddTeamUsers(ctx, orgID, teamID, UsersToAtlas(usersToAdd)).Execute()
	return err
}

func (tm *TeamsAPI) RemoveUser(ctx context.Context, orgID, teamID, userID string) error {
	_, err := tm.teamsAPI.RemoveUserFromTeam(ctx, orgID, teamID, userID).Execute()
	return err
}
