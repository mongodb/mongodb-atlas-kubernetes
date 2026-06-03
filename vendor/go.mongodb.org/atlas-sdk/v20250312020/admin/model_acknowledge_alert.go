// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// AcknowledgeAlert Acknowledging an alert prevents successive notifications. Specify the `acknowledgeUntil` date and optional comment or `unacknowledgeAlert` boolean.
type AcknowledgeAlert struct {
	// Date and time until which this alert has been acknowledged. This parameter expresses its value in the ISO 8601 timestamp format in UTC. The resource returns this parameter if a MongoDB User previously acknowledged this alert.
	AcknowledgedUntil *time.Time `json:"acknowledgedUntil,omitempty"`
	// Comment that a MongoDB Cloud user submitted when acknowledging the alert.
	AcknowledgementComment *string `json:"acknowledgementComment,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Flag that indicates to unacknowledge a previously acknowledged alert. By default this value is set to false. If set to true, it will override the `acknowledgedUntil` parameter.
	UnacknowledgeAlert *bool `json:"unacknowledgeAlert,omitempty"`
}

// NewAcknowledgeAlert instantiates a new AcknowledgeAlert object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAcknowledgeAlert() *AcknowledgeAlert {
	this := AcknowledgeAlert{}
	return &this
}

// NewAcknowledgeAlertWithDefaults instantiates a new AcknowledgeAlert object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAcknowledgeAlertWithDefaults() *AcknowledgeAlert {
	this := AcknowledgeAlert{}
	return &this
}

// GetAcknowledgedUntil returns the AcknowledgedUntil field value if set, zero value otherwise
func (o *AcknowledgeAlert) GetAcknowledgedUntil() time.Time {
	if o == nil || IsNil(o.AcknowledgedUntil) {
		var ret time.Time
		return ret
	}
	return *o.AcknowledgedUntil
}

// GetAcknowledgedUntilOk returns a tuple with the AcknowledgedUntil field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AcknowledgeAlert) GetAcknowledgedUntilOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AcknowledgedUntil) {
		return nil, false
	}

	return o.AcknowledgedUntil, true
}

// HasAcknowledgedUntil returns a boolean if a field has been set.
func (o *AcknowledgeAlert) HasAcknowledgedUntil() bool {
	if o != nil && !IsNil(o.AcknowledgedUntil) {
		return true
	}

	return false
}

// SetAcknowledgedUntil gets a reference to the given time.Time and assigns it to the AcknowledgedUntil field.
func (o *AcknowledgeAlert) SetAcknowledgedUntil(v time.Time) {
	o.AcknowledgedUntil = &v
}

// GetAcknowledgementComment returns the AcknowledgementComment field value if set, zero value otherwise
func (o *AcknowledgeAlert) GetAcknowledgementComment() string {
	if o == nil || IsNil(o.AcknowledgementComment) {
		var ret string
		return ret
	}
	return *o.AcknowledgementComment
}

// GetAcknowledgementCommentOk returns a tuple with the AcknowledgementComment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AcknowledgeAlert) GetAcknowledgementCommentOk() (*string, bool) {
	if o == nil || IsNil(o.AcknowledgementComment) {
		return nil, false
	}

	return o.AcknowledgementComment, true
}

// HasAcknowledgementComment returns a boolean if a field has been set.
func (o *AcknowledgeAlert) HasAcknowledgementComment() bool {
	if o != nil && !IsNil(o.AcknowledgementComment) {
		return true
	}

	return false
}

// SetAcknowledgementComment gets a reference to the given string and assigns it to the AcknowledgementComment field.
func (o *AcknowledgeAlert) SetAcknowledgementComment(v string) {
	o.AcknowledgementComment = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *AcknowledgeAlert) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AcknowledgeAlert) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *AcknowledgeAlert) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *AcknowledgeAlert) SetLinks(v []Link) {
	o.Links = &v
}

// GetUnacknowledgeAlert returns the UnacknowledgeAlert field value if set, zero value otherwise
func (o *AcknowledgeAlert) GetUnacknowledgeAlert() bool {
	if o == nil || IsNil(o.UnacknowledgeAlert) {
		var ret bool
		return ret
	}
	return *o.UnacknowledgeAlert
}

// GetUnacknowledgeAlertOk returns a tuple with the UnacknowledgeAlert field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AcknowledgeAlert) GetUnacknowledgeAlertOk() (*bool, bool) {
	if o == nil || IsNil(o.UnacknowledgeAlert) {
		return nil, false
	}

	return o.UnacknowledgeAlert, true
}

// HasUnacknowledgeAlert returns a boolean if a field has been set.
func (o *AcknowledgeAlert) HasUnacknowledgeAlert() bool {
	if o != nil && !IsNil(o.UnacknowledgeAlert) {
		return true
	}

	return false
}

// SetUnacknowledgeAlert gets a reference to the given bool and assigns it to the UnacknowledgeAlert field.
func (o *AcknowledgeAlert) SetUnacknowledgeAlert(v bool) {
	o.UnacknowledgeAlert = &v
}
