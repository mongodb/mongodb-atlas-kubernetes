// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskGBAutoScaling Setting that enables disk auto-scaling.
type DiskGBAutoScaling struct {
	// Flag that indicates whether this cluster enables disk auto-scaling. The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.
	Enabled *bool `json:"enabled,omitempty"`
}

// NewDiskGBAutoScaling instantiates a new DiskGBAutoScaling object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskGBAutoScaling() *DiskGBAutoScaling {
	this := DiskGBAutoScaling{}
	return &this
}

// NewDiskGBAutoScalingWithDefaults instantiates a new DiskGBAutoScaling object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskGBAutoScalingWithDefaults() *DiskGBAutoScaling {
	this := DiskGBAutoScaling{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *DiskGBAutoScaling) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskGBAutoScaling) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *DiskGBAutoScaling) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *DiskGBAutoScaling) SetEnabled(v bool) {
	o.Enabled = &v
}
