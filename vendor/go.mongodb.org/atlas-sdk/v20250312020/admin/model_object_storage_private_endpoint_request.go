// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ObjectStoragePrivateEndpointRequest struct for ObjectStoragePrivateEndpointRequest
type ObjectStoragePrivateEndpointRequest struct {
	// Human-readable label that identifies the cloud provider.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// Cloud provider region in which the Object Storage private endpoint is located.
	RegionName *string `json:"regionName,omitempty"`
}

// NewObjectStoragePrivateEndpointRequest instantiates a new ObjectStoragePrivateEndpointRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewObjectStoragePrivateEndpointRequest() *ObjectStoragePrivateEndpointRequest {
	this := ObjectStoragePrivateEndpointRequest{}
	return &this
}

// NewObjectStoragePrivateEndpointRequestWithDefaults instantiates a new ObjectStoragePrivateEndpointRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewObjectStoragePrivateEndpointRequestWithDefaults() *ObjectStoragePrivateEndpointRequest {
	this := ObjectStoragePrivateEndpointRequest{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *ObjectStoragePrivateEndpointRequest) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ObjectStoragePrivateEndpointRequest) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *ObjectStoragePrivateEndpointRequest) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *ObjectStoragePrivateEndpointRequest) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *ObjectStoragePrivateEndpointRequest) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ObjectStoragePrivateEndpointRequest) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *ObjectStoragePrivateEndpointRequest) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *ObjectStoragePrivateEndpointRequest) SetRegionName(v string) {
	o.RegionName = &v
}
