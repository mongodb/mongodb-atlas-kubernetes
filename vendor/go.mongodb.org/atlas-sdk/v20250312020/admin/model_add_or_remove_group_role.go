// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AddOrRemoveGroupRole struct for AddOrRemoveGroupRole
type AddOrRemoveGroupRole struct {
	// Project-level role to assign to or remove from the MongoDB Cloud user.
	GroupRole string `json:"groupRole"`
}

// NewAddOrRemoveGroupRole instantiates a new AddOrRemoveGroupRole object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAddOrRemoveGroupRole(groupRole string) *AddOrRemoveGroupRole {
	this := AddOrRemoveGroupRole{}
	this.GroupRole = groupRole
	return &this
}

// NewAddOrRemoveGroupRoleWithDefaults instantiates a new AddOrRemoveGroupRole object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAddOrRemoveGroupRoleWithDefaults() *AddOrRemoveGroupRole {
	this := AddOrRemoveGroupRole{}
	return &this
}

// GetGroupRole returns the GroupRole field value
func (o *AddOrRemoveGroupRole) GetGroupRole() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GroupRole
}

// GetGroupRoleOk returns a tuple with the GroupRole field value
// and a boolean to check if the value has been set.
func (o *AddOrRemoveGroupRole) GetGroupRoleOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GroupRole, true
}

// SetGroupRole sets field value
func (o *AddOrRemoveGroupRole) SetGroupRole(v string) {
	o.GroupRole = v
}
