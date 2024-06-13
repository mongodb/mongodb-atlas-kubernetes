package atlasdeployment

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestSyncComputeConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name         string
		deployment   *akov2.AdvancedDeploymentSpec
		atlasCluster *mongodbatlas.AdvancedCluster
		expected     *akov2.AdvancedDeploymentSpec
	}{
		{
			name: "should not modify new region when there's no cluster in Atlas",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should not modify new region without autoscaling",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								RegionName: "EU_WEST2",
							},
						},
					},
				},
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should not modify when removing a region",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should not modify existing region without autoscaling",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								RegionName: "EU_WEST2",
							},
						},
					},
				},
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should set Atlas instance sizes for existing region with compute autoscaling enabled",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should unset compute autoscaling for existing region when it is disabled",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should unset disc size for existing region with disc autoscaling enabled and disk size has not be changed",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
		{
			name: "should keep disc size for existing region with disc autoscaling enabled but disk size has be changed",
			deployment: &akov2.AdvancedDeploymentSpec{
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
			},
			atlasCluster: &mongodbatlas.AdvancedCluster{
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
			},
			expected: &akov2.AdvancedDeploymentSpec{
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
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			syncRegionConfiguration(tc.deployment, tc.atlasCluster)
			assert.Equal(t, tc.expected, tc.deployment)
		})
	}
}

func TestNormalizeSpecs2(t *testing.T) {
	for _, tc := range []struct {
		name     string
		regions  []*akov2.AdvancedRegionConfig
		expected []*akov2.AdvancedRegionConfig
	}{
		{
			name:     "should do no action for a nil slice",
			regions:  nil,
			expected: nil,
		},
		{
			name:     "should do no action for a nil entry in the slice",
			regions:  []*akov2.AdvancedRegionConfig{nil},
			expected: []*akov2.AdvancedRegionConfig{nil},
		},
		{
			name: "should do no action when all specs are not nil",
			regions: []*akov2.AdvancedRegionConfig{
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
			},
			expected: []*akov2.AdvancedRegionConfig{
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
			},
		},
		{
			name: "should use electable spec as base when not nil",
			regions: []*akov2.AdvancedRegionConfig{
				{
					ElectableSpecs: &akov2.Specs{
						InstanceSize: "M10",
						NodeCount:    pointer.MakePtr(3),
					},
				},
			},
			expected: []*akov2.AdvancedRegionConfig{
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
			},
		},
		{
			name: "should use read only spec as base when not nil",
			regions: []*akov2.AdvancedRegionConfig{
				{
					ReadOnlySpecs: &akov2.Specs{
						InstanceSize: "M10",
						NodeCount:    pointer.MakePtr(3),
					},
				},
			},
			expected: []*akov2.AdvancedRegionConfig{
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
			},
		},
		{
			name: "should use analytics spec as base when not nil",
			regions: []*akov2.AdvancedRegionConfig{
				{
					AnalyticsSpecs: &akov2.Specs{
						InstanceSize: "M10",
						NodeCount:    pointer.MakePtr(3),
					},
				},
			},
			expected: []*akov2.AdvancedRegionConfig{
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
			},
		},
		{
			name: "should use read only spec as base when analytics is also not nil",
			regions: []*akov2.AdvancedRegionConfig{
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
			},
			expected: []*akov2.AdvancedRegionConfig{
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
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			normalizeSpecs(tc.regions)
			if !reflect.DeepEqual(tc.expected, tc.regions) {
				t.Errorf("expected: %v, got: %v", tc.expected, tc.regions)
			}
		})
	}
}
