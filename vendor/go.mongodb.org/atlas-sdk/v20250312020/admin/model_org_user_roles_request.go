// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgUserRolesRequest Organization and project level roles to assign the MongoDB Cloud user within one organization.
type OrgUserRolesRequest struct {
	// List of project level role assignments to assign the MongoDB Cloud user.
	GroupRoleAssignments *[]GroupRoleAssignment `json:"groupRoleAssignments,omitempty"`
	// One or more organization level roles to assign the MongoDB Cloud user.
	OrgRoles []string `json:"orgRoles"`
}

// NewOrgUserRolesRequest instantiates a new OrgUserRolesRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgUserRolesRequest(orgRoles []string) *OrgUserRolesRequest {
	this := OrgUserRolesRequest{}
	this.OrgRoles = orgRoles
	return &this
}

// NewOrgUserRolesRequestWithDefaults instantiates a new OrgUserRolesRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgUserRolesRequestWithDefaults() *OrgUserRolesRequest {
	this := OrgUserRolesRequest{}
	return &this
}

// GetGroupRoleAssignments returns the GroupRoleAssignments field value if set, zero value otherwise
func (o *OrgUserRolesRequest) GetGroupRoleAssignments() []GroupRoleAssignment {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		var ret []GroupRoleAssignment
		return ret
	}
	return *o.GroupRoleAssignments
}

// GetGroupRoleAssignmentsOk returns a tuple with the GroupRoleAssignments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserRolesRequest) GetGroupRoleAssignmentsOk() (*[]GroupRoleAssignment, bool) {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		return nil, false
	}

	return o.GroupRoleAssignments, true
}

// HasGroupRoleAssignments returns a boolean if a field has been set.
func (o *OrgUserRolesRequest) HasGroupRoleAssignments() bool {
	if o != nil && !IsNil(o.GroupRoleAssignments) {
		return true
	}

	return false
}

// SetGroupRoleAssignments gets a reference to the given []GroupRoleAssignment and assigns it to the GroupRoleAssignments field.
func (o *OrgUserRolesRequest) SetGroupRoleAssignments(v []GroupRoleAssignment) {
	o.GroupRoleAssignments = &v
}

// GetOrgRoles returns the OrgRoles field value
func (o *OrgUserRolesRequest) GetOrgRoles() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.OrgRoles
}

// GetOrgRolesOk returns a tuple with the OrgRoles field value
// and a boolean to check if the value has been set.
func (o *OrgUserRolesRequest) GetOrgRolesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrgRoles, true
}

// SetOrgRoles sets field value
func (o *OrgUserRolesRequest) SetOrgRoles(v []string) {
	o.OrgRoles = v
}
