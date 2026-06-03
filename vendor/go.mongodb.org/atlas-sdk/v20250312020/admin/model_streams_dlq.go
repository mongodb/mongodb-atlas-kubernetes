// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsDLQ Dead letter queue for the stream processor.
type StreamsDLQ struct {
	// Name of the collection to use for the DLQ.
	Coll *string `json:"coll,omitempty"`
	// Name of the connection to write DLQ messages to. Must be an Atlas connection.
	ConnectionName *string `json:"connectionName,omitempty"`
	// Name of the database to use for the DLQ.
	Db *string `json:"db,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewStreamsDLQ instantiates a new StreamsDLQ object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsDLQ() *StreamsDLQ {
	this := StreamsDLQ{}
	return &this
}

// NewStreamsDLQWithDefaults instantiates a new StreamsDLQ object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsDLQWithDefaults() *StreamsDLQ {
	this := StreamsDLQ{}
	return &this
}

// GetColl returns the Coll field value if set, zero value otherwise
func (o *StreamsDLQ) GetColl() string {
	if o == nil || IsNil(o.Coll) {
		var ret string
		return ret
	}
	return *o.Coll
}

// GetCollOk returns a tuple with the Coll field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsDLQ) GetCollOk() (*string, bool) {
	if o == nil || IsNil(o.Coll) {
		return nil, false
	}

	return o.Coll, true
}

// HasColl returns a boolean if a field has been set.
func (o *StreamsDLQ) HasColl() bool {
	if o != nil && !IsNil(o.Coll) {
		return true
	}

	return false
}

// SetColl gets a reference to the given string and assigns it to the Coll field.
func (o *StreamsDLQ) SetColl(v string) {
	o.Coll = &v
}

// GetConnectionName returns the ConnectionName field value if set, zero value otherwise
func (o *StreamsDLQ) GetConnectionName() string {
	if o == nil || IsNil(o.ConnectionName) {
		var ret string
		return ret
	}
	return *o.ConnectionName
}

// GetConnectionNameOk returns a tuple with the ConnectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsDLQ) GetConnectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectionName) {
		return nil, false
	}

	return o.ConnectionName, true
}

// HasConnectionName returns a boolean if a field has been set.
func (o *StreamsDLQ) HasConnectionName() bool {
	if o != nil && !IsNil(o.ConnectionName) {
		return true
	}

	return false
}

// SetConnectionName gets a reference to the given string and assigns it to the ConnectionName field.
func (o *StreamsDLQ) SetConnectionName(v string) {
	o.ConnectionName = &v
}

// GetDb returns the Db field value if set, zero value otherwise
func (o *StreamsDLQ) GetDb() string {
	if o == nil || IsNil(o.Db) {
		var ret string
		return ret
	}
	return *o.Db
}

// GetDbOk returns a tuple with the Db field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsDLQ) GetDbOk() (*string, bool) {
	if o == nil || IsNil(o.Db) {
		return nil, false
	}

	return o.Db, true
}

// HasDb returns a boolean if a field has been set.
func (o *StreamsDLQ) HasDb() bool {
	if o != nil && !IsNil(o.Db) {
		return true
	}

	return false
}

// SetDb gets a reference to the given string and assigns it to the Db field.
func (o *StreamsDLQ) SetDb(v string) {
	o.Db = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsDLQ) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsDLQ) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsDLQ) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsDLQ) SetLinks(v []Link) {
	o.Links = &v
}
