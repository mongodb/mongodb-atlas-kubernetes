// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ApiAtlasCheckpoint struct for ApiAtlasCheckpoint
type ApiAtlasCheckpoint struct {
	// Unique 24-hexadecimal digit string that identifies the cluster that contains the checkpoint.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Date and time when the checkpoint completed and the balancer restarted. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Completed *time.Time `json:"completed,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that owns the checkpoints.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies checkpoint.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Metadata that describes the complete snapshot.  - For a replica set, this array contains a single document. - For a sharded cluster, this array contains one document for each shard plus one document for the config host.
	// Read only field.
	Parts *[]ApiCheckpointPart `json:"parts,omitempty"`
	// Flag that indicates whether MongoDB Cloud can use the checkpoint for a restore.
	// Read only field.
	Restorable *bool `json:"restorable,omitempty"`
	// Date and time when the balancer stopped and began the checkpoint. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Started *time.Time `json:"started,omitempty"`
	// Date and time to which the checkpoint restores. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// NewApiAtlasCheckpoint instantiates a new ApiAtlasCheckpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasCheckpoint() *ApiAtlasCheckpoint {
	this := ApiAtlasCheckpoint{}
	return &this
}

// NewApiAtlasCheckpointWithDefaults instantiates a new ApiAtlasCheckpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasCheckpointWithDefaults() *ApiAtlasCheckpoint {
	this := ApiAtlasCheckpoint{}
	return &this
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *ApiAtlasCheckpoint) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetCompleted returns the Completed field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetCompleted() time.Time {
	if o == nil || IsNil(o.Completed) {
		var ret time.Time
		return ret
	}
	return *o.Completed
}

// GetCompletedOk returns a tuple with the Completed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetCompletedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Completed) {
		return nil, false
	}

	return o.Completed, true
}

// HasCompleted returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasCompleted() bool {
	if o != nil && !IsNil(o.Completed) {
		return true
	}

	return false
}

// SetCompleted gets a reference to the given time.Time and assigns it to the Completed field.
func (o *ApiAtlasCheckpoint) SetCompleted(v time.Time) {
	o.Completed = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ApiAtlasCheckpoint) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiAtlasCheckpoint) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ApiAtlasCheckpoint) SetLinks(v []Link) {
	o.Links = &v
}

// GetParts returns the Parts field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetParts() []ApiCheckpointPart {
	if o == nil || IsNil(o.Parts) {
		var ret []ApiCheckpointPart
		return ret
	}
	return *o.Parts
}

// GetPartsOk returns a tuple with the Parts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetPartsOk() (*[]ApiCheckpointPart, bool) {
	if o == nil || IsNil(o.Parts) {
		return nil, false
	}

	return o.Parts, true
}

// HasParts returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasParts() bool {
	if o != nil && !IsNil(o.Parts) {
		return true
	}

	return false
}

// SetParts gets a reference to the given []ApiCheckpointPart and assigns it to the Parts field.
func (o *ApiAtlasCheckpoint) SetParts(v []ApiCheckpointPart) {
	o.Parts = &v
}

// GetRestorable returns the Restorable field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetRestorable() bool {
	if o == nil || IsNil(o.Restorable) {
		var ret bool
		return ret
	}
	return *o.Restorable
}

// GetRestorableOk returns a tuple with the Restorable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetRestorableOk() (*bool, bool) {
	if o == nil || IsNil(o.Restorable) {
		return nil, false
	}

	return o.Restorable, true
}

// HasRestorable returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasRestorable() bool {
	if o != nil && !IsNil(o.Restorable) {
		return true
	}

	return false
}

// SetRestorable gets a reference to the given bool and assigns it to the Restorable field.
func (o *ApiAtlasCheckpoint) SetRestorable(v bool) {
	o.Restorable = &v
}

// GetStarted returns the Started field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetStarted() time.Time {
	if o == nil || IsNil(o.Started) {
		var ret time.Time
		return ret
	}
	return *o.Started
}

// GetStartedOk returns a tuple with the Started field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetStartedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Started) {
		return nil, false
	}

	return o.Started, true
}

// HasStarted returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasStarted() bool {
	if o != nil && !IsNil(o.Started) {
		return true
	}

	return false
}

// SetStarted gets a reference to the given time.Time and assigns it to the Started field.
func (o *ApiAtlasCheckpoint) SetStarted(v time.Time) {
	o.Started = &v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise
func (o *ApiAtlasCheckpoint) GetTimestamp() time.Time {
	if o == nil || IsNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasCheckpoint) GetTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}

	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *ApiAtlasCheckpoint) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *ApiAtlasCheckpoint) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}
