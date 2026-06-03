// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BackupRestoreJob struct for BackupRestoreJob
type BackupRestoreJob struct {
	// Unique 24-hexadecimal digit string that identifies the batch to which this restore job belongs. This parameter exists only for a sharded cluster restore.
	// Read only field.
	BatchId *string `json:"batchId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the sharded cluster checkpoint. The checkpoint represents the point in time back to which you want to restore you data. This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`. Use this parameter with sharded clusters only.  - If you set `checkpointId`, you can't set `oplogInc`, `oplogTs`, `snapshotId`, or `pointInTimeUTCMillis`. - If you provide this parameter, this endpoint restores all data up to this checkpoint to the database you specify in the `delivery` object.
	// Write only field.
	CheckpointId *string `json:"checkpointId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the cluster with the snapshot you want to return. This parameter returns for restore clusters.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Human-readable label that identifies the cluster containing the snapshots you want to retrieve.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Date and time when someone requested this restore job. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created  *time.Time               `json:"created,omitempty"`
	Delivery BackupRestoreJobDelivery `json:"delivery"`
	// Unique 24-hexadecimal digit string that identifies the an imported deployment job. This parameter exists when restoring from an imported snapshot/cluster shot.
	// Read only field.
	DeploymentJobId *string `json:"deploymentJobId,omitempty"`
	// Flag that indicates whether someone encrypted the data in the restored snapshot.
	// Read only field.
	EncryptionEnabled *bool `json:"encryptionEnabled,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that owns the snapshots.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List that contains documents mapping each restore file to a hashed checksum. This parameter applies after you download the corresponding `delivery.url`. If `\"methodName\" : \"HTTP\"`, this list contains one object that represents the hash of the `.tar.gz` file.
	// Read only field.
	Hashes *[]RestoreJobFileHash `json:"hashes,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the restore job.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Universally Unique Identifier (UUID) that identifies the Key Management Interoperability (KMIP) master key used to encrypt the snapshot data. This parameter applies only when `\"encryptionEnabled\" : \"true\"`.
	// Read only field.
	MasterKeyUUID *string `json:"masterKeyUUID,omitempty"`
	// Thirty-two-bit incrementing ordinal that represents operations within a given second. When paired with `oplogTs`, this represents the point in time to which MongoDB Cloud restores your data. This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`.  - If you set `oplogInc`, you must set `oplogTs`, and can't set `checkpointId`, `snapshotId`, or `pointInTimeUTCMillis`. - If you provide this parameter, this endpoint restores all data up to and including this Oplog timestamp to the database you specified in the `delivery` object.
	// Write only field.
	OplogInc *int `json:"oplogInc,omitempty"`
	// Date and time from which you want to restore this snapshot. This parameter expresses its value in ISO 8601 format in UTC. This represents the first part of an Oplog timestamp. When paired with `oplogInc`, they represent the last database operation to which you want to restore your data. This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`. Run a query against `local.oplog.rs` on your replica set to find the desired timestamp.  - If you set `oplogTs`, you must set `oplogInc`, and you can't set `checkpointId`, `snapshotId`, or `pointInTimeUTCMillis`. - If you provide this parameter, this endpoint restores all data up to and including this Oplog timestamp to the database you specified in the `delivery` object.
	// Write only field.
	OplogTs *string `json:"oplogTs,omitempty"`
	// Timestamp from which you want to restore this snapshot. This parameter expresses its value in the number of milliseconds elapsed since the UNIX epoch. This timestamp must fall within the last 24 hours of the current time. This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`.  - If you provide this parameter, this endpoint restores all data up to this point in time to the database you specified in the `delivery` object. - If you set `pointInTimeUTCMillis`, you can't set `oplogInc`, `oplogTs`, `snapshotId`, or `checkpointId`.
	// Write only field.
	PointInTimeUTCMillis *int64 `json:"pointInTimeUTCMillis,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot to restore. If you set `snapshotId`, you can't set `oplogInc`, `oplogTs`, `pointInTimeUTCMillis`, or `checkpointId`.
	SnapshotId *string `json:"snapshotId,omitempty"`
	// Human-readable label that identifies the status of the downloadable file at the time of the request.
	// Read only field.
	StatusName *string           `json:"statusName,omitempty"`
	Timestamp  *ApiBSONTimestamp `json:"timestamp,omitempty"`
}

// NewBackupRestoreJob instantiates a new BackupRestoreJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupRestoreJob(delivery BackupRestoreJobDelivery) *BackupRestoreJob {
	this := BackupRestoreJob{}
	this.Delivery = delivery
	return &this
}

// NewBackupRestoreJobWithDefaults instantiates a new BackupRestoreJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupRestoreJobWithDefaults() *BackupRestoreJob {
	this := BackupRestoreJob{}
	return &this
}

// GetBatchId returns the BatchId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetBatchId() string {
	if o == nil || IsNil(o.BatchId) {
		var ret string
		return ret
	}
	return *o.BatchId
}

// GetBatchIdOk returns a tuple with the BatchId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetBatchIdOk() (*string, bool) {
	if o == nil || IsNil(o.BatchId) {
		return nil, false
	}

	return o.BatchId, true
}

// HasBatchId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasBatchId() bool {
	if o != nil && !IsNil(o.BatchId) {
		return true
	}

	return false
}

// SetBatchId gets a reference to the given string and assigns it to the BatchId field.
func (o *BackupRestoreJob) SetBatchId(v string) {
	o.BatchId = &v
}

// GetCheckpointId returns the CheckpointId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetCheckpointId() string {
	if o == nil || IsNil(o.CheckpointId) {
		var ret string
		return ret
	}
	return *o.CheckpointId
}

// GetCheckpointIdOk returns a tuple with the CheckpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetCheckpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.CheckpointId) {
		return nil, false
	}

	return o.CheckpointId, true
}

// HasCheckpointId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasCheckpointId() bool {
	if o != nil && !IsNil(o.CheckpointId) {
		return true
	}

	return false
}

// SetCheckpointId gets a reference to the given string and assigns it to the CheckpointId field.
func (o *BackupRestoreJob) SetCheckpointId(v string) {
	o.CheckpointId = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *BackupRestoreJob) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *BackupRestoreJob) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *BackupRestoreJob) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BackupRestoreJob) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *BackupRestoreJob) SetCreated(v time.Time) {
	o.Created = &v
}

// GetDelivery returns the Delivery field value
func (o *BackupRestoreJob) GetDelivery() BackupRestoreJobDelivery {
	if o == nil {
		var ret BackupRestoreJobDelivery
		return ret
	}

	return o.Delivery
}

// GetDeliveryOk returns a tuple with the Delivery field value
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetDeliveryOk() (*BackupRestoreJobDelivery, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Delivery, true
}

// SetDelivery sets field value
func (o *BackupRestoreJob) SetDelivery(v BackupRestoreJobDelivery) {
	o.Delivery = v
}

// GetDeploymentJobId returns the DeploymentJobId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetDeploymentJobId() string {
	if o == nil || IsNil(o.DeploymentJobId) {
		var ret string
		return ret
	}
	return *o.DeploymentJobId
}

// GetDeploymentJobIdOk returns a tuple with the DeploymentJobId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetDeploymentJobIdOk() (*string, bool) {
	if o == nil || IsNil(o.DeploymentJobId) {
		return nil, false
	}

	return o.DeploymentJobId, true
}

// HasDeploymentJobId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasDeploymentJobId() bool {
	if o != nil && !IsNil(o.DeploymentJobId) {
		return true
	}

	return false
}

// SetDeploymentJobId gets a reference to the given string and assigns it to the DeploymentJobId field.
func (o *BackupRestoreJob) SetDeploymentJobId(v string) {
	o.DeploymentJobId = &v
}

// GetEncryptionEnabled returns the EncryptionEnabled field value if set, zero value otherwise
func (o *BackupRestoreJob) GetEncryptionEnabled() bool {
	if o == nil || IsNil(o.EncryptionEnabled) {
		var ret bool
		return ret
	}
	return *o.EncryptionEnabled
}

// GetEncryptionEnabledOk returns a tuple with the EncryptionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetEncryptionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.EncryptionEnabled) {
		return nil, false
	}

	return o.EncryptionEnabled, true
}

// HasEncryptionEnabled returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasEncryptionEnabled() bool {
	if o != nil && !IsNil(o.EncryptionEnabled) {
		return true
	}

	return false
}

// SetEncryptionEnabled gets a reference to the given bool and assigns it to the EncryptionEnabled field.
func (o *BackupRestoreJob) SetEncryptionEnabled(v bool) {
	o.EncryptionEnabled = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *BackupRestoreJob) SetGroupId(v string) {
	o.GroupId = &v
}

// GetHashes returns the Hashes field value if set, zero value otherwise
func (o *BackupRestoreJob) GetHashes() []RestoreJobFileHash {
	if o == nil || IsNil(o.Hashes) {
		var ret []RestoreJobFileHash
		return ret
	}
	return *o.Hashes
}

// GetHashesOk returns a tuple with the Hashes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetHashesOk() (*[]RestoreJobFileHash, bool) {
	if o == nil || IsNil(o.Hashes) {
		return nil, false
	}

	return o.Hashes, true
}

// HasHashes returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasHashes() bool {
	if o != nil && !IsNil(o.Hashes) {
		return true
	}

	return false
}

// SetHashes gets a reference to the given []RestoreJobFileHash and assigns it to the Hashes field.
func (o *BackupRestoreJob) SetHashes(v []RestoreJobFileHash) {
	o.Hashes = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupRestoreJob) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupRestoreJob) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *BackupRestoreJob) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *BackupRestoreJob) SetLinks(v []Link) {
	o.Links = &v
}

// GetMasterKeyUUID returns the MasterKeyUUID field value if set, zero value otherwise
func (o *BackupRestoreJob) GetMasterKeyUUID() string {
	if o == nil || IsNil(o.MasterKeyUUID) {
		var ret string
		return ret
	}
	return *o.MasterKeyUUID
}

// GetMasterKeyUUIDOk returns a tuple with the MasterKeyUUID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetMasterKeyUUIDOk() (*string, bool) {
	if o == nil || IsNil(o.MasterKeyUUID) {
		return nil, false
	}

	return o.MasterKeyUUID, true
}

// HasMasterKeyUUID returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasMasterKeyUUID() bool {
	if o != nil && !IsNil(o.MasterKeyUUID) {
		return true
	}

	return false
}

// SetMasterKeyUUID gets a reference to the given string and assigns it to the MasterKeyUUID field.
func (o *BackupRestoreJob) SetMasterKeyUUID(v string) {
	o.MasterKeyUUID = &v
}

// GetOplogInc returns the OplogInc field value if set, zero value otherwise
func (o *BackupRestoreJob) GetOplogInc() int {
	if o == nil || IsNil(o.OplogInc) {
		var ret int
		return ret
	}
	return *o.OplogInc
}

// GetOplogIncOk returns a tuple with the OplogInc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetOplogIncOk() (*int, bool) {
	if o == nil || IsNil(o.OplogInc) {
		return nil, false
	}

	return o.OplogInc, true
}

// HasOplogInc returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasOplogInc() bool {
	if o != nil && !IsNil(o.OplogInc) {
		return true
	}

	return false
}

// SetOplogInc gets a reference to the given int and assigns it to the OplogInc field.
func (o *BackupRestoreJob) SetOplogInc(v int) {
	o.OplogInc = &v
}

// GetOplogTs returns the OplogTs field value if set, zero value otherwise
func (o *BackupRestoreJob) GetOplogTs() string {
	if o == nil || IsNil(o.OplogTs) {
		var ret string
		return ret
	}
	return *o.OplogTs
}

// GetOplogTsOk returns a tuple with the OplogTs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetOplogTsOk() (*string, bool) {
	if o == nil || IsNil(o.OplogTs) {
		return nil, false
	}

	return o.OplogTs, true
}

// HasOplogTs returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasOplogTs() bool {
	if o != nil && !IsNil(o.OplogTs) {
		return true
	}

	return false
}

// SetOplogTs gets a reference to the given string and assigns it to the OplogTs field.
func (o *BackupRestoreJob) SetOplogTs(v string) {
	o.OplogTs = &v
}

// GetPointInTimeUTCMillis returns the PointInTimeUTCMillis field value if set, zero value otherwise
func (o *BackupRestoreJob) GetPointInTimeUTCMillis() int64 {
	if o == nil || IsNil(o.PointInTimeUTCMillis) {
		var ret int64
		return ret
	}
	return *o.PointInTimeUTCMillis
}

// GetPointInTimeUTCMillisOk returns a tuple with the PointInTimeUTCMillis field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetPointInTimeUTCMillisOk() (*int64, bool) {
	if o == nil || IsNil(o.PointInTimeUTCMillis) {
		return nil, false
	}

	return o.PointInTimeUTCMillis, true
}

// HasPointInTimeUTCMillis returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasPointInTimeUTCMillis() bool {
	if o != nil && !IsNil(o.PointInTimeUTCMillis) {
		return true
	}

	return false
}

// SetPointInTimeUTCMillis gets a reference to the given int64 and assigns it to the PointInTimeUTCMillis field.
func (o *BackupRestoreJob) SetPointInTimeUTCMillis(v int64) {
	o.PointInTimeUTCMillis = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise
func (o *BackupRestoreJob) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}

	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *BackupRestoreJob) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BackupRestoreJob) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BackupRestoreJob) SetStatusName(v string) {
	o.StatusName = &v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise
func (o *BackupRestoreJob) GetTimestamp() ApiBSONTimestamp {
	if o == nil || IsNil(o.Timestamp) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJob) GetTimestampOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}

	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *BackupRestoreJob) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given ApiBSONTimestamp and assigns it to the Timestamp field.
func (o *BackupRestoreJob) SetTimestamp(v ApiBSONTimestamp) {
	o.Timestamp = &v
}
