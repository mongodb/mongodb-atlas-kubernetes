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
	"go.mongodb.org/atlas-sdk/v20250312011/admin"

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
						BackupEnabled: pointer.MakePtr(true),
						BiConnector: &akov2.BiConnectorSpec{
							Enabled:        pointer.MakePtr(true),
							ReadPreference: "secondary",
						},
						DiskSizeGB:               pointer.MakePtr(20),
						EncryptionAtRestProvider: "AWS",
						Paused:                   pointer.MakePtr(false),
						PitEnabled:               pointer.MakePtr(true),
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
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
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
						JavascriptEnabled:         pointer.MakePtr(true),
					},
				},
			},
			expected: &Cluster{
				ProjectID: "project-id",
				//nolint:dupl
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:          "cluster0",
					ClusterType:   "REPLICASET",
					BackupEnabled: pointer.MakePtr(true),
					BiConnector: &akov2.BiConnectorSpec{
						Enabled:        pointer.MakePtr(true),
						ReadPreference: "secondary",
					},
					DiskSizeGB:               pointer.MakePtr(20),
					EncryptionAtRestProvider: "AWS",
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(true),
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
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
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
					JavascriptEnabled:         pointer.MakePtr(true),
					NoTableScan:               pointer.MakePtr(false),
				},
				//nolint:dupl
				customResource: &akov2.AtlasDeployment{
					Spec: akov2.AtlasDeploymentSpec{
						//nolint:dupl
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{
							Name:          "cluster0",
							ClusterType:   "REPLICASET",
							BackupEnabled: pointer.MakePtr(true),
							BiConnector: &akov2.BiConnectorSpec{
								Enabled:        pointer.MakePtr(true),
								ReadPreference: "secondary",
							},
							DiskSizeGB:               pointer.MakePtr(20),
							EncryptionAtRestProvider: "AWS",
							Paused:                   pointer.MakePtr(false),
							PitEnabled:               pointer.MakePtr(true),
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
											Priority:     pointer.MakePtr(7),
											ElectableSpecs: &akov2.Specs{
												InstanceSize: "M10",
												NodeCount:    pointer.MakePtr(3),
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
							JavascriptEnabled:         pointer.MakePtr(true),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
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
											Enabled:          pointer.MakePtr(false),
											ScaleDownEnabled: pointer.MakePtr(false),
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
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					BackupEnabled:            pointer.MakePtr(false),
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
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
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
									Priority:            pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
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
					BackupEnabled:            pointer.MakePtr(false),
					VersionReleaseSystem:     "LTS",
					EncryptionAtRestProvider: "NONE",
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
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
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
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
											NodeCount:    pointer.MakePtr(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(1),
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
											NodeCount:    pointer.MakePtr(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M2",
											NodeCount:    pointer.MakePtr(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M2",
											NodeCount:    pointer.MakePtr(1),
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
											NodeCount:    pointer.MakePtr(1),
										},
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M5",
											NodeCount:    pointer.MakePtr(1),
										},
										ReadOnlySpecs: &akov2.Specs{
											InstanceSize: "M5",
											NodeCount:    pointer.MakePtr(1),
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
						FailIndexKeyTooLong: pointer.MakePtr(true),
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
