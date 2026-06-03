// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SchemaAdvisorItemRecommendation struct for SchemaAdvisorItemRecommendation
type SchemaAdvisorItemRecommendation struct {
	// List that contains the namespaces and information on why those namespaces triggered the recommendation.
	// Read only field.
	AffectedNamespaces *[]SchemaAdvisorNamespaceTriggers `json:"affectedNamespaces,omitempty"`
	// Description of the specified recommendation.
	// Read only field.
	Description *string `json:"description,omitempty"`
	// Type of recommendation.
	// Read only field.
	Recommendation *string `json:"recommendation,omitempty"`
}

// NewSchemaAdvisorItemRecommendation instantiates a new SchemaAdvisorItemRecommendation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSchemaAdvisorItemRecommendation() *SchemaAdvisorItemRecommendation {
	this := SchemaAdvisorItemRecommendation{}
	return &this
}

// NewSchemaAdvisorItemRecommendationWithDefaults instantiates a new SchemaAdvisorItemRecommendation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSchemaAdvisorItemRecommendationWithDefaults() *SchemaAdvisorItemRecommendation {
	this := SchemaAdvisorItemRecommendation{}
	return &this
}

// GetAffectedNamespaces returns the AffectedNamespaces field value if set, zero value otherwise
func (o *SchemaAdvisorItemRecommendation) GetAffectedNamespaces() []SchemaAdvisorNamespaceTriggers {
	if o == nil || IsNil(o.AffectedNamespaces) {
		var ret []SchemaAdvisorNamespaceTriggers
		return ret
	}
	return *o.AffectedNamespaces
}

// GetAffectedNamespacesOk returns a tuple with the AffectedNamespaces field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorItemRecommendation) GetAffectedNamespacesOk() (*[]SchemaAdvisorNamespaceTriggers, bool) {
	if o == nil || IsNil(o.AffectedNamespaces) {
		return nil, false
	}

	return o.AffectedNamespaces, true
}

// HasAffectedNamespaces returns a boolean if a field has been set.
func (o *SchemaAdvisorItemRecommendation) HasAffectedNamespaces() bool {
	if o != nil && !IsNil(o.AffectedNamespaces) {
		return true
	}

	return false
}

// SetAffectedNamespaces gets a reference to the given []SchemaAdvisorNamespaceTriggers and assigns it to the AffectedNamespaces field.
func (o *SchemaAdvisorItemRecommendation) SetAffectedNamespaces(v []SchemaAdvisorNamespaceTriggers) {
	o.AffectedNamespaces = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *SchemaAdvisorItemRecommendation) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorItemRecommendation) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *SchemaAdvisorItemRecommendation) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *SchemaAdvisorItemRecommendation) SetDescription(v string) {
	o.Description = &v
}

// GetRecommendation returns the Recommendation field value if set, zero value otherwise
func (o *SchemaAdvisorItemRecommendation) GetRecommendation() string {
	if o == nil || IsNil(o.Recommendation) {
		var ret string
		return ret
	}
	return *o.Recommendation
}

// GetRecommendationOk returns a tuple with the Recommendation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorItemRecommendation) GetRecommendationOk() (*string, bool) {
	if o == nil || IsNil(o.Recommendation) {
		return nil, false
	}

	return o.Recommendation, true
}

// HasRecommendation returns a boolean if a field has been set.
func (o *SchemaAdvisorItemRecommendation) HasRecommendation() bool {
	if o != nil && !IsNil(o.Recommendation) {
		return true
	}

	return false
}

// SetRecommendation gets a reference to the given string and assigns it to the Recommendation field.
func (o *SchemaAdvisorItemRecommendation) SetRecommendation(v string) {
	o.Recommendation = &v
}
