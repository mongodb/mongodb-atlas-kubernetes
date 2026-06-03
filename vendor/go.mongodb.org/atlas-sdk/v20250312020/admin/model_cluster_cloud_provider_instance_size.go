// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterCloudProviderInstanceSize struct for ClusterCloudProviderInstanceSize
type ClusterCloudProviderInstanceSize struct {
	// List of regions that this cloud provider supports for this instance size.
	// Read only field.
	AvailableRegions *[]AvailableCloudProviderRegion `json:"availableRegions,omitempty"`
	// Human-readable label that identifies the instance size or cluster tier.
	// Read only field.
	Name *string `json:"name,omitempty"`
}

// NewClusterCloudProviderInstanceSize instantiates a new ClusterCloudProviderInstanceSize object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterCloudProviderInstanceSize() *ClusterCloudProviderInstanceSize {
	this := ClusterCloudProviderInstanceSize{}
	return &this
}

// NewClusterCloudProviderInstanceSizeWithDefaults instantiates a new ClusterCloudProviderInstanceSize object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterCloudProviderInstanceSizeWithDefaults() *ClusterCloudProviderInstanceSize {
	this := ClusterCloudProviderInstanceSize{}
	return &this
}

// GetAvailableRegions returns the AvailableRegions field value if set, zero value otherwise
func (o *ClusterCloudProviderInstanceSize) GetAvailableRegions() []AvailableCloudProviderRegion {
	if o == nil || IsNil(o.AvailableRegions) {
		var ret []AvailableCloudProviderRegion
		return ret
	}
	return *o.AvailableRegions
}

// GetAvailableRegionsOk returns a tuple with the AvailableRegions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterCloudProviderInstanceSize) GetAvailableRegionsOk() (*[]AvailableCloudProviderRegion, bool) {
	if o == nil || IsNil(o.AvailableRegions) {
		return nil, false
	}

	return o.AvailableRegions, true
}

// HasAvailableRegions returns a boolean if a field has been set.
func (o *ClusterCloudProviderInstanceSize) HasAvailableRegions() bool {
	if o != nil && !IsNil(o.AvailableRegions) {
		return true
	}

	return false
}

// SetAvailableRegions gets a reference to the given []AvailableCloudProviderRegion and assigns it to the AvailableRegions field.
func (o *ClusterCloudProviderInstanceSize) SetAvailableRegions(v []AvailableCloudProviderRegion) {
	o.AvailableRegions = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *ClusterCloudProviderInstanceSize) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterCloudProviderInstanceSize) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ClusterCloudProviderInstanceSize) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ClusterCloudProviderInstanceSize) SetName(v string) {
	o.Name = &v
}
