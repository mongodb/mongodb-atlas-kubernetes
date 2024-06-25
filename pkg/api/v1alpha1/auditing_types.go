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

package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	status "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1/status"
)

// TODO: uncomment when publishing API
// // Important:
// // The procedure working with this file:
// // 1. Edit the file
// // 1. Run "make generate" to regenerate code
// // 2. Run "make manifests" to regenerate the CRD
func init() {
	SchemeBuilder.Register(&AtlasAuditing{}, &AtlasAuditingList{})
}

// +k8s:deepcopy-gen=package

type AuditingSpecTypes string

const (
	// Standalone operation mode for the Auditing Config
	Standalone AuditingSpecTypes = "standalone"

	// Linked operation mode for the Auditing Config
	Linked AuditingSpecTypes = "linked"
)

func (auditType AuditingSpecTypes) Valid() bool {
	return auditType == Standalone || auditType == Linked
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AtlasAuditing is the Schema for the Atlas Auditing API
// +k8s:deepcopy-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasAuditing struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasAuditingSpec          `json:"spec,omitempty"`
	Status status.AtlasAuditingStatus `json:"status,omitempty"`
}

func (in *AtlasAuditing) GetStatus() api.Status {
	return in.Status
}

func (in *AtlasAuditing) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	in.Status.Conditions = conditions
	in.Status.ObservedGeneration = in.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasAuditingStatusOption)
		v(&in.Status)
	}
}

// AtlasAuditingSpec defines the desired state of Database Auditing in Atlas
// +k8s:deepcopy-gen=true
type AtlasAuditingSpec struct {
	AtlasAuditingConfig `json:",inline"`

	// Type of the Auditing config definition
	// +kubebuilder:default:=standalone
	// +kubebuilder:validation:Enum:=standalone;linked
	Type AuditingSpecTypes `json:"type"`

	// ProjectIDs is a list of projects using this auditing config
	// This can NOT be used when type is "linked"
	ProjectIDs []string `json:"projectIDs,omitempty"`
}

// AtlasAuditingConfig represents the actual fields that can bet set by Auditing
type AtlasAuditingConfig struct {
	// Enabled is true when database auditing is on for the given projects
	Enabled bool `json:"enabled"`

	// AuditAuthorizationSuccess is true when auth successes are to be logged
	AuditAuthorizationSuccess bool `json:"auditAuthorizationSuccess,omitempty"`

	// AuditFilter contains the JSON/YAML definition of the audit logging filter
	AuditFilter *apiextensionsv1.JSON `json:"auditFilter,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasAuditingList contains a list of AtlasAuditing
type AtlasAuditingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasAuditing `json:"items"`
}
