package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func (r *AtlasProjectReconciler) ensureIntegration(workflowCtx *workflow.Context, akoProject *mdbv1.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canIntegrationsReconcile(workflowCtx, protected, akoProject)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.IntegrationReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.IntegrationReadyType, result)

		return result
	}

	result := r.createOrDeleteIntegrations(workflowCtx, akoProject.ID(), akoProject)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.IntegrationReadyType, result)
		return result
	}

	if len(akoProject.Spec.Integrations) == 0 {
		workflowCtx.UnsetCondition(status.IntegrationReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(status.IntegrationReadyType)
	return workflow.OK()
}

func (r *AtlasProjectReconciler) createOrDeleteIntegrations(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	integrationsInAtlas, err := fetchIntegrations(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInternal, err.Error())
	}
	integrationsInAtlasAlias := toAliasThirdPartyIntegration(integrationsInAtlas.Results)

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

func fetchIntegrations(ctx *workflow.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, error) {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got Integrations From Atlas: %v", *integrationsInAtlas)
	return integrationsInAtlas, nil
}

func (r *AtlasProjectReconciler) updateIntegrationsAtlas(ctx *workflow.Context, projectID string, integrationsToUpdate [][]set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrationsToUpdate {
		atlasIntegration := item[0].(aliasThirdPartyIntegration)
		kubeIntegration, err := item[1].(project.Integration).ToAtlas(r.Client, namespace)
		if kubeIntegration == nil {
			ctx.Log.Warnw("Update Integrations", "Can not convert kube integration", err)
			return workflow.Terminate(workflow.ProjectIntegrationInternal, "Update Integrations: Can not convert kube integration")
		}
		t := mongodbatlas.ThirdPartyIntegration(atlasIntegration)
		if &t != kubeIntegration {
			ctx.Log.Debugf("Try to update integration: %s", kubeIntegration.Type)
			if _, _, err := ctx.Client.Integrations.Replace(context.Background(), projectID, kubeIntegration.Type, kubeIntegration); err != nil {
				return workflow.Terminate(workflow.ProjectIntegrationRequest, "Can not convert integration")
			}
		}
	}
	return workflow.OK()
}

func deleteIntegrationsFromAtlas(ctx *workflow.Context, projectID string, integrationsToRemove []set.Identifiable) error {
	for _, integration := range integrationsToRemove {
		if _, err := ctx.Client.Integrations.Delete(context.Background(), projectID, integration.Identifier().(string)); err != nil {
			return err
		}
		ctx.Log.Debugf("Third Party Integration deleted: %s", integration.Identifier())
	}
	return nil
}

func (r *AtlasProjectReconciler) createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrations []set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrations {
		integration, err := item.(project.Integration).ToAtlas(r.Client, namespace)
		if err != nil || integration == nil {
			return workflow.Terminate(workflow.ProjectIntegrationInternal, fmt.Sprintf("cannot convert integration: %s", err.Error()))
		}

		_, resp, err := ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)
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
		if isPrometheusType(atlas.Type) {
			areEqual = arePrometheusesEqual(atlas, spec)
		} else {
			specAsAtlas, _ := spec.ToAtlas(r.Client, namespace)
			specAlias := aliasThirdPartyIntegration(*specAsAtlas)
			areEqual = AreIntegrationsEqual(&atlas, &specAlias)
		}
		ctx.Log.Debugw("checkIntegrationsReady", "atlas", atlas, "spec", spec, "areEqual", areEqual)

		if !areEqual {
			return false
		}
	}

	return true
}

func AreIntegrationsEqual(atlas, specAsAtlas *aliasThirdPartyIntegration) bool {
	return reflect.DeepEqual(cleanCopyToCompare(atlas), cleanCopyToCompare(specAsAtlas))
}

func cleanCopyToCompare(input *aliasThirdPartyIntegration) *aliasThirdPartyIntegration {
	if input == nil {
		return input
	}

	result := *input
	keepLastFourChars(&result.APIKey)
	keepLastFourChars(&result.APIToken)
	keepLastFourChars(&result.LicenseKey)
	keepLastFourChars(&result.Password)
	keepLastFourChars(&result.ReadToken)
	keepLastFourChars(&result.RoutingKey)
	keepLastFourChars(&result.Secret)
	keepLastFourChars(&result.ServiceKey)
	keepLastFourChars(&result.WriteToken)

	return &result
}

func keepLastFourChars(strPtr *string) {
	if strPtr == nil {
		return
	}

	charCount := 4
	str := *strPtr
	if len(str) <= charCount {
		return
	}

	*strPtr = str[len(str)-charCount:]
}

type aliasThirdPartyIntegration mongodbatlas.ThirdPartyIntegration

func (i aliasThirdPartyIntegration) Identifier() interface{} {
	return i.Type
}

func toAliasThirdPartyIntegration(list []*mongodbatlas.ThirdPartyIntegration) []aliasThirdPartyIntegration {
	result := make([]aliasThirdPartyIntegration, len(list))
	for i, item := range list {
		result[i] = aliasThirdPartyIntegration(*item)
	}
	return result
}

func syncPrometheusStatus(ctx *workflow.Context, project *mdbv1.AtlasProject, integrationPairs [][]set.Identifiable) {
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

func searchAtlasIntegration(integrationPairs [][]set.Identifiable, filterFunc func(typeName string) bool) (integration mongodbatlas.ThirdPartyIntegration, found bool) {
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

func canIntegrationsReconcile(workflowCtx *workflow.Context, protected bool, akoProject *mdbv1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &mdbv1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	list, _, err := workflowCtx.Client.Integrations.List(workflowCtx.Context, akoProject.ID())
	if err != nil {
		return false, err
	}

	if list.TotalCount == 0 {
		return true, nil
	}

	atlasIntegrations := toAliasThirdPartyIntegration(list.Results)
	diff := set.Difference(atlasIntegrations, latestConfig.Integrations)

	if len(diff) == 0 {
		return true, nil
	}

	diff = set.Difference(akoProject.Spec.Integrations, atlasIntegrations)

	return len(diff) == 0, nil
}
