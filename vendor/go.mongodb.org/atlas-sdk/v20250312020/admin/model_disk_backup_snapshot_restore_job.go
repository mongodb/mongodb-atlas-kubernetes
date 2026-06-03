// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// DiskBackupSnapshotRestoreJob struct for DiskBackupSnapshotRestoreJob
type DiskBackupSnapshotRestoreJob struct {
	// Flag that indicates whether someone canceled this restore job.
	// Read only field.
	Cancelled *bool `json:"cancelled,omitempty"`
	// Information on the restore job for each replica set in the sharded cluster.
	// Read only field.
	Components *[]DiskBackupRestoreMember `json:"components,omitempty"`
	// Human-readable label that categorizes the restore job to create.
	DeliveryType string `json:"deliveryType"`
	// One or more Uniform Resource Locators (URLs) that point to the compressed snapshot files for manual download. MongoDB Cloud returns this parameter when `\"deliveryType\" : \"download\"`.
	// Read only field.
	DeliveryUrl      *[]string         `json:"deliveryUrl,omitempty"`
	DesiredTimestamp *ApiBSONTimestamp `json:"desiredTimestamp,omitempty"`
	// Flag that indicates whether the restore job expired.
	// Read only field.
	Expired *bool `json:"expired,omitempty"`
	// Date and time when the restore job expires. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Flag that indicates whether the restore job failed.
	// Read only field.
	Failed *bool `json:"failed,omitempty"`
	// Date and time when the restore job completed. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	FinishedAt *time.Time `json:"finishedAt,omitempty"`
	// Unique 24-hexadecimal character string that identifies the restore job.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Oplog operation number from which you want to restore this snapshot. This number represents the second part of an Oplog timestamp. The resource returns this parameter when `\"deliveryType\" : \"pointInTime\"` and `oplogTs` exceeds `0`.
	OplogInc *int `json:"oplogInc,omitempty"`
	// Date and time from which you want to restore this snapshot. This parameter expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch. This number represents the first part of an Oplog timestamp. The resource returns this parameter when `\"deliveryType\" : \"pointInTime\"` and `oplogTs` exceeds `0`.
	OplogTs *int `json:"oplogTs,omitempty"`
	// Date and time from which MongoDB Cloud restored this snapshot. This parameter expresses this timestamp in the number of seconds that have elapsed since the UNIX epoch. The resource returns this parameter when `\"deliveryType\" : \"pointInTime\"` and `pointInTimeUTCSeconds` exceeds `0`.
	PointInTimeUTCSeconds *int `json:"pointInTimeUTCSeconds,omitempty"`
	// One or more Uniform Resource Locators (URLs) that point to the compressed snapshot files for manual download and the corresponding private endpoint(s). MongoDB Cloud returns this parameter when `\"deliveryType\" : \"download\"` and the download can be performed privately.
	// Read only field.
	PrivateDownloadDeliveryUrls *[]ApiPrivateDownloadDeliveryUrl `json:"privateDownloadDeliveryUrls,omitempty"`
	// Unique 24-hexadecimal character string that identifies the snapshot.
	SnapshotId *string `json:"snapshotId,omitempty"`
	// Human-readable label that identifies the target cluster to which the restore job restores the snapshot. The resource returns this parameter when `\"deliveryType\":` `\"automated\"`. Required for `automated` and `pointInTime` restore types.
	TargetClusterName *string `json:"targetClusterName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the target project for the specified `targetClusterName`. Required for `automated` and `pointInTime` restore types.
	TargetGroupId *string `json:"targetGroupId,omitempty"`
	// Date and time when MongoDB Cloud took the snapshot associated with `snapshotId`. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// NewDiskBackupSnapshotRestoreJob instantiates a new DiskBackupSnapshotRestoreJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiskBackupSnapshotRestoreJob(deliveryType string) *DiskBackupSnapshotRestoreJob {
	this := DiskBackupSnapshotRestoreJob{}
	this.DeliveryType = deliveryType
	return &this
}

// NewDiskBackupSnapshotRestoreJobWithDefaults instantiates a new DiskBackupSnapshotRestoreJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiskBackupSnapshotRestoreJobWithDefaults() *DiskBackupSnapshotRestoreJob {
	this := DiskBackupSnapshotRestoreJob{}
	return &this
}

// GetCancelled returns the Cancelled field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetCancelled() bool {
	if o == nil || IsNil(o.Cancelled) {
		var ret bool
		return ret
	}
	return *o.Cancelled
}

// GetCancelledOk returns a tuple with the Cancelled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetCancelledOk() (*bool, bool) {
	if o == nil || IsNil(o.Cancelled) {
		return nil, false
	}

	return o.Cancelled, true
}

// HasCancelled returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasCancelled() bool {
	if o != nil && !IsNil(o.Cancelled) {
		return true
	}

	return false
}

// SetCancelled gets a reference to the given bool and assigns it to the Cancelled field.
func (o *DiskBackupSnapshotRestoreJob) SetCancelled(v bool) {
	o.Cancelled = &v
}

// GetComponents returns the Components field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetComponents() []DiskBackupRestoreMember {
	if o == nil || IsNil(o.Components) {
		var ret []DiskBackupRestoreMember
		return ret
	}
	return *o.Components
}

// GetComponentsOk returns a tuple with the Components field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetComponentsOk() (*[]DiskBackupRestoreMember, bool) {
	if o == nil || IsNil(o.Components) {
		return nil, false
	}

	return o.Components, true
}

// HasComponents returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasComponents() bool {
	if o != nil && !IsNil(o.Components) {
		return true
	}

	return false
}

// SetComponents gets a reference to the given []DiskBackupRestoreMember and assigns it to the Components field.
func (o *DiskBackupSnapshotRestoreJob) SetComponents(v []DiskBackupRestoreMember) {
	o.Components = &v
}

// GetDeliveryType returns the DeliveryType field value
func (o *DiskBackupSnapshotRestoreJob) GetDeliveryType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DeliveryType
}

// GetDeliveryTypeOk returns a tuple with the DeliveryType field value
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetDeliveryTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DeliveryType, true
}

// SetDeliveryType sets field value
func (o *DiskBackupSnapshotRestoreJob) SetDeliveryType(v string) {
	o.DeliveryType = v
}

// GetDeliveryUrl returns the DeliveryUrl field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetDeliveryUrl() []string {
	if o == nil || IsNil(o.DeliveryUrl) {
		var ret []string
		return ret
	}
	return *o.DeliveryUrl
}

// GetDeliveryUrlOk returns a tuple with the DeliveryUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetDeliveryUrlOk() (*[]string, bool) {
	if o == nil || IsNil(o.DeliveryUrl) {
		return nil, false
	}

	return o.DeliveryUrl, true
}

// HasDeliveryUrl returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasDeliveryUrl() bool {
	if o != nil && !IsNil(o.DeliveryUrl) {
		return true
	}

	return false
}

// SetDeliveryUrl gets a reference to the given []string and assigns it to the DeliveryUrl field.
func (o *DiskBackupSnapshotRestoreJob) SetDeliveryUrl(v []string) {
	o.DeliveryUrl = &v
}

// GetDesiredTimestamp returns the DesiredTimestamp field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetDesiredTimestamp() ApiBSONTimestamp {
	if o == nil || IsNil(o.DesiredTimestamp) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.DesiredTimestamp
}

// GetDesiredTimestampOk returns a tuple with the DesiredTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetDesiredTimestampOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.DesiredTimestamp) {
		return nil, false
	}

	return o.DesiredTimestamp, true
}

// HasDesiredTimestamp returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasDesiredTimestamp() bool {
	if o != nil && !IsNil(o.DesiredTimestamp) {
		return true
	}

	return false
}

// SetDesiredTimestamp gets a reference to the given ApiBSONTimestamp and assigns it to the DesiredTimestamp field.
func (o *DiskBackupSnapshotRestoreJob) SetDesiredTimestamp(v ApiBSONTimestamp) {
	o.DesiredTimestamp = &v
}

// GetExpired returns the Expired field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetExpired() bool {
	if o == nil || IsNil(o.Expired) {
		var ret bool
		return ret
	}
	return *o.Expired
}

// GetExpiredOk returns a tuple with the Expired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetExpiredOk() (*bool, bool) {
	if o == nil || IsNil(o.Expired) {
		return nil, false
	}

	return o.Expired, true
}

// HasExpired returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasExpired() bool {
	if o != nil && !IsNil(o.Expired) {
		return true
	}

	return false
}

// SetExpired gets a reference to the given bool and assigns it to the Expired field.
func (o *DiskBackupSnapshotRestoreJob) SetExpired(v bool) {
	o.Expired = &v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetExpiresAt() time.Time {
	if o == nil || IsNil(o.ExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpiresAt) {
		return nil, false
	}

	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasExpiresAt() bool {
	if o != nil && !IsNil(o.ExpiresAt) {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *DiskBackupSnapshotRestoreJob) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

// GetFailed returns the Failed field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetFailed() bool {
	if o == nil || IsNil(o.Failed) {
		var ret bool
		return ret
	}
	return *o.Failed
}

// GetFailedOk returns a tuple with the Failed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetFailedOk() (*bool, bool) {
	if o == nil || IsNil(o.Failed) {
		return nil, false
	}

	return o.Failed, true
}

// HasFailed returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasFailed() bool {
	if o != nil && !IsNil(o.Failed) {
		return true
	}

	return false
}

// SetFailed gets a reference to the given bool and assigns it to the Failed field.
func (o *DiskBackupSnapshotRestoreJob) SetFailed(v bool) {
	o.Failed = &v
}

// GetFinishedAt returns the FinishedAt field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetFinishedAt() time.Time {
	if o == nil || IsNil(o.FinishedAt) {
		var ret time.Time
		return ret
	}
	return *o.FinishedAt
}

// GetFinishedAtOk returns a tuple with the FinishedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetFinishedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FinishedAt) {
		return nil, false
	}

	return o.FinishedAt, true
}

// HasFinishedAt returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasFinishedAt() bool {
	if o != nil && !IsNil(o.FinishedAt) {
		return true
	}

	return false
}

// SetFinishedAt gets a reference to the given time.Time and assigns it to the FinishedAt field.
func (o *DiskBackupSnapshotRestoreJob) SetFinishedAt(v time.Time) {
	o.FinishedAt = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *DiskBackupSnapshotRestoreJob) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *DiskBackupSnapshotRestoreJob) SetLinks(v []Link) {
	o.Links = &v
}

// GetOplogInc returns the OplogInc field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetOplogInc() int {
	if o == nil || IsNil(o.OplogInc) {
		var ret int
		return ret
	}
	return *o.OplogInc
}

// GetOplogIncOk returns a tuple with the OplogInc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetOplogIncOk() (*int, bool) {
	if o == nil || IsNil(o.OplogInc) {
		return nil, false
	}

	return o.OplogInc, true
}

// HasOplogInc returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasOplogInc() bool {
	if o != nil && !IsNil(o.OplogInc) {
		return true
	}

	return false
}

// SetOplogInc gets a reference to the given int and assigns it to the OplogInc field.
func (o *DiskBackupSnapshotRestoreJob) SetOplogInc(v int) {
	o.OplogInc = &v
}

// GetOplogTs returns the OplogTs field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetOplogTs() int {
	if o == nil || IsNil(o.OplogTs) {
		var ret int
		return ret
	}
	return *o.OplogTs
}

// GetOplogTsOk returns a tuple with the OplogTs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetOplogTsOk() (*int, bool) {
	if o == nil || IsNil(o.OplogTs) {
		return nil, false
	}

	return o.OplogTs, true
}

// HasOplogTs returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasOplogTs() bool {
	if o != nil && !IsNil(o.OplogTs) {
		return true
	}

	return false
}

// SetOplogTs gets a reference to the given int and assigns it to the OplogTs field.
func (o *DiskBackupSnapshotRestoreJob) SetOplogTs(v int) {
	o.OplogTs = &v
}

// GetPointInTimeUTCSeconds returns the PointInTimeUTCSeconds field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetPointInTimeUTCSeconds() int {
	if o == nil || IsNil(o.PointInTimeUTCSeconds) {
		var ret int
		return ret
	}
	return *o.PointInTimeUTCSeconds
}

// GetPointInTimeUTCSecondsOk returns a tuple with the PointInTimeUTCSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetPointInTimeUTCSecondsOk() (*int, bool) {
	if o == nil || IsNil(o.PointInTimeUTCSeconds) {
		return nil, false
	}

	return o.PointInTimeUTCSeconds, true
}

// HasPointInTimeUTCSeconds returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasPointInTimeUTCSeconds() bool {
	if o != nil && !IsNil(o.PointInTimeUTCSeconds) {
		return true
	}

	return false
}

// SetPointInTimeUTCSeconds gets a reference to the given int and assigns it to the PointInTimeUTCSeconds field.
func (o *DiskBackupSnapshotRestoreJob) SetPointInTimeUTCSeconds(v int) {
	o.PointInTimeUTCSeconds = &v
}

// GetPrivateDownloadDeliveryUrls returns the PrivateDownloadDeliveryUrls field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetPrivateDownloadDeliveryUrls() []ApiPrivateDownloadDeliveryUrl {
	if o == nil || IsNil(o.PrivateDownloadDeliveryUrls) {
		var ret []ApiPrivateDownloadDeliveryUrl
		return ret
	}
	return *o.PrivateDownloadDeliveryUrls
}

// GetPrivateDownloadDeliveryUrlsOk returns a tuple with the PrivateDownloadDeliveryUrls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetPrivateDownloadDeliveryUrlsOk() (*[]ApiPrivateDownloadDeliveryUrl, bool) {
	if o == nil || IsNil(o.PrivateDownloadDeliveryUrls) {
		return nil, false
	}

	return o.PrivateDownloadDeliveryUrls, true
}

// HasPrivateDownloadDeliveryUrls returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasPrivateDownloadDeliveryUrls() bool {
	if o != nil && !IsNil(o.PrivateDownloadDeliveryUrls) {
		return true
	}

	return false
}

// SetPrivateDownloadDeliveryUrls gets a reference to the given []ApiPrivateDownloadDeliveryUrl and assigns it to the PrivateDownloadDeliveryUrls field.
func (o *DiskBackupSnapshotRestoreJob) SetPrivateDownloadDeliveryUrls(v []ApiPrivateDownloadDeliveryUrl) {
	o.PrivateDownloadDeliveryUrls = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}

	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *DiskBackupSnapshotRestoreJob) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetTargetClusterName returns the TargetClusterName field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetTargetClusterName() string {
	if o == nil || IsNil(o.TargetClusterName) {
		var ret string
		return ret
	}
	return *o.TargetClusterName
}

// GetTargetClusterNameOk returns a tuple with the TargetClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetTargetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.TargetClusterName) {
		return nil, false
	}

	return o.TargetClusterName, true
}

// HasTargetClusterName returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasTargetClusterName() bool {
	if o != nil && !IsNil(o.TargetClusterName) {
		return true
	}

	return false
}

// SetTargetClusterName gets a reference to the given string and assigns it to the TargetClusterName field.
func (o *DiskBackupSnapshotRestoreJob) SetTargetClusterName(v string) {
	o.TargetClusterName = &v
}

// GetTargetGroupId returns the TargetGroupId field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetTargetGroupId() string {
	if o == nil || IsNil(o.TargetGroupId) {
		var ret string
		return ret
	}
	return *o.TargetGroupId
}

// GetTargetGroupIdOk returns a tuple with the TargetGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetTargetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TargetGroupId) {
		return nil, false
	}

	return o.TargetGroupId, true
}

// HasTargetGroupId returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasTargetGroupId() bool {
	if o != nil && !IsNil(o.TargetGroupId) {
		return true
	}

	return false
}

// SetTargetGroupId gets a reference to the given string and assigns it to the TargetGroupId field.
func (o *DiskBackupSnapshotRestoreJob) SetTargetGroupId(v string) {
	o.TargetGroupId = &v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise
func (o *DiskBackupSnapshotRestoreJob) GetTimestamp() time.Time {
	if o == nil || IsNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiskBackupSnapshotRestoreJob) GetTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}

	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *DiskBackupSnapshotRestoreJob) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *DiskBackupSnapshotRestoreJob) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}
