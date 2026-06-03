// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// UserAccessListResponse struct for UserAccessListResponse
type UserAccessListResponse struct {
	// Range of IP addresses in Classless Inter-Domain Routing (CIDR) notation in the access list for the API key.
	CidrBlock *string `json:"cidrBlock,omitempty"`
	// Total number of requests that have originated from the Internet Protocol (IP) address given as the value of the `lastUsedAddress` parameter.
	// Read only field.
	Count *int `json:"count,omitempty"`
	// Date and time when someone added the network addresses to the specified API access list. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Network address in the access list for the API key.
	IpAddress *string `json:"ipAddress,omitempty"`
	// Date and time when MongoDB Cloud received the most recent request that originated from this Internet Protocol version 4 or version 6 address. The resource returns this parameter when at least one request has originated from this IP address. MongoDB Cloud updates this parameter each time a client accesses the permitted resource. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastUsed *time.Time `json:"lastUsed,omitempty"`
	// Network address that issued the most recent request to the API. This parameter requires the address to be expressed as one Internet Protocol version 4 or version 6 address. The resource returns this parameter after this IP address made at least one request.
	// Read only field.
	LastUsedAddress *string `json:"lastUsedAddress,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewUserAccessListResponse instantiates a new UserAccessListResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAccessListResponse() *UserAccessListResponse {
	this := UserAccessListResponse{}
	return &this
}

// NewUserAccessListResponseWithDefaults instantiates a new UserAccessListResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAccessListResponseWithDefaults() *UserAccessListResponse {
	this := UserAccessListResponse{}
	return &this
}

// GetCidrBlock returns the CidrBlock field value if set, zero value otherwise
func (o *UserAccessListResponse) GetCidrBlock() string {
	if o == nil || IsNil(o.CidrBlock) {
		var ret string
		return ret
	}
	return *o.CidrBlock
}

// GetCidrBlockOk returns a tuple with the CidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.CidrBlock) {
		return nil, false
	}

	return o.CidrBlock, true
}

// HasCidrBlock returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasCidrBlock() bool {
	if o != nil && !IsNil(o.CidrBlock) {
		return true
	}

	return false
}

// SetCidrBlock gets a reference to the given string and assigns it to the CidrBlock field.
func (o *UserAccessListResponse) SetCidrBlock(v string) {
	o.CidrBlock = &v
}

// GetCount returns the Count field value if set, zero value otherwise
func (o *UserAccessListResponse) GetCount() int {
	if o == nil || IsNil(o.Count) {
		var ret int
		return ret
	}
	return *o.Count
}

// GetCountOk returns a tuple with the Count field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetCountOk() (*int, bool) {
	if o == nil || IsNil(o.Count) {
		return nil, false
	}

	return o.Count, true
}

// HasCount returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasCount() bool {
	if o != nil && !IsNil(o.Count) {
		return true
	}

	return false
}

// SetCount gets a reference to the given int and assigns it to the Count field.
func (o *UserAccessListResponse) SetCount(v int) {
	o.Count = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *UserAccessListResponse) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *UserAccessListResponse) SetCreated(v time.Time) {
	o.Created = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *UserAccessListResponse) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *UserAccessListResponse) SetIpAddress(v string) {
	o.IpAddress = &v
}

// GetLastUsed returns the LastUsed field value if set, zero value otherwise
func (o *UserAccessListResponse) GetLastUsed() time.Time {
	if o == nil || IsNil(o.LastUsed) {
		var ret time.Time
		return ret
	}
	return *o.LastUsed
}

// GetLastUsedOk returns a tuple with the LastUsed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetLastUsedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUsed) {
		return nil, false
	}

	return o.LastUsed, true
}

// HasLastUsed returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasLastUsed() bool {
	if o != nil && !IsNil(o.LastUsed) {
		return true
	}

	return false
}

// SetLastUsed gets a reference to the given time.Time and assigns it to the LastUsed field.
func (o *UserAccessListResponse) SetLastUsed(v time.Time) {
	o.LastUsed = &v
}

// GetLastUsedAddress returns the LastUsedAddress field value if set, zero value otherwise
func (o *UserAccessListResponse) GetLastUsedAddress() string {
	if o == nil || IsNil(o.LastUsedAddress) {
		var ret string
		return ret
	}
	return *o.LastUsedAddress
}

// GetLastUsedAddressOk returns a tuple with the LastUsedAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetLastUsedAddressOk() (*string, bool) {
	if o == nil || IsNil(o.LastUsedAddress) {
		return nil, false
	}

	return o.LastUsedAddress, true
}

// HasLastUsedAddress returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasLastUsedAddress() bool {
	if o != nil && !IsNil(o.LastUsedAddress) {
		return true
	}

	return false
}

// SetLastUsedAddress gets a reference to the given string and assigns it to the LastUsedAddress field.
func (o *UserAccessListResponse) SetLastUsedAddress(v string) {
	o.LastUsedAddress = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *UserAccessListResponse) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListResponse) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *UserAccessListResponse) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *UserAccessListResponse) SetLinks(v []Link) {
	o.Links = &v
}
