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

package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
)

func (c *Cleaner) listTeamsByOrg(ctx context.Context, orgID string) []admin.TeamResponse {
	teamsList, _, err := c.client.TeamsApi.
		ListOrgTeams(ctx, orgID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("Failed to list teams of organization %s: %s", orgID, err))

		return nil
	}

	if teamsList.GetTotalCount() == 0 {
		fmt.Println(text.FgYellow.Sprintf("No teams found in organization %s", orgID))

		return nil
	}

	return *teamsList.Results
}

func (c *Cleaner) deleteTeam(ctx context.Context, orgID string, team *admin.TeamResponse) {
	_, err := c.client.TeamsApi.DeleteOrgTeam(ctx, orgID, team.GetId()).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to request deletion of team %s(%s): %s", team.GetName(), team.GetId(), err))

		return
	}

	fmt.Println(text.FgGreen.Sprintf("\tRequested deletion of team %s(%s)", team.GetName(), team.GetId()))
}
