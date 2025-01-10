/*
Copyright 2025 MongoDB.

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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func init() {
	SchemeBuilder.Register(&AtlasIPAccessList{}, &AtlasIPAccessListList{})
}

// AtlasIPAccessListSpec defines the desired state of AtlasIPAccessList.
type AtlasIPAccessListSpec struct {
	// Project is a reference to AtlasProject resource the user belongs to
	// +kubebuilder:validation:Optional
	Project *common.ResourceRefNamespaced `json:"projectRef,omitempty"`
	// ExternalProject holds the Atlas project ID the user belongs to
	// +kubebuilder:validation:Optional
	ExternalProject *ExternalProjectReference `json:"externalProjectRef,omitempty"`
	// Local credentials
	api.LocalCredentialHolder `json:",inline"`
	// Entry using an IP address in this access list entry.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
	// Range of IP addresses in CIDR notation in this access list entry.
	// +optional
	CIDRBlock string `json:"cidrBlock,omitempty"`
	// Unique identifier of AWS security group in this access list entry.
	// +optional
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"`
	// Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry.
	// +optional
	DeleteAfterDate *metav1.Time `json:"deleteAfterDate,omitempty"`
	// Comment associated with this access list entry.
	// +optional
	Comment string `json:"comment,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=aial
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// AtlasIPAccessList is the Schema for the atlasipaccesslists API.
type AtlasIPAccessList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasIPAccessListSpec          `json:"spec,omitempty"`
	Status status.AtlasIPAccessListStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasIPAccessListList contains a list of AtlasIPAccessList.
type AtlasIPAccessListList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasIPAccessList `json:"items"`
}

func (ial *AtlasIPAccessList) AtlasProjectObjectKey() client.ObjectKey {
	ns := ial.Namespace
	if ial.Spec.Project.Namespace != "" {
		ns = ial.Spec.Project.Namespace
	}

	return kube.ObjectKey(ns, ial.Spec.Project.Name)
}

func (ial *AtlasIPAccessList) Credentials() *api.LocalObjectReference {
	return ial.Spec.Credentials()
}

func (ial *AtlasIPAccessList) GetStatus() api.Status {
	return ial.Status
}

func (ial *AtlasIPAccessList) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	ial.Status.Conditions = conditions
	ial.Status.ObservedGeneration = ial.ObjectMeta.Generation

	for _, o := range options {
		v := o.(status.AtlasIPAccessListStatusOption)
		v(&ial.Status)
	}
}
