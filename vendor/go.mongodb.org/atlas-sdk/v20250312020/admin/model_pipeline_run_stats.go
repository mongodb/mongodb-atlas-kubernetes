// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PipelineRunStats Runtime statistics for this Data Lake Pipeline run.
type PipelineRunStats struct {
	// Total data size in bytes exported for this pipeline run.
	// Read only field.
	BytesExported *int64 `json:"bytesExported,omitempty"`
	// Number of docs ingested for a this pipeline run.
	// Read only field.
	NumDocs *int64 `json:"numDocs,omitempty"`
}

// NewPipelineRunStats instantiates a new PipelineRunStats object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineRunStats() *PipelineRunStats {
	this := PipelineRunStats{}
	return &this
}

// NewPipelineRunStatsWithDefaults instantiates a new PipelineRunStats object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineRunStatsWithDefaults() *PipelineRunStats {
	this := PipelineRunStats{}
	return &this
}

// GetBytesExported returns the BytesExported field value if set, zero value otherwise
func (o *PipelineRunStats) GetBytesExported() int64 {
	if o == nil || IsNil(o.BytesExported) {
		var ret int64
		return ret
	}
	return *o.BytesExported
}

// GetBytesExportedOk returns a tuple with the BytesExported field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineRunStats) GetBytesExportedOk() (*int64, bool) {
	if o == nil || IsNil(o.BytesExported) {
		return nil, false
	}

	return o.BytesExported, true
}

// HasBytesExported returns a boolean if a field has been set.
func (o *PipelineRunStats) HasBytesExported() bool {
	if o != nil && !IsNil(o.BytesExported) {
		return true
	}

	return false
}

// SetBytesExported gets a reference to the given int64 and assigns it to the BytesExported field.
func (o *PipelineRunStats) SetBytesExported(v int64) {
	o.BytesExported = &v
}

// GetNumDocs returns the NumDocs field value if set, zero value otherwise
func (o *PipelineRunStats) GetNumDocs() int64 {
	if o == nil || IsNil(o.NumDocs) {
		var ret int64
		return ret
	}
	return *o.NumDocs
}

// GetNumDocsOk returns a tuple with the NumDocs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineRunStats) GetNumDocsOk() (*int64, bool) {
	if o == nil || IsNil(o.NumDocs) {
		return nil, false
	}

	return o.NumDocs, true
}

// HasNumDocs returns a boolean if a field has been set.
func (o *PipelineRunStats) HasNumDocs() bool {
	if o != nil && !IsNil(o.NumDocs) {
		return true
	}

	return false
}

// SetNumDocs gets a reference to the given int64 and assigns it to the NumDocs field.
func (o *PipelineRunStats) SetNumDocs(v int64) {
	o.NumDocs = &v
}
