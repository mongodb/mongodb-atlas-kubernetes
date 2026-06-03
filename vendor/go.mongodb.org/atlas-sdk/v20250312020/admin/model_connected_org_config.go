// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ConnectedOrgConfig struct for ConnectedOrgConfig
type ConnectedOrgConfig struct {
	// The collection of unique ids representing the identity providers that can be used for data access in this organization.
	DataAccessIdentityProviderIds *[]string `json:"dataAccessIdentityProviderIds,omitempty"`
	// Approved domains that restrict users who can join the organization based on their email address.
	DomainAllowList *[]string `json:"domainAllowList,omitempty"`
	// Value that indicates whether domain restriction is enabled for this connected organization.
	DomainRestrictionEnabled bool `json:"domainRestrictionEnabled"`
	// Legacy 20-hexadecimal digit string that identifies the UI access identity provider that this connected organization configuration is associated with. This id can be found within the Federation Management Console > Identity Providers tab by clicking the info icon in the IdP ID row of a configured identity provider.
	IdentityProviderId *string `json:"identityProviderId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the connected organization configuration.
	// Read only field.
	OrgId string `json:"orgId"`
	// Atlas roles that are granted to a user in this organization after authenticating. Roles are a human-readable label that identifies the collection of privileges that MongoDB Cloud grants a specific MongoDB Cloud user. These roles can only be organization specific roles.
	PostAuthRoleGrants *[]string `json:"postAuthRoleGrants,omitempty"`
	// Role mappings that are configured in this organization.
	RoleMappings *[]AuthFederationRoleMapping `json:"roleMappings,omitempty"`
	// List that contains the users who have an email address that doesn't match any domain on the allowed list.
	UserConflicts *[]FederatedUser `json:"userConflicts,omitempty"`
}

// NewConnectedOrgConfig instantiates a new ConnectedOrgConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectedOrgConfig(domainRestrictionEnabled bool, orgId string) *ConnectedOrgConfig {
	this := ConnectedOrgConfig{}
	this.DomainRestrictionEnabled = domainRestrictionEnabled
	this.OrgId = orgId
	return &this
}

// NewConnectedOrgConfigWithDefaults instantiates a new ConnectedOrgConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectedOrgConfigWithDefaults() *ConnectedOrgConfig {
	this := ConnectedOrgConfig{}
	return &this
}

// GetDataAccessIdentityProviderIds returns the DataAccessIdentityProviderIds field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetDataAccessIdentityProviderIds() []string {
	if o == nil || IsNil(o.DataAccessIdentityProviderIds) {
		var ret []string
		return ret
	}
	return *o.DataAccessIdentityProviderIds
}

// GetDataAccessIdentityProviderIdsOk returns a tuple with the DataAccessIdentityProviderIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetDataAccessIdentityProviderIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.DataAccessIdentityProviderIds) {
		return nil, false
	}

	return o.DataAccessIdentityProviderIds, true
}

// HasDataAccessIdentityProviderIds returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasDataAccessIdentityProviderIds() bool {
	if o != nil && !IsNil(o.DataAccessIdentityProviderIds) {
		return true
	}

	return false
}

// SetDataAccessIdentityProviderIds gets a reference to the given []string and assigns it to the DataAccessIdentityProviderIds field.
func (o *ConnectedOrgConfig) SetDataAccessIdentityProviderIds(v []string) {
	o.DataAccessIdentityProviderIds = &v
}

// GetDomainAllowList returns the DomainAllowList field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetDomainAllowList() []string {
	if o == nil || IsNil(o.DomainAllowList) {
		var ret []string
		return ret
	}
	return *o.DomainAllowList
}

// GetDomainAllowListOk returns a tuple with the DomainAllowList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetDomainAllowListOk() (*[]string, bool) {
	if o == nil || IsNil(o.DomainAllowList) {
		return nil, false
	}

	return o.DomainAllowList, true
}

// HasDomainAllowList returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasDomainAllowList() bool {
	if o != nil && !IsNil(o.DomainAllowList) {
		return true
	}

	return false
}

// SetDomainAllowList gets a reference to the given []string and assigns it to the DomainAllowList field.
func (o *ConnectedOrgConfig) SetDomainAllowList(v []string) {
	o.DomainAllowList = &v
}

// GetDomainRestrictionEnabled returns the DomainRestrictionEnabled field value
func (o *ConnectedOrgConfig) GetDomainRestrictionEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.DomainRestrictionEnabled
}

// GetDomainRestrictionEnabledOk returns a tuple with the DomainRestrictionEnabled field value
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetDomainRestrictionEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DomainRestrictionEnabled, true
}

// SetDomainRestrictionEnabled sets field value
func (o *ConnectedOrgConfig) SetDomainRestrictionEnabled(v bool) {
	o.DomainRestrictionEnabled = v
}

// GetIdentityProviderId returns the IdentityProviderId field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetIdentityProviderId() string {
	if o == nil || IsNil(o.IdentityProviderId) {
		var ret string
		return ret
	}
	return *o.IdentityProviderId
}

// GetIdentityProviderIdOk returns a tuple with the IdentityProviderId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetIdentityProviderIdOk() (*string, bool) {
	if o == nil || IsNil(o.IdentityProviderId) {
		return nil, false
	}

	return o.IdentityProviderId, true
}

// HasIdentityProviderId returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasIdentityProviderId() bool {
	if o != nil && !IsNil(o.IdentityProviderId) {
		return true
	}

	return false
}

// SetIdentityProviderId gets a reference to the given string and assigns it to the IdentityProviderId field.
func (o *ConnectedOrgConfig) SetIdentityProviderId(v string) {
	o.IdentityProviderId = &v
}

// GetOrgId returns the OrgId field value
func (o *ConnectedOrgConfig) GetOrgId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetOrgIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrgId, true
}

// SetOrgId sets field value
func (o *ConnectedOrgConfig) SetOrgId(v string) {
	o.OrgId = v
}

// GetPostAuthRoleGrants returns the PostAuthRoleGrants field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetPostAuthRoleGrants() []string {
	if o == nil || IsNil(o.PostAuthRoleGrants) {
		var ret []string
		return ret
	}
	return *o.PostAuthRoleGrants
}

// GetPostAuthRoleGrantsOk returns a tuple with the PostAuthRoleGrants field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetPostAuthRoleGrantsOk() (*[]string, bool) {
	if o == nil || IsNil(o.PostAuthRoleGrants) {
		return nil, false
	}

	return o.PostAuthRoleGrants, true
}

// HasPostAuthRoleGrants returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasPostAuthRoleGrants() bool {
	if o != nil && !IsNil(o.PostAuthRoleGrants) {
		return true
	}

	return false
}

// SetPostAuthRoleGrants gets a reference to the given []string and assigns it to the PostAuthRoleGrants field.
func (o *ConnectedOrgConfig) SetPostAuthRoleGrants(v []string) {
	o.PostAuthRoleGrants = &v
}

// GetRoleMappings returns the RoleMappings field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetRoleMappings() []AuthFederationRoleMapping {
	if o == nil || IsNil(o.RoleMappings) {
		var ret []AuthFederationRoleMapping
		return ret
	}
	return *o.RoleMappings
}

// GetRoleMappingsOk returns a tuple with the RoleMappings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetRoleMappingsOk() (*[]AuthFederationRoleMapping, bool) {
	if o == nil || IsNil(o.RoleMappings) {
		return nil, false
	}

	return o.RoleMappings, true
}

// HasRoleMappings returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasRoleMappings() bool {
	if o != nil && !IsNil(o.RoleMappings) {
		return true
	}

	return false
}

// SetRoleMappings gets a reference to the given []AuthFederationRoleMapping and assigns it to the RoleMappings field.
func (o *ConnectedOrgConfig) SetRoleMappings(v []AuthFederationRoleMapping) {
	o.RoleMappings = &v
}

// GetUserConflicts returns the UserConflicts field value if set, zero value otherwise
func (o *ConnectedOrgConfig) GetUserConflicts() []FederatedUser {
	if o == nil || IsNil(o.UserConflicts) {
		var ret []FederatedUser
		return ret
	}
	return *o.UserConflicts
}

// GetUserConflictsOk returns a tuple with the UserConflicts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectedOrgConfig) GetUserConflictsOk() (*[]FederatedUser, bool) {
	if o == nil || IsNil(o.UserConflicts) {
		return nil, false
	}

	return o.UserConflicts, true
}

// HasUserConflicts returns a boolean if a field has been set.
func (o *ConnectedOrgConfig) HasUserConflicts() bool {
	if o != nil && !IsNil(o.UserConflicts) {
		return true
	}

	return false
}

// SetUserConflicts gets a reference to the given []FederatedUser and assigns it to the UserConflicts field.
func (o *ConnectedOrgConfig) SetUserConflicts(v []FederatedUser) {
	o.UserConflicts = &v
}
