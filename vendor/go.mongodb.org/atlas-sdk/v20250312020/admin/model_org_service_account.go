// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// OrgServiceAccount Organization Service Account that Atlas created for the organization.
type OrgServiceAccount struct {
	// The Client ID of the Service Account.
	ClientId *string `json:"clientId,omitempty"`
	// The date that the Service Account was created on. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// Human readable description for the Service Account.
	Description *string `json:"description,omitempty"`
	// Human-readable name for the Service Account.
	Name *string `json:"name,omitempty"`
	// A list of Organization roles associated with the Service Account.
	Roles *[]string `json:"roles,omitempty"`
	// A list of secrets associated with the specified Service Account.
	Secrets *[]ServiceAccountSecret `json:"secrets,omitempty"`
}

// NewOrgServiceAccount instantiates a new OrgServiceAccount object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgServiceAccount() *OrgServiceAccount {
	this := OrgServiceAccount{}
	return &this
}

// NewOrgServiceAccountWithDefaults instantiates a new OrgServiceAccount object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgServiceAccountWithDefaults() *OrgServiceAccount {
	this := OrgServiceAccount{}
	return &this
}

// GetClientId returns the ClientId field value if set, zero value otherwise
func (o *OrgServiceAccount) GetClientId() string {
	if o == nil || IsNil(o.ClientId) {
		var ret string
		return ret
	}
	return *o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetClientIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClientId) {
		return nil, false
	}

	return o.ClientId, true
}

// HasClientId returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasClientId() bool {
	if o != nil && !IsNil(o.ClientId) {
		return true
	}

	return false
}

// SetClientId gets a reference to the given string and assigns it to the ClientId field.
func (o *OrgServiceAccount) SetClientId(v string) {
	o.ClientId = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *OrgServiceAccount) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *OrgServiceAccount) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *OrgServiceAccount) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *OrgServiceAccount) SetDescription(v string) {
	o.Description = &v
}

// GetName returns the Name field value if set, zero value otherwise
func (o *OrgServiceAccount) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}

	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *OrgServiceAccount) SetName(v string) {
	o.Name = &v
}

// GetRoles returns the Roles field value if set, zero value otherwise
func (o *OrgServiceAccount) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return *o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetRolesOk() (*[]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}

	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *OrgServiceAccount) SetRoles(v []string) {
	o.Roles = &v
}

// GetSecrets returns the Secrets field value if set, zero value otherwise
func (o *OrgServiceAccount) GetSecrets() []ServiceAccountSecret {
	if o == nil || IsNil(o.Secrets) {
		var ret []ServiceAccountSecret
		return ret
	}
	return *o.Secrets
}

// GetSecretsOk returns a tuple with the Secrets field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgServiceAccount) GetSecretsOk() (*[]ServiceAccountSecret, bool) {
	if o == nil || IsNil(o.Secrets) {
		return nil, false
	}

	return o.Secrets, true
}

// HasSecrets returns a boolean if a field has been set.
func (o *OrgServiceAccount) HasSecrets() bool {
	if o != nil && !IsNil(o.Secrets) {
		return true
	}

	return false
}

// SetSecrets gets a reference to the given []ServiceAccountSecret and assigns it to the Secrets field.
func (o *OrgServiceAccount) SetSecrets(v []ServiceAccountSecret) {
	o.Secrets = &v
}
