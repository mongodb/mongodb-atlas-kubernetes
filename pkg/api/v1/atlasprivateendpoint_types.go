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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
// +kubebuilder:validation:XValidation:rule="(oldSelf.provider == 'AWS' && self.provider == 'AWS' && has(self.awsConfiguration)) || self.provider != 'AWS'",message="must define a configuration for AWS provider"
// +kubebuilder:validation:XValidation:rule="(oldSelf.provider == 'AZURE' && self.provider == 'AZURE' && has(self.azureConfiguration)) || self.provider != 'AZURE'",message="must define a configuration for Azure provider"
// +kubebuilder:validation:XValidation:rule="(oldSelf.provider == 'AWS' && self.provider == 'GCP' && has(self.gcpConfiguration)) || self.provider != 'GCP'",message="must define a configuration for GCP provider"

// AtlasPrivateEndpointSpec is the specification of the desired configuration of a project private endpoint
type AtlasPrivateEndpointSpec struct {
	// Project is a reference to AtlasProject resource the user belongs to
	// +kubebuilder:validation:Optional
	Project *common.ResourceRefNamespaced `json:"projectRef,omitempty"`
	// ExternalProject holds the Atlas project ID the user belongs to
	// +kubebuilder:validation:Optional
	ExternalProject *ExternalProjectReference `json:"externalProjectRef,omitempty"`

	// Local credentials
	api.LocalCredentialHolder `json:",inline"`

	// Name of the cloud service provider for which you want to create the private endpoint service.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`
	// Region of the chosen cloud provider in which you want to create the private endpoint service.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// AWSConfiguration is the specific AWS settings for the private endpoint
	// +kubebuilder:validation:Optional
	AWSConfiguration *AWSPrivateEndpointConfiguration `json:"awsConfiguration,omitempty"`
	// AzureConfiguration is the specific Azure settings for the private endpoint
	// +kubebuilder:validation:Optional
	AzureConfiguration *AzurePrivateEndpointConfiguration `json:"azureConfiguration,omitempty"`
	// GCPConfiguration is the specific Google Cloud settings for the private endpoint
	// +kubebuilder:validation:Optional
	GCPConfiguration *GCPPrivateEndpointConfiguration `json:"gcpConfiguration,omitempty"`
}

type AWSPrivateEndpointConfiguration struct {
	// ID that identifies the private endpoint's network interface that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	ID string `json:"id"`
}

type AzurePrivateEndpointConfiguration struct {
	// ID that identifies the private endpoint's network interface that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	ID string `json:"id"`
	// IP address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	IP string `json:"ipAddress"`
}

type GCPPrivateEndpointConfiguration struct {
	// ProjectID that identifies the Google Cloud project in which you created the endpoints.
	// +kubebuilder:validation:Required
	ProjectID string `json:"projectId"`
	// GroupName is the label that identifies a set of endpoints.
	// +kubebuilder:validation:Required
	GroupName string `json:"groupName"`
	// Endpoints is the list of individual private endpoints that comprise this endpoint group.
	// +kubebuilder:validation:Required
	Endpoints []GCPPrivateEndpoint `json:"endpoints"`
}

type GCPPrivateEndpoint struct {
	// Name that identifies the Google Cloud consumer forwarding rule that you created.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// IP address to which this Google Cloud consumer forwarding rule resolves.
	// +kubebuilder:validation:Required
	IP string `json:"ipAddress"`
}

// AtlasPrivateEndpointStatus is the most recent observed status of the AtlasPrivateEndpoint cluster. Read-only.
type AtlasPrivateEndpointStatus struct {
	api.Common `json:",inline"`
	// ServiceID is the unique identifier of the private endpoint service in Atlas
	ServiceID string `json:"serviceId,omitempty"`
	// InterfaceID is the unique identifier of the private endpoint interface in Atlas
	InterfaceID string `json:"interfaceId,omitempty"`
	// Status is the state of the private endpoint connection
	Status string `json:"status,omitempty"`
	// Error is the description of the failure occurred when configuring the private endpoint
	Error string `json:"error,omitempty"`
	// ServiceName is the unique identifier of the Amazon Web Services (AWS) PrivateLink endpoint service or Azure Private Link Service managed by Atlas
	ServiceName string `json:"serviceName,omitempty"`
	// ServiceResourceID is the root-relative path that identifies of the Azure Private Link Service
	ServiceResourceID string `json:"serviceResourceId,omitempty"`
	// ServiceAttachmentNames is the list of URLs that identifies endpoints that Atlas can use to access one service across the private connection
	ServiceAttachmentNames []string `json:"serviceAttachmentNames,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=pe
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Region",type=string,JSONPath=`.spec.region`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// The AtlasPrivateEndpoint custom resource definition (CRD) defines a desired [Private Endpoint](https://www.mongodb.com/docs/atlas/security-private-endpoint/#std-label-private-endpoint-overview) configuration for an Atlas project.
// It allows a private connection between your cloud provider and Atlas that doesn't send information through a public network.
//
// You can use private endpoints to create a unidirectional connection to Atlas clusters from your virtual network.
type AtlasPrivateEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasPrivateEndpointSpec   `json:"spec,omitempty"`
	Status AtlasPrivateEndpointStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasPrivateEndpointList contains a list of AtlasPrivateEndpoint
type AtlasPrivateEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasPrivateEndpoint `json:"items"`
}

// +kubebuilder:object:generate=false

type AtlasPrivateEndpointStatusOption func(s *AtlasPrivateEndpointStatus)

func (pe *AtlasPrivateEndpoint) GetStatus() api.Status {
	return pe.Status
}

func (pe *AtlasPrivateEndpoint) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	pe.Status.Conditions = conditions
	pe.Status.ObservedGeneration = pe.ObjectMeta.Generation

	for _, o := range options {
		v := o.(AtlasPrivateEndpointStatusOption)
		v(&pe.Status)
	}
}

func init() {
	SchemeBuilder.Register(&AtlasPrivateEndpoint{}, &AtlasPrivateEndpointList{})
}
