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
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
)

func (r *AtlasProjectReconciler) ensureIntegration(workflowCtx *workflow.Context, akoProject *akov2.AtlasProject) workflow.Result {
	result := r.createOrDeleteIntegrations(workflowCtx, akoProject.ID(), akoProject)
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

func (r *AtlasProjectReconciler) createOrDeleteIntegrations(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) workflow.Result {
	integrationsInAtlas, err := fetchIntegrations(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err)
	}
	integrationsInAtlasAlias := toAliasThirdPartyIntegration(integrationsInAtlas.Results)

	candidatesForDelete := set.DeprecatedDifference(integrationsInAtlasAlias, project.Spec.Integrations)
	lastApplied, err := mapLastAppliedProjectIntegrations(project)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err)
	}
	identifiersForDelete := filterOwnedIntegrations(candidatesForDelete, lastApplied)
	ctx.Log.Debugf("identifiersForDelete: %v", identifiersForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, identifiersForDelete); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err)
	}

	integrationsToUpdate := set.DeprecatedIntersection(integrationsInAtlasAlias, project.Spec.Integrations)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := r.updateIntegrationsAtlas(ctx, projectID, integrationsToUpdate, project.Namespace); !result.IsOk() {
		return result
	}

	identifiersForCreate := set.DeprecatedDifference(project.Spec.Integrations, integrationsInAtlasAlias)
	ctx.Log.Debugf("identifiersForCreate: %v", identifiersForCreate)
	if result := r.createIntegrationsInAtlas(ctx, projectID, identifiersForCreate, project.Namespace); !result.IsOk() {
		return result
	}

	syncPrometheusStatus(ctx, project, integrationsToUpdate)
	if ready := r.checkIntegrationsReady(ctx, integrationsToUpdate, project.Spec.Integrations); !ready {
		return workflow.InProgress(workflow.ProjectIntegrationReady, "in progress")
	}

	return workflow.OK()
}

func fetchIntegrations(ctx *workflow.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, error) {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(ctx.Context, projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got Integrations From Atlas: %v", *integrationsInAtlas)
	return integrationsInAtlas, nil
}

func (r *AtlasProjectReconciler) updateIntegrationsAtlas(ctx *workflow.Context, projectID string, integrationsToUpdate [][]set.DeprecatedIdentifiable, namespace string) workflow.Result {
	for _, item := range integrationsToUpdate {
		kubeIntegration, err := item[1].(project.Integration).ToAtlas(ctx.Context, r.Client, namespace)
		if kubeIntegration == nil {
			ctx.Log.Warnw("Update Integrations", "Can not convert kube integration", err)
			return workflow.Terminate(workflow.ProjectIntegrationInternal, errors.New("update Integrations: Can not convert kube integration"))
		}
		// As integration secrets are redacted from Atlas, we cannot properly compare them,
		// so as a simple fix we assume changes are always needed at evaluation time
		ctx.Log.Debugf("Try to update integration: %s", kubeIntegration.Type)
		if _, _, err := ctx.Client.Integrations.Replace(ctx.Context, projectID, kubeIntegration.Type, kubeIntegration); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Errorf("cannot apply integration: %w", err))
		}
	}
	return workflow.OK()
}

func deleteIntegrationsFromAtlas(ctx *workflow.Context, projectID string, integrationsToRemove []set.DeprecatedIdentifiable) error {
	for _, integration := range integrationsToRemove {
		if _, err := ctx.Client.Integrations.Delete(ctx.Context, projectID, integration.Identifier().(string)); err != nil {
			return err
		}
		ctx.Log.Debugf("Third Party Integration deleted: %s", integration.Identifier())
	}
	return nil
}

func (r *AtlasProjectReconciler) createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrations []set.DeprecatedIdentifiable, namespace string) workflow.Result {
	for _, item := range integrations {
		integration, err := item.(project.Integration).ToAtlas(ctx.Context, r.Client, namespace)
		if err != nil || integration == nil {
			return workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Errorf("cannot convert integration: %w", err))
		}

		_, resp, err := ctx.Client.Integrations.Create(ctx.Context, projectID, integration.Type, integration)
		if resp != nil && resp.StatusCode != http.StatusOK {
			ctx.Log.Debugw("Create request failed", "Status", resp.Status, "Integration", integration)
		}
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, err)
		}
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) checkIntegrationsReady(ctx *workflow.Context, integrationsIntersection [][]set.DeprecatedIdentifiable, requestedIntegrations []project.Integration) bool {
	if len(integrationsIntersection) != len(requestedIntegrations) {
		return false
	}

	for _, integrationPair := range integrationsIntersection {
		atlas := integrationPair[0].(aliasThirdPartyIntegration)
		spec := integrationPair[1].(project.Integration)

		var areEqual bool
		if isPrometheusType(atlas.Type) {
			areEqual = arePrometheusesEqual(atlas, spec)
		} else {
			// As integration secrets are redacted from Atlas, we cannot properly compare them,
			// so as a simple fix we assume changes were applied correctly as we would
			// have otherwise hit an error at apply time
			areEqual = true
		}
		ctx.Log.Debugw("checkIntegrationsReady", "atlas", atlas, "spec", spec, "areEqual", areEqual)

		if !areEqual {
			return false
		}
	}

	return true
}

type aliasThirdPartyIntegration mongodbatlas.ThirdPartyIntegration

func (i aliasThirdPartyIntegration) Identifier() interface{} {
	return i.Type
}

func toAliasThirdPartyIntegration(list []*mongodbatlas.ThirdPartyIntegration) []aliasThirdPartyIntegration {
	result := make([]aliasThirdPartyIntegration, len(list))
	for i, item := range list {
		if item == nil {
			continue
		}
		result[i] = aliasThirdPartyIntegration(*item)
	}
	return result
}

func syncPrometheusStatus(ctx *workflow.Context, project *akov2.AtlasProject, integrationPairs [][]set.DeprecatedIdentifiable) {
	prometheusIntegration, found := searchAtlasIntegration(integrationPairs, isPrometheusType)
	if !found {
		ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(nil))
		return
	}

	ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(&status.Prometheus{
		Scheme:       prometheusIntegration.Scheme,
		DiscoveryURL: buildPrometheusDiscoveryURL(ctx.Client.BaseURL, project.ID()),
	}))
}

func searchAtlasIntegration(integrationPairs [][]set.DeprecatedIdentifiable, filterFunc func(typeName string) bool) (integration mongodbatlas.ThirdPartyIntegration, found bool) {
	for _, pair := range integrationPairs {
		integrationAlias := pair[0].(aliasThirdPartyIntegration)
		if filterFunc(integrationAlias.Type) {
			return mongodbatlas.ThirdPartyIntegration(integrationAlias), true
		}
	}

	return integration, false
}

func arePrometheusesEqual(atlas aliasThirdPartyIntegration, spec project.Integration) bool {
	return atlas.Type == spec.Type &&
		atlas.UserName == spec.UserName &&
		atlas.ServiceDiscovery == spec.ServiceDiscovery &&
		atlas.Enabled == spec.Enabled
}

func isPrometheusType(typeName string) bool {
	return typeName == "PROMETHEUS"
}

func buildPrometheusDiscoveryURL(baseURL *url.URL, projectID string) string {
	api := fmt.Sprintf("https://%s/prometheus/v1.0", baseURL.Host)
	return fmt.Sprintf("%s/groups/%s/discovery", api, projectID)
}

func mapLastAppliedProjectIntegrations(atlasProject *akov2.AtlasProject) ([]project.Integration, error) {
	lastApplied, err := lastAppliedSpecFrom(atlasProject)
	if err != nil {
		return nil, err
	}

	if lastApplied == nil || len(lastApplied.Integrations) == 0 {
		return nil, nil
	}

	return lastApplied.Integrations, nil
}

func filterOwnedIntegrations(integrationIDs []set.DeprecatedIdentifiable, lastApplied []project.Integration) []set.DeprecatedIdentifiable {
	if len(integrationIDs) == 0 {
		return nil
	}
	filteredIDs := make([]set.DeprecatedIdentifiable, 0, len(integrationIDs))
	for _, integration := range integrationIDs {
		id := integration.Identifier().(string)
		if isIntegrationOwned(lastApplied, id) {
			filteredIDs = append(filteredIDs, integration)
		}
	}
	return filteredIDs
}

func isIntegrationOwned(lastApplied []project.Integration, integrationType string) bool {
	for _, ownedIntegration := range lastApplied {
		if ownedIntegration.Type == integrationType {
			return true
		}
	}
	return false
}
