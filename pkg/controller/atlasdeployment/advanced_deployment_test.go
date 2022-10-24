package atlasdeployment

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

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

func TestAdvancedDeployment_handleAutoscaling(t *testing.T) {
	testCases := []struct {
		desiredDeployment *v1.AdvancedDeploymentSpec
		currentDeployment *mongodbatlas.AdvancedCluster
		expected          *v1.AdvancedDeploymentSpec
		shouldFail        bool
		testName          string
		err               error
	}{
		{
			testName: "One region and autoscaling ENABLED for compute AND disk",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M30",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M30",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One region and autoscaling ENABLED for compute ONLY",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M40",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(false),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M40",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(false),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One region and autoscaling ENABLED for diskGB ONLY",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M40",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(false),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M40",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(false),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "Two regions and autoscaling ENABLED for compute AND disk in different regions",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M30",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
							{
								RegionName: "EASTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M30",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M30",
								},
							},
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M30",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
							{
								RegionName: "EASTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One region and autoscaling DISABLED for diskGB AND compute",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M20",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(false),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(false),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M20",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M20",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(false),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(false),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One regions and autoscaling ENABLED for compute and InstanceSize outside of min boundary",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M10",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M20",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M10",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M20",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M20",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One regions and autoscaling ENABLED for compute and InstanceSize outside of max boundary",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M50",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M50",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								RegionName: "WESTERN_EUROPE",
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M40",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: false,
		},
		{
			testName: "One region and autoscaling with wrong configuration",
			desiredDeployment: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: toptr.MakePtr(15),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M30",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "S10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			currentDeployment: &mongodbatlas.AdvancedCluster{
				ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
					{
						RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
							{
								ElectableSpecs: &mongodbatlas.Specs{
									InstanceSize: "M30",
								},
							},
						},
					},
				},
			},
			expected: &v1.AdvancedDeploymentSpec{
				DiskSizeGB: nil,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "us-east-1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									DiskIOPS:      nil,
									EbsVolumeType: "",
									InstanceSize:  "M30",
									NodeCount:     toptr.MakePtr(1),
								},
								AutoScaling: &v1.AdvancedAutoScalingSpec{
									DiskGB: &v1.DiskGB{
										Enabled: toptr.MakePtr(true),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "S10",
										MaxInstanceSize:  "M40",
									},
								},
							},
						},
					},
				},
			},
			shouldFail: true,
			err:        errors.New("instance size is invalid. instance family should be M or R"),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			ctx := workflow.NewContext(zap.S(), []status.Condition{})
			err := handleAutoscaling(ctx, tt.desiredDeployment, tt.currentDeployment)

			assert.Equal(t, tt.err, err)
			if !reflect.DeepEqual(tt.desiredDeployment, tt.expected) && !tt.shouldFail {
				expJSON, err := json.MarshalIndent(tt.expected, "", " ")
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				inpJSON, err := json.MarshalIndent(tt.desiredDeployment, "", " ")
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				t.Errorf("FAIL. Expected: %v, Got: %v", string(expJSON), string(inpJSON))
			}
		})
	}
}

func TestNormalizeInstanceSize(t *testing.T) {
	t.Run("InstanceSizeName should not change when inside of autoscaling configuration boundaries", func(t *testing.T) {
		ctx := workflow.NewContext(zap.S(), []status.Condition{})
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          toptr.MakePtr(true),
				ScaleDownEnabled: toptr.MakePtr(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M30",
			},
		}

		size, err := normalizeInstanceSize(ctx, "M10", autoscaling)

		assert.NoError(t, err)
		assert.Equal(t, "M10", size)
	})
	t.Run("InstanceSizeName should change to minimum size when outside of the bottom autoscaling configuration boundaries", func(t *testing.T) {
		ctx := workflow.NewContext(zap.S(), []status.Condition{})
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          toptr.MakePtr(true),
				ScaleDownEnabled: toptr.MakePtr(true),
				MinInstanceSize:  "M20",
				MaxInstanceSize:  "M30",
			},
		}

		size, err := normalizeInstanceSize(ctx, "M10", autoscaling)

		assert.NoError(t, err)
		assert.Equal(t, "M20", size)
	})
	t.Run("InstanceSizeName should change to maximum size when outside of the top autoscaling configuration boundaries", func(t *testing.T) {
		ctx := workflow.NewContext(zap.S(), []status.Condition{})
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          toptr.MakePtr(true),
				ScaleDownEnabled: toptr.MakePtr(true),
				MinInstanceSize:  "M20",
				MaxInstanceSize:  "M30",
			},
		}

		size, err := normalizeInstanceSize(ctx, "M40", autoscaling)

		assert.NoError(t, err)
		assert.Equal(t, "M30", size)
	})
}
