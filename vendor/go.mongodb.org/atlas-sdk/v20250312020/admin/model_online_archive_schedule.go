// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OnlineArchiveSchedule Regular frequency and duration when archiving process occurs.
type OnlineArchiveSchedule struct {
	// Type of schedule.
	Type string `json:"type"`
	// Hour of the day when the scheduled window to run one online archive ends. This field uses the UTC time zone. The window must have a duration of at least two hours. If the end time is before or equal to the start time, the window extends to the next day.
	EndHour *int `json:"endHour,omitempty"`
	// Minute of the hour when the scheduled window to run one online archive ends. This field uses the UTC time zone. The window must have a duration of at least two hours. If the end time is before or equal to the start time, the window extends to the next day.
	EndMinute *int `json:"endMinute,omitempty"`
	// Hour of the day when the scheduled window to run one online archive starts. This field uses the UTC time zone.
	StartHour *int `json:"startHour,omitempty"`
	// Minute of the hour when the scheduled window to run one online archive starts. This field uses the UTC time zone.
	StartMinute *int `json:"startMinute,omitempty"`
	// Day of the week when the scheduled archive starts. The week starts with Monday (`1`) and ends with Sunday (`7`).
	DayOfWeek *int `json:"dayOfWeek,omitempty"`
	// Day of the month when the scheduled archive starts.
	DayOfMonth *int `json:"dayOfMonth,omitempty"`
}

// NewOnlineArchiveSchedule instantiates a new OnlineArchiveSchedule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOnlineArchiveSchedule(type_ string) *OnlineArchiveSchedule {
	this := OnlineArchiveSchedule{}
	this.Type = type_
	return &this
}

// NewOnlineArchiveScheduleWithDefaults instantiates a new OnlineArchiveSchedule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOnlineArchiveScheduleWithDefaults() *OnlineArchiveSchedule {
	this := OnlineArchiveSchedule{}
	return &this
}

// GetType returns the Type field value
func (o *OnlineArchiveSchedule) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *OnlineArchiveSchedule) SetType(v string) {
	o.Type = v
}

// GetEndHour returns the EndHour field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetEndHour() int {
	if o == nil || IsNil(o.EndHour) {
		var ret int
		return ret
	}
	return *o.EndHour
}

// GetEndHourOk returns a tuple with the EndHour field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetEndHourOk() (*int, bool) {
	if o == nil || IsNil(o.EndHour) {
		return nil, false
	}

	return o.EndHour, true
}

// HasEndHour returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasEndHour() bool {
	if o != nil && !IsNil(o.EndHour) {
		return true
	}

	return false
}

// SetEndHour gets a reference to the given int and assigns it to the EndHour field.
func (o *OnlineArchiveSchedule) SetEndHour(v int) {
	o.EndHour = &v
}

// GetEndMinute returns the EndMinute field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetEndMinute() int {
	if o == nil || IsNil(o.EndMinute) {
		var ret int
		return ret
	}
	return *o.EndMinute
}

// GetEndMinuteOk returns a tuple with the EndMinute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetEndMinuteOk() (*int, bool) {
	if o == nil || IsNil(o.EndMinute) {
		return nil, false
	}

	return o.EndMinute, true
}

// HasEndMinute returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasEndMinute() bool {
	if o != nil && !IsNil(o.EndMinute) {
		return true
	}

	return false
}

// SetEndMinute gets a reference to the given int and assigns it to the EndMinute field.
func (o *OnlineArchiveSchedule) SetEndMinute(v int) {
	o.EndMinute = &v
}

// GetStartHour returns the StartHour field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetStartHour() int {
	if o == nil || IsNil(o.StartHour) {
		var ret int
		return ret
	}
	return *o.StartHour
}

// GetStartHourOk returns a tuple with the StartHour field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetStartHourOk() (*int, bool) {
	if o == nil || IsNil(o.StartHour) {
		return nil, false
	}

	return o.StartHour, true
}

// HasStartHour returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasStartHour() bool {
	if o != nil && !IsNil(o.StartHour) {
		return true
	}

	return false
}

// SetStartHour gets a reference to the given int and assigns it to the StartHour field.
func (o *OnlineArchiveSchedule) SetStartHour(v int) {
	o.StartHour = &v
}

// GetStartMinute returns the StartMinute field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetStartMinute() int {
	if o == nil || IsNil(o.StartMinute) {
		var ret int
		return ret
	}
	return *o.StartMinute
}

// GetStartMinuteOk returns a tuple with the StartMinute field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetStartMinuteOk() (*int, bool) {
	if o == nil || IsNil(o.StartMinute) {
		return nil, false
	}

	return o.StartMinute, true
}

// HasStartMinute returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasStartMinute() bool {
	if o != nil && !IsNil(o.StartMinute) {
		return true
	}

	return false
}

// SetStartMinute gets a reference to the given int and assigns it to the StartMinute field.
func (o *OnlineArchiveSchedule) SetStartMinute(v int) {
	o.StartMinute = &v
}

// GetDayOfWeek returns the DayOfWeek field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetDayOfWeek() int {
	if o == nil || IsNil(o.DayOfWeek) {
		var ret int
		return ret
	}
	return *o.DayOfWeek
}

// GetDayOfWeekOk returns a tuple with the DayOfWeek field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetDayOfWeekOk() (*int, bool) {
	if o == nil || IsNil(o.DayOfWeek) {
		return nil, false
	}

	return o.DayOfWeek, true
}

// HasDayOfWeek returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasDayOfWeek() bool {
	if o != nil && !IsNil(o.DayOfWeek) {
		return true
	}

	return false
}

// SetDayOfWeek gets a reference to the given int and assigns it to the DayOfWeek field.
func (o *OnlineArchiveSchedule) SetDayOfWeek(v int) {
	o.DayOfWeek = &v
}

// GetDayOfMonth returns the DayOfMonth field value if set, zero value otherwise
func (o *OnlineArchiveSchedule) GetDayOfMonth() int {
	if o == nil || IsNil(o.DayOfMonth) {
		var ret int
		return ret
	}
	return *o.DayOfMonth
}

// GetDayOfMonthOk returns a tuple with the DayOfMonth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OnlineArchiveSchedule) GetDayOfMonthOk() (*int, bool) {
	if o == nil || IsNil(o.DayOfMonth) {
		return nil, false
	}

	return o.DayOfMonth, true
}

// HasDayOfMonth returns a boolean if a field has been set.
func (o *OnlineArchiveSchedule) HasDayOfMonth() bool {
	if o != nil && !IsNil(o.DayOfMonth) {
		return true
	}

	return false
}

// SetDayOfMonth gets a reference to the given int and assigns it to the DayOfMonth field.
func (o *OnlineArchiveSchedule) SetDayOfMonth(v int) {
	o.DayOfMonth = &v
}
