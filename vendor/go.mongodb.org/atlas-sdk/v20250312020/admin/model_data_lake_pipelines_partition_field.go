// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DataLakePipelinesPartitionField Partition Field in the Data Lake Storage provider for a Data Lake Pipeline.
type DataLakePipelinesPartitionField struct {
	// Human-readable label that identifies the field name used to partition data.
	FieldName string `json:"fieldName"`
	// Sequence in which MongoDB Cloud slices the collection data to create partitions. The resource expresses this sequence starting with zero.
	Order int `json:"order"`
}

// NewDataLakePipelinesPartitionField instantiates a new DataLakePipelinesPartitionField object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataLakePipelinesPartitionField(fieldName string, order int) *DataLakePipelinesPartitionField {
	this := DataLakePipelinesPartitionField{}
	this.FieldName = fieldName
	this.Order = order
	return &this
}

// NewDataLakePipelinesPartitionFieldWithDefaults instantiates a new DataLakePipelinesPartitionField object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataLakePipelinesPartitionFieldWithDefaults() *DataLakePipelinesPartitionField {
	this := DataLakePipelinesPartitionField{}
	var order int = 0
	this.Order = order
	return &this
}

// GetFieldName returns the FieldName field value
func (o *DataLakePipelinesPartitionField) GetFieldName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FieldName
}

// GetFieldNameOk returns a tuple with the FieldName field value
// and a boolean to check if the value has been set.
func (o *DataLakePipelinesPartitionField) GetFieldNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FieldName, true
}

// SetFieldName sets field value
func (o *DataLakePipelinesPartitionField) SetFieldName(v string) {
	o.FieldName = v
}

// GetOrder returns the Order field value
func (o *DataLakePipelinesPartitionField) GetOrder() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.Order
}

// GetOrderOk returns a tuple with the Order field value
// and a boolean to check if the value has been set.
func (o *DataLakePipelinesPartitionField) GetOrderOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Order, true
}

// SetOrder sets field value
func (o *DataLakePipelinesPartitionField) SetOrder(v int) {
	o.Order = v
}
