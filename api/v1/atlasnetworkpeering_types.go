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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasNetworkPeering{}, &AtlasNetworkPeeringList{})
}

// AtlasNetworkPeering is the Schema for the AtlasNetworkPeering API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Id",type=string,JSONPath=`.status.id`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Container Provisioned",type=string,JSONPath=`.status.containerProvisioned`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=anp
type AtlasNetworkPeering struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasNetworkPeeringSpec          `json:"spec,omitempty"`
	Status status.AtlasNetworkPeeringStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasNetworkPeeringList contains a list of AtlasNetworkPeering
type AtlasNetworkPeeringList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasNetworkPeering `json:"items"`
}

// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"

// AtlasNetworkPeeringSpec defines the desired state of AtlasNetworkPeering
type AtlasNetworkPeeringSpec struct {
	ProjectDualReference `json:",inline"`

	AtlasNetworkPeeringConfig `json:",inline"`

	AtlasProviderContainerConfig `json:",inline"`
}

// AtlasNetworkPeeringConfig defines the Atlas specifics of the desired state of Peering Connections
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

// AtlasProviderContainerConfig defines the Atlas specifics of the desired state of Peering Container
type AtlasProviderContainerConfig struct {
	// ContainerRegion is the provider region name of Atlas network peer container in Atlas region format
	// +optional
	ContainerRegion string `json:"containerRegion"`
	// Atlas CIDR. It needs to be set if ContainerID is not set.
	// +optional
	AtlasCIDRBlock string `json:"atlasCidrBlock"`
}

// AWSNetworkPeeringConfiguration defines tha Atlas desired state for AWS
type AWSNetworkPeeringConfiguration struct {
	// AccepterRegionName is the provider region name of user's vpc in AWS native region format
	AccepterRegionName string `json:"accepterRegionName"`
	// AccountID of the user's vpc.
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// User VPC CIDR.
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	// AWS VPC ID.
	VpcID string `json:"vpcId,omitempty"`
}

// AzureNetworkPeeringConfiguration defines tha Atlas desired state for Azure
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

// GCPNetworkPeeringConfiguration defines tha Atlas desired state for Google
type GCPNetworkPeeringConfiguration struct {
	// User GCP Project ID. Its applicable only for GCP.
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// GCP Network Peer Name. Its applicable only for GCP.
	NetworkName string `json:"networkName,omitempty"`
}

func (np *AtlasNetworkPeering) GetStatus() api.Status {
	return np.Status
}

func (np *AtlasNetworkPeering) Credentials() *api.LocalObjectReference {
	return np.Spec.ConnectionSecret
}

func (np *AtlasNetworkPeering) ProjectDualRef() *ProjectDualReference {
	return &np.Spec.ProjectDualReference
}

func (np *AtlasNetworkPeering) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	np.Status.Conditions = conditions
	np.Status.ObservedGeneration = np.ObjectMeta.Generation

	for _, o := range options {
		v := o.(status.AtlasNetworkPeeringStatusOption)
		v(&np.Status)
	}
}
