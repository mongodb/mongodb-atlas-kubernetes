// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UserAccessListRequest struct for UserAccessListRequest
type UserAccessListRequest struct {
	// Range of network addresses that you want to add to the access list for the API key. This parameter requires the range to be expressed in classless inter-domain routing (CIDR) notation of Internet Protocol version 4 or version 6 addresses. You can set a value for this parameter or `ipAddress` but not both in the same request.
	CidrBlock *string `json:"cidrBlock,omitempty"`
	// Network address that you want to add to the access list for the API key. This parameter requires the address to be expressed as one Internet Protocol version 4 or version 6 address. You can set a value for this parameter or `cidrBlock` but not both in the same request.
	IpAddress *string `json:"ipAddress,omitempty"`
}

// NewUserAccessListRequest instantiates a new UserAccessListRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAccessListRequest() *UserAccessListRequest {
	this := UserAccessListRequest{}
	return &this
}

// NewUserAccessListRequestWithDefaults instantiates a new UserAccessListRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAccessListRequestWithDefaults() *UserAccessListRequest {
	this := UserAccessListRequest{}
	return &this
}

// GetCidrBlock returns the CidrBlock field value if set, zero value otherwise
func (o *UserAccessListRequest) GetCidrBlock() string {
	if o == nil || IsNil(o.CidrBlock) {
		var ret string
		return ret
	}
	return *o.CidrBlock
}

// GetCidrBlockOk returns a tuple with the CidrBlock field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListRequest) GetCidrBlockOk() (*string, bool) {
	if o == nil || IsNil(o.CidrBlock) {
		return nil, false
	}

	return o.CidrBlock, true
}

// HasCidrBlock returns a boolean if a field has been set.
func (o *UserAccessListRequest) HasCidrBlock() bool {
	if o != nil && !IsNil(o.CidrBlock) {
		return true
	}

	return false
}

// SetCidrBlock gets a reference to the given string and assigns it to the CidrBlock field.
func (o *UserAccessListRequest) SetCidrBlock(v string) {
	o.CidrBlock = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *UserAccessListRequest) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessListRequest) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *UserAccessListRequest) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *UserAccessListRequest) SetIpAddress(v string) {
	o.IpAddress = &v
}
