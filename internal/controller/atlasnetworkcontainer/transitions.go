// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasnetworkcontainer

import (
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

func (r *AtlasNetworkContainerReconciler) create(workflowCtx *workflow.Context, req *reconcileRequest) (ctrl.Result, error) {
	cfg := networkcontainer.NewNetworkContainerConfig(
		req.networkContainer.Spec.Provider,
		&req.networkContainer.Spec.AtlasNetworkContainerConfig,
	)
	createdContainer, err := req.service.Create(workflowCtx.Context, req.projectID, cfg)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to create container: %w", err)
		return r.terminate(workflowCtx, req.networkContainer, workflow.NetworkContainerNotConfigured, wrappedErr), nil
	}
	return r.ready(workflowCtx, req.networkContainer, createdContainer)
}

func (r *AtlasNetworkContainerReconciler) sync(workflowCtx *workflow.Context, req *reconcileRequest, atlasContainer *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	desiredConfig := networkcontainer.NewNetworkContainerConfig(
		req.networkContainer.Spec.Provider, &req.networkContainer.Spec.AtlasNetworkContainerConfig)
	// only the CIDR block can be updated in a container
	if desiredConfig.CIDRBlock != atlasContainer.NetworkContainerConfig.CIDRBlock {
		return r.update(workflowCtx, req, atlasContainer.ID, desiredConfig)
	}
	return r.ready(workflowCtx, req.networkContainer, atlasContainer)
}

func (r *AtlasNetworkContainerReconciler) update(workflowCtx *workflow.Context, req *reconcileRequest, id string, config *networkcontainer.NetworkContainerConfig) (ctrl.Result, error) {
	updatedContainer, err := req.service.Update(workflowCtx.Context, req.projectID, id, config)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update container: %w", err)
		return r.terminate(workflowCtx, req.networkContainer, workflow.NetworkContainerNotConfigured, wrappedErr), nil
	}
	return r.ready(workflowCtx, req.networkContainer, updatedContainer)
}

func (r *AtlasNetworkContainerReconciler) delete(workflowCtx *workflow.Context, req *reconcileRequest, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(req.networkContainer, r.ObjectDeletionProtection) {
		return r.unmanage(workflowCtx, req.networkContainer)
	}
	err := req.service.Delete(workflowCtx.Context, req.projectID, container.ID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete container: %w", err)
		return r.terminate(workflowCtx, req.networkContainer, workflow.NetworkContainerNotDeleted, wrappedErr), nil
	}
	return r.unmanage(workflowCtx, req.networkContainer)
}

func (r *AtlasNetworkContainerReconciler) ready(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer, container *networkcontainer.NetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.SetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotSet, err), nil
	}

	workflowCtx.SetConditionTrueMsg(api.NetworkContainerReady, fmt.Sprintf("Network Container %s is ready", container.ID)).
		SetConditionTrue(api.ReadyType).EnsureStatusOption(updateNetworkContainerStatusOption(container))

	if networkContainer.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult(), nil
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasNetworkContainerReconciler) unmanage(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.UnsetFinalizer); err != nil {
		return r.terminate(workflowCtx, networkContainer, workflow.AtlasFinalizerNotRemoved, err), nil
	}
	return workflow.Deleted().ReconcileResult(), nil
}

func (r *AtlasNetworkContainerReconciler) release(workflowCtx *workflow.Context, networkContainer *akov2.AtlasNetworkContainer, err error) ctrl.Result {
	if errors.Is(err, reconciler.ErrMissingKubeProject) {
		if finalizerErr := customresource.ManageFinalizer(workflowCtx.Context, r.Client, networkContainer, customresource.UnsetFinalizer); finalizerErr != nil {
			err = errors.Join(err, finalizerErr)
		}
	}
	return r.terminate(workflowCtx, networkContainer, workflow.NetworkContainerNotConfigured, err)
}

func (r *AtlasNetworkContainerReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) ctrl.Result {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).SetConditionFromResult(condition, result)

	return result.ReconcileResult()
}

func updateNetworkContainerStatusOption(container *networkcontainer.NetworkContainer) status.AtlasNetworkContainerStatusOption {
	return func(containerStatus *status.AtlasNetworkContainerStatus) {
		networkcontainer.ApplyNetworkContainerStatus(containerStatus, container)
	}
}
