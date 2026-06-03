// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsSampleConnections Sample connections to add to SPI.
type StreamsSampleConnections struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Flag that indicates whether to add a `sample_stream_solar` connection.
	Solar *bool `json:"solar,omitempty"`
}

// NewStreamsSampleConnections instantiates a new StreamsSampleConnections object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsSampleConnections() *StreamsSampleConnections {
	this := StreamsSampleConnections{}
	var solar bool = false
	this.Solar = &solar
	return &this
}

// NewStreamsSampleConnectionsWithDefaults instantiates a new StreamsSampleConnections object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsSampleConnectionsWithDefaults() *StreamsSampleConnections {
	this := StreamsSampleConnections{}
	var solar bool = false
	this.Solar = &solar
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsSampleConnections) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsSampleConnections) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsSampleConnections) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsSampleConnections) SetLinks(v []Link) {
	o.Links = &v
}

// GetSolar returns the Solar field value if set, zero value otherwise
func (o *StreamsSampleConnections) GetSolar() bool {
	if o == nil || IsNil(o.Solar) {
		var ret bool
		return ret
	}
	return *o.Solar
}

// GetSolarOk returns a tuple with the Solar field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsSampleConnections) GetSolarOk() (*bool, bool) {
	if o == nil || IsNil(o.Solar) {
		return nil, false
	}

	return o.Solar, true
}

// HasSolar returns a boolean if a field has been set.
func (o *StreamsSampleConnections) HasSolar() bool {
	if o != nil && !IsNil(o.Solar) {
		return true
	}

	return false
}

// SetSolar gets a reference to the given bool and assigns it to the Solar field.
func (o *StreamsSampleConnections) SetSolar(v bool) {
	o.Solar = &v
}
