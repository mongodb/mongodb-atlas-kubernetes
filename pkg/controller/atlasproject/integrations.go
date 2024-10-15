package atlasproject

import (
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
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
	integrationsInAtlas, err := fetchIntegrations(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err.Error())
	}
	integrationsInAtlasAlias := toAliasThirdPartyIntegration(integrationsInAtlas.GetResults())

	identifiersForDelete := set.Difference(integrationsInAtlasAlias, project.Spec.Integrations)
	ctx.Log.Debugf("identifiersForDelete: %v", identifiersForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, identifiersForDelete); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err.Error())
	}

	integrationsToUpdate := set.Intersection(integrationsInAtlasAlias, project.Spec.Integrations)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := r.updateIntegrationsAtlas(ctx, projectID, integrationsToUpdate, project.Namespace); !result.IsOk() {
		return result
	}

	identifiersForCreate := set.Difference(project.Spec.Integrations, integrationsInAtlasAlias)
	ctx.Log.Debugf("identifiersForCreate: %v", identifiersForCreate)
	if result := r.createIntegrationsInAtlas(ctx, projectID, identifiersForCreate, project.Namespace); !result.IsOk() {
		return result
	}

	syncPrometheusStatus(ctx, project, integrationsToUpdate)
	if ready := r.checkIntegrationsReady(ctx, project.Namespace, integrationsToUpdate, project.Spec.Integrations); !ready {
		return workflow.InProgress(workflow.ProjectIntegrationReady, "in progress")
	}

	return workflow.OK()
}

func fetchIntegrations(ctx *workflow.Context, projectID string) (*admin.PaginatedIntegration, error) {
	integrationsInAtlas, _, err := ctx.SdkClient.ThirdPartyIntegrationsApi.ListThirdPartyIntegrations(ctx.Context, projectID).Execute()
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got Integrations From Atlas: %v", *integrationsInAtlas)
	return integrationsInAtlas, nil
}

func (r *AtlasProjectReconciler) updateIntegrationsAtlas(ctx *workflow.Context, projectID string, integrationsToUpdate [][]set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrationsToUpdate {
		kubeIntegration, err := item[1].(project.Integration).ToAtlas(ctx.Context, r.Client, namespace)
		if kubeIntegration == nil {
			ctx.Log.Warnw("Update Integrations", "Can not convert kube integration", err)
			return workflow.Terminate(workflow.ProjectIntegrationInternal, "Update Integrations: Can not convert kube integration")
		}
		// As integration secrets are redacted from Atlas, we cannot properly compare them,
		// so as a simple fix we assume changes are always needed at evaluation time
		ctx.Log.Debugf("Try to update integration: %s", kubeIntegration.Type)
		if _, _, err := ctx.Client.Integrations.Replace(ctx.Context, projectID, kubeIntegration.Type, kubeIntegration); err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, fmt.Sprintf("Can not apply integration: %v", err))
		}
	}
	return workflow.OK()
}

func deleteIntegrationsFromAtlas(ctx *workflow.Context, projectID string, integrationsToRemove []set.Identifiable) error {
	for _, integration := range integrationsToRemove {
		if _, err := ctx.Client.Integrations.Delete(ctx.Context, projectID, integration.Identifier().(string)); err != nil {
			return err
		}
		ctx.Log.Debugf("Third Party Integration deleted: %s", integration.Identifier())
	}
	return nil
}

func (r *AtlasProjectReconciler) createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrations []set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrations {
		integration, err := item.(project.Integration).ToAtlas(ctx.Context, r.Client, namespace)
		if err != nil || integration == nil {
			return workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Sprintf("cannot convert integration: %s", err.Error()))
		}

		_, resp, err := ctx.SdkClient.ThirdPartyIntegrationsApi.CreateThirdPartyIntegration(ctx.Context, integration.GetType(), projectID, integration).Execute()
		if resp.StatusCode != http.StatusOK {
			ctx.Log.Debugw("Create request failed", "Status", resp.Status, "Integration", integration)
		}
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationRequest, err.Error())
		}
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) checkIntegrationsReady(ctx *workflow.Context, namespace string, integrationsIntersection [][]set.Identifiable, requestedIntegrations []project.Integration) bool {
	if len(integrationsIntersection) != len(requestedIntegrations) {
		return false
	}

	for _, integrationPair := range integrationsIntersection {
		atlas := integrationPair[0].(aliasThirdPartyIntegration)
		spec := integrationPair[1].(project.Integration)

		var areEqual bool
		if isPrometheusType(atlas.GetType()) {
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

type aliasThirdPartyIntegration struct {
	admin.ThirdPartyIntegration
}

func (i aliasThirdPartyIntegration) Identifier() interface{} {
	return i.GetType()
}

func toAliasThirdPartyIntegration(list []admin.ThirdPartyIntegration) []aliasThirdPartyIntegration {
	result := make([]aliasThirdPartyIntegration, len(list))
	for i, item := range list {
		result[i] = aliasThirdPartyIntegration{item}
	}
	return result
}

func syncPrometheusStatus(ctx *workflow.Context, project *akov2.AtlasProject, integrationPairs [][]set.Identifiable) {
	ctx.EnsureStatusOption(status.AtlasProjectPrometheusOption(&status.Prometheus{
		DiscoveryURL: buildPrometheusDiscoveryURL(ctx.SdkClient.GetConfig().Servers[0].URL, project.ID()),
	}))
}

func searchAtlasIntegration(integrationPairs [][]set.Identifiable, filterFunc func(typeName string) bool) (integration admin.ThirdPartyIntegration, found bool) {
	for _, pair := range integrationPairs {
		integrationAlias := pair[0].(aliasThirdPartyIntegration)
		if filterFunc(integrationAlias.GetType()) {
			return integrationAlias.ThirdPartyIntegration, true
		}
	}

	return integration, false
}

func arePrometheusesEqual(atlas aliasThirdPartyIntegration, spec project.Integration) bool {
	return atlas.GetType() == spec.Type &&
		atlas.GetUsername() == spec.UserName &&
		atlas.GetServiceDiscovery() == spec.ServiceDiscovery &&
		atlas.GetEnabled() == spec.Enabled
}

func isPrometheusType(typeName string) bool {
	return typeName == "PROMETHEUS"
}

func buildPrometheusDiscoveryURL(baseURL string, projectID string) string {
	api := fmt.Sprintf("%s/prometheus/v1.0", baseURL)
	return fmt.Sprintf("%s/groups/%s/discovery", api, projectID)
}
