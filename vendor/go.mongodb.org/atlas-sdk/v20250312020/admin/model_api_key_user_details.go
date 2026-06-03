// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiKeyUserDetails Details of the Programmatic API Keys.
type ApiKeyUserDetails struct {
	// Purpose or explanation provided when someone created this organization API key.
	Desc *string `json:"desc,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this organization API key assigned to this project.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Redacted private key returned for this organization API key. This key displays unredacted when first created.
	// Read only field.
	PrivateKey *string `json:"privateKey,omitempty"`
	// Public API key value set for the specified organization API key.
	// Read only field.
	PublicKey *string `json:"publicKey,omitempty"`
	// List that contains the roles that the API key needs to have. All roles you provide must be valid for the specified project or organization. Each request must include a minimum of one valid role. The resource returns all project and organization roles assigned to the API key.
	Roles *[]CloudAccessRoleAssignment `json:"roles,omitempty"`
}

// NewApiKeyUserDetails instantiates a new ApiKeyUserDetails object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiKeyUserDetails() *ApiKeyUserDetails {
	this := ApiKeyUserDetails{}
	return &this
}

// NewApiKeyUserDetailsWithDefaults instantiates a new ApiKeyUserDetails object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiKeyUserDetailsWithDefaults() *ApiKeyUserDetails {
	this := ApiKeyUserDetails{}
	return &this
}

// GetDesc returns the Desc field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetDesc() string {
	if o == nil || IsNil(o.Desc) {
		var ret string
		return ret
	}
	return *o.Desc
}

// GetDescOk returns a tuple with the Desc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetDescOk() (*string, bool) {
	if o == nil || IsNil(o.Desc) {
		return nil, false
	}

	return o.Desc, true
}

// HasDesc returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasDesc() bool {
	if o != nil && !IsNil(o.Desc) {
		return true
	}

	return false
}

// SetDesc gets a reference to the given string and assigns it to the Desc field.
func (o *ApiKeyUserDetails) SetDesc(v string) {
	o.Desc = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ApiKeyUserDetails) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ApiKeyUserDetails) SetLinks(v []Link) {
	o.Links = &v
}

// GetPrivateKey returns the PrivateKey field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetPrivateKey() string {
	if o == nil || IsNil(o.PrivateKey) {
		var ret string
		return ret
	}
	return *o.PrivateKey
}

// GetPrivateKeyOk returns a tuple with the PrivateKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetPrivateKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateKey) {
		return nil, false
	}

	return o.PrivateKey, true
}

// HasPrivateKey returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasPrivateKey() bool {
	if o != nil && !IsNil(o.PrivateKey) {
		return true
	}

	return false
}

// SetPrivateKey gets a reference to the given string and assigns it to the PrivateKey field.
func (o *ApiKeyUserDetails) SetPrivateKey(v string) {
	o.PrivateKey = &v
}

// GetPublicKey returns the PublicKey field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetPublicKey() string {
	if o == nil || IsNil(o.PublicKey) {
		var ret string
		return ret
	}
	return *o.PublicKey
}

// GetPublicKeyOk returns a tuple with the PublicKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetPublicKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PublicKey) {
		return nil, false
	}

	return o.PublicKey, true
}

// HasPublicKey returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasPublicKey() bool {
	if o != nil && !IsNil(o.PublicKey) {
		return true
	}

	return false
}

// SetPublicKey gets a reference to the given string and assigns it to the PublicKey field.
func (o *ApiKeyUserDetails) SetPublicKey(v string) {
	o.PublicKey = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *ApiKeyUserDetails) GetRoles() []CloudAccessRoleAssignment {
	if o == nil || IsNil(o.Roles) {
		var ret []CloudAccessRoleAssignment
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiKeyUserDetails) GetRolesOk() (*[]CloudAccessRoleAssignment, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *ApiKeyUserDetails) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []CloudAccessRoleAssignment and assigns it to the Roles field.
func (o *ApiKeyUserDetails) SetRoles(v []CloudAccessRoleAssignment) {
	o.Roles = &v
}
