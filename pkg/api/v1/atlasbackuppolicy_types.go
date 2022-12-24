/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
)

// AtlasBackupPolicySpec defines the desired state of AtlasBackupPolicy
type AtlasBackupPolicySpec struct {
	// A list of BackupPolicy items
	Items []AtlasBackupPolicyItem `json:"items"`
}

type AtlasBackupPolicyItem struct {
	// Frequency associated with the backup policy item. One of the following values: hourly, daily, weekly or monthly. You cannot specify multiple hourly and daily backup policy items.
	// +kubebuilder:validation:Enum:=hourly;daily;weekly;monthly
	FrequencyType string `json:"frequencyType"`

	// Desired frequency of the new backup policy item specified by FrequencyType. A value of 1 specifies the first instance of the corresponding FrequencyType.
	// The only accepted value you can set for frequency interval with NVMe clusters is 12.
	// +kubebuilder:validation:Enum:=1;2;3;4;5;6;7;8;9;10;11;12;13;14;15;16;17;18;19;20;21;22;23;24;25;26;27;28;40
	FrequencyInterval int `json:"frequencyInterval"`

	// Scope of the backup policy item: days, weeks, or months
	// +kubebuilder:validation:Enum:=days;weeks;months
	RetentionUnit string `json:"retentionUnit"`

	// Value to associate with RetentionUnit
	RetentionValue int `json:"retentionValue"`
}

// AtlasBackupPolicy is the Schema for the atlasbackuppolicies API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasBackupPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasBackupPolicySpec   `json:"spec,omitempty"`
	Status AtlasBackupPolicyStatus `json:"status,omitempty"`
}

func (in *AtlasBackupPolicy) GetStatus() status.Status {
	return nil
}

func (in *AtlasBackupPolicy) UpdateStatus(_ []status.Condition, _ ...status.Option) {
}

type AtlasBackupPolicyStatus struct {
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
