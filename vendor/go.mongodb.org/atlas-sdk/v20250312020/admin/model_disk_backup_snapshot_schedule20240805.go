// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DiskBackupSnapshotSchedule20240805 struct for DiskBackupSnapshotSchedule20240805
type DiskBackupSnapshotSchedule20240805 struct {
	// Flag that indicates whether the copy settings are automatically managed by MongoDB Cloud and sync to the cluster topology.
	AutoCopySettingsEnabled *bool `json:"autoCopySettingsEnabled,omitempty"`
	// Flag that indicates whether MongoDB Cloud automatically exports Cloud Backup Snapshots to the Export Bucket.
	AutoExportEnabled *bool `json:"autoExportEnabled,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the cluster with the Snapshot you want to return.
	// Read only field.
	ClusterId *string `json:"clusterId,omitempty"`
	// Human-readable label that identifies the cluster with the Snapshot you want to return.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Flag that indicates whether copy settings use `copyPolicyItems` instead of `frequencies`. When true, requests must supply `copyPolicyItems` and responses return `copyPolicyItems` only. When false or omitted, requests must supply `frequencies` and responses return `frequencies` only.
	CopyPolicyItemsEnabled *bool `json:"copyPolicyItemsEnabled,omitempty"`
	// List that contains a document for each copy setting item in the desired backup policy.
	CopySettings *[]DiskBackupCopySetting20240805 `json:"copySettings,omitempty"`
	// List that contains a document for each deleted copy setting whose backup copies you want to delete.
	// Write only field.
	DeleteCopiedBackups *[]DeleteCopiedBackups20240805 `json:"deleteCopiedBackups,omitempty"`
	// Flag that indicates whether to delete Snapshot copies that MongoDB Cloud took previously when their associated `copyPolicyItems` are removed from a `copySetting`. This option requires `copyPolicyItemsEnabled` to be true.
	// Write only field.
	DeleteCopySnapshots *bool `json:"deleteCopySnapshots,omitempty"`
	// Flag that indicates whether to delete Snapshots that MongoDB Cloud took previously when deleting the associated backup policy.
	// Write only field.
	DeleteSnapshots *bool             `json:"deleteSnapshots,omitempty"`
	Export          *AutoExportPolicy `json:"export,omitempty"`
	// List that contains a document for each extra retention setting item in the desired backup policy.
	ExtraRetentionSettings *[]ExtraRetentionSetting `json:"extraRetentionSettings,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Date and time when MongoDB Cloud takes the next Snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	NextSnapshot *time.Time `json:"nextSnapshot,omitempty"`
	// Rules set for this backup schedule.
	Policies []AdvancedDiskBackupSnapshotSchedulePolicy `json:"policies"`
	// Hour of day in Coordinated Universal Time (UTC) that represents when MongoDB Cloud takes the Snapshot.
	ReferenceHourOfDay *int `json:"referenceHourOfDay,omitempty"`
	// Minute of the `referenceHourOfDay` that represents when MongoDB Cloud takes the Snapshot.
	ReferenceMinuteOfHour *int `json:"referenceMinuteOfHour,omitempty"`
	// Number of previous days that you can restore back to with Continuous Cloud Backup accuracy. You must specify a positive, non-zero integer. This parameter applies to continuous Cloud Backups only.
	RestoreWindowDays *int `json:"restoreWindowDays,omitempty"`
	// Flag that indicates whether to apply the retention changes for updated copy policy items to Snapshot copies that MongoDB Cloud took previously.
	// Write only field.
	UpdateCopySnapshots *bool `json:"updateCopySnapshots,omitempty"`
	// Flag that indicates whether to apply the retention changes in the updated backup policy to Snapshots that MongoDB Cloud took previously.
	// Write only field.
	UpdateSnapshots *bool `json:"updateSnapshots,omitempty"`
	// Flag that indicates whether to use organization and project names instead of organization and project UUIDs in the path to the metadata files that MongoDB Cloud uploads to your Export Bucket.
	UseOrgAndGroupNamesInExportPrefix *bool `json:"useOrgAndGroupNamesInExportPrefix,omitempty"`
}

// NewDiskBackupSnapshotSchedule20240805 instantiates a new DiskBackupSnapshotSchedule20240805 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshotSchedule20240805(policies []AdvancedDiskBackupSnapshotSchedulePolicy) *DiskBackupSnapshotSchedule20240805 {
	this := DiskBackupSnapshotSchedule20240805{}
	this.Policies = policies
	return &this
}

// NewDiskBackupSnapshotSchedule20240805WithDefaults instantiates a new DiskBackupSnapshotSchedule20240805 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotSchedule20240805WithDefaults() *DiskBackupSnapshotSchedule20240805 {
	this := DiskBackupSnapshotSchedule20240805{}
	return &this
}

// GetAutoCopySettingsEnabled returns the AutoCopySettingsEnabled field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetAutoCopySettingsEnabled() bool {
	if o == nil || IsNil(o.AutoCopySettingsEnabled) {
		var ret bool
		return ret
	}
	return *o.AutoCopySettingsEnabled
}

// GetAutoCopySettingsEnabledOk returns a tuple with the AutoCopySettingsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetAutoCopySettingsEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AutoCopySettingsEnabled) {
		return nil, false
	}

	return o.AutoCopySettingsEnabled, true
}

// HasAutoCopySettingsEnabled returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasAutoCopySettingsEnabled() bool {
	if o != nil && !IsNil(o.AutoCopySettingsEnabled) {
		return true
	}

	return false
}

// SetAutoCopySettingsEnabled gets a reference to the given bool and assigns it to the AutoCopySettingsEnabled field.
func (o *DiskBackupSnapshotSchedule20240805) SetAutoCopySettingsEnabled(v bool) {
	o.AutoCopySettingsEnabled = &v
}

// GetAutoExportEnabled returns the AutoExportEnabled field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetAutoExportEnabled() bool {
	if o == nil || IsNil(o.AutoExportEnabled) {
		var ret bool
		return ret
	}
	return *o.AutoExportEnabled
}

// GetAutoExportEnabledOk returns a tuple with the AutoExportEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetAutoExportEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.AutoExportEnabled) {
		return nil, false
	}

	return o.AutoExportEnabled, true
}

// HasAutoExportEnabled returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasAutoExportEnabled() bool {
	if o != nil && !IsNil(o.AutoExportEnabled) {
		return true
	}

	return false
}

// SetAutoExportEnabled gets a reference to the given bool and assigns it to the AutoExportEnabled field.
func (o *DiskBackupSnapshotSchedule20240805) SetAutoExportEnabled(v bool) {
	o.AutoExportEnabled = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}

	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *DiskBackupSnapshotSchedule20240805) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *DiskBackupSnapshotSchedule20240805) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetCopyPolicyItemsEnabled returns the CopyPolicyItemsEnabled field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetCopyPolicyItemsEnabled() bool {
	if o == nil || IsNil(o.CopyPolicyItemsEnabled) {
		var ret bool
		return ret
	}
	return *o.CopyPolicyItemsEnabled
}

// GetCopyPolicyItemsEnabledOk returns a tuple with the CopyPolicyItemsEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetCopyPolicyItemsEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.CopyPolicyItemsEnabled) {
		return nil, false
	}

	return o.CopyPolicyItemsEnabled, true
}

// HasCopyPolicyItemsEnabled returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasCopyPolicyItemsEnabled() bool {
	if o != nil && !IsNil(o.CopyPolicyItemsEnabled) {
		return true
	}

	return false
}

// SetCopyPolicyItemsEnabled gets a reference to the given bool and assigns it to the CopyPolicyItemsEnabled field.
func (o *DiskBackupSnapshotSchedule20240805) SetCopyPolicyItemsEnabled(v bool) {
	o.CopyPolicyItemsEnabled = &v
}

// GetCopySettings returns the CopySettings field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetCopySettings() []DiskBackupCopySetting20240805 {
	if o == nil || IsNil(o.CopySettings) {
		var ret []DiskBackupCopySetting20240805
		return ret
	}
	return *o.CopySettings
}

// GetCopySettingsOk returns a tuple with the CopySettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetCopySettingsOk() (*[]DiskBackupCopySetting20240805, bool) {
	if o == nil || IsNil(o.CopySettings) {
		return nil, false
	}

	return o.CopySettings, true
}

// HasCopySettings returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasCopySettings() bool {
	if o != nil && !IsNil(o.CopySettings) {
		return true
	}

	return false
}

// SetCopySettings gets a reference to the given []DiskBackupCopySetting20240805 and assigns it to the CopySettings field.
func (o *DiskBackupSnapshotSchedule20240805) SetCopySettings(v []DiskBackupCopySetting20240805) {
	o.CopySettings = &v
}

// GetDeleteCopiedBackups returns the DeleteCopiedBackups field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteCopiedBackups() []DeleteCopiedBackups20240805 {
	if o == nil || IsNil(o.DeleteCopiedBackups) {
		var ret []DeleteCopiedBackups20240805
		return ret
	}
	return *o.DeleteCopiedBackups
}

// GetDeleteCopiedBackupsOk returns a tuple with the DeleteCopiedBackups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteCopiedBackupsOk() (*[]DeleteCopiedBackups20240805, bool) {
	if o == nil || IsNil(o.DeleteCopiedBackups) {
		return nil, false
	}

	return o.DeleteCopiedBackups, true
}

// HasDeleteCopiedBackups returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasDeleteCopiedBackups() bool {
	if o != nil && !IsNil(o.DeleteCopiedBackups) {
		return true
	}

	return false
}

// SetDeleteCopiedBackups gets a reference to the given []DeleteCopiedBackups20240805 and assigns it to the DeleteCopiedBackups field.
func (o *DiskBackupSnapshotSchedule20240805) SetDeleteCopiedBackups(v []DeleteCopiedBackups20240805) {
	o.DeleteCopiedBackups = &v
}

// GetDeleteCopySnapshots returns the DeleteCopySnapshots field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteCopySnapshots() bool {
	if o == nil || IsNil(o.DeleteCopySnapshots) {
		var ret bool
		return ret
	}
	return *o.DeleteCopySnapshots
}

// GetDeleteCopySnapshotsOk returns a tuple with the DeleteCopySnapshots field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteCopySnapshotsOk() (*bool, bool) {
	if o == nil || IsNil(o.DeleteCopySnapshots) {
		return nil, false
	}

	return o.DeleteCopySnapshots, true
}

// HasDeleteCopySnapshots returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasDeleteCopySnapshots() bool {
	if o != nil && !IsNil(o.DeleteCopySnapshots) {
		return true
	}

	return false
}

// SetDeleteCopySnapshots gets a reference to the given bool and assigns it to the DeleteCopySnapshots field.
func (o *DiskBackupSnapshotSchedule20240805) SetDeleteCopySnapshots(v bool) {
	o.DeleteCopySnapshots = &v
}

// GetDeleteSnapshots returns the DeleteSnapshots field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteSnapshots() bool {
	if o == nil || IsNil(o.DeleteSnapshots) {
		var ret bool
		return ret
	}
	return *o.DeleteSnapshots
}

// GetDeleteSnapshotsOk returns a tuple with the DeleteSnapshots field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetDeleteSnapshotsOk() (*bool, bool) {
	if o == nil || IsNil(o.DeleteSnapshots) {
		return nil, false
	}

	return o.DeleteSnapshots, true
}

// HasDeleteSnapshots returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasDeleteSnapshots() bool {
	if o != nil && !IsNil(o.DeleteSnapshots) {
		return true
	}

	return false
}

// SetDeleteSnapshots gets a reference to the given bool and assigns it to the DeleteSnapshots field.
func (o *DiskBackupSnapshotSchedule20240805) SetDeleteSnapshots(v bool) {
	o.DeleteSnapshots = &v
}

// GetExport returns the Export field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetExport() AutoExportPolicy {
	if o == nil || IsNil(o.Export) {
		var ret AutoExportPolicy
		return ret
	}
	return *o.Export
}

// GetExportOk returns a tuple with the Export field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetExportOk() (*AutoExportPolicy, bool) {
	if o == nil || IsNil(o.Export) {
		return nil, false
	}

	return o.Export, true
}

// HasExport returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasExport() bool {
	if o != nil && !IsNil(o.Export) {
		return true
	}

	return false
}

// SetExport gets a reference to the given AutoExportPolicy and assigns it to the Export field.
func (o *DiskBackupSnapshotSchedule20240805) SetExport(v AutoExportPolicy) {
	o.Export = &v
}

// GetExtraRetentionSettings returns the ExtraRetentionSettings field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetExtraRetentionSettings() []ExtraRetentionSetting {
	if o == nil || IsNil(o.ExtraRetentionSettings) {
		var ret []ExtraRetentionSetting
		return ret
	}
	return *o.ExtraRetentionSettings
}

// GetExtraRetentionSettingsOk returns a tuple with the ExtraRetentionSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetExtraRetentionSettingsOk() (*[]ExtraRetentionSetting, bool) {
	if o == nil || IsNil(o.ExtraRetentionSettings) {
		return nil, false
	}

	return o.ExtraRetentionSettings, true
}

// HasExtraRetentionSettings returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasExtraRetentionSettings() bool {
	if o != nil && !IsNil(o.ExtraRetentionSettings) {
		return true
	}

	return false
}

// SetExtraRetentionSettings gets a reference to the given []ExtraRetentionSetting and assigns it to the ExtraRetentionSettings field.
func (o *DiskBackupSnapshotSchedule20240805) SetExtraRetentionSettings(v []ExtraRetentionSetting) {
	o.ExtraRetentionSettings = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshotSchedule20240805) SetLinks(v []Link) {
	o.Links = &v
}

// GetNextSnapshot returns the NextSnapshot field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetNextSnapshot() time.Time {
	if o == nil || IsNil(o.NextSnapshot) {
		var ret time.Time
		return ret
	}
	return *o.NextSnapshot
}

// GetNextSnapshotOk returns a tuple with the NextSnapshot field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetNextSnapshotOk() (*time.Time, bool) {
	if o == nil || IsNil(o.NextSnapshot) {
		return nil, false
	}

	return o.NextSnapshot, true
}

// HasNextSnapshot returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasNextSnapshot() bool {
	if o != nil && !IsNil(o.NextSnapshot) {
		return true
	}

	return false
}

// SetNextSnapshot gets a reference to the given time.Time and assigns it to the NextSnapshot field.
func (o *DiskBackupSnapshotSchedule20240805) SetNextSnapshot(v time.Time) {
	o.NextSnapshot = &v
}

// GetPolicies returns the Policies field value
func (o *DiskBackupSnapshotSchedule20240805) GetPolicies() []AdvancedDiskBackupSnapshotSchedulePolicy {
	if o == nil {
		var ret []AdvancedDiskBackupSnapshotSchedulePolicy
		return ret
	}

	return o.Policies
}

// GetPoliciesOk returns a tuple with the Policies field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetPoliciesOk() (*[]AdvancedDiskBackupSnapshotSchedulePolicy, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Policies, true
}

// SetPolicies sets field value
func (o *DiskBackupSnapshotSchedule20240805) SetPolicies(v []AdvancedDiskBackupSnapshotSchedulePolicy) {
	o.Policies = v
}

// GetReferenceHourOfDay returns the ReferenceHourOfDay field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetReferenceHourOfDay() int {
	if o == nil || IsNil(o.ReferenceHourOfDay) {
		var ret int
		return ret
	}
	return *o.ReferenceHourOfDay
}

// GetReferenceHourOfDayOk returns a tuple with the ReferenceHourOfDay field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetReferenceHourOfDayOk() (*int, bool) {
	if o == nil || IsNil(o.ReferenceHourOfDay) {
		return nil, false
	}

	return o.ReferenceHourOfDay, true
}

// HasReferenceHourOfDay returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasReferenceHourOfDay() bool {
	if o != nil && !IsNil(o.ReferenceHourOfDay) {
		return true
	}

	return false
}

// SetReferenceHourOfDay gets a reference to the given int and assigns it to the ReferenceHourOfDay field.
func (o *DiskBackupSnapshotSchedule20240805) SetReferenceHourOfDay(v int) {
	o.ReferenceHourOfDay = &v
}

// GetReferenceMinuteOfHour returns the ReferenceMinuteOfHour field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetReferenceMinuteOfHour() int {
	if o == nil || IsNil(o.ReferenceMinuteOfHour) {
		var ret int
		return ret
	}
	return *o.ReferenceMinuteOfHour
}

// GetReferenceMinuteOfHourOk returns a tuple with the ReferenceMinuteOfHour field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetReferenceMinuteOfHourOk() (*int, bool) {
	if o == nil || IsNil(o.ReferenceMinuteOfHour) {
		return nil, false
	}

	return o.ReferenceMinuteOfHour, true
}

// HasReferenceMinuteOfHour returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasReferenceMinuteOfHour() bool {
	if o != nil && !IsNil(o.ReferenceMinuteOfHour) {
		return true
	}

	return false
}

// SetReferenceMinuteOfHour gets a reference to the given int and assigns it to the ReferenceMinuteOfHour field.
func (o *DiskBackupSnapshotSchedule20240805) SetReferenceMinuteOfHour(v int) {
	o.ReferenceMinuteOfHour = &v
}

// GetRestoreWindowDays returns the RestoreWindowDays field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetRestoreWindowDays() int {
	if o == nil || IsNil(o.RestoreWindowDays) {
		var ret int
		return ret
	}
	return *o.RestoreWindowDays
}

// GetRestoreWindowDaysOk returns a tuple with the RestoreWindowDays field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetRestoreWindowDaysOk() (*int, bool) {
	if o == nil || IsNil(o.RestoreWindowDays) {
		return nil, false
	}

	return o.RestoreWindowDays, true
}

// HasRestoreWindowDays returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasRestoreWindowDays() bool {
	if o != nil && !IsNil(o.RestoreWindowDays) {
		return true
	}

	return false
}

// SetRestoreWindowDays gets a reference to the given int and assigns it to the RestoreWindowDays field.
func (o *DiskBackupSnapshotSchedule20240805) SetRestoreWindowDays(v int) {
	o.RestoreWindowDays = &v
}

// GetUpdateCopySnapshots returns the UpdateCopySnapshots field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetUpdateCopySnapshots() bool {
	if o == nil || IsNil(o.UpdateCopySnapshots) {
		var ret bool
		return ret
	}
	return *o.UpdateCopySnapshots
}

// GetUpdateCopySnapshotsOk returns a tuple with the UpdateCopySnapshots field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetUpdateCopySnapshotsOk() (*bool, bool) {
	if o == nil || IsNil(o.UpdateCopySnapshots) {
		return nil, false
	}

	return o.UpdateCopySnapshots, true
}

// HasUpdateCopySnapshots returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasUpdateCopySnapshots() bool {
	if o != nil && !IsNil(o.UpdateCopySnapshots) {
		return true
	}

	return false
}

// SetUpdateCopySnapshots gets a reference to the given bool and assigns it to the UpdateCopySnapshots field.
func (o *DiskBackupSnapshotSchedule20240805) SetUpdateCopySnapshots(v bool) {
	o.UpdateCopySnapshots = &v
}

// GetUpdateSnapshots returns the UpdateSnapshots field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetUpdateSnapshots() bool {
	if o == nil || IsNil(o.UpdateSnapshots) {
		var ret bool
		return ret
	}
	return *o.UpdateSnapshots
}

// GetUpdateSnapshotsOk returns a tuple with the UpdateSnapshots field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetUpdateSnapshotsOk() (*bool, bool) {
	if o == nil || IsNil(o.UpdateSnapshots) {
		return nil, false
	}

	return o.UpdateSnapshots, true
}

// HasUpdateSnapshots returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasUpdateSnapshots() bool {
	if o != nil && !IsNil(o.UpdateSnapshots) {
		return true
	}

	return false
}

// SetUpdateSnapshots gets a reference to the given bool and assigns it to the UpdateSnapshots field.
func (o *DiskBackupSnapshotSchedule20240805) SetUpdateSnapshots(v bool) {
	o.UpdateSnapshots = &v
}

// GetUseOrgAndGroupNamesInExportPrefix returns the UseOrgAndGroupNamesInExportPrefix field value if set, zero value otherwise
func (o *DiskBackupSnapshotSchedule20240805) GetUseOrgAndGroupNamesInExportPrefix() bool {
	if o == nil || IsNil(o.UseOrgAndGroupNamesInExportPrefix) {
		var ret bool
		return ret
	}
	return *o.UseOrgAndGroupNamesInExportPrefix
}

// GetUseOrgAndGroupNamesInExportPrefixOk returns a tuple with the UseOrgAndGroupNamesInExportPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotSchedule20240805) GetUseOrgAndGroupNamesInExportPrefixOk() (*bool, bool) {
	if o == nil || IsNil(o.UseOrgAndGroupNamesInExportPrefix) {
		return nil, false
	}

	return o.UseOrgAndGroupNamesInExportPrefix, true
}

// HasUseOrgAndGroupNamesInExportPrefix returns a boolean if a field has been set.
func (o *DiskBackupSnapshotSchedule20240805) HasUseOrgAndGroupNamesInExportPrefix() bool {
	if o != nil && !IsNil(o.UseOrgAndGroupNamesInExportPrefix) {
		return true
	}

	return false
}

// SetUseOrgAndGroupNamesInExportPrefix gets a reference to the given bool and assigns it to the UseOrgAndGroupNamesInExportPrefix field.
func (o *DiskBackupSnapshotSchedule20240805) SetUseOrgAndGroupNamesInExportPrefix(v bool) {
	o.UseOrgAndGroupNamesInExportPrefix = &v
}
