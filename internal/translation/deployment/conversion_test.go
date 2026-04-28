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

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestNewDeployment(t *testing.T) {
	tests := map[string]struct {
		cr       *akov2.AtlasDeployment
		expected Deployment
	}{
		"should create a new serverless deployment as flex": {
			cr: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        "SERVERLESS",
							BackingProviderName: "AWS",
							RegionName:          "US_EAST_1",
						},
						PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
							{
								Name:                    "spe1",
								CloudProviderEndpointID: "1234567890",
							},
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "name",
								Value: "test",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: true,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			expected: &Flex{
				FlexSpec: &akov2.FlexSpec{
					Name: "instance0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "name",
							Value: "test",
						},
					},
					TerminationProtectionEnabled: true,
				},
				ProjectID: "project-id",
				customResource: &akov2.AtlasDeployment{
					Spec: akov2.AtlasDeploymentSpec{
						ServerlessSpec: &akov2.ServerlessSpec{
							Name: "instance0",
							ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
								ProviderName:        "SERVERLESS",
								BackingProviderName: "AWS",
								RegionName:          "US_EAST_1",
							},
							PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
								{
									Name:                    "spe1",
									CloudProviderEndpointID: "1234567890",
								},
							},
							Tags: []*akov2.TagSpec{
								{
									Key:   "name",
									Value: "test",
								},
							},
							BackupOptions: akov2.ServerlessBackupOptions{
								ServerlessContinuousBackupEnabled: true,
							},
							TerminationProtectionEnabled: true,
						},
					},
				},
			},
		},
		"should create a new regular deployment": {
			//nolint:dupl
			cr: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					//nolint:dupl
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          "cluster0",
						ClusterType:   "REPLICASET",
						BackupEnabled: new(true),
						BiConnector: &akov2.BiConnectorSpec{
							Enabled:        new(true),
							ReadPreference: "secondary",
						},
						DiskSizeGB:               new(20),
						EncryptionAtRestProvider: "AWS",
						Paused:                   new(false),
						PitEnabled:               new(true),
						MongoDBMajorVersion:      "8.0",
						VersionReleaseSystem:     "LTS",
						RootCertType:             "ISRGROOTX1",
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
						TerminationProtectionEnabled: true,
						ConfigServerManagementMode:   "ATLAS_MANAGED",
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
					},
					ProcessArgs: &akov2.ProcessArgs{
						MinimumEnabledTLSProtocol: "TLS1_2",
						JavascriptEnabled:         new(true),
					},
				},
			},
			expected: &Cluster{
				ProjectID: "project-id",
				//nolint:dupl
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: new(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        new(true),
						ReadPreference: "secondary",
					},
					DiskSizeGB:               new(20),
					EncryptionAtRestProvider: "AWS",
					Paused:                   new(false),
					PitEnabled:               new(true),
					MongoDBMajorVersion:      "8.0",
					VersionReleaseSystem:     "LTS",
					RootCertType:             "ISRGROOTX1",
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
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: new(false),
										},
									},
								},
							},
						},
					},
					TerminationProtectionEnabled: true,
					ConfigServerManagementMode:   "ATLAS_MANAGED",
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
				},
				ProcessArgs: &akov2.ProcessArgs{
					MinimumEnabledTLSProtocol: "TLS1_2",
					JavascriptEnabled:         new(true),
					NoTableScan:               new(false),
				},
				//nolint:dupl
				customResource: &akov2.AtlasDeployment{
					Spec: akov2.AtlasDeploymentSpec{
						//nolint:dupl
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{
							Name:          "cluster0",
							ClusterType:   "REPLICASET",
							BackupEnabled: new(true),
							BiConnector: &akov2.BiConnectorSpec{
								Enabled:        new(true),
								ReadPreference: "secondary",
							},
							DiskSizeGB:               new(20),
							EncryptionAtRestProvider: "AWS",
							Paused:                   new(false),
							PitEnabled:               new(true),
							MongoDBMajorVersion:      "8.0",
							VersionReleaseSystem:     "LTS",
							RootCertType:             "ISRGROOTX1",
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
							TerminationProtectionEnabled: true,
							ConfigServerManagementMode:   "ATLAS_MANAGED",
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
						},
						ProcessArgs: &akov2.ProcessArgs{
							MinimumEnabledTLSProtocol: "TLS1_2",
							JavascriptEnabled:         new(true),
						},
					},
				},
				computeAutoscalingEnabled: false,
				isTenant:                  false,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, NewDeployment("project-id", tt.cr))
		})
	}
}

func TestNormalizeClusterDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *Cluster
		expected   *Cluster
	}{
		"nil replication spec": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: nil,
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: nil,
				},
			},
		},
		"nil replication spec entries": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						nil, nil, nil,
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						nil, nil, nil,
					},
				},
			},
		},
		"nil region configs": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: nil,
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							NumShards:     1,
							ZoneName:      "Zone 1",
							RegionConfigs: nil,
						},
					},
				},
			},
		},
		"nil region config entries": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								nil, nil, nil,
							},
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							NumShards: 1,
							ZoneName:  "Zone 1",
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								nil, nil, nil,
							},
						},
					},
				},
			},
		},
		"nil regionconfig specs": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: nil,
									ElectableSpecs: nil,
									ReadOnlySpecs:  nil,
								},
							},
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							NumShards: 1,
							ZoneName:  "Zone 1",
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: nil,
									ElectableSpecs: nil,
									ReadOnlySpecs:  nil,
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
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
		"empty regionconfig specs": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: &akov2.Specs{},
									ElectableSpecs: &akov2.Specs{},
									ReadOnlySpecs:  &akov2.Specs{},
								},
							},
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							NumShards: 1,
							ZoneName:  "Zone 1",
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: nil,
									ElectableSpecs: nil,
									ReadOnlySpecs:  nil,
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
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
		"regionconfig specs with autoscaling": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: &akov2.Specs{},
									ElectableSpecs: &akov2.Specs{},
									ReadOnlySpecs:  &akov2.Specs{},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled:          new(false),
											ScaleDownEnabled: new(false),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M10",
										},
									},
								},
							},
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ClusterType:              "REPLICASET",
					EncryptionAtRestProvider: "NONE",
					VersionReleaseSystem:     "LTS",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            new(false),
					Tags:                     []*akov2.TagSpec{},

					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							NumShards: 1,
							ZoneName:  "Zone 1",
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									AnalyticsSpecs: nil,
									ElectableSpecs: nil,
									ReadOnlySpecs:  nil,
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
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
		"normalize deployment configuration": {
			deployment: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					Labels: []common.LabelSpec{
						{
							Key:   "b",
							Value: "b",
						},
						{
							Key:   "a",
							Value: "a",
						},
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "b",
							Value: "b",
						},
						{
							Key:   "a",
							Value: "a",
						},
					},
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName:        "AWS",
									BackingProviderName: "AWS",
									RegionName:          "US_EAST_1",
									Priority:            new(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    new(3),
									},
									ReadOnlySpecs: &akov2.Specs{
										InstanceSize: "M10",
									},
								},
							},
						},
					},
				},
			},
			expected: &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            new(false),
					VersionReleaseSystem:     "LTS",
					EncryptionAtRestProvider: "NONE",
					Paused:                   new(false),
					PitEnabled:               new(false),
					RootCertType:             "ISRGROOTX1",
					Labels: []common.LabelSpec{

						{
							Key:   "a",
							Value: "a",
						},
						{
							Key:   "b",
							Value: "b",
						},
					},
					Tags: []*akov2.TagSpec{

						{
							Key:   "a",
							Value: "a",
						},
						{
							Key:   "b",
							Value: "b",
						},
					},
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
										Compute: &akov2.ComputeSpec{
											Enabled: new(false),
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
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			normalizeClusterDeployment(tt.deployment)

			assert.Equal(t, tt.expected, tt.deployment)
		})
	}
}

func TestNormalizeRegionConfigsPriorityOrdering(t *testing.T) {
	tests := map[string]struct {
		regionConfigs []*akov2.AdvancedRegionConfig
		expectedOrder []int // Expected priority values in order
	}{
		"should sort priorities in descending order": {
			regionConfigs: []*akov2.AdvancedRegionConfig{
				{
					ProviderName: "AWS",
					RegionName:   "US_EAST_1",
					Priority:     new(5),
				},
				{
					ProviderName: "AWS",
					RegionName:   "US_WEST_1",
					Priority:     new(7),
				},
				{
					ProviderName: "AWS",
					RegionName:   "EU_WEST_1",
					Priority:     new(6),
				},
			},
			expectedOrder: []int{7, 6, 5},
		},
		"should sort by priority descending first, then provider+region": {
			regionConfigs: []*akov2.AdvancedRegionConfig{
				{
					ProviderName: "AWS",
					RegionName:   "US_EAST_1",
					Priority:     new(5),
				},
				{
					ProviderName: "AWS",
					RegionName:   "US_EAST_1",
					Priority:     new(7),
				},
				{
					ProviderName: "GCP",
					RegionName:   "US_CENTRAL_1",
					Priority:     new(6),
				},
			},
			expectedOrder: []int{7, 6, 5},
		},
		"should handle nil priorities as 0": {
			regionConfigs: []*akov2.AdvancedRegionConfig{
				{
					ProviderName: "AWS",
					RegionName:   "US_EAST_1",
					Priority:     nil,
				},
				{
					ProviderName: "AWS",
					RegionName:   "US_WEST_1",
					Priority:     new(7),
				},
				{
					ProviderName: "AWS",
					RegionName:   "EU_WEST_1",
					Priority:     new(5),
				},
			},
			expectedOrder: []int{7, 5, 0},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a cluster with the region configs
			cluster := &Cluster{
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: tt.regionConfigs,
						},
					},
				},
			}

			normalizeClusterDeployment(cluster)

			// Verify the ordering
			actualOrder := make([]int, len(cluster.ReplicationSpecs[0].RegionConfigs))
			for i, config := range cluster.ReplicationSpecs[0].RegionConfigs {
				if config.Priority != nil {
					actualOrder[i] = *config.Priority
				} else {
					actualOrder[i] = 0
				}
			}

			assert.Equal(t, tt.expectedOrder, actualOrder, "Region configs should be ordered by priority in descending order")
		})
	}
}

func TestConnSet(t *testing.T) {
	testCases := []struct {
		title    string
		inputs   [][]Connection
		expected []Connection
	}{
		{
			title: "Disjoint lists concatenate",
			inputs: [][]Connection{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "D"}, {Name: "E"}, {Name: "F"}},
			},
			expected: []Connection{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"}, {Name: "E"}, {Name: "F"},
			},
		},
		{
			title: "Common items get merged away",
			inputs: [][]Connection{
				{{Name: "A"}, {Name: "B"}, {Name: "C"}},
				{{Name: "B"}, {Name: "C"}, {Name: "D"}},
			},
			expected: []Connection{
				{Name: "A"}, {Name: "B"}, {Name: "C"}, {Name: "D"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := connectionSet(tc.inputs...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRoundtrip_ManagedNamespace(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &admin.GeoSharding20240805{}
		f.Fuzz(fuzzed)
		fuzzed.CustomZoneMapping = nil
		t.Log(fuzzed.ManagedNamespaces)

		fromAtlasResult := managedNamespacesFromAtlas(fuzzed)
		for i, r := range fromAtlasResult {
			toAtlasResult := managedNamespaceToAtlas(&r)
			ns := fuzzed.GetManagedNamespaces()[i]
			equals := (ns.GetCollection() == toAtlasResult.GetCollection() &&
				ns.GetCustomShardKey() == toAtlasResult.GetCustomShardKey() &&
				ns.GetDb() == toAtlasResult.GetDb() &&
				ns.GetIsCustomShardKeyHashed() == toAtlasResult.GetIsCustomShardKeyHashed() &&
				ns.GetIsShardKeyUnique() == toAtlasResult.GetIsShardKeyUnique() &&
				ns.GetNumInitialChunks() == toAtlasResult.GetNumInitialChunks() &&
				ns.GetPresplitHashedZones() == toAtlasResult.GetPresplitHashedZones())
			if !equals {
				t.Log(cmp.Diff(fuzzed.GetManagedNamespaces()[i], *toAtlasResult))
			}
			require.True(t, equals)
		}
	}
}

func TestRoundtrip_CustomZone(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &admin.GeoSharding20240805{}
		f.Fuzz(fuzzed)
		fuzzed.ManagedNamespaces = nil

		fromAtlasResult := customZonesFromAtlas(fuzzed)
		toAtlasResult := customZonesToAtlas(fromAtlasResult)

		require.Equal(t, len(fuzzed.GetCustomZoneMapping()), len(toAtlasResult.GetCustomZoneMappings()))

		for _, r := range toAtlasResult.GetCustomZoneMappings() {
			equals := fuzzed.GetCustomZoneMapping()[r.Location] == r.Zone
			if !equals {
				t.Log(cmp.Diff(fuzzed.GetCustomZoneMapping()[r.Location], r.Zone))
			}
			require.True(t, equals)
		}
	}
}

func TestDeprecated(t *testing.T) {
	for _, tc := range []struct {
		name           string
		deployment     *akov2.AtlasDeployment
		wantDeprecated bool
		wantReason     string
		wantMsg        string
	}{
		{
			name: "nil replication specs",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: nil,
					},
				},
			},
		},
		{
			name: "nil replication spec entries",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							nil, nil, nil,
						},
					},
				},
			},
		},
		{
			name: "nil region configs",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "nil region config entries",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									nil, nil, nil,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "nil regionconfig specs",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										AnalyticsSpecs: nil,
										ElectableSpecs: nil,
										ReadOnlySpecs:  nil,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "empty regionconfig specs",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										AnalyticsSpecs: &akov2.Specs{},
										ElectableSpecs: &akov2.Specs{},
										ReadOnlySpecs:  &akov2.Specs{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "non deprecated M10 instance",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										AnalyticsSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    new(1),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "deprecated M2 instance",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										AnalyticsSpecs: &akov2.Specs{
											InstanceSize: "M2",
											NodeCount:    new(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M2",
											NodeCount:    new(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M2",
											NodeCount:    new(1),
										},
									},
								},
							},
						},
					},
				},
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_DEPRECATION,
			wantMsg:        "WARNING: M2 and M5 instance sizes are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.",
		},
		{
			name: "deprecated M2 instance",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										AnalyticsSpecs: &akov2.Specs{
											InstanceSize: "M5",
											NodeCount:    new(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M5",
											NodeCount:    new(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M5",
											NodeCount:    new(1),
										},
									},
								},
							},
						},
					},
				},
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_DEPRECATION,
			wantMsg:        "WARNING: M2 and M5 instance sizes are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.",
		},
		{
			name: "default read concern set",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
					ProcessArgs: &akov2.ProcessArgs{
						DefaultReadConcern: "true",
					},
				},
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_DEPRECATION,
			wantMsg:        "Process Arg DefaultReadConcern is no longer available in Atlas. Setting this will have no effect.",
		},
		{
			name: "fail index key too long set",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
					ProcessArgs: &akov2.ProcessArgs{
						FailIndexKeyTooLong: new(true),
					},
				},
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_DEPRECATION,
			wantMsg:        "Process Arg FailIndexKeyTooLong is no longer available in Atlas. Setting this will have no effect.",
		},
		{
			name: "empty serverless instance",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{},
				},
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_DEPRECATION,
			wantMsg:        "WARNING: Serverless is deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.",
		},
		{
			name: "remove upgrade flag",
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					UpgradeToDedicated: true,
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
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
			},
			wantDeprecated: true,
			wantReason:     NOTIFICATION_REASON_RECOMMENDATION,
			wantMsg:        "Cluster is already dedicated. It’s recommended to remove or set the upgrade flag to false",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDeployment("123", tc.deployment)
			gotDeprecated, gotReason, gotMsg := d.Notifications()
			require.Equal(t, tc.wantDeprecated, gotDeprecated)
			require.Equal(t, tc.wantReason, gotReason)
			require.Equal(t, tc.wantMsg, gotMsg)
		})
	}
}

func TestIsType(t *testing.T) {
	tests := map[string]struct {
		deployment     *akov2.AtlasDeployment
		wantServerless bool
		wantFlex       bool
		wantTenant     bool
		wantDedicated  bool
	}{
		"Cluster is serverless": {
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{},
				},
			},
			wantServerless: false,
			wantFlex:       true,
			wantTenant:     false,
			wantDedicated:  false,
		},
		"Cluster is flex": {
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					FlexSpec: &akov2.FlexSpec{},
				},
			},
			wantServerless: false,
			wantFlex:       true,
			wantTenant:     false,
			wantDedicated:  false,
		},
		"Cluster is tenant": {
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName:        "TENANT",
										BackingProviderName: "AWS",
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantServerless: false,
			wantFlex:       false,
			wantTenant:     true,
			wantDedicated:  false,
		},
		"Cluster is dedicated": {
			deployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
										},
									},
								},
							},
						},
					},
				},
			},
			wantServerless: false,
			wantFlex:       false,
			wantTenant:     false,
			wantDedicated:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := NewDeployment("123", tt.deployment)
			assert.Equal(t, tt.wantServerless, d.IsServerless())
			assert.Equal(t, tt.wantFlex, d.IsFlex())
			assert.Equal(t, tt.wantTenant, d.IsTenant())
			assert.Equal(t, tt.wantDedicated, d.IsDedicated())
		})
	}
}

// Regression test for https://github.com/mongodb/mongodb-atlas-kubernetes/issues/3142
// terminationProtectionEnabled=false must be sent to Atlas as an explicit false,
// not omitted. Otherwise, when Atlas has it true (e.g. set via UI), the PATCH can
// never disable it and the controller reconciles forever.
func TestClusterUpdateToAtlas_TerminationProtectionFalseIsSent(t *testing.T) {
	cluster := &Cluster{
		AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:                         "cluster0",
			ClusterType:                  "REPLICASET",
			TerminationProtectionEnabled: false,
		},
	}

	req := clusterUpdateToAtlas(cluster)

	require.NotNil(t, req.TerminationProtectionEnabled,
		"expected TerminationProtectionEnabled to be sent as explicit false, got nil (field would be omitted from PATCH)")
	assert.False(t, *req.TerminationProtectionEnabled,
		"expected TerminationProtectionEnabled to be false")
}

func TestClusterCreateToAtlas_TerminationProtectionFalseIsSent(t *testing.T) {
	cluster := &Cluster{
		AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:                         "cluster0",
			ClusterType:                  "REPLICASET",
			TerminationProtectionEnabled: false,
		},
	}

	req := clusterCreateToAtlas(cluster)

	require.NotNil(t, req.TerminationProtectionEnabled,
		"expected TerminationProtectionEnabled to be sent as explicit false, got nil")
	assert.False(t, *req.TerminationProtectionEnabled)
}

// Regression test for https://github.com/mongodb/mongodb-atlas-kubernetes/issues/3142
// When the user leaves processArgs fields unset in the CR, Atlas returns its own
// defaults (e.g. DefaultWriteConcern="majority", OplogSizeMB=990). The reconciler's
// equality check (reflect.DeepEqual on *akov2.ProcessArgs) must treat those
// CR-unset fields as "no opinion" — otherwise the round-trip loops forever:
//   - DeepEqual says desired != current
//   - processArgsToAtlas drops zero values via MakePtrOrNil → PATCH body is empty
//   - Atlas keeps its defaults → next reconcile reports the same diff → UPDATE
//     is re-issued endlessly.
//
// This test pins the symmetry: a pair where AKO has an unset field and Atlas
// has a server default should not require an update request. The fix may be
// either in the compare step or in how processArgsFromAtlas normalizes defaults.
func TestProcessArgs_UnsetFieldShouldNotDivergeFromAtlasDefaults(t *testing.T) {
	// What the user's CR yields after normalization: mostly defaults, no
	// opinion on DefaultWriteConcern or OplogSizeMB.
	akoArgs := &akov2.ProcessArgs{}
	normalizeProcessArgs(akoArgs)

	// What Atlas returns for an unconfigured cluster: populated server defaults.
	atlasArgs := processArgsFromAtlas(&admin.ClusterDescriptionProcessArgs20240805{
		DefaultWriteConcern:       pointer.MakePtr("majority"),
		MinimumEnabledTlsProtocol: pointer.MakePtr("TLS1_2"),
		JavascriptEnabled:         pointer.MakePtr(true),
		NoTableScan:               pointer.MakePtr(false),
		OplogSizeMB:               pointer.MakePtr(990),
	})

	// Drives the controller's UpdateProcessArgs call in
	// internal/controller/atlasdeployment/advanced_deployment.go.
	require.True(t,
		ProcessArgsEqual(akoArgs, atlasArgs),
		"expected AKO processArgs (CR-unset defaults) to be considered equal to "+
			"Atlas processArgs (server defaults), otherwise the controller will "+
			"call UpdateProcessArgs every reconcile — and the PATCH body omits "+
			"zero values, so Atlas never converges. Got:\n  ako:   %+v\n  atlas: %+v",
		akoArgs, atlasArgs)
}
