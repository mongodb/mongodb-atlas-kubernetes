// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BillingInvoice struct for BillingInvoice
type BillingInvoice struct {
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
	// List that contains individual services included in this invoice.
	// Read only field.
	LineItems *[]InvoiceLineItem `json:"lineItems,omitempty"`
	// List that contains the invoices for organizations linked to the paying organization.
	// Read only field.
	LinkedInvoices *[]BillingInvoice `json:"linkedInvoices,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization charged for services consumed from MongoDB Cloud.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// List that contains funds transferred to MongoDB to cover the specified service noted in this invoice.
	// Read only field.
	Payments *[]BillingPayment `json:"payments,omitempty"`
	// List that contains payments that MongoDB returned to the organization for this invoice.
	// Read only field.
	Refunds *[]BillingRefund `json:"refunds,omitempty"`
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

// NewBillingInvoice instantiates a new BillingInvoice object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBillingInvoice() *BillingInvoice {
	this := BillingInvoice{}
	return &this
}

// NewBillingInvoiceWithDefaults instantiates a new BillingInvoice object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBillingInvoiceWithDefaults() *BillingInvoice {
	this := BillingInvoice{}
	return &this
}

// GetAmountBilledCents returns the AmountBilledCents field value if set, zero value otherwise
func (o *BillingInvoice) GetAmountBilledCents() int64 {
	if o == nil || IsNil(o.AmountBilledCents) {
		var ret int64
		return ret
	}
	return *o.AmountBilledCents
}

// GetAmountBilledCentsOk returns a tuple with the AmountBilledCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetAmountBilledCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountBilledCents) {
		return nil, false
	}

	return o.AmountBilledCents, true
}

// HasAmountBilledCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasAmountBilledCents() bool {
	if o != nil && !IsNil(o.AmountBilledCents) {
		return true
	}

	return false
}

// SetAmountBilledCents gets a reference to the given int64 and assigns it to the AmountBilledCents field.
func (o *BillingInvoice) SetAmountBilledCents(v int64) {
	o.AmountBilledCents = &v
}

// GetAmountPaidCents returns the AmountPaidCents field value if set, zero value otherwise
func (o *BillingInvoice) GetAmountPaidCents() int64 {
	if o == nil || IsNil(o.AmountPaidCents) {
		var ret int64
		return ret
	}
	return *o.AmountPaidCents
}

// GetAmountPaidCentsOk returns a tuple with the AmountPaidCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetAmountPaidCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountPaidCents) {
		return nil, false
	}

	return o.AmountPaidCents, true
}

// HasAmountPaidCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasAmountPaidCents() bool {
	if o != nil && !IsNil(o.AmountPaidCents) {
		return true
	}

	return false
}

// SetAmountPaidCents gets a reference to the given int64 and assigns it to the AmountPaidCents field.
func (o *BillingInvoice) SetAmountPaidCents(v int64) {
	o.AmountPaidCents = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BillingInvoice) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BillingInvoice) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *BillingInvoice) SetCreated(v time.Time) {
	o.Created = &v
}

// GetCreditsCents returns the CreditsCents field value if set, zero value otherwise
func (o *BillingInvoice) GetCreditsCents() int64 {
	if o == nil || IsNil(o.CreditsCents) {
		var ret int64
		return ret
	}
	return *o.CreditsCents
}

// GetCreditsCentsOk returns a tuple with the CreditsCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetCreditsCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.CreditsCents) {
		return nil, false
	}

	return o.CreditsCents, true
}

// HasCreditsCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasCreditsCents() bool {
	if o != nil && !IsNil(o.CreditsCents) {
		return true
	}

	return false
}

// SetCreditsCents gets a reference to the given int64 and assigns it to the CreditsCents field.
func (o *BillingInvoice) SetCreditsCents(v int64) {
	o.CreditsCents = &v
}

// GetEndDate returns the EndDate field value if set, zero value otherwise
func (o *BillingInvoice) GetEndDate() time.Time {
	if o == nil || IsNil(o.EndDate) {
		var ret time.Time
		return ret
	}
	return *o.EndDate
}

// GetEndDateOk returns a tuple with the EndDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetEndDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.EndDate) {
		return nil, false
	}

	return o.EndDate, true
}

// HasEndDate returns a boolean if a field has been set.
func (o *BillingInvoice) HasEndDate() bool {
	if o != nil && !IsNil(o.EndDate) {
		return true
	}

	return false
}

// SetEndDate gets a reference to the given time.Time and assigns it to the EndDate field.
func (o *BillingInvoice) SetEndDate(v time.Time) {
	o.EndDate = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BillingInvoice) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BillingInvoice) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BillingInvoice) SetId(v string) {
	o.Id = &v
}

// GetLineItems returns the LineItems field value if set, zero value otherwise
func (o *BillingInvoice) GetLineItems() []InvoiceLineItem {
	if o == nil || IsNil(o.LineItems) {
		var ret []InvoiceLineItem
		return ret
	}
	return *o.LineItems
}

// GetLineItemsOk returns a tuple with the LineItems field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetLineItemsOk() (*[]InvoiceLineItem, bool) {
	if o == nil || IsNil(o.LineItems) {
		return nil, false
	}

	return o.LineItems, true
}

// HasLineItems returns a boolean if a field has been set.
func (o *BillingInvoice) HasLineItems() bool {
	if o != nil && !IsNil(o.LineItems) {
		return true
	}

	return false
}

// SetLineItems gets a reference to the given []InvoiceLineItem and assigns it to the LineItems field.
func (o *BillingInvoice) SetLineItems(v []InvoiceLineItem) {
	o.LineItems = &v
}

// GetLinkedInvoices returns the LinkedInvoices field value if set, zero value otherwise
func (o *BillingInvoice) GetLinkedInvoices() []BillingInvoice {
	if o == nil || IsNil(o.LinkedInvoices) {
		var ret []BillingInvoice
		return ret
	}
	return *o.LinkedInvoices
}

// GetLinkedInvoicesOk returns a tuple with the LinkedInvoices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetLinkedInvoicesOk() (*[]BillingInvoice, bool) {
	if o == nil || IsNil(o.LinkedInvoices) {
		return nil, false
	}

	return o.LinkedInvoices, true
}

// HasLinkedInvoices returns a boolean if a field has been set.
func (o *BillingInvoice) HasLinkedInvoices() bool {
	if o != nil && !IsNil(o.LinkedInvoices) {
		return true
	}

	return false
}

// SetLinkedInvoices gets a reference to the given []BillingInvoice and assigns it to the LinkedInvoices field.
func (o *BillingInvoice) SetLinkedInvoices(v []BillingInvoice) {
	o.LinkedInvoices = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *BillingInvoice) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *BillingInvoice) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *BillingInvoice) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *BillingInvoice) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *BillingInvoice) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *BillingInvoice) SetOrgId(v string) {
	o.OrgId = &v
}

// GetPayments returns the Payments field value if set, zero value otherwise
func (o *BillingInvoice) GetPayments() []BillingPayment {
	if o == nil || IsNil(o.Payments) {
		var ret []BillingPayment
		return ret
	}
	return *o.Payments
}

// GetPaymentsOk returns a tuple with the Payments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetPaymentsOk() (*[]BillingPayment, bool) {
	if o == nil || IsNil(o.Payments) {
		return nil, false
	}

	return o.Payments, true
}

// HasPayments returns a boolean if a field has been set.
func (o *BillingInvoice) HasPayments() bool {
	if o != nil && !IsNil(o.Payments) {
		return true
	}

	return false
}

// SetPayments gets a reference to the given []BillingPayment and assigns it to the Payments field.
func (o *BillingInvoice) SetPayments(v []BillingPayment) {
	o.Payments = &v
}

// GetRefunds returns the Refunds field value if set, zero value otherwise
func (o *BillingInvoice) GetRefunds() []BillingRefund {
	if o == nil || IsNil(o.Refunds) {
		var ret []BillingRefund
		return ret
	}
	return *o.Refunds
}

// GetRefundsOk returns a tuple with the Refunds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetRefundsOk() (*[]BillingRefund, bool) {
	if o == nil || IsNil(o.Refunds) {
		return nil, false
	}

	return o.Refunds, true
}

// HasRefunds returns a boolean if a field has been set.
func (o *BillingInvoice) HasRefunds() bool {
	if o != nil && !IsNil(o.Refunds) {
		return true
	}

	return false
}

// SetRefunds gets a reference to the given []BillingRefund and assigns it to the Refunds field.
func (o *BillingInvoice) SetRefunds(v []BillingRefund) {
	o.Refunds = &v
}

// GetSalesTaxCents returns the SalesTaxCents field value if set, zero value otherwise
func (o *BillingInvoice) GetSalesTaxCents() int64 {
	if o == nil || IsNil(o.SalesTaxCents) {
		var ret int64
		return ret
	}
	return *o.SalesTaxCents
}

// GetSalesTaxCentsOk returns a tuple with the SalesTaxCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetSalesTaxCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SalesTaxCents) {
		return nil, false
	}

	return o.SalesTaxCents, true
}

// HasSalesTaxCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasSalesTaxCents() bool {
	if o != nil && !IsNil(o.SalesTaxCents) {
		return true
	}

	return false
}

// SetSalesTaxCents gets a reference to the given int64 and assigns it to the SalesTaxCents field.
func (o *BillingInvoice) SetSalesTaxCents(v int64) {
	o.SalesTaxCents = &v
}

// GetStartDate returns the StartDate field value if set, zero value otherwise
func (o *BillingInvoice) GetStartDate() time.Time {
	if o == nil || IsNil(o.StartDate) {
		var ret time.Time
		return ret
	}
	return *o.StartDate
}

// GetStartDateOk returns a tuple with the StartDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetStartDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartDate) {
		return nil, false
	}

	return o.StartDate, true
}

// HasStartDate returns a boolean if a field has been set.
func (o *BillingInvoice) HasStartDate() bool {
	if o != nil && !IsNil(o.StartDate) {
		return true
	}

	return false
}

// SetStartDate gets a reference to the given time.Time and assigns it to the StartDate field.
func (o *BillingInvoice) SetStartDate(v time.Time) {
	o.StartDate = &v
}

// GetStartingBalanceCents returns the StartingBalanceCents field value if set, zero value otherwise
func (o *BillingInvoice) GetStartingBalanceCents() int64 {
	if o == nil || IsNil(o.StartingBalanceCents) {
		var ret int64
		return ret
	}
	return *o.StartingBalanceCents
}

// GetStartingBalanceCentsOk returns a tuple with the StartingBalanceCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetStartingBalanceCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.StartingBalanceCents) {
		return nil, false
	}

	return o.StartingBalanceCents, true
}

// HasStartingBalanceCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasStartingBalanceCents() bool {
	if o != nil && !IsNil(o.StartingBalanceCents) {
		return true
	}

	return false
}

// SetStartingBalanceCents gets a reference to the given int64 and assigns it to the StartingBalanceCents field.
func (o *BillingInvoice) SetStartingBalanceCents(v int64) {
	o.StartingBalanceCents = &v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BillingInvoice) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BillingInvoice) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BillingInvoice) SetStatusName(v string) {
	o.StatusName = &v
}

// GetSubtotalCents returns the SubtotalCents field value if set, zero value otherwise
func (o *BillingInvoice) GetSubtotalCents() int64 {
	if o == nil || IsNil(o.SubtotalCents) {
		var ret int64
		return ret
	}
	return *o.SubtotalCents
}

// GetSubtotalCentsOk returns a tuple with the SubtotalCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetSubtotalCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SubtotalCents) {
		return nil, false
	}

	return o.SubtotalCents, true
}

// HasSubtotalCents returns a boolean if a field has been set.
func (o *BillingInvoice) HasSubtotalCents() bool {
	if o != nil && !IsNil(o.SubtotalCents) {
		return true
	}

	return false
}

// SetSubtotalCents gets a reference to the given int64 and assigns it to the SubtotalCents field.
func (o *BillingInvoice) SetSubtotalCents(v int64) {
	o.SubtotalCents = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise
func (o *BillingInvoice) GetUpdated() time.Time {
	if o == nil || IsNil(o.Updated) {
		var ret time.Time
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingInvoice) GetUpdatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}

	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *BillingInvoice) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given time.Time and assigns it to the Updated field.
func (o *BillingInvoice) SetUpdated(v time.Time) {
	o.Updated = &v
}
