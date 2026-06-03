// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// GroupUserResponse struct for GroupUserResponse
type GroupUserResponse struct {
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.
	// Read only field.
	Id string `json:"id"`
	// String enum that indicates the user's organization membership status: ACTIVE (member), PENDING (invited), `INVITATION_EXPIRED` (invitation expired), or `INVITATION_REJECTED` (invitation declined).
	// Read only field.
	OrgMembershipStatus string `json:"orgMembershipStatus"`
	// One or more project-level roles assigned to the MongoDB Cloud user.
	// Read only field.
	Roles []string `json:"roles"`
	// Email address that represents the username of the MongoDB Cloud user.
	// Read only field.
	Username string `json:"username"`
	// Date and time when MongoDB Cloud sent the invitation. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC. This field is absent for active users.
	// Read only field.
	InvitationCreatedAt *time.Time `json:"invitationCreatedAt,omitempty"`
	// Date and time when the invitation from MongoDB Cloud expires. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC. This field is absent for active users and null for rejected invitations.
	// Read only field.
	InvitationExpiresAt *time.Time `json:"invitationExpiresAt,omitempty"`
	// Username of the MongoDB Cloud user who sent the invitation to join the organization.
	// Read only field.
	InviterUsername *string `json:"inviterUsername,omitempty"`
	// Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.
	// Read only field.
	Country *string `json:"country,omitempty"`
	// Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// First or given name that belongs to the MongoDB Cloud user.
	// Read only field.
	FirstName *string `json:"firstName,omitempty"`
	// Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastAuth *time.Time `json:"lastAuth,omitempty"`
	// Last name, family name, or surname that belongs to the MongoDB Cloud user.
	// Read only field.
	LastName *string `json:"lastName,omitempty"`
	// Mobile phone number that belongs to the MongoDB Cloud user.
	// Read only field.
	MobileNumber *string `json:"mobileNumber,omitempty"`
}

// NewGroupUserResponse instantiates a new GroupUserResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupUserResponse(id string, orgMembershipStatus string, roles []string, username string) *GroupUserResponse {
	this := GroupUserResponse{}
	this.Id = id
	this.OrgMembershipStatus = orgMembershipStatus
	this.Roles = roles
	this.Username = username
	return &this
}

// NewGroupUserResponseWithDefaults instantiates a new GroupUserResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupUserResponseWithDefaults() *GroupUserResponse {
	this := GroupUserResponse{}
	return &this
}

// GetId returns the Id field value
func (o *GroupUserResponse) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *GroupUserResponse) SetId(v string) {
	o.Id = v
}

// GetOrgMembershipStatus returns the OrgMembershipStatus field value
func (o *GroupUserResponse) GetOrgMembershipStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrgMembershipStatus
}

// GetOrgMembershipStatusOk returns a tuple with the OrgMembershipStatus field value
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetOrgMembershipStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrgMembershipStatus, true
}

// SetOrgMembershipStatus sets field value
func (o *GroupUserResponse) SetOrgMembershipStatus(v string) {
	o.OrgMembershipStatus = v
}

// GetRoles returns the Roles field value
func (o *GroupUserResponse) GetRoles() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetRolesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *GroupUserResponse) SetRoles(v []string) {
	o.Roles = v
}

// GetUsername returns the Username field value
func (o *GroupUserResponse) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *GroupUserResponse) SetUsername(v string) {
	o.Username = v
}

// GetInvitationCreatedAt returns the InvitationCreatedAt field value if set, zero value otherwise
func (o *GroupUserResponse) GetInvitationCreatedAt() time.Time {
	if o == nil || IsNil(o.InvitationCreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.InvitationCreatedAt
}

// GetInvitationCreatedAtOk returns a tuple with the InvitationCreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetInvitationCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.InvitationCreatedAt) {
		return nil, false
	}

	return o.InvitationCreatedAt, true
}

// HasInvitationCreatedAt returns a boolean if a field has been set.
func (o *GroupUserResponse) HasInvitationCreatedAt() bool {
	if o != nil && !IsNil(o.InvitationCreatedAt) {
		return true
	}

	return false
}

// SetInvitationCreatedAt gets a reference to the given time.Time and assigns it to the InvitationCreatedAt field.
func (o *GroupUserResponse) SetInvitationCreatedAt(v time.Time) {
	o.InvitationCreatedAt = &v
}

// GetInvitationExpiresAt returns the InvitationExpiresAt field value if set, zero value otherwise
func (o *GroupUserResponse) GetInvitationExpiresAt() time.Time {
	if o == nil || IsNil(o.InvitationExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.InvitationExpiresAt
}

// GetInvitationExpiresAtOk returns a tuple with the InvitationExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetInvitationExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.InvitationExpiresAt) {
		return nil, false
	}

	return o.InvitationExpiresAt, true
}

// HasInvitationExpiresAt returns a boolean if a field has been set.
func (o *GroupUserResponse) HasInvitationExpiresAt() bool {
	if o != nil && !IsNil(o.InvitationExpiresAt) {
		return true
	}

	return false
}

// SetInvitationExpiresAt gets a reference to the given time.Time and assigns it to the InvitationExpiresAt field.
func (o *GroupUserResponse) SetInvitationExpiresAt(v time.Time) {
	o.InvitationExpiresAt = &v
}

// GetInviterUsername returns the InviterUsername field value if set, zero value otherwise
func (o *GroupUserResponse) GetInviterUsername() string {
	if o == nil || IsNil(o.InviterUsername) {
		var ret string
		return ret
	}
	return *o.InviterUsername
}

// GetInviterUsernameOk returns a tuple with the InviterUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetInviterUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.InviterUsername) {
		return nil, false
	}

	return o.InviterUsername, true
}

// HasInviterUsername returns a boolean if a field has been set.
func (o *GroupUserResponse) HasInviterUsername() bool {
	if o != nil && !IsNil(o.InviterUsername) {
		return true
	}

	return false
}

// SetInviterUsername gets a reference to the given string and assigns it to the InviterUsername field.
func (o *GroupUserResponse) SetInviterUsername(v string) {
	o.InviterUsername = &v
}

// GetCountry returns the Country field value if set, zero value otherwise
func (o *GroupUserResponse) GetCountry() string {
	if o == nil || IsNil(o.Country) {
		var ret string
		return ret
	}
	return *o.Country
}

// GetCountryOk returns a tuple with the Country field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetCountryOk() (*string, bool) {
	if o == nil || IsNil(o.Country) {
		return nil, false
	}

	return o.Country, true
}

// HasCountry returns a boolean if a field has been set.
func (o *GroupUserResponse) HasCountry() bool {
	if o != nil && !IsNil(o.Country) {
		return true
	}

	return false
}

// SetCountry gets a reference to the given string and assigns it to the Country field.
func (o *GroupUserResponse) SetCountry(v string) {
	o.Country = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *GroupUserResponse) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *GroupUserResponse) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *GroupUserResponse) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetFirstName returns the FirstName field value if set, zero value otherwise
func (o *GroupUserResponse) GetFirstName() string {
	if o == nil || IsNil(o.FirstName) {
		var ret string
		return ret
	}
	return *o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetFirstNameOk() (*string, bool) {
	if o == nil || IsNil(o.FirstName) {
		return nil, false
	}

	return o.FirstName, true
}

// HasFirstName returns a boolean if a field has been set.
func (o *GroupUserResponse) HasFirstName() bool {
	if o != nil && !IsNil(o.FirstName) {
		return true
	}

	return false
}

// SetFirstName gets a reference to the given string and assigns it to the FirstName field.
func (o *GroupUserResponse) SetFirstName(v string) {
	o.FirstName = &v
}

// GetLastAuth returns the LastAuth field value if set, zero value otherwise
func (o *GroupUserResponse) GetLastAuth() time.Time {
	if o == nil || IsNil(o.LastAuth) {
		var ret time.Time
		return ret
	}
	return *o.LastAuth
}

// GetLastAuthOk returns a tuple with the LastAuth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetLastAuthOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastAuth) {
		return nil, false
	}

	return o.LastAuth, true
}

// HasLastAuth returns a boolean if a field has been set.
func (o *GroupUserResponse) HasLastAuth() bool {
	if o != nil && !IsNil(o.LastAuth) {
		return true
	}

	return false
}

// SetLastAuth gets a reference to the given time.Time and assigns it to the LastAuth field.
func (o *GroupUserResponse) SetLastAuth(v time.Time) {
	o.LastAuth = &v
}

// GetLastName returns the LastName field value if set, zero value otherwise
func (o *GroupUserResponse) GetLastName() string {
	if o == nil || IsNil(o.LastName) {
		var ret string
		return ret
	}
	return *o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetLastNameOk() (*string, bool) {
	if o == nil || IsNil(o.LastName) {
		return nil, false
	}

	return o.LastName, true
}

// HasLastName returns a boolean if a field has been set.
func (o *GroupUserResponse) HasLastName() bool {
	if o != nil && !IsNil(o.LastName) {
		return true
	}

	return false
}

// SetLastName gets a reference to the given string and assigns it to the LastName field.
func (o *GroupUserResponse) SetLastName(v string) {
	o.LastName = &v
}

// GetMobileNumber returns the MobileNumber field value if set, zero value otherwise
func (o *GroupUserResponse) GetMobileNumber() string {
	if o == nil || IsNil(o.MobileNumber) {
		var ret string
		return ret
	}
	return *o.MobileNumber
}

// GetMobileNumberOk returns a tuple with the MobileNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GroupUserResponse) GetMobileNumberOk() (*string, bool) {
	if o == nil || IsNil(o.MobileNumber) {
		return nil, false
	}

	return o.MobileNumber, true
}

// HasMobileNumber returns a boolean if a field has been set.
func (o *GroupUserResponse) HasMobileNumber() bool {
	if o != nil && !IsNil(o.MobileNumber) {
		return true
	}

	return false
}

// SetMobileNumber gets a reference to the given string and assigns it to the MobileNumber field.
func (o *GroupUserResponse) SetMobileNumber(v string) {
	o.MobileNumber = &v
}
