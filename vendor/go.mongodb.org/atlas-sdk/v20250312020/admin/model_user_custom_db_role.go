// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UserCustomDBRole struct for UserCustomDBRole
type UserCustomDBRole struct {
	// List of the individual privilege actions that the role grants.
	Actions *[]DatabasePrivilegeAction `json:"actions,omitempty"`
	// List of the built-in roles that this custom role inherits.
	InheritedRoles *[]DatabaseInheritedRole `json:"inheritedRoles,omitempty"`
	// Human-readable label that identifies the role for the request. This name must be unique for this custom role in this project.
	RoleName string `json:"roleName"`
}

// NewUserCustomDBRole instantiates a new UserCustomDBRole object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserCustomDBRole(roleName string) *UserCustomDBRole {
	this := UserCustomDBRole{}
	this.RoleName = roleName
	return &this
}

// NewUserCustomDBRoleWithDefaults instantiates a new UserCustomDBRole object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserCustomDBRoleWithDefaults() *UserCustomDBRole {
	this := UserCustomDBRole{}
	return &this
}

// GetActions returns the Actions field value if set, zero value otherwise
func (o *UserCustomDBRole) GetActions() []DatabasePrivilegeAction {
	if o == nil || IsNil(o.Actions) {
		var ret []DatabasePrivilegeAction
		return ret
	}
	return *o.Actions
}

// GetActionsOk returns a tuple with the Actions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCustomDBRole) GetActionsOk() (*[]DatabasePrivilegeAction, bool) {
	if o == nil || IsNil(o.Actions) {
		return nil, false
	}

	return o.Actions, true
}

// HasActions returns a boolean if a field has been set.
func (o *UserCustomDBRole) HasActions() bool {
	if o != nil && !IsNil(o.Actions) {
		return true
	}

	return false
}

// SetActions gets a reference to the given []DatabasePrivilegeAction and assigns it to the Actions field.
func (o *UserCustomDBRole) SetActions(v []DatabasePrivilegeAction) {
	o.Actions = &v
}

// GetInheritedRoles returns the InheritedRoles field value if set, zero value otherwise
func (o *UserCustomDBRole) GetInheritedRoles() []DatabaseInheritedRole {
	if o == nil || IsNil(o.InheritedRoles) {
		var ret []DatabaseInheritedRole
		return ret
	}
	return *o.InheritedRoles
}

// GetInheritedRolesOk returns a tuple with the InheritedRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserCustomDBRole) GetInheritedRolesOk() (*[]DatabaseInheritedRole, bool) {
	if o == nil || IsNil(o.InheritedRoles) {
		return nil, false
	}

	return o.InheritedRoles, true
}

// HasInheritedRoles returns a boolean if a field has been set.
func (o *UserCustomDBRole) HasInheritedRoles() bool {
	if o != nil && !IsNil(o.InheritedRoles) {
		return true
	}

	return false
}

// SetInheritedRoles gets a reference to the given []DatabaseInheritedRole and assigns it to the InheritedRoles field.
func (o *UserCustomDBRole) SetInheritedRoles(v []DatabaseInheritedRole) {
	o.InheritedRoles = &v
}

// GetRoleName returns the RoleName field value
func (o *UserCustomDBRole) GetRoleName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RoleName
}

// GetRoleNameOk returns a tuple with the RoleName field value
// and a boolean to check if the value has been set.
func (o *UserCustomDBRole) GetRoleNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RoleName, true
}

// SetRoleName sets field value
func (o *UserCustomDBRole) SetRoleName(v string) {
	o.RoleName = v
}
