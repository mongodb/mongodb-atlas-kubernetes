// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CreateDataProcessRegion Settings to configure the region where you wish to store your archived data.
type CreateDataProcessRegion struct {
	// Human-readable label that identifies the Cloud service provider where you wish to store your archived data. `AZURE` or `GCP` may be selected only if it is the Cloud service provider for the cluster and no archives for any other cloud provider have been created for the cluster.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Human-readable label that identifies the geographic location of the region where you wish to store your archived data.
	Region *string `json:"region,omitempty"`
}

// NewCreateDataProcessRegion instantiates a new CreateDataProcessRegion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateDataProcessRegion() *CreateDataProcessRegion {
	this := CreateDataProcessRegion{}
	return &this
}

// NewCreateDataProcessRegionWithDefaults instantiates a new CreateDataProcessRegion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateDataProcessRegionWithDefaults() *CreateDataProcessRegion {
	this := CreateDataProcessRegion{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *CreateDataProcessRegion) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateDataProcessRegion) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *CreateDataProcessRegion) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *CreateDataProcessRegion) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetRegion returns the Region field value if set, zero value otherwise
func (o *CreateDataProcessRegion) GetRegion() string {
	if o == nil || IsNil(o.Region) {
		var ret string
		return ret
	}
	return *o.Region
}

// GetRegionOk returns a tuple with the Region field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateDataProcessRegion) GetRegionOk() (*string, bool) {
	if o == nil || IsNil(o.Region) {
		return nil, false
	}

	return o.Region, true
}

// HasRegion returns a boolean if a field has been set.
func (o *CreateDataProcessRegion) HasRegion() bool {
	if o != nil && !IsNil(o.Region) {
		return true
	}

	return false
}

// SetRegion gets a reference to the given string and assigns it to the Region field.
func (o *CreateDataProcessRegion) SetRegion(v string) {
	o.Region = &v
}
