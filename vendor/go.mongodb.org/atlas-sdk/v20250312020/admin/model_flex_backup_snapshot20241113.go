// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// FlexBackupSnapshot20241113 Details for one snapshot of a flex cluster.
type FlexBackupSnapshot20241113 struct {
	// Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Expiration *time.Time `json:"expiration,omitempty"`
	// Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	FinishTime *time.Time `json:"finishTime,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// MongoDB host version that the snapshot runs.
	// Read only field.
	MongoDBVersion *string `json:"mongoDBVersion,omitempty"`
	// Date and time when MongoDB Cloud will take the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	ScheduledTime *time.Time `json:"scheduledTime,omitempty"`
	// Date and time when MongoDB Cloud began taking the snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	StartTime *time.Time `json:"startTime,omitempty"`
	// Phase of the workflow for this snapshot at the time this resource made this request.
	// Read only field.
	Status *string `json:"status,omitempty"`
}

// NewFlexBackupSnapshot20241113 instantiates a new FlexBackupSnapshot20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexBackupSnapshot20241113() *FlexBackupSnapshot20241113 {
	this := FlexBackupSnapshot20241113{}
	return &this
}

// NewFlexBackupSnapshot20241113WithDefaults instantiates a new FlexBackupSnapshot20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexBackupSnapshot20241113WithDefaults() *FlexBackupSnapshot20241113 {
	this := FlexBackupSnapshot20241113{}
	return &this
}

// GetExpiration returns the Expiration field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetExpiration() time.Time {
	if o == nil || IsNil(o.Expiration) {
		var ret time.Time
		return ret
	}
	return *o.Expiration
}

// GetExpirationOk returns a tuple with the Expiration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetExpirationOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Expiration) {
		return nil, false
	}

	return o.Expiration, true
}

// HasExpiration returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasExpiration() bool {
	if o != nil && !IsNil(o.Expiration) {
		return true
	}

	return false
}

// SetExpiration gets a reference to the given time.Time and assigns it to the Expiration field.
func (o *FlexBackupSnapshot20241113) SetExpiration(v time.Time) {
	o.Expiration = &v
}

// GetFinishTime returns the FinishTime field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetFinishTime() time.Time {
	if o == nil || IsNil(o.FinishTime) {
		var ret time.Time
		return ret
	}
	return *o.FinishTime
}

// GetFinishTimeOk returns a tuple with the FinishTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetFinishTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FinishTime) {
		return nil, false
	}

	return o.FinishTime, true
}

// HasFinishTime returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasFinishTime() bool {
	if o != nil && !IsNil(o.FinishTime) {
		return true
	}

	return false
}

// SetFinishTime gets a reference to the given time.Time and assigns it to the FinishTime field.
func (o *FlexBackupSnapshot20241113) SetFinishTime(v time.Time) {
	o.FinishTime = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *FlexBackupSnapshot20241113) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *FlexBackupSnapshot20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetMongoDBVersion returns the MongoDBVersion field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetMongoDBVersion() string {
	if o == nil || IsNil(o.MongoDBVersion) {
		var ret string
		return ret
	}
	return *o.MongoDBVersion
}

// GetMongoDBVersionOk returns a tuple with the MongoDBVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetMongoDBVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MongoDBVersion) {
		return nil, false
	}

	return o.MongoDBVersion, true
}

// HasMongoDBVersion returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasMongoDBVersion() bool {
	if o != nil && !IsNil(o.MongoDBVersion) {
		return true
	}

	return false
}

// SetMongoDBVersion gets a reference to the given string and assigns it to the MongoDBVersion field.
func (o *FlexBackupSnapshot20241113) SetMongoDBVersion(v string) {
	o.MongoDBVersion = &v
}

// GetScheduledTime returns the ScheduledTime field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetScheduledTime() time.Time {
	if o == nil || IsNil(o.ScheduledTime) {
		var ret time.Time
		return ret
	}
	return *o.ScheduledTime
}

// GetScheduledTimeOk returns a tuple with the ScheduledTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetScheduledTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ScheduledTime) {
		return nil, false
	}

	return o.ScheduledTime, true
}

// HasScheduledTime returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasScheduledTime() bool {
	if o != nil && !IsNil(o.ScheduledTime) {
		return true
	}

	return false
}

// SetScheduledTime gets a reference to the given time.Time and assigns it to the ScheduledTime field.
func (o *FlexBackupSnapshot20241113) SetScheduledTime(v time.Time) {
	o.ScheduledTime = &v
}

// GetStartTime returns the StartTime field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetStartTime() time.Time {
	if o == nil || IsNil(o.StartTime) {
		var ret time.Time
		return ret
	}
	return *o.StartTime
}

// GetStartTimeOk returns a tuple with the StartTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetStartTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartTime) {
		return nil, false
	}

	return o.StartTime, true
}

// HasStartTime returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasStartTime() bool {
	if o != nil && !IsNil(o.StartTime) {
		return true
	}

	return false
}

// SetStartTime gets a reference to the given time.Time and assigns it to the StartTime field.
func (o *FlexBackupSnapshot20241113) SetStartTime(v time.Time) {
	o.StartTime = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *FlexBackupSnapshot20241113) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSnapshot20241113) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *FlexBackupSnapshot20241113) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *FlexBackupSnapshot20241113) SetStatus(v string) {
	o.Status = &v
}
