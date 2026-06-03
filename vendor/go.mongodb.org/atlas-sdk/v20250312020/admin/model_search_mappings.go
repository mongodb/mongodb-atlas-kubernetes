// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SearchMappings Index specifications for the collection's fields.
type SearchMappings struct {
	// Indicates whether the index uses static, default dynamic, or configurable dynamic mappings. Set to `true` to enable dynamic mapping with default type set or define object to specify the name of the configured type sets for dynamic mapping. If you specify configurable dynamic mappings, you must define the referred type sets in the `typeSets` field. Set to `false` to use only static mappings through `mappings.fields`.
	Dynamic any `json:"dynamic,omitempty"`
	// One or more field specifications for the Atlas Search index. Required if `mappings.dynamic` is omitted or set to `false`.
	Fields *map[string]any `json:"fields,omitempty"`
}

// NewSearchMappings instantiates a new SearchMappings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSearchMappings() *SearchMappings {
	this := SearchMappings{}
	return &this
}

// NewSearchMappingsWithDefaults instantiates a new SearchMappings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSearchMappingsWithDefaults() *SearchMappings {
	this := SearchMappings{}
	return &this
}

// GetDynamic returns the Dynamic field value if set, zero value otherwise
func (o *SearchMappings) GetDynamic() any {
	if o == nil || IsNil(o.Dynamic) {
		var ret any
		return ret
	}
	return o.Dynamic
}

// GetDynamicOk returns a tuple with the Dynamic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchMappings) GetDynamicOk() (any, bool) {
	if o == nil || IsNil(o.Dynamic) {
		var ret any
		return ret, false
	}

	return o.Dynamic, true
}

// HasDynamic returns a boolean if a field has been set.
func (o *SearchMappings) HasDynamic() bool {
	if o != nil && !IsNil(o.Dynamic) {
		return true
	}

	return false
}

// SetDynamic gets a reference to the given any and assigns it to the Dynamic field.
func (o *SearchMappings) SetDynamic(v any) {
	o.Dynamic = v
}

// GetFields returns the Fields field value if set, zero value otherwise
func (o *SearchMappings) GetFields() map[string]any {
	if o == nil || IsNil(o.Fields) {
		var ret map[string]any
		return ret
	}
	return *o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchMappings) GetFieldsOk() (*map[string]any, bool) {
	if o == nil || IsNil(o.Fields) {
		return nil, false
	}

	return o.Fields, true
}

// HasFields returns a boolean if a field has been set.
func (o *SearchMappings) HasFields() bool {
	if o != nil && !IsNil(o.Fields) {
		return true
	}

	return false
}

// SetFields gets a reference to the given map[string]any and assigns it to the Fields field.
func (o *SearchMappings) SetFields(v map[string]any) {
	o.Fields = &v
}
