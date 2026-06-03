// Code based on the AtlasAPI V2 OpenAPI file

package admin

// Team struct for Team
type Team struct {
	// Unique 24-hexadecimal digit string that identifies this team.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Human-readable label that identifies the team.
	Name string `json:"name"`
	// List that contains the MongoDB Cloud users in this team.
	Usernames []string `json:"usernames"`
}

// NewTeam instantiates a new Team object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTeam(name string, usernames []string) *Team {
	this := Team{}
	this.Name = name
	this.Usernames = usernames
	return &this
}

// NewTeamWithDefaults instantiates a new Team object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTeamWithDefaults() *Team {
	this := Team{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *Team) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Team) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Team) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Team) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *Team) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Team) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *Team) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *Team) SetLinks(v []Link) {
	o.Links = &v
}

// GetName returns the Name field value
func (o *Team) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Team) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Team) SetName(v string) {
	o.Name = v
}

// GetUsernames returns the Usernames field value
func (o *Team) GetUsernames() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Usernames
}

// GetUsernamesOk returns a tuple with the Usernames field value
// and a boolean to check if the value has been set.
func (o *Team) GetUsernamesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Usernames, true
}

// SetUsernames sets field value
func (o *Team) SetUsernames(v []string) {
	o.Usernames = v
}
