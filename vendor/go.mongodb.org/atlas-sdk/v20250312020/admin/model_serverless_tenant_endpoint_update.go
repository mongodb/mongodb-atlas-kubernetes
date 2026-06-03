// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessTenantEndpointUpdate Update view for a serverless tenant endpoint.
type ServerlessTenantEndpointUpdate struct {
	// Human-readable comment associated with the private endpoint.
	// Write only field.
	Comment *string `json:"comment,omitempty"`
	// Human-readable label that identifies the cloud provider of the tenant endpoint.
	// Write only field.
	ProviderName string `json:"providerName"`
	// Unique string that identifies the Azure private endpoint's network interface for this private endpoint service.
	// Write only field.
	CloudProviderEndpointId *string `json:"cloudProviderEndpointId,omitempty"`
	// IPv4 address of the private endpoint in your Azure VNet that someone added to this private endpoint service.
	// Write only field.
	PrivateEndpointIpAddress *string `json:"privateEndpointIpAddress,omitempty"`
}

// NewServerlessTenantEndpointUpdate instantiates a new ServerlessTenantEndpointUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessTenantEndpointUpdate(providerName string) *ServerlessTenantEndpointUpdate {
	this := ServerlessTenantEndpointUpdate{}
	this.ProviderName = providerName
	return &this
}

// NewServerlessTenantEndpointUpdateWithDefaults instantiates a new ServerlessTenantEndpointUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessTenantEndpointUpdateWithDefaults() *ServerlessTenantEndpointUpdate {
	this := ServerlessTenantEndpointUpdate{}
	return &this
}

// GetComment returns the Comment field value if set, zero value otherwise
func (o *ServerlessTenantEndpointUpdate) GetComment() string {
	if o == nil || IsNil(o.Comment) {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpointUpdate) GetCommentOk() (*string, bool) {
	if o == nil || IsNil(o.Comment) {
		return nil, false
	}

	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *ServerlessTenantEndpointUpdate) HasComment() bool {
	if o != nil && !IsNil(o.Comment) {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *ServerlessTenantEndpointUpdate) SetComment(v string) {
	o.Comment = &v
}

// GetProviderName returns the ProviderName field value
func (o *ServerlessTenantEndpointUpdate) GetProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpointUpdate) GetProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderName, true
}

// SetProviderName sets field value
func (o *ServerlessTenantEndpointUpdate) SetProviderName(v string) {
	o.ProviderName = v
}

// GetCloudProviderEndpointId returns the CloudProviderEndpointId field value if set, zero value otherwise
func (o *ServerlessTenantEndpointUpdate) GetCloudProviderEndpointId() string {
	if o == nil || IsNil(o.CloudProviderEndpointId) {
		var ret string
		return ret
	}
	return *o.CloudProviderEndpointId
}

// GetCloudProviderEndpointIdOk returns a tuple with the CloudProviderEndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpointUpdate) GetCloudProviderEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProviderEndpointId) {
		return nil, false
	}

	return o.CloudProviderEndpointId, true
}

// HasCloudProviderEndpointId returns a boolean if a field has been set.
func (o *ServerlessTenantEndpointUpdate) HasCloudProviderEndpointId() bool {
	if o != nil && !IsNil(o.CloudProviderEndpointId) {
		return true
	}

	return false
}

// SetCloudProviderEndpointId gets a reference to the given string and assigns it to the CloudProviderEndpointId field.
func (o *ServerlessTenantEndpointUpdate) SetCloudProviderEndpointId(v string) {
	o.CloudProviderEndpointId = &v
}

// GetPrivateEndpointIpAddress returns the PrivateEndpointIpAddress field value if set, zero value otherwise
func (o *ServerlessTenantEndpointUpdate) GetPrivateEndpointIpAddress() string {
	if o == nil || IsNil(o.PrivateEndpointIpAddress) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointIpAddress
}

// GetPrivateEndpointIpAddressOk returns a tuple with the PrivateEndpointIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantEndpointUpdate) GetPrivateEndpointIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointIpAddress) {
		return nil, false
	}

	return o.PrivateEndpointIpAddress, true
}

// HasPrivateEndpointIpAddress returns a boolean if a field has been set.
func (o *ServerlessTenantEndpointUpdate) HasPrivateEndpointIpAddress() bool {
	if o != nil && !IsNil(o.PrivateEndpointIpAddress) {
		return true
	}

	return false
}

// SetPrivateEndpointIpAddress gets a reference to the given string and assigns it to the PrivateEndpointIpAddress field.
func (o *ServerlessTenantEndpointUpdate) SetPrivateEndpointIpAddress(v string) {
	o.PrivateEndpointIpAddress = &v
}
