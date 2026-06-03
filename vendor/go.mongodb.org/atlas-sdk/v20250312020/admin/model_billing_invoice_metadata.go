// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BillingInvoiceMetadata struct for BillingInvoiceMetadata
type BillingInvoiceMetadata struct {
	// Sum of services that the specified organization consumed in the period covered in this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	AmountBilledCents *int64 `json:"amountBilledCents,omitempty"`
	// Sum that the specified organization paid toward this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	AmountPaidCents *int64 `json:"amountPaidCents,omitempty"`
	// Date and time when MongoDB Cloud created this invoice. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Sum that MongoDB credited the specified organization toward this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	CreditsCents *int64 `json:"creditsCents,omitempty"`
	// Date and time when MongoDB Cloud finished the billing period that this invoice covers. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	EndDate *time.Time `json:"endDate,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List that contains the invoices for organizations linked to the paying organization.
	// Read only field.
	LinkedInvoices *[]BillingInvoiceMetadata `json:"linkedInvoices,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization charged for services consumed from MongoDB Cloud.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Sum of sales tax applied to this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	SalesTaxCents *int64 `json:"salesTaxCents,omitempty"`
	// Date and time when MongoDB Cloud began the billing period that this invoice covers. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	StartDate *time.Time `json:"startDate,omitempty"`
	// Sum that the specified organization owed to MongoDB when MongoDB issued this invoice. This parameter expresses its value in US Dollars.
	// Read only field.
	StartingBalanceCents *int64 `json:"startingBalanceCents,omitempty"`
	// Phase of payment processing in which this invoice exists when you made this request. Accepted phases include:  - `CLOSED`: MongoDB finalized all charges in the billing cycle but has yet to charge the customer. - `FAILED`: MongoDB attempted to charge the provided credit card but charge for that amount failed. - `FORGIVEN`: Customer initiated payment which MongoDB later forgave. - `FREE`: All charges totalled zero so the customer won't be charged. - `INVOICED`: MongoDB handled these charges using elastic invoicing. - `PAID`: MongoDB succeeded in charging the provided credit card. - `PENDING`: Invoice includes charges for the current billing cycle. - `PREPAID`: Customer has a pre-paid plan so they won't be charged.
	StatusName *string `json:"statusName,omitempty"`
	// Sum of all positive invoice line items contained in this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	SubtotalCents *int64 `json:"subtotalCents,omitempty"`
	// Date and time when MongoDB Cloud last updated the value of this payment. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Updated *time.Time `json:"updated,omitempty"`
}

// NewBillingInvoiceMetadata instantiates a new BillingInvoiceMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBillingInvoiceMetadata() *BillingInvoiceMetadata {
	this := BillingInvoiceMetadata{}
	return &this
}

// NewBillingInvoiceMetadataWithDefaults instantiates a new BillingInvoiceMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBillingInvoiceMetadataWithDefaults() *BillingInvoiceMetadata {
	this := BillingInvoiceMetadata{}
	return &this
}

// GetAmountBilledCents returns the AmountBilledCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetAmountBilledCents() int64 {
	if o == nil || IsNil(o.AmountBilledCents) {
		var ret int64
		return ret
	}
	return *o.AmountBilledCents
}

// GetAmountBilledCentsOk returns a tuple with the AmountBilledCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetAmountBilledCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountBilledCents) {
		return nil, false
	}

	return o.AmountBilledCents, true
}

// HasAmountBilledCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasAmountBilledCents() bool {
	if o != nil && !IsNil(o.AmountBilledCents) {
		return true
	}

	return false
}

// SetAmountBilledCents gets a reference to the given int64 and assigns it to the AmountBilledCents field.
func (o *BillingInvoiceMetadata) SetAmountBilledCents(v int64) {
	o.AmountBilledCents = &v
}

// GetAmountPaidCents returns the AmountPaidCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetAmountPaidCents() int64 {
	if o == nil || IsNil(o.AmountPaidCents) {
		var ret int64
		return ret
	}
	return *o.AmountPaidCents
}

// GetAmountPaidCentsOk returns a tuple with the AmountPaidCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetAmountPaidCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountPaidCents) {
		return nil, false
	}

	return o.AmountPaidCents, true
}

// HasAmountPaidCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasAmountPaidCents() bool {
	if o != nil && !IsNil(o.AmountPaidCents) {
		return true
	}

	return false
}

// SetAmountPaidCents gets a reference to the given int64 and assigns it to the AmountPaidCents field.
func (o *BillingInvoiceMetadata) SetAmountPaidCents(v int64) {
	o.AmountPaidCents = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *BillingInvoiceMetadata) SetCreated(v time.Time) {
	o.Created = &v
}

// GetCreditsCents returns the CreditsCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetCreditsCents() int64 {
	if o == nil || IsNil(o.CreditsCents) {
		var ret int64
		return ret
	}
	return *o.CreditsCents
}

// GetCreditsCentsOk returns a tuple with the CreditsCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetCreditsCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.CreditsCents) {
		return nil, false
	}

	return o.CreditsCents, true
}

// HasCreditsCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasCreditsCents() bool {
	if o != nil && !IsNil(o.CreditsCents) {
		return true
	}

	return false
}

// SetCreditsCents gets a reference to the given int64 and assigns it to the CreditsCents field.
func (o *BillingInvoiceMetadata) SetCreditsCents(v int64) {
	o.CreditsCents = &v
}

// GetEndDate returns the EndDate field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetEndDate() time.Time {
	if o == nil || IsNil(o.EndDate) {
		var ret time.Time
		return ret
	}
	return *o.EndDate
}

// GetEndDateOk returns a tuple with the EndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetEndDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.EndDate) {
		return nil, false
	}

	return o.EndDate, true
}

// HasEndDate returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasEndDate() bool {
	if o != nil && !IsNil(o.EndDate) {
		return true
	}

	return false
}

// SetEndDate gets a reference to the given time.Time and assigns it to the EndDate field.
func (o *BillingInvoiceMetadata) SetEndDate(v time.Time) {
	o.EndDate = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BillingInvoiceMetadata) SetId(v string) {
	o.Id = &v
}

// GetLinkedInvoices returns the LinkedInvoices field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetLinkedInvoices() []BillingInvoiceMetadata {
	if o == nil || IsNil(o.LinkedInvoices) {
		var ret []BillingInvoiceMetadata
		return ret
	}
	return *o.LinkedInvoices
}

// GetLinkedInvoicesOk returns a tuple with the LinkedInvoices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetLinkedInvoicesOk() (*[]BillingInvoiceMetadata, bool) {
	if o == nil || IsNil(o.LinkedInvoices) {
		return nil, false
	}

	return o.LinkedInvoices, true
}

// HasLinkedInvoices returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasLinkedInvoices() bool {
	if o != nil && !IsNil(o.LinkedInvoices) {
		return true
	}

	return false
}

// SetLinkedInvoices gets a reference to the given []BillingInvoiceMetadata and assigns it to the LinkedInvoices field.
func (o *BillingInvoiceMetadata) SetLinkedInvoices(v []BillingInvoiceMetadata) {
	o.LinkedInvoices = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *BillingInvoiceMetadata) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *BillingInvoiceMetadata) SetOrgId(v string) {
	o.OrgId = &v
}

// GetSalesTaxCents returns the SalesTaxCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetSalesTaxCents() int64 {
	if o == nil || IsNil(o.SalesTaxCents) {
		var ret int64
		return ret
	}
	return *o.SalesTaxCents
}

// GetSalesTaxCentsOk returns a tuple with the SalesTaxCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetSalesTaxCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SalesTaxCents) {
		return nil, false
	}

	return o.SalesTaxCents, true
}

// HasSalesTaxCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasSalesTaxCents() bool {
	if o != nil && !IsNil(o.SalesTaxCents) {
		return true
	}

	return false
}

// SetSalesTaxCents gets a reference to the given int64 and assigns it to the SalesTaxCents field.
func (o *BillingInvoiceMetadata) SetSalesTaxCents(v int64) {
	o.SalesTaxCents = &v
}

// GetStartDate returns the StartDate field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetStartDate() time.Time {
	if o == nil || IsNil(o.StartDate) {
		var ret time.Time
		return ret
	}
	return *o.StartDate
}

// GetStartDateOk returns a tuple with the StartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetStartDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartDate) {
		return nil, false
	}

	return o.StartDate, true
}

// HasStartDate returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasStartDate() bool {
	if o != nil && !IsNil(o.StartDate) {
		return true
	}

	return false
}

// SetStartDate gets a reference to the given time.Time and assigns it to the StartDate field.
func (o *BillingInvoiceMetadata) SetStartDate(v time.Time) {
	o.StartDate = &v
}

// GetStartingBalanceCents returns the StartingBalanceCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetStartingBalanceCents() int64 {
	if o == nil || IsNil(o.StartingBalanceCents) {
		var ret int64
		return ret
	}
	return *o.StartingBalanceCents
}

// GetStartingBalanceCentsOk returns a tuple with the StartingBalanceCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetStartingBalanceCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.StartingBalanceCents) {
		return nil, false
	}

	return o.StartingBalanceCents, true
}

// HasStartingBalanceCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasStartingBalanceCents() bool {
	if o != nil && !IsNil(o.StartingBalanceCents) {
		return true
	}

	return false
}

// SetStartingBalanceCents gets a reference to the given int64 and assigns it to the StartingBalanceCents field.
func (o *BillingInvoiceMetadata) SetStartingBalanceCents(v int64) {
	o.StartingBalanceCents = &v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BillingInvoiceMetadata) SetStatusName(v string) {
	o.StatusName = &v
}

// GetSubtotalCents returns the SubtotalCents field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetSubtotalCents() int64 {
	if o == nil || IsNil(o.SubtotalCents) {
		var ret int64
		return ret
	}
	return *o.SubtotalCents
}

// GetSubtotalCentsOk returns a tuple with the SubtotalCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetSubtotalCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SubtotalCents) {
		return nil, false
	}

	return o.SubtotalCents, true
}

// HasSubtotalCents returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasSubtotalCents() bool {
	if o != nil && !IsNil(o.SubtotalCents) {
		return true
	}

	return false
}

// SetSubtotalCents gets a reference to the given int64 and assigns it to the SubtotalCents field.
func (o *BillingInvoiceMetadata) SetSubtotalCents(v int64) {
	o.SubtotalCents = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise
func (o *BillingInvoiceMetadata) GetUpdated() time.Time {
	if o == nil || IsNil(o.Updated) {
		var ret time.Time
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoiceMetadata) GetUpdatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}

	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *BillingInvoiceMetadata) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given time.Time and assigns it to the Updated field.
func (o *BillingInvoiceMetadata) SetUpdated(v time.Time) {
	o.Updated = &v
}
