/*
Copyright 2024 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

type AtlasNetworkPeeringConfig struct {
	// Name of the cloud service provider for which you want to create the network peering service.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`
	// ID of the network peer container. If not set, operator will create a new container with ContainerRegion and AtlasCIDRBlock input.
	// +optional
	ContainerID string `json:"containerId"`

	// AWSConfiguration is the specific AWS settings for network peering
	// +kubebuilder:validation:Optional
	AWSConfiguration *AWSNetworkPeeringConfiguration `json:"awsConfiguration,omitempty"`
	// AzureConfiguration is the specific Azure settings for network peering
	// +kubebuilder:validation:Optional
	AzureConfiguration *AzureNetworkPeeringConfiguration `json:"azureConfiguration,omitempty"`
	// GCPConfiguration is the specific Google Cloud settings for network peering
	// +kubebuilder:validation:Optional
	GCPConfiguration *GCPNetworkPeeringConfiguration `json:"gcpConfiguration,omitempty"`
}

type AtlasProviderContainerConfig struct {
	// ContainerRegion is the provider region name of Atlas network peer container. If not set, AccepterRegionName is used.
	// +optional
	ContainerRegion string `json:"containerRegion"`
	// Atlas CIDR. It needs to be set if ContainerID is not set.
	// +optional
	AtlasCIDRBlock string `json:"atlasCidrBlock"`
}

type AWSNetworkPeeringConfiguration struct {
	// AccepterRegionName is the provider region name of user's vpc.
	AccepterRegionName string `json:"accepterRegionName"`
	// AccountID of the user's vpc.
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// User VPC CIDR.
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	// AWS VPC ID.
	VpcID string `json:"vpcId,omitempty"`
}

type AzureNetworkPeeringConfiguration struct {
	//AzureDirectoryID is the unique identifier for an Azure AD directory.
	AzureDirectoryID string `json:"azureDirectoryId,omitempty"`
	// AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	//ResourceGroupName is the name of your Azure resource group.
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// VNetName is name of your Azure VNet. Its applicable only for Azure.
	VNetName string `json:"vnetName,omitempty"`
}

type GCPNetworkPeeringConfiguration struct {
	// User GCP Project ID. Its applicable only for GCP.
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// GCP Network Peer Name. Its applicable only for GCP.
	NetworkName string `json:"networkName,omitempty"`
}
