package v1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestBackupCompliancePolicyToAtlas(t *testing.T) {
	t.Run("Can convert Compliance Policy to Atlas", func(t *testing.T) {
		in := AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-bcp",
				Namespace: "test-ns",
				Labels: map[string]string{
					"test": "label",
				},
			},
			Spec: AtlasBackupCompliancePolicySpec{
				AuthorizedEmail:         "example@test.com",
				AuthorizedUserFirstName: "James",
				AuthorizedUserLastName:  "Bond",
				CopyProtectionEnabled:   true,
				EncryptionAtRestEnabled: false,
				PITEnabled:              true,
				RestoreWindowDays:       24,
				ScheduledPolicyItems: []AtlasBackupPolicyItem{
					{
						FrequencyType:     "monthly",
						FrequencyInterval: 2,
						RetentionUnit:     "months",
						RetentionValue:    4,
					},
				},
				OnDemandPolicy: AtlasOnDemandPolicy{
					RetentionUnit:  "days",
					RetentionValue: 14,
				},
			},
		}
		out := in.ToAtlas("testProjectID")

		want := admin.DataProtectionSettings20231001{
			AuthorizedEmail:         "example@test.com",
			AuthorizedUserFirstName: "James",
			AuthorizedUserLastName:  "Bond",
			CopyProtectionEnabled:   pointer.MakePtr(true),
			EncryptionAtRestEnabled: pointer.MakePtr(false),
			PitEnabled:              pointer.MakePtr(true),
			ProjectId:               pointer.MakePtr("testProjectID"),
			RestoreWindowDays:       pointer.MakePtr(24),
			ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
				{
					FrequencyType:     "monthly",
					FrequencyInterval: 2,
					RetentionUnit:     "months",
					RetentionValue:    4,
				},
			},
			OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
				FrequencyInterval: 0,
				FrequencyType:     "ondemand",
				RetentionUnit:     "days",
				RetentionValue:    14,
			},
		}

		assert.Equal(t, "", cmp.Diff(*out, want))
	})
}

func TestBackupCompliancePolicyFromAtlas(t *testing.T) {
	t.Run("Can convert Compliance Policy from Atlas", func(t *testing.T) {
		in := &admin.DataProtectionSettings20231001{
			AuthorizedEmail:         "example@test.com",
			AuthorizedUserFirstName: "James",
			AuthorizedUserLastName:  "Bond",
			CopyProtectionEnabled:   pointer.MakePtr(true),
			EncryptionAtRestEnabled: pointer.MakePtr(false),
			PitEnabled:              pointer.MakePtr(true),
			ProjectId:               pointer.MakePtr("testProjectID"),
			RestoreWindowDays:       pointer.MakePtr(24),
			ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
				{
					FrequencyType:     "monthly",
					FrequencyInterval: 2,
					RetentionUnit:     "months",
					RetentionValue:    4,
				},
			},
			OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
				FrequencyInterval: 0,
				FrequencyType:     "ondemand",
				RetentionUnit:     "days",
				RetentionValue:    14,
			},
		}

		out := NewBCPFromAtlas(in)

		want := AtlasBackupCompliancePolicySpec{
			AuthorizedEmail:         "example@test.com",
			AuthorizedUserFirstName: "James",
			AuthorizedUserLastName:  "Bond",
			CopyProtectionEnabled:   true,
			EncryptionAtRestEnabled: false,
			PITEnabled:              true,
			RestoreWindowDays:       24,
			ScheduledPolicyItems: []AtlasBackupPolicyItem{
				{
					FrequencyType:     "monthly",
					FrequencyInterval: 2,
					RetentionUnit:     "months",
					RetentionValue:    4,
				},
			},
			OnDemandPolicy: AtlasOnDemandPolicy{
				RetentionUnit:  "days",
				RetentionValue: 14,
			},
		}

		assert.Equal(t, "", cmp.Diff(*out, want))
	})
}
