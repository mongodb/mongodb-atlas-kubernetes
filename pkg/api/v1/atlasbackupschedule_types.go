/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package v1

import (
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
)

// AtlasBackupScheduleSpec defines the desired state of AtlasBackupSchedule
type AtlasBackupScheduleSpec struct {
	// Specify true to enable automatic export of cloud backup snapshots to the AWS bucket. You must also define the export policy using export. If omitted, defaults to false.
	// +optional
	// +kubebuilder:default:=true
	AutoExportEnabled bool `json:"autoExportEnabled,omitempty"`

	// Export policy for automatically exporting cloud backup snapshots to AWS bucket.
	// +optional
	Export *AtlasBackupExportSpec `json:"export,omitempty"`

	// A reference (name & namespace) for backup policy in the desired updated backup policy.
	PolicyRef common.ResourceRefNamespaced `json:"policy"`

	// UTC Hour of day between 0 and 23, inclusive, representing which hour of the day that Atlas takes snapshots for backup policy items
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=23
	// +optional
	ReferenceHourOfDay int64 `json:"referenceHourOfDay,omitempty"`

	// UTC Minutes after ReferenceHourOfDay that Atlas takes snapshots for backup policy items. Must be between 0 and 59, inclusive.
	// +kubebuilder:validation:Minimum:=0
	// +kubebuilder:validation:Maximum:=59
	// +optional
	ReferenceMinuteOfHour int64 `json:"referenceMinuteOfHour,omitempty"`

	// Number of days back in time you can restore to with Continuous Cloud Backup accuracy. Must be a positive, non-zero integer. Applies to continuous cloud backups only.
	// +optional
	// +kubebuilder:default:=1
	RestoreWindowDays int64 `json:"restoreWindowDays,omitempty"`

	// Specify true to apply the retention changes in the updated backup policy to snapshots that Atlas took previously.
	// +optional
	UpdateSnapshots bool `json:"updateSnapshots,omitempty"`

	// Specify true to use organization and project names instead of organization and project UUIDs in the path for the metadata files that Atlas uploads to your S3 bucket after it finishes exporting the snapshots
	// +optional
	UseOrgAndGroupNamesInExportPrefix bool `json:"useOrgAndGroupNamesInExportPrefix,omitempty"`
}

type AtlasBackupExportSpec struct {
	// Unique identifier of the AWS bucket to export the cloud backup snapshot to.
	ExportBucketID string `json:"exportBucketId"`
	// +kubebuilder:validation:Enum:=MONTHLY
	// +kubebuilder:default:=MONTHLY
	FrequencyType string `json:"frequencyType"`
}

// AtlasBackupSchedule is the Schema for the atlasbackupschedules API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasBackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AtlasBackupScheduleSpec `json:"spec,omitempty"`

	Status AtlasBackupScheduleStatus `json:"status,omitempty"`
}

type AtlasBackupScheduleStatus struct {
}

//+kubebuilder:object:root=true

// AtlasBackupScheduleList contains a list of AtlasBackupSchedule
type AtlasBackupScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasBackupSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasBackupSchedule{}, &AtlasBackupScheduleList{})
}

// BackupScheduleFromAtlas converts specs of a backup schedule in native atlas format to AtlasBackupSchedule
func BackupScheduleFromAtlas(policy *mongodbatlas.CloudProviderSnapshotBackupPolicy) (*AtlasBackupScheduleSpec, *AtlasBackupPolicySpec, error) {
	scheduleSpec := &AtlasBackupScheduleSpec{}
	err := compat.JSONCopy(&scheduleSpec, policy)
	if err != nil {
		return nil, nil, err
	}

	policySpec := &AtlasBackupPolicySpec{}

	// Atlas backup schedule doesn't contain any policy
	if len(policy.Policies) < 1 {
		return scheduleSpec, policySpec, nil
	}

	if len(policy.Policies) > 1 {
		return nil, nil, errors.New("more than one policy found in Atlas Backup Schedule")
	}

	// There must be only one policy
	firstPolicy := policy.Policies[0]

	policyItems := make([]AtlasBackupPolicyItem, 0, len(firstPolicy.PolicyItems))

	policySpec.Items = policyItems

	for _, atlasPolicyItem := range firstPolicy.PolicyItems {
		newElem := AtlasBackupPolicyItem{}
		err := compat.JSONCopy(&newElem, atlasPolicyItem)
		if err != nil {
			return nil, nil, err
		}
		policySpec.Items = append(policySpec.Items, newElem)
	}

	return scheduleSpec, policySpec, err
}
