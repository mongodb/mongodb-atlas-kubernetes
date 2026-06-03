// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BackupComplianceOnDemandPolicyItem Specifications for on-demand policy.
type BackupComplianceOnDemandPolicyItem struct {
	// Number that indicates the frequency interval for a set of snapshots. MongoDB Cloud ignores this setting for non-hourly policy items in Backup Compliance Policy settings.
	FrequencyInterval int `json:"frequencyInterval"`
	// Human-readable label that identifies the frequency type associated with the backup policy.
	FrequencyType string `json:"frequencyType"`
	// Unique 24-hexadecimal digit string that identifies this backup policy item.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Unit of time in which MongoDB Cloud measures snapshot retention.
	RetentionUnit string `json:"retentionUnit"`
	// Duration in days, weeks, months, or years that MongoDB Cloud retains the snapshot. For less frequent policy items, MongoDB Cloud requires that you specify a value greater than or equal to the value specified for more frequent policy items.  For example: If the hourly policy item specifies a retention of two days, you must specify two days or greater for the retention of the weekly policy item.
	RetentionValue int `json:"retentionValue"`
}

// NewBackupComplianceOnDemandPolicyItem instantiates a new BackupComplianceOnDemandPolicyItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupComplianceOnDemandPolicyItem(frequencyInterval int, frequencyType string, retentionUnit string, retentionValue int) *BackupComplianceOnDemandPolicyItem {
	this := BackupComplianceOnDemandPolicyItem{}
	this.FrequencyInterval = frequencyInterval
	this.FrequencyType = frequencyType
	this.RetentionUnit = retentionUnit
	this.RetentionValue = retentionValue
	return &this
}

// NewBackupComplianceOnDemandPolicyItemWithDefaults instantiates a new BackupComplianceOnDemandPolicyItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupComplianceOnDemandPolicyItemWithDefaults() *BackupComplianceOnDemandPolicyItem {
	this := BackupComplianceOnDemandPolicyItem{}
	return &this
}

// GetFrequencyInterval returns the FrequencyInterval field value
func (o *BackupComplianceOnDemandPolicyItem) GetFrequencyInterval() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.FrequencyInterval
}

// GetFrequencyIntervalOk returns a tuple with the FrequencyInterval field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceOnDemandPolicyItem) GetFrequencyIntervalOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FrequencyInterval, true
}

// SetFrequencyInterval sets field value
func (o *BackupComplianceOnDemandPolicyItem) SetFrequencyInterval(v int) {
	o.FrequencyInterval = v
}

// GetFrequencyType returns the FrequencyType field value
func (o *BackupComplianceOnDemandPolicyItem) GetFrequencyType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceOnDemandPolicyItem) GetFrequencyTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FrequencyType, true
}

// SetFrequencyType sets field value
func (o *BackupComplianceOnDemandPolicyItem) SetFrequencyType(v string) {
	o.FrequencyType = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupComplianceOnDemandPolicyItem) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupComplianceOnDemandPolicyItem) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupComplianceOnDemandPolicyItem) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupComplianceOnDemandPolicyItem) SetId(v string) {
	o.Id = &v
}

// GetRetentionUnit returns the RetentionUnit field value
func (o *BackupComplianceOnDemandPolicyItem) GetRetentionUnit() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RetentionUnit
}

// GetRetentionUnitOk returns a tuple with the RetentionUnit field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceOnDemandPolicyItem) GetRetentionUnitOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RetentionUnit, true
}

// SetRetentionUnit sets field value
func (o *BackupComplianceOnDemandPolicyItem) SetRetentionUnit(v string) {
	o.RetentionUnit = v
}

// GetRetentionValue returns the RetentionValue field value
func (o *BackupComplianceOnDemandPolicyItem) GetRetentionValue() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.RetentionValue
}

// GetRetentionValueOk returns a tuple with the RetentionValue field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceOnDemandPolicyItem) GetRetentionValueOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RetentionValue, true
}

// SetRetentionValue sets field value
func (o *BackupComplianceOnDemandPolicyItem) SetRetentionValue(v int) {
	o.RetentionValue = v
}
