// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// X509Certificate struct for X509Certificate
type X509Certificate struct {
	// Latest date that the certificate is valid. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	NotAfter *time.Time `json:"notAfter,omitempty"`
	// Earliest date that the certificate is valid. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	NotBefore *time.Time `json:"notBefore,omitempty"`
}

// NewX509Certificate instantiates a new X509Certificate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewX509Certificate() *X509Certificate {
	this := X509Certificate{}
	return &this
}

// NewX509CertificateWithDefaults instantiates a new X509Certificate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewX509CertificateWithDefaults() *X509Certificate {
	this := X509Certificate{}
	return &this
}

// GetNotAfter returns the NotAfter field value if set, zero value otherwise
func (o *X509Certificate) GetNotAfter() time.Time {
	if o == nil || IsNil(o.NotAfter) {
		var ret time.Time
		return ret
	}
	return *o.NotAfter
}

// GetNotAfterOk returns a tuple with the NotAfter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *X509Certificate) GetNotAfterOk() (*time.Time, bool) {
	if o == nil || IsNil(o.NotAfter) {
		return nil, false
	}

	return o.NotAfter, true
}

// HasNotAfter returns a boolean if a field has been set.
func (o *X509Certificate) HasNotAfter() bool {
	if o != nil && !IsNil(o.NotAfter) {
		return true
	}

	return false
}

// SetNotAfter gets a reference to the given time.Time and assigns it to the NotAfter field.
func (o *X509Certificate) SetNotAfter(v time.Time) {
	o.NotAfter = &v
}

// GetNotBefore returns the NotBefore field value if set, zero value otherwise
func (o *X509Certificate) GetNotBefore() time.Time {
	if o == nil || IsNil(o.NotBefore) {
		var ret time.Time
		return ret
	}
	return *o.NotBefore
}

// GetNotBeforeOk returns a tuple with the NotBefore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *X509Certificate) GetNotBeforeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.NotBefore) {
		return nil, false
	}

	return o.NotBefore, true
}

// HasNotBefore returns a boolean if a field has been set.
func (o *X509Certificate) HasNotBefore() bool {
	if o != nil && !IsNil(o.NotBefore) {
		return true
	}

	return false
}

// SetNotBefore gets a reference to the given time.Time and assigns it to the NotBefore field.
func (o *X509Certificate) SetNotBefore(v time.Time) {
	o.NotBefore = &v
}
