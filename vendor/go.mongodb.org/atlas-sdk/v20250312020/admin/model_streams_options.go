// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsOptions Optional configuration for the stream processor.
type StreamsOptions struct {
	Dlq *StreamsDLQ `json:"dlq,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewStreamsOptions instantiates a new StreamsOptions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsOptions() *StreamsOptions {
	this := StreamsOptions{}
	return &this
}

// NewStreamsOptionsWithDefaults instantiates a new StreamsOptions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsOptionsWithDefaults() *StreamsOptions {
	this := StreamsOptions{}
	return &this
}

// GetDlq returns the Dlq field value if set, zero value otherwise
func (o *StreamsOptions) GetDlq() StreamsDLQ {
	if o == nil || IsNil(o.Dlq) {
		var ret StreamsDLQ
		return ret
	}
	return *o.Dlq
}

// GetDlqOk returns a tuple with the Dlq field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsOptions) GetDlqOk() (*StreamsDLQ, bool) {
	if o == nil || IsNil(o.Dlq) {
		return nil, false
	}

	return o.Dlq, true
}

// HasDlq returns a boolean if a field has been set.
func (o *StreamsOptions) HasDlq() bool {
	if o != nil && !IsNil(o.Dlq) {
		return true
	}

	return false
}

// SetDlq gets a reference to the given StreamsDLQ and assigns it to the Dlq field.
func (o *StreamsOptions) SetDlq(v StreamsDLQ) {
	o.Dlq = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsOptions) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsOptions) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsOptions) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsOptions) SetLinks(v []Link) {
	o.Links = &v
}
