//Copyright 2022 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

// AtlasBackupPolicySpec defines the desired state of AtlasBackupPolicy
type AtlasBackupPolicySpec struct {
	// A list of BackupPolicy items.
	Items []AtlasBackupPolicyItem `json:"items"`
}

type AtlasBackupPolicyItem struct {
	// Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly. You cannot specify multiple hourly and daily backup policy items.
	// +kubebuilder:validation:Enum:=hourly;daily;weekly;monthly;yearly
	FrequencyType string `json:"frequencyType"`

	// Desired frequency of the new backup policy item specified by FrequencyType. A value of 1 specifies the first instance of the corresponding FrequencyType.
	// The only accepted value you can set for frequency interval with NVMe clusters is 12.
	// +kubebuilder:validation:Enum:=1;2;3;4;5;6;7;8;9;10;11;12;13;14;15;16;17;18;19;20;21;22;23;24;25;26;27;28;40
	FrequencyInterval int `json:"frequencyInterval"`

	// Scope of the backup policy item: days, weeks, or months.
	// +kubebuilder:validation:Enum:=days;weeks;months;years
	RetentionUnit string `json:"retentionUnit"`

	// Value to associate with RetentionUnit.
	RetentionValue int `json:"retentionValue"`
}

var _ api.AtlasCustomResource = &AtlasBackupPolicy{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=atlas,shortName=abp
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//
// AtlasBackupPolicy is the Schema for the atlasbackuppolicies API
type AtlasBackupPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasBackupPolicySpec     `json:"spec,omitempty"`
	Status status.BackupPolicyStatus `json:"status,omitempty"`
}

func (in *AtlasBackupPolicy) GetStatus() api.Status {
	return in.Status
}

func (in *AtlasBackupPolicy) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	in.Status.Conditions = conditions
	in.Status.ObservedGeneration = in.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasBackupPolicyStatusOption)
		v(&in.Status)
	}
}

//+kubebuilder:object:root=true

// AtlasBackupPolicyList contains a list of AtlasBackupPolicy
type AtlasBackupPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasBackupPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasBackupPolicy{}, &AtlasBackupPolicyList{})
}
