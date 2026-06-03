// Code based on the AtlasAPI V2 OpenAPI file

package admin

// GroupUserRequest struct for GroupUserRequest
type GroupUserRequest struct {
	// One or more project-level roles to assign the MongoDB Cloud user.
	// Write only field.
	Roles []string `json:"roles"`
	// Email address that represents the username of the MongoDB Cloud user.
	// Write only field.
	Username string `json:"username"`
}

// NewGroupUserRequest instantiates a new GroupUserRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupUserRequest(roles []string, username string) *GroupUserRequest {
	this := GroupUserRequest{}
	this.Roles = roles
	this.Username = username
	return &this
}

// NewGroupUserRequestWithDefaults instantiates a new GroupUserRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupUserRequestWithDefaults() *GroupUserRequest {
	this := GroupUserRequest{}
	return &this
}

// GetRoles returns the Roles field value
func (o *GroupUserRequest) GetRoles() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *GroupUserRequest) GetRolesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *GroupUserRequest) SetRoles(v []string) {
	o.Roles = v
}

// GetUsername returns the Username field value
func (o *GroupUserRequest) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *GroupUserRequest) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *GroupUserRequest) SetUsername(v string) {
	o.Username = v
}
