// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// StreamsStartStreamProcessorWith A request to start a stream processor.
type StreamsStartStreamProcessorWith struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// When true or not specified, the stream processor resumes from its last checkpoint. When false, the stream processor starts fresh.
	ResumeFromCheckpoint *bool `json:"resumeFromCheckpoint,omitempty"`
	// The operation time after which the change stream source should begin reporting. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	StartAtOperationTime *time.Time `json:"startAtOperationTime,omitempty"`
	// Selected tier for the Stream Workspace. Configures Memory / VCPU allowances.
	Tier *string `json:"tier,omitempty"`
}

// NewStreamsStartStreamProcessorWith instantiates a new StreamsStartStreamProcessorWith object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsStartStreamProcessorWith() *StreamsStartStreamProcessorWith {
	this := StreamsStartStreamProcessorWith{}
	return &this
}

// NewStreamsStartStreamProcessorWithWithDefaults instantiates a new StreamsStartStreamProcessorWith object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsStartStreamProcessorWithWithDefaults() *StreamsStartStreamProcessorWith {
	this := StreamsStartStreamProcessorWith{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsStartStreamProcessorWith) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsStartStreamProcessorWith) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsStartStreamProcessorWith) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsStartStreamProcessorWith) SetLinks(v []Link) {
	o.Links = &v
}

// GetResumeFromCheckpoint returns the ResumeFromCheckpoint field value if set, zero value otherwise
func (o *StreamsStartStreamProcessorWith) GetResumeFromCheckpoint() bool {
	if o == nil || IsNil(o.ResumeFromCheckpoint) {
		var ret bool
		return ret
	}
	return *o.ResumeFromCheckpoint
}

// GetResumeFromCheckpointOk returns a tuple with the ResumeFromCheckpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsStartStreamProcessorWith) GetResumeFromCheckpointOk() (*bool, bool) {
	if o == nil || IsNil(o.ResumeFromCheckpoint) {
		return nil, false
	}

	return o.ResumeFromCheckpoint, true
}

// HasResumeFromCheckpoint returns a boolean if a field has been set.
func (o *StreamsStartStreamProcessorWith) HasResumeFromCheckpoint() bool {
	if o != nil && !IsNil(o.ResumeFromCheckpoint) {
		return true
	}

	return false
}

// SetResumeFromCheckpoint gets a reference to the given bool and assigns it to the ResumeFromCheckpoint field.
func (o *StreamsStartStreamProcessorWith) SetResumeFromCheckpoint(v bool) {
	o.ResumeFromCheckpoint = &v
}

// GetStartAtOperationTime returns the StartAtOperationTime field value if set, zero value otherwise
func (o *StreamsStartStreamProcessorWith) GetStartAtOperationTime() time.Time {
	if o == nil || IsNil(o.StartAtOperationTime) {
		var ret time.Time
		return ret
	}
	return *o.StartAtOperationTime
}

// GetStartAtOperationTimeOk returns a tuple with the StartAtOperationTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsStartStreamProcessorWith) GetStartAtOperationTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.StartAtOperationTime) {
		return nil, false
	}

	return o.StartAtOperationTime, true
}

// HasStartAtOperationTime returns a boolean if a field has been set.
func (o *StreamsStartStreamProcessorWith) HasStartAtOperationTime() bool {
	if o != nil && !IsNil(o.StartAtOperationTime) {
		return true
	}

	return false
}

// SetStartAtOperationTime gets a reference to the given time.Time and assigns it to the StartAtOperationTime field.
func (o *StreamsStartStreamProcessorWith) SetStartAtOperationTime(v time.Time) {
	o.StartAtOperationTime = &v
}

// GetTier returns the Tier field value if set, zero value otherwise
func (o *StreamsStartStreamProcessorWith) GetTier() string {
	if o == nil || IsNil(o.Tier) {
		var ret string
		return ret
	}
	return *o.Tier
}

// GetTierOk returns a tuple with the Tier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsStartStreamProcessorWith) GetTierOk() (*string, bool) {
	if o == nil || IsNil(o.Tier) {
		return nil, false
	}

	return o.Tier, true
}

// HasTier returns a boolean if a field has been set.
func (o *StreamsStartStreamProcessorWith) HasTier() bool {
	if o != nil && !IsNil(o.Tier) {
		return true
	}

	return false
}

// SetTier gets a reference to the given string and assigns it to the Tier field.
func (o *StreamsStartStreamProcessorWith) SetTier(v string) {
	o.Tier = &v
}
