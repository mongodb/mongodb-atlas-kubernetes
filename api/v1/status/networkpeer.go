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

package status

import (
	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

type AtlasNetworkPeer struct {
	// Unique identifier for NetworkPeer.
	ID string `json:"id"`
	// Cloud provider for which you want to retrieve a network peer.
	ProviderName provider.ProviderName `json:"providerName"`
	// Region for which you want to create the network peer. It isn't needed for GCP
	Region string `json:"region"`
	// Status of the network peer. Applicable only for AWS.
	StatusName string `json:"statusName,omitempty"`
	// Error state of the network peer. Applicable only for Azure.
	ErrorState string `json:"errorState,omitempty"`
	// Error state of the network peer. Applicable only for AWS.
	ErrorStateName string `json:"errorStateName,omitempty"`
	// Error state of the network peer. Applicable only for GCP.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Unique identifier of the network peer connection. Applicable only for AWS.
	ConnectionID string `json:"connectionId,omitempty"`
	// Status of the network peer. Applicable only for GCP and Azure.
	Status string `json:"status,omitempty"`
	// VPC is general purpose field for storing the name of the VPC.
	// VPC is vpcID for AWS, user networkName for GCP, and vnetName for Azure.
	VPC string `json:"vpc,omitempty"`
	// ProjectID of the user's vpc. Applicable only for GCP.
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// Atlas Network Name. Applicable only for GCP. It's needed to add network peer connection.
	AtlasNetworkName string `json:"atlasNetworkName,omitempty"`
	// ProjectID of Atlas container. Applicable only for GCP. It's needed to add network peer connection.
	AtlasGCPProjectID string `json:"atlasGcpProjectId,omitempty"`
	// ContainerID of Atlas network peer container.
	ContainerID string `json:"containerId,omitempty"`
}

func NewNetworkPeerStatus(atlasPeer admin.BaseNetworkPeeringConnectionSettings, providerName provider.ProviderName, vpcName string, container admin.CloudProviderContainer) AtlasNetworkPeer {
	return AtlasNetworkPeer{
		ID:                atlasPeer.GetId(),
		ProviderName:      providerName,
		Region:            atlasPeer.GetAccepterRegionName(),
		StatusName:        atlasPeer.GetStatusName(),
		ErrorMessage:      atlasPeer.GetErrorMessage(),
		ErrorState:        atlasPeer.GetErrorState(),
		ErrorStateName:    atlasPeer.GetErrorStateName(),
		ConnectionID:      atlasPeer.GetConnectionId(),
		Status:            atlasPeer.GetStatus(),
		VPC:               vpcName,
		AtlasNetworkName:  container.GetNetworkName(),
		AtlasGCPProjectID: container.GetGcpProjectId(),
		ContainerID:       container.GetId(),
		GCPProjectID:      atlasPeer.GetGcpProjectId(),
	}
}

func (in *AtlasNetworkPeer) GetStatus() string {
	if in.StatusName == "" {
		return in.Status
	}
	return in.StatusName
}
