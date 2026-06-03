// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupPaginatedEvent struct for GroupPaginatedEvent
type GroupPaginatedEvent struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of returned documents that MongoDB Cloud provides when completing this request.
	// Read only field.
	Results []EventViewForNdsGroup `json:"results"`
	// Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`. The total number is an estimate and may not be exact.
	// Read only field.
	TotalCount *int `json:"totalCount,omitempty"`
}

// NewGroupPaginatedEvent instantiates a new GroupPaginatedEvent object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupPaginatedEvent(results []EventViewForNdsGroup) *GroupPaginatedEvent {
	this := GroupPaginatedEvent{}
	this.Results = results
	return &this
}

// NewGroupPaginatedEventWithDefaults instantiates a new GroupPaginatedEvent object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupPaginatedEventWithDefaults() *GroupPaginatedEvent {
	this := GroupPaginatedEvent{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *GroupPaginatedEvent) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPaginatedEvent) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *GroupPaginatedEvent) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *GroupPaginatedEvent) SetLinks(v []Link) {
	o.Links = &v
}

// GetResults returns the Results field value
func (o *GroupPaginatedEvent) GetResults() []EventViewForNdsGroup {
	if o == nil {
		var ret []EventViewForNdsGroup
		return ret
	}

	return o.Results
}

// GetResultsOk returns a tuple with the Results field value
// and a boolean to check if the value has been set.
func (o *GroupPaginatedEvent) GetResultsOk() (*[]EventViewForNdsGroup, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Results, true
}

// SetResults sets field value
func (o *GroupPaginatedEvent) SetResults(v []EventViewForNdsGroup) {
	o.Results = v
}

// GetTotalCount returns the TotalCount field value if set, zero value otherwise
func (o *GroupPaginatedEvent) GetTotalCount() int {
	if o == nil || IsNil(o.TotalCount) {
		var ret int
		return ret
	}
	return *o.TotalCount
}

// GetTotalCountOk returns a tuple with the TotalCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupPaginatedEvent) GetTotalCountOk() (*int, bool) {
	if o == nil || IsNil(o.TotalCount) {
		return nil, false
	}

	return o.TotalCount, true
}

// HasTotalCount returns a boolean if a field has been set.
func (o *GroupPaginatedEvent) HasTotalCount() bool {
	if o != nil && !IsNil(o.TotalCount) {
		return true
	}

	return false
}

// SetTotalCount gets a reference to the given int and assigns it to the TotalCount field.
func (o *GroupPaginatedEvent) SetTotalCount(v int) {
	o.TotalCount = &v
}
