// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ExportStatus State of the Export Job.
type ExportStatus struct {
	// Count of collections whose documents were exported to the Export Bucket.
	// Read only field.
	ExportedCollections *int `json:"exportedCollections,omitempty"`
	// Total count of collections whose documents will be exported to the Export Bucket.
	// Read only field.
	TotalCollections *int `json:"totalCollections,omitempty"`
}

// NewExportStatus instantiates a new ExportStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExportStatus() *ExportStatus {
	this := ExportStatus{}
	return &this
}

// NewExportStatusWithDefaults instantiates a new ExportStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExportStatusWithDefaults() *ExportStatus {
	this := ExportStatus{}
	return &this
}

// GetExportedCollections returns the ExportedCollections field value if set, zero value otherwise
func (o *ExportStatus) GetExportedCollections() int {
	if o == nil || IsNil(o.ExportedCollections) {
		var ret int
		return ret
	}
	return *o.ExportedCollections
}

// GetExportedCollectionsOk returns a tuple with the ExportedCollections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportStatus) GetExportedCollectionsOk() (*int, bool) {
	if o == nil || IsNil(o.ExportedCollections) {
		return nil, false
	}

	return o.ExportedCollections, true
}

// HasExportedCollections returns a boolean if a field has been set.
func (o *ExportStatus) HasExportedCollections() bool {
	if o != nil && !IsNil(o.ExportedCollections) {
		return true
	}

	return false
}

// SetExportedCollections gets a reference to the given int and assigns it to the ExportedCollections field.
func (o *ExportStatus) SetExportedCollections(v int) {
	o.ExportedCollections = &v
}

// GetTotalCollections returns the TotalCollections field value if set, zero value otherwise
func (o *ExportStatus) GetTotalCollections() int {
	if o == nil || IsNil(o.TotalCollections) {
		var ret int
		return ret
	}
	return *o.TotalCollections
}

// GetTotalCollectionsOk returns a tuple with the TotalCollections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExportStatus) GetTotalCollectionsOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCollections) {
		return nil, false
	}

	return o.TotalCollections, true
}

// HasTotalCollections returns a boolean if a field has been set.
func (o *ExportStatus) HasTotalCollections() bool {
	if o != nil && !IsNil(o.TotalCollections) {
		return true
	}

	return false
}

// SetTotalCollections gets a reference to the given int and assigns it to the TotalCollections field.
func (o *ExportStatus) SetTotalCollections(v int) {
	o.TotalCollections = &v
}
