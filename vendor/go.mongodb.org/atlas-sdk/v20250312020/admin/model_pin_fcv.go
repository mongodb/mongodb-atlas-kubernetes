// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// PinFCV struct for PinFCV
type PinFCV struct {
	// Expiration date of the fixed FCV. If not specified, the expiration date will default to 4 weeks from the date FCV was originally pinned. Note that this field cannot exceed 4 weeks from the pinned date. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewPinFCV instantiates a new PinFCV object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPinFCV() *PinFCV {
	this := PinFCV{}
	return &this
}

// NewPinFCVWithDefaults instantiates a new PinFCV object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPinFCVWithDefaults() *PinFCV {
	this := PinFCV{}
	return &this
}

// GetExpirationDate returns the ExpirationDate field value if set, zero value otherwise
func (o *PinFCV) GetExpirationDate() time.Time {
	if o == nil || IsNil(o.ExpirationDate) {
		var ret time.Time
		return ret
	}
	return *o.ExpirationDate
}

// GetExpirationDateOk returns a tuple with the ExpirationDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinFCV) GetExpirationDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpirationDate) {
		return nil, false
	}

	return o.ExpirationDate, true
}

// HasExpirationDate returns a boolean if a field has been set.
func (o *PinFCV) HasExpirationDate() bool {
	if o != nil && !IsNil(o.ExpirationDate) {
		return true
	}

	return false
}

// SetExpirationDate gets a reference to the given time.Time and assigns it to the ExpirationDate field.
func (o *PinFCV) SetExpirationDate(v time.Time) {
	o.ExpirationDate = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *PinFCV) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinFCV) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *PinFCV) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *PinFCV) SetLinks(v []Link) {
	o.Links = &v
}
