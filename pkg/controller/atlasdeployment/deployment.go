package atlasdeployment

import (
	"errors"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func ConvertLegacyDeployment(deploymentSpec *mdbv1.AtlasDeploymentSpec) error {
	legacy := deploymentSpec.DeploymentSpec

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
		Labels:                   []common.LabelSpec{},
		MongoDBMajorVersion:      legacy.MongoDBMajorVersion,
		Name:                     legacy.Name,
		Paused:                   legacy.Paused,
		PitEnabled:               legacy.PitEnabled,
		ReplicationSpecs:         replicationSpecs,
	}

	deploymentSpec.DeploymentSpec = nil

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

	legacyReplicatonSpecs := legacy.ReplicationSpecs
	if len(legacyReplicatonSpecs) == 0 {
		legacyReplicatonSpecs = append(legacyReplicatonSpecs, mdbv1.ReplicationSpec{
			NumShards: toptr.MakePtr[int64](1),
			ZoneName:  "Zone 1",
			RegionsConfig: map[string]mdbv1.RegionsConfig{
				legacy.ProviderSettings.RegionName: {
					AnalyticsNodes: toptr.MakePtr(int64(0)),
					ElectableNodes: toptr.MakePtr(int64(3)),
					ReadOnlyNodes:  toptr.MakePtr(int64(0)),
					Priority:       toptr.MakePtr(int64(7)),
				},
			},
		})
	}

	for _, legacyResplicationSpec := range legacyReplicatonSpecs {
		resplicationSpec := &mdbv1.AdvancedReplicationSpec{
			NumShards:     int(*legacyResplicationSpec.NumShards),
			ZoneName:      legacyResplicationSpec.ZoneName,
			RegionConfigs: []*mdbv1.AdvancedRegionConfig{},
		}

		for legacyRegion, legacyRegionConfig := range legacyResplicationSpec.RegionsConfig {
			regionConfig := mdbv1.AdvancedRegionConfig{
				AnalyticsSpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     toptr.MakePtr(int(*legacyRegionConfig.AnalyticsNodes)),
				},
				ElectableSpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     toptr.MakePtr(int(*legacyRegionConfig.ElectableNodes)),
				},
				ReadOnlySpecs: &mdbv1.Specs{
					DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
					EbsVolumeType: legacy.ProviderSettings.VolumeType,
					InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
					NodeCount:     toptr.MakePtr(int(*legacyRegionConfig.ReadOnlyNodes)),
				},
				AutoScaling:         convertLegacyAutoScaling(legacy.AutoScaling),
				BackingProviderName: legacy.ProviderSettings.BackingProviderName,
				Priority:            toptr.MakePtr(int(*legacyRegionConfig.Priority)),
				ProviderName:        string(legacy.ProviderSettings.ProviderName),
				RegionName:          legacyRegion,
			}

			resplicationSpec.RegionConfigs = append(resplicationSpec.RegionConfigs, &regionConfig)
		}

		result = append(result, resplicationSpec)
	}

	return result, nil
}

func convertLegacyAutoScaling(legacy *mdbv1.AutoScalingSpec) *mdbv1.AdvancedAutoScalingSpec {
	if legacy == nil {
		return nil
	}

	return &mdbv1.AdvancedAutoScalingSpec{
		DiskGB: &mdbv1.DiskGB{
			Enabled: legacy.DiskGBEnabled,
		},
		Compute: legacy.Compute,
	}
}

func getDefaultClusterType(legacyType mdbv1.DeploymentType) string {
	clusterType := mdbv1.TypeReplicaSet

	if legacyType != "" {
		clusterType = legacyType
	}

	return string(clusterType)
}
