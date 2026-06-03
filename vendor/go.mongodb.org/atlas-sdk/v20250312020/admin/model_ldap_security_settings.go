// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LDAPSecuritySettings Lightweight Directory Access Protocol (LDAP) over Transport Layer Security (TLS) configuration details that apply to the specified project.
type LDAPSecuritySettings struct {
	// Flag that indicates whether users can authenticate using an Lightweight Directory Access Protocol (LDAP) host.
	AuthenticationEnabled *bool `json:"authenticationEnabled,omitempty"`
	// Flag that indicates whether users can authorize access to MongoDB Cloud resources using an Lightweight Directory Access Protocol (LDAP) host.
	AuthorizationEnabled *bool `json:"authorizationEnabled,omitempty"`
	// Lightweight Directory Access Protocol (LDAP) query template that MongoDB Cloud runs to obtain the LDAP groups associated with the authenticated user. MongoDB Cloud uses this parameter only for user authorization. Use the `{USER}` placeholder in the Uniform Resource Locator (URL) to substitute the authenticated username. The query relates to the host specified with the hostname. Format this query according to [RFC 4515](https://datatracker.ietf.org/doc/html/rfc4515) and [RFC 4516](https://datatracker.ietf.org/doc/html/rfc4516).
	AuthzQueryTemplate *string `json:"authzQueryTemplate,omitempty"`
	// Password that MongoDB Cloud uses to authenticate the `bindUsername`.
	// Write only field.
	BindPassword *string `json:"bindPassword,omitempty"`
	// Full Distinguished Name (DN) of the Lightweight Directory Access Protocol (LDAP) user that MongoDB Cloud uses to connect to the LDAP host. LDAP distinguished names must be formatted according to RFC 2253.
	BindUsername *string `json:"bindUsername,omitempty"`
	// Certificate Authority (CA) certificate that MongoDB Cloud uses to verify the identity of the Lightweight Directory Access Protocol (LDAP) host. MongoDB Cloud allows self-signed certificates. To delete an assigned value, pass an empty string: `\"caCertificate\": \"\"`.
	CaCertificate *string `json:"caCertificate,omitempty"`
	// Human-readable label that identifies the hostname or Internet Protocol (IP) address of the Lightweight Directory Access Protocol (LDAP) host. This host must have access to the internet or have a Virtual Private Cloud (VPC) peering connection to your cluster.
	Hostname *string `json:"hostname,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Port to which the Lightweight Directory Access Protocol (LDAP) host listens for client connections.
	Port *int `json:"port,omitempty"`
	// User-to-Distinguished Name (DN) map that MongoDB Cloud uses to transform a Lightweight Directory Access Protocol (LDAP) username into an LDAP DN.
	UserToDNMapping *[]UserToDNMapping `json:"userToDNMapping,omitempty"`
}

// NewLDAPSecuritySettings instantiates a new LDAPSecuritySettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLDAPSecuritySettings() *LDAPSecuritySettings {
	this := LDAPSecuritySettings{}
	var authzQueryTemplate string = "{USER}?memberOf?base"
	this.AuthzQueryTemplate = &authzQueryTemplate
	var port int = 636
	this.Port = &port
	return &this
}

// NewLDAPSecuritySettingsWithDefaults instantiates a new LDAPSecuritySettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLDAPSecuritySettingsWithDefaults() *LDAPSecuritySettings {
	this := LDAPSecuritySettings{}
	var authzQueryTemplate string = "{USER}?memberOf?base"
	this.AuthzQueryTemplate = &authzQueryTemplate
	var port int = 636
	this.Port = &port
	return &this
}

// GetAuthenticationEnabled returns the AuthenticationEnabled field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetAuthenticationEnabled() bool {
	if o == nil || IsNil(o.AuthenticationEnabled) {
		var ret bool
		return ret
	}
	return *o.AuthenticationEnabled
}

// GetAuthenticationEnabledOk returns a tuple with the AuthenticationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetAuthenticationEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AuthenticationEnabled) {
		return nil, false
	}

	return o.AuthenticationEnabled, true
}

// HasAuthenticationEnabled returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasAuthenticationEnabled() bool {
	if o != nil && !IsNil(o.AuthenticationEnabled) {
		return true
	}

	return false
}

// SetAuthenticationEnabled gets a reference to the given bool and assigns it to the AuthenticationEnabled field.
func (o *LDAPSecuritySettings) SetAuthenticationEnabled(v bool) {
	o.AuthenticationEnabled = &v
}

// GetAuthorizationEnabled returns the AuthorizationEnabled field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetAuthorizationEnabled() bool {
	if o == nil || IsNil(o.AuthorizationEnabled) {
		var ret bool
		return ret
	}
	return *o.AuthorizationEnabled
}

// GetAuthorizationEnabledOk returns a tuple with the AuthorizationEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetAuthorizationEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AuthorizationEnabled) {
		return nil, false
	}

	return o.AuthorizationEnabled, true
}

// HasAuthorizationEnabled returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasAuthorizationEnabled() bool {
	if o != nil && !IsNil(o.AuthorizationEnabled) {
		return true
	}

	return false
}

// SetAuthorizationEnabled gets a reference to the given bool and assigns it to the AuthorizationEnabled field.
func (o *LDAPSecuritySettings) SetAuthorizationEnabled(v bool) {
	o.AuthorizationEnabled = &v
}

// GetAuthzQueryTemplate returns the AuthzQueryTemplate field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetAuthzQueryTemplate() string {
	if o == nil || IsNil(o.AuthzQueryTemplate) {
		var ret string
		return ret
	}
	return *o.AuthzQueryTemplate
}

// GetAuthzQueryTemplateOk returns a tuple with the AuthzQueryTemplate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetAuthzQueryTemplateOk() (*string, bool) {
	if o == nil || IsNil(o.AuthzQueryTemplate) {
		return nil, false
	}

	return o.AuthzQueryTemplate, true
}

// HasAuthzQueryTemplate returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasAuthzQueryTemplate() bool {
	if o != nil && !IsNil(o.AuthzQueryTemplate) {
		return true
	}

	return false
}

// SetAuthzQueryTemplate gets a reference to the given string and assigns it to the AuthzQueryTemplate field.
func (o *LDAPSecuritySettings) SetAuthzQueryTemplate(v string) {
	o.AuthzQueryTemplate = &v
}

// GetBindPassword returns the BindPassword field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetBindPassword() string {
	if o == nil || IsNil(o.BindPassword) {
		var ret string
		return ret
	}
	return *o.BindPassword
}

// GetBindPasswordOk returns a tuple with the BindPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetBindPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.BindPassword) {
		return nil, false
	}

	return o.BindPassword, true
}

// HasBindPassword returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasBindPassword() bool {
	if o != nil && !IsNil(o.BindPassword) {
		return true
	}

	return false
}

// SetBindPassword gets a reference to the given string and assigns it to the BindPassword field.
func (o *LDAPSecuritySettings) SetBindPassword(v string) {
	o.BindPassword = &v
}

// GetBindUsername returns the BindUsername field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetBindUsername() string {
	if o == nil || IsNil(o.BindUsername) {
		var ret string
		return ret
	}
	return *o.BindUsername
}

// GetBindUsernameOk returns a tuple with the BindUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetBindUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.BindUsername) {
		return nil, false
	}

	return o.BindUsername, true
}

// HasBindUsername returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasBindUsername() bool {
	if o != nil && !IsNil(o.BindUsername) {
		return true
	}

	return false
}

// SetBindUsername gets a reference to the given string and assigns it to the BindUsername field.
func (o *LDAPSecuritySettings) SetBindUsername(v string) {
	o.BindUsername = &v
}

// GetCaCertificate returns the CaCertificate field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetCaCertificate() string {
	if o == nil || IsNil(o.CaCertificate) {
		var ret string
		return ret
	}
	return *o.CaCertificate
}

// GetCaCertificateOk returns a tuple with the CaCertificate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetCaCertificateOk() (*string, bool) {
	if o == nil || IsNil(o.CaCertificate) {
		return nil, false
	}

	return o.CaCertificate, true
}

// HasCaCertificate returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasCaCertificate() bool {
	if o != nil && !IsNil(o.CaCertificate) {
		return true
	}

	return false
}

// SetCaCertificate gets a reference to the given string and assigns it to the CaCertificate field.
func (o *LDAPSecuritySettings) SetCaCertificate(v string) {
	o.CaCertificate = &v
}

// GetHostname returns the Hostname field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetHostname() string {
	if o == nil || IsNil(o.Hostname) {
		var ret string
		return ret
	}
	return *o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetHostnameOk() (*string, bool) {
	if o == nil || IsNil(o.Hostname) {
		return nil, false
	}

	return o.Hostname, true
}

// HasHostname returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasHostname() bool {
	if o != nil && !IsNil(o.Hostname) {
		return true
	}

	return false
}

// SetHostname gets a reference to the given string and assigns it to the Hostname field.
func (o *LDAPSecuritySettings) SetHostname(v string) {
	o.Hostname = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *LDAPSecuritySettings) SetLinks(v []Link) {
	o.Links = &v
}

// GetPort returns the Port field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetPort() int {
	if o == nil || IsNil(o.Port) {
		var ret int
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetPortOk() (*int, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}

	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int and assigns it to the Port field.
func (o *LDAPSecuritySettings) SetPort(v int) {
	o.Port = &v
}

// GetUserToDNMapping returns the UserToDNMapping field value if set, zero value otherwise
func (o *LDAPSecuritySettings) GetUserToDNMapping() []UserToDNMapping {
	if o == nil || IsNil(o.UserToDNMapping) {
		var ret []UserToDNMapping
		return ret
	}
	return *o.UserToDNMapping
}

// GetUserToDNMappingOk returns a tuple with the UserToDNMapping field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPSecuritySettings) GetUserToDNMappingOk() (*[]UserToDNMapping, bool) {
	if o == nil || IsNil(o.UserToDNMapping) {
		return nil, false
	}

	return o.UserToDNMapping, true
}

// HasUserToDNMapping returns a boolean if a field has been set.
func (o *LDAPSecuritySettings) HasUserToDNMapping() bool {
	if o != nil && !IsNil(o.UserToDNMapping) {
		return true
	}

	return false
}

// SetUserToDNMapping gets a reference to the given []UserToDNMapping and assigns it to the UserToDNMapping field.
func (o *LDAPSecuritySettings) SetUserToDNMapping(v []UserToDNMapping) {
	o.UserToDNMapping = &v
}
