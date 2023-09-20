package validate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: nil}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("Instance size not empty when serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "M10",
					ProviderName:     "SERVERLESS",
				},
			}}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("Instance size unset when not serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "",
					ProviderName:     "AWS",
				},
			}}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
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
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
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
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
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
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
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
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
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
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
			})
		})
	})
	t.Run("Valid cluster specs", func(t *testing.T) {
		t.Run("Advanced cluster spec specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: nil}
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("Regular cluster specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
		})

		t.Run("Serverless Cluster", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					ProviderName: "SERVERLESS",
				},
			}}
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
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
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
		})
	})
}

func TestDeploymentForGov(t *testing.T) {
	t.Run("should fail when deployment is configured to non-gov region", func(t *testing.T) {
		deploy := mdbv1.AtlasDeploymentSpec{
			DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					RegionName: "EU_EAST_1",
				},
			},
		}

		assert.ErrorContains(t, deploymentForGov(&deploy, "GOV_REGIONS_ONLY"), "deployment in atlas for government support a restricted set of regions: EU_EAST_1 is not part of AWS for government regions")
	})

	t.Run("should fail when advanced deployment is configured to non-gov region", func(t *testing.T) {
		deploy := mdbv1.AtlasDeploymentSpec{
			AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
				ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mdbv1.AdvancedRegionConfig{
							{
								RegionName: "EU_EAST_1",
							},
						},
					},
				},
			},
		}

		assert.ErrorContains(t, deploymentForGov(&deploy, "COMMERCIAL_FEDRAMP_REGIONS_ONLY"), "advanced deployment in atlas for government support a restricted set of regions: EU_EAST_1 is not part of AWS FedRAMP regions")
	})
}

func TestProjectValidation(t *testing.T) {
	t.Run("should fail when commercial Atlas sets region restriction field to GOV_REGIONS_ONLY", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
			},
		}

		assert.ErrorContains(t, Project(akoProject, false), "regionUsageRestriction can be used only with Atlas for government")
	})

	t.Run("should fail when commercial Atlas sets region restriction field to COMMERCIAL_FEDRAMP_REGIONS_ONLY", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "COMMERCIAL_FEDRAMP_REGIONS_ONLY",
			},
		}

		assert.ErrorContains(t, Project(akoProject, false), "regionUsageRestriction can be used only with Atlas for government")
	})

	t.Run("should not fail if commercial Atlas sets region restriction field to empty", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{},
		}

		assert.NoError(t, Project(akoProject, false))
	})

	t.Run("should not fail if commercial Atlas sets region restriction field to NONE", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "NONE",
			},
		}

		assert.NoError(t, Project(akoProject, false))
	})

	t.Run("custom roles spec", func(t *testing.T) {
		t.Run("empty custom roles spec", func(t *testing.T) {
			spec := &mdbv1.AtlasProject{
				Spec: mdbv1.AtlasProjectSpec{},
			}
			assert.NoError(t, Project(spec, false))
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
			assert.NoError(t, Project(spec, false))
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
			assert.Error(t, Project(spec, false))
		})
	})
}

func TestProjectForGov(t *testing.T) {
	t.Run("should fail if there's non AWS network peering config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []mdbv1.NetworkPeer{
					{
						ProviderName:        "GCP",
						AccepterRegionName:  "europe-west-1",
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/18",
						NetworkName:         "my-gcp-peer",
						GCPProjectID:        "my-gcp-project",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. one or more network peers are not set to AWS")
	})

	t.Run("should fail if there's no gov region in network peering config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []mdbv1.NetworkPeer{
					{
						ProviderName:        "AWS",
						AccepterRegionName:  "us-east-1",
						ContainerRegion:     "us-east-1",
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "network peering in atlas for government support a restricted set of regions: us-east-1 is not part of AWS for government regions")
	})

	t.Run("should fail if there's a GCP encryption at rest config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					GoogleCloudKms: mdbv1.GoogleCloudKms{
						Enabled: toptr.MakePtr(true),
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. disable encryption at rest for Google Cloud")
	})

	t.Run("should fail if there's a Azure encryption at rest config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AzureKeyVault: mdbv1.AzureKeyVault{
						Enabled: toptr.MakePtr(true),
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. disable encryption at rest for Azure")
	})

	t.Run("should fail if there's a AWS encryption at rest config with wrong region", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "us-east-1",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "encryption at rest in atlas for government support a restricted set of regions: us-east-1 is not part of AWS for government regions")
	})

	t.Run("should fail if there's non AWS private endpoint config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: "GCP",
						Region:   "europe-west-1",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. one or more private endpoints are not set to AWS")
	})

	t.Run("should fail if there's no gov region in private endpoint config", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "COMMERCIAL_FEDRAMP_REGIONS_ONLY",
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: "AWS",
						Region:   "eu-east-1",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "private endpoint in atlas for government support a restricted set of regions: eu-east-1 is not part of AWS FedRAMP regions")
	})

	t.Run("should succeed if resources are properly configured", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []mdbv1.NetworkPeer{
					{
						ProviderName:        "AWS",
						AccepterRegionName:  "us-gov-east-1",
						ContainerRegion:     "us-gov-east-1",
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
				EncryptionAtRest: &mdbv1.EncryptionAtRest{
					AwsKms: mdbv1.AwsKms{
						Enabled: toptr.MakePtr(true),
						Region:  "us-gov-east-1",
					},
				},
				PrivateEndpoints: []mdbv1.PrivateEndpoint{
					{
						Provider: "AWS",
						Region:   "us-gov-east-1",
					},
				},
			},
		}

		assert.NoError(t, Project(akoProject, true))
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

const sampleSAKeyFmt = `{
	"type": "service_account",
	"project_id": "some-project",
	"private_key_id": "dc6c401f0acd0147ca70e3169f579b570583b58f",
	"private_key": "%s",
	"client_email": "619108922856-compute@developer.gserviceaccount.com",
	"client_id": "117865750705662546099",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-computedeveloper.gserviceaccount.com",
	"universe_domain": "googleapis.com"
}`

const sampleSAKeyOneLineFmt = `{	"type": "service_account",	"project_id": "some-project",	"private_key_id": "dc6c401f0acd0147ca70e3169f579b570583b58f",	"private_key": "%s",	"client_email": "619108922856-compute@developer.gserviceaccount.com",	"client_id": "117865750705662546099",	"auth_uri": "https://accounts.google.com/o/oauth2/auth",	"token_uri": "https://oauth2.googleapis.com/token",	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-computedeveloper.gserviceaccount.com",	"universe_domain": "googleapis.com"}`

func testEncryptionAtRest(enabled bool) *mdbv1.EncryptionAtRest {
	flag := enabled
	return &mdbv1.EncryptionAtRest{
		GoogleCloudKms: mdbv1.GoogleCloudKms{
			Enabled:           &flag,
			ServiceAccountKey: sampleSAKey(),
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
		enc.GoogleCloudKms.ServiceAccountKey = sampleSAKey()
		fmt.Printf("g key=%q", enc.GoogleCloudKms.ServiceAccountKey)
		assert.NoError(t, encryptionAtRest(enc))
	})

	t.Run("google service account key validation succeeds for a good key in a single line", func(t *testing.T) {
		enc := testEncryptionAtRest(true)
		enc.GoogleCloudKms.ServiceAccountKey = sampleSAKeyOneLine()
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

func TestProjectIpAccessList(t *testing.T) {
	t.Run("should return no error for empty list", func(t *testing.T) {
		assert.NoError(t, projectIPAccessList([]project.IPAccessList{}))
	})

	t.Run("should return error when multiple ways were configured", func(t *testing.T) {
		data := map[string]struct {
			ipAccessList []project.IPAccessList
			err          string
		}{
			"for CIDRBlock with IPAddress": {
				ipAccessList: []project.IPAccessList{
					{
						IPAddress: "10.0.0.1",
						CIDRBlock: "10.0.0.0/24",
					},
				},
				err: "don't set ipAddress or awsSecurityGroup when configuring cidrBlock",
			},
			"for CIDRBlock with awsSecurityGroup": {
				ipAccessList: []project.IPAccessList{
					{
						AwsSecurityGroup: "sg-0129d834cbf03bc6d",
						CIDRBlock:        "10.0.0.0/24",
					},
				},
				err: "don't set ipAddress or awsSecurityGroup when configuring cidrBlock",
			},
			"for IPAddress with awsSecurityGroup": {
				ipAccessList: []project.IPAccessList{
					{
						AwsSecurityGroup: "sg-0129d834cbf03bc6d",
						IPAddress:        "10.0.0.1",
					},
				},
				err: "don't set cidrBlock or awsSecurityGroup when configuring ipAddress",
			},
		}

		for desc, item := range data {
			t.Run(desc, func(t *testing.T) {
				assert.ErrorContains(t, projectIPAccessList(item.ipAccessList), item.err)
			})
		}
	})

	t.Run("should return error when configuration is invalid", func(t *testing.T) {
		data := map[string]struct {
			ipAccessList []project.IPAccessList
			err          string
		}{
			"for empty config": {
				ipAccessList: []project.IPAccessList{{}},
				err:          "invalid config! one of option must be configured",
			},
			"for CIDRBlock": {
				ipAccessList: []project.IPAccessList{
					{
						CIDRBlock: "10.0.0.0",
					},
				},
				err: "invalid cidrBlock: 10.0.0.0",
			},
			"for IPAddress": {
				ipAccessList: []project.IPAccessList{
					{
						IPAddress: "10.0.0.350",
					},
				},
				err: "invalid ipAddress: 10.0.0.350",
			},
			"for awsSecurityGroup": {
				ipAccessList: []project.IPAccessList{
					{
						AwsSecurityGroup: "invalid0129d834cbf03bc6d",
					},
				},
				err: "invalid awsSecurityGroup: invalid0129d834cbf03bc6d",
			},
			"for DeleteAfterDate": {
				ipAccessList: []project.IPAccessList{
					{
						IPAddress:       "10.0.0.10",
						DeleteAfterDate: "2020-01-02T15:04:05-07000",
					},
				},
				err: "invalid deleteAfterDate: 2020-01-02T15:04:05-07000. value should follow ISO8601 format",
			},
		}

		for desc, item := range data {
			t.Run(desc, func(t *testing.T) {
				assert.ErrorContains(t, projectIPAccessList(item.ipAccessList), item.err)
			})
		}
	})
}

func newPrivateKeyPEM() string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pemPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return string(pemPrivateKey)
}

func sampleSAKey() string {
	return fmt.Sprintf(sampleSAKeyFmt, wrappedKey())
}

func sampleSAKeyOneLine() string {
	return fmt.Sprintf(sampleSAKeyOneLineFmt, wrapKey(wrappedKey()))
}

func wrapKey(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

func wrappedKey() string {
	return wrapKey(newPrivateKeyPEM())
}
