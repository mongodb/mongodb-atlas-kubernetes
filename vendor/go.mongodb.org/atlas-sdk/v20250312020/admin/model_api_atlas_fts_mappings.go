// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasFTSMappings Index specifications for the collection's fields.
type ApiAtlasFTSMappings struct {
	// Flag that indicates whether the index uses dynamic or static mappings. Required if `mappings.fields` is omitted.
	Dynamic *bool `json:"dynamic,omitempty"`
	// One or more field specifications for the Atlas Search index. Required if `mappings.dynamic` is omitted or set to `false`.
	Fields any `json:"fields,omitempty"`
}

// NewApiAtlasFTSMappings instantiates a new ApiAtlasFTSMappings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasFTSMappings() *ApiAtlasFTSMappings {
	this := ApiAtlasFTSMappings{}
	var dynamic bool = false
	this.Dynamic = &dynamic
	return &this
}

// NewApiAtlasFTSMappingsWithDefaults instantiates a new ApiAtlasFTSMappings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasFTSMappingsWithDefaults() *ApiAtlasFTSMappings {
	this := ApiAtlasFTSMappings{}
	var dynamic bool = false
	this.Dynamic = &dynamic
	return &this
}

// GetDynamic returns the Dynamic field value if set, zero value otherwise
func (o *ApiAtlasFTSMappings) GetDynamic() bool {
	if o == nil || IsNil(o.Dynamic) {
		var ret bool
		return ret
	}
	return *o.Dynamic
}

// GetDynamicOk returns a tuple with the Dynamic field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSMappings) GetDynamicOk() (*bool, bool) {
	if o == nil || IsNil(o.Dynamic) {
		return nil, false
	}

	return o.Dynamic, true
}

// HasDynamic returns a boolean if a field has been set.
func (o *ApiAtlasFTSMappings) HasDynamic() bool {
	if o != nil && !IsNil(o.Dynamic) {
		return true
	}

	return false
}

// SetDynamic gets a reference to the given bool and assigns it to the Dynamic field.
func (o *ApiAtlasFTSMappings) SetDynamic(v bool) {
	o.Dynamic = &v
}

// GetFields returns the Fields field value if set, zero value otherwise
func (o *ApiAtlasFTSMappings) GetFields() any {
	if o == nil || IsNil(o.Fields) {
		var ret any
		return ret
	}
	return o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasFTSMappings) GetFieldsOk() (any, bool) {
	if o == nil || IsNil(o.Fields) {
		var ret any
		return ret, false
	}

	return o.Fields, true
}

// HasFields returns a boolean if a field has been set.
func (o *ApiAtlasFTSMappings) HasFields() bool {
	if o != nil && !IsNil(o.Fields) {
		return true
	}

	return false
}

// SetFields gets a reference to the given any and assigns it to the Fields field.
func (o *ApiAtlasFTSMappings) SetFields(v any) {
	o.Fields = v
}
