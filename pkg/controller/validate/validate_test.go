package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{}, ServerlessSpec: &akov2.ServerlessSpec{}}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := akov2.AtlasDeploymentSpec{DeploymentSpec: nil}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("different instance sizes for advanced deployment", func(t *testing.T) {
			t.Run("different instance size in the same region", func(t *testing.T) {
				spec := akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M20"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
			})
			t.Run("different instance size in different regions", func(t *testing.T) {
				spec := akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									},
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M20"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
			})
			t.Run("different instance size in different replications", func(t *testing.T) {
				spec := akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									},
								},
							},
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M20"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
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
				spec := akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M30",
											},
										},
									},
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									},
								},
							},
						},
					},
				}
				assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
			})
			t.Run("different autoscaling in different replications", func(t *testing.T) {
				spec := akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M10",
												MaxInstanceSize:  "M30",
											},
										},
									},
								},
							},
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{InstanceSize: "M20"},
										ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
										AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
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
			spec := akov2.AtlasDeploymentSpec{DeploymentSpec: &akov2.AdvancedDeploymentSpec{}, ServerlessSpec: nil}
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("Advanced cluster with replication config", func(t *testing.T) {
			spec := akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
									ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
									AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
											MinInstanceSize:  "M10",
											MaxInstanceSize:  "M30",
										},
									},
								},
							},
						},
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
									ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M10"},
									AnalyticsSpecs: &akov2.Specs{InstanceSize: "M10"},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled:          pointer.MakePtr(true),
											ScaleDownEnabled: pointer.MakePtr(true),
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
		deploy := akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						ZoneName:  "Zone EU",
						NumShards: 1,
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								RegionName: "EU_EAST_1",
							},
						},
					},
				},
			},
		}

		assert.ErrorContains(t, deploymentForGov(&deploy, "GOV_REGIONS_ONLY"), "deployment in atlas for government support a restricted set of regions: EU_EAST_1 is not part of AWS for government regions")
	})

	t.Run("should fail when advanced deployment is configured to non-gov region", func(t *testing.T) {
		deploy := akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
			},
		}

		assert.ErrorContains(t, Project(akoProject, false), "regionUsageRestriction can be used only with Atlas for government")
	})

	t.Run("should fail when commercial Atlas sets region restriction field to COMMERCIAL_FEDRAMP_REGIONS_ONLY", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "COMMERCIAL_FEDRAMP_REGIONS_ONLY",
			},
		}

		assert.ErrorContains(t, Project(akoProject, false), "regionUsageRestriction can be used only with Atlas for government")
	})

	t.Run("should not fail if commercial Atlas sets region restriction field to empty", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{},
		}

		assert.NoError(t, Project(akoProject, false))
	})

	t.Run("should not fail if commercial Atlas sets region restriction field to NONE", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "NONE",
			},
		}

		assert.NoError(t, Project(akoProject, false))
	})

	t.Run("custom roles spec", func(t *testing.T) {
		t.Run("empty custom roles spec", func(t *testing.T) {
			spec := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{},
			}
			assert.NoError(t, Project(spec, false))
		})
		t.Run("valid custom roles spec", func(t *testing.T) {
			spec := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					CustomRoles: []akov2.CustomRole{
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
			spec := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					CustomRoles: []akov2.CustomRole{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []akov2.NetworkPeer{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []akov2.NetworkPeer{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &akov2.EncryptionAtRest{
					GoogleCloudKms: akov2.GoogleCloudKms{
						Enabled: pointer.MakePtr(true),
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. disable encryption at rest for Google Cloud")
	})

	t.Run("should fail if there's a Azure encryption at rest config", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &akov2.EncryptionAtRest{
					AzureKeyVault: akov2.AzureKeyVault{
						Enabled: pointer.MakePtr(true),
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "atlas for government only supports AWS provider. disable encryption at rest for Azure")
	})

	t.Run("should fail if there's a AWS encryption at rest config with wrong region", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				EncryptionAtRest: &akov2.EncryptionAtRest{
					AwsKms: akov2.AwsKms{
						Enabled: pointer.MakePtr(true),
						Region:  "us-east-1",
					},
				},
			},
		}

		assert.ErrorContains(t, Project(akoProject, true), "encryption at rest in atlas for government support a restricted set of regions: us-east-1 is not part of AWS for government regions")
	})

	t.Run("should fail if there's non AWS private endpoint config", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				PrivateEndpoints: []akov2.PrivateEndpoint{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "COMMERCIAL_FEDRAMP_REGIONS_ONLY",
				PrivateEndpoints: []akov2.PrivateEndpoint{
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				NetworkPeers: []akov2.NetworkPeer{
					{
						ProviderName:        "AWS",
						AccepterRegionName:  "us-gov-east-1",
						ContainerRegion:     "us-gov-east-1",
						RouteTableCIDRBlock: "192.168.0.0/16",
						AtlasCIDRBlock:      "10.8.0.0/22",
					},
				},
				EncryptionAtRest: &akov2.EncryptionAtRest{
					AwsKms: akov2.AwsKms{
						Enabled: pointer.MakePtr(true),
						Region:  "us-gov-east-1",
					},
				},
				PrivateEndpoints: []akov2.PrivateEndpoint{
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
		bSchedule := &akov2.AtlasBackupSchedule{
			Spec: akov2.AtlasBackupScheduleSpec{
				AutoExportEnabled: true,
			},
		}
		deployment := &akov2.AtlasDeployment{
			Status: status.AtlasDeploymentStatus{},
		}
		assert.Error(t, BackupSchedule(bSchedule, deployment))
	})

	t.Run("copy setting is set but replica-set id is not available", func(t *testing.T) {
		bSchedule := &akov2.AtlasBackupSchedule{
			Spec: akov2.AtlasBackupScheduleSpec{
				CopySettings: []akov2.CopySetting{
					{
						RegionName:       pointer.MakePtr("US_WEST_1"),
						CloudProvider:    pointer.MakePtr("AWS"),
						ShouldCopyOplogs: pointer.MakePtr(true),
						Frequencies:      []string{"WEEKLY"},
					},
				},
			},
		}
		deployment := &akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					PitEnabled: pointer.MakePtr(true),
				},
			},
		}
		assert.Error(t, BackupSchedule(bSchedule, deployment))
	})

	t.Run("copy settings on advanced deployment", func(t *testing.T) {
		t.Run("copy settings is valid", func(t *testing.T) {
			bSchedule := &akov2.AtlasBackupSchedule{
				Spec: akov2.AtlasBackupScheduleSpec{
					CopySettings: []akov2.CopySetting{
						{
							RegionName:       pointer.MakePtr("US_WEST_1"),
							CloudProvider:    pointer.MakePtr("AWS"),
							ShouldCopyOplogs: pointer.MakePtr(true),
							Frequencies:      []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						PitEnabled: pointer.MakePtr(true),
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
			bSchedule := &akov2.AtlasBackupSchedule{
				Spec: akov2.AtlasBackupScheduleSpec{
					CopySettings: []akov2.CopySetting{
						{
							ShouldCopyOplogs: pointer.MakePtr(true),
						},
						{
							RegionName: pointer.MakePtr("US_WEST_1"),
						},
					},
				},
			}
			deployment := &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
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
			bSchedule := &akov2.AtlasBackupSchedule{
				Spec: akov2.AtlasBackupScheduleSpec{
					CopySettings: []akov2.CopySetting{
						{
							RegionName:       pointer.MakePtr("US_WEST_1"),
							CloudProvider:    pointer.MakePtr("AWS"),
							ShouldCopyOplogs: pointer.MakePtr(true),
							Frequencies:      []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						PitEnabled: pointer.MakePtr(true),
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
			bSchedule := &akov2.AtlasBackupSchedule{
				Spec: akov2.AtlasBackupScheduleSpec{
					CopySettings: []akov2.CopySetting{
						{
							ShouldCopyOplogs: pointer.MakePtr(true),
						},
						{
							RegionName: pointer.MakePtr("US_WEST_1"),
						},
					},
				},
			}
			deployment := &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
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

func TestProjectAlertConfigs(t *testing.T) {
	t.Run("should not fail on duplications when alert config is disabled", func(t *testing.T) {
		prj := akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				AlertConfigurations: []akov2.AlertConfiguration{
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
				},
				AlertConfigurationSyncEnabled: false,
			},
		}
		assert.NoError(t, Project(&prj, false /*isGov*/))
	})

	t.Run("should fail on duplications when alert config is enabled", func(t *testing.T) {
		prj := akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				AlertConfigurations: []akov2.AlertConfiguration{
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
				},
				AlertConfigurationSyncEnabled: true,
			},
		}
		assert.ErrorContains(t, Project(&prj, false /*isGov*/),
			"alert config at position 1 is a duplicate of alert config at position 0")
	})

	t.Run("should fail on first duplication in when alert config is enabled", func(t *testing.T) {
		prj := akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				AlertConfigurations: []akov2.AlertConfiguration{
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("JOINED_GROUP"),
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("JOINED_GROUP"),
				},
				AlertConfigurationSyncEnabled: true,
			},
		}
		assert.ErrorContains(t, Project(&prj, false /*isGov*/),
			"alert config at position 2 is a duplicate of alert config at position 0")
	})

	t.Run("should succeed on absence of duplications in when alert config is enabled", func(t *testing.T) {
		prj := akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				AlertConfigurations: []akov2.AlertConfiguration{
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("JOINED_GROUP"),
					sampleAlertConfig("invented_event_3"),
					sampleAlertConfig("invented_event_4"),
				},
				AlertConfigurationSyncEnabled: true,
			},
		}
		assert.NoError(t, Project(&prj, false /*isGov*/))
	})
}

func sampleAlertConfig(typeName string) akov2.AlertConfiguration {
	return akov2.AlertConfiguration{
		EventTypeName: typeName,
		Enabled:       true,
		Threshold: &akov2.Threshold{
			Operator:  "LESS_THAN",
			Threshold: "1",
			Units:     "HOURS",
		},
		Notifications: []akov2.Notification{
			{
				IntervalMin:  5,
				DelayMin:     pointer.MakePtr(5),
				EmailEnabled: pointer.MakePtr(true),
				SMSEnabled:   pointer.MakePtr(false),
				Roles: []string{
					"GROUP_OWNER",
				},
				TypeName: "GROUP",
			},
		},
	}
}

func TestInstanceSizeForAdvancedDeployment(t *testing.T) {
	t.Run("should succeed when instance size are the same for all node types", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(0),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
					},
				},
			},
		}

		assert.NoError(t, instanceSizeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should fail when instance size are different between node types", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(0),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M20",
							NodeCount:    pointer.MakePtr(1),
						},
					},
				},
			},
		}

		assert.EqualError(t, instanceSizeForAdvancedDeployment(replicationSpecs), "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment")
	})

	t.Run("should fail when instance size are different across regions", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(0),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
					},
					{
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M20",
							NodeCount:    pointer.MakePtr(0),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M20",
							NodeCount:    pointer.MakePtr(1),
						},
					},
				},
			},
		}

		assert.EqualError(t, instanceSizeForAdvancedDeployment(replicationSpecs), "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment")
	})
}

func TestInstanceSizeRangeForAdvancedDeployment(t *testing.T) {
	t.Run("should succeed when region has no autoscaling config", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
					},
				},
			},
		}

		assert.NoError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should succeed when instance size is with autoscaling range", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M30",
							},
						},
					},
				},
			},
		}

		assert.NoError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should fail when instance size is below autoscaling range", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(1),
						},
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M20",
								MaxInstanceSize:  "M40",
							},
						},
					},
				},
			},
		}

		assert.EqualError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs), "the instance size is below the minimum autoscaling configuration")
	})

	t.Run("should fail when instance size is above autoscaling range", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						ElectableSpecs: &akov2.Specs{
							InstanceSize: "M40",
							NodeCount:    pointer.MakePtr(3),
						},
						ReadOnlySpecs: &akov2.Specs{
							InstanceSize: "M40",
							NodeCount:    pointer.MakePtr(1),
						},
						AnalyticsSpecs: &akov2.Specs{
							InstanceSize: "M40",
							NodeCount:    pointer.MakePtr(1),
						},
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M30",
							},
						},
					},
				},
			},
		}

		assert.EqualError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs), "the instance size is above the maximum autoscaling configuration")
	})
}

func TestAutoscalingForAdvancedDeployment(t *testing.T) {
	t.Run("should fail when different compute autoscaling config are set", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
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
					{
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							DiskGB: &akov2.DiskGB{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(false),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M40",
							},
						},
					},
				},
			},
		}

		assert.EqualError(t, autoscalingForAdvancedDeployment(replicationSpecs), "autoscaling must be the same for all regions and across all replication specs for advanced deployment")
	})

	t.Run("should fail when different disc autoscaling config are set", func(t *testing.T) {
		replicationSpecs := []*akov2.AdvancedReplicationSpec{
			{
				RegionConfigs: []*akov2.AdvancedRegionConfig{
					{
						AutoScaling: &akov2.AdvancedAutoScalingSpec{
							DiskGB: &akov2.DiskGB{
								Enabled: pointer.MakePtr(false),
							},
							Compute: &akov2.ComputeSpec{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M40",
							},
						},
					},
					{
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
		}

		assert.EqualError(t, autoscalingForAdvancedDeployment(replicationSpecs), "autoscaling must be the same for all regions and across all replication specs for advanced deployment")
	})
}

func TestServerlessPrivateEndpoints(t *testing.T) {
	t.Run("should pass when there are no private endpoints with the same name", func(t *testing.T) {
		privateEndpoints := []akov2.ServerlessPrivateEndpoint{
			{
				Name: "spe-1",
			},
			{
				Name: "spe-2",
			},
			{
				Name: "spe-3",
			},
		}

		err := serverlessPrivateEndpoints(privateEndpoints)

		assert.NoError(t, err)
	})

	t.Run("should fail when there are private endpoints with duplicated name", func(t *testing.T) {
		privateEndpoints := []akov2.ServerlessPrivateEndpoint{
			{
				Name: "spe-1",
			},
			{
				Name: "spe-2",
			},
			{
				Name: "spe-1",
			},
			{
				Name: "spe-3",
			},
			{
				Name: "spe-2",
			},
		}

		err := serverlessPrivateEndpoints(privateEndpoints)

		assert.ErrorContains(t, err, "serverless private endpoint should have a unique name: spe-1 is duplicated\nserverless private endpoint should have a unique name: spe-2 is duplicated")
	})
}
