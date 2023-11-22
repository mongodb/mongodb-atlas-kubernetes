package atlas

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
)

func (c *Cleaner) listClusters(ctx context.Context, projectID string) []admin.AdvancedClusterDescription {
	clusters, _, err := c.client.ClustersApi.
		ListClusters(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list clusters for project %s: %s", projectID, err))

		return nil
	}

	return clusters.Results
}

func (c *Cleaner) deleteClusters(ctx context.Context, projectID string, clusters []admin.AdvancedClusterDescription) {
	for _, cluster := range clusters {
		if cluster.GetStateName() == "DELETING" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tCluster %s is being deleted...", cluster.GetName()))

			return
		}

		_, err := c.client.ClustersApi.DeleteCluster(ctx, projectID, cluster.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of cluster %s: %s", cluster.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of cluster %s", cluster.GetName()))
	}
}

func (c *Cleaner) listServerlessClusters(ctx context.Context, projectID string) []admin.ServerlessInstanceDescription {
	clusters, _, err := c.client.ServerlessInstancesApi.
		ListServerlessInstances(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list serverless clusters for project %s: %s", projectID, err))

		return nil
	}

	return clusters.Results
}

func (c *Cleaner) deleteServerlessClusters(ctx context.Context, projectID string, clusters []admin.ServerlessInstanceDescription) {
	for _, cluster := range clusters {
		if cluster.GetStateName() == "DELETING" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tServerless Cluster %s is being deleted...", cluster.GetName()))

			return
		}

		_, _, err := c.client.ServerlessInstancesApi.DeleteServerlessInstance(ctx, projectID, cluster.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of serverless cluster %s: %s", cluster.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of serverless cluster %s", cluster.GetName()))
	}
}
