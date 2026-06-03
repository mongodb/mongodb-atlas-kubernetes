// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakeApiBase An aggregation pipeline that applies to the collection.
type DataLakeApiBase struct {
	// Human-readable label that identifies the view, which corresponds to an aggregation pipeline on a collection.
	Name *string `json:"name,omitempty"`
	// Aggregation pipeline stages to apply to the source collection.
	Pipeline *string `json:"pipeline,omitempty"`
	// Human-readable label that identifies the source collection for the view.
	Source *string `json:"source,omitempty"`
}

// NewDataLakeApiBase instantiates a new DataLakeApiBase object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakeApiBase() *DataLakeApiBase {
	this := DataLakeApiBase{}
	return &this
}

// NewDataLakeApiBaseWithDefaults instantiates a new DataLakeApiBase object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakeApiBaseWithDefaults() *DataLakeApiBase {
	this := DataLakeApiBase{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise
func (o *DataLakeApiBase) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeApiBase) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataLakeApiBase) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataLakeApiBase) SetName(v string) {
	o.Name = &v
}

// GetPipeline returns the Pipeline field value if set, zero value otherwise
func (o *DataLakeApiBase) GetPipeline() string {
	if o == nil || IsNil(o.Pipeline) {
		var ret string
		return ret
	}
	return *o.Pipeline
}

// GetPipelineOk returns a tuple with the Pipeline field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeApiBase) GetPipelineOk() (*string, bool) {
	if o == nil || IsNil(o.Pipeline) {
		return nil, false
	}

	return o.Pipeline, true
}

// HasPipeline returns a boolean if a field has been set.
func (o *DataLakeApiBase) HasPipeline() bool {
	if o != nil && !IsNil(o.Pipeline) {
		return true
	}

	return false
}

// SetPipeline gets a reference to the given string and assigns it to the Pipeline field.
func (o *DataLakeApiBase) SetPipeline(v string) {
	o.Pipeline = &v
}

// GetSource returns the Source field value if set, zero value otherwise
func (o *DataLakeApiBase) GetSource() string {
	if o == nil || IsNil(o.Source) {
		var ret string
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataLakeApiBase) GetSourceOk() (*string, bool) {
	if o == nil || IsNil(o.Source) {
		return nil, false
	}

	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *DataLakeApiBase) HasSource() bool {
	if o != nil && !IsNil(o.Source) {
		return true
	}

	return false
}

// SetSource gets a reference to the given string and assigns it to the Source field.
func (o *DataLakeApiBase) SetSource(v string) {
	o.Source = &v
}
