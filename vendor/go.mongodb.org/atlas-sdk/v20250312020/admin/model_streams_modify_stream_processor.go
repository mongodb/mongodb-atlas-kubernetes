// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsModifyStreamProcessor A request to modify an existing stream processor.
type StreamsModifyStreamProcessor struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// New name for the stream processor.
	Name    *string                              `json:"name,omitempty"`
	Options *StreamsModifyStreamProcessorOptions `json:"options,omitempty"`
	// New pipeline for the stream processor.
	Pipeline *[]any `json:"pipeline,omitempty"`
}

// NewStreamsModifyStreamProcessor instantiates a new StreamsModifyStreamProcessor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsModifyStreamProcessor() *StreamsModifyStreamProcessor {
	this := StreamsModifyStreamProcessor{}
	return &this
}

// NewStreamsModifyStreamProcessorWithDefaults instantiates a new StreamsModifyStreamProcessor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsModifyStreamProcessorWithDefaults() *StreamsModifyStreamProcessor {
	this := StreamsModifyStreamProcessor{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsModifyStreamProcessor) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsModifyStreamProcessor) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsModifyStreamProcessor) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsModifyStreamProcessor) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *StreamsModifyStreamProcessor) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsModifyStreamProcessor) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *StreamsModifyStreamProcessor) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *StreamsModifyStreamProcessor) SetName(v string) {
	o.Name = &v
}

// GetOptions returns the Options field value if set, zero value otherwise
func (o *StreamsModifyStreamProcessor) GetOptions() StreamsModifyStreamProcessorOptions {
	if o == nil || IsNil(o.Options) {
		var ret StreamsModifyStreamProcessorOptions
		return ret
	}
	return *o.Options
}

// GetOptionsOk returns a tuple with the Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsModifyStreamProcessor) GetOptionsOk() (*StreamsModifyStreamProcessorOptions, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}

	return o.Options, true
}

// HasOptions returns a boolean if a field has been set.
func (o *StreamsModifyStreamProcessor) HasOptions() bool {
	if o != nil && !IsNil(o.Options) {
		return true
	}

	return false
}

// SetOptions gets a reference to the given StreamsModifyStreamProcessorOptions and assigns it to the Options field.
func (o *StreamsModifyStreamProcessor) SetOptions(v StreamsModifyStreamProcessorOptions) {
	o.Options = &v
}

// GetPipeline returns the Pipeline field value if set, zero value otherwise
func (o *StreamsModifyStreamProcessor) GetPipeline() []any {
	if o == nil || IsNil(o.Pipeline) {
		var ret []any
		return ret
	}
	return *o.Pipeline
}

// GetPipelineOk returns a tuple with the Pipeline field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsModifyStreamProcessor) GetPipelineOk() (*[]any, bool) {
	if o == nil || IsNil(o.Pipeline) {
		return nil, false
	}

	return o.Pipeline, true
}

// HasPipeline returns a boolean if a field has been set.
func (o *StreamsModifyStreamProcessor) HasPipeline() bool {
	if o != nil && !IsNil(o.Pipeline) {
		return true
	}

	return false
}

// SetPipeline gets a reference to the given []any and assigns it to the Pipeline field.
func (o *StreamsModifyStreamProcessor) SetPipeline(v []any) {
	o.Pipeline = &v
}
