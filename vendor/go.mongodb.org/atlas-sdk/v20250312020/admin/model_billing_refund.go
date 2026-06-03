// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BillingRefund One payment that MongoDB returned to the organization for this invoice.
type BillingRefund struct {
	// Sum of the funds returned to the specified organization expressed in cents (100th of US Dollar).
	// Read only field.
	AmountCents *int64 `json:"amountCents,omitempty"`
	// Date and time when MongoDB Cloud created this refund. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the payment that the organization had made.
	// Read only field.
	PaymentId *string `json:"paymentId,omitempty"`
	// Justification that MongoDB accepted to return funds to the organization.
	// Read only field.
	Reason *string `json:"reason,omitempty"`
}

// NewBillingRefund instantiates a new BillingRefund object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBillingRefund() *BillingRefund {
	this := BillingRefund{}
	return &this
}

// NewBillingRefundWithDefaults instantiates a new BillingRefund object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBillingRefundWithDefaults() *BillingRefund {
	this := BillingRefund{}
	return &this
}

// GetAmountCents returns the AmountCents field value if set, zero value otherwise
func (o *BillingRefund) GetAmountCents() int64 {
	if o == nil || IsNil(o.AmountCents) {
		var ret int64
		return ret
	}
	return *o.AmountCents
}

// GetAmountCentsOk returns a tuple with the AmountCents field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingRefund) GetAmountCentsOk() (*int64, bool) {
	if o == nil || IsNil(o.AmountCents) {
		return nil, false
	}

	return o.AmountCents, true
}

// HasAmountCents returns a boolean if a field has been set.
func (o *BillingRefund) HasAmountCents() bool {
	if o != nil && !IsNil(o.AmountCents) {
		return true
	}

	return false
}

// SetAmountCents gets a reference to the given int64 and assigns it to the AmountCents field.
func (o *BillingRefund) SetAmountCents(v int64) {
	o.AmountCents = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *BillingRefund) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingRefund) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *BillingRefund) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *BillingRefund) SetCreated(v time.Time) {
	o.Created = &v
}

// GetPaymentId returns the PaymentId field value if set, zero value otherwise
func (o *BillingRefund) GetPaymentId() string {
	if o == nil || IsNil(o.PaymentId) {
		var ret string
		return ret
	}
	return *o.PaymentId
}

// GetPaymentIdOk returns a tuple with the PaymentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingRefund) GetPaymentIdOk() (*string, bool) {
	if o == nil || IsNil(o.PaymentId) {
		return nil, false
	}

	return o.PaymentId, true
}

// HasPaymentId returns a boolean if a field has been set.
func (o *BillingRefund) HasPaymentId() bool {
	if o != nil && !IsNil(o.PaymentId) {
		return true
	}

	return false
}

// SetPaymentId gets a reference to the given string and assigns it to the PaymentId field.
func (o *BillingRefund) SetPaymentId(v string) {
	o.PaymentId = &v
}

// GetReason returns the Reason field value if set, zero value otherwise
func (o *BillingRefund) GetReason() string {
	if o == nil || IsNil(o.Reason) {
		var ret string
		return ret
	}
	return *o.Reason
}

// GetReasonOk returns a tuple with the Reason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BillingRefund) GetReasonOk() (*string, bool) {
	if o == nil || IsNil(o.Reason) {
		return nil, false
	}

	return o.Reason, true
}

// HasReason returns a boolean if a field has been set.
func (o *BillingRefund) HasReason() bool {
	if o != nil && !IsNil(o.Reason) {
		return true
	}

	return false
}

// SetReason gets a reference to the given string and assigns it to the Reason field.
func (o *BillingRefund) SetReason(v string) {
	o.Reason = &v
}
