package atlasdeployment

import (
	"errors"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func ConvertLegacyDeployment(deployment *mdbv1.AtlasDeployment) error {
	legacy := deployment.Spec.DeploymentSpec

	replicationSpecs, err := convertLegacyReplicationSpecs(legacy)
	if err != nil {
		return err
	}

	deployment.Spec.AdvancedDeploymentSpec = &mdbv1.AdvancedDeploymentSpec{
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

	return nil
}

func convertLegacyReplicationSpecs(legacy *mdbv1.DeploymentSpec) ([]*mdbv1.AdvancedReplicationSpec, error) {
	result := []*mdbv1.AdvancedReplicationSpec{}

	if legacy.ProviderSettings == nil {
		return nil, errors.New("ProviderSettings should not be empty")
	}

	regionConfig := &mdbv1.AdvancedRegionConfig{
		AutoScaling: convertLegacyAutoScaling(legacy.AutoScaling),
		ElectableSpecs: &mdbv1.Specs{
			DiskIOPS:      legacy.ProviderSettings.DiskIOPS,
			EbsVolumeType: legacy.ProviderSettings.VolumeType,
			InstanceSize:  legacy.ProviderSettings.InstanceSizeName,
			NodeCount:     toptr.MakePtr(3),
		},
		BackingProviderName: legacy.ProviderSettings.BackingProviderName,
		Priority:            toptr.MakePtr(7),
		ProviderName:        string(legacy.ProviderSettings.ProviderName),
		RegionName:          legacy.ProviderSettings.RegionName,
	}

	resplicationSpec := &mdbv1.AdvancedReplicationSpec{
		RegionConfigs: []*mdbv1.AdvancedRegionConfig{regionConfig},
	}

	result = append(result, resplicationSpec)

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
