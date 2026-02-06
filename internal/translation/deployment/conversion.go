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
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/tag"
)

const NOTIFICATION_REASON_DEPRECATION = "DeprecationWarning"
const NOTIFICATION_REASON_RECOMMENDATION = "RecommendationWarning"

type Deployment interface {
	GetName() string
	GetProjectID() string
	GetCustomResource() *akov2.AtlasDeployment
	GetState() string
	GetMongoDBVersion() string
	GetConnection() *status.ConnectionStrings
	GetReplicaSet() []status.ReplicaSet
	IsServerless() bool
	IsFlex() bool
	IsTenant() bool
	IsDedicated() bool
	Notifications() (bool, string, string)
}

type Cluster struct {
	*akov2.AdvancedDeploymentSpec
	ProjectID      string
	State          string
	MongoDBVersion string
	Connection     *status.ConnectionStrings
	ProcessArgs    *akov2.ProcessArgs
	ReplicaSet     []status.ReplicaSet
	ZoneID         string

	customResource            *akov2.AtlasDeployment
	computeAutoscalingEnabled bool
	instanceSizeOverride      string
	isTenant                  bool
}

func (c *Cluster) GetName() string {
	return c.Name
}

func (c *Cluster) GetProjectID() string {
	return c.ProjectID
}

func (c *Cluster) GetState() string {
	return c.State
}

func (c *Cluster) GetMongoDBVersion() string {
	return c.MongoDBVersion
}

func (c *Cluster) GetConnection() *status.ConnectionStrings {
	return c.Connection
}

func (c *Cluster) GetReplicaSet() []status.ReplicaSet {
	return c.ReplicaSet
}

func (c *Cluster) GetCustomResource() *akov2.AtlasDeployment {
	return c.customResource
}

func (c *Cluster) IsServerless() bool {
	return false
}

func (c *Cluster) IsFlex() bool {
	return false
}

func (c *Cluster) IsTenant() bool {
	return c.isTenant
}

func (c *Cluster) IsDedicated() bool {
	return !c.IsTenant()
}

func (c *Cluster) Notifications() (bool, string, string) {
	for _, replicationSpec := range c.ReplicationSpecs {
		if replicationSpec == nil {
			continue
		}

		for _, regionConfig := range replicationSpec.RegionConfigs {
			if regionConfig == nil {
				continue
			}

			if deprecatedSpecs(regionConfig.ElectableSpecs) ||
				deprecatedSpecs(regionConfig.ReadOnlySpecs) ||
				deprecatedSpecs(regionConfig.AnalyticsSpecs) {
				return true, NOTIFICATION_REASON_DEPRECATION, "WARNING: M2 and M5 instance sizes are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details."
			}
		}
	}

	processArgs := c.customResource.Spec.ProcessArgs
	if processArgs != nil {
		if processArgs.DefaultReadConcern != "" {
			return true, NOTIFICATION_REASON_DEPRECATION, "Process Arg DefaultReadConcern is no longer available in Atlas. Setting this will have no effect."
		}
		if processArgs.FailIndexKeyTooLong != nil {
			return true, NOTIFICATION_REASON_DEPRECATION, "Process Arg FailIndexKeyTooLong is no longer available in Atlas. Setting this will have no effect."
		}
	}

	if c.IsDedicated() && c.customResource.Spec.UpgradeToDedicated {
		return true, NOTIFICATION_REASON_RECOMMENDATION, "Cluster is already dedicated. It’s recommended to remove or set the upgrade flag to false"
	}
	return false, "", ""
}

func deprecatedSpecs(specs *akov2.Specs) bool {
	if specs == nil {
		return false
	}
	if specs.InstanceSize == "M2" || specs.InstanceSize == "M5" {
		return true
	}
	return false
}

type Flex struct {
	*akov2.FlexSpec
	ProjectID      string
	State          string
	MongoDBVersion string
	Connection     *status.ConnectionStrings

	customResource *akov2.AtlasDeployment
}

func (f *Flex) GetName() string {
	return f.Name
}

func (f *Flex) GetProjectID() string {
	return f.ProjectID
}

func (f *Flex) GetState() string {
	return f.State
}

func (f *Flex) GetMongoDBVersion() string {
	return f.MongoDBVersion
}

func (f *Flex) GetConnection() *status.ConnectionStrings {
	return f.Connection
}

func (f *Flex) GetReplicaSet() []status.ReplicaSet {
	return nil
}

func (f *Flex) GetCustomResource() *akov2.AtlasDeployment {
	return f.customResource
}

func (f *Flex) IsServerless() bool {
	return false
}

func (f *Flex) IsFlex() bool {
	return true
}

func (f *Flex) IsTenant() bool {
	return false
}

func (f *Flex) IsDedicated() bool {
	return false
}

func (f *Flex) Notifications() (bool, string, string) {
	if f.customResource.IsServerless() {
		return true, NOTIFICATION_REASON_DEPRECATION, "WARNING: Serverless is deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details."
	}

	return false, "", ""
}

type Connection struct {
	Name             string
	ConnURL          string
	SrvConnURL       string
	PrivateURL       string
	SrvPrivateURL    string
	Serverless       bool
	PrivateEndpoints []PrivateEndpoint
}

type PrivateEndpoint struct {
	URL       string
	ServerURL string
	ShardURL  string
	Endpoint  []Endpoint
}

type Endpoint struct {
	ID       string
	Provider string
	Region   string
}

func NewDeployment(projectID string, atlasDeployment *akov2.AtlasDeployment) Deployment {
	if atlasDeployment.IsServerless() {
		flex := &Flex{
			customResource: atlasDeployment,
			ProjectID:      projectID,
			FlexSpec:       serverlessToFlexSpec(atlasDeployment.Spec.ServerlessSpec.DeepCopy()),
		}
		normalizeFlexDeployment(flex)

		return flex
	}

	if atlasDeployment.IsFlex() {
		flex := &Flex{
			customResource: atlasDeployment,
			ProjectID:      projectID,
			FlexSpec:       atlasDeployment.Spec.FlexSpec.DeepCopy(),
		}
		normalizeFlexDeployment(flex)
		return flex
	}

	cluster := &Cluster{
		customResource:         atlasDeployment,
		ProjectID:              projectID,
		AdvancedDeploymentSpec: atlasDeployment.Spec.DeploymentSpec.DeepCopy(),
		ProcessArgs:            atlasDeployment.Spec.ProcessArgs.DeepCopy(),
	}
	normalizeClusterDeployment(cluster)

	return cluster
}

func serverlessToFlexSpec(serverless *akov2.ServerlessSpec) *akov2.FlexSpec {
	settings := &akov2.FlexProviderSettings{}
	if serverless.ProviderSettings != nil {
		settings.BackingProviderName = serverless.ProviderSettings.BackingProviderName
		settings.RegionName = serverless.ProviderSettings.RegionName
	}

	return &akov2.FlexSpec{
		Name:                         serverless.Name,
		Tags:                         serverless.Tags,
		TerminationProtectionEnabled: serverless.TerminationProtectionEnabled,
		ProviderSettings:             settings,
	}
}

func normalizeFlexDeployment(flex *Flex) {
	if flex.FlexSpec.Tags == nil {
		flex.FlexSpec.Tags = []*akov2.TagSpec{}
	}
	cmp.NormalizeSlice(flex.Tags, func(a, b *akov2.TagSpec) int {
		return strings.Compare(a.Key, b.Key)
	})
}

func normalizeClusterDeployment(cluster *Cluster) {
	isTenant, computeAutoscalingEnabled, instanceSizeOverride := getAutoscalingOverride(cluster.ReplicationSpecs)
	cluster.computeAutoscalingEnabled = computeAutoscalingEnabled
	cluster.instanceSizeOverride = instanceSizeOverride
	cluster.isTenant = isTenant

	if cluster.ClusterType == "" {
		cluster.ClusterType = "REPLICASET"
	}

	cluster.Paused = pointer.GetOrPointerToDefault(cluster.Paused, false)
	if cluster.VersionReleaseSystem == "" {
		cluster.VersionReleaseSystem = "LTS"
	}

	if cluster.RootCertType == "" {
		cluster.RootCertType = "ISRGROOTX1"
	}

	if !isTenant {
		cluster.BackupEnabled = pointer.GetOrPointerToDefault(cluster.BackupEnabled, false)
		cluster.PitEnabled = pointer.GetOrPointerToDefault(cluster.PitEnabled, false)

		if cluster.EncryptionAtRestProvider == "" {
			cluster.EncryptionAtRestProvider = "NONE"
		}
	}

	if cluster.BiConnector != nil && cluster.BiConnector.Enabled != nil && !*cluster.BiConnector.Enabled {
		cluster.BiConnector = nil
	}

	if cluster.AdvancedDeploymentSpec.Tags == nil {
		cluster.AdvancedDeploymentSpec.Tags = []*akov2.TagSpec{}
	}

	cmp.NormalizeSlice(cluster.Tags, func(a, b *akov2.TagSpec) int {
		return strings.Compare(a.Key, b.Key)
	})

	cmp.NormalizeSlice(cluster.Labels, func(a, b common.LabelSpec) int {
		return strings.Compare(a.Key, b.Key)
	})

	normalizeReplicationSpecs(cluster, isTenant)
	normalizeProcessArgs(cluster.ProcessArgs)
}

func normalizeReplicationSpecs(cluster *Cluster, isTenant bool) {
	for ix, replicationSpec := range cluster.ReplicationSpecs {
		if replicationSpec == nil {
			continue
		}
		if replicationSpec.NumShards == 0 {
			replicationSpec.NumShards = 1
		}
		if replicationSpec.ZoneName == "" {
			replicationSpec.ZoneName = fmt.Sprintf("Zone %d", ix+1)
		}

		normalizeRegionConfigs(replicationSpec.RegionConfigs, isTenant)
	}
	cmp.NormalizeSlice(cluster.ReplicationSpecs, func(a, b *akov2.AdvancedReplicationSpec) int {
		var zoneA, zoneB string
		if a != nil {
			zoneA = a.ZoneName
		}
		if b != nil {
			zoneB = b.ZoneName
		}
		return strings.Compare(zoneA, zoneB)
	})
}

func compareRegionConfigs(a, b *akov2.AdvancedRegionConfig) int {
	aPriority := 0
	if a != nil && a.Priority != nil {
		aPriority = *a.Priority
	}
	bPriority := 0
	if b != nil && b.Priority != nil {
		bPriority = *b.Priority
	}
	var aProviderRegion, bProviderRegion string
	if a != nil {
		aProviderRegion = a.ProviderName + a.RegionName
	}
	if b != nil {
		bProviderRegion = b.ProviderName + b.RegionName
	}
	if aPriority < bPriority {
		return 1
	}
	if aPriority > bPriority {
		return -1
	}
	return strings.Compare(bProviderRegion, aProviderRegion)
}

func normalizeRegionConfigs(regionConfigs []*akov2.AdvancedRegionConfig, isTenant bool) {
	cmp.NormalizeSlice(regionConfigs, compareRegionConfigs)

	for _, regionConfig := range regionConfigs {
		if regionConfig == nil {
			continue
		}
		if regionConfig.ProviderName != string(provider.ProviderTenant) {
			regionConfig.BackingProviderName = ""
		}
		if regionConfig.ElectableSpecs != nil {
			if !isTenant && (regionConfig.ElectableSpecs.NodeCount == nil || *regionConfig.ElectableSpecs.NodeCount == 0) {
				regionConfig.ElectableSpecs = nil
			}
		}
		if regionConfig.ReadOnlySpecs != nil && (regionConfig.ReadOnlySpecs.NodeCount == nil || *regionConfig.ReadOnlySpecs.NodeCount == 0) {
			regionConfig.ReadOnlySpecs = nil
		}
		if regionConfig.AnalyticsSpecs != nil && (regionConfig.AnalyticsSpecs.NodeCount == nil || *regionConfig.AnalyticsSpecs.NodeCount == 0) {
			regionConfig.AnalyticsSpecs = nil
		}

		computeUnsetOrDisabled := regionConfig.AutoScaling == nil || regionConfig.AutoScaling.Compute == nil ||
			regionConfig.AutoScaling.Compute.Enabled == nil || !*regionConfig.AutoScaling.Compute.Enabled
		diskUnsetOrDisabled := regionConfig.AutoScaling == nil || regionConfig.AutoScaling.DiskGB == nil ||
			regionConfig.AutoScaling.DiskGB.Enabled == nil || !*regionConfig.AutoScaling.DiskGB.Enabled

		if regionConfig.AutoScaling == nil {
			regionConfig.AutoScaling = &akov2.AdvancedAutoScalingSpec{}
		}

		if computeUnsetOrDisabled {
			regionConfig.AutoScaling.Compute = &akov2.ComputeSpec{
				Enabled: pointer.MakePtr(false),
			}
		}

		if diskUnsetOrDisabled {
			regionConfig.AutoScaling.DiskGB = &akov2.DiskGB{
				Enabled: pointer.MakePtr(false),
			}
		}
	}
}

func normalizeProcessArgs(args *akov2.ProcessArgs) {
	if args == nil {
		return
	}

	if args.JavascriptEnabled == nil {
		args.JavascriptEnabled = pointer.MakePtr(true)
	}

	if args.MinimumEnabledTLSProtocol == "" {
		args.MinimumEnabledTLSProtocol = "TLS1_2"
	}

	if args.NoTableScan == nil {
		args.NoTableScan = pointer.MakePtr(false)
	}

	// those are ignored fields nowadays
	args.DefaultReadConcern = ""
	args.FailIndexKeyTooLong = nil
}

func getAutoscalingOverride(replications []*akov2.AdvancedReplicationSpec) (bool, bool, string) {
	var instanceSize string
	var isTenant bool
	for _, replica := range replications {
		if replica == nil {
			continue
		}
		for _, region := range replica.RegionConfigs {
			if region == nil {
				continue
			}
			if region.ProviderName == string(provider.ProviderTenant) {
				isTenant = true
			}

			if region.ElectableSpecs != nil {
				instanceSize = region.ElectableSpecs.InstanceSize
			}
			if region.ReadOnlySpecs != nil {
				instanceSize = region.ReadOnlySpecs.InstanceSize
			}
			if region.AnalyticsSpecs != nil {
				instanceSize = region.AnalyticsSpecs.InstanceSize
			}

			if region.AutoScaling != nil &&
				region.AutoScaling.Compute != nil &&
				region.AutoScaling.Compute.Enabled != nil &&
				*region.AutoScaling.Compute.Enabled {
				return isTenant, true, instanceSize
			}
		}
	}

	return isTenant, false, ""
}

func clusterFromAtlas(clusterDesc *admin.ClusterDescription20240805) *Cluster {
	connectionStrings := clusterDesc.GetConnectionStrings()
	pes := make([]status.PrivateEndpoint, 0, len(connectionStrings.GetPrivateEndpoint()))
	for _, pe := range connectionStrings.GetPrivateEndpoint() {
		eps := make([]status.Endpoint, 0, len(pe.GetEndpoints()))
		for _, ep := range pe.GetEndpoints() {
			eps = append(
				eps,
				status.Endpoint{
					EndpointID:   ep.GetEndpointId(),
					ProviderName: ep.GetProviderName(),
					Region:       ep.GetRegion(),
				})
		}

		pes = append(
			pes,
			status.PrivateEndpoint{
				ConnectionString:                  pe.GetConnectionString(),
				SRVConnectionString:               pe.GetSrvConnectionString(),
				SRVShardOptimizedConnectionString: pe.GetSrvShardOptimizedConnectionString(),
				Endpoints:                         eps,
			})
	}

	cluster := &Cluster{
		ProjectID:      clusterDesc.GetGroupId(),
		State:          clusterDesc.GetStateName(),
		MongoDBVersion: clusterDesc.GetMongoDBVersion(),
		Connection: &status.ConnectionStrings{
			Standard:        connectionStrings.GetStandard(),
			StandardSrv:     connectionStrings.GetStandardSrv(),
			Private:         connectionStrings.GetPrivate(),
			PrivateSrv:      connectionStrings.GetPrivateSrv(),
			PrivateEndpoint: pes,
		},
		ReplicaSet: replicaSetFromAtlas(clusterDesc.GetReplicationSpecs()),
		AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:                         clusterDesc.GetName(),
			ClusterType:                  clusterDesc.GetClusterType(),
			MongoDBMajorVersion:          clusterDesc.GetMongoDBMajorVersion(),
			MongoDBVersion:               clusterDesc.GetMongoDBVersion(),
			VersionReleaseSystem:         clusterDesc.GetVersionReleaseSystem(),
			DiskSizeGB:                   diskSizeFromAtlas(clusterDesc),
			BackupEnabled:                clusterDesc.BackupEnabled,
			BiConnector:                  biConnectFromAtlas(clusterDesc.BiConnector),
			EncryptionAtRestProvider:     clusterDesc.GetEncryptionAtRestProvider(),
			Labels:                       labelsFromAtlas(clusterDesc.GetLabels()),
			Paused:                       clusterDesc.Paused,
			PitEnabled:                   clusterDesc.PitEnabled,
			ReplicationSpecs:             replicationSpecFromAtlas(clusterDesc.GetReplicationSpecs()),
			RootCertType:                 clusterDesc.GetRootCertType(),
			Tags:                         tag.FromAtlas(clusterDesc.GetTags()),
			CustomZoneMapping:            nil,
			ManagedNamespaces:            nil,
			TerminationProtectionEnabled: clusterDesc.GetTerminationProtectionEnabled(),
			SearchNodes:                  nil,
			SearchIndexes:                nil,
			ConfigServerManagementMode:   clusterDesc.GetConfigServerManagementMode(),
		},
	}
	normalizeClusterDeployment(cluster)

	if len(clusterDesc.GetReplicationSpecs()) > 0 {
		cluster.ZoneID = clusterDesc.GetReplicationSpecs()[0].GetZoneId()
	}

	return cluster
}

func clusterCreateToAtlas(cluster *Cluster) *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
		Name:                         pointer.MakePtrOrNil(cluster.Name),
		ClusterType:                  pointer.MakePtrOrNil(cluster.ClusterType),
		MongoDBMajorVersion:          pointer.MakePtrOrNil(cluster.MongoDBMajorVersion),
		VersionReleaseSystem:         pointer.MakePtrOrNil(cluster.VersionReleaseSystem),
		BackupEnabled:                cluster.BackupEnabled,
		BiConnector:                  biConnectToAtlas(cluster.BiConnector),
		EncryptionAtRestProvider:     pointer.MakePtrOrNil(cluster.EncryptionAtRestProvider),
		Labels:                       labelsToAtlas(cluster.Labels),
		Paused:                       cluster.Paused,
		PitEnabled:                   cluster.PitEnabled,
		ReplicationSpecs:             replicationSpecToAtlas(cluster.ReplicationSpecs, cluster.ClusterType, cluster.DiskSizeGB),
		RootCertType:                 pointer.MakePtrOrNil(cluster.RootCertType),
		Tags:                         tag.ToAtlas(cluster.Tags),
		TerminationProtectionEnabled: pointer.MakePtrOrNil(cluster.TerminationProtectionEnabled),
		ConfigServerManagementMode:   pointer.MakePtrOrNil(cluster.ConfigServerManagementMode),
	}
}

func clusterUpdateToAtlas(cluster *Cluster) *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
		ClusterType:                  pointer.MakePtrOrNil(cluster.ClusterType),
		MongoDBMajorVersion:          pointer.MakePtrOrNil(cluster.MongoDBMajorVersion),
		VersionReleaseSystem:         pointer.MakePtrOrNil(cluster.VersionReleaseSystem),
		BackupEnabled:                cluster.BackupEnabled,
		BiConnector:                  biConnectToAtlas(cluster.BiConnector),
		EncryptionAtRestProvider:     pointer.MakePtrOrNil(cluster.EncryptionAtRestProvider),
		Labels:                       labelsToAtlas(cluster.Labels),
		Paused:                       cluster.Paused,
		PitEnabled:                   cluster.PitEnabled,
		ReplicationSpecs:             replicationSpecToAtlas(cluster.ReplicationSpecs, cluster.ClusterType, cluster.DiskSizeGB),
		RootCertType:                 pointer.MakePtrOrNil(cluster.RootCertType),
		Tags:                         tag.ToAtlas(cluster.Tags),
		TerminationProtectionEnabled: pointer.MakePtrOrNil(cluster.TerminationProtectionEnabled),
		ConfigServerManagementMode:   pointer.MakePtrOrNil(cluster.ConfigServerManagementMode),
	}
}

func replicaSetFromAtlas(replicationSpecs []admin.ReplicationSpec20240805) []status.ReplicaSet {
	replicaSet := make([]status.ReplicaSet, 0, len(replicationSpecs))
	for _, replicationSpec := range replicationSpecs {
		replicaSet = append(
			replicaSet,
			status.ReplicaSet{
				ID:       replicationSpec.GetId(),
				ZoneName: replicationSpec.GetZoneName(),
			},
		)
	}

	return replicaSet
}

func diskSizeFromAtlas(cluster *admin.ClusterDescription20240805) *int {
	var value float64

	if specs := cluster.GetReplicationSpecs(); len(specs) > 0 {
		if configs := specs[0].GetRegionConfigs(); len(configs) > 0 {
			if e, ok := configs[0].GetElectableSpecsOk(); ok {
				value = e.GetDiskSizeGB()
			} else if r, ok := configs[0].GetReadOnlySpecsOk(); ok {
				value = r.GetDiskSizeGB()
			} else if a, ok := configs[0].GetAnalyticsSpecsOk(); ok {
				value = a.GetDiskSizeGB()
			}
		}
	}

	if value >= 1 {
		return pointer.MakePtr(int(value))
	}

	return nil
}

func biConnectFromAtlas(conn *admin.BiConnector) *akov2.BiConnectorSpec {
	if conn == nil {
		return nil
	}

	return &akov2.BiConnectorSpec{
		Enabled:        conn.Enabled,
		ReadPreference: conn.GetReadPreference(),
	}
}

func labelsFromAtlas(cLabels []admin.ComponentLabel) []common.LabelSpec {
	if len(cLabels) == 0 {
		return nil
	}
	labels := make([]common.LabelSpec, 0, len(cLabels))
	for _, cLabel := range cLabels {
		labels = append(
			labels,
			common.LabelSpec{
				Key:   cLabel.GetKey(),
				Value: cLabel.GetValue(),
			},
		)
	}

	return labels
}

func replicationSpecFromAtlas(replicationSpecs []admin.ReplicationSpec20240805) []*akov2.AdvancedReplicationSpec {
	hSpecOrDefault := func(spec admin.HardwareSpec20240805, providerName string) *akov2.Specs {
		if spec.GetNodeCount() == 0 && providerName != string(provider.ProviderTenant) {
			return nil
		}

		return &akov2.Specs{
			InstanceSize:  spec.GetInstanceSize(),
			NodeCount:     spec.NodeCount,
			EbsVolumeType: spec.GetEbsVolumeType(),
			DiskIOPS:      pointer.MakePtrOrNil(int64(spec.GetDiskIOPS())),
		}
	}
	dHSpecOrDefault := func(spec admin.DedicatedHardwareSpec20240805) *akov2.Specs {
		if spec.GetNodeCount() == 0 {
			return nil
		}

		return &akov2.Specs{
			InstanceSize:  spec.GetInstanceSize(),
			NodeCount:     spec.NodeCount,
			EbsVolumeType: spec.GetEbsVolumeType(),
			DiskIOPS:      pointer.MakePtrOrNil(int64(spec.GetDiskIOPS())),
		}
	}
	autoScalingOrDefault := func(spec admin.AdvancedAutoScalingSettings) *akov2.AdvancedAutoScalingSpec {
		compute := spec.GetCompute()
		diskGB := spec.GetDiskGB()
		if !compute.GetEnabled() && !diskGB.GetEnabled() {
			return nil
		}

		autoscaling := &akov2.AdvancedAutoScalingSpec{}

		if compute.GetEnabled() {
			autoscaling.Compute = &akov2.ComputeSpec{
				Enabled:          compute.Enabled,
				ScaleDownEnabled: compute.ScaleDownEnabled,
				MinInstanceSize:  compute.GetMinInstanceSize(),
				MaxInstanceSize:  compute.GetMaxInstanceSize(),
			}
		}

		if diskGB.GetEnabled() {
			autoscaling.DiskGB = &akov2.DiskGB{
				Enabled: diskGB.Enabled,
			}
		}

		return autoscaling
	}

	specs := make([]*akov2.AdvancedReplicationSpec, 0, len(replicationSpecs))
	for _, spec := range replicationSpecs {
		regionConfigs := make([]*akov2.AdvancedRegionConfig, 0, len(spec.GetRegionConfigs()))
		for _, regionConfig := range spec.GetRegionConfigs() {
			regionConfigs = append(
				regionConfigs,
				&akov2.AdvancedRegionConfig{
					ProviderName:        regionConfig.GetProviderName(),
					BackingProviderName: regionConfig.GetBackingProviderName(),
					RegionName:          regionConfig.GetRegionName(),
					Priority:            regionConfig.Priority,
					ElectableSpecs:      hSpecOrDefault(regionConfig.GetElectableSpecs(), regionConfig.GetProviderName()),
					ReadOnlySpecs:       dHSpecOrDefault(regionConfig.GetReadOnlySpecs()),
					AnalyticsSpecs:      dHSpecOrDefault(regionConfig.GetAnalyticsSpecs()),
					AutoScaling:         autoScalingOrDefault(regionConfig.GetAutoScaling()),
				},
			)
		}

		specs = append(
			specs,
			&akov2.AdvancedReplicationSpec{
				ZoneName:      spec.GetZoneName(),
				RegionConfigs: regionConfigs,
			},
		)
	}

	return specs
}

func processArgsFromAtlas(config *admin.ClusterDescriptionProcessArgs20240805) *akov2.ProcessArgs {
	oplogMinRetentionHours := ""
	if config.GetOplogMinRetentionHours() > 0 {
		oplogMinRetentionHours = strconv.FormatFloat(config.GetOplogMinRetentionHours(), 'f', -1, 64)
	}

	args := akov2.ProcessArgs{
		DefaultWriteConcern:              config.GetDefaultWriteConcern(),
		MinimumEnabledTLSProtocol:        config.GetMinimumEnabledTlsProtocol(),
		JavascriptEnabled:                config.JavascriptEnabled,
		NoTableScan:                      pointer.MakePtr(pointer.GetOrDefault(config.NoTableScan, false)),
		OplogSizeMB:                      pointer.MakePtrOrNil(int64(pointer.GetOrDefault(config.OplogSizeMB, 0))),
		SampleSizeBIConnector:            pointer.MakePtrOrNil(int64(pointer.GetOrDefault(config.SampleSizeBIConnector, 0))),
		SampleRefreshIntervalBIConnector: pointer.MakePtrOrNil(int64(pointer.GetOrDefault(config.SampleRefreshIntervalBIConnector, 0))),
		OplogMinRetentionHours:           oplogMinRetentionHours,
	}
	normalizeProcessArgs(&args)

	return &args
}

func biConnectToAtlas(conn *akov2.BiConnectorSpec) *admin.BiConnector {
	if conn == nil {
		return nil
	}

	return &admin.BiConnector{
		Enabled:        conn.Enabled,
		ReadPreference: &conn.ReadPreference,
	}
}

func labelsToAtlas(labels []common.LabelSpec) *[]admin.ComponentLabel {
	if len(labels) == 0 {
		return nil
	}

	cLabels := make([]admin.ComponentLabel, 0, len(labels))
	for _, label := range labels {
		labelKey := label.Key
		labelValue := label.Value
		cLabels = append(
			cLabels,
			admin.ComponentLabel{
				Key:   &labelKey,
				Value: &labelValue,
			},
		)
	}

	return &cLabels
}

func replicationSpecToAtlas(replicationSpecs []*akov2.AdvancedReplicationSpec, clusterType string, diskSize *int) *[]admin.ReplicationSpec20240805 {
	if len(replicationSpecs) == 0 {
		return nil
	}

	var diskSizeGB *float64
	if diskSize != nil {
		diskSizeGB = pointer.MakePtr(float64(*diskSize))
	}

	hSpecOrDefault := func(spec *akov2.Specs) *admin.HardwareSpec20240805 {
		if spec == nil {
			return nil
		}

		var diskIOPs *int
		if spec.DiskIOPS != nil {
			diskIOPs = pointer.MakePtr(int(*spec.DiskIOPS))
		}

		return &admin.HardwareSpec20240805{
			InstanceSize:  &spec.InstanceSize,
			NodeCount:     spec.NodeCount,
			EbsVolumeType: pointer.NonZeroOrDefault(spec.EbsVolumeType, "STANDARD"),
			DiskIOPS:      diskIOPs,
			DiskSizeGB:    diskSizeGB,
		}
	}
	dHSpecOrDefault := func(spec *akov2.Specs) *admin.DedicatedHardwareSpec20240805 {
		if spec == nil || *spec.NodeCount == 0 {
			return nil
		}

		var diskIOPs *int
		if spec.DiskIOPS != nil {
			diskIOPs = pointer.MakePtr(int(*spec.DiskIOPS))
		}

		return &admin.DedicatedHardwareSpec20240805{
			InstanceSize:  &spec.InstanceSize,
			NodeCount:     spec.NodeCount,
			EbsVolumeType: pointer.NonZeroOrDefault(spec.EbsVolumeType, "STANDARD"),
			DiskIOPS:      diskIOPs,
			DiskSizeGB:    diskSizeGB,
		}
	}
	autoScalingOrDefault := func(spec *akov2.AdvancedAutoScalingSpec) *admin.AdvancedAutoScalingSettings {
		computeExist := spec != nil && spec.Compute != nil
		diskGBExist := spec != nil && spec.DiskGB != nil

		if !computeExist && !diskGBExist {
			return &admin.AdvancedAutoScalingSettings{
				Compute: &admin.AdvancedComputeAutoScaling{
					Enabled: pointer.MakePtr(false),
				},
				DiskGB: &admin.DiskGBAutoScaling{
					Enabled: pointer.MakePtr(false),
				},
			}
		}

		autoscaling := admin.NewAdvancedAutoScalingSettings()

		if computeExist {
			autoscaling.Compute = &admin.AdvancedComputeAutoScaling{
				Enabled:          spec.Compute.Enabled,
				ScaleDownEnabled: spec.Compute.ScaleDownEnabled,
			}

			if spec.Compute.Enabled != nil && *spec.Compute.Enabled {
				autoscaling.Compute.MaxInstanceSize = &spec.Compute.MaxInstanceSize
			}

			if spec.Compute.ScaleDownEnabled != nil && *spec.Compute.ScaleDownEnabled {
				autoscaling.Compute.MinInstanceSize = &spec.Compute.MinInstanceSize
			}
		}

		if diskGBExist {
			autoscaling.DiskGB = &admin.DiskGBAutoScaling{
				Enabled: spec.DiskGB.Enabled,
			}
		}

		return autoscaling
	}

	specs := make([]admin.ReplicationSpec20240805, 0, len(replicationSpecs))
	for _, spec := range replicationSpecs {
		regionConfigs := make([]admin.CloudRegionConfig20240805, 0, len(spec.RegionConfigs))
		for _, regionConfig := range spec.RegionConfigs {
			regionConfigs = append(
				regionConfigs,
				admin.CloudRegionConfig20240805{
					ProviderName:        pointer.MakePtrOrNil(regionConfig.ProviderName),
					BackingProviderName: pointer.MakePtrOrNil(regionConfig.BackingProviderName),
					RegionName:          pointer.MakePtrOrNil(regionConfig.RegionName),
					Priority:            regionConfig.Priority,
					ElectableSpecs:      hSpecOrDefault(regionConfig.ElectableSpecs),
					ReadOnlySpecs:       dHSpecOrDefault(regionConfig.ReadOnlySpecs),
					AnalyticsSpecs:      dHSpecOrDefault(regionConfig.AnalyticsSpecs),
					AutoScaling:         autoScalingOrDefault(regionConfig.AutoScaling),
				},
			)
		}

		specs = append(
			specs,
			admin.ReplicationSpec20240805{
				ZoneName:      pointer.MakePtrOrNil(spec.ZoneName),
				RegionConfigs: &regionConfigs,
			},
		)
	}
	if clusterType == string(akov2.TypeSharded) {
		for i := 1; i < replicationSpecs[0].NumShards; i++ {
			specs = append(specs, specs[0])
		}
	}

	return &specs
}

func processArgsToAtlas(config *akov2.ProcessArgs) (*admin.ClusterDescriptionProcessArgs20240805, error) {
	var oplogMinRetentionHours *float64
	if config.OplogMinRetentionHours != "" {
		parsed, err := strconv.ParseFloat(config.OplogMinRetentionHours, 64)
		if err != nil {
			return nil, err
		}

		oplogMinRetentionHours = &parsed
	}

	return &admin.ClusterDescriptionProcessArgs20240805{
		DefaultWriteConcern:              pointer.MakePtrOrNil(config.DefaultWriteConcern),
		MinimumEnabledTlsProtocol:        pointer.MakePtrOrNil(config.MinimumEnabledTLSProtocol),
		JavascriptEnabled:                config.JavascriptEnabled,
		NoTableScan:                      config.NoTableScan,
		OplogSizeMB:                      pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.OplogSizeMB, 0))),
		SampleSizeBIConnector:            pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.SampleSizeBIConnector, 0))),
		SampleRefreshIntervalBIConnector: pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.SampleRefreshIntervalBIConnector, 0))),
		OplogMinRetentionHours:           oplogMinRetentionHours,
	}, nil
}

func clustersToConnections(clusters []admin.ClusterDescription20240805) []Connection {
	conns := []Connection{}
	for _, c := range clusters {
		conns = append(conns, Connection{
			Name:             c.GetName(),
			ConnURL:          c.ConnectionStrings.GetStandard(),
			SrvConnURL:       c.ConnectionStrings.GetStandardSrv(),
			PrivateURL:       c.ConnectionStrings.GetPrivate(),
			SrvPrivateURL:    c.ConnectionStrings.GetPrivateSrv(),
			Serverless:       false,
			PrivateEndpoints: fillClusterPrivateEndpoints(c.ConnectionStrings.GetPrivateEndpoint()),
		})
	}
	return conns
}

func fillClusterPrivateEndpoints(cpeList []admin.ClusterDescriptionConnectionStringsPrivateEndpoint) []PrivateEndpoint {
	pes := []PrivateEndpoint{}
	for _, cpe := range cpeList {
		pes = append(pes, PrivateEndpoint{
			URL:       cpe.GetConnectionString(),
			ServerURL: cpe.GetSrvConnectionString(),
			ShardURL:  cpe.GetSrvShardOptimizedConnectionString(),
		})
	}
	return pes
}

func flexToConnections(flex []admin.FlexClusterDescription20241113) []Connection {
	conns := []Connection{}
	for _, f := range flex {
		conns = append(conns, Connection{
			Name:       f.GetName(),
			ConnURL:    f.ConnectionStrings.GetStandard(),
			SrvConnURL: f.ConnectionStrings.GetStandardSrv(),
			Serverless: false,
		})
	}
	return conns
}

func connectionSet(conns ...[]Connection) []Connection {
	return set(func(conn Connection) string { return conn.Name }, conns...)
}

func set[T any](nameFn func(T) string, lists ...[]T) []T {
	hash := map[string]struct{}{}
	result := []T{}
	for _, list := range lists {
		for _, item := range list {
			name := nameFn(item)
			if _, found := hash[name]; !found {
				hash[name] = struct{}{}
				result = append(result, item)
			}
		}
	}
	return result
}

func customZonesFromAtlas(gs *admin.GeoSharding20240805) *[]akov2.CustomZoneMapping {
	if gs == nil {
		return nil
	}
	out := make([]akov2.CustomZoneMapping, 0, len(gs.GetCustomZoneMapping()))
	for k, v := range gs.GetCustomZoneMapping() {
		out = append(out, akov2.CustomZoneMapping{Location: k, Zone: v})
	}
	return &out
}

func managedNamespacesFromAtlas(gs *admin.GeoSharding20240805) []akov2.ManagedNamespace {
	if gs == nil {
		return nil
	}
	out := make([]akov2.ManagedNamespace, len(gs.GetManagedNamespaces()))
	for i, ns := range gs.GetManagedNamespaces() {
		out[i] = akov2.ManagedNamespace{
			Db:                     ns.Db,
			Collection:             ns.Collection,
			CustomShardKey:         ns.CustomShardKey,
			IsCustomShardKeyHashed: ns.IsCustomShardKeyHashed,
			IsShardKeyUnique:       ns.IsShardKeyUnique,
			NumInitialChunks:       int(ns.GetNumInitialChunks()),
			PresplitHashedZones:    ns.PresplitHashedZones,
		}
	}
	return out
}

func customZonesToAtlas(in *[]akov2.CustomZoneMapping) *admin.CustomZoneMappings {
	if in == nil {
		return nil
	}

	out := make([]admin.ZoneMapping, len(*in))
	for i, zone := range *in {
		out[i] = *admin.NewZoneMapping(zone.Location, zone.Zone)
	}

	return &admin.CustomZoneMappings{
		CustomZoneMappings: &out,
	}
}

func managedNamespaceToAtlas(in *akov2.ManagedNamespace) *admin.ManagedNamespaces {
	if in == nil {
		return nil
	}
	return &admin.ManagedNamespaces{
		Db:                     in.Db,
		Collection:             in.Collection,
		CustomShardKey:         in.CustomShardKey,
		IsCustomShardKeyHashed: in.IsCustomShardKeyHashed,
		IsShardKeyUnique:       in.IsShardKeyUnique,
		NumInitialChunks:       pointer.MakePtr(int64(in.NumInitialChunks)),
		PresplitHashedZones:    in.PresplitHashedZones,
	}
}

func flexFromAtlas(instance *admin.FlexClusterDescription20241113) *Flex {
	connectionStrings := instance.GetConnectionStrings()
	providerSettings := instance.GetProviderSettings()

	// Copying this because the existing tags.FromAtlas uses a different SDK version
	t := instance.GetTags()
	tags := make([]*akov2.TagSpec, 0, len(t))
	for _, rTag := range t {
		tags = append(
			tags,
			&akov2.TagSpec{
				Key:   rTag.GetKey(),
				Value: rTag.GetValue(),
			},
		)
	}

	f := &Flex{
		ProjectID: instance.GetGroupId(),
		FlexSpec: &akov2.FlexSpec{
			Name:                         instance.GetName(),
			Tags:                         tags,
			TerminationProtectionEnabled: instance.GetTerminationProtectionEnabled(),
			ProviderSettings: &akov2.FlexProviderSettings{
				BackingProviderName: providerSettings.GetBackingProviderName(),
				RegionName:          providerSettings.GetRegionName(),
			},
		},
		State:          instance.GetStateName(),
		MongoDBVersion: instance.GetMongoDBVersion(),
		Connection: &status.ConnectionStrings{
			Standard:    connectionStrings.GetStandard(),
			StandardSrv: connectionStrings.GetStandardSrv(),
		},
	}
	normalizeFlexDeployment(f)

	return f
}

func flexCreateToAtlas(flex *Flex) *admin.FlexClusterDescriptionCreate20241113 {
	return &admin.FlexClusterDescriptionCreate20241113{
		Name:                         flex.Name,
		Tags:                         tag.FlexToAtlas(flex.Tags),
		TerminationProtectionEnabled: &flex.TerminationProtectionEnabled,
		ProviderSettings: admin.FlexProviderSettingsCreate20241113{
			BackingProviderName: flex.ProviderSettings.BackingProviderName,
			RegionName:          flex.ProviderSettings.RegionName,
		},
	}
}

func flexUpdateToAtlas(flex *Flex) *admin.FlexClusterDescriptionUpdate20241113 {
	return &admin.FlexClusterDescriptionUpdate20241113{
		Tags:                         tag.FlexToAtlas(flex.Tags),
		TerminationProtectionEnabled: &flex.TerminationProtectionEnabled,
	}
}

func flexUpgradeToAtlas(cluster *Cluster) *admin.AtlasTenantClusterUpgradeRequest20240805 {
	spec := cluster.GetCustomResource().Spec.DeploymentSpec
	return &admin.AtlasTenantClusterUpgradeRequest20240805{
		Name:                         spec.Name,
		ClusterType:                  pointer.MakePtrOrNil(spec.ClusterType),
		MongoDBMajorVersion:          pointer.MakePtrOrNil(spec.MongoDBMajorVersion),
		VersionReleaseSystem:         pointer.MakePtrOrNil(spec.VersionReleaseSystem),
		BackupEnabled:                spec.BackupEnabled,
		BiConnector:                  biConnectToAtlas(spec.BiConnector),
		EncryptionAtRestProvider:     pointer.MakePtrOrNil(spec.EncryptionAtRestProvider),
		Labels:                       labelsToAtlas(spec.Labels),
		Paused:                       spec.Paused,
		PitEnabled:                   spec.PitEnabled,
		ReplicationSpecs:             replicationSpecToAtlas(spec.ReplicationSpecs, spec.ClusterType, spec.DiskSizeGB),
		RootCertType:                 pointer.MakePtrOrNil(spec.RootCertType),
		Tags:                         tag.ToAtlas(spec.Tags),
		TerminationProtectionEnabled: pointer.MakePtrOrNil(spec.TerminationProtectionEnabled),
		ConfigServerManagementMode:   pointer.MakePtrOrNil(spec.ConfigServerManagementMode),
	}
}
