// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupOnDemandSnapshotRequest struct for DiskBackupOnDemandSnapshotRequest
type DiskBackupOnDemandSnapshotRequest struct {
	// Human-readable phrase or sentence that explains the purpose of the snapshot. The resource returns this parameter when `\"status\" : \"onDemand\"`.
	Description *string `json:"description,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Number of days that MongoDB Cloud should retain the on-demand snapshot. Must be at least **1**.
	RetentionInDays *int `json:"retentionInDays,omitempty"`
}

// NewDiskBackupOnDemandSnapshotRequest instantiates a new DiskBackupOnDemandSnapshotRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupOnDemandSnapshotRequest() *DiskBackupOnDemandSnapshotRequest {
	this := DiskBackupOnDemandSnapshotRequest{}
	return &this
}

// NewDiskBackupOnDemandSnapshotRequestWithDefaults instantiates a new DiskBackupOnDemandSnapshotRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupOnDemandSnapshotRequestWithDefaults() *DiskBackupOnDemandSnapshotRequest {
	this := DiskBackupOnDemandSnapshotRequest{}
	return &this
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *DiskBackupOnDemandSnapshotRequest) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupOnDemandSnapshotRequest) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *DiskBackupOnDemandSnapshotRequest) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *DiskBackupOnDemandSnapshotRequest) SetDescription(v string) {
	o.Description = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupOnDemandSnapshotRequest) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupOnDemandSnapshotRequest) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupOnDemandSnapshotRequest) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupOnDemandSnapshotRequest) SetLinks(v []Link) {
	o.Links = &v
}

// GetRetentionInDays returns the RetentionInDays field value if set, zero value otherwise
func (o *DiskBackupOnDemandSnapshotRequest) GetRetentionInDays() int {
	if o == nil || IsNil(o.RetentionInDays) {
		var ret int
		return ret
	}
	return *o.RetentionInDays
}

// GetRetentionInDaysOk returns a tuple with the RetentionInDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupOnDemandSnapshotRequest) GetRetentionInDaysOk() (*int, bool) {
	if o == nil || IsNil(o.RetentionInDays) {
		return nil, false
	}

	return o.RetentionInDays, true
}

// HasRetentionInDays returns a boolean if a field has been set.
func (o *DiskBackupOnDemandSnapshotRequest) HasRetentionInDays() bool {
	if o != nil && !IsNil(o.RetentionInDays) {
		return true
	}

	return false
}

// SetRetentionInDays gets a reference to the given int and assigns it to the RetentionInDays field.
func (o *DiskBackupOnDemandSnapshotRequest) SetRetentionInDays(v int) {
	o.RetentionInDays = &v
}
