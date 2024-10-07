package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	testTeamID = "team-id"
)

func TestNewTeam(t *testing.T) {
	for _, tc := range []struct {
		title        string
		projTeamSpec *akov2.Team
		teamID       string
		expectedTeam *Team
	}{
		{
			title: "Nil spec returns nil user",
		},
		{
			title:        "Empty spec returns Empty user",
			projTeamSpec: &akov2.Team{},
			expectedTeam: &Team{},
		},
		{
			title: "Populated spec is properly created",
			projTeamSpec: &akov2.Team{
				Roles: []akov2.TeamRole{"role1", "role2"},
			},
			teamID: testTeamID,
			expectedTeam: &Team{
				Roles:  []string{"role1", "role2"},
				TeamID: testTeamID,
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			team := NewTeam(tc.projTeamSpec, tc.teamID)
			assert.Equal(t, tc.expectedTeam, team)
		})
	}
}

func TestNewAssignedTeam(t *testing.T) {
	for _, tc := range []struct {
		title        string
		teamSpec     *akov2.TeamSpec
		teamID       string
		expectedTeam *AssignedTeam
	}{
		{
			title: "Nil spec returns nil user",
		},
		{
			title:        "Empty spec returns Empty user",
			teamSpec:     &akov2.TeamSpec{},
			expectedTeam: &AssignedTeam{},
		},
		{
			title: "Populated spec is properly created",
			teamSpec: &akov2.TeamSpec{
				Name:      testTeamName,
				Usernames: []akov2.TeamUser{"user1", "user2"},
			},
			teamID: testTeamID,
			expectedTeam: &AssignedTeam{
				TeamName:  testTeamName,
				TeamID:    testTeamID,
				Usernames: []string{"user1", "user2"},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			team := NewAssignedTeam(tc.teamSpec, tc.teamID)
			assert.Equal(t, tc.expectedTeam, team)
		})
	}
}
