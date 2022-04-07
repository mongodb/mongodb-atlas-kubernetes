package atlasproject

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

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
	integrationsInAtlas, err := getIntegrationsAndConvert(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}

	indentificatorsForDelete := set.Difference(integrationsInAtlas, requestedIntegrations)
	ctx.Log.Debugf("indentificatorsForDelete: %v", indentificatorsForDelete)
	if err := deleteIntegrationsFromAtlas(ctx, projectID, indentificatorsForDelete, ctx.Log); err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}

	integrationsToUpdate := set.Intersection(integrationsInAtlas, requestedIntegrations)
	ctx.Log.Debugf("integrationsToUpdate: %v", integrationsToUpdate)
	if result := updateIntegrationsAtlas(ctx, projectID, integrationsToUpdate, ctx.Log); !result.IsOk() {
		return result
	}

	indentificatorsForCreate := set.Difference(requestedIntegrations, integrationsInAtlas)
	ctx.Log.Debugf("indentificatorsForCreate: %v", indentificatorsForCreate)
	if result := createIntegrationsInAtlas(ctx, projectID, indentificatorsForCreate); !result.IsOk() {
		return result
	}

	ready, err := checkIntegrationsReady(ctx, projectID, requestedIntegrations)
	if err != nil {
		return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
	}
	if !ready {
		ctx.SetConditionFalse(status.IntegrationReadyType)
		workflow.InProgress(workflow.ProjectIntegrationInAtlasInternal, "....in progress")
	}
	ctx.SetConditionTrue(status.IntegrationReadyType)
	return workflow.OK()
}

func updateIntegrationsAtlas(ctx *workflow.Context, projectID string, integrationsToUpdate [][]set.Identifiable, log *zap.SugaredLogger) workflow.Result {
	for _, item := range integrationsToUpdate {
		atlasIntegration, err := item[0].(project.Integration).ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
		}

		kubeIntegration, err := item[1].(project.Integration).ToAtlas()
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
		}
		if &atlasIntegration != &kubeIntegration {
			log.Debugf("Try to update integration: %s", kubeIntegration.Type)
			if _, _, err = ctx.Client.Integrations.Replace(context.Background(), projectID, kubeIntegration.Type, kubeIntegration); err != nil {
				return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, err.Error())
			}
		}
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
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasInternal, err.Error())
		}

		_, _, err = ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)

		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlasRequest, err.Error())
		}
	}
	return workflow.OK()
}

func checkIntegrationsReady(ctx *workflow.Context, projectID string, requestedIntegrations []project.Integration) (bool, error) {
	atlasIntegrations, err := getIntegrationsAndConvert(ctx, projectID)
	if err != nil {
		return false, err
	}
	if reflect.DeepEqual(atlasIntegrations, requestedIntegrations) {
		ctx.SetConditionTrue(status.IntegrationReadyType)
		return true, nil
	}
	return false, err
}

func getIntegrationsAndConvert(ctx *workflow.Context, projectID string) ([]project.Integration, error) {
	integrationsInAtlas, _, err := ctx.Client.Integrations.List(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	ctx.Log.Debugf("Got Integrations From Atlas: %v", &integrationsInAtlas)

	convertedIntegrationsInAtlas := fromAtlas(integrationsInAtlas.Results)
	if err != nil {
		return nil, err
	}
	return convertedIntegrationsInAtlas, nil
}

func fromAtlas(source []*mongodbatlas.ThirdPartyIntegration) []project.Integration {
	result := make([]project.Integration, len(source))
	for i, item := range source {
		result[i] = project.Integration(*item)
		fmt.Print(result[i].Type)
	}
	return result
}
