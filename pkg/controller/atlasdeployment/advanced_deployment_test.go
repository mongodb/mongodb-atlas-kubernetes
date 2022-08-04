package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestMergedAdvancedDeployment(t *testing.T) {
	defaultAtlas := v1.DefaultAwsAdvancedDeployment("default", "my-project")
	defaultAtlas.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ProviderName = "TENANT"
	defaultAtlas.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName = "AWS"
	defaultAtlas.Spec.AdvancedDeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs = &v1.Specs{
		InstanceSize: "M5",
	}

	t.Run("Test merging clusters removes backing provider name if empty", func(t *testing.T) {
		advancedCluster := mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ID:        "123",
					ZoneName:  "Zone1",
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName:          "US_EAST_1",
							BackingProviderName: "",
						},
					},
				},
			},
		}

		merged, _, err := MergedAdvancedDeployment(advancedCluster, *defaultAtlas.Spec.AdvancedDeploymentSpec)
		assert.NoError(t, err)
		assert.Empty(t, merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})

	t.Run("Test merging clusters does not remove backing provider name if it is present in the atlas type", func(t *testing.T) {
		advancedCluster := mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ID:        "123",
					ZoneName:  "Zone1",
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName:          "US_EAST_1",
							BackingProviderName: "AWS",
						},
					},
				},
			},
		}

		merged, _, err := MergedAdvancedDeployment(advancedCluster, *defaultAtlas.Spec.AdvancedDeploymentSpec)
		assert.NoError(t, err)
		assert.Equal(t, "AWS", merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})
}

func TestAdvancedDeploymentOutdatedFields(t *testing.T) {
	autoScalingEnabled := true
	autoScalingDisabled := !autoScalingEnabled

	t.Run("Operator unset instanceSize of autoscaled deployment", func(t *testing.T) {
		autoScalingDisabledInstanceSize := "M30"

		advancedCluster := v1.AdvancedDeploymentSpec{
			ReplicationSpecs: []*v1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*v1.AdvancedRegionConfig{
						{
							RegionName: "US_EAST_1",
							ElectableSpecs: &v1.Specs{
								InstanceSize: "M10",
							},
							AnalyticsSpecs: &v1.Specs{
								InstanceSize: "M10",
							},
							AutoScaling: &v1.AdvancedAutoScalingSpec{
								Compute: &v1.ComputeSpec{
									Enabled: &autoScalingEnabled,
								},
							},
						},
						{
							RegionName: "US_EAST_2",
							ElectableSpecs: &v1.Specs{
								InstanceSize: "M20",
							},
							AutoScaling: &v1.AdvancedAutoScalingSpec{
								Compute: &v1.ComputeSpec{
									Enabled: &autoScalingEnabled,
								},
							},
						},
						{
							RegionName: "US_EAST_1",
							ElectableSpecs: &v1.Specs{
								InstanceSize: autoScalingDisabledInstanceSize,
							},
							AutoScaling: &v1.AdvancedAutoScalingSpec{
								Compute: &v1.ComputeSpec{
									Enabled: &autoScalingDisabled,
								},
							},
						},
					},
				},
			},
		}

		c := cleanupTheSpec(advancedCluster)
		// first regionConfig
		assert.Equal(t, c.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize, "")
		assert.Equal(t, c.ReplicationSpecs[0].RegionConfigs[0].AnalyticsSpecs.InstanceSize, "")
		// second regionConfig
		assert.Equal(t, c.ReplicationSpecs[0].RegionConfigs[1].ElectableSpecs.InstanceSize, "")
		// third regionConfig with disabled autoscaling
		assert.Equal(t, c.ReplicationSpecs[0].RegionConfigs[2].ElectableSpecs.InstanceSize, autoScalingDisabledInstanceSize)
	})

	t.Run("Operator unset diskSizeGB of autoscaled deployment", func(t *testing.T) {
		diskSize := 40
		advancedCluster := v1.AdvancedDeploymentSpec{
			DiskSizeGB: &diskSize,
			ReplicationSpecs: []*v1.AdvancedReplicationSpec{
				{
					RegionConfigs: []*v1.AdvancedRegionConfig{
						{
							RegionName: "US_EAST_1",
							ElectableSpecs: &v1.Specs{
								InstanceSize: "M10",
							},
							AnalyticsSpecs: &v1.Specs{
								InstanceSize: "M10",
							},
							AutoScaling: &v1.AdvancedAutoScalingSpec{
								DiskGB: &v1.DiskGB{
									Enabled: &autoScalingEnabled,
								},
							},
						},
					},
				},
			},
		}
		c := cleanupTheSpec(advancedCluster)
		assert.Nil(t, c.DiskSizeGB)
	})
}
