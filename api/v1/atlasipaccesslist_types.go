//Copyright 2025 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	SchemeBuilder.Register(&AtlasIPAccessList{}, &AtlasIPAccessListList{})
}

// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"

// AtlasIPAccessListSpec defines the desired state of AtlasIPAccessList.
type AtlasIPAccessListSpec struct {
	// ProjectReference is the dual external or kubernetes reference with access credentials.
	ProjectDualReference `json:",inline"`
	// Entries is the list of IP Access to be managed.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Entries []IPAccessEntry `json:"entries"`
}

// +kubebuilder:validation:XValidation:rule="!(has(self.ipAddress) && (has(self.cidrBlock) || has(self.awsSecurityGroup))) && !(has(self.cidrBlock) && has(self.awsSecurityGroup))",message="Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set."

type IPAccessEntry struct {
	// Entry using an IP address in this access list entry.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
	// Range of IP addresses in CIDR notation in this access list entry.
	// +optional
	CIDRBlock string `json:"cidrBlock,omitempty"`
	// Unique identifier of AWS security group in this access list entry.
	// +optional
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"`
	// Date and time after which Atlas deletes the temporary access list entry.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
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
// +kubebuilder:resource:categories=atlas,shortName=aip
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

func (ial *AtlasIPAccessList) Credentials() *api.LocalObjectReference {
	return ial.Spec.ConnectionSecret
}

func (ial *AtlasIPAccessList) ProjectDualRef() *ProjectDualReference {
	return &ial.Spec.ProjectDualReference
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
