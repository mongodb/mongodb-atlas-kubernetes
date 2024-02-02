package atlasdeployment

import (
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func syncRegionConfiguration(deploymentSpec *mdbv1.AdvancedDeploymentSpec, atlasCluster *mongodbatlas.AdvancedCluster) {
	// When there's no config to handle, do nothing
	if deploymentSpec == nil || len(deploymentSpec.ReplicationSpecs) == 0 {
		return
	}

	// When there's no cluster in Atlas, we need to keep configuration
	if atlasCluster == nil || len(atlasCluster.ReplicationSpecs) == 0 {
		return
	}

	for _, regionSpec := range deploymentSpec.ReplicationSpecs[0].RegionConfigs {
		// When disc auto-scaling is enabled and there's no updated on disk size, unset disc size letting auto-scaling config control it
		if isDiskAutoScalingEnabled(regionSpec.AutoScaling) && !hasDiskSizeChanged(deploymentSpec.DiskSizeGB, atlasCluster.DiskSizeGB) {
			deploymentSpec.DiskSizeGB = nil
		}
	}

	// when editing a region, normalize change compute configuration
	regionsHasChanged := false
	if regionsConfigHasChanged(deploymentSpec.ReplicationSpecs[0].RegionConfigs, atlasCluster.ReplicationSpecs[0].RegionConfigs) {
		regionsHasChanged = true
		normalizeSpecs(deploymentSpec.ReplicationSpecs[0].RegionConfigs)
	}

	for _, regionSpec := range deploymentSpec.ReplicationSpecs[0].RegionConfigs {
		// When compute auto-scaling is enabled, unset instance size to avoid override production workload
		if isComputeAutoScalingEnabled(regionSpec.AutoScaling) {
			if !regionsHasChanged {
				regionSpec.ElectableSpecs.InstanceSize = ""
				regionSpec.ReadOnlySpecs.InstanceSize = ""
				regionSpec.AnalyticsSpecs.InstanceSize = ""
			}
		} else {
			if regionSpec.AutoScaling != nil {
				regionSpec.AutoScaling.Compute = nil
			}
		}
	}
}

func isComputeAutoScalingEnabled(autoScalingSpec *mdbv1.AdvancedAutoScalingSpec) bool {
	return autoScalingSpec != nil && autoScalingSpec.Compute != nil && autoScalingSpec.Compute.Enabled != nil && *autoScalingSpec.Compute.Enabled
}

func isDiskAutoScalingEnabled(autoScalingSpec *mdbv1.AdvancedAutoScalingSpec) bool {
	return autoScalingSpec != nil && autoScalingSpec.DiskGB != nil && autoScalingSpec.DiskGB.Enabled != nil && *autoScalingSpec.DiskGB.Enabled
}

func hasDiskSizeChanged(deploymentDiskSize *int, clusterDiskSize *float64) bool {
	if deploymentDiskSize == nil && clusterDiskSize == nil {
		return false
	}

	if deploymentDiskSize == nil && clusterDiskSize != nil {
		return true
	}

	if deploymentDiskSize != nil && clusterDiskSize == nil {
		return true
	}

	if *clusterDiskSize < 0 {
		return true
	}

	return *deploymentDiskSize != int(*clusterDiskSize)
}

func regionsConfigHasChanged(deploymentRegions []*mdbv1.AdvancedRegionConfig, atlasRegions []*mongodbatlas.AdvancedRegionConfig) bool {
	if len(deploymentRegions) != len(atlasRegions) {
		return true
	}

	mapDeploymentRegions := map[string]*mdbv1.AdvancedRegionConfig{}
	for _, region := range deploymentRegions {
		mapDeploymentRegions[region.RegionName] = region
	}

	for _, region := range atlasRegions {
		if _, ok := mapDeploymentRegions[region.RegionName]; !ok {
			return true
		}

		var atlasAsOperatorRegion mdbv1.AdvancedRegionConfig
		err := compat.JSONCopy(&atlasAsOperatorRegion, region)
		if err != nil {
			return true
		}

		if cmp.Diff(mapDeploymentRegions[region.RegionName], &atlasAsOperatorRegion) != "" {
			return true
		}
	}

	return false
}

func normalizeSpecs(regions []*mdbv1.AdvancedRegionConfig) {
	for _, region := range regions {
		if region == nil {
			return
		}

		var notNilSpecs mdbv1.Specs
		if region.ElectableSpecs != nil {
			notNilSpecs = *region.ElectableSpecs
		} else if region.ReadOnlySpecs != nil {
			notNilSpecs = *region.ReadOnlySpecs
		} else if region.AnalyticsSpecs != nil {
			notNilSpecs = *region.AnalyticsSpecs
		}

		if region.ElectableSpecs == nil {
			region.ElectableSpecs = &notNilSpecs
			region.ElectableSpecs.NodeCount = pointer.MakePtr(0)
		}

		if region.ReadOnlySpecs == nil {
			region.ReadOnlySpecs = &notNilSpecs
			region.ReadOnlySpecs.NodeCount = pointer.MakePtr(0)
		}

		if region.AnalyticsSpecs == nil {
			region.AnalyticsSpecs = &notNilSpecs
			region.AnalyticsSpecs.NodeCount = pointer.MakePtr(0)
		}
	}
}
