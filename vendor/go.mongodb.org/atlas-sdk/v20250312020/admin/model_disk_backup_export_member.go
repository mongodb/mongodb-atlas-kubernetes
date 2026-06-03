// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupExportMember struct for DiskBackupExportMember
type DiskBackupExportMember struct {
	// Unique 24-hexadecimal character string that identifies the the Cloud Backup snapshot export job for each shard in a sharded cluster.
	// Read only field.
	ExportId *string `json:"exportId,omitempty"`
	// Human-readable label that identifies the replica set on the sharded cluster.
	// Read only field.
	ReplicaSetName *string `json:"replicaSetName,omitempty"`
}

// NewDiskBackupExportMember instantiates a new DiskBackupExportMember object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupExportMember() *DiskBackupExportMember {
	this := DiskBackupExportMember{}
	return &this
}

// NewDiskBackupExportMemberWithDefaults instantiates a new DiskBackupExportMember object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupExportMemberWithDefaults() *DiskBackupExportMember {
	this := DiskBackupExportMember{}
	return &this
}

// GetExportId returns the ExportId field value if set, zero value otherwise
func (o *DiskBackupExportMember) GetExportId() string {
	if o == nil || IsNil(o.ExportId) {
		var ret string
		return ret
	}
	return *o.ExportId
}

// GetExportIdOk returns a tuple with the ExportId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportMember) GetExportIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExportId) {
		return nil, false
	}

	return o.ExportId, true
}

// HasExportId returns a boolean if a field has been set.
func (o *DiskBackupExportMember) HasExportId() bool {
	if o != nil && !IsNil(o.ExportId) {
		return true
	}

	return false
}

// SetExportId gets a reference to the given string and assigns it to the ExportId field.
func (o *DiskBackupExportMember) SetExportId(v string) {
	o.ExportId = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *DiskBackupExportMember) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportMember) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *DiskBackupExportMember) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *DiskBackupExportMember) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}
