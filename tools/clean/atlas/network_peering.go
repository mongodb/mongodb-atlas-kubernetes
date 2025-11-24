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
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

func (c *Cleaner) listNetworkPeering(ctx context.Context, projectID string, providers []string) []admin.BaseNetworkPeeringConnectionSettings {
	peers := []admin.BaseNetworkPeeringConnectionSettings{}
	for _, providerName := range providers {
		peers = append(peers, c.listNetworkPeeringForProvider(ctx, projectID, providerName)...)
	}
	if len(peers) == 0 {
		return nil
	}
	return peers
}

func (c *Cleaner) listNetworkPeeringForProvider(ctx context.Context, projectID, providerName string) []admin.BaseNetworkPeeringConnectionSettings {
	queryArgs := admin.ListPeeringConnectionsApiParams{GroupId: projectID, ProviderName: &providerName}
	peers, _, err := c.client.NetworkPeeringApi.ListPeeringConnectionsWithParams(ctx, &queryArgs).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list %s networking peering for project %s: %s", providerName, projectID, err))
		return []admin.BaseNetworkPeeringConnectionSettings{}
	}
	return *peers.Results
}

func (c *Cleaner) getNetworkPeeringContainer(ctx context.Context, projectID, ID string) *admin.CloudProviderContainer {
	container, _, err := c.client.NetworkPeeringApi.GetPeeringContainer(ctx, projectID, ID).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to get network peering container %s: %s", ID, err))

		return nil
	}

	return container
}

func (c *Cleaner) deleteNetworkPeering(ctx context.Context, projectID string, peers []admin.BaseNetworkPeeringConnectionSettings) {
	for _, peer := range peers {
		switch peer.GetProviderName() {
		case CloudProviderAWS:
			container := c.getNetworkPeeringContainer(ctx, projectID, peer.GetContainerId())
			if container == nil {
				continue
			}

			err := c.aws.DeleteVpc(ctx, peer.GetVpcId(), container.GetRegionName())
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC %s at region %s from AWS: %s", peer.GetVpcId(), container.GetRegionName(), err))

				continue
			}
		case CloudProviderGCP:
			err := c.gcp.DeleteVpc(ctx, peer.GetNetworkName())
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC %s at project %s from GCP: %s", peer.GetNetworkName(), peer.GetGcpProjectId(), err))

				continue
			}
		case CloudProviderAZURE:
			err := c.azure.DeleteVpc(ctx, peer.GetVnetName())
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC %s from Azure: %s", peer.GetVnetName(), err))

				continue
			}
		}

		_, _, err := c.client.NetworkPeeringApi.DeletePeeringConnection(ctx, projectID, peer.GetId()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of network peering %s: %s", peer.GetId(), err))

			continue
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of network peering %s", peer.GetId()))
	}
}
