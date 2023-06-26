package validate

import (
	"fmt"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: nil}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("Instance size not empty when serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "M10",
					ProviderName:     "SERVERLESS",
				},
			}}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("Instance size unset when not serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "",
					ProviderName:     "AWS",
				},
			}}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("different instance sizes for advanced deployment", func(t *testing.T) {
			t.Run("different instance size in the same region", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M20"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(spec))
			})
			t.Run("different instance size in different regions", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M20"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(spec))
			})
			t.Run("different instance size in different replications", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
								},
							},
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M20"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(spec))
			})
		})
		t.Run("different autoscaling for advanced deployment", func(t *testing.T) {
			t.Run("different instance size in different regions", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
											Compute: &mdbv1.ComputeSpec{
												Enabled:          toptr.MakePtr(true),
												ScaleDownEnabled: toptr.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M30",
											},
										},
									},
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(spec))
			})
			t.Run("different autoscaling in different replications", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
										AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
											Compute: &mdbv1.ComputeSpec{
												Enabled:          toptr.MakePtr(true),
												ScaleDownEnabled: toptr.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M30",
											},
										},
									},
								},
							},
							{
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{InstanceSize: "M20"},
										ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(spec))
			})
		})
	})
	t.Run("Valid cluster specs", func(t *testing.T) {
		t.Run("Advanced cluster spec specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: nil}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})
		t.Run("Regular cluster specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})

		t.Run("Serverless Cluster", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					ProviderName: "SERVERLESS",
				},
			}}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})
		t.Run("Advanced cluster with replication config", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{
				AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
					ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
						{
							RegionConfigs: []*mdbv1.AdvancedRegionConfig{
								{
									ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
									AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
										Compute: &mdbv1.ComputeSpec{
											Enabled:          toptr.MakePtr(true),
											ScaleDownEnabled: toptr.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M30",
										},
									},
								},
							},
						},
						{
							RegionConfigs: []*mdbv1.AdvancedRegionConfig{
								{
									ElectableSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									ReadOnlySpecs:  &mdbv1.Specs{InstanceSize: "M10"},
									AnalyticsSpecs: &mdbv1.Specs{InstanceSize: "M10"},
									AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
										Compute: &mdbv1.ComputeSpec{
											Enabled:          toptr.MakePtr(true),
											ScaleDownEnabled: toptr.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M30",
										},
									},
								},
							},
						},
					},
				},
			}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})
	})
}

func TestProjectValidation(t *testing.T) {
	t.Run("custom roles spec", func(t *testing.T) {
		t.Run("empty custom roles spec", func(t *testing.T) {
			spec := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{},
			}
			assert.NoError(t, Project(spec))
		})
		t.Run("valid custom roles spec", func(t *testing.T) {
			spec := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					CustomRoles: []mdbv1.CustomRole{
						{
							Name: "cr-1",
						},
						{
							Name: "cr-2",
						},
						{
							Name: "cr-3",
						},
					},
				},
			}
			assert.NoError(t, Project(spec))
		})
		t.Run("invalid custom roles spec", func(t *testing.T) {
			spec := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{
					CustomRoles: []mdbv1.CustomRole{
						{
							Name: "cr-1",
						},
						{
							Name: "cr-1",
						},
						{
							Name: "cr-1",
						},
						{
							Name: "cr-2",
						},
						{
							Name: "cr-2",
						},
					},
				},
			}
			assert.Error(t, Project(spec))
		})
	})
}

func TestBackupScheduleValidation(t *testing.T) {
	t.Run("auto export is enabled without export policy", func(t *testing.T) {
		bSchedule := &mdbv1.AtlasBackupSchedule{
			Spec: mdbv1.AtlasBackupScheduleSpec{
				AutoExportEnabled: true,
			},
		}
		deployment := &mdbv1.AtlasDeployment{
			Status: status.AtlasDeploymentStatus{},
		}
		assert.Error(t, BackupSchedule(bSchedule, deployment))
	})

	t.Run("copy settings on advanced deployment", func(t *testing.T) {
		t.Run("copy settings is valid", func(t *testing.T) {
			bSchedule := &mdbv1.AtlasBackupSchedule{
				Spec: mdbv1.AtlasBackupScheduleSpec{
					CopySettings: []mdbv1.CopySetting{
						{
							RegionName:        toptr.MakePtr("US_WEST_1"),
							ReplicationSpecID: toptr.MakePtr("123"),
							CloudProvider:     toptr.MakePtr("AWS"),
							ShouldCopyOplogs:  toptr.MakePtr(true),
							Frequencies:       []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						PitEnabled: toptr.MakePtr(true),
					},
				},
				Status: status.AtlasDeploymentStatus{
					ReplicaSets: []status.ReplicaSet{
						{
							ID:       "123",
							ZoneName: "Zone 1",
						},
					},
				},
			}
			assert.NoError(t, BackupSchedule(bSchedule, deployment))
		})

		t.Run("copy settings is invalid", func(t *testing.T) {
			bSchedule := &mdbv1.AtlasBackupSchedule{
				Spec: mdbv1.AtlasBackupScheduleSpec{
					CopySettings: []mdbv1.CopySetting{
						{
							ShouldCopyOplogs: toptr.MakePtr(true),
						},
						{
							RegionName:        toptr.MakePtr("US_WEST_1"),
							ReplicationSpecID: toptr.MakePtr("123"),
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{},
				},
				Status: status.AtlasDeploymentStatus{
					ReplicaSets: []status.ReplicaSet{
						{
							ID:       "321",
							ZoneName: "Zone 1",
						},
					},
				},
			}
			assert.Error(t, BackupSchedule(bSchedule, deployment))
		})
	})

	t.Run("copy settings on legacy deployment", func(t *testing.T) {
		t.Run("copy settings is valid", func(t *testing.T) {
			bSchedule := &mdbv1.AtlasBackupSchedule{
				Spec: mdbv1.AtlasBackupScheduleSpec{
					CopySettings: []mdbv1.CopySetting{
						{
							RegionName:        toptr.MakePtr("US_WEST_1"),
							ReplicationSpecID: toptr.MakePtr("123"),
							CloudProvider:     toptr.MakePtr("AWS"),
							ShouldCopyOplogs:  toptr.MakePtr(true),
							Frequencies:       []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.DeploymentSpec{
						PitEnabled: toptr.MakePtr(true),
					},
				},
				Status: status.AtlasDeploymentStatus{
					ReplicaSets: []status.ReplicaSet{
						{
							ID:       "123",
							ZoneName: "Zone 1",
						},
					},
				},
			}
			assert.NoError(t, BackupSchedule(bSchedule, deployment))
		})

		t.Run("copy settings is invalid", func(t *testing.T) {
			bSchedule := &mdbv1.AtlasBackupSchedule{
				Spec: mdbv1.AtlasBackupScheduleSpec{
					CopySettings: []mdbv1.CopySetting{
						{
							ShouldCopyOplogs: toptr.MakePtr(true),
						},
						{
							RegionName:        toptr.MakePtr("US_WEST_1"),
							ReplicationSpecID: toptr.MakePtr("123"),
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.DeploymentSpec{},
				},
				Status: status.AtlasDeploymentStatus{
					ReplicaSets: []status.ReplicaSet{
						{
							ID:       "321",
							ZoneName: "Zone 1",
						},
					},
				},
			}
			assert.Error(t, BackupSchedule(bSchedule, deployment))
		})
	})
}

// this key was removed immediately after download, so don't bother
const sampleSAKey = `{
	"type": "service_account",
	"project_id": "some-project",
	"private_key_id": "dc6c401f0acd0147ca70e3169f579b570583b58f",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCviL4+Bnn759sV\nNrtyexwHtR/5JYzwivupqOMdz1zuscqKfJOo6RXzq7Em3saoxLpHwwJzPq+HX1D+\ndaE0fB2hO0Mfjkcgmro0jBbLnRJgFCx7NwkyDfj8z6i4zx2CxZLiwdqo3Q3hEMNy\n5oZs52tFlTkOALCxa76Aq4eUIblupyhlhETQB1fb4D9+U57b4eeeFRyccwR7Cg0x\nfZUE1udV7YwKphLS8dCbLoqAQ3jmaxv+Qjo8e5Mj5oaMxzRAdOO1VtG9GaLU96Rc\nKGS75w+E/eR3Pm7b2dFo3jba2jvgm3U++EJi5/0zL/TQnOMwNHASPUsYQAhZezB5\nh5/MDwmjAgMBAAECggEAIOccZer33Zipz9GpFD3wVJ+GZUC9KO+cWcJ/A/z1Ggb4\nhLnyQbSjOUAjHjqe+U6a7k2m/WwwIctjlrb85yYmtayymc0lFv75zVS/Bx6jrZ/K\ncLQxxJCq7dSM90tXaEZZkKiusH1zFw9522VLqEk+qdXdUnsdo7wjAuJkMQebRxq8\n1lp3UGqAraBXLRrYUnQwBRezSYh93nZ+u+etjRCBjMYoy06PGDrJG05/FFNgQy1y\nGkBVKnmf9WNyDuX1ePyoXCAvRkwUW5ixTNGoGrRwbwleaFTvYPBlPM+TyCw4rdnU\nzRVNpoCqkf6S9tXjHKmGgwkyMqmYMCOk/5xCSb2p4QKBgQD14xcXIf9lnyHJCllU\nQdUIAmIp9/91Rdjpf/y98LCNrD/cyTvVr0+xSnN9ksBEGwvTIaaBcvdCMNBG8sI9\nsnO8W8GG1Bs80D3XIbJFaGmTmOngvrfie3tbP77wfcPgn659Q0I1+bNZGNVX+WEM\nn+f7rfGPa8Br/zHpQf2gaGWHAwKBgQC2wOrInJ+ndpqU49JVaT6YtQj0OKeLhSpc\n0N5DdW0jjD0WrnhOsLUMPC5V8R5fo4tFPfXVIfE2k8J5xxgopJzGLlHmhxq0ltmF\nbSoM2uHKf0UiRFmVTwZmzDwn+Ym/H9J/6L7Gd86u8kmWfJYFa3OzqJTvw7e+k4kD\nITb/NlEg4QKBgQCfW4AZg/Ur/Ug+LTDbxJa2TCUmog20CYKdQk+hIh6qktoI03qt\n8KKrel8DIVruSMEPIp3xA3twMIarlKWCqucLSkRQh6LndOa/SJ1rElJqUA4zlCdE\n51Z5OwUag8exCoxhrnd4183+jnOmQn89WV1V5dPKacEZvRix3gzsKvyx1QKBgFsH\nlOsAOPYtOapYIHiyx59A7YjYf3wbhJJe55cqcoZ2YCdgGET59/R0NZBRXhO9Xq3K\nwxy6n2/UAdauuPXlqMF+aQUu3rp9OTQgwAVPMZCv/DupWAXrKwEhUgWHYnl03GEi\nCYTKQIUb4lO3EvL4JtWiby1Oi8O9sU2ByectoxOBAoGASm9BXSN8Ru+dP5/E55mb\nd//aQlxlIvROhWSnotGzhyQ6DVk2fRRZQAuTEFVEprBX87gckdzb5cdaomziO9be\nhmtv1ValgmOCnta2AYw1blvfGK7B5FEpFckMniLjWap08aironIImj6ligLWDqc0\nNbdyAvc6N/5qG8gu4f8C2Q4=\n-----END PRIVATE KEY-----\n",
	"client_email": "619108922856-compute@developer.gserviceaccount.com",
	"client_id": "117865750705662546099",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-compute%40developer.gserviceaccount.com",
	"universe_domain": "googleapis.com"
}`

const sampleSAKeyOneLine = `{	"type": "service_account",	"project_id": "some-project",	"private_key_id": "dc6c401f0acd0147ca70e3169f579b570583b58f",	"private_key": "-----BEGIN PRIVATE KEY-----\\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCviL4+Bnn759sV\\nNrtyexwHtR/5JYzwivupqOMdz1zuscqKfJOo6RXzq7Em3saoxLpHwwJzPq+HX1D+\\ndaE0fB2hO0Mfjkcgmro0jBbLnRJgFCx7NwkyDfj8z6i4zx2CxZLiwdqo3Q3hEMNy\\n5oZs52tFlTkOALCxa76Aq4eUIblupyhlhETQB1fb4D9+U57b4eeeFRyccwR7Cg0x\\nfZUE1udV7YwKphLS8dCbLoqAQ3jmaxv+Qjo8e5Mj5oaMxzRAdOO1VtG9GaLU96Rc\\nKGS75w+E/eR3Pm7b2dFo3jba2jvgm3U++EJi5/0zL/TQnOMwNHASPUsYQAhZezB5\\nh5/MDwmjAgMBAAECggEAIOccZer33Zipz9GpFD3wVJ+GZUC9KO+cWcJ/A/z1Ggb4\\nhLnyQbSjOUAjHjqe+U6a7k2m/WwwIctjlrb85yYmtayymc0lFv75zVS/Bx6jrZ/K\\ncLQxxJCq7dSM90tXaEZZkKiusH1zFw9522VLqEk+qdXdUnsdo7wjAuJkMQebRxq8\\n1lp3UGqAraBXLRrYUnQwBRezSYh93nZ+u+etjRCBjMYoy06PGDrJG05/FFNgQy1y\\nGkBVKnmf9WNyDuX1ePyoXCAvRkwUW5ixTNGoGrRwbwleaFTvYPBlPM+TyCw4rdnU\\nzRVNpoCqkf6S9tXjHKmGgwkyMqmYMCOk/5xCSb2p4QKBgQD14xcXIf9lnyHJCllU\\nQdUIAmIp9/91Rdjpf/y98LCNrD/cyTvVr0+xSnN9ksBEGwvTIaaBcvdCMNBG8sI9\\nsnO8W8GG1Bs80D3XIbJFaGmTmOngvrfie3tbP77wfcPgn659Q0I1+bNZGNVX+WEM\\nn+f7rfGPa8Br/zHpQf2gaGWHAwKBgQC2wOrInJ+ndpqU49JVaT6YtQj0OKeLhSpc\\n0N5DdW0jjD0WrnhOsLUMPC5V8R5fo4tFPfXVIfE2k8J5xxgopJzGLlHmhxq0ltmF\\nbSoM2uHKf0UiRFmVTwZmzDwn+Ym/H9J/6L7Gd86u8kmWfJYFa3OzqJTvw7e+k4kD\\nITb/NlEg4QKBgQCfW4AZg/Ur/Ug+LTDbxJa2TCUmog20CYKdQk+hIh6qktoI03qt\\n8KKrel8DIVruSMEPIp3xA3twMIarlKWCqucLSkRQh6LndOa/SJ1rElJqUA4zlCdE\\n51Z5OwUag8exCoxhrnd4183+jnOmQn89WV1V5dPKacEZvRix3gzsKvyx1QKBgFsH\\nlOsAOPYtOapYIHiyx59A7YjYf3wbhJJe55cqcoZ2YCdgGET59/R0NZBRXhO9Xq3K\\nwxy6n2/UAdauuPXlqMF+aQUu3rp9OTQgwAVPMZCv/DupWAXrKwEhUgWHYnl03GEi\\nCYTKQIUb4lO3EvL4JtWiby1Oi8O9sU2ByectoxOBAoGASm9BXSN8Ru+dP5/E55mb\\nd//aQlxlIvROhWSnotGzhyQ6DVk2fRRZQAuTEFVEprBX87gckdzb5cdaomziO9be\\nhmtv1ValgmOCnta2AYw1blvfGK7B5FEpFckMniLjWap08aironIImj6ligLWDqc0\\nNbdyAvc6N/5qG8gu4f8C2Q4=\\n-----END PRIVATE KEY-----\\n",	"client_email": "619108922856-compute@developer.gserviceaccount.com",	"client_id": "117865750705662546099",	"auth_uri": "https://accounts.google.com/o/oauth2/auth",	"token_uri": "https://oauth2.googleapis.com/token",	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-compute%40developer.gserviceaccount.com",	"universe_domain": "googleapis.com"}`

func testEncryptionAtRest(enabled bool) *mdbv1.EncryptionAtRest {
	flag := enabled
	return &mdbv1.EncryptionAtRest{
		GoogleCloudKms: mdbv1.GoogleCloudKms{
			Enabled:           &flag,
			ServiceAccountKey: sampleSAKey,
		},
	}
}

func withProperUrls(properties string) string {
	urls := `"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-compute%40developer.gserviceaccount.com"`
	return fmt.Sprintf(`{%s, %s}`, urls, properties)
}

func TestEncryptionAtRestValidation(t *testing.T) {
	t.Run("google service account key validation succeeds if no encryption at rest is used", func(t *testing.T) {
		assert.NoError(t, encryptionAtRest(&mdbv1.EncryptionAtRest{}))
	})

	t.Run("google service account key validation succeeds if encryption at rest is disabled", func(t *testing.T) {
		assert.NoError(t, encryptionAtRest(testEncryptionAtRest(false)))
	})

	t.Run("google service account key validation succeeds if encryption is enabled but the key is empty", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = ""
		assert.ErrorContains(t, encryptionAtRest(enc), "missing Google Service Account Key but GCP KMS is enabled")
	})

	t.Run("google service account key validation succeeds for a good key", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = sampleSAKey
		assert.NoError(t, encryptionAtRest(enc))
	})

	t.Run("google service account key validation succeeds for a good key in a single line", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = sampleSAKeyOneLine
		assert.NoError(t, encryptionAtRest(enc))
	})

	t.Run("google service account key validation fails for an empty json key", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = "{}"
		assert.ErrorContains(t, encryptionAtRest(enc), "invalid empty service account key")
	})

	t.Run("google service account key validation fails for an empty array json as key", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = "[]"
		assert.ErrorContains(t, encryptionAtRest(enc), "cannot unmarshal array into Go value")
	})

	t.Run("google service account key validation fails for a json object with a wrong field type", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = `{"type":true}`
		assert.ErrorContains(t, encryptionAtRest(enc), "cannot unmarshal bool")
	})

	t.Run("google service account key validation fails for a bad pem key", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = withProperUrls(`"private_key":"-----BEGIN PRIVATE KEY-----\nMIIEvQblah\n-----END PRIVATE KEY-----\n"`)
		assert.ErrorContains(t, encryptionAtRest(enc), "failed to decode PEM")
	})

	t.Run("google service account key validation fails for a bad URL", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = withProperUrls(`"token_uri": "http//badurl.example"`)
		assert.ErrorContains(t, encryptionAtRest(enc), "invalid URL address")
	})
}
