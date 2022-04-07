package atlasproject

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func (r *AtlasProjectReconciler) ensureIntegration(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	integrationsInAtlas, err := fetchIntegrations(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}
	integrationsInAtlasAlias := toAliasThirdPartyIntegration(integrationsInAtlas.Results)

	indentificatorsForDelete := set.Difference(integrationsInAtlasAlias, project.Spec.Integrations)
	ctx.Log.Debugf("indentificatorsForDelete: %v", indentificatorsForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, indentificatorsForDelete); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}

	integrationsToUpdate := set.Intersection(integrationsInAtlasAlias, project.Spec.Integrations)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := r.updateIntegrationsAtlas(ctx, projectID, integrationsToUpdate, project.Namespace); !result.IsOk() {
		return result
	}

	indentificatorsForCreate := set.Difference(project.Spec.Integrations, integrationsInAtlasAlias)
	ctx.Log.Debugf("indentificatorsForCreate: %v", indentificatorsForCreate)
	if result := r.createIntegrationsInAtlas(ctx, projectID, indentificatorsForCreate, project.Namespace); !result.IsOk() {
		return result
	}

	setPrometheusStatus(project, integrationsInAtlas)
	if ready := r.checkIntegrationsReady(ctx, project.Namespace, integrationsToUpdate, project.Spec.Integrations); !ready {
		ctx.SetConditionFalse(status.IntegrationReadyType)
		return workflow.InProgress(workflow.ProjectIntegrationInAtlasInternal, "in progress")
	}
	if len(project.Spec.Integrations) > 0 {
		ctx.SetConditionTrue(status.IntegrationReadyType)
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
		kubeIntegration, _ := item[1].(project.Integration).ToAtlas(r.Client, namespace)
		if kubeIntegration == nil {
			ctx.Log.Warn("Update Integrations: Can not convert kube integration")
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, "Update Integrations: Can not convert kube integration")
		}
		t := mongodbatlas.ThirdPartyIntegration(atlasIntegration)
		if &t != kubeIntegration {
			ctx.Log.Debugf("Try to update integration: %s", kubeIntegration.Type)
			if _, _, err := ctx.Client.Integrations.Replace(context.Background(), projectID, kubeIntegration.Type, kubeIntegration); err != nil {
				return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, "Can not convert integration")
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
		ctx.Log.Debugw("Third Party Integration deleted: ", integration.Identifier())
	}
	return nil
}

func (r *AtlasProjectReconciler) createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrations []set.Identifiable, namespace string) workflow.Result {
	for _, item := range integrations {
		integration, err := item.(project.Integration).ToAtlas(r.Client, namespace)
		if err != nil {
			ctx.Log.Debugw("cannot convert integration", "err", err)
		}
		if integration == nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, "Can not convert integration")
		}

		_, _, err = ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, err.Error())
		}
	}
	return workflow.OK()
}

func (r *AtlasProjectReconciler) checkIntegrationsReady(ctx *workflow.Context, namespace string, integrationsIntersection [][]set.Identifiable, requestedIntegrations []project.Integration) bool {
	if len(integrationsIntersection) != len(requestedIntegrations) {
		return false
	}

	ctx.Log.Debugw("checkIntegrationsReady", "integrationsIntersection", integrationsIntersection)
	for _, integrationPair := range integrationsIntersection {
		atlas := integrationPair[0].(aliasThirdPartyIntegration)
		spec := integrationPair[1].(project.Integration)
		var areEqual bool
		if isPrometheusType(atlas.Type) {
			areEqual = arePrometheusesEqual(atlas, spec)
		} else {
			specAsAtlas, _ := spec.ToAtlas(r.Client, namespace)
			areEqual = reflect.DeepEqual(atlas, specAsAtlas)
		}

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
		result[i] = aliasThirdPartyIntegration(*item)
	}
	return result
}

func setPrometheusStatus(project *mdbv1.AtlasProject, atlasIntegrations *mongodbatlas.ThirdPartyIntegrations) {
	for _, atlasIntegration := range atlasIntegrations.Results {
		if isPrometheusType(atlasIntegration.Type) {
			project.Status.Prometheus.DiscoveryURL = buildPrometheusDiscoveryURL(project.ID())
		}
	}
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

func buildPrometheusDiscoveryURL(projectID string) string {
	api := "https://cloud.mongodb.com/api/atlas/v1.0"
	return fmt.Sprintf("%s/groups/%s/discovery", api, projectID)
}
