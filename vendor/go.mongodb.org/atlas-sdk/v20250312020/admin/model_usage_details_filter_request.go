// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UsageDetailsFilterRequest Request body which contains various fields to filter line items as part of certain Invoice Usage Details queries.
type UsageDetailsFilterRequest struct {
	// The inclusive billing start date for usage details filter.
	BillEndDate *string `json:"billEndDate,omitempty"`
	// The inclusive billing start date for usage details filter.
	BillStartDate *string `json:"billStartDate,omitempty"`
	// The list of unique cluster ids to be included in the Usage Details filter.
	ClusterIds *[]string `json:"clusterIds,omitempty"`
	// The list of groups to be included in the Usage Details filter.
	GroupIds *[]string `json:"groupIds,omitempty"`
	// Whether zero cent line items should be included.
	IncludeZeroCentLineItems *bool `json:"includeZeroCentLineItems,omitempty"`
	// The list of projects to be included in the Cost Explorer Query.
	SkuServices *[]string `json:"skuServices,omitempty"`
	// The inclusive billing start date for usage details filter.
	UsageEndDate *string `json:"usageEndDate,omitempty"`
	// The inclusive usage start date for usage details filter.
	UsageStartDate *string `json:"usageStartDate,omitempty"`
}

// NewUsageDetailsFilterRequest instantiates a new UsageDetailsFilterRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUsageDetailsFilterRequest() *UsageDetailsFilterRequest {
	this := UsageDetailsFilterRequest{}
	return &this
}

// NewUsageDetailsFilterRequestWithDefaults instantiates a new UsageDetailsFilterRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUsageDetailsFilterRequestWithDefaults() *UsageDetailsFilterRequest {
	this := UsageDetailsFilterRequest{}
	return &this
}

// GetBillEndDate returns the BillEndDate field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetBillEndDate() string {
	if o == nil || IsNil(o.BillEndDate) {
		var ret string
		return ret
	}
	return *o.BillEndDate
}

// GetBillEndDateOk returns a tuple with the BillEndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetBillEndDateOk() (*string, bool) {
	if o == nil || IsNil(o.BillEndDate) {
		return nil, false
	}

	return o.BillEndDate, true
}

// HasBillEndDate returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasBillEndDate() bool {
	if o != nil && !IsNil(o.BillEndDate) {
		return true
	}

	return false
}

// SetBillEndDate gets a reference to the given string and assigns it to the BillEndDate field.
func (o *UsageDetailsFilterRequest) SetBillEndDate(v string) {
	o.BillEndDate = &v
}

// GetBillStartDate returns the BillStartDate field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetBillStartDate() string {
	if o == nil || IsNil(o.BillStartDate) {
		var ret string
		return ret
	}
	return *o.BillStartDate
}

// GetBillStartDateOk returns a tuple with the BillStartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetBillStartDateOk() (*string, bool) {
	if o == nil || IsNil(o.BillStartDate) {
		return nil, false
	}

	return o.BillStartDate, true
}

// HasBillStartDate returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasBillStartDate() bool {
	if o != nil && !IsNil(o.BillStartDate) {
		return true
	}

	return false
}

// SetBillStartDate gets a reference to the given string and assigns it to the BillStartDate field.
func (o *UsageDetailsFilterRequest) SetBillStartDate(v string) {
	o.BillStartDate = &v
}

// GetClusterIds returns the ClusterIds field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetClusterIds() []string {
	if o == nil || IsNil(o.ClusterIds) {
		var ret []string
		return ret
	}
	return *o.ClusterIds
}

// GetClusterIdsOk returns a tuple with the ClusterIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetClusterIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.ClusterIds) {
		return nil, false
	}

	return o.ClusterIds, true
}

// HasClusterIds returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasClusterIds() bool {
	if o != nil && !IsNil(o.ClusterIds) {
		return true
	}

	return false
}

// SetClusterIds gets a reference to the given []string and assigns it to the ClusterIds field.
func (o *UsageDetailsFilterRequest) SetClusterIds(v []string) {
	o.ClusterIds = &v
}

// GetGroupIds returns the GroupIds field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetGroupIds() []string {
	if o == nil || IsNil(o.GroupIds) {
		var ret []string
		return ret
	}
	return *o.GroupIds
}

// GetGroupIdsOk returns a tuple with the GroupIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetGroupIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.GroupIds) {
		return nil, false
	}

	return o.GroupIds, true
}

// HasGroupIds returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasGroupIds() bool {
	if o != nil && !IsNil(o.GroupIds) {
		return true
	}

	return false
}

// SetGroupIds gets a reference to the given []string and assigns it to the GroupIds field.
func (o *UsageDetailsFilterRequest) SetGroupIds(v []string) {
	o.GroupIds = &v
}

// GetIncludeZeroCentLineItems returns the IncludeZeroCentLineItems field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetIncludeZeroCentLineItems() bool {
	if o == nil || IsNil(o.IncludeZeroCentLineItems) {
		var ret bool
		return ret
	}
	return *o.IncludeZeroCentLineItems
}

// GetIncludeZeroCentLineItemsOk returns a tuple with the IncludeZeroCentLineItems field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetIncludeZeroCentLineItemsOk() (*bool, bool) {
	if o == nil || IsNil(o.IncludeZeroCentLineItems) {
		return nil, false
	}

	return o.IncludeZeroCentLineItems, true
}

// HasIncludeZeroCentLineItems returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasIncludeZeroCentLineItems() bool {
	if o != nil && !IsNil(o.IncludeZeroCentLineItems) {
		return true
	}

	return false
}

// SetIncludeZeroCentLineItems gets a reference to the given bool and assigns it to the IncludeZeroCentLineItems field.
func (o *UsageDetailsFilterRequest) SetIncludeZeroCentLineItems(v bool) {
	o.IncludeZeroCentLineItems = &v
}

// GetSkuServices returns the SkuServices field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetSkuServices() []string {
	if o == nil || IsNil(o.SkuServices) {
		var ret []string
		return ret
	}
	return *o.SkuServices
}

// GetSkuServicesOk returns a tuple with the SkuServices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetSkuServicesOk() (*[]string, bool) {
	if o == nil || IsNil(o.SkuServices) {
		return nil, false
	}

	return o.SkuServices, true
}

// HasSkuServices returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasSkuServices() bool {
	if o != nil && !IsNil(o.SkuServices) {
		return true
	}

	return false
}

// SetSkuServices gets a reference to the given []string and assigns it to the SkuServices field.
func (o *UsageDetailsFilterRequest) SetSkuServices(v []string) {
	o.SkuServices = &v
}

// GetUsageEndDate returns the UsageEndDate field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetUsageEndDate() string {
	if o == nil || IsNil(o.UsageEndDate) {
		var ret string
		return ret
	}
	return *o.UsageEndDate
}

// GetUsageEndDateOk returns a tuple with the UsageEndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetUsageEndDateOk() (*string, bool) {
	if o == nil || IsNil(o.UsageEndDate) {
		return nil, false
	}

	return o.UsageEndDate, true
}

// HasUsageEndDate returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasUsageEndDate() bool {
	if o != nil && !IsNil(o.UsageEndDate) {
		return true
	}

	return false
}

// SetUsageEndDate gets a reference to the given string and assigns it to the UsageEndDate field.
func (o *UsageDetailsFilterRequest) SetUsageEndDate(v string) {
	o.UsageEndDate = &v
}

// GetUsageStartDate returns the UsageStartDate field value if set, zero value otherwise
func (o *UsageDetailsFilterRequest) GetUsageStartDate() string {
	if o == nil || IsNil(o.UsageStartDate) {
		var ret string
		return ret
	}
	return *o.UsageStartDate
}

// GetUsageStartDateOk returns a tuple with the UsageStartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UsageDetailsFilterRequest) GetUsageStartDateOk() (*string, bool) {
	if o == nil || IsNil(o.UsageStartDate) {
		return nil, false
	}

	return o.UsageStartDate, true
}

// HasUsageStartDate returns a boolean if a field has been set.
func (o *UsageDetailsFilterRequest) HasUsageStartDate() bool {
	if o != nil && !IsNil(o.UsageStartDate) {
		return true
	}

	return false
}

// SetUsageStartDate gets a reference to the given string and assigns it to the UsageStartDate field.
func (o *UsageDetailsFilterRequest) SetUsageStartDate(v string) {
	o.UsageStartDate = &v
}
