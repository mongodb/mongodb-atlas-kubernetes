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

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
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
					BackupEnabled: new(false),
					Paused:        new(true),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					Paused:        new(false),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:   "cluster0",
					Paused: new(true),
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
					BackupEnabled: new(false),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					DiskSizeGB:    new(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    new(false),
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
					BackupEnabled: new(false),
					DiskSizeGB:    new(20),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					DiskSizeGB:    new(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    new(false),
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
					BackupEnabled: new(false),
					DiskSizeGB:    new(30),
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					DiskSizeGB:    new(20),
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:             "cluster0",
					ClusterType:      "REPLICASET",
					BackupEnabled:    new(false),
					DiskSizeGB:       new(30),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{},
				},
			},
			changed: true,
		},
		"should update all spec when there are changes and disabling autoscaling": {
			//nolint:dupl
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               new(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        new(true),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
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
					BackupEnabled: new(true),
					DiskSizeGB:    new(20),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
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
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               new(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        new(true),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(false),
										},
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
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
		"should update all spec when there are changes and cnhage autoscaling": {
			//nolint:dupl
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               new(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        new(true),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
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
					BackupEnabled: new(true),
					DiskSizeGB:    new(20),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
								},
							},
						},
					},
				},
			},
			//nolint:dupl
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
					EncryptionAtRestProvider: "AWS",
					MongoDBMajorVersion:      "8.0",
					RootCertType:             "ISRGROOTX1",
					PitEnabled:               new(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        new(true),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
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
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
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
					BackupEnabled:            new(false),
					DiskSizeGB:               new(30),
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
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										DiskGB: &akov2.DiskGB{
											Enabled: new(true),
										},
										Compute: &akov2.ComputeSpec{
											Enabled:          new(true),
											ScaleDownEnabled: new(true),
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
		// Bug exposure: AutoScaling should NOT be included when it hasn't changed (both desired and current are nil).
		// This exposes the bug where getAutoScalingChanges always returns a non-nil value with default/empty fields
		// even when desired is nil, causing AutoScaling to always be included in changes.
		// Expected: AutoScaling should be nil in changes when unchanged
		// Actual (buggy): AutoScaling is included with {DiskGB: {Enabled: false}, Compute: {Enabled: false}} even when unchanged
		// This causes reconciliation loops because we send unnecessary updates to Atlas.
		"BUG_EXPOSURE: should not include AutoScaling in changes when AutoScaling hasn't changed": {
			akoCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "aws-cluster",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize:  "M10",
										NodeCount:     new(3),
										EbsVolumeType: "PROVISIONED", // Explicitly specified, different from Atlas
									},
									// AutoScaling is nil (not specified)
								},
							},
						},
					},
				},
			},
			atlasCluster: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "aws-cluster",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize:  "M10",
										NodeCount:     new(3),
										EbsVolumeType: "STANDARD", // Different from desired
									},
									// AutoScaling is also nil (no change)
								},
							},
						},
					},
				},
			},
			expectedChanges: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "aws-cluster",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_EAST_1",
									Priority:     new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize:  "M10",
										NodeCount:     new(3),
										EbsVolumeType: "PROVISIONED", // Only this should change
									},
									// AutoScaling should be nil (not included when unchanged)
									// Bug: getAutoScalingChanges returns non-nil with default values even when desired is nil
								},
							},
						},
					},
				},
			},
			changed: true,
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

func TestSpecAreEqual(t *testing.T) {
	tests := map[string]struct {
		ako      *akov2.AtlasDeployment
		atlas    *admin.ClusterDescription20240805
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
			atlas: &admin.ClusterDescription20240805{
				ClusterType: new("REPLICASET"),
			},
		},
		"should return false when backup enabled flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						BackupEnabled: new(true),
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				BackupEnabled: new(false),
			},
		},
		"should return false when BI connector config are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						BiConnector: &akov2.BiConnectorSpec{
							Enabled:        new(true),
							ReadPreference: "secondary",
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				BiConnector: &admin.BiConnector{
					Enabled:        new(false),
					ReadPreference: new("secondary"),
				},
			},
		},
		"should return false when disk size are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						DiskSizeGB: new(20),
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ElectableSpecs: &admin.HardwareSpec20240805{
									DiskSizeGB: new(10.0),
								},
							},
						},
					},
				},
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
			atlas: &admin.ClusterDescription20240805{
				EncryptionAtRestProvider: new("NONE"),
			},
		},
		"should return false when config server management are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ConfigServerManagementMode: "ATLAS_MANAGED",
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				ConfigServerManagementMode: new("FIXED_TO_DEDICATED"),
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
			atlas: &admin.ClusterDescription20240805{
				MongoDBMajorVersion: new("7.0"),
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
			atlas: &admin.ClusterDescription20240805{
				VersionReleaseSystem: new("LTS"),
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
			atlas: &admin.ClusterDescription20240805{
				RootCertType: new("NONE"),
			},
		},
		"should return false when paused flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Paused: new(true),
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				Paused: new(false),
			},
		},
		"should return false when pit flag are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						PitEnabled: new(true),
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				PitEnabled: new(false),
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
			atlas: &admin.ClusterDescription20240805{
				TerminationProtectionEnabled: new(false),
			},
		},
		"should return false when num of shards are different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ClusterType: "SHARDED",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 3,
							},
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				ClusterType: new("SHARDED"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						Id: new("abc123"),
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
			atlas: &admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
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
												Enabled:          new(true),
												ScaleDownEnabled: new(true),
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
			atlas: &admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled: new(false),
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
			atlas: &admin.ClusterDescription20240805{
				Labels: &[]admin.ComponentLabel{
					{
						Key:   new("label2"),
						Value: new("label2"),
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
			atlas: &admin.ClusterDescription20240805{
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag2",
						Value: "tag2",
					},
				},
			},
		},
		"should return false when instance size of shared cluster changed": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName:        "TENANT",
										BackingProviderName: "AWS",
										RegionName:          "US_EAST_1",
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M2",
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				Name:        new("cluster0"),
				ClusterType: new("REPLICASET"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName:        new("TENANT"),
								BackingProviderName: new("AWS"),
								RegionName:          new("US_EAST_1"),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: new("M0"),
								},
							},
						},
					},
				},
			},
			expected: false,
		},
		"should return false when shard count is different": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "SHARDED",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     new(7),
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(5),
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				ClusterType: new("SHARDED"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
								},
							},
						},
					},
				},
			},
			expected: false,
		},
		"should return true when sharded cluster are the same": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "SHARDED",
						BackupEnabled: new(true),
						DiskSizeGB:    new(20),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          new(true),
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
								NumShards: 3,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     new(7),
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(5),
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				Name:                     new("cluster0"),
				ClusterType:              new("SHARDED"),
				BackupEnabled:            new(true),
				EncryptionAtRestProvider: new("NONE"),
				Paused:                   new(false),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   new("label1"),
						Value: new("label1"),
					},
				},
				MongoDBMajorVersion: new("7.0"),
				MongoDBVersion:      new("7.1.5"),
				PitEnabled:          new(true),
				RootCertType:        new("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         new("LTS"),
				TerminationProtectionEnabled: new(false),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		"should return true when cluster are the same": {
			ako: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "REPLICASET",
						BackupEnabled: new(true),
						DiskSizeGB:    new(20),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          new(true),
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
										Priority:     new(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(5),
										},
									},
								},
							},
						},
					},
				},
			},
			atlas: &admin.ClusterDescription20240805{
				Name:                     new("cluster0"),
				ClusterType:              new("REPLICASET"),
				BackupEnabled:            new(true),
				EncryptionAtRestProvider: new("NONE"),
				Paused:                   new(false),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   new("label1"),
						Value: new("label1"),
					},
				},
				MongoDBMajorVersion: new("7.0"),
				MongoDBVersion:      new("7.1.5"),
				PitEnabled:          new(true),
				RootCertType:        new("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         new("LTS"),
				TerminationProtectionEnabled: new(false),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(3),
									DiskSizeGB:   new(20.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(20.0),
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
						BackupEnabled: new(true),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          new(true),
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
										Priority:     new(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          new(true),
												ScaleDownEnabled: new(true),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     new("cluster0"),
				ClusterType:              new("REPLICASET"),
				BackupEnabled:            new(true),
				EncryptionAtRestProvider: new("NONE"),
				Paused:                   new(false),
				Labels: &[]admin.ComponentLabel{
					{
						Key:   new("label1"),
						Value: new("label1"),
					},
				},
				MongoDBMajorVersion: new("7.0"),
				MongoDBVersion:      new("7.1.5"),
				PitEnabled:          new(true),
				RootCertType:        new("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         new("LTS"),
				TerminationProtectionEnabled: new(false),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(3),
									DiskSizeGB:   new(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(10.0),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          new(true),
										ScaleDownEnabled: new(true),
										MinInstanceSize:  new("M10"),
										MaxInstanceSize:  new("M40"),
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
						BackupEnabled: new(true),
						Labels: []common.LabelSpec{
							{
								Key:   "label1",
								Value: "label1",
							},
						},
						MongoDBMajorVersion: "7.0",
						PitEnabled:          new(true),
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
										Priority:     new(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          new(true),
												ScaleDownEnabled: new(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M40",
											},
										},
									},
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     new(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(3),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    new(5),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          new(true),
												ScaleDownEnabled: new(true),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     new("cluster0"),
				ClusterType:              new("REPLICASET"),
				BackupEnabled:            new(true),
				EncryptionAtRestProvider: new("NONE"),
				Paused:                   new(false),

				Labels: &[]admin.ComponentLabel{
					{
						Key:   new("label1"),
						Value: new("label1"),
					},
				},
				MongoDBMajorVersion: new("7.0"),
				MongoDBVersion:      new("7.1.5"),
				PitEnabled:          new(true),
				RootCertType:        new("ISRGROOTX1"),
				Tags: &[]admin.ResourceTag{
					{
						Key:   "tag1",
						Value: "tag1",
					},
				},
				VersionReleaseSystem:         new("LTS"),
				TerminationProtectionEnabled: new(false),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_WEST_1"),
								Priority:     new(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(3),
									DiskSizeGB:   new(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(10.0),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          new(true),
										ScaleDownEnabled: new(true),
										MinInstanceSize:  new("M10"),
										MaxInstanceSize:  new("M40"),
									},
								},
							},
							{
								ProviderName: new("AWS"),
								RegionName:   new("US_EAST_1"),
								Priority:     new(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(3),
									DiskSizeGB:   new(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: new("M10"),
									NodeCount:    new(5),
									DiskSizeGB:   new(10.0),
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          new(true),
										ScaleDownEnabled: new(true),
										MinInstanceSize:  new("M10"),
										MaxInstanceSize:  new("M40"),
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
			assert.Equal(t, tt.expected, specAreEqual(NewDeployment("project-id", tt.ako).(*Cluster), clusterFromAtlas(tt.atlas)))
		})
	}
}

func TestComputeChangesEbsVolumeTypeBugGCP(t *testing.T) {
	// This test reproduces the bug where GCP clusters get stuck in a reconcile loop
	// because getSpecsChanges always includes EbsVolumeType (even when empty),
	// which gets converted to "STANDARD" by replicationSpecToAtlas, causing Atlas
	// to reject/ignore it, and the operator to detect it as a change again.

	tests := map[string]struct {
		akoCRD          *akov2.AtlasDeployment
		atlasResponse   *admin.ClusterDescription20240805
		expectedChanges *Cluster
		expectedChanged bool
		description     string
	}{
		"BUG_REPRODUCTION: GCP cluster with NodeCount change should not include EbsVolumeType": {
			description: "When a GCP cluster has a real change (NodeCount), getSpecsChanges includes EbsVolumeType " +
				"as empty string. This gets converted to STANDARD by replicationSpecToAtlas, causing Atlas to reject it for GCP. " +
				"This is the actual bug flow that causes infinite reconcile loops.",
			akoCRD: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "credit-service-uat",
						ClusterType:   "REPLICASET",
						BackupEnabled: new(true),
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								ZoneName:  "Zone 1",
								NumShards: 1,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "GCP",
										RegionName:   "EASTERN_US",
										Priority:     new(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(3), // Different from Atlas to force changes
											// EbsVolumeType is NOT specified (correct for GCP)
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          new(true),
												MaxInstanceSize:  "M30",
												MinInstanceSize:  "M10",
												ScaleDownEnabled: new(true),
											},
											DiskGB: &akov2.DiskGB{
												Enabled: new(false),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			atlasResponse: &admin.ClusterDescription20240805{
				Name:          new("credit-service-uat"),
				ClusterType:   new("REPLICASET"),
				BackupEnabled: new(true),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						ZoneName: new("Zone 1"),
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: new("GCP"),
								RegionName:   new("EASTERN_US"),
								Priority:     new(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize:  new("M10"),
									NodeCount:     new(2), // Different from desired
									EbsVolumeType: nil,    // Atlas doesn't return EbsVolumeType for GCP
								},
								AutoScaling: &admin.AdvancedAutoScalingSettings{
									Compute: &admin.AdvancedComputeAutoScaling{
										Enabled:          new(true),
										MaxInstanceSize:  new("M30"),
										MinInstanceSize:  new("M10"),
										ScaleDownEnabled: new(true),
									},
									DiskGB: &admin.DiskGBAutoScaling{
										Enabled: new(false),
									},
								},
							},
						},
					},
				},
			},
			expectedChanged: true, // Changes expected due to NodeCount difference
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("Test description: %s", tt.description)

			// Normalize inputs the same way the code does before comparison
			akoCluster := NewDeployment("project-id", tt.akoCRD).(*Cluster)
			atlasCluster := clusterFromAtlas(tt.atlasResponse)

			// Compute changes
			changes, changed := ComputeChanges(akoCluster, atlasCluster)

			// Verify expectations
			assert.Equal(t, tt.expectedChanged, changed, "changed flag mismatch")
			assert.NotNil(t, changes, "changes should not be nil when changed is true")

			atlasUpdateRequest := clusterUpdateToAtlas(changes)
			assert.NotNil(t, atlasUpdateRequest, "atlasUpdateRequest should not be nil")
			assert.NotNil(t, atlasUpdateRequest.ReplicationSpecs, "atlasUpdateRequest.ReplicationSpecs should not be nil")
			for _, repSpec := range *atlasUpdateRequest.ReplicationSpecs {
				for _, regionConfig := range repSpec.GetRegionConfigs() {
					assert.NotEqual(t, "AWS", regionConfig.GetProviderName(), "This test is for non AWS configs")
					assert.NotNil(t, regionConfig.ElectableSpecs, "regionConfig.ElectableSpecs should not be nil")
					assert.Nil(t, regionConfig.ElectableSpecs.EbsVolumeType, "regionConfig.ElectableSpecs.EbsVolumeType should not be nil")
				}
			}
		})
	}
}

func TestReplicationSpecAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoReplicationSpec   *akov2.AdvancedReplicationSpec
		atlasReplicationSpec *akov2.AdvancedReplicationSpec
		autoscalingEnabled   bool
		expected             bool
	}{
		"should return false when zone name has changed": {
			akoReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName: "Zone 1",
			},
			atlasReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName: "First Zone",
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when new region was added": {
			akoReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{},
					{},
				},
			},
			atlasReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{},
				},
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when new region config has changed": {
			akoReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ProviderName: "AWS",
						RegionName:   "EU_CENTRAL_1",
					},
				},
			},
			atlasReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ProviderName: "AWS",
						RegionName:   "EU_WEST_1",
					},
				},
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return true when spec are equal": {
			akoReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ProviderName: "AWS",
						RegionName:   "EU_CENTRAL_1",
						Priority:     new(7),
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    new(3),
						},
					},
				},
			},
			atlasReplicationSpec: &akov2.AdvancedReplicationSpec{
				ZoneName:  "Zone 1",
				NumShards: 1,
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ProviderName: "AWS",
						RegionName:   "EU_CENTRAL_1",
						Priority:     new(7),
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    new(3),
						},
					},
				},
			},
			autoscalingEnabled: false,
			expected:           true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, replicationSpecAreEqual(tt.akoReplicationSpec, tt.atlasReplicationSpec, tt.autoscalingEnabled))
		})
	}
}

func TestRegionConfigAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoRegionConfig    *akov2.AdvancedRegionConfig
		atlasRegionConfig  *akov2.AdvancedRegionConfig
		autoscalingEnabled bool
		expected           bool
	}{
		"should return false when provider has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "GCP",
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when region has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_WEST_1",
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when priority has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(6),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when electable spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when read-only spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when analytics spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when autoscaling has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: new(true),
					},
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when backing provider has changed for tenant instances": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName:        "TENANT",
				RegionName:          "US_EAST_1",
				BackingProviderName: "AWS",
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName:        "TENANT",
				RegionName:          "US_EAST_1",
				BackingProviderName: "GCP",
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return true when region config are equal": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: new(true),
					},
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     new(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    new(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: new(true),
					},
				},
			},
			autoscalingEnabled: false,
			expected:           true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, regionConfigAreEqual(tt.akoRegionConfig, tt.atlasRegionConfig, tt.autoscalingEnabled))
		})
	}
}

func TestNodeSpecAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoNodeSpec        *akov2.Specs
		atlasNodeSpec      *akov2.Specs
		autoscalingEnabled bool
		expected           bool
	}{
		"should return true when both specs are unset": {
			expected: true,
		},
		"should return false when ako spec is set and atlas spec is unset": {
			akoNodeSpec: &akov2.Specs{},
		},
		"should return false when ako spec is unset and atlas spec is set": {
			atlasNodeSpec: &akov2.Specs{},
		},
		"should return false when instance size has changed and autoscaling is disabled": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M10",
			},
		},
		"should return false when node count has changed": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
				NodeCount:    new(3),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
				NodeCount:    new(1),
			},
		},
		"should return false when ebs volume has changed": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     new(3),
				EbsVolumeType: "STANDARD",
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
				NodeCount:    new(3),
			},
		},
		"should return false when disk iop has changed": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     new(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      new(int64(3000)),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     new(3),
				EbsVolumeType: "STANDARD",
			},
		},
		"should return true when instance size has changed and autoscaling is enabled": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M10",
			},
			autoscalingEnabled: true,
			expected:           true,
		},
		"should return true when specs are equal": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize:  "M10",
				NodeCount:     new(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      new(int64(5000)),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize:  "M10",
				NodeCount:     new(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      new(int64(5000)),
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, nodeSpecAreEqual(tt.akoNodeSpec, tt.atlasNodeSpec, tt.autoscalingEnabled))
		})
	}
}

func TestAutoscalingConfigAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoAutoscaling   *akov2.AdvancedAutoScalingSpec
		atlasAutoscaling *akov2.AdvancedAutoScalingSpec
		expected         bool
	}{
		"should return true both autoscaling are unset": {
			expected: true,
		},
		"should return false when ako autoscaling is set and atlas autoscaling is unset": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{},
		},
		"should return false when ako autoscaling is unset and atlas autoscaling is set": {
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{},
		},
		"should return false when disk autoscaling has changed": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: new(true),
				},
			},
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: new(false),
				},
			},
		},
		"should return false when compute autoscaling has changed": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled: new(true),
				},
			},
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled: new(false),
				},
			},
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: new(true),
				},
				Compute: &akov2.ComputeSpec{
					Enabled:          new(true),
					ScaleDownEnabled: new(true),
					MinInstanceSize:  "M10",
					MaxInstanceSize:  "M40",
				},
			},
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: new(true),
				},
				Compute: &akov2.ComputeSpec{
					Enabled:          new(true),
					ScaleDownEnabled: new(true),
					MinInstanceSize:  "M10",
					MaxInstanceSize:  "M40",
				},
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, autoscalingConfigAreEqual(tt.akoAutoscaling, tt.atlasAutoscaling))
		})
	}
}

func TestDiskAutoscalingConfigAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoAutoscaling   *akov2.DiskGB
		atlasAutoscaling *akov2.DiskGB
		expected         bool
	}{
		"should return true both autoscaling are unset": {
			expected: true,
		},
		"should return false when ako autoscaling is set and atlas autoscaling is unset": {
			akoAutoscaling: &akov2.DiskGB{},
		},
		"should return false when ako autoscaling is unset and atlas autoscaling is set": {
			atlasAutoscaling: &akov2.DiskGB{},
		},
		"should return false when autoscaling has changed": {
			akoAutoscaling: &akov2.DiskGB{
				Enabled: new(true),
			},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: new(false),
			},
		},
		"should return true when autoscaling enabled flag is unset": {
			akoAutoscaling: &akov2.DiskGB{},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: new(true),
			},
			expected: true,
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.DiskGB{
				Enabled: new(true),
			},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: new(true),
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, diskAutoscalingConfigAreEqual(tt.akoAutoscaling, tt.atlasAutoscaling))
		})
	}
}

func TestComputeAutoscalingConfigAreEqual(t *testing.T) {
	tests := map[string]struct {
		akoAutoscaling   *akov2.ComputeSpec
		atlasAutoscaling *akov2.ComputeSpec
		expected         bool
	}{
		"should return true both autoscaling are unset": {
			expected: true,
		},
		"should return false when ako autoscaling is set and atlas autoscaling is unset": {
			akoAutoscaling: &akov2.ComputeSpec{},
		},
		"should return false when ako autoscaling is unset and atlas autoscaling is set": {
			atlasAutoscaling: &akov2.ComputeSpec{},
		},
		"should return false when enabled flag has changed": {
			akoAutoscaling: &akov2.ComputeSpec{
				Enabled: new(true),
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				Enabled: new(false),
			},
		},
		"should return false when scale down enabled flag has changed": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(true),
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(false),
			},
		},
		"should return false when min instance has changed and scale down is enabled": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(true),
				MinInstanceSize:  "M20",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(true),
				MinInstanceSize:  "M10",
			},
		},
		"should return true when min instance has changed and scale down is disabled": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(false),
				MinInstanceSize:  "M20",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: new(false),
				MinInstanceSize:  "M10",
			},
			expected: true,
		},
		"should return false when max instance has changed": {
			akoAutoscaling: &akov2.ComputeSpec{
				MaxInstanceSize: "M20",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				MaxInstanceSize: "M10",
			},
		},
		"should return true when autoscaling enabled flags are unset": {
			akoAutoscaling: &akov2.ComputeSpec{
				MinInstanceSize: "M10",
				MaxInstanceSize: "M40",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				Enabled:          new(true),
				ScaleDownEnabled: new(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M40",
			},
			expected: true,
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.ComputeSpec{
				Enabled:          new(true),
				ScaleDownEnabled: new(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M40",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				Enabled:          new(true),
				ScaleDownEnabled: new(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M40",
			},
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, computeAutoscalingConfigAreEqual(tt.akoAutoscaling, tt.atlasAutoscaling))
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
				current:  new(false),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  new(false),
				current:  new(false),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  new(false),
				expected: true,
			},
			"desired and current are true": {
				desired:  new(true),
				current:  new(true),
				expected: true,
			},
			"desired is nil and current is true": {
				current:  new(true),
				expected: false,
			},
			"desired is true and current is nil": {
				desired:  new(true),
				expected: false,
			},
			"desired is false and current is true": {
				desired:  new(false),
				current:  new(true),
				expected: false,
			},
			"desired is true and current is false": {
				desired:  new(true),
				current:  new(false),
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
				current:  new(""),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  new(""),
				current:  new(""),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  new(""),
				expected: true,
			},
			"desired and current have value": {
				desired:  new("value"),
				current:  new("value"),
				expected: true,
			},
			"desired is nil and current has value": {
				current:  new("value"),
				expected: false,
			},
			"desired has value and current is nil": {
				desired:  new("value"),
				expected: false,
			},
			"desired has different value from current": {
				desired:  new("value"),
				current:  new("other-value"),
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
				current:  new(0),
				expected: true,
			},
			"desired and current are pointer to zero value": {
				desired:  new(0),
				current:  new(0),
				expected: true,
			},
			"desired is pointer to zero value and current is nil": {
				desired:  new(0),
				expected: true,
			},
			"desired and current have value": {
				desired:  new(10),
				current:  new(10),
				expected: true,
			},
			"desired is nil and current has value": {
				current:  new(10),
				expected: false,
			},
			"desired has value and current is nil": {
				desired:  new(10),
				expected: false,
			},
			"desired has different value from current": {
				desired:  new(10),
				current:  new(11),
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
