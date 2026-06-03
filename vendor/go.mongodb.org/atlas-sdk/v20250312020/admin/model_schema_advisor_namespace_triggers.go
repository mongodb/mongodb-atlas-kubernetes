// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SchemaAdvisorNamespaceTriggers struct for SchemaAdvisorNamespaceTriggers
type SchemaAdvisorNamespaceTriggers struct {
	// Namespace of the affected collection. Will be null for `REDUCE_NUMBER_OF_NAMESPACE` recommendation.
	// Read only field.
	Namespace *string `json:"namespace,omitempty"`
	// List of triggers that specify why the collection activated the recommendation.
	// Read only field.
	Triggers *[]SchemaAdvisorTriggerDetails `json:"triggers,omitempty"`
}

// NewSchemaAdvisorNamespaceTriggers instantiates a new SchemaAdvisorNamespaceTriggers object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSchemaAdvisorNamespaceTriggers() *SchemaAdvisorNamespaceTriggers {
	this := SchemaAdvisorNamespaceTriggers{}
	return &this
}

// NewSchemaAdvisorNamespaceTriggersWithDefaults instantiates a new SchemaAdvisorNamespaceTriggers object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSchemaAdvisorNamespaceTriggersWithDefaults() *SchemaAdvisorNamespaceTriggers {
	this := SchemaAdvisorNamespaceTriggers{}
	return &this
}

// GetNamespace returns the Namespace field value if set, zero value otherwise
func (o *SchemaAdvisorNamespaceTriggers) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorNamespaceTriggers) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}

	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *SchemaAdvisorNamespaceTriggers) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *SchemaAdvisorNamespaceTriggers) SetNamespace(v string) {
	o.Namespace = &v
}

// GetTriggers returns the Triggers field value if set, zero value otherwise
func (o *SchemaAdvisorNamespaceTriggers) GetTriggers() []SchemaAdvisorTriggerDetails {
	if o == nil || IsNil(o.Triggers) {
		var ret []SchemaAdvisorTriggerDetails
		return ret
	}
	return *o.Triggers
}

// GetTriggersOk returns a tuple with the Triggers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorNamespaceTriggers) GetTriggersOk() (*[]SchemaAdvisorTriggerDetails, bool) {
	if o == nil || IsNil(o.Triggers) {
		return nil, false
	}

	return o.Triggers, true
}

// HasTriggers returns a boolean if a field has been set.
func (o *SchemaAdvisorNamespaceTriggers) HasTriggers() bool {
	if o != nil && !IsNil(o.Triggers) {
		return true
	}

	return false
}

// SetTriggers gets a reference to the given []SchemaAdvisorTriggerDetails and assigns it to the Triggers field.
func (o *SchemaAdvisorNamespaceTriggers) SetTriggers(v []SchemaAdvisorTriggerDetails) {
	o.Triggers = &v
}
