// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CloudAccessRoleAssignment MongoDB Cloud user's roles and the corresponding organization or project to which that role applies. Each role can apply to one organization or one project but not both.
type CloudAccessRoleAssignment struct {
	// Unique 24-hexadecimal digit string that identifies the project to which this role belongs. You can set a value for this parameter or `orgId` but not both in the same request.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the organization to which this role belongs. You can set a value for this parameter or `groupId` but not both in the same request.
	OrgId *string `json:"orgId,omitempty"`
	// Human-readable label that identifies the collection of privileges that MongoDB Cloud grants a specific API key, MongoDB Cloud user, or MongoDB Cloud team. These roles include organization- and project-level roles.
	RoleName *string `json:"roleName,omitempty"`
}

// NewCloudAccessRoleAssignment instantiates a new CloudAccessRoleAssignment object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudAccessRoleAssignment() *CloudAccessRoleAssignment {
	this := CloudAccessRoleAssignment{}
	return &this
}

// NewCloudAccessRoleAssignmentWithDefaults instantiates a new CloudAccessRoleAssignment object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudAccessRoleAssignmentWithDefaults() *CloudAccessRoleAssignment {
	this := CloudAccessRoleAssignment{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *CloudAccessRoleAssignment) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAccessRoleAssignment) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *CloudAccessRoleAssignment) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *CloudAccessRoleAssignment) SetGroupId(v string) {
	o.GroupId = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *CloudAccessRoleAssignment) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAccessRoleAssignment) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *CloudAccessRoleAssignment) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *CloudAccessRoleAssignment) SetOrgId(v string) {
	o.OrgId = &v
}

// GetRoleName returns the RoleName field value if set, zero value otherwise
func (o *CloudAccessRoleAssignment) GetRoleName() string {
	if o == nil || IsNil(o.RoleName) {
		var ret string
		return ret
	}
	return *o.RoleName
}

// GetRoleNameOk returns a tuple with the RoleName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAccessRoleAssignment) GetRoleNameOk() (*string, bool) {
	if o == nil || IsNil(o.RoleName) {
		return nil, false
	}

	return o.RoleName, true
}

// HasRoleName returns a boolean if a field has been set.
func (o *CloudAccessRoleAssignment) HasRoleName() bool {
	if o != nil && !IsNil(o.RoleName) {
		return true
	}

	return false
}

// SetRoleName gets a reference to the given string and assigns it to the RoleName field.
func (o *CloudAccessRoleAssignment) SetRoleName(v string) {
	o.RoleName = &v
}
