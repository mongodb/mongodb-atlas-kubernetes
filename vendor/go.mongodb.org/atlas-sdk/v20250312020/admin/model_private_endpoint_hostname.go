// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PrivateEndpointHostname Set of Private endpoint and hostnames.
type PrivateEndpointHostname struct {
	// Human-readable label that identifies the hostname.
	// Read only field.
	Hostname *string `json:"hostname,omitempty"`
	// Human-readable label that identifies private endpoint.
	// Read only field.
	PrivateEndpoint *string `json:"privateEndpoint,omitempty"`
}

// NewPrivateEndpointHostname instantiates a new PrivateEndpointHostname object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateEndpointHostname() *PrivateEndpointHostname {
	this := PrivateEndpointHostname{}
	return &this
}

// NewPrivateEndpointHostnameWithDefaults instantiates a new PrivateEndpointHostname object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateEndpointHostnameWithDefaults() *PrivateEndpointHostname {
	this := PrivateEndpointHostname{}
	return &this
}

// GetHostname returns the Hostname field value if set, zero value otherwise
func (o *PrivateEndpointHostname) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointHostname) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}

	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *PrivateEndpointHostname) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *PrivateEndpointHostname) SetHostname(v string) {
	o.Hostname = &v
}

// GetPrivateEndpoint returns the PrivateEndpoint field value if set, zero value otherwise
func (o *PrivateEndpointHostname) GetPrivateEndpoint() string {
	if o == nil || IsNil(o.PrivateEndpoint) {
		var ret string
		return ret
	}
	return *o.PrivateEndpoint
}

// GetPrivateEndpointOk returns a tuple with the PrivateEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointHostname) GetPrivateEndpointOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpoint) {
		return nil, false
	}

	return o.PrivateEndpoint, true
}

// HasPrivateEndpoint returns a boolean if a field has been set.
func (o *PrivateEndpointHostname) HasPrivateEndpoint() bool {
	if o != nil && !IsNil(o.PrivateEndpoint) {
		return true
	}

	return false
}

// SetPrivateEndpoint gets a reference to the given string and assigns it to the PrivateEndpoint field.
func (o *PrivateEndpointHostname) SetPrivateEndpoint(v string) {
	o.PrivateEndpoint = &v
}
