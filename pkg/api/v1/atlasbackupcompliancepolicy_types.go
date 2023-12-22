package v1

import (
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
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
	AuthorizedEmail         string                  `json:"authorizedEmail"`
	AuthorizedUserFirstName string                  `json:"authorizedUserFirstName"`
	AuthorizedUserLastName  string                  `json:"authorizedUserLastName"`
	CopyProtectionEnabled   bool                    `json:"copyProtectionEnabled"`
	EncryptionAtRestEnabled bool                    `json:"encryptionAtRestEnabled"`
	PITEnabled              bool                    `json:"pointInTimeEnabled"`
	RestoreWindowDays       int64                   `json:"restoreWindowDays"`
	ScheduledPolicyItems    []AtlasBackupPolicyItem `json:"scheduledPolicyItems"`
	OnDemandPolicy          AtlasBackupPolicyItem   `json:"onDemandPolicy"`
}

func (b *AtlasBackupCompliancePolicy) ToAtlas() *mongodbatlas.BackupCompliancePolicy {
	result := &mongodbatlas.BackupCompliancePolicy{
		AuthorizedEmail:         b.Spec.AuthorizedEmail,
		CopyProtectionEnabled:   &b.Spec.CopyProtectionEnabled,
		EncryptionAtRestEnabled: &b.Spec.EncryptionAtRestEnabled,
		PitEnabled:              &b.Spec.PITEnabled,
		RestoreWindowDays:       toptr.MakePtr[int64](b.Spec.RestoreWindowDays),
	}

	result.OnDemandPolicyItem = mongodbatlas.PolicyItem{
		FrequencyInterval: b.Spec.OnDemandPolicy.FrequencyInterval,
		FrequencyType:     strings.ToLower(b.Spec.OnDemandPolicy.FrequencyType),
		RetentionValue:    b.Spec.OnDemandPolicy.RetentionValue,
		RetentionUnit:     strings.ToLower(b.Spec.OnDemandPolicy.RetentionUnit),
	}

	for _, policy := range b.Spec.ScheduledPolicyItems {
		result.ScheduledPolicyItems = append(result.ScheduledPolicyItems, mongodbatlas.ScheduledPolicyItem{
			FrequencyInterval: policy.FrequencyInterval,
			FrequencyType:     policy.FrequencyType,
			RetentionUnit:     policy.RetentionUnit,
			RetentionValue:    policy.RetentionValue,
		})
	}

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
