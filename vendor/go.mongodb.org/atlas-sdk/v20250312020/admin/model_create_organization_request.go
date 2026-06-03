// Code based on the AtlasAPI V2 OpenAPI file

package admin

// CreateOrganizationRequest struct for CreateOrganizationRequest
type CreateOrganizationRequest struct {
	ApiKey *CreateAtlasOrganizationApiKey `json:"apiKey,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the federation to link the newly created organization to. If specified, the proposed Organization Owner of the new organization must have the Organization Owner role in an organization associated with the federation.
	FederationSettingsId *string `json:"federationSettingsId,omitempty"`
	// Human-readable label that identifies the organization.
	Name string `json:"name"`
	// Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user that you want to assign the Organization Owner role. This user must be a member of the same organization as the calling API key. If you provide `federationSettingsId`,  this user must instead have the Organization Owner role on an organization in the specified federation. This parameter is required only when you authenticate with Programmatic API Keys.
	OrgOwnerId     *string                   `json:"orgOwnerId,omitempty"`
	ServiceAccount *OrgServiceAccountRequest `json:"serviceAccount,omitempty"`
	// Disables automatic alert creation. When set to true, no organization level alerts will be created automatically.
	SkipDefaultAlertsSettings *bool `json:"skipDefaultAlertsSettings,omitempty"`
}

// NewCreateOrganizationRequest instantiates a new CreateOrganizationRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateOrganizationRequest(name string) *CreateOrganizationRequest {
	this := CreateOrganizationRequest{}
	this.Name = name
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// NewCreateOrganizationRequestWithDefaults instantiates a new CreateOrganizationRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateOrganizationRequestWithDefaults() *CreateOrganizationRequest {
	this := CreateOrganizationRequest{}
	var skipDefaultAlertsSettings bool = false
	this.SkipDefaultAlertsSettings = &skipDefaultAlertsSettings
	return &this
}

// GetApiKey returns the ApiKey field value if set, zero value otherwise
func (o *CreateOrganizationRequest) GetApiKey() CreateAtlasOrganizationApiKey {
	if o == nil || IsNil(o.ApiKey) {
		var ret CreateAtlasOrganizationApiKey
		return ret
	}
	return *o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetApiKeyOk() (*CreateAtlasOrganizationApiKey, bool) {
	if o == nil || IsNil(o.ApiKey) {
		return nil, false
	}

	return o.ApiKey, true
}

// HasApiKey returns a boolean if a field has been set.
func (o *CreateOrganizationRequest) HasApiKey() bool {
	if o != nil && !IsNil(o.ApiKey) {
		return true
	}

	return false
}

// SetApiKey gets a reference to the given CreateAtlasOrganizationApiKey and assigns it to the ApiKey field.
func (o *CreateOrganizationRequest) SetApiKey(v CreateAtlasOrganizationApiKey) {
	o.ApiKey = &v
}

// GetFederationSettingsId returns the FederationSettingsId field value if set, zero value otherwise
func (o *CreateOrganizationRequest) GetFederationSettingsId() string {
	if o == nil || IsNil(o.FederationSettingsId) {
		var ret string
		return ret
	}
	return *o.FederationSettingsId
}

// GetFederationSettingsIdOk returns a tuple with the FederationSettingsId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetFederationSettingsIdOk() (*string, bool) {
	if o == nil || IsNil(o.FederationSettingsId) {
		return nil, false
	}

	return o.FederationSettingsId, true
}

// HasFederationSettingsId returns a boolean if a field has been set.
func (o *CreateOrganizationRequest) HasFederationSettingsId() bool {
	if o != nil && !IsNil(o.FederationSettingsId) {
		return true
	}

	return false
}

// SetFederationSettingsId gets a reference to the given string and assigns it to the FederationSettingsId field.
func (o *CreateOrganizationRequest) SetFederationSettingsId(v string) {
	o.FederationSettingsId = &v
}

// GetName returns the Name field value
func (o *CreateOrganizationRequest) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *CreateOrganizationRequest) SetName(v string) {
	o.Name = v
}

// GetOrgOwnerId returns the OrgOwnerId field value if set, zero value otherwise
func (o *CreateOrganizationRequest) GetOrgOwnerId() string {
	if o == nil || IsNil(o.OrgOwnerId) {
		var ret string
		return ret
	}
	return *o.OrgOwnerId
}

// GetOrgOwnerIdOk returns a tuple with the OrgOwnerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetOrgOwnerIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgOwnerId) {
		return nil, false
	}

	return o.OrgOwnerId, true
}

// HasOrgOwnerId returns a boolean if a field has been set.
func (o *CreateOrganizationRequest) HasOrgOwnerId() bool {
	if o != nil && !IsNil(o.OrgOwnerId) {
		return true
	}

	return false
}

// SetOrgOwnerId gets a reference to the given string and assigns it to the OrgOwnerId field.
func (o *CreateOrganizationRequest) SetOrgOwnerId(v string) {
	o.OrgOwnerId = &v
}

// GetServiceAccount returns the ServiceAccount field value if set, zero value otherwise
func (o *CreateOrganizationRequest) GetServiceAccount() OrgServiceAccountRequest {
	if o == nil || IsNil(o.ServiceAccount) {
		var ret OrgServiceAccountRequest
		return ret
	}
	return *o.ServiceAccount
}

// GetServiceAccountOk returns a tuple with the ServiceAccount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetServiceAccountOk() (*OrgServiceAccountRequest, bool) {
	if o == nil || IsNil(o.ServiceAccount) {
		return nil, false
	}

	return o.ServiceAccount, true
}

// HasServiceAccount returns a boolean if a field has been set.
func (o *CreateOrganizationRequest) HasServiceAccount() bool {
	if o != nil && !IsNil(o.ServiceAccount) {
		return true
	}

	return false
}

// SetServiceAccount gets a reference to the given OrgServiceAccountRequest and assigns it to the ServiceAccount field.
func (o *CreateOrganizationRequest) SetServiceAccount(v OrgServiceAccountRequest) {
	o.ServiceAccount = &v
}

// GetSkipDefaultAlertsSettings returns the SkipDefaultAlertsSettings field value if set, zero value otherwise
func (o *CreateOrganizationRequest) GetSkipDefaultAlertsSettings() bool {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		var ret bool
		return ret
	}
	return *o.SkipDefaultAlertsSettings
}

// GetSkipDefaultAlertsSettingsOk returns a tuple with the SkipDefaultAlertsSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateOrganizationRequest) GetSkipDefaultAlertsSettingsOk() (*bool, bool) {
	if o == nil || IsNil(o.SkipDefaultAlertsSettings) {
		return nil, false
	}

	return o.SkipDefaultAlertsSettings, true
}

// HasSkipDefaultAlertsSettings returns a boolean if a field has been set.
func (o *CreateOrganizationRequest) HasSkipDefaultAlertsSettings() bool {
	if o != nil && !IsNil(o.SkipDefaultAlertsSettings) {
		return true
	}

	return false
}

// SetSkipDefaultAlertsSettings gets a reference to the given bool and assigns it to the SkipDefaultAlertsSettings field.
func (o *CreateOrganizationRequest) SetSkipDefaultAlertsSettings(v bool) {
	o.SkipDefaultAlertsSettings = &v
}
