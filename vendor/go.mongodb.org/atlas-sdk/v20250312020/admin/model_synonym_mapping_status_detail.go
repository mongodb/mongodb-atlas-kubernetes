// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SynonymMappingStatusDetail Contains the status of the index's synonym mappings on each search host. This field (and its subfields) only appear if the index has synonyms defined.
type SynonymMappingStatusDetail struct {
	// Optional message describing an error.
	Message *string `json:"message,omitempty"`
	// Flag that indicates whether the synonym mapping is queryable on a host.
	Queryable *bool `json:"queryable,omitempty"`
	// Status that describes this index's synonym mappings. This status appears only if the index has synonyms defined.
	Status *string `json:"status,omitempty"`
}

// NewSynonymMappingStatusDetail instantiates a new SynonymMappingStatusDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSynonymMappingStatusDetail() *SynonymMappingStatusDetail {
	this := SynonymMappingStatusDetail{}
	return &this
}

// NewSynonymMappingStatusDetailWithDefaults instantiates a new SynonymMappingStatusDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSynonymMappingStatusDetailWithDefaults() *SynonymMappingStatusDetail {
	this := SynonymMappingStatusDetail{}
	return &this
}

// GetMessage returns the Message field value if set, zero value otherwise
func (o *SynonymMappingStatusDetail) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SynonymMappingStatusDetail) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}

	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *SynonymMappingStatusDetail) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *SynonymMappingStatusDetail) SetMessage(v string) {
	o.Message = &v
}

// GetQueryable returns the Queryable field value if set, zero value otherwise
func (o *SynonymMappingStatusDetail) GetQueryable() bool {
	if o == nil || IsNil(o.Queryable) {
		var ret bool
		return ret
	}
	return *o.Queryable
}

// GetQueryableOk returns a tuple with the Queryable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SynonymMappingStatusDetail) GetQueryableOk() (*bool, bool) {
	if o == nil || IsNil(o.Queryable) {
		return nil, false
	}

	return o.Queryable, true
}

// HasQueryable returns a boolean if a field has been set.
func (o *SynonymMappingStatusDetail) HasQueryable() bool {
	if o != nil && !IsNil(o.Queryable) {
		return true
	}

	return false
}

// SetQueryable gets a reference to the given bool and assigns it to the Queryable field.
func (o *SynonymMappingStatusDetail) SetQueryable(v bool) {
	o.Queryable = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *SynonymMappingStatusDetail) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SynonymMappingStatusDetail) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *SynonymMappingStatusDetail) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *SynonymMappingStatusDetail) SetStatus(v string) {
	o.Status = &v
}
