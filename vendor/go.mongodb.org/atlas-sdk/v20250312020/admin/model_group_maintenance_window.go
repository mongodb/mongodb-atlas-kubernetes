// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupMaintenanceWindow struct for GroupMaintenanceWindow
type GroupMaintenanceWindow struct {
	// Flag that indicates whether MongoDB Cloud should defer all maintenance windows for one week after you enable them. This setting controls the same underlying auto-deferral feature as the `/maintenanceWindow/autoDefer` endpoint. Use either this field (to set a specific value) or that endpoint (to toggle the current value). For most use cases, this field in the PATCH request is preferred because it allows setting an explicit value rather than toggling.
	AutoDeferOnceEnabled *bool `json:"autoDeferOnceEnabled,omitempty"`
	// One-based integer that represents the day of the week that the maintenance window starts.  - `1`: Sunday. - `2`: Monday. - `3`: Tuesday. - `4`: Wednesday. - `5`: Thursday. - `6`: Friday. - `7`: Saturday.
	DayOfWeek int `json:"dayOfWeek"`
	// Zero-based integer that represents the hour of the of the day that the maintenance window starts according to a 24-hour clock. Use `0` for midnight and `12` for noon.
	HourOfDay *int `json:"hourOfDay,omitempty"`
	// Number of times the current maintenance event for this project has been deferred.
	// Read only field.
	NumberOfDeferrals *int            `json:"numberOfDeferrals,omitempty"`
	ProtectedHours    *ProtectedHours `json:"protectedHours,omitempty"`
	// Flag that indicates whether MongoDB Cloud starts the maintenance window immediately upon receiving this request. To start the maintenance window immediately for your project, MongoDB Cloud must have maintenance scheduled and you must set a maintenance window. This flag resets to `false` after MongoDB Cloud completes maintenance.
	StartASAP *bool `json:"startASAP,omitempty"`
	// Identifier for the current time zone of the maintenance window. This can only be updated via the Project Settings UI.
	// Read only field.
	TimeZoneId *string `json:"timeZoneId,omitempty"`
}

// NewGroupMaintenanceWindow instantiates a new GroupMaintenanceWindow object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupMaintenanceWindow(dayOfWeek int) *GroupMaintenanceWindow {
	this := GroupMaintenanceWindow{}
	this.DayOfWeek = dayOfWeek
	return &this
}

// NewGroupMaintenanceWindowWithDefaults instantiates a new GroupMaintenanceWindow object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupMaintenanceWindowWithDefaults() *GroupMaintenanceWindow {
	this := GroupMaintenanceWindow{}
	return &this
}

// GetAutoDeferOnceEnabled returns the AutoDeferOnceEnabled field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetAutoDeferOnceEnabled() bool {
	if o == nil || IsNil(o.AutoDeferOnceEnabled) {
		var ret bool
		return ret
	}
	return *o.AutoDeferOnceEnabled
}

// GetAutoDeferOnceEnabledOk returns a tuple with the AutoDeferOnceEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetAutoDeferOnceEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AutoDeferOnceEnabled) {
		return nil, false
	}

	return o.AutoDeferOnceEnabled, true
}

// HasAutoDeferOnceEnabled returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasAutoDeferOnceEnabled() bool {
	if o != nil && !IsNil(o.AutoDeferOnceEnabled) {
		return true
	}

	return false
}

// SetAutoDeferOnceEnabled gets a reference to the given bool and assigns it to the AutoDeferOnceEnabled field.
func (o *GroupMaintenanceWindow) SetAutoDeferOnceEnabled(v bool) {
	o.AutoDeferOnceEnabled = &v
}

// GetDayOfWeek returns the DayOfWeek field value
func (o *GroupMaintenanceWindow) GetDayOfWeek() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.DayOfWeek
}

// GetDayOfWeekOk returns a tuple with the DayOfWeek field value
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetDayOfWeekOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DayOfWeek, true
}

// SetDayOfWeek sets field value
func (o *GroupMaintenanceWindow) SetDayOfWeek(v int) {
	o.DayOfWeek = v
}

// GetHourOfDay returns the HourOfDay field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetHourOfDay() int {
	if o == nil || IsNil(o.HourOfDay) {
		var ret int
		return ret
	}
	return *o.HourOfDay
}

// GetHourOfDayOk returns a tuple with the HourOfDay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetHourOfDayOk() (*int, bool) {
	if o == nil || IsNil(o.HourOfDay) {
		return nil, false
	}

	return o.HourOfDay, true
}

// HasHourOfDay returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasHourOfDay() bool {
	if o != nil && !IsNil(o.HourOfDay) {
		return true
	}

	return false
}

// SetHourOfDay gets a reference to the given int and assigns it to the HourOfDay field.
func (o *GroupMaintenanceWindow) SetHourOfDay(v int) {
	o.HourOfDay = &v
}

// GetNumberOfDeferrals returns the NumberOfDeferrals field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetNumberOfDeferrals() int {
	if o == nil || IsNil(o.NumberOfDeferrals) {
		var ret int
		return ret
	}
	return *o.NumberOfDeferrals
}

// GetNumberOfDeferralsOk returns a tuple with the NumberOfDeferrals field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetNumberOfDeferralsOk() (*int, bool) {
	if o == nil || IsNil(o.NumberOfDeferrals) {
		return nil, false
	}

	return o.NumberOfDeferrals, true
}

// HasNumberOfDeferrals returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasNumberOfDeferrals() bool {
	if o != nil && !IsNil(o.NumberOfDeferrals) {
		return true
	}

	return false
}

// SetNumberOfDeferrals gets a reference to the given int and assigns it to the NumberOfDeferrals field.
func (o *GroupMaintenanceWindow) SetNumberOfDeferrals(v int) {
	o.NumberOfDeferrals = &v
}

// GetProtectedHours returns the ProtectedHours field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetProtectedHours() ProtectedHours {
	if o == nil || IsNil(o.ProtectedHours) {
		var ret ProtectedHours
		return ret
	}
	return *o.ProtectedHours
}

// GetProtectedHoursOk returns a tuple with the ProtectedHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetProtectedHoursOk() (*ProtectedHours, bool) {
	if o == nil || IsNil(o.ProtectedHours) {
		return nil, false
	}

	return o.ProtectedHours, true
}

// HasProtectedHours returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasProtectedHours() bool {
	if o != nil && !IsNil(o.ProtectedHours) {
		return true
	}

	return false
}

// SetProtectedHours gets a reference to the given ProtectedHours and assigns it to the ProtectedHours field.
func (o *GroupMaintenanceWindow) SetProtectedHours(v ProtectedHours) {
	o.ProtectedHours = &v
}

// GetStartASAP returns the StartASAP field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetStartASAP() bool {
	if o == nil || IsNil(o.StartASAP) {
		var ret bool
		return ret
	}
	return *o.StartASAP
}

// GetStartASAPOk returns a tuple with the StartASAP field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetStartASAPOk() (*bool, bool) {
	if o == nil || IsNil(o.StartASAP) {
		return nil, false
	}

	return o.StartASAP, true
}

// HasStartASAP returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasStartASAP() bool {
	if o != nil && !IsNil(o.StartASAP) {
		return true
	}

	return false
}

// SetStartASAP gets a reference to the given bool and assigns it to the StartASAP field.
func (o *GroupMaintenanceWindow) SetStartASAP(v bool) {
	o.StartASAP = &v
}

// GetTimeZoneId returns the TimeZoneId field value if set, zero value otherwise
func (o *GroupMaintenanceWindow) GetTimeZoneId() string {
	if o == nil || IsNil(o.TimeZoneId) {
		var ret string
		return ret
	}
	return *o.TimeZoneId
}

// GetTimeZoneIdOk returns a tuple with the TimeZoneId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupMaintenanceWindow) GetTimeZoneIdOk() (*string, bool) {
	if o == nil || IsNil(o.TimeZoneId) {
		return nil, false
	}

	return o.TimeZoneId, true
}

// HasTimeZoneId returns a boolean if a field has been set.
func (o *GroupMaintenanceWindow) HasTimeZoneId() bool {
	if o != nil && !IsNil(o.TimeZoneId) {
		return true
	}

	return false
}

// SetTimeZoneId gets a reference to the given string and assigns it to the TimeZoneId field.
func (o *GroupMaintenanceWindow) SetTimeZoneId(v string) {
	o.TimeZoneId = &v
}
