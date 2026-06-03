// Code based on the AtlasAPI V2 OpenAPI file

package admin

// StreamsMatcher Rules to apply when comparing a stream processing workspace or stream processor against this alert configuration.
type StreamsMatcher struct {
	// Name of the parameter in the target object that MongoDB Cloud checks. The parameter must match all rules for MongoDB Cloud to check for alert configurations.
	FieldName string `json:"fieldName"`
	// Comparison operator to apply when checking the current metric value against **matcher[n].value**. The `REGEX` operator only supports inclusive matches. Use the `NOT_CONTAINS` operator to exclude values.
	Operator string `json:"operator"`
	// Value to match or exceed using the specified `matchers.operator`.
	Value string `json:"value"`
}

// NewStreamsMatcher instantiates a new StreamsMatcher object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreamsMatcher(fieldName string, operator string, value string) *StreamsMatcher {
	this := StreamsMatcher{}
	this.FieldName = fieldName
	this.Operator = operator
	this.Value = value
	return &this
}

// NewStreamsMatcherWithDefaults instantiates a new StreamsMatcher object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreamsMatcherWithDefaults() *StreamsMatcher {
	this := StreamsMatcher{}
	return &this
}

// GetFieldName returns the FieldName field value
func (o *StreamsMatcher) GetFieldName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FieldName
}

// GetFieldNameOk returns a tuple with the FieldName field value
// and a boolean to check if the value has been set.
func (o *StreamsMatcher) GetFieldNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FieldName, true
}

// SetFieldName sets field value
func (o *StreamsMatcher) SetFieldName(v string) {
	o.FieldName = v
}

// GetOperator returns the Operator field value
func (o *StreamsMatcher) GetOperator() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Operator
}

// GetOperatorOk returns a tuple with the Operator field value
// and a boolean to check if the value has been set.
func (o *StreamsMatcher) GetOperatorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Operator, true
}

// SetOperator sets field value
func (o *StreamsMatcher) SetOperator(v string) {
	o.Operator = v
}

// GetValue returns the Value field value
func (o *StreamsMatcher) GetValue() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *StreamsMatcher) GetValueOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *StreamsMatcher) SetValue(v string) {
	o.Value = v
}
