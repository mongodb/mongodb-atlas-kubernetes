// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsProcessorWithStats An atlas stream processor with optional stats.
type StreamsProcessorWithStats struct {
	// Unique 24-hexadecimal character string that identifies the stream processor.
	// Read only field.
	Id string `json:"_id"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable name of the stream processor.
	// Read only field.
	Name    string          `json:"name"`
	Options *StreamsOptions `json:"options,omitempty"`
	// Stream aggregation pipeline you want to apply to your streaming data.
	// Read only field.
	Pipeline []any `json:"pipeline"`
	// The state of the stream processor. Commonly occurring states are 'CREATED', 'STARTED', 'STOPPED' and 'FAILED'.
	// Read only field.
	State string `json:"state"`
	// The stats associated with the stream processor.
	// Read only field.
	Stats any `json:"stats,omitempty"`
	// Selected tier for the Stream Workspace. Configures Memory / VCPU allowances.
	Tier *string `json:"tier,omitempty"`
}

// NewStreamsProcessorWithStats instantiates a new StreamsProcessorWithStats object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsProcessorWithStats(id string, name string, pipeline []any, state string) *StreamsProcessorWithStats {
	this := StreamsProcessorWithStats{}
	this.Id = id
	this.Name = name
	this.Pipeline = pipeline
	this.State = state
	return &this
}

// NewStreamsProcessorWithStatsWithDefaults instantiates a new StreamsProcessorWithStats object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsProcessorWithStatsWithDefaults() *StreamsProcessorWithStats {
	this := StreamsProcessorWithStats{}
	return &this
}

// GetId returns the Id field value
func (o *StreamsProcessorWithStats) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *StreamsProcessorWithStats) SetId(v string) {
	o.Id = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *StreamsProcessorWithStats) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *StreamsProcessorWithStats) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *StreamsProcessorWithStats) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value
func (o *StreamsProcessorWithStats) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *StreamsProcessorWithStats) SetName(v string) {
	o.Name = v
}

// GetOptions returns the Options field value if set, zero value otherwise
func (o *StreamsProcessorWithStats) GetOptions() StreamsOptions {
	if o == nil || IsNil(o.Options) {
		var ret StreamsOptions
		return ret
	}
	return *o.Options
}

// GetOptionsOk returns a tuple with the Options field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetOptionsOk() (*StreamsOptions, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}

	return o.Options, true
}

// HasOptions returns a boolean if a field has been set.
func (o *StreamsProcessorWithStats) HasOptions() bool {
	if o != nil && !IsNil(o.Options) {
		return true
	}

	return false
}

// SetOptions gets a reference to the given StreamsOptions and assigns it to the Options field.
func (o *StreamsProcessorWithStats) SetOptions(v StreamsOptions) {
	o.Options = &v
}

// GetPipeline returns the Pipeline field value
func (o *StreamsProcessorWithStats) GetPipeline() []any {
	if o == nil {
		var ret []any
		return ret
	}

	return o.Pipeline
}

// GetPipelineOk returns a tuple with the Pipeline field value
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetPipelineOk() (*[]any, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Pipeline, true
}

// SetPipeline sets field value
func (o *StreamsProcessorWithStats) SetPipeline(v []any) {
	o.Pipeline = v
}

// GetState returns the State field value
func (o *StreamsProcessorWithStats) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *StreamsProcessorWithStats) SetState(v string) {
	o.State = v
}

// GetStats returns the Stats field value if set, zero value otherwise
func (o *StreamsProcessorWithStats) GetStats() any {
	if o == nil || IsNil(o.Stats) {
		var ret any
		return ret
	}
	return o.Stats
}

// GetStatsOk returns a tuple with the Stats field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetStatsOk() (any, bool) {
	if o == nil || IsNil(o.Stats) {
		var ret any
		return ret, false
	}

	return o.Stats, true
}

// HasStats returns a boolean if a field has been set.
func (o *StreamsProcessorWithStats) HasStats() bool {
	if o != nil && !IsNil(o.Stats) {
		return true
	}

	return false
}

// SetStats gets a reference to the given any and assigns it to the Stats field.
func (o *StreamsProcessorWithStats) SetStats(v any) {
	o.Stats = v
}

// GetTier returns the Tier field value if set, zero value otherwise
func (o *StreamsProcessorWithStats) GetTier() string {
	if o == nil || IsNil(o.Tier) {
		var ret string
		return ret
	}
	return *o.Tier
}

// GetTierOk returns a tuple with the Tier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreamsProcessorWithStats) GetTierOk() (*string, bool) {
	if o == nil || IsNil(o.Tier) {
		return nil, false
	}

	return o.Tier, true
}

// HasTier returns a boolean if a field has been set.
func (o *StreamsProcessorWithStats) HasTier() bool {
	if o != nil && !IsNil(o.Tier) {
		return true
	}

	return false
}

// SetTier gets a reference to the given string and assigns it to the Tier field.
func (o *StreamsProcessorWithStats) SetTier(v string) {
	o.Tier = &v
}
