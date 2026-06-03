// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FieldViolation struct for FieldViolation
type FieldViolation struct {
	// A description of why the request element is bad.
	Description string `json:"description"`
	// A path that leads to a field in the request body.
	Field string `json:"field"`
}

// NewFieldViolation instantiates a new FieldViolation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFieldViolation(description string, field string) *FieldViolation {
	this := FieldViolation{}
	this.Description = description
	this.Field = field
	return &this
}

// NewFieldViolationWithDefaults instantiates a new FieldViolation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFieldViolationWithDefaults() *FieldViolation {
	this := FieldViolation{}
	return &this
}

// GetDescription returns the Description field value
func (o *FieldViolation) GetDescription() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Description
}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
func (o *FieldViolation) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Description, true
}

// SetDescription sets field value
func (o *FieldViolation) SetDescription(v string) {
	o.Description = v
}

// GetField returns the Field field value
func (o *FieldViolation) GetField() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Field
}

// GetFieldOk returns a tuple with the Field field value
// and a boolean to check if the value has been set.
func (o *FieldViolation) GetFieldOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Field, true
}

// SetField sets field value
func (o *FieldViolation) SetField(v string) {
	o.Field = v
}
