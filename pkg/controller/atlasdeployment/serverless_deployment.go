package atlasdeployment

import (
	"fmt"
	"net/http"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"go.mongodb.org/atlas/mongodbatlas"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDeploymentReconciler) ensureServerlessInstanceState(workflowCtx *workflow.Context, project *akov2.AtlasProject, deployment *akov2.AtlasDeployment) (atlasDeployment *mongodbatlas.Cluster, _ workflow.Result) {
	if deployment == nil || deployment.Spec.ServerlessSpec == nil {
		return nil, workflow.Terminate(workflow.ServerlessPrivateEndpointReady, "deployment spec is empty")
	}
	serverlessSpec := deployment.Spec.ServerlessSpec
	atlasDeployment, resp, err := workflowCtx.Client.ServerlessInstances.Get(workflowCtx.Context, project.Status.ID, serverlessSpec.Name)
	if err != nil {
		if resp == nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return atlasDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}

		atlasDeployment, err = serverlessSpec.ToAtlas()
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}
		workflowCtx.Log.Infof("Serverless Instance %s doesn't exist in Atlas - creating", serverlessSpec.Name)
		atlasDeployment, _, err = workflowCtx.Client.ServerlessInstances.Create(workflowCtx.Context, project.Status.ID, &mongodbatlas.ServerlessCreateRequestParams{
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
		convertedDeployment, err := serverlessSpec.ToAtlas()
		if err != nil {
			return atlasDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}
		if convertedDeployment.Tags == nil {
			convertedDeployment.Tags = &[]*mongodbatlas.Tag{}
		}
		if !isTagsEqual(*(atlasDeployment.Tags), *(convertedDeployment.Tags)) {
			atlasDeployment, _, err = workflowCtx.Client.ServerlessInstances.Update(workflowCtx.Context, project.Status.ID, serverlessSpec.Name, &mongodbatlas.ServerlessUpdateRequestParams{
				Tag: convertedDeployment.Tags,
				ServerlessBackupOptions: &mongodbatlas.ServerlessBackupOptions{
					ServerlessContinuousBackupEnabled: &serverlessSpec.BackupOptions.ServerlessContinuousBackupEnabled,
				},
				TerminationProtectionEnabled: &serverlessSpec.TerminationProtectionEnabled,
			})
			if err != nil {
				return atlasDeployment, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
			}
			return atlasDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
		}
		result := ensureServerlessPrivateEndpoints(workflowCtx, project.ID(), deployment, atlasDeployment.Name, r.SubObjectDeletionProtection)
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

func isTagsEqual(a []*mongodbatlas.Tag, c []*mongodbatlas.Tag) bool {
	if len(a) == len(c) {
		for i, aTags := range a {
			if aTags.Key != c[i].Key || aTags.Value != c[i].Value {
				return false
			}
		}
		return true
	}
	return false
}
