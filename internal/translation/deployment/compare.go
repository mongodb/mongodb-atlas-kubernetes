package deployment

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func ComputeChanges(desired, current *Cluster) (*Cluster, bool) {

	// Paused is special case that must be handled individually from other changes
	if !areEqual(desired.Paused, current.Paused) {
		return &Cluster{
			ProjectID: desired.ProjectID,
			AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:   desired.Name,
				Paused: pointer.MakePtr(pointer.GetOrDefault(desired.Paused, false)),
			},
		}, true
	}

	if specAreEqual(desired, current) {
		return nil, false
	}

	changes := &Cluster{
		ProjectID:                 desired.GetProjectID(),
		computeAutoscalingEnabled: desired.computeAutoscalingEnabled,
		AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:                         desired.Name,
			ClusterType:                  desired.ClusterType,
			MongoDBMajorVersion:          desired.MongoDBMajorVersion,
			VersionReleaseSystem:         desired.VersionReleaseSystem,
			BackupEnabled:                desired.BackupEnabled,
			EncryptionAtRestProvider:     desired.EncryptionAtRestProvider,
			BiConnector:                  desired.BiConnector,
			PitEnabled:                   desired.PitEnabled,
			RootCertType:                 desired.RootCertType,
			TerminationProtectionEnabled: desired.TerminationProtectionEnabled,
			Labels:                       desired.Labels,
			Tags:                         desired.Tags,
		},
	}

	if desired.DiskSizeGB != nil && !areEqual(desired.DiskSizeGB, current.DiskSizeGB) {
		changes.DiskSizeGB = desired.DiskSizeGB
	}

	changesReplicationSpecs := make([]*akov2.AdvancedReplicationSpec, 0, len(desired.ReplicationSpecs))
	for _, desiredReplicationSpec := range desired.ReplicationSpecs {
		changesRegionConfig := make([]*akov2.AdvancedRegionConfig, 0, len(desiredReplicationSpec.RegionConfigs))
		for _, desiredRegionConfig := range desiredReplicationSpec.RegionConfigs {
			changesRegionConfig = append(
				changesRegionConfig,
				&akov2.AdvancedRegionConfig{
					ProviderName:        desiredRegionConfig.ProviderName,
					BackingProviderName: desiredRegionConfig.BackingProviderName,
					RegionName:          desiredRegionConfig.RegionName,
					Priority:            desiredRegionConfig.Priority,
					ElectableSpecs:      getSpecsChanges(desiredRegionConfig.ElectableSpecs),
					ReadOnlySpecs:       getSpecsChanges(desiredRegionConfig.ReadOnlySpecs),
					AnalyticsSpecs:      getSpecsChanges(desiredRegionConfig.AnalyticsSpecs),
					AutoScaling:         getAutoScalingChanges(desiredRegionConfig.AutoScaling),
				},
			)
		}

		changedReplicationSpec := &akov2.AdvancedReplicationSpec{
			ZoneName:      desiredReplicationSpec.ZoneName,
			NumShards:     desiredReplicationSpec.NumShards,
			RegionConfigs: changesRegionConfig,
		}
		changesReplicationSpecs = append(changesReplicationSpecs, changedReplicationSpec)
	}

	changes.ReplicationSpecs = changesReplicationSpecs

	d, _ := json.MarshalIndent(desired, "", "  ")
	c, _ := json.MarshalIndent(current, "", "  ")
	fmt.Println("DEBUG >>> ", "DESIRED ", string(d))
	fmt.Println("DEBUG >>> ", "ACTUAL ", string(c))

	b, _ := json.MarshalIndent(changes, "", "  ")
	fmt.Println("DEBUG >>> ", "changes ", string(b))
	os.Exit(0)
	return changes, true
}

func getSpecsChanges(desired *akov2.Specs) *akov2.Specs {
	if desired == nil {
		return nil
	}

	return &akov2.Specs{
		InstanceSize:  desired.InstanceSize,
		NodeCount:     desired.NodeCount,
		EbsVolumeType: pointer.GetOrDefault(&desired.EbsVolumeType, "STANDARD"),
		DiskIOPS:      desired.DiskIOPS,
	}
}

func getAutoScalingChanges(desired *akov2.AdvancedAutoScalingSpec) *akov2.AdvancedAutoScalingSpec {
	if desired == nil {
		return &akov2.AdvancedAutoScalingSpec{
			DiskGB: &akov2.DiskGB{
				Enabled: pointer.MakePtr(false),
			},
			Compute: &akov2.ComputeSpec{
				Enabled: pointer.MakePtr(false),
			},
		}
	}

	return &akov2.AdvancedAutoScalingSpec{
		DiskGB:  desired.DiskGB,
		Compute: desired.Compute,
	}
}

func specAreEqual(desired, current *Cluster) bool {
	if desired.ClusterType != current.ClusterType {
		fmt.Println("DEBUG", "ClusterType", desired.ClusterType, current.ClusterType)
		return false
	}

	if desired.BackupEnabled != nil && !areEqual(desired.BackupEnabled, current.BackupEnabled) {
		fmt.Println("DEBUG", "BackupEnabled", desired.BackupEnabled, current.BackupEnabled)
		return false
	}

	if !reflect.DeepEqual(desired.BiConnector, current.BiConnector) {
		fmt.Println("DEBUG", "BiConnector", desired.BiConnector, current.BiConnector)
		return false
	}

	if desired.DiskSizeGB != nil && !areEqual(desired.DiskSizeGB, current.DiskSizeGB) {
		fmt.Println("DEBUG", "DiskSizeGB", desired.DiskSizeGB, current.DiskSizeGB)
		return false
	}

	if desired.EncryptionAtRestProvider != "" && !areEqual(&desired.EncryptionAtRestProvider, &current.EncryptionAtRestProvider) {
		fmt.Println("DEBUG", "EncryptionAtRestProvider", desired.EncryptionAtRestProvider, current.EncryptionAtRestProvider)
		return false
	}

	if desired.MongoDBMajorVersion != "" && !areEqual(&desired.MongoDBMajorVersion, &current.MongoDBMajorVersion) {
		fmt.Println("DEBUG", "MongoDBMajorVersion", desired.MongoDBMajorVersion, current.MongoDBMajorVersion)
		return false
	}

	if desired.VersionReleaseSystem != "" && !areEqual(&desired.VersionReleaseSystem, &current.VersionReleaseSystem) {
		fmt.Println("DEBUG", "VersionReleaseSystem", desired.VersionReleaseSystem, current.VersionReleaseSystem)
		return false
	}

	if desired.RootCertType != "" && !areEqual(&desired.RootCertType, &current.RootCertType) {
		fmt.Println("DEBUG", "RootCertType", desired.RootCertType, current.RootCertType)
		return false
	}

	if !areEqual(desired.Paused, current.Paused) {
		fmt.Println("DEBUG", "Paused", desired.Paused, current.Paused)
		return false
	}

	if !areEqual(desired.PitEnabled, current.PitEnabled) {
		fmt.Println("DEBUG", "PitEnabled", desired.PitEnabled, current.PitEnabled)
		return false
	}

	if !areEqual(&desired.TerminationProtectionEnabled, &current.TerminationProtectionEnabled) {
		fmt.Println("DEBUG", "TerminationProtectionEnabled", desired.TerminationProtectionEnabled, current.TerminationProtectionEnabled)
		return false
	}

	if !reflect.DeepEqual(desired.Tags, current.Tags) {
		fmt.Println("DEBUG", "Tags", desired.Tags, current.Tags)
		return false
	}

	if !reflect.DeepEqual(desired.Labels, current.Labels) {
		fmt.Println("DEBUG", "Labels", desired.Labels, current.Labels)
		return false
	}

	for ix, desiredReplicationSpec := range desired.ReplicationSpecs {
		if !replicationSpecAreEqual(desiredReplicationSpec, current.ReplicationSpecs[ix], desired.computeAutoscalingEnabled) {
			fmt.Println("DEBUG", "ReplicationSpecs", desiredReplicationSpec, current.ReplicationSpecs[ix])
			return false
		}
	}

	return true
}

func replicationSpecAreEqual(desired, current *akov2.AdvancedReplicationSpec, autoscalingEnabled bool) bool {
	if desired.ZoneName != current.ZoneName {
		fmt.Println("DEBUG", "ZoneName", desired.ZoneName, current.ZoneName)
		return false
	}

	if desired.NumShards != current.NumShards {
		fmt.Println("DEBUG", "NumShards", desired.NumShards, current.NumShards)
		return false
	}

	if len(desired.RegionConfigs) != len(current.RegionConfigs) {
		fmt.Println("DEBUG", "RegionConfigs", desired.RegionConfigs, current.RegionConfigs)
		return false
	}

	for regIx, desiredRegionConfig := range desired.RegionConfigs {
		currentRegionConfig := current.RegionConfigs[regIx]

		if !regionConfigAreEqual(desiredRegionConfig, currentRegionConfig, autoscalingEnabled) {
			fmt.Println("DEBUG", "RegionConfigs", desiredRegionConfig, currentRegionConfig)
			return false
		}
	}

	return true
}

func regionConfigAreEqual(desired, current *akov2.AdvancedRegionConfig, autoscalingEnabled bool) bool {
	if desired.ProviderName != current.ProviderName {
		fmt.Println("DEBUG", "ProviderName", desired.ProviderName, current.ProviderName)
		return false
	}

	if desired.ProviderName == string(provider.ProviderTenant) {
		return (desired.BackingProviderName == current.BackingProviderName) &&
			(desired.ElectableSpecs.InstanceSize == current.ElectableSpecs.InstanceSize)
	}

	if desired.RegionName != current.RegionName {
		fmt.Println("DEBUG", "RegionName", desired.RegionName, current.RegionName)
		return false
	}

	if desired.Priority != nil && !areEqual(desired.Priority, current.Priority) {
		fmt.Println("DEBUG", "Priority", desired.Priority, current.Priority)
		return false
	}

	if !nodeSpecAreEqual(desired.ElectableSpecs, current.ElectableSpecs, autoscalingEnabled) {
		fmt.Println("DEBUG", "ElectableSpecs", desired.ElectableSpecs, current.ElectableSpecs)
		return false
	}

	if !nodeSpecAreEqual(desired.ReadOnlySpecs, current.ReadOnlySpecs, autoscalingEnabled) {
		fmt.Println("DEBUG", "ReadOnlySpecs", desired.ReadOnlySpecs, current.ReadOnlySpecs)
		return false
	}

	if !nodeSpecAreEqual(desired.AnalyticsSpecs, current.AnalyticsSpecs, autoscalingEnabled) {
		fmt.Println("DEBUG", "AnalyticsSpecs", desired.AnalyticsSpecs, current.AnalyticsSpecs)
		return false
	}

	if !autoscalingConfigAreEqual(desired.AutoScaling, current.AutoScaling) {
		fmt.Println("DEBUG", "AutoScaling", desired.AutoScaling, current.AutoScaling)
		return false
	}

	return true
}

func nodeSpecAreEqual(desired, current *akov2.Specs, autoscalingEnabled bool) bool {
	if desired == nil && current == nil {
		fmt.Println("DEBUG", "Specs root", desired, current)
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		fmt.Println("DEBUG", "Specs value", desired, current)
		return false
	}

	if !autoscalingEnabled && desired.InstanceSize != current.InstanceSize {
		fmt.Println("DEBUG", "InstanceSize", desired.InstanceSize, current.InstanceSize)
		return false
	}

	if !areEqual(desired.NodeCount, current.NodeCount) {
		fmt.Println("DEBUG", "NodeCount", desired.NodeCount, current.NodeCount)
		return false
	}

	if desired.EbsVolumeType != "" && desired.EbsVolumeType != current.EbsVolumeType {
		fmt.Println("DEBUG", "EbsVolumeType", desired.EbsVolumeType, current.EbsVolumeType)
		return false
	}

	if desired.DiskIOPS != nil && !areEqual(desired.DiskIOPS, current.DiskIOPS) {
		fmt.Println("DEBUG", "DiskIOPS", desired.DiskIOPS, current.DiskIOPS)
		return false
	}

	return true
}

func autoscalingConfigAreEqual(desired, current *akov2.AdvancedAutoScalingSpec) bool {
	if desired == nil && current == nil {
		fmt.Println("DEBUG", "AutoScalingSpec", desired, current)
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		fmt.Println("DEBUG", "AutoScalingSpec", desired, current)
		return false
	}

	if !diskAutoscalingConfigAreEqual(desired.DiskGB, current.DiskGB) {
		fmt.Println("DEBUG", "DiskGB", desired.DiskGB, current.DiskGB)
		return false
	}

	if !computeAutoscalingConfigAreEqual(desired.Compute, current.Compute) {
		fmt.Println("DEBUG", "ComputeSpec", desired.Compute, current.Compute)
		return false
	}

	return true
}

func diskAutoscalingConfigAreEqual(desired, current *akov2.DiskGB) bool {
	if desired == nil && current == nil {
		fmt.Println("DEBUG", "DiskGB desired & current", desired, current)
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		fmt.Println("DEBUG", "DiskGB values", desired, current)
		return false
	}

	if desired.Enabled != nil && !areEqual(desired.Enabled, current.Enabled) {
		fmt.Println("DEBUG", "DiskGB enabled", desired.Enabled, current.Enabled)
		return false
	}

	return true
}

func computeAutoscalingConfigAreEqual(desired, current *akov2.ComputeSpec) bool {
	if desired == nil && current == nil {
		fmt.Println("DEBUG", "ComputeSpec", desired, current)
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		fmt.Println("DEBUG", "ComputeSpec", desired, current)
		return false
	}

	if desired.Enabled != nil && !areEqual(desired.Enabled, current.Enabled) {
		fmt.Println("DEBUG", "ComputeSpec", desired.Enabled, current.Enabled)
		return false
	}

	if desired.ScaleDownEnabled != nil && !areEqual(desired.ScaleDownEnabled, current.ScaleDownEnabled) {
		fmt.Println("DEBUG", "ComputeSpec", desired.ScaleDownEnabled, current.ScaleDownEnabled)
		return false
	}

	if desired.MinInstanceSize != current.MinInstanceSize {
		fmt.Println("DEBUG", "ComputeSpec", desired.MinInstanceSize, current.MinInstanceSize)
		return false
	}

	if desired.MaxInstanceSize != current.MaxInstanceSize {
		fmt.Println("DEBUG", "ComputeSpec", desired.MaxInstanceSize, current.MaxInstanceSize)
		return false
	}

	return true
}

func areEqual[T comparable](desired, current *T) bool {
	var val1, val2 T

	if desired != nil {
		val1 = *desired
	}

	if current != nil {
		val2 = *current
	}

	return val1 == val2
}
