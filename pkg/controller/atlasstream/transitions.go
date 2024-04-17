package atlasstream

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *InstanceReconciler) create(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	streamInstance *akov2.AtlasStreamInstance,
	mapper streamConnectionMapper,
) (ctrl.Result, error) {
	connections := make([]admin.StreamsConnection, 0, len(streamInstance.Spec.ConnectionRegistry))
	streamTenant := admin.StreamsTenant{
		Name: &streamInstance.Name,
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			CloudProvider: streamInstance.Spec.Config.Provider,
			Region:        streamInstance.Spec.Config.Region,
		},
		GroupId:     pointer.MakePtr(project.ID()),
		Connections: &connections,
	}

	for _, connectionRef := range streamInstance.Spec.ConnectionRegistry {
		streamConnection := akov2.AtlasStreamConnection{}
		err := r.Client.Get(ctx.Context, *connectionRef.GetObject(connectionRef.Namespace), &streamConnection)
		if err != nil {
			return r.terminate(ctx, workflow.StreamInstanceNotCreated, fmt.Errorf("failed to retrieve connection %v: %w", connectionRef, err))
		}

		connection, err := mapper(&streamConnection)
		if err != nil {
			return r.terminate(ctx, workflow.StreamInstanceNotCreated, err)
		}

		connections = append(connections, *connection)
	}

	atlasStreamInstance, _, err := ctx.SdkClient.StreamsApi.
		CreateStreamInstance(ctx.Context, project.ID(), &streamTenant).
		Execute()

	if err != nil {
		return r.terminate(ctx, workflow.StreamInstanceNotCreated, err)
	}

	return r.ready(ctx, atlasStreamInstance)
}

func (r *InstanceReconciler) update(ctx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) error {
	dataProcessRegion := admin.StreamsDataProcessRegion{
		CloudProvider: streamInstance.Spec.Config.Provider,
		Region:        streamInstance.Spec.Config.Region,
	}

	_, _, err := ctx.SdkClient.StreamsApi.
		UpdateStreamInstance(ctx.Context, project.ID(), streamInstance.Spec.Name, &dataProcessRegion).
		Execute()

	if err != nil {
		return err
	}

	return nil
}

func (r *InstanceReconciler) delete(ctx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(streamInstance, r.ObjectDeletionProtection) {
		ctx.Log.Info("Not removing AtlasStreamInstance from Atlas as per configuration")
	} else {
		if err := deleteStreamInstance(ctx, project, streamInstance); err != nil {
			return r.terminate(ctx, workflow.StreamInstanceNotRemoved, err)
		}
	}
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, streamInstance, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.OK().ReconcileResult(), nil
}

func deleteStreamInstance(ctx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) error {
	_, _, err := ctx.SdkClient.StreamsApi.
		DeleteStreamInstance(ctx.Context, project.ID(), streamInstance.Spec.Name).
		Execute()

	if err != nil && !admin.IsErrorCode(err, atlas.ResourceNotFound) {
		return err
	}

	return nil
}

func createConnections(
	workflowCtx *workflow.Context,
	project *akov2.AtlasProject,
	akoStreamInstance *akov2.AtlasStreamInstance,
	akoStreamConnections []*akov2.AtlasStreamConnection,
	mapper streamConnectionMapper,
) error {
	for _, akoStreamConnection := range akoStreamConnections {
		connection, err := mapper(akoStreamConnection)
		if err != nil {
			return err
		}

		_, _, err = workflowCtx.SdkClient.StreamsApi.
			CreateStreamConnection(workflowCtx.Context, project.ID(), akoStreamInstance.Spec.Name, connection).
			Execute()

		if err != nil {
			return err
		}

		workflowCtx.EnsureStatusOption(
			status.AtlasStreamInstanceAddConnection(
				connection.GetName(),
				common.ResourceRefNamespaced{
					Name:      akoStreamConnection.Name,
					Namespace: akoStreamConnection.Namespace,
				},
			),
		)
	}

	return nil
}

func updateConnections(
	workflowCtx *workflow.Context,
	project *akov2.AtlasProject,
	akoStreamInstance *akov2.AtlasStreamInstance,
	akoStreamConnections []*akov2.AtlasStreamConnection,
	mapper streamConnectionMapper,
) error {
	for _, akoStreamConnection := range akoStreamConnections {
		connection, err := mapper(akoStreamConnection)
		if err != nil {
			return err
		}

		_, _, err = workflowCtx.SdkClient.StreamsApi.
			UpdateStreamConnection(workflowCtx.Context, project.ID(), akoStreamInstance.Spec.Name, akoStreamConnection.Spec.Name, connection).
			Execute()

		if err != nil {
			return err
		}

		workflowCtx.EnsureStatusOption(
			status.AtlasStreamInstanceAddConnection(
				connection.GetName(),
				common.ResourceRefNamespaced{
					Name:      akoStreamConnection.Name,
					Namespace: akoStreamConnection.Namespace,
				},
			),
		)
	}

	return nil
}

func deleteConnections(
	workflowCtx *workflow.Context,
	project *akov2.AtlasProject,
	streamInstance *akov2.AtlasStreamInstance,
	atlasStreamConnections []*admin.StreamsConnection,
) error {
	for _, atlasStreamConnection := range atlasStreamConnections {
		_, _, err := workflowCtx.SdkClient.StreamsApi.
			DeleteStreamConnection(workflowCtx.Context, project.ID(), streamInstance.Spec.Name, atlasStreamConnection.GetName()).
			Execute()

		if err != nil && !admin.IsErrorCode(err, atlas.ResourceNotFound) {
			return err
		}

		workflowCtx.EnsureStatusOption(status.AtlasStreamInstanceRemoveConnection(atlasStreamConnection.GetName()))
	}

	return nil
}

// transitions back to pending state
// also terminates if a "terminate" occurred
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
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
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

func (r *InstanceReconciler) ready(ctx *workflow.Context, streamInstance *admin.StreamsTenant) (ctrl.Result, error) {
	ctx.EnsureStatusOption(status.AtlasStreamInstanceDetails(streamInstance.GetId(), streamInstance.GetHostnames()))
	result := workflow.OK()
	ctx.SetConditionFromResult(status.StreamInstanceReadyType, result)
	return result.ReconcileResult(), nil
}
