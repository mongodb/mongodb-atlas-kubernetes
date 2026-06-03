// Code based on the AtlasAPI V2 OpenAPI file

package admin

// SchemaRegistryAuthentication Authentication configuration for Schema Registry.
type SchemaRegistryAuthentication struct {
	// Authentication type discriminator. Specifies the authentication mechanism for Confluent Schema Registry.
	Type string `json:"type"`
	// Password or Private Key for authentication.
	// Write only field.
	Password *string `json:"password,omitempty"`
	// Username or Public Key for authentication.
	Username *string `json:"username,omitempty"`
}

// NewSchemaRegistryAuthentication instantiates a new SchemaRegistryAuthentication object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSchemaRegistryAuthentication(type_ string) *SchemaRegistryAuthentication {
	this := SchemaRegistryAuthentication{}
	this.Type = type_
	return &this
}

// NewSchemaRegistryAuthenticationWithDefaults instantiates a new SchemaRegistryAuthentication object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSchemaRegistryAuthenticationWithDefaults() *SchemaRegistryAuthentication {
	this := SchemaRegistryAuthentication{}
	return &this
}

// GetType returns the Type field value
func (o *SchemaRegistryAuthentication) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *SchemaRegistryAuthentication) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *SchemaRegistryAuthentication) SetType(v string) {
	o.Type = v
}

// GetPassword returns the Password field value if set, zero value otherwise
func (o *SchemaRegistryAuthentication) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaRegistryAuthentication) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}

	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *SchemaRegistryAuthentication) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *SchemaRegistryAuthentication) SetPassword(v string) {
	o.Password = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *SchemaRegistryAuthentication) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SchemaRegistryAuthentication) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *SchemaRegistryAuthentication) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *SchemaRegistryAuthentication) SetUsername(v string) {
	o.Username = &v
}
