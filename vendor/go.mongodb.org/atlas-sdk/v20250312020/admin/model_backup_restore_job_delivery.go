// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// BackupRestoreJobDelivery Method and details that indicate how to deliver the restored snapshot data.
type BackupRestoreJobDelivery struct {
	// Header name to use when downloading the restore, used with `\"delivery.methodName\" : \"HTTP\"`.
	// Read only field.
	AuthHeader *string `json:"authHeader,omitempty"`
	// Header value to use when downloading the restore, used with `\"delivery.methodName\" : \"HTTP\"`.
	// Read only field.
	AuthValue *string `json:"authValue,omitempty"`
	// Number of hours after the restore job completes that indicates when the Uniform Resource Locator (URL) for the snapshot download file expires. The resource returns this parameter when `\"delivery.methodName\" : \"HTTP\"`.
	ExpirationHours *int `json:"expirationHours,omitempty"`
	// Date and time when the Uniform Resource Locator (URL) for the snapshot download file expires. This parameter expresses its value in the ISO 8601 timestamp format in UTC. The resource returns this parameter when `\"delivery.methodName\" : \"HTTP\"`.
	// Read only field.
	Expires *time.Time `json:"expires,omitempty"`
	// Positive integer that indicates how many times you can use the Uniform Resource Locator (URL) for the snapshot download file. The resource returns this parameter when `\"delivery.methodName\" : \"HTTP\"`.
	MaxDownloads *int `json:"maxDownloads,omitempty"`
	// Human-readable label that identifies the means for delivering the data. If you set `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`, you must also set: `delivery.targetGroupId` and `delivery.targetClusterName` or `delivery.targetClusterId`. The response returns `\"delivery.methodName\" : \"HTTP\"` as an automated restore uses HyperText Transport Protocol (HTTP) to deliver the restore job to the target host.
	MethodName string `json:"methodName"`
	// State of the downloadable snapshot file when MongoDB Cloud received this request.
	// Read only field.
	StatusName *string `json:"statusName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the target cluster. Use the `clusterId` returned in the response body of the **Get All Snapshots** and **Get a Snapshot** endpoints. This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`.   If the target cluster doesn't have backup enabled, two resources return parameters with empty values:  - **Get All Snapshots** endpoint returns an empty results array without `clusterId` elements - **Get a Snapshot** endpoint doesn't return a `clusterId` parameter.  To return a response with the `clusterId` parameter, either use the `delivery.targetClusterName` parameter or enable backup on the target cluster.
	TargetClusterId *string `json:"targetClusterId,omitempty"`
	// Human-readable label that identifies the target cluster. Use the `clusterName` returned in the response body of the **Get All Snapshots** and **Get a Snapshot** endpoints.  This parameter applies when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`.  If the target cluster doesn't have backup enabled, two resources return parameters with empty values:  - **Get All Snapshots** endpoint returns an empty results array without `clusterId` elements - **Get a Snapshot** endpoint doesn't return a `clusterId` parameter.  To return a response with the `clusterId` parameter, either use the `delivery.targetClusterName` parameter or enable backup on the target cluster.
	TargetClusterName *string `json:"targetClusterName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that contains the destination cluster for the restore job. The resource returns this parameter when `\"delivery.methodName\" : \"AUTOMATED_RESTORE\"`.
	TargetGroupId *string `json:"targetGroupId,omitempty"`
	// Uniform Resource Locator (URL) from which you can download the restored snapshot data. URL includes the verification key. The resource returns this parameter when `\"delivery.methodName\" : \"HTTP\"`.
	// Read only field.
	// Deprecated
	Url *string `json:"url,omitempty"`
	// Uniform Resource Locator (URL) from which you can download the restored snapshot data. This should be preferred over `url`. The verification key must be sent as an HTTP header. The resource returns this parameter when `\"delivery.methodName\" : \"HTTP\"`.
	// Read only field.
	UrlV2 *string `json:"urlV2,omitempty"`
}

// NewBackupRestoreJobDelivery instantiates a new BackupRestoreJobDelivery object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupRestoreJobDelivery(methodName string) *BackupRestoreJobDelivery {
	this := BackupRestoreJobDelivery{}
	this.MethodName = methodName
	return &this
}

// NewBackupRestoreJobDeliveryWithDefaults instantiates a new BackupRestoreJobDelivery object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupRestoreJobDeliveryWithDefaults() *BackupRestoreJobDelivery {
	this := BackupRestoreJobDelivery{}
	return &this
}

// GetAuthHeader returns the AuthHeader field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetAuthHeader() string {
	if o == nil || IsNil(o.AuthHeader) {
		var ret string
		return ret
	}
	return *o.AuthHeader
}

// GetAuthHeaderOk returns a tuple with the AuthHeader field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetAuthHeaderOk() (*string, bool) {
	if o == nil || IsNil(o.AuthHeader) {
		return nil, false
	}

	return o.AuthHeader, true
}

// HasAuthHeader returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasAuthHeader() bool {
	if o != nil && !IsNil(o.AuthHeader) {
		return true
	}

	return false
}

// SetAuthHeader gets a reference to the given string and assigns it to the AuthHeader field.
func (o *BackupRestoreJobDelivery) SetAuthHeader(v string) {
	o.AuthHeader = &v
}

// GetAuthValue returns the AuthValue field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetAuthValue() string {
	if o == nil || IsNil(o.AuthValue) {
		var ret string
		return ret
	}
	return *o.AuthValue
}

// GetAuthValueOk returns a tuple with the AuthValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetAuthValueOk() (*string, bool) {
	if o == nil || IsNil(o.AuthValue) {
		return nil, false
	}

	return o.AuthValue, true
}

// HasAuthValue returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasAuthValue() bool {
	if o != nil && !IsNil(o.AuthValue) {
		return true
	}

	return false
}

// SetAuthValue gets a reference to the given string and assigns it to the AuthValue field.
func (o *BackupRestoreJobDelivery) SetAuthValue(v string) {
	o.AuthValue = &v
}

// GetExpirationHours returns the ExpirationHours field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetExpirationHours() int {
	if o == nil || IsNil(o.ExpirationHours) {
		var ret int
		return ret
	}
	return *o.ExpirationHours
}

// GetExpirationHoursOk returns a tuple with the ExpirationHours field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetExpirationHoursOk() (*int, bool) {
	if o == nil || IsNil(o.ExpirationHours) {
		return nil, false
	}

	return o.ExpirationHours, true
}

// HasExpirationHours returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasExpirationHours() bool {
	if o != nil && !IsNil(o.ExpirationHours) {
		return true
	}

	return false
}

// SetExpirationHours gets a reference to the given int and assigns it to the ExpirationHours field.
func (o *BackupRestoreJobDelivery) SetExpirationHours(v int) {
	o.ExpirationHours = &v
}

// GetExpires returns the Expires field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetExpires() time.Time {
	if o == nil || IsNil(o.Expires) {
		var ret time.Time
		return ret
	}
	return *o.Expires
}

// GetExpiresOk returns a tuple with the Expires field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetExpiresOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Expires) {
		return nil, false
	}

	return o.Expires, true
}

// HasExpires returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasExpires() bool {
	if o != nil && !IsNil(o.Expires) {
		return true
	}

	return false
}

// SetExpires gets a reference to the given time.Time and assigns it to the Expires field.
func (o *BackupRestoreJobDelivery) SetExpires(v time.Time) {
	o.Expires = &v
}

// GetMaxDownloads returns the MaxDownloads field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetMaxDownloads() int {
	if o == nil || IsNil(o.MaxDownloads) {
		var ret int
		return ret
	}
	return *o.MaxDownloads
}

// GetMaxDownloadsOk returns a tuple with the MaxDownloads field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetMaxDownloadsOk() (*int, bool) {
	if o == nil || IsNil(o.MaxDownloads) {
		return nil, false
	}

	return o.MaxDownloads, true
}

// HasMaxDownloads returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasMaxDownloads() bool {
	if o != nil && !IsNil(o.MaxDownloads) {
		return true
	}

	return false
}

// SetMaxDownloads gets a reference to the given int and assigns it to the MaxDownloads field.
func (o *BackupRestoreJobDelivery) SetMaxDownloads(v int) {
	o.MaxDownloads = &v
}

// GetMethodName returns the MethodName field value
func (o *BackupRestoreJobDelivery) GetMethodName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.MethodName
}

// GetMethodNameOk returns a tuple with the MethodName field value
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetMethodNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MethodName, true
}

// SetMethodName sets field value
func (o *BackupRestoreJobDelivery) SetMethodName(v string) {
	o.MethodName = v
}

// GetStatusName returns the StatusName field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetStatusName() string {
	if o == nil || IsNil(o.StatusName) {
		var ret string
		return ret
	}
	return *o.StatusName
}

// GetStatusNameOk returns a tuple with the StatusName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetStatusNameOk() (*string, bool) {
	if o == nil || IsNil(o.StatusName) {
		return nil, false
	}

	return o.StatusName, true
}

// HasStatusName returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasStatusName() bool {
	if o != nil && !IsNil(o.StatusName) {
		return true
	}

	return false
}

// SetStatusName gets a reference to the given string and assigns it to the StatusName field.
func (o *BackupRestoreJobDelivery) SetStatusName(v string) {
	o.StatusName = &v
}

// GetTargetClusterId returns the TargetClusterId field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetTargetClusterId() string {
	if o == nil || IsNil(o.TargetClusterId) {
		var ret string
		return ret
	}
	return *o.TargetClusterId
}

// GetTargetClusterIdOk returns a tuple with the TargetClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetTargetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.TargetClusterId) {
		return nil, false
	}

	return o.TargetClusterId, true
}

// HasTargetClusterId returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasTargetClusterId() bool {
	if o != nil && !IsNil(o.TargetClusterId) {
		return true
	}

	return false
}

// SetTargetClusterId gets a reference to the given string and assigns it to the TargetClusterId field.
func (o *BackupRestoreJobDelivery) SetTargetClusterId(v string) {
	o.TargetClusterId = &v
}

// GetTargetClusterName returns the TargetClusterName field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetTargetClusterName() string {
	if o == nil || IsNil(o.TargetClusterName) {
		var ret string
		return ret
	}
	return *o.TargetClusterName
}

// GetTargetClusterNameOk returns a tuple with the TargetClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetTargetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.TargetClusterName) {
		return nil, false
	}

	return o.TargetClusterName, true
}

// HasTargetClusterName returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasTargetClusterName() bool {
	if o != nil && !IsNil(o.TargetClusterName) {
		return true
	}

	return false
}

// SetTargetClusterName gets a reference to the given string and assigns it to the TargetClusterName field.
func (o *BackupRestoreJobDelivery) SetTargetClusterName(v string) {
	o.TargetClusterName = &v
}

// GetTargetGroupId returns the TargetGroupId field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetTargetGroupId() string {
	if o == nil || IsNil(o.TargetGroupId) {
		var ret string
		return ret
	}
	return *o.TargetGroupId
}

// GetTargetGroupIdOk returns a tuple with the TargetGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetTargetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TargetGroupId) {
		return nil, false
	}

	return o.TargetGroupId, true
}

// HasTargetGroupId returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasTargetGroupId() bool {
	if o != nil && !IsNil(o.TargetGroupId) {
		return true
	}

	return false
}

// SetTargetGroupId gets a reference to the given string and assigns it to the TargetGroupId field.
func (o *BackupRestoreJobDelivery) SetTargetGroupId(v string) {
	o.TargetGroupId = &v
}

// GetUrl returns the Url field value if set, zero value otherwise
// Deprecated
func (o *BackupRestoreJobDelivery) GetUrl() string {
	if o == nil || IsNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
// Deprecated
func (o *BackupRestoreJobDelivery) GetUrlOk() (*string, bool) {
	if o == nil || IsNil(o.Url) {
		return nil, false
	}

	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasUrl() bool {
	if o != nil && !IsNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
// Deprecated
func (o *BackupRestoreJobDelivery) SetUrl(v string) {
	o.Url = &v
}

// GetUrlV2 returns the UrlV2 field value if set, zero value otherwise
func (o *BackupRestoreJobDelivery) GetUrlV2() string {
	if o == nil || IsNil(o.UrlV2) {
		var ret string
		return ret
	}
	return *o.UrlV2
}

// GetUrlV2Ok returns a tuple with the UrlV2 field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupRestoreJobDelivery) GetUrlV2Ok() (*string, bool) {
	if o == nil || IsNil(o.UrlV2) {
		return nil, false
	}

	return o.UrlV2, true
}

// HasUrlV2 returns a boolean if a field has been set.
func (o *BackupRestoreJobDelivery) HasUrlV2() bool {
	if o != nil && !IsNil(o.UrlV2) {
		return true
	}

	return false
}

// SetUrlV2 gets a reference to the given string and assigns it to the UrlV2 field.
func (o *BackupRestoreJobDelivery) SetUrlV2(v string) {
	o.UrlV2 = &v
}
