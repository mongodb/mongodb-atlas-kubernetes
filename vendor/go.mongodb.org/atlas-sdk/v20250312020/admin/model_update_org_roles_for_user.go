// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UpdateOrgRolesForUser struct for UpdateOrgRolesForUser
type UpdateOrgRolesForUser struct {
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// One or more organization level roles to assign to the MongoDB Cloud user.
	OrgRoles *[]string `json:"orgRoles,omitempty"`
}

// NewUpdateOrgRolesForUser instantiates a new UpdateOrgRolesForUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateOrgRolesForUser() *UpdateOrgRolesForUser {
	this := UpdateOrgRolesForUser{}
	return &this
}

// NewUpdateOrgRolesForUserWithDefaults instantiates a new UpdateOrgRolesForUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateOrgRolesForUserWithDefaults() *UpdateOrgRolesForUser {
	this := UpdateOrgRolesForUser{}
	return &this
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *UpdateOrgRolesForUser) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateOrgRolesForUser) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *UpdateOrgRolesForUser) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *UpdateOrgRolesForUser) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgRoles returns the OrgRoles field value if set, zero value otherwise
func (o *UpdateOrgRolesForUser) GetOrgRoles() []string {
	if o == nil || IsNil(o.OrgRoles) {
		var ret []string
		return ret
	}
	return *o.OrgRoles
}

// GetOrgRolesOk returns a tuple with the OrgRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateOrgRolesForUser) GetOrgRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.OrgRoles) {
		return nil, false
	}

	return o.OrgRoles, true
}

// HasOrgRoles returns a boolean if a field has been set.
func (o *UpdateOrgRolesForUser) HasOrgRoles() bool {
	if o != nil && !IsNil(o.OrgRoles) {
		return true
	}

	return false
}

// SetOrgRoles gets a reference to the given []string and assigns it to the OrgRoles field.
func (o *UpdateOrgRolesForUser) SetOrgRoles(v []string) {
	o.OrgRoles = &v
}
