package validate

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

func AtlasDeployment(atlasDeployment *akov2.AtlasDeployment, isGov bool, regionUsageRestrictions string) error {
	isRegularDeployment := atlasDeployment.Spec.DeploymentSpec != nil
	isServerlessDeployment := atlasDeployment.Spec.ServerlessSpec != nil
	var err error
	var tagsSpec []*akov2.TagSpec

	switch {
	case !isRegularDeployment && !isServerlessDeployment:
		return errors.New("expected exactly one of spec.deploymentSpec or spec.serverlessSpec to be present, but none were")
	case isRegularDeployment && isServerlessDeployment:
		return errors.New("expected exactly one of spec.deploymentSpec or spec.serverlessSpec to be present, but none were")
	case !isRegularDeployment && isServerlessDeployment:
		tagsSpec = atlasDeployment.Spec.ServerlessSpec.Tags
		err = serverlessDeployment(atlasDeployment.Spec.ServerlessSpec)
	default:
		tagsSpec = atlasDeployment.Spec.DeploymentSpec.Tags
		err = regularDeployment(atlasDeployment.Spec.DeploymentSpec, isGov, regionUsageRestrictions)
	}

	if err != nil {
		return err
	}

	if err = Tags(tagsSpec); err != nil {
		return err
	}

	return nil
}

func regularDeployment(spec *akov2.AdvancedDeploymentSpec, isGov bool, regionUsageRestrictions string) error {
	if isGov {
		if err := deploymentForGov(spec, regionUsageRestrictions); err != nil {
			return err
		}
	}

	var autoscaling akov2.AdvancedAutoScalingSpec
	var instanceSize string
	for _, replicaSetSpec := range spec.ReplicationSpecs {
		for _, regionConfig := range replicaSetSpec.RegionConfigs {
			if err := providerConfig(regionConfig); err != nil {
				return err
			}

			if err := autoscalingForDeployment(regionConfig.AutoScaling, firstSetAutoscaling(&autoscaling, regionConfig)); err != nil {
				return err
			}

			if err := instanceSizeForDeployment(regionConfig, firstNonEmptyInstanceSize(&instanceSize, regionConfig)); err != nil {
				return err
			}

			if err := instanceSizeRangeForAdvancedDeployment(regionConfig); err != nil {
				return err
			}
		}
	}

	return nil
}

func providerConfig(regionConfig *akov2.AdvancedRegionConfig) error {
	supportedProviders := provider.SupportedProviders()

	switch {
	case regionConfig.ProviderName == string(provider.ProviderServerless):
		return errors.New("deployment cannot be configured as serverless. use dedicated configuration for serverless instance")
	case regionConfig.ProviderName == string(provider.ProviderTenant):
		if !supportedProviders.IsSupported(provider.ProviderName(regionConfig.BackingProviderName)) {
			return errors.New("backing provider name is not supported")
		}
	default:
		if !supportedProviders.IsSupported(provider.ProviderName(regionConfig.ProviderName)) {
			return errors.New("provider name is not supported")
		}
	}

	return nil
}

func autoscalingForDeployment(autoscaling, previousAutoscaling *akov2.AdvancedAutoScalingSpec) error {
	if autoscaling == nil && previousAutoscaling == nil {
		return nil
	}

	if cmp.Diff(autoscaling, previousAutoscaling, cmpopts.EquateEmpty()) != "" {
		return errors.New("autoscaling must be the same for all regions and across all replication specs for advanced deployment")
	}

	return nil
}

func instanceSizeForDeployment(regionConfig *akov2.AdvancedRegionConfig, instanceSize string) error {
	err := errors.New("instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment")

	if regionConfig.ElectableSpecs != nil && regionConfig.ElectableSpecs.InstanceSize != instanceSize {
		return err
	}

	if regionConfig.ReadOnlySpecs != nil && regionConfig.ReadOnlySpecs.InstanceSize != instanceSize {
		return err
	}

	if regionConfig.AnalyticsSpecs != nil && regionConfig.AnalyticsSpecs.InstanceSize != instanceSize {
		return err
	}

	return nil
}

func firstNonEmptyInstanceSize(currentInstanceSize *string, regionConfig *akov2.AdvancedRegionConfig) string {
	if currentInstanceSize != nil && *currentInstanceSize != "" {
		return *currentInstanceSize
	}

	if regionConfig.ElectableSpecs != nil && regionConfig.ElectableSpecs.InstanceSize != "" {
		*currentInstanceSize = regionConfig.ElectableSpecs.InstanceSize
		return *currentInstanceSize
	}

	if regionConfig.ReadOnlySpecs != nil && regionConfig.ReadOnlySpecs.InstanceSize != "" {
		*currentInstanceSize = regionConfig.ReadOnlySpecs.InstanceSize
		return *currentInstanceSize
	}

	if regionConfig.AnalyticsSpecs != nil && regionConfig.AnalyticsSpecs.InstanceSize != "" {
		*currentInstanceSize = regionConfig.AnalyticsSpecs.InstanceSize
		return *currentInstanceSize
	}

	return ""
}

func firstSetAutoscaling(autoscaling *akov2.AdvancedAutoScalingSpec, regionConfig *akov2.AdvancedRegionConfig) *akov2.AdvancedAutoScalingSpec {
	if autoscaling.DiskGB != nil || autoscaling.Compute != nil {
		return autoscaling
	}

	if regionConfig.AutoScaling != nil {
		*autoscaling = *regionConfig.AutoScaling
		return autoscaling
	}

	return nil
}

func instanceSizeRangeForAdvancedDeployment(regionConfig *akov2.AdvancedRegionConfig) error {
	if regionConfig.AutoScaling == nil || regionConfig.AutoScaling.Compute == nil || regionConfig.AutoScaling.Compute.Enabled == nil || !*regionConfig.AutoScaling.Compute.Enabled {
		return nil
	}

	if regionConfig.ElectableSpecs != nil {
		if err := advancedInstanceSizeInRange(
			regionConfig.ElectableSpecs.InstanceSize,
			regionConfig.AutoScaling.Compute.MinInstanceSize,
			regionConfig.AutoScaling.Compute.MaxInstanceSize); err != nil {
			return err
		}
	}

	if regionConfig.ReadOnlySpecs != nil {
		if err := advancedInstanceSizeInRange(
			regionConfig.ReadOnlySpecs.InstanceSize,
			regionConfig.AutoScaling.Compute.MinInstanceSize,
			regionConfig.AutoScaling.Compute.MaxInstanceSize); err != nil {
			return err
		}
	}

	if regionConfig.AnalyticsSpecs != nil {
		if err := advancedInstanceSizeInRange(
			regionConfig.AnalyticsSpecs.InstanceSize,
			regionConfig.AutoScaling.Compute.MinInstanceSize,
			regionConfig.AutoScaling.Compute.MaxInstanceSize); err != nil {
			return err
		}
	}

	return nil
}

func advancedInstanceSizeInRange(currentInstanceSize, minInstanceSize, maxInstanceSize string) error {
	minSize, err := NewFromInstanceSizeName(minInstanceSize)
	if err != nil {
		return err
	}

	maxSize, err := NewFromInstanceSizeName(maxInstanceSize)
	if err != nil {
		return err
	}

	currentSize, err := NewFromInstanceSizeName(currentInstanceSize)
	if err != nil {
		return err
	}

	if CompareInstanceSizes(currentSize, minSize) == -1 {
		return errors.New("the instance size is below the minimum autoscaling configuration")
	}

	if CompareInstanceSizes(currentSize, maxSize) == 1 {
		return errors.New("the instance size is above the maximum autoscaling configuration")
	}

	return nil
}

func deploymentForGov(spec *akov2.AdvancedDeploymentSpec, regionUsageRestrictions string) error {
	for _, replication := range spec.ReplicationSpecs {
		for _, region := range replication.RegionConfigs {
			regionErr := validCloudGovRegion(regionUsageRestrictions, region.RegionName)
			if regionErr != nil {
				return fmt.Errorf("deployment in atlas for government support a restricted set of regions: %w", regionErr)
			}
		}
	}

	return nil
}

func serverlessDeployment(spec *akov2.ServerlessSpec) error {
	supportedProviders := provider.SupportedProviders()
	switch {
	case spec.ProviderSettings == nil:
		return errors.New("provider settings cannot be empty")
	case spec.ProviderSettings.ProviderName != provider.ProviderServerless:
		return errors.New("provider name must be SERVERLESS")
	case !supportedProviders.IsSupported(provider.ProviderName(spec.ProviderSettings.BackingProviderName)):
		return errors.New("backing provider name is not supported")
	}

	err := serverlessPrivateEndpoints(spec.PrivateEndpoints)
	if err != nil {
		return err
	}

	return nil
}

func serverlessPrivateEndpoints(privateEndpoints []akov2.ServerlessPrivateEndpoint) error {
	namesMap := map[string]struct{}{}

	for _, privateEndpoint := range privateEndpoints {
		if _, ok := namesMap[privateEndpoint.Name]; ok {
			return fmt.Errorf("serverless private endpoint should have a unique name: %s is duplicated", privateEndpoint.Name)
		}

		namesMap[privateEndpoint.Name] = struct{}{}
	}

	return nil
}
