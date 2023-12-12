package v1

import (
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AtlasBackupCompliancePolicy defines the desired state of a compliance policy in Atlas.
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

func (b *AtlasBackupCompliancePolicySpec) ToAtlas() *mongodbatlas.BackupCompliancePolicy {
	result := &mongodbatlas.BackupCompliancePolicy{
		AuthorizedEmail:         b.AuthorizedEmail,
		CopyProtectionEnabled:   &b.CopyProtectionEnabled,
		EncryptionAtRestEnabled: &b.EncryptionAtRestEnabled,
		PitEnabled:              &b.PITEnabled,
		RestoreWindowDays:       toptr.MakePtr[int64](b.RestoreWindowDays),
	}

	result.OnDemandPolicyItem = mongodbatlas.PolicyItem{
		FrequencyInterval: b.OnDemandPolicy.FrequencyInterval,
		FrequencyType:     strings.ToLower(b.OnDemandPolicy.FrequencyType),
		RetentionValue:    b.OnDemandPolicy.RetentionValue,
		RetentionUnit:     strings.ToLower(b.OnDemandPolicy.RetentionUnit),
	}

	for _, policy := range b.ScheduledPolicyItems {
		result.ScheduledPolicyItems = append(result.ScheduledPolicyItems, mongodbatlas.ScheduledPolicyItem{
			FrequencyInterval: policy.FrequencyInterval,
			FrequencyType:     policy.FrequencyType,
			RetentionUnit:     policy.RetentionUnit,
			RetentionValue:    policy.RetentionValue,
		})
	}

	return result
}
