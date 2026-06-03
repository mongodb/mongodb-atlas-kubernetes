// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SchemaAdvisorResponse struct for SchemaAdvisorResponse
type SchemaAdvisorResponse struct {
	// List that contains the documents with information about the schema advice that Performance Advisor suggests.
	// Read only field.
	Recommendations *[]SchemaAdvisorItemRecommendation `json:"recommendations,omitempty"`
}

// NewSchemaAdvisorResponse instantiates a new SchemaAdvisorResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSchemaAdvisorResponse() *SchemaAdvisorResponse {
	this := SchemaAdvisorResponse{}
	return &this
}

// NewSchemaAdvisorResponseWithDefaults instantiates a new SchemaAdvisorResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSchemaAdvisorResponseWithDefaults() *SchemaAdvisorResponse {
	this := SchemaAdvisorResponse{}
	return &this
}

// GetRecommendations returns the Recommendations field value if set, zero value otherwise
func (o *SchemaAdvisorResponse) GetRecommendations() []SchemaAdvisorItemRecommendation {
	if o == nil || IsNil(o.Recommendations) {
		var ret []SchemaAdvisorItemRecommendation
		return ret
	}
	return *o.Recommendations
}

// GetRecommendationsOk returns a tuple with the Recommendations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaAdvisorResponse) GetRecommendationsOk() (*[]SchemaAdvisorItemRecommendation, bool) {
	if o == nil || IsNil(o.Recommendations) {
		return nil, false
	}

	return o.Recommendations, true
}

// HasRecommendations returns a boolean if a field has been set.
func (o *SchemaAdvisorResponse) HasRecommendations() bool {
	if o != nil && !IsNil(o.Recommendations) {
		return true
	}

	return false
}

// SetRecommendations gets a reference to the given []SchemaAdvisorItemRecommendation and assigns it to the Recommendations field.
func (o *SchemaAdvisorResponse) SetRecommendations(v []SchemaAdvisorItemRecommendation) {
	o.Recommendations = &v
}
