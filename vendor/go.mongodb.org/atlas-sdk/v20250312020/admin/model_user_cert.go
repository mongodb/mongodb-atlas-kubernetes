// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// UserCert struct for UserCert
type UserCert struct {
	// Unique 24-hexadecimal character string that identifies this certificate.
	// Read only field.
	Id *int64 `json:"_id,omitempty"`
	// Date and time when MongoDB Cloud created this certificate. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Number of months that the certificate remains valid until it expires.
	// Write only field.
	MonthsUntilExpiration *int `json:"monthsUntilExpiration,omitempty"`
	// Date and time when this certificate expires. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	NotAfter *time.Time `json:"notAfter,omitempty"`
	// Subject Alternative Name associated with this certificate. This parameter expresses its value as a distinguished name as defined in RFC 2253.
	// Read only field.
	Subject *string `json:"subject,omitempty"`
}

// NewUserCert instantiates a new UserCert object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserCert() *UserCert {
	this := UserCert{}
	var monthsUntilExpiration int = 3
	this.MonthsUntilExpiration = &monthsUntilExpiration
	return &this
}

// NewUserCertWithDefaults instantiates a new UserCert object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserCertWithDefaults() *UserCert {
	this := UserCert{}
	var monthsUntilExpiration int = 3
	this.MonthsUntilExpiration = &monthsUntilExpiration
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *UserCert) GetId() int64 {
	if o == nil || IsNil(o.Id) {
		var ret int64
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetIdOk() (*int64, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *UserCert) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int64 and assigns it to the Id field.
func (o *UserCert) SetId(v int64) {
	o.Id = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *UserCert) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *UserCert) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *UserCert) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *UserCert) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *UserCert) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *UserCert) SetGroupId(v string) {
	o.GroupId = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *UserCert) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *UserCert) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *UserCert) SetLinks(v []Link) {
	o.Links = &v
}

// GetMonthsUntilExpiration returns the MonthsUntilExpiration field value if set, zero value otherwise
func (o *UserCert) GetMonthsUntilExpiration() int {
	if o == nil || IsNil(o.MonthsUntilExpiration) {
		var ret int
		return ret
	}
	return *o.MonthsUntilExpiration
}

// GetMonthsUntilExpirationOk returns a tuple with the MonthsUntilExpiration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetMonthsUntilExpirationOk() (*int, bool) {
	if o == nil || IsNil(o.MonthsUntilExpiration) {
		return nil, false
	}

	return o.MonthsUntilExpiration, true
}

// HasMonthsUntilExpiration returns a boolean if a field has been set.
func (o *UserCert) HasMonthsUntilExpiration() bool {
	if o != nil && !IsNil(o.MonthsUntilExpiration) {
		return true
	}

	return false
}

// SetMonthsUntilExpiration gets a reference to the given int and assigns it to the MonthsUntilExpiration field.
func (o *UserCert) SetMonthsUntilExpiration(v int) {
	o.MonthsUntilExpiration = &v
}

// GetNotAfter returns the NotAfter field value if set, zero value otherwise
func (o *UserCert) GetNotAfter() time.Time {
	if o == nil || IsNil(o.NotAfter) {
		var ret time.Time
		return ret
	}
	return *o.NotAfter
}

// GetNotAfterOk returns a tuple with the NotAfter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetNotAfterOk() (*time.Time, bool) {
	if o == nil || IsNil(o.NotAfter) {
		return nil, false
	}

	return o.NotAfter, true
}

// HasNotAfter returns a boolean if a field has been set.
func (o *UserCert) HasNotAfter() bool {
	if o != nil && !IsNil(o.NotAfter) {
		return true
	}

	return false
}

// SetNotAfter gets a reference to the given time.Time and assigns it to the NotAfter field.
func (o *UserCert) SetNotAfter(v time.Time) {
	o.NotAfter = &v
}

// GetSubject returns the Subject field value if set, zero value otherwise
func (o *UserCert) GetSubject() string {
	if o == nil || IsNil(o.Subject) {
		var ret string
		return ret
	}
	return *o.Subject
}

// GetSubjectOk returns a tuple with the Subject field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCert) GetSubjectOk() (*string, bool) {
	if o == nil || IsNil(o.Subject) {
		return nil, false
	}

	return o.Subject, true
}

// HasSubject returns a boolean if a field has been set.
func (o *UserCert) HasSubject() bool {
	if o != nil && !IsNil(o.Subject) {
		return true
	}

	return false
}

// SetSubject gets a reference to the given string and assigns it to the Subject field.
func (o *UserCert) SetSubject(v string) {
	o.Subject = &v
}
