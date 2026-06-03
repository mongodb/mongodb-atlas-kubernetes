// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SchemaAdvisorTriggerDetails struct for SchemaAdvisorTriggerDetails
type SchemaAdvisorTriggerDetails struct {
	// Description of the trigger type.
	// Read only field.
	Description *string `json:"description,omitempty"`
	// Type of trigger.
	// Read only field.
	TriggerType *string `json:"triggerType,omitempty"`
}

// NewSchemaAdvisorTriggerDetails instantiates a new SchemaAdvisorTriggerDetails object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSchemaAdvisorTriggerDetails() *SchemaAdvisorTriggerDetails {
	this := SchemaAdvisorTriggerDetails{}
	return &this
}

// NewSchemaAdvisorTriggerDetailsWithDefaults instantiates a new SchemaAdvisorTriggerDetails object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSchemaAdvisorTriggerDetailsWithDefaults() *SchemaAdvisorTriggerDetails {
	this := SchemaAdvisorTriggerDetails{}
	return &this
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *SchemaAdvisorTriggerDetails) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorTriggerDetails) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *SchemaAdvisorTriggerDetails) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *SchemaAdvisorTriggerDetails) SetDescription(v string) {
	o.Description = &v
}

// GetTriggerType returns the TriggerType field value if set, zero value otherwise
func (o *SchemaAdvisorTriggerDetails) GetTriggerType() string {
	if o == nil || IsNil(o.TriggerType) {
		var ret string
		return ret
	}
	return *o.TriggerType
}

// GetTriggerTypeOk returns a tuple with the TriggerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorTriggerDetails) GetTriggerTypeOk() (*string, bool) {
	if o == nil || IsNil(o.TriggerType) {
		return nil, false
	}

	return o.TriggerType, true
}

// HasTriggerType returns a boolean if a field has been set.
func (o *SchemaAdvisorTriggerDetails) HasTriggerType() bool {
	if o != nil && !IsNil(o.TriggerType) {
		return true
	}

	return false
}

// SetTriggerType gets a reference to the given string and assigns it to the TriggerType field.
func (o *SchemaAdvisorTriggerDetails) SetTriggerType(v string) {
	o.TriggerType = &v
}
