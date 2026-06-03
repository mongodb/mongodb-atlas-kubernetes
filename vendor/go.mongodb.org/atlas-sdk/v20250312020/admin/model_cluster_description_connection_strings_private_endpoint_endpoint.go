// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterDescriptionConnectionStringsPrivateEndpointEndpoint Details of a private endpoint deployed for this cluster.
type ClusterDescriptionConnectionStringsPrivateEndpointEndpoint struct {
	// Unique string that the cloud provider uses to identify the private endpoint.
	// Read only field.
	EndpointId *string `json:"endpointId,omitempty"`
	// Cloud provider in which MongoDB Cloud deploys the private endpoint.
	// Read only field.
	ProviderName *string `json:"providerName,omitempty"`
	// Region where the private endpoint is deployed.
	// Read only field.
	Region *string `json:"region,omitempty"`
}

// NewClusterDescriptionConnectionStringsPrivateEndpointEndpoint instantiates a new ClusterDescriptionConnectionStringsPrivateEndpointEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterDescriptionConnectionStringsPrivateEndpointEndpoint() *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint {
	this := ClusterDescriptionConnectionStringsPrivateEndpointEndpoint{}
	return &this
}

// NewClusterDescriptionConnectionStringsPrivateEndpointEndpointWithDefaults instantiates a new ClusterDescriptionConnectionStringsPrivateEndpointEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterDescriptionConnectionStringsPrivateEndpointEndpointWithDefaults() *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint {
	this := ClusterDescriptionConnectionStringsPrivateEndpointEndpoint{}
	return &this
}

// GetEndpointId returns the EndpointId field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetEndpointId() string {
	if o == nil || IsNil(o.EndpointId) {
		var ret string
		return ret
	}
	return *o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetEndpointIdOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointId) {
		return nil, false
	}

	return o.EndpointId, true
}

// HasEndpointId returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) HasEndpointId() bool {
	if o != nil && !IsNil(o.EndpointId) {
		return true
	}

	return false
}

// SetEndpointId gets a reference to the given string and assigns it to the EndpointId field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) SetEndpointId(v string) {
	o.EndpointId = &v
}

// GetProviderName returns the ProviderName field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetProviderName() string {
	if o == nil || IsNil(o.ProviderName) {
		var ret string
		return ret
	}
	return *o.ProviderName
}

// GetProviderNameOk returns a tuple with the ProviderName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetProviderNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderName) {
		return nil, false
	}

	return o.ProviderName, true
}

// HasProviderName returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) HasProviderName() bool {
	if o != nil && !IsNil(o.ProviderName) {
		return true
	}

	return false
}

// SetProviderName gets a reference to the given string and assigns it to the ProviderName field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) SetProviderName(v string) {
	o.ProviderName = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *ClusterDescriptionConnectionStringsPrivateEndpointEndpoint) SetRegion(v string) {
	o.Region = &v
}
