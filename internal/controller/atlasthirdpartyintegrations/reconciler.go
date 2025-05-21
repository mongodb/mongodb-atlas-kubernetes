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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/dualref"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/result"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/autogen/secret"
	internalbuilder "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/builder"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type Reconciler struct {
	client client.Client
}

func NewAtlasThirdPartyIntegrationsReconciler() *Reconciler {
	return &Reconciler{}
}

func (r *Reconciler) NewBuilderWithManager(mgr ctrl.Manager) *builder.Builder {
	r.client = mgr.GetClient()
	obj := &akov2next.AtlasThirdPartyIntegration{}
	return internalbuilder.NewDefaultControllerBuilder(mgr, obj)
}

func (r *Reconciler) HandleInitial(ctx context.Context, integration *akov2next.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	atlasClientSet := atlas.FromContext(ctx)
	service := thirdpartyintegration.NewThirdPartyIntegrationServiceFromClientSet(atlasClientSet)
	project, err := dualref.ResolveProject(ctx, r.client, integration)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to fetch referenced project: %w", err))
	}

	requestedIntegration, err := populateIntegration(ctx, r.client, integration)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to populate inetgration: %w", err))
	}
	createdIntegration, err := service.Create(ctx, project.ID, requestedIntegration)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create Atlas Third Party Integration for %s: %w",
			integration.Spec.Type, err))
	}
	integration.Status.ID = createdIntegration.ID
	return result.NextState(state.StateCreating,
		fmt.Sprintf("Creating Atlas Third Party Integration for %s", integration.Spec.Type))
}

func populateIntegration(ctx context.Context, kubeClient client.Client, integration *akov2next.AtlasThirdPartyIntegration) (*thirdpartyintegration.ThirdPartyIntegration, error) {
	secrets, err := fetchIntegrationSecrets(ctx, kubeClient, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch integration secrets: %w", err)
	}
	internalIntegration, err := thirdpartyintegration.NewFromSpec(integration, secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to populate integration: %w", err)
	}
	return internalIntegration, err
}

func fetchIntegrationSecrets(ctx context.Context, kubeClient client.Client, integration *akov2next.AtlasThirdPartyIntegration) (map[string][]byte, error) {
	switch integration.Spec.Type {
	case "DATADOG":
		return secret.Fetch(ctx, kubeClient, integration.Spec.Datadog.APIKeySecret.Name)
	case "MICROSOFT_TEAMS":
		return secret.Fetch(ctx, kubeClient, integration.Spec.MicrosoftTeams.URLSecret.Name)
	case "NEW_RELIC":
		return secret.Fetch(ctx, kubeClient, integration.Spec.NewRelic.CredentialsSecret.Name)
	case "OPS_GENIE":
		return secret.Fetch(ctx, kubeClient, integration.Spec.OpsGenie.APIKeySecret.Name)
	case "PAGER_DUTY":
		return secret.Fetch(ctx, kubeClient, integration.Spec.PagerDuty.ServiceKeySecret.Name)
	case "PROMETHEUS":
		return secret.Fetch(ctx, kubeClient, integration.Spec.Prometheus.PrometheusCredentialsSecret.Name)
	case "SLACK":
		return secret.Fetch(ctx, kubeClient, integration.Spec.Slack.APITokenSecret.Name)
	case "VICTOR_OPS":
		return secret.Fetch(ctx, kubeClient, integration.Spec.VictorOps.APIKeySecret.Name)
	case "WEBHOOK":
		return secret.Fetch(ctx, kubeClient, integration.Spec.Webhook.URLSecret.Name)
	default:
		return nil, fmt.Errorf("%w %v", thirdpartyintegration.ErrUnsupportedIntegrationType, integration.Spec.Type)
	}
}
