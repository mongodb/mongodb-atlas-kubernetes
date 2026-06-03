// Code based on the AtlasAPI V2 OpenAPI file

package admin

// VectorSearchIndexDefinition The vector search index definition set by the user.
type VectorSearchIndexDefinition struct {
	// Settings that configure the fields, one per object, to index. You must define at least one \"vector\" type field. You can optionally define \"filter\" type fields also.
	Fields []any `json:"fields"`
	// Top-level path to the array that contains vector fields. When provided, vector fields under this path are treated as nested.
	NestedRoot *string `json:"nestedRoot,omitempty"`
	// Number of index partitions. Allowed values are [1, 2, 4].
	NumPartitions *int `json:"numPartitions,omitempty"`
	// Flag that indicates whether to store all fields (true) on Atlas Search. By default, Atlas doesn't store (false) the fields on Atlas Search.  Alternatively, you can specify an object that only contains the list of fields to store (include) or not store (exclude) on Atlas Search. Note that storing all fields (true) is not allowed for vector search indexes. To learn more, see Stored Source Fields.
	StoredSource any `json:"storedSource,omitempty"`
}

// NewVectorSearchIndexDefinition instantiates a new VectorSearchIndexDefinition object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVectorSearchIndexDefinition(fields []any) *VectorSearchIndexDefinition {
	this := VectorSearchIndexDefinition{}
	this.Fields = fields
	var numPartitions int = 1
	this.NumPartitions = &numPartitions
	return &this
}

// NewVectorSearchIndexDefinitionWithDefaults instantiates a new VectorSearchIndexDefinition object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVectorSearchIndexDefinitionWithDefaults() *VectorSearchIndexDefinition {
	this := VectorSearchIndexDefinition{}
	var numPartitions int = 1
	this.NumPartitions = &numPartitions
	return &this
}

// GetFields returns the Fields field value
func (o *VectorSearchIndexDefinition) GetFields() []any {
	if o == nil {
		var ret []any
		return ret
	}

	return o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexDefinition) GetFieldsOk() (*[]any, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Fields, true
}

// SetFields sets field value
func (o *VectorSearchIndexDefinition) SetFields(v []any) {
	o.Fields = v
}

// GetNestedRoot returns the NestedRoot field value if set, zero value otherwise
func (o *VectorSearchIndexDefinition) GetNestedRoot() string {
	if o == nil || IsNil(o.NestedRoot) {
		var ret string
		return ret
	}
	return *o.NestedRoot
}

// GetNestedRootOk returns a tuple with the NestedRoot field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexDefinition) GetNestedRootOk() (*string, bool) {
	if o == nil || IsNil(o.NestedRoot) {
		return nil, false
	}

	return o.NestedRoot, true
}

// HasNestedRoot returns a boolean if a field has been set.
func (o *VectorSearchIndexDefinition) HasNestedRoot() bool {
	if o != nil && !IsNil(o.NestedRoot) {
		return true
	}

	return false
}

// SetNestedRoot gets a reference to the given string and assigns it to the NestedRoot field.
func (o *VectorSearchIndexDefinition) SetNestedRoot(v string) {
	o.NestedRoot = &v
}

// GetNumPartitions returns the NumPartitions field value if set, zero value otherwise
func (o *VectorSearchIndexDefinition) GetNumPartitions() int {
	if o == nil || IsNil(o.NumPartitions) {
		var ret int
		return ret
	}
	return *o.NumPartitions
}

// GetNumPartitionsOk returns a tuple with the NumPartitions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexDefinition) GetNumPartitionsOk() (*int, bool) {
	if o == nil || IsNil(o.NumPartitions) {
		return nil, false
	}

	return o.NumPartitions, true
}

// HasNumPartitions returns a boolean if a field has been set.
func (o *VectorSearchIndexDefinition) HasNumPartitions() bool {
	if o != nil && !IsNil(o.NumPartitions) {
		return true
	}

	return false
}

// SetNumPartitions gets a reference to the given int and assigns it to the NumPartitions field.
func (o *VectorSearchIndexDefinition) SetNumPartitions(v int) {
	o.NumPartitions = &v
}

// GetStoredSource returns the StoredSource field value if set, zero value otherwise
func (o *VectorSearchIndexDefinition) GetStoredSource() any {
	if o == nil || IsNil(o.StoredSource) {
		var ret any
		return ret
	}
	return o.StoredSource
}

// GetStoredSourceOk returns a tuple with the StoredSource field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *VectorSearchIndexDefinition) GetStoredSourceOk() (any, bool) {
	if o == nil || IsNil(o.StoredSource) {
		var ret any
		return ret, false
	}

	return o.StoredSource, true
}

// HasStoredSource returns a boolean if a field has been set.
func (o *VectorSearchIndexDefinition) HasStoredSource() bool {
	if o != nil && !IsNil(o.StoredSource) {
		return true
	}

	return false
}

// SetStoredSource gets a reference to the given any and assigns it to the StoredSource field.
func (o *VectorSearchIndexDefinition) SetStoredSource(v any) {
	o.StoredSource = v
}
