// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// GroupInvitation struct for GroupInvitation
type GroupInvitation struct {
	// Date and time when MongoDB Cloud sent the invitation. This parameter expresses its value in ISO 8601 format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Date and time when MongoDB Cloud expires the invitation. This parameter expresses its value in ISO 8601 format in UTC.
	// Read only field.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Unique 24-hexadecimal character string that identifies the project.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Human-readable label that identifies the project to which you invited the MongoDB Cloud user.
	// Read only field.
	GroupName *string `json:"groupName,omitempty"`
	// Unique 24-hexadecimal character string that identifies the invitation.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Email address of the MongoDB Cloud user who sent the invitation.
	// Read only field.
	InviterUsername *string `json:"inviterUsername,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// One or more organization or project level roles to assign to the MongoDB Cloud user.
	Roles *[]string `json:"roles,omitempty"`
	// Email address of the MongoDB Cloud user invited to join the project.
	// Read only field.
	Username *string `json:"username,omitempty"`
}

// NewGroupInvitation instantiates a new GroupInvitation object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupInvitation() *GroupInvitation {
	this := GroupInvitation{}
	return &this
}

// NewGroupInvitationWithDefaults instantiates a new GroupInvitation object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupInvitationWithDefaults() *GroupInvitation {
	this := GroupInvitation{}
	return &this
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *GroupInvitation) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *GroupInvitation) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *GroupInvitation) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise
func (o *GroupInvitation) GetExpiresAt() time.Time {
	if o == nil || IsNil(o.ExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpiresAt) {
		return nil, false
	}

	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *GroupInvitation) HasExpiresAt() bool {
	if o != nil && !IsNil(o.ExpiresAt) {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *GroupInvitation) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *GroupInvitation) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *GroupInvitation) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *GroupInvitation) SetGroupId(v string) {
	o.GroupId = &v
}

// GetGroupName returns the GroupName field value if set, zero value otherwise
func (o *GroupInvitation) GetGroupName() string {
	if o == nil || IsNil(o.GroupName) {
		var ret string
		return ret
	}
	return *o.GroupName
}

// GetGroupNameOk returns a tuple with the GroupName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetGroupNameOk() (*string, bool) {
	if o == nil || IsNil(o.GroupName) {
		return nil, false
	}

	return o.GroupName, true
}

// HasGroupName returns a boolean if a field has been set.
func (o *GroupInvitation) HasGroupName() bool {
	if o != nil && !IsNil(o.GroupName) {
		return true
	}

	return false
}

// SetGroupName gets a reference to the given string and assigns it to the GroupName field.
func (o *GroupInvitation) SetGroupName(v string) {
	o.GroupName = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *GroupInvitation) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *GroupInvitation) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *GroupInvitation) SetId(v string) {
	o.Id = &v
}

// GetInviterUsername returns the InviterUsername field value if set, zero value otherwise
func (o *GroupInvitation) GetInviterUsername() string {
	if o == nil || IsNil(o.InviterUsername) {
		var ret string
		return ret
	}
	return *o.InviterUsername
}

// GetInviterUsernameOk returns a tuple with the InviterUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetInviterUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.InviterUsername) {
		return nil, false
	}

	return o.InviterUsername, true
}

// HasInviterUsername returns a boolean if a field has been set.
func (o *GroupInvitation) HasInviterUsername() bool {
	if o != nil && !IsNil(o.InviterUsername) {
		return true
	}

	return false
}

// SetInviterUsername gets a reference to the given string and assigns it to the InviterUsername field.
func (o *GroupInvitation) SetInviterUsername(v string) {
	o.InviterUsername = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *GroupInvitation) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *GroupInvitation) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *GroupInvitation) SetLinks(v []Link) {
	o.Links = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *GroupInvitation) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *GroupInvitation) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *GroupInvitation) SetRoles(v []string) {
	o.Roles = &v
}

// GetUsername returns the Username field value if set, zero value otherwise
func (o *GroupInvitation) GetUsername() string {
	if o == nil || IsNil(o.Username) {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupInvitation) GetUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.Username) {
		return nil, false
	}

	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *GroupInvitation) HasUsername() bool {
	if o != nil && !IsNil(o.Username) {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *GroupInvitation) SetUsername(v string) {
	o.Username = &v
}
