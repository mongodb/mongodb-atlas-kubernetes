package atlasstream

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasStreamsInstanceReconciler) create(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	streamInstance *akov2.AtlasStreamInstance,
) (ctrl.Result, error) {
	streamTenant := admin.StreamsTenant{
		Name: &streamInstance.Spec.Name,
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			CloudProvider: streamInstance.Spec.Config.Provider,
			Region:        streamInstance.Spec.Config.Region,
		},
		GroupId: pointer.MakePtr(project.ID()),
	}

	atlasStreamInstance, _, err := ctx.SdkClient.StreamsApi.
		CreateStreamInstance(ctx.Context, project.ID(), &streamTenant).
		Execute()

	if err != nil {
		return r.terminate(ctx, workflow.StreamInstanceNotCreated, err)
	}

	return r.inProgress(ctx, atlasStreamInstance)
}

func (r *AtlasStreamsInstanceReconciler) update(ctx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	dataProcessRegion := admin.StreamsDataProcessRegion{
		CloudProvider: streamInstance.Spec.Config.Provider,
		Region:        streamInstance.Spec.Config.Region,
	}

	atlasStreamInstance, _, err := ctx.SdkClient.StreamsApi.
		UpdateStreamInstance(ctx.Context, project.ID(), streamInstance.Spec.Name, &dataProcessRegion).
		Execute()

	if err != nil {
		return r.terminate(ctx, workflow.StreamInstanceNotCreated, err)
	}

	return r.inProgress(ctx, atlasStreamInstance)
}

func (r *AtlasStreamsInstanceReconciler) delete(ctx *workflow.Context, project *akov2.AtlasProject, streamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(streamInstance, r.ObjectDeletionProtection) {
		r.Log.Info("Not removing AtlasStreamInstance from Atlas as per configuration")
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

	if err != nil && !admin.IsErrorCode(err, instanceNotFound) {
		return err
	}

	return nil
}

func createConnections(
	ctx *workflow.Context,
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

		_, _, err = ctx.SdkClient.StreamsApi.
			CreateStreamConnection(ctx.Context, project.ID(), akoStreamInstance.Spec.Name, connection).
			Execute()

		if err != nil {
			return err
		}

		ctx.EnsureStatusOption(
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
	ctx *workflow.Context,
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

		_, _, err = ctx.SdkClient.StreamsApi.
			UpdateStreamConnection(ctx.Context, project.ID(), akoStreamInstance.Spec.Name, akoStreamConnection.Spec.Name, connection).
			Execute()

		if err != nil {
			return err
		}

		ctx.EnsureStatusOption(
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
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	streamInstance *akov2.AtlasStreamInstance,
	atlasStreamConnections []*admin.StreamsConnection,
) error {
	for _, atlasStreamConnection := range atlasStreamConnections {
		_, _, err := ctx.SdkClient.StreamsApi.
			DeleteStreamConnection(ctx.Context, project.ID(), streamInstance.Spec.Name, atlasStreamConnection.GetName()).
			Execute()

		if err != nil && !admin.IsErrorCode(err, instanceNotFound) {
			return err
		}

		ctx.EnsureStatusOption(status.AtlasStreamInstanceRemoveConnection(atlasStreamConnection.GetName()))
	}

	return nil
}

// transitions back to pending state
// also terminates if a "terminate" occurred
func (r *AtlasStreamsInstanceReconciler) skip(ctx context.Context, log *zap.SugaredLogger, streamInstance *akov2.AtlasStreamInstance) ctrl.Result {
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
func (r *AtlasStreamsInstanceReconciler) invalidate(invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasStreamInstance is invalid: %v", invalid)
	return invalid.ReconcileResult(), nil
}

// transitions back to pending setting unsupported state
func (r *AtlasStreamsInstanceReconciler) unsupport(ctx *workflow.Context) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported, "the AtlasStreamInstance is not supported by Atlas for government").
		WithoutRetry()
	ctx.SetConditionFromResult(api.StreamInstanceReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

// transitions back to pending state setting an error status
func (r *AtlasStreamsInstanceReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err.Error())
	ctx.SetConditionFromResult(api.StreamInstanceReadyType, terminated)
	return terminated.ReconcileResult(), nil
}

func (r *AtlasStreamsInstanceReconciler) ready(ctx *workflow.Context, streamInstance *admin.StreamsTenant) (ctrl.Result, error) {
	ctx.EnsureStatusOption(status.AtlasStreamInstanceDetails(streamInstance.GetId(), streamInstance.GetHostnames()))
	result := workflow.OK()
	ctx.SetConditionFromResult(api.ReadyType, result)
	ctx.SetConditionFromResult(api.StreamInstanceReadyType, result)
	return result.ReconcileResult(), nil
}

func (r *AtlasStreamsInstanceReconciler) inProgress(ctx *workflow.Context, streamInstance *admin.StreamsTenant) (ctrl.Result, error) {
	ctx.EnsureStatusOption(status.AtlasStreamInstanceDetails(streamInstance.GetId(), streamInstance.GetHostnames()))
	result := workflow.InProgress(workflow.StreamInstanceSetupInProgress, "configuring stream instance in Atlas")
	ctx.SetConditionFromResult(api.StreamInstanceReadyType, result)
	return result.ReconcileResult(), nil
}
