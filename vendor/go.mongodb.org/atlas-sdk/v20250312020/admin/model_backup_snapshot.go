// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BackupSnapshot struct for BackupSnapshot
type BackupSnapshot struct {
	// Unique 24-hexadecimal digit string that identifies the cluster with the snapshots you want to return.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Human-readable label that identifies the cluster.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Flag that indicates whether the snapshot exists. This flag returns `false` while MongoDB Cloud creates the snapshot.
	// Read only field.
	Complete *bool             `json:"complete,omitempty"`
	Created  *ApiBSONTimestamp `json:"created,omitempty"`
	// Flag that indicates whether someone can delete this snapshot. You can't set `\"doNotDelete\" : true` and set a timestamp for **expires** in the same request.
	DoNotDelete *bool `json:"doNotDelete,omitempty"`
	// Date and time when MongoDB Cloud deletes the snapshot. If `\"doNotDelete\" : true`, MongoDB Cloud removes any value set for this parameter. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	Expires *time.Time `json:"expires,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that owns the snapshots.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Flag indicating if this is an incremental or a full snapshot.
	// Read only field.
	Incremental               *bool             `json:"incremental,omitempty"`
	LastOplogAppliedTimestamp *ApiBSONTimestamp `json:"lastOplogAppliedTimestamp,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Metadata that describes the complete snapshot.  - For a replica set, this array contains a single document. - For a sharded cluster, this array contains one document for each shard plus one document for the config host.
	// Read only field.
	Parts *[]BackupSnapshotPart `json:"parts,omitempty"`
}

// NewBackupSnapshot instantiates a new BackupSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupSnapshot() *BackupSnapshot {
	this := BackupSnapshot{}
	return &this
}

// NewBackupSnapshotWithDefaults instantiates a new BackupSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupSnapshotWithDefaults() *BackupSnapshot {
	this := BackupSnapshot{}
	return &this
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *BackupSnapshot) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *BackupSnapshot) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *BackupSnapshot) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *BackupSnapshot) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *BackupSnapshot) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *BackupSnapshot) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetComplete returns the Complete field value if set, zero value otherwise
func (o *BackupSnapshot) GetComplete() bool {
	if o == nil || IsNil(o.Complete) {
		var ret bool
		return ret
	}
	return *o.Complete
}

// GetCompleteOk returns a tuple with the Complete field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetCompleteOk() (*bool, bool) {
	if o == nil || IsNil(o.Complete) {
		return nil, false
	}

	return o.Complete, true
}

// HasComplete returns a boolean if a field has been set.
func (o *BackupSnapshot) HasComplete() bool {
	if o != nil && !IsNil(o.Complete) {
		return true
	}

	return false
}

// SetComplete gets a reference to the given bool and assigns it to the Complete field.
func (o *BackupSnapshot) SetComplete(v bool) {
	o.Complete = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BackupSnapshot) GetCreated() ApiBSONTimestamp {
	if o == nil || IsNil(o.Created) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetCreatedOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BackupSnapshot) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given ApiBSONTimestamp and assigns it to the Created field.
func (o *BackupSnapshot) SetCreated(v ApiBSONTimestamp) {
	o.Created = &v
}

// GetDoNotDelete returns the DoNotDelete field value if set, zero value otherwise
func (o *BackupSnapshot) GetDoNotDelete() bool {
	if o == nil || IsNil(o.DoNotDelete) {
		var ret bool
		return ret
	}
	return *o.DoNotDelete
}

// GetDoNotDeleteOk returns a tuple with the DoNotDelete field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetDoNotDeleteOk() (*bool, bool) {
	if o == nil || IsNil(o.DoNotDelete) {
		return nil, false
	}

	return o.DoNotDelete, true
}

// HasDoNotDelete returns a boolean if a field has been set.
func (o *BackupSnapshot) HasDoNotDelete() bool {
	if o != nil && !IsNil(o.DoNotDelete) {
		return true
	}

	return false
}

// SetDoNotDelete gets a reference to the given bool and assigns it to the DoNotDelete field.
func (o *BackupSnapshot) SetDoNotDelete(v bool) {
	o.DoNotDelete = &v
}

// GetExpires returns the Expires field value if set, zero value otherwise
func (o *BackupSnapshot) GetExpires() time.Time {
	if o == nil || IsNil(o.Expires) {
		var ret time.Time
		return ret
	}
	return *o.Expires
}

// GetExpiresOk returns a tuple with the Expires field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetExpiresOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Expires) {
		return nil, false
	}

	return o.Expires, true
}

// HasExpires returns a boolean if a field has been set.
func (o *BackupSnapshot) HasExpires() bool {
	if o != nil && !IsNil(o.Expires) {
		return true
	}

	return false
}

// SetExpires gets a reference to the given time.Time and assigns it to the Expires field.
func (o *BackupSnapshot) SetExpires(v time.Time) {
	o.Expires = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *BackupSnapshot) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *BackupSnapshot) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *BackupSnapshot) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupSnapshot) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupSnapshot) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupSnapshot) SetId(v string) {
	o.Id = &v
}

// GetIncremental returns the Incremental field value if set, zero value otherwise
func (o *BackupSnapshot) GetIncremental() bool {
	if o == nil || IsNil(o.Incremental) {
		var ret bool
		return ret
	}
	return *o.Incremental
}

// GetIncrementalOk returns a tuple with the Incremental field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetIncrementalOk() (*bool, bool) {
	if o == nil || IsNil(o.Incremental) {
		return nil, false
	}

	return o.Incremental, true
}

// HasIncremental returns a boolean if a field has been set.
func (o *BackupSnapshot) HasIncremental() bool {
	if o != nil && !IsNil(o.Incremental) {
		return true
	}

	return false
}

// SetIncremental gets a reference to the given bool and assigns it to the Incremental field.
func (o *BackupSnapshot) SetIncremental(v bool) {
	o.Incremental = &v
}

// GetLastOplogAppliedTimestamp returns the LastOplogAppliedTimestamp field value if set, zero value otherwise
func (o *BackupSnapshot) GetLastOplogAppliedTimestamp() ApiBSONTimestamp {
	if o == nil || IsNil(o.LastOplogAppliedTimestamp) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.LastOplogAppliedTimestamp
}

// GetLastOplogAppliedTimestampOk returns a tuple with the LastOplogAppliedTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetLastOplogAppliedTimestampOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.LastOplogAppliedTimestamp) {
		return nil, false
	}

	return o.LastOplogAppliedTimestamp, true
}

// HasLastOplogAppliedTimestamp returns a boolean if a field has been set.
func (o *BackupSnapshot) HasLastOplogAppliedTimestamp() bool {
	if o != nil && !IsNil(o.LastOplogAppliedTimestamp) {
		return true
	}

	return false
}

// SetLastOplogAppliedTimestamp gets a reference to the given ApiBSONTimestamp and assigns it to the LastOplogAppliedTimestamp field.
func (o *BackupSnapshot) SetLastOplogAppliedTimestamp(v ApiBSONTimestamp) {
	o.LastOplogAppliedTimestamp = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *BackupSnapshot) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *BackupSnapshot) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *BackupSnapshot) SetLinks(v []Link) {
	o.Links = &v
}

// GetParts returns the Parts field value if set, zero value otherwise
func (o *BackupSnapshot) GetParts() []BackupSnapshotPart {
	if o == nil || IsNil(o.Parts) {
		var ret []BackupSnapshotPart
		return ret
	}
	return *o.Parts
}

// GetPartsOk returns a tuple with the Parts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupSnapshot) GetPartsOk() (*[]BackupSnapshotPart, bool) {
	if o == nil || IsNil(o.Parts) {
		return nil, false
	}

	return o.Parts, true
}

// HasParts returns a boolean if a field has been set.
func (o *BackupSnapshot) HasParts() bool {
	if o != nil && !IsNil(o.Parts) {
		return true
	}

	return false
}

// SetParts gets a reference to the given []BackupSnapshotPart and assigns it to the Parts field.
func (o *BackupSnapshot) SetParts(v []BackupSnapshotPart) {
	o.Parts = &v
}
