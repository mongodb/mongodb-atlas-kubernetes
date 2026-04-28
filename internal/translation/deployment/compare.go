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

package deployment

import (
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
				Paused: new(pointer.GetOrDefault(desired.Paused, false)),
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
			ConfigServerManagementMode:   desired.ConfigServerManagementMode,
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
	for specIdx, desiredReplicationSpec := range desired.ReplicationSpecs {
		var currentReplicationSpec *akov2.AdvancedReplicationSpec
		if specIdx < len(current.ReplicationSpecs) {
			currentReplicationSpec = current.ReplicationSpecs[specIdx]
		}
		changesRegionConfig := make([]*akov2.AdvancedRegionConfig, 0, len(desiredReplicationSpec.RegionConfigs))
		for regIdx, desiredRegionConfig := range desiredReplicationSpec.RegionConfigs {
			var currentRegionConfig *akov2.AdvancedRegionConfig
			if currentReplicationSpec != nil && regIdx < len(currentReplicationSpec.RegionConfigs) {
				currentRegionConfig = currentReplicationSpec.RegionConfigs[regIdx]
			}
			regionConfig := &akov2.AdvancedRegionConfig{
				ProviderName:        desiredRegionConfig.ProviderName,
				BackingProviderName: desiredRegionConfig.BackingProviderName,
				RegionName:          desiredRegionConfig.RegionName,
				Priority:            desiredRegionConfig.Priority,
				ElectableSpecs:      getSpecsChanges(desiredRegionConfig.ElectableSpecs, desiredRegionConfig.ProviderName),
				ReadOnlySpecs:       getSpecsChanges(desiredRegionConfig.ReadOnlySpecs, desiredRegionConfig.ProviderName),
				AnalyticsSpecs:      getSpecsChanges(desiredRegionConfig.AnalyticsSpecs, desiredRegionConfig.ProviderName),
			}
			// Only include AutoScaling if it has changed
			if autoScalingChanges := getAutoScalingChanges(desiredRegionConfig.AutoScaling, currentRegionConfig); autoScalingChanges != nil {
				regionConfig.AutoScaling = autoScalingChanges
			}
			changesRegionConfig = append(changesRegionConfig, regionConfig)
		}

		changedReplicationSpec := &akov2.AdvancedReplicationSpec{
			ZoneName:      desiredReplicationSpec.ZoneName,
			NumShards:     desiredReplicationSpec.NumShards,
			RegionConfigs: changesRegionConfig,
		}
		changesReplicationSpecs = append(changesReplicationSpecs, changedReplicationSpec)
	}

	changes.ReplicationSpecs = changesReplicationSpecs

	return changes, true
}

func getSpecsChanges(desired *akov2.Specs, providerName string) *akov2.Specs {
	if desired == nil {
		return nil
	}

	specs := &akov2.Specs{
		InstanceSize: desired.InstanceSize,
		NodeCount:    desired.NodeCount,
		DiskIOPS:     desired.DiskIOPS,
	}

	// Only include EbsVolumeType when:
	// 1. It's explicitly set (non-empty), OR
	// 2. Provider is AWS (EbsVolumeType is only valid for AWS)
	// This prevents sending EbsVolumeType="STANDARD" to GCP/Azure clusters, which causes reconcile loops
	if desired.EbsVolumeType != "" || providerName == "AWS" {
		specs.EbsVolumeType = pointer.GetOrDefault(&desired.EbsVolumeType, "STANDARD")
	}

	return specs
}

func getAutoScalingChanges(desired *akov2.AdvancedAutoScalingSpec, current *akov2.AdvancedRegionConfig) *akov2.AdvancedAutoScalingSpec {
	var currentAutoScaling *akov2.AdvancedAutoScalingSpec
	if current != nil {
		currentAutoScaling = current.AutoScaling
	}

	// If desired and current are equal, return nil (no changes)
	if autoscalingConfigAreEqual(desired, currentAutoScaling) {
		return nil
	}

	// If desired is nil but current is not, we need to disable autoscaling
	if desired == nil {
		return &akov2.AdvancedAutoScalingSpec{
			DiskGB: &akov2.DiskGB{
				Enabled: new(false),
			},
			Compute: &akov2.ComputeSpec{
				Enabled: new(false),
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
		return false
	}

	if desired.BackupEnabled != nil && !areEqual(desired.BackupEnabled, current.BackupEnabled) {
		return false
	}

	if !reflect.DeepEqual(desired.BiConnector, current.BiConnector) {
		return false
	}

	if desired.DiskSizeGB != nil && !areEqual(desired.DiskSizeGB, current.DiskSizeGB) {
		return false
	}

	if desired.EncryptionAtRestProvider != "" && !areEqual(&desired.EncryptionAtRestProvider, &current.EncryptionAtRestProvider) {
		return false
	}

	if desired.ConfigServerManagementMode != "" && !areEqual(&desired.ConfigServerManagementMode, &current.ConfigServerManagementMode) {
		return false
	}

	if desired.MongoDBMajorVersion != "" && !areEqual(&desired.MongoDBMajorVersion, &current.MongoDBMajorVersion) {
		return false
	}

	if desired.VersionReleaseSystem != "" && !areEqual(&desired.VersionReleaseSystem, &current.VersionReleaseSystem) {
		return false
	}

	if desired.RootCertType != "" && !areEqual(&desired.RootCertType, &current.RootCertType) {
		return false
	}

	if !areEqual(desired.Paused, current.Paused) {
		return false
	}

	if !areEqual(desired.PitEnabled, current.PitEnabled) {
		return false
	}

	if !areEqual(&desired.TerminationProtectionEnabled, &current.TerminationProtectionEnabled) {
		return false
	}

	if !reflect.DeepEqual(desired.Tags, current.Tags) {
		return false
	}

	if !reflect.DeepEqual(desired.Labels, current.Labels) {
		return false
	}

	for ix, desiredReplicationSpec := range desired.ReplicationSpecs {
		if desired.ClusterType == string(akov2.TypeSharded) {
			if desiredReplicationSpec.NumShards != len(current.ReplicationSpecs) {
				return false
			}
		}
		if !replicationSpecAreEqual(desiredReplicationSpec, current.ReplicationSpecs[ix], desired.computeAutoscalingEnabled) {
			return false
		}
	}

	return true
}

func replicationSpecAreEqual(desired, current *akov2.AdvancedReplicationSpec, autoscalingEnabled bool) bool {
	if desired.ZoneName != current.ZoneName {
		return false
	}

	if len(desired.RegionConfigs) != len(current.RegionConfigs) {
		return false
	}

	for regIx, desiredRegionConfig := range desired.RegionConfigs {
		currentRegionConfig := current.RegionConfigs[regIx]

		if !regionConfigAreEqual(desiredRegionConfig, currentRegionConfig, autoscalingEnabled) {
			return false
		}
	}

	return true
}

func regionConfigAreEqual(desired, current *akov2.AdvancedRegionConfig, autoscalingEnabled bool) bool {
	if desired.ProviderName != current.ProviderName {
		return false
	}

	if desired.ProviderName == string(provider.ProviderTenant) {
		return (desired.BackingProviderName == current.BackingProviderName) &&
			(desired.ElectableSpecs.InstanceSize == current.ElectableSpecs.InstanceSize)
	}

	if desired.RegionName != current.RegionName {
		return false
	}

	if desired.Priority != nil && !areEqual(desired.Priority, current.Priority) {
		return false
	}

	if !nodeSpecAreEqual(desired.ElectableSpecs, current.ElectableSpecs, autoscalingEnabled) {
		return false
	}

	if !nodeSpecAreEqual(desired.ReadOnlySpecs, current.ReadOnlySpecs, autoscalingEnabled) {
		return false
	}

	if !nodeSpecAreEqual(desired.AnalyticsSpecs, current.AnalyticsSpecs, autoscalingEnabled) {
		return false
	}

	if !autoscalingConfigAreEqual(desired.AutoScaling, current.AutoScaling) {
		return false
	}

	return true
}

func nodeSpecAreEqual(desired, current *akov2.Specs, autoscalingEnabled bool) bool {
	if desired == nil && current == nil {
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		return false
	}

	if !autoscalingEnabled && desired.InstanceSize != current.InstanceSize {
		return false
	}

	if !areEqual(desired.NodeCount, current.NodeCount) {
		return false
	}

	if desired.EbsVolumeType != "" && desired.EbsVolumeType != current.EbsVolumeType {
		return false
	}

	if desired.DiskIOPS != nil && !areEqual(desired.DiskIOPS, current.DiskIOPS) {
		return false
	}

	return true
}

func autoscalingConfigAreEqual(desired, current *akov2.AdvancedAutoScalingSpec) bool {
	if desired == nil && current == nil {
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		return false
	}

	if !diskAutoscalingConfigAreEqual(desired.DiskGB, current.DiskGB) {
		return false
	}

	if !computeAutoscalingConfigAreEqual(desired.Compute, current.Compute) {
		return false
	}

	return true
}

func diskAutoscalingConfigAreEqual(desired, current *akov2.DiskGB) bool {
	if desired == nil && current == nil {
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		return false
	}

	if desired.Enabled != nil && !areEqual(desired.Enabled, current.Enabled) {
		return false
	}

	return true
}

func computeAutoscalingConfigAreEqual(desired, current *akov2.ComputeSpec) bool {
	if desired == nil && current == nil {
		return true
	}

	if (desired != nil && current == nil) || (desired == nil && current != nil) {
		return false
	}

	if desired.Enabled != nil && !areEqual(desired.Enabled, current.Enabled) {
		return false
	}

	if desired.ScaleDownEnabled != nil && !areEqual(desired.ScaleDownEnabled, current.ScaleDownEnabled) {
		return false
	}

	scaleDown := current.ScaleDownEnabled != nil && *current.ScaleDownEnabled
	if scaleDown && (desired.MinInstanceSize != current.MinInstanceSize) {
		return false
	}

	if desired.MaxInstanceSize != current.MaxInstanceSize {
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

// ProcessArgsEqual reports whether the AKO-side processArgs and the
// Atlas-side processArgs are semantically equivalent for reconcile purposes.
//
// Fields that are the zero value on the AKO side are treated as "no opinion"
// and not compared. This mirrors processArgsToAtlas, which omits zero values
// from the PATCH body via MakePtrOrNil — if a field would not be sent over
// the wire, it must not drive an update either, or the reconciler loops on
// Atlas-populated server defaults (issue #3142).
func ProcessArgsEqual(ako, atlas *akov2.ProcessArgs) bool {
	if ako == nil {
		return true
	}
	if atlas == nil {
		atlas = &akov2.ProcessArgs{}
	}
	if ako.DefaultWriteConcern != "" && ako.DefaultWriteConcern != atlas.DefaultWriteConcern {
		return false
	}
	if ako.MinimumEnabledTLSProtocol != "" && ako.MinimumEnabledTLSProtocol != atlas.MinimumEnabledTLSProtocol {
		return false
	}
	if ako.OplogMinRetentionHours != "" && ako.OplogMinRetentionHours != atlas.OplogMinRetentionHours {
		return false
	}
	if ako.JavascriptEnabled != nil && !areEqual(ako.JavascriptEnabled, atlas.JavascriptEnabled) {
		return false
	}
	if ako.NoTableScan != nil && !areEqual(ako.NoTableScan, atlas.NoTableScan) {
		return false
	}
	if ako.OplogSizeMB != nil && !areEqual(ako.OplogSizeMB, atlas.OplogSizeMB) {
		return false
	}
	if ako.SampleSizeBIConnector != nil && !areEqual(ako.SampleSizeBIConnector, atlas.SampleSizeBIConnector) {
		return false
	}
	if ako.SampleRefreshIntervalBIConnector != nil && !areEqual(ako.SampleRefreshIntervalBIConnector, atlas.SampleRefreshIntervalBIConnector) {
		return false
	}
	return true
}
