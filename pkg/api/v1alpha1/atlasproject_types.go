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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	v1status "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasProject{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.experimental.mongodb.com

// AtlasProject is the Schema for the atlasprojects API
type AtlasProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasProjectSpec            `json:"spec,omitempty"`
	Status v1status.AtlasProjectStatus `json:"status,omitempty"`
}

type AtlasProjectSpec struct {
	v1.AtlasProjectSpec `json:",inline"`

	AtlasExperimentalProjectSpec `json:",inline"`
}

type AtlasExperimentalProjectSpec struct {
	// NewField is a new field in the Atlas project and completely breaks compatibility with v1
	// +kubebuilder:validation:Required
	NewField string `json:"newRequiredField"`
}
