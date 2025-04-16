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

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/api"

// AtlasNetworkPeeringStatus is a status for the AtlasNetworkPeering Custom resource.
// Not the one included in the AtlasProject
type AtlasNetworkPeeringStatus struct {
	api.Common `json:",inline"`

	// ID recrods the identified of the peer created by Atlas
	ID string `json:"id,omitempty"`

	// Status describes the last status seen for the network peering setup
	Status string `json:"status,omitempty"`

	// AWSStatus contains AWS only related status information
	AWSStatus *AWSPeeringStatus `json:"awsStatus,omitempty"`

	// AzureStatus contains Azure only related status information
	AzureStatus *AzurePeeringStatus `json:"azureStatus,omitempty"`

	// GCPStatus contains GCP only related status information
	GCPStatus *GCPPeeringStatus `json:"gcpStatus,omitempty"`
}

// AWSPeeringStatus contains AWS only related status for network peering & container
type AWSPeeringStatus struct {
	// VpcID is AWS VPC id on the Atlas side
	VpcID string `json:"vpcId,omitempty"`

	// ConnectionID is the AWS VPC peering connection ID
	ConnectionID string `json:"connectionId,omitempty"`
}

// AzurePeeringStatus contains Azure only related status information
type AzurePeeringStatus struct {
	// AzureSubscriptionID is Azure Subscription id on the Atlas side
	AzureSubscriptionID string `json:"azureSubscriptionIDpcId,omitempty"`

	// VnetName is Azure network on the Atlas side
	VnetName string `json:"vNetName,omitempty"`
}

// GCPPeeringStatus contains GCP only related status information
type GCPPeeringStatus struct {
	// GCPProjectID is GCP project on the Atlas side
	GCPProjectID string `json:"gcpProjectID,omitempty"`

	// NetworkName is GCP network on the Atlas side
	NetworkName string `json:"networkName,omitempty"`
}

// +kubebuilder:object:generate=false

type AtlasNetworkPeeringStatusOption func(s *AtlasNetworkPeeringStatus)
