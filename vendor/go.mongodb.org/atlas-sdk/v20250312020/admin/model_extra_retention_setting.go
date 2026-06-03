// Code based on the AtlasAPI V2 OpenAPI file

package admin

// ExtraRetentionSetting Extra retention setting item in the desired backup policy.
type ExtraRetentionSetting struct {
	// The frequency type for the extra retention settings for the cluster.
	FrequencyType *string `json:"frequencyType,omitempty"`
	// The number of extra retention days for the cluster.
	RetentionDays *int `json:"retentionDays,omitempty"`
}

// NewExtraRetentionSetting instantiates a new ExtraRetentionSetting object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExtraRetentionSetting() *ExtraRetentionSetting {
	this := ExtraRetentionSetting{}
	return &this
}

// NewExtraRetentionSettingWithDefaults instantiates a new ExtraRetentionSetting object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExtraRetentionSettingWithDefaults() *ExtraRetentionSetting {
	this := ExtraRetentionSetting{}
	return &this
}

// GetFrequencyType returns the FrequencyType field value if set, zero value otherwise
func (o *ExtraRetentionSetting) GetFrequencyType() string {
	if o == nil || IsNil(o.FrequencyType) {
		var ret string
		return ret
	}
	return *o.FrequencyType
}

// GetFrequencyTypeOk returns a tuple with the FrequencyType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtraRetentionSetting) GetFrequencyTypeOk() (*string, bool) {
	if o == nil || IsNil(o.FrequencyType) {
		return nil, false
	}

	return o.FrequencyType, true
}

// HasFrequencyType returns a boolean if a field has been set.
func (o *ExtraRetentionSetting) HasFrequencyType() bool {
	if o != nil && !IsNil(o.FrequencyType) {
		return true
	}

	return false
}

// SetFrequencyType gets a reference to the given string and assigns it to the FrequencyType field.
func (o *ExtraRetentionSetting) SetFrequencyType(v string) {
	o.FrequencyType = &v
}

// GetRetentionDays returns the RetentionDays field value if set, zero value otherwise
func (o *ExtraRetentionSetting) GetRetentionDays() int {
	if o == nil || IsNil(o.RetentionDays) {
		var ret int
		return ret
	}
	return *o.RetentionDays
}

// GetRetentionDaysOk returns a tuple with the RetentionDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtraRetentionSetting) GetRetentionDaysOk() (*int, bool) {
	if o == nil || IsNil(o.RetentionDays) {
		return nil, false
	}

	return o.RetentionDays, true
}

// HasRetentionDays returns a boolean if a field has been set.
func (o *ExtraRetentionSetting) HasRetentionDays() bool {
	if o != nil && !IsNil(o.RetentionDays) {
		return true
	}

	return false
}

// SetRetentionDays gets a reference to the given int and assigns it to the RetentionDays field.
func (o *ExtraRetentionSetting) SetRetentionDays(v int) {
	o.RetentionDays = &v
}
