package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
)

func (c *Cleaner) listFederatedDatabases(ctx context.Context, projectID string) []admin.DataLakeTenant {
	federatedDBs, _, err := c.client.DataFederationApi.
		ListFederatedDatabases(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list federated databases for project %s: %s", projectID, err))

		return nil
	}

	return federatedDBs
}

func (c *Cleaner) deleteFederatedDatabases(ctx context.Context, projectID string, dbs []admin.DataLakeTenant) {
	for _, fedDB := range dbs {
		if fedDB.GetState() == "DELETED" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tFederated Database %s is being deleted...", fedDB.GetName()))

			continue
		}

		_, _, err := c.client.DataFederationApi.DeleteFederatedDatabase(ctx, projectID, fedDB.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of Federated Database %s: %s", fedDB.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of Federated Database %s", fedDB.GetName()))
	}
}
