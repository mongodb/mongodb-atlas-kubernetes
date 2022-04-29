package atlasproject

import (
	"context"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

type integrations struct {
	list             []project.Integration
	projectNamespace string
}

func (r *AtlasProjectReconciler) ensureIntegration(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	integrationList := integrations{
		list:             project.Spec.Integrations,
		projectNamespace: project.Namespace,
	}
	if result := createOrDeleteIntegrationsInAtlas(ctx, r.Client, projectID, integrationList); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

func createOrDeleteIntegrationsInAtlas(ctx *workflow.Context, c client.Client, projectID string, requestedIntegrations integrations) workflow.Result {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(context.Background(), projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}
	ctx.Log.Debugf("Got Integrations From Atlas: %v", *integrationsInAtlas)
	integrationsInAtlasAlias := toAliasThirdPartyIntegration(integrationsInAtlas.Results)

	indentificatorsForDelete := set.Difference(integrationsInAtlasAlias, requestedIntegrations.list)
	ctx.Log.Debugf("indentificatorsForDelete: %v", indentificatorsForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, indentificatorsForDelete); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}

	integrationsToUpdate := set.Intersection(integrationsInAtlasAlias, requestedIntegrations.list)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := updateIntegrationsAtlas(ctx, c, projectID, integrationsToUpdate, requestedIntegrations.projectNamespace); !result.IsOk() {
		return result
	}

	indentificatorsForCreate := set.Difference(requestedIntegrations.list, integrationsInAtlasAlias)
	ctx.Log.Debugf("indentificatorsForCreate: %v", indentificatorsForCreate)
	if result := createIntegrationsInAtlas(ctx, c, projectID, indentificatorsForCreate, requestedIntegrations.projectNamespace); !result.IsOk() {
		return result
	}

	ready, err := checkIntegrationsReady(ctx, c, projectID, requestedIntegrations)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}
	if !ready {
		ctx.SetConditionFalse(status.IntegrationReadyType)
		return workflow.InProgress(workflow.ProjectIntegrationInAtlasInternal, "in progress")
	}
	if len(requestedIntegrations.list) > 0 {
		ctx.SetConditionTrue(status.IntegrationReadyType)
	}
	return workflow.OK()
}

func updateIntegrationsAtlas(ctx *workflow.Context, c client.Client, projectID string, integrationsToUpdate [][]set.Identifiable, defaultNS string) workflow.Result {
	for _, item := range integrationsToUpdate {
		atlasIntegration := item[0].(aliasThirdPartyIntegration)
		kubeIntegration := item[1].(project.Integration).ToAtlas(defaultNS, c)
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

func createIntegrationsInAtlas(ctx *workflow.Context, c client.Client, projectID string, integrations []set.Identifiable, defaultNS string) workflow.Result {
	for _, item := range integrations {
		integration := item.(project.Integration).ToAtlas(defaultNS, c)
		if integration == nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, "Can not convert integration")
		}

		_, _, err := ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, err.Error())
		}
	}
	return workflow.OK()
}

func checkIntegrationsReady(ctx *workflow.Context, c client.Client, projectID string, requestedIntegrations integrations) (bool, error) {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(context.Background(), projectID)
	if err != nil {
		return false, err
	}

	requestedIntegrationsConverted := convertToAtlasIntegrationList(requestedIntegrations, c)
	if reflect.DeepEqual(integrationsInAtlas.Results, requestedIntegrationsConverted) {
		return true, nil
	}
	return false, err
}

func convertToAtlasIntegrationList(list integrations, c client.Client) []*mongodbatlas.ThirdPartyIntegration {
	result := make([]*mongodbatlas.ThirdPartyIntegration, len(list.list))
	for i, item := range list.list {
		result[i] = item.ToAtlas(list.projectNamespace, c)
	}
	return result
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
