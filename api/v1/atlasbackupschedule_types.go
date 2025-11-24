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
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// AtlasBackupScheduleSpec defines the desired state of AtlasBackupSchedule.
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
	// Unique Atlas identifier of the AWS bucket which was granted access to export backup snapshot.
	ExportBucketID string `json:"exportBucketId"`
	// Human-readable label that indicates the rate at which the export policy item occurs.
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
	// Flag that indicates whether to copy the oplogs to the target region.
	ShouldCopyOplogs *bool `json:"shouldCopyOplogs,omitempty"`
	// List that describes which types of snapshots to copy.
	// +kubebuilder:validation:MinItems=1
	Frequencies []string `json:"frequencies,omitempty"`
}

var _ api.AtlasCustomResource = &AtlasBackupSchedule{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=atlas,shortName=abs
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
//
// AtlasBackupSchedule is the Schema for the atlasbackupschedules API.
type AtlasBackupSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AtlasBackupScheduleSpec `json:"spec,omitempty"`

	Status status.BackupScheduleStatus `json:"status,omitempty"`
}

func (in *AtlasBackupSchedule) ToAtlas(clusterID, clusterName, zoneID string, policy *AtlasBackupPolicy) *admin.DiskBackupSnapshotSchedule20240805 {
	atlasPolicy := admin.AdvancedDiskBackupSnapshotSchedulePolicy{}

	items := make([]admin.DiskBackupApiPolicyItem, 0, len(policy.Spec.Items))
	for _, bpItem := range policy.Spec.Items {
		items = append(items, admin.DiskBackupApiPolicyItem{
			FrequencyInterval: bpItem.FrequencyInterval,
			FrequencyType:     strings.ToLower(bpItem.FrequencyType),
			RetentionUnit:     strings.ToLower(bpItem.RetentionUnit),
			RetentionValue:    bpItem.RetentionValue,
		})
	}
	atlasPolicy.PolicyItems = &items

	result := &admin.DiskBackupSnapshotSchedule20240805{
		ClusterName:                       pointer.MakePtrOrNil(clusterName),
		ClusterId:                         pointer.MakePtrOrNil(clusterID),
		ReferenceHourOfDay:                pointer.MakePtr(int(in.Spec.ReferenceHourOfDay)),
		ReferenceMinuteOfHour:             pointer.MakePtr(int(in.Spec.ReferenceMinuteOfHour)),
		RestoreWindowDays:                 pointer.MakePtr(int(in.Spec.RestoreWindowDays)),
		UpdateSnapshots:                   pointer.MakePtr(in.Spec.UpdateSnapshots),
		Policies:                          &[]admin.AdvancedDiskBackupSnapshotSchedulePolicy{atlasPolicy},
		AutoExportEnabled:                 pointer.MakePtr(in.Spec.AutoExportEnabled),
		UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(in.Spec.UseOrgAndGroupNamesInExportPrefix),
	}

	if in.Spec.Export != nil {
		result.Export = &admin.AutoExportPolicy{
			ExportBucketId: pointer.MakePtr(in.Spec.Export.ExportBucketID),
			FrequencyType:  pointer.MakePtr(in.Spec.Export.FrequencyType),
		}
	}

	copySettings := make([]admin.DiskBackupCopySetting20240805, 0, len(in.Spec.CopySettings))
	for _, copySetting := range in.Spec.CopySettings {
		copySettings = append(copySettings, admin.DiskBackupCopySetting20240805{
			CloudProvider:    copySetting.CloudProvider,
			RegionName:       copySetting.RegionName,
			ZoneId:           zoneID,
			ShouldCopyOplogs: copySetting.ShouldCopyOplogs,
			Frequencies:      &copySetting.Frequencies,
		})
	}
	result.CopySettings = &copySettings

	return result
}

func (in *AtlasBackupSchedule) GetStatus() api.Status {
	return in.Status
}

func (in *AtlasBackupSchedule) UpdateStatus(conditions []api.Condition, options ...api.Option) {
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
