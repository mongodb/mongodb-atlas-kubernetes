package teams

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type Team struct {
	Roles    []string
	TeamName string
	TeamID   string
}

type AssignedTeam struct {
	Usernames []string
	TeamName  string
	TeamID    string
}

type TeamUser struct {
	Username string
	UserID   string
}

func NewTeam(projTeamSpec *akov2.Team, teamID string) *Team {
	if projTeamSpec == nil {
		return nil
	}

	roles := make([]string, 0)
	for _, role := range projTeamSpec.Roles {
		roles = append(roles, string(role))
	}

	team := &Team{
		Roles:  roles,
		TeamID: teamID,
	}

	return team
}

func NewAssignedTeam(teamSpec *akov2.TeamSpec, teamID string) *AssignedTeam {
	if teamSpec == nil {
		return nil
	}
	usernames := make([]string, 0)
	for _, username := range teamSpec.Usernames {
		usernames = append(usernames, string(username))
	}

	team := &AssignedTeam{
		TeamID:    teamID,
		TeamName:  teamSpec.Name,
		Usernames: usernames,
	}

	return team
}

func TeamFromAtlas(team *admin.TeamResponse) *Team {
	if team == nil {
		return nil
	}
	tm := &Team{
		TeamID:   team.GetId(),
		TeamName: team.GetName(),
	}
	return tm
}

func TeamToAtlas(team *akov2.Team, teamID string) *Team {
	roleNames := make([]string, 0)
	for _, role := range team.Roles {
		roleNames = append(roleNames, string(role))
	}
	return &Team{
		TeamID: teamID,
		Roles:  roleNames,
	}
}

func TeamRoleFromAtlas(atlasTeams []admin.TeamRole) []Team {
	teams := make([]Team, 0)
	for _, team := range atlasTeams {
		teams = append(teams, Team{Roles: team.GetRoleNames(), TeamID: team.GetTeamId()})
	}
	return teams
}

func TeamRoleToAtlas(atlasTeams []Team) []admin.TeamRole {
	teams := make([]admin.TeamRole, 0)

	for _, team := range atlasTeams {
		result := admin.TeamRole{
			TeamId:    &team.TeamID,
			RoleNames: &team.Roles,
		}
		teams = append(teams, result)
	}
	return teams
}

func AssignedTeamFromAtlas(assignedTeam *admin.TeamResponse) *AssignedTeam {
	return &AssignedTeam{
		TeamID:   assignedTeam.GetId(),
		TeamName: assignedTeam.GetName(),
	}
}

func AssignedTeamToAtlas(assignedTeam *AssignedTeam) *admin.Team {
	return &admin.Team{
		Id:        &assignedTeam.TeamID,
		Name:      assignedTeam.TeamName,
		Usernames: &assignedTeam.Usernames,
	}
}

func UsersFromAtlas(users *admin.PaginatedApiAppUser) []TeamUser {
	teamUsers := make([]TeamUser, 0)
	for _, user := range users.GetResults() {
		teamUsers = append(teamUsers, TeamUser{Username: user.Username, UserID: user.GetId()})
	}
	return teamUsers
}
