// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PaginatedHostViewAtlas struct for PaginatedHostViewAtlas
type PaginatedHostViewAtlas struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]LinkAtlas `json:"links,omitempty"`
	// List of returned documents that MongoDB Cloud provides when completing this request.
	// Read only field.
	Results []ApiHostViewAtlas `json:"results"`
	// Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`. The total number is an estimate and may not be exact.
	// Read only field.
	TotalCount *int `json:"totalCount,omitempty"`
}

// NewPaginatedHostViewAtlas instantiates a new PaginatedHostViewAtlas object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginatedHostViewAtlas(results []ApiHostViewAtlas) *PaginatedHostViewAtlas {
	this := PaginatedHostViewAtlas{}
	this.Results = results
	return &this
}

// NewPaginatedHostViewAtlasWithDefaults instantiates a new PaginatedHostViewAtlas object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginatedHostViewAtlasWithDefaults() *PaginatedHostViewAtlas {
	this := PaginatedHostViewAtlas{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PaginatedHostViewAtlas) GetLinks() []LinkAtlas {
	if o == nil || IsNil(o.Links) {
		var ret []LinkAtlas
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedHostViewAtlas) GetLinksOk() (*[]LinkAtlas, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PaginatedHostViewAtlas) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []LinkAtlas and assigns it to the Links field.
func (o *PaginatedHostViewAtlas) SetLinks(v []LinkAtlas) {
	o.Links = &v
}

// GetResults returns the Results field value
func (o *PaginatedHostViewAtlas) GetResults() []ApiHostViewAtlas {
	if o == nil {
		var ret []ApiHostViewAtlas
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginatedHostViewAtlas) GetResultsOk() (*[]ApiHostViewAtlas, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Results, true
}

// SetResults sets field value
func (o *PaginatedHostViewAtlas) SetResults(v []ApiHostViewAtlas) {
	o.Results = v
}

// GetTotalCount returns the TotalCount field value if set, zero value otherwise
func (o *PaginatedHostViewAtlas) GetTotalCount() int {
	if o == nil || IsNil(o.TotalCount) {
		var ret int
		return ret
	}
	return *o.TotalCount
}

// GetTotalCountOk returns a tuple with the TotalCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedHostViewAtlas) GetTotalCountOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCount) {
		return nil, false
	}

	return o.TotalCount, true
}

// HasTotalCount returns a boolean if a field has been set.
func (o *PaginatedHostViewAtlas) HasTotalCount() bool {
	if o != nil && !IsNil(o.TotalCount) {
		return true
	}

	return false
}

// SetTotalCount gets a reference to the given int and assigns it to the TotalCount field.
func (o *PaginatedHostViewAtlas) SetTotalCount(v int) {
	o.TotalCount = &v
}
