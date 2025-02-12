package networkcontainer

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type NetworkContainerConfig struct {
	Provider string
	akov2.AtlasNetworkContainerConfig
}

type NetworkContainer struct {
	NetworkContainerConfig
	ID          string
	Provisioned bool
	AWSStatus   *status.AWSContainerStatus
	AzureStatus *status.AzureContainerStatus
	GCPStatus   *status.GCPContainerStatus
}

func NewNetworkContainerConfig(provider string, config *akov2.AtlasNetworkContainerConfig) *NetworkContainerConfig {
	return &NetworkContainerConfig{
		Provider:                    provider,
		AtlasNetworkContainerConfig: *config,
	}
}

func ApplyNetworkContainerStatus(containerStatus *status.AtlasNetworkContainerStatus, container *NetworkContainer) {
	containerStatus.ID = container.ID
	containerStatus.Provisioned = container.Provisioned
	containerStatus.AWSStatus = container.AWSStatus
	containerStatus.AzureStatus = container.AzureStatus
	containerStatus.GCPStatus = container.GCPStatus
}

func toAtlas(container *NetworkContainer) *admin.CloudProviderContainer {
	cpc := toAtlasConfig(&container.NetworkContainerConfig)
	cpc.Id = pointer.SetOrNil(container.ID, "")
	return cpc
}

func toAtlasConfig(cfg *NetworkContainerConfig) *admin.CloudProviderContainer {
	cpc := &admin.CloudProviderContainer{
		ProviderName:   pointer.SetOrNil(cfg.Provider, ""),
		AtlasCidrBlock: pointer.SetOrNil(cfg.CIDRBlock, ""),
	}
	if cpc.GetProviderName() == string(provider.ProviderAWS) {
		cpc.RegionName = pointer.SetOrNil(cfg.Region, "")
	} else {
		cpc.Region = pointer.SetOrNil(cfg.Region, "")
	}
	return cpc
}

func fromAtlas(container *admin.CloudProviderContainer) *NetworkContainer {
	pc := fromAtlasNoStatus(container)
	pc.Provisioned = container.GetProvisioned()
	switch provider.ProviderName(pc.Provider) {
	case provider.ProviderAWS:
		pc.AWSStatus = fromAtlasAWSStatus(container)
	case provider.ProviderAzure:
		pc.AzureStatus = fromAtlasAzureStatus(container)
	case provider.ProviderGCP:
		pc.GCPStatus = fromAtlasGCPStatus(container)
	}
	return pc
}

func fromAtlasNoStatus(container *admin.CloudProviderContainer) *NetworkContainer {
	region := container.GetRegion()
	if container.GetProviderName() == string(provider.ProviderAWS) {
		region = container.GetRegionName()
	}
	return &NetworkContainer{
		NetworkContainerConfig: NetworkContainerConfig{
			Provider: container.GetProviderName(),
			AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
				CIDRBlock: container.GetAtlasCidrBlock(),
				Region:    region,
			},
		},
		ID: container.GetId(),
	}
}

func fromAtlasAWSStatus(container *admin.CloudProviderContainer) *status.AWSContainerStatus {
	if container.VpcId == nil {
		return nil
	}
	return &status.AWSContainerStatus{
		VpcID: container.GetVpcId(),
	}
}

func fromAtlasAzureStatus(container *admin.CloudProviderContainer) *status.AzureContainerStatus {
	if container.AzureSubscriptionId == nil && container.VnetName == nil {
		return nil
	}
	return &status.AzureContainerStatus{
		AzureSubscriptionID: container.GetAzureSubscriptionId(),
		VnetName:            container.GetVnetName(),
	}
}

func fromAtlasGCPStatus(container *admin.CloudProviderContainer) *status.GCPContainerStatus {
	if container.GcpProjectId == nil && container.NetworkName == nil {
		return nil
	}
	return &status.GCPContainerStatus{
		GCPProjectID: container.GetGcpProjectId(),
		NetworkName:  container.GetNetworkName(),
	}
}
