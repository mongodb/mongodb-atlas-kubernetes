// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ServerlessTenantCreateRequest struct for ServerlessTenantCreateRequest
type ServerlessTenantCreateRequest struct {
	// Human-readable comment associated with the private endpoint.
	// Write only field.
	Comment *string `json:"comment,omitempty"`
}

// NewServerlessTenantCreateRequest instantiates a new ServerlessTenantCreateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessTenantCreateRequest() *ServerlessTenantCreateRequest {
	this := ServerlessTenantCreateRequest{}
	return &this
}

// NewServerlessTenantCreateRequestWithDefaults instantiates a new ServerlessTenantCreateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessTenantCreateRequestWithDefaults() *ServerlessTenantCreateRequest {
	this := ServerlessTenantCreateRequest{}
	return &this
}

// GetComment returns the Comment field value if set, zero value otherwise
func (o *ServerlessTenantCreateRequest) GetComment() string {
	if o == nil || IsNil(o.Comment) {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessTenantCreateRequest) GetCommentOk() (*string, bool) {
	if o == nil || IsNil(o.Comment) {
		return nil, false
	}

	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *ServerlessTenantCreateRequest) HasComment() bool {
	if o != nil && !IsNil(o.Comment) {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *ServerlessTenantCreateRequest) SetComment(v string) {
	o.Comment = &v
}
