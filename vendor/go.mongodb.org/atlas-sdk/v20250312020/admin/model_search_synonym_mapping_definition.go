// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SearchSynonymMappingDefinition Synonyms used for this full text index.
type SearchSynonymMappingDefinition struct {
	// Specific pre-defined method chosen to apply to the synonyms to be searched.
	Analyzer string `json:"analyzer"`
	// Label that identifies the synonym definition. Each `synonym.name` must be unique within the same index definition.
	Name   string        `json:"name"`
	Source SynonymSource `json:"source"`
}

// NewSearchSynonymMappingDefinition instantiates a new SearchSynonymMappingDefinition object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSearchSynonymMappingDefinition(analyzer string, name string, source SynonymSource) *SearchSynonymMappingDefinition {
	this := SearchSynonymMappingDefinition{}
	this.Analyzer = analyzer
	this.Name = name
	this.Source = source
	return &this
}

// NewSearchSynonymMappingDefinitionWithDefaults instantiates a new SearchSynonymMappingDefinition object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSearchSynonymMappingDefinitionWithDefaults() *SearchSynonymMappingDefinition {
	this := SearchSynonymMappingDefinition{}
	return &this
}

// GetAnalyzer returns the Analyzer field value
func (o *SearchSynonymMappingDefinition) GetAnalyzer() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Analyzer
}

// GetAnalyzerOk returns a tuple with the Analyzer field value
// and a boolean to check if the value has been set.
func (o *SearchSynonymMappingDefinition) GetAnalyzerOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Analyzer, true
}

// SetAnalyzer sets field value
func (o *SearchSynonymMappingDefinition) SetAnalyzer(v string) {
	o.Analyzer = v
}

// GetName returns the Name field value
func (o *SearchSynonymMappingDefinition) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *SearchSynonymMappingDefinition) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *SearchSynonymMappingDefinition) SetName(v string) {
	o.Name = v
}

// GetSource returns the Source field value
func (o *SearchSynonymMappingDefinition) GetSource() SynonymSource {
	if o == nil {
		var ret SynonymSource
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *SearchSynonymMappingDefinition) GetSourceOk() (*SynonymSource, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value
func (o *SearchSynonymMappingDefinition) SetSource(v SynonymSource) {
	o.Source = v
}
