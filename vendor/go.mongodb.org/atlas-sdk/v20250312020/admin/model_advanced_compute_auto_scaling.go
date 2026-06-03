// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AdvancedComputeAutoScaling Options that determine how this cluster handles CPU scaling.
type AdvancedComputeAutoScaling struct {
	// Flag that indicates whether instance size reactive auto-scaling is enabled.  - Set to `true` to enable instance size reactive auto-scaling. If enabled, you must specify a value for `replicationSpecs[n].regionConfigs[m].autoScaling.compute.maxInstanceSize`. - Set to `false` to disable instance size reactive auto-scaling.
	Enabled *bool `json:"enabled,omitempty"`
	// Instance size boundary to which your cluster can automatically scale.
	// Read only field.
	MaxInstanceSize *string `json:"maxInstanceSize,omitempty"`
	// Instance size boundary to which your cluster can automatically scale.
	// Read only field.
	MinInstanceSize *string `json:"minInstanceSize,omitempty"`
	// Flag that indicates whether the instance size may scale down via reactive auto-scaling. MongoDB Cloud requires this parameter if `replicationSpecs[n].regionConfigs[m].autoScaling.compute.enabled` is `true`. If you enable this option, specify a value for `replicationSpecs[n].regionConfigs[m].autoScaling.compute.minInstanceSize`.
	ScaleDownEnabled *bool `json:"scaleDownEnabled,omitempty"`
}

// NewAdvancedComputeAutoScaling instantiates a new AdvancedComputeAutoScaling object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAdvancedComputeAutoScaling() *AdvancedComputeAutoScaling {
	this := AdvancedComputeAutoScaling{}
	return &this
}

// NewAdvancedComputeAutoScalingWithDefaults instantiates a new AdvancedComputeAutoScaling object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAdvancedComputeAutoScalingWithDefaults() *AdvancedComputeAutoScaling {
	this := AdvancedComputeAutoScaling{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *AdvancedComputeAutoScaling) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedComputeAutoScaling) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *AdvancedComputeAutoScaling) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *AdvancedComputeAutoScaling) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetMaxInstanceSize returns the MaxInstanceSize field value if set, zero value otherwise
func (o *AdvancedComputeAutoScaling) GetMaxInstanceSize() string {
	if o == nil || IsNil(o.MaxInstanceSize) {
		var ret string
		return ret
	}
	return *o.MaxInstanceSize
}

// GetMaxInstanceSizeOk returns a tuple with the MaxInstanceSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedComputeAutoScaling) GetMaxInstanceSizeOk() (*string, bool) {
	if o == nil || IsNil(o.MaxInstanceSize) {
		return nil, false
	}

	return o.MaxInstanceSize, true
}

// HasMaxInstanceSize returns a boolean if a field has been set.
func (o *AdvancedComputeAutoScaling) HasMaxInstanceSize() bool {
	if o != nil && !IsNil(o.MaxInstanceSize) {
		return true
	}

	return false
}

// SetMaxInstanceSize gets a reference to the given string and assigns it to the MaxInstanceSize field.
func (o *AdvancedComputeAutoScaling) SetMaxInstanceSize(v string) {
	o.MaxInstanceSize = &v
}

// GetMinInstanceSize returns the MinInstanceSize field value if set, zero value otherwise
func (o *AdvancedComputeAutoScaling) GetMinInstanceSize() string {
	if o == nil || IsNil(o.MinInstanceSize) {
		var ret string
		return ret
	}
	return *o.MinInstanceSize
}

// GetMinInstanceSizeOk returns a tuple with the MinInstanceSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedComputeAutoScaling) GetMinInstanceSizeOk() (*string, bool) {
	if o == nil || IsNil(o.MinInstanceSize) {
		return nil, false
	}

	return o.MinInstanceSize, true
}

// HasMinInstanceSize returns a boolean if a field has been set.
func (o *AdvancedComputeAutoScaling) HasMinInstanceSize() bool {
	if o != nil && !IsNil(o.MinInstanceSize) {
		return true
	}

	return false
}

// SetMinInstanceSize gets a reference to the given string and assigns it to the MinInstanceSize field.
func (o *AdvancedComputeAutoScaling) SetMinInstanceSize(v string) {
	o.MinInstanceSize = &v
}

// GetScaleDownEnabled returns the ScaleDownEnabled field value if set, zero value otherwise
func (o *AdvancedComputeAutoScaling) GetScaleDownEnabled() bool {
	if o == nil || IsNil(o.ScaleDownEnabled) {
		var ret bool
		return ret
	}
	return *o.ScaleDownEnabled
}

// GetScaleDownEnabledOk returns a tuple with the ScaleDownEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedComputeAutoScaling) GetScaleDownEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.ScaleDownEnabled) {
		return nil, false
	}

	return o.ScaleDownEnabled, true
}

// HasScaleDownEnabled returns a boolean if a field has been set.
func (o *AdvancedComputeAutoScaling) HasScaleDownEnabled() bool {
	if o != nil && !IsNil(o.ScaleDownEnabled) {
		return true
	}

	return false
}

// SetScaleDownEnabled gets a reference to the given bool and assigns it to the ScaleDownEnabled field.
func (o *AdvancedComputeAutoScaling) SetScaleDownEnabled(v bool) {
	o.ScaleDownEnabled = &v
}
