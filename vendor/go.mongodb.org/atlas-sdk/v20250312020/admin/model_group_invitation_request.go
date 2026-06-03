// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupInvitationRequest struct for GroupInvitationRequest
type GroupInvitationRequest struct {
	// One or more project level roles to assign to the MongoDB Cloud user.
	Roles *[]string `json:"roles,omitempty"`
	// Email address of the MongoDB Cloud user invited to the specified project.
	Username *string `json:"username,omitempty"`
}

// NewGroupInvitationRequest instantiates a new GroupInvitationRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupInvitationRequest() *GroupInvitationRequest {
	this := GroupInvitationRequest{}
	return &this
}

// NewGroupInvitationRequestWithDefaults instantiates a new GroupInvitationRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupInvitationRequestWithDefaults() *GroupInvitationRequest {
	this := GroupInvitationRequest{}
	return &this
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *GroupInvitationRequest) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitationRequest) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *GroupInvitationRequest) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *GroupInvitationRequest) SetRoles(v []string) {
	o.Roles = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *GroupInvitationRequest) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitationRequest) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *GroupInvitationRequest) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *GroupInvitationRequest) SetUsername(v string) {
	o.Username = &v
}
