package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestSyncComputeConfiguration(t *testing.T) {
	t.Run("should not modify new region when there's no cluster in Atlas", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, nil)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should not modify new region without autoscaling", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
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
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
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
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
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
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
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

	t.Run("should unset instance size for existing region with compute autoscaling enabled", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mdbv1.Specs{
								NodeCount: toptr.MakePtr(3),
							},
							ReadOnlySpecs: &mdbv1.Specs{
								NodeCount: toptr.MakePtr(1),
							},
							AnalyticsSpecs: &mdbv1.Specs{
								NodeCount: toptr.MakePtr(1),
							},
							AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
								Compute: &mdbv1.ComputeSpec{
									Enabled:          toptr.MakePtr(true),
									ScaleDownEnabled: toptr.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
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
								NodeCount:    toptr.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								Compute: &mongodbatlas.Compute{
									Enabled:          toptr.MakePtr(true),
									ScaleDownEnabled: toptr.MakePtr(true),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should unset compute autoscaling for existing region when it is disabled", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
							AutoScaling: &mdbv1.AdvancedAutoScalingSpec{},
							Priority:    toptr.MakePtr(7),
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
								NodeCount:    toptr.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								Compute: &mongodbatlas.Compute{
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should unset disc size for existing region with disc autoscaling enabled and disk size has not be changed", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(20),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &mdbv1.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
								Compute: &mdbv1.ComputeSpec{
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &mdbv1.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			DiskSizeGB: toptr.MakePtr(20.0),
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								DiskGB: &mongodbatlas.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
								Compute: &mongodbatlas.Compute{
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}

		syncRegionConfiguration(advancedDeployment, atlasCluster)
		assert.Equal(t, expected, advancedDeployment)
	})

	t.Run("should keep disc size for existing region with disc autoscaling enabled but disk size has be changed", func(t *testing.T) {
		advancedDeployment := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(30),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &mdbv1.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
								Compute: &mdbv1.ComputeSpec{
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		expected := &mdbv1.AdvancedDeploymentSpec{
			DiskSizeGB: toptr.MakePtr(30),
			ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mdbv1.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
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
								DiskGB: &mdbv1.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
							},
							Priority: toptr.MakePtr(7),
						},
					},
				},
			},
		}
		atlasCluster := &mongodbatlas.AdvancedCluster{
			DiskSizeGB: toptr.MakePtr(20.0),
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST2",
							ElectableSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(3),
							},
							ReadOnlySpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AnalyticsSpecs: &mongodbatlas.Specs{
								InstanceSize: "M10",
								NodeCount:    toptr.MakePtr(1),
							},
							AutoScaling: &mongodbatlas.AdvancedAutoScaling{
								DiskGB: &mongodbatlas.DiskGB{
									Enabled: toptr.MakePtr(true),
								},
								Compute: &mongodbatlas.Compute{
									Enabled:          toptr.MakePtr(false),
									ScaleDownEnabled: toptr.MakePtr(false),
									MinInstanceSize:  "M10",
									MaxInstanceSize:  "M30",
								},
							},
							Priority: toptr.MakePtr(7),
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
		var regions []*mdbv1.AdvancedRegionConfig
		normalizeSpecs(regions)

		assert.Nil(t, regions)
	})

	t.Run("should do no action for a nil entry in the slice", func(t *testing.T) {
		regions := []*mdbv1.AdvancedRegionConfig{
			nil,
		}
		normalizeSpecs(regions)

		assert.Equal(
			t,
			[]*mdbv1.AdvancedRegionConfig{
				nil,
			},
			regions,
		)
	})

	t.Run("should do no action when all specs are not nil", func(t *testing.T) {
		regions := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(2),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(1),
				},
			},
		}
		expected := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(2),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(1),
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
		regions := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
			},
		}
		expected := []*mdbv1.AdvancedRegionConfig{
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
					NodeCount:    toptr.MakePtr(0),
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
		regions := []*mdbv1.AdvancedRegionConfig{
			{
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
			},
		}
		expected := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(0),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(0),
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
		regions := []*mdbv1.AdvancedRegionConfig{
			{
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
			},
		}
		expected := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(0),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(0),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
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
		regions := []*mdbv1.AdvancedRegionConfig{
			{
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(2),
				},
			},
		}
		expected := []*mdbv1.AdvancedRegionConfig{
			{
				ElectableSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(0),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(3),
				},
				AnalyticsSpecs: &mdbv1.Specs{
					InstanceSize: "M10",
					NodeCount:    toptr.MakePtr(2),
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
