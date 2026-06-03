// Code based on the AtlasAPI V2 OpenAPI file

package admin

// FlexBackupSettings20241113 Flex backup configuration.
type FlexBackupSettings20241113 struct {
	// Flag that indicates whether backups are performed for this flex cluster. Backup uses flex cluster backups.
	// Read only field.
	Enabled *bool `json:"enabled,omitempty"`
}

// NewFlexBackupSettings20241113 instantiates a new FlexBackupSettings20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexBackupSettings20241113() *FlexBackupSettings20241113 {
	this := FlexBackupSettings20241113{}
	return &this
}

// NewFlexBackupSettings20241113WithDefaults instantiates a new FlexBackupSettings20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexBackupSettings20241113WithDefaults() *FlexBackupSettings20241113 {
	this := FlexBackupSettings20241113{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise
func (o *FlexBackupSettings20241113) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupSettings20241113) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}

	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *FlexBackupSettings20241113) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *FlexBackupSettings20241113) SetEnabled(v bool) {
	o.Enabled = &v
}
