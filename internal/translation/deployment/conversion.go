package deployment

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/tag"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

type Deployment interface {
	GetName() string
	GetProjectID() string
	GetCustomResource() *akov2.AtlasDeployment
	GetState() string
	GetMongoDBVersion() string
	GetConnection() *status.ConnectionStrings
	GetReplicaSet() []status.ReplicaSet
}

type Cluster struct {
	*akov2.AdvancedDeploymentSpec
	ProjectID      string
	State          string
	MongoDBVersion string
	Connection     *status.ConnectionStrings
	ProcessArgs    *akov2.ProcessArgs
	ReplicaSet     []status.ReplicaSet

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

func (c *Cluster) IsTenant() bool {
	return c.isTenant
}

type Serverless struct {
	*akov2.ServerlessSpec
	ProjectID      string
	State          string
	MongoDBVersion string
	Connection     *status.ConnectionStrings

	customResource *akov2.AtlasDeployment
}

func (s *Serverless) GetName() string {
	return s.Name
}

func (s *Serverless) GetProjectID() string {
	return s.ProjectID
}

func (s *Serverless) GetState() string {
	return s.State
}

func (s *Serverless) GetMongoDBVersion() string {
	return s.MongoDBVersion
}

func (s *Serverless) GetConnection() *status.ConnectionStrings {
	return s.Connection
}

func (s *Serverless) GetReplicaSet() []status.ReplicaSet {
	return nil
}

func (s *Serverless) GetCustomResource() *akov2.AtlasDeployment {
	return s.customResource
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
		serverless := &Serverless{
			customResource: atlasDeployment,
			ProjectID:      projectID,
			ServerlessSpec: atlasDeployment.Spec.ServerlessSpec.DeepCopy(),
		}
		normalizeServerlessDeployment(serverless)

		return serverless
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

func normalizeServerlessDeployment(serverless *Serverless) {
	serverless.ServerlessSpec.PrivateEndpoints = nil
	if serverless.ServerlessSpec.Tags == nil {
		serverless.ServerlessSpec.Tags = []*akov2.TagSpec{}
	}

	cmp.NormalizeSlice(serverless.Tags, func(a, b *akov2.TagSpec) int {
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
		if cluster.MongoDBMajorVersion == "" {
			cluster.MongoDBMajorVersion = "7.0"
		}

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
		if replicationSpec.NumShards == 0 {
			replicationSpec.NumShards = 1
		}
		if replicationSpec.ZoneName == "" {
			replicationSpec.ZoneName = fmt.Sprintf("Zone %d", ix+1)
		}
		cmp.NormalizeSlice(replicationSpec.RegionConfigs, func(a, b *akov2.AdvancedRegionConfig) int {
			aPriority := "0"
			if a.Priority != nil {
				aPriority = strconv.Itoa(*a.Priority)
			}
			bPriority := "0"
			if b.Priority != nil {
				bPriority = strconv.Itoa(*b.Priority)
			}
			return strings.Compare(a.ProviderName+a.RegionName+aPriority, b.ProviderName+b.RegionName+bPriority)
		})
		for _, regionConfig := range replicationSpec.RegionConfigs {
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
			if computeUnsetOrDisabled && diskUnsetOrDisabled {
				regionConfig.AutoScaling = nil
			}
		}
	}
	cmp.NormalizeSlice(cluster.ReplicationSpecs, func(a, b *akov2.AdvancedReplicationSpec) int {
		return strings.Compare(a.ZoneName, b.ZoneName)
	})
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
}

func getAutoscalingOverride(replications []*akov2.AdvancedReplicationSpec) (bool, bool, string) {
	var instanceSize string
	var isTenant bool
	for _, replica := range replications {
		for _, region := range replica.RegionConfigs {
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

func clusterFromAtlas(clusterDesc *admin.AdvancedClusterDescription) *Cluster {
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
			DiskSizeGB:                   diskSizeFromAtlas(clusterDesc.GetDiskSizeGB()),
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
		},
	}
	normalizeClusterDeployment(cluster)

	return cluster
}

func clusterCreateToAtlas(cluster *Cluster) *admin.AdvancedClusterDescription {
	return &admin.AdvancedClusterDescription{
		Name:                         pointer.MakePtrOrNil(cluster.Name),
		ClusterType:                  pointer.MakePtrOrNil(cluster.ClusterType),
		MongoDBMajorVersion:          pointer.MakePtrOrNil(cluster.MongoDBMajorVersion),
		VersionReleaseSystem:         pointer.MakePtrOrNil(cluster.VersionReleaseSystem),
		DiskSizeGB:                   diskSizeToAtlas(cluster.DiskSizeGB),
		BackupEnabled:                cluster.BackupEnabled,
		BiConnector:                  biConnectToAtlas(cluster.BiConnector),
		EncryptionAtRestProvider:     pointer.MakePtrOrNil(cluster.EncryptionAtRestProvider),
		Labels:                       labelsToAtlas(cluster.Labels),
		Paused:                       cluster.Paused,
		PitEnabled:                   cluster.PitEnabled,
		ReplicationSpecs:             replicationSpecToAtlas(cluster.ReplicationSpecs),
		RootCertType:                 pointer.MakePtrOrNil(cluster.RootCertType),
		Tags:                         tag.ToAtlas(cluster.Tags),
		TerminationProtectionEnabled: pointer.MakePtrOrNil(cluster.TerminationProtectionEnabled),
	}
}

func clusterUpdateToAtlas(cluster *Cluster) *admin.AdvancedClusterDescription {
	return &admin.AdvancedClusterDescription{
		ClusterType:                  pointer.MakePtrOrNil(cluster.ClusterType),
		MongoDBMajorVersion:          pointer.MakePtrOrNil(cluster.MongoDBMajorVersion),
		VersionReleaseSystem:         pointer.MakePtrOrNil(cluster.VersionReleaseSystem),
		DiskSizeGB:                   diskSizeToAtlas(cluster.DiskSizeGB),
		BackupEnabled:                cluster.BackupEnabled,
		BiConnector:                  biConnectToAtlas(cluster.BiConnector),
		EncryptionAtRestProvider:     pointer.MakePtrOrNil(cluster.EncryptionAtRestProvider),
		Labels:                       labelsToAtlas(cluster.Labels),
		Paused:                       cluster.Paused,
		PitEnabled:                   cluster.PitEnabled,
		ReplicationSpecs:             replicationSpecToAtlas(cluster.ReplicationSpecs),
		RootCertType:                 pointer.MakePtrOrNil(cluster.RootCertType),
		Tags:                         tag.ToAtlas(cluster.Tags),
		TerminationProtectionEnabled: pointer.MakePtrOrNil(cluster.TerminationProtectionEnabled),
	}
}

func replicaSetFromAtlas(replicationSpecs []admin.ReplicationSpec) []status.ReplicaSet {
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

func diskSizeFromAtlas(value float64) *int {
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

func replicationSpecFromAtlas(replicationSpecs []admin.ReplicationSpec) []*akov2.AdvancedReplicationSpec {
	hSpecOrDefault := func(spec admin.HardwareSpec) *akov2.Specs {
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
	dHSpecOrDefault := func(spec admin.DedicatedHardwareSpec) *akov2.Specs {
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
					ElectableSpecs:      hSpecOrDefault(regionConfig.GetElectableSpecs()),
					ReadOnlySpecs:       dHSpecOrDefault(regionConfig.GetReadOnlySpecs()),
					AnalyticsSpecs:      dHSpecOrDefault(regionConfig.GetAnalyticsSpecs()),
					AutoScaling:         autoScalingOrDefault(regionConfig.GetAutoScaling()),
				},
			)
		}

		specs = append(
			specs,
			&akov2.AdvancedReplicationSpec{
				NumShards:     spec.GetNumShards(),
				ZoneName:      spec.GetZoneName(),
				RegionConfigs: regionConfigs,
			},
		)
	}

	return specs
}

func processArgsFromAtlas(config *admin.ClusterDescriptionProcessArgs) *akov2.ProcessArgs {
	oplogMinRetentionHours := ""
	if config.GetOplogMinRetentionHours() > 0 {
		oplogMinRetentionHours = strconv.FormatFloat(config.GetOplogMinRetentionHours(), 'f', -1, 64)
	}

	args := akov2.ProcessArgs{
		DefaultReadConcern:               config.GetDefaultReadConcern(),
		DefaultWriteConcern:              config.GetDefaultWriteConcern(),
		MinimumEnabledTLSProtocol:        config.GetMinimumEnabledTlsProtocol(),
		FailIndexKeyTooLong:              config.FailIndexKeyTooLong,
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

func diskSizeToAtlas(value *int) *float64 {
	if value != nil {
		return pointer.MakePtr(float64(*value))
	}

	return nil
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

func replicationSpecToAtlas(replicationSpecs []*akov2.AdvancedReplicationSpec) *[]admin.ReplicationSpec {
	if len(replicationSpecs) == 0 {
		return nil
	}

	hSpecOrDefault := func(spec *akov2.Specs) *admin.HardwareSpec {
		if spec == nil {
			return nil
		}

		var diskIOPs *int
		if spec.DiskIOPS != nil {
			diskIOPs = pointer.MakePtr(int(*spec.DiskIOPS))
		}

		return &admin.HardwareSpec{
			InstanceSize:  &spec.InstanceSize,
			NodeCount:     spec.NodeCount,
			EbsVolumeType: pointer.NonZeroOrDefault(spec.EbsVolumeType, "STANDARD"),
			DiskIOPS:      diskIOPs,
		}
	}
	dHSpecOrDefault := func(spec *akov2.Specs) *admin.DedicatedHardwareSpec {
		if spec == nil || *spec.NodeCount == 0 {
			return nil
		}

		var diskIOPs *int
		if spec.DiskIOPS != nil {
			diskIOPs = pointer.MakePtr(int(*spec.DiskIOPS))
		}

		return &admin.DedicatedHardwareSpec{
			InstanceSize:  &spec.InstanceSize,
			NodeCount:     spec.NodeCount,
			EbsVolumeType: pointer.NonZeroOrDefault(spec.EbsVolumeType, "STANDARD"),
			DiskIOPS:      diskIOPs,
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

	specs := make([]admin.ReplicationSpec, 0, len(replicationSpecs))
	for _, spec := range replicationSpecs {
		regionConfigs := make([]admin.CloudRegionConfig, 0, len(spec.RegionConfigs))
		for _, regionConfig := range spec.RegionConfigs {
			regionConfigs = append(
				regionConfigs,
				admin.CloudRegionConfig{
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
			admin.ReplicationSpec{
				NumShards:     pointer.NonZeroOrDefault(spec.NumShards, 1),
				ZoneName:      pointer.MakePtrOrNil(spec.ZoneName),
				RegionConfigs: &regionConfigs,
			},
		)
	}

	return &specs
}

func processArgsToAtlas(config *akov2.ProcessArgs) (*admin.ClusterDescriptionProcessArgs, error) {
	var oplogMinRetentionHours *float64
	if config.OplogMinRetentionHours != "" {
		parsed, err := strconv.ParseFloat(config.OplogMinRetentionHours, 64)
		if err != nil {
			return nil, err
		}

		oplogMinRetentionHours = &parsed
	}

	return &admin.ClusterDescriptionProcessArgs{
		DefaultReadConcern:               pointer.MakePtrOrNil(config.DefaultReadConcern),
		DefaultWriteConcern:              pointer.MakePtrOrNil(config.DefaultWriteConcern),
		MinimumEnabledTlsProtocol:        pointer.MakePtrOrNil(config.MinimumEnabledTLSProtocol),
		FailIndexKeyTooLong:              config.FailIndexKeyTooLong,
		JavascriptEnabled:                config.JavascriptEnabled,
		NoTableScan:                      config.NoTableScan,
		OplogSizeMB:                      pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.OplogSizeMB, 0))),
		SampleSizeBIConnector:            pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.SampleSizeBIConnector, 0))),
		SampleRefreshIntervalBIConnector: pointer.MakePtrOrNil(int(pointer.GetOrDefault(config.SampleRefreshIntervalBIConnector, 0))),
		OplogMinRetentionHours:           oplogMinRetentionHours,
	}, nil
}

func serverlessFromAtlas(instance *admin.ServerlessInstanceDescription) *Serverless {
	providerSettings := instance.GetProviderSettings()
	serverlessBackupOptions := instance.GetServerlessBackupOptions()
	connectionStrings := instance.GetConnectionStrings()

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
				SRVConnectionString: pe.GetSrvConnectionString(),
				Endpoints:           eps,
			})
	}

	s := &Serverless{
		ProjectID: instance.GetGroupId(),
		ServerlessSpec: &akov2.ServerlessSpec{
			Name: instance.GetName(),
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				ProviderName:        provider.ProviderName(providerSettings.GetProviderName()),
				BackingProviderName: providerSettings.GetBackingProviderName(),
				RegionName:          providerSettings.GetRegionName(),
			},
			Tags: tag.FromAtlas(instance.GetTags()),
			BackupOptions: akov2.ServerlessBackupOptions{
				ServerlessContinuousBackupEnabled: serverlessBackupOptions.GetServerlessContinuousBackupEnabled(),
			},
			TerminationProtectionEnabled: instance.GetTerminationProtectionEnabled(),
		},
		State:          instance.GetStateName(),
		MongoDBVersion: instance.GetMongoDBVersion(),
		Connection: &status.ConnectionStrings{
			StandardSrv:     connectionStrings.GetStandardSrv(),
			PrivateEndpoint: pes,
		},
	}
	normalizeServerlessDeployment(s)

	return s
}

func serverlessCreateToAtlas(serverless *Serverless) *admin.ServerlessInstanceDescriptionCreate {
	return &admin.ServerlessInstanceDescriptionCreate{
		Name: serverless.Name,
		ProviderSettings: admin.ServerlessProviderSettings{
			ProviderName:        pointer.MakePtr(string(serverless.ProviderSettings.ProviderName)),
			BackingProviderName: serverless.ProviderSettings.BackingProviderName,
			RegionName:          serverless.ProviderSettings.RegionName,
		},
		ServerlessBackupOptions: &admin.ClusterServerlessBackupOptions{
			ServerlessContinuousBackupEnabled: &serverless.BackupOptions.ServerlessContinuousBackupEnabled,
		},
		Tags:                         tag.ToAtlas(serverless.Tags),
		TerminationProtectionEnabled: &serverless.TerminationProtectionEnabled,
	}
}

func serverlessUpdateToAtlas(serverless *Serverless) *admin.ServerlessInstanceDescriptionUpdate {
	return &admin.ServerlessInstanceDescriptionUpdate{
		ServerlessBackupOptions: &admin.ClusterServerlessBackupOptions{
			ServerlessContinuousBackupEnabled: &serverless.BackupOptions.ServerlessContinuousBackupEnabled,
		},
		Tags:                         tag.ToAtlas(serverless.Tags),
		TerminationProtectionEnabled: &serverless.TerminationProtectionEnabled,
	}
}

func clustersToConnections(clusters []admin.AdvancedClusterDescription) []Connection {
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

func serverlessToConnections(serverless []admin.ServerlessInstanceDescription) []Connection {
	conns := []Connection{}
	for _, s := range serverless {
		conns = append(conns, Connection{
			Name:             s.GetName(),
			ConnURL:          "",
			SrvConnURL:       s.ConnectionStrings.GetStandardSrv(),
			Serverless:       true,
			PrivateEndpoints: fillServerlessPrivateEndpoints(s.ConnectionStrings.GetPrivateEndpoint()),
		})
	}
	return conns
}

func fillServerlessPrivateEndpoints(cpeList []admin.ServerlessConnectionStringsPrivateEndpointList) []PrivateEndpoint {
	pes := []PrivateEndpoint{}
	for _, cpe := range cpeList {
		pes = append(pes, PrivateEndpoint{
			ServerURL: cpe.GetSrvConnectionString(),
		})
	}
	return pes
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
