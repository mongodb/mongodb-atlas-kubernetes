// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudProviderEndpointServiceRequest struct for CloudProviderEndpointServiceRequest
type CloudProviderEndpointServiceRequest struct {
	// Flag that indicates whether this endpoint service uses PSC port-mapping. This is only applicable for GCP Private Endpoint Services.
	// Write only field.
	PortMappingEnabled *bool `json:"portMappingEnabled,omitempty"`
	// Human-readable label that identifies the cloud service provider for which you want to create the private endpoint service.
	// Write only field.
	ProviderName string `json:"providerName"`
	// Cloud provider region in which you want to create the private endpoint service. Regions accepted as values differ for [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/), [Google Cloud Platform](https://docs.atlas.mongodb.com/reference/google-gcp/), and [Microsoft Azure](https://docs.atlas.mongodb.com/reference/microsoft-azure/).
	// Write only field.
	Region string `json:"region"`
	// List of regions that the endpoint service supports. Native cross region support is implemented for AWS only.
	// Write only field.
	SupportedRemoteRegions *[]string `json:"supportedRemoteRegions,omitempty"`
}

// NewCloudProviderEndpointServiceRequest instantiates a new CloudProviderEndpointServiceRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudProviderEndpointServiceRequest(providerName string, region string) *CloudProviderEndpointServiceRequest {
	this := CloudProviderEndpointServiceRequest{}
	var portMappingEnabled bool = false
	this.PortMappingEnabled = &portMappingEnabled
	this.ProviderName = providerName
	this.Region = region
	return &this
}

// NewCloudProviderEndpointServiceRequestWithDefaults instantiates a new CloudProviderEndpointServiceRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudProviderEndpointServiceRequestWithDefaults() *CloudProviderEndpointServiceRequest {
	this := CloudProviderEndpointServiceRequest{}
	var portMappingEnabled bool = false
	this.PortMappingEnabled = &portMappingEnabled
	return &this
}

// GetPortMappingEnabled returns the PortMappingEnabled field value if set, zero value otherwise
func (o *CloudProviderEndpointServiceRequest) GetPortMappingEnabled() bool {
	if o == nil || IsNil(o.PortMappingEnabled) {
		var ret bool
		return ret
	}
	return *o.PortMappingEnabled
}

// GetPortMappingEnabledOk returns a tuple with the PortMappingEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderEndpointServiceRequest) GetPortMappingEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.PortMappingEnabled) {
		return nil, false
	}

	return o.PortMappingEnabled, true
}

// HasPortMappingEnabled returns a boolean if a field has been set.
func (o *CloudProviderEndpointServiceRequest) HasPortMappingEnabled() bool {
	if o != nil && !IsNil(o.PortMappingEnabled) {
		return true
	}

	return false
}

// SetPortMappingEnabled gets a reference to the given bool and assigns it to the PortMappingEnabled field.
func (o *CloudProviderEndpointServiceRequest) SetPortMappingEnabled(v bool) {
	o.PortMappingEnabled = &v
}

// GetProviderName returns the ProviderName field value
func (o *CloudProviderEndpointServiceRequest) GetProviderName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value
// and a boolean to check if the value has been set.
func (o *CloudProviderEndpointServiceRequest) GetProviderNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ProviderName, true
}

// SetProviderName sets field value
func (o *CloudProviderEndpointServiceRequest) SetProviderName(v string) {
	o.ProviderName = v
}

// GetRegion returns the Region field value
func (o *CloudProviderEndpointServiceRequest) GetRegion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Region
}

// GetRegionOk returns a tuple with the Region field value
// and a boolean to check if the value has been set.
func (o *CloudProviderEndpointServiceRequest) GetRegionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Region, true
}

// SetRegion sets field value
func (o *CloudProviderEndpointServiceRequest) SetRegion(v string) {
	o.Region = v
}

// GetSupportedRemoteRegions returns the SupportedRemoteRegions field value if set, zero value otherwise
func (o *CloudProviderEndpointServiceRequest) GetSupportedRemoteRegions() []string {
	if o == nil || IsNil(o.SupportedRemoteRegions) {
		var ret []string
		return ret
	}
	return *o.SupportedRemoteRegions
}

// GetSupportedRemoteRegionsOk returns a tuple with the SupportedRemoteRegions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudProviderEndpointServiceRequest) GetSupportedRemoteRegionsOk() (*[]string, bool) {
	if o == nil || IsNil(o.SupportedRemoteRegions) {
		return nil, false
	}

	return o.SupportedRemoteRegions, true
}

// HasSupportedRemoteRegions returns a boolean if a field has been set.
func (o *CloudProviderEndpointServiceRequest) HasSupportedRemoteRegions() bool {
	if o != nil && !IsNil(o.SupportedRemoteRegions) {
		return true
	}

	return false
}

// SetSupportedRemoteRegions gets a reference to the given []string and assigns it to the SupportedRemoteRegions field.
func (o *CloudProviderEndpointServiceRequest) SetSupportedRemoteRegions(v []string) {
	o.SupportedRemoteRegions = &v
}
