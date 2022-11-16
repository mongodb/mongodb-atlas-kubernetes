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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"

	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status

// AtlasTeam is the Schema for the Atlas Teams API
type AtlasTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TeamSpec          `json:"spec"`
	Status status.TeamStatus `json:"status,omitempty"`
}

// +kubebuilder:validation:Format=email

type TeamUser string

// TeamSpec defines the desired state of a Team in Atlas
type TeamSpec struct {
	// The name of the team you want to create.
	Name string `json:"name"`
	// Valid email addresses of users to add to the new team
	Usernames []TeamUser `json:"usernames"`
}

// +kubebuilder:object:root=true

// AtlasTeamList contains a list of AtlasTeam
type AtlasTeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasTeam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasTeam{}, &AtlasTeamList{})
}

func (in *AtlasTeam) Identifier() interface{} {
	return in.Status.ID
}

func (in *AtlasTeam) GetStatus() status.Status {
	return in.Status
}

func (in *AtlasTeam) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	in.Status.Conditions = conditions
	in.Status.ObservedGeneration = in.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasTeamStatusOption)
		v(&in.Status)
	}
}

func (in *AtlasTeam) ToAtlas() (*mongodbatlas.Team, error) {
	result := &mongodbatlas.Team{}
	err := compat.JSONCopy(result, in.Spec)

	return result, err
}
