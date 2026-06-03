// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BackupComplianceScheduledPolicyItem Specifications for scheduled policy.
type BackupComplianceScheduledPolicyItem struct {
	// Number that indicates the frequency interval for a set of Snapshots. A value of `1` specifies the first instance of the corresponding `frequencyType`.  - In a yearly policy item, `1` indicates that the yearly Snapshot occurs on the first day of January and `12` indicates the first day of December.  - In a monthly policy item, `1` indicates that the monthly Snapshot occurs on the first day of the month and `40` indicates the last day of the month.  - In a weekly policy item, `1` indicates that the weekly Snapshot occurs on Monday and `7` indicates Sunday.  - In an hourly policy item, you can set the frequency interval to `1`, `2`, `4`, `6`, `8`, or `12`. For hourly policy items for NVMe clusters, MongoDB Cloud accepts only `12` as the frequency interval value.   MongoDB Cloud ignores this setting for non-hourly policy items in Backup Compliance Policy settings.
	FrequencyInterval int `json:"frequencyInterval"`
	// Human-readable label that identifies the frequency type associated with the backup policy.
	FrequencyType string `json:"frequencyType"`
	// Unique 24-hexadecimal digit string that identifies this backup policy item.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Unit of time in which MongoDB Cloud measures Snapshot retention.
	RetentionUnit string `json:"retentionUnit"`
	// Duration in days, weeks, months, or years that MongoDB Cloud retains the Snapshot. For less frequent policy items, MongoDB Cloud requires that you specify a value greater than or equal to the value specified for more frequent policy items.  For example: If the hourly policy item specifies a retention of two days, you must specify two days or greater for the retention of the weekly policy item.
	RetentionValue int `json:"retentionValue"`
}

// NewBackupComplianceScheduledPolicyItem instantiates a new BackupComplianceScheduledPolicyItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupComplianceScheduledPolicyItem(frequencyInterval int, frequencyType string, retentionUnit string, retentionValue int) *BackupComplianceScheduledPolicyItem {
	this := BackupComplianceScheduledPolicyItem{}
	this.FrequencyInterval = frequencyInterval
	this.FrequencyType = frequencyType
	this.RetentionUnit = retentionUnit
	this.RetentionValue = retentionValue
	return &this
}

// NewBackupComplianceScheduledPolicyItemWithDefaults instantiates a new BackupComplianceScheduledPolicyItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupComplianceScheduledPolicyItemWithDefaults() *BackupComplianceScheduledPolicyItem {
	this := BackupComplianceScheduledPolicyItem{}
	return &this
}

// GetFrequencyInterval returns the FrequencyInterval field value
func (o *BackupComplianceScheduledPolicyItem) GetFrequencyInterval() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.FrequencyInterval
}

// GetFrequencyIntervalOk returns a tuple with the FrequencyInterval field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceScheduledPolicyItem) GetFrequencyIntervalOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FrequencyInterval, true
}

// SetFrequencyInterval sets field value
func (o *BackupComplianceScheduledPolicyItem) SetFrequencyInterval(v int) {
	o.FrequencyInterval = v
}

// GetFrequencyType returns the FrequencyType field value
func (o *BackupComplianceScheduledPolicyItem) GetFrequencyType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceScheduledPolicyItem) GetFrequencyTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FrequencyType, true
}

// SetFrequencyType sets field value
func (o *BackupComplianceScheduledPolicyItem) SetFrequencyType(v string) {
	o.FrequencyType = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *BackupComplianceScheduledPolicyItem) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupComplianceScheduledPolicyItem) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *BackupComplianceScheduledPolicyItem) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *BackupComplianceScheduledPolicyItem) SetId(v string) {
	o.Id = &v
}

// GetRetentionUnit returns the RetentionUnit field value
func (o *BackupComplianceScheduledPolicyItem) GetRetentionUnit() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RetentionUnit
}

// GetRetentionUnitOk returns a tuple with the RetentionUnit field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceScheduledPolicyItem) GetRetentionUnitOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RetentionUnit, true
}

// SetRetentionUnit sets field value
func (o *BackupComplianceScheduledPolicyItem) SetRetentionUnit(v string) {
	o.RetentionUnit = v
}

// GetRetentionValue returns the RetentionValue field value
func (o *BackupComplianceScheduledPolicyItem) GetRetentionValue() int {
	if o == nil {
		var ret int
		return ret
	}

	return o.RetentionValue
}

// GetRetentionValueOk returns a tuple with the RetentionValue field value
// and a boolean to check if the value has been set.
func (o *BackupComplianceScheduledPolicyItem) GetRetentionValueOk() (*int, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RetentionValue, true
}

// SetRetentionValue sets field value
func (o *BackupComplianceScheduledPolicyItem) SetRetentionValue(v int) {
	o.RetentionValue = v
}
