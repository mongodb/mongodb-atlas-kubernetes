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

package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasThirdPartyIntegrationBySecretsIndex = "atlasthirdpartyintegration.spec.secrets"
)

type AtlasThirdPartyIntegrationBySecretsIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasThirdPartyIntegrationBySecretsIndexer(logger *zap.Logger) *AtlasThirdPartyIntegrationBySecretsIndexer {
	return &AtlasThirdPartyIntegrationBySecretsIndexer{
		logger: logger.Named(AtlasThirdPartyIntegrationBySecretsIndex).Sugar(),
	}
}

func (*AtlasThirdPartyIntegrationBySecretsIndexer) Object() client.Object {
	return &akov2.AtlasThirdPartyIntegration{}
}

func (*AtlasThirdPartyIntegrationBySecretsIndexer) Name() string {
	return AtlasThirdPartyIntegrationBySecretsIndex
}

func (a *AtlasThirdPartyIntegrationBySecretsIndexer) Keys(object client.Object) []string {
	tpi, ok := object.(*akov2.AtlasThirdPartyIntegration)
	if !ok {
		a.logger.Errorf("expected %T but got %T", &akov2.AtlasThirdPartyIntegration{}, object)
		return nil
	}
	name := keyName(tpi)
	if name == "" {
		return nil
	}
	return []string{client.ObjectKey{Name: name, Namespace: tpi.Namespace}.String()}
}

func keyName(tpi *akov2.AtlasThirdPartyIntegration) string {
	switch tpi.Spec.Type {
	case "DATADOG":
		if tpi.Spec.Datadog != nil && tpi.Spec.Datadog.APIKeySecretRef.Name != "" {
			return tpi.Spec.Datadog.APIKeySecretRef.Name
		}
	case "MICROSOFT_TEAMS":
		if tpi.Spec.MicrosoftTeams != nil && tpi.Spec.MicrosoftTeams.URLSecretRef.Name != "" {
			return tpi.Spec.MicrosoftTeams.URLSecretRef.Name
		}
	case "NEW_RELIC":
		if tpi.Spec.NewRelic != nil && tpi.Spec.NewRelic.CredentialsSecretRef.Name != "" {
			return tpi.Spec.NewRelic.CredentialsSecretRef.Name
		}
	case "OPS_GENIE":
		if tpi.Spec.OpsGenie != nil && tpi.Spec.OpsGenie.APIKeySecretRef.Name != "" {
			return tpi.Spec.OpsGenie.APIKeySecretRef.Name
		}
	case "PAGER_DUTY":
		if tpi.Spec.PagerDuty != nil && tpi.Spec.PagerDuty.ServiceKeySecretRef.Name != "" {
			return tpi.Spec.PagerDuty.ServiceKeySecretRef.Name
		}
	case "PROMETHEUS":
		if tpi.Spec.Prometheus != nil && tpi.Spec.Prometheus.PrometheusCredentialsSecretRef.Name != "" {
			return tpi.Spec.Prometheus.PrometheusCredentialsSecretRef.Name
		}
	case "SLACK":
		if tpi.Spec.Slack != nil && tpi.Spec.Slack.APITokenSecretRef.Name != "" {
			return tpi.Spec.Slack.APITokenSecretRef.Name
		}
	case "VICTOR_OPS":
		if tpi.Spec.VictorOps != nil && tpi.Spec.VictorOps.APIKeySecretRef.Name != "" {
			return tpi.Spec.VictorOps.APIKeySecretRef.Name
		}
	case "WEBHOOK":
		if tpi.Spec.Webhook != nil && tpi.Spec.Webhook.URLSecretRef.Name != "" {
			return tpi.Spec.Webhook.URLSecretRef.Name
		}
	}
	return ""
}
