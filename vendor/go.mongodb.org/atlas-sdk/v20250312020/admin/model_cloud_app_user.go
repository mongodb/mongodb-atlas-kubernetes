// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// CloudAppUser struct for CloudAppUser
type CloudAppUser struct {
	// Two alphabet characters that identifies MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.
	Country string `json:"country"`
	// Date and time when the current account is created. This value is in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Email address that belongs to the MongoDB Cloud user.
	// Read only field.
	// Deprecated
	EmailAddress string `json:"emailAddress"`
	// First or given name that belongs to the MongoDB Cloud user.
	FirstName string `json:"firstName"`
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.
	// Read only field.
	LastAuth *time.Time `json:"lastAuth,omitempty"`
	// Last name, family name, or surname that belongs to the MongoDB Cloud user.
	LastName string `json:"lastName"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Mobile phone number that belongs to the MongoDB Cloud user.
	MobileNumber string `json:"mobileNumber"`
	// Password applied with the username to log in to MongoDB Cloud. MongoDB Cloud does not return this parameter except in response to creating a new MongoDB Cloud user. Only the MongoDB Cloud user can update their password after it has been set from the MongoDB Cloud console.
	Password string `json:"password"`
	// List of objects that display the MongoDB Cloud user's roles and the corresponding organization or project to which that role applies. A role can apply to one organization or one project but not both.
	Roles []CloudAccessRoleAssignment `json:"roles"`
	// List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.
	// Read only field.
	TeamIds *[]string `json:"teamIds,omitempty"`
	// Email address that represents the username of the MongoDB Cloud user.
	Username string `json:"username"`
}

// NewCloudAppUser instantiates a new CloudAppUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCloudAppUser(country string, emailAddress string, firstName string, lastName string, mobileNumber string, password string, roles []CloudAccessRoleAssignment, username string) *CloudAppUser {
	this := CloudAppUser{}
	this.Country = country
	this.EmailAddress = emailAddress
	this.FirstName = firstName
	this.LastName = lastName
	this.MobileNumber = mobileNumber
	this.Password = password
	this.Roles = roles
	this.Username = username
	return &this
}

// NewCloudAppUserWithDefaults instantiates a new CloudAppUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCloudAppUserWithDefaults() *CloudAppUser {
	this := CloudAppUser{}
	return &this
}

// GetCountry returns the Country field value
func (o *CloudAppUser) GetCountry() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Country
}

// GetCountryOk returns a tuple with the Country field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetCountryOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Country, true
}

// SetCountry sets field value
func (o *CloudAppUser) SetCountry(v string) {
	o.Country = v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *CloudAppUser) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *CloudAppUser) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *CloudAppUser) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetEmailAddress returns the EmailAddress field value
// Deprecated
func (o *CloudAppUser) GetEmailAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EmailAddress
}

// GetEmailAddressOk returns a tuple with the EmailAddress field value
// and a boolean to check if the value has been set.
// Deprecated
func (o *CloudAppUser) GetEmailAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EmailAddress, true
}

// SetEmailAddress sets field value
// Deprecated
func (o *CloudAppUser) SetEmailAddress(v string) {
	o.EmailAddress = v
}

// GetFirstName returns the FirstName field value
func (o *CloudAppUser) GetFirstName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetFirstNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FirstName, true
}

// SetFirstName sets field value
func (o *CloudAppUser) SetFirstName(v string) {
	o.FirstName = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *CloudAppUser) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CloudAppUser) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CloudAppUser) SetId(v string) {
	o.Id = &v
}

// GetLastAuth returns the LastAuth field value if set, zero value otherwise
func (o *CloudAppUser) GetLastAuth() time.Time {
	if o == nil || IsNil(o.LastAuth) {
		var ret time.Time
		return ret
	}
	return *o.LastAuth
}

// GetLastAuthOk returns a tuple with the LastAuth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetLastAuthOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastAuth) {
		return nil, false
	}

	return o.LastAuth, true
}

// HasLastAuth returns a boolean if a field has been set.
func (o *CloudAppUser) HasLastAuth() bool {
	if o != nil && !IsNil(o.LastAuth) {
		return true
	}

	return false
}

// SetLastAuth gets a reference to the given time.Time and assigns it to the LastAuth field.
func (o *CloudAppUser) SetLastAuth(v time.Time) {
	o.LastAuth = &v
}

// GetLastName returns the LastName field value
func (o *CloudAppUser) GetLastName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetLastNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastName, true
}

// SetLastName sets field value
func (o *CloudAppUser) SetLastName(v string) {
	o.LastName = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *CloudAppUser) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *CloudAppUser) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *CloudAppUser) SetLinks(v []Link) {
	o.Links = &v
}

// GetMobileNumber returns the MobileNumber field value
func (o *CloudAppUser) GetMobileNumber() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MobileNumber
}

// GetMobileNumberOk returns a tuple with the MobileNumber field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetMobileNumberOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MobileNumber, true
}

// SetMobileNumber sets field value
func (o *CloudAppUser) SetMobileNumber(v string) {
	o.MobileNumber = v
}

// GetPassword returns the Password field value
func (o *CloudAppUser) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *CloudAppUser) SetPassword(v string) {
	o.Password = v
}

// GetRoles returns the Roles field value
func (o *CloudAppUser) GetRoles() []CloudAccessRoleAssignment {
	if o == nil {
		var ret []CloudAccessRoleAssignment
		return ret
	}

	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetRolesOk() (*[]CloudAccessRoleAssignment, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Roles, true
}

// SetRoles sets field value
func (o *CloudAppUser) SetRoles(v []CloudAccessRoleAssignment) {
	o.Roles = v
}

// GetTeamIds returns the TeamIds field value if set, zero value otherwise
func (o *CloudAppUser) GetTeamIds() []string {
	if o == nil || IsNil(o.TeamIds) {
		var ret []string
		return ret
	}
	return *o.TeamIds
}

// GetTeamIdsOk returns a tuple with the TeamIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetTeamIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.TeamIds) {
		return nil, false
	}

	return o.TeamIds, true
}

// HasTeamIds returns a boolean if a field has been set.
func (o *CloudAppUser) HasTeamIds() bool {
	if o != nil && !IsNil(o.TeamIds) {
		return true
	}

	return false
}

// SetTeamIds gets a reference to the given []string and assigns it to the TeamIds field.
func (o *CloudAppUser) SetTeamIds(v []string) {
	o.TeamIds = &v
}

// GetUsername returns the Username field value
func (o *CloudAppUser) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *CloudAppUser) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *CloudAppUser) SetUsername(v string) {
	o.Username = v
}
