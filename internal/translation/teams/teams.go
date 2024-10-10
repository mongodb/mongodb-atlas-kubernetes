package teams

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
	atlasAssignedTeams, _, err := tm.teamsAPI.ListProjectTeams(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get project team list from Atlas: %w", err)
	}
	return TeamRolesFromAtlas(atlasAssignedTeams.GetResults()), err
}

func (tm *TeamsAPI) GetTeamByName(ctx context.Context, orgID, teamName string) (*AssignedTeam, error) {
	atlasTeam, resp, err := tm.teamsAPI.GetTeamByName(ctx, orgID, teamName).Execute()
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by name from Atlas: %w", err)
	}
	return AssignedTeamFromAtlas(atlasTeam), err
}

func (tm *TeamsAPI) GetTeamByID(ctx context.Context, orgID, teamID string) (*AssignedTeam, error) {
	atlasTeam, _, err := tm.teamsAPI.GetTeamById(ctx, orgID, teamID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get team by ID from Atlas: %w", err)
	}
	return AssignedTeamFromAtlas(atlasTeam), err
}

func (tm *TeamsAPI) Assign(ctx context.Context, at *[]AssignedTeam, projectID string) error {
	desiredRoles := TeamRolesToAtlas(*at)
	_, _, err := tm.teamsAPI.AddAllTeamsToProject(ctx, projectID, &desiredRoles).Execute()
	return err
}

func (tm *TeamsAPI) Unassign(ctx context.Context, projectID, teamID string) error {
	_, err := tm.teamsAPI.RemoveProjectTeam(ctx, projectID, teamID).Execute()
	return err
}

func (tm *TeamsAPI) Create(ctx context.Context, at *Team, orgID string) (*Team, error) {
	desiredTeam := TeamToAtlas(at)
	atlasTeam, _, err := tm.teamsAPI.CreateTeam(ctx, orgID, desiredTeam).Execute()
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
	atlasTeam, _, err := tm.teamsAPI.RenameTeam(ctx, orgID, at.TeamID, teamUpdate).Execute()
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

	_, _, err := tm.teamsAPI.UpdateTeamRoles(ctx, projectID, at.TeamID, &admin.TeamRole{RoleNames: &roles}).Execute()
	return err
}

func (tm *TeamsAPI) GetTeamUsers(ctx context.Context, orgID, teamID string) ([]TeamUser, error) {
	atlasUsers, _, err := tm.teamsAPI.ListTeamUsers(ctx, orgID, teamID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get team users from Atlas: %w", err)
	}

	return UsersFromAtlas(atlasUsers), err
}

func (tm *TeamsAPI) AddUsers(ctx context.Context, usersToAdd *[]TeamUser, orgID, teamID string) error {
	_, _, err := tm.teamsAPI.AddTeamUser(ctx, orgID, teamID, UsersToAtlas(usersToAdd)).Execute()
	return err
}

func (tm *TeamsAPI) RemoveUser(ctx context.Context, orgID, teamID, userID string) error {
	_, err := tm.teamsAPI.RemoveTeamUser(ctx, orgID, teamID, userID).Execute()
	return err
}
