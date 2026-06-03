// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgUserRequest struct for OrgUserRequest
type OrgUserRequest struct {
	Roles OrgUserRolesRequest `json:"roles"`
	// List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.
	// Write only field.
	TeamIds *[]string `json:"teamIds,omitempty"`
	// Email address that represents the username of the MongoDB Cloud user.
	// Write only field.
	Username string `json:"username"`
}

// NewOrgUserRequest instantiates a new OrgUserRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgUserRequest(roles OrgUserRolesRequest, username string) *OrgUserRequest {
	this := OrgUserRequest{}
	this.Roles = roles
	this.Username = username
	return &this
}

// NewOrgUserRequestWithDefaults instantiates a new OrgUserRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgUserRequestWithDefaults() *OrgUserRequest {
	this := OrgUserRequest{}
	return &this
}

// GetRoles returns the Roles field value
func (o *OrgUserRequest) GetRoles() OrgUserRolesRequest {
	if o == nil {
		var ret OrgUserRolesRequest
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *OrgUserRequest) GetRolesOk() (*OrgUserRolesRequest, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *OrgUserRequest) SetRoles(v OrgUserRolesRequest) {
	o.Roles = v
}

// GetTeamIds returns the TeamIds field value if set, zero value otherwise
func (o *OrgUserRequest) GetTeamIds() []string {
	if o == nil || IsNil(o.TeamIds) {
		var ret []string
		return ret
	}
	return *o.TeamIds
}

// GetTeamIdsOk returns a tuple with the TeamIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserRequest) GetTeamIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.TeamIds) {
		return nil, false
	}

	return o.TeamIds, true
}

// HasTeamIds returns a boolean if a field has been set.
func (o *OrgUserRequest) HasTeamIds() bool {
	if o != nil && !IsNil(o.TeamIds) {
		return true
	}

	return false
}

// SetTeamIds gets a reference to the given []string and assigns it to the TeamIds field.
func (o *OrgUserRequest) SetTeamIds(v []string) {
	o.TeamIds = &v
}

// GetUsername returns the Username field value
func (o *OrgUserRequest) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *OrgUserRequest) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *OrgUserRequest) SetUsername(v string) {
	o.Username = v
}
