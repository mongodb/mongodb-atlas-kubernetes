// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DiskBackupSnapshot struct for DiskBackupSnapshot
type DiskBackupSnapshot struct {
	// Date and time when MongoDB Cloud took the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Human-readable phrase or sentence that explains the purpose of the snapshot. The resource returns this parameter when `\"status\": \"onDemand\"`.
	// Read only field.
	Description *string `json:"description,omitempty"`
	// Date and time when MongoDB Cloud deletes the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Human-readable label that identifies how often this snapshot triggers.
	// Read only field.
	FrequencyType *string `json:"frequencyType,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique string that identifies the Amazon Web Services (AWS) Key Management Service (KMS) Customer Master Key (CMK) used to encrypt the snapshot. The resource returns this value when `\"encryptionEnabled\" : true`.
	// Read only field.
	MasterKeyUUID *string `json:"masterKeyUUID,omitempty"`
	// Version of the MongoDB host that this snapshot backs up.
	// Read only field.
	MongodVersion *string `json:"mongodVersion,omitempty"`
	// List that contains unique identifiers for the policy items.
	// Read only field.
	PolicyItems *[]string `json:"policyItems,omitempty"`
	// Human-readable label that identifies when this snapshot triggers.
	// Read only field.
	SnapshotType *string `json:"snapshotType,omitempty"`
	// Human-readable label that indicates the stage of the backup process for this snapshot.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Number of bytes taken to store the backup at time of snapshot.
	// Read only field.
	StorageSizeBytes *int64 `json:"storageSizeBytes,omitempty"`
	// Human-readable label that categorizes the cluster as a replica set or sharded cluster.
	// Read only field.
	Type *string `json:"type,omitempty"`
	// Human-readable label that identifies the cloud provider.
	// Read only field.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// List that identifies the regions to which MongoDB Cloud copies the snapshot.
	// Read only field.
	CopyRegions *[]string `json:"copyRegions,omitempty"`
	// Human-readable label that identifies the replica set from which MongoDB Cloud took this snapshot. The resource returns this parameter when `\"type\": \"replicaSet\"`.
	// Read only field.
	ReplicaSetName *string `json:"replicaSetName,omitempty"`
	// Describes a sharded cluster's config server type.
	// Read only field.
	ConfigServerType *string `json:"configServerType,omitempty"`
	// List that includes the snapshots and the cloud provider that stores the snapshots. The resource returns this parameter when `\"type\" : \"SHARDED_CLUSTER\"`.
	// Read only field.
	Members *[]DiskBackupShardedClusterSnapshotMember `json:"members,omitempty"`
	// List that contains the unique identifiers of the snapshots created for the shards and config host for a sharded cluster. The resource returns this parameter when `\"type\": \"SHARDED_CLUSTER\"`. These identifiers should match the ones specified in the **members[n].id** parameters. This allows you to map a snapshot to its shard or config host name.
	// Read only field.
	SnapshotIds *[]string `json:"snapshotIds,omitempty"`
}

// NewDiskBackupSnapshot instantiates a new DiskBackupSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshot() *DiskBackupSnapshot {
	this := DiskBackupSnapshot{}
	return &this
}

// NewDiskBackupSnapshotWithDefaults instantiates a new DiskBackupSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotWithDefaults() *DiskBackupSnapshot {
	this := DiskBackupSnapshot{}
	return &this
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *DiskBackupSnapshot) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *DiskBackupSnapshot) SetDescription(v string) {
	o.Description = &v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetExpiresAt() time.Time {
	if o == nil || IsNil(o.ExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpiresAt) {
		return nil, false
	}

	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasExpiresAt() bool {
	if o != nil && !IsNil(o.ExpiresAt) {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *DiskBackupSnapshot) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

// GetFrequencyType returns the FrequencyType field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetFrequencyType() string {
	if o == nil || IsNil(o.FrequencyType) {
		var ret string
		return ret
	}
	return *o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetFrequencyTypeOk() (*string, bool) {
	if o == nil || IsNil(o.FrequencyType) {
		return nil, false
	}

	return o.FrequencyType, true
}

// HasFrequencyType returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasFrequencyType() bool {
	if o != nil && !IsNil(o.FrequencyType) {
		return true
	}

	return false
}

// SetFrequencyType gets a reference to the given string and assigns it to the FrequencyType field.
func (o *DiskBackupSnapshot) SetFrequencyType(v string) {
	o.FrequencyType = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DiskBackupSnapshot) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshot) SetLinks(v []Link) {
	o.Links = &v
}

// GetMasterKeyUUID returns the MasterKeyUUID field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetMasterKeyUUID() string {
	if o == nil || IsNil(o.MasterKeyUUID) {
		var ret string
		return ret
	}
	return *o.MasterKeyUUID
}

// GetMasterKeyUUIDOk returns a tuple with the MasterKeyUUID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetMasterKeyUUIDOk() (*string, bool) {
	if o == nil || IsNil(o.MasterKeyUUID) {
		return nil, false
	}

	return o.MasterKeyUUID, true
}

// HasMasterKeyUUID returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasMasterKeyUUID() bool {
	if o != nil && !IsNil(o.MasterKeyUUID) {
		return true
	}

	return false
}

// SetMasterKeyUUID gets a reference to the given string and assigns it to the MasterKeyUUID field.
func (o *DiskBackupSnapshot) SetMasterKeyUUID(v string) {
	o.MasterKeyUUID = &v
}

// GetMongodVersion returns the MongodVersion field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetMongodVersion() string {
	if o == nil || IsNil(o.MongodVersion) {
		var ret string
		return ret
	}
	return *o.MongodVersion
}

// GetMongodVersionOk returns a tuple with the MongodVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetMongodVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongodVersion) {
		return nil, false
	}

	return o.MongodVersion, true
}

// HasMongodVersion returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasMongodVersion() bool {
	if o != nil && !IsNil(o.MongodVersion) {
		return true
	}

	return false
}

// SetMongodVersion gets a reference to the given string and assigns it to the MongodVersion field.
func (o *DiskBackupSnapshot) SetMongodVersion(v string) {
	o.MongodVersion = &v
}

// GetPolicyItems returns the PolicyItems field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetPolicyItems() []string {
	if o == nil || IsNil(o.PolicyItems) {
		var ret []string
		return ret
	}
	return *o.PolicyItems
}

// GetPolicyItemsOk returns a tuple with the PolicyItems field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetPolicyItemsOk() (*[]string, bool) {
	if o == nil || IsNil(o.PolicyItems) {
		return nil, false
	}

	return o.PolicyItems, true
}

// HasPolicyItems returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasPolicyItems() bool {
	if o != nil && !IsNil(o.PolicyItems) {
		return true
	}

	return false
}

// SetPolicyItems gets a reference to the given []string and assigns it to the PolicyItems field.
func (o *DiskBackupSnapshot) SetPolicyItems(v []string) {
	o.PolicyItems = &v
}

// GetSnapshotType returns the SnapshotType field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetSnapshotType() string {
	if o == nil || IsNil(o.SnapshotType) {
		var ret string
		return ret
	}
	return *o.SnapshotType
}

// GetSnapshotTypeOk returns a tuple with the SnapshotType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetSnapshotTypeOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotType) {
		return nil, false
	}

	return o.SnapshotType, true
}

// HasSnapshotType returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasSnapshotType() bool {
	if o != nil && !IsNil(o.SnapshotType) {
		return true
	}

	return false
}

// SetSnapshotType gets a reference to the given string and assigns it to the SnapshotType field.
func (o *DiskBackupSnapshot) SetSnapshotType(v string) {
	o.SnapshotType = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *DiskBackupSnapshot) SetStatus(v string) {
	o.Status = &v
}

// GetStorageSizeBytes returns the StorageSizeBytes field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetStorageSizeBytes() int64 {
	if o == nil || IsNil(o.StorageSizeBytes) {
		var ret int64
		return ret
	}
	return *o.StorageSizeBytes
}

// GetStorageSizeBytesOk returns a tuple with the StorageSizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetStorageSizeBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.StorageSizeBytes) {
		return nil, false
	}

	return o.StorageSizeBytes, true
}

// HasStorageSizeBytes returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasStorageSizeBytes() bool {
	if o != nil && !IsNil(o.StorageSizeBytes) {
		return true
	}

	return false
}

// SetStorageSizeBytes gets a reference to the given int64 and assigns it to the StorageSizeBytes field.
func (o *DiskBackupSnapshot) SetStorageSizeBytes(v int64) {
	o.StorageSizeBytes = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *DiskBackupSnapshot) SetType(v string) {
	o.Type = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *DiskBackupSnapshot) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetCopyRegions returns the CopyRegions field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetCopyRegions() []string {
	if o == nil || IsNil(o.CopyRegions) {
		var ret []string
		return ret
	}
	return *o.CopyRegions
}

// GetCopyRegionsOk returns a tuple with the CopyRegions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetCopyRegionsOk() (*[]string, bool) {
	if o == nil || IsNil(o.CopyRegions) {
		return nil, false
	}

	return o.CopyRegions, true
}

// HasCopyRegions returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasCopyRegions() bool {
	if o != nil && !IsNil(o.CopyRegions) {
		return true
	}

	return false
}

// SetCopyRegions gets a reference to the given []string and assigns it to the CopyRegions field.
func (o *DiskBackupSnapshot) SetCopyRegions(v []string) {
	o.CopyRegions = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *DiskBackupSnapshot) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}

// GetConfigServerType returns the ConfigServerType field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetConfigServerType() string {
	if o == nil || IsNil(o.ConfigServerType) {
		var ret string
		return ret
	}
	return *o.ConfigServerType
}

// GetConfigServerTypeOk returns a tuple with the ConfigServerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetConfigServerTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ConfigServerType) {
		return nil, false
	}

	return o.ConfigServerType, true
}

// HasConfigServerType returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasConfigServerType() bool {
	if o != nil && !IsNil(o.ConfigServerType) {
		return true
	}

	return false
}

// SetConfigServerType gets a reference to the given string and assigns it to the ConfigServerType field.
func (o *DiskBackupSnapshot) SetConfigServerType(v string) {
	o.ConfigServerType = &v
}

// GetMembers returns the Members field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetMembers() []DiskBackupShardedClusterSnapshotMember {
	if o == nil || IsNil(o.Members) {
		var ret []DiskBackupShardedClusterSnapshotMember
		return ret
	}
	return *o.Members
}

// GetMembersOk returns a tuple with the Members field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetMembersOk() (*[]DiskBackupShardedClusterSnapshotMember, bool) {
	if o == nil || IsNil(o.Members) {
		return nil, false
	}

	return o.Members, true
}

// HasMembers returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasMembers() bool {
	if o != nil && !IsNil(o.Members) {
		return true
	}

	return false
}

// SetMembers gets a reference to the given []DiskBackupShardedClusterSnapshotMember and assigns it to the Members field.
func (o *DiskBackupSnapshot) SetMembers(v []DiskBackupShardedClusterSnapshotMember) {
	o.Members = &v
}

// GetSnapshotIds returns the SnapshotIds field value if set, zero value otherwise
func (o *DiskBackupSnapshot) GetSnapshotIds() []string {
	if o == nil || IsNil(o.SnapshotIds) {
		var ret []string
		return ret
	}
	return *o.SnapshotIds
}

// GetSnapshotIdsOk returns a tuple with the SnapshotIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshot) GetSnapshotIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.SnapshotIds) {
		return nil, false
	}

	return o.SnapshotIds, true
}

// HasSnapshotIds returns a boolean if a field has been set.
func (o *DiskBackupSnapshot) HasSnapshotIds() bool {
	if o != nil && !IsNil(o.SnapshotIds) {
		return true
	}

	return false
}

// SetSnapshotIds gets a reference to the given []string and assigns it to the SnapshotIds field.
func (o *DiskBackupSnapshot) SetSnapshotIds(v []string) {
	o.SnapshotIds = &v
}
