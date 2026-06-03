// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ServiceAccountIPAccessListEntry struct for ServiceAccountIPAccessListEntry
type ServiceAccountIPAccessListEntry struct {
	// Range of network addresses in the access list for the Service Account. This parameter requires the range to be expressed in Classless Inter-Domain Routing (CIDR) notation of Internet Protocol version 4 or version 6 addresses. You can set a value for this parameter or `ipAddress`, but not for both in the same request.
	CidrBlock *string `json:"cidrBlock,omitempty"`
	// Date MongoDB Cloud added the entry was added to the Access List. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Network address in the access list for the Service Account. This parameter requires the address to be expressed as one Internet Protocol version 4 or version 6 address. You can set a value for this parameter or `cidrBlock`, but not for both in the same request.
	IpAddress *string `json:"ipAddress,omitempty"`
	// Network address that issued the most recent request to the API. This parameter requires the address to be expressed as one Internet Protocol version 4 or version 6 address. The resource returns this parameter after this IP address makes at least one request.
	// Read only field.
	LastUsedAddress *string `json:"lastUsedAddress,omitempty"`
	// Date when MongoDB Cloud received the most recent request that originated from this Internet Protocol version 4 or version 6 address. The resource returns this parameter when at least one request originates from this IP address. MongoDB Cloud updates this parameter each time a client accesses the permitted resource, with a delay of up to 5 minutes. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	// The number of requests that has originated from this network address.
	// Read only field.
	RequestCount *int `json:"requestCount,omitempty"`
}

// NewServiceAccountIPAccessListEntry instantiates a new ServiceAccountIPAccessListEntry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServiceAccountIPAccessListEntry() *ServiceAccountIPAccessListEntry {
	this := ServiceAccountIPAccessListEntry{}
	return &this
}

// NewServiceAccountIPAccessListEntryWithDefaults instantiates a new ServiceAccountIPAccessListEntry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServiceAccountIPAccessListEntryWithDefaults() *ServiceAccountIPAccessListEntry {
	this := ServiceAccountIPAccessListEntry{}
	return &this
}

// GetCidrBlock returns the CidrBlock field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetCidrBlock() string {
	if o == nil || IsNil(o.CidrBlock) {
		var ret string
		return ret
	}
	return *o.CidrBlock
}

// GetCidrBlockOk returns a tuple with the CidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.CidrBlock) {
		return nil, false
	}

	return o.CidrBlock, true
}

// HasCidrBlock returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasCidrBlock() bool {
	if o != nil && !IsNil(o.CidrBlock) {
		return true
	}

	return false
}

// SetCidrBlock gets a reference to the given string and assigns it to the CidrBlock field.
func (o *ServiceAccountIPAccessListEntry) SetCidrBlock(v string) {
	o.CidrBlock = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *ServiceAccountIPAccessListEntry) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *ServiceAccountIPAccessListEntry) SetIpAddress(v string) {
	o.IpAddress = &v
}

// GetLastUsedAddress returns the LastUsedAddress field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetLastUsedAddress() string {
	if o == nil || IsNil(o.LastUsedAddress) {
		var ret string
		return ret
	}
	return *o.LastUsedAddress
}

// GetLastUsedAddressOk returns a tuple with the LastUsedAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetLastUsedAddressOk() (*string, bool) {
	if o == nil || IsNil(o.LastUsedAddress) {
		return nil, false
	}

	return o.LastUsedAddress, true
}

// HasLastUsedAddress returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasLastUsedAddress() bool {
	if o != nil && !IsNil(o.LastUsedAddress) {
		return true
	}

	return false
}

// SetLastUsedAddress gets a reference to the given string and assigns it to the LastUsedAddress field.
func (o *ServiceAccountIPAccessListEntry) SetLastUsedAddress(v string) {
	o.LastUsedAddress = &v
}

// GetLastUsedAt returns the LastUsedAt field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetLastUsedAt() time.Time {
	if o == nil || IsNil(o.LastUsedAt) {
		var ret time.Time
		return ret
	}
	return *o.LastUsedAt
}

// GetLastUsedAtOk returns a tuple with the LastUsedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetLastUsedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastUsedAt) {
		return nil, false
	}

	return o.LastUsedAt, true
}

// HasLastUsedAt returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasLastUsedAt() bool {
	if o != nil && !IsNil(o.LastUsedAt) {
		return true
	}

	return false
}

// SetLastUsedAt gets a reference to the given time.Time and assigns it to the LastUsedAt field.
func (o *ServiceAccountIPAccessListEntry) SetLastUsedAt(v time.Time) {
	o.LastUsedAt = &v
}

// GetRequestCount returns the RequestCount field value if set, zero value otherwise
func (o *ServiceAccountIPAccessListEntry) GetRequestCount() int {
	if o == nil || IsNil(o.RequestCount) {
		var ret int
		return ret
	}
	return *o.RequestCount
}

// GetRequestCountOk returns a tuple with the RequestCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServiceAccountIPAccessListEntry) GetRequestCountOk() (*int, bool) {
	if o == nil || IsNil(o.RequestCount) {
		return nil, false
	}

	return o.RequestCount, true
}

// HasRequestCount returns a boolean if a field has been set.
func (o *ServiceAccountIPAccessListEntry) HasRequestCount() bool {
	if o != nil && !IsNil(o.RequestCount) {
		return true
	}

	return false
}

// SetRequestCount gets a reference to the given int and assigns it to the RequestCount field.
func (o *ServiceAccountIPAccessListEntry) SetRequestCount(v int) {
	o.RequestCount = &v
}
