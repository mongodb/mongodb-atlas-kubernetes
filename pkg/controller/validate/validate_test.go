package validate

import (
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
