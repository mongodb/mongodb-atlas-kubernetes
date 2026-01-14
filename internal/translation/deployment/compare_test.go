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
	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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
		"should update all spec when there are changes and disabling autoscaling": {
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
											Enabled: pointer.MakePtr(false),
										},
										Compute: &akov2.ComputeSpec{
											Enabled: pointer.MakePtr(false),
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ElectableSpecs: &admin.HardwareSpec20240805{
									DiskSizeGB: pointer.MakePtr(10.0),
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
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
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
				ConfigServerManagementMode: pointer.MakePtr("FIXED_TO_DEDICATED"),
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
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
			atlas: &admin.ClusterDescription20240805{
				TerminationProtectionEnabled: pointer.MakePtr(false),
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
				ClusterType: pointer.MakePtr("SHARDED"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						Id: pointer.MakePtr("abc123"),
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
			atlas: &admin.ClusterDescription20240805{
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
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
			atlas: &admin.ClusterDescription20240805{
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
				Name:        pointer.MakePtr("cluster0"),
				ClusterType: pointer.MakePtr("REPLICASET"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName:        pointer.MakePtr("TENANT"),
								BackingProviderName: pointer.MakePtr("AWS"),
								RegionName:          pointer.MakePtr("US_EAST_1"),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M0"),
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
										Priority:     pointer.MakePtr(7),
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
			atlas: &admin.ClusterDescription20240805{
				ClusterType: pointer.MakePtr("SHARDED"),
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
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
						BackupEnabled: pointer.MakePtr(true),
						DiskSizeGB:    pointer.MakePtr(20),
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
								NumShards: 3,
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
										Priority:     pointer.MakePtr(7),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("SHARDED"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
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
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
								},
							},
						},
					},
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
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
						BackupEnabled: pointer.MakePtr(true),
						DiskSizeGB:    pointer.MakePtr(20),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
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
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
									DiskSizeGB:   pointer.MakePtr(20.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(20.0),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),
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
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_EAST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
									DiskSizeGB:   pointer.MakePtr(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(10.0),
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
			atlas: &admin.ClusterDescription20240805{
				Name:                     pointer.MakePtr("cluster0"),
				ClusterType:              pointer.MakePtr("REPLICASET"),
				BackupEnabled:            pointer.MakePtr(true),
				EncryptionAtRestProvider: pointer.MakePtr("NONE"),
				Paused:                   pointer.MakePtr(false),

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
				ReplicationSpecs: &[]admin.ReplicationSpec20240805{
					{
						RegionConfigs: &[]admin.CloudRegionConfig20240805{
							{
								ProviderName: pointer.MakePtr("AWS"),
								RegionName:   pointer.MakePtr("US_WEST_1"),
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
									DiskSizeGB:   pointer.MakePtr(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(10.0),
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
								ElectableSpecs: &admin.HardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(3),
									DiskSizeGB:   pointer.MakePtr(10.0),
								},
								ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
									InstanceSize: pointer.MakePtr("M10"),
									NodeCount:    pointer.MakePtr(5),
									DiskSizeGB:   pointer.MakePtr(10.0),
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
			assert.Equal(t, tt.expected, specAreEqual(NewDeployment("project-id", tt.ako).(*Cluster), clusterFromAtlas(tt.atlas)))
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
						Priority:     pointer.MakePtr(7),
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
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
						Priority:     pointer.MakePtr(7),
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
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
				Priority:     pointer.MakePtr(7),
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(6),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when electable spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when read-only spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when analytics spec has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
			},
			autoscalingEnabled: false,
			expected:           false,
		},
		"should return false when autoscaling has changed": {
			akoRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: pointer.MakePtr(true),
					},
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
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
				Priority:     pointer.MakePtr(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: pointer.MakePtr(true),
					},
				},
			},
			atlasRegionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AWS",
				RegionName:   "EU_CENTRAL_1",
				Priority:     pointer.MakePtr(7),
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					DiskGB: &akov2.DiskGB{
						Enabled: pointer.MakePtr(true),
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
				NodeCount:    pointer.MakePtr(3),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
				NodeCount:    pointer.MakePtr(1),
			},
		},
		"should return false when ebs volume has changed": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     pointer.MakePtr(3),
				EbsVolumeType: "STANDARD",
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize: "M20",
				NodeCount:    pointer.MakePtr(3),
			},
		},
		"should return false when disk iop has changed": {
			akoNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     pointer.MakePtr(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      pointer.MakePtr(int64(3000)),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize:  "M20",
				NodeCount:     pointer.MakePtr(3),
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
				NodeCount:     pointer.MakePtr(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      pointer.MakePtr(int64(5000)),
			},
			atlasNodeSpec: &akov2.Specs{
				InstanceSize:  "M10",
				NodeCount:     pointer.MakePtr(3),
				EbsVolumeType: "STANDARD",
				DiskIOPS:      pointer.MakePtr(int64(5000)),
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
					Enabled: pointer.MakePtr(true),
				},
			},
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
				DiskGB: &akov2.DiskGB{
					Enabled: pointer.MakePtr(false),
				},
			},
		},
		"should return false when compute autoscaling has changed": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled: pointer.MakePtr(true),
				},
			},
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled: pointer.MakePtr(false),
				},
			},
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.AdvancedAutoScalingSpec{
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
			atlasAutoscaling: &akov2.AdvancedAutoScalingSpec{
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
				Enabled: pointer.MakePtr(true),
			},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: pointer.MakePtr(false),
			},
		},
		"should return true when autoscaling enabled flag is unset": {
			akoAutoscaling: &akov2.DiskGB{},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: pointer.MakePtr(true),
			},
			expected: true,
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.DiskGB{
				Enabled: pointer.MakePtr(true),
			},
			atlasAutoscaling: &akov2.DiskGB{
				Enabled: pointer.MakePtr(true),
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
				Enabled: pointer.MakePtr(true),
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				Enabled: pointer.MakePtr(false),
			},
		},
		"should return false when scale down enabled flag has changed": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(true),
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(false),
			},
		},
		"should return false when min instance has changed and scale down is enabled": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(true),
				MinInstanceSize:  "M20",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(true),
				MinInstanceSize:  "M10",
			},
		},
		"should return true when min instance has changed and scale down is disabled": {
			akoAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(false),
				MinInstanceSize:  "M20",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				ScaleDownEnabled: pointer.MakePtr(false),
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
				Enabled:          pointer.MakePtr(true),
				ScaleDownEnabled: pointer.MakePtr(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M40",
			},
			expected: true,
		},
		"should return true when autoscaling are equal": {
			akoAutoscaling: &akov2.ComputeSpec{
				Enabled:          pointer.MakePtr(true),
				ScaleDownEnabled: pointer.MakePtr(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M40",
			},
			atlasAutoscaling: &akov2.ComputeSpec{
				Enabled:          pointer.MakePtr(true),
				ScaleDownEnabled: pointer.MakePtr(true),
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
