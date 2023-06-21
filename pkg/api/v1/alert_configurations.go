package v1

import (
	"fmt"
	"strconv"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type AlertConfiguration struct {
	// If omitted, the configuration is disabled.
	Enabled bool `json:"enabled,omitempty"`
	// The type of event that will trigger an alert.
	EventTypeName string `json:"eventTypeName,omitempty"`
	// You can filter using the matchers array only when the EventTypeName specifies an event for a host, replica set, or sharded cluster.
	Matchers []Matcher `json:"matchers,omitempty"`
	// Threshold  causes an alert to be triggered.
	Threshold *Threshold `json:"threshold,omitempty"`
	// Notifications are sending when an alert condition is detected.
	Notifications []Notification `json:"notifications,omitempty"`
	// MetricThreshold  causes an alert to be triggered.
	MetricThreshold *MetricThreshold `json:"metricThreshold,omitempty"`
}

func (in *AlertConfiguration) ToAtlas() (*mongodbatlas.AlertConfiguration, error) {
	if in == nil {
		return nil, nil
	}
	// Some field can be converted directly
	result := &mongodbatlas.AlertConfiguration{
		Enabled:       &in.Enabled,
		EventTypeName: in.EventTypeName,
	}

	for _, m := range in.Matchers {
		matcher := mongodbatlas.Matcher{
			FieldName: m.FieldName,
			Operator:  m.Operator,
			Value:     m.Value,
		}
		result.Matchers = append(result.Matchers, matcher)
	}

	for _, n := range in.Notifications {
		notification := &mongodbatlas.Notification{}
		err := compat.JSONCopy(&notification, n)
		if err != nil {
			return nil, err
		}
		result.Notifications = append(result.Notifications, *notification)
	}

	// Some fields require special conversion
	tr, err := in.Threshold.ToAtlas()
	if err != nil {
		return nil, err
	}
	result.Threshold = tr
	metricThreshold, err := in.MetricThreshold.ToAtlas()
	if err != nil {
		return nil, err
	}
	result.MetricThreshold = metricThreshold
	return result, err
}

type Matcher struct {
	// Name of the field in the target object to match on.
	FieldName string `json:"fieldName,omitempty"`
	// The operator to test the field’s value.
	Operator string `json:"operator,omitempty"`
	// Value to test with the specified operator.
	Value string `json:"value,omitempty"`
}

func (in *Matcher) IsEqual(matcher mongodbatlas.Matcher) bool {
	if in == nil {
		return false
	}
	return in.FieldName == matcher.FieldName &&
		in.Operator == matcher.Operator &&
		in.Value == matcher.Value
}

type Threshold struct {
	// Operator to apply when checking the current metric value against the threshold value. it accepts the following values: GREATER_THAN, LESS_THAN
	Operator string `json:"operator,omitempty"`
	// The units for the threshold value
	Units string `json:"units,omitempty"`
	// Threshold value outside which an alert will be triggered.
	Threshold string `json:"threshold,omitempty"`
}

func (in *Threshold) IsEqual(threshold *mongodbatlas.Threshold) bool {
	logger := zap.NewExample().Sugar()
	if in == nil {
		return threshold == nil
	}
	if threshold == nil {
		return false
	}
	logger.Debugf("threshold: %v", threshold)
	if in.Operator != threshold.Operator {
		logger.Debugf("operator: %s != %s", in.Operator, threshold.Operator)
		return false
	}
	if in.Units != threshold.Units {
		logger.Debugf("units: %s != %s", in.Units, threshold.Units)
		return false
	}
	if in.Threshold != strconv.FormatFloat(threshold.Threshold, 'f', -1, 64) {
		logger.Debugf("threshold: %s != %s", in.Threshold, strconv.FormatFloat(threshold.Threshold, 'f', -1, 64))
		return false
	}
	return true
}

func (in *Threshold) ToAtlas() (*mongodbatlas.Threshold, error) {
	if in == nil {
		return nil, nil
	}
	tr, err := strconv.ParseFloat(in.Threshold, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse threshold value: %w. should be float", err)
	}
	result := &mongodbatlas.Threshold{
		Operator:  in.Operator,
		Units:     in.Units,
		Threshold: tr,
	}
	return result, nil
}

type Notification struct {
	// Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	APIToken string `json:"apiToken,omitempty"`
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// Slack channel name. Populated for the SLACK notifications type.
	ChannelName string `json:"channelName,omitempty"`
	// Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.
	DatadogAPIKey string `json:"datadogApiKey,omitempty"`
	// +optional
	DatadogAPIKeyRef common.ResourceRefNamespaced `json:"datadogAPIKeyRef,omitempty"`
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
	// +optional
	FlowdockAPITokenRef common.ResourceRefNamespaced `json:"flowdockApiTokenRef,omitempty"`
	// Flowdock flow namse in lower-case letters.
	FlowName string `json:"flowName,omitempty"`
	// Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.
	IntervalMin int `json:"intervalMin,omitempty"`
	// Mobile number to which alert notifications are sent. Populated for the SMS notifications type.
	MobileNumber string `json:"mobileNumber,omitempty"`
	// Opsgenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	OpsGenieAPIKey string `json:"opsGenieApiKey,omitempty"`
	// +optional
	OpsGenieAPIKeyRef common.ResourceRefNamespaced `json:"opsGenieApiKeyRef,omitempty"`
	// Region that indicates which API URL to use.
	OpsGenieRegion string `json:"opsGenieRegion,omitempty"`
	// Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.
	OrgName string `json:"orgName,omitempty"`
	// PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	ServiceKey string `json:"serviceKey,omitempty"`
	// +optinal
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
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
	// +optional
	// Secret for VictorOps should contain both APIKey and RoutingKey values
	VictorOpsSecretRef common.ResourceRefNamespaced `json:"victorOpsSecretRef,omitempty"`
	// VictorOps routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	VictorOpsRoutingKey string `json:"victorOpsRoutingKey,omitempty"`
	// The following roles grant privileges within a project.
	Roles []string `json:"roles,omitempty"`
}

func (in *Notification) IsEqual(notification mongodbatlas.Notification) bool {
	if in == nil {
		return false
	}
	if in.APIToken != notification.APIToken ||
		in.ChannelName != notification.ChannelName ||
		in.DatadogAPIKey != notification.DatadogAPIKey ||
		in.DatadogRegion != notification.DatadogRegion ||
		!util.PtrValuesEqual(in.DelayMin, notification.DelayMin) ||
		in.EmailAddress != notification.EmailAddress ||
		!util.PtrValuesEqual(in.EmailEnabled, notification.EmailEnabled) ||
		in.FlowdockAPIToken != notification.FlowdockAPIToken ||
		in.FlowName != notification.FlowName ||
		in.IntervalMin != notification.IntervalMin ||
		in.MobileNumber != notification.MobileNumber ||
		in.OpsGenieAPIKey != notification.OpsGenieAPIKey ||
		in.OpsGenieRegion != notification.OpsGenieRegion ||
		in.OrgName != notification.OrgName ||
		in.ServiceKey != notification.ServiceKey ||
		!util.PtrValuesEqual(in.SMSEnabled, notification.SMSEnabled) ||
		in.TeamID != notification.TeamID ||
		in.TeamName != notification.TeamName ||
		in.TypeName != notification.TypeName ||
		in.Username != notification.Username ||
		in.VictorOpsAPIKey != notification.VictorOpsAPIKey ||
		in.VictorOpsRoutingKey != notification.VictorOpsRoutingKey {
		return false
	}

	if !util.IsEqualWithoutOrder(in.Roles, notification.Roles) {
		return false
	}

	return true
}

// MetricThreshold  causes an alert to be triggered. Required if "eventTypeName" : "OUTSIDE_METRIC_THRESHOLD".
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

func (in *MetricThreshold) IsEqual(threshold *mongodbatlas.MetricThreshold) bool {
	if in == nil {
		return threshold == nil
	}
	if threshold == nil {
		return false
	}
	return in.MetricName == threshold.MetricName &&
		in.Operator == threshold.Operator &&
		in.Threshold == strconv.FormatFloat(threshold.Threshold, 'f', -1, 64) &&
		in.Units == threshold.Units &&
		in.Mode == threshold.Mode
}

func (in *MetricThreshold) ToAtlas() (*mongodbatlas.MetricThreshold, error) {
	if in == nil {
		return nil, nil
	}
	tr, err := strconv.ParseFloat(in.Threshold, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse threshold value: %w. should be float", err)
	}
	result := &mongodbatlas.MetricThreshold{
		MetricName: in.MetricName,
		Operator:   in.Operator,
		Threshold:  tr,
		Units:      in.Units,
		Mode:       in.Mode,
	}
	return result, nil
}
