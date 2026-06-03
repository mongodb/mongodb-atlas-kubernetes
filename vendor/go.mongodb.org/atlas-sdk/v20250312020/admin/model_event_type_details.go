// Code based on the AtlasAPI V2 OpenAPI file

package admin

// EventTypeDetails A singular type of event.
type EventTypeDetails struct {
	// Whether or not this event type can be configured as an alert via the API.
	// Read only field.
	Alertable *bool `json:"alertable,omitempty"`
	// Description of the event type.
	// Read only field.
	Description *string `json:"description,omitempty"`
	// Enum representation of the event type.
	// Read only field.
	EventType *string `json:"eventType,omitempty"`
}

// NewEventTypeDetails instantiates a new EventTypeDetails object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEventTypeDetails() *EventTypeDetails {
	this := EventTypeDetails{}
	return &this
}

// NewEventTypeDetailsWithDefaults instantiates a new EventTypeDetails object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEventTypeDetailsWithDefaults() *EventTypeDetails {
	this := EventTypeDetails{}
	return &this
}

// GetAlertable returns the Alertable field value if set, zero value otherwise
func (o *EventTypeDetails) GetAlertable() bool {
	if o == nil || IsNil(o.Alertable) {
		var ret bool
		return ret
	}
	return *o.Alertable
}

// GetAlertableOk returns a tuple with the Alertable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventTypeDetails) GetAlertableOk() (*bool, bool) {
	if o == nil || IsNil(o.Alertable) {
		return nil, false
	}

	return o.Alertable, true
}

// HasAlertable returns a boolean if a field has been set.
func (o *EventTypeDetails) HasAlertable() bool {
	if o != nil && !IsNil(o.Alertable) {
		return true
	}

	return false
}

// SetAlertable gets a reference to the given bool and assigns it to the Alertable field.
func (o *EventTypeDetails) SetAlertable(v bool) {
	o.Alertable = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *EventTypeDetails) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventTypeDetails) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *EventTypeDetails) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *EventTypeDetails) SetDescription(v string) {
	o.Description = &v
}

// GetEventType returns the EventType field value if set, zero value otherwise
func (o *EventTypeDetails) GetEventType() string {
	if o == nil || IsNil(o.EventType) {
		var ret string
		return ret
	}
	return *o.EventType
}

// GetEventTypeOk returns a tuple with the EventType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EventTypeDetails) GetEventTypeOk() (*string, bool) {
	if o == nil || IsNil(o.EventType) {
		return nil, false
	}

	return o.EventType, true
}

// HasEventType returns a boolean if a field has been set.
func (o *EventTypeDetails) HasEventType() bool {
	if o != nil && !IsNil(o.EventType) {
		return true
	}

	return false
}

// SetEventType gets a reference to the given string and assigns it to the EventType field.
func (o *EventTypeDetails) SetEventType(v string) {
	o.EventType = &v
}
