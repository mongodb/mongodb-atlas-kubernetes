// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DBUserTLSX509Settings Settings to configure TLS Certificates for database users.
type DBUserTLSX509Settings struct {
	// Concatenated list of customer certificate authority (CA) certificates needed to authenticate database users. MongoDB Cloud expects this as a PEM-formatted certificate.
	Cas *string `json:"cas,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
}

// NewDBUserTLSX509Settings instantiates a new DBUserTLSX509Settings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDBUserTLSX509Settings() *DBUserTLSX509Settings {
	this := DBUserTLSX509Settings{}
	return &this
}

// NewDBUserTLSX509SettingsWithDefaults instantiates a new DBUserTLSX509Settings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDBUserTLSX509SettingsWithDefaults() *DBUserTLSX509Settings {
	this := DBUserTLSX509Settings{}
	return &this
}

// GetCas returns the Cas field value if set, zero value otherwise
func (o *DBUserTLSX509Settings) GetCas() string {
	if o == nil || IsNil(o.Cas) {
		var ret string
		return ret
	}
	return *o.Cas
}

// GetCasOk returns a tuple with the Cas field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DBUserTLSX509Settings) GetCasOk() (*string, bool) {
	if o == nil || IsNil(o.Cas) {
		return nil, false
	}

	return o.Cas, true
}

// HasCas returns a boolean if a field has been set.
func (o *DBUserTLSX509Settings) HasCas() bool {
	if o != nil && !IsNil(o.Cas) {
		return true
	}

	return false
}

// SetCas gets a reference to the given string and assigns it to the Cas field.
func (o *DBUserTLSX509Settings) SetCas(v string) {
	o.Cas = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DBUserTLSX509Settings) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DBUserTLSX509Settings) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DBUserTLSX509Settings) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DBUserTLSX509Settings) SetLinks(v []Link) {
	o.Links = &v
}
