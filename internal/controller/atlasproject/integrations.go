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

package atlasproject

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	integration "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

func (r *AtlasProjectReconciler) fromAKO(ctx context.Context, project *akov2.AtlasProject) ([]*integration.ThirdPartyIntegration, error) {
	result := make([]*integration.ThirdPartyIntegration, 0, len(project.Spec.Integrations))

	for _, i := range project.Spec.Integrations {
		tpi := &integration.ThirdPartyIntegration{
			AtlasThirdPartyIntegrationSpec: akov2.AtlasThirdPartyIntegrationSpec{
				Type: i.Type,
			},
		}
		switch i.Type {
		case "DATADOG":
			apiKey, err := i.APIKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read API key for Datadog integration: %w", err)
			}

			tpi.DatadogSecrets = &integration.DatadogSecrets{
				APIKey: apiKey,
			}
			tpi.Datadog = &akov2.DatadogIntegration{
				Region:                       i.Region,
				SendCollectionLatencyMetrics: pointer.MakePtr("disabled"),
				SendDatabaseMetrics:          pointer.MakePtr("disabled"),
			}
		case "MICROSOFT_TEAMS":
			tpi.MicrosoftTeamsSecrets = &integration.MicrosoftTeamsSecrets{
				WebhookUrl: i.MicrosoftTeamsWebhookURL,
			}
			tpi.MicrosoftTeams = &akov2.MicrosoftTeamsIntegration{}
		case "NEW_RELIC":
			licenseKey, err := i.LicenseKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read license key for NewRelic integration: %w", err)
			}

			readToken, err := i.ReadTokenRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read read token for NewRelic integration: %w", err)
			}

			writeToken, err := i.WriteTokenRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read write token for NewRelic integration: %w", err)
			}

			tpi.NewRelicSecrets = &integration.NewRelicSecrets{
				AccountID:  i.AccountID,
				LicenseKey: licenseKey,
				ReadToken:  readToken,
				WriteToken: writeToken,
			}
			tpi.NewRelic = &akov2.NewRelicIntegration{}
		case "OPS_GENIE":
			apiKey, err := i.APIKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read API key for OpsGenie integration: %w", err)
			}

			tpi.OpsGenieSecrets = &integration.OpsGenieSecrets{
				APIKey: apiKey,
			}
			tpi.OpsGenie = &akov2.OpsGenieIntegration{
				Region: i.Region,
			}
		case "PAGER_DUTY":
			serviceKey, err := i.ServiceKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read service key for PagerDuty integration: %w", err)
			}
			tpi.PagerDutySecrets = &integration.PagerDutySecrets{
				ServiceKey: serviceKey,
			}
			tpi.PagerDuty = &akov2.PagerDutyIntegration{
				Region: i.Region,
			}
		case "PROMETHEUS":
			password, err := i.PasswordRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read password for Prometheus integration: %w", err)
			}

			enabled := "enabled"
			if !i.Enabled {
				enabled = "disabled"
			}

			tpi.Prometheus = &akov2.PrometheusIntegration{
				Enabled:          pointer.MakePtr(enabled),
				ServiceDiscovery: i.ServiceDiscovery,
			}
			tpi.PrometheusSecrets = &integration.PrometheusSecrets{
				Username: i.UserName,
				Password: password,
			}
		case "SLACK":
			apiToken, err := i.APITokenRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read API token for Slack integration: %w", err)
			}
			tpi.Slack = &akov2.SlackIntegration{
				ChannelName: i.ChannelName,
				TeamName:    i.TeamName,
			}
			tpi.SlackSecrets = &integration.SlackSecrets{
				APIToken: apiToken,
			}
		case "VICTOR_OPS":
			apiKey, err := i.APIKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read API key for VictorOps integration: %w", err)
			}

			routingKey, err := i.RoutingKeyRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read routing key for VictorOps integration: %w", err)
			}

			tpi.VictorOps = &akov2.VictorOpsIntegration{
				RoutingKey: routingKey,
			}
			tpi.VictorOpsSecrets = &integration.VictorOpsSecrets{
				APIKey: apiKey,
			}
		case "WEBHOOK":
			secret, err := i.SecretRef.ReadPassword(ctx, r.Client, project.Namespace)
			if err != nil {
				return nil, fmt.Errorf("failed to read secret for Webhook integration: %w", err)
			}

			tpi.Webhook = &akov2.WebhookIntegration{}
			tpi.WebhookSecrets = &integration.WebhookSecrets{
				URL:    i.URL,
				Secret: secret,
			}
		}

		result = append(result, tpi)
	}

	return result, nil
}

func (r *AtlasProjectReconciler) ensureIntegration(workflowCtx *workflow.Context, akoProject *akov2.AtlasProject) workflow.DeprecatedResult {
	integrationsInAKO, err := r.fromAKO(workflowCtx.Context, akoProject)
	if err != nil {
		result := workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Errorf("failed to convert integrations from AKO: %w", err))
		workflowCtx.SetConditionFromResult(api.IntegrationReadyType, result)

		return result
	}

	lastAppliedIntegrations, err := mapLastAppliedProjectIntegrations(akoProject)
	if err != nil {
		result := workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Errorf("failed to map last applied integrations: %w", err))
		workflowCtx.SetConditionFromResult(api.IntegrationReadyType, result)

		return result
	}

	reconciler := NewIntegrationReconciler(workflowCtx, akoProject, integrationsInAKO, lastAppliedIntegrations)
	result := reconciler.reconcile(workflowCtx)

	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.IntegrationReadyType, result)
		return result
	}

	if len(akoProject.Spec.Integrations) == 0 {
		workflowCtx.UnsetCondition(api.IntegrationReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(api.IntegrationReadyType)
	return workflow.OK()
}

type IntegrationReconciler struct {
	project                     *akov2.AtlasProject
	lasAppliedIntegrationsTypes map[string]struct{}
	integrationsInAKO           map[string]*integration.ThirdPartyIntegration
	service                     integration.ThirdPartyIntegrationService
}

func (ir IntegrationReconciler) reconcile(ctx *workflow.Context) workflow.DeprecatedResult {
	list, err := ir.service.List(ctx.Context, ir.project.ID())
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err)
	}

	integrationInAtlas := mapIntegrationsPerType(list)

	for _, inAtlas := range integrationInAtlas {
		if _, found := ir.integrationsInAKO[inAtlas.Type]; found {
			continue
		}

		if _, found := ir.lasAppliedIntegrationsTypes[inAtlas.Type]; !found {
			ctx.Log.Debugf("integration %s is not owned by AKO, skipping", inAtlas.Type)
			continue
		}

		ctx.Log.Debugf("deleting integration %s", inAtlas.Type)
		if err = ir.service.Delete(ctx.Context, ir.project.ID(), inAtlas.Type); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Errorf("failed to remove integration %s: %w", inAtlas.Type, err))
		}
	}

	for _, inAKO := range ir.integrationsInAKO {
		if _, found := integrationInAtlas[inAKO.Type]; found {
			ctx.Log.Debugf("updating integration %s", inAKO.Type)
			if _, err = ir.service.Update(ctx.Context, ir.project.ID(), inAKO); err != nil {
				return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Errorf("failed to update integration %s: %w", inAKO.Type, err))
			}

			continue
		}

		ctx.Log.Debugf("creating integration %s", inAKO.Type)
		if _, err = ir.service.Create(ctx.Context, ir.project.ID(), inAKO); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Errorf("failed to create integration %s: %w", inAKO.Type, err))
		}
	}

	if _, found := ir.integrationsInAKO["PROMETHEUS"]; found {
		ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(&status.Prometheus{
			Scheme:       "https",
			DiscoveryURL: fmt.Sprintf("https://%s/prometheus/v1.0/groups/%s/discovery", ctx.SdkClientSet.SdkClient20250312009.GetConfig().Host, ir.project.ID()),
		}))
	} else {
		ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(nil))
	}

	return workflow.OK()
}

func NewIntegrationReconciler(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	integrations []*integration.ThirdPartyIntegration,
	lastAppliedIntegrationsTypes map[string]struct{},
) *IntegrationReconciler {
	return &IntegrationReconciler{
		project:                     project,
		integrationsInAKO:           mapIntegrationsPerType(integrations),
		lasAppliedIntegrationsTypes: lastAppliedIntegrationsTypes,
		service:                     integration.NewThirdPartyIntegrationService(ctx.SdkClientSet.SdkClient20250312009.ThirdPartyIntegrationsApi),
	}
}

func mapIntegrationsPerType(integrations []*integration.ThirdPartyIntegration) map[string]*integration.ThirdPartyIntegration {
	integrationsPerType := make(map[string]*integration.ThirdPartyIntegration)
	for _, i := range integrations {
		integrationsPerType[i.Type] = i
	}

	return integrationsPerType
}

func mapLastAppliedProjectIntegrations(atlasProject *akov2.AtlasProject) (map[string]struct{}, error) {
	lastApplied, err := lastAppliedSpecFrom(atlasProject)
	if err != nil {
		return nil, err
	}

	if lastApplied == nil || len(lastApplied.Integrations) == 0 {
		return nil, nil
	}

	result := make(map[string]struct{}, len(lastApplied.Integrations))
	for _, i := range lastApplied.Integrations {
		result[i.Type] = struct{}{}
	}

	return result, nil
}
