// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// InvoiceLineItem One service included in this invoice.
type InvoiceLineItem struct {
	// Human-readable label that identifies the cluster that incurred the charge.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Date and time when MongoDB Cloud created this line item. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Sum by which MongoDB discounted this line item. MongoDB Cloud expresses this value in cents (100ths of one US Dollar). The resource returns this parameter when a discount applies.
	// Read only field.
	DiscountCents *int64 `json:"discountCents,omitempty"`
	// Date and time when when MongoDB Cloud finished charging for this line item. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	EndDate *time.Time `json:"endDate,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project associated to this line item.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Human-readable label that identifies the project.
	GroupName *string `json:"groupName,omitempty"`
	// Comment that applies to this line item.
	// Read only field.
	Note *string `json:"note,omitempty"`
	// Percentage by which MongoDB discounted this line item. The resource returns this parameter when a discount applies.
	// Read only field.
	PercentDiscount *float32 `json:"percentDiscount,omitempty"`
	// Number of units included for the line item. These can be expressions of storage (GB), time (hours), or other units.
	// Read only field.
	Quantity *float64 `json:"quantity,omitempty"`
	// Human-readable description of the service that this line item provided. This Stock Keeping Unit (SKU) could be the instance type, a support charge, advanced security, or another service.
	// Read only field.
	Sku *string `json:"sku,omitempty"`
	// Date and time when MongoDB Cloud began charging for this line item. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	StartDate *time.Time `json:"startDate,omitempty"`
	// Human-readable label that identifies the Atlas App Services application associated with this line item.
	// Read only field.
	StitchAppName *string `json:"stitchAppName,omitempty"`
	// A map of key-value pairs corresponding to the tags associated with the line item resource.
	// Read only field.
	Tags *map[string][]string `json:"tags,omitempty"`
	// Lower bound for usage amount range in current SKU tier.   **NOTE**: `lineItems[n].tierLowerBound` appears only if your `lineItems[n].sku` is tiered.
	// Read only field.
	TierLowerBound *float64 `json:"tierLowerBound,omitempty"`
	// Upper bound for usage amount range in current SKU tier.   **NOTE**: `lineItems[n].tierUpperBound` appears only if your `lineItems[n].sku` is tiered.
	// Read only field.
	TierUpperBound *float64 `json:"tierUpperBound,omitempty"`
	// Sum of the cost set for this line item. MongoDB Cloud expresses this value in cents (100ths of one US Dollar) and calculates this value as `unitPriceDollars` * `quantity` * 100.
	// Read only field.
	TotalPriceCents *int64 `json:"totalPriceCents,omitempty"`
	// Element used to express what **quantity** this line item measures. This value can be elements of time, storage capacity, and the like.
	// Read only field.
	Unit *string `json:"unit,omitempty"`
	// Value per **unit** for this line item expressed in US Dollars.
	// Read only field.
	UnitPriceDollars *float64 `json:"unitPriceDollars,omitempty"`
}

// NewInvoiceLineItem instantiates a new InvoiceLineItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInvoiceLineItem() *InvoiceLineItem {
	this := InvoiceLineItem{}
	return &this
}

// NewInvoiceLineItemWithDefaults instantiates a new InvoiceLineItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInvoiceLineItemWithDefaults() *InvoiceLineItem {
	this := InvoiceLineItem{}
	return &this
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *InvoiceLineItem) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *InvoiceLineItem) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *InvoiceLineItem) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *InvoiceLineItem) SetCreated(v time.Time) {
	o.Created = &v
}

// GetDiscountCents returns the DiscountCents field value if set, zero value otherwise
func (o *InvoiceLineItem) GetDiscountCents() int64 {
	if o == nil || IsNil(o.DiscountCents) {
		var ret int64
		return ret
	}
	return *o.DiscountCents
}

// GetDiscountCentsOk returns a tuple with the DiscountCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetDiscountCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.DiscountCents) {
		return nil, false
	}

	return o.DiscountCents, true
}

// HasDiscountCents returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasDiscountCents() bool {
	if o != nil && !IsNil(o.DiscountCents) {
		return true
	}

	return false
}

// SetDiscountCents gets a reference to the given int64 and assigns it to the DiscountCents field.
func (o *InvoiceLineItem) SetDiscountCents(v int64) {
	o.DiscountCents = &v
}

// GetEndDate returns the EndDate field value if set, zero value otherwise
func (o *InvoiceLineItem) GetEndDate() time.Time {
	if o == nil || IsNil(o.EndDate) {
		var ret time.Time
		return ret
	}
	return *o.EndDate
}

// GetEndDateOk returns a tuple with the EndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetEndDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.EndDate) {
		return nil, false
	}

	return o.EndDate, true
}

// HasEndDate returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasEndDate() bool {
	if o != nil && !IsNil(o.EndDate) {
		return true
	}

	return false
}

// SetEndDate gets a reference to the given time.Time and assigns it to the EndDate field.
func (o *InvoiceLineItem) SetEndDate(v time.Time) {
	o.EndDate = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *InvoiceLineItem) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *InvoiceLineItem) SetGroupId(v string) {
	o.GroupId = &v
}

// GetGroupName returns the GroupName field value if set, zero value otherwise
func (o *InvoiceLineItem) GetGroupName() string {
	if o == nil || IsNil(o.GroupName) {
		var ret string
		return ret
	}
	return *o.GroupName
}

// GetGroupNameOk returns a tuple with the GroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.GroupName) {
		return nil, false
	}

	return o.GroupName, true
}

// HasGroupName returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasGroupName() bool {
	if o != nil && !IsNil(o.GroupName) {
		return true
	}

	return false
}

// SetGroupName gets a reference to the given string and assigns it to the GroupName field.
func (o *InvoiceLineItem) SetGroupName(v string) {
	o.GroupName = &v
}

// GetNote returns the Note field value if set, zero value otherwise
func (o *InvoiceLineItem) GetNote() string {
	if o == nil || IsNil(o.Note) {
		var ret string
		return ret
	}
	return *o.Note
}

// GetNoteOk returns a tuple with the Note field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetNoteOk() (*string, bool) {
	if o == nil || IsNil(o.Note) {
		return nil, false
	}

	return o.Note, true
}

// HasNote returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasNote() bool {
	if o != nil && !IsNil(o.Note) {
		return true
	}

	return false
}

// SetNote gets a reference to the given string and assigns it to the Note field.
func (o *InvoiceLineItem) SetNote(v string) {
	o.Note = &v
}

// GetPercentDiscount returns the PercentDiscount field value if set, zero value otherwise
func (o *InvoiceLineItem) GetPercentDiscount() float32 {
	if o == nil || IsNil(o.PercentDiscount) {
		var ret float32
		return ret
	}
	return *o.PercentDiscount
}

// GetPercentDiscountOk returns a tuple with the PercentDiscount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetPercentDiscountOk() (*float32, bool) {
	if o == nil || IsNil(o.PercentDiscount) {
		return nil, false
	}

	return o.PercentDiscount, true
}

// HasPercentDiscount returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasPercentDiscount() bool {
	if o != nil && !IsNil(o.PercentDiscount) {
		return true
	}

	return false
}

// SetPercentDiscount gets a reference to the given float32 and assigns it to the PercentDiscount field.
func (o *InvoiceLineItem) SetPercentDiscount(v float32) {
	o.PercentDiscount = &v
}

// GetQuantity returns the Quantity field value if set, zero value otherwise
func (o *InvoiceLineItem) GetQuantity() float64 {
	if o == nil || IsNil(o.Quantity) {
		var ret float64
		return ret
	}
	return *o.Quantity
}

// GetQuantityOk returns a tuple with the Quantity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetQuantityOk() (*float64, bool) {
	if o == nil || IsNil(o.Quantity) {
		return nil, false
	}

	return o.Quantity, true
}

// HasQuantity returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasQuantity() bool {
	if o != nil && !IsNil(o.Quantity) {
		return true
	}

	return false
}

// SetQuantity gets a reference to the given float64 and assigns it to the Quantity field.
func (o *InvoiceLineItem) SetQuantity(v float64) {
	o.Quantity = &v
}

// GetSku returns the Sku field value if set, zero value otherwise
func (o *InvoiceLineItem) GetSku() string {
	if o == nil || IsNil(o.Sku) {
		var ret string
		return ret
	}
	return *o.Sku
}

// GetSkuOk returns a tuple with the Sku field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetSkuOk() (*string, bool) {
	if o == nil || IsNil(o.Sku) {
		return nil, false
	}

	return o.Sku, true
}

// HasSku returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasSku() bool {
	if o != nil && !IsNil(o.Sku) {
		return true
	}

	return false
}

// SetSku gets a reference to the given string and assigns it to the Sku field.
func (o *InvoiceLineItem) SetSku(v string) {
	o.Sku = &v
}

// GetStartDate returns the StartDate field value if set, zero value otherwise
func (o *InvoiceLineItem) GetStartDate() time.Time {
	if o == nil || IsNil(o.StartDate) {
		var ret time.Time
		return ret
	}
	return *o.StartDate
}

// GetStartDateOk returns a tuple with the StartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetStartDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartDate) {
		return nil, false
	}

	return o.StartDate, true
}

// HasStartDate returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasStartDate() bool {
	if o != nil && !IsNil(o.StartDate) {
		return true
	}

	return false
}

// SetStartDate gets a reference to the given time.Time and assigns it to the StartDate field.
func (o *InvoiceLineItem) SetStartDate(v time.Time) {
	o.StartDate = &v
}

// GetStitchAppName returns the StitchAppName field value if set, zero value otherwise
func (o *InvoiceLineItem) GetStitchAppName() string {
	if o == nil || IsNil(o.StitchAppName) {
		var ret string
		return ret
	}
	return *o.StitchAppName
}

// GetStitchAppNameOk returns a tuple with the StitchAppName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetStitchAppNameOk() (*string, bool) {
	if o == nil || IsNil(o.StitchAppName) {
		return nil, false
	}

	return o.StitchAppName, true
}

// HasStitchAppName returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasStitchAppName() bool {
	if o != nil && !IsNil(o.StitchAppName) {
		return true
	}

	return false
}

// SetStitchAppName gets a reference to the given string and assigns it to the StitchAppName field.
func (o *InvoiceLineItem) SetStitchAppName(v string) {
	o.StitchAppName = &v
}

// GetTags returns the Tags field value if set, zero value otherwise
func (o *InvoiceLineItem) GetTags() map[string][]string {
	if o == nil || IsNil(o.Tags) {
		var ret map[string][]string
		return ret
	}
	return *o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetTagsOk() (*map[string][]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}

	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given map[string][]string and assigns it to the Tags field.
func (o *InvoiceLineItem) SetTags(v map[string][]string) {
	o.Tags = &v
}

// GetTierLowerBound returns the TierLowerBound field value if set, zero value otherwise
func (o *InvoiceLineItem) GetTierLowerBound() float64 {
	if o == nil || IsNil(o.TierLowerBound) {
		var ret float64
		return ret
	}
	return *o.TierLowerBound
}

// GetTierLowerBoundOk returns a tuple with the TierLowerBound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetTierLowerBoundOk() (*float64, bool) {
	if o == nil || IsNil(o.TierLowerBound) {
		return nil, false
	}

	return o.TierLowerBound, true
}

// HasTierLowerBound returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasTierLowerBound() bool {
	if o != nil && !IsNil(o.TierLowerBound) {
		return true
	}

	return false
}

// SetTierLowerBound gets a reference to the given float64 and assigns it to the TierLowerBound field.
func (o *InvoiceLineItem) SetTierLowerBound(v float64) {
	o.TierLowerBound = &v
}

// GetTierUpperBound returns the TierUpperBound field value if set, zero value otherwise
func (o *InvoiceLineItem) GetTierUpperBound() float64 {
	if o == nil || IsNil(o.TierUpperBound) {
		var ret float64
		return ret
	}
	return *o.TierUpperBound
}

// GetTierUpperBoundOk returns a tuple with the TierUpperBound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetTierUpperBoundOk() (*float64, bool) {
	if o == nil || IsNil(o.TierUpperBound) {
		return nil, false
	}

	return o.TierUpperBound, true
}

// HasTierUpperBound returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasTierUpperBound() bool {
	if o != nil && !IsNil(o.TierUpperBound) {
		return true
	}

	return false
}

// SetTierUpperBound gets a reference to the given float64 and assigns it to the TierUpperBound field.
func (o *InvoiceLineItem) SetTierUpperBound(v float64) {
	o.TierUpperBound = &v
}

// GetTotalPriceCents returns the TotalPriceCents field value if set, zero value otherwise
func (o *InvoiceLineItem) GetTotalPriceCents() int64 {
	if o == nil || IsNil(o.TotalPriceCents) {
		var ret int64
		return ret
	}
	return *o.TotalPriceCents
}

// GetTotalPriceCentsOk returns a tuple with the TotalPriceCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetTotalPriceCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.TotalPriceCents) {
		return nil, false
	}

	return o.TotalPriceCents, true
}

// HasTotalPriceCents returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasTotalPriceCents() bool {
	if o != nil && !IsNil(o.TotalPriceCents) {
		return true
	}

	return false
}

// SetTotalPriceCents gets a reference to the given int64 and assigns it to the TotalPriceCents field.
func (o *InvoiceLineItem) SetTotalPriceCents(v int64) {
	o.TotalPriceCents = &v
}

// GetUnit returns the Unit field value if set, zero value otherwise
func (o *InvoiceLineItem) GetUnit() string {
	if o == nil || IsNil(o.Unit) {
		var ret string
		return ret
	}
	return *o.Unit
}

// GetUnitOk returns a tuple with the Unit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetUnitOk() (*string, bool) {
	if o == nil || IsNil(o.Unit) {
		return nil, false
	}

	return o.Unit, true
}

// HasUnit returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasUnit() bool {
	if o != nil && !IsNil(o.Unit) {
		return true
	}

	return false
}

// SetUnit gets a reference to the given string and assigns it to the Unit field.
func (o *InvoiceLineItem) SetUnit(v string) {
	o.Unit = &v
}

// GetUnitPriceDollars returns the UnitPriceDollars field value if set, zero value otherwise
func (o *InvoiceLineItem) GetUnitPriceDollars() float64 {
	if o == nil || IsNil(o.UnitPriceDollars) {
		var ret float64
		return ret
	}
	return *o.UnitPriceDollars
}

// GetUnitPriceDollarsOk returns a tuple with the UnitPriceDollars field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InvoiceLineItem) GetUnitPriceDollarsOk() (*float64, bool) {
	if o == nil || IsNil(o.UnitPriceDollars) {
		return nil, false
	}

	return o.UnitPriceDollars, true
}

// HasUnitPriceDollars returns a boolean if a field has been set.
func (o *InvoiceLineItem) HasUnitPriceDollars() bool {
	if o != nil && !IsNil(o.UnitPriceDollars) {
		return true
	}

	return false
}

// SetUnitPriceDollars gets a reference to the given float64 and assigns it to the UnitPriceDollars field.
func (o *InvoiceLineItem) SetUnitPriceDollars(v float64) {
	o.UnitPriceDollars = &v
}
