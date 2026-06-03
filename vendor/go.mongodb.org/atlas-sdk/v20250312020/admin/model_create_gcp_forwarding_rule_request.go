// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CreateGCPForwardingRuleRequest struct for CreateGCPForwardingRuleRequest
type CreateGCPForwardingRuleRequest struct {
	// Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.
	// Write only field.
	EndpointName *string `json:"endpointName,omitempty"`
	// One Private Internet Protocol version 4 (IPv4) address to which this Google Cloud consumer forwarding rule resolves.
	// Write only field.
	IpAddress *string `json:"ipAddress,omitempty"`
}

// NewCreateGCPForwardingRuleRequest instantiates a new CreateGCPForwardingRuleRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateGCPForwardingRuleRequest() *CreateGCPForwardingRuleRequest {
	this := CreateGCPForwardingRuleRequest{}
	return &this
}

// NewCreateGCPForwardingRuleRequestWithDefaults instantiates a new CreateGCPForwardingRuleRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateGCPForwardingRuleRequestWithDefaults() *CreateGCPForwardingRuleRequest {
	this := CreateGCPForwardingRuleRequest{}
	return &this
}

// GetEndpointName returns the EndpointName field value if set, zero value otherwise
func (o *CreateGCPForwardingRuleRequest) GetEndpointName() string {
	if o == nil || IsNil(o.EndpointName) {
		var ret string
		return ret
	}
	return *o.EndpointName
}

// GetEndpointNameOk returns a tuple with the EndpointName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateGCPForwardingRuleRequest) GetEndpointNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointName) {
		return nil, false
	}

	return o.EndpointName, true
}

// HasEndpointName returns a boolean if a field has been set.
func (o *CreateGCPForwardingRuleRequest) HasEndpointName() bool {
	if o != nil && !IsNil(o.EndpointName) {
		return true
	}

	return false
}

// SetEndpointName gets a reference to the given string and assigns it to the EndpointName field.
func (o *CreateGCPForwardingRuleRequest) SetEndpointName(v string) {
	o.EndpointName = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *CreateGCPForwardingRuleRequest) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateGCPForwardingRuleRequest) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *CreateGCPForwardingRuleRequest) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *CreateGCPForwardingRuleRequest) SetIpAddress(v string) {
	o.IpAddress = &v
}
