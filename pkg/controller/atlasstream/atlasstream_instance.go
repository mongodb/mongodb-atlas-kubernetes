package atlasstream

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *InstanceReconciler) ensureStreamInstance(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) workflow.Result {
	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		GetStreamInstance(workflowCtx.Context, project.ID(), streamInstance.Spec.Name).
		Execute()
	// Fail when API doesn't succeed and failure is not resource not found
	if err != nil && !admin.IsErrorCode(err, atlas.ResourceNotFound) {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	// Create resource if instance doesn't exist and k8s resource wasn't deleted
	if err != nil && admin.IsErrorCode(err, atlas.ResourceNotFound) && streamInstance.GetDeletionTimestamp().IsZero() {
		err = createStreamInstance(workflowCtx, project, streamInstance)
		if err != nil {
			return workflow.Terminate(workflow.StreamInstanceNotCreated, err.Error())
		}

		return workflow.OK()
	}

	if !streamInstance.GetDeletionTimestamp().IsZero() {
		if customresource.IsResourcePolicyKeepOrDefault(streamInstance, r.ObjectDeletionProtection) {
			workflowCtx.Log.Info("Not removing AtlasStreamInstance from Atlas as per configuration")
		} else {
			err = deleteStreamInstance(workflowCtx, project, streamInstance)
			if err != nil {
				return workflow.Terminate(workflow.StreamInstanceNotRemoved, err.Error())
			}
		}

		if err = customresource.ManageFinalizer(workflowCtx.Context, r.Client, streamInstance, customresource.UnsetFinalizer); err != nil {
			workflowCtx.Log.Errorw("failed to remove finalizer", "error", err)

			return workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
		}
	}

	if HasStreamInstanceChanged(streamInstance, atlasStreamInstance) {
		err = updateStreamInstance(workflowCtx, project, streamInstance)
		if err != nil {
			return workflow.Terminate(workflow.StreamInstanceNotUpdated, err.Error())
		}
	}

	return workflow.OK()
}

func createStreamInstance(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) error {
	streamTenant := admin.StreamsTenant{
		Name: &streamInstance.Name,
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			CloudProvider: streamInstance.Spec.Config.Provider,
			Region:        streamInstance.Spec.Config.Region,
		},
		GroupId: pointer.MakePtr(project.ID()),
	}

	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		CreateStreamInstance(workflowCtx.Context, project.ID(), &streamTenant).
		Execute()

	if err != nil {
		return err
	}

	workflowCtx.EnsureStatusOption(status.AtlasStreamInstanceDetails(atlasStreamInstance.GetId(), atlasStreamInstance.GetHostnames()))

	return nil
}

func deleteStreamInstance(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) error {
	_, _, err := workflowCtx.SdkClient.StreamsApi.
		DeleteStreamInstance(workflowCtx.Context, project.ID(), streamInstance.Spec.Name).
		Execute()

	if err != nil && !admin.IsErrorCode(err, atlas.ResourceNotFound) {
		return err
	}

	return nil
}

func updateStreamInstance(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) error {
	dataProcessRegion := admin.StreamsDataProcessRegion{
		CloudProvider: streamInstance.Spec.Config.Provider,
		Region:        streamInstance.Spec.Config.Region,
	}

	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		UpdateStreamInstance(workflowCtx.Context, project.ID(), streamInstance.Spec.Name, &dataProcessRegion).
		Execute()

	if err != nil {
		return err
	}

	workflowCtx.EnsureStatusOption(status.AtlasStreamInstanceDetails(atlasStreamInstance.GetId(), atlasStreamInstance.GetHostnames()))

	return nil
}

func HasStreamInstanceChanged(streamInstance *akov2.AtlasStreamInstance, atlasStreamInstance *admin.StreamsTenant) bool {
	config := streamInstance.Spec.Config
	dataProcessRegion := atlasStreamInstance.GetDataProcessRegion()

	return config.Provider != dataProcessRegion.GetCloudProvider() || config.Region == dataProcessRegion.GetRegion()
}
