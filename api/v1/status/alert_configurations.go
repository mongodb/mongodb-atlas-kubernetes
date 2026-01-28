// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package status

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

type AlertConfiguration struct {
	// Unique identifier.
	ID string `json:"id,omitempty"`
	// ErrorMessage is massage if the alert configuration is in an incorrect state.
	ErrorMessage string `json:"errorMessage,omitempty"`
	// Unique identifier of the project that owns this alert configuration.
	GroupID string `json:"groupId,omitempty"`
	// ID of the alert configuration that triggered this alert.
	AlertConfigID string `json:"alertConfigId,omitempty"`
	// The type of event that will trigger an alert.
	EventTypeName string `json:"eventTypeName,omitempty"`
	// Timestamp in ISO 8601 date and time format in UTC when this alert configuration was created.
	Created string `json:"created,omitempty"`
	// The current state of the alert. Possible values are: TRACKING, OPEN, CLOSED, CANCELED
	Status string `json:"status,omitempty"`
	// The date through which the alert has been acknowledged. Will not be present if the alert has never been acknowledged.
	AcknowledgedUntil string `json:"acknowledgedUntil,omitempty"`
	// The comment left by the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	AcknowledgementComment string `json:"acknowledgementComment,omitempty"`
	// The username of the user who acknowledged the alert. Will not be present if the alert has never been acknowledged.
	AcknowledgingUsername string `json:"acknowledgingUsername,omitempty"`
	// Timestamp in ISO 8601 date and time format in UTC when this alert configuration was last updated.
	Updated string `json:"updated,omitempty"`
	// When the alert was closed. Only present if the status is CLOSED.
	Resolved string `json:"resolved,omitempty"`
	// When the last notification was sent for this alert. Only present if notifications have been sent.
	LastNotified string `json:"lastNotified,omitempty"`
	// The hostname and port of each host to which the alert applies. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.
	HostnameAndPort string `json:"hostnameAndPort,omitempty"`
	// ID of the host to which the metric pertains. Only present for alerts of type HOST, HOST_METRIC, and REPLICA_SET.
	HostID string `json:"hostId,omitempty"`
	// Name of the replica set. Only present for alerts of type HOST, HOST_METRIC, BACKUP, and REPLICA_SET.
	ReplicaSetName string `json:"replicaSetName,omitempty"`
	// The name of the measurement whose value went outside the threshold. Only present if eventTypeName is set to OUTSIDE_METRIC_THRESHOLD.
	MetricName string `json:"metricName,omitempty"`
	// If omitted, the configuration is disabled.
	Enabled *bool `json:"enabled,omitempty"`
	// The ID of the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.
	ClusterID string `json:"clusterId,omitempty"`
	// The name the cluster to which this alert applies. Only present for alerts of type BACKUP, REPLICA_SET, and CLUSTER.
	ClusterName string `json:"clusterName,omitempty"`
	// For alerts of the type BACKUP, the type of server being backed up.
	SourceTypeName string `json:"sourceTypeName,omitempty"`
	// CurrentValue represents current value of the metric that triggered the alert. Only present for alerts of type HOST_METRIC.
	CurrentValue *CurrentValue `json:"currentValue,omitempty"`
	// You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.
	Matchers []Matcher `json:"matchers,omitempty"`
	// MetricThreshold  causes an alert to be triggered.
	MetricThreshold *MetricThreshold `json:"metricThreshold,omitempty"`
	// Threshold  causes an alert to be triggered.
	Threshold *Threshold `json:"threshold,omitempty"`
	// Notifications are sending when an alert condition is detected.
	Notifications []Notification `json:"notifications,omitempty"`
	// Severity of the alert.
	SeverityOverride string `json:"severityOverride,omitempty"`
}

type Notification struct {
	// Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	APIToken string `json:"apiToken,omitempty"`
	// Slack channel name. Populated for the SLACK notifications type.
	ChannelName string `json:"channelName,omitempty"`
	// Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.
	DatadogAPIKey string `json:"datadogApiKey,omitempty"`
	// Region that indicates which API URL to use
	DatadogRegion string `json:"datadogRegion,omitempty"`
	// Number of minutes to wait after an alert condition is detected before sending out the first notification.
	DelayMin *int `json:"delayMin,omitempty"`
	// Email address to which alert notifications are sent. Populated for the EMAIL notifications type.
	EmailAddress string `json:"emailAddress,omitempty"`
	// Flag indicating if email notifications should be sent. Populated for ORG, GROUP, and USER notifications types.
	EmailEnabled *bool `json:"emailEnabled,omitempty"`
	// The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	FlowdockAPIToken string `json:"flowdockApiToken,omitempty"`
	// Flowdock flow namse in lower-case letters.
	FlowName string `json:"flowName,omitempty"`
	// Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.
	IntervalMin int `json:"intervalMin,omitempty"`
	// Mobile number to which alert notifications are sent. Populated for the SMS notifications type.
	MobileNumber string `json:"mobileNumber,omitempty"`
	// Opsgenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	OpsGenieAPIKey string `json:"opsGenieApiKey,omitempty"`
	// Region that indicates which API URL to use.
	OpsGenieRegion string `json:"opsGenieRegion,omitempty"`
	// Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.
	OrgName string `json:"orgName,omitempty"`
	// PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	ServiceKey string `json:"serviceKey,omitempty"`
	// Flag indicating if text message notifications should be sent. Populated for ORG, GROUP, and USER notifications types.
	SMSEnabled *bool `json:"smsEnabled,omitempty"`
	// Unique identifier of a team.
	TeamID string `json:"teamId,omitempty"`
	// Label for the team that receives this notification.
	TeamName string `json:"teamName,omitempty"`
	// Type of alert notification.
	TypeName string `json:"typeName,omitempty"`
	// Name of the Atlas user to which to send notifications. Only a user in the project that owns the alert configuration is allowed here. Populated for the USER notifications type.
	Username string `json:"username,omitempty"`
	// VictorOps API key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	VictorOpsAPIKey string `json:"victorOpsApiKey,omitempty"`
	// VictorOps routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	VictorOpsRoutingKey string `json:"victorOpsRoutingKey,omitempty"`
	// The following roles grant privileges within a project.
	Roles []string `json:"roles,omitempty"`
}

func NotificationFromAtlas(notification admin.AlertsNotificationRootForGroup) Notification {
	return Notification{
		APIToken:       notification.GetApiToken(),
		ChannelName:    notification.GetChannelName(),
		DatadogRegion:  notification.GetDatadogRegion(),
		DelayMin:       notification.DelayMin,
		EmailAddress:   notification.GetEmailAddress(),
		EmailEnabled:   notification.EmailEnabled,
		IntervalMin:    notification.GetIntervalMin(),
		MobileNumber:   notification.GetMobileNumber(),
		OpsGenieRegion: notification.GetOpsGenieRegion(),
		ServiceKey:     notification.GetServiceKey(),
		SMSEnabled:     notification.SmsEnabled,
		TeamID:         notification.GetTeamId(),
		TeamName:       notification.GetTeamName(),
		TypeName:       notification.GetTypeName(),
		Username:       notification.GetUsername(),
		Roles:          notification.GetRoles(),
	}
}

type Threshold struct {
	// Operator to apply when checking the current metric value against the threshold value. it accepts the following values: GREATER_THAN, LESS_THAN
	Operator string `json:"operator,omitempty"`
	// The units for the threshold value
	Units string `json:"units,omitempty"`
	// Threshold value outside which an alert will be triggered.
	Threshold string `json:"threshold,omitempty"`
}

func ThresholdFromAtlas(threshold *admin.StreamProcessorMetricThreshold) *Threshold {
	if threshold == nil {
		return nil
	}

	return &Threshold{
		Operator:  threshold.GetOperator(),
		Units:     threshold.GetUnits(),
		Threshold: strconv.FormatInt(int64(threshold.GetThreshold()), 10),
	}
}

type MetricThreshold struct {
	// Name of the metric to check.
	MetricName string `json:"metricName,omitempty"`
	// Operator to apply when checking the current metric value against the threshold value.
	Operator string `json:"operator,omitempty"`
	// Threshold value outside which an alert will be triggered.
	Threshold string `json:"threshold"`
	// The units for the threshold value.
	Units string `json:"units,omitempty"`
	// This must be set to AVERAGE. Atlas computes the current metric value as an average.
	Mode string `json:"mode,omitempty"`
}

func MetricThresholdFromAtlas(threshold *admin.FlexClusterMetricThreshold) *MetricThreshold {
	if threshold == nil {
		return nil
	}

	metricThreshold := &MetricThreshold{}
	metricThreshold.MetricName = threshold.GetMetricName()
	metricThreshold.Operator = threshold.GetOperator()
	metricThreshold.Threshold = strconv.FormatFloat(threshold.GetThreshold(), 'f', -1, 64)
	metricThreshold.Units = threshold.GetUnits()
	metricThreshold.Mode = threshold.GetMode()

	return metricThreshold
}

type Matcher struct {
	// Name of the field in the target object to match on.
	FieldName string `json:"fieldName,omitempty"`
	// The operator to test the field’s value.
	Operator string `json:"operator,omitempty"`
	// Value to test with the specified operator.
	Value string `json:"value,omitempty"`
}

func ParseAlertConfiguration(alertConfiguration admin.GroupAlertsConfig, logger *zap.SugaredLogger) AlertConfiguration {
	status := AlertConfiguration{
		ID:               alertConfiguration.GetId(),
		GroupID:          alertConfiguration.GetGroupId(),
		EventTypeName:    alertConfiguration.GetEventTypeName(),
		Created:          timeutil.FormatISO8601(alertConfiguration.GetCreated()),
		Updated:          timeutil.FormatISO8601(alertConfiguration.GetUpdated()),
		Enabled:          alertConfiguration.Enabled,
		SeverityOverride: alertConfiguration.GetSeverityOverride(),
	}

	if unstructuredMatchers, ok := alertConfiguration.GetMatchersOk(); ok {
		var matchers []Matcher
		err := compat.JSONCopy(matchers, *unstructuredMatchers)
		if err != nil {
			logger.Errorf("unable to convert matchers to structured type: %s", err)
		}
		status.Matchers = matchers
	}

	mThreshold := alertConfiguration.GetMetricThreshold()
	status.MetricThreshold = MetricThresholdFromAtlas(&mThreshold)

	threshold := alertConfiguration.GetThreshold()
	status.Threshold = ThresholdFromAtlas(&threshold)

	if notifications, ok := alertConfiguration.GetNotificationsOk(); ok {
		status.Notifications = make([]Notification, 0, len(*notifications))
		for _, notification := range *notifications {
			notificationFromAtlas := NotificationFromAtlas(notification)
			status.Notifications = append(status.Notifications, notificationFromAtlas)
		}
	}

	return status
}

type CurrentValue struct {
	// The value of the metric.
	Number string `json:"number,omitempty"`
	// The units for the value. Depends on the type of metric.
	Units string `json:"units,omitempty"`
}

func NewFailedParseAlertConfigStatus(errorMessage string, jsonSpec string) AlertConfiguration {
	result := AlertConfiguration{}
	err := json.Unmarshal([]byte(jsonSpec), &result)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Error parsing jsonSpec: %s. error %s", jsonSpec, err)
		return result
	}
	result.ErrorMessage = errorMessage
	return result
}

func NewIncorrectAlertConfigStatus(errorMessage string, alertConfig *admin.GroupAlertsConfig, logger *zap.SugaredLogger) AlertConfiguration {
	if alertConfig == nil {
		return AlertConfiguration{
			ErrorMessage: fmt.Sprintf("Error: %s. alertConfig is nil", errorMessage),
		}
	}
	result := ParseAlertConfiguration(*alertConfig, logger)
	result.ErrorMessage = errorMessage
	return result
}
