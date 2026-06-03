// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrganizationInvitationRequest struct for OrganizationInvitationRequest
type OrganizationInvitationRequest struct {
	// List of projects that the user will be added to when they accept their invitation to the organization.
	GroupRoleAssignments *[]OrganizationInvitationGroupRoleAssignmentsRequest `json:"groupRoleAssignments,omitempty"`
	// One or more organization level roles to assign to the MongoDB Cloud user.
	Roles *[]string `json:"roles,omitempty"`
	// List of teams to which you want to invite the desired MongoDB Cloud user.
	TeamIds *[]string `json:"teamIds,omitempty"`
	// Email address that belongs to the desired MongoDB Cloud user.
	Username *string `json:"username,omitempty"`
}

// NewOrganizationInvitationRequest instantiates a new OrganizationInvitationRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrganizationInvitationRequest() *OrganizationInvitationRequest {
	this := OrganizationInvitationRequest{}
	return &this
}

// NewOrganizationInvitationRequestWithDefaults instantiates a new OrganizationInvitationRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrganizationInvitationRequestWithDefaults() *OrganizationInvitationRequest {
	this := OrganizationInvitationRequest{}
	return &this
}

// GetGroupRoleAssignments returns the GroupRoleAssignments field value if set, zero value otherwise
func (o *OrganizationInvitationRequest) GetGroupRoleAssignments() []OrganizationInvitationGroupRoleAssignmentsRequest {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		var ret []OrganizationInvitationGroupRoleAssignmentsRequest
		return ret
	}
	return *o.GroupRoleAssignments
}

// GetGroupRoleAssignmentsOk returns a tuple with the GroupRoleAssignments field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationInvitationRequest) GetGroupRoleAssignmentsOk() (*[]OrganizationInvitationGroupRoleAssignmentsRequest, bool) {
	if o == nil || IsNil(o.GroupRoleAssignments) {
		return nil, false
	}

	return o.GroupRoleAssignments, true
}

// HasGroupRoleAssignments returns a boolean if a field has been set.
func (o *OrganizationInvitationRequest) HasGroupRoleAssignments() bool {
	if o != nil && !IsNil(o.GroupRoleAssignments) {
		return true
	}

	return false
}

// SetGroupRoleAssignments gets a reference to the given []OrganizationInvitationGroupRoleAssignmentsRequest and assigns it to the GroupRoleAssignments field.
func (o *OrganizationInvitationRequest) SetGroupRoleAssignments(v []OrganizationInvitationGroupRoleAssignmentsRequest) {
	o.GroupRoleAssignments = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *OrganizationInvitationRequest) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationInvitationRequest) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *OrganizationInvitationRequest) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *OrganizationInvitationRequest) SetRoles(v []string) {
	o.Roles = &v
}

// GetTeamIds returns the TeamIds field value if set, zero value otherwise
func (o *OrganizationInvitationRequest) GetTeamIds() []string {
	if o == nil || IsNil(o.TeamIds) {
		var ret []string
		return ret
	}
	return *o.TeamIds
}

// GetTeamIdsOk returns a tuple with the TeamIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationInvitationRequest) GetTeamIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.TeamIds) {
		return nil, false
	}

	return o.TeamIds, true
}

// HasTeamIds returns a boolean if a field has been set.
func (o *OrganizationInvitationRequest) HasTeamIds() bool {
	if o != nil && !IsNil(o.TeamIds) {
		return true
	}

	return false
}

// SetTeamIds gets a reference to the given []string and assigns it to the TeamIds field.
func (o *OrganizationInvitationRequest) SetTeamIds(v []string) {
	o.TeamIds = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *OrganizationInvitationRequest) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationInvitationRequest) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *OrganizationInvitationRequest) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *OrganizationInvitationRequest) SetUsername(v string) {
	o.Username = &v
}
