// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CustomSessionTimeouts Defines the session timeout settings for managing user sessions at the organization level. When set to null, the field's value is unset, and the default timeout settings are applied.
type CustomSessionTimeouts struct {
	// Specifies the absolute session timeout duration in seconds. When set to null, the field's value is unset, and the default value of 43,200 seconds (12 hours) is applied. Accepted values range between a minimum of 3,600 seconds (1 hour) and a maximum of 43,200 seconds (12 hours).
	AbsoluteSessionTimeoutInSeconds *int `json:"absoluteSessionTimeoutInSeconds,omitempty"`
	// Specifies the idle session timeout duration in seconds. When set to null, the field's value is unset, and the default behavior depends on the context: no timeout for Atlas Commercial, and 600 seconds (10 minutes) for Atlas for Government. Accepted values start at a minimum of 300 seconds (5 minutes). For Atlas Commercial, the maximum value cannot exceed the configured absolute session timeout. For Atlas for Government, the maximum value is capped at 600 seconds (10 minutes).
	IdleSessionTimeoutInSeconds *int `json:"idleSessionTimeoutInSeconds,omitempty"`
}

// NewCustomSessionTimeouts instantiates a new CustomSessionTimeouts object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCustomSessionTimeouts() *CustomSessionTimeouts {
	this := CustomSessionTimeouts{}
	return &this
}

// NewCustomSessionTimeoutsWithDefaults instantiates a new CustomSessionTimeouts object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCustomSessionTimeoutsWithDefaults() *CustomSessionTimeouts {
	this := CustomSessionTimeouts{}
	return &this
}

// GetAbsoluteSessionTimeoutInSeconds returns the AbsoluteSessionTimeoutInSeconds field value if set, zero value otherwise
func (o *CustomSessionTimeouts) GetAbsoluteSessionTimeoutInSeconds() int {
	if o == nil || IsNil(o.AbsoluteSessionTimeoutInSeconds) {
		var ret int
		return ret
	}
	return *o.AbsoluteSessionTimeoutInSeconds
}

// GetAbsoluteSessionTimeoutInSecondsOk returns a tuple with the AbsoluteSessionTimeoutInSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CustomSessionTimeouts) GetAbsoluteSessionTimeoutInSecondsOk() (*int, bool) {
	if o == nil || IsNil(o.AbsoluteSessionTimeoutInSeconds) {
		return nil, false
	}

	return o.AbsoluteSessionTimeoutInSeconds, true
}

// HasAbsoluteSessionTimeoutInSeconds returns a boolean if a field has been set.
func (o *CustomSessionTimeouts) HasAbsoluteSessionTimeoutInSeconds() bool {
	if o != nil && !IsNil(o.AbsoluteSessionTimeoutInSeconds) {
		return true
	}

	return false
}

// SetAbsoluteSessionTimeoutInSeconds gets a reference to the given int and assigns it to the AbsoluteSessionTimeoutInSeconds field.
func (o *CustomSessionTimeouts) SetAbsoluteSessionTimeoutInSeconds(v int) {
	o.AbsoluteSessionTimeoutInSeconds = &v
}

// GetIdleSessionTimeoutInSeconds returns the IdleSessionTimeoutInSeconds field value if set, zero value otherwise
func (o *CustomSessionTimeouts) GetIdleSessionTimeoutInSeconds() int {
	if o == nil || IsNil(o.IdleSessionTimeoutInSeconds) {
		var ret int
		return ret
	}
	return *o.IdleSessionTimeoutInSeconds
}

// GetIdleSessionTimeoutInSecondsOk returns a tuple with the IdleSessionTimeoutInSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CustomSessionTimeouts) GetIdleSessionTimeoutInSecondsOk() (*int, bool) {
	if o == nil || IsNil(o.IdleSessionTimeoutInSeconds) {
		return nil, false
	}

	return o.IdleSessionTimeoutInSeconds, true
}

// HasIdleSessionTimeoutInSeconds returns a boolean if a field has been set.
func (o *CustomSessionTimeouts) HasIdleSessionTimeoutInSeconds() bool {
	if o != nil && !IsNil(o.IdleSessionTimeoutInSeconds) {
		return true
	}

	return false
}

// SetIdleSessionTimeoutInSeconds gets a reference to the given int and assigns it to the IdleSessionTimeoutInSeconds field.
func (o *CustomSessionTimeouts) SetIdleSessionTimeoutInSeconds(v int) {
	o.IdleSessionTimeoutInSeconds = &v
}
