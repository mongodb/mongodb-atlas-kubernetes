// Code based on the AtlasAPI V2 OpenAPI file

package admin

// UserToDNMapping User-to-Distinguished Name (DN) map that MongoDB Cloud uses to transform a Lightweight Directory Access Protocol (LDAP) username into an LDAP DN.
type UserToDNMapping struct {
	// Lightweight Directory Access Protocol (LDAP) query template that inserts the LDAP name that the regular expression matches into an LDAP query Uniform Resource Identifier (URI). The formatting for the query must conform to [RFC 4515](https://datatracker.ietf.org/doc/html/rfc4515) and [RFC 4516](https://datatracker.ietf.org/doc/html/rfc4516).
	LdapQuery *string `json:"ldapQuery,omitempty"`
	// Regular expression that MongoDB Cloud uses to match against the provided Lightweight Directory Access Protocol (LDAP) username. Each parenthesis-enclosed section represents a regular expression capture group that the substitution or `ldapQuery` template uses.
	Match string `json:"match"`
	// Lightweight Directory Access Protocol (LDAP) Distinguished Name (DN) template that converts the LDAP username that matches regular expression in the *match* parameter into an LDAP Distinguished Name (DN).
	Substitution *string `json:"substitution,omitempty"`
}

// NewUserToDNMapping instantiates a new UserToDNMapping object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserToDNMapping(match string) *UserToDNMapping {
	this := UserToDNMapping{}
	this.Match = match
	return &this
}

// NewUserToDNMappingWithDefaults instantiates a new UserToDNMapping object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserToDNMappingWithDefaults() *UserToDNMapping {
	this := UserToDNMapping{}
	return &this
}

// GetLdapQuery returns the LdapQuery field value if set, zero value otherwise
func (o *UserToDNMapping) GetLdapQuery() string {
	if o == nil || IsNil(o.LdapQuery) {
		var ret string
		return ret
	}
	return *o.LdapQuery
}

// GetLdapQueryOk returns a tuple with the LdapQuery field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserToDNMapping) GetLdapQueryOk() (*string, bool) {
	if o == nil || IsNil(o.LdapQuery) {
		return nil, false
	}

	return o.LdapQuery, true
}

// HasLdapQuery returns a boolean if a field has been set.
func (o *UserToDNMapping) HasLdapQuery() bool {
	if o != nil && !IsNil(o.LdapQuery) {
		return true
	}

	return false
}

// SetLdapQuery gets a reference to the given string and assigns it to the LdapQuery field.
func (o *UserToDNMapping) SetLdapQuery(v string) {
	o.LdapQuery = &v
}

// GetMatch returns the Match field value
func (o *UserToDNMapping) GetMatch() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Match
}

// GetMatchOk returns a tuple with the Match field value
// and a boolean to check if the value has been set.
func (o *UserToDNMapping) GetMatchOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Match, true
}

// SetMatch sets field value
func (o *UserToDNMapping) SetMatch(v string) {
	o.Match = v
}

// GetSubstitution returns the Substitution field value if set, zero value otherwise
func (o *UserToDNMapping) GetSubstitution() string {
	if o == nil || IsNil(o.Substitution) {
		var ret string
		return ret
	}
	return *o.Substitution
}

// GetSubstitutionOk returns a tuple with the Substitution field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserToDNMapping) GetSubstitutionOk() (*string, bool) {
	if o == nil || IsNil(o.Substitution) {
		return nil, false
	}

	return o.Substitution, true
}

// HasSubstitution returns a boolean if a field has been set.
func (o *UserToDNMapping) HasSubstitution() bool {
	if o != nil && !IsNil(o.Substitution) {
		return true
	}

	return false
}

// SetSubstitution gets a reference to the given string and assigns it to the Substitution field.
func (o *UserToDNMapping) SetSubstitution(v string) {
	o.Substitution = &v
}
