// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SearchIndexResponse struct for SearchIndexResponse
type SearchIndexResponse struct {
	// Label that identifies the collection that contains one or more Atlas Search indexes.
	CollectionName *string `json:"collectionName,omitempty"`
	// Label that identifies the database that contains the collection with one or more Atlas Search indexes.
	Database *string `json:"database,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this Atlas Search index.
	IndexID                 *string                                  `json:"indexID,omitempty"`
	LatestDefinition        *BaseSearchIndexResponseLatestDefinition `json:"latestDefinition,omitempty"`
	LatestDefinitionVersion *SearchIndexDefinitionVersion            `json:"latestDefinitionVersion,omitempty"`
	// Label that identifies this index. Within each namespace, the names of all indexes must be unique.
	Name *string `json:"name,omitempty"`
	// Flag that indicates whether the index is queryable on all hosts.
	Queryable *bool `json:"queryable,omitempty"`
	// Condition of the search index when you made this request.  - `DELETING`: The index is being deleted. - `FAILED` The index build failed. Indexes can enter the FAILED state due to an invalid index definition. - `STALE`: The index is queryable but has stopped replicating data from the indexed collection. Searches on the index may return out-of-date data. - `PENDING`: Atlas has not yet started building the index. - `BUILDING`: Atlas is building or re-building the index after an edit. - `READY`: The index is ready and can support queries.
	Status *string `json:"status,omitempty"`
	// List of documents detailing index status on each host.
	StatusDetail *[]VectorSearchHostStatusDetail `json:"statusDetail,omitempty"`
	// Type of the index. The default type is search.
	Type *string `json:"type,omitempty"`
	// Status that describes this index's synonym mappings. This status appears only if the index has synonyms defined.
	SynonymMappingStatus *string `json:"synonymMappingStatus,omitempty"`
	// A list of documents describing the status of the index's synonym mappings on each search host. Only appears if the index has synonyms defined.
	SynonymMappingStatusDetail *[]map[string]SynonymMappingStatusDetail `json:"synonymMappingStatusDetail,omitempty"`
}

// NewSearchIndexResponse instantiates a new SearchIndexResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSearchIndexResponse() *SearchIndexResponse {
	this := SearchIndexResponse{}
	return &this
}

// NewSearchIndexResponseWithDefaults instantiates a new SearchIndexResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSearchIndexResponseWithDefaults() *SearchIndexResponse {
	this := SearchIndexResponse{}
	return &this
}

// GetCollectionName returns the CollectionName field value if set, zero value otherwise
func (o *SearchIndexResponse) GetCollectionName() string {
	if o == nil || IsNil(o.CollectionName) {
		var ret string
		return ret
	}
	return *o.CollectionName
}

// GetCollectionNameOk returns a tuple with the CollectionName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetCollectionNameOk() (*string, bool) {
	if o == nil || IsNil(o.CollectionName) {
		return nil, false
	}

	return o.CollectionName, true
}

// HasCollectionName returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasCollectionName() bool {
	if o != nil && !IsNil(o.CollectionName) {
		return true
	}

	return false
}

// SetCollectionName gets a reference to the given string and assigns it to the CollectionName field.
func (o *SearchIndexResponse) SetCollectionName(v string) {
	o.CollectionName = &v
}

// GetDatabase returns the Database field value if set, zero value otherwise
func (o *SearchIndexResponse) GetDatabase() string {
	if o == nil || IsNil(o.Database) {
		var ret string
		return ret
	}
	return *o.Database
}

// GetDatabaseOk returns a tuple with the Database field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetDatabaseOk() (*string, bool) {
	if o == nil || IsNil(o.Database) {
		return nil, false
	}

	return o.Database, true
}

// HasDatabase returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasDatabase() bool {
	if o != nil && !IsNil(o.Database) {
		return true
	}

	return false
}

// SetDatabase gets a reference to the given string and assigns it to the Database field.
func (o *SearchIndexResponse) SetDatabase(v string) {
	o.Database = &v
}

// GetIndexID returns the IndexID field value if set, zero value otherwise
func (o *SearchIndexResponse) GetIndexID() string {
	if o == nil || IsNil(o.IndexID) {
		var ret string
		return ret
	}
	return *o.IndexID
}

// GetIndexIDOk returns a tuple with the IndexID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetIndexIDOk() (*string, bool) {
	if o == nil || IsNil(o.IndexID) {
		return nil, false
	}

	return o.IndexID, true
}

// HasIndexID returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasIndexID() bool {
	if o != nil && !IsNil(o.IndexID) {
		return true
	}

	return false
}

// SetIndexID gets a reference to the given string and assigns it to the IndexID field.
func (o *SearchIndexResponse) SetIndexID(v string) {
	o.IndexID = &v
}

// GetLatestDefinition returns the LatestDefinition field value if set, zero value otherwise
func (o *SearchIndexResponse) GetLatestDefinition() BaseSearchIndexResponseLatestDefinition {
	if o == nil || IsNil(o.LatestDefinition) {
		var ret BaseSearchIndexResponseLatestDefinition
		return ret
	}
	return *o.LatestDefinition
}

// GetLatestDefinitionOk returns a tuple with the LatestDefinition field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetLatestDefinitionOk() (*BaseSearchIndexResponseLatestDefinition, bool) {
	if o == nil || IsNil(o.LatestDefinition) {
		return nil, false
	}

	return o.LatestDefinition, true
}

// HasLatestDefinition returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasLatestDefinition() bool {
	if o != nil && !IsNil(o.LatestDefinition) {
		return true
	}

	return false
}

// SetLatestDefinition gets a reference to the given BaseSearchIndexResponseLatestDefinition and assigns it to the LatestDefinition field.
func (o *SearchIndexResponse) SetLatestDefinition(v BaseSearchIndexResponseLatestDefinition) {
	o.LatestDefinition = &v
}

// GetLatestDefinitionVersion returns the LatestDefinitionVersion field value if set, zero value otherwise
func (o *SearchIndexResponse) GetLatestDefinitionVersion() SearchIndexDefinitionVersion {
	if o == nil || IsNil(o.LatestDefinitionVersion) {
		var ret SearchIndexDefinitionVersion
		return ret
	}
	return *o.LatestDefinitionVersion
}

// GetLatestDefinitionVersionOk returns a tuple with the LatestDefinitionVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetLatestDefinitionVersionOk() (*SearchIndexDefinitionVersion, bool) {
	if o == nil || IsNil(o.LatestDefinitionVersion) {
		return nil, false
	}

	return o.LatestDefinitionVersion, true
}

// HasLatestDefinitionVersion returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasLatestDefinitionVersion() bool {
	if o != nil && !IsNil(o.LatestDefinitionVersion) {
		return true
	}

	return false
}

// SetLatestDefinitionVersion gets a reference to the given SearchIndexDefinitionVersion and assigns it to the LatestDefinitionVersion field.
func (o *SearchIndexResponse) SetLatestDefinitionVersion(v SearchIndexDefinitionVersion) {
	o.LatestDefinitionVersion = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *SearchIndexResponse) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *SearchIndexResponse) SetName(v string) {
	o.Name = &v
}

// GetQueryable returns the Queryable field value if set, zero value otherwise
func (o *SearchIndexResponse) GetQueryable() bool {
	if o == nil || IsNil(o.Queryable) {
		var ret bool
		return ret
	}
	return *o.Queryable
}

// GetQueryableOk returns a tuple with the Queryable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetQueryableOk() (*bool, bool) {
	if o == nil || IsNil(o.Queryable) {
		return nil, false
	}

	return o.Queryable, true
}

// HasQueryable returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasQueryable() bool {
	if o != nil && !IsNil(o.Queryable) {
		return true
	}

	return false
}

// SetQueryable gets a reference to the given bool and assigns it to the Queryable field.
func (o *SearchIndexResponse) SetQueryable(v bool) {
	o.Queryable = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *SearchIndexResponse) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *SearchIndexResponse) SetStatus(v string) {
	o.Status = &v
}

// GetStatusDetail returns the StatusDetail field value if set, zero value otherwise
func (o *SearchIndexResponse) GetStatusDetail() []VectorSearchHostStatusDetail {
	if o == nil || IsNil(o.StatusDetail) {
		var ret []VectorSearchHostStatusDetail
		return ret
	}
	return *o.StatusDetail
}

// GetStatusDetailOk returns a tuple with the StatusDetail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetStatusDetailOk() (*[]VectorSearchHostStatusDetail, bool) {
	if o == nil || IsNil(o.StatusDetail) {
		return nil, false
	}

	return o.StatusDetail, true
}

// HasStatusDetail returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasStatusDetail() bool {
	if o != nil && !IsNil(o.StatusDetail) {
		return true
	}

	return false
}

// SetStatusDetail gets a reference to the given []VectorSearchHostStatusDetail and assigns it to the StatusDetail field.
func (o *SearchIndexResponse) SetStatusDetail(v []VectorSearchHostStatusDetail) {
	o.StatusDetail = &v
}

// GetType returns the Type field value if set, zero value otherwise
func (o *SearchIndexResponse) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}

	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *SearchIndexResponse) SetType(v string) {
	o.Type = &v
}

// GetSynonymMappingStatus returns the SynonymMappingStatus field value if set, zero value otherwise
func (o *SearchIndexResponse) GetSynonymMappingStatus() string {
	if o == nil || IsNil(o.SynonymMappingStatus) {
		var ret string
		return ret
	}
	return *o.SynonymMappingStatus
}

// GetSynonymMappingStatusOk returns a tuple with the SynonymMappingStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetSynonymMappingStatusOk() (*string, bool) {
	if o == nil || IsNil(o.SynonymMappingStatus) {
		return nil, false
	}

	return o.SynonymMappingStatus, true
}

// HasSynonymMappingStatus returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasSynonymMappingStatus() bool {
	if o != nil && !IsNil(o.SynonymMappingStatus) {
		return true
	}

	return false
}

// SetSynonymMappingStatus gets a reference to the given string and assigns it to the SynonymMappingStatus field.
func (o *SearchIndexResponse) SetSynonymMappingStatus(v string) {
	o.SynonymMappingStatus = &v
}

// GetSynonymMappingStatusDetail returns the SynonymMappingStatusDetail field value if set, zero value otherwise
func (o *SearchIndexResponse) GetSynonymMappingStatusDetail() []map[string]SynonymMappingStatusDetail {
	if o == nil || IsNil(o.SynonymMappingStatusDetail) {
		var ret []map[string]SynonymMappingStatusDetail
		return ret
	}
	return *o.SynonymMappingStatusDetail
}

// GetSynonymMappingStatusDetailOk returns a tuple with the SynonymMappingStatusDetail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchIndexResponse) GetSynonymMappingStatusDetailOk() (*[]map[string]SynonymMappingStatusDetail, bool) {
	if o == nil || IsNil(o.SynonymMappingStatusDetail) {
		return nil, false
	}

	return o.SynonymMappingStatusDetail, true
}

// HasSynonymMappingStatusDetail returns a boolean if a field has been set.
func (o *SearchIndexResponse) HasSynonymMappingStatusDetail() bool {
	if o != nil && !IsNil(o.SynonymMappingStatusDetail) {
		return true
	}

	return false
}

// SetSynonymMappingStatusDetail gets a reference to the given []map[string]SynonymMappingStatusDetail and assigns it to the SynonymMappingStatusDetail field.
func (o *SearchIndexResponse) SetSynonymMappingStatusDetail(v []map[string]SynonymMappingStatusDetail) {
	o.SynonymMappingStatusDetail = &v
}
