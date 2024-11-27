package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

func (c *Cleaner) listFederatedDBPrivateEndpoints(ctx context.Context, projectID string) []admin.PrivateNetworkEndpointIdEntry {
	federatedDBPEs, _, err := c.client.DataFederationApi.ListDataFederationPrivateEndpoints(ctx, projectID).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list federated databases for project %s: %s", projectID, err))

		return nil
	}

	return federatedDBPEs.GetResults()
}

func (c *Cleaner) deleteFederatedDBPrivateEndpoints(ctx context.Context, projectID string, dbpes []admin.PrivateNetworkEndpointIdEntry) {
	for _, fedDBPE := range dbpes {
		_, _, err := c.client.DataFederationApi.DeleteDataFederationPrivateEndpoint(ctx, projectID, fedDBPE.GetEndpointId()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of Federated DB private endpoint %s: %s", fedDBPE.GetEndpointId(), err))
		}
		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of Federated DB private endpoint %s", fedDBPE.GetEndpointId()))
	}
}
