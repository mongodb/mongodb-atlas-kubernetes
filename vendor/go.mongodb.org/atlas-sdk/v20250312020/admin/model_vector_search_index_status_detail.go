// Code based on the AtlasAPI V2 OpenAPI file

package admin

// VectorSearchIndexStatusDetail Contains status information about a vector search index.
type VectorSearchIndexStatusDetail struct {
	Definition        *VectorSearchIndexDefinition  `json:"definition,omitempty"`
	DefinitionVersion *SearchIndexDefinitionVersion `json:"definitionVersion,omitempty"`
	// Optional message describing an error.
	Message *string `json:"message,omitempty"`
	// Flag that indicates whether the index generation is queryable on the host.
	Queryable *bool `json:"queryable,omitempty"`
	// Condition of the search index when you made this request.  - `DELETING`: The index is being deleted. - `FAILED` The index build failed. Indexes can enter the FAILED state due to an invalid index definition. - `STALE`: The index is queryable but has stopped replicating data from the indexed collection. Searches on the index may return out-of-date data. - `PENDING`: Atlas has not yet started building the index. - `BUILDING`: Atlas is building or re-building the index after an edit. - `READY`: The index is ready and can support queries.
	Status *string `json:"status,omitempty"`
}

// NewVectorSearchIndexStatusDetail instantiates a new VectorSearchIndexStatusDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVectorSearchIndexStatusDetail() *VectorSearchIndexStatusDetail {
	this := VectorSearchIndexStatusDetail{}
	return &this
}

// NewVectorSearchIndexStatusDetailWithDefaults instantiates a new VectorSearchIndexStatusDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVectorSearchIndexStatusDetailWithDefaults() *VectorSearchIndexStatusDetail {
	this := VectorSearchIndexStatusDetail{}
	return &this
}

// GetDefinition returns the Definition field value if set, zero value otherwise
func (o *VectorSearchIndexStatusDetail) GetDefinition() VectorSearchIndexDefinition {
	if o == nil || IsNil(o.Definition) {
		var ret VectorSearchIndexDefinition
		return ret
	}
	return *o.Definition
}

// GetDefinitionOk returns a tuple with the Definition field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexStatusDetail) GetDefinitionOk() (*VectorSearchIndexDefinition, bool) {
	if o == nil || IsNil(o.Definition) {
		return nil, false
	}

	return o.Definition, true
}

// HasDefinition returns a boolean if a field has been set.
func (o *VectorSearchIndexStatusDetail) HasDefinition() bool {
	if o != nil && !IsNil(o.Definition) {
		return true
	}

	return false
}

// SetDefinition gets a reference to the given VectorSearchIndexDefinition and assigns it to the Definition field.
func (o *VectorSearchIndexStatusDetail) SetDefinition(v VectorSearchIndexDefinition) {
	o.Definition = &v
}

// GetDefinitionVersion returns the DefinitionVersion field value if set, zero value otherwise
func (o *VectorSearchIndexStatusDetail) GetDefinitionVersion() SearchIndexDefinitionVersion {
	if o == nil || IsNil(o.DefinitionVersion) {
		var ret SearchIndexDefinitionVersion
		return ret
	}
	return *o.DefinitionVersion
}

// GetDefinitionVersionOk returns a tuple with the DefinitionVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexStatusDetail) GetDefinitionVersionOk() (*SearchIndexDefinitionVersion, bool) {
	if o == nil || IsNil(o.DefinitionVersion) {
		return nil, false
	}

	return o.DefinitionVersion, true
}

// HasDefinitionVersion returns a boolean if a field has been set.
func (o *VectorSearchIndexStatusDetail) HasDefinitionVersion() bool {
	if o != nil && !IsNil(o.DefinitionVersion) {
		return true
	}

	return false
}

// SetDefinitionVersion gets a reference to the given SearchIndexDefinitionVersion and assigns it to the DefinitionVersion field.
func (o *VectorSearchIndexStatusDetail) SetDefinitionVersion(v SearchIndexDefinitionVersion) {
	o.DefinitionVersion = &v
}

// GetMessage returns the Message field value if set, zero value otherwise
func (o *VectorSearchIndexStatusDetail) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexStatusDetail) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}

	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *VectorSearchIndexStatusDetail) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *VectorSearchIndexStatusDetail) SetMessage(v string) {
	o.Message = &v
}

// GetQueryable returns the Queryable field value if set, zero value otherwise
func (o *VectorSearchIndexStatusDetail) GetQueryable() bool {
	if o == nil || IsNil(o.Queryable) {
		var ret bool
		return ret
	}
	return *o.Queryable
}

// GetQueryableOk returns a tuple with the Queryable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexStatusDetail) GetQueryableOk() (*bool, bool) {
	if o == nil || IsNil(o.Queryable) {
		return nil, false
	}

	return o.Queryable, true
}

// HasQueryable returns a boolean if a field has been set.
func (o *VectorSearchIndexStatusDetail) HasQueryable() bool {
	if o != nil && !IsNil(o.Queryable) {
		return true
	}

	return false
}

// SetQueryable gets a reference to the given bool and assigns it to the Queryable field.
func (o *VectorSearchIndexStatusDetail) SetQueryable(v bool) {
	o.Queryable = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *VectorSearchIndexStatusDetail) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexStatusDetail) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *VectorSearchIndexStatusDetail) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *VectorSearchIndexStatusDetail) SetStatus(v string) {
	o.Status = &v
}
