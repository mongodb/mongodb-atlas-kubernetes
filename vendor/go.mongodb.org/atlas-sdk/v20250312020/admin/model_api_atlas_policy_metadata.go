// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasPolicyMetadata struct for ApiAtlasPolicyMetadata
type ApiAtlasPolicyMetadata struct {
	// Unique 24-hexadecimal character string that identifies the policy.
	// Read only field.
	PolicyId *string `json:"policyId,omitempty"`
}

// NewApiAtlasPolicyMetadata instantiates a new ApiAtlasPolicyMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasPolicyMetadata() *ApiAtlasPolicyMetadata {
	this := ApiAtlasPolicyMetadata{}
	return &this
}

// NewApiAtlasPolicyMetadataWithDefaults instantiates a new ApiAtlasPolicyMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasPolicyMetadataWithDefaults() *ApiAtlasPolicyMetadata {
	this := ApiAtlasPolicyMetadata{}
	return &this
}

// GetPolicyId returns the PolicyId field value if set, zero value otherwise
func (o *ApiAtlasPolicyMetadata) GetPolicyId() string {
	if o == nil || IsNil(o.PolicyId) {
		var ret string
		return ret
	}
	return *o.PolicyId
}

// GetPolicyIdOk returns a tuple with the PolicyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasPolicyMetadata) GetPolicyIdOk() (*string, bool) {
	if o == nil || IsNil(o.PolicyId) {
		return nil, false
	}

	return o.PolicyId, true
}

// HasPolicyId returns a boolean if a field has been set.
func (o *ApiAtlasPolicyMetadata) HasPolicyId() bool {
	if o != nil && !IsNil(o.PolicyId) {
		return true
	}

	return false
}

// SetPolicyId gets a reference to the given string and assigns it to the PolicyId field.
func (o *ApiAtlasPolicyMetadata) SetPolicyId(v string) {
	o.PolicyId = &v
}
