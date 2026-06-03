// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupExportJobRequest struct for DiskBackupExportJobRequest
type DiskBackupExportJobRequest struct {
	// Collection of key-value pairs that represent custom data to add to the metadata file that MongoDB Cloud uploads to the bucket when the export job finishes.
	CustomData *[]BackupLabel `json:"customData,omitempty"`
	// Unique 24-hexadecimal character string that identifies the Export Bucket.
	// Write only field.
	ExportBucketId string `json:"exportBucketId"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal character string that identifies the Cloud Backup Snapshot to export.
	// Write only field.
	SnapshotId string `json:"snapshotId"`
}

// NewDiskBackupExportJobRequest instantiates a new DiskBackupExportJobRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupExportJobRequest(exportBucketId string, snapshotId string) *DiskBackupExportJobRequest {
	this := DiskBackupExportJobRequest{}
	this.ExportBucketId = exportBucketId
	this.SnapshotId = snapshotId
	return &this
}

// NewDiskBackupExportJobRequestWithDefaults instantiates a new DiskBackupExportJobRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupExportJobRequestWithDefaults() *DiskBackupExportJobRequest {
	this := DiskBackupExportJobRequest{}
	return &this
}

// GetCustomData returns the CustomData field value if set, zero value otherwise
func (o *DiskBackupExportJobRequest) GetCustomData() []BackupLabel {
	if o == nil || IsNil(o.CustomData) {
		var ret []BackupLabel
		return ret
	}
	return *o.CustomData
}

// GetCustomDataOk returns a tuple with the CustomData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJobRequest) GetCustomDataOk() (*[]BackupLabel, bool) {
	if o == nil || IsNil(o.CustomData) {
		return nil, false
	}

	return o.CustomData, true
}

// HasCustomData returns a boolean if a field has been set.
func (o *DiskBackupExportJobRequest) HasCustomData() bool {
	if o != nil && !IsNil(o.CustomData) {
		return true
	}

	return false
}

// SetCustomData gets a reference to the given []BackupLabel and assigns it to the CustomData field.
func (o *DiskBackupExportJobRequest) SetCustomData(v []BackupLabel) {
	o.CustomData = &v
}

// GetExportBucketId returns the ExportBucketId field value
func (o *DiskBackupExportJobRequest) GetExportBucketId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExportBucketId
}

// GetExportBucketIdOk returns a tuple with the ExportBucketId field value
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJobRequest) GetExportBucketIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExportBucketId, true
}

// SetExportBucketId sets field value
func (o *DiskBackupExportJobRequest) SetExportBucketId(v string) {
	o.ExportBucketId = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupExportJobRequest) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJobRequest) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupExportJobRequest) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupExportJobRequest) SetLinks(v []Link) {
	o.Links = &v
}

// GetSnapshotId returns the SnapshotId field value
func (o *DiskBackupExportJobRequest) GetSnapshotId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJobRequest) GetSnapshotIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SnapshotId, true
}

// SetSnapshotId sets field value
func (o *DiskBackupExportJobRequest) SetSnapshotId(v string) {
	o.SnapshotId = v
}
