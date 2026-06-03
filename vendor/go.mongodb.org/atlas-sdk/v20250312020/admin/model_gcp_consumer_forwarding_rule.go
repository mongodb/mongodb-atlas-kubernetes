// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GCPConsumerForwardingRule struct for GCPConsumerForwardingRule
type GCPConsumerForwardingRule struct {
	// Human-readable label that identifies the Google Cloud consumer forwarding rule that you created.
	// Read only field.
	EndpointName *string `json:"endpointName,omitempty"`
	// One Private Internet Protocol version 4 (IPv4) address to which this Google Cloud consumer forwarding rule resolves.
	// Read only field.
	IpAddress *string `json:"ipAddress,omitempty"`
	// State of the MongoDB Cloud endpoint group when MongoDB Cloud received this request.
	// Read only field.
	Status *string `json:"status,omitempty"`
}

// NewGCPConsumerForwardingRule instantiates a new GCPConsumerForwardingRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGCPConsumerForwardingRule() *GCPConsumerForwardingRule {
	this := GCPConsumerForwardingRule{}
	return &this
}

// NewGCPConsumerForwardingRuleWithDefaults instantiates a new GCPConsumerForwardingRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGCPConsumerForwardingRuleWithDefaults() *GCPConsumerForwardingRule {
	this := GCPConsumerForwardingRule{}
	return &this
}

// GetEndpointName returns the EndpointName field value if set, zero value otherwise
func (o *GCPConsumerForwardingRule) GetEndpointName() string {
	if o == nil || IsNil(o.EndpointName) {
		var ret string
		return ret
	}
	return *o.EndpointName
}

// GetEndpointNameOk returns a tuple with the EndpointName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GCPConsumerForwardingRule) GetEndpointNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointName) {
		return nil, false
	}

	return o.EndpointName, true
}

// HasEndpointName returns a boolean if a field has been set.
func (o *GCPConsumerForwardingRule) HasEndpointName() bool {
	if o != nil && !IsNil(o.EndpointName) {
		return true
	}

	return false
}

// SetEndpointName gets a reference to the given string and assigns it to the EndpointName field.
func (o *GCPConsumerForwardingRule) SetEndpointName(v string) {
	o.EndpointName = &v
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise
func (o *GCPConsumerForwardingRule) GetIpAddress() string {
	if o == nil || IsNil(o.IpAddress) {
		var ret string
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GCPConsumerForwardingRule) GetIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}

	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *GCPConsumerForwardingRule) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given string and assigns it to the IpAddress field.
func (o *GCPConsumerForwardingRule) SetIpAddress(v string) {
	o.IpAddress = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *GCPConsumerForwardingRule) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GCPConsumerForwardingRule) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *GCPConsumerForwardingRule) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *GCPConsumerForwardingRule) SetStatus(v string) {
	o.Status = &v
}
