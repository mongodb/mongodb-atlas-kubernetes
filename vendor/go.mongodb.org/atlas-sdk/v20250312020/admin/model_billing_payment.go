// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BillingPayment Funds transferred to MongoDB to cover the specified service in this invoice.
type BillingPayment struct {
	// Sum of services that the specified organization consumed in the period covered in this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	AmountBilledCents *int64 `json:"amountBilledCents,omitempty"`
	// Sum that the specified organization paid toward the associated invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	AmountPaidCents *int64 `json:"amountPaidCents,omitempty"`
	// Date and time when the customer made this payment attempt. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// The currency in which payment was paid. This parameter expresses its value in 3-letter ISO 4217 currency code.
	// Read only field.
	Currency *string `json:"currency,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this payment toward the associated invoice.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Sum of sales tax applied to this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	SalesTaxCents *int64 `json:"salesTaxCents,omitempty"`
	// Phase of payment processing for the associated invoice when you made this request. These phases include:  - `CANCELLED`: Customer or MongoDB cancelled the payment. - `ERROR`: Issue arose when attempting to complete payment. - `FAILED`: MongoDB tried to charge the credit card without success. - `FAILED_AUTHENTICATION`: Strong Customer Authentication has failed. Confirm that your payment method is authenticated. - `FORGIVEN`: Customer initiated payment which MongoDB later forgave. - `INVOICED`: MongoDB issued an invoice that included this line item. - `NEW`: Customer provided a method of payment, but MongoDB hasn't tried to charge the credit card. - `PAID`: Customer submitted a successful payment. - `PARTIAL_PAID`: Customer paid for part of this line item.
	StatusName *string `json:"statusName,omitempty"`
	// Sum of all positive invoice line items contained in this invoice. This parameter expresses its value in cents (100ths of one US Dollar).
	// Read only field.
	SubtotalCents *int64 `json:"subtotalCents,omitempty"`
	// The unit price applied to `amountBilledCents` to compute total payment amount. This value is represented as a decimal string.
	// Read only field.
	UnitPrice *string `json:"unitPrice,omitempty"`
	// Date and time when the customer made an update to this payment attempt. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Updated *time.Time `json:"updated,omitempty"`
}

// NewBillingPayment instantiates a new BillingPayment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBillingPayment() *BillingPayment {
	this := BillingPayment{}
	return &this
}

// NewBillingPaymentWithDefaults instantiates a new BillingPayment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBillingPaymentWithDefaults() *BillingPayment {
	this := BillingPayment{}
	return &this
}

// GetAmountBilledCents returns the AmountBilledCents field value if set, zero value otherwise
func (o *BillingPayment) GetAmountBilledCents() int64 {
	if o == nil || IsNil(o.AmountBilledCents) {
		var ret int64
		return ret
	}
	return *o.AmountBilledCents
}

// GetAmountBilledCentsOk returns a tuple with the AmountBilledCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetAmountBilledCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountBilledCents) {
		return nil, false
	}

	return o.AmountBilledCents, true
}

// HasAmountBilledCents returns a boolean if a field has been set.
func (o *BillingPayment) HasAmountBilledCents() bool {
	if o != nil && !IsNil(o.AmountBilledCents) {
		return true
	}

	return false
}

// SetAmountBilledCents gets a reference to the given int64 and assigns it to the AmountBilledCents field.
func (o *BillingPayment) SetAmountBilledCents(v int64) {
	o.AmountBilledCents = &v
}

// GetAmountPaidCents returns the AmountPaidCents field value if set, zero value otherwise
func (o *BillingPayment) GetAmountPaidCents() int64 {
	if o == nil || IsNil(o.AmountPaidCents) {
		var ret int64
		return ret
	}
	return *o.AmountPaidCents
}

// GetAmountPaidCentsOk returns a tuple with the AmountPaidCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetAmountPaidCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountPaidCents) {
		return nil, false
	}

	return o.AmountPaidCents, true
}

// HasAmountPaidCents returns a boolean if a field has been set.
func (o *BillingPayment) HasAmountPaidCents() bool {
	if o != nil && !IsNil(o.AmountPaidCents) {
		return true
	}

	return false
}

// SetAmountPaidCents gets a reference to the given int64 and assigns it to the AmountPaidCents field.
func (o *BillingPayment) SetAmountPaidCents(v int64) {
	o.AmountPaidCents = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BillingPayment) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BillingPayment) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *BillingPayment) SetCreated(v time.Time) {
	o.Created = &v
}

// GetCurrency returns the Currency field value if set, zero value otherwise
func (o *BillingPayment) GetCurrency() string {
	if o == nil || IsNil(o.Currency) {
		var ret string
		return ret
	}
	return *o.Currency
}

// GetCurrencyOk returns a tuple with the Currency field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetCurrencyOk() (*string, bool) {
	if o == nil || IsNil(o.Currency) {
		return nil, false
	}

	return o.Currency, true
}

// HasCurrency returns a boolean if a field has been set.
func (o *BillingPayment) HasCurrency() bool {
	if o != nil && !IsNil(o.Currency) {
		return true
	}

	return false
}

// SetCurrency gets a reference to the given string and assigns it to the Currency field.
func (o *BillingPayment) SetCurrency(v string) {
	o.Currency = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BillingPayment) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BillingPayment) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BillingPayment) SetId(v string) {
	o.Id = &v
}

// GetSalesTaxCents returns the SalesTaxCents field value if set, zero value otherwise
func (o *BillingPayment) GetSalesTaxCents() int64 {
	if o == nil || IsNil(o.SalesTaxCents) {
		var ret int64
		return ret
	}
	return *o.SalesTaxCents
}

// GetSalesTaxCentsOk returns a tuple with the SalesTaxCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetSalesTaxCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SalesTaxCents) {
		return nil, false
	}

	return o.SalesTaxCents, true
}

// HasSalesTaxCents returns a boolean if a field has been set.
func (o *BillingPayment) HasSalesTaxCents() bool {
	if o != nil && !IsNil(o.SalesTaxCents) {
		return true
	}

	return false
}

// SetSalesTaxCents gets a reference to the given int64 and assigns it to the SalesTaxCents field.
func (o *BillingPayment) SetSalesTaxCents(v int64) {
	o.SalesTaxCents = &v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BillingPayment) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BillingPayment) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BillingPayment) SetStatusName(v string) {
	o.StatusName = &v
}

// GetSubtotalCents returns the SubtotalCents field value if set, zero value otherwise
func (o *BillingPayment) GetSubtotalCents() int64 {
	if o == nil || IsNil(o.SubtotalCents) {
		var ret int64
		return ret
	}
	return *o.SubtotalCents
}

// GetSubtotalCentsOk returns a tuple with the SubtotalCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetSubtotalCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.SubtotalCents) {
		return nil, false
	}

	return o.SubtotalCents, true
}

// HasSubtotalCents returns a boolean if a field has been set.
func (o *BillingPayment) HasSubtotalCents() bool {
	if o != nil && !IsNil(o.SubtotalCents) {
		return true
	}

	return false
}

// SetSubtotalCents gets a reference to the given int64 and assigns it to the SubtotalCents field.
func (o *BillingPayment) SetSubtotalCents(v int64) {
	o.SubtotalCents = &v
}

// GetUnitPrice returns the UnitPrice field value if set, zero value otherwise
func (o *BillingPayment) GetUnitPrice() string {
	if o == nil || IsNil(o.UnitPrice) {
		var ret string
		return ret
	}
	return *o.UnitPrice
}

// GetUnitPriceOk returns a tuple with the UnitPrice field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetUnitPriceOk() (*string, bool) {
	if o == nil || IsNil(o.UnitPrice) {
		return nil, false
	}

	return o.UnitPrice, true
}

// HasUnitPrice returns a boolean if a field has been set.
func (o *BillingPayment) HasUnitPrice() bool {
	if o != nil && !IsNil(o.UnitPrice) {
		return true
	}

	return false
}

// SetUnitPrice gets a reference to the given string and assigns it to the UnitPrice field.
func (o *BillingPayment) SetUnitPrice(v string) {
	o.UnitPrice = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise
func (o *BillingPayment) GetUpdated() time.Time {
	if o == nil || IsNil(o.Updated) {
		var ret time.Time
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingPayment) GetUpdatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}

	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *BillingPayment) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given time.Time and assigns it to the Updated field.
func (o *BillingPayment) SetUpdated(v time.Time) {
	o.Updated = &v
}
