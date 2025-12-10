// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func init() {
	SchemeBuilder.Register(&AtlasBackupCompliancePolicy{}, &AtlasBackupCompliancePolicyList{})
}

// AtlasBackupCompliancePolicy defines the desired state of a compliance policy in Atlas.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=atlas,shortName=abcp
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`

// The AtlasBackupCompliancePolicy is a configuration that enforces specific backup and retention requirements
type AtlasBackupCompliancePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasBackupCompliancePolicySpec     `json:"spec,omitempty"`
	Status status.BackupCompliancePolicyStatus `json:"status,omitempty"`
}

// AtlasBackupCompliancePolicySpec is the specification of the desired configuration of backup compliance policy
type AtlasBackupCompliancePolicySpec struct {
	// Email address of the user who authorized to update the Backup Compliance Policy settings.
	// +kubebuilder:validation:Required
	AuthorizedEmail string `json:"authorizedEmail"`
	// First name of the user who authorized to updated the Backup Compliance Policy settings.
	// +kubebuilder:validation:Required
	AuthorizedUserFirstName string `json:"authorizedUserFirstName"`
	// Last name of the user who authorized to updated the Backup Compliance Policy settings.
	// +kubebuilder:validation:Required
	AuthorizedUserLastName string `json:"authorizedUserLastName"`
	// Flag that indicates whether to prevent cluster users from deleting backups copied to other regions, even if those additional snapshot regions are removed.
	// +kubebuilder:validation:default:=false
	CopyProtectionEnabled bool `json:"copyProtectionEnabled,omitempty"`
	// Flag that indicates whether Encryption at Rest using Customer Key Management is required for all clusters with a Backup Compliance Policy.
	// +kubebuilder:validation:default:=false
	EncryptionAtRestEnabled bool `json:"encryptionAtRestEnabled,omitempty"`
	// Flag that indicates whether to overwrite non-complying backup policies with the new data protection settings or not.
	OverwriteBackupPolicies bool `json:"overwriteBackupPolicies,omitempty"`
	// Flag that indicates whether the cluster uses Continuous Cloud Backups with a Backup Compliance Policy.
	// +kubebuilder:validation:default:=false
	PITEnabled bool `json:"pointInTimeEnabled,omitempty"`
	// Number of previous days that you can restore back to with Continuous Cloud Backup with a Backup Compliance Policy. This parameter applies only to Continuous Cloud Backups with a Backup Compliance Policy.
	RestoreWindowDays int `json:"restoreWindowDays,omitempty"`
	// List that contains the specifications for one scheduled policy.
	ScheduledPolicyItems []AtlasBackupPolicyItem `json:"scheduledPolicyItems,omitempty"`
	// Specifications for on-demand policy.
	OnDemandPolicy AtlasOnDemandPolicy `json:"onDemandPolicy,omitempty"`
}

type AtlasOnDemandPolicy struct {
	// Scope of the backup policy item: days, weeks, or months.
	// +kubebuilder:validation:Enum:=days;weeks;months
	// +kubebuilder:validation:Required
	RetentionUnit string `json:"retentionUnit"`

	// Value to associate with RetentionUnit.
	// +kubebuilder:validation:Required
	RetentionValue int `json:"retentionValue"`
}

func (b *AtlasBackupCompliancePolicy) ToAtlas(projectID string) *admin.DataProtectionSettings20231001 {
	result := &admin.DataProtectionSettings20231001{
		AuthorizedEmail:         b.Spec.AuthorizedEmail,
		AuthorizedUserFirstName: b.Spec.AuthorizedUserFirstName,
		AuthorizedUserLastName:  b.Spec.AuthorizedUserLastName,
		CopyProtectionEnabled:   &b.Spec.CopyProtectionEnabled,
		EncryptionAtRestEnabled: &b.Spec.EncryptionAtRestEnabled,
		PitEnabled:              &b.Spec.PITEnabled,
		ProjectId:               pointer.MakePtr(projectID),
		RestoreWindowDays:       pointer.MakePtr(b.Spec.RestoreWindowDays),
	}

	var emptyPolicy AtlasOnDemandPolicy
	if b.Spec.OnDemandPolicy != emptyPolicy {
		result.OnDemandPolicyItem = &admin.BackupComplianceOnDemandPolicyItem{
			FrequencyInterval: 0,
			FrequencyType:     "ondemand",
			RetentionValue:    b.Spec.OnDemandPolicy.RetentionValue,
			RetentionUnit:     strings.ToLower(b.Spec.OnDemandPolicy.RetentionUnit),
		}
	}

	temp := make([]admin.BackupComplianceScheduledPolicyItem, len(b.Spec.ScheduledPolicyItems))
	for i, policy := range b.Spec.ScheduledPolicyItems {
		temp[i] = admin.BackupComplianceScheduledPolicyItem{
			FrequencyInterval: policy.FrequencyInterval,
			FrequencyType:     policy.FrequencyType,
			RetentionUnit:     policy.RetentionUnit,
			RetentionValue:    policy.RetentionValue,
		}
	}
	result.ScheduledPolicyItems = &temp

	return result
}

func NewBCPFromAtlas(in *admin.DataProtectionSettings20231001) *AtlasBackupCompliancePolicySpec {
	if in == nil {
		return nil
	}

	out := &AtlasBackupCompliancePolicySpec{
		AuthorizedEmail:         in.AuthorizedEmail,
		AuthorizedUserFirstName: in.AuthorizedUserFirstName,
		AuthorizedUserLastName:  in.AuthorizedUserLastName,
		CopyProtectionEnabled:   admin.GetOrDefault(in.CopyProtectionEnabled, false),
		EncryptionAtRestEnabled: admin.GetOrDefault(in.EncryptionAtRestEnabled, false),
		PITEnabled:              admin.GetOrDefault(in.PitEnabled, false),
		RestoreWindowDays:       admin.GetOrDefault(in.RestoreWindowDays, 0),
		OnDemandPolicy: AtlasOnDemandPolicy{
			RetentionUnit:  in.GetOnDemandPolicyItem().RetentionUnit,
			RetentionValue: in.GetOnDemandPolicyItem().RetentionValue,
		},
	}

	temp := make([]AtlasBackupPolicyItem, len(in.GetScheduledPolicyItems()))
	for i, policy := range *in.ScheduledPolicyItems {
		temp[i] = AtlasBackupPolicyItem{
			FrequencyInterval: policy.FrequencyInterval,
			FrequencyType:     policy.FrequencyType,
			RetentionUnit:     policy.RetentionUnit,
			RetentionValue:    policy.RetentionValue,
		}
	}
	out.ScheduledPolicyItems = temp

	return out
}

func (s *AtlasBackupCompliancePolicySpec) Normalize() (*AtlasBackupCompliancePolicySpec, error) {
	err := cmp.Normalize(s)
	if s != nil {
		s.OverwriteBackupPolicies = false
	}
	return s, err
}

func (b *AtlasBackupCompliancePolicy) GetStatus() api.Status {
	return b.Status
}

func (b *AtlasBackupCompliancePolicy) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	b.Status.Conditions = conditions
	b.Status.ObservedGeneration = b.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasBackupCompliancePolicyStatusOption)
		v(&b.Status)
	}
}

// AtlasBackupCompliancePolicyList contains a list of AtlasBackupCompliancePolicy
// +kubebuilder:object:root=true
type AtlasBackupCompliancePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*AtlasBackupCompliancePolicy `json:"items"`
}
