/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package v1

import (
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
)

// AtlasBackupScheduleSpec defines the desired state of AtlasBackupSchedule
type AtlasBackupScheduleSpec struct {
	// Specify true to enable automatic export of cloud backup snapshots to the AWS bucket. You must also define the export policy using export. If omitted, defaults to false.
	// +optional
	// +kubebuilder:default:=false
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

	// Copy backups to other regions for increased resiliency and faster restores.
	// +optional
	CopySettings []CopySetting `json:"copySettings,omitempty"`
}

type AtlasBackupExportSpec struct {
	// Unique Atlas identifier of the AWS bucket which was granted access to export backup snapshot
	ExportBucketID string `json:"exportBucketId"`
	// +kubebuilder:validation:Enum:=monthly
	// +kubebuilder:default:=monthly
	FrequencyType string `json:"frequencyType"`
}

type CopySetting struct {
	// Identifies the cloud provider that stores the snapshot copy.
	// +kubebuilder:validation:Enum:=AWS;GCP;AZURE
	// +kubebuilder:default:=AWS
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Target region to copy snapshots belonging to replicationSpecId to.
	RegionName *string `json:"regionName,omitempty"`
	// Unique identifier that identifies the replication object for a zone in a cluster.
	ReplicationSpecID *string `json:"replicationSpecId,omitempty"`
	// Flag that indicates whether to copy the oplogs to the target region.
	ShouldCopyOplogs *bool `json:"shouldCopyOplogs,omitempty"`
	// List that describes which types of snapshots to copy.
	// +kubebuilder:validation:MinItems=1
	Frequencies []string `json:"frequencies,omitempty"`
}

// AtlasBackupSchedule is the Schema for the atlasbackupschedules API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasBackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AtlasBackupScheduleSpec `json:"spec,omitempty"`

	Status status.BackupScheduleStatus `json:"status,omitempty"`
}

func (in *AtlasBackupSchedule) ToAtlas(clusterID, clusterName string, policy *AtlasBackupPolicy) *mongodbatlas.CloudProviderSnapshotBackupPolicy {
	atlasPolicy := mongodbatlas.Policy{}

	for _, bpItem := range policy.Spec.Items {
		atlasPolicy.PolicyItems = append(atlasPolicy.PolicyItems, mongodbatlas.PolicyItem{
			FrequencyInterval: bpItem.FrequencyInterval,
			FrequencyType:     strings.ToLower(bpItem.FrequencyType),
			RetentionValue:    bpItem.RetentionValue,
			RetentionUnit:     strings.ToLower(bpItem.RetentionUnit),
		})
	}

	result := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
		ClusterName:                       clusterName,
		ReferenceHourOfDay:                &in.Spec.ReferenceHourOfDay,
		ReferenceMinuteOfHour:             &in.Spec.ReferenceMinuteOfHour,
		RestoreWindowDays:                 &in.Spec.RestoreWindowDays,
		UpdateSnapshots:                   &in.Spec.UpdateSnapshots,
		Policies:                          []mongodbatlas.Policy{atlasPolicy},
		AutoExportEnabled:                 &in.Spec.AutoExportEnabled,
		UseOrgAndGroupNamesInExportPrefix: &in.Spec.UseOrgAndGroupNamesInExportPrefix,
		CopySettings:                      make([]mongodbatlas.CopySetting, 0, len(in.Spec.CopySettings)),
	}

	if in.Spec.Export != nil {
		result.Export = &mongodbatlas.Export{
			ExportBucketID: in.Spec.Export.ExportBucketID,
			FrequencyType:  in.Spec.Export.FrequencyType,
		}
	}

	for _, copySetting := range in.Spec.CopySettings {
		result.CopySettings = append(result.CopySettings, mongodbatlas.CopySetting{
			CloudProvider:     copySetting.CloudProvider,
			RegionName:        copySetting.RegionName,
			ReplicationSpecID: copySetting.ReplicationSpecID,
			ShouldCopyOplogs:  copySetting.ShouldCopyOplogs,
			Frequencies:       copySetting.Frequencies,
		})
	}

	result.ClusterID = clusterID
	return result
}

func (in *AtlasBackupSchedule) GetStatus() status.Status {
	return in.Status
}

func (in *AtlasBackupSchedule) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	in.Status.Conditions = conditions
	in.Status.ObservedGeneration = in.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasBackupScheduleStatusOption)
		v(&in.Status)
	}
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
