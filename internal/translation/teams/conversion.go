package teams

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type Team struct {
	Usernames []string
	TeamID    string
	TeamName  string
}

type AssignedTeam struct {
	Roles    []string
	TeamID   string
	TeamName string
}

type TeamUser struct {
	Username string
	UserID   string
}

func NewTeam(teamSpec *akov2.TeamSpec, teamID string) *Team {
	if teamSpec == nil {
		return nil
	}
	usernames := make([]string, 0, len(teamSpec.Usernames))
	for _, username := range teamSpec.Usernames {
		usernames = append(usernames, string(username))
	}

	team := &Team{
		TeamID:    teamID,
		TeamName:  teamSpec.Name,
		Usernames: usernames,
	}

	return team
}

func NewAssignedTeam(projTeamSpec *akov2.Team, teamID string) *AssignedTeam {
	if projTeamSpec == nil {
		return nil
	}

	roles := make([]string, 0, len(projTeamSpec.Roles))
	for _, role := range projTeamSpec.Roles {
		roles = append(roles, string(role))
	}

	team := &AssignedTeam{
		Roles:  roles,
		TeamID: teamID,
	}

	return team
}

func TeamFromAtlas(assignedTeam *admin.TeamResponse) *Team {
	return &Team{
		TeamID:   assignedTeam.GetId(),
		TeamName: assignedTeam.GetName(),
	}
}

func TeamToAtlas(team *Team) *admin.Team {
	return &admin.Team{
		Id:        pointer.MakePtrOrNil(team.TeamID),
		Name:      team.TeamName,
		Usernames: &team.Usernames,
	}
}

func AssignedTeamFromAtlas(team *admin.TeamResponse) *AssignedTeam {
	if team == nil {
		return nil
	}

	tm := &AssignedTeam{
		TeamID:   team.GetId(),
		TeamName: team.GetName(),
	}
	return tm
}

func TeamRolesFromAtlas(atlasTeams []admin.TeamRole) []AssignedTeam {
	teams := make([]AssignedTeam, 0, len(atlasTeams))
	for _, team := range atlasTeams {
		teams = append(teams, AssignedTeam{Roles: team.GetRoleNames(), TeamID: team.GetTeamId()})
	}
	return teams
}

func TeamRolesToAtlas(atlasTeams []AssignedTeam) []admin.TeamRole {
	if atlasTeams == nil {
		return nil
	}
	teams := make([]admin.TeamRole, 0, len(atlasTeams))

	for _, team := range atlasTeams {
		result := admin.TeamRole{
			TeamId:    pointer.MakePtrOrNil(team.TeamID),
			RoleNames: &team.Roles,
		}
		teams = append(teams, result)
	}
	return teams
}

func UsersFromAtlas(users *admin.PaginatedApiAppUser) []TeamUser {
	teamUsers := make([]TeamUser, 0, len(users.GetResults()))
	for _, user := range users.GetResults() {
		teamUsers = append(teamUsers, TeamUser{
			Username: user.Username,
			UserID:   user.GetId(),
		})
	}
	return teamUsers
}

func UsersToAtlas(teamUsers *[]TeamUser) *[]admin.AddUserToTeam {
	users := *teamUsers
	desiredUsers := make([]admin.AddUserToTeam, 0, len(users))
	for _, user := range users {
		desiredUsers = append(desiredUsers, admin.AddUserToTeam{
			Id: user.UserID,
		})
	}
	return &desiredUsers
}
