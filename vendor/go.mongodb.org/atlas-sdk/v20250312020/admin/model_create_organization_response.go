// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CreateOrganizationResponse struct for CreateOrganizationResponse
type CreateOrganizationResponse struct {
	ApiKey *ApiKeyUserDetails `json:"apiKey,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the federation that you linked the newly created organization to.
	// Read only field.
	FederationSettingsId *string `json:"federationSettingsId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user that you assigned the Organization Owner role in the new organization.
	// Read only field.
	OrgOwnerId     *string            `json:"orgOwnerId,omitempty"`
	Organization   *AtlasOrganization `json:"organization,omitempty"`
	ServiceAccount *OrgServiceAccount `json:"serviceAccount,omitempty"`
	// Disables automatic alert creation. When set to true, no organization level alerts will be created automatically.
	SkipDefaultAlertsSettings *bool `json:"skipDefaultAlertsSettings,omitempty"`
}

// NewCreateOrganizationResponse instantiates a new CreateOrganizationResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateOrganizationResponse() *CreateOrganizationResponse {
	this := CreateOrganizationResponse{}
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// NewCreateOrganizationResponseWithDefaults instantiates a new CreateOrganizationResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateOrganizationResponseWithDefaults() *CreateOrganizationResponse {
	this := CreateOrganizationResponse{}
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// GetApiKey returns the ApiKey field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetApiKey() ApiKeyUserDetails {
	if o == nil || IsNil(o.ApiKey) {
		var ret ApiKeyUserDetails
		return ret
	}
	return *o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetApiKeyOk() (*ApiKeyUserDetails, bool) {
	if o == nil || IsNil(o.ApiKey) {
		return nil, false
	}

	return o.ApiKey, true
}

// HasApiKey returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasApiKey() bool {
	if o != nil && !IsNil(o.ApiKey) {
		return true
	}

	return false
}

// SetApiKey gets a reference to the given ApiKeyUserDetails and assigns it to the ApiKey field.
func (o *CreateOrganizationResponse) SetApiKey(v ApiKeyUserDetails) {
	o.ApiKey = &v
}

// GetFederationSettingsId returns the FederationSettingsId field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetFederationSettingsId() string {
	if o == nil || IsNil(o.FederationSettingsId) {
		var ret string
		return ret
	}
	return *o.FederationSettingsId
}

// GetFederationSettingsIdOk returns a tuple with the FederationSettingsId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetFederationSettingsIdOk() (*string, bool) {
	if o == nil || IsNil(o.FederationSettingsId) {
		return nil, false
	}

	return o.FederationSettingsId, true
}

// HasFederationSettingsId returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasFederationSettingsId() bool {
	if o != nil && !IsNil(o.FederationSettingsId) {
		return true
	}

	return false
}

// SetFederationSettingsId gets a reference to the given string and assigns it to the FederationSettingsId field.
func (o *CreateOrganizationResponse) SetFederationSettingsId(v string) {
	o.FederationSettingsId = &v
}

// GetOrgOwnerId returns the OrgOwnerId field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetOrgOwnerId() string {
	if o == nil || IsNil(o.OrgOwnerId) {
		var ret string
		return ret
	}
	return *o.OrgOwnerId
}

// GetOrgOwnerIdOk returns a tuple with the OrgOwnerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetOrgOwnerIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgOwnerId) {
		return nil, false
	}

	return o.OrgOwnerId, true
}

// HasOrgOwnerId returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasOrgOwnerId() bool {
	if o != nil && !IsNil(o.OrgOwnerId) {
		return true
	}

	return false
}

// SetOrgOwnerId gets a reference to the given string and assigns it to the OrgOwnerId field.
func (o *CreateOrganizationResponse) SetOrgOwnerId(v string) {
	o.OrgOwnerId = &v
}

// GetOrganization returns the Organization field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetOrganization() AtlasOrganization {
	if o == nil || IsNil(o.Organization) {
		var ret AtlasOrganization
		return ret
	}
	return *o.Organization
}

// GetOrganizationOk returns a tuple with the Organization field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetOrganizationOk() (*AtlasOrganization, bool) {
	if o == nil || IsNil(o.Organization) {
		return nil, false
	}

	return o.Organization, true
}

// HasOrganization returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasOrganization() bool {
	if o != nil && !IsNil(o.Organization) {
		return true
	}

	return false
}

// SetOrganization gets a reference to the given AtlasOrganization and assigns it to the Organization field.
func (o *CreateOrganizationResponse) SetOrganization(v AtlasOrganization) {
	o.Organization = &v
}

// GetServiceAccount returns the ServiceAccount field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetServiceAccount() OrgServiceAccount {
	if o == nil || IsNil(o.ServiceAccount) {
		var ret OrgServiceAccount
		return ret
	}
	return *o.ServiceAccount
}

// GetServiceAccountOk returns a tuple with the ServiceAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetServiceAccountOk() (*OrgServiceAccount, bool) {
	if o == nil || IsNil(o.ServiceAccount) {
		return nil, false
	}

	return o.ServiceAccount, true
}

// HasServiceAccount returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasServiceAccount() bool {
	if o != nil && !IsNil(o.ServiceAccount) {
		return true
	}

	return false
}

// SetServiceAccount gets a reference to the given OrgServiceAccount and assigns it to the ServiceAccount field.
func (o *CreateOrganizationResponse) SetServiceAccount(v OrgServiceAccount) {
	o.ServiceAccount = &v
}

// GetSkipDefaultAlertsSettings returns the SkipDefaultAlertsSettings field value if set, zero value otherwise
func (o *CreateOrganizationResponse) GetSkipDefaultAlertsSettings() bool {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		var ret bool
		return ret
	}
	return *o.SkipDefaultAlertsSettings
}

// GetSkipDefaultAlertsSettingsOk returns a tuple with the SkipDefaultAlertsSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationResponse) GetSkipDefaultAlertsSettingsOk() (*bool, bool) {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		return nil, false
	}

	return o.SkipDefaultAlertsSettings, true
}

// HasSkipDefaultAlertsSettings returns a boolean if a field has been set.
func (o *CreateOrganizationResponse) HasSkipDefaultAlertsSettings() bool {
	if o != nil && !IsNil(o.SkipDefaultAlertsSettings) {
		return true
	}

	return false
}

// SetSkipDefaultAlertsSettings gets a reference to the given bool and assigns it to the SkipDefaultAlertsSettings field.
func (o *CreateOrganizationResponse) SetSkipDefaultAlertsSettings(v bool) {
	o.SkipDefaultAlertsSettings = &v
}
