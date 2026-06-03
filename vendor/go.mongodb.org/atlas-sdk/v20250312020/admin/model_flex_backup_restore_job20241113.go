// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// FlexBackupRestoreJob20241113 Details for one restore job of a flex cluster.
type FlexBackupRestoreJob20241113 struct {
	// Means by which this resource returns the snapshot to the requesting MongoDB Cloud user.
	// Read only field.
	DeliveryType *string `json:"deliveryType,omitempty"`
	// Date and time when the download link no longer works. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the restore job.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Human-readable label that identifies the source instance.
	// Read only field.
	InstanceName *string `json:"instanceName,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project from which the restore job originated.
	// Read only field.
	ProjectId *string `json:"projectId,omitempty"`
	// Date and time when MongoDB Cloud completed writing this snapshot. MongoDB Cloud changes the status of the restore job to `CLOSED`. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	RestoreFinishedDate *time.Time `json:"restoreFinishedDate,omitempty"`
	// Date and time when MongoDB Cloud will restore this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	RestoreScheduledDate *time.Time `json:"restoreScheduledDate,omitempty"`
	// Date and time when MongoDB Cloud completed writing this snapshot. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	SnapshotFinishedDate *time.Time `json:"snapshotFinishedDate,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the snapshot to restore.
	// Read only field.
	SnapshotId *string `json:"snapshotId,omitempty"`
	// Internet address from which you can download the compressed snapshot files. The resource returns this parameter when  `\"deliveryType\" : \"DOWNLOAD\"`.
	// Read only field.
	SnapshotUrl *string `json:"snapshotUrl,omitempty"`
	// Phase of the restore workflow for this job at the time this resource made this request.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Human-readable label that identifies the instance or cluster on the target project to which you want to restore the snapshot. You can restore the snapshot to another flex or dedicated cluster tier.
	// Read only field.
	TargetDeploymentItemName *string `json:"targetDeploymentItemName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the instance or cluster to which you want to restore the snapshot.
	// Read only field.
	TargetProjectId *string `json:"targetProjectId,omitempty"`
}

// NewFlexBackupRestoreJob20241113 instantiates a new FlexBackupRestoreJob20241113 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFlexBackupRestoreJob20241113() *FlexBackupRestoreJob20241113 {
	this := FlexBackupRestoreJob20241113{}
	return &this
}

// NewFlexBackupRestoreJob20241113WithDefaults instantiates a new FlexBackupRestoreJob20241113 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFlexBackupRestoreJob20241113WithDefaults() *FlexBackupRestoreJob20241113 {
	this := FlexBackupRestoreJob20241113{}
	return &this
}

// GetDeliveryType returns the DeliveryType field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetDeliveryType() string {
	if o == nil || IsNil(o.DeliveryType) {
		var ret string
		return ret
	}
	return *o.DeliveryType
}

// GetDeliveryTypeOk returns a tuple with the DeliveryType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetDeliveryTypeOk() (*string, bool) {
	if o == nil || IsNil(o.DeliveryType) {
		return nil, false
	}

	return o.DeliveryType, true
}

// HasDeliveryType returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasDeliveryType() bool {
	if o != nil && !IsNil(o.DeliveryType) {
		return true
	}

	return false
}

// SetDeliveryType gets a reference to the given string and assigns it to the DeliveryType field.
func (o *FlexBackupRestoreJob20241113) SetDeliveryType(v string) {
	o.DeliveryType = &v
}

// GetExpirationDate returns the ExpirationDate field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetExpirationDate() time.Time {
	if o == nil || IsNil(o.ExpirationDate) {
		var ret time.Time
		return ret
	}
	return *o.ExpirationDate
}

// GetExpirationDateOk returns a tuple with the ExpirationDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetExpirationDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpirationDate) {
		return nil, false
	}

	return o.ExpirationDate, true
}

// HasExpirationDate returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasExpirationDate() bool {
	if o != nil && !IsNil(o.ExpirationDate) {
		return true
	}

	return false
}

// SetExpirationDate gets a reference to the given time.Time and assigns it to the ExpirationDate field.
func (o *FlexBackupRestoreJob20241113) SetExpirationDate(v time.Time) {
	o.ExpirationDate = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *FlexBackupRestoreJob20241113) SetId(v string) {
	o.Id = &v
}

// GetInstanceName returns the InstanceName field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetInstanceName() string {
	if o == nil || IsNil(o.InstanceName) {
		var ret string
		return ret
	}
	return *o.InstanceName
}

// GetInstanceNameOk returns a tuple with the InstanceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetInstanceNameOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceName) {
		return nil, false
	}

	return o.InstanceName, true
}

// HasInstanceName returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasInstanceName() bool {
	if o != nil && !IsNil(o.InstanceName) {
		return true
	}

	return false
}

// SetInstanceName gets a reference to the given string and assigns it to the InstanceName field.
func (o *FlexBackupRestoreJob20241113) SetInstanceName(v string) {
	o.InstanceName = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *FlexBackupRestoreJob20241113) SetLinks(v []Link) {
	o.Links = &v
}

// GetProjectId returns the ProjectId field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetProjectId() string {
	if o == nil || IsNil(o.ProjectId) {
		var ret string
		return ret
	}
	return *o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProjectId) {
		return nil, false
	}

	return o.ProjectId, true
}

// HasProjectId returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasProjectId() bool {
	if o != nil && !IsNil(o.ProjectId) {
		return true
	}

	return false
}

// SetProjectId gets a reference to the given string and assigns it to the ProjectId field.
func (o *FlexBackupRestoreJob20241113) SetProjectId(v string) {
	o.ProjectId = &v
}

// GetRestoreFinishedDate returns the RestoreFinishedDate field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetRestoreFinishedDate() time.Time {
	if o == nil || IsNil(o.RestoreFinishedDate) {
		var ret time.Time
		return ret
	}
	return *o.RestoreFinishedDate
}

// GetRestoreFinishedDateOk returns a tuple with the RestoreFinishedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetRestoreFinishedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.RestoreFinishedDate) {
		return nil, false
	}

	return o.RestoreFinishedDate, true
}

// HasRestoreFinishedDate returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasRestoreFinishedDate() bool {
	if o != nil && !IsNil(o.RestoreFinishedDate) {
		return true
	}

	return false
}

// SetRestoreFinishedDate gets a reference to the given time.Time and assigns it to the RestoreFinishedDate field.
func (o *FlexBackupRestoreJob20241113) SetRestoreFinishedDate(v time.Time) {
	o.RestoreFinishedDate = &v
}

// GetRestoreScheduledDate returns the RestoreScheduledDate field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetRestoreScheduledDate() time.Time {
	if o == nil || IsNil(o.RestoreScheduledDate) {
		var ret time.Time
		return ret
	}
	return *o.RestoreScheduledDate
}

// GetRestoreScheduledDateOk returns a tuple with the RestoreScheduledDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetRestoreScheduledDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.RestoreScheduledDate) {
		return nil, false
	}

	return o.RestoreScheduledDate, true
}

// HasRestoreScheduledDate returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasRestoreScheduledDate() bool {
	if o != nil && !IsNil(o.RestoreScheduledDate) {
		return true
	}

	return false
}

// SetRestoreScheduledDate gets a reference to the given time.Time and assigns it to the RestoreScheduledDate field.
func (o *FlexBackupRestoreJob20241113) SetRestoreScheduledDate(v time.Time) {
	o.RestoreScheduledDate = &v
}

// GetSnapshotFinishedDate returns the SnapshotFinishedDate field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetSnapshotFinishedDate() time.Time {
	if o == nil || IsNil(o.SnapshotFinishedDate) {
		var ret time.Time
		return ret
	}
	return *o.SnapshotFinishedDate
}

// GetSnapshotFinishedDateOk returns a tuple with the SnapshotFinishedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetSnapshotFinishedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.SnapshotFinishedDate) {
		return nil, false
	}

	return o.SnapshotFinishedDate, true
}

// HasSnapshotFinishedDate returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasSnapshotFinishedDate() bool {
	if o != nil && !IsNil(o.SnapshotFinishedDate) {
		return true
	}

	return false
}

// SetSnapshotFinishedDate gets a reference to the given time.Time and assigns it to the SnapshotFinishedDate field.
func (o *FlexBackupRestoreJob20241113) SetSnapshotFinishedDate(v time.Time) {
	o.SnapshotFinishedDate = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}

	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *FlexBackupRestoreJob20241113) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetSnapshotUrl returns the SnapshotUrl field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetSnapshotUrl() string {
	if o == nil || IsNil(o.SnapshotUrl) {
		var ret string
		return ret
	}
	return *o.SnapshotUrl
}

// GetSnapshotUrlOk returns a tuple with the SnapshotUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetSnapshotUrlOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotUrl) {
		return nil, false
	}

	return o.SnapshotUrl, true
}

// HasSnapshotUrl returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasSnapshotUrl() bool {
	if o != nil && !IsNil(o.SnapshotUrl) {
		return true
	}

	return false
}

// SetSnapshotUrl gets a reference to the given string and assigns it to the SnapshotUrl field.
func (o *FlexBackupRestoreJob20241113) SetSnapshotUrl(v string) {
	o.SnapshotUrl = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *FlexBackupRestoreJob20241113) SetStatus(v string) {
	o.Status = &v
}

// GetTargetDeploymentItemName returns the TargetDeploymentItemName field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetTargetDeploymentItemName() string {
	if o == nil || IsNil(o.TargetDeploymentItemName) {
		var ret string
		return ret
	}
	return *o.TargetDeploymentItemName
}

// GetTargetDeploymentItemNameOk returns a tuple with the TargetDeploymentItemName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetTargetDeploymentItemNameOk() (*string, bool) {
	if o == nil || IsNil(o.TargetDeploymentItemName) {
		return nil, false
	}

	return o.TargetDeploymentItemName, true
}

// HasTargetDeploymentItemName returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasTargetDeploymentItemName() bool {
	if o != nil && !IsNil(o.TargetDeploymentItemName) {
		return true
	}

	return false
}

// SetTargetDeploymentItemName gets a reference to the given string and assigns it to the TargetDeploymentItemName field.
func (o *FlexBackupRestoreJob20241113) SetTargetDeploymentItemName(v string) {
	o.TargetDeploymentItemName = &v
}

// GetTargetProjectId returns the TargetProjectId field value if set, zero value otherwise
func (o *FlexBackupRestoreJob20241113) GetTargetProjectId() string {
	if o == nil || IsNil(o.TargetProjectId) {
		var ret string
		return ret
	}
	return *o.TargetProjectId
}

// GetTargetProjectIdOk returns a tuple with the TargetProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FlexBackupRestoreJob20241113) GetTargetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.TargetProjectId) {
		return nil, false
	}

	return o.TargetProjectId, true
}

// HasTargetProjectId returns a boolean if a field has been set.
func (o *FlexBackupRestoreJob20241113) HasTargetProjectId() bool {
	if o != nil && !IsNil(o.TargetProjectId) {
		return true
	}

	return false
}

// SetTargetProjectId gets a reference to the given string and assigns it to the TargetProjectId field.
func (o *FlexBackupRestoreJob20241113) SetTargetProjectId(v string) {
	o.TargetProjectId = &v
}
