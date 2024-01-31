package atlasdeployment

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
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
				for _, arc := range atlasCluster.ReplicationSpecs[0].RegionConfigs {
					if arc.RegionName == regionSpec.RegionName {
						if regionSpec.ElectableSpecs != nil && arc.ElectableSpecs != nil {
							regionSpec.ElectableSpecs.InstanceSize = arc.ElectableSpecs.InstanceSize
						}
						if regionSpec.ReadOnlySpecs != nil && arc.ReadOnlySpecs != nil {
							regionSpec.ReadOnlySpecs.InstanceSize = arc.ReadOnlySpecs.InstanceSize
						}
						if regionSpec.AnalyticsSpecs != nil && arc.AnalyticsSpecs != nil {
							regionSpec.AnalyticsSpecs.InstanceSize = arc.AnalyticsSpecs.InstanceSize
						}
						break
					}
				}
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
		// do not compare instance sizes when autoscaling is ON
		if isComputeAutoScalingEnabled(region.AutoScaling) {
			region = ignoreInstanceSize(region.DeepCopy())
		}
		mapDeploymentRegions[region.RegionName] = region
	}

	for _, region := range atlasRegions {
		k8sRegion, ok := mapDeploymentRegions[region.RegionName]
		if !ok {
			return true
		}

		var atlasAsOperatorRegion mdbv1.AdvancedRegionConfig
		var atlasRegion = &atlasAsOperatorRegion
		err := compat.JSONCopy(atlasRegion, region)
		if err != nil {
			return true
		}

		// do not compare instance sizes when autoscaling is ON
		if isComputeAutoScalingEnabled(k8sRegion.AutoScaling) {
			atlasRegion = ignoreInstanceSize(atlasRegion.DeepCopy())
		}

		if diff := cmp.Diff(k8sRegion, atlasRegion); diff != "" {
			return true
		}
	}

	return false
}

func ignoreInstanceSize(rc *mdbv1.AdvancedRegionConfig) *mdbv1.AdvancedRegionConfig {
	if rc == nil {
		return rc
	}
	if rc.ElectableSpecs != nil {
		rc.ElectableSpecs.InstanceSize = ""
	}
	if rc.ReadOnlySpecs != nil {
		rc.ReadOnlySpecs.InstanceSize = ""
	}
	if rc.AnalyticsSpecs != nil {
		rc.AnalyticsSpecs.InstanceSize = ""
	}
	return rc
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
			region.ElectableSpecs.NodeCount = toptr.MakePtr(0)
		}

		if region.ReadOnlySpecs == nil {
			region.ReadOnlySpecs = &notNilSpecs
			region.ReadOnlySpecs.NodeCount = toptr.MakePtr(0)
		}

		if region.AnalyticsSpecs == nil {
			region.AnalyticsSpecs = &notNilSpecs
			region.AnalyticsSpecs.NodeCount = toptr.MakePtr(0)
		}
	}
}
