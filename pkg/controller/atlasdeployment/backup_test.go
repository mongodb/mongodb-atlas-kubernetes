package atlasdeployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

type backupsClient struct {
	GetFunc func() (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error)
}

func (bs *backupsClient) Get(_ context.Context, _ string, _ string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	return bs.GetFunc()
}

func (bs *backupsClient) Update(_ context.Context, _ string, _ string, _ *mongodbatlas.CloudProviderSnapshotBackupPolicy) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (bs *backupsClient) Delete(_ context.Context, _ string, _ string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

const (
	projectID   = "testProjectID"
	clusterName = "testClusterName"
	clusterID   = "testClusterID"
)

func Test_backupScheduleManagedByAtlas(t *testing.T) {
	t.Run("should return err when wrong resource passed", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{}, projectID, clusterName, &mdbv1.AtlasBackupPolicy{})
		result, err := validator(&mdbv1.AtlasDeployment{})
		assert.EqualError(t, err, errArgIsNotBackupSchedule.Error())
		assert.False(t, result)
	})

	t.Run("should return false if backupschedule is not in atlas", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{
			CloudProviderSnapshotBackupPolicies: &backupsClient{
				GetFunc: func() (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound}
				},
			},
		}, projectID, clusterName, &mdbv1.AtlasBackupPolicy{})
		result, err := validator(&mdbv1.AtlasBackupSchedule{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true if resources are not equal", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{
			CloudProviderSnapshotBackupPolicies: &backupsClient{
				GetFunc: func() (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderSnapshotBackupPolicy{
							ClusterID:             clusterID,
							ClusterName:           clusterName,
							ReferenceHourOfDay:    new(int64),
							ReferenceMinuteOfHour: new(int64),
							RestoreWindowDays:     new(int64),
							UpdateSnapshots:       new(bool),
							NextSnapshot:          "",
							Policies: []mongodbatlas.Policy{
								{
									ID: "policy-id",
									PolicyItems: []mongodbatlas.PolicyItem{
										{
											ID:                "policy-item-id",
											FrequencyInterval: 10,
											FrequencyType:     "hours",
											RetentionUnit:     "days",
											RetentionValue:    10,
										},
									},
								},
							},
							AutoExportEnabled:                 toptr.MakePtr(false),
							Export:                            &mongodbatlas.Export{},
							UseOrgAndGroupNamesInExportPrefix: toptr.MakePtr(false),
							Links:                             []*mongodbatlas.Link{},
							CopySettings: []mongodbatlas.CopySetting{
								{
									CloudProvider:     toptr.MakePtr[string]("AWS"),
									RegionName:        toptr.MakePtr[string]("us-east-1"),
									ReplicationSpecID: toptr.MakePtr[string]("test-id"),
									ShouldCopyOplogs:  new(bool),
									Frequencies:       []string{},
								},
							},
							DeleteCopiedBackups: []mongodbatlas.DeleteCopiedBackup{},
						},
						&mongodbatlas.Response{}, nil
				},
			},
		}, projectID, clusterName, &mdbv1.AtlasBackupPolicy{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       mdbv1.AtlasBackupPolicySpec{},
			Status:     status.BackupPolicyStatus{},
		})
		result, err := validator(&mdbv1.AtlasBackupSchedule{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if resources are equal", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{
			CloudProviderSnapshotBackupPolicies: &backupsClient{
				GetFunc: func() (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
					return &mongodbatlas.CloudProviderSnapshotBackupPolicy{
							ClusterID:             clusterID,
							ClusterName:           clusterName,
							ReferenceHourOfDay:    toptr.MakePtr[int64](10),
							ReferenceMinuteOfHour: toptr.MakePtr[int64](10),
							RestoreWindowDays:     toptr.MakePtr[int64](10),
							UpdateSnapshots:       toptr.MakePtr[bool](false),
							NextSnapshot:          "",
							Policies: []mongodbatlas.Policy{
								{
									ID: "policy-id",
									PolicyItems: []mongodbatlas.PolicyItem{
										{
											ID:                "policy-item-id",
											FrequencyInterval: 10,
											FrequencyType:     "hours",
											RetentionUnit:     "days",
											RetentionValue:    10,
										},
									},
								},
							},
							AutoExportEnabled:                 toptr.MakePtr(false),
							Export:                            &mongodbatlas.Export{},
							UseOrgAndGroupNamesInExportPrefix: toptr.MakePtr(false),
							Links:                             []*mongodbatlas.Link{},
							CopySettings: []mongodbatlas.CopySetting{
								{
									CloudProvider:     toptr.MakePtr[string]("AWS"),
									RegionName:        toptr.MakePtr[string]("us-east-1"),
									ReplicationSpecID: toptr.MakePtr[string]("test-id"),
									ShouldCopyOplogs:  toptr.MakePtr(false),
									Frequencies:       []string{},
								},
							},
							DeleteCopiedBackups: []mongodbatlas.DeleteCopiedBackup{},
						},
						&mongodbatlas.Response{}, nil
				},
			},
		}, projectID, clusterName, &mdbv1.AtlasBackupPolicy{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: mdbv1.AtlasBackupPolicySpec{
				Items: []mdbv1.AtlasBackupPolicyItem{
					{
						FrequencyType:     "hours",
						FrequencyInterval: 10,
						RetentionUnit:     "days",
						RetentionValue:    10,
					},
				},
			},
		})
		result, err := validator(&mdbv1.AtlasBackupSchedule{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: mdbv1.AtlasBackupScheduleSpec{
				AutoExportEnabled:                 false,
				Export:                            &mdbv1.AtlasBackupExportSpec{},
				PolicyRef:                         common.ResourceRefNamespaced{},
				ReferenceHourOfDay:                10,
				ReferenceMinuteOfHour:             10,
				RestoreWindowDays:                 10,
				UpdateSnapshots:                   false,
				UseOrgAndGroupNamesInExportPrefix: false,
				CopySettings: []mdbv1.CopySetting{
					{
						CloudProvider:     toptr.MakePtr[string]("AWS"),
						RegionName:        toptr.MakePtr[string]("us-east-1"),
						ReplicationSpecID: toptr.MakePtr[string]("test-id"),
						ShouldCopyOplogs:  toptr.MakePtr(false),
						Frequencies:       []string{},
					},
				},
			},
			Status: status.BackupScheduleStatus{},
		})
		assert.NoError(t, err)
		assert.False(t, result)
	})
}
