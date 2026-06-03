// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AddOrRemoveUserFromTeam struct for AddOrRemoveUserFromTeam
type AddOrRemoveUserFromTeam struct {
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.
	// Write only field.
	Id string `json:"id"`
}

// NewAddOrRemoveUserFromTeam instantiates a new AddOrRemoveUserFromTeam object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAddOrRemoveUserFromTeam(id string) *AddOrRemoveUserFromTeam {
	this := AddOrRemoveUserFromTeam{}
	this.Id = id
	return &this
}

// NewAddOrRemoveUserFromTeamWithDefaults instantiates a new AddOrRemoveUserFromTeam object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAddOrRemoveUserFromTeamWithDefaults() *AddOrRemoveUserFromTeam {
	this := AddOrRemoveUserFromTeam{}
	return &this
}

// GetId returns the Id field value
func (o *AddOrRemoveUserFromTeam) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AddOrRemoveUserFromTeam) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AddOrRemoveUserFromTeam) SetId(v string) {
	o.Id = v
}
