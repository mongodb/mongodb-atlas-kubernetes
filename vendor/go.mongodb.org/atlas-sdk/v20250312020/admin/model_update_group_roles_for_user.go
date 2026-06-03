// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UpdateGroupRolesForUser struct for UpdateGroupRolesForUser
type UpdateGroupRolesForUser struct {
	// One or more project-level roles to assign to the MongoDB Cloud user.
	GroupRoles *[]string `json:"groupRoles,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewUpdateGroupRolesForUser instantiates a new UpdateGroupRolesForUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateGroupRolesForUser() *UpdateGroupRolesForUser {
	this := UpdateGroupRolesForUser{}
	return &this
}

// NewUpdateGroupRolesForUserWithDefaults instantiates a new UpdateGroupRolesForUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateGroupRolesForUserWithDefaults() *UpdateGroupRolesForUser {
	this := UpdateGroupRolesForUser{}
	return &this
}

// GetGroupRoles returns the GroupRoles field value if set, zero value otherwise
func (o *UpdateGroupRolesForUser) GetGroupRoles() []string {
	if o == nil || IsNil(o.GroupRoles) {
		var ret []string
		return ret
	}
	return *o.GroupRoles
}

// GetGroupRolesOk returns a tuple with the GroupRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateGroupRolesForUser) GetGroupRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.GroupRoles) {
		return nil, false
	}

	return o.GroupRoles, true
}

// HasGroupRoles returns a boolean if a field has been set.
func (o *UpdateGroupRolesForUser) HasGroupRoles() bool {
	if o != nil && !IsNil(o.GroupRoles) {
		return true
	}

	return false
}

// SetGroupRoles gets a reference to the given []string and assigns it to the GroupRoles field.
func (o *UpdateGroupRolesForUser) SetGroupRoles(v []string) {
	o.GroupRoles = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *UpdateGroupRolesForUser) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateGroupRolesForUser) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *UpdateGroupRolesForUser) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *UpdateGroupRolesForUser) SetLinks(v []Link) {
	o.Links = &v
}
