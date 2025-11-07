//Copyright 2024 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

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

// AtlasNetworkPeeringSpec defines the desired state of AtlasNetworkPeering
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
// +kubebuilder:validation:XValidation:rule="(has(self.containerRef.name) && !has(self.containerRef.id)) || (!has(self.containerRef.name) && has(self.containerRef.id))",message="must either have a container Atlas id or Kubernetes name, but not both (or neither)"
// +kubebuilder:validation:XValidation:rule="(self.containerRef.name == oldSelf.containerRef.name) || (!has(self.containerRef.name) && !has(oldSelf.containerRef.name))",message="container ref name is immutable"
// +kubebuilder:validation:XValidation:rule="(self.containerRef.id == oldSelf.containerRef.id) || (!has(self.containerRef.id) && !has(oldSelf.containerRef.id))",message="container ref id is immutable"
// +kubebuilder:validation:XValidation:rule="(self.id == oldSelf.id) || (!has(self.id) && !has(oldSelf.id))",message="id is immutable"
type AtlasNetworkPeeringSpec struct {
	ProjectDualReference `json:",inline"`

	ContainerRef ContainerDualReference `json:"containerRef"`

	AtlasNetworkPeeringConfig `json:",inline"`
}

// ContainerDualReference refers to a Network Container either by Kubernetes name or Atlas ID.
type ContainerDualReference struct {
	// Name of the container Kubernetes resource, must be present in the same namespace.
	// Use either name or ID, not both.
	// +optional
	Name string `json:"name,omitempty"`

	// ID is the Atlas identifier of the Network Container Atlas resource this Peering Connection relies on.
	// Use either name or ID, not both.
	// +optional
	ID string `json:"id,omitempty"`
}

// AtlasNetworkPeeringConfig defines the Atlas specifics of the desired state of Peering Connections
type AtlasNetworkPeeringConfig struct {
	// ID is the peering identifier for an already existent network peering to be managed by the operator.
	// This field is immutable.
	// +optional
	ID string `json:"id,omitempty"`

	// Name of the cloud service provider for which you want to create the network peering service.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`

	// AWSConfiguration is the specific AWS settings for network peering.
	// +kubebuilder:validation:Optional
	AWSConfiguration *AWSNetworkPeeringConfiguration `json:"awsConfiguration,omitempty"`
	// AzureConfiguration is the specific Azure settings for network peering.
	// +kubebuilder:validation:Optional
	AzureConfiguration *AzureNetworkPeeringConfiguration `json:"azureConfiguration,omitempty"`
	// GCPConfiguration is the specific Google Cloud settings for network peering.
	// +kubebuilder:validation:Optional
	GCPConfiguration *GCPNetworkPeeringConfiguration `json:"gcpConfiguration,omitempty"`
}

// AWSNetworkPeeringConfiguration defines tha Atlas desired state for AWS
type AWSNetworkPeeringConfiguration struct {
	// AccepterRegionName is the provider region name of user's vpc in AWS native region format.
	// +kubebuilder:validation:Required
	AccepterRegionName string `json:"accepterRegionName"`
	// AccountID of the user's vpc.
	// +kubebuilder:validation:Required
	AWSAccountID string `json:"awsAccountId,omitempty"`
	// User VPC CIDR.
	// +kubebuilder:validation:Required
	RouteTableCIDRBlock string `json:"routeTableCidrBlock,omitempty"`
	// AWS VPC ID.
	// +kubebuilder:validation:Required
	VpcID string `json:"vpcId,omitempty"`
}

// AzureNetworkPeeringConfiguration defines tha Atlas desired state for Azure
type AzureNetworkPeeringConfiguration struct {
	//AzureDirectoryID is the unique identifier for an Azure AD directory.
	// +kubebuilder:validation:Required
	AzureDirectoryID string `json:"azureDirectoryId,omitempty"`
	// AzureSubscriptionID is the unique identifier of the Azure subscription in which the VNet resides.
	// +kubebuilder:validation:Required
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	//ResourceGroupName is the name of your Azure resource group.
	// +kubebuilder:validation:Required
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// VNetName is name of your Azure VNet. Its applicable only for Azure.
	// +kubebuilder:validation:Required
	VNetName string `json:"vNetName,omitempty"`
}

// GCPNetworkPeeringConfiguration defines tha Atlas desired state for Google
type GCPNetworkPeeringConfiguration struct {
	// User GCP Project ID. Its applicable only for GCP.
	// +kubebuilder:validation:Required
	GCPProjectID string `json:"gcpProjectId,omitempty"`
	// GCP Network Peer Name. Its applicable only for GCP.
	// +kubebuilder:validation:Required
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
