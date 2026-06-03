// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasUserMetadata The user that last updated the atlas resource policy.
type ApiAtlasUserMetadata struct {
	// Unique 24-hexadecimal character string that identifies a user.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Human-readable label that describes a user.
	// Read only field.
	Name *string `json:"name,omitempty"`
}

// NewApiAtlasUserMetadata instantiates a new ApiAtlasUserMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasUserMetadata() *ApiAtlasUserMetadata {
	this := ApiAtlasUserMetadata{}
	return &this
}

// NewApiAtlasUserMetadataWithDefaults instantiates a new ApiAtlasUserMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasUserMetadataWithDefaults() *ApiAtlasUserMetadata {
	this := ApiAtlasUserMetadata{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiAtlasUserMetadata) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasUserMetadata) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiAtlasUserMetadata) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiAtlasUserMetadata) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *ApiAtlasUserMetadata) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasUserMetadata) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ApiAtlasUserMetadata) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ApiAtlasUserMetadata) SetName(v string) {
	o.Name = &v
}
