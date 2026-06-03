// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PrivateLinkEndpoint struct for PrivateLinkEndpoint
type PrivateLinkEndpoint struct {
	// Cloud service provider that serves the requested endpoint.
	// Read only field.
	CloudProvider string `json:"cloudProvider"`
	// Flag that indicates whether MongoDB Cloud received a request to remove the specified private endpoint from the private endpoint service.
	// Read only field.
	DeleteRequested *bool `json:"deleteRequested,omitempty"`
	// Error message returned when requesting private connection resource. The resource returns `null` if the request succeeded.
	// Read only field.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// State of the Amazon Web Service PrivateLink connection when MongoDB Cloud received this request.
	// Read only field.
	ConnectionStatus *string `json:"connectionStatus,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the interface endpoint.
	// Read only field.
	InterfaceEndpointId *string `json:"interfaceEndpointId,omitempty"`
	// Human-readable label that MongoDB Cloud generates that identifies the private endpoint connection.
	// Read only field.
	PrivateEndpointConnectionName *string `json:"privateEndpointConnectionName,omitempty"`
	// IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	PrivateEndpointIPAddress *string `json:"privateEndpointIPAddress,omitempty"`
	// Unique string that identifies the Azure private endpoint's network interface that someone added to this private endpoint service.
	// Read only field.
	PrivateEndpointResourceId *string `json:"privateEndpointResourceId,omitempty"`
	// State of the Google Cloud network endpoint group when MongoDB Cloud received this request.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Human-readable label that identifies a set of endpoints. If this private endpoint belongs to a port-mapped endpoint service, this field is the private endpoint name.
	// Read only field.
	EndpointGroupName *string `json:"endpointGroupName,omitempty"`
	// List of individual private endpoints that comprise this endpoint group. If this endpoint belongs to a port-mapped endpoint service, this field will only contain a list of one private endpoint.
	// Read only field.
	Endpoints *[]GCPConsumerForwardingRule `json:"endpoints,omitempty"`
	// Unique string that identifies the Google Cloud project in which you created the endpoints.
	// Read only field.
	GcpProjectId *string `json:"gcpProjectId,omitempty"`
	// Flag that indicates whether the endpoint service for this endpoint group uses PSC port-mapping.
	PortMappingEnabled *bool `json:"portMappingEnabled,omitempty"`
}

// NewPrivateLinkEndpoint instantiates a new PrivateLinkEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateLinkEndpoint(cloudProvider string) *PrivateLinkEndpoint {
	this := PrivateLinkEndpoint{}
	this.CloudProvider = cloudProvider
	return &this
}

// NewPrivateLinkEndpointWithDefaults instantiates a new PrivateLinkEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateLinkEndpointWithDefaults() *PrivateLinkEndpoint {
	this := PrivateLinkEndpoint{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *PrivateLinkEndpoint) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *PrivateLinkEndpoint) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetDeleteRequested returns the DeleteRequested field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetDeleteRequested() bool {
	if o == nil || IsNil(o.DeleteRequested) {
		var ret bool
		return ret
	}
	return *o.DeleteRequested
}

// GetDeleteRequestedOk returns a tuple with the DeleteRequested field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetDeleteRequestedOk() (*bool, bool) {
	if o == nil || IsNil(o.DeleteRequested) {
		return nil, false
	}

	return o.DeleteRequested, true
}

// HasDeleteRequested returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasDeleteRequested() bool {
	if o != nil && !IsNil(o.DeleteRequested) {
		return true
	}

	return false
}

// SetDeleteRequested gets a reference to the given bool and assigns it to the DeleteRequested field.
func (o *PrivateLinkEndpoint) SetDeleteRequested(v bool) {
	o.DeleteRequested = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *PrivateLinkEndpoint) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetConnectionStatus returns the ConnectionStatus field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetConnectionStatus() string {
	if o == nil || IsNil(o.ConnectionStatus) {
		var ret string
		return ret
	}
	return *o.ConnectionStatus
}

// GetConnectionStatusOk returns a tuple with the ConnectionStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetConnectionStatusOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionStatus) {
		return nil, false
	}

	return o.ConnectionStatus, true
}

// HasConnectionStatus returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasConnectionStatus() bool {
	if o != nil && !IsNil(o.ConnectionStatus) {
		return true
	}

	return false
}

// SetConnectionStatus gets a reference to the given string and assigns it to the ConnectionStatus field.
func (o *PrivateLinkEndpoint) SetConnectionStatus(v string) {
	o.ConnectionStatus = &v
}

// GetInterfaceEndpointId returns the InterfaceEndpointId field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetInterfaceEndpointId() string {
	if o == nil || IsNil(o.InterfaceEndpointId) {
		var ret string
		return ret
	}
	return *o.InterfaceEndpointId
}

// GetInterfaceEndpointIdOk returns a tuple with the InterfaceEndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetInterfaceEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.InterfaceEndpointId) {
		return nil, false
	}

	return o.InterfaceEndpointId, true
}

// HasInterfaceEndpointId returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasInterfaceEndpointId() bool {
	if o != nil && !IsNil(o.InterfaceEndpointId) {
		return true
	}

	return false
}

// SetInterfaceEndpointId gets a reference to the given string and assigns it to the InterfaceEndpointId field.
func (o *PrivateLinkEndpoint) SetInterfaceEndpointId(v string) {
	o.InterfaceEndpointId = &v
}

// GetPrivateEndpointConnectionName returns the PrivateEndpointConnectionName field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetPrivateEndpointConnectionName() string {
	if o == nil || IsNil(o.PrivateEndpointConnectionName) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointConnectionName
}

// GetPrivateEndpointConnectionNameOk returns a tuple with the PrivateEndpointConnectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetPrivateEndpointConnectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointConnectionName) {
		return nil, false
	}

	return o.PrivateEndpointConnectionName, true
}

// HasPrivateEndpointConnectionName returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasPrivateEndpointConnectionName() bool {
	if o != nil && !IsNil(o.PrivateEndpointConnectionName) {
		return true
	}

	return false
}

// SetPrivateEndpointConnectionName gets a reference to the given string and assigns it to the PrivateEndpointConnectionName field.
func (o *PrivateLinkEndpoint) SetPrivateEndpointConnectionName(v string) {
	o.PrivateEndpointConnectionName = &v
}

// GetPrivateEndpointIPAddress returns the PrivateEndpointIPAddress field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetPrivateEndpointIPAddress() string {
	if o == nil || IsNil(o.PrivateEndpointIPAddress) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointIPAddress
}

// GetPrivateEndpointIPAddressOk returns a tuple with the PrivateEndpointIPAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetPrivateEndpointIPAddressOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointIPAddress) {
		return nil, false
	}

	return o.PrivateEndpointIPAddress, true
}

// HasPrivateEndpointIPAddress returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasPrivateEndpointIPAddress() bool {
	if o != nil && !IsNil(o.PrivateEndpointIPAddress) {
		return true
	}

	return false
}

// SetPrivateEndpointIPAddress gets a reference to the given string and assigns it to the PrivateEndpointIPAddress field.
func (o *PrivateLinkEndpoint) SetPrivateEndpointIPAddress(v string) {
	o.PrivateEndpointIPAddress = &v
}

// GetPrivateEndpointResourceId returns the PrivateEndpointResourceId field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetPrivateEndpointResourceId() string {
	if o == nil || IsNil(o.PrivateEndpointResourceId) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointResourceId
}

// GetPrivateEndpointResourceIdOk returns a tuple with the PrivateEndpointResourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetPrivateEndpointResourceIdOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointResourceId) {
		return nil, false
	}

	return o.PrivateEndpointResourceId, true
}

// HasPrivateEndpointResourceId returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasPrivateEndpointResourceId() bool {
	if o != nil && !IsNil(o.PrivateEndpointResourceId) {
		return true
	}

	return false
}

// SetPrivateEndpointResourceId gets a reference to the given string and assigns it to the PrivateEndpointResourceId field.
func (o *PrivateLinkEndpoint) SetPrivateEndpointResourceId(v string) {
	o.PrivateEndpointResourceId = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *PrivateLinkEndpoint) SetStatus(v string) {
	o.Status = &v
}

// GetEndpointGroupName returns the EndpointGroupName field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetEndpointGroupName() string {
	if o == nil || IsNil(o.EndpointGroupName) {
		var ret string
		return ret
	}
	return *o.EndpointGroupName
}

// GetEndpointGroupNameOk returns a tuple with the EndpointGroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetEndpointGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointGroupName) {
		return nil, false
	}

	return o.EndpointGroupName, true
}

// HasEndpointGroupName returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasEndpointGroupName() bool {
	if o != nil && !IsNil(o.EndpointGroupName) {
		return true
	}

	return false
}

// SetEndpointGroupName gets a reference to the given string and assigns it to the EndpointGroupName field.
func (o *PrivateLinkEndpoint) SetEndpointGroupName(v string) {
	o.EndpointGroupName = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetEndpoints() []GCPConsumerForwardingRule {
	if o == nil || IsNil(o.Endpoints) {
		var ret []GCPConsumerForwardingRule
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetEndpointsOk() (*[]GCPConsumerForwardingRule, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}

	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []GCPConsumerForwardingRule and assigns it to the Endpoints field.
func (o *PrivateLinkEndpoint) SetEndpoints(v []GCPConsumerForwardingRule) {
	o.Endpoints = &v
}

// GetGcpProjectId returns the GcpProjectId field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetGcpProjectId() string {
	if o == nil || IsNil(o.GcpProjectId) {
		var ret string
		return ret
	}
	return *o.GcpProjectId
}

// GetGcpProjectIdOk returns a tuple with the GcpProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetGcpProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.GcpProjectId) {
		return nil, false
	}

	return o.GcpProjectId, true
}

// HasGcpProjectId returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasGcpProjectId() bool {
	if o != nil && !IsNil(o.GcpProjectId) {
		return true
	}

	return false
}

// SetGcpProjectId gets a reference to the given string and assigns it to the GcpProjectId field.
func (o *PrivateLinkEndpoint) SetGcpProjectId(v string) {
	o.GcpProjectId = &v
}

// GetPortMappingEnabled returns the PortMappingEnabled field value if set, zero value otherwise
func (o *PrivateLinkEndpoint) GetPortMappingEnabled() bool {
	if o == nil || IsNil(o.PortMappingEnabled) {
		var ret bool
		return ret
	}
	return *o.PortMappingEnabled
}

// GetPortMappingEnabledOk returns a tuple with the PortMappingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateLinkEndpoint) GetPortMappingEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.PortMappingEnabled) {
		return nil, false
	}

	return o.PortMappingEnabled, true
}

// HasPortMappingEnabled returns a boolean if a field has been set.
func (o *PrivateLinkEndpoint) HasPortMappingEnabled() bool {
	if o != nil && !IsNil(o.PortMappingEnabled) {
		return true
	}

	return false
}

// SetPortMappingEnabled gets a reference to the given bool and assigns it to the PortMappingEnabled field.
func (o *PrivateLinkEndpoint) SetPortMappingEnabled(v bool) {
	o.PortMappingEnabled = &v
}
