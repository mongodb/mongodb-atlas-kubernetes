// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterAutoScalingSettings Range of instance sizes to which your cluster can scale.
type ClusterAutoScalingSettings struct {
	Compute *ClusterComputeAutoScaling `json:"compute,omitempty"`
	// Flag that indicates whether someone enabled disk auto-scaling for this cluster.
	DiskGBEnabled *bool `json:"diskGBEnabled,omitempty"`
}

// NewClusterAutoScalingSettings instantiates a new ClusterAutoScalingSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterAutoScalingSettings() *ClusterAutoScalingSettings {
	this := ClusterAutoScalingSettings{}
	var diskGBEnabled bool = false
	this.DiskGBEnabled = &diskGBEnabled
	return &this
}

// NewClusterAutoScalingSettingsWithDefaults instantiates a new ClusterAutoScalingSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterAutoScalingSettingsWithDefaults() *ClusterAutoScalingSettings {
	this := ClusterAutoScalingSettings{}
	var diskGBEnabled bool = false
	this.DiskGBEnabled = &diskGBEnabled
	return &this
}

// GetCompute returns the Compute field value if set, zero value otherwise
func (o *ClusterAutoScalingSettings) GetCompute() ClusterComputeAutoScaling {
	if o == nil || IsNil(o.Compute) {
		var ret ClusterComputeAutoScaling
		return ret
	}
	return *o.Compute
}

// GetComputeOk returns a tuple with the Compute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterAutoScalingSettings) GetComputeOk() (*ClusterComputeAutoScaling, bool) {
	if o == nil || IsNil(o.Compute) {
		return nil, false
	}

	return o.Compute, true
}

// HasCompute returns a boolean if a field has been set.
func (o *ClusterAutoScalingSettings) HasCompute() bool {
	if o != nil && !IsNil(o.Compute) {
		return true
	}

	return false
}

// SetCompute gets a reference to the given ClusterComputeAutoScaling and assigns it to the Compute field.
func (o *ClusterAutoScalingSettings) SetCompute(v ClusterComputeAutoScaling) {
	o.Compute = &v
}

// GetDiskGBEnabled returns the DiskGBEnabled field value if set, zero value otherwise
func (o *ClusterAutoScalingSettings) GetDiskGBEnabled() bool {
	if o == nil || IsNil(o.DiskGBEnabled) {
		var ret bool
		return ret
	}
	return *o.DiskGBEnabled
}

// GetDiskGBEnabledOk returns a tuple with the DiskGBEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterAutoScalingSettings) GetDiskGBEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.DiskGBEnabled) {
		return nil, false
	}

	return o.DiskGBEnabled, true
}

// HasDiskGBEnabled returns a boolean if a field has been set.
func (o *ClusterAutoScalingSettings) HasDiskGBEnabled() bool {
	if o != nil && !IsNil(o.DiskGBEnabled) {
		return true
	}

	return false
}

// SetDiskGBEnabled gets a reference to the given bool and assigns it to the DiskGBEnabled field.
func (o *ClusterAutoScalingSettings) SetDiskGBEnabled(v bool) {
	o.DiskGBEnabled = &v
}
