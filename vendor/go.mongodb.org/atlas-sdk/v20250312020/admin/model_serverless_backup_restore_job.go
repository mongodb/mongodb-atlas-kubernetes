// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// ServerlessBackupRestoreJob struct for ServerlessBackupRestoreJob
type ServerlessBackupRestoreJob struct {
	// Flag that indicates whether someone canceled this restore job.
	// Read only field.
	Cancelled *bool `json:"cancelled,omitempty"`
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
	// Unique 24-hexadecimal character string that identifies the snapshot.
	SnapshotId *string `json:"snapshotId,omitempty"`
	// Human-readable label that identifies the target cluster to which the restore job restores the snapshot. The resource returns this parameter when `\"deliveryType\":` `\"automated\"`.
	TargetClusterName string `json:"targetClusterName"`
	// Unique 24-hexadecimal digit string that identifies the target project for the specified `targetClusterName`.
	TargetGroupId string `json:"targetGroupId"`
	// Date and time when MongoDB Cloud took the snapshot associated with `snapshotId`. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// NewServerlessBackupRestoreJob instantiates a new ServerlessBackupRestoreJob object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerlessBackupRestoreJob(deliveryType string, targetClusterName string, targetGroupId string) *ServerlessBackupRestoreJob {
	this := ServerlessBackupRestoreJob{}
	this.DeliveryType = deliveryType
	this.TargetClusterName = targetClusterName
	this.TargetGroupId = targetGroupId
	return &this
}

// NewServerlessBackupRestoreJobWithDefaults instantiates a new ServerlessBackupRestoreJob object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerlessBackupRestoreJobWithDefaults() *ServerlessBackupRestoreJob {
	this := ServerlessBackupRestoreJob{}
	return &this
}

// GetCancelled returns the Cancelled field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetCancelled() bool {
	if o == nil || IsNil(o.Cancelled) {
		var ret bool
		return ret
	}
	return *o.Cancelled
}

// GetCancelledOk returns a tuple with the Cancelled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetCancelledOk() (*bool, bool) {
	if o == nil || IsNil(o.Cancelled) {
		return nil, false
	}

	return o.Cancelled, true
}

// HasCancelled returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasCancelled() bool {
	if o != nil && !IsNil(o.Cancelled) {
		return true
	}

	return false
}

// SetCancelled gets a reference to the given bool and assigns it to the Cancelled field.
func (o *ServerlessBackupRestoreJob) SetCancelled(v bool) {
	o.Cancelled = &v
}

// GetDeliveryType returns the DeliveryType field value
func (o *ServerlessBackupRestoreJob) GetDeliveryType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DeliveryType
}

// GetDeliveryTypeOk returns a tuple with the DeliveryType field value
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetDeliveryTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DeliveryType, true
}

// SetDeliveryType sets field value
func (o *ServerlessBackupRestoreJob) SetDeliveryType(v string) {
	o.DeliveryType = v
}

// GetDeliveryUrl returns the DeliveryUrl field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetDeliveryUrl() []string {
	if o == nil || IsNil(o.DeliveryUrl) {
		var ret []string
		return ret
	}
	return *o.DeliveryUrl
}

// GetDeliveryUrlOk returns a tuple with the DeliveryUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetDeliveryUrlOk() (*[]string, bool) {
	if o == nil || IsNil(o.DeliveryUrl) {
		return nil, false
	}

	return o.DeliveryUrl, true
}

// HasDeliveryUrl returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasDeliveryUrl() bool {
	if o != nil && !IsNil(o.DeliveryUrl) {
		return true
	}

	return false
}

// SetDeliveryUrl gets a reference to the given []string and assigns it to the DeliveryUrl field.
func (o *ServerlessBackupRestoreJob) SetDeliveryUrl(v []string) {
	o.DeliveryUrl = &v
}

// GetDesiredTimestamp returns the DesiredTimestamp field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetDesiredTimestamp() ApiBSONTimestamp {
	if o == nil || IsNil(o.DesiredTimestamp) {
		var ret ApiBSONTimestamp
		return ret
	}
	return *o.DesiredTimestamp
}

// GetDesiredTimestampOk returns a tuple with the DesiredTimestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetDesiredTimestampOk() (*ApiBSONTimestamp, bool) {
	if o == nil || IsNil(o.DesiredTimestamp) {
		return nil, false
	}

	return o.DesiredTimestamp, true
}

// HasDesiredTimestamp returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasDesiredTimestamp() bool {
	if o != nil && !IsNil(o.DesiredTimestamp) {
		return true
	}

	return false
}

// SetDesiredTimestamp gets a reference to the given ApiBSONTimestamp and assigns it to the DesiredTimestamp field.
func (o *ServerlessBackupRestoreJob) SetDesiredTimestamp(v ApiBSONTimestamp) {
	o.DesiredTimestamp = &v
}

// GetExpired returns the Expired field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetExpired() bool {
	if o == nil || IsNil(o.Expired) {
		var ret bool
		return ret
	}
	return *o.Expired
}

// GetExpiredOk returns a tuple with the Expired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetExpiredOk() (*bool, bool) {
	if o == nil || IsNil(o.Expired) {
		return nil, false
	}

	return o.Expired, true
}

// HasExpired returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasExpired() bool {
	if o != nil && !IsNil(o.Expired) {
		return true
	}

	return false
}

// SetExpired gets a reference to the given bool and assigns it to the Expired field.
func (o *ServerlessBackupRestoreJob) SetExpired(v bool) {
	o.Expired = &v
}

// GetExpiresAt returns the ExpiresAt field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetExpiresAt() time.Time {
	if o == nil || IsNil(o.ExpiresAt) {
		var ret time.Time
		return ret
	}
	return *o.ExpiresAt
}

// GetExpiresAtOk returns a tuple with the ExpiresAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetExpiresAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpiresAt) {
		return nil, false
	}

	return o.ExpiresAt, true
}

// HasExpiresAt returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasExpiresAt() bool {
	if o != nil && !IsNil(o.ExpiresAt) {
		return true
	}

	return false
}

// SetExpiresAt gets a reference to the given time.Time and assigns it to the ExpiresAt field.
func (o *ServerlessBackupRestoreJob) SetExpiresAt(v time.Time) {
	o.ExpiresAt = &v
}

// GetFailed returns the Failed field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetFailed() bool {
	if o == nil || IsNil(o.Failed) {
		var ret bool
		return ret
	}
	return *o.Failed
}

// GetFailedOk returns a tuple with the Failed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetFailedOk() (*bool, bool) {
	if o == nil || IsNil(o.Failed) {
		return nil, false
	}

	return o.Failed, true
}

// HasFailed returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasFailed() bool {
	if o != nil && !IsNil(o.Failed) {
		return true
	}

	return false
}

// SetFailed gets a reference to the given bool and assigns it to the Failed field.
func (o *ServerlessBackupRestoreJob) SetFailed(v bool) {
	o.Failed = &v
}

// GetFinishedAt returns the FinishedAt field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetFinishedAt() time.Time {
	if o == nil || IsNil(o.FinishedAt) {
		var ret time.Time
		return ret
	}
	return *o.FinishedAt
}

// GetFinishedAtOk returns a tuple with the FinishedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetFinishedAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.FinishedAt) {
		return nil, false
	}

	return o.FinishedAt, true
}

// HasFinishedAt returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasFinishedAt() bool {
	if o != nil && !IsNil(o.FinishedAt) {
		return true
	}

	return false
}

// SetFinishedAt gets a reference to the given time.Time and assigns it to the FinishedAt field.
func (o *ServerlessBackupRestoreJob) SetFinishedAt(v time.Time) {
	o.FinishedAt = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ServerlessBackupRestoreJob) SetId(v string) {
	o.Id = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *ServerlessBackupRestoreJob) SetLinks(v []Link) {
	o.Links = &v
}

// GetOplogInc returns the OplogInc field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetOplogInc() int {
	if o == nil || IsNil(o.OplogInc) {
		var ret int
		return ret
	}
	return *o.OplogInc
}

// GetOplogIncOk returns a tuple with the OplogInc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetOplogIncOk() (*int, bool) {
	if o == nil || IsNil(o.OplogInc) {
		return nil, false
	}

	return o.OplogInc, true
}

// HasOplogInc returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasOplogInc() bool {
	if o != nil && !IsNil(o.OplogInc) {
		return true
	}

	return false
}

// SetOplogInc gets a reference to the given int and assigns it to the OplogInc field.
func (o *ServerlessBackupRestoreJob) SetOplogInc(v int) {
	o.OplogInc = &v
}

// GetOplogTs returns the OplogTs field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetOplogTs() int {
	if o == nil || IsNil(o.OplogTs) {
		var ret int
		return ret
	}
	return *o.OplogTs
}

// GetOplogTsOk returns a tuple with the OplogTs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetOplogTsOk() (*int, bool) {
	if o == nil || IsNil(o.OplogTs) {
		return nil, false
	}

	return o.OplogTs, true
}

// HasOplogTs returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasOplogTs() bool {
	if o != nil && !IsNil(o.OplogTs) {
		return true
	}

	return false
}

// SetOplogTs gets a reference to the given int and assigns it to the OplogTs field.
func (o *ServerlessBackupRestoreJob) SetOplogTs(v int) {
	o.OplogTs = &v
}

// GetPointInTimeUTCSeconds returns the PointInTimeUTCSeconds field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetPointInTimeUTCSeconds() int {
	if o == nil || IsNil(o.PointInTimeUTCSeconds) {
		var ret int
		return ret
	}
	return *o.PointInTimeUTCSeconds
}

// GetPointInTimeUTCSecondsOk returns a tuple with the PointInTimeUTCSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetPointInTimeUTCSecondsOk() (*int, bool) {
	if o == nil || IsNil(o.PointInTimeUTCSeconds) {
		return nil, false
	}

	return o.PointInTimeUTCSeconds, true
}

// HasPointInTimeUTCSeconds returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasPointInTimeUTCSeconds() bool {
	if o != nil && !IsNil(o.PointInTimeUTCSeconds) {
		return true
	}

	return false
}

// SetPointInTimeUTCSeconds gets a reference to the given int and assigns it to the PointInTimeUTCSeconds field.
func (o *ServerlessBackupRestoreJob) SetPointInTimeUTCSeconds(v int) {
	o.PointInTimeUTCSeconds = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}

	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *ServerlessBackupRestoreJob) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetTargetClusterName returns the TargetClusterName field value
func (o *ServerlessBackupRestoreJob) GetTargetClusterName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TargetClusterName
}

// GetTargetClusterNameOk returns a tuple with the TargetClusterName field value
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetTargetClusterNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TargetClusterName, true
}

// SetTargetClusterName sets field value
func (o *ServerlessBackupRestoreJob) SetTargetClusterName(v string) {
	o.TargetClusterName = v
}

// GetTargetGroupId returns the TargetGroupId field value
func (o *ServerlessBackupRestoreJob) GetTargetGroupId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TargetGroupId
}

// GetTargetGroupIdOk returns a tuple with the TargetGroupId field value
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetTargetGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TargetGroupId, true
}

// SetTargetGroupId sets field value
func (o *ServerlessBackupRestoreJob) SetTargetGroupId(v string) {
	o.TargetGroupId = v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise
func (o *ServerlessBackupRestoreJob) GetTimestamp() time.Time {
	if o == nil || IsNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ServerlessBackupRestoreJob) GetTimestampOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Timestamp) {
		return nil, false
	}

	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *ServerlessBackupRestoreJob) HasTimestamp() bool {
	if o != nil && !IsNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *ServerlessBackupRestoreJob) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}
