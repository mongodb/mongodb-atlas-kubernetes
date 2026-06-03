// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupInvitationUpdateRequest struct for GroupInvitationUpdateRequest
type GroupInvitationUpdateRequest struct {
	// One or more project-level roles to assign to the MongoDB Cloud user.
	Roles *[]string `json:"roles,omitempty"`
}

// NewGroupInvitationUpdateRequest instantiates a new GroupInvitationUpdateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupInvitationUpdateRequest() *GroupInvitationUpdateRequest {
	this := GroupInvitationUpdateRequest{}
	return &this
}

// NewGroupInvitationUpdateRequestWithDefaults instantiates a new GroupInvitationUpdateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupInvitationUpdateRequestWithDefaults() *GroupInvitationUpdateRequest {
	this := GroupInvitationUpdateRequest{}
	return &this
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *GroupInvitationUpdateRequest) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitationUpdateRequest) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *GroupInvitationUpdateRequest) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *GroupInvitationUpdateRequest) SetRoles(v []string) {
	o.Roles = &v
}
