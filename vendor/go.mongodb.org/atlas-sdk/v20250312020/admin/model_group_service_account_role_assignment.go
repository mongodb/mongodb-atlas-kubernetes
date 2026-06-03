// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupServiceAccountRoleAssignment struct for GroupServiceAccountRoleAssignment
type GroupServiceAccountRoleAssignment struct {
	// The Project permissions for the Service Account in the specified Project.
	Roles []string `json:"roles"`
}

// NewGroupServiceAccountRoleAssignment instantiates a new GroupServiceAccountRoleAssignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupServiceAccountRoleAssignment(roles []string) *GroupServiceAccountRoleAssignment {
	this := GroupServiceAccountRoleAssignment{}
	this.Roles = roles
	return &this
}

// NewGroupServiceAccountRoleAssignmentWithDefaults instantiates a new GroupServiceAccountRoleAssignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupServiceAccountRoleAssignmentWithDefaults() *GroupServiceAccountRoleAssignment {
	this := GroupServiceAccountRoleAssignment{}
	return &this
}

// GetRoles returns the Roles field value
func (o *GroupServiceAccountRoleAssignment) GetRoles() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *GroupServiceAccountRoleAssignment) GetRolesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *GroupServiceAccountRoleAssignment) SetRoles(v []string) {
	o.Roles = v
}
