package atlasstream

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	ctrl "sigs.k8s.io/controller-runtime"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *InstanceReconciler) handlePendingOrReady(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		GetStreamInstance(workflowCtx.Context, project.ID(), streamInstance.Spec.Name).
		Execute()
	// Fail when API doesn't succeed and failure is not resource not found
	if err != nil && !admin.IsErrorCode(err, atlas.ResourceNotFound) {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	isMarkedAsDeleted := !streamInstance.GetDeletionTimestamp().IsZero()
	isNotInAtlas := err != nil && admin.IsErrorCode(err, atlas.ResourceNotFound)

	switch {
	case isNotInAtlas && !isMarkedAsDeleted:
		// if no streams processing instance is not in atlas and is not marked as deleted - create
		return r.create(workflowCtx, project, streamInstance)
	case isMarkedAsDeleted:
		// if a streams processing instance is marked as deleted,
		// independently whether it exists in Atlas or not - delete
		return r.delete(workflowCtx, project, streamInstance)
	case hasChanged(streamInstance, atlasStreamInstance):
		// if a streams processing instance is ready and has changed - update
		return r.update(workflowCtx, project, streamInstance)
	default:
		// no change, streams processing instance stays in ready or pending state
		return workflow.OK().ReconcileResult(), nil
	}
}

func hasChanged(streamInstance *akov2.AtlasStreamInstance, atlasStreamInstance *admin.StreamsTenant) bool {
	config := streamInstance.Spec.Config
	dataProcessRegion := atlasStreamInstance.GetDataProcessRegion()

	return config.Provider != dataProcessRegion.GetCloudProvider() || config.Region == dataProcessRegion.GetRegion()
}
