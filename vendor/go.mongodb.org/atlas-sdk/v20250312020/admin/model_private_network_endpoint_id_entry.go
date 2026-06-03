// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PrivateNetworkEndpointIdEntry struct for PrivateNetworkEndpointIdEntry
type PrivateNetworkEndpointIdEntry struct {
	// Link ID that identifies the Azure private endpoint connection.
	AzureLinkId *string `json:"azureLinkId,omitempty"`
	// Human-readable string to associate with this private endpoint.
	Comment *string `json:"comment,omitempty"`
	// Human-readable label to identify customer's VPC endpoint DNS name. If defined, you must also specify a value for `region`.
	CustomerEndpointDNSName *string `json:"customerEndpointDNSName,omitempty"`
	// IP address used to connect to the Azure private endpoint.
	CustomerEndpointIPAddress *string `json:"customerEndpointIPAddress,omitempty"`
	// Unique string that identifies the private endpoint. For AWS, this is a 22-character alphanumeric string in the format `vpce-<17 hex characters>`. For Azure, this is the full resource ID in the format `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/privateEndpoints/{endpointName}`.
	EndpointId string `json:"endpointId"`
	// Error message describing a failure approving the private endpoint request.
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// Human-readable label that identifies the cloud service provider. Atlas Data Lake supports Amazon Web Services and Azure.
	Provider *string `json:"provider,omitempty"`
	// Human-readable label to identify the region of customer's VPC endpoint. If defined, you must also specify a value for `customerEndpointDNSName`.
	Region *string `json:"region,omitempty"`
	// Status of the private endpoint connection request.
	Status *string `json:"status,omitempty"`
	// Human-readable label that identifies the resource type associated with this private endpoint.
	Type *string `json:"type,omitempty"`
}

// NewPrivateNetworkEndpointIdEntry instantiates a new PrivateNetworkEndpointIdEntry object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateNetworkEndpointIdEntry(endpointId string) *PrivateNetworkEndpointIdEntry {
	this := PrivateNetworkEndpointIdEntry{}
	this.EndpointId = endpointId
	var provider string = "AWS"
	this.Provider = &provider
	var type_ string = "DATA_LAKE"
	this.Type = &type_
	return &this
}

// NewPrivateNetworkEndpointIdEntryWithDefaults instantiates a new PrivateNetworkEndpointIdEntry object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateNetworkEndpointIdEntryWithDefaults() *PrivateNetworkEndpointIdEntry {
	this := PrivateNetworkEndpointIdEntry{}
	var provider string = "AWS"
	this.Provider = &provider
	var type_ string = "DATA_LAKE"
	this.Type = &type_
	return &this
}

// GetAzureLinkId returns the AzureLinkId field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetAzureLinkId() string {
	if o == nil || IsNil(o.AzureLinkId) {
		var ret string
		return ret
	}
	return *o.AzureLinkId
}

// GetAzureLinkIdOk returns a tuple with the AzureLinkId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetAzureLinkIdOk() (*string, bool) {
	if o == nil || IsNil(o.AzureLinkId) {
		return nil, false
	}

	return o.AzureLinkId, true
}

// HasAzureLinkId returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasAzureLinkId() bool {
	if o != nil && !IsNil(o.AzureLinkId) {
		return true
	}

	return false
}

// SetAzureLinkId gets a reference to the given string and assigns it to the AzureLinkId field.
func (o *PrivateNetworkEndpointIdEntry) SetAzureLinkId(v string) {
	o.AzureLinkId = &v
}

// GetComment returns the Comment field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetComment() string {
	if o == nil || IsNil(o.Comment) {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetCommentOk() (*string, bool) {
	if o == nil || IsNil(o.Comment) {
		return nil, false
	}

	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasComment() bool {
	if o != nil && !IsNil(o.Comment) {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *PrivateNetworkEndpointIdEntry) SetComment(v string) {
	o.Comment = &v
}

// GetCustomerEndpointDNSName returns the CustomerEndpointDNSName field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetCustomerEndpointDNSName() string {
	if o == nil || IsNil(o.CustomerEndpointDNSName) {
		var ret string
		return ret
	}
	return *o.CustomerEndpointDNSName
}

// GetCustomerEndpointDNSNameOk returns a tuple with the CustomerEndpointDNSName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetCustomerEndpointDNSNameOk() (*string, bool) {
	if o == nil || IsNil(o.CustomerEndpointDNSName) {
		return nil, false
	}

	return o.CustomerEndpointDNSName, true
}

// HasCustomerEndpointDNSName returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasCustomerEndpointDNSName() bool {
	if o != nil && !IsNil(o.CustomerEndpointDNSName) {
		return true
	}

	return false
}

// SetCustomerEndpointDNSName gets a reference to the given string and assigns it to the CustomerEndpointDNSName field.
func (o *PrivateNetworkEndpointIdEntry) SetCustomerEndpointDNSName(v string) {
	o.CustomerEndpointDNSName = &v
}

// GetCustomerEndpointIPAddress returns the CustomerEndpointIPAddress field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetCustomerEndpointIPAddress() string {
	if o == nil || IsNil(o.CustomerEndpointIPAddress) {
		var ret string
		return ret
	}
	return *o.CustomerEndpointIPAddress
}

// GetCustomerEndpointIPAddressOk returns a tuple with the CustomerEndpointIPAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetCustomerEndpointIPAddressOk() (*string, bool) {
	if o == nil || IsNil(o.CustomerEndpointIPAddress) {
		return nil, false
	}

	return o.CustomerEndpointIPAddress, true
}

// HasCustomerEndpointIPAddress returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasCustomerEndpointIPAddress() bool {
	if o != nil && !IsNil(o.CustomerEndpointIPAddress) {
		return true
	}

	return false
}

// SetCustomerEndpointIPAddress gets a reference to the given string and assigns it to the CustomerEndpointIPAddress field.
func (o *PrivateNetworkEndpointIdEntry) SetCustomerEndpointIPAddress(v string) {
	o.CustomerEndpointIPAddress = &v
}

// GetEndpointId returns the EndpointId field value
func (o *PrivateNetworkEndpointIdEntry) GetEndpointId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetEndpointIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndpointId, true
}

// SetEndpointId sets field value
func (o *PrivateNetworkEndpointIdEntry) SetEndpointId(v string) {
	o.EndpointId = v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetErrorMessage() string {
	if o == nil || IsNil(o.ErrorMessage) {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetErrorMessageOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorMessage) {
		return nil, false
	}

	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasErrorMessage() bool {
	if o != nil && !IsNil(o.ErrorMessage) {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *PrivateNetworkEndpointIdEntry) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

// GetProvider returns the Provider field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetProvider() string {
	if o == nil || IsNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetProviderOk() (*string, bool) {
	if o == nil || IsNil(o.Provider) {
		return nil, false
	}

	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasProvider() bool {
	if o != nil && !IsNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *PrivateNetworkEndpointIdEntry) SetProvider(v string) {
	o.Provider = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *PrivateNetworkEndpointIdEntry) SetRegion(v string) {
	o.Region = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *PrivateNetworkEndpointIdEntry) SetStatus(v string) {
	o.Status = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *PrivateNetworkEndpointIdEntry) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateNetworkEndpointIdEntry) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *PrivateNetworkEndpointIdEntry) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *PrivateNetworkEndpointIdEntry) SetType(v string) {
	o.Type = &v
}
