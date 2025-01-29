package networkcontainer

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type NetworkContainer struct {
	akov2.AtlasNetworkContainerConfig
	ID          string
	Provider    string
	Provisioned bool
	AWSStatus   *AWSContainerStatus
	AzureStatus *AzureContainerStatus
	GCPStatus   *GoogleContainerStatus
}

type AWSContainerStatus struct {
	VpcID string
}

type AzureContainerStatus struct {
	AzureSubscriptionID string
	VnetName            string
}

type GoogleContainerStatus struct {
	GCPProjectID string
	NetworkName  string
}

func NewNetworkContainerSpec(provider string, config *akov2.AtlasNetworkContainerConfig) *NetworkContainer {
	return &NetworkContainer{
		Provider:                    provider,
		AtlasNetworkContainerConfig: *config,
	}
}

func ApplyNetworkContainerStatus(containerStatus *status.AtlasNetworkContainerStatus, container *NetworkContainer) {
	containerStatus.ID = container.ID
	containerStatus.Provisioned = container.Provisioned
}

func toAtlas(container *NetworkContainer) *admin.CloudProviderContainer {
	cpc := &admin.CloudProviderContainer{
		Id:             pointer.SetOrNil(container.ID, ""),
		ProviderName:   pointer.SetOrNil(container.Provider, ""),
		AtlasCidrBlock: pointer.SetOrNil(container.CIDRBlock, ""),
	}
	if cpc.GetProviderName() == string(provider.ProviderAWS) {
		cpc.RegionName = pointer.SetOrNil(container.Region, "")
	} else {
		cpc.Region = pointer.SetOrNil(container.Region, "")
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
		pc.GCPStatus = fromAtlasGoogleStatus(container)
	}
	return pc
}

func fromAtlasNoStatus(container *admin.CloudProviderContainer) *NetworkContainer {
	region := container.GetRegion()
	if container.GetProviderName() == string(provider.ProviderAWS) {
		region = container.GetRegionName()
	}
	return &NetworkContainer{
		ID:       container.GetId(),
		Provider: container.GetProviderName(),
		AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
			CIDRBlock: container.GetAtlasCidrBlock(),
			Region:    region,
		},
	}
}

func fromAtlasAWSStatus(container *admin.CloudProviderContainer) *AWSContainerStatus {
	if container.VpcId == nil {
		return nil
	}
	return &AWSContainerStatus{
		VpcID: container.GetVpcId(),
	}
}

func fromAtlasAzureStatus(container *admin.CloudProviderContainer) *AzureContainerStatus {
	if container.AzureSubscriptionId == nil && container.VnetName == nil {
		return nil
	}
	return &AzureContainerStatus{
		AzureSubscriptionID: container.GetAzureSubscriptionId(),
		VnetName:            container.GetVnetName(),
	}
}

func fromAtlasGoogleStatus(container *admin.CloudProviderContainer) *GoogleContainerStatus {
	if container.GcpProjectId == nil && container.NetworkName == nil {
		return nil
	}
	return &GoogleContainerStatus{
		GCPProjectID: container.GetGcpProjectId(),
		NetworkName:  container.GetNetworkName(),
	}
}
