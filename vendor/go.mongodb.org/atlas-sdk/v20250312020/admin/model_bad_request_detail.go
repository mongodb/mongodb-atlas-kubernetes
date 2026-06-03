// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BadRequestDetail Bad request detail.
type BadRequestDetail struct {
	// Describes all violations in a client request.
	Fields *[]FieldViolation `json:"fields,omitempty"`
}

// NewBadRequestDetail instantiates a new BadRequestDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBadRequestDetail() *BadRequestDetail {
	this := BadRequestDetail{}
	return &this
}

// NewBadRequestDetailWithDefaults instantiates a new BadRequestDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBadRequestDetailWithDefaults() *BadRequestDetail {
	this := BadRequestDetail{}
	return &this
}

// GetFields returns the Fields field value if set, zero value otherwise
func (o *BadRequestDetail) GetFields() []FieldViolation {
	if o == nil || IsNil(o.Fields) {
		var ret []FieldViolation
		return ret
	}
	return *o.Fields
}

// GetFieldsOk returns a tuple with the Fields field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BadRequestDetail) GetFieldsOk() (*[]FieldViolation, bool) {
	if o == nil || IsNil(o.Fields) {
		return nil, false
	}

	return o.Fields, true
}

// HasFields returns a boolean if a field has been set.
func (o *BadRequestDetail) HasFields() bool {
	if o != nil && !IsNil(o.Fields) {
		return true
	}

	return false
}

// SetFields gets a reference to the given []FieldViolation and assigns it to the Fields field.
func (o *BadRequestDetail) SetFields(v []FieldViolation) {
	o.Fields = &v
}
