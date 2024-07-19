package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestComputeChanges(t *testing.T) {
	tests := map[string]struct {
		akoCluster      *Cluster
		atlasCluster    *Cluster
		expectedChanges *Cluster
		changed         bool
	}{
		"should handle pause change as special case": {
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(false),
					Paused:        pointer.MakePtr(true),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					Paused:        pointer.MakePtr(false),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:   "cluster0",
					Paused: pointer.MakePtr(true),
				},
			},
			changed: true,
		},
		"should not update disk size when unset": {
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(false),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					DiskSizeGB:    pointer.MakePtr(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    pointer.MakePtr(false),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{},
				},
			},
			changed: true,
		},
		"should not update disk size when they are the same": {
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(false),
					DiskSizeGB:    pointer.MakePtr(20),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					DiskSizeGB:    pointer.MakePtr(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    pointer.MakePtr(false),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{},
				},
			},
			changed: true,
		},
		"should update disk size when they are different": {
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(false),
					DiskSizeGB:    pointer.MakePtr(30),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					DiskSizeGB:    pointer.MakePtr(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    pointer.MakePtr(false),
					DiskSizeGB:       pointer.MakePtr(30),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{},
				},
			},
			changed: true,
		},
		"should update all spec when there are changes": {
			//nolint:dupl
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					DiskSizeGB:               pointer.MakePtr(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               pointer.MakePtr(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        pointer.MakePtr(true),
						ReadPreference: "secondary",
					},
					Labels: []common.LabelSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					VersionReleaseSystem:         "LTS",
					TerminationProtectionEnabled: true,
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M40",
										},
									},
								},
							},
						},
					},
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					DiskSizeGB:    pointer.MakePtr(20),
				},
			},
			//nolint:dupl
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					DiskSizeGB:               pointer.MakePtr(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               pointer.MakePtr(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        pointer.MakePtr(true),
						ReadPreference: "secondary",
					},
					Labels: []common.LabelSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					VersionReleaseSystem:         "LTS",
					TerminationProtectionEnabled: true,
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M40",
										},
									},
								},
							},
						},
					},
				},
			},
			changed: true,
		},
		"should return nil when there are no changes": {
			//nolint:dupl
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					DiskSizeGB:               pointer.MakePtr(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "",
					Tags: []*akov2.TagSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					VersionReleaseSystem:         "LTS",
					TerminationProtectionEnabled: false,
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M40",
										},
									},
								},
							},
						},
					},
				},
			},
			//nolint:dupl
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					DiskSizeGB:               pointer.MakePtr(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "",
					Tags: []*akov2.TagSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					VersionReleaseSystem:         "LTS",
					TerminationProtectionEnabled: false,
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M40",
										},
									},
								},
							},
						},
					},
				},
			},
			expectedChanges: nil,
			changed:         false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			changes, changed := ComputeChanges(tt.akoCluster, tt.atlasCluster)
			assert.Equal(t, tt.changed, changed)
			assert.Equal(t, tt.expectedChanges, changes)
		})
	}
}

func TestAreEqual(t *testing.T) {
	t.Run("should compare booleans", func(t *testing.T) {
		tests := map[string]struct {
			desired  *bool
			current  *bool
			expected bool
		}{
			"both are nil": {
				expected: true,
			},
			"desired is nil and current is pointer to zero value": {
				current:  pointer.MakePtr(false),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  pointer.MakePtr(false),
				current:  pointer.MakePtr(false),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  pointer.MakePtr(false),
				expected: true,
			},
			"desired and current are true": {
				desired:  pointer.MakePtr(true),
				current:  pointer.MakePtr(true),
				expected: true,
			},
			"desired is nil and current is true": {
				current:  pointer.MakePtr(true),
				expected: false,
			},
			"desired is true and current is nil": {
				desired:  pointer.MakePtr(true),
				expected: false,
			},
			"desired is false and current is true": {
				desired:  pointer.MakePtr(false),
				current:  pointer.MakePtr(true),
				expected: false,
			},
			"desired is true and current is false": {
				desired:  pointer.MakePtr(true),
				current:  pointer.MakePtr(false),
				expected: false,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, tt.expected, areEqual(tt.desired, tt.current))
			})
		}
	})

	// nolint:dupl
	t.Run("should compare strings", func(t *testing.T) {
		tests := map[string]struct {
			desired  *string
			current  *string
			expected bool
		}{
			"both are nil": {
				expected: true,
			},
			"desired is nil and current is pointer to zero value": {
				current:  pointer.MakePtr(""),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  pointer.MakePtr(""),
				current:  pointer.MakePtr(""),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  pointer.MakePtr(""),
				expected: true,
			},
			"desired and current have value": {
				desired:  pointer.MakePtr("value"),
				current:  pointer.MakePtr("value"),
				expected: true,
			},
			"desired is nil and current has value": {
				current:  pointer.MakePtr("value"),
				expected: false,
			},
			"desired has value and current is nil": {
				desired:  pointer.MakePtr("value"),
				expected: false,
			},
			"desired has different value from current": {
				desired:  pointer.MakePtr("value"),
				current:  pointer.MakePtr("other-value"),
				expected: false,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, tt.expected, areEqual(tt.desired, tt.current))
			})
		}
	})

	// nolint:dupl
	t.Run("should compare integers", func(t *testing.T) {
		tests := map[string]struct {
			desired  *int
			current  *int
			expected bool
		}{
			"both are nil": {
				expected: true,
			},
			"desired is nil and current is pointer to zero value": {
				current:  pointer.MakePtr(0),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  pointer.MakePtr(0),
				current:  pointer.MakePtr(0),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  pointer.MakePtr(0),
				expected: true,
			},
			"desired and current have value": {
				desired:  pointer.MakePtr(10),
				current:  pointer.MakePtr(10),
				expected: true,
			},
			"desired is nil and current has value": {
				current:  pointer.MakePtr(10),
				expected: false,
			},
			"desired has value and current is nil": {
				desired:  pointer.MakePtr(10),
				expected: false,
			},
			"desired has different value from current": {
				desired:  pointer.MakePtr(10),
				current:  pointer.MakePtr(11),
				expected: false,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, tt.expected, areEqual(tt.desired, tt.current))
			})
		}
	})
}

func TestSpecAreEqual(t *testing.T) {
	tests := map[string]struct {
		ako      *akov2.AtlasDeployment
		atlas    *admin.AdvancedClusterDescription
		expected bool
	}{
		"should return false when cluster type are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ClusterType: "SHARDED",
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				ClusterType: pointer.MakePtr("REPLICASET"),
			},
		},
		"should return false when backup enabled flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						BackupEnabled: pointer.MakePtr(true),
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				BackupEnabled: pointer.MakePtr(false),
			},
		},
		"should return false when BI connector config are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						BiConnector: &akov2.BiConnectorSpec{
							Enabled:        pointer.MakePtr(true),
							ReadPreference: "secondary",
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				BiConnector: &admin.BiConnector{
					Enabled:        pointer.MakePtr(false),
					ReadPreference: pointer.MakePtr("secondary"),
				},
			},
		},
		"should return false when disk size are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						DiskSizeGB: pointer.MakePtr(20),
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				DiskSizeGB: pointer.MakePtr(10.0),
			},
		},
		"should return false when encryption at rest config are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						EncryptionAtRestProvider: "AWS",
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
			},
		},
		"should return false when mongodb version are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						MongoDBMajorVersion: "8.0",
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				MongoDBMajorVersion: pointer.MakePtr("7.0"),
			},
		},
		"should return false when version release system are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						VersionReleaseSystem: "CONTINUOUS",
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				VersionReleaseSystem: pointer.MakePtr("LTS"),
			},
		},
		"should return false when root cert type are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						RootCertType: "ISRGROOTX1",
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				RootCertType: pointer.MakePtr("NONE"),
			},
		},
		"should return false when paused flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Paused: pointer.MakePtr(true),
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Paused: pointer.MakePtr(false),
			},
		},
		"should return false when pit flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						PitEnabled: pointer.MakePtr(true),
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				PitEnabled: pointer.MakePtr(false),
			},
		},
		"should return false when termination protection flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						TerminationProtectionEnabled: true,
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				TerminationProtectionEnabled: pointer.MakePtr(false),
			},
		},
		"should return false when num of shards are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 3,
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
					},
				},
			},
		},
		"should return false when region are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
						RegionConfigs: &[]admin.CloudRegionConfig{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
							},
						},
					},
				},
			},
		},
		"should return false when autoscaling are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M40",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
						RegionConfigs: &[]admin.CloudRegionConfig{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled: pointer.MakePtr(false),
									},
								},
							},
						},
					},
				},
			},
		},
		"should return false when labels are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Labels: &[]admin.ComponentLabel{
					{
						Key:   pointer.MakePtr("label2"),
						Value: pointer.MakePtr("label2"),
					},
				},
			},
		},
		"should return false when tags are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Tags: []*akov2.TagSpec{
							{
								Key:   "tag1",
								Value: "tag1",
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag2",
						Value: "tag2",
					},
				},
			},
		},
		"should return true when cluster are the same": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "REPLICASET",
						BackupEnabled: pointer.MakePtr(true),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          pointer.MakePtr(true),
						RootCertType:        "ISRGROOTX1",
						Tags: []*akov2.TagSpec{
							{
								Key:   "tag1",
								Value: "tag1",
							},
						},
						VersionReleaseSystem:         "LTS",
						TerminationProtectionEnabled: false,
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(5),
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
				DiskSizeGB:               pointer.MakePtr(10.0),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   pointer.MakePtr("label1"),
						Value: pointer.MakePtr("label1"),
					},
				},
				MongoDBMajorVersion: pointer.MakePtr("7.0"),
				MongoDBVersion:      pointer.MakePtr("7.1.5"),
				PitEnabled:          pointer.MakePtr(true),
				RootCertType:        pointer.MakePtr("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         pointer.MakePtr("LTS"),
				TerminationProtectionEnabled: pointer.MakePtr(false),
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
						RegionConfigs: &[]admin.CloudRegionConfig{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		"should return true when instance size has changed but autoscaling is enabled": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "REPLICASET",
						BackupEnabled: pointer.MakePtr(true),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          pointer.MakePtr(true),
						RootCertType:        "ISRGROOTX1",
						Tags: []*akov2.TagSpec{
							{
								Key:   "tag1",
								Value: "tag1",
							},
						},
						VersionReleaseSystem:         "LTS",
						TerminationProtectionEnabled: false,
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M40",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
				DiskSizeGB:               pointer.MakePtr(10.0),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   pointer.MakePtr("label1"),
						Value: pointer.MakePtr("label1"),
					},
				},
				MongoDBMajorVersion: pointer.MakePtr("7.0"),
				MongoDBVersion:      pointer.MakePtr("7.1.5"),
				PitEnabled:          pointer.MakePtr(true),
				RootCertType:        pointer.MakePtr("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         pointer.MakePtr("LTS"),
				TerminationProtectionEnabled: pointer.MakePtr(false),
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
						RegionConfigs: &[]admin.CloudRegionConfig{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  pointer.MakePtr("M10"),
										MaxInstanceSize:  pointer.MakePtr("M40"),
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		"should return true when cluster are the same with a unordered region": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "REPLICASET",
						BackupEnabled: pointer.MakePtr(true),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          pointer.MakePtr(true),
						RootCertType:        "ISRGROOTX1",
						Tags: []*akov2.TagSpec{
							{
								Key:   "tag1",
								Value: "tag1",
							},
						},
						VersionReleaseSystem:         "LTS",
						TerminationProtectionEnabled: false,
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M40",
											},
										},
									},
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M40",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.AdvancedClusterDescription{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
				DiskSizeGB:               pointer.MakePtr(10.0),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   pointer.MakePtr("label1"),
						Value: pointer.MakePtr("label1"),
					},
				},
				MongoDBMajorVersion: pointer.MakePtr("7.0"),
				MongoDBVersion:      pointer.MakePtr("7.1.5"),
				PitEnabled:          pointer.MakePtr(true),
				RootCertType:        pointer.MakePtr("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         pointer.MakePtr("LTS"),
				TerminationProtectionEnabled: pointer.MakePtr(false),
				ReplicationSpecs: &[]admin.ReplicationSpec{
					{
						NumShards: pointer.MakePtr(1),
						RegionConfigs: &[]admin.CloudRegionConfig{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_WEST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  pointer.MakePtr("M10"),
										MaxInstanceSize:  pointer.MakePtr("M40"),
									},
								},
							},
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  pointer.MakePtr("M10"),
										MaxInstanceSize:  pointer.MakePtr("M40"),
									},
								},
							},
						},
					},
				},
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//assert.Equal(t, NewDeployment(tt.ako, "project-id").(*Cluster).AdvancedDeploymentSpec, clusterFromAtlas(tt.atlas).AdvancedDeploymentSpec)
			assert.Equal(t, tt.expected, specAreEqual(NewDeployment("project-id", tt.ako).(*Cluster), clusterFromAtlas(tt.atlas)))
		})
	}
}
