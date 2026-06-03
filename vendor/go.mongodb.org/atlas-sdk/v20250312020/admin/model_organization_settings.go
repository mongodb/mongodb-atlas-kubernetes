// Code based on the AtlasAPI V2 OpenAPI file

package admin

// OrganizationSettings Collection of settings that configures the organization.
type OrganizationSettings struct {
	// Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
	ApiAccessListRequired *bool                  `json:"apiAccessListRequired,omitempty"`
	CustomSessionTimeouts *CustomSessionTimeouts `json:"customSessionTimeouts,omitempty"`
	// Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and is enabled by default. Once this setting is turned on, Project Owners may be able to enable or disable individual AI features at the project level.
	GenAIFeaturesEnabled *bool `json:"genAIFeaturesEnabled,omitempty"`
	// Number that represents the maximum period before expiry in hours for new Atlas Admin API Service Account secrets within the specified organization.
	MaxServiceAccountSecretValidityInHours *int `json:"maxServiceAccountSecretValidityInHours,omitempty"`
	// Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization. To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
	MultiFactorAuthRequired *bool `json:"multiFactorAuthRequired,omitempty"`
	// Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure and cluster logs for any deployment in the specified organization without explicit permission. Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues. To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.
	RestrictEmployeeAccess *bool `json:"restrictEmployeeAccess,omitempty"`
	// String that specifies a single email address for the specified organization to receive security-related notifications. Specifying a security contact does not grant them authorization or access to Atlas for security decisions or approvals. An empty string is valid and clears the existing security contact (if any).
	SecurityContact *string `json:"securityContact,omitempty"`
	// Flag that indicates whether a group's Atlas Stream Processing workspaces in this organization can create connections to other group's clusters in the same organization.
	StreamsCrossGroupEnabled *bool `json:"streamsCrossGroupEnabled,omitempty"`
}

// NewOrganizationSettings instantiates a new OrganizationSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrganizationSettings() *OrganizationSettings {
	this := OrganizationSettings{}
	var genAIFeaturesEnabled bool = true
	this.GenAIFeaturesEnabled = &genAIFeaturesEnabled
	return &this
}

// NewOrganizationSettingsWithDefaults instantiates a new OrganizationSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrganizationSettingsWithDefaults() *OrganizationSettings {
	this := OrganizationSettings{}
	var genAIFeaturesEnabled bool = true
	this.GenAIFeaturesEnabled = &genAIFeaturesEnabled
	return &this
}

// GetApiAccessListRequired returns the ApiAccessListRequired field value if set, zero value otherwise
func (o *OrganizationSettings) GetApiAccessListRequired() bool {
	if o == nil || IsNil(o.ApiAccessListRequired) {
		var ret bool
		return ret
	}
	return *o.ApiAccessListRequired
}

// GetApiAccessListRequiredOk returns a tuple with the ApiAccessListRequired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetApiAccessListRequiredOk() (*bool, bool) {
	if o == nil || IsNil(o.ApiAccessListRequired) {
		return nil, false
	}

	return o.ApiAccessListRequired, true
}

// HasApiAccessListRequired returns a boolean if a field has been set.
func (o *OrganizationSettings) HasApiAccessListRequired() bool {
	if o != nil && !IsNil(o.ApiAccessListRequired) {
		return true
	}

	return false
}

// SetApiAccessListRequired gets a reference to the given bool and assigns it to the ApiAccessListRequired field.
func (o *OrganizationSettings) SetApiAccessListRequired(v bool) {
	o.ApiAccessListRequired = &v
}

// GetCustomSessionTimeouts returns the CustomSessionTimeouts field value if set, zero value otherwise
func (o *OrganizationSettings) GetCustomSessionTimeouts() CustomSessionTimeouts {
	if o == nil || IsNil(o.CustomSessionTimeouts) {
		var ret CustomSessionTimeouts
		return ret
	}
	return *o.CustomSessionTimeouts
}

// GetCustomSessionTimeoutsOk returns a tuple with the CustomSessionTimeouts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetCustomSessionTimeoutsOk() (*CustomSessionTimeouts, bool) {
	if o == nil || IsNil(o.CustomSessionTimeouts) {
		return nil, false
	}

	return o.CustomSessionTimeouts, true
}

// HasCustomSessionTimeouts returns a boolean if a field has been set.
func (o *OrganizationSettings) HasCustomSessionTimeouts() bool {
	if o != nil && !IsNil(o.CustomSessionTimeouts) {
		return true
	}

	return false
}

// SetCustomSessionTimeouts gets a reference to the given CustomSessionTimeouts and assigns it to the CustomSessionTimeouts field.
func (o *OrganizationSettings) SetCustomSessionTimeouts(v CustomSessionTimeouts) {
	o.CustomSessionTimeouts = &v
}

// GetGenAIFeaturesEnabled returns the GenAIFeaturesEnabled field value if set, zero value otherwise
func (o *OrganizationSettings) GetGenAIFeaturesEnabled() bool {
	if o == nil || IsNil(o.GenAIFeaturesEnabled) {
		var ret bool
		return ret
	}
	return *o.GenAIFeaturesEnabled
}

// GetGenAIFeaturesEnabledOk returns a tuple with the GenAIFeaturesEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetGenAIFeaturesEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.GenAIFeaturesEnabled) {
		return nil, false
	}

	return o.GenAIFeaturesEnabled, true
}

// HasGenAIFeaturesEnabled returns a boolean if a field has been set.
func (o *OrganizationSettings) HasGenAIFeaturesEnabled() bool {
	if o != nil && !IsNil(o.GenAIFeaturesEnabled) {
		return true
	}

	return false
}

// SetGenAIFeaturesEnabled gets a reference to the given bool and assigns it to the GenAIFeaturesEnabled field.
func (o *OrganizationSettings) SetGenAIFeaturesEnabled(v bool) {
	o.GenAIFeaturesEnabled = &v
}

// GetMaxServiceAccountSecretValidityInHours returns the MaxServiceAccountSecretValidityInHours field value if set, zero value otherwise
func (o *OrganizationSettings) GetMaxServiceAccountSecretValidityInHours() int {
	if o == nil || IsNil(o.MaxServiceAccountSecretValidityInHours) {
		var ret int
		return ret
	}
	return *o.MaxServiceAccountSecretValidityInHours
}

// GetMaxServiceAccountSecretValidityInHoursOk returns a tuple with the MaxServiceAccountSecretValidityInHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetMaxServiceAccountSecretValidityInHoursOk() (*int, bool) {
	if o == nil || IsNil(o.MaxServiceAccountSecretValidityInHours) {
		return nil, false
	}

	return o.MaxServiceAccountSecretValidityInHours, true
}

// HasMaxServiceAccountSecretValidityInHours returns a boolean if a field has been set.
func (o *OrganizationSettings) HasMaxServiceAccountSecretValidityInHours() bool {
	if o != nil && !IsNil(o.MaxServiceAccountSecretValidityInHours) {
		return true
	}

	return false
}

// SetMaxServiceAccountSecretValidityInHours gets a reference to the given int and assigns it to the MaxServiceAccountSecretValidityInHours field.
func (o *OrganizationSettings) SetMaxServiceAccountSecretValidityInHours(v int) {
	o.MaxServiceAccountSecretValidityInHours = &v
}

// GetMultiFactorAuthRequired returns the MultiFactorAuthRequired field value if set, zero value otherwise
func (o *OrganizationSettings) GetMultiFactorAuthRequired() bool {
	if o == nil || IsNil(o.MultiFactorAuthRequired) {
		var ret bool
		return ret
	}
	return *o.MultiFactorAuthRequired
}

// GetMultiFactorAuthRequiredOk returns a tuple with the MultiFactorAuthRequired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetMultiFactorAuthRequiredOk() (*bool, bool) {
	if o == nil || IsNil(o.MultiFactorAuthRequired) {
		return nil, false
	}

	return o.MultiFactorAuthRequired, true
}

// HasMultiFactorAuthRequired returns a boolean if a field has been set.
func (o *OrganizationSettings) HasMultiFactorAuthRequired() bool {
	if o != nil && !IsNil(o.MultiFactorAuthRequired) {
		return true
	}

	return false
}

// SetMultiFactorAuthRequired gets a reference to the given bool and assigns it to the MultiFactorAuthRequired field.
func (o *OrganizationSettings) SetMultiFactorAuthRequired(v bool) {
	o.MultiFactorAuthRequired = &v
}

// GetRestrictEmployeeAccess returns the RestrictEmployeeAccess field value if set, zero value otherwise
func (o *OrganizationSettings) GetRestrictEmployeeAccess() bool {
	if o == nil || IsNil(o.RestrictEmployeeAccess) {
		var ret bool
		return ret
	}
	return *o.RestrictEmployeeAccess
}

// GetRestrictEmployeeAccessOk returns a tuple with the RestrictEmployeeAccess field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetRestrictEmployeeAccessOk() (*bool, bool) {
	if o == nil || IsNil(o.RestrictEmployeeAccess) {
		return nil, false
	}

	return o.RestrictEmployeeAccess, true
}

// HasRestrictEmployeeAccess returns a boolean if a field has been set.
func (o *OrganizationSettings) HasRestrictEmployeeAccess() bool {
	if o != nil && !IsNil(o.RestrictEmployeeAccess) {
		return true
	}

	return false
}

// SetRestrictEmployeeAccess gets a reference to the given bool and assigns it to the RestrictEmployeeAccess field.
func (o *OrganizationSettings) SetRestrictEmployeeAccess(v bool) {
	o.RestrictEmployeeAccess = &v
}

// GetSecurityContact returns the SecurityContact field value if set, zero value otherwise
func (o *OrganizationSettings) GetSecurityContact() string {
	if o == nil || IsNil(o.SecurityContact) {
		var ret string
		return ret
	}
	return *o.SecurityContact
}

// GetSecurityContactOk returns a tuple with the SecurityContact field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetSecurityContactOk() (*string, bool) {
	if o == nil || IsNil(o.SecurityContact) {
		return nil, false
	}

	return o.SecurityContact, true
}

// HasSecurityContact returns a boolean if a field has been set.
func (o *OrganizationSettings) HasSecurityContact() bool {
	if o != nil && !IsNil(o.SecurityContact) {
		return true
	}

	return false
}

// SetSecurityContact gets a reference to the given string and assigns it to the SecurityContact field.
func (o *OrganizationSettings) SetSecurityContact(v string) {
	o.SecurityContact = &v
}

// GetStreamsCrossGroupEnabled returns the StreamsCrossGroupEnabled field value if set, zero value otherwise
func (o *OrganizationSettings) GetStreamsCrossGroupEnabled() bool {
	if o == nil || IsNil(o.StreamsCrossGroupEnabled) {
		var ret bool
		return ret
	}
	return *o.StreamsCrossGroupEnabled
}

// GetStreamsCrossGroupEnabledOk returns a tuple with the StreamsCrossGroupEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OrganizationSettings) GetStreamsCrossGroupEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.StreamsCrossGroupEnabled) {
		return nil, false
	}

	return o.StreamsCrossGroupEnabled, true
}

// HasStreamsCrossGroupEnabled returns a boolean if a field has been set.
func (o *OrganizationSettings) HasStreamsCrossGroupEnabled() bool {
	if o != nil && !IsNil(o.StreamsCrossGroupEnabled) {
		return true
	}

	return false
}

// SetStreamsCrossGroupEnabled gets a reference to the given bool and assigns it to the StreamsCrossGroupEnabled field.
func (o *OrganizationSettings) SetStreamsCrossGroupEnabled(v bool) {
	o.StreamsCrossGroupEnabled = &v
}
