// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PaginatedApiAtlasServerlessBackupSnapshot struct for PaginatedApiAtlasServerlessBackupSnapshot
type PaginatedApiAtlasServerlessBackupSnapshot struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of returned documents that MongoDB Cloud provides when completing this request.
	// Read only field.
	Results []ServerlessBackupSnapshot `json:"results"`
	// Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`. The total number is an estimate and may not be exact.
	// Read only field.
	TotalCount *int `json:"totalCount,omitempty"`
}

// NewPaginatedApiAtlasServerlessBackupSnapshot instantiates a new PaginatedApiAtlasServerlessBackupSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginatedApiAtlasServerlessBackupSnapshot(results []ServerlessBackupSnapshot) *PaginatedApiAtlasServerlessBackupSnapshot {
	this := PaginatedApiAtlasServerlessBackupSnapshot{}
	this.Results = results
	return &this
}

// NewPaginatedApiAtlasServerlessBackupSnapshotWithDefaults instantiates a new PaginatedApiAtlasServerlessBackupSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginatedApiAtlasServerlessBackupSnapshotWithDefaults() *PaginatedApiAtlasServerlessBackupSnapshot {
	this := PaginatedApiAtlasServerlessBackupSnapshot{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) SetLinks(v []Link) {
	o.Links = &v
}

// GetResults returns the Results field value
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetResults() []ServerlessBackupSnapshot {
	if o == nil {
		var ret []ServerlessBackupSnapshot
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetResultsOk() (*[]ServerlessBackupSnapshot, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Results, true
}

// SetResults sets field value
func (o *PaginatedApiAtlasServerlessBackupSnapshot) SetResults(v []ServerlessBackupSnapshot) {
	o.Results = v
}

// GetTotalCount returns the TotalCount field value if set, zero value otherwise
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetTotalCount() int {
	if o == nil || IsNil(o.TotalCount) {
		var ret int
		return ret
	}
	return *o.TotalCount
}

// GetTotalCountOk returns a tuple with the TotalCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) GetTotalCountOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCount) {
		return nil, false
	}

	return o.TotalCount, true
}

// HasTotalCount returns a boolean if a field has been set.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) HasTotalCount() bool {
	if o != nil && !IsNil(o.TotalCount) {
		return true
	}

	return false
}

// SetTotalCount gets a reference to the given int and assigns it to the TotalCount field.
func (o *PaginatedApiAtlasServerlessBackupSnapshot) SetTotalCount(v int) {
	o.TotalCount = &v
}
