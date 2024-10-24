package networkpeering

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

type ProviderContainer struct {
	ID             string
	ProviderName   provider.ProviderName
	AtlasCIDRBlock string
	RegionName     string // AWS
	Region         string // Azure
}

func toAtlasConnection(peer *akov2.NetworkPeer) *admin.BaseNetworkPeeringConnectionSettings {
	switch peer.ProviderName {
	case provider.ProviderAWS:
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(string(peer.ProviderName), ""),
			AccepterRegionName:  pointer.SetOrNil(peer.AccepterRegionName, ""),
			AwsAccountId:        pointer.SetOrNil(peer.AWSAccountID, ""),
			RouteTableCidrBlock: pointer.SetOrNil(peer.RouteTableCIDRBlock, ""),
			VpcId:               pointer.SetOrNil(peer.VpcID, ""),
		}
	case provider.ProviderGCP:
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:  peer.ContainerID,
			ProviderName: pointer.SetOrNil(string(peer.ProviderName), ""),
			GcpProjectId: pointer.SetOrNil(peer.GCPProjectID, ""),
			NetworkName:  pointer.SetOrNil(peer.NetworkName, ""),
		}
	case provider.ProviderAzure:
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(string(peer.ProviderName), ""),
			AzureDirectoryId:    pointer.SetOrNil(peer.AzureDirectoryID, ""),
			AzureSubscriptionId: pointer.SetOrNil(peer.AzureSubscriptionID, ""),
			ResourceGroupName:   pointer.SetOrNil(peer.ResourceGroupName, ""),
			VnetName:            pointer.SetOrNil(peer.VNetName, ""),
		}
	default:
		panic(fmt.Errorf("unsupported provider %q", peer.ProviderName))
	}
}

func fromAtlasConnection(conn *admin.BaseNetworkPeeringConnectionSettings) *akov2.NetworkPeer {
	switch provider.ProviderName(conn.GetProviderName()) {
	case provider.ProviderAWS:
		return &akov2.NetworkPeer{
			ContainerID:         conn.GetContainerId(),
			ProviderName:        provider.ProviderName(conn.GetProviderName()),
			AccepterRegionName:  conn.GetAccepterRegionName(),
			AWSAccountID:        conn.GetAwsAccountId(),
			RouteTableCIDRBlock: conn.GetRouteTableCidrBlock(),
			VpcID:               conn.GetVpcId(),
		}
	case provider.ProviderGCP:
		return &akov2.NetworkPeer{
			ContainerID:  conn.GetContainerId(),
			ProviderName: provider.ProviderName(conn.GetProviderName()),
			GCPProjectID: conn.GetGcpProjectId(),
			NetworkName:  conn.GetNetworkName(),
		}
	case provider.ProviderAzure:
		return &akov2.NetworkPeer{
			ContainerID:         conn.GetContainerId(),
			ProviderName:        provider.ProviderName(conn.GetProviderName()),
			AzureDirectoryID:    conn.GetAzureDirectoryId(),
			AzureSubscriptionID: conn.GetAzureSubscriptionId(),
			ResourceGroupName:   conn.GetResourceGroupName(),
			VNetName:            conn.GetVnetName(),
		}
	default:
		panic(fmt.Errorf("unsupported provider %q", conn.GetProviderName()))
	}
}

func fromAtlasConnectionList(list []admin.BaseNetworkPeeringConnectionSettings) []akov2.NetworkPeer {
	if list == nil {
		return nil
	}
	peers := make([]akov2.NetworkPeer, 0, len(list))
	for _, conn := range list {
		peers = append(peers, *fromAtlasConnection(&conn))
	}
	return peers
}

func toAtlasContainer(container *ProviderContainer) *admin.CloudProviderContainer {
	return &admin.CloudProviderContainer{
		Id:             pointer.SetOrNil(container.ID, ""),
		ProviderName:   pointer.SetOrNil(string(container.ProviderName), ""),
		AtlasCidrBlock: pointer.SetOrNil(container.AtlasCIDRBlock, ""),
		RegionName:     pointer.SetOrNil(container.RegionName, ""),
		Region:         pointer.SetOrNil(container.Region, ""),
	}
}

func fromAtlasContainer(container *admin.CloudProviderContainer) *ProviderContainer {
	return &ProviderContainer{
		ID:             container.GetId(),
		ProviderName:   provider.ProviderName(container.GetProviderName()),
		AtlasCIDRBlock: container.GetAtlasCidrBlock(),
		RegionName:     container.GetRegionName(),
		Region:         container.GetRegion(),
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
