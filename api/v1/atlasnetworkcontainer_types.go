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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasNetworkContainer{}, &AtlasNetworkContainerList{})
}

// AtlasNetworkContainer is the Schema for the AtlasNetworkContainer API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Id",type=string,JSONPath=`.status.id`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=anc
type AtlasNetworkContainer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasNetworkContainerSpec          `json:"spec,omitempty"`
	Status status.AtlasNetworkContainerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasNetworkContainerList contains a list of AtlasNetworkContainer.
type AtlasNetworkContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasNetworkContainer `json:"items"`
}

// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
// +kubebuilder:validation:XValidation:rule="(self.provider == 'GCP' && !has(self.region)) || (self.provider != 'GCP')",message="must not set region for GCP containers"
// +kubebuilder:validation:XValidation:rule="((self.provider == 'AWS' || self.provider == 'AZURE') && has(self.region)) || (self.provider == 'GCP')",message="must set region for AWS and Azure containers"
// +kubebuilder:validation:XValidation:rule="(self.id == oldSelf.id) || (!has(self.id) && !has(oldSelf.id))",message="id is immutable"
// +kubebuilder:validation:XValidation:rule="(self.region == oldSelf.region) || (!has(self.region) && !has(oldSelf.region))",message="region is immutable"

// AtlasNetworkContainerSpec defines the desired state of an AtlasNetworkContainer.
type AtlasNetworkContainerSpec struct {
	ProjectDualReference `json:",inline"`

	// Provider is the name of the cloud provider hosting the network container.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`

	// Collection of settings that configures the network container for a virtual private connection in a cloud provider.
	AtlasNetworkContainerConfig `json:",inline"`
}

// AtlasNetworkContainerConfig defines the Atlas specifics of the desired state of a Network Container.
type AtlasNetworkContainerConfig struct {
	// ID is the container identifier for an already existent network container to be managed by the operator.
	// This field can be used in conjunction with cidrBlock to update the cidrBlock of an existing container.
	// This field is immutable.
	// +optional
	ID string `json:"id,omitempty"`

	// ContainerRegion is the provider region name of Atlas network peer container in Atlas region format
	// This is required by AWS and Azure, but not used by GCP.
	// This field is immutable, Atlas does not admit network container changes.
	// +optional
	Region string `json:"region,omitempty"`

	// Atlas CIDR. It needs to be set if ContainerID is not set.
	// +optional
	CIDRBlock string `json:"cidrBlock"`
}

func (np *AtlasNetworkContainer) GetStatus() api.Status {
	return np.Status
}

func (np *AtlasNetworkContainer) Credentials() *api.LocalObjectReference {
	return np.Spec.ConnectionSecret
}

func (np *AtlasNetworkContainer) ProjectDualRef() *ProjectDualReference {
	return &np.Spec.ProjectDualReference
}

func (np *AtlasNetworkContainer) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	np.Status.Conditions = conditions
	np.Status.ObservedGeneration = np.ObjectMeta.Generation

	for _, o := range options {
		v := o.(status.AtlasNetworkContainerStatusOption)
		v(&np.Status)
	}
}
