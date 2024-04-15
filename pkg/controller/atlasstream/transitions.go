package atlasstream

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *InstanceReconciler) create(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
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
		return r.terminate(workflowCtx, workflow.StreamInstanceNotCreated, err)
	}

	return r.ready(workflowCtx, atlasStreamInstance)
}

func (r *InstanceReconciler) update(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	dataProcessRegion := admin.StreamsDataProcessRegion{
		CloudProvider: streamInstance.Spec.Config.Provider,
		Region:        streamInstance.Spec.Config.Region,
	}

	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		UpdateStreamInstance(workflowCtx.Context, project.ID(), streamInstance.Spec.Name, &dataProcessRegion).
		Execute()

	if err != nil {
		r.terminate(workflowCtx, workflow.StreamInstanceNotUpdated, err)
	}

	return r.ready(workflowCtx, atlasStreamInstance)
}

func (r *InstanceReconciler) delete(workflowCtx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(streamInstance, r.ObjectDeletionProtection) {
		workflowCtx.Log.Info("Not removing AtlasStreamInstance from Atlas as per configuration")
	} else {
		if err := deleteStreamInstance(workflowCtx, project, streamInstance); err != nil {
			return r.terminate(workflowCtx, workflow.StreamInstanceNotRemoved, err)
		}
	}
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, streamInstance, customresource.UnsetFinalizer); err != nil {
		return r.terminate(workflowCtx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.OK().ReconcileResult(), nil
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

// transitions back to pending state
// also terminates if an terminate occured
func (r *InstanceReconciler) skip(ctx context.Context, log *zap.SugaredLogger, streamInstance *akov2.AtlasStreamInstance) ctrl.Result {
	log.Infow(fmt.Sprintf("-> Skipping AtlasStreamInstance reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", streamInstance.Spec)
	if !streamInstance.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, streamInstance, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err.Error())
			log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult()
		}
	}

	return workflow.OK().ReconcileResult()
}

// transitions back to pending state setting an terminate state
func (r *InstanceReconciler) invalidate(invalid workflow.Result) (ctrl.Result, error) {
	r.Log.Debugf("AtlasStreamInstance is invalid: %v", invalid)
	return invalid.ReconcileResult(), nil
}

// transitions back to pending setting unsupported state
func (r *InstanceReconciler) unsupport(ctx *workflow.Context) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, "the AtlasStreamInstance is not supported by Atlas for government").
		WithoutRetry()
	ctx.SetConditionFromResult(status.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

// transitions back to pending state setting an error status
func (r *InstanceReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err.Error())
	ctx.SetConditionFromResult(status.StreamInstanceReadyType, terminated)
	return terminated.ReconcileResult(), nil
}

func (r *InstanceReconciler) ready(workflowCtx *workflow.Context, streamInstance *admin.StreamsTenant) (ctrl.Result, error) {
	workflowCtx.EnsureStatusOption(status.AtlasStreamInstanceDetails(streamInstance.GetId(), streamInstance.GetHostnames()))
	result := workflow.OK()
	workflowCtx.SetConditionFromResult(status.StreamInstanceReadyType, result)
	return result.ReconcileResult(), nil
}
