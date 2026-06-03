// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgUserUpdateRequest struct for OrgUserUpdateRequest
type OrgUserUpdateRequest struct {
	Roles *OrgUserRolesRequest `json:"roles,omitempty"`
	// List of unique 24-hexadecimal digit strings that identifies the teams to assign the MongoDB Cloud user.
	// Write only field.
	TeamIds *[]string `json:"teamIds,omitempty"`
}

// NewOrgUserUpdateRequest instantiates a new OrgUserUpdateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgUserUpdateRequest() *OrgUserUpdateRequest {
	this := OrgUserUpdateRequest{}
	return &this
}

// NewOrgUserUpdateRequestWithDefaults instantiates a new OrgUserUpdateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgUserUpdateRequestWithDefaults() *OrgUserUpdateRequest {
	this := OrgUserUpdateRequest{}
	return &this
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *OrgUserUpdateRequest) GetRoles() OrgUserRolesRequest {
	if o == nil || IsNil(o.Roles) {
		var ret OrgUserRolesRequest
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserUpdateRequest) GetRolesOk() (*OrgUserRolesRequest, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *OrgUserUpdateRequest) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given OrgUserRolesRequest and assigns it to the Roles field.
func (o *OrgUserUpdateRequest) SetRoles(v OrgUserRolesRequest) {
	o.Roles = &v
}

// GetTeamIds returns the TeamIds field value if set, zero value otherwise
func (o *OrgUserUpdateRequest) GetTeamIds() []string {
	if o == nil || IsNil(o.TeamIds) {
		var ret []string
		return ret
	}
	return *o.TeamIds
}

// GetTeamIdsOk returns a tuple with the TeamIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgUserUpdateRequest) GetTeamIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.TeamIds) {
		return nil, false
	}

	return o.TeamIds, true
}

// HasTeamIds returns a boolean if a field has been set.
func (o *OrgUserUpdateRequest) HasTeamIds() bool {
	if o != nil && !IsNil(o.TeamIds) {
		return true
	}

	return false
}

// SetTeamIds gets a reference to the given []string and assigns it to the TeamIds field.
func (o *OrgUserUpdateRequest) SetTeamIds(v []string) {
	o.TeamIds = &v
}
