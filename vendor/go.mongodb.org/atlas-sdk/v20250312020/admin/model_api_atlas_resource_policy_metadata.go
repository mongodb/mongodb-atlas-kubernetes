// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasResourcePolicyMetadata struct for ApiAtlasResourcePolicyMetadata
type ApiAtlasResourcePolicyMetadata struct {
	// List of policies that are in conflict with the current state of the resource.
	// Read only field.
	PoliciesCausingNonCompliance *[]ApiAtlasPolicyMetadata `json:"policiesCausingNonCompliance,omitempty"`
	// Unique 24-hexadecimal character string that identifies the atlas resource policy.
	// Read only field.
	ResourcePolicyId *string `json:"resourcePolicyId,omitempty"`
	// Human-readable label that describes the atlas resource policy.
	// Read only field.
	ResourcePolicyName *string `json:"resourcePolicyName,omitempty"`
}

// NewApiAtlasResourcePolicyMetadata instantiates a new ApiAtlasResourcePolicyMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasResourcePolicyMetadata() *ApiAtlasResourcePolicyMetadata {
	this := ApiAtlasResourcePolicyMetadata{}
	return &this
}

// NewApiAtlasResourcePolicyMetadataWithDefaults instantiates a new ApiAtlasResourcePolicyMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasResourcePolicyMetadataWithDefaults() *ApiAtlasResourcePolicyMetadata {
	this := ApiAtlasResourcePolicyMetadata{}
	return &this
}

// GetPoliciesCausingNonCompliance returns the PoliciesCausingNonCompliance field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicyMetadata) GetPoliciesCausingNonCompliance() []ApiAtlasPolicyMetadata {
	if o == nil || IsNil(o.PoliciesCausingNonCompliance) {
		var ret []ApiAtlasPolicyMetadata
		return ret
	}
	return *o.PoliciesCausingNonCompliance
}

// GetPoliciesCausingNonComplianceOk returns a tuple with the PoliciesCausingNonCompliance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicyMetadata) GetPoliciesCausingNonComplianceOk() (*[]ApiAtlasPolicyMetadata, bool) {
	if o == nil || IsNil(o.PoliciesCausingNonCompliance) {
		return nil, false
	}

	return o.PoliciesCausingNonCompliance, true
}

// HasPoliciesCausingNonCompliance returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicyMetadata) HasPoliciesCausingNonCompliance() bool {
	if o != nil && !IsNil(o.PoliciesCausingNonCompliance) {
		return true
	}

	return false
}

// SetPoliciesCausingNonCompliance gets a reference to the given []ApiAtlasPolicyMetadata and assigns it to the PoliciesCausingNonCompliance field.
func (o *ApiAtlasResourcePolicyMetadata) SetPoliciesCausingNonCompliance(v []ApiAtlasPolicyMetadata) {
	o.PoliciesCausingNonCompliance = &v
}

// GetResourcePolicyId returns the ResourcePolicyId field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicyMetadata) GetResourcePolicyId() string {
	if o == nil || IsNil(o.ResourcePolicyId) {
		var ret string
		return ret
	}
	return *o.ResourcePolicyId
}

// GetResourcePolicyIdOk returns a tuple with the ResourcePolicyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicyMetadata) GetResourcePolicyIdOk() (*string, bool) {
	if o == nil || IsNil(o.ResourcePolicyId) {
		return nil, false
	}

	return o.ResourcePolicyId, true
}

// HasResourcePolicyId returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicyMetadata) HasResourcePolicyId() bool {
	if o != nil && !IsNil(o.ResourcePolicyId) {
		return true
	}

	return false
}

// SetResourcePolicyId gets a reference to the given string and assigns it to the ResourcePolicyId field.
func (o *ApiAtlasResourcePolicyMetadata) SetResourcePolicyId(v string) {
	o.ResourcePolicyId = &v
}

// GetResourcePolicyName returns the ResourcePolicyName field value if set, zero value otherwise
func (o *ApiAtlasResourcePolicyMetadata) GetResourcePolicyName() string {
	if o == nil || IsNil(o.ResourcePolicyName) {
		var ret string
		return ret
	}
	return *o.ResourcePolicyName
}

// GetResourcePolicyNameOk returns a tuple with the ResourcePolicyName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasResourcePolicyMetadata) GetResourcePolicyNameOk() (*string, bool) {
	if o == nil || IsNil(o.ResourcePolicyName) {
		return nil, false
	}

	return o.ResourcePolicyName, true
}

// HasResourcePolicyName returns a boolean if a field has been set.
func (o *ApiAtlasResourcePolicyMetadata) HasResourcePolicyName() bool {
	if o != nil && !IsNil(o.ResourcePolicyName) {
		return true
	}

	return false
}

// SetResourcePolicyName gets a reference to the given string and assigns it to the ResourcePolicyName field.
func (o *ApiAtlasResourcePolicyMetadata) SetResourcePolicyName(v string) {
	o.ResourcePolicyName = &v
}
