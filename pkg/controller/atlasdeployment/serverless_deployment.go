package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func ensureServerlessInstanceState(ctx *workflow.Context, project *mdbv1.AtlasProject, serverlessSpec *mdbv1.ServerlessSpec) (atlasDeployment *mongodbatlas.Cluster, _ workflow.Result) {
	atlasDeployment, resp, err := ctx.Client.ServerlessInstances.Get(context.Background(), project.Status.ID, serverlessSpec.Name)
	if err != nil {
		if resp == nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return atlasDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}

		atlasDeployment, err = serverlessSpec.ServerlessToAtlas()
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}
		ctx.Log.Infof("Serverless Instance %s doesn't exist in Atlas - creating", serverlessSpec.Name)
		atlasDeployment, _, err = ctx.Client.ServerlessInstances.Create(context.Background(), project.Status.ID, &mongodbatlas.ServerlessCreateRequestParams{
			Name: serverlessSpec.Name,
			ProviderSettings: &mongodbatlas.ServerlessProviderSettings{
				BackingProviderName: serverlessSpec.ProviderSettings.BackingProviderName,
				ProviderName:        string(serverlessSpec.ProviderSettings.ProviderName),
				RegionName:          serverlessSpec.ProviderSettings.RegionName,
			},
			Tag: atlasDeployment.Tags,
		})
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}
	}

	switch atlasDeployment.StateName {
	case status.StateIDLE:
		convertedDeployment, err := serverlessSpec.ServerlessToAtlas()
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}
		fmt.Println("TAGS >>> operator:", convertedDeployment.Tags, " \n atlas", atlasDeployment.Tags)
		if convertedDeployment.Tags == nil {
			convertedDeployment.Tags = []*mongodbatlas.Tag{}
		}
		if (reflect.DeepEqual(convertedDeployment.Tags, atlasDeployment.Tags)) != true /* TODO : add || if backupOptions aren't equal */ {

			fmt.Println(reflect.DeepEqual(convertedDeployment.Tags, atlasDeployment.Tags))

			atlasDeployment, _, err = ctx.Client.ServerlessInstances.Update(context.Background(), project.Status.ID, serverlessSpec.Name, &mongodbatlas.ServerlessUpdateRequestParams{
				ServerlessBackupOptions: convertedDeployment.ServerlessBackupOptions,
				Tag:                     convertedDeployment.Tags,
			})
			if err != nil {
				return atlasDeployment, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
			}
			return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
		}
		result := ensureServerlessPrivateEndpoints(ctx, project.ID(), serverlessSpec, atlasDeployment.Name)
		return atlasDeployment, result

	case status.StateCREATING:
		return atlasDeployment, workflow.InProgress(workflow.DeploymentCreating, "deployment is provisioning")

	case status.StateUPDATING, status.StateREPAIRING:
		return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return atlasDeployment, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown deployment state %q", atlasDeployment.StateName))
	}
}
