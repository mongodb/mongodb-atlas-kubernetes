// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validate

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
)

const (
	DeploymentSet = 1 << iota
	ServerlessSet
	FlexSet
)

func deploymentSpecMask(atlasDeployment *akov2.AtlasDeployment) int {
	mask := 0
	if atlasDeployment.Spec.DeploymentSpec != nil {
		mask = mask + DeploymentSet
	}
	if atlasDeployment.Spec.ServerlessSpec != nil {
		mask = mask + ServerlessSet
	}
	if atlasDeployment.Spec.FlexSpec != nil {
		mask = mask + FlexSet
	}
	return mask
}

func AtlasDeployment(atlasDeployment *akov2.AtlasDeployment) error {
	var err error
	var tagsSpec []*akov2.TagSpec
	switch deploymentSpecMask(atlasDeployment) {
	case 0:
		return errors.New("expected exactly one of spec.deploymentSpec or spec.serverlessSpec or spec.flexSpec to be present, but none were")
	case DeploymentSet:
		tagsSpec = atlasDeployment.Spec.DeploymentSpec.Tags
		err = regularDeployment(atlasDeployment.Spec.DeploymentSpec)
	case ServerlessSet:
		tagsSpec = atlasDeployment.Spec.ServerlessSpec.Tags
		err = serverlessDeployment(atlasDeployment.Spec.ServerlessSpec)
	case FlexSet:
		tagsSpec = atlasDeployment.Spec.FlexSpec.Tags
		err = flexDeployment(atlasDeployment.Spec.FlexSpec)
	default:
		return errors.New("expected exactly one of spec.deploymentSpec or spec.serverlessSpec or spec.flexSpec to be present, but multiple were")
	}

	if err != nil {
		return err
	}

	if err = Tags(tagsSpec); err != nil {
		return err
	}

	return nil
}

func regularDeployment(spec *akov2.AdvancedDeploymentSpec) error {
	var autoscaling akov2.AdvancedAutoScalingSpec
	for _, replicaSetSpec := range spec.ReplicationSpecs {
		for _, regionConfig := range replicaSetSpec.RegionConfigs {
			if err := providerConfig(regionConfig); err != nil {
				return err
			}

			if err := autoscalingForDeployment(regionConfig.AutoScaling, firstSetAutoscaling(&autoscaling, regionConfig)); err != nil {
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

func flexDeployment(spec *akov2.FlexSpec) error {
	supportedProviders := provider.SupportedProviders()
	switch {
	case spec.ProviderSettings == nil:
		return errors.New("provider settings cannot be empty")
	case !supportedProviders.IsSupported(provider.ProviderName(spec.ProviderSettings.BackingProviderName)):
		return errors.New("backing provider name is not supported")
	case spec.ProviderSettings.RegionName == "":
		return errors.New("regionName cannot be empty")
	}

	return nil
}
