package atlasproject

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/integrations"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
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
	integrationsInAtlas, err := r.integrationsService.List(ctx.Context, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err.Error())
	}
	ctx.Log.Debugf("Got Integrations from Atlas: %v", integrationsInAtlas)

	integrationsInKube, err := integrations.NewIntegrations(project.Spec.Integrations)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err.Error())
	}

	identifiersForDelete := set.Difference(integrationsInAtlas, integrationsInKube)
	ctx.Log.Debugf("identifiersForDelete: %v", identifiersForDelete)
	if result := r.deleteIntegrationsFromAtlas(ctx, projectID, identifiersForDelete); !result.IsOk() {
		return result
	}

	integrationsToUpdate := set.Intersection(integrationsInAtlas, integrationsInKube)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := r.updateIntegrationsAtlas(ctx, projectID, integrationsToUpdate, project.Namespace); !result.IsOk() {
		return result
	}

	identifiersForCreate := set.Difference(integrationsInKube, integrationsInAtlas)
	ctx.Log.Debugf("identifiersForCreate: %v", identifiersForCreate)
	if result := r.createIntegrationsInAtlas(ctx, projectID, identifiersForCreate, project.Namespace); !result.IsOk() {
		return result
	}

	syncPrometheusStatus(ctx, project, integrationsToUpdate)
	if ready := r.checkIntegrationsReady(ctx, integrationsToUpdate, integrationsInKube); !ready {
		return workflow.InProgress(workflow.ProjectIntegrationReady, "in progress")
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) updateIntegrationsAtlas(ctx *workflow.Context, projectID string, integrationsToUpdate [][]set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrationsToUpdate {
		kubeIntegration, err := item[1].(project.Integration)
		if err != nil {
			ctx.Log.Warnw("Update Integrations", "Can not convert kube integration", err)
			return workflow.Terminate(workflow.ProjectIntegrationInternal, "Update Integrations: Can not convert kube integration")
		}
		// As integration secrets are redacted from Atlas, we cannot properly compare them,
		// so as a simple fix we assume changes are always needed at evaluation time
		ctx.Log.Debugf("Try to update integration: %s", kubeIntegration.Type)
		if err := r.integrationsService.Update(ctx.Context, projectID, kubeIntegration.Type, kubeIntegration, r.Client, namespace); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot apply integration %s: %v", kubeIntegration.Type, err.Error()))
		}
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) deleteIntegrationsFromAtlas(ctx *workflow.Context, projectID string, integrationsToRemove []set.Identifiable) workflow.Result {
	for _, integration := range integrationsToRemove {
		if err := r.integrationsService.Delete(ctx.Context, projectID, integration.Identifier().(string)); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot delete integration %s: %v", integration.Identifier().(string), err.Error()))
		}
		ctx.Log.Debugf("Third Party Integration deleted: %s", integration.Identifier())
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrationsToCreate []set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrationsToCreate {
		integration := item.(integrations.Integration)

		if err := r.integrationsService.Create(ctx.Context, projectID, integration.Type, integration, r.Client, namespace); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Cannot create integration %s: %v", integration.Type, err.Error()))
		}
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) checkIntegrationsReady(ctx *workflow.Context, integrationsIntersection [][]set.Identifiable, requestedIntegrations []integrations.Integration) bool {
	if len(integrationsIntersection) != len(requestedIntegrations) {
		return false
	}

	for _, integrationPair := range integrationsIntersection {
		atlas := integrationPair[0].(integrations.Integration)
		spec := integrationPair[1].(integrations.Integration)

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

func syncPrometheusStatus(ctx *workflow.Context, project *akov2.AtlasProject, integrationPairs [][]set.Identifiable) {
	prometheusIntegration, found := searchAtlasIntegration(integrationPairs, isPrometheusType)
	if !found {
		ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(nil))
		return
	}
	ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(&status.Prometheus{
		Scheme:       prometheusIntegration.Scheme,
		DiscoveryURL: buildPrometheusDiscoveryURL(ctx.SdkClient.GetConfig().Servers[0].URL, project.ID()),
	}))
}

func searchAtlasIntegration(integrationPairs [][]set.Identifiable, filterFunc func(typeName string) bool) (integration integrations.Integration, found bool) {
	for _, pair := range integrationPairs {
		integrationAlias := pair[0].(integrations.Integration)
		if filterFunc(integrationAlias.Type) {
			return integrationAlias, true
		}
	}

	return integration, false
}

func arePrometheusesEqual(atlas, spec integrations.Integration) bool {
	return atlas.Type == spec.Type &&
		atlas.UserName == spec.UserName &&
		atlas.ServiceDiscovery == spec.ServiceDiscovery &&
		atlas.Enabled == spec.Enabled
}

func isPrometheusType(typeName string) bool {
	return typeName == "PROMETHEUS"
}

func buildPrometheusDiscoveryURL(baseURL string, projectID string) string {
	api := fmt.Sprintf("%s/prometheus/v1.0", baseURL)
	return fmt.Sprintf("%s/groups/%s/discovery", api, projectID)
}
