// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// GroupAlertsConfig struct for GroupAlertsConfig
type GroupAlertsConfig struct {
	// Date and time when MongoDB Cloud created the alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Flag that indicates whether someone enabled this alert configuration for the specified project.
	Enabled *bool `json:"enabled,omitempty"`
	// Event type that triggers an alert.
	EventTypeName *string `json:"eventTypeName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that owns this alert configuration.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this alert configuration.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// List of rules that determine whether MongoDB Cloud checks an object for the alert configuration.
	Matchers *[]StreamsMatcher `json:"matchers,omitempty"`
	// List that contains the targets that MongoDB Cloud sends notifications.
	Notifications *[]AlertsNotificationRootForGroup `json:"notifications,omitempty"`
	// Severity of the event.
	SeverityOverride *string `json:"severityOverride,omitempty"`
	// Date and time when someone last updated this alert configuration. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Updated         *time.Time                      `json:"updated,omitempty"`
	MetricThreshold *StreamProcessorMetricThreshold `json:"metricThreshold,omitempty"`
	Threshold       *StreamProcessorMetricThreshold `json:"threshold,omitempty"`
}

// NewGroupAlertsConfig instantiates a new GroupAlertsConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupAlertsConfig() *GroupAlertsConfig {
	this := GroupAlertsConfig{}
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// NewGroupAlertsConfigWithDefaults instantiates a new GroupAlertsConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupAlertsConfigWithDefaults() *GroupAlertsConfig {
	this := GroupAlertsConfig{}
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *GroupAlertsConfig) SetCreated(v time.Time) {
	o.Created = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *GroupAlertsConfig) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetEventTypeName returns the EventTypeName field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetEventTypeName() string {
	if o == nil || IsNil(o.EventTypeName) {
		var ret string
		return ret
	}
	return *o.EventTypeName
}

// GetEventTypeNameOk returns a tuple with the EventTypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetEventTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EventTypeName) {
		return nil, false
	}

	return o.EventTypeName, true
}

// HasEventTypeName returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasEventTypeName() bool {
	if o != nil && !IsNil(o.EventTypeName) {
		return true
	}

	return false
}

// SetEventTypeName gets a reference to the given string and assigns it to the EventTypeName field.
func (o *GroupAlertsConfig) SetEventTypeName(v string) {
	o.EventTypeName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *GroupAlertsConfig) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *GroupAlertsConfig) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *GroupAlertsConfig) SetLinks(v []Link) {
	o.Links = &v
}

// GetMatchers returns the Matchers field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetMatchers() []StreamsMatcher {
	if o == nil || IsNil(o.Matchers) {
		var ret []StreamsMatcher
		return ret
	}
	return *o.Matchers
}

// GetMatchersOk returns a tuple with the Matchers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetMatchersOk() (*[]StreamsMatcher, bool) {
	if o == nil || IsNil(o.Matchers) {
		return nil, false
	}

	return o.Matchers, true
}

// HasMatchers returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasMatchers() bool {
	if o != nil && !IsNil(o.Matchers) {
		return true
	}

	return false
}

// SetMatchers gets a reference to the given []StreamsMatcher and assigns it to the Matchers field.
func (o *GroupAlertsConfig) SetMatchers(v []StreamsMatcher) {
	o.Matchers = &v
}

// GetNotifications returns the Notifications field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetNotifications() []AlertsNotificationRootForGroup {
	if o == nil || IsNil(o.Notifications) {
		var ret []AlertsNotificationRootForGroup
		return ret
	}
	return *o.Notifications
}

// GetNotificationsOk returns a tuple with the Notifications field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetNotificationsOk() (*[]AlertsNotificationRootForGroup, bool) {
	if o == nil || IsNil(o.Notifications) {
		return nil, false
	}

	return o.Notifications, true
}

// HasNotifications returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasNotifications() bool {
	if o != nil && !IsNil(o.Notifications) {
		return true
	}

	return false
}

// SetNotifications gets a reference to the given []AlertsNotificationRootForGroup and assigns it to the Notifications field.
func (o *GroupAlertsConfig) SetNotifications(v []AlertsNotificationRootForGroup) {
	o.Notifications = &v
}

// GetSeverityOverride returns the SeverityOverride field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetSeverityOverride() string {
	if o == nil || IsNil(o.SeverityOverride) {
		var ret string
		return ret
	}
	return *o.SeverityOverride
}

// GetSeverityOverrideOk returns a tuple with the SeverityOverride field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetSeverityOverrideOk() (*string, bool) {
	if o == nil || IsNil(o.SeverityOverride) {
		return nil, false
	}

	return o.SeverityOverride, true
}

// HasSeverityOverride returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasSeverityOverride() bool {
	if o != nil && !IsNil(o.SeverityOverride) {
		return true
	}

	return false
}

// SetSeverityOverride gets a reference to the given string and assigns it to the SeverityOverride field.
func (o *GroupAlertsConfig) SetSeverityOverride(v string) {
	o.SeverityOverride = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetUpdated() time.Time {
	if o == nil || IsNil(o.Updated) {
		var ret time.Time
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetUpdatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}

	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given time.Time and assigns it to the Updated field.
func (o *GroupAlertsConfig) SetUpdated(v time.Time) {
	o.Updated = &v
}

// GetMetricThreshold returns the MetricThreshold field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetMetricThreshold() StreamProcessorMetricThreshold {
	if o == nil || IsNil(o.MetricThreshold) {
		var ret StreamProcessorMetricThreshold
		return ret
	}
	return *o.MetricThreshold
}

// GetMetricThresholdOk returns a tuple with the MetricThreshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetMetricThresholdOk() (*StreamProcessorMetricThreshold, bool) {
	if o == nil || IsNil(o.MetricThreshold) {
		return nil, false
	}

	return o.MetricThreshold, true
}

// HasMetricThreshold returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasMetricThreshold() bool {
	if o != nil && !IsNil(o.MetricThreshold) {
		return true
	}

	return false
}

// SetMetricThreshold gets a reference to the given StreamProcessorMetricThreshold and assigns it to the MetricThreshold field.
func (o *GroupAlertsConfig) SetMetricThreshold(v StreamProcessorMetricThreshold) {
	o.MetricThreshold = &v
}

// GetThreshold returns the Threshold field value if set, zero value otherwise
func (o *GroupAlertsConfig) GetThreshold() StreamProcessorMetricThreshold {
	if o == nil || IsNil(o.Threshold) {
		var ret StreamProcessorMetricThreshold
		return ret
	}
	return *o.Threshold
}

// GetThresholdOk returns a tuple with the Threshold field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupAlertsConfig) GetThresholdOk() (*StreamProcessorMetricThreshold, bool) {
	if o == nil || IsNil(o.Threshold) {
		return nil, false
	}

	return o.Threshold, true
}

// HasThreshold returns a boolean if a field has been set.
func (o *GroupAlertsConfig) HasThreshold() bool {
	if o != nil && !IsNil(o.Threshold) {
		return true
	}

	return false
}

// SetThreshold gets a reference to the given StreamProcessorMetricThreshold and assigns it to the Threshold field.
func (o *GroupAlertsConfig) SetThreshold(v StreamProcessorMetricThreshold) {
	o.Threshold = &v
}
