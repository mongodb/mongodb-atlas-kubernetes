// Code based on the AtlasAPI V2 OpenAPI file

package admin

// PartitionField Metadata to partition this online archive.
type PartitionField struct {
	// Human-readable label that identifies the parameter that MongoDB Cloud uses to partition data. To specify a nested parameter, use the dot notation.
	FieldName string `json:"fieldName"`
	// Data type of the parameter that that MongoDB Cloud uses to partition data. Partition parameters of type UUID must be of binary subtype 4. MongoDB Cloud skips partition parameters of type UUID with subtype 3.
	// Read only field.
	FieldType *string `json:"fieldType,omitempty"`
	// Sequence in which MongoDB Cloud slices the collection data to create partitions. The resource expresses this sequence starting with zero. The value of the `criteria.dateField` parameter defaults as the first item in the partition sequence.
	Order int `json:"order"`
}

// NewPartitionField instantiates a new PartitionField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPartitionField(fieldName string, order int) *PartitionField {
	this := PartitionField{}
	this.FieldName = fieldName
	this.Order = order
	return &this
}

// NewPartitionFieldWithDefaults instantiates a new PartitionField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPartitionFieldWithDefaults() *PartitionField {
	this := PartitionField{}
	var order int = 0
	this.Order = order
	return &this
}

// GetFieldName returns the FieldName field value
func (o *PartitionField) GetFieldName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FieldName
}

// GetFieldNameOk returns a tuple with the FieldName field value
// and a boolean to check if the value has been set.
func (o *PartitionField) GetFieldNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FieldName, true
}

// SetFieldName sets field value
func (o *PartitionField) SetFieldName(v string) {
	o.FieldName = v
}

// GetFieldType returns the FieldType field value if set, zero value otherwise
func (o *PartitionField) GetFieldType() string {
	if o == nil || IsNil(o.FieldType) {
		var ret string
		return ret
	}
	return *o.FieldType
}

// GetFieldTypeOk returns a tuple with the FieldType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PartitionField) GetFieldTypeOk() (*string, bool) {
	if o == nil || IsNil(o.FieldType) {
		return nil, false
	}

	return o.FieldType, true
}

// HasFieldType returns a boolean if a field has been set.
func (o *PartitionField) HasFieldType() bool {
	if o != nil && !IsNil(o.FieldType) {
		return true
	}

	return false
}

// SetFieldType gets a reference to the given string and assigns it to the FieldType field.
func (o *PartitionField) SetFieldType(v string) {
	o.FieldType = &v
}

// GetOrder returns the Order field value
func (o *PartitionField) GetOrder() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *PartitionField) GetOrderOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *PartitionField) SetOrder(v int) {
	o.Order = v
}
