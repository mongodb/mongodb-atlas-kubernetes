// Code based on the AtlasAPI V2 OpenAPI file

package admin

// BackupLabel Collection of key-value pairs that represent custom data to add to the metadata file that MongoDB Cloud uploads to the bucket when the export job finishes.
type BackupLabel struct {
	// Key for the metadata file that MongoDB Cloud uploads to the bucket when the export job finishes.
	Key *string `json:"key,omitempty"`
	// Value for the key to include in file that MongoDB Cloud uploads to the bucket when the export job finishes.
	Value *string `json:"value,omitempty"`
}

// NewBackupLabel instantiates a new BackupLabel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupLabel() *BackupLabel {
	this := BackupLabel{}
	return &this
}

// NewBackupLabelWithDefaults instantiates a new BackupLabel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupLabelWithDefaults() *BackupLabel {
	this := BackupLabel{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise
func (o *BackupLabel) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupLabel) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}

	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *BackupLabel) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *BackupLabel) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise
func (o *BackupLabel) GetValue() string {
	if o == nil || IsNil(o.Value) {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupLabel) GetValueOk() (*string, bool) {
	if o == nil || IsNil(o.Value) {
		return nil, false
	}

	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *BackupLabel) HasValue() bool {
	if o != nil && !IsNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *BackupLabel) SetValue(v string) {
	o.Value = &v
}
