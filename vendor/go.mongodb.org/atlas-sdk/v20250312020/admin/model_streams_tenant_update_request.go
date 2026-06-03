// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsTenantUpdateRequest Details to update a stream tenant.
type StreamsTenantUpdateRequest struct {
	// Human-readable label that identifies the cloud provider.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Name of the cloud provider region hosting Atlas Stream Processing.
	Region       *string       `json:"region,omitempty"`
	StreamConfig *StreamConfig `json:"streamConfig,omitempty"`
}

// NewStreamsTenantUpdateRequest instantiates a new StreamsTenantUpdateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsTenantUpdateRequest() *StreamsTenantUpdateRequest {
	this := StreamsTenantUpdateRequest{}
	return &this
}

// NewStreamsTenantUpdateRequestWithDefaults instantiates a new StreamsTenantUpdateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsTenantUpdateRequestWithDefaults() *StreamsTenantUpdateRequest {
	this := StreamsTenantUpdateRequest{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *StreamsTenantUpdateRequest) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenantUpdateRequest) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *StreamsTenantUpdateRequest) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *StreamsTenantUpdateRequest) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsTenantUpdateRequest) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenantUpdateRequest) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsTenantUpdateRequest) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsTenantUpdateRequest) SetLinks(v []Link) {
	o.Links = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *StreamsTenantUpdateRequest) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenantUpdateRequest) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *StreamsTenantUpdateRequest) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *StreamsTenantUpdateRequest) SetRegion(v string) {
	o.Region = &v
}

// GetStreamConfig returns the StreamConfig field value if set, zero value otherwise
func (o *StreamsTenantUpdateRequest) GetStreamConfig() StreamConfig {
	if o == nil || IsNil(o.StreamConfig) {
		var ret StreamConfig
		return ret
	}
	return *o.StreamConfig
}

// GetStreamConfigOk returns a tuple with the StreamConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsTenantUpdateRequest) GetStreamConfigOk() (*StreamConfig, bool) {
	if o == nil || IsNil(o.StreamConfig) {
		return nil, false
	}

	return o.StreamConfig, true
}

// HasStreamConfig returns a boolean if a field has been set.
func (o *StreamsTenantUpdateRequest) HasStreamConfig() bool {
	if o != nil && !IsNil(o.StreamConfig) {
		return true
	}

	return false
}

// SetStreamConfig gets a reference to the given StreamConfig and assigns it to the StreamConfig field.
func (o *StreamsTenantUpdateRequest) SetStreamConfig(v StreamConfig) {
	o.StreamConfig = &v
}
