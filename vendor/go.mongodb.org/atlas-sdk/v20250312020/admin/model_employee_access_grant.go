// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// EmployeeAccessGrant MongoDB employee granted access level and expiration for a cluster.
type EmployeeAccessGrant struct {
	// Expiration date for the employee access grant. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	ExpirationTime time.Time `json:"expirationTime"`
	// Level of access to grant to MongoDB Employees.
	GrantType string `json:"grantType"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewEmployeeAccessGrant instantiates a new EmployeeAccessGrant object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEmployeeAccessGrant(expirationTime time.Time, grantType string) *EmployeeAccessGrant {
	this := EmployeeAccessGrant{}
	this.ExpirationTime = expirationTime
	this.GrantType = grantType
	return &this
}

// NewEmployeeAccessGrantWithDefaults instantiates a new EmployeeAccessGrant object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEmployeeAccessGrantWithDefaults() *EmployeeAccessGrant {
	this := EmployeeAccessGrant{}
	return &this
}

// GetExpirationTime returns the ExpirationTime field value
func (o *EmployeeAccessGrant) GetExpirationTime() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.ExpirationTime
}

// GetExpirationTimeOk returns a tuple with the ExpirationTime field value
// and a boolean to check if the value has been set.
func (o *EmployeeAccessGrant) GetExpirationTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpirationTime, true
}

// SetExpirationTime sets field value
func (o *EmployeeAccessGrant) SetExpirationTime(v time.Time) {
	o.ExpirationTime = v
}

// GetGrantType returns the GrantType field value
func (o *EmployeeAccessGrant) GetGrantType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GrantType
}

// GetGrantTypeOk returns a tuple with the GrantType field value
// and a boolean to check if the value has been set.
func (o *EmployeeAccessGrant) GetGrantTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GrantType, true
}

// SetGrantType sets field value
func (o *EmployeeAccessGrant) SetGrantType(v string) {
	o.GrantType = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *EmployeeAccessGrant) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EmployeeAccessGrant) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *EmployeeAccessGrant) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *EmployeeAccessGrant) SetLinks(v []Link) {
	o.Links = &v
}
