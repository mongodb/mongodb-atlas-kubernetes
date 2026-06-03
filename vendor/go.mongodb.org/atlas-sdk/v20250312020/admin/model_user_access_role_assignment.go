// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UserAccessRoleAssignment struct for UserAccessRoleAssignment
type UserAccessRoleAssignment struct {
	// List of roles to grant this API key. If you provide this list, provide a minimum of one role and ensure each role applies to this project.
	Roles *[]string `json:"roles,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization API key.
	// Read only field.
	UserId *string `json:"userId,omitempty"`
}

// NewUserAccessRoleAssignment instantiates a new UserAccessRoleAssignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAccessRoleAssignment() *UserAccessRoleAssignment {
	this := UserAccessRoleAssignment{}
	return &this
}

// NewUserAccessRoleAssignmentWithDefaults instantiates a new UserAccessRoleAssignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAccessRoleAssignmentWithDefaults() *UserAccessRoleAssignment {
	this := UserAccessRoleAssignment{}
	return &this
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *UserAccessRoleAssignment) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessRoleAssignment) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *UserAccessRoleAssignment) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *UserAccessRoleAssignment) SetRoles(v []string) {
	o.Roles = &v
}

// GetUserId returns the UserId field value if set, zero value otherwise
func (o *UserAccessRoleAssignment) GetUserId() string {
	if o == nil || IsNil(o.UserId) {
		var ret string
		return ret
	}
	return *o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessRoleAssignment) GetUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.UserId) {
		return nil, false
	}

	return o.UserId, true
}

// HasUserId returns a boolean if a field has been set.
func (o *UserAccessRoleAssignment) HasUserId() bool {
	if o != nil && !IsNil(o.UserId) {
		return true
	}

	return false
}

// SetUserId gets a reference to the given string and assigns it to the UserId field.
func (o *UserAccessRoleAssignment) SetUserId(v string) {
	o.UserId = &v
}
