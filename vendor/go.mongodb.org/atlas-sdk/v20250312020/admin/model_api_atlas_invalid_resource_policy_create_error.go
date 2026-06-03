// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasInvalidResourcePolicyCreateError struct for ApiAtlasInvalidResourcePolicyCreateError
type ApiAtlasInvalidResourcePolicyCreateError struct {
	// Human-readable label that displays the type of an error.
	ErrorType *string `json:"errorType,omitempty"`
	// List of invalid policies containing details of their validation errors.
	// Read only field.
	InvalidPolicies *[]ApiAtlasInvalidPolicy `json:"invalidPolicies,omitempty"`
}

// NewApiAtlasInvalidResourcePolicyCreateError instantiates a new ApiAtlasInvalidResourcePolicyCreateError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasInvalidResourcePolicyCreateError() *ApiAtlasInvalidResourcePolicyCreateError {
	this := ApiAtlasInvalidResourcePolicyCreateError{}
	return &this
}

// NewApiAtlasInvalidResourcePolicyCreateErrorWithDefaults instantiates a new ApiAtlasInvalidResourcePolicyCreateError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasInvalidResourcePolicyCreateErrorWithDefaults() *ApiAtlasInvalidResourcePolicyCreateError {
	this := ApiAtlasInvalidResourcePolicyCreateError{}
	return &this
}

// GetErrorType returns the ErrorType field value if set, zero value otherwise
func (o *ApiAtlasInvalidResourcePolicyCreateError) GetErrorType() string {
	if o == nil || IsNil(o.ErrorType) {
		var ret string
		return ret
	}
	return *o.ErrorType
}

// GetErrorTypeOk returns a tuple with the ErrorType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasInvalidResourcePolicyCreateError) GetErrorTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ErrorType) {
		return nil, false
	}

	return o.ErrorType, true
}

// HasErrorType returns a boolean if a field has been set.
func (o *ApiAtlasInvalidResourcePolicyCreateError) HasErrorType() bool {
	if o != nil && !IsNil(o.ErrorType) {
		return true
	}

	return false
}

// SetErrorType gets a reference to the given string and assigns it to the ErrorType field.
func (o *ApiAtlasInvalidResourcePolicyCreateError) SetErrorType(v string) {
	o.ErrorType = &v
}

// GetInvalidPolicies returns the InvalidPolicies field value if set, zero value otherwise
func (o *ApiAtlasInvalidResourcePolicyCreateError) GetInvalidPolicies() []ApiAtlasInvalidPolicy {
	if o == nil || IsNil(o.InvalidPolicies) {
		var ret []ApiAtlasInvalidPolicy
		return ret
	}
	return *o.InvalidPolicies
}

// GetInvalidPoliciesOk returns a tuple with the InvalidPolicies field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasInvalidResourcePolicyCreateError) GetInvalidPoliciesOk() (*[]ApiAtlasInvalidPolicy, bool) {
	if o == nil || IsNil(o.InvalidPolicies) {
		return nil, false
	}

	return o.InvalidPolicies, true
}

// HasInvalidPolicies returns a boolean if a field has been set.
func (o *ApiAtlasInvalidResourcePolicyCreateError) HasInvalidPolicies() bool {
	if o != nil && !IsNil(o.InvalidPolicies) {
		return true
	}

	return false
}

// SetInvalidPolicies gets a reference to the given []ApiAtlasInvalidPolicy and assigns it to the InvalidPolicies field.
func (o *ApiAtlasInvalidResourcePolicyCreateError) SetInvalidPolicies(v []ApiAtlasInvalidPolicy) {
	o.InvalidPolicies = &v
}
