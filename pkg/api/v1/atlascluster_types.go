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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AtlasClusterSpec defines the desired state of AtlasCluster
type AtlasClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AtlasCluster. Edit AtlasCluster_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// AtlasClusterStatus defines the observed state of AtlasCluster
type AtlasClusterStatus struct {
	status.Common `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AtlasCluster is the Schema for the atlasclusters API
type AtlasCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasClusterSpec   `json:"spec,omitempty"`
	Status AtlasClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasClusterList contains a list of AtlasCluster
type AtlasClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasCluster{}, &AtlasClusterList{})
}

func (c AtlasCluster) GetStatus() interface{} {
	return c.Status
}
