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
	"go.mongodb.org/atlas-sdk/v20250312022/admin"
)

func (c *Cleaner) listClusters(ctx context.Context, projectID string) []admin.ClusterDescription20240805 {
	clusters, _, err := c.client.ClustersAPI.
		ListClusters(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list clusters for project %s: %s", projectID, err))

		return nil
	}

	return clusters.Results
}

func (c *Cleaner) deleteClusters(ctx context.Context, projectID string, clusters []admin.ClusterDescription20240805) {
	for _, cluster := range clusters {
		if cluster.GetStateName() == "DELETING" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tCluster %s is being deleted...", cluster.GetName()))

			continue
		}

		if cluster.GetTerminationProtectionEnabled() {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tDisabling termination protection for Cluster %s...", cluster.GetName()))
			params := new(admin.ClusterDescription20240805)
			params.SetTerminationProtectionEnabled(false)
			_, _, err := c.client.ClustersAPI.UpdateCluster(ctx, projectID, cluster.GetName(), params).Execute()
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to update cluster %s: %s", cluster.GetName(), err))
			}
		}

		_, err := c.client.ClustersAPI.DeleteCluster(ctx, projectID, cluster.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of cluster %s: %s", cluster.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of cluster %s", cluster.GetName()))
	}
}

func (c *Cleaner) listServerlessClusters(ctx context.Context, projectID string) []admin.ServerlessInstanceDescription {
	clusters, _, err := c.client.ServerlessInstancesAPI.
		ListServerlessInstances(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list serverless clusters for project %s: %s", projectID, err))

		return nil
	}

	return clusters.Results
}

func (c *Cleaner) deleteServerlessClusters(_ context.Context, _ string, clusters []admin.ServerlessInstanceDescription) {
	for _, cluster := range clusters {
		fmt.Println(text.FgYellow.Sprintf("\t\t\tServerless Cluster %s skipped (DeleteServerlessInstance API removed in SDK v20250312022)", cluster.GetName()))
	}
}

func (c *Cleaner) listFlexClusters(ctx context.Context, projectID string) []admin.FlexClusterDescription20241113 {
	clusters, _, err := c.client.FlexClustersAPI.
		ListFlexClusters(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list flex clusters for project %s: %s", projectID, err))

		return nil
	}

	return clusters.Results
}

func (c *Cleaner) deleteFlexClusters(ctx context.Context, projectID string, clusters []admin.FlexClusterDescription20241113) {
	for _, cluster := range clusters {
		if cluster.GetStateName() == "DELETING" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tServerless Cluster %s is being deleted...", cluster.GetName()))

			continue
		}

		if cluster.GetTerminationProtectionEnabled() {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tDisabling termination protection for flex Cluster %s...", cluster.GetName()))
			params := new(admin.FlexClusterDescriptionUpdate20241113)
			params.SetTerminationProtectionEnabled(false)
			_, _, err := c.client.FlexClustersAPI.UpdateFlexCluster(ctx, projectID, cluster.GetName(), params).Execute()
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to update flex cluster %s: %s", cluster.GetName(), err))
			}
		}

		_, err := c.client.FlexClustersAPI.DeleteFlexCluster(ctx, projectID, cluster.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of flex cluster %s: %s", cluster.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of flex cluster %s", cluster.GetName()))
	}
}
