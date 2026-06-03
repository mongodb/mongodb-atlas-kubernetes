// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ServerlessBackupSnapshot struct for ServerlessBackupSnapshot
type ServerlessBackupSnapshot struct {
	// Date and time when MongoDB Cloud took the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
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
	// Version of the MongoDB host that this snapshot backs up.
	// Read only field.
	MongodVersion *string `json:"mongodVersion,omitempty"`
	// Human-readable label given to the serverless instance from which MongoDB Cloud took this snapshot.
	// Read only field.
	ServerlessInstanceName *string `json:"serverlessInstanceName,omitempty"`
	// Human-readable label that identifies when this snapshot triggers.
	// Read only field.
	SnapshotType *string `json:"snapshotType,omitempty"`
	// Human-readable label that indicates the stage of the backup process for this snapshot.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Number of bytes taken to store the backup snapshot.
	// Read only field.
	StorageSizeBytes *int64 `json:"storageSizeBytes,omitempty"`
}

// NewServerlessBackupSnapshot instantiates a new ServerlessBackupSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessBackupSnapshot() *ServerlessBackupSnapshot {
	this := ServerlessBackupSnapshot{}
	return &this
}

// NewServerlessBackupSnapshotWithDefaults instantiates a new ServerlessBackupSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessBackupSnapshotWithDefaults() *ServerlessBackupSnapshot {
	this := ServerlessBackupSnapshot{}
	return &this
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *ServerlessBackupSnapshot) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetExpiresAt() time.Time {
	if o == nil || IsNil(o.ExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpiresAt) {
		return nil, false
	}

	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasExpiresAt() bool {
	if o != nil && !IsNil(o.ExpiresAt) {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *ServerlessBackupSnapshot) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

// GetFrequencyType returns the FrequencyType field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetFrequencyType() string {
	if o == nil || IsNil(o.FrequencyType) {
		var ret string
		return ret
	}
	return *o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetFrequencyTypeOk() (*string, bool) {
	if o == nil || IsNil(o.FrequencyType) {
		return nil, false
	}

	return o.FrequencyType, true
}

// HasFrequencyType returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasFrequencyType() bool {
	if o != nil && !IsNil(o.FrequencyType) {
		return true
	}

	return false
}

// SetFrequencyType gets a reference to the given string and assigns it to the FrequencyType field.
func (o *ServerlessBackupSnapshot) SetFrequencyType(v string) {
	o.FrequencyType = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ServerlessBackupSnapshot) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ServerlessBackupSnapshot) SetLinks(v []Link) {
	o.Links = &v
}

// GetMongodVersion returns the MongodVersion field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetMongodVersion() string {
	if o == nil || IsNil(o.MongodVersion) {
		var ret string
		return ret
	}
	return *o.MongodVersion
}

// GetMongodVersionOk returns a tuple with the MongodVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetMongodVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongodVersion) {
		return nil, false
	}

	return o.MongodVersion, true
}

// HasMongodVersion returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasMongodVersion() bool {
	if o != nil && !IsNil(o.MongodVersion) {
		return true
	}

	return false
}

// SetMongodVersion gets a reference to the given string and assigns it to the MongodVersion field.
func (o *ServerlessBackupSnapshot) SetMongodVersion(v string) {
	o.MongodVersion = &v
}

// GetServerlessInstanceName returns the ServerlessInstanceName field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetServerlessInstanceName() string {
	if o == nil || IsNil(o.ServerlessInstanceName) {
		var ret string
		return ret
	}
	return *o.ServerlessInstanceName
}

// GetServerlessInstanceNameOk returns a tuple with the ServerlessInstanceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetServerlessInstanceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServerlessInstanceName) {
		return nil, false
	}

	return o.ServerlessInstanceName, true
}

// HasServerlessInstanceName returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasServerlessInstanceName() bool {
	if o != nil && !IsNil(o.ServerlessInstanceName) {
		return true
	}

	return false
}

// SetServerlessInstanceName gets a reference to the given string and assigns it to the ServerlessInstanceName field.
func (o *ServerlessBackupSnapshot) SetServerlessInstanceName(v string) {
	o.ServerlessInstanceName = &v
}

// GetSnapshotType returns the SnapshotType field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetSnapshotType() string {
	if o == nil || IsNil(o.SnapshotType) {
		var ret string
		return ret
	}
	return *o.SnapshotType
}

// GetSnapshotTypeOk returns a tuple with the SnapshotType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetSnapshotTypeOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotType) {
		return nil, false
	}

	return o.SnapshotType, true
}

// HasSnapshotType returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasSnapshotType() bool {
	if o != nil && !IsNil(o.SnapshotType) {
		return true
	}

	return false
}

// SetSnapshotType gets a reference to the given string and assigns it to the SnapshotType field.
func (o *ServerlessBackupSnapshot) SetSnapshotType(v string) {
	o.SnapshotType = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *ServerlessBackupSnapshot) SetStatus(v string) {
	o.Status = &v
}

// GetStorageSizeBytes returns the StorageSizeBytes field value if set, zero value otherwise
func (o *ServerlessBackupSnapshot) GetStorageSizeBytes() int64 {
	if o == nil || IsNil(o.StorageSizeBytes) {
		var ret int64
		return ret
	}
	return *o.StorageSizeBytes
}

// GetStorageSizeBytesOk returns a tuple with the StorageSizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupSnapshot) GetStorageSizeBytesOk() (*int64, bool) {
	if o == nil || IsNil(o.StorageSizeBytes) {
		return nil, false
	}

	return o.StorageSizeBytes, true
}

// HasStorageSizeBytes returns a boolean if a field has been set.
func (o *ServerlessBackupSnapshot) HasStorageSizeBytes() bool {
	if o != nil && !IsNil(o.StorageSizeBytes) {
		return true
	}

	return false
}

// SetStorageSizeBytes gets a reference to the given int64 and assigns it to the StorageSizeBytes field.
func (o *ServerlessBackupSnapshot) SetStorageSizeBytes(v int64) {
	o.StorageSizeBytes = &v
}
