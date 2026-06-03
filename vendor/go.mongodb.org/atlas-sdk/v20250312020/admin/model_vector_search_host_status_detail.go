// Code based on the AtlasAPI V2 OpenAPI file

package admin

// VectorSearchHostStatusDetail struct for VectorSearchHostStatusDetail
type VectorSearchHostStatusDetail struct {
	// Hostname that corresponds to the status detail.
	Hostname  *string                        `json:"hostname,omitempty"`
	MainIndex *VectorSearchIndexStatusDetail `json:"mainIndex,omitempty"`
	// Flag that indicates whether the index is queryable on the host.
	Queryable   *bool                          `json:"queryable,omitempty"`
	StagedIndex *VectorSearchIndexStatusDetail `json:"stagedIndex,omitempty"`
	// Condition of the search index when you made this request.  - `DELETING`: The index is being deleted. - `FAILED` The index build failed. Indexes can enter the FAILED state due to an invalid index definition. - `STALE`: The index is queryable but has stopped replicating data from the indexed collection. Searches on the index may return out-of-date data. - `PENDING`: Atlas has not yet started building the index. - `BUILDING`: Atlas is building or re-building the index after an edit. - `READY`: The index is ready and can support queries.
	Status *string `json:"status,omitempty"`
}

// NewVectorSearchHostStatusDetail instantiates a new VectorSearchHostStatusDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVectorSearchHostStatusDetail() *VectorSearchHostStatusDetail {
	this := VectorSearchHostStatusDetail{}
	return &this
}

// NewVectorSearchHostStatusDetailWithDefaults instantiates a new VectorSearchHostStatusDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVectorSearchHostStatusDetailWithDefaults() *VectorSearchHostStatusDetail {
	this := VectorSearchHostStatusDetail{}
	return &this
}

// GetHostname returns the Hostname field value if set, zero value otherwise
func (o *VectorSearchHostStatusDetail) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchHostStatusDetail) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}

	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *VectorSearchHostStatusDetail) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *VectorSearchHostStatusDetail) SetHostname(v string) {
	o.Hostname = &v
}

// GetMainIndex returns the MainIndex field value if set, zero value otherwise
func (o *VectorSearchHostStatusDetail) GetMainIndex() VectorSearchIndexStatusDetail {
	if o == nil || IsNil(o.MainIndex) {
		var ret VectorSearchIndexStatusDetail
		return ret
	}
	return *o.MainIndex
}

// GetMainIndexOk returns a tuple with the MainIndex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchHostStatusDetail) GetMainIndexOk() (*VectorSearchIndexStatusDetail, bool) {
	if o == nil || IsNil(o.MainIndex) {
		return nil, false
	}

	return o.MainIndex, true
}

// HasMainIndex returns a boolean if a field has been set.
func (o *VectorSearchHostStatusDetail) HasMainIndex() bool {
	if o != nil && !IsNil(o.MainIndex) {
		return true
	}

	return false
}

// SetMainIndex gets a reference to the given VectorSearchIndexStatusDetail and assigns it to the MainIndex field.
func (o *VectorSearchHostStatusDetail) SetMainIndex(v VectorSearchIndexStatusDetail) {
	o.MainIndex = &v
}

// GetQueryable returns the Queryable field value if set, zero value otherwise
func (o *VectorSearchHostStatusDetail) GetQueryable() bool {
	if o == nil || IsNil(o.Queryable) {
		var ret bool
		return ret
	}
	return *o.Queryable
}

// GetQueryableOk returns a tuple with the Queryable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchHostStatusDetail) GetQueryableOk() (*bool, bool) {
	if o == nil || IsNil(o.Queryable) {
		return nil, false
	}

	return o.Queryable, true
}

// HasQueryable returns a boolean if a field has been set.
func (o *VectorSearchHostStatusDetail) HasQueryable() bool {
	if o != nil && !IsNil(o.Queryable) {
		return true
	}

	return false
}

// SetQueryable gets a reference to the given bool and assigns it to the Queryable field.
func (o *VectorSearchHostStatusDetail) SetQueryable(v bool) {
	o.Queryable = &v
}

// GetStagedIndex returns the StagedIndex field value if set, zero value otherwise
func (o *VectorSearchHostStatusDetail) GetStagedIndex() VectorSearchIndexStatusDetail {
	if o == nil || IsNil(o.StagedIndex) {
		var ret VectorSearchIndexStatusDetail
		return ret
	}
	return *o.StagedIndex
}

// GetStagedIndexOk returns a tuple with the StagedIndex field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchHostStatusDetail) GetStagedIndexOk() (*VectorSearchIndexStatusDetail, bool) {
	if o == nil || IsNil(o.StagedIndex) {
		return nil, false
	}

	return o.StagedIndex, true
}

// HasStagedIndex returns a boolean if a field has been set.
func (o *VectorSearchHostStatusDetail) HasStagedIndex() bool {
	if o != nil && !IsNil(o.StagedIndex) {
		return true
	}

	return false
}

// SetStagedIndex gets a reference to the given VectorSearchIndexStatusDetail and assigns it to the StagedIndex field.
func (o *VectorSearchHostStatusDetail) SetStagedIndex(v VectorSearchIndexStatusDetail) {
	o.StagedIndex = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *VectorSearchHostStatusDetail) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchHostStatusDetail) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *VectorSearchHostStatusDetail) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *VectorSearchHostStatusDetail) SetStatus(v string) {
	o.Status = &v
}
