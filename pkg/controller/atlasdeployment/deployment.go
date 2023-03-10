package atlasdeployment

import (
	"errors"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func ConvertLegacyDeployment(deploymentSpec *mdbv1.AtlasDeploymentSpec) error {
	legacy := deploymentSpec.DeploymentSpec

	if legacy == nil {
		return nil
	}

	replicationSpecs, err := convertLegacyReplicationSpecs(legacy)
	if err != nil {
		return err
	}

	deploymentSpec.AdvancedDeploymentSpec = &mdbv1.AdvancedDeploymentSpec{
		BackupEnabled:            legacy.ProviderBackupEnabled,
		BiConnector:              legacy.BIConnector,
		ClusterType:              getDefaultClusterType(legacy.ClusterType),
		DiskSizeGB:               legacy.DiskSizeGB,
		EncryptionAtRestProvider: legacy.EncryptionAtRestProvider,
		Labels:                   legacy.Labels,
		MongoDBMajorVersion:      legacy.MongoDBMajorVersion,
		Name:                     legacy.Name,
		Paused:                   legacy.Paused,
		PitEnabled:               legacy.PitEnabled,
		ReplicationSpecs:         replicationSpecs,
		CustomZoneMapping:        legacy.CustomZoneMapping,
		ManagedNamespaces:        legacy.ManagedNamespaces,
	}

	return nil
}

func convertLegacyReplicationSpecs(legacy *mdbv1.DeploymentSpec) ([]*mdbv1.AdvancedReplicationSpec, error) {
	result := []*mdbv1.AdvancedReplicationSpec{}

	if legacy == nil {
		return result, nil
	}

	if legacy.ProviderSettings == nil {
		return nil, errors.New("ProviderSettings should not be empty")
	}

	if len(legacy.ReplicationSpecs) == 0 {
		fillDefaultReplicationSpec(legacy)
	}

	for _, legacyResplicationSpec := range legacy.ReplicationSpecs {
		resplicationSpec := &mdbv1.AdvancedReplicationSpec{
			NumShards:     *convertLegacyInt64(legacyResplicationSpec.NumShards),
			ZoneName:      legacyResplicationSpec.ZoneName,
			RegionConfigs: []*mdbv1.AdvancedRegionConfig{},
		}

		for legacyRegion, legacyRegionConfig := range legacyResplicationSpec.RegionsConfig {
			regionConfig := mdbv1.AdvancedRegionConfig{
				AnalyticsSpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     convertLegacyInt64(legacyRegionConfig.AnalyticsNodes),
				},
				ElectableSpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     convertLegacyInt64(legacyRegionConfig.ElectableNodes),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     convertLegacyInt64(legacyRegionConfig.ReadOnlyNodes),
				},
				AutoScaling:         convertLegacyAutoScaling(legacy.AutoScaling, legacy.ProviderSettings.AutoScaling),
				BackingProviderName: legacy.ProviderSettings.BackingProviderName,
				Priority:            convertLegacyInt64(legacyRegionConfig.Priority),
				ProviderName:        string(legacy.ProviderSettings.ProviderName),
				RegionName:          legacyRegion,
			}

			resplicationSpec.RegionConfigs = append(resplicationSpec.RegionConfigs, &regionConfig)
		}

		result = append(result, resplicationSpec)
	}

	return result, nil
}

func convertLegacyAutoScaling(legacyRoot, legacyPS *mdbv1.AutoScalingSpec) *mdbv1.AdvancedAutoScalingSpec {
	if legacyRoot == nil || legacyPS == nil {
		return nil
	}

	autoScaling := &mdbv1.AdvancedAutoScalingSpec{
		DiskGB: &mdbv1.DiskGB{
			Enabled: legacyRoot.DiskGBEnabled,
		},
	}

	if legacyRoot.Compute != nil && legacyRoot.Compute.Enabled != nil {
		autoScaling.Compute = &mdbv1.ComputeSpec{
			Enabled:          legacyRoot.Compute.Enabled,
			ScaleDownEnabled: legacyRoot.Compute.ScaleDownEnabled,
			MinInstanceSize:  emptyIfDisabled(legacyPS.Compute.MinInstanceSize, legacyRoot.Compute.Enabled),
			MaxInstanceSize:  emptyIfDisabled(legacyPS.Compute.MaxInstanceSize, legacyRoot.Compute.Enabled),
		}
	}

	return autoScaling
}

func fillDefaultReplicationSpec(legacy *mdbv1.DeploymentSpec) {
	replicationSpec := mdbv1.ReplicationSpec{
		NumShards: toptr.MakePtr[int64](1),
		RegionsConfig: map[string]mdbv1.RegionsConfig{
			legacy.ProviderSettings.RegionName: {
				AnalyticsNodes: toptr.MakePtr(int64(0)),
				ElectableNodes: toptr.MakePtr(int64(3)),
				ReadOnlyNodes:  toptr.MakePtr(int64(0)),
				Priority:       toptr.MakePtr(int64(7)),
			},
		},
	}

	if legacy.ClusterType == mdbv1.TypeGeoSharded {
		replicationSpec.ZoneName = "Zone 1"
	}

	legacy.ReplicationSpecs = append(legacy.ReplicationSpecs, replicationSpec)
}

func getDefaultClusterType(legacyType mdbv1.DeploymentType) string {
	clusterType := mdbv1.TypeReplicaSet

	if legacyType != "" {
		clusterType = legacyType
	}

	return string(clusterType)
}

func convertLegacyInt64(input *int64) *int {
	if input == nil {
		return nil
	}

	return toptr.MakePtr(int(*input))
}

func emptyIfDisabled(value string, flag *bool) string {
	if flag == nil || !*flag {
		return ""
	}

	return value
}
