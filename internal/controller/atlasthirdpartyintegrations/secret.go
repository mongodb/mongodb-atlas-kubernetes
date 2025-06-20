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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

func (h *AtlasThirdPartyIntegrationHandler) secretChanged(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (bool, error) {
	secretName, err := secretName(integration) // at this point the secret should have filed at populateIntegration
	if err != nil {
		return false, fmt.Errorf("failed to check for secret changes: %w", err)
	}
	secret, err := fetchSecret(ctx, h.Client, secretName, integration.Namespace)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve secret %s to evaluate changes: %w", secretName, err)
	}
	currentValue := hashSecret(secret.Data)
	if v, ok := secret.GetAnnotations()[AnnotationContentHash]; ok {
		if v == currentValue {
			return false, nil
		}
	}
	if err := patchSecretAnnotation(ctx, h.Client, secret, AnnotationContentHash, currentValue); err != nil {
		return false, fmt.Errorf("failed to record current secret %s value hash: %w", secretName, err)
	}
	return true, nil
}

func (h *AtlasThirdPartyIntegrationHandler) ensureSecretHash(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) error {
	if _, err := h.secretChanged(ctx, integration); err != nil {
		return fmt.Errorf("failed to ensure secret hash annotation: %w", err)
	}
	return nil
}

func hashSecret(secretData map[string][]byte) string {
	keys := make([]string, 0, len(secretData))
	for k := range secretData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, key := range keys {
		h.Write(([]byte)(key))
		h.Write(secretData[key])
	}
	return hex.EncodeToString(h.Sum(nil))
}

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

func patchSecretAnnotation(ctx context.Context, kubeClient client.Client, secret *v1.Secret, annotation, value string) error {
	updatedSecret := secret.DeepCopy()
	if updatedSecret.Annotations == nil {
		updatedSecret.Annotations = map[string]string{}
	}
	updatedSecret.Annotations[annotation] = value

	secretJSON, err := json.Marshal(updatedSecret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}
	patchErr := kubeClient.Patch(ctx, secret, client.RawPatch(types.MergePatchType, secretJSON))
	if patchErr != nil {
		return fmt.Errorf("failed to patch secret: %w", patchErr)
	}
	return nil
}

func fetchSecretData(ctx context.Context, kubeClient client.Client, name, namespace string) (map[string][]byte, error) {
	secret, err := fetchSecret(ctx, kubeClient, name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret value: %w", err)
	}
	return secret.Data, nil
}
