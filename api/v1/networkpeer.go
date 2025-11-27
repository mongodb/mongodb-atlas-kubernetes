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

package v1

import (
	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// NetworkPeer configured for the current Project.
// Deprecated: Migrate to the AtlasNetworkPeering and AtlasNetworkContainer custom resources in accordance with
// the migration guide at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
// +optional
type NetworkPeer struct {
	// AccepterRegionName is the provider region name of user's VPC.
	// +optional
	AccepterRegionName string `json:"accepterRegionName"`
	// ContainerRegion is the provider region name of Atlas network peer container. If not set, AccepterRegionName is used.
	// +optional
	ContainerRegion string `json:"containerRegion"`
	// AccountID of the user's VPC.
	// +optional
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// ID of the network peer container. If not set, operator will create a new container with ContainerRegion and AtlasCIDRBlock input.
	// +optional
	ContainerID string `json:"containerId"`
	// ProviderName is the name of the provider. If not set, it will be set to "AWS".
	// +optional
	ProviderName provider.ProviderName `json:"providerName,omitempty"`
	// User VPC CIDR.
	// +optional
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	// AWS VPC ID.
	// +optional
	VpcID string `json:"vpcId,omitempty"`
	// Atlas CIDR. It needs to be set if ContainerID is not set.
	// +optional
	AtlasCIDRBlock string `json:"atlasCidrBlock"`
	// AzureDirectoryID is the unique identifier for an Azure AD directory.
	// +optional
	AzureDirectoryID string `json:"azureDirectoryId,omitempty"`
	// AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.
	// +optional
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	// ResourceGroupName is the name of your Azure resource group.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// VNetName is name of your Azure VNet. Its applicable only for Azure.
	// +optional
	VNetName string `json:"vnetName,omitempty"`
	// User GCP Project ID. Its applicable only for GCP.
	// +optional
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// GCP Network Peer Name. Its applicable only for GCP.
	// +optional
	NetworkName string `json:"networkName,omitempty"`
}

// NewNetworkPeerFromAtlas creates a network peer based off a network peering connection from Atlas.
// Note: ContainerRegion and AtlasCIDRBlock are unset
// as this information is not provided by Atlas for a peering connection.
func NewNetworkPeerFromAtlas(atlasPeer admin.BaseNetworkPeeringConnectionSettings) *NetworkPeer {
	return &NetworkPeer{
		AccepterRegionName:  atlasPeer.GetAccepterRegionName(),
		AWSAccountID:        atlasPeer.GetAwsAccountId(),
		ContainerID:         atlasPeer.GetContainerId(),
		ProviderName:        provider.ProviderName(atlasPeer.GetProviderName()),
		RouteTableCIDRBlock: atlasPeer.GetRouteTableCidrBlock(),
		VpcID:               atlasPeer.GetVpcId(),
		AzureDirectoryID:    atlasPeer.GetAzureDirectoryId(),
		AzureSubscriptionID: atlasPeer.GetAzureSubscriptionId(),
		ResourceGroupName:   atlasPeer.GetResourceGroupName(),
		VNetName:            atlasPeer.GetVnetName(),
		GCPProjectID:        atlasPeer.GetGcpProjectId(),
		NetworkName:         atlasPeer.GetNetworkName(),
	}
}
func (in *NetworkPeer) ToAtlas() (*admin.BaseNetworkPeeringConnectionSettings, error) {
	result := &admin.BaseNetworkPeeringConnectionSettings{}
	err := compat.JSONCopy(result, in)
	return result, err
}

func (in *NetworkPeer) ToAtlasPeer() *admin.BaseNetworkPeeringConnectionSettings {
	switch in.ProviderName {
	case provider.ProviderAWS:
		return &admin.BaseNetworkPeeringConnectionSettings{
			AccepterRegionName:  pointer.SetOrNil(in.AccepterRegionName, ""),
			AwsAccountId:        pointer.SetOrNil(in.AWSAccountID, ""),
			ContainerId:         in.ContainerID,
			ProviderName:        pointer.SetOrNil(string(in.ProviderName), ""),
			RouteTableCidrBlock: pointer.SetOrNil(in.RouteTableCIDRBlock, ""),
			VpcId:               pointer.SetOrNil(in.VpcID, ""),
		}
	case provider.ProviderGCP:
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:  in.ContainerID,
			ProviderName: pointer.SetOrNil(string(in.ProviderName), ""),
			GcpProjectId: pointer.SetOrNil(in.GCPProjectID, ""),
			NetworkName:  pointer.SetOrNil(in.NetworkName, ""),
		}
	case provider.ProviderAzure:
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         in.ContainerID,
			ProviderName:        pointer.SetOrNil(string(in.ProviderName), ""),
			AzureDirectoryId:    pointer.SetOrNil(in.AzureDirectoryID, ""),
			AzureSubscriptionId: pointer.SetOrNil(in.AzureSubscriptionID, ""),
			ResourceGroupName:   pointer.SetOrNil(in.ResourceGroupName, ""),
			VnetName:            pointer.SetOrNil(in.VNetName, ""),
		}
	}

	return nil
}

func (in *NetworkPeer) GetContainerRegion() string {
	if in.ContainerRegion != "" {
		return in.ContainerRegion
	}
	return in.AccepterRegionName
}
