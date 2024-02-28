package v1

import (
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

// AtlasBackupCompliancePolicy defines the desired state of a compliance policy in Atlas.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasBackupCompliancePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasBackupCompliancePolicySpec     `json:"spec,omitempty"`
	Status status.BackupCompliancePolicyStatus `json:"status,omitempty"`
}

type AtlasBackupCompliancePolicySpec struct {
	AuthorizedEmail         string `json:"authorizedEmail"`
	AuthorizedUserFirstName string `json:"authorizedUserFirstName"`
	AuthorizedUserLastName  string `json:"authorizedUserLastName"`
	// +kubebuilder:validation:default:=false
	CopyProtectionEnabled bool `json:"copyProtectionEnabled"`
	// +kubebuilder:validation:default:=false
	EncryptionAtRestEnabled bool `json:"encryptionAtRestEnabled"`
	Enforce                 bool `json:"enforce"`
	// +kubebuilder:validation:default:=false
	PITEnabled           bool                    `json:"pointInTimeEnabled"`
	RestoreWindowDays    int                     `json:"restoreWindowDays"`
	ScheduledPolicyItems []AtlasBackupPolicyItem `json:"scheduledPolicyItems"`
	OnDemandPolicy       AtlasOnDemandPolicy     `json:"onDemandPolicy"`
}

type AtlasOnDemandPolicy struct {
	// Scope of the backup policy item: days, weeks, or months
	// +kubebuilder:validation:Enum:=days;weeks;months
	RetentionUnit string `json:"retentionUnit"`

	// Value to associate with RetentionUnit
	RetentionValue int `json:"retentionValue"`
}

func (b *AtlasBackupCompliancePolicy) ToAtlas(projectID string) *admin.DataProtectionSettings20231001 {
	// TODO: add enforce flag once present in the API
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

	result.OnDemandPolicyItem = &admin.BackupComplianceOnDemandPolicyItem{
		FrequencyInterval: 0,
		FrequencyType:     "ondemand",
		RetentionValue:    b.Spec.OnDemandPolicy.RetentionValue,
		RetentionUnit:     strings.ToLower(b.Spec.OnDemandPolicy.RetentionUnit),
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

func (b *AtlasBackupCompliancePolicy) GetStatus() status.Status {
	return b.Status
}

func (b *AtlasBackupCompliancePolicy) UpdateStatus(conditions []status.Condition, options ...status.Option) {
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

func init() {
	SchemeBuilder.Register(&AtlasBackupCompliancePolicy{}, &AtlasBackupCompliancePolicyList{})
}
