// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func init() {
	SchemeBuilder.Register(&AtlasPrivateEndpoint{}, &AtlasPrivateEndpointList{})
}

// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"

// AtlasPrivateEndpointSpec is the specification of the desired configuration of a project private endpoint
type AtlasPrivateEndpointSpec struct {
	// ProjectReference is the dual external or kubernetes reference with access credentials.
	ProjectDualReference `json:",inline"`

	// Name of the cloud service provider for which you want to create the private endpoint service.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:Required
	Provider string `json:"provider"`
	// Region of the chosen cloud provider in which you want to create the private endpoint service.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// AWSConfiguration is the specific AWS settings for the private endpoint.
	// +listType=map
	// +listMapKey=id
	// +kubebuilder:validation:Optional
	AWSConfiguration []AWSPrivateEndpointConfiguration `json:"awsConfiguration,omitempty"`
	// AzureConfiguration is the specific Azure settings for the private endpoint.
	// +listType=map
	// +listMapKey=id
	// +kubebuilder:validation:Optional
	AzureConfiguration []AzurePrivateEndpointConfiguration `json:"azureConfiguration,omitempty"`
	// GCPConfiguration is the specific Google Cloud settings for the private endpoint.
	// +listType=map
	// +listMapKey=groupName
	// +kubebuilder:validation:Optional
	GCPConfiguration []GCPPrivateEndpointConfiguration `json:"gcpConfiguration,omitempty"`
}

// AWSPrivateEndpointConfiguration holds the AWS configuration done on customer network.
type AWSPrivateEndpointConfiguration struct {
	// ID that identifies the private endpoint's network interface that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	ID string `json:"id"`
}

// AzurePrivateEndpointConfiguration holds the Azure configuration done on customer network.
type AzurePrivateEndpointConfiguration struct {
	// ID that identifies the private endpoint's network interface that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	ID string `json:"id"`
	// IP address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	// +kubebuilder:validation:Required
	IP string `json:"ipAddress"`
}

// GCPPrivateEndpointConfiguration holds the GCP configuration done on customer network.
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

// GCPPrivateEndpoint holds the GCP forwarding rules configured on customer network.
type GCPPrivateEndpoint struct {
	// Name that identifies the Google Cloud consumer forwarding rule that you created.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// IP address to which this Google Cloud consumer forwarding rule resolves.
	// +kubebuilder:validation:Required
	IP string `json:"ipAddress"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=ape
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

	Spec   AtlasPrivateEndpointSpec          `json:"spec,omitempty"`
	Status status.AtlasPrivateEndpointStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasPrivateEndpointList contains a list of AtlasPrivateEndpoint
type AtlasPrivateEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasPrivateEndpoint `json:"items"`
}

func (pe *AtlasPrivateEndpoint) AtlasProjectObjectKey() client.ObjectKey {
	ns := pe.Namespace
	if pe.Spec.ProjectRef.Namespace != "" {
		ns = pe.Spec.ProjectRef.Namespace
	}

	return kube.ObjectKey(ns, pe.Spec.ProjectRef.Name)
}

func (pe *AtlasPrivateEndpoint) GetStatus() api.Status {
	return pe.Status
}

func (pe *AtlasPrivateEndpoint) Credentials() *api.LocalObjectReference {
	return pe.Spec.ConnectionSecret
}

func (pe *AtlasPrivateEndpoint) ProjectDualRef() *ProjectDualReference {
	return &pe.Spec.ProjectDualReference
}

func (pe *AtlasPrivateEndpoint) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	pe.Status.Conditions = conditions
	pe.Status.ObservedGeneration = pe.ObjectMeta.Generation

	for _, o := range options {
		v := o.(status.AtlasPrivateEndpointStatusOption)
		v(&pe.Status)
	}
}
