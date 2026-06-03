// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasInvalidPolicy struct for ApiAtlasInvalidPolicy
type ApiAtlasInvalidPolicy struct {
	// A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.
	// Read only field.
	Body *string `json:"body,omitempty"`
	// List of validation errors.
	// Read only field.
	Errors *[]ApiAtlasInvalidPolicyErrorDetail `json:"errors,omitempty"`
}

// NewApiAtlasInvalidPolicy instantiates a new ApiAtlasInvalidPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasInvalidPolicy() *ApiAtlasInvalidPolicy {
	this := ApiAtlasInvalidPolicy{}
	return &this
}

// NewApiAtlasInvalidPolicyWithDefaults instantiates a new ApiAtlasInvalidPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasInvalidPolicyWithDefaults() *ApiAtlasInvalidPolicy {
	this := ApiAtlasInvalidPolicy{}
	return &this
}

// GetBody returns the Body field value if set, zero value otherwise
func (o *ApiAtlasInvalidPolicy) GetBody() string {
	if o == nil || IsNil(o.Body) {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasInvalidPolicy) GetBodyOk() (*string, bool) {
	if o == nil || IsNil(o.Body) {
		return nil, false
	}

	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *ApiAtlasInvalidPolicy) HasBody() bool {
	if o != nil && !IsNil(o.Body) {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *ApiAtlasInvalidPolicy) SetBody(v string) {
	o.Body = &v
}

// GetErrors returns the Errors field value if set, zero value otherwise
func (o *ApiAtlasInvalidPolicy) GetErrors() []ApiAtlasInvalidPolicyErrorDetail {
	if o == nil || IsNil(o.Errors) {
		var ret []ApiAtlasInvalidPolicyErrorDetail
		return ret
	}
	return *o.Errors
}

// GetErrorsOk returns a tuple with the Errors field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasInvalidPolicy) GetErrorsOk() (*[]ApiAtlasInvalidPolicyErrorDetail, bool) {
	if o == nil || IsNil(o.Errors) {
		return nil, false
	}

	return o.Errors, true
}

// HasErrors returns a boolean if a field has been set.
func (o *ApiAtlasInvalidPolicy) HasErrors() bool {
	if o != nil && !IsNil(o.Errors) {
		return true
	}

	return false
}

// SetErrors gets a reference to the given []ApiAtlasInvalidPolicyErrorDetail and assigns it to the Errors field.
func (o *ApiAtlasInvalidPolicy) SetErrors(v []ApiAtlasInvalidPolicyErrorDetail) {
	o.Errors = &v
}
