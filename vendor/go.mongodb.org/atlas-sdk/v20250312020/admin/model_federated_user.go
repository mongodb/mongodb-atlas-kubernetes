// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FederatedUser MongoDB Cloud user linked to this federated authentication.
type FederatedUser struct {
	// Email address of the MongoDB Cloud user linked to the federated organization.
	EmailAddress string `json:"emailAddress"`
	// Unique 24-hexadecimal digit string that identifies the federation to which this MongoDB Cloud user belongs.
	FederationSettingsId string `json:"federationSettingsId"`
	// First or given name that belongs to the MongoDB Cloud user.
	FirstName string `json:"firstName"`
	// Last name, family name, or surname that belongs to the MongoDB Cloud user.
	LastName string `json:"lastName"`
	// Unique 24-hexadecimal digit string that identifies this user.
	// Read only field.
	UserId *string `json:"userId,omitempty"`
}

// NewFederatedUser instantiates a new FederatedUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFederatedUser(emailAddress string, federationSettingsId string, firstName string, lastName string) *FederatedUser {
	this := FederatedUser{}
	this.EmailAddress = emailAddress
	this.FederationSettingsId = federationSettingsId
	this.FirstName = firstName
	this.LastName = lastName
	return &this
}

// NewFederatedUserWithDefaults instantiates a new FederatedUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFederatedUserWithDefaults() *FederatedUser {
	this := FederatedUser{}
	return &this
}

// GetEmailAddress returns the EmailAddress field value
func (o *FederatedUser) GetEmailAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EmailAddress
}

// GetEmailAddressOk returns a tuple with the EmailAddress field value
// and a boolean to check if the value has been set.
func (o *FederatedUser) GetEmailAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EmailAddress, true
}

// SetEmailAddress sets field value
func (o *FederatedUser) SetEmailAddress(v string) {
	o.EmailAddress = v
}

// GetFederationSettingsId returns the FederationSettingsId field value
func (o *FederatedUser) GetFederationSettingsId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FederationSettingsId
}

// GetFederationSettingsIdOk returns a tuple with the FederationSettingsId field value
// and a boolean to check if the value has been set.
func (o *FederatedUser) GetFederationSettingsIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FederationSettingsId, true
}

// SetFederationSettingsId sets field value
func (o *FederatedUser) SetFederationSettingsId(v string) {
	o.FederationSettingsId = v
}

// GetFirstName returns the FirstName field value
func (o *FederatedUser) GetFirstName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FirstName
}

// GetFirstNameOk returns a tuple with the FirstName field value
// and a boolean to check if the value has been set.
func (o *FederatedUser) GetFirstNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FirstName, true
}

// SetFirstName sets field value
func (o *FederatedUser) SetFirstName(v string) {
	o.FirstName = v
}

// GetLastName returns the LastName field value
func (o *FederatedUser) GetLastName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LastName
}

// GetLastNameOk returns a tuple with the LastName field value
// and a boolean to check if the value has been set.
func (o *FederatedUser) GetLastNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastName, true
}

// SetLastName sets field value
func (o *FederatedUser) SetLastName(v string) {
	o.LastName = v
}

// GetUserId returns the UserId field value if set, zero value otherwise
func (o *FederatedUser) GetUserId() string {
	if o == nil || IsNil(o.UserId) {
		var ret string
		return ret
	}
	return *o.UserId
}

// GetUserIdOk returns a tuple with the UserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederatedUser) GetUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.UserId) {
		return nil, false
	}

	return o.UserId, true
}

// HasUserId returns a boolean if a field has been set.
func (o *FederatedUser) HasUserId() bool {
	if o != nil && !IsNil(o.UserId) {
		return true
	}

	return false
}

// SetUserId gets a reference to the given string and assigns it to the UserId field.
func (o *FederatedUser) SetUserId(v string) {
	o.UserId = &v
}
