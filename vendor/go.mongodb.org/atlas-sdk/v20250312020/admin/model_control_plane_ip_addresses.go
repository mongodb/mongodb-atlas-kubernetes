// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ControlPlaneIPAddresses List of IP addresses in the Atlas control plane.
type ControlPlaneIPAddresses struct {
	Inbound  *InboundControlPlaneCloudProviderIPAddresses  `json:"inbound,omitempty"`
	Outbound *OutboundControlPlaneCloudProviderIPAddresses `json:"outbound,omitempty"`
}

// NewControlPlaneIPAddresses instantiates a new ControlPlaneIPAddresses object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewControlPlaneIPAddresses() *ControlPlaneIPAddresses {
	this := ControlPlaneIPAddresses{}
	return &this
}

// NewControlPlaneIPAddressesWithDefaults instantiates a new ControlPlaneIPAddresses object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewControlPlaneIPAddressesWithDefaults() *ControlPlaneIPAddresses {
	this := ControlPlaneIPAddresses{}
	return &this
}

// GetInbound returns the Inbound field value if set, zero value otherwise
func (o *ControlPlaneIPAddresses) GetInbound() InboundControlPlaneCloudProviderIPAddresses {
	if o == nil || IsNil(o.Inbound) {
		var ret InboundControlPlaneCloudProviderIPAddresses
		return ret
	}
	return *o.Inbound
}

// GetInboundOk returns a tuple with the Inbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ControlPlaneIPAddresses) GetInboundOk() (*InboundControlPlaneCloudProviderIPAddresses, bool) {
	if o == nil || IsNil(o.Inbound) {
		return nil, false
	}

	return o.Inbound, true
}

// HasInbound returns a boolean if a field has been set.
func (o *ControlPlaneIPAddresses) HasInbound() bool {
	if o != nil && !IsNil(o.Inbound) {
		return true
	}

	return false
}

// SetInbound gets a reference to the given InboundControlPlaneCloudProviderIPAddresses and assigns it to the Inbound field.
func (o *ControlPlaneIPAddresses) SetInbound(v InboundControlPlaneCloudProviderIPAddresses) {
	o.Inbound = &v
}

// GetOutbound returns the Outbound field value if set, zero value otherwise
func (o *ControlPlaneIPAddresses) GetOutbound() OutboundControlPlaneCloudProviderIPAddresses {
	if o == nil || IsNil(o.Outbound) {
		var ret OutboundControlPlaneCloudProviderIPAddresses
		return ret
	}
	return *o.Outbound
}

// GetOutboundOk returns a tuple with the Outbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ControlPlaneIPAddresses) GetOutboundOk() (*OutboundControlPlaneCloudProviderIPAddresses, bool) {
	if o == nil || IsNil(o.Outbound) {
		return nil, false
	}

	return o.Outbound, true
}

// HasOutbound returns a boolean if a field has been set.
func (o *ControlPlaneIPAddresses) HasOutbound() bool {
	if o != nil && !IsNil(o.Outbound) {
		return true
	}

	return false
}

// SetOutbound gets a reference to the given OutboundControlPlaneCloudProviderIPAddresses and assigns it to the Outbound field.
func (o *ControlPlaneIPAddresses) SetOutbound(v OutboundControlPlaneCloudProviderIPAddresses) {
	o.Outbound = &v
}
