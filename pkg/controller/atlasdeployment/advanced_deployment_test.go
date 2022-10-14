package atlasdeployment

import (
	"encoding/json"
	"reflect"
	"testing"

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
		input      *v1.AdvancedDeploymentSpec
		expected   *v1.AdvancedDeploymentSpec
		shouldFail bool
		testName   string
	}{
		{
			testName: "One region and autoscaling ENABLED for compute AND disk",
			input: &v1.AdvancedDeploymentSpec{
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
			input: &v1.AdvancedDeploymentSpec{
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
			input: &v1.AdvancedDeploymentSpec{
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
			input: &v1.AdvancedDeploymentSpec{
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
										Enabled:          toptr.MakePtr(false),
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
									InstanceSize:  "M20",
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
										MaxInstanceSize:  "M30",
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
										Enabled:          toptr.MakePtr(false),
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
										Enabled: toptr.MakePtr(false),
									},
									Compute: &v1.ComputeSpec{
										Enabled:          toptr.MakePtr(true),
										ScaleDownEnabled: nil,
										MinInstanceSize:  "M10",
										MaxInstanceSize:  "M30",
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
			input: &v1.AdvancedDeploymentSpec{
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
	}
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			handleAutoscaling(tt.input)
			if !reflect.DeepEqual(tt.input, tt.expected) && !tt.shouldFail {
				expJSON, err := json.MarshalIndent(tt.expected, "", " ")
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				inpJSON, err := json.MarshalIndent(tt.input, "", " ")
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
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          boolptr(true),
				ScaleDownEnabled: boolptr(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M30",
			},
		}

		assert.Equal(t, "M10", normalizeInstanceSize("M10", autoscaling))
	})
	t.Run("InstanceSizeName should change to minimum size when outside of the bottom autoscaling configuration boundaries", func(t *testing.T) {
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          boolptr(true),
				ScaleDownEnabled: boolptr(true),
				MinInstanceSize:  "M20",
				MaxInstanceSize:  "M30",
			},
		}

		assert.Equal(t, "M20", normalizeInstanceSize("M10", autoscaling))
	})
	t.Run("InstanceSizeName should change to maximum size when outside of the top autoscaling configuration boundaries", func(t *testing.T) {
		autoscaling := &v1.AdvancedAutoScalingSpec{
			Compute: &v1.ComputeSpec{
				Enabled:          boolptr(true),
				ScaleDownEnabled: boolptr(true),
				MinInstanceSize:  "M20",
				MaxInstanceSize:  "M30",
			},
		}

		assert.Equal(t, "M30", normalizeInstanceSize("M40", autoscaling))
	})
}

func TestExtractNumberFromInstanceTypeName(t *testing.T) {
	data := map[string]struct {
		Name   string
		Number int
	}{
		"Should extract number from M instance": {
			Name:   "M5",
			Number: 5,
		},
		"Should extract number from R instance": {
			Name:   "R50",
			Number: 50,
		},
		"Should extract number from NVME instance": {
			Name:   "R500_NVME",
			Number: 500,
		},
	}

	for name, test := range data {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.Number, extractNumberFromInstanceTypeName(test.Name))
		})
	}
}

func boolptr(b bool) *bool {
	return &b
}
