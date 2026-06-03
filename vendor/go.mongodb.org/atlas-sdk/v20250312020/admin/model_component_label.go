// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ComponentLabel Human-readable labels applied to this MongoDB Cloud component.
type ComponentLabel struct {
	// Key applied to tag and categorize this component.
	Key *string `json:"key,omitempty"`
	// Value set to the Key applied to tag and categorize this component.
	Value *string `json:"value,omitempty"`
}

// NewComponentLabel instantiates a new ComponentLabel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewComponentLabel() *ComponentLabel {
	this := ComponentLabel{}
	return &this
}

// NewComponentLabelWithDefaults instantiates a new ComponentLabel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewComponentLabelWithDefaults() *ComponentLabel {
	this := ComponentLabel{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise
func (o *ComponentLabel) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComponentLabel) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}

	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *ComponentLabel) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *ComponentLabel) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise
func (o *ComponentLabel) GetValue() string {
	if o == nil || IsNil(o.Value) {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ComponentLabel) GetValueOk() (*string, bool) {
	if o == nil || IsNil(o.Value) {
		return nil, false
	}

	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *ComponentLabel) HasValue() bool {
	if o != nil && !IsNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *ComponentLabel) SetValue(v string) {
	o.Value = &v
}
