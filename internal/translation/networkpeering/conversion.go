package networkpeering

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

type NetworkPeer struct {
	akov2.AtlasNetworkPeeringConfig
	ID string
}

func NewNetworkPeer(id string, cfg *akov2.AtlasNetworkPeeringConfig) *NetworkPeer {
	return &NetworkPeer{
		AtlasNetworkPeeringConfig: *cfg,
		ID:                        id,
	}
}

type ProviderContainer struct {
	akov2.AtlasProviderContainerConfig
	ID       string
	Provider string
}

func NewProviderContainer(id string, provider string, cfg *akov2.AtlasProviderContainerConfig) *ProviderContainer {
	return &ProviderContainer{
		AtlasProviderContainerConfig: *cfg,
		ID:                           id,
		Provider:                     provider,
	}
}

func toAtlasConnection(peer *NetworkPeer) (*admin.BaseNetworkPeeringConnectionSettings, error) {
	switch peer.Provider {
	case string(provider.ProviderAWS):
		if peer.AWSConfiguration == nil {
			return nil, fmt.Errorf("unsupported AWS peer with AWSConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(peer.Provider, ""),
			AccepterRegionName:  pointer.SetOrNil(peer.AWSConfiguration.AccepterRegionName, ""),
			AwsAccountId:        pointer.SetOrNil(peer.AWSConfiguration.AWSAccountID, ""),
			RouteTableCidrBlock: pointer.SetOrNil(peer.AWSConfiguration.RouteTableCIDRBlock, ""),
			VpcId:               pointer.SetOrNil(peer.AWSConfiguration.VpcID, ""),
		}, nil
	case string(provider.ProviderGCP):
		if peer.GCPConfiguration == nil {
			return nil, fmt.Errorf("unsupported Google peer with GCPConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:  peer.ContainerID,
			ProviderName: pointer.SetOrNil(peer.Provider, ""),
			GcpProjectId: pointer.SetOrNil(peer.GCPConfiguration.GCPProjectID, ""),
			NetworkName:  pointer.SetOrNil(peer.GCPConfiguration.NetworkName, ""),
		}, nil
	case string(provider.ProviderAzure):
		if peer.AzureConfiguration == nil {
			return nil, fmt.Errorf("unsupported Azure peer with AzureConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(peer.Provider, ""),
			AzureDirectoryId:    pointer.SetOrNil(peer.AzureConfiguration.AzureDirectoryID, ""),
			AzureSubscriptionId: pointer.SetOrNil(peer.AzureConfiguration.AzureSubscriptionID, ""),
			ResourceGroupName:   pointer.SetOrNil(peer.AzureConfiguration.ResourceGroupName, ""),
			VnetName:            pointer.SetOrNil(peer.AzureConfiguration.VNetName, ""),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", peer.Provider)
	}
}

func fromAtlasConnection(conn *admin.BaseNetworkPeeringConnectionSettings) (*NetworkPeer, error) {
	switch provider.ProviderName(conn.GetProviderName()) {
	case provider.ProviderAWS:
		return &NetworkPeer{
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  conn.GetAccepterRegionName(),
					AWSAccountID:        conn.GetAwsAccountId(),
					RouteTableCIDRBlock: conn.GetRouteTableCidrBlock(),
					VpcID:               conn.GetVpcId(),
				},
			},
		}, nil
	case provider.ProviderGCP:
		return &NetworkPeer{
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: conn.GetGcpProjectId(),
					NetworkName:  conn.GetNetworkName(),
				},
			},
		}, nil
	case provider.ProviderAzure:
		return &NetworkPeer{
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
				AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
					AzureDirectoryID:    conn.GetAzureDirectoryId(),
					AzureSubscriptionID: conn.GetAzureSubscriptionId(),
					ResourceGroupName:   conn.GetResourceGroupName(),
					VNetName:            conn.GetVnetName(),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", conn.GetProviderName())
	}
}

func fromAtlasConnectionList(list []admin.BaseNetworkPeeringConnectionSettings) ([]NetworkPeer, error) {
	if list == nil {
		return nil, nil
	}
	peers := make([]NetworkPeer, 0, len(list))
	for i, conn := range list {
		c, err := fromAtlasConnection(&conn)
		if err != nil {
			return nil, fmt.Errorf("failed to convert connection list item %d: %w", i, err)
		}
		peers = append(peers, *c)
	}
	return peers, nil
}

func toAtlasContainer(container *ProviderContainer) *admin.CloudProviderContainer {
	cpc := &admin.CloudProviderContainer{
		Id:             pointer.SetOrNil(container.ID, ""),
		ProviderName:   pointer.SetOrNil(container.Provider, ""),
		AtlasCidrBlock: pointer.SetOrNil(container.AtlasCIDRBlock, ""),
	}
	if cpc.GetProviderName() == string(provider.ProviderAWS) {
		cpc.RegionName = pointer.SetOrNil(container.ContainerRegion, "")
	} else {
		cpc.Region = pointer.SetOrNil(container.ContainerRegion, "")
	}
	return cpc
}

func fromAtlasContainer(container *admin.CloudProviderContainer) *ProviderContainer {
	region := container.GetRegion()
	if container.GetProviderName() == string(provider.ProviderAWS) {
		region = container.GetRegionName()
	}
	return &ProviderContainer{
		ID:       container.GetId(),
		Provider: container.GetProviderName(),
		AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
			AtlasCIDRBlock:  container.GetAtlasCidrBlock(),
			ContainerRegion: region,
		},
	}
}

func fromAtlasContainerList(list []admin.CloudProviderContainer) []ProviderContainer {
	if list == nil {
		return nil
	}
	containers := make([]ProviderContainer, 0, len(list))
	for _, container := range list {
		containers = append(containers, *fromAtlasContainer(&container))
	}
	return containers
}
