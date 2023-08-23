package v1

import (
	"testing"

	"github.com/go-test/deep"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func Test_BackupScheduleToAtlas(t *testing.T) {
	testData := []struct {
		name        string
		inSchedule  *AtlasBackupSchedule
		inPolicy    *AtlasBackupPolicy
		clusterName string
		output      *mongodbatlas.CloudProviderSnapshotBackupPolicy
		shouldFail  bool
	}{
		{
			name: "Correct data",
			inSchedule: &AtlasBackupSchedule{
				Spec: AtlasBackupScheduleSpec{
					AutoExportEnabled:                 true,
					ReferenceHourOfDay:                10,
					ReferenceMinuteOfHour:             10,
					RestoreWindowDays:                 7,
					UpdateSnapshots:                   false,
					UseOrgAndGroupNamesInExportPrefix: false,
				},
			},
			inPolicy: &AtlasBackupPolicy{
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
			},
			clusterName: "testCluster",
			output: &mongodbatlas.CloudProviderSnapshotBackupPolicy{
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
			},
			shouldFail: false,
		},
	}

	for _, tt := range testData {
		result := tt.inSchedule.ToAtlas(tt.output.ClusterID, tt.clusterName, tt.inPolicy)
		if diff := deep.Equal(result, tt.output); diff != nil {
			t.Error(diff)
		}
	}
}
