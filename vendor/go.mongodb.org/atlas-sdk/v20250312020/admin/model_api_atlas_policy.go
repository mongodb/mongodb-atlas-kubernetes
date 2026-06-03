// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasPolicy struct for ApiAtlasPolicy
type ApiAtlasPolicy struct {
	// A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.
	// Read only field.
	Body *string `json:"body,omitempty"`
	// Unique 24-hexadecimal character string that identifies the policy.
	// Read only field.
	Id *string `json:"id,omitempty"`
}

// NewApiAtlasPolicy instantiates a new ApiAtlasPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasPolicy() *ApiAtlasPolicy {
	this := ApiAtlasPolicy{}
	return &this
}

// NewApiAtlasPolicyWithDefaults instantiates a new ApiAtlasPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasPolicyWithDefaults() *ApiAtlasPolicy {
	this := ApiAtlasPolicy{}
	return &this
}

// GetBody returns the Body field value if set, zero value otherwise
func (o *ApiAtlasPolicy) GetBody() string {
	if o == nil || IsNil(o.Body) {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasPolicy) GetBodyOk() (*string, bool) {
	if o == nil || IsNil(o.Body) {
		return nil, false
	}

	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *ApiAtlasPolicy) HasBody() bool {
	if o != nil && !IsNil(o.Body) {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *ApiAtlasPolicy) SetBody(v string) {
	o.Body = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiAtlasPolicy) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasPolicy) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiAtlasPolicy) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiAtlasPolicy) SetId(v string) {
	o.Id = &v
}
