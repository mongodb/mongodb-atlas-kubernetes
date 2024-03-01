package v1

import (
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	internalcmp "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compare"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
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

func (in AlertConfiguration) Key() string {
	return strconv.FormatBool(in.Enabled) +
		in.EventTypeName + "|" +
		internalcmp.SliceKey(in.Matchers) + "|" +
		internalcmp.PointerKey(in.Threshold) + "|" +
		internalcmp.SliceKey(in.Notifications) + "|" +
		internalcmp.PointerKey(in.MetricThreshold)
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
		notification, err := n.ToAtlas()
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

func (in Matcher) Key() string {
	return in.FieldName + "|" + in.Operator + "|" + in.Value
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

func (in Threshold) Key() string {
	return in.Operator + "|" + in.Units + "|" + in.Threshold
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
	apiToken string
	// Secret containing a Slack API token or Bot token. Populated for the SLACK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// Slack channel name. Populated for the SLACK notifications type.
	ChannelName   string `json:"channelName,omitempty"`
	datadogAPIKey string
	// Secret containing a Datadog API Key. Found in the Datadog dashboard. Populated for the DATADOG notifications type.
	// +optional
	DatadogAPIKeyRef common.ResourceRefNamespaced `json:"datadogAPIKeyRef,omitempty"`
	// Region that indicates which API URL to use
	DatadogRegion string `json:"datadogRegion,omitempty"`
	// Number of minutes to wait after an alert condition is detected before sending out the first notification.
	DelayMin *int `json:"delayMin,omitempty"`
	// Email address to which alert notifications are sent. Populated for the EMAIL notifications type.
	EmailAddress string `json:"emailAddress,omitempty"`
	// Flag indicating if email notifications should be sent. Populated for ORG, GROUP, and USER notifications types.
	EmailEnabled     *bool `json:"emailEnabled,omitempty"`
	flowdockAPIToken string
	// The Flowdock personal API token. Populated for the FLOWDOCK notifications type. If the token later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	// +optional
	FlowdockAPITokenRef common.ResourceRefNamespaced `json:"flowdockApiTokenRef,omitempty"`
	// Flowdock flow name in lower-case letters.
	FlowName string `json:"flowName,omitempty"`
	// Number of minutes to wait between successive notifications for unacknowledged alerts that are not resolved.
	IntervalMin int `json:"intervalMin,omitempty"`
	// Mobile number to which alert notifications are sent. Populated for the SMS notifications type.
	MobileNumber   string `json:"mobileNumber,omitempty"`
	opsGenieAPIKey string
	// OpsGenie API Key. Populated for the OPS_GENIE notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the token.
	// +optional
	OpsGenieAPIKeyRef common.ResourceRefNamespaced `json:"opsGenieApiKeyRef,omitempty"`
	// Region that indicates which API URL to use.
	OpsGenieRegion string `json:"opsGenieRegion,omitempty"`
	// Flowdock organization name in lower-case letters. This is the name that appears after www.flowdock.com/app/ in the URL string. Populated for the FLOWDOCK notifications type.
	OrgName    string `json:"orgName,omitempty"`
	serviceKey string
	// PagerDuty service key. Populated for the PAGER_DUTY notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	// +optional
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
	Username            string `json:"username,omitempty"`
	victorOpsAPIKey     string
	victorOpsRoutingKey string
	// Secret containing a VictorOps API key and Routing key. Populated for the VICTOR_OPS notifications type. If the key later becomes invalid, Atlas sends an email to the project owner and eventually removes the key.
	// +optional
	VictorOpsSecretRef common.ResourceRefNamespaced `json:"victorOpsSecretRef,omitempty"`
	// The following roles grant privileges within a project.
	Roles []string `json:"roles,omitempty"`
}

func (in Notification) Key() string {
	return in.APITokenRef.Key() + "|" +
		in.ChannelName + "|" +
		in.DatadogAPIKeyRef.Key() + "|" +
		in.DatadogRegion + "|" +
		strconv.Itoa(admin.GetOrDefault(in.DelayMin, 0)) + "|" +
		in.EmailAddress + "|" +
		strconv.FormatBool(admin.GetOrDefault(in.EmailEnabled, false)) + "|" +
		in.FlowdockAPITokenRef.Key() + "|" +
		in.FlowName + "|" +
		strconv.Itoa(in.IntervalMin) + "|" +
		in.MobileNumber + "|" +
		in.OpsGenieAPIKeyRef.Key() + "|" +
		in.OpsGenieRegion + "|" +
		in.OrgName + "|" +
		in.ServiceKeyRef.Key() + "|" +
		strconv.FormatBool(admin.GetOrDefault(in.SMSEnabled, false)) + "|" +
		in.TeamID + "|" +
		in.TeamName + "|" +
		in.TypeName + "|" +
		in.Username + "|" +
		in.VictorOpsSecretRef.Key() + "|" +
		"[" + strings.Join(in.Roles, ",") + "]"
}

func (in *Notification) SetAPIToken(token string) {
	in.apiToken = token
}

func (in *Notification) SetDatadogAPIKey(token string) {
	in.datadogAPIKey = token
}

func (in *Notification) SetFlowdockAPIToken(token string) {
	in.flowdockAPIToken = token
}

func (in *Notification) SetOpsGenieAPIKey(token string) {
	in.opsGenieAPIKey = token
}

func (in *Notification) SetServiceKey(token string) {
	in.serviceKey = token
}

func (in *Notification) SetVictorOpsAPIKey(token string) {
	in.victorOpsAPIKey = token
}

func (in *Notification) SetVictorOpsRoutingKey(token string) {
	in.victorOpsRoutingKey = token
}

func (in *Notification) IsEqual(notification mongodbatlas.Notification) bool {
	if in == nil {
		return false
	}
	if in.apiToken != notification.APIToken ||
		in.ChannelName != notification.ChannelName ||
		in.datadogAPIKey != notification.DatadogAPIKey ||
		in.DatadogRegion != notification.DatadogRegion ||
		!compare.PtrValuesEqual(in.DelayMin, notification.DelayMin) ||
		in.EmailAddress != notification.EmailAddress ||
		!compare.PtrValuesEqual(in.EmailEnabled, notification.EmailEnabled) ||
		in.flowdockAPIToken != notification.FlowdockAPIToken ||
		in.FlowName != notification.FlowName ||
		in.IntervalMin != notification.IntervalMin ||
		in.MobileNumber != notification.MobileNumber ||
		in.opsGenieAPIKey != notification.OpsGenieAPIKey ||
		in.OpsGenieRegion != notification.OpsGenieRegion ||
		in.OrgName != notification.OrgName ||
		in.serviceKey != notification.ServiceKey ||
		!compare.PtrValuesEqual(in.SMSEnabled, notification.SMSEnabled) ||
		in.TeamID != notification.TeamID ||
		in.TeamName != notification.TeamName ||
		in.TypeName != notification.TypeName ||
		in.Username != notification.Username ||
		in.victorOpsAPIKey != notification.VictorOpsAPIKey ||
		in.victorOpsRoutingKey != notification.VictorOpsRoutingKey {
		return false
	}

	if !compare.IsEqualWithoutOrder(in.Roles, notification.Roles) {
		return false
	}

	return true
}

func (in *Notification) ToAtlas() (*mongodbatlas.Notification, error) {
	result := &mongodbatlas.Notification{}
	err := compat.JSONCopy(result, in)
	if err != nil {
		return nil, err
	}
	result.APIToken = in.apiToken
	result.DatadogAPIKey = in.datadogAPIKey
	result.FlowdockAPIToken = in.flowdockAPIToken
	result.OpsGenieAPIKey = in.opsGenieAPIKey
	result.ServiceKey = in.serviceKey
	result.VictorOpsAPIKey = in.victorOpsAPIKey
	result.VictorOpsRoutingKey = in.victorOpsRoutingKey
	return result, nil
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

func (in MetricThreshold) Key() string {
	return in.MetricName + "|" +
		in.Operator + "|" +
		in.Threshold + "|" +
		in.Units + "|" +
		in.Mode
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
