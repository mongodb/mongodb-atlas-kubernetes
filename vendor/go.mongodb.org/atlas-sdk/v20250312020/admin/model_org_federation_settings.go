// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrgFederationSettings Details that define how to connect one MongoDB Cloud organization to one federated authentication service.
type OrgFederationSettings struct {
	// List of domains associated with the organization's identity provider.
	FederatedDomains []string `json:"federatedDomains"`
	// Flag that indicates whether this organization has role mappings configured.
	HasRoleMappings *bool `json:"hasRoleMappings,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this federation.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Legacy 20-hexadecimal digit string that identifies the identity provider connected to this organization.
	IdentityProviderId *string `json:"identityProviderId,omitempty"`
	// String enum that indicates whether the identity provider is active.
	IdentityProviderStatus *string `json:"identityProviderStatus,omitempty"`
}

// NewOrgFederationSettings instantiates a new OrgFederationSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrgFederationSettings(federatedDomains []string) *OrgFederationSettings {
	this := OrgFederationSettings{}
	this.FederatedDomains = federatedDomains
	return &this
}

// NewOrgFederationSettingsWithDefaults instantiates a new OrgFederationSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrgFederationSettingsWithDefaults() *OrgFederationSettings {
	this := OrgFederationSettings{}
	return &this
}

// GetFederatedDomains returns the FederatedDomains field value
func (o *OrgFederationSettings) GetFederatedDomains() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.FederatedDomains
}

// GetFederatedDomainsOk returns a tuple with the FederatedDomains field value
// and a boolean to check if the value has been set.
func (o *OrgFederationSettings) GetFederatedDomainsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FederatedDomains, true
}

// SetFederatedDomains sets field value
func (o *OrgFederationSettings) SetFederatedDomains(v []string) {
	o.FederatedDomains = v
}

// GetHasRoleMappings returns the HasRoleMappings field value if set, zero value otherwise
func (o *OrgFederationSettings) GetHasRoleMappings() bool {
	if o == nil || IsNil(o.HasRoleMappings) {
		var ret bool
		return ret
	}
	return *o.HasRoleMappings
}

// GetHasRoleMappingsOk returns a tuple with the HasRoleMappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgFederationSettings) GetHasRoleMappingsOk() (*bool, bool) {
	if o == nil || IsNil(o.HasRoleMappings) {
		return nil, false
	}

	return o.HasRoleMappings, true
}

// HasHasRoleMappings returns a boolean if a field has been set.
func (o *OrgFederationSettings) HasHasRoleMappings() bool {
	if o != nil && !IsNil(o.HasRoleMappings) {
		return true
	}

	return false
}

// SetHasRoleMappings gets a reference to the given bool and assigns it to the HasRoleMappings field.
func (o *OrgFederationSettings) SetHasRoleMappings(v bool) {
	o.HasRoleMappings = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *OrgFederationSettings) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgFederationSettings) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *OrgFederationSettings) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *OrgFederationSettings) SetId(v string) {
	o.Id = &v
}

// GetIdentityProviderId returns the IdentityProviderId field value if set, zero value otherwise
func (o *OrgFederationSettings) GetIdentityProviderId() string {
	if o == nil || IsNil(o.IdentityProviderId) {
		var ret string
		return ret
	}
	return *o.IdentityProviderId
}

// GetIdentityProviderIdOk returns a tuple with the IdentityProviderId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgFederationSettings) GetIdentityProviderIdOk() (*string, bool) {
	if o == nil || IsNil(o.IdentityProviderId) {
		return nil, false
	}

	return o.IdentityProviderId, true
}

// HasIdentityProviderId returns a boolean if a field has been set.
func (o *OrgFederationSettings) HasIdentityProviderId() bool {
	if o != nil && !IsNil(o.IdentityProviderId) {
		return true
	}

	return false
}

// SetIdentityProviderId gets a reference to the given string and assigns it to the IdentityProviderId field.
func (o *OrgFederationSettings) SetIdentityProviderId(v string) {
	o.IdentityProviderId = &v
}

// GetIdentityProviderStatus returns the IdentityProviderStatus field value if set, zero value otherwise
func (o *OrgFederationSettings) GetIdentityProviderStatus() string {
	if o == nil || IsNil(o.IdentityProviderStatus) {
		var ret string
		return ret
	}
	return *o.IdentityProviderStatus
}

// GetIdentityProviderStatusOk returns a tuple with the IdentityProviderStatus field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrgFederationSettings) GetIdentityProviderStatusOk() (*string, bool) {
	if o == nil || IsNil(o.IdentityProviderStatus) {
		return nil, false
	}

	return o.IdentityProviderStatus, true
}

// HasIdentityProviderStatus returns a boolean if a field has been set.
func (o *OrgFederationSettings) HasIdentityProviderStatus() bool {
	if o != nil && !IsNil(o.IdentityProviderStatus) {
		return true
	}

	return false
}

// SetIdentityProviderStatus gets a reference to the given string and assigns it to the IdentityProviderStatus field.
func (o *OrgFederationSettings) SetIdentityProviderStatus(v string) {
	o.IdentityProviderStatus = &v
}
