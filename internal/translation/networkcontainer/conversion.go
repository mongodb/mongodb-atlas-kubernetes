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

package networkcontainer

import (
	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type NetworkContainerConfig struct {
	Provider string
	akov2.AtlasNetworkContainerConfig
}

type AWSContainerStatus struct {
	VpcID       string
	ContainerID string
}

type AzureContainerStatus struct {
	AzureSubscriptionID string
	VnetName            string
}

type GCPContainerStatus struct {
	GCPProjectID string
	NetworkName  string
}

type NetworkContainer struct {
	NetworkContainerConfig
	ID          string
	Provisioned bool
	AWSStatus   *AWSContainerStatus
	AzureStatus *AzureContainerStatus
	GCPStatus   *GCPContainerStatus
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

func fromAtlasAWSStatus(container *admin.CloudProviderContainer) *AWSContainerStatus {
	if vpcID, ok := container.GetVpcIdOk(); ok {
		return &AWSContainerStatus{
			VpcID: *vpcID,
		}
	}
	return nil
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

func fromAtlasGCPStatus(container *admin.CloudProviderContainer) *GCPContainerStatus {
	if container.GcpProjectId == nil && container.NetworkName == nil {
		return nil
	}
	return &GCPContainerStatus{
		GCPProjectID: container.GetGcpProjectId(),
		NetworkName:  container.GetNetworkName(),
	}
}
