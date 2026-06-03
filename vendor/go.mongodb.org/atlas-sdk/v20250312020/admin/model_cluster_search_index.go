// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ClusterSearchIndex struct for ClusterSearchIndex
type ClusterSearchIndex struct {
	// Human-readable label that identifies the collection that contains one or more Atlas Search indexes.
	CollectionName string `json:"collectionName"`
	// Human-readable label that identifies the database that contains the collection with one or more Atlas Search indexes.
	Database string `json:"database"`
	// Unique 24-hexadecimal digit string that identifies this Atlas Search index.
	// Read only field.
	IndexID *string `json:"indexID,omitempty"`
	// Human-readable label that identifies this index. Within each namespace, names of all indexes in the namespace must be unique.
	Name string `json:"name"`
	// Number of index partitions. Allowed values are [1, 2, 4].
	NumPartitions *int `json:"numPartitions,omitempty"`
	// Condition of the search index when you made this request.  - `IN_PROGRESS`: Atlas is building or re-building the index after an edit. - `STEADY`: You can use this search index. - `FAILED`: Atlas could not build the index. - `MIGRATING`: Atlas is upgrading the underlying cluster tier and migrating indexes. - `PAUSED`: The cluster is paused.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Flag that indicates whether to store all fields (true) on Atlas Search. By default, Atlas doesn't store (false) the fields on Atlas Search.  Alternatively, you can specify an object that only contains the list of fields to store (include) or not store (exclude) on Atlas Search. Note that storing all fields (true) is not allowed for vector search indexes. To learn more, see documentation.
	StoredSource any `json:"storedSource,omitempty"`
	// Type of the index. Default type is search.
	Type *string `json:"type,omitempty"`
	// Specific pre-defined method chosen to convert database field text into searchable words. This conversion reduces the text of fields into the smallest units of text. These units are called a **term** or **token**. This process, known as tokenization, involves a variety of changes made to the text in fields:  - extracting words - removing punctuation - removing accents - changing to lowercase - removing common words - reducing words to their root form (stemming) - changing words to their base form (lemmatization)  MongoDB Cloud uses the selected process to build the Atlas Search index.
	Analyzer *string `json:"analyzer,omitempty"`
	// List of user-defined methods to convert database field text into searchable words.
	Analyzers *[]ApiAtlasFTSAnalyzers `json:"analyzers,omitempty"`
	Mappings  *ApiAtlasFTSMappings    `json:"mappings,omitempty"`
	// Method applied to identify words when searching this index.
	SearchAnalyzer *string `json:"searchAnalyzer,omitempty"`
	// Rule sets that map words to their synonyms in this index.
	Synonyms *[]SearchSynonymMappingDefinition `json:"synonyms,omitempty"`
	// Settings that configure the fields, one per object, to index. You must define at least one \"vector\" type field. You can optionally define \"filter\" type fields also.
	Fields *[]any `json:"fields,omitempty"`
	// Top-level path to the array that contains vector fields. When provided, vector fields under this path are treated as nested.
	NestedRoot *string `json:"nestedRoot,omitempty"`
}

// NewClusterSearchIndex instantiates a new ClusterSearchIndex object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterSearchIndex(collectionName string, database string, name string) *ClusterSearchIndex {
	this := ClusterSearchIndex{}
	this.CollectionName = collectionName
	this.Database = database
	this.Name = name
	var numPartitions int = 1
	this.NumPartitions = &numPartitions
	var analyzer string = "lucene.standard"
	this.Analyzer = &analyzer
	var searchAnalyzer string = "lucene.standard"
	this.SearchAnalyzer = &searchAnalyzer
	return &this
}

// NewClusterSearchIndexWithDefaults instantiates a new ClusterSearchIndex object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterSearchIndexWithDefaults() *ClusterSearchIndex {
	this := ClusterSearchIndex{}
	var numPartitions int = 1
	this.NumPartitions = &numPartitions
	var analyzer string = "lucene.standard"
	this.Analyzer = &analyzer
	var searchAnalyzer string = "lucene.standard"
	this.SearchAnalyzer = &searchAnalyzer
	return &this
}

// GetCollectionName returns the CollectionName field value
func (o *ClusterSearchIndex) GetCollectionName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CollectionName
}

// GetCollectionNameOk returns a tuple with the CollectionName field value
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetCollectionNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CollectionName, true
}

// SetCollectionName sets field value
func (o *ClusterSearchIndex) SetCollectionName(v string) {
	o.CollectionName = v
}

// GetDatabase returns the Database field value
func (o *ClusterSearchIndex) GetDatabase() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Database
}

// GetDatabaseOk returns a tuple with the Database field value
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetDatabaseOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Database, true
}

// SetDatabase sets field value
func (o *ClusterSearchIndex) SetDatabase(v string) {
	o.Database = v
}

// GetIndexID returns the IndexID field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetIndexID() string {
	if o == nil || IsNil(o.IndexID) {
		var ret string
		return ret
	}
	return *o.IndexID
}

// GetIndexIDOk returns a tuple with the IndexID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetIndexIDOk() (*string, bool) {
	if o == nil || IsNil(o.IndexID) {
		return nil, false
	}

	return o.IndexID, true
}

// HasIndexID returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasIndexID() bool {
	if o != nil && !IsNil(o.IndexID) {
		return true
	}

	return false
}

// SetIndexID gets a reference to the given string and assigns it to the IndexID field.
func (o *ClusterSearchIndex) SetIndexID(v string) {
	o.IndexID = &v
}

// GetName returns the Name field value
func (o *ClusterSearchIndex) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *ClusterSearchIndex) SetName(v string) {
	o.Name = v
}

// GetNumPartitions returns the NumPartitions field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetNumPartitions() int {
	if o == nil || IsNil(o.NumPartitions) {
		var ret int
		return ret
	}
	return *o.NumPartitions
}

// GetNumPartitionsOk returns a tuple with the NumPartitions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetNumPartitionsOk() (*int, bool) {
	if o == nil || IsNil(o.NumPartitions) {
		return nil, false
	}

	return o.NumPartitions, true
}

// HasNumPartitions returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasNumPartitions() bool {
	if o != nil && !IsNil(o.NumPartitions) {
		return true
	}

	return false
}

// SetNumPartitions gets a reference to the given int and assigns it to the NumPartitions field.
func (o *ClusterSearchIndex) SetNumPartitions(v int) {
	o.NumPartitions = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *ClusterSearchIndex) SetStatus(v string) {
	o.Status = &v
}

// GetStoredSource returns the StoredSource field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetStoredSource() any {
	if o == nil || IsNil(o.StoredSource) {
		var ret any
		return ret
	}
	return o.StoredSource
}

// GetStoredSourceOk returns a tuple with the StoredSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetStoredSourceOk() (any, bool) {
	if o == nil || IsNil(o.StoredSource) {
		var ret any
		return ret, false
	}

	return o.StoredSource, true
}

// HasStoredSource returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasStoredSource() bool {
	if o != nil && !IsNil(o.StoredSource) {
		return true
	}

	return false
}

// SetStoredSource gets a reference to the given any and assigns it to the StoredSource field.
func (o *ClusterSearchIndex) SetStoredSource(v any) {
	o.StoredSource = v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ClusterSearchIndex) SetType(v string) {
	o.Type = &v
}

// GetAnalyzer returns the Analyzer field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetAnalyzer() string {
	if o == nil || IsNil(o.Analyzer) {
		var ret string
		return ret
	}
	return *o.Analyzer
}

// GetAnalyzerOk returns a tuple with the Analyzer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetAnalyzerOk() (*string, bool) {
	if o == nil || IsNil(o.Analyzer) {
		return nil, false
	}

	return o.Analyzer, true
}

// HasAnalyzer returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasAnalyzer() bool {
	if o != nil && !IsNil(o.Analyzer) {
		return true
	}

	return false
}

// SetAnalyzer gets a reference to the given string and assigns it to the Analyzer field.
func (o *ClusterSearchIndex) SetAnalyzer(v string) {
	o.Analyzer = &v
}

// GetAnalyzers returns the Analyzers field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetAnalyzers() []ApiAtlasFTSAnalyzers {
	if o == nil || IsNil(o.Analyzers) {
		var ret []ApiAtlasFTSAnalyzers
		return ret
	}
	return *o.Analyzers
}

// GetAnalyzersOk returns a tuple with the Analyzers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetAnalyzersOk() (*[]ApiAtlasFTSAnalyzers, bool) {
	if o == nil || IsNil(o.Analyzers) {
		return nil, false
	}

	return o.Analyzers, true
}

// HasAnalyzers returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasAnalyzers() bool {
	if o != nil && !IsNil(o.Analyzers) {
		return true
	}

	return false
}

// SetAnalyzers gets a reference to the given []ApiAtlasFTSAnalyzers and assigns it to the Analyzers field.
func (o *ClusterSearchIndex) SetAnalyzers(v []ApiAtlasFTSAnalyzers) {
	o.Analyzers = &v
}

// GetMappings returns the Mappings field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetMappings() ApiAtlasFTSMappings {
	if o == nil || IsNil(o.Mappings) {
		var ret ApiAtlasFTSMappings
		return ret
	}
	return *o.Mappings
}

// GetMappingsOk returns a tuple with the Mappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetMappingsOk() (*ApiAtlasFTSMappings, bool) {
	if o == nil || IsNil(o.Mappings) {
		return nil, false
	}

	return o.Mappings, true
}

// HasMappings returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasMappings() bool {
	if o != nil && !IsNil(o.Mappings) {
		return true
	}

	return false
}

// SetMappings gets a reference to the given ApiAtlasFTSMappings and assigns it to the Mappings field.
func (o *ClusterSearchIndex) SetMappings(v ApiAtlasFTSMappings) {
	o.Mappings = &v
}

// GetSearchAnalyzer returns the SearchAnalyzer field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetSearchAnalyzer() string {
	if o == nil || IsNil(o.SearchAnalyzer) {
		var ret string
		return ret
	}
	return *o.SearchAnalyzer
}

// GetSearchAnalyzerOk returns a tuple with the SearchAnalyzer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetSearchAnalyzerOk() (*string, bool) {
	if o == nil || IsNil(o.SearchAnalyzer) {
		return nil, false
	}

	return o.SearchAnalyzer, true
}

// HasSearchAnalyzer returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasSearchAnalyzer() bool {
	if o != nil && !IsNil(o.SearchAnalyzer) {
		return true
	}

	return false
}

// SetSearchAnalyzer gets a reference to the given string and assigns it to the SearchAnalyzer field.
func (o *ClusterSearchIndex) SetSearchAnalyzer(v string) {
	o.SearchAnalyzer = &v
}

// GetSynonyms returns the Synonyms field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetSynonyms() []SearchSynonymMappingDefinition {
	if o == nil || IsNil(o.Synonyms) {
		var ret []SearchSynonymMappingDefinition
		return ret
	}
	return *o.Synonyms
}

// GetSynonymsOk returns a tuple with the Synonyms field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetSynonymsOk() (*[]SearchSynonymMappingDefinition, bool) {
	if o == nil || IsNil(o.Synonyms) {
		return nil, false
	}

	return o.Synonyms, true
}

// HasSynonyms returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasSynonyms() bool {
	if o != nil && !IsNil(o.Synonyms) {
		return true
	}

	return false
}

// SetSynonyms gets a reference to the given []SearchSynonymMappingDefinition and assigns it to the Synonyms field.
func (o *ClusterSearchIndex) SetSynonyms(v []SearchSynonymMappingDefinition) {
	o.Synonyms = &v
}

// GetFields returns the Fields field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetFields() []any {
	if o == nil || IsNil(o.Fields) {
		var ret []any
		return ret
	}
	return *o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetFieldsOk() (*[]any, bool) {
	if o == nil || IsNil(o.Fields) {
		return nil, false
	}

	return o.Fields, true
}

// HasFields returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasFields() bool {
	if o != nil && !IsNil(o.Fields) {
		return true
	}

	return false
}

// SetFields gets a reference to the given []any and assigns it to the Fields field.
func (o *ClusterSearchIndex) SetFields(v []any) {
	o.Fields = &v
}

// GetNestedRoot returns the NestedRoot field value if set, zero value otherwise
func (o *ClusterSearchIndex) GetNestedRoot() string {
	if o == nil || IsNil(o.NestedRoot) {
		var ret string
		return ret
	}
	return *o.NestedRoot
}

// GetNestedRootOk returns a tuple with the NestedRoot field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterSearchIndex) GetNestedRootOk() (*string, bool) {
	if o == nil || IsNil(o.NestedRoot) {
		return nil, false
	}

	return o.NestedRoot, true
}

// HasNestedRoot returns a boolean if a field has been set.
func (o *ClusterSearchIndex) HasNestedRoot() bool {
	if o != nil && !IsNil(o.NestedRoot) {
		return true
	}

	return false
}

// SetNestedRoot gets a reference to the given string and assigns it to the NestedRoot field.
func (o *ClusterSearchIndex) SetNestedRoot(v string) {
	o.NestedRoot = &v
}
