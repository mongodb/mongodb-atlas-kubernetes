// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ClusterDescription20240805 struct for ClusterDescription20240805
type ClusterDescription20240805 struct {
	// If reconfiguration is necessary to regain a primary due to a regional outage, submit this field alongside your topology reconfiguration to request a new regional outage resistant topology. Forced reconfigurations during an outage of the majority of electable nodes carry a risk of data loss if replicated writes (even majority committed writes) have not been replicated to the new primary node. MongoDB Atlas docs contain more information. To proceed with an operation which carries that risk, set `acceptDataRisksAndForceReplicaSetReconfig` to the current date. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	AcceptDataRisksAndForceReplicaSetReconfig *time.Time `json:"acceptDataRisksAndForceReplicaSetReconfig,omitempty"`
	// Governs adaptive capacity behavior of Azure nodes in single-cloud Azure clusters or multi-cloud clusters that include Azure nodes. Adaptive capacity enables fallback hardware selection when the primary instance family is unavailable. ``ENABLED`` means the cluster explicitly opts in to adaptive capacity. ``DISABLED`` means the cluster explicitly opts out; the cluster receives capacity errors instead of being placed on fallback hardware. ``null`` means the field is unset; Azure clusters use adaptive capacity by default when the feature is enabled at the group level. Setting this field for single-cloud AWS or GCP clusters is a no-op.
	AdaptiveCapacity      *string                               `json:"adaptiveCapacity,omitempty"`
	AdvancedConfiguration *ApiAtlasClusterAdvancedConfiguration `json:"advancedConfiguration,omitempty"`
	// Flag that indicates whether the cluster can perform backups. If set to `true`, the cluster can perform backups. You must set this value to `true` for NVMe clusters. Backup uses Cloud Backups for dedicated clusters and [Shared Cluster Backups](https://docs.atlas.mongodb.com/backup/shared-tier/overview/) for tenant clusters. If set to `false`, the cluster doesn't use backups.
	BackupEnabled *bool        `json:"backupEnabled,omitempty"`
	BiConnector   *BiConnector `json:"biConnector,omitempty"`
	// Configuration of nodes that comprise the cluster.
	ClusterType *string `json:"clusterType,omitempty"`
	// Config Server Management Mode for creating or updating a sharded cluster. When configured as `ATLAS_MANAGED`, Atlas may automatically switch the cluster's config server type for optimal performance and savings. When configured as `FIXED_TO_DEDICATED`, the cluster will always use a dedicated config server.
	ConfigServerManagementMode *string `json:"configServerManagementMode,omitempty"`
	// Describes a sharded cluster's config server type.
	// Read only field.
	ConfigServerType  *string                   `json:"configServerType,omitempty"`
	ConnectionStrings *ClusterConnectionStrings `json:"connectionStrings,omitempty"`
	// Date and time when MongoDB Cloud created this cluster. This parameter expresses its value in ISO 8601 format in UTC.
	// Read only field.
	CreateDate *time.Time `json:"createDate,omitempty"`
	// Disk warming mode selection.
	DiskWarmingMode *string `json:"diskWarmingMode,omitempty"`
	// List of settings that represent the actual cluster state. This is read-only and always returned in the response. It reflects the current cluster configuration, which may differ from `replicationSpecs` due to system-managed changes.
	// Read only field.
	EffectiveReplicationSpecs *[]ReplicationSpec20240805 `json:"effectiveReplicationSpecs,omitempty"`
	// Cloud service provider that manages your customer keys to provide an additional layer of encryption at rest for the cluster. To enable customer key management for encryption at rest, the cluster `replicationSpecs[n].regionConfigs[m].{type}Specs.instanceSize` setting must be `M10` or higher and `\"backupEnabled\" : false` or omitted entirely.
	EncryptionAtRestProvider *string `json:"encryptionAtRestProvider,omitempty"`
	// Feature compatibility version of the cluster. This will always appear regardless of whether FCV is pinned.
	// Read only field.
	FeatureCompatibilityVersion *string `json:"featureCompatibilityVersion,omitempty"`
	// Feature compatibility version expiration date. Will only appear if FCV is pinned. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	FeatureCompatibilityVersionExpirationDate *time.Time `json:"featureCompatibilityVersionExpirationDate,omitempty"`
	// Set this field to configure the Sharding Management Mode when creating a new Global Cluster.  When set to false, the management mode is set to Atlas-Managed Sharding. This mode fully manages the sharding of your Global Cluster and is built to provide a seamless deployment experience.  When set to true, the management mode is set to Self-Managed Sharding. This mode leaves the management of shards in your hands and is built to provide an advanced and flexible deployment experience.  This setting cannot be changed once the cluster is deployed.
	GlobalClusterSelfManagedSharding *bool `json:"globalClusterSelfManagedSharding,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the cluster.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Internal classification of the cluster's role. Possible values: `NONE` (regular user cluster), `SYSTEM_CLUSTER` (system cluster for backup), `INTERNAL_SHADOW_CLUSTER` (internal use shadow cluster for testing).
	// Read only field.
	InternalClusterRole *string `json:"internalClusterRole,omitempty"`
	// Collection of key-value pairs between 1 to 255 characters in length that tag and categorize the cluster. The MongoDB Cloud console doesn't display your labels.  Cluster labels are deprecated and will be removed in a future release. We strongly recommend that you use Resource Tags instead.
	// Deprecated
	Labels *[]ComponentLabel `json:"labels,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links                      *[]Link              `json:"links,omitempty"`
	MongoDBEmployeeAccessGrant *EmployeeAccessGrant `json:"mongoDBEmployeeAccessGrant,omitempty"`
	// MongoDB major version of the cluster. Set to the binary major version.   On creation: Choose from the available versions of MongoDB, or leave unspecified for the current recommended default in the MongoDB Cloud platform. The recommended version is a recent Long Term Support version. The default is not guaranteed to be the most recently released version throughout the entire release cycle. For versions available in a specific project, see the linked documentation or use the API endpoint for [project LTS versions endpoint](#tag/Projects/operation/getProjectLtsVersions).   On update: Increase version only by 1 major version at a time. If the cluster is pinned to a MongoDB feature compatibility version exactly one major version below the current MongoDB version, the MongoDB version can be downgraded to the previous major version.
	MongoDBMajorVersion *string `json:"mongoDBMajorVersion,omitempty"`
	// Version of MongoDB that the cluster runs.
	// Read only field.
	MongoDBVersion *string `json:"mongoDBVersion,omitempty"`
	// Human-readable label that identifies the cluster.
	Name *string `json:"name,omitempty"`
	// Flag that indicates whether the cluster is paused.
	Paused *bool `json:"paused,omitempty"`
	// Flag that indicates whether the cluster uses continuous cloud backups.
	PitEnabled *bool `json:"pitEnabled,omitempty"`
	// Enable or disable log redaction.  This setting configures the ``mongod`` or ``mongos`` to redact any document field contents from a message accompanying a given log event before logging. This prevents the program from writing potentially sensitive data stored on the database to the diagnostic log. Metadata such as error or operation codes, line numbers, and source file names are still visible in the logs.  Use ``redactClientLogData`` in conjunction with Encryption at Rest and TLS/SSL (Transport Encryption) to assist compliance with regulatory requirements.  *Note*: changing this setting on a cluster will trigger a rolling restart as soon as the cluster is updated.
	RedactClientLogData *bool `json:"redactClientLogData,omitempty"`
	// Set this field to configure the replica set scaling mode for your cluster.  By default, Atlas scales under `WORKLOAD_TYPE`. This mode allows Atlas to scale your analytics nodes in parallel to your operational nodes.  When configured as `SEQUENTIAL`, Atlas scales all nodes sequentially. This mode is intended for steady-state workloads and applications performing latency-sensitive secondary reads.  When configured as `NODE_TYPE`, Atlas scales your electable nodes in parallel with your read-only and analytics nodes. This mode is intended for large, dynamic workloads requiring frequent and timely cluster tier scaling. This is the fastest scaling strategy, but it might impact latency of workloads when performing extensive secondary reads.
	ReplicaSetScalingStrategy *string `json:"replicaSetScalingStrategy,omitempty"`
	// List of settings that configure your cluster regions. This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.
	ReplicationSpecs *[]ReplicationSpec20240805 `json:"replicationSpecs,omitempty"`
	// Flag that indicates whether the cluster retains backups.
	RetainBackups *bool `json:"retainBackups,omitempty"`
	// Root Certificate Authority that MongoDB Atlas cluster uses. MongoDB Cloud supports Internet Security Research Group.
	RootCertType *string `json:"rootCertType,omitempty"`
	// Human-readable label that indicates any current activity being taken on this cluster by the Atlas control plane. With the exception of CREATING and DELETING states, clusters should always be available and have a Primary node even when in states indicating ongoing activity.   - `IDLE`: Atlas is making no changes to this cluster and all changes requested via the UI or API can be assumed to have been applied.  - `CREATING`: A cluster being provisioned for the very first time returns state CREATING until it is ready for connections. Ensure IP Access List and DB Users are configured before attempting to connect.  - `UPDATING`: A change requested via the UI, API, AutoScaling, or other scheduled activity is taking place.  - `DELETING`: The cluster is in the process of deletion and will soon be deleted.  - `REPAIRING`: One or more nodes in the cluster are being returned to service by the Atlas control plane. Other nodes should continue to provide service as normal.
	// Read only field.
	StateName *string `json:"stateName,omitempty"`
	// List that contains key-value pairs between 1 to 255 characters in length for tagging and categorizing the cluster.
	Tags *[]ResourceTag `json:"tags,omitempty"`
	// Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, MongoDB Cloud won't delete the cluster. If set to `false`, MongoDB Cloud will delete the cluster.
	TerminationProtectionEnabled *bool `json:"terminationProtectionEnabled,omitempty"`
	// Flag that indicates whether AWS time-based snapshot copies will be used instead of slower standard snapshot copies during fast Atlas cross-region initial syncs. This flag is only relevant for clusters containing AWS nodes.
	UseAwsTimeBasedSnapshotCopyForFastInitialSync *bool `json:"useAwsTimeBasedSnapshotCopyForFastInitialSync,omitempty"`
	// Method by which the cluster maintains the MongoDB versions. If value is `CONTINUOUS`, you must not specify `mongoDBMajorVersion`.
	VersionReleaseSystem *string `json:"versionReleaseSystem,omitempty"`
}

// NewClusterDescription20240805 instantiates a new ClusterDescription20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterDescription20240805() *ClusterDescription20240805 {
	this := ClusterDescription20240805{}
	var backupEnabled bool = false
	this.BackupEnabled = &backupEnabled
	var configServerManagementMode string = "ATLAS_MANAGED"
	this.ConfigServerManagementMode = &configServerManagementMode
	var diskWarmingMode string = "FULLY_WARMED"
	this.DiskWarmingMode = &diskWarmingMode
	var replicaSetScalingStrategy string = "WORKLOAD_TYPE"
	this.ReplicaSetScalingStrategy = &replicaSetScalingStrategy
	var retainBackups bool = false
	this.RetainBackups = &retainBackups
	var rootCertType string = "ISRGROOTX1"
	this.RootCertType = &rootCertType
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	var useAwsTimeBasedSnapshotCopyForFastInitialSync bool = false
	this.UseAwsTimeBasedSnapshotCopyForFastInitialSync = &useAwsTimeBasedSnapshotCopyForFastInitialSync
	var versionReleaseSystem string = "LTS"
	this.VersionReleaseSystem = &versionReleaseSystem
	return &this
}

// NewClusterDescription20240805WithDefaults instantiates a new ClusterDescription20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterDescription20240805WithDefaults() *ClusterDescription20240805 {
	this := ClusterDescription20240805{}
	var backupEnabled bool = false
	this.BackupEnabled = &backupEnabled
	var configServerManagementMode string = "ATLAS_MANAGED"
	this.ConfigServerManagementMode = &configServerManagementMode
	var diskWarmingMode string = "FULLY_WARMED"
	this.DiskWarmingMode = &diskWarmingMode
	var replicaSetScalingStrategy string = "WORKLOAD_TYPE"
	this.ReplicaSetScalingStrategy = &replicaSetScalingStrategy
	var retainBackups bool = false
	this.RetainBackups = &retainBackups
	var rootCertType string = "ISRGROOTX1"
	this.RootCertType = &rootCertType
	var terminationProtectionEnabled bool = false
	this.TerminationProtectionEnabled = &terminationProtectionEnabled
	var useAwsTimeBasedSnapshotCopyForFastInitialSync bool = false
	this.UseAwsTimeBasedSnapshotCopyForFastInitialSync = &useAwsTimeBasedSnapshotCopyForFastInitialSync
	var versionReleaseSystem string = "LTS"
	this.VersionReleaseSystem = &versionReleaseSystem
	return &this
}

// GetAcceptDataRisksAndForceReplicaSetReconfig returns the AcceptDataRisksAndForceReplicaSetReconfig field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetAcceptDataRisksAndForceReplicaSetReconfig() time.Time {
	if o == nil || IsNil(o.AcceptDataRisksAndForceReplicaSetReconfig) {
		var ret time.Time
		return ret
	}
	return *o.AcceptDataRisksAndForceReplicaSetReconfig
}

// GetAcceptDataRisksAndForceReplicaSetReconfigOk returns a tuple with the AcceptDataRisksAndForceReplicaSetReconfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetAcceptDataRisksAndForceReplicaSetReconfigOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AcceptDataRisksAndForceReplicaSetReconfig) {
		return nil, false
	}

	return o.AcceptDataRisksAndForceReplicaSetReconfig, true
}

// HasAcceptDataRisksAndForceReplicaSetReconfig returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasAcceptDataRisksAndForceReplicaSetReconfig() bool {
	if o != nil && !IsNil(o.AcceptDataRisksAndForceReplicaSetReconfig) {
		return true
	}

	return false
}

// SetAcceptDataRisksAndForceReplicaSetReconfig gets a reference to the given time.Time and assigns it to the AcceptDataRisksAndForceReplicaSetReconfig field.
func (o *ClusterDescription20240805) SetAcceptDataRisksAndForceReplicaSetReconfig(v time.Time) {
	o.AcceptDataRisksAndForceReplicaSetReconfig = &v
}

// GetAdaptiveCapacity returns the AdaptiveCapacity field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetAdaptiveCapacity() string {
	if o == nil || IsNil(o.AdaptiveCapacity) {
		var ret string
		return ret
	}
	return *o.AdaptiveCapacity
}

// GetAdaptiveCapacityOk returns a tuple with the AdaptiveCapacity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetAdaptiveCapacityOk() (*string, bool) {
	if o == nil || IsNil(o.AdaptiveCapacity) {
		return nil, false
	}

	return o.AdaptiveCapacity, true
}

// HasAdaptiveCapacity returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasAdaptiveCapacity() bool {
	if o != nil && !IsNil(o.AdaptiveCapacity) {
		return true
	}

	return false
}

// SetAdaptiveCapacity gets a reference to the given string and assigns it to the AdaptiveCapacity field.
func (o *ClusterDescription20240805) SetAdaptiveCapacity(v string) {
	o.AdaptiveCapacity = &v
}

// GetAdvancedConfiguration returns the AdvancedConfiguration field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetAdvancedConfiguration() ApiAtlasClusterAdvancedConfiguration {
	if o == nil || IsNil(o.AdvancedConfiguration) {
		var ret ApiAtlasClusterAdvancedConfiguration
		return ret
	}
	return *o.AdvancedConfiguration
}

// GetAdvancedConfigurationOk returns a tuple with the AdvancedConfiguration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetAdvancedConfigurationOk() (*ApiAtlasClusterAdvancedConfiguration, bool) {
	if o == nil || IsNil(o.AdvancedConfiguration) {
		return nil, false
	}

	return o.AdvancedConfiguration, true
}

// HasAdvancedConfiguration returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasAdvancedConfiguration() bool {
	if o != nil && !IsNil(o.AdvancedConfiguration) {
		return true
	}

	return false
}

// SetAdvancedConfiguration gets a reference to the given ApiAtlasClusterAdvancedConfiguration and assigns it to the AdvancedConfiguration field.
func (o *ClusterDescription20240805) SetAdvancedConfiguration(v ApiAtlasClusterAdvancedConfiguration) {
	o.AdvancedConfiguration = &v
}

// GetBackupEnabled returns the BackupEnabled field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetBackupEnabled() bool {
	if o == nil || IsNil(o.BackupEnabled) {
		var ret bool
		return ret
	}
	return *o.BackupEnabled
}

// GetBackupEnabledOk returns a tuple with the BackupEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetBackupEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.BackupEnabled) {
		return nil, false
	}

	return o.BackupEnabled, true
}

// HasBackupEnabled returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasBackupEnabled() bool {
	if o != nil && !IsNil(o.BackupEnabled) {
		return true
	}

	return false
}

// SetBackupEnabled gets a reference to the given bool and assigns it to the BackupEnabled field.
func (o *ClusterDescription20240805) SetBackupEnabled(v bool) {
	o.BackupEnabled = &v
}

// GetBiConnector returns the BiConnector field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetBiConnector() BiConnector {
	if o == nil || IsNil(o.BiConnector) {
		var ret BiConnector
		return ret
	}
	return *o.BiConnector
}

// GetBiConnectorOk returns a tuple with the BiConnector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetBiConnectorOk() (*BiConnector, bool) {
	if o == nil || IsNil(o.BiConnector) {
		return nil, false
	}

	return o.BiConnector, true
}

// HasBiConnector returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasBiConnector() bool {
	if o != nil && !IsNil(o.BiConnector) {
		return true
	}

	return false
}

// SetBiConnector gets a reference to the given BiConnector and assigns it to the BiConnector field.
func (o *ClusterDescription20240805) SetBiConnector(v BiConnector) {
	o.BiConnector = &v
}

// GetClusterType returns the ClusterType field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetClusterType() string {
	if o == nil || IsNil(o.ClusterType) {
		var ret string
		return ret
	}
	return *o.ClusterType
}

// GetClusterTypeOk returns a tuple with the ClusterType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetClusterTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterType) {
		return nil, false
	}

	return o.ClusterType, true
}

// HasClusterType returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasClusterType() bool {
	if o != nil && !IsNil(o.ClusterType) {
		return true
	}

	return false
}

// SetClusterType gets a reference to the given string and assigns it to the ClusterType field.
func (o *ClusterDescription20240805) SetClusterType(v string) {
	o.ClusterType = &v
}

// GetConfigServerManagementMode returns the ConfigServerManagementMode field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetConfigServerManagementMode() string {
	if o == nil || IsNil(o.ConfigServerManagementMode) {
		var ret string
		return ret
	}
	return *o.ConfigServerManagementMode
}

// GetConfigServerManagementModeOk returns a tuple with the ConfigServerManagementMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetConfigServerManagementModeOk() (*string, bool) {
	if o == nil || IsNil(o.ConfigServerManagementMode) {
		return nil, false
	}

	return o.ConfigServerManagementMode, true
}

// HasConfigServerManagementMode returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasConfigServerManagementMode() bool {
	if o != nil && !IsNil(o.ConfigServerManagementMode) {
		return true
	}

	return false
}

// SetConfigServerManagementMode gets a reference to the given string and assigns it to the ConfigServerManagementMode field.
func (o *ClusterDescription20240805) SetConfigServerManagementMode(v string) {
	o.ConfigServerManagementMode = &v
}

// GetConfigServerType returns the ConfigServerType field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetConfigServerType() string {
	if o == nil || IsNil(o.ConfigServerType) {
		var ret string
		return ret
	}
	return *o.ConfigServerType
}

// GetConfigServerTypeOk returns a tuple with the ConfigServerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetConfigServerTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ConfigServerType) {
		return nil, false
	}

	return o.ConfigServerType, true
}

// HasConfigServerType returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasConfigServerType() bool {
	if o != nil && !IsNil(o.ConfigServerType) {
		return true
	}

	return false
}

// SetConfigServerType gets a reference to the given string and assigns it to the ConfigServerType field.
func (o *ClusterDescription20240805) SetConfigServerType(v string) {
	o.ConfigServerType = &v
}

// GetConnectionStrings returns the ConnectionStrings field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetConnectionStrings() ClusterConnectionStrings {
	if o == nil || IsNil(o.ConnectionStrings) {
		var ret ClusterConnectionStrings
		return ret
	}
	return *o.ConnectionStrings
}

// GetConnectionStringsOk returns a tuple with the ConnectionStrings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetConnectionStringsOk() (*ClusterConnectionStrings, bool) {
	if o == nil || IsNil(o.ConnectionStrings) {
		return nil, false
	}

	return o.ConnectionStrings, true
}

// HasConnectionStrings returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasConnectionStrings() bool {
	if o != nil && !IsNil(o.ConnectionStrings) {
		return true
	}

	return false
}

// SetConnectionStrings gets a reference to the given ClusterConnectionStrings and assigns it to the ConnectionStrings field.
func (o *ClusterDescription20240805) SetConnectionStrings(v ClusterConnectionStrings) {
	o.ConnectionStrings = &v
}

// GetCreateDate returns the CreateDate field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetCreateDate() time.Time {
	if o == nil || IsNil(o.CreateDate) {
		var ret time.Time
		return ret
	}
	return *o.CreateDate
}

// GetCreateDateOk returns a tuple with the CreateDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetCreateDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateDate) {
		return nil, false
	}

	return o.CreateDate, true
}

// HasCreateDate returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasCreateDate() bool {
	if o != nil && !IsNil(o.CreateDate) {
		return true
	}

	return false
}

// SetCreateDate gets a reference to the given time.Time and assigns it to the CreateDate field.
func (o *ClusterDescription20240805) SetCreateDate(v time.Time) {
	o.CreateDate = &v
}

// GetDiskWarmingMode returns the DiskWarmingMode field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetDiskWarmingMode() string {
	if o == nil || IsNil(o.DiskWarmingMode) {
		var ret string
		return ret
	}
	return *o.DiskWarmingMode
}

// GetDiskWarmingModeOk returns a tuple with the DiskWarmingMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetDiskWarmingModeOk() (*string, bool) {
	if o == nil || IsNil(o.DiskWarmingMode) {
		return nil, false
	}

	return o.DiskWarmingMode, true
}

// HasDiskWarmingMode returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasDiskWarmingMode() bool {
	if o != nil && !IsNil(o.DiskWarmingMode) {
		return true
	}

	return false
}

// SetDiskWarmingMode gets a reference to the given string and assigns it to the DiskWarmingMode field.
func (o *ClusterDescription20240805) SetDiskWarmingMode(v string) {
	o.DiskWarmingMode = &v
}

// GetEffectiveReplicationSpecs returns the EffectiveReplicationSpecs field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetEffectiveReplicationSpecs() []ReplicationSpec20240805 {
	if o == nil || IsNil(o.EffectiveReplicationSpecs) {
		var ret []ReplicationSpec20240805
		return ret
	}
	return *o.EffectiveReplicationSpecs
}

// GetEffectiveReplicationSpecsOk returns a tuple with the EffectiveReplicationSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetEffectiveReplicationSpecsOk() (*[]ReplicationSpec20240805, bool) {
	if o == nil || IsNil(o.EffectiveReplicationSpecs) {
		return nil, false
	}

	return o.EffectiveReplicationSpecs, true
}

// HasEffectiveReplicationSpecs returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasEffectiveReplicationSpecs() bool {
	if o != nil && !IsNil(o.EffectiveReplicationSpecs) {
		return true
	}

	return false
}

// SetEffectiveReplicationSpecs gets a reference to the given []ReplicationSpec20240805 and assigns it to the EffectiveReplicationSpecs field.
func (o *ClusterDescription20240805) SetEffectiveReplicationSpecs(v []ReplicationSpec20240805) {
	o.EffectiveReplicationSpecs = &v
}

// GetEncryptionAtRestProvider returns the EncryptionAtRestProvider field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetEncryptionAtRestProvider() string {
	if o == nil || IsNil(o.EncryptionAtRestProvider) {
		var ret string
		return ret
	}
	return *o.EncryptionAtRestProvider
}

// GetEncryptionAtRestProviderOk returns a tuple with the EncryptionAtRestProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetEncryptionAtRestProviderOk() (*string, bool) {
	if o == nil || IsNil(o.EncryptionAtRestProvider) {
		return nil, false
	}

	return o.EncryptionAtRestProvider, true
}

// HasEncryptionAtRestProvider returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasEncryptionAtRestProvider() bool {
	if o != nil && !IsNil(o.EncryptionAtRestProvider) {
		return true
	}

	return false
}

// SetEncryptionAtRestProvider gets a reference to the given string and assigns it to the EncryptionAtRestProvider field.
func (o *ClusterDescription20240805) SetEncryptionAtRestProvider(v string) {
	o.EncryptionAtRestProvider = &v
}

// GetFeatureCompatibilityVersion returns the FeatureCompatibilityVersion field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetFeatureCompatibilityVersion() string {
	if o == nil || IsNil(o.FeatureCompatibilityVersion) {
		var ret string
		return ret
	}
	return *o.FeatureCompatibilityVersion
}

// GetFeatureCompatibilityVersionOk returns a tuple with the FeatureCompatibilityVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetFeatureCompatibilityVersionOk() (*string, bool) {
	if o == nil || IsNil(o.FeatureCompatibilityVersion) {
		return nil, false
	}

	return o.FeatureCompatibilityVersion, true
}

// HasFeatureCompatibilityVersion returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasFeatureCompatibilityVersion() bool {
	if o != nil && !IsNil(o.FeatureCompatibilityVersion) {
		return true
	}

	return false
}

// SetFeatureCompatibilityVersion gets a reference to the given string and assigns it to the FeatureCompatibilityVersion field.
func (o *ClusterDescription20240805) SetFeatureCompatibilityVersion(v string) {
	o.FeatureCompatibilityVersion = &v
}

// GetFeatureCompatibilityVersionExpirationDate returns the FeatureCompatibilityVersionExpirationDate field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetFeatureCompatibilityVersionExpirationDate() time.Time {
	if o == nil || IsNil(o.FeatureCompatibilityVersionExpirationDate) {
		var ret time.Time
		return ret
	}
	return *o.FeatureCompatibilityVersionExpirationDate
}

// GetFeatureCompatibilityVersionExpirationDateOk returns a tuple with the FeatureCompatibilityVersionExpirationDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetFeatureCompatibilityVersionExpirationDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FeatureCompatibilityVersionExpirationDate) {
		return nil, false
	}

	return o.FeatureCompatibilityVersionExpirationDate, true
}

// HasFeatureCompatibilityVersionExpirationDate returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasFeatureCompatibilityVersionExpirationDate() bool {
	if o != nil && !IsNil(o.FeatureCompatibilityVersionExpirationDate) {
		return true
	}

	return false
}

// SetFeatureCompatibilityVersionExpirationDate gets a reference to the given time.Time and assigns it to the FeatureCompatibilityVersionExpirationDate field.
func (o *ClusterDescription20240805) SetFeatureCompatibilityVersionExpirationDate(v time.Time) {
	o.FeatureCompatibilityVersionExpirationDate = &v
}

// GetGlobalClusterSelfManagedSharding returns the GlobalClusterSelfManagedSharding field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetGlobalClusterSelfManagedSharding() bool {
	if o == nil || IsNil(o.GlobalClusterSelfManagedSharding) {
		var ret bool
		return ret
	}
	return *o.GlobalClusterSelfManagedSharding
}

// GetGlobalClusterSelfManagedShardingOk returns a tuple with the GlobalClusterSelfManagedSharding field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetGlobalClusterSelfManagedShardingOk() (*bool, bool) {
	if o == nil || IsNil(o.GlobalClusterSelfManagedSharding) {
		return nil, false
	}

	return o.GlobalClusterSelfManagedSharding, true
}

// HasGlobalClusterSelfManagedSharding returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasGlobalClusterSelfManagedSharding() bool {
	if o != nil && !IsNil(o.GlobalClusterSelfManagedSharding) {
		return true
	}

	return false
}

// SetGlobalClusterSelfManagedSharding gets a reference to the given bool and assigns it to the GlobalClusterSelfManagedSharding field.
func (o *ClusterDescription20240805) SetGlobalClusterSelfManagedSharding(v bool) {
	o.GlobalClusterSelfManagedSharding = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ClusterDescription20240805) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ClusterDescription20240805) SetId(v string) {
	o.Id = &v
}

// GetInternalClusterRole returns the InternalClusterRole field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetInternalClusterRole() string {
	if o == nil || IsNil(o.InternalClusterRole) {
		var ret string
		return ret
	}
	return *o.InternalClusterRole
}

// GetInternalClusterRoleOk returns a tuple with the InternalClusterRole field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetInternalClusterRoleOk() (*string, bool) {
	if o == nil || IsNil(o.InternalClusterRole) {
		return nil, false
	}

	return o.InternalClusterRole, true
}

// HasInternalClusterRole returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasInternalClusterRole() bool {
	if o != nil && !IsNil(o.InternalClusterRole) {
		return true
	}

	return false
}

// SetInternalClusterRole gets a reference to the given string and assigns it to the InternalClusterRole field.
func (o *ClusterDescription20240805) SetInternalClusterRole(v string) {
	o.InternalClusterRole = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise
// Deprecated
func (o *ClusterDescription20240805) GetLabels() []ComponentLabel {
	if o == nil || IsNil(o.Labels) {
		var ret []ComponentLabel
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *ClusterDescription20240805) GetLabelsOk() (*[]ComponentLabel, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}

	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given []ComponentLabel and assigns it to the Labels field.
// Deprecated
func (o *ClusterDescription20240805) SetLabels(v []ComponentLabel) {
	o.Labels = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ClusterDescription20240805) SetLinks(v []Link) {
	o.Links = &v
}

// GetMongoDBEmployeeAccessGrant returns the MongoDBEmployeeAccessGrant field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetMongoDBEmployeeAccessGrant() EmployeeAccessGrant {
	if o == nil || IsNil(o.MongoDBEmployeeAccessGrant) {
		var ret EmployeeAccessGrant
		return ret
	}
	return *o.MongoDBEmployeeAccessGrant
}

// GetMongoDBEmployeeAccessGrantOk returns a tuple with the MongoDBEmployeeAccessGrant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetMongoDBEmployeeAccessGrantOk() (*EmployeeAccessGrant, bool) {
	if o == nil || IsNil(o.MongoDBEmployeeAccessGrant) {
		return nil, false
	}

	return o.MongoDBEmployeeAccessGrant, true
}

// HasMongoDBEmployeeAccessGrant returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasMongoDBEmployeeAccessGrant() bool {
	if o != nil && !IsNil(o.MongoDBEmployeeAccessGrant) {
		return true
	}

	return false
}

// SetMongoDBEmployeeAccessGrant gets a reference to the given EmployeeAccessGrant and assigns it to the MongoDBEmployeeAccessGrant field.
func (o *ClusterDescription20240805) SetMongoDBEmployeeAccessGrant(v EmployeeAccessGrant) {
	o.MongoDBEmployeeAccessGrant = &v
}

// GetMongoDBMajorVersion returns the MongoDBMajorVersion field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetMongoDBMajorVersion() string {
	if o == nil || IsNil(o.MongoDBMajorVersion) {
		var ret string
		return ret
	}
	return *o.MongoDBMajorVersion
}

// GetMongoDBMajorVersionOk returns a tuple with the MongoDBMajorVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetMongoDBMajorVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongoDBMajorVersion) {
		return nil, false
	}

	return o.MongoDBMajorVersion, true
}

// HasMongoDBMajorVersion returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasMongoDBMajorVersion() bool {
	if o != nil && !IsNil(o.MongoDBMajorVersion) {
		return true
	}

	return false
}

// SetMongoDBMajorVersion gets a reference to the given string and assigns it to the MongoDBMajorVersion field.
func (o *ClusterDescription20240805) SetMongoDBMajorVersion(v string) {
	o.MongoDBMajorVersion = &v
}

// GetMongoDBVersion returns the MongoDBVersion field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetMongoDBVersion() string {
	if o == nil || IsNil(o.MongoDBVersion) {
		var ret string
		return ret
	}
	return *o.MongoDBVersion
}

// GetMongoDBVersionOk returns a tuple with the MongoDBVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetMongoDBVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongoDBVersion) {
		return nil, false
	}

	return o.MongoDBVersion, true
}

// HasMongoDBVersion returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasMongoDBVersion() bool {
	if o != nil && !IsNil(o.MongoDBVersion) {
		return true
	}

	return false
}

// SetMongoDBVersion gets a reference to the given string and assigns it to the MongoDBVersion field.
func (o *ClusterDescription20240805) SetMongoDBVersion(v string) {
	o.MongoDBVersion = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ClusterDescription20240805) SetName(v string) {
	o.Name = &v
}

// GetPaused returns the Paused field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetPaused() bool {
	if o == nil || IsNil(o.Paused) {
		var ret bool
		return ret
	}
	return *o.Paused
}

// GetPausedOk returns a tuple with the Paused field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetPausedOk() (*bool, bool) {
	if o == nil || IsNil(o.Paused) {
		return nil, false
	}

	return o.Paused, true
}

// HasPaused returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasPaused() bool {
	if o != nil && !IsNil(o.Paused) {
		return true
	}

	return false
}

// SetPaused gets a reference to the given bool and assigns it to the Paused field.
func (o *ClusterDescription20240805) SetPaused(v bool) {
	o.Paused = &v
}

// GetPitEnabled returns the PitEnabled field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetPitEnabled() bool {
	if o == nil || IsNil(o.PitEnabled) {
		var ret bool
		return ret
	}
	return *o.PitEnabled
}

// GetPitEnabledOk returns a tuple with the PitEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetPitEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.PitEnabled) {
		return nil, false
	}

	return o.PitEnabled, true
}

// HasPitEnabled returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasPitEnabled() bool {
	if o != nil && !IsNil(o.PitEnabled) {
		return true
	}

	return false
}

// SetPitEnabled gets a reference to the given bool and assigns it to the PitEnabled field.
func (o *ClusterDescription20240805) SetPitEnabled(v bool) {
	o.PitEnabled = &v
}

// GetRedactClientLogData returns the RedactClientLogData field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetRedactClientLogData() bool {
	if o == nil || IsNil(o.RedactClientLogData) {
		var ret bool
		return ret
	}
	return *o.RedactClientLogData
}

// GetRedactClientLogDataOk returns a tuple with the RedactClientLogData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetRedactClientLogDataOk() (*bool, bool) {
	if o == nil || IsNil(o.RedactClientLogData) {
		return nil, false
	}

	return o.RedactClientLogData, true
}

// HasRedactClientLogData returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasRedactClientLogData() bool {
	if o != nil && !IsNil(o.RedactClientLogData) {
		return true
	}

	return false
}

// SetRedactClientLogData gets a reference to the given bool and assigns it to the RedactClientLogData field.
func (o *ClusterDescription20240805) SetRedactClientLogData(v bool) {
	o.RedactClientLogData = &v
}

// GetReplicaSetScalingStrategy returns the ReplicaSetScalingStrategy field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetReplicaSetScalingStrategy() string {
	if o == nil || IsNil(o.ReplicaSetScalingStrategy) {
		var ret string
		return ret
	}
	return *o.ReplicaSetScalingStrategy
}

// GetReplicaSetScalingStrategyOk returns a tuple with the ReplicaSetScalingStrategy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetReplicaSetScalingStrategyOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetScalingStrategy) {
		return nil, false
	}

	return o.ReplicaSetScalingStrategy, true
}

// HasReplicaSetScalingStrategy returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasReplicaSetScalingStrategy() bool {
	if o != nil && !IsNil(o.ReplicaSetScalingStrategy) {
		return true
	}

	return false
}

// SetReplicaSetScalingStrategy gets a reference to the given string and assigns it to the ReplicaSetScalingStrategy field.
func (o *ClusterDescription20240805) SetReplicaSetScalingStrategy(v string) {
	o.ReplicaSetScalingStrategy = &v
}

// GetReplicationSpecs returns the ReplicationSpecs field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetReplicationSpecs() []ReplicationSpec20240805 {
	if o == nil || IsNil(o.ReplicationSpecs) {
		var ret []ReplicationSpec20240805
		return ret
	}
	return *o.ReplicationSpecs
}

// GetReplicationSpecsOk returns a tuple with the ReplicationSpecs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetReplicationSpecsOk() (*[]ReplicationSpec20240805, bool) {
	if o == nil || IsNil(o.ReplicationSpecs) {
		return nil, false
	}

	return o.ReplicationSpecs, true
}

// HasReplicationSpecs returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasReplicationSpecs() bool {
	if o != nil && !IsNil(o.ReplicationSpecs) {
		return true
	}

	return false
}

// SetReplicationSpecs gets a reference to the given []ReplicationSpec20240805 and assigns it to the ReplicationSpecs field.
func (o *ClusterDescription20240805) SetReplicationSpecs(v []ReplicationSpec20240805) {
	o.ReplicationSpecs = &v
}

// GetRetainBackups returns the RetainBackups field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetRetainBackups() bool {
	if o == nil || IsNil(o.RetainBackups) {
		var ret bool
		return ret
	}
	return *o.RetainBackups
}

// GetRetainBackupsOk returns a tuple with the RetainBackups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetRetainBackupsOk() (*bool, bool) {
	if o == nil || IsNil(o.RetainBackups) {
		return nil, false
	}

	return o.RetainBackups, true
}

// HasRetainBackups returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasRetainBackups() bool {
	if o != nil && !IsNil(o.RetainBackups) {
		return true
	}

	return false
}

// SetRetainBackups gets a reference to the given bool and assigns it to the RetainBackups field.
func (o *ClusterDescription20240805) SetRetainBackups(v bool) {
	o.RetainBackups = &v
}

// GetRootCertType returns the RootCertType field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetRootCertType() string {
	if o == nil || IsNil(o.RootCertType) {
		var ret string
		return ret
	}
	return *o.RootCertType
}

// GetRootCertTypeOk returns a tuple with the RootCertType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetRootCertTypeOk() (*string, bool) {
	if o == nil || IsNil(o.RootCertType) {
		return nil, false
	}

	return o.RootCertType, true
}

// HasRootCertType returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasRootCertType() bool {
	if o != nil && !IsNil(o.RootCertType) {
		return true
	}

	return false
}

// SetRootCertType gets a reference to the given string and assigns it to the RootCertType field.
func (o *ClusterDescription20240805) SetRootCertType(v string) {
	o.RootCertType = &v
}

// GetStateName returns the StateName field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetStateName() string {
	if o == nil || IsNil(o.StateName) {
		var ret string
		return ret
	}
	return *o.StateName
}

// GetStateNameOk returns a tuple with the StateName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetStateNameOk() (*string, bool) {
	if o == nil || IsNil(o.StateName) {
		return nil, false
	}

	return o.StateName, true
}

// HasStateName returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasStateName() bool {
	if o != nil && !IsNil(o.StateName) {
		return true
	}

	return false
}

// SetStateName gets a reference to the given string and assigns it to the StateName field.
func (o *ClusterDescription20240805) SetStateName(v string) {
	o.StateName = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetTags() []ResourceTag {
	if o == nil || IsNil(o.Tags) {
		var ret []ResourceTag
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetTagsOk() (*[]ResourceTag, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []ResourceTag and assigns it to the Tags field.
func (o *ClusterDescription20240805) SetTags(v []ResourceTag) {
	o.Tags = &v
}

// GetTerminationProtectionEnabled returns the TerminationProtectionEnabled field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetTerminationProtectionEnabled() bool {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		var ret bool
		return ret
	}
	return *o.TerminationProtectionEnabled
}

// GetTerminationProtectionEnabledOk returns a tuple with the TerminationProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetTerminationProtectionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.TerminationProtectionEnabled) {
		return nil, false
	}

	return o.TerminationProtectionEnabled, true
}

// HasTerminationProtectionEnabled returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasTerminationProtectionEnabled() bool {
	if o != nil && !IsNil(o.TerminationProtectionEnabled) {
		return true
	}

	return false
}

// SetTerminationProtectionEnabled gets a reference to the given bool and assigns it to the TerminationProtectionEnabled field.
func (o *ClusterDescription20240805) SetTerminationProtectionEnabled(v bool) {
	o.TerminationProtectionEnabled = &v
}

// GetUseAwsTimeBasedSnapshotCopyForFastInitialSync returns the UseAwsTimeBasedSnapshotCopyForFastInitialSync field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetUseAwsTimeBasedSnapshotCopyForFastInitialSync() bool {
	if o == nil || IsNil(o.UseAwsTimeBasedSnapshotCopyForFastInitialSync) {
		var ret bool
		return ret
	}
	return *o.UseAwsTimeBasedSnapshotCopyForFastInitialSync
}

// GetUseAwsTimeBasedSnapshotCopyForFastInitialSyncOk returns a tuple with the UseAwsTimeBasedSnapshotCopyForFastInitialSync field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetUseAwsTimeBasedSnapshotCopyForFastInitialSyncOk() (*bool, bool) {
	if o == nil || IsNil(o.UseAwsTimeBasedSnapshotCopyForFastInitialSync) {
		return nil, false
	}

	return o.UseAwsTimeBasedSnapshotCopyForFastInitialSync, true
}

// HasUseAwsTimeBasedSnapshotCopyForFastInitialSync returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasUseAwsTimeBasedSnapshotCopyForFastInitialSync() bool {
	if o != nil && !IsNil(o.UseAwsTimeBasedSnapshotCopyForFastInitialSync) {
		return true
	}

	return false
}

// SetUseAwsTimeBasedSnapshotCopyForFastInitialSync gets a reference to the given bool and assigns it to the UseAwsTimeBasedSnapshotCopyForFastInitialSync field.
func (o *ClusterDescription20240805) SetUseAwsTimeBasedSnapshotCopyForFastInitialSync(v bool) {
	o.UseAwsTimeBasedSnapshotCopyForFastInitialSync = &v
}

// GetVersionReleaseSystem returns the VersionReleaseSystem field value if set, zero value otherwise
func (o *ClusterDescription20240805) GetVersionReleaseSystem() string {
	if o == nil || IsNil(o.VersionReleaseSystem) {
		var ret string
		return ret
	}
	return *o.VersionReleaseSystem
}

// GetVersionReleaseSystemOk returns a tuple with the VersionReleaseSystem field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescription20240805) GetVersionReleaseSystemOk() (*string, bool) {
	if o == nil || IsNil(o.VersionReleaseSystem) {
		return nil, false
	}

	return o.VersionReleaseSystem, true
}

// HasVersionReleaseSystem returns a boolean if a field has been set.
func (o *ClusterDescription20240805) HasVersionReleaseSystem() bool {
	if o != nil && !IsNil(o.VersionReleaseSystem) {
		return true
	}

	return false
}

// SetVersionReleaseSystem gets a reference to the given string and assigns it to the VersionReleaseSystem field.
func (o *ClusterDescription20240805) SetVersionReleaseSystem(v string) {
	o.VersionReleaseSystem = &v
}
