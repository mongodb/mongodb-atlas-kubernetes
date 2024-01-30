package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func (c *Cleaner) listTeamsByOrg(ctx context.Context, orgID string) []admin.TeamResponse {
	teamsList, _, err := c.client.TeamsApi.
		ListOrganizationTeams(ctx, orgID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("Failed to list teams of organization %s: %s", orgID, err))

		return nil
	}

	if teamsList.GetTotalCount() == 0 {
		fmt.Println(text.FgYellow.Sprintf("No teams found in organization %s", orgID))

		return nil
	}

	return teamsList.Results
}

func (c *Cleaner) deleteTeam(ctx context.Context, orgID string, team *admin.TeamResponse) {
	_, _, err := c.client.TeamsApi.DeleteTeam(ctx, orgID, team.GetId()).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to request deletion of team %s(%s): %s", team.GetName(), team.GetId(), err))

		return
	}

	fmt.Println(text.FgGreen.Sprintf("\tRequested deletion of team %s(%s)", team.GetName(), team.GetId()))
}
