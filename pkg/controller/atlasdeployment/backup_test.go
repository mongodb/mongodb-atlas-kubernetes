package atlasdeployment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	atlas_mock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
)

const (
	projectID   = "testProjectID"
	clusterName = "testClusterName"
	clusterID   = "testClusterID"
)

func Test_backupScheduleManagedByAtlas(t *testing.T) {
	deploment := &mdbv1.AtlasDeployment{
		Spec: mdbv1.AtlasDeploymentSpec{
			DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
				Name: clusterName,
			},
		},
	}

	t.Run("should return err when wrong resource passed", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{}, projectID, deploment, &mdbv1.AtlasBackupPolicy{})
		result, err := validator(&mdbv1.AtlasDeployment{})
		assert.EqualError(t, err, errArgIsNotBackupSchedule.Error())
		assert.False(t, result)
	})

	t.Run("should return false if backupschedule is not in atlas", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{
			CloudProviderSnapshotBackupPolicies: &atlas_mock.CloudProviderSnapshotBackupPoliciesClientMock{
				GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
					return nil, &mongodbatlas.Response{}, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ResourceNotFound}
				},
			},
		}, projectID, deploment, &mdbv1.AtlasBackupPolicy{})
		result, err := validator(&mdbv1.AtlasBackupSchedule{})
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return true if resources are not equal", func(t *testing.T) {
		validator := backupScheduleManagedByAtlas(context.TODO(), mongodbatlas.Client{
			CloudProviderSnapshotBackupPolicies: &atlas_mock.CloudProviderSnapshotBackupPoliciesClientMock{
				GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
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
		}, projectID, deploment, &mdbv1.AtlasBackupPolicy{
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
			CloudProviderSnapshotBackupPolicies: &atlas_mock.CloudProviderSnapshotBackupPoliciesClientMock{
				GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
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
		}, projectID, deploment, &mdbv1.AtlasBackupPolicy{
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
						CloudProvider:    toptr.MakePtr[string]("AWS"),
						RegionName:       toptr.MakePtr[string]("us-east-1"),
						ShouldCopyOplogs: toptr.MakePtr(false),
						Frequencies:      []string{},
					},
				},
			},
			Status: status.BackupScheduleStatus{},
		})
		assert.NoError(t, err)
		assert.False(t, result)
	})
}

func Test_backupSchedulesAreEqual(t *testing.T) {
	examplePolicy := mongodbatlas.CloudProviderSnapshotBackupPolicy{
		ClusterID:             "testID",
		ClusterName:           "testName",
		ReferenceHourOfDay:    toptr.MakePtr[int64](12),
		ReferenceMinuteOfHour: toptr.MakePtr[int64](59),
		RestoreWindowDays:     toptr.MakePtr[int64](4),
		UpdateSnapshots:       toptr.MakePtr[bool](false),
		NextSnapshot:          "test123",
		Policies: []mongodbatlas.Policy{
			{
				ID: "testID",
				PolicyItems: []mongodbatlas.PolicyItem{
					{
						ID:                "testID1",
						FrequencyInterval: 10,
						FrequencyType:     "testFreq1",
						RetentionUnit:     "testRet1",
						RetentionValue:    21,
					},
					{
						ID:                "testID2",
						FrequencyInterval: 20,
						FrequencyType:     "testFreq2",
						RetentionUnit:     "testRet2",
						RetentionValue:    450,
					},
				},
			},
		},
		AutoExportEnabled: toptr.MakePtr[bool](true),
		Export: &mongodbatlas.Export{
			ExportBucketID: "testID",
			FrequencyType:  "testFreq",
		},
		UseOrgAndGroupNamesInExportPrefix: toptr.MakePtr[bool](false),
		Links: []*mongodbatlas.Link{
			{
				Rel:  "abc",
				Href: "xyz",
			},
		},
		CopySettings: []mongodbatlas.CopySetting{
			{
				CloudProvider:     toptr.MakePtr[string]("testString"),
				RegionName:        toptr.MakePtr[string]("testString"),
				ReplicationSpecID: toptr.MakePtr[string]("testString"),
				ShouldCopyOplogs:  toptr.MakePtr[bool](true),
				Frequencies:       []string{"testString"},
			},
		},
		DeleteCopiedBackups: []mongodbatlas.DeleteCopiedBackup{
			{
				CloudProvider:     toptr.MakePtr[string]("testString"),
				RegionName:        toptr.MakePtr[string]("testString"),
				ReplicationSpecID: toptr.MakePtr[string]("testString"),
			},
		},
	}

	t.Run("should return true when backups are both empty", func(t *testing.T) {
		res, err := backupSchedulesAreEqual(&mongodbatlas.CloudProviderSnapshotBackupPolicy{}, &mongodbatlas.CloudProviderSnapshotBackupPolicy{})
		assert.NoError(t, err)
		assert.Equal(t, res, bsEqual)
	})
	t.Run("should return true when backups are identical", func(t *testing.T) {
		res, err := backupSchedulesAreEqual(&examplePolicy, &examplePolicy)
		assert.NoError(t, err)
		assert.Equal(t, res, bsEqual)
	})
	t.Run("should return true when backups are identical after normalization", func(t *testing.T) {
		firstPolicy := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
			ClusterID: clusterID,
			Links: []*mongodbatlas.Link{
				{
					Href: "policy1",
					Rel:  "policy1",
				},
			},
			NextSnapshot: "policy1 NextSnapshot",
			Policies: []mongodbatlas.Policy{
				{
					ID: "policy ID",
					PolicyItems: []mongodbatlas.PolicyItem{
						{
							ID:                "policy1 item 1 id",
							FrequencyInterval: 1,
							FrequencyType:     "testFreq1",
							RetentionUnit:     "testRet1",
							RetentionValue:    1,
						},
						{
							ID:                "policy 1 item 2 id",
							FrequencyInterval: 2,
							FrequencyType:     "testFreq2",
							RetentionUnit:     "testRet2",
							RetentionValue:    2,
						},
					},
				},
			},
		}
		secondPolicy := &mongodbatlas.CloudProviderSnapshotBackupPolicy{
			ClusterID: clusterID,
			Links: []*mongodbatlas.Link{
				{
					Href: "policy2",
					Rel:  "policy2",
				},
			},
			NextSnapshot: "policy2 NextSnapshot",
			Policies: []mongodbatlas.Policy{
				{
					ID: "policy ID",
					PolicyItems: []mongodbatlas.PolicyItem{
						{
							ID:                "policy2 item 1 id",
							FrequencyInterval: 1,
							FrequencyType:     "testFreq1",
							RetentionUnit:     "testRet1",
							RetentionValue:    1,
						},
						{
							ID:                "policy 2 item 2 id",
							FrequencyInterval: 2,
							FrequencyType:     "testFreq2",
							RetentionUnit:     "testRet2",
							RetentionValue:    2,
						},
					},
				},
			},
		}
		res, err := backupSchedulesAreEqual(firstPolicy, secondPolicy)
		assert.NoError(t, err)
		assert.Equal(t, res, bsEqual)
	})
	t.Run("should return false when backups differ", func(t *testing.T) {
		changedPolicy := examplePolicy
		changedPolicy.ClusterName = "different name"
		res, err := backupSchedulesAreEqual(&examplePolicy, &changedPolicy)
		assert.NoError(t, err)
		assert.Equal(t, res, bsNotEqual)
	})
}
