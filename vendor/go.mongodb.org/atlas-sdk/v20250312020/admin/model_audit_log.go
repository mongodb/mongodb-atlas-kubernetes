// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AuditLog struct for AuditLog
type AuditLog struct {
	// Flag that indicates whether someone set auditing to track successful authentications. This only applies to the `\"atype\" : \"authCheck\"` audit filter. Setting this parameter to `true` degrades cluster performance.
	AuditAuthorizationSuccess *bool `json:"auditAuthorizationSuccess,omitempty"`
	// JSON document that specifies which events to record. Escape any characters that may prevent parsing, such as single or double quotes, using a backslash (`\\`).
	AuditFilter *string `json:"auditFilter,omitempty"`
	// Human-readable label that displays how to configure the audit filter.
	// Read only field.
	ConfigurationType *string `json:"configurationType,omitempty"`
	// Flag that indicates whether someone enabled database auditing for the specified project.
	Enabled *bool `json:"enabled,omitempty"`
}

// NewAuditLog instantiates a new AuditLog object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuditLog() *AuditLog {
	this := AuditLog{}
	var auditAuthorizationSuccess bool = false
	this.AuditAuthorizationSuccess = &auditAuthorizationSuccess
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// NewAuditLogWithDefaults instantiates a new AuditLog object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuditLogWithDefaults() *AuditLog {
	this := AuditLog{}
	var auditAuthorizationSuccess bool = false
	this.AuditAuthorizationSuccess = &auditAuthorizationSuccess
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// GetAuditAuthorizationSuccess returns the AuditAuthorizationSuccess field value if set, zero value otherwise
func (o *AuditLog) GetAuditAuthorizationSuccess() bool {
	if o == nil || IsNil(o.AuditAuthorizationSuccess) {
		var ret bool
		return ret
	}
	return *o.AuditAuthorizationSuccess
}

// GetAuditAuthorizationSuccessOk returns a tuple with the AuditAuthorizationSuccess field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLog) GetAuditAuthorizationSuccessOk() (*bool, bool) {
	if o == nil || IsNil(o.AuditAuthorizationSuccess) {
		return nil, false
	}

	return o.AuditAuthorizationSuccess, true
}

// HasAuditAuthorizationSuccess returns a boolean if a field has been set.
func (o *AuditLog) HasAuditAuthorizationSuccess() bool {
	if o != nil && !IsNil(o.AuditAuthorizationSuccess) {
		return true
	}

	return false
}

// SetAuditAuthorizationSuccess gets a reference to the given bool and assigns it to the AuditAuthorizationSuccess field.
func (o *AuditLog) SetAuditAuthorizationSuccess(v bool) {
	o.AuditAuthorizationSuccess = &v
}

// GetAuditFilter returns the AuditFilter field value if set, zero value otherwise
func (o *AuditLog) GetAuditFilter() string {
	if o == nil || IsNil(o.AuditFilter) {
		var ret string
		return ret
	}
	return *o.AuditFilter
}

// GetAuditFilterOk returns a tuple with the AuditFilter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLog) GetAuditFilterOk() (*string, bool) {
	if o == nil || IsNil(o.AuditFilter) {
		return nil, false
	}

	return o.AuditFilter, true
}

// HasAuditFilter returns a boolean if a field has been set.
func (o *AuditLog) HasAuditFilter() bool {
	if o != nil && !IsNil(o.AuditFilter) {
		return true
	}

	return false
}

// SetAuditFilter gets a reference to the given string and assigns it to the AuditFilter field.
func (o *AuditLog) SetAuditFilter(v string) {
	o.AuditFilter = &v
}

// GetConfigurationType returns the ConfigurationType field value if set, zero value otherwise
func (o *AuditLog) GetConfigurationType() string {
	if o == nil || IsNil(o.ConfigurationType) {
		var ret string
		return ret
	}
	return *o.ConfigurationType
}

// GetConfigurationTypeOk returns a tuple with the ConfigurationType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLog) GetConfigurationTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ConfigurationType) {
		return nil, false
	}

	return o.ConfigurationType, true
}

// HasConfigurationType returns a boolean if a field has been set.
func (o *AuditLog) HasConfigurationType() bool {
	if o != nil && !IsNil(o.ConfigurationType) {
		return true
	}

	return false
}

// SetConfigurationType gets a reference to the given string and assigns it to the ConfigurationType field.
func (o *AuditLog) SetConfigurationType(v string) {
	o.ConfigurationType = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *AuditLog) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuditLog) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *AuditLog) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *AuditLog) SetEnabled(v bool) {
	o.Enabled = &v
}
