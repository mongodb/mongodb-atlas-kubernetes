package v1

import (
	"testing"

	"github.com/go-test/deep"
	"go.mongodb.org/atlas/mongodbatlas"

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
		output := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
			ClusterID:                         "test-id",
			ClusterName:                       "testCluster",
			AutoExportEnabled:                 pointer.MakePtr(true),
			ReferenceHourOfDay:                pointer.MakePtr[int64](10),
			ReferenceMinuteOfHour:             pointer.MakePtr[int64](10),
			RestoreWindowDays:                 pointer.MakePtr[int64](7),
			UpdateSnapshots:                   pointer.MakePtr(false),
			UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(false),
			Policies: []mongodbatlas.Policy{
				{
					ID: "",
					PolicyItems: []mongodbatlas.PolicyItem{
						{
							ID:                "",
							FrequencyType:     "hourly",
							FrequencyInterval: 10,
							RetentionUnit:     "weeks",
							RetentionValue:    1,
						},
					},
				},
			},
			CopySettings: []mongodbatlas.CopySetting{},
		}

		result := inSchedule.ToAtlas(output.ClusterID, clusterName, replicaSetID, inPolicy)
		if diff := deep.Equal(result, output); diff != nil {
			t.Error(diff)
		}
	})
}
