package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

func TestAtlasDeployment(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment         *akov2.AtlasDeployment
		isGov                   bool
		regionUsageRestrictions string
		expectedError           string
	}{
		"Both specs present": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
					ServerlessSpec: &akov2.ServerlessSpec{},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "expected exactly one of spec.deploymentSpec or spec.serverlessSpec to be present, but none were",
		},
		"Neither spec present": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "expected exactly one of spec.deploymentSpec or spec.serverlessSpec to be present, but none were",
		},
		"Only DeploymentSpec present": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "us-east-1",
									},
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "",
		},
		"Only ServerlessSpec present": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AZURE",
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "",
		},
		"ServerlessSpec with provider config error": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderAWS,
							BackingProviderName: "AZURE",
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "provider name must be SERVERLESS",
		},
		"Regular deployment with config errors": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{RegionName: "invalid-region"},
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "provider name is not supported",
		},
		"Deployment with misconfigured tags": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "us-east-1",
									},
								},
							},
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "tag1",
								Value: "value1",
							},
							{
								Key:   "tag1",
								Value: "value2",
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "duplicate keys found in tags, this is forbidden",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := AtlasDeployment(tt.atlasDeployment, tt.isGov, tt.regionUsageRestrictions)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegularDeployment(t *testing.T) {
	tests := map[string]struct {
		spec                    *akov2.AdvancedDeploymentSpec
		isGov                   bool
		regionUsageRestrictions string
		expectedError           string
	}{
		"Valid regular deployment": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-east-1",
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "",
		},
		"Valid government deployment": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-gov-west-1",
							},
						},
					},
				},
			},
			isGov:                   true,
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "",
		},
		"Invalid government deployment": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{RegionName: "us-east-1"},
						},
					},
				},
			},
			isGov:                   true,
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "deployment in atlas for government support a restricted set of regions: us-east-1 is not part of AWS for government regions",
		},
		"Provider config error": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{ProviderName: "AliCloud"},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "provider name is not supported",
		},
		"AutoScaling is misconfigured across regions": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-east-1",
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									Compute: &akov2.ComputeSpec{
										Enabled:         pointer.MakePtr(true),
										MinInstanceSize: "M10",
										MaxInstanceSize: "M40",
									},
								},
							},
							{
								ProviderName: "AWS",
								RegionName:   "us-west-1",
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									Compute: &akov2.ComputeSpec{
										Enabled:         pointer.MakePtr(true),
										MinInstanceSize: "M20",
										MaxInstanceSize: "M40",
									},
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "autoscaling must be the same for all regions and across all replication specs for advanced deployment",
		},
		"AutoScaling is misconfigured across replications": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-east-1",
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									Compute: &akov2.ComputeSpec{
										Enabled:         pointer.MakePtr(true),
										MinInstanceSize: "M10",
										MaxInstanceSize: "M40",
									},
								},
							},
						},
					},
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-west-1",
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									Compute: &akov2.ComputeSpec{
										Enabled:         pointer.MakePtr(true),
										MinInstanceSize: "M20",
										MaxInstanceSize: "M40",
									},
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "autoscaling must be the same for all regions and across all replication specs for advanced deployment",
		},
		"Instance size is misconfigured": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-east-1",
								ElectableSpecs: &akov2.Specs{
									InstanceSize: "M10",
								},
								ReadOnlySpecs: &akov2.Specs{
									InstanceSize: "M20",
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment",
		},
		"Instance size is out of autoscaling range": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "us-east-1",
								ElectableSpecs: &akov2.Specs{
									InstanceSize: "M10",
								},
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									Compute: &akov2.ComputeSpec{
										Enabled:         pointer.MakePtr(true),
										MinInstanceSize: "M20",
										MaxInstanceSize: "M50",
									},
								},
							},
						},
					},
				},
			},
			isGov:                   false,
			regionUsageRestrictions: "",
			expectedError:           "the instance size is below the minimum autoscaling configuration",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := regularDeployment(tt.spec, tt.isGov, tt.regionUsageRestrictions)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProviderConfig(t *testing.T) {
	tests := map[string]struct {
		regionConfig  *akov2.AdvancedRegionConfig
		expectedError string
	}{
		"Serverless provider name": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: string(provider.ProviderServerless),
			},
			expectedError: "deployment cannot be configured as serverless. use dedicated configuration for serverless instance",
		},
		"Tenant with unsupported backing provider name": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ProviderName:        string(provider.ProviderTenant),
				BackingProviderName: "AliCloud",
			},
			expectedError: "backing provider name is not supported",
		},
		"Tenant with supported backing provider name": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ProviderName:        string(provider.ProviderTenant),
				BackingProviderName: string(provider.ProviderAWS),
			},
			expectedError: "",
		},
		"Unsupported provider name": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: "AliCloud",
			},
			expectedError: "provider name is not supported",
		},
		"Supported provider name": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ProviderName: string(provider.ProviderAWS),
			},
			expectedError: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := providerConfig(tt.regionConfig)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAutoscalingForDeployment(t *testing.T) {
	tests := map[string]struct {
		autoscaling         *akov2.AdvancedAutoScalingSpec
		previousAutoscaling *akov2.AdvancedAutoScalingSpec
		expectedError       string
	}{
		"Both AutoScaling are nil": {
			autoscaling:         nil,
			previousAutoscaling: nil,
			expectedError:       "",
		},
		"Both AutoScaling are the same": {
			autoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			previousAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			expectedError: "",
		},
		"AutoScaling.Compute.Enabled are different": {
			autoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			previousAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(false),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			expectedError: "autoscaling must be the same for all regions and across all replication specs for advanced deployment",
		},
		"AutoScaling.Compute.MinInstanceSize are different": {
			autoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M30",
					MaxInstanceSize: "M40",
				},
			},
			previousAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			expectedError: "autoscaling must be the same for all regions and across all replication specs for advanced deployment",
		},
		"AutoScaling.Compute.MaxInstanceSize are different": {
			autoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M50",
				},
			},
			previousAutoscaling: &akov2.AdvancedAutoScalingSpec{
				Compute: &akov2.ComputeSpec{
					Enabled:         pointer.MakePtr(true),
					MinInstanceSize: "M20",
					MaxInstanceSize: "M40",
				},
			},
			expectedError: "autoscaling must be the same for all regions and across all replication specs for advanced deployment",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := autoscalingForDeployment(tt.autoscaling, tt.previousAutoscaling)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInstanceSizeForDeployment(t *testing.T) {
	tests := map[string]struct {
		regionConfig  *akov2.AdvancedRegionConfig
		instanceSize  string
		expectedError string
	}{
		"All specs are nil": {
			regionConfig:  &akov2.AdvancedRegionConfig{},
			instanceSize:  "M3",
			expectedError: "",
		},
		"ElectableSpecs instance size mismatch": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ElectableSpecs: &akov2.Specs{InstanceSize: "M20"},
			},
			instanceSize:  "M30",
			expectedError: "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment",
		},
		"ReadOnlySpecs instance size mismatch": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ReadOnlySpecs: &akov2.Specs{InstanceSize: "M20"},
			},
			instanceSize:  "M30",
			expectedError: "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment",
		},
		"AnalyticsSpecs instance size mismatch": {
			regionConfig: &akov2.AdvancedRegionConfig{
				AnalyticsSpecs: &akov2.Specs{InstanceSize: "M20"},
			},
			instanceSize:  "M30",
			expectedError: "instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment",
		},
		"All specs match the instance size": {
			regionConfig: &akov2.AdvancedRegionConfig{
				ElectableSpecs: &akov2.Specs{InstanceSize: "M30"},
				ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M30"},
				AnalyticsSpecs: &akov2.Specs{InstanceSize: "M30"},
			},
			instanceSize:  "M30",
			expectedError: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := instanceSizeForDeployment(tt.regionConfig, tt.instanceSize)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFirstNonEmptyInstanceSize(t *testing.T) {
	tests := map[string]struct {
		currentInstanceSize string
		regionConfig        *akov2.AdvancedRegionConfig
		expectedResult      string
	}{
		"Non-empty current instance size": {
			currentInstanceSize: "M30",
			regionConfig:        &akov2.AdvancedRegionConfig{},
			expectedResult:      "M30",
		},
		"ElectableSpecs instance size is used": {
			currentInstanceSize: "",
			regionConfig: &akov2.AdvancedRegionConfig{
				ElectableSpecs: &akov2.Specs{InstanceSize: "M20"},
			},
			expectedResult: "M20",
		},
		"ReadOnlySpecs instance size is used": {
			currentInstanceSize: "",
			regionConfig: &akov2.AdvancedRegionConfig{
				ReadOnlySpecs: &akov2.Specs{InstanceSize: "M30"},
			},
			expectedResult: "M30",
		},
		"AnalyticsSpecs instance size is used": {
			currentInstanceSize: "",
			regionConfig: &akov2.AdvancedRegionConfig{
				AnalyticsSpecs: &akov2.Specs{InstanceSize: "M40"},
			},
			expectedResult: "M40",
		},
		"All instance sizes are empty": {
			currentInstanceSize: "",
			regionConfig: &akov2.AdvancedRegionConfig{
				ElectableSpecs: &akov2.Specs{InstanceSize: ""},
				ReadOnlySpecs:  &akov2.Specs{InstanceSize: ""},
				AnalyticsSpecs: &akov2.Specs{InstanceSize: ""},
			},
			expectedResult: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			currentInstanceSize := tt.currentInstanceSize
			result := firstNonEmptyInstanceSize(&currentInstanceSize, tt.regionConfig)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestInstanceSizeRangeForAdvancedDeployment(t *testing.T) {
	tests := map[string]struct {
		regionConfig  *akov2.AdvancedRegionConfig
		expectedError string
	}{
		"AutoScaling is nil": {
			regionConfig:  &akov2.AdvancedRegionConfig{},
			expectedError: "",
		},
		"AutoScaling.Compute is nil": {
			regionConfig:  &akov2.AdvancedRegionConfig{AutoScaling: &akov2.AdvancedAutoScalingSpec{}},
			expectedError: "",
		},
		"AutoScaling.Compute.Enabled is nil": {
			regionConfig:  &akov2.AdvancedRegionConfig{AutoScaling: &akov2.AdvancedAutoScalingSpec{Compute: &akov2.ComputeSpec{}}},
			expectedError: "",
		},
		"AutoScaling.Compute.Enabled is false": {
			regionConfig:  &akov2.AdvancedRegionConfig{AutoScaling: &akov2.AdvancedAutoScalingSpec{Compute: &akov2.ComputeSpec{Enabled: pointer.MakePtr(false)}}},
			expectedError: "",
		},
		"ElectableSpecs instance size below minimum": {
			regionConfig: &akov2.AdvancedRegionConfig{
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					Compute: &akov2.ComputeSpec{
						Enabled:         pointer.MakePtr(true),
						MinInstanceSize: "M20",
						MaxInstanceSize: "M40",
					},
				},
				ElectableSpecs: &akov2.Specs{InstanceSize: "M10"},
			},
			expectedError: "the instance size is below the minimum autoscaling configuration",
		},
		"ReadOnlySpecs instance size above maximum": {
			regionConfig: &akov2.AdvancedRegionConfig{
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					Compute: &akov2.ComputeSpec{
						Enabled:         pointer.MakePtr(true),
						MinInstanceSize: "M20",
						MaxInstanceSize: "M40",
					},
				},
				ReadOnlySpecs: &akov2.Specs{InstanceSize: "M50"},
			},
			expectedError: "the instance size is above the maximum autoscaling configuration",
		},
		"AnalyticsSpecs instance size above maximum": {
			regionConfig: &akov2.AdvancedRegionConfig{
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					Compute: &akov2.ComputeSpec{
						Enabled:         pointer.MakePtr(true),
						MinInstanceSize: "M20",
						MaxInstanceSize: "M40",
					},
				},
				AnalyticsSpecs: &akov2.Specs{InstanceSize: "M50"},
			},
			expectedError: "the instance size is above the maximum autoscaling configuration",
		},
		"All specs within range": {
			regionConfig: &akov2.AdvancedRegionConfig{
				AutoScaling: &akov2.AdvancedAutoScalingSpec{
					Compute: &akov2.ComputeSpec{
						Enabled:         pointer.MakePtr(true),
						MinInstanceSize: "M20",
						MaxInstanceSize: "M40",
					},
				},
				ElectableSpecs: &akov2.Specs{InstanceSize: "M30"},
				ReadOnlySpecs:  &akov2.Specs{InstanceSize: "M30"},
				AnalyticsSpecs: &akov2.Specs{InstanceSize: "M30"},
			},
			expectedError: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := instanceSizeRangeForAdvancedDeployment(tt.regionConfig)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdvancedInstanceSizeInRange(t *testing.T) {
	tests := map[string]struct {
		currentInstanceSize string
		minInstanceSize     string
		maxInstanceSize     string
		expectedError       string
	}{
		"Error on currentInstanceSize": {
			currentInstanceSize: "invalid",
			minInstanceSize:     "M10",
			maxInstanceSize:     "M30",
			expectedError:       "instance size is invalid. instance family should be M or R",
		},
		"Error on minInstanceSize": {
			currentInstanceSize: "M20",
			minInstanceSize:     "invalid",
			maxInstanceSize:     "M30",
			expectedError:       "instance size is invalid. instance family should be M or R",
		},
		"Error on maxInstanceSize": {
			currentInstanceSize: "M20",
			minInstanceSize:     "M10",
			maxInstanceSize:     "invalid",
			expectedError:       "instance size is invalid. instance family should be M or R",
		},
		"Instance size below minimum": {
			currentInstanceSize: "M10",
			minInstanceSize:     "M20",
			maxInstanceSize:     "M30",
			expectedError:       "the instance size is below the minimum autoscaling configuration",
		},
		"Instance size above maximum": {
			currentInstanceSize: "M40",
			minInstanceSize:     "M20",
			maxInstanceSize:     "M30",
			expectedError:       "the instance size is above the maximum autoscaling configuration",
		},
		"Instance size within range": {
			currentInstanceSize: "M20",
			minInstanceSize:     "M10",
			maxInstanceSize:     "M30",
			expectedError:       "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := advancedInstanceSizeInRange(tt.currentInstanceSize, tt.minInstanceSize, tt.maxInstanceSize)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeploymentForGov(t *testing.T) {
	tests := map[string]struct {
		spec                    *akov2.AdvancedDeploymentSpec
		regionUsageRestrictions string
		expectedError           string
	}{
		"All regions supported": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{RegionName: "us-gov-west-1"},
						},
					},
				},
			},
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "",
		},
		"One region not supported": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{RegionName: "us-gov-west-1"},
							{RegionName: "us-east-1"},
						},
					},
				},
			},
			regionUsageRestrictions: "COMMERCIAL_FEDRAMP_REGIONS_ONLY",
			expectedError:           "deployment in atlas for government support a restricted set of regions: us-gov-west-1 is not part of AWS FedRAMP regions",
		},
		"All regions not supported": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{RegionName: "us-east-1"},
							{RegionName: "eu-west-1"},
						},
					},
				},
			},
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "deployment in atlas for government support a restricted set of regions: us-east-1 is not part of AWS for government regions",
		},
		"No regions specified": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						RegionConfigs: []*akov2.AdvancedRegionConfig{},
					},
				},
			},
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "",
		},
		"No replication specs specified": {
			spec: &akov2.AdvancedDeploymentSpec{
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{},
			},
			regionUsageRestrictions: "GOV_REGIONS_ONLY",
			expectedError:           "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := deploymentForGov(tt.spec, tt.regionUsageRestrictions)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServerlessDeployment(t *testing.T) {
	tests := map[string]struct {
		spec          *akov2.ServerlessSpec
		expectedError string
	}{
		"Provider settings nil": {
			spec: &akov2.ServerlessSpec{
				ProviderSettings: nil,
			},
			expectedError: "provider settings cannot be empty",
		},
		"Provider name not SERVERLESS": {
			spec: &akov2.ServerlessSpec{
				ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
					ProviderName:        provider.ProviderAWS,
					BackingProviderName: "",
				},
			},
			expectedError: "provider name must be SERVERLESS",
		},
		"Backing provider name not supported": {
			spec: &akov2.ServerlessSpec{
				ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
					ProviderName:        provider.ProviderServerless,
					BackingProviderName: "AliCloud",
				},
			},
			expectedError: "backing provider name is not supported",
		},
		"Serverless private endpoint are wrong": {
			spec: &akov2.ServerlessSpec{
				ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
					ProviderName:        provider.ProviderServerless,
					BackingProviderName: "AWS",
				},
				PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
					{Name: "sp1"},
					{Name: "sp1"},
				},
			},
			expectedError: "serverless private endpoint should have a unique name: sp1 is duplicated",
		},
		"Valid serverless spec": {
			spec: &akov2.ServerlessSpec{
				ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
					ProviderName:        provider.ProviderServerless,
					BackingProviderName: "AWS",
				},
			},
			expectedError: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := serverlessDeployment(tt.spec)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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

		assert.ErrorContains(t, err, "serverless private endpoint should have a unique name: spe-1 is duplicated")
	})
}
