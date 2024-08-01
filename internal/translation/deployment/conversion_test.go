package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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
