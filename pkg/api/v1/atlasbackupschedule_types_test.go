package v1

import (
	"testing"

	"github.com/go-test/deep"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
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
			AutoExportEnabled:                 toptr.MakePtr(true),
			ReferenceHourOfDay:                toptr.MakePtr[int64](10),
			ReferenceMinuteOfHour:             toptr.MakePtr[int64](10),
			RestoreWindowDays:                 toptr.MakePtr[int64](7),
			UpdateSnapshots:                   toptr.MakePtr(false),
			UseOrgAndGroupNamesInExportPrefix: toptr.MakePtr(false),
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
