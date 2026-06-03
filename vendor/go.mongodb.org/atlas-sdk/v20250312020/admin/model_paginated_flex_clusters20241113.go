// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PaginatedFlexClusters20241113 struct for PaginatedFlexClusters20241113
type PaginatedFlexClusters20241113 struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of returned documents that MongoDB Cloud provides when completing this request.
	// Read only field.
	Results []FlexClusterDescription20241113 `json:"results"`
	// Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`. The total number is an estimate and may not be exact.
	// Read only field.
	TotalCount *int `json:"totalCount,omitempty"`
}

// NewPaginatedFlexClusters20241113 instantiates a new PaginatedFlexClusters20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaginatedFlexClusters20241113(results []FlexClusterDescription20241113) *PaginatedFlexClusters20241113 {
	this := PaginatedFlexClusters20241113{}
	this.Results = results
	return &this
}

// NewPaginatedFlexClusters20241113WithDefaults instantiates a new PaginatedFlexClusters20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginatedFlexClusters20241113WithDefaults() *PaginatedFlexClusters20241113 {
	this := PaginatedFlexClusters20241113{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PaginatedFlexClusters20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedFlexClusters20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PaginatedFlexClusters20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *PaginatedFlexClusters20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetResults returns the Results field value
func (o *PaginatedFlexClusters20241113) GetResults() []FlexClusterDescription20241113 {
	if o == nil {
		var ret []FlexClusterDescription20241113
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *PaginatedFlexClusters20241113) GetResultsOk() (*[]FlexClusterDescription20241113, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Results, true
}

// SetResults sets field value
func (o *PaginatedFlexClusters20241113) SetResults(v []FlexClusterDescription20241113) {
	o.Results = v
}

// GetTotalCount returns the TotalCount field value if set, zero value otherwise
func (o *PaginatedFlexClusters20241113) GetTotalCount() int {
	if o == nil || IsNil(o.TotalCount) {
		var ret int
		return ret
	}
	return *o.TotalCount
}

// GetTotalCountOk returns a tuple with the TotalCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaginatedFlexClusters20241113) GetTotalCountOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCount) {
		return nil, false
	}

	return o.TotalCount, true
}

// HasTotalCount returns a boolean if a field has been set.
func (o *PaginatedFlexClusters20241113) HasTotalCount() bool {
	if o != nil && !IsNil(o.TotalCount) {
		return true
	}

	return false
}

// SetTotalCount gets a reference to the given int and assigns it to the TotalCount field.
func (o *PaginatedFlexClusters20241113) SetTotalCount(v int) {
	o.TotalCount = &v
}
