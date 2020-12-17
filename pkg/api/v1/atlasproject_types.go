/*
Copyright 2020 MongoDB.

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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Important:
// The procedure working with this file:
// 1. Edit the file
// 1. Run "make generate" to regenerate code
// 2. Run "make manifests" to regenerate the CRD

// Dev note: this file should be placed in "v1" package (not the nested one) as 'make manifests' doesn't generate the proper
// CRD - this may be addressed later as having a subpackage may get a much nicer code

func init() {
	SchemeBuilder.Register(&AtlasProject{}, &AtlasProjectList{})
}

// AtlasProjectSpec defines the desired state of Project in Atlas
type AtlasProjectSpec struct {

	// Name is the name of the Project that is created in Atlas by the Operator if it doesn't exist yet.
	Name string `json:"name"`

	// ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
	// Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.
	// +optional
	ConnectionSecret *SecretRef `json:"connectionSecretRef,omitempty"`

	// ProjectIPAccessList allows to enable the IP Access List for the Project. See more information at
	// https://docs.atlas.mongodb.com/reference/api/ip-access-list/add-entries-to-access-list/
	// +optional
	ProjectIPAccessList []ProjectIPAccessList `json:"projectIpAccessList,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com

// AtlasProject is the Schema for the atlasprojects API
type AtlasProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasProjectSpec          `json:"spec,omitempty"`
	Status status.AtlasProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// AtlasProjectList contains a list of AtlasProject
type AtlasProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasProject `json:"items"`
}

type ProjectIPAccessList struct {
	// Unique identifier of AWS security group in this access list entry.
	// +optional
	AwsSecurityGroup string `json:"awsSecurityGroup,omitempty"`
	// Range of IP addresses in CIDR notation in this access list entry.
	// +optional
	CIDRBlock string `json:"cidrBlock,omitempty"`
	// Comment associated with this access list entry.
	// +optional
	Comment string `json:"comment,omitempty"`
	// Timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the temporary access list entry.
	// +optional
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`
	// Entry using an IP address in this access list entry.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`
}

func (p AtlasProject) ConnectionSecretObjectKey() *client.ObjectKey {
	if p.Spec.ConnectionSecret != nil {
		key := kube.ObjectKey(p.Namespace, p.Spec.ConnectionSecret.Name)
		return &key
	}
	return nil
}

func (p AtlasProject) GetStatus() interface{} {
	return p.Status
}

func (p *AtlasProject) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	p.Status.Conditions = conditions

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasProjectStatusOption)
		v(&p.Status)
	}
}
