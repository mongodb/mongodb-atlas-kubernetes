// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataProcessRegion Settings to configure the region where you wish to store your archived data.
type DataProcessRegion struct {
	// Human-readable label that identifies the Cloud service provider where you store your archived data.
	// Read only field.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Human-readable label that identifies the geographic location of the region where you store your archived data.
	// Read only field.
	Region *string `json:"region,omitempty"`
}

// NewDataProcessRegion instantiates a new DataProcessRegion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataProcessRegion() *DataProcessRegion {
	this := DataProcessRegion{}
	return &this
}

// NewDataProcessRegionWithDefaults instantiates a new DataProcessRegion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataProcessRegionWithDefaults() *DataProcessRegion {
	this := DataProcessRegion{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *DataProcessRegion) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProcessRegion) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *DataProcessRegion) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *DataProcessRegion) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *DataProcessRegion) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProcessRegion) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *DataProcessRegion) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *DataProcessRegion) SetRegion(v string) {
	o.Region = &v
}
