package atlasproject

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func ensureIntegration(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	integrationList := project.Spec.DeepCopy().Integrations
	if result := createOrDeleteIntegrationInAtlas(ctx, projectID, integrationList); !result.IsOk() {
		return result
	}
	return workflow.OK()
}

func createOrDeleteIntegrationInAtlas(ctx *workflow.Context, projectID string, requestedIntegrations []project.Integration) workflow.Result {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(context.Background(), projectID)
	ctx.Log.Debugf("integrationsInAtlas: %v", integrationsInAtlas)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasOthers, err.Error())
	}

	currentIntegrationsInAtlas := fromAtlas(integrationsInAtlas.Results) // TODO rename ^
	ctx.Log.Debugf("currentIntegrationsInAtlas: %v", currentIntegrationsInAtlas)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasOthers, err.Error())
	}

	indentificatorsForDelete := set.Difference(currentIntegrationsInAtlas, requestedIntegrations)
	ctx.Log.Debugf("indentificatorsForDelete: %v", indentificatorsForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, indentificatorsForDelete, ctx.Log); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasOthers, err.Error())
	}

	// integrationsToUpdate := set.Intersection(currentIntegrationsInAtlas, requestedIntegrations) // TODO ??

	indentificatorsForCreate := set.Difference(requestedIntegrations, currentIntegrationsInAtlas)
	ctx.Log.Debugf("indentificatorsForCreate: %v", indentificatorsForCreate)
	if result := createIntegrationsInAtlas(ctx, projectID, indentificatorsForCreate); !result.IsOk() {
		return result
	}

	return workflow.OK()
}

func deleteIntegrationsFromAtlas(ctx *workflow.Context, projectID string, integrationsToRemove []set.Identifiable, log *zap.SugaredLogger) error {
	for _, integration := range integrationsToRemove {
		if _, err := ctx.Client.Integrations.Delete(context.Background(), projectID, integration.Identifier().(string)); err != nil {
			return err
		}
		log.Debugw("Third Party Integration deleted: ", integration.Identifier())
	}
	return nil
}

func createIntegrationsInAtlas(ctx *workflow.Context, projectID string, integrations []set.Identifiable) workflow.Result {
	for _, item := range integrations {
		integration, err := item.(project.Integration).ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasOthers, err.Error())
		}
		ctx.Log.Warnf("WARNING!!!!!!!!!!!!", integration.Type, integration.Region, integration.APIKey)
		// TODO do we need thirdPartIntegration results here?
		_, _, err = ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)

		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, err.Error())
		}
	}
	return workflow.OK()
}

// =======================
func buildMap(integrations []project.Integration) map[string]project.Integration {
	newMap := map[string]project.Integration{}
	for _, item := range integrations {
		newMap[item.Identifier().(string)] = item
	}
	return newMap
}

func fromAtlas(source []*mongodbatlas.ThirdPartyIntegration) []project.Integration {
	result := make([]project.Integration, len(source))
	for i, item := range source {
		result[i] = project.Integration(*item)
		fmt.Print(result[i].Type)
	}
	return result
}
