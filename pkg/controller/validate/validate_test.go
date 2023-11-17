package validate

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, ServerlessSpec: &mdbv1.ServerlessSpec{}}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{DeploymentSpec: nil}
			assert.Error(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("different instance sizes for advanced deployment", func(t *testing.T) {
			t.Run("different instance size in the same region", func(t *testing.T) {
				spec := mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
			spec := mdbv1.AtlasDeploymentSpec{DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, ServerlessSpec: nil}
			assert.NoError(t, DeploymentSpec(&spec, false, "NONE"))
			assert.Nil(t, DeploymentSpec(&spec, false, "NONE"))
		})
		t.Run("Advanced cluster with replication config", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{
				DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
			DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
				ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
					{
						ZoneName:  "Zone EU",
						NumShards: 1,
						RegionConfigs: []*mdbv1.AdvancedRegionConfig{
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
		deploy := mdbv1.AtlasDeploymentSpec{
			DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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

	t.Run("copy setting is set but replica-set id is not available", func(t *testing.T) {
		bSchedule := &mdbv1.AtlasBackupSchedule{
			Spec: mdbv1.AtlasBackupScheduleSpec{
				CopySettings: []mdbv1.CopySetting{
					{
						RegionName:       toptr.MakePtr("US_WEST_1"),
						CloudProvider:    toptr.MakePtr("AWS"),
						ShouldCopyOplogs: toptr.MakePtr(true),
						Frequencies:      []string{"WEEKLY"},
					},
				},
			},
		}
		deployment := &mdbv1.AtlasDeployment{
			Spec: mdbv1.AtlasDeploymentSpec{
				DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
					PitEnabled: toptr.MakePtr(true),
				},
			},
		}
		assert.Error(t, BackupSchedule(bSchedule, deployment))
	})

	t.Run("copy settings on advanced deployment", func(t *testing.T) {
		t.Run("copy settings is valid", func(t *testing.T) {
			bSchedule := &mdbv1.AtlasBackupSchedule{
				Spec: mdbv1.AtlasBackupScheduleSpec{
					CopySettings: []mdbv1.CopySetting{
						{
							RegionName:       toptr.MakePtr("US_WEST_1"),
							CloudProvider:    toptr.MakePtr("AWS"),
							ShouldCopyOplogs: toptr.MakePtr(true),
							Frequencies:      []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
							RegionName: toptr.MakePtr("US_WEST_1"),
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{},
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
							RegionName:       toptr.MakePtr("US_WEST_1"),
							CloudProvider:    toptr.MakePtr("AWS"),
							ShouldCopyOplogs: toptr.MakePtr(true),
							Frequencies:      []string{"WEEKLY"},
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
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
							RegionName: toptr.MakePtr("US_WEST_1"),
						},
					},
				},
			}
			deployment := &mdbv1.AtlasDeployment{
				Spec: mdbv1.AtlasDeploymentSpec{
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{},
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
		prj := mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
					sampleAlertConfig("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
				},
				AlertConfigurationSyncEnabled: false,
			},
		}
		assert.NoError(t, Project(&prj, false /*isGov*/))
	})

	t.Run("should fail on duplications when alert config is enabled", func(t *testing.T) {
		prj := mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
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
		prj := mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
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
		prj := mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				AlertConfigurations: []mdbv1.AlertConfiguration{
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

func sampleAlertConfig(typeName string) mdbv1.AlertConfiguration {
	return mdbv1.AlertConfiguration{
		EventTypeName: typeName,
		Enabled:       true,
		Threshold: &mdbv1.Threshold{
			Operator:  "LESS_THAN",
			Threshold: "1",
			Units:     "HOURS",
		},
		Notifications: []mdbv1.Notification{
			{
				IntervalMin:  5,
				DelayMin:     toptr.MakePtr(5),
				EmailEnabled: toptr.MakePtr(true),
				SMSEnabled:   toptr.MakePtr(false),
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
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(0),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
					},
				},
			},
		}

		assert.NoError(t, instanceSizeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should fail when instance size are different between node types", func(t *testing.T) {
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(0),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M20",
							NodeCount:    toptr.MakePtr(1),
						},
					},
				},
			},
		}

		assert.EqualError(t, instanceSizeForAdvancedDeployment(replicationSpecs), "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment")
	})

	t.Run("should fail when instance size are different across regions", func(t *testing.T) {
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(0),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
					},
					{
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M20",
							NodeCount:    toptr.MakePtr(0),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M20",
							NodeCount:    toptr.MakePtr(1),
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
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
					},
				},
			},
		}

		assert.NoError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should succeed when instance size is with autoscaling range", func(t *testing.T) {
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
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
		}

		assert.NoError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs))
	})

	t.Run("should fail when instance size is below autoscaling range", func(t *testing.T) {
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M10",
							NodeCount:    toptr.MakePtr(1),
						},
						AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
							Compute: &mdbv1.ComputeSpec{
								Enabled:          toptr.MakePtr(true),
								ScaleDownEnabled: toptr.MakePtr(true),
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
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						ElectableSpecs: &mdbv1.Specs{
							InstanceSize: "M40",
							NodeCount:    toptr.MakePtr(3),
						},
						ReadOnlySpecs: &mdbv1.Specs{
							InstanceSize: "M40",
							NodeCount:    toptr.MakePtr(1),
						},
						AnalyticsSpecs: &mdbv1.Specs{
							InstanceSize: "M40",
							NodeCount:    toptr.MakePtr(1),
						},
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
		}

		assert.EqualError(t, instanceSizeRangeForAdvancedDeployment(replicationSpecs), "the instance size is above the maximum autoscaling configuration")
	})
}

func TestAutoscalingForAdvancedDeployment(t *testing.T) {
	t.Run("should fail when different compute autoscaling config are set", func(t *testing.T) {
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
							DiskGB: &mdbv1.DiskGB{
								Enabled: toptr.MakePtr(true),
							},
							Compute: &mdbv1.ComputeSpec{
								Enabled:          toptr.MakePtr(true),
								ScaleDownEnabled: toptr.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M40",
							},
						},
					},
					{
						AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
							DiskGB: &mdbv1.DiskGB{
								Enabled: toptr.MakePtr(true),
							},
							Compute: &mdbv1.ComputeSpec{
								Enabled:          toptr.MakePtr(true),
								ScaleDownEnabled: toptr.MakePtr(false),
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
		replicationSpecs := []*mdbv1.AdvancedReplicationSpec{
			{
				RegionConfigs: []*mdbv1.AdvancedRegionConfig{
					{
						AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
							DiskGB: &mdbv1.DiskGB{
								Enabled: toptr.MakePtr(false),
							},
							Compute: &mdbv1.ComputeSpec{
								Enabled:          toptr.MakePtr(true),
								ScaleDownEnabled: toptr.MakePtr(true),
								MinInstanceSize:  "M10",
								MaxInstanceSize:  "M40",
							},
						},
					},
					{
						AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
							DiskGB: &mdbv1.DiskGB{
								Enabled: toptr.MakePtr(true),
							},
							Compute: &mdbv1.ComputeSpec{
								Enabled:          toptr.MakePtr(true),
								ScaleDownEnabled: toptr.MakePtr(true),
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
