// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AtlasClusterOutageSimulationOutageFilter struct for AtlasClusterOutageSimulationOutageFilter
type AtlasClusterOutageSimulationOutageFilter struct {
	// The cloud provider of the region that undergoes the outage simulation.
	CloudProvider *string `json:"cloudProvider,omitempty"`
	// The name of the region to undergo an outage simulation.
	RegionName *string `json:"regionName,omitempty"`
	// The type of cluster outage to simulate. `REGION` simulates a cluster outage for a region.
	Type *string `json:"type,omitempty"`
}

// NewAtlasClusterOutageSimulationOutageFilter instantiates a new AtlasClusterOutageSimulationOutageFilter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAtlasClusterOutageSimulationOutageFilter() *AtlasClusterOutageSimulationOutageFilter {
	this := AtlasClusterOutageSimulationOutageFilter{}
	return &this
}

// NewAtlasClusterOutageSimulationOutageFilterWithDefaults instantiates a new AtlasClusterOutageSimulationOutageFilter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAtlasClusterOutageSimulationOutageFilterWithDefaults() *AtlasClusterOutageSimulationOutageFilter {
	this := AtlasClusterOutageSimulationOutageFilter{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise
func (o *AtlasClusterOutageSimulationOutageFilter) GetCloudProvider() string {
	if o == nil || IsNil(o.CloudProvider) {
		var ret string
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) GetCloudProviderOk() (*string, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}

	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given string and assigns it to the CloudProvider field.
func (o *AtlasClusterOutageSimulationOutageFilter) SetCloudProvider(v string) {
	o.CloudProvider = &v
}

// GetRegionName returns the RegionName field value if set, zero value otherwise
func (o *AtlasClusterOutageSimulationOutageFilter) GetRegionName() string {
	if o == nil || IsNil(o.RegionName) {
		var ret string
		return ret
	}
	return *o.RegionName
}

// GetRegionNameOk returns a tuple with the RegionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) GetRegionNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionName) {
		return nil, false
	}

	return o.RegionName, true
}

// HasRegionName returns a boolean if a field has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) HasRegionName() bool {
	if o != nil && !IsNil(o.RegionName) {
		return true
	}

	return false
}

// SetRegionName gets a reference to the given string and assigns it to the RegionName field.
func (o *AtlasClusterOutageSimulationOutageFilter) SetRegionName(v string) {
	o.RegionName = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *AtlasClusterOutageSimulationOutageFilter) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *AtlasClusterOutageSimulationOutageFilter) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *AtlasClusterOutageSimulationOutageFilter) SetType(v string) {
	o.Type = &v
}
