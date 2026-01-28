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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
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

		assert.True(t, reflect.DeepEqual(*out, want), cmp.Diff(*out, want))
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

		assert.True(t, reflect.DeepEqual(*out, want), cmp.Diff(*out, want))
	})
}

func TestBackupCompliancePolicyFromAtlasNilOndemandPolicy(t *testing.T) {
	t.Run("Can convert from Atlas when OndemandPolicyItem is nil", func(t *testing.T) {
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
			OnDemandPolicyItem: nil,
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
			OnDemandPolicy: AtlasOnDemandPolicy{},
		}

		assert.True(t, reflect.DeepEqual(*out, want), cmp.Diff(*out, want))
	})
}

func TestBackupCompliancePolicyToAtlasNilOndemandPolicy(t *testing.T) {
	t.Run("Can convert to Atlas when OndemandPolicyItem is nil", func(t *testing.T) {
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
			OnDemandPolicyItem: nil,
		}

		assert.True(t, reflect.DeepEqual(*out, want), cmp.Diff(*out, want))
	})
}
