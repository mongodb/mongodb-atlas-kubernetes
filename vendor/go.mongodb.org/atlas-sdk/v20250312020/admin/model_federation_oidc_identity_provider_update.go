// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FederationOidcIdentityProviderUpdate struct for FederationOidcIdentityProviderUpdate
type FederationOidcIdentityProviderUpdate struct {
	// Identifier of the intended recipient of the token.
	Audience *string `json:"audience,omitempty"`
	// Indicates whether authorization is granted based on group membership or user ID.
	AuthorizationType *string `json:"authorizationType,omitempty"`
	// The description of the identity provider.
	Description *string `json:"description,omitempty"`
	// Human-readable label that identifies the identity provider.
	DisplayName *string `json:"displayName,omitempty"`
	// Identifier of the claim which contains IdP Group IDs in the token.
	GroupsClaim *string `json:"groupsClaim,omitempty"`
	// String enum that indicates the type of the identity provider. Default is WORKFORCE.
	IdpType *string `json:"idpType,omitempty"`
	// Unique string that identifies the issuer of the SAML Assertion or OIDC metadata/discovery document URL.
	IssuerUri *string `json:"issuerUri,omitempty"`
	// String enum that indicates the protocol of the identity provider. Either SAML or OIDC.
	Protocol *string `json:"protocol,omitempty"`
	// Identifier of the claim which contains the user ID in the token.
	UserClaim *string `json:"userClaim,omitempty"`
	// List that contains the domains associated with the identity provider.
	AssociatedDomains *[]string `json:"associatedDomains,omitempty"`
	// Client identifier that is assigned to an application by the Identity Provider.
	ClientId *string `json:"clientId,omitempty"`
	// Scopes that MongoDB applications will request from the authorization endpoint.
	RequestedScopes *[]string `json:"requestedScopes,omitempty"`
}

// NewFederationOidcIdentityProviderUpdate instantiates a new FederationOidcIdentityProviderUpdate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFederationOidcIdentityProviderUpdate() *FederationOidcIdentityProviderUpdate {
	this := FederationOidcIdentityProviderUpdate{}
	return &this
}

// NewFederationOidcIdentityProviderUpdateWithDefaults instantiates a new FederationOidcIdentityProviderUpdate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFederationOidcIdentityProviderUpdateWithDefaults() *FederationOidcIdentityProviderUpdate {
	this := FederationOidcIdentityProviderUpdate{}
	return &this
}

// GetAudience returns the Audience field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetAudience() string {
	if o == nil || IsNil(o.Audience) {
		var ret string
		return ret
	}
	return *o.Audience
}

// GetAudienceOk returns a tuple with the Audience field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetAudienceOk() (*string, bool) {
	if o == nil || IsNil(o.Audience) {
		return nil, false
	}

	return o.Audience, true
}

// HasAudience returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasAudience() bool {
	if o != nil && !IsNil(o.Audience) {
		return true
	}

	return false
}

// SetAudience gets a reference to the given string and assigns it to the Audience field.
func (o *FederationOidcIdentityProviderUpdate) SetAudience(v string) {
	o.Audience = &v
}

// GetAuthorizationType returns the AuthorizationType field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetAuthorizationType() string {
	if o == nil || IsNil(o.AuthorizationType) {
		var ret string
		return ret
	}
	return *o.AuthorizationType
}

// GetAuthorizationTypeOk returns a tuple with the AuthorizationType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetAuthorizationTypeOk() (*string, bool) {
	if o == nil || IsNil(o.AuthorizationType) {
		return nil, false
	}

	return o.AuthorizationType, true
}

// HasAuthorizationType returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasAuthorizationType() bool {
	if o != nil && !IsNil(o.AuthorizationType) {
		return true
	}

	return false
}

// SetAuthorizationType gets a reference to the given string and assigns it to the AuthorizationType field.
func (o *FederationOidcIdentityProviderUpdate) SetAuthorizationType(v string) {
	o.AuthorizationType = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *FederationOidcIdentityProviderUpdate) SetDescription(v string) {
	o.Description = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}

	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *FederationOidcIdentityProviderUpdate) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetGroupsClaim returns the GroupsClaim field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetGroupsClaim() string {
	if o == nil || IsNil(o.GroupsClaim) {
		var ret string
		return ret
	}
	return *o.GroupsClaim
}

// GetGroupsClaimOk returns a tuple with the GroupsClaim field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetGroupsClaimOk() (*string, bool) {
	if o == nil || IsNil(o.GroupsClaim) {
		return nil, false
	}

	return o.GroupsClaim, true
}

// HasGroupsClaim returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasGroupsClaim() bool {
	if o != nil && !IsNil(o.GroupsClaim) {
		return true
	}

	return false
}

// SetGroupsClaim gets a reference to the given string and assigns it to the GroupsClaim field.
func (o *FederationOidcIdentityProviderUpdate) SetGroupsClaim(v string) {
	o.GroupsClaim = &v
}

// GetIdpType returns the IdpType field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetIdpType() string {
	if o == nil || IsNil(o.IdpType) {
		var ret string
		return ret
	}
	return *o.IdpType
}

// GetIdpTypeOk returns a tuple with the IdpType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetIdpTypeOk() (*string, bool) {
	if o == nil || IsNil(o.IdpType) {
		return nil, false
	}

	return o.IdpType, true
}

// HasIdpType returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasIdpType() bool {
	if o != nil && !IsNil(o.IdpType) {
		return true
	}

	return false
}

// SetIdpType gets a reference to the given string and assigns it to the IdpType field.
func (o *FederationOidcIdentityProviderUpdate) SetIdpType(v string) {
	o.IdpType = &v
}

// GetIssuerUri returns the IssuerUri field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetIssuerUri() string {
	if o == nil || IsNil(o.IssuerUri) {
		var ret string
		return ret
	}
	return *o.IssuerUri
}

// GetIssuerUriOk returns a tuple with the IssuerUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetIssuerUriOk() (*string, bool) {
	if o == nil || IsNil(o.IssuerUri) {
		return nil, false
	}

	return o.IssuerUri, true
}

// HasIssuerUri returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasIssuerUri() bool {
	if o != nil && !IsNil(o.IssuerUri) {
		return true
	}

	return false
}

// SetIssuerUri gets a reference to the given string and assigns it to the IssuerUri field.
func (o *FederationOidcIdentityProviderUpdate) SetIssuerUri(v string) {
	o.IssuerUri = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}

	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *FederationOidcIdentityProviderUpdate) SetProtocol(v string) {
	o.Protocol = &v
}

// GetUserClaim returns the UserClaim field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetUserClaim() string {
	if o == nil || IsNil(o.UserClaim) {
		var ret string
		return ret
	}
	return *o.UserClaim
}

// GetUserClaimOk returns a tuple with the UserClaim field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetUserClaimOk() (*string, bool) {
	if o == nil || IsNil(o.UserClaim) {
		return nil, false
	}

	return o.UserClaim, true
}

// HasUserClaim returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasUserClaim() bool {
	if o != nil && !IsNil(o.UserClaim) {
		return true
	}

	return false
}

// SetUserClaim gets a reference to the given string and assigns it to the UserClaim field.
func (o *FederationOidcIdentityProviderUpdate) SetUserClaim(v string) {
	o.UserClaim = &v
}

// GetAssociatedDomains returns the AssociatedDomains field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetAssociatedDomains() []string {
	if o == nil || IsNil(o.AssociatedDomains) {
		var ret []string
		return ret
	}
	return *o.AssociatedDomains
}

// GetAssociatedDomainsOk returns a tuple with the AssociatedDomains field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetAssociatedDomainsOk() (*[]string, bool) {
	if o == nil || IsNil(o.AssociatedDomains) {
		return nil, false
	}

	return o.AssociatedDomains, true
}

// HasAssociatedDomains returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasAssociatedDomains() bool {
	if o != nil && !IsNil(o.AssociatedDomains) {
		return true
	}

	return false
}

// SetAssociatedDomains gets a reference to the given []string and assigns it to the AssociatedDomains field.
func (o *FederationOidcIdentityProviderUpdate) SetAssociatedDomains(v []string) {
	o.AssociatedDomains = &v
}

// GetClientId returns the ClientId field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetClientId() string {
	if o == nil || IsNil(o.ClientId) {
		var ret string
		return ret
	}
	return *o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetClientIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClientId) {
		return nil, false
	}

	return o.ClientId, true
}

// HasClientId returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasClientId() bool {
	if o != nil && !IsNil(o.ClientId) {
		return true
	}

	return false
}

// SetClientId gets a reference to the given string and assigns it to the ClientId field.
func (o *FederationOidcIdentityProviderUpdate) SetClientId(v string) {
	o.ClientId = &v
}

// GetRequestedScopes returns the RequestedScopes field value if set, zero value otherwise
func (o *FederationOidcIdentityProviderUpdate) GetRequestedScopes() []string {
	if o == nil || IsNil(o.RequestedScopes) {
		var ret []string
		return ret
	}
	return *o.RequestedScopes
}

// GetRequestedScopesOk returns a tuple with the RequestedScopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProviderUpdate) GetRequestedScopesOk() (*[]string, bool) {
	if o == nil || IsNil(o.RequestedScopes) {
		return nil, false
	}

	return o.RequestedScopes, true
}

// HasRequestedScopes returns a boolean if a field has been set.
func (o *FederationOidcIdentityProviderUpdate) HasRequestedScopes() bool {
	if o != nil && !IsNil(o.RequestedScopes) {
		return true
	}

	return false
}

// SetRequestedScopes gets a reference to the given []string and assigns it to the RequestedScopes field.
func (o *FederationOidcIdentityProviderUpdate) SetRequestedScopes(v []string) {
	o.RequestedScopes = &v
}
