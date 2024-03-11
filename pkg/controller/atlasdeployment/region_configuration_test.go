package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestSyncComputeConfiguration(t *testing.T) {
	t.Run("should not modify new region when there's no cluster in Atlas", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, nil)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should not modify new region without autoscaling", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName: "EU_WEST2",
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should not modify new region with autoscaling", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName: "EU_WEST2",
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should not modify when removing a region", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName: "EU_WEST2",
						},
						{
							RegionName: "EU_WEST1",
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should not modify existing region without autoscaling", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName: "EU_WEST2",
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should set Atlas instance sizes for existing region with compute autoscaling enabled", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: pointer.MakePtr(7),
						},
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST1",
							ElectableSpecs: &akov2.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(3),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(6),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &akov2.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(3),
							},
							ReadOnlySpecs: &akov2.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(1),
							},
							AnalyticsSpecs: &akov2.Specs{
								InstanceSize: "M30",
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
							Priority: pointer.MakePtr(7),
						},
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST1",
							ElectableSpecs: &akov2.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(3),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(6),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								Compute: &mongodbatlas.Compute{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST1",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M30",
								NodeCount:    pointer.MakePtr(3),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								Compute: &mongodbatlas.Compute{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(6),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should unset compute autoscaling for existing region when it is disabled", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
									Enabled:          pointer.MakePtr(false),
									ScaleDownEnabled: pointer.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							AutoScaling: &akov2.AdvancedAutoScalingSpec{},
							Priority:    pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								Compute: &mongodbatlas.Compute{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should unset disc size for existing region with disc autoscaling enabled and disk size has not be changed", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(20),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(false),
									ScaleDownEnabled: pointer.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			DiskSizeGB: pointer.MakePtr(20.0),
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								DiskGB: &mongodbatlas.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &mongodbatlas.Compute{
									Enabled:          pointer.MakePtr(false),
									ScaleDownEnabled: pointer.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should keep disc size for existing region with disc autoscaling enabled but disk size has be changed", func(t *testing.T) {
		advancedDeployment := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(30),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(false),
									ScaleDownEnabled: pointer.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &akov2.AdvancedDeploymentSpec{
			DiskSizeGB: pointer.MakePtr(30),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			DiskSizeGB: pointer.MakePtr(20.0),
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    pointer.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								DiskGB: &mongodbatlas.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &mongodbatlas.Compute{
									Enabled:          pointer.MakePtr(false),
									ScaleDownEnabled: pointer.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: pointer.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})
}

func TestNormalizeSpecs(t *testing.T) {
	t.Run("should do no action for a nil slice", func(t *testing.T) {
		var regions []*akov2.AdvancedRegionConfig
		normalizeSpecs(regions)

		assert.Nil(t, regions)
	})

	t.Run("should do no action for a nil entry in the slice", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			nil,
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			[]*akov2.AdvancedRegionConfig{
				nil,
			},
			regions,
		)
	})

	t.Run("should do no action when all specs are not nil", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(2),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(1),
				},
			},
		}
		expected := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(2),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(1),
				},
			},
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			expected,
			regions,
		)
	})

	t.Run("should use electable spec as base when not nil", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
		}
		expected := []*akov2.AdvancedRegionConfig{
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
					NodeCount:    pointer.MakePtr(0),
				},
			},
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			expected,
			regions,
		)
	})

	t.Run("should use read only spec as base when not nil", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			{
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
		}
		expected := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(0),
				},
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(0),
				},
			},
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			expected,
			regions,
		)
	})

	t.Run("should use analytics spec as base when not nil", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			{
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
		}
		expected := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(0),
				},
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(0),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
			},
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			expected,
			regions,
		)
	})

	t.Run("should use read only spec as base when analytics is also not nil", func(t *testing.T) {
		regions := []*akov2.AdvancedRegionConfig{
			{
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(2),
				},
			},
		}
		expected := []*akov2.AdvancedRegionConfig{
			{
				ElectableSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(0),
				},
				ReadOnlySpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(3),
				},
				AnalyticsSpecs: &akov2.Specs{
					InstanceSize: "M10",
					NodeCount:    pointer.MakePtr(2),
				},
			},
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			expected,
			regions,
		)
	})
}
