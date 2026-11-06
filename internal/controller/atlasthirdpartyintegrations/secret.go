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

package integrations

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

func fetchIntegrationSecrets(ctx context.Context, kubeClient client.Client, integration *akov2.AtlasThirdPartyIntegration) (map[string][]byte, error) {
	name, err := secretName(integration)
	if err != nil {
		return nil, fmt.Errorf("failed to solve integration secret name: %w", err)
	}
	return fetchSecretData(ctx, kubeClient, name, integration.Namespace)
}

func secretName(integration *akov2.AtlasThirdPartyIntegration) (string, error) {
	switch integration.Spec.Type {
	case "DATADOG":
		return integration.Spec.Datadog.APIKeySecretRef.Name, nil
	case "MICROSOFT_TEAMS":
		return integration.Spec.MicrosoftTeams.URLSecretRef.Name, nil
	case "NEW_RELIC":
		return integration.Spec.NewRelic.CredentialsSecretRef.Name, nil
	case "OPS_GENIE":
		return integration.Spec.OpsGenie.APIKeySecretRef.Name, nil
	case "PAGER_DUTY":
		return integration.Spec.PagerDuty.ServiceKeySecretRef.Name, nil
	case "PROMETHEUS":
		return integration.Spec.Prometheus.PrometheusCredentialsSecretRef.Name, nil
	case "SLACK":
		return integration.Spec.Slack.APITokenSecretRef.Name, nil
	case "VICTOR_OPS":
		return integration.Spec.VictorOps.APIKeySecretRef.Name, nil
	case "WEBHOOK":
		return integration.Spec.Webhook.URLSecretRef.Name, nil
	default:
		return "", fmt.Errorf("%w %v", thirdpartyintegration.ErrUnsupportedIntegrationType, integration.Spec.Type)
	}
}

func fetchSecret(ctx context.Context, kubeClient client.Client, name, namespace string) (*v1.Secret, error) {
	secret := v1.Secret{}
	err := kubeClient.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret: %w", err)
	}
	return &secret, nil
}

func fetchSecretData(ctx context.Context, kubeClient client.Client, name, namespace string) (map[string][]byte, error) {
	secret, err := fetchSecret(ctx, kubeClient, name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret value: %w", err)
	}
	return secret.Data, nil
}
