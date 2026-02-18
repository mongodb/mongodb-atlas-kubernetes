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

package thirdpartyintegration

import (
	"fmt"
	"testing"

	gofuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const fuzzIterations = 100

var integrationTypes = []string{
	"DATADOG",
	"MICROSOFT_TEAMS",
	"NEW_RELIC",
	"OPS_GENIE",
	"PAGER_DUTY",
	"PROMETHEUS",
	"SLACK",
	"VICTOR_OPS",
	"WEBHOOK",
}

var enabledValues = []*string{
	pointer.MakePtr("enabled"),
	pointer.MakePtr("disabled"),
}

func FuzzConvertIntegrations(f *testing.F) {
	for i := range uint(fuzzIterations) {
		f.Add(fmt.Appendf(nil, "seed sample %x", i), i)
	}
	f.Fuzz(func(t *testing.T, data []byte, index uint) {
		integration := ThirdPartyIntegration{}
		fuzzIntegration(gofuzz.NewFromGoFuzz(data), index, &integration)
		atlasIntegration, err := toAtlas(&integration)
		require.NoError(t, err)
		result, err := fromAtlas(atlasIntegration)
		require.NoError(t, err)
		assert.Equal(t, &integration, result, "failed for index=%d", index)
	})
}

func fuzzIntegration(fuzzer *gofuzz.Fuzzer, index uint, integration *ThirdPartyIntegration) {
	fuzzer.NilChance(0).Fuzz(integration)
	integration.ID = "" // ID is provided by Atlas, cannot complete a roundtrip
	integration.ProjectDualReference.ExternalProjectRef = nil
	integration.ProjectDualReference.ProjectRef = nil
	integration.ProjectDualReference.ConnectionSecret = nil

	integration.Type = integrationTypes[index%uint(len(integrationTypes))] // type by index

	if integration.Type == "DATADOG" {
		index2 := index + 1
		integration.Datadog.SendCollectionLatencyMetrics = enabledValues[index%uint(len(enabledValues))]
		integration.Datadog.SendDatabaseMetrics = enabledValues[index2%uint(len(enabledValues))]
		integration.Datadog.APIKeySecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.Datadog = nil
		integration.DatadogSecrets = nil
	}

	if integration.Type == "MICROSOFT_TEAMS" {
		integration.MicrosoftTeams.URLSecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.MicrosoftTeams = nil
		integration.MicrosoftTeamsSecrets = nil
	}

	if integration.Type == "NEW_RELIC" {
		integration.NewRelic.CredentialsSecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.NewRelic = nil
		integration.NewRelicSecrets = nil
	}

	if integration.Type == "OPS_GENIE" {
		integration.OpsGenie.APIKeySecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.OpsGenie = nil
		integration.OpsGenieSecrets = nil
	}

	if integration.Type == "PAGER_DUTY" {
		integration.PagerDuty.ServiceKeySecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.PagerDuty = nil
		integration.PagerDutySecrets = nil
	}

	if integration.Type == "PROMETHEUS" {
		integration.Prometheus.Enabled = enabledValues[index%uint(len(enabledValues))]
		integration.Prometheus.PrometheusCredentialsSecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.Prometheus = nil
		integration.PrometheusSecrets = nil
	}

	if integration.Type == "SLACK" {
		integration.Slack.APITokenSecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.Slack = nil
		integration.SlackSecrets = nil
	}

	if integration.Type == "VICTOR_OPS" {
		integration.VictorOps.APIKeySecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.VictorOps = nil
		integration.VictorOpsSecrets = nil
	}

	if integration.Type == "WEBHOOK" {
		integration.Webhook.URLSecretRef.Name = "" // not part of the atlas conversion roundtrip
	} else {
		integration.Webhook = nil
		integration.WebhookSecrets = nil
	}
}
