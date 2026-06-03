// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"time"
)

// AlertViewForNdsGroup struct for AlertViewForNdsGroup
type AlertViewForNdsGroup struct {
	// Date and time until which this alert has been acknowledged. This parameter expresses its value in the ISO 8601 timestamp format in UTC. The resource returns this parameter if a MongoDB User previously acknowledged this alert.  - To acknowledge this alert forever, set the parameter value to 100 years in the future.  - To unacknowledge a previously acknowledged alert, do not set this parameter value.
	AcknowledgedUntil *time.Time `json:"acknowledgedUntil,omitempty"`
	// Comment that a MongoDB Cloud user submitted when acknowledging the alert.
	AcknowledgementComment *string `json:"acknowledgementComment,omitempty"`
	// MongoDB Cloud username of the person who acknowledged the alert. The response returns this parameter if a MongoDB Cloud user previously acknowledged this alert.
	// Read only field.
	AcknowledgingUsername *string `json:"acknowledgingUsername,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the alert configuration that sets this alert.
	// Read only field.
	AlertConfigId *string `json:"alertConfigId,omitempty"`
	// Date and time when MongoDB Cloud created this alert. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Created *time.Time `json:"created,omitempty"`
	// Incident that triggered this alert.
	// Read only field.
	EventTypeName *string `json:"eventTypeName,omitempty"`
	// Unique 24-hexadecimal digit string that identifies the project that owns this alert.
	// Read only field.
	GroupId *string `json:"groupId,omitempty"`
	// Unique 24-hexadecimal digit string that identifies this alert.
	// Read only field.
	Id *string `json:"id,omitempty"`
	// Date and time that any notifications were last sent for this alert. This parameter expresses its value in the ISO 8601 timestamp format in UTC. The resource returns this parameter if MongoDB Cloud has sent notifications for this alert.
	// Read only field.
	LastNotified *time.Time `json:"lastNotified,omitempty"`
	// List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.
	// Read only field.
	Links *[]Link `json:"links,omitempty"`
	// Unique 24-hexadecimal character string that identifies the organization that owns the project to which this alert applies.
	// Read only field.
	OrgId *string `json:"orgId,omitempty"`
	// Date and time that this alert changed to `\"status\" : \"CLOSED\"`. This parameter expresses its value in the ISO 8601 timestamp format in UTC. The resource returns this parameter once `\"status\" : \"CLOSED\"`.
	// Read only field.
	Resolved *time.Time `json:"resolved,omitempty"`
	// State of this alert at the time you requested its details. TRACKING indicates the alert condition exists but has not persisted for the minimum notification delay. OPEN indicates the alert condition currently exists. CLOSED indicates the alert condition has been resolved.
	// Read only field.
	Status *string `json:"status,omitempty"`
	// Date and time when someone last updated this alert. This parameter expresses its value in the ISO 8601 timestamp format in UTC.
	// Read only field.
	Updated *time.Time `json:"updated,omitempty"`
	// Human-readable label that identifies the cluster to which this alert applies. This resource returns this parameter for alerts of events impacting backups, replica sets, or sharded clusters.
	// Read only field.
	ClusterName *string `json:"clusterName,omitempty"`
	// Hostname and port of the host to which this alert applies. The resource returns this parameter for alerts of events impacting hosts or replica sets.
	// Read only field.
	HostnameAndPort *string `json:"hostnameAndPort,omitempty"`
	// Name of the replica set to which this alert applies. The response returns this parameter for alerts of events impacting backups, hosts, or replica sets.
	// Read only field.
	ReplicaSetName *string            `json:"replicaSetName,omitempty"`
	CurrentValue   *NumberMetricValue `json:"currentValue,omitempty"`
	// Name of the metric against which Atlas checks the configured `metricThreshold.threshold`.  To learn more about the available metrics, see <a href=\"https://www.mongodb.com/docs/atlas/reference/alert-host-metrics/#std-label-measurement-types\" target=\"_blank\">Host Metrics</a>.  **NOTE**: If you set `eventTypeName` to `OUTSIDE_SERVERLESS_METRIC_THRESHOLD`, you can specify only metrics available for serverless. To learn more, see <a href=\"https://dochub.mongodb.org/core/alert-config-serverless-measurements\" target=\"_blank\">Serverless Measurements</a>.
	// Read only field.
	MetricName *string `json:"metricName,omitempty"`
	// List of unique 24-hexadecimal character strings that identify the replica set members that are not in PRIMARY nor SECONDARY state.
	// Read only field.
	NonRunningHostIds *[]string `json:"nonRunningHostIds,omitempty"`
	// Unique 24-hexadecimal character string that identifies the parent cluster to which this alert applies. The parent cluster contains the sharded nodes. MongoDB Cloud returns this parameter only for alerts of events impacting sharded clusters.
	// Read only field.
	ParentClusterId *string `json:"parentClusterId,omitempty"`
	// The name of the Stream Processing Workspace to which this alert applies. The resource returns this parameter for alerts of events impacting Stream Processing Workspaces.
	// Read only field.
	InstanceName *string `json:"instanceName,omitempty"`
	// The error message associated with the Stream Processor to which this alert applies.
	// Read only field.
	ProcessorErrorMsg *string `json:"processorErrorMsg,omitempty"`
	// The name of the Stream Processor to which this alert applies. The resource returns this parameter for alerts of events impacting Stream Processors.
	// Read only field.
	ProcessorName *string `json:"processorName,omitempty"`
	// The state of the Stream Processor to which this alert applies. The resource returns this parameter for alerts of events impacting Stream Processors.
	// Read only field.
	ProcessorState *string `json:"processorState,omitempty"`
}

// NewAlertViewForNdsGroup instantiates a new AlertViewForNdsGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAlertViewForNdsGroup() *AlertViewForNdsGroup {
	this := AlertViewForNdsGroup{}
	return &this
}

// NewAlertViewForNdsGroupWithDefaults instantiates a new AlertViewForNdsGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAlertViewForNdsGroupWithDefaults() *AlertViewForNdsGroup {
	this := AlertViewForNdsGroup{}
	return &this
}

// GetAcknowledgedUntil returns the AcknowledgedUntil field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetAcknowledgedUntil() time.Time {
	if o == nil || IsNil(o.AcknowledgedUntil) {
		var ret time.Time
		return ret
	}
	return *o.AcknowledgedUntil
}

// GetAcknowledgedUntilOk returns a tuple with the AcknowledgedUntil field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetAcknowledgedUntilOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AcknowledgedUntil) {
		return nil, false
	}

	return o.AcknowledgedUntil, true
}

// HasAcknowledgedUntil returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasAcknowledgedUntil() bool {
	if o != nil && !IsNil(o.AcknowledgedUntil) {
		return true
	}

	return false
}

// SetAcknowledgedUntil gets a reference to the given time.Time and assigns it to the AcknowledgedUntil field.
func (o *AlertViewForNdsGroup) SetAcknowledgedUntil(v time.Time) {
	o.AcknowledgedUntil = &v
}

// GetAcknowledgementComment returns the AcknowledgementComment field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetAcknowledgementComment() string {
	if o == nil || IsNil(o.AcknowledgementComment) {
		var ret string
		return ret
	}
	return *o.AcknowledgementComment
}

// GetAcknowledgementCommentOk returns a tuple with the AcknowledgementComment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetAcknowledgementCommentOk() (*string, bool) {
	if o == nil || IsNil(o.AcknowledgementComment) {
		return nil, false
	}

	return o.AcknowledgementComment, true
}

// HasAcknowledgementComment returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasAcknowledgementComment() bool {
	if o != nil && !IsNil(o.AcknowledgementComment) {
		return true
	}

	return false
}

// SetAcknowledgementComment gets a reference to the given string and assigns it to the AcknowledgementComment field.
func (o *AlertViewForNdsGroup) SetAcknowledgementComment(v string) {
	o.AcknowledgementComment = &v
}

// GetAcknowledgingUsername returns the AcknowledgingUsername field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetAcknowledgingUsername() string {
	if o == nil || IsNil(o.AcknowledgingUsername) {
		var ret string
		return ret
	}
	return *o.AcknowledgingUsername
}

// GetAcknowledgingUsernameOk returns a tuple with the AcknowledgingUsername field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetAcknowledgingUsernameOk() (*string, bool) {
	if o == nil || IsNil(o.AcknowledgingUsername) {
		return nil, false
	}

	return o.AcknowledgingUsername, true
}

// HasAcknowledgingUsername returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasAcknowledgingUsername() bool {
	if o != nil && !IsNil(o.AcknowledgingUsername) {
		return true
	}

	return false
}

// SetAcknowledgingUsername gets a reference to the given string and assigns it to the AcknowledgingUsername field.
func (o *AlertViewForNdsGroup) SetAcknowledgingUsername(v string) {
	o.AcknowledgingUsername = &v
}

// GetAlertConfigId returns the AlertConfigId field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetAlertConfigId() string {
	if o == nil || IsNil(o.AlertConfigId) {
		var ret string
		return ret
	}
	return *o.AlertConfigId
}

// GetAlertConfigIdOk returns a tuple with the AlertConfigId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetAlertConfigIdOk() (*string, bool) {
	if o == nil || IsNil(o.AlertConfigId) {
		return nil, false
	}

	return o.AlertConfigId, true
}

// HasAlertConfigId returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasAlertConfigId() bool {
	if o != nil && !IsNil(o.AlertConfigId) {
		return true
	}

	return false
}

// SetAlertConfigId gets a reference to the given string and assigns it to the AlertConfigId field.
func (o *AlertViewForNdsGroup) SetAlertConfigId(v string) {
	o.AlertConfigId = &v
}

// GetCreated returns the Created field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetCreated() time.Time {
	if o == nil || IsNil(o.Created) {
		var ret time.Time
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetCreatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Created) {
		return nil, false
	}

	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasCreated() bool {
	if o != nil && !IsNil(o.Created) {
		return true
	}

	return false
}

// SetCreated gets a reference to the given time.Time and assigns it to the Created field.
func (o *AlertViewForNdsGroup) SetCreated(v time.Time) {
	o.Created = &v
}

// GetEventTypeName returns the EventTypeName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetEventTypeName() string {
	if o == nil || IsNil(o.EventTypeName) {
		var ret string
		return ret
	}
	return *o.EventTypeName
}

// GetEventTypeNameOk returns a tuple with the EventTypeName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetEventTypeNameOk() (*string, bool) {
	if o == nil || IsNil(o.EventTypeName) {
		return nil, false
	}

	return o.EventTypeName, true
}

// HasEventTypeName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasEventTypeName() bool {
	if o != nil && !IsNil(o.EventTypeName) {
		return true
	}

	return false
}

// SetEventTypeName gets a reference to the given string and assigns it to the EventTypeName field.
func (o *AlertViewForNdsGroup) SetEventTypeName(v string) {
	o.EventTypeName = &v
}

// GetGroupId returns the GroupId field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetGroupId() string {
	if o == nil || IsNil(o.GroupId) {
		var ret string
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}

	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given string and assigns it to the GroupId field.
func (o *AlertViewForNdsGroup) SetGroupId(v string) {
	o.GroupId = &v
}

// GetId returns the Id field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}

	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *AlertViewForNdsGroup) SetId(v string) {
	o.Id = &v
}

// GetLastNotified returns the LastNotified field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetLastNotified() time.Time {
	if o == nil || IsNil(o.LastNotified) {
		var ret time.Time
		return ret
	}
	return *o.LastNotified
}

// GetLastNotifiedOk returns a tuple with the LastNotified field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetLastNotifiedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastNotified) {
		return nil, false
	}

	return o.LastNotified, true
}

// HasLastNotified returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasLastNotified() bool {
	if o != nil && !IsNil(o.LastNotified) {
		return true
	}

	return false
}

// SetLastNotified gets a reference to the given time.Time and assigns it to the LastNotified field.
func (o *AlertViewForNdsGroup) SetLastNotified(v time.Time) {
	o.LastNotified = &v
}

// GetLinks returns the Links field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetLinks() []Link {
	if o == nil || IsNil(o.Links) {
		var ret []Link
		return ret
	}
	return *o.Links
}

// GetLinksOk returns a tuple with the Links field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetLinksOk() (*[]Link, bool) {
	if o == nil || IsNil(o.Links) {
		return nil, false
	}

	return o.Links, true
}

// HasLinks returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasLinks() bool {
	if o != nil && !IsNil(o.Links) {
		return true
	}

	return false
}

// SetLinks gets a reference to the given []Link and assigns it to the Links field.
func (o *AlertViewForNdsGroup) SetLinks(v []Link) {
	o.Links = &v
}

// GetOrgId returns the OrgId field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetOrgId() string {
	if o == nil || IsNil(o.OrgId) {
		var ret string
		return ret
	}
	return *o.OrgId
}

// GetOrgIdOk returns a tuple with the OrgId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetOrgIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrgId) {
		return nil, false
	}

	return o.OrgId, true
}

// HasOrgId returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasOrgId() bool {
	if o != nil && !IsNil(o.OrgId) {
		return true
	}

	return false
}

// SetOrgId gets a reference to the given string and assigns it to the OrgId field.
func (o *AlertViewForNdsGroup) SetOrgId(v string) {
	o.OrgId = &v
}

// GetResolved returns the Resolved field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetResolved() time.Time {
	if o == nil || IsNil(o.Resolved) {
		var ret time.Time
		return ret
	}
	return *o.Resolved
}

// GetResolvedOk returns a tuple with the Resolved field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetResolvedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Resolved) {
		return nil, false
	}

	return o.Resolved, true
}

// HasResolved returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasResolved() bool {
	if o != nil && !IsNil(o.Resolved) {
		return true
	}

	return false
}

// SetResolved gets a reference to the given time.Time and assigns it to the Resolved field.
func (o *AlertViewForNdsGroup) SetResolved(v time.Time) {
	o.Resolved = &v
}

// GetStatus returns the Status field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}

	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *AlertViewForNdsGroup) SetStatus(v string) {
	o.Status = &v
}

// GetUpdated returns the Updated field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetUpdated() time.Time {
	if o == nil || IsNil(o.Updated) {
		var ret time.Time
		return ret
	}
	return *o.Updated
}

// GetUpdatedOk returns a tuple with the Updated field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetUpdatedOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Updated) {
		return nil, false
	}

	return o.Updated, true
}

// HasUpdated returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasUpdated() bool {
	if o != nil && !IsNil(o.Updated) {
		return true
	}

	return false
}

// SetUpdated gets a reference to the given time.Time and assigns it to the Updated field.
func (o *AlertViewForNdsGroup) SetUpdated(v time.Time) {
	o.Updated = &v
}

// GetClusterName returns the ClusterName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetClusterName() string {
	if o == nil || IsNil(o.ClusterName) {
		var ret string
		return ret
	}
	return *o.ClusterName
}

// GetClusterNameOk returns a tuple with the ClusterName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetClusterNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterName) {
		return nil, false
	}

	return o.ClusterName, true
}

// HasClusterName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasClusterName() bool {
	if o != nil && !IsNil(o.ClusterName) {
		return true
	}

	return false
}

// SetClusterName gets a reference to the given string and assigns it to the ClusterName field.
func (o *AlertViewForNdsGroup) SetClusterName(v string) {
	o.ClusterName = &v
}

// GetHostnameAndPort returns the HostnameAndPort field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetHostnameAndPort() string {
	if o == nil || IsNil(o.HostnameAndPort) {
		var ret string
		return ret
	}
	return *o.HostnameAndPort
}

// GetHostnameAndPortOk returns a tuple with the HostnameAndPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetHostnameAndPortOk() (*string, bool) {
	if o == nil || IsNil(o.HostnameAndPort) {
		return nil, false
	}

	return o.HostnameAndPort, true
}

// HasHostnameAndPort returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasHostnameAndPort() bool {
	if o != nil && !IsNil(o.HostnameAndPort) {
		return true
	}

	return false
}

// SetHostnameAndPort gets a reference to the given string and assigns it to the HostnameAndPort field.
func (o *AlertViewForNdsGroup) SetHostnameAndPort(v string) {
	o.HostnameAndPort = &v
}

// GetReplicaSetName returns the ReplicaSetName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetReplicaSetName() string {
	if o == nil || IsNil(o.ReplicaSetName) {
		var ret string
		return ret
	}
	return *o.ReplicaSetName
}

// GetReplicaSetNameOk returns a tuple with the ReplicaSetName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetReplicaSetNameOk() (*string, bool) {
	if o == nil || IsNil(o.ReplicaSetName) {
		return nil, false
	}

	return o.ReplicaSetName, true
}

// HasReplicaSetName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasReplicaSetName() bool {
	if o != nil && !IsNil(o.ReplicaSetName) {
		return true
	}

	return false
}

// SetReplicaSetName gets a reference to the given string and assigns it to the ReplicaSetName field.
func (o *AlertViewForNdsGroup) SetReplicaSetName(v string) {
	o.ReplicaSetName = &v
}

// GetCurrentValue returns the CurrentValue field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetCurrentValue() NumberMetricValue {
	if o == nil || IsNil(o.CurrentValue) {
		var ret NumberMetricValue
		return ret
	}
	return *o.CurrentValue
}

// GetCurrentValueOk returns a tuple with the CurrentValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetCurrentValueOk() (*NumberMetricValue, bool) {
	if o == nil || IsNil(o.CurrentValue) {
		return nil, false
	}

	return o.CurrentValue, true
}

// HasCurrentValue returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasCurrentValue() bool {
	if o != nil && !IsNil(o.CurrentValue) {
		return true
	}

	return false
}

// SetCurrentValue gets a reference to the given NumberMetricValue and assigns it to the CurrentValue field.
func (o *AlertViewForNdsGroup) SetCurrentValue(v NumberMetricValue) {
	o.CurrentValue = &v
}

// GetMetricName returns the MetricName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetMetricName() string {
	if o == nil || IsNil(o.MetricName) {
		var ret string
		return ret
	}
	return *o.MetricName
}

// GetMetricNameOk returns a tuple with the MetricName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetMetricNameOk() (*string, bool) {
	if o == nil || IsNil(o.MetricName) {
		return nil, false
	}

	return o.MetricName, true
}

// HasMetricName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasMetricName() bool {
	if o != nil && !IsNil(o.MetricName) {
		return true
	}

	return false
}

// SetMetricName gets a reference to the given string and assigns it to the MetricName field.
func (o *AlertViewForNdsGroup) SetMetricName(v string) {
	o.MetricName = &v
}

// GetNonRunningHostIds returns the NonRunningHostIds field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetNonRunningHostIds() []string {
	if o == nil || IsNil(o.NonRunningHostIds) {
		var ret []string
		return ret
	}
	return *o.NonRunningHostIds
}

// GetNonRunningHostIdsOk returns a tuple with the NonRunningHostIds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetNonRunningHostIdsOk() (*[]string, bool) {
	if o == nil || IsNil(o.NonRunningHostIds) {
		return nil, false
	}

	return o.NonRunningHostIds, true
}

// HasNonRunningHostIds returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasNonRunningHostIds() bool {
	if o != nil && !IsNil(o.NonRunningHostIds) {
		return true
	}

	return false
}

// SetNonRunningHostIds gets a reference to the given []string and assigns it to the NonRunningHostIds field.
func (o *AlertViewForNdsGroup) SetNonRunningHostIds(v []string) {
	o.NonRunningHostIds = &v
}

// GetParentClusterId returns the ParentClusterId field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetParentClusterId() string {
	if o == nil || IsNil(o.ParentClusterId) {
		var ret string
		return ret
	}
	return *o.ParentClusterId
}

// GetParentClusterIdOk returns a tuple with the ParentClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetParentClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ParentClusterId) {
		return nil, false
	}

	return o.ParentClusterId, true
}

// HasParentClusterId returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasParentClusterId() bool {
	if o != nil && !IsNil(o.ParentClusterId) {
		return true
	}

	return false
}

// SetParentClusterId gets a reference to the given string and assigns it to the ParentClusterId field.
func (o *AlertViewForNdsGroup) SetParentClusterId(v string) {
	o.ParentClusterId = &v
}

// GetInstanceName returns the InstanceName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetInstanceName() string {
	if o == nil || IsNil(o.InstanceName) {
		var ret string
		return ret
	}
	return *o.InstanceName
}

// GetInstanceNameOk returns a tuple with the InstanceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetInstanceNameOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceName) {
		return nil, false
	}

	return o.InstanceName, true
}

// HasInstanceName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasInstanceName() bool {
	if o != nil && !IsNil(o.InstanceName) {
		return true
	}

	return false
}

// SetInstanceName gets a reference to the given string and assigns it to the InstanceName field.
func (o *AlertViewForNdsGroup) SetInstanceName(v string) {
	o.InstanceName = &v
}

// GetProcessorErrorMsg returns the ProcessorErrorMsg field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetProcessorErrorMsg() string {
	if o == nil || IsNil(o.ProcessorErrorMsg) {
		var ret string
		return ret
	}
	return *o.ProcessorErrorMsg
}

// GetProcessorErrorMsgOk returns a tuple with the ProcessorErrorMsg field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetProcessorErrorMsgOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorErrorMsg) {
		return nil, false
	}

	return o.ProcessorErrorMsg, true
}

// HasProcessorErrorMsg returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasProcessorErrorMsg() bool {
	if o != nil && !IsNil(o.ProcessorErrorMsg) {
		return true
	}

	return false
}

// SetProcessorErrorMsg gets a reference to the given string and assigns it to the ProcessorErrorMsg field.
func (o *AlertViewForNdsGroup) SetProcessorErrorMsg(v string) {
	o.ProcessorErrorMsg = &v
}

// GetProcessorName returns the ProcessorName field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetProcessorName() string {
	if o == nil || IsNil(o.ProcessorName) {
		var ret string
		return ret
	}
	return *o.ProcessorName
}

// GetProcessorNameOk returns a tuple with the ProcessorName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetProcessorNameOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorName) {
		return nil, false
	}

	return o.ProcessorName, true
}

// HasProcessorName returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasProcessorName() bool {
	if o != nil && !IsNil(o.ProcessorName) {
		return true
	}

	return false
}

// SetProcessorName gets a reference to the given string and assigns it to the ProcessorName field.
func (o *AlertViewForNdsGroup) SetProcessorName(v string) {
	o.ProcessorName = &v
}

// GetProcessorState returns the ProcessorState field value if set, zero value otherwise
func (o *AlertViewForNdsGroup) GetProcessorState() string {
	if o == nil || IsNil(o.ProcessorState) {
		var ret string
		return ret
	}
	return *o.ProcessorState
}

// GetProcessorStateOk returns a tuple with the ProcessorState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AlertViewForNdsGroup) GetProcessorStateOk() (*string, bool) {
	if o == nil || IsNil(o.ProcessorState) {
		return nil, false
	}

	return o.ProcessorState, true
}

// HasProcessorState returns a boolean if a field has been set.
func (o *AlertViewForNdsGroup) HasProcessorState() bool {
	if o != nil && !IsNil(o.ProcessorState) {
		return true
	}

	return false
}

// SetProcessorState gets a reference to the given string and assigns it to the ProcessorState field.
func (o *AlertViewForNdsGroup) SetProcessorState(v string) {
	o.ProcessorState = &v
}
