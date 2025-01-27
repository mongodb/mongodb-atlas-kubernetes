package deployment

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestNewDeployment(t *testing.T) {
	tests := map[string]struct {
		cr       *akov2.AtlasDeployment
		expected Deployment
	}{
		"should create a new serverless deployment": {
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
			expected: &Serverless{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
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
								},
							},
						},
					},
					TerminationProtectionEnabled: true,
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

func TestNormalizeServerlessDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *Serverless
		expected   *Serverless
	}{
		"normalize deployment without tags": {
			deployment: &Serverless{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: true,
					},
					TerminationProtectionEnabled: true,
				},
			},
			expected: &Serverless{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					Tags: []*akov2.TagSpec{},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: true,
					},
					TerminationProtectionEnabled: true,
				},
			},
		},
		"normalize deployment with tags": {
			deployment: &Serverless{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
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
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: true,
					},
					TerminationProtectionEnabled: true,
				},
			},
			expected: &Serverless{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
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
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: true,
					},
					TerminationProtectionEnabled: true,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			normalizeServerlessDeployment(tt.deployment)

			assert.Equal(t, tt.expected, tt.deployment)
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
					MongoDBMajorVersion:      "7.0",
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
		fuzzed := &admin.GeoSharding{}
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
		fuzzed := &admin.GeoSharding{}
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
