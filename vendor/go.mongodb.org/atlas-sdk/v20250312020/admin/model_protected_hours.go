// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ProtectedHours Defines the a window where maintenance will not begin within.
type ProtectedHours struct {
	// Zero-based integer that represents the end hour of the of the day that the maintenance will not begin in.
	EndHourOfDay *int `json:"endHourOfDay,omitempty"`
	// Zero-based integer that represents the beginning hour of the of the day that the maintenance will not begin in.
	StartHourOfDay *int `json:"startHourOfDay,omitempty"`
}

// NewProtectedHours instantiates a new ProtectedHours object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProtectedHours() *ProtectedHours {
	this := ProtectedHours{}
	return &this
}

// NewProtectedHoursWithDefaults instantiates a new ProtectedHours object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProtectedHoursWithDefaults() *ProtectedHours {
	this := ProtectedHours{}
	return &this
}

// GetEndHourOfDay returns the EndHourOfDay field value if set, zero value otherwise
func (o *ProtectedHours) GetEndHourOfDay() int {
	if o == nil || IsNil(o.EndHourOfDay) {
		var ret int
		return ret
	}
	return *o.EndHourOfDay
}

// GetEndHourOfDayOk returns a tuple with the EndHourOfDay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProtectedHours) GetEndHourOfDayOk() (*int, bool) {
	if o == nil || IsNil(o.EndHourOfDay) {
		return nil, false
	}

	return o.EndHourOfDay, true
}

// HasEndHourOfDay returns a boolean if a field has been set.
func (o *ProtectedHours) HasEndHourOfDay() bool {
	if o != nil && !IsNil(o.EndHourOfDay) {
		return true
	}

	return false
}

// SetEndHourOfDay gets a reference to the given int and assigns it to the EndHourOfDay field.
func (o *ProtectedHours) SetEndHourOfDay(v int) {
	o.EndHourOfDay = &v
}

// GetStartHourOfDay returns the StartHourOfDay field value if set, zero value otherwise
func (o *ProtectedHours) GetStartHourOfDay() int {
	if o == nil || IsNil(o.StartHourOfDay) {
		var ret int
		return ret
	}
	return *o.StartHourOfDay
}

// GetStartHourOfDayOk returns a tuple with the StartHourOfDay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProtectedHours) GetStartHourOfDayOk() (*int, bool) {
	if o == nil || IsNil(o.StartHourOfDay) {
		return nil, false
	}

	return o.StartHourOfDay, true
}

// HasStartHourOfDay returns a boolean if a field has been set.
func (o *ProtectedHours) HasStartHourOfDay() bool {
	if o != nil && !IsNil(o.StartHourOfDay) {
		return true
	}

	return false
}

// SetStartHourOfDay gets a reference to the given int and assigns it to the StartHourOfDay field.
func (o *ProtectedHours) SetStartHourOfDay(v int) {
	o.StartHourOfDay = &v
}
