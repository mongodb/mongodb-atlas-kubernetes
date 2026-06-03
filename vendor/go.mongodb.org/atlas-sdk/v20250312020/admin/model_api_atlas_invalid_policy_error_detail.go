// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasInvalidPolicyErrorDetail struct for ApiAtlasInvalidPolicyErrorDetail
type ApiAtlasInvalidPolicyErrorDetail struct {
	// A string that provides a detailed description of a validation error.
	// Read only field.
	Detail *string `json:"detail,omitempty"`
}

// NewApiAtlasInvalidPolicyErrorDetail instantiates a new ApiAtlasInvalidPolicyErrorDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasInvalidPolicyErrorDetail() *ApiAtlasInvalidPolicyErrorDetail {
	this := ApiAtlasInvalidPolicyErrorDetail{}
	return &this
}

// NewApiAtlasInvalidPolicyErrorDetailWithDefaults instantiates a new ApiAtlasInvalidPolicyErrorDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasInvalidPolicyErrorDetailWithDefaults() *ApiAtlasInvalidPolicyErrorDetail {
	this := ApiAtlasInvalidPolicyErrorDetail{}
	return &this
}

// GetDetail returns the Detail field value if set, zero value otherwise
func (o *ApiAtlasInvalidPolicyErrorDetail) GetDetail() string {
	if o == nil || IsNil(o.Detail) {
		var ret string
		return ret
	}
	return *o.Detail
}

// GetDetailOk returns a tuple with the Detail field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasInvalidPolicyErrorDetail) GetDetailOk() (*string, bool) {
	if o == nil || IsNil(o.Detail) {
		return nil, false
	}

	return o.Detail, true
}

// HasDetail returns a boolean if a field has been set.
func (o *ApiAtlasInvalidPolicyErrorDetail) HasDetail() bool {
	if o != nil && !IsNil(o.Detail) {
		return true
	}

	return false
}

// SetDetail gets a reference to the given string and assigns it to the Detail field.
func (o *ApiAtlasInvalidPolicyErrorDetail) SetDetail(v string) {
	o.Detail = &v
}
