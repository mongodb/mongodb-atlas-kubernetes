// Code based on the AtlasAPI V2 OpenAPI file

package admin

// LDAPVerifyConnectivityJobRequestParams Request information needed to verify an Lightweight Directory Access Protocol (LDAP) over Transport Layer Security (TLS) configuration. The response does not return the `bindPassword`.
type LDAPVerifyConnectivityJobRequestParams struct {
	// Lightweight Directory Access Protocol (LDAP) query template that MongoDB Cloud applies to create an LDAP query to return the LDAP groups associated with the authenticated MongoDB user. MongoDB Cloud uses this parameter only for user authorization.  Use the `{USER}` placeholder in the Uniform Resource Locator (URL) to substitute the authenticated username. The query relates to the host specified with the hostname. Format this query per [RFC 4515](https://datatracker.ietf.org/doc/html/rfc4515) and [RFC 4516](https://datatracker.ietf.org/doc/html/rfc4516).
	// Write only field.
	AuthzQueryTemplate *string `json:"authzQueryTemplate,omitempty"`
	// Password that MongoDB Cloud uses to authenticate the `bindUsername`.
	// Write only field.
	BindPassword string `json:"bindPassword"`
	// Full Distinguished Name (DN) of the Lightweight Directory Access Protocol (LDAP) user that MongoDB Cloud uses to connect to the LDAP host. LDAP distinguished names must be formatted according to RFC 2253.
	BindUsername string `json:"bindUsername"`
	// Certificate Authority (CA) certificate that MongoDB Cloud uses to verify the identity of the Lightweight Directory Access Protocol (LDAP) host. MongoDB Cloud allows self-signed certificates. To delete an assigned value, pass an empty string: `\"caCertificate\": \"\"`.
	CaCertificate *string `json:"caCertificate,omitempty"`
	// Human-readable label that identifies the hostname or Internet Protocol (IP) address of the Lightweight Directory Access Protocol (LDAP) host. This host must have access to the internet or have a Virtual Private Cloud (VPC) peering connection to your cluster.
	Hostname string `json:"hostname"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// IANA port to which the Lightweight Directory Access Protocol (LDAP) host listens for client connections.
	Port int `json:"port"`
}

// NewLDAPVerifyConnectivityJobRequestParams instantiates a new LDAPVerifyConnectivityJobRequestParams object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLDAPVerifyConnectivityJobRequestParams(bindPassword string, bindUsername string, hostname string, port int) *LDAPVerifyConnectivityJobRequestParams {
	this := LDAPVerifyConnectivityJobRequestParams{}
	var authzQueryTemplate string = "{USER}?memberOf?base"
	this.AuthzQueryTemplate = &authzQueryTemplate
	this.BindPassword = bindPassword
	this.BindUsername = bindUsername
	this.Hostname = hostname
	this.Port = port
	return &this
}

// NewLDAPVerifyConnectivityJobRequestParamsWithDefaults instantiates a new LDAPVerifyConnectivityJobRequestParams object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLDAPVerifyConnectivityJobRequestParamsWithDefaults() *LDAPVerifyConnectivityJobRequestParams {
	this := LDAPVerifyConnectivityJobRequestParams{}
	var authzQueryTemplate string = "{USER}?memberOf?base"
	this.AuthzQueryTemplate = &authzQueryTemplate
	var port int = 636
	this.Port = port
	return &this
}

// GetAuthzQueryTemplate returns the AuthzQueryTemplate field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequestParams) GetAuthzQueryTemplate() string {
	if o == nil || IsNil(o.AuthzQueryTemplate) {
		var ret string
		return ret
	}
	return *o.AuthzQueryTemplate
}

// GetAuthzQueryTemplateOk returns a tuple with the AuthzQueryTemplate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetAuthzQueryTemplateOk() (*string, bool) {
	if o == nil || IsNil(o.AuthzQueryTemplate) {
		return nil, false
	}

	return o.AuthzQueryTemplate, true
}

// HasAuthzQueryTemplate returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) HasAuthzQueryTemplate() bool {
	if o != nil && !IsNil(o.AuthzQueryTemplate) {
		return true
	}

	return false
}

// SetAuthzQueryTemplate gets a reference to the given string and assigns it to the AuthzQueryTemplate field.
func (o *LDAPVerifyConnectivityJobRequestParams) SetAuthzQueryTemplate(v string) {
	o.AuthzQueryTemplate = &v
}

// GetBindPassword returns the BindPassword field value
func (o *LDAPVerifyConnectivityJobRequestParams) GetBindPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BindPassword
}

// GetBindPasswordOk returns a tuple with the BindPassword field value
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetBindPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BindPassword, true
}

// SetBindPassword sets field value
func (o *LDAPVerifyConnectivityJobRequestParams) SetBindPassword(v string) {
	o.BindPassword = v
}

// GetBindUsername returns the BindUsername field value
func (o *LDAPVerifyConnectivityJobRequestParams) GetBindUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BindUsername
}

// GetBindUsernameOk returns a tuple with the BindUsername field value
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetBindUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BindUsername, true
}

// SetBindUsername sets field value
func (o *LDAPVerifyConnectivityJobRequestParams) SetBindUsername(v string) {
	o.BindUsername = v
}

// GetCaCertificate returns the CaCertificate field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequestParams) GetCaCertificate() string {
	if o == nil || IsNil(o.CaCertificate) {
		var ret string
		return ret
	}
	return *o.CaCertificate
}

// GetCaCertificateOk returns a tuple with the CaCertificate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetCaCertificateOk() (*string, bool) {
	if o == nil || IsNil(o.CaCertificate) {
		return nil, false
	}

	return o.CaCertificate, true
}

// HasCaCertificate returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) HasCaCertificate() bool {
	if o != nil && !IsNil(o.CaCertificate) {
		return true
	}

	return false
}

// SetCaCertificate gets a reference to the given string and assigns it to the CaCertificate field.
func (o *LDAPVerifyConnectivityJobRequestParams) SetCaCertificate(v string) {
	o.CaCertificate = &v
}

// GetHostname returns the Hostname field value
func (o *LDAPVerifyConnectivityJobRequestParams) GetHostname() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hostname
}

// GetHostnameOk returns a tuple with the Hostname field value
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetHostnameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Hostname, true
}

// SetHostname sets field value
func (o *LDAPVerifyConnectivityJobRequestParams) SetHostname(v string) {
	o.Hostname = v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *LDAPVerifyConnectivityJobRequestParams) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *LDAPVerifyConnectivityJobRequestParams) SetLinks(v []Link) {
	o.Links = &v
}

// GetPort returns the Port field value
func (o *LDAPVerifyConnectivityJobRequestParams) GetPort() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.Port
}

// GetPortOk returns a tuple with the Port field value
// and a boolean to check if the value has been set.
func (o *LDAPVerifyConnectivityJobRequestParams) GetPortOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Port, true
}

// SetPort sets field value
func (o *LDAPVerifyConnectivityJobRequestParams) SetPort(v int) {
	o.Port = v
}
