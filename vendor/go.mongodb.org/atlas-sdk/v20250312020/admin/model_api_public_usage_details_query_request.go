// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiPublicUsageDetailsQueryRequest Request body for an Invoice Usage Details query with filtering, pagination, and sort.
type ApiPublicUsageDetailsQueryRequest struct {
	Filters *UsageDetailsFilterRequest `json:"filters,omitempty"`
	// Specify the field used to specify how to sort query results. Default to bill date.
	SortField *string `json:"sortField,omitempty"`
	// Specify the sort order (ascending / descending) used to specify how to sort query results. Defaults to descending.
	SortOrder *string `json:"sortOrder,omitempty"`
}

// NewApiPublicUsageDetailsQueryRequest instantiates a new ApiPublicUsageDetailsQueryRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiPublicUsageDetailsQueryRequest() *ApiPublicUsageDetailsQueryRequest {
	this := ApiPublicUsageDetailsQueryRequest{}
	return &this
}

// NewApiPublicUsageDetailsQueryRequestWithDefaults instantiates a new ApiPublicUsageDetailsQueryRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiPublicUsageDetailsQueryRequestWithDefaults() *ApiPublicUsageDetailsQueryRequest {
	this := ApiPublicUsageDetailsQueryRequest{}
	return &this
}

// GetFilters returns the Filters field value if set, zero value otherwise
func (o *ApiPublicUsageDetailsQueryRequest) GetFilters() UsageDetailsFilterRequest {
	if o == nil || IsNil(o.Filters) {
		var ret UsageDetailsFilterRequest
		return ret
	}
	return *o.Filters
}

// GetFiltersOk returns a tuple with the Filters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiPublicUsageDetailsQueryRequest) GetFiltersOk() (*UsageDetailsFilterRequest, bool) {
	if o == nil || IsNil(o.Filters) {
		return nil, false
	}

	return o.Filters, true
}

// HasFilters returns a boolean if a field has been set.
func (o *ApiPublicUsageDetailsQueryRequest) HasFilters() bool {
	if o != nil && !IsNil(o.Filters) {
		return true
	}

	return false
}

// SetFilters gets a reference to the given UsageDetailsFilterRequest and assigns it to the Filters field.
func (o *ApiPublicUsageDetailsQueryRequest) SetFilters(v UsageDetailsFilterRequest) {
	o.Filters = &v
}

// GetSortField returns the SortField field value if set, zero value otherwise
func (o *ApiPublicUsageDetailsQueryRequest) GetSortField() string {
	if o == nil || IsNil(o.SortField) {
		var ret string
		return ret
	}
	return *o.SortField
}

// GetSortFieldOk returns a tuple with the SortField field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiPublicUsageDetailsQueryRequest) GetSortFieldOk() (*string, bool) {
	if o == nil || IsNil(o.SortField) {
		return nil, false
	}

	return o.SortField, true
}

// HasSortField returns a boolean if a field has been set.
func (o *ApiPublicUsageDetailsQueryRequest) HasSortField() bool {
	if o != nil && !IsNil(o.SortField) {
		return true
	}

	return false
}

// SetSortField gets a reference to the given string and assigns it to the SortField field.
func (o *ApiPublicUsageDetailsQueryRequest) SetSortField(v string) {
	o.SortField = &v
}

// GetSortOrder returns the SortOrder field value if set, zero value otherwise
func (o *ApiPublicUsageDetailsQueryRequest) GetSortOrder() string {
	if o == nil || IsNil(o.SortOrder) {
		var ret string
		return ret
	}
	return *o.SortOrder
}

// GetSortOrderOk returns a tuple with the SortOrder field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiPublicUsageDetailsQueryRequest) GetSortOrderOk() (*string, bool) {
	if o == nil || IsNil(o.SortOrder) {
		return nil, false
	}

	return o.SortOrder, true
}

// HasSortOrder returns a boolean if a field has been set.
func (o *ApiPublicUsageDetailsQueryRequest) HasSortOrder() bool {
	if o != nil && !IsNil(o.SortOrder) {
		return true
	}

	return false
}

// SetSortOrder gets a reference to the given string and assigns it to the SortOrder field.
func (o *ApiPublicUsageDetailsQueryRequest) SetSortOrder(v string) {
	o.SortOrder = &v
}
