// Code based on the AtlasAPI V2 OpenAPI file

package admin

// AdvancedDiskBackupSnapshotSchedulePolicy List that contains a document for each backup policy item in the desired backup policy.
type AdvancedDiskBackupSnapshotSchedulePolicy struct {
	// Unique 24-hexadecimal digit string that identifies this backup policy.
	Id *string `json:"id,omitempty"`
	// List that contains the specifications for one policy.
	PolicyItems []DiskBackupApiPolicyItem `json:"policyItems"`
}

// NewAdvancedDiskBackupSnapshotSchedulePolicy instantiates a new AdvancedDiskBackupSnapshotSchedulePolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAdvancedDiskBackupSnapshotSchedulePolicy(policyItems []DiskBackupApiPolicyItem) *AdvancedDiskBackupSnapshotSchedulePolicy {
	this := AdvancedDiskBackupSnapshotSchedulePolicy{}
	this.PolicyItems = policyItems
	return &this
}

// NewAdvancedDiskBackupSnapshotSchedulePolicyWithDefaults instantiates a new AdvancedDiskBackupSnapshotSchedulePolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAdvancedDiskBackupSnapshotSchedulePolicyWithDefaults() *AdvancedDiskBackupSnapshotSchedulePolicy {
	this := AdvancedDiskBackupSnapshotSchedulePolicy{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) SetId(v string) {
	o.Id = &v
}

// GetPolicyItems returns the PolicyItems field value
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) GetPolicyItems() []DiskBackupApiPolicyItem {
	if o == nil {
		var ret []DiskBackupApiPolicyItem
		return ret
	}

	return o.PolicyItems
}

// GetPolicyItemsOk returns a tuple with the PolicyItems field value
// and a boolean to check if the value has been set.
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) GetPolicyItemsOk() (*[]DiskBackupApiPolicyItem, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PolicyItems, true
}

// SetPolicyItems sets field value
func (o *AdvancedDiskBackupSnapshotSchedulePolicy) SetPolicyItems(v []DiskBackupApiPolicyItem) {
	o.PolicyItems = v
}
