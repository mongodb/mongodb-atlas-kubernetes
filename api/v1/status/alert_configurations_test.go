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

package status_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func alertConfigWithNotification() *admin.GroupAlertsConfig {
	return &admin.GroupAlertsConfig{
		Id:            admin.PtrString("alert-config-id"),
		GroupId:       admin.PtrString("group-id"),
		EventTypeName: admin.PtrString("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
		Enabled:       admin.PtrBool(true),
		Notifications: &[]admin.AlertsNotificationRootForGroup{
			{
				TypeName:            admin.PtrString("SLACK"),
				ChannelName:         admin.PtrString("#alerts"),
				ApiToken:            admin.PtrString("xoxb-plaintext-slack-token"),
				ServiceKey:          admin.PtrString("plaintext-pagerduty-service-key"),
				DatadogApiKey:       admin.PtrString("plaintext-datadog-key"),
				OpsGenieApiKey:      admin.PtrString("plaintext-opsgenie-key"),
				VictorOpsApiKey:     admin.PtrString("plaintext-victorops-key"),
				VictorOpsRoutingKey: admin.PtrString("plaintext-victorops-routing-key"),
			},
		},
	}
}

func TestNewIncorrectAlertConfigStatus_RedactsCredentials(t *testing.T) {
	result := status.NewIncorrectAlertConfigStatus(
		"failed to create atlas alert configuration: boom",
		alertConfigWithNotification(),
		zaptest.NewLogger(t).Sugar(),
	)

	require.Len(t, result.Notifications, 1)
	notification := result.Notifications[0]

	assert.Empty(t, notification.APIToken)
	assert.Empty(t, notification.ServiceKey)
	assert.Empty(t, notification.DatadogAPIKey)
	assert.Empty(t, notification.OpsGenieAPIKey)
	assert.Empty(t, notification.VictorOpsAPIKey)
	assert.Empty(t, notification.VictorOpsRoutingKey)

	assert.Equal(t, "SLACK", notification.TypeName)
	assert.Equal(t, "#alerts", notification.ChannelName)
	assert.Equal(t, "failed to create atlas alert configuration: boom", result.ErrorMessage)
}

func TestNewIncorrectAlertConfigStatus_NoPlaintextInSerializedStatus(t *testing.T) {
	result := status.NewIncorrectAlertConfigStatus(
		"failed to create atlas alert configuration: boom",
		alertConfigWithNotification(),
		zaptest.NewLogger(t).Sugar(),
	)

	serialized, err := json.Marshal(result)
	require.NoError(t, err)

	for _, secret := range []string{
		"xoxb-plaintext-slack-token",
		"plaintext-pagerduty-service-key",
		"plaintext-datadog-key",
		"plaintext-opsgenie-key",
		"plaintext-victorops-key",
		"plaintext-victorops-routing-key",
	} {
		assert.NotContains(t, string(serialized), secret)
	}
}

func TestParseAlertConfiguration_KeepsMaskedCredentials(t *testing.T) {
	response := admin.GroupAlertsConfig{
		Id:            admin.PtrString("alert-config-id"),
		EventTypeName: admin.PtrString("REPLICATION_OPLOG_WINDOW_RUNNING_OUT"),
		Notifications: &[]admin.AlertsNotificationRootForGroup{
			{
				TypeName: admin.PtrString("SLACK"),
				ApiToken: admin.PtrString("xoxb-*********************"),
			},
		},
	}

	result := status.ParseAlertConfiguration(response, zaptest.NewLogger(t).Sugar())

	require.Len(t, result.Notifications, 1)
	assert.Equal(t, "xoxb-*********************", result.Notifications[0].APIToken)
}
