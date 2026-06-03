// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ConnectedOrgConfigRoleAssignment struct for ConnectedOrgConfigRoleAssignment
type ConnectedOrgConfigRoleAssignment struct {
	// Unique 24-hexadecimal digit string that identifies the project to which this role belongs. Each element within `roleAssignments` can have a value for `groupId` or `orgId`, but not both.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization to which this role belongs. Each element within `roleAssignments` can have a value for `orgId` or `groupId`, but not both.
	OrgId *string `json:"orgId,omitempty"`
	// Human-readable label that identifies the collection of privileges that MongoDB Cloud grants a specific API key, MongoDB Cloud user, or MongoDB Cloud team. These roles include organization- and project-level roles.
	Role *string `json:"role,omitempty"`
}

// NewConnectedOrgConfigRoleAssignment instantiates a new ConnectedOrgConfigRoleAssignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectedOrgConfigRoleAssignment() *ConnectedOrgConfigRoleAssignment {
	this := ConnectedOrgConfigRoleAssignment{}
	return &this
}

// NewConnectedOrgConfigRoleAssignmentWithDefaults instantiates a new ConnectedOrgConfigRoleAssignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectedOrgConfigRoleAssignmentWithDefaults() *ConnectedOrgConfigRoleAssignment {
	this := ConnectedOrgConfigRoleAssignment{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *ConnectedOrgConfigRoleAssignment) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfigRoleAssignment) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *ConnectedOrgConfigRoleAssignment) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *ConnectedOrgConfigRoleAssignment) SetGroupId(v string) {
	o.GroupId = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *ConnectedOrgConfigRoleAssignment) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfigRoleAssignment) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *ConnectedOrgConfigRoleAssignment) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *ConnectedOrgConfigRoleAssignment) SetOrgId(v string) {
	o.OrgId = &v
}

// GetRole returns the Role field value if set, zero value otherwise
func (o *ConnectedOrgConfigRoleAssignment) GetRole() string {
	if o == nil || IsNil(o.Role) {
		var ret string
		return ret
	}
	return *o.Role
}

// GetRoleOk returns a tuple with the Role field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfigRoleAssignment) GetRoleOk() (*string, bool) {
	if o == nil || IsNil(o.Role) {
		return nil, false
	}

	return o.Role, true
}

// HasRole returns a boolean if a field has been set.
func (o *ConnectedOrgConfigRoleAssignment) HasRole() bool {
	if o != nil && !IsNil(o.Role) {
		return true
	}

	return false
}

// SetRole gets a reference to the given string and assigns it to the Role field.
func (o *ConnectedOrgConfigRoleAssignment) SetRole(v string) {
	o.Role = &v
}
