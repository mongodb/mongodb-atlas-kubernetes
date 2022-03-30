package atlasproject

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func validateIntegrationLists(integrationList []project.Intergation) error {
	for _, item := range integrationList {
		if err := validateSingleIntegration(item); err != nil {
			return err
		}
	}
	return nil
}

func validateSingleIntegration(integration project.Intergation) error {

	return nil
}

func ensureIntegration(ctx *workflow.Context, projectID string, project *mdbv1.AtlasProject) workflow.Result {
	integrationList := project.Spec.DeepCopy().Integrations
	if result := createIntegrations(ctx, projectID, integrationList); !result.IsOk() {
		return result
	}
	// ctx.EnsureStatusOption(status.)
	// integrations := project.Integration
	return workflow.OK()
}

func createOrDeleteIntegration2(client mongodbatlas.Client, projectID string) workflow.Result {
	panic("unimplemented")
}

// func createOrDeleteIntegration(ctx *workflow.Context, projectID string, integrations []project.Intergation) workflow.Result {
// 	return workflow.InProgress("test")
// }

func createIntegrations(ctx *workflow.Context, projectID string, integrations []project.Intergation) workflow.Result {

	// newIntegration := []mongodbatlas.ThirdPartyIntegration{}

	for _, item := range integrations {
		integration, err := item.ToAtlas()
		
		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlas, err.Error())
		}

		// TODO do we need thirdPartIntegration?
		_, _, err = ctx.Client.Integrations.Create(context.Background(), projectID, integration.Type, integration)

		if err != nil {
			return workflow.Terminate(workflow.ProjectIntegrationInAtlas, err.Error())
		}

	}

	return workflow.OK()
}
