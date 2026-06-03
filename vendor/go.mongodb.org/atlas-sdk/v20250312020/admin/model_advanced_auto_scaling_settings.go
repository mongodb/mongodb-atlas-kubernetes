// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AdvancedAutoScalingSettings Options that determine how this cluster handles resource scaling.
type AdvancedAutoScalingSettings struct {
	Compute *AdvancedComputeAutoScaling `json:"compute,omitempty"`
	DiskGB  *DiskGBAutoScaling          `json:"diskGB,omitempty"`
}

// NewAdvancedAutoScalingSettings instantiates a new AdvancedAutoScalingSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAdvancedAutoScalingSettings() *AdvancedAutoScalingSettings {
	this := AdvancedAutoScalingSettings{}
	return &this
}

// NewAdvancedAutoScalingSettingsWithDefaults instantiates a new AdvancedAutoScalingSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAdvancedAutoScalingSettingsWithDefaults() *AdvancedAutoScalingSettings {
	this := AdvancedAutoScalingSettings{}
	return &this
}

// GetCompute returns the Compute field value if set, zero value otherwise
func (o *AdvancedAutoScalingSettings) GetCompute() AdvancedComputeAutoScaling {
	if o == nil || IsNil(o.Compute) {
		var ret AdvancedComputeAutoScaling
		return ret
	}
	return *o.Compute
}

// GetComputeOk returns a tuple with the Compute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedAutoScalingSettings) GetComputeOk() (*AdvancedComputeAutoScaling, bool) {
	if o == nil || IsNil(o.Compute) {
		return nil, false
	}

	return o.Compute, true
}

// HasCompute returns a boolean if a field has been set.
func (o *AdvancedAutoScalingSettings) HasCompute() bool {
	if o != nil && !IsNil(o.Compute) {
		return true
	}

	return false
}

// SetCompute gets a reference to the given AdvancedComputeAutoScaling and assigns it to the Compute field.
func (o *AdvancedAutoScalingSettings) SetCompute(v AdvancedComputeAutoScaling) {
	o.Compute = &v
}

// GetDiskGB returns the DiskGB field value if set, zero value otherwise
func (o *AdvancedAutoScalingSettings) GetDiskGB() DiskGBAutoScaling {
	if o == nil || IsNil(o.DiskGB) {
		var ret DiskGBAutoScaling
		return ret
	}
	return *o.DiskGB
}

// GetDiskGBOk returns a tuple with the DiskGB field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedAutoScalingSettings) GetDiskGBOk() (*DiskGBAutoScaling, bool) {
	if o == nil || IsNil(o.DiskGB) {
		return nil, false
	}

	return o.DiskGB, true
}

// HasDiskGB returns a boolean if a field has been set.
func (o *AdvancedAutoScalingSettings) HasDiskGB() bool {
	if o != nil && !IsNil(o.DiskGB) {
		return true
	}

	return false
}

// SetDiskGB gets a reference to the given DiskGBAutoScaling and assigns it to the DiskGB field.
func (o *AdvancedAutoScalingSettings) SetDiskGB(v DiskGBAutoScaling) {
	o.DiskGB = &v
}
