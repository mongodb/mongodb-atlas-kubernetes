// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DataProtectionSettings20231001 struct for DataProtectionSettings20231001
type DataProtectionSettings20231001 struct {
	// Email address of the user who authorized to update the Backup Compliance Policy  settings.
	AuthorizedEmail string `json:"authorizedEmail"`
	// First name of the user who authorized to updated the Backup Compliance Policy  settings.
	AuthorizedUserFirstName string `json:"authorizedUserFirstName"`
	// Last name of the user who authorized to updated the Backup Compliance Policy  settings.
	AuthorizedUserLastName string `json:"authorizedUserLastName"`
	// Flag that indicates whether to prevent cluster users from deleting backups copied to other regions, even if those additional snapshot regions are removed. If unspecified, this value defaults to false.
	CopyProtectionEnabled *bool `json:"copyProtectionEnabled,omitempty"`
	// Flag that indicates whether the Backup Compliance Policy is allowed to be disabled. It is default to false and a support ticket needs to be filed to request setting to true.
	// Read only field.
	Deletable *bool `json:"deletable,omitempty"`
	// Flag that indicates whether Encryption at Rest using Customer Key  Management is required for all clusters with a Backup Compliance Policy. If unspecified, this value defaults to false.
	EncryptionAtRestEnabled *bool                               `json:"encryptionAtRestEnabled,omitempty"`
	OnDemandPolicyItem      *BackupComplianceOnDemandPolicyItem `json:"onDemandPolicyItem,omitempty"`
	// Flag that indicates whether the cluster uses Continuous Cloud Backups with a Backup Compliance Policy. If unspecified, this value defaults to false.
	PitEnabled *bool `json:"pitEnabled,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project for the Backup Compliance Policy.
	ProjectId *string `json:"projectId,omitempty"`
	// Number of previous days that you can restore back to with Continuous Cloud Backup with a Backup Compliance Policy. You must specify a positive, non-zero integer, and the maximum retention window can't exceed the hourly retention time. This parameter applies only to Continuous Cloud Backups with a Backup Compliance Policy.
	RestoreWindowDays *int `json:"restoreWindowDays,omitempty"`
	// List that contains the specifications for one scheduled policy.
	ScheduledPolicyItems *[]BackupComplianceScheduledPolicyItem `json:"scheduledPolicyItems,omitempty"`
	// Label that indicates the state of the Backup Compliance Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
	// Read only field.
	State *string `json:"state,omitempty"`
	// ISO 8601 timestamp format in UTC that indicates when the user updated the Data Protection Policy settings. MongoDB Cloud ignores this setting when you enable or update the Backup Compliance Policy settings.
	// Read only field.
	UpdatedDate *time.Time `json:"updatedDate,omitempty"`
	// Email address that identifies the user who updated the Backup Compliance Policy settings. MongoDB Cloud ignores this email setting when you enable or update the Backup Compliance Policy settings.
	// Read only field.
	UpdatedUser *string `json:"updatedUser,omitempty"`
}

// NewDataProtectionSettings20231001 instantiates a new DataProtectionSettings20231001 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataProtectionSettings20231001(authorizedEmail string, authorizedUserFirstName string, authorizedUserLastName string) *DataProtectionSettings20231001 {
	this := DataProtectionSettings20231001{}
	this.AuthorizedEmail = authorizedEmail
	this.AuthorizedUserFirstName = authorizedUserFirstName
	this.AuthorizedUserLastName = authorizedUserLastName
	var copyProtectionEnabled bool = false
	this.CopyProtectionEnabled = &copyProtectionEnabled
	var encryptionAtRestEnabled bool = false
	this.EncryptionAtRestEnabled = &encryptionAtRestEnabled
	var pitEnabled bool = false
	this.PitEnabled = &pitEnabled
	return &this
}

// NewDataProtectionSettings20231001WithDefaults instantiates a new DataProtectionSettings20231001 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataProtectionSettings20231001WithDefaults() *DataProtectionSettings20231001 {
	this := DataProtectionSettings20231001{}
	var copyProtectionEnabled bool = false
	this.CopyProtectionEnabled = &copyProtectionEnabled
	var encryptionAtRestEnabled bool = false
	this.EncryptionAtRestEnabled = &encryptionAtRestEnabled
	var pitEnabled bool = false
	this.PitEnabled = &pitEnabled
	return &this
}

// GetAuthorizedEmail returns the AuthorizedEmail field value
func (o *DataProtectionSettings20231001) GetAuthorizedEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AuthorizedEmail
}

// GetAuthorizedEmailOk returns a tuple with the AuthorizedEmail field value
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetAuthorizedEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthorizedEmail, true
}

// SetAuthorizedEmail sets field value
func (o *DataProtectionSettings20231001) SetAuthorizedEmail(v string) {
	o.AuthorizedEmail = v
}

// GetAuthorizedUserFirstName returns the AuthorizedUserFirstName field value
func (o *DataProtectionSettings20231001) GetAuthorizedUserFirstName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AuthorizedUserFirstName
}

// GetAuthorizedUserFirstNameOk returns a tuple with the AuthorizedUserFirstName field value
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetAuthorizedUserFirstNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthorizedUserFirstName, true
}

// SetAuthorizedUserFirstName sets field value
func (o *DataProtectionSettings20231001) SetAuthorizedUserFirstName(v string) {
	o.AuthorizedUserFirstName = v
}

// GetAuthorizedUserLastName returns the AuthorizedUserLastName field value
func (o *DataProtectionSettings20231001) GetAuthorizedUserLastName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AuthorizedUserLastName
}

// GetAuthorizedUserLastNameOk returns a tuple with the AuthorizedUserLastName field value
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetAuthorizedUserLastNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.AuthorizedUserLastName, true
}

// SetAuthorizedUserLastName sets field value
func (o *DataProtectionSettings20231001) SetAuthorizedUserLastName(v string) {
	o.AuthorizedUserLastName = v
}

// GetCopyProtectionEnabled returns the CopyProtectionEnabled field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetCopyProtectionEnabled() bool {
	if o == nil || IsNil(o.CopyProtectionEnabled) {
		var ret bool
		return ret
	}
	return *o.CopyProtectionEnabled
}

// GetCopyProtectionEnabledOk returns a tuple with the CopyProtectionEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetCopyProtectionEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.CopyProtectionEnabled) {
		return nil, false
	}

	return o.CopyProtectionEnabled, true
}

// HasCopyProtectionEnabled returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasCopyProtectionEnabled() bool {
	if o != nil && !IsNil(o.CopyProtectionEnabled) {
		return true
	}

	return false
}

// SetCopyProtectionEnabled gets a reference to the given bool and assigns it to the CopyProtectionEnabled field.
func (o *DataProtectionSettings20231001) SetCopyProtectionEnabled(v bool) {
	o.CopyProtectionEnabled = &v
}

// GetDeletable returns the Deletable field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetDeletable() bool {
	if o == nil || IsNil(o.Deletable) {
		var ret bool
		return ret
	}
	return *o.Deletable
}

// GetDeletableOk returns a tuple with the Deletable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetDeletableOk() (*bool, bool) {
	if o == nil || IsNil(o.Deletable) {
		return nil, false
	}

	return o.Deletable, true
}

// HasDeletable returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasDeletable() bool {
	if o != nil && !IsNil(o.Deletable) {
		return true
	}

	return false
}

// SetDeletable gets a reference to the given bool and assigns it to the Deletable field.
func (o *DataProtectionSettings20231001) SetDeletable(v bool) {
	o.Deletable = &v
}

// GetEncryptionAtRestEnabled returns the EncryptionAtRestEnabled field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetEncryptionAtRestEnabled() bool {
	if o == nil || IsNil(o.EncryptionAtRestEnabled) {
		var ret bool
		return ret
	}
	return *o.EncryptionAtRestEnabled
}

// GetEncryptionAtRestEnabledOk returns a tuple with the EncryptionAtRestEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetEncryptionAtRestEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.EncryptionAtRestEnabled) {
		return nil, false
	}

	return o.EncryptionAtRestEnabled, true
}

// HasEncryptionAtRestEnabled returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasEncryptionAtRestEnabled() bool {
	if o != nil && !IsNil(o.EncryptionAtRestEnabled) {
		return true
	}

	return false
}

// SetEncryptionAtRestEnabled gets a reference to the given bool and assigns it to the EncryptionAtRestEnabled field.
func (o *DataProtectionSettings20231001) SetEncryptionAtRestEnabled(v bool) {
	o.EncryptionAtRestEnabled = &v
}

// GetOnDemandPolicyItem returns the OnDemandPolicyItem field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetOnDemandPolicyItem() BackupComplianceOnDemandPolicyItem {
	if o == nil || IsNil(o.OnDemandPolicyItem) {
		var ret BackupComplianceOnDemandPolicyItem
		return ret
	}
	return *o.OnDemandPolicyItem
}

// GetOnDemandPolicyItemOk returns a tuple with the OnDemandPolicyItem field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetOnDemandPolicyItemOk() (*BackupComplianceOnDemandPolicyItem, bool) {
	if o == nil || IsNil(o.OnDemandPolicyItem) {
		return nil, false
	}

	return o.OnDemandPolicyItem, true
}

// HasOnDemandPolicyItem returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasOnDemandPolicyItem() bool {
	if o != nil && !IsNil(o.OnDemandPolicyItem) {
		return true
	}

	return false
}

// SetOnDemandPolicyItem gets a reference to the given BackupComplianceOnDemandPolicyItem and assigns it to the OnDemandPolicyItem field.
func (o *DataProtectionSettings20231001) SetOnDemandPolicyItem(v BackupComplianceOnDemandPolicyItem) {
	o.OnDemandPolicyItem = &v
}

// GetPitEnabled returns the PitEnabled field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetPitEnabled() bool {
	if o == nil || IsNil(o.PitEnabled) {
		var ret bool
		return ret
	}
	return *o.PitEnabled
}

// GetPitEnabledOk returns a tuple with the PitEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetPitEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.PitEnabled) {
		return nil, false
	}

	return o.PitEnabled, true
}

// HasPitEnabled returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasPitEnabled() bool {
	if o != nil && !IsNil(o.PitEnabled) {
		return true
	}

	return false
}

// SetPitEnabled gets a reference to the given bool and assigns it to the PitEnabled field.
func (o *DataProtectionSettings20231001) SetPitEnabled(v bool) {
	o.PitEnabled = &v
}

// GetProjectId returns the ProjectId field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetProjectId() string {
	if o == nil || IsNil(o.ProjectId) {
		var ret string
		return ret
	}
	return *o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProjectId) {
		return nil, false
	}

	return o.ProjectId, true
}

// HasProjectId returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasProjectId() bool {
	if o != nil && !IsNil(o.ProjectId) {
		return true
	}

	return false
}

// SetProjectId gets a reference to the given string and assigns it to the ProjectId field.
func (o *DataProtectionSettings20231001) SetProjectId(v string) {
	o.ProjectId = &v
}

// GetRestoreWindowDays returns the RestoreWindowDays field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetRestoreWindowDays() int {
	if o == nil || IsNil(o.RestoreWindowDays) {
		var ret int
		return ret
	}
	return *o.RestoreWindowDays
}

// GetRestoreWindowDaysOk returns a tuple with the RestoreWindowDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetRestoreWindowDaysOk() (*int, bool) {
	if o == nil || IsNil(o.RestoreWindowDays) {
		return nil, false
	}

	return o.RestoreWindowDays, true
}

// HasRestoreWindowDays returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasRestoreWindowDays() bool {
	if o != nil && !IsNil(o.RestoreWindowDays) {
		return true
	}

	return false
}

// SetRestoreWindowDays gets a reference to the given int and assigns it to the RestoreWindowDays field.
func (o *DataProtectionSettings20231001) SetRestoreWindowDays(v int) {
	o.RestoreWindowDays = &v
}

// GetScheduledPolicyItems returns the ScheduledPolicyItems field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetScheduledPolicyItems() []BackupComplianceScheduledPolicyItem {
	if o == nil || IsNil(o.ScheduledPolicyItems) {
		var ret []BackupComplianceScheduledPolicyItem
		return ret
	}
	return *o.ScheduledPolicyItems
}

// GetScheduledPolicyItemsOk returns a tuple with the ScheduledPolicyItems field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetScheduledPolicyItemsOk() (*[]BackupComplianceScheduledPolicyItem, bool) {
	if o == nil || IsNil(o.ScheduledPolicyItems) {
		return nil, false
	}

	return o.ScheduledPolicyItems, true
}

// HasScheduledPolicyItems returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasScheduledPolicyItems() bool {
	if o != nil && !IsNil(o.ScheduledPolicyItems) {
		return true
	}

	return false
}

// SetScheduledPolicyItems gets a reference to the given []BackupComplianceScheduledPolicyItem and assigns it to the ScheduledPolicyItems field.
func (o *DataProtectionSettings20231001) SetScheduledPolicyItems(v []BackupComplianceScheduledPolicyItem) {
	o.ScheduledPolicyItems = &v
}

// GetState returns the State field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}

	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *DataProtectionSettings20231001) SetState(v string) {
	o.State = &v
}

// GetUpdatedDate returns the UpdatedDate field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetUpdatedDate() time.Time {
	if o == nil || IsNil(o.UpdatedDate) {
		var ret time.Time
		return ret
	}
	return *o.UpdatedDate
}

// GetUpdatedDateOk returns a tuple with the UpdatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetUpdatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdatedDate) {
		return nil, false
	}

	return o.UpdatedDate, true
}

// HasUpdatedDate returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasUpdatedDate() bool {
	if o != nil && !IsNil(o.UpdatedDate) {
		return true
	}

	return false
}

// SetUpdatedDate gets a reference to the given time.Time and assigns it to the UpdatedDate field.
func (o *DataProtectionSettings20231001) SetUpdatedDate(v time.Time) {
	o.UpdatedDate = &v
}

// GetUpdatedUser returns the UpdatedUser field value if set, zero value otherwise
func (o *DataProtectionSettings20231001) GetUpdatedUser() string {
	if o == nil || IsNil(o.UpdatedUser) {
		var ret string
		return ret
	}
	return *o.UpdatedUser
}

// GetUpdatedUserOk returns a tuple with the UpdatedUser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataProtectionSettings20231001) GetUpdatedUserOk() (*string, bool) {
	if o == nil || IsNil(o.UpdatedUser) {
		return nil, false
	}

	return o.UpdatedUser, true
}

// HasUpdatedUser returns a boolean if a field has been set.
func (o *DataProtectionSettings20231001) HasUpdatedUser() bool {
	if o != nil && !IsNil(o.UpdatedUser) {
		return true
	}

	return false
}

// SetUpdatedUser gets a reference to the given string and assigns it to the UpdatedUser field.
func (o *DataProtectionSettings20231001) SetUpdatedUser(v string) {
	o.UpdatedUser = &v
}
