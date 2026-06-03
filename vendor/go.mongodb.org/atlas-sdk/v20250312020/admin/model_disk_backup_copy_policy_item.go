// Code based on the AtlasAPI V2 OpenAPI file

package admin

// DiskBackupCopyPolicyItem Specifications for one copy policy item.
type DiskBackupCopyPolicyItem struct {
	// Human-readable label that identifies the frequency type associated with the copy policy.
	FrequencyType string `json:"frequencyType"`
	// Unique 24-hexadecimal digit string that identifies this copy policy item.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Unit of time in which MongoDB Cloud measures snapshot copy retention.
	RetentionUnit *string `json:"retentionUnit,omitempty"`
	// Duration in days, weeks, months, or years that MongoDB Cloud retains the snapshot copy.
	RetentionValue *int `json:"retentionValue,omitempty"`
}

// NewDiskBackupCopyPolicyItem instantiates a new DiskBackupCopyPolicyItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupCopyPolicyItem(frequencyType string) *DiskBackupCopyPolicyItem {
	this := DiskBackupCopyPolicyItem{}
	this.FrequencyType = frequencyType
	return &this
}

// NewDiskBackupCopyPolicyItemWithDefaults instantiates a new DiskBackupCopyPolicyItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupCopyPolicyItemWithDefaults() *DiskBackupCopyPolicyItem {
	this := DiskBackupCopyPolicyItem{}
	return &this
}

// GetFrequencyType returns the FrequencyType field value
func (o *DiskBackupCopyPolicyItem) GetFrequencyType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value
// and a boolean to check if the value has been set.
func (o *DiskBackupCopyPolicyItem) GetFrequencyTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.FrequencyType, true
}

// SetFrequencyType sets field value
func (o *DiskBackupCopyPolicyItem) SetFrequencyType(v string) {
	o.FrequencyType = v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *DiskBackupCopyPolicyItem) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupCopyPolicyItem) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DiskBackupCopyPolicyItem) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DiskBackupCopyPolicyItem) SetId(v string) {
	o.Id = &v
}

// GetRetentionUnit returns the RetentionUnit field value if set, zero value otherwise
func (o *DiskBackupCopyPolicyItem) GetRetentionUnit() string {
	if o == nil || IsNil(o.RetentionUnit) {
		var ret string
		return ret
	}
	return *o.RetentionUnit
}

// GetRetentionUnitOk returns a tuple with the RetentionUnit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupCopyPolicyItem) GetRetentionUnitOk() (*string, bool) {
	if o == nil || IsNil(o.RetentionUnit) {
		return nil, false
	}

	return o.RetentionUnit, true
}

// HasRetentionUnit returns a boolean if a field has been set.
func (o *DiskBackupCopyPolicyItem) HasRetentionUnit() bool {
	if o != nil && !IsNil(o.RetentionUnit) {
		return true
	}

	return false
}

// SetRetentionUnit gets a reference to the given string and assigns it to the RetentionUnit field.
func (o *DiskBackupCopyPolicyItem) SetRetentionUnit(v string) {
	o.RetentionUnit = &v
}

// GetRetentionValue returns the RetentionValue field value if set, zero value otherwise
func (o *DiskBackupCopyPolicyItem) GetRetentionValue() int {
	if o == nil || IsNil(o.RetentionValue) {
		var ret int
		return ret
	}
	return *o.RetentionValue
}

// GetRetentionValueOk returns a tuple with the RetentionValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupCopyPolicyItem) GetRetentionValueOk() (*int, bool) {
	if o == nil || IsNil(o.RetentionValue) {
		return nil, false
	}

	return o.RetentionValue, true
}

// HasRetentionValue returns a boolean if a field has been set.
func (o *DiskBackupCopyPolicyItem) HasRetentionValue() bool {
	if o != nil && !IsNil(o.RetentionValue) {
		return true
	}

	return false
}

// SetRetentionValue gets a reference to the given int and assigns it to the RetentionValue field.
func (o *DiskBackupCopyPolicyItem) SetRetentionValue(v int) {
	o.RetentionValue = &v
}
