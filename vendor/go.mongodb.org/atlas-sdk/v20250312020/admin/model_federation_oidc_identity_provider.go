// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// FederationOidcIdentityProvider struct for FederationOidcIdentityProvider
type FederationOidcIdentityProvider struct {
	// List that contains the connected organization configurations associated with the identity provider.
	AssociatedOrgs *[]ConnectedOrgConfig `json:"associatedOrgs,omitempty"`
	// Identifier of the intended recipient of the token.
	Audience *string `json:"audience,omitempty"`
	// Indicates whether authorization is granted based on group membership or user ID.
	AuthorizationType *string `json:"authorizationType,omitempty"`
	// Date that the identity provider was created on. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// The description of the identity provider.
	Description *string `json:"description,omitempty"`
	// Human-readable label that identifies the identity provider.
	DisplayName *string `json:"displayName,omitempty"`
	// Identifier of the claim which contains IdP Group IDs in the token.
	GroupsClaim *string `json:"groupsClaim,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the identity provider.
	// Read only field.
	Id string `json:"id"`
	// String enum that indicates the type of the identity provider. Default is WORKFORCE.
	IdpType *string `json:"idpType,omitempty"`
	// Unique string that identifies the issuer of the SAML Assertion or OIDC metadata/discovery document URL.
	IssuerUri *string `json:"issuerUri,omitempty"`
	// Legacy 20-hexadecimal digit string that identifies the identity provider.
	OktaIdpId string `json:"oktaIdpId"`
	// String enum that indicates the protocol of the identity provider. Either SAML or OIDC.
	Protocol *string `json:"protocol,omitempty"`
	// Date that the identity provider was last updated on. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	// Identifier of the claim which contains the user ID in the token.
	UserClaim *string `json:"userClaim,omitempty"`
	// List that contains the domains associated with the identity provider.
	AssociatedDomains *[]string `json:"associatedDomains,omitempty"`
	// Client identifier that is assigned to an application by the Identity Provider.
	ClientId *string `json:"clientId,omitempty"`
	// Scopes that MongoDB applications will request from the authorization endpoint.
	RequestedScopes *[]string `json:"requestedScopes,omitempty"`
}

// NewFederationOidcIdentityProvider instantiates a new FederationOidcIdentityProvider object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFederationOidcIdentityProvider(id string, oktaIdpId string) *FederationOidcIdentityProvider {
	this := FederationOidcIdentityProvider{}
	this.Id = id
	this.OktaIdpId = oktaIdpId
	return &this
}

// NewFederationOidcIdentityProviderWithDefaults instantiates a new FederationOidcIdentityProvider object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFederationOidcIdentityProviderWithDefaults() *FederationOidcIdentityProvider {
	this := FederationOidcIdentityProvider{}
	return &this
}

// GetAssociatedOrgs returns the AssociatedOrgs field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetAssociatedOrgs() []ConnectedOrgConfig {
	if o == nil || IsNil(o.AssociatedOrgs) {
		var ret []ConnectedOrgConfig
		return ret
	}
	return *o.AssociatedOrgs
}

// GetAssociatedOrgsOk returns a tuple with the AssociatedOrgs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetAssociatedOrgsOk() (*[]ConnectedOrgConfig, bool) {
	if o == nil || IsNil(o.AssociatedOrgs) {
		return nil, false
	}

	return o.AssociatedOrgs, true
}

// HasAssociatedOrgs returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasAssociatedOrgs() bool {
	if o != nil && !IsNil(o.AssociatedOrgs) {
		return true
	}

	return false
}

// SetAssociatedOrgs gets a reference to the given []ConnectedOrgConfig and assigns it to the AssociatedOrgs field.
func (o *FederationOidcIdentityProvider) SetAssociatedOrgs(v []ConnectedOrgConfig) {
	o.AssociatedOrgs = &v
}

// GetAudience returns the Audience field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetAudience() string {
	if o == nil || IsNil(o.Audience) {
		var ret string
		return ret
	}
	return *o.Audience
}

// GetAudienceOk returns a tuple with the Audience field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetAudienceOk() (*string, bool) {
	if o == nil || IsNil(o.Audience) {
		return nil, false
	}

	return o.Audience, true
}

// HasAudience returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasAudience() bool {
	if o != nil && !IsNil(o.Audience) {
		return true
	}

	return false
}

// SetAudience gets a reference to the given string and assigns it to the Audience field.
func (o *FederationOidcIdentityProvider) SetAudience(v string) {
	o.Audience = &v
}

// GetAuthorizationType returns the AuthorizationType field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetAuthorizationType() string {
	if o == nil || IsNil(o.AuthorizationType) {
		var ret string
		return ret
	}
	return *o.AuthorizationType
}

// GetAuthorizationTypeOk returns a tuple with the AuthorizationType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetAuthorizationTypeOk() (*string, bool) {
	if o == nil || IsNil(o.AuthorizationType) {
		return nil, false
	}

	return o.AuthorizationType, true
}

// HasAuthorizationType returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasAuthorizationType() bool {
	if o != nil && !IsNil(o.AuthorizationType) {
		return true
	}

	return false
}

// SetAuthorizationType gets a reference to the given string and assigns it to the AuthorizationType field.
func (o *FederationOidcIdentityProvider) SetAuthorizationType(v string) {
	o.AuthorizationType = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetCreatedAt() time.Time {
	if o == nil || IsNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}

	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *FederationOidcIdentityProvider) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetDescription returns the Description field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}

	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *FederationOidcIdentityProvider) SetDescription(v string) {
	o.Description = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}

	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *FederationOidcIdentityProvider) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetGroupsClaim returns the GroupsClaim field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetGroupsClaim() string {
	if o == nil || IsNil(o.GroupsClaim) {
		var ret string
		return ret
	}
	return *o.GroupsClaim
}

// GetGroupsClaimOk returns a tuple with the GroupsClaim field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetGroupsClaimOk() (*string, bool) {
	if o == nil || IsNil(o.GroupsClaim) {
		return nil, false
	}

	return o.GroupsClaim, true
}

// HasGroupsClaim returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasGroupsClaim() bool {
	if o != nil && !IsNil(o.GroupsClaim) {
		return true
	}

	return false
}

// SetGroupsClaim gets a reference to the given string and assigns it to the GroupsClaim field.
func (o *FederationOidcIdentityProvider) SetGroupsClaim(v string) {
	o.GroupsClaim = &v
}

// GetId returns the Id field value
func (o *FederationOidcIdentityProvider) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *FederationOidcIdentityProvider) SetId(v string) {
	o.Id = v
}

// GetIdpType returns the IdpType field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetIdpType() string {
	if o == nil || IsNil(o.IdpType) {
		var ret string
		return ret
	}
	return *o.IdpType
}

// GetIdpTypeOk returns a tuple with the IdpType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetIdpTypeOk() (*string, bool) {
	if o == nil || IsNil(o.IdpType) {
		return nil, false
	}

	return o.IdpType, true
}

// HasIdpType returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasIdpType() bool {
	if o != nil && !IsNil(o.IdpType) {
		return true
	}

	return false
}

// SetIdpType gets a reference to the given string and assigns it to the IdpType field.
func (o *FederationOidcIdentityProvider) SetIdpType(v string) {
	o.IdpType = &v
}

// GetIssuerUri returns the IssuerUri field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetIssuerUri() string {
	if o == nil || IsNil(o.IssuerUri) {
		var ret string
		return ret
	}
	return *o.IssuerUri
}

// GetIssuerUriOk returns a tuple with the IssuerUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetIssuerUriOk() (*string, bool) {
	if o == nil || IsNil(o.IssuerUri) {
		return nil, false
	}

	return o.IssuerUri, true
}

// HasIssuerUri returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasIssuerUri() bool {
	if o != nil && !IsNil(o.IssuerUri) {
		return true
	}

	return false
}

// SetIssuerUri gets a reference to the given string and assigns it to the IssuerUri field.
func (o *FederationOidcIdentityProvider) SetIssuerUri(v string) {
	o.IssuerUri = &v
}

// GetOktaIdpId returns the OktaIdpId field value
func (o *FederationOidcIdentityProvider) GetOktaIdpId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OktaIdpId
}

// GetOktaIdpIdOk returns a tuple with the OktaIdpId field value
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetOktaIdpIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OktaIdpId, true
}

// SetOktaIdpId sets field value
func (o *FederationOidcIdentityProvider) SetOktaIdpId(v string) {
	o.OktaIdpId = v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}

	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *FederationOidcIdentityProvider) SetProtocol(v string) {
	o.Protocol = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetUpdatedAt() time.Time {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret time.Time
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}

	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given time.Time and assigns it to the UpdatedAt field.
func (o *FederationOidcIdentityProvider) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = &v
}

// GetUserClaim returns the UserClaim field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetUserClaim() string {
	if o == nil || IsNil(o.UserClaim) {
		var ret string
		return ret
	}
	return *o.UserClaim
}

// GetUserClaimOk returns a tuple with the UserClaim field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetUserClaimOk() (*string, bool) {
	if o == nil || IsNil(o.UserClaim) {
		return nil, false
	}

	return o.UserClaim, true
}

// HasUserClaim returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasUserClaim() bool {
	if o != nil && !IsNil(o.UserClaim) {
		return true
	}

	return false
}

// SetUserClaim gets a reference to the given string and assigns it to the UserClaim field.
func (o *FederationOidcIdentityProvider) SetUserClaim(v string) {
	o.UserClaim = &v
}

// GetAssociatedDomains returns the AssociatedDomains field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetAssociatedDomains() []string {
	if o == nil || IsNil(o.AssociatedDomains) {
		var ret []string
		return ret
	}
	return *o.AssociatedDomains
}

// GetAssociatedDomainsOk returns a tuple with the AssociatedDomains field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetAssociatedDomainsOk() (*[]string, bool) {
	if o == nil || IsNil(o.AssociatedDomains) {
		return nil, false
	}

	return o.AssociatedDomains, true
}

// HasAssociatedDomains returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasAssociatedDomains() bool {
	if o != nil && !IsNil(o.AssociatedDomains) {
		return true
	}

	return false
}

// SetAssociatedDomains gets a reference to the given []string and assigns it to the AssociatedDomains field.
func (o *FederationOidcIdentityProvider) SetAssociatedDomains(v []string) {
	o.AssociatedDomains = &v
}

// GetClientId returns the ClientId field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetClientId() string {
	if o == nil || IsNil(o.ClientId) {
		var ret string
		return ret
	}
	return *o.ClientId
}

// GetClientIdOk returns a tuple with the ClientId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetClientIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClientId) {
		return nil, false
	}

	return o.ClientId, true
}

// HasClientId returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasClientId() bool {
	if o != nil && !IsNil(o.ClientId) {
		return true
	}

	return false
}

// SetClientId gets a reference to the given string and assigns it to the ClientId field.
func (o *FederationOidcIdentityProvider) SetClientId(v string) {
	o.ClientId = &v
}

// GetRequestedScopes returns the RequestedScopes field value if set, zero value otherwise
func (o *FederationOidcIdentityProvider) GetRequestedScopes() []string {
	if o == nil || IsNil(o.RequestedScopes) {
		var ret []string
		return ret
	}
	return *o.RequestedScopes
}

// GetRequestedScopesOk returns a tuple with the RequestedScopes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FederationOidcIdentityProvider) GetRequestedScopesOk() (*[]string, bool) {
	if o == nil || IsNil(o.RequestedScopes) {
		return nil, false
	}

	return o.RequestedScopes, true
}

// HasRequestedScopes returns a boolean if a field has been set.
func (o *FederationOidcIdentityProvider) HasRequestedScopes() bool {
	if o != nil && !IsNil(o.RequestedScopes) {
		return true
	}

	return false
}

// SetRequestedScopes gets a reference to the given []string and assigns it to the RequestedScopes field.
func (o *FederationOidcIdentityProvider) SetRequestedScopes(v []string) {
	o.RequestedScopes = &v
}
