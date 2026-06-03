// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DiskBackupExportJob struct for DiskBackupExportJob
type DiskBackupExportJob struct {
	// Information on the export job for each replica set in the sharded cluster.
	// Read only field.
	Components *[]DiskBackupExportMember `json:"components,omitempty"`
	// Date and time when a user or Atlas created the Export Job. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Collection of key-value pairs that represent custom data for the metadata file that MongoDB Cloud uploads when the Export Job finishes.
	CustomData *[]BackupLabel `json:"customData,omitempty"`
	// Unique 24-hexadecimal character string that identifies the Export Bucket.
	// Read only field.
	ExportBucketId string        `json:"exportBucketId"`
	ExportStatus   *ExportStatus `json:"exportStatus,omitempty"`
	// Date and time when this Export Job completed. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.
	// Read only field.
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	// Unique 24-hexadecimal character string that identifies the restore job.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Prefix used for all blob storage objects uploaded as part of the Export Job.
	// Read only field.
	Prefix *string `json:"prefix,omitempty"`
	// Unique 24-hexadecimal character string that identifies the snapshot.
	SnapshotId *string `json:"snapshotId,omitempty"`
	// State of the Export Job.
	// Read only field.
	State       *string      `json:"state,omitempty"`
	StateReason *StateReason `json:"stateReason,omitempty"`
}

// NewDiskBackupExportJob instantiates a new DiskBackupExportJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupExportJob(exportBucketId string) *DiskBackupExportJob {
	this := DiskBackupExportJob{}
	this.ExportBucketId = exportBucketId
	return &this
}

// NewDiskBackupExportJobWithDefaults instantiates a new DiskBackupExportJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupExportJobWithDefaults() *DiskBackupExportJob {
	this := DiskBackupExportJob{}
	return &this
}

// GetComponents returns the Components field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetComponents() []DiskBackupExportMember {
	if o == nil || IsNil(o.Components) {
		var ret []DiskBackupExportMember
		return ret
	}
	return *o.Components
}

// GetComponentsOk returns a tuple with the Components field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetComponentsOk() (*[]DiskBackupExportMember, bool) {
	if o == nil || IsNil(o.Components) {
		return nil, false
	}

	return o.Components, true
}

// HasComponents returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasComponents() bool {
	if o != nil && !IsNil(o.Components) {
		return true
	}

	return false
}

// SetComponents gets a reference to the given []DiskBackupExportMember and assigns it to the Components field.
func (o *DiskBackupExportJob) SetComponents(v []DiskBackupExportMember) {
	o.Components = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *DiskBackupExportJob) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetCustomData returns the CustomData field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetCustomData() []BackupLabel {
	if o == nil || IsNil(o.CustomData) {
		var ret []BackupLabel
		return ret
	}
	return *o.CustomData
}

// GetCustomDataOk returns a tuple with the CustomData field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetCustomDataOk() (*[]BackupLabel, bool) {
	if o == nil || IsNil(o.CustomData) {
		return nil, false
	}

	return o.CustomData, true
}

// HasCustomData returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasCustomData() bool {
	if o != nil && !IsNil(o.CustomData) {
		return true
	}

	return false
}

// SetCustomData gets a reference to the given []BackupLabel and assigns it to the CustomData field.
func (o *DiskBackupExportJob) SetCustomData(v []BackupLabel) {
	o.CustomData = &v
}

// GetExportBucketId returns the ExportBucketId field value
func (o *DiskBackupExportJob) GetExportBucketId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExportBucketId
}

// GetExportBucketIdOk returns a tuple with the ExportBucketId field value
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetExportBucketIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExportBucketId, true
}

// SetExportBucketId sets field value
func (o *DiskBackupExportJob) SetExportBucketId(v string) {
	o.ExportBucketId = v
}

// GetExportStatus returns the ExportStatus field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetExportStatus() ExportStatus {
	if o == nil || IsNil(o.ExportStatus) {
		var ret ExportStatus
		return ret
	}
	return *o.ExportStatus
}

// GetExportStatusOk returns a tuple with the ExportStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetExportStatusOk() (*ExportStatus, bool) {
	if o == nil || IsNil(o.ExportStatus) {
		return nil, false
	}

	return o.ExportStatus, true
}

// HasExportStatus returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasExportStatus() bool {
	if o != nil && !IsNil(o.ExportStatus) {
		return true
	}

	return false
}

// SetExportStatus gets a reference to the given ExportStatus and assigns it to the ExportStatus field.
func (o *DiskBackupExportJob) SetExportStatus(v ExportStatus) {
	o.ExportStatus = &v
}

// GetFinishedAt returns the FinishedAt field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetFinishedAt() time.Time {
	if o == nil || IsNil(o.FinishedAt) {
		var ret time.Time
		return ret
	}
	return *o.FinishedAt
}

// GetFinishedAtOk returns a tuple with the FinishedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetFinishedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FinishedAt) {
		return nil, false
	}

	return o.FinishedAt, true
}

// HasFinishedAt returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasFinishedAt() bool {
	if o != nil && !IsNil(o.FinishedAt) {
		return true
	}

	return false
}

// SetFinishedAt gets a reference to the given time.Time and assigns it to the FinishedAt field.
func (o *DiskBackupExportJob) SetFinishedAt(v time.Time) {
	o.FinishedAt = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DiskBackupExportJob) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupExportJob) SetLinks(v []Link) {
	o.Links = &v
}

// GetPrefix returns the Prefix field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetPrefix() string {
	if o == nil || IsNil(o.Prefix) {
		var ret string
		return ret
	}
	return *o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetPrefixOk() (*string, bool) {
	if o == nil || IsNil(o.Prefix) {
		return nil, false
	}

	return o.Prefix, true
}

// HasPrefix returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasPrefix() bool {
	if o != nil && !IsNil(o.Prefix) {
		return true
	}

	return false
}

// SetPrefix gets a reference to the given string and assigns it to the Prefix field.
func (o *DiskBackupExportJob) SetPrefix(v string) {
	o.Prefix = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}

	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *DiskBackupExportJob) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *DiskBackupExportJob) SetState(v string) {
	o.State = &v
}

// GetStateReason returns the StateReason field value if set, zero value otherwise
func (o *DiskBackupExportJob) GetStateReason() StateReason {
	if o == nil || IsNil(o.StateReason) {
		var ret StateReason
		return ret
	}
	return *o.StateReason
}

// GetStateReasonOk returns a tuple with the StateReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupExportJob) GetStateReasonOk() (*StateReason, bool) {
	if o == nil || IsNil(o.StateReason) {
		return nil, false
	}

	return o.StateReason, true
}

// HasStateReason returns a boolean if a field has been set.
func (o *DiskBackupExportJob) HasStateReason() bool {
	if o != nil && !IsNil(o.StateReason) {
		return true
	}

	return false
}

// SetStateReason gets a reference to the given StateReason and assigns it to the StateReason field.
func (o *DiskBackupExportJob) SetStateReason(v StateReason) {
	o.StateReason = &v
}
