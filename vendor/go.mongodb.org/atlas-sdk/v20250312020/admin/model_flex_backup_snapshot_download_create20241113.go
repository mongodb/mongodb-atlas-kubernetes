// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FlexBackupSnapshotDownloadCreate20241113 Details for one backup snapshot download of a flex cluster.
type FlexBackupSnapshotDownloadCreate20241113 struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot to download.
	// Write only field.
	SnapshotId string `json:"snapshotId"`
}

// NewFlexBackupSnapshotDownloadCreate20241113 instantiates a new FlexBackupSnapshotDownloadCreate20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexBackupSnapshotDownloadCreate20241113(snapshotId string) *FlexBackupSnapshotDownloadCreate20241113 {
	this := FlexBackupSnapshotDownloadCreate20241113{}
	this.SnapshotId = snapshotId
	return &this
}

// NewFlexBackupSnapshotDownloadCreate20241113WithDefaults instantiates a new FlexBackupSnapshotDownloadCreate20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexBackupSnapshotDownloadCreate20241113WithDefaults() *FlexBackupSnapshotDownloadCreate20241113 {
	this := FlexBackupSnapshotDownloadCreate20241113{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *FlexBackupSnapshotDownloadCreate20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshotDownloadCreate20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *FlexBackupSnapshotDownloadCreate20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *FlexBackupSnapshotDownloadCreate20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetSnapshotId returns the SnapshotId field value
func (o *FlexBackupSnapshotDownloadCreate20241113) GetSnapshotId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshotDownloadCreate20241113) GetSnapshotIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SnapshotId, true
}

// SetSnapshotId sets field value
func (o *FlexBackupSnapshotDownloadCreate20241113) SetSnapshotId(v string) {
	o.SnapshotId = v
}
