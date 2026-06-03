// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsDataProcessRegion Information about the cloud provider region in which MongoDB Cloud processes the stream.
type StreamsDataProcessRegion struct {
	// Human-readable label that identifies the cloud provider.
	CloudProvider string `json:"cloudProvider"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Name of the cloud provider region hosting Atlas Stream Processing.
	Region string `json:"region"`
}

// NewStreamsDataProcessRegion instantiates a new StreamsDataProcessRegion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsDataProcessRegion(cloudProvider string, region string) *StreamsDataProcessRegion {
	this := StreamsDataProcessRegion{}
	this.CloudProvider = cloudProvider
	this.Region = region
	return &this
}

// NewStreamsDataProcessRegionWithDefaults instantiates a new StreamsDataProcessRegion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsDataProcessRegionWithDefaults() *StreamsDataProcessRegion {
	this := StreamsDataProcessRegion{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *StreamsDataProcessRegion) GetCloudProvider() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *StreamsDataProcessRegion) GetCloudProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *StreamsDataProcessRegion) SetCloudProvider(v string) {
	o.CloudProvider = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsDataProcessRegion) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsDataProcessRegion) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsDataProcessRegion) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsDataProcessRegion) SetLinks(v []Link) {
	o.Links = &v
}

// GetRegion returns the Region field value
func (o *StreamsDataProcessRegion) GetRegion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Region
}

// GetRegionOk returns a tuple with the Region field value
// and a boolean to check if the value has been set.
func (o *StreamsDataProcessRegion) GetRegionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Region, true
}

// SetRegion sets field value
func (o *StreamsDataProcessRegion) SetRegion(v string) {
	o.Region = v
}
