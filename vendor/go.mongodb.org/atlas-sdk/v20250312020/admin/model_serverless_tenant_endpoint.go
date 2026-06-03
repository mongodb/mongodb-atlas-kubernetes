// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessTenantEndpoint struct for ServerlessTenantEndpoint
type ServerlessTenantEndpoint struct {
	// Unique 24-hexadecimal digit string that identifies the private endpoint.
	// Read only field.
	Id *string `json:"_id,omitempty"`
	// Unique string that identifies the Azure private endpoint's network interface that someone added to this private endpoint service.
	// Read only field.
	CloudProviderEndpointId *string `json:"cloudProviderEndpointId,omitempty"`
	// Human-readable comment associated with the private endpoint.
	// Read only field.
	Comment *string `json:"comment,omitempty"`
	// Unique string that identifies the Azure private endpoint service. MongoDB Cloud returns null while it creates the endpoint service.
	// Read only field.
	EndpointServiceName *string `json:"endpointServiceName,omitempty"`
	// Human-readable error message that indicates error condition associated with establishing the private endpoint connection.
	// Read only field.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Human-readable label that indicates the current operating status of the private endpoint.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Human-readable label that identifies the cloud service provider.
	// Read only field.
	ProviderName *string `json:"providerName,omitempty"`
	// IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	// Read only field.
	PrivateEndpointIpAddress *string `json:"privateEndpointIpAddress,omitempty"`
	// Root-relative path that identifies the Azure Private Link Service that MongoDB Cloud manages. MongoDB Cloud returns null while it creates the endpoint service.
	// Read only field.
	PrivateLinkServiceResourceId *string `json:"privateLinkServiceResourceId,omitempty"`
}

// NewServerlessTenantEndpoint instantiates a new ServerlessTenantEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessTenantEndpoint() *ServerlessTenantEndpoint {
	this := ServerlessTenantEndpoint{}
	return &this
}

// NewServerlessTenantEndpointWithDefaults instantiates a new ServerlessTenantEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessTenantEndpointWithDefaults() *ServerlessTenantEndpoint {
	this := ServerlessTenantEndpoint{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ServerlessTenantEndpoint) SetId(v string) {
	o.Id = &v
}

// GetCloudProviderEndpointId returns the CloudProviderEndpointId field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetCloudProviderEndpointId() string {
	if o == nil || IsNil(o.CloudProviderEndpointId) {
		var ret string
		return ret
	}
	return *o.CloudProviderEndpointId
}

// GetCloudProviderEndpointIdOk returns a tuple with the CloudProviderEndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetCloudProviderEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProviderEndpointId) {
		return nil, false
	}

	return o.CloudProviderEndpointId, true
}

// HasCloudProviderEndpointId returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasCloudProviderEndpointId() bool {
	if o != nil && !IsNil(o.CloudProviderEndpointId) {
		return true
	}

	return false
}

// SetCloudProviderEndpointId gets a reference to the given string and assigns it to the CloudProviderEndpointId field.
func (o *ServerlessTenantEndpoint) SetCloudProviderEndpointId(v string) {
	o.CloudProviderEndpointId = &v
}

// GetComment returns the Comment field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetComment() string {
	if o == nil || IsNil(o.Comment) {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetCommentOk() (*string, bool) {
	if o == nil || IsNil(o.Comment) {
		return nil, false
	}

	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasComment() bool {
	if o != nil && !IsNil(o.Comment) {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *ServerlessTenantEndpoint) SetComment(v string) {
	o.Comment = &v
}

// GetEndpointServiceName returns the EndpointServiceName field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetEndpointServiceName() string {
	if o == nil || IsNil(o.EndpointServiceName) {
		var ret string
		return ret
	}
	return *o.EndpointServiceName
}

// GetEndpointServiceNameOk returns a tuple with the EndpointServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetEndpointServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointServiceName) {
		return nil, false
	}

	return o.EndpointServiceName, true
}

// HasEndpointServiceName returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasEndpointServiceName() bool {
	if o != nil && !IsNil(o.EndpointServiceName) {
		return true
	}

	return false
}

// SetEndpointServiceName gets a reference to the given string and assigns it to the EndpointServiceName field.
func (o *ServerlessTenantEndpoint) SetEndpointServiceName(v string) {
	o.EndpointServiceName = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *ServerlessTenantEndpoint) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *ServerlessTenantEndpoint) SetStatus(v string) {
	o.Status = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *ServerlessTenantEndpoint) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetPrivateEndpointIpAddress returns the PrivateEndpointIpAddress field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetPrivateEndpointIpAddress() string {
	if o == nil || IsNil(o.PrivateEndpointIpAddress) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointIpAddress
}

// GetPrivateEndpointIpAddressOk returns a tuple with the PrivateEndpointIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetPrivateEndpointIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointIpAddress) {
		return nil, false
	}

	return o.PrivateEndpointIpAddress, true
}

// HasPrivateEndpointIpAddress returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasPrivateEndpointIpAddress() bool {
	if o != nil && !IsNil(o.PrivateEndpointIpAddress) {
		return true
	}

	return false
}

// SetPrivateEndpointIpAddress gets a reference to the given string and assigns it to the PrivateEndpointIpAddress field.
func (o *ServerlessTenantEndpoint) SetPrivateEndpointIpAddress(v string) {
	o.PrivateEndpointIpAddress = &v
}

// GetPrivateLinkServiceResourceId returns the PrivateLinkServiceResourceId field value if set, zero value otherwise
func (o *ServerlessTenantEndpoint) GetPrivateLinkServiceResourceId() string {
	if o == nil || IsNil(o.PrivateLinkServiceResourceId) {
		var ret string
		return ret
	}
	return *o.PrivateLinkServiceResourceId
}

// GetPrivateLinkServiceResourceIdOk returns a tuple with the PrivateLinkServiceResourceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpoint) GetPrivateLinkServiceResourceIdOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceResourceId) {
		return nil, false
	}

	return o.PrivateLinkServiceResourceId, true
}

// HasPrivateLinkServiceResourceId returns a boolean if a field has been set.
func (o *ServerlessTenantEndpoint) HasPrivateLinkServiceResourceId() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceResourceId) {
		return true
	}

	return false
}

// SetPrivateLinkServiceResourceId gets a reference to the given string and assigns it to the PrivateLinkServiceResourceId field.
func (o *ServerlessTenantEndpoint) SetPrivateLinkServiceResourceId(v string) {
	o.PrivateLinkServiceResourceId = &v
}
