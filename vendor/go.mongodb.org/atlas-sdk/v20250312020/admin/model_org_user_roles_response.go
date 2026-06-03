// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgUserRolesResponse Organization- and project-level roles assigned to one MongoDB Cloud user within one organization.
type OrgUserRolesResponse struct {
	// List of project-level role assignments assigned to the MongoDB Cloud user.
	GroupRoleAssignments *[]GroupRoleAssignment `json:"groupRoleAssignments,omitempty"`
	// One or more organization-level roles assigned to the MongoDB Cloud user.
	OrgRoles *[]string `json:"orgRoles,omitempty"`
}

// NewOrgUserRolesResponse instantiates a new OrgUserRolesResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgUserRolesResponse() *OrgUserRolesResponse {
	this := OrgUserRolesResponse{}
	return &this
}

// NewOrgUserRolesResponseWithDefaults instantiates a new OrgUserRolesResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgUserRolesResponseWithDefaults() *OrgUserRolesResponse {
	this := OrgUserRolesResponse{}
	return &this
}

// GetGroupRoleAssignments returns the GroupRoleAssignments field value if set, zero value otherwise
func (o *OrgUserRolesResponse) GetGroupRoleAssignments() []GroupRoleAssignment {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		var ret []GroupRoleAssignment
		return ret
	}
	return *o.GroupRoleAssignments
}

// GetGroupRoleAssignmentsOk returns a tuple with the GroupRoleAssignments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserRolesResponse) GetGroupRoleAssignmentsOk() (*[]GroupRoleAssignment, bool) {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		return nil, false
	}

	return o.GroupRoleAssignments, true
}

// HasGroupRoleAssignments returns a boolean if a field has been set.
func (o *OrgUserRolesResponse) HasGroupRoleAssignments() bool {
	if o != nil && !IsNil(o.GroupRoleAssignments) {
		return true
	}

	return false
}

// SetGroupRoleAssignments gets a reference to the given []GroupRoleAssignment and assigns it to the GroupRoleAssignments field.
func (o *OrgUserRolesResponse) SetGroupRoleAssignments(v []GroupRoleAssignment) {
	o.GroupRoleAssignments = &v
}

// GetOrgRoles returns the OrgRoles field value if set, zero value otherwise
func (o *OrgUserRolesResponse) GetOrgRoles() []string {
	if o == nil || IsNil(o.OrgRoles) {
		var ret []string
		return ret
	}
	return *o.OrgRoles
}

// GetOrgRolesOk returns a tuple with the OrgRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserRolesResponse) GetOrgRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.OrgRoles) {
		return nil, false
	}

	return o.OrgRoles, true
}

// HasOrgRoles returns a boolean if a field has been set.
func (o *OrgUserRolesResponse) HasOrgRoles() bool {
	if o != nil && !IsNil(o.OrgRoles) {
		return true
	}

	return false
}

// SetOrgRoles gets a reference to the given []string and assigns it to the OrgRoles field.
func (o *OrgUserRolesResponse) SetOrgRoles(v []string) {
	o.OrgRoles = &v
}
