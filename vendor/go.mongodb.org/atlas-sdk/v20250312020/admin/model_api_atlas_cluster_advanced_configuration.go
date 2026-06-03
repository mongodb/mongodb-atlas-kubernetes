// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ApiAtlasClusterAdvancedConfiguration Group of settings that configures a subset of the advanced configuration details.
type ApiAtlasClusterAdvancedConfiguration struct {
	// The custom OpenSSL cipher suite list for TLS 1.2. This field is only valid when `tlsCipherConfigMode` is set to `CUSTOM`.
	CustomOpensslCipherConfigTls12 *[]string `json:"customOpensslCipherConfigTls12,omitempty"`
	// The custom OpenSSL cipher suite list for TLS 1.3. This field is only valid when `tlsCipherConfigMode` is set to `CUSTOM`.
	CustomOpensslCipherConfigTls13 *[]string `json:"customOpensslCipherConfigTls13,omitempty"`
	// Minimum Transport Layer Security (TLS) version that the cluster accepts for incoming connections. Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.
	MinimumEnabledTlsProtocol *string `json:"minimumEnabledTlsProtocol,omitempty"`
	// The TLS cipher suite configuration mode. The default mode uses the default cipher suites. The custom mode allows you to specify custom cipher suites for both TLS 1.2 and TLS 1.3.
	TlsCipherConfigMode *string `json:"tlsCipherConfigMode,omitempty"`
}

// NewApiAtlasClusterAdvancedConfiguration instantiates a new ApiAtlasClusterAdvancedConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAtlasClusterAdvancedConfiguration() *ApiAtlasClusterAdvancedConfiguration {
	this := ApiAtlasClusterAdvancedConfiguration{}
	return &this
}

// NewApiAtlasClusterAdvancedConfigurationWithDefaults instantiates a new ApiAtlasClusterAdvancedConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAtlasClusterAdvancedConfigurationWithDefaults() *ApiAtlasClusterAdvancedConfiguration {
	this := ApiAtlasClusterAdvancedConfiguration{}
	return &this
}

// GetCustomOpensslCipherConfigTls12 returns the CustomOpensslCipherConfigTls12 field value if set, zero value otherwise
func (o *ApiAtlasClusterAdvancedConfiguration) GetCustomOpensslCipherConfigTls12() []string {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls12) {
		var ret []string
		return ret
	}
	return *o.CustomOpensslCipherConfigTls12
}

// GetCustomOpensslCipherConfigTls12Ok returns a tuple with the CustomOpensslCipherConfigTls12 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) GetCustomOpensslCipherConfigTls12Ok() (*[]string, bool) {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls12) {
		return nil, false
	}

	return o.CustomOpensslCipherConfigTls12, true
}

// HasCustomOpensslCipherConfigTls12 returns a boolean if a field has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) HasCustomOpensslCipherConfigTls12() bool {
	if o != nil && !IsNil(o.CustomOpensslCipherConfigTls12) {
		return true
	}

	return false
}

// SetCustomOpensslCipherConfigTls12 gets a reference to the given []string and assigns it to the CustomOpensslCipherConfigTls12 field.
func (o *ApiAtlasClusterAdvancedConfiguration) SetCustomOpensslCipherConfigTls12(v []string) {
	o.CustomOpensslCipherConfigTls12 = &v
}

// GetCustomOpensslCipherConfigTls13 returns the CustomOpensslCipherConfigTls13 field value if set, zero value otherwise
func (o *ApiAtlasClusterAdvancedConfiguration) GetCustomOpensslCipherConfigTls13() []string {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls13) {
		var ret []string
		return ret
	}
	return *o.CustomOpensslCipherConfigTls13
}

// GetCustomOpensslCipherConfigTls13Ok returns a tuple with the CustomOpensslCipherConfigTls13 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) GetCustomOpensslCipherConfigTls13Ok() (*[]string, bool) {
	if o == nil || IsNil(o.CustomOpensslCipherConfigTls13) {
		return nil, false
	}

	return o.CustomOpensslCipherConfigTls13, true
}

// HasCustomOpensslCipherConfigTls13 returns a boolean if a field has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) HasCustomOpensslCipherConfigTls13() bool {
	if o != nil && !IsNil(o.CustomOpensslCipherConfigTls13) {
		return true
	}

	return false
}

// SetCustomOpensslCipherConfigTls13 gets a reference to the given []string and assigns it to the CustomOpensslCipherConfigTls13 field.
func (o *ApiAtlasClusterAdvancedConfiguration) SetCustomOpensslCipherConfigTls13(v []string) {
	o.CustomOpensslCipherConfigTls13 = &v
}

// GetMinimumEnabledTlsProtocol returns the MinimumEnabledTlsProtocol field value if set, zero value otherwise
func (o *ApiAtlasClusterAdvancedConfiguration) GetMinimumEnabledTlsProtocol() string {
	if o == nil || IsNil(o.MinimumEnabledTlsProtocol) {
		var ret string
		return ret
	}
	return *o.MinimumEnabledTlsProtocol
}

// GetMinimumEnabledTlsProtocolOk returns a tuple with the MinimumEnabledTlsProtocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) GetMinimumEnabledTlsProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.MinimumEnabledTlsProtocol) {
		return nil, false
	}

	return o.MinimumEnabledTlsProtocol, true
}

// HasMinimumEnabledTlsProtocol returns a boolean if a field has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) HasMinimumEnabledTlsProtocol() bool {
	if o != nil && !IsNil(o.MinimumEnabledTlsProtocol) {
		return true
	}

	return false
}

// SetMinimumEnabledTlsProtocol gets a reference to the given string and assigns it to the MinimumEnabledTlsProtocol field.
func (o *ApiAtlasClusterAdvancedConfiguration) SetMinimumEnabledTlsProtocol(v string) {
	o.MinimumEnabledTlsProtocol = &v
}

// GetTlsCipherConfigMode returns the TlsCipherConfigMode field value if set, zero value otherwise
func (o *ApiAtlasClusterAdvancedConfiguration) GetTlsCipherConfigMode() string {
	if o == nil || IsNil(o.TlsCipherConfigMode) {
		var ret string
		return ret
	}
	return *o.TlsCipherConfigMode
}

// GetTlsCipherConfigModeOk returns a tuple with the TlsCipherConfigMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) GetTlsCipherConfigModeOk() (*string, bool) {
	if o == nil || IsNil(o.TlsCipherConfigMode) {
		return nil, false
	}

	return o.TlsCipherConfigMode, true
}

// HasTlsCipherConfigMode returns a boolean if a field has been set.
func (o *ApiAtlasClusterAdvancedConfiguration) HasTlsCipherConfigMode() bool {
	if o != nil && !IsNil(o.TlsCipherConfigMode) {
		return true
	}

	return false
}

// SetTlsCipherConfigMode gets a reference to the given string and assigns it to the TlsCipherConfigMode field.
func (o *ApiAtlasClusterAdvancedConfiguration) SetTlsCipherConfigMode(v string) {
	o.TlsCipherConfigMode = &v
}
