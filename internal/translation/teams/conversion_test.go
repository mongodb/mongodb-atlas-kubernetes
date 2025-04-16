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
	"testing"

	"github.com/stretchr/testify/assert"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	testTeamID = "team-id"
)

func TestNewTeam(t *testing.T) {
	for _, tc := range []struct {
		title        string
		teamSpec     *akov2.TeamSpec
		teamID       string
		expectedTeam *Team
	}{
		{
			title: "Nil spec returns nil user",
		},
		{
			title:        "Empty spec returns Empty user",
			teamSpec:     &akov2.TeamSpec{},
			expectedTeam: &Team{Usernames: []string{}},
		},
		{
			title: "Populated spec is properly created",
			teamSpec: &akov2.TeamSpec{
				Name:      testTeamName,
				Usernames: []akov2.TeamUser{"user1", "user2"},
			},
			teamID: testTeamID,
			expectedTeam: &Team{
				TeamName:  testTeamName,
				TeamID:    testTeamID,
				Usernames: []string{"user1", "user2"},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			team := NewTeam(tc.teamSpec, tc.teamID)
			assert.Equal(t, tc.expectedTeam, team)
		})
	}
}

func TestNewAssignedTeam(t *testing.T) {
	for _, tc := range []struct {
		title        string
		projTeamSpec *akov2.Team
		teamID       string
		expectedTeam *AssignedTeam
	}{
		{
			title: "Nil spec returns nil user",
		},
		{
			title:        "Empty spec returns Empty user",
			projTeamSpec: &akov2.Team{},
			expectedTeam: &AssignedTeam{Roles: []string{}},
		},
		{
			title: "Populated spec is properly created",
			projTeamSpec: &akov2.Team{
				Roles: []akov2.TeamRole{"role1", "role2"},
			},
			teamID: testTeamID,
			expectedTeam: &AssignedTeam{
				Roles:  []string{"role1", "role2"},
				TeamID: testTeamID,
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			team := NewAssignedTeam(tc.projTeamSpec, tc.teamID)
			assert.Equal(t, tc.expectedTeam, team)
		})
	}
}
