// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterIPAddresses List of IP addresses in a cluster.
type ClusterIPAddresses struct {
	// Human-readable label that identifies the cluster.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// List of future inbound IP addresses associated with the cluster. If your network allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your application can connect to your Atlas cluster.
	// Read only field.
	FutureInbound *[]string `json:"futureInbound,omitempty"`
	// List of future outbound IP addresses associated with the cluster. If your network allows inbound HTTP requests only from specific IP addresses, you must allow access from the following IP addresses so that your Atlas cluster can communicate with your webhooks and KMS.
	// Read only field.
	FutureOutbound *[]string `json:"futureOutbound,omitempty"`
	// List of inbound IP addresses associated with the cluster. If your network allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your application can connect to your Atlas cluster.
	// Read only field.
	Inbound *[]string `json:"inbound,omitempty"`
	// List of outbound IP addresses associated with the cluster. If your network allows inbound HTTP requests only from specific IP addresses, you must allow access from the following IP addresses so that your Atlas cluster can communicate with your webhooks and KMS.
	// Read only field.
	Outbound *[]string `json:"outbound,omitempty"`
}

// NewClusterIPAddresses instantiates a new ClusterIPAddresses object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterIPAddresses() *ClusterIPAddresses {
	this := ClusterIPAddresses{}
	return &this
}

// NewClusterIPAddressesWithDefaults instantiates a new ClusterIPAddresses object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterIPAddressesWithDefaults() *ClusterIPAddresses {
	this := ClusterIPAddresses{}
	return &this
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *ClusterIPAddresses) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterIPAddresses) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *ClusterIPAddresses) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *ClusterIPAddresses) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetFutureInbound returns the FutureInbound field value if set, zero value otherwise
func (o *ClusterIPAddresses) GetFutureInbound() []string {
	if o == nil || IsNil(o.FutureInbound) {
		var ret []string
		return ret
	}
	return *o.FutureInbound
}

// GetFutureInboundOk returns a tuple with the FutureInbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterIPAddresses) GetFutureInboundOk() (*[]string, bool) {
	if o == nil || IsNil(o.FutureInbound) {
		return nil, false
	}

	return o.FutureInbound, true
}

// HasFutureInbound returns a boolean if a field has been set.
func (o *ClusterIPAddresses) HasFutureInbound() bool {
	if o != nil && !IsNil(o.FutureInbound) {
		return true
	}

	return false
}

// SetFutureInbound gets a reference to the given []string and assigns it to the FutureInbound field.
func (o *ClusterIPAddresses) SetFutureInbound(v []string) {
	o.FutureInbound = &v
}

// GetFutureOutbound returns the FutureOutbound field value if set, zero value otherwise
func (o *ClusterIPAddresses) GetFutureOutbound() []string {
	if o == nil || IsNil(o.FutureOutbound) {
		var ret []string
		return ret
	}
	return *o.FutureOutbound
}

// GetFutureOutboundOk returns a tuple with the FutureOutbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterIPAddresses) GetFutureOutboundOk() (*[]string, bool) {
	if o == nil || IsNil(o.FutureOutbound) {
		return nil, false
	}

	return o.FutureOutbound, true
}

// HasFutureOutbound returns a boolean if a field has been set.
func (o *ClusterIPAddresses) HasFutureOutbound() bool {
	if o != nil && !IsNil(o.FutureOutbound) {
		return true
	}

	return false
}

// SetFutureOutbound gets a reference to the given []string and assigns it to the FutureOutbound field.
func (o *ClusterIPAddresses) SetFutureOutbound(v []string) {
	o.FutureOutbound = &v
}

// GetInbound returns the Inbound field value if set, zero value otherwise
func (o *ClusterIPAddresses) GetInbound() []string {
	if o == nil || IsNil(o.Inbound) {
		var ret []string
		return ret
	}
	return *o.Inbound
}

// GetInboundOk returns a tuple with the Inbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterIPAddresses) GetInboundOk() (*[]string, bool) {
	if o == nil || IsNil(o.Inbound) {
		return nil, false
	}

	return o.Inbound, true
}

// HasInbound returns a boolean if a field has been set.
func (o *ClusterIPAddresses) HasInbound() bool {
	if o != nil && !IsNil(o.Inbound) {
		return true
	}

	return false
}

// SetInbound gets a reference to the given []string and assigns it to the Inbound field.
func (o *ClusterIPAddresses) SetInbound(v []string) {
	o.Inbound = &v
}

// GetOutbound returns the Outbound field value if set, zero value otherwise
func (o *ClusterIPAddresses) GetOutbound() []string {
	if o == nil || IsNil(o.Outbound) {
		var ret []string
		return ret
	}
	return *o.Outbound
}

// GetOutboundOk returns a tuple with the Outbound field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterIPAddresses) GetOutboundOk() (*[]string, bool) {
	if o == nil || IsNil(o.Outbound) {
		return nil, false
	}

	return o.Outbound, true
}

// HasOutbound returns a boolean if a field has been set.
func (o *ClusterIPAddresses) HasOutbound() bool {
	if o != nil && !IsNil(o.Outbound) {
		return true
	}

	return false
}

// SetOutbound gets a reference to the given []string and assigns it to the Outbound field.
func (o *ClusterIPAddresses) SetOutbound(v []string) {
	o.Outbound = &v
}
