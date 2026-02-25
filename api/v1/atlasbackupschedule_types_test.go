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
	"testing"

	"github.com/go-test/deep"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func Test_BackupScheduleToAtlas(t *testing.T) {
	t.Run("Can convert BackupSchedule to Atlas", func(t *testing.T) {
		inSchedule := &AtlasBackupSchedule{
			Spec: AtlasBackupScheduleSpec{
				AutoExportEnabled:                 true,
				ReferenceHourOfDay:                10,
				ReferenceMinuteOfHour:             10,
				RestoreWindowDays:                 7,
				UpdateSnapshots:                   false,
				UseOrgAndGroupNamesInExportPrefix: false,
			},
		}
		inPolicy := &AtlasBackupPolicy{
			Spec: AtlasBackupPolicySpec{
				Items: []AtlasBackupPolicyItem{
					{
						FrequencyType:     "hourly",
						FrequencyInterval: 10,
						RetentionUnit:     "weeks",
						RetentionValue:    1,
					},
				},
			},
		}
		clusterName := "testCluster"
		replicaSetID := "test-cluster-replica-set-id"
		output := &admin.DiskBackupSnapshotSchedule20240805{
			ClusterId:                         pointer.MakePtr("test-id"),
			ClusterName:                       pointer.MakePtr("testCluster"),
			AutoExportEnabled:                 pointer.MakePtr(true),
			ReferenceHourOfDay:                pointer.MakePtr(10),
			ReferenceMinuteOfHour:             pointer.MakePtr(10),
			RestoreWindowDays:                 pointer.MakePtr(7),
			UpdateSnapshots:                   pointer.MakePtr(false),
			UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(false),
			Policies: &[]admin.AdvancedDiskBackupSnapshotSchedulePolicy{
				{
					PolicyItems: &[]admin.DiskBackupApiPolicyItem{
						{
							FrequencyType:     "hourly",
							FrequencyInterval: 10,
							RetentionUnit:     "weeks",
							RetentionValue:    1,
						},
					},
				},
			},
			CopySettings: &[]admin.DiskBackupCopySetting20240805{},
		}

		result := inSchedule.ToAtlas(output.GetClusterId(), clusterName, replicaSetID, inPolicy)
		if diff := deep.Equal(result, output); diff != nil {
			t.Error(diff)
		}
	})
}
