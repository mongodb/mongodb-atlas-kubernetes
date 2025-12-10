// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

import (
	"context"
	"errors"
	"fmt"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312010/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	StateCreating  = "CREATING"
	StateUpdating  = "UPDATING"
	StateDeleting  = "DELETING"
	StateRepairing = "REPAIRING"

	SharedImmutableError = "TENANT_CLUSTER_UPDATE_UNSUPPORTED"
)

type Handlerv20250312 struct {
	kubeClient         client.Client
	atlasClient        *v20250312sdk.APIClient
	translator         crapi.Translator
	deletionProtection bool
}

func NewHandlerv20250312(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translator crapi.Translator, deletionProtection bool) *Handlerv20250312 {
	return &Handlerv20250312{
		atlasClient:        atlasClient,
		deletionProtection: deletionProtection,
		kubeClient:         kubeClient,
		translator:         translator,
	}
}

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	body := &v20250312sdk.ClusterDescription20240805{}
	params := &v20250312sdk.CreateClusterApiParams{
		ClusterDescription20240805: body,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	err = h.translator.ToAPI(body, cluster, deps...)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate cluster API body to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.ClustersApi.CreateClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create cluster: %w", err))
	}

	err = h.patchStatus(ctx, cluster, response, deps...)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}

	return result.NextState(state.StateCreating, "Cluster is being created.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	id, ok := cluster.GetAnnotations()["mongodb.com/external-id"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-id"))
	}

	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.GetClusterApiParams{
		ClusterName:                *cluster.Spec.V20250312.Entry.Name,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.ClustersApi.GetClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to get Cluster with id %s: %w", id, err))
	}

	err = h.patchStatus(ctx, cluster, response, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to patch Cluster status: %w", err))
	}

	return result.NextState(state.StateImported, "Cluster is being imported.")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateImported, cluster)
}

// HandleCreating handles the creating state for version v20250312
func (h *Handlerv20250312) HandleCreating(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.GetClusterApiParams{
		ClusterName:                *cluster.Spec.V20250312.Entry.Name,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	atlasCluster, _, err := h.atlasClient.ClustersApi.GetClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to get Cluster with name %s: %w", params.ClusterName, err))
	}

	return h.handleAtlasClusterState(ctx, cluster, atlasCluster, state.StateCreating, state.StateCreated, deps...)
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateCreated, cluster)
}

// HandleUpdating handles the updating state for version v20250312
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.GetClusterApiParams{
		ClusterName:                *cluster.Spec.V20250312.Entry.Name,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	atlasCluster, _, err := h.atlasClient.ClustersApi.GetClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to get Cluster with name %s: %w", params.ClusterName, err))
	}

	return h.handleAtlasClusterState(ctx, cluster, atlasCluster, state.StateUpdating, state.StateUpdated, deps...)
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateUpdated, cluster)
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(cluster, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "Cluster deleted.")
	}

	if cluster.Status.V20250312 == nil || cluster.Status.V20250312.Id == nil {
		return result.NextState(state.StateDeleted, "Cluster deleted.")
	}

	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.DeleteClusterApiParams{
		ClusterName: *cluster.Spec.V20250312.Entry.Name,
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	_, err = h.atlasClient.ClustersApi.DeleteClusterWithParams(ctx, params).Execute()
	if v20250312sdk.IsErrorCode(err, "CLUSTER_NOT_FOUND") {
		return result.NextState(state.StateDeleted, "Cluster deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete Cluster: %w", err))
	}

	return result.NextState(state.StateDeleting, "Deleting Cluster.")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(state.StateDeleting, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.GetClusterApiParams{
		ClusterName:                *cluster.Spec.V20250312.Entry.Name,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	atlasCluster, _, err := h.atlasClient.ClustersApi.GetClusterWithParams(ctx, params).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, "CLUSTER_NOT_FOUND"):
		return result.NextState(state.StateDeleted, "Deleted")
	case err != nil:
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete Cluster: %w", err))
	}

	return h.handleAtlasClusterState(ctx, cluster, atlasCluster, state.StateDeleting, state.StateDeleted, deps...)
}

func (h *Handlerv20250312) handleUpserted(ctx context.Context, currentState state.ResourceState, cluster *akov2generated.Cluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, cluster)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	update, err := ctrlstate.ShouldUpdate(cluster, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "Cluster is up to date. No update required.")
	}

	body := &v20250312sdk.ClusterDescription20240805{}
	params := &v20250312sdk.UpdateClusterApiParams{
		ClusterName:                *cluster.Spec.V20250312.Entry.Name,
		UseEffectiveInstanceFields: pointer.MakePtr(true),
		ClusterDescription20240805: body,
	}
	err = h.translator.ToAPI(params, cluster, deps...)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	err = h.translator.ToAPI(body, cluster, deps...)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate cluster API body to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.ClustersApi.UpdateClusterWithParams(ctx, params).Execute()
	if err != nil {
		if v20250312sdk.IsErrorCode(err, SharedImmutableError) {
			return result.NextState(currentState, "Shared Cluster is immutable. No update performed.")
		}

		return result.Error(currentState, err)
	}

	err = h.patchStatus(ctx, cluster, response, deps...)
	if err != nil {
		return result.Error(currentState, err)
	}

	return result.NextState(state.StateUpdating, "Cluster is being updated.")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.Cluster{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(_ controllerruntime.Manager, _ reconcile.Reconciler, _ controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}

func (h *Handlerv20250312) patchStatus(ctx context.Context, cluster *akov2generated.Cluster, atlasCluster *v20250312sdk.ClusterDescription20240805, deps ...client.Object) error {
	clusterCopy := cluster.DeepCopy()
	_, err := h.translator.FromAPI(clusterCopy, atlasCluster)
	if err != nil {
		return fmt.Errorf("failed to translate Cluster from Atlas: %w", err)
	}

	return ctrlstate.NewPatcher(clusterCopy).
		UpdateStateTracker(deps...).
		UpdateStatus().
		Patch(ctx, h.kubeClient)
}

func (h *Handlerv20250312) handleAtlasClusterState(ctx context.Context, cluster *akov2generated.Cluster, atlasCluster *v20250312sdk.ClusterDescription20240805, currentState, nextState state.ResourceState, deps ...client.Object) (ctrlstate.Result, error) {
	err := h.patchStatus(ctx, cluster, atlasCluster, deps...)
	if err != nil {
		return result.Error(currentState, err)
	}

	switch atlasCluster.GetStateName() {
	case StateCreating:
		return result.NextState(state.StateCreating, "Cluster is being created.")
	case StateUpdating, StateRepairing:
		return result.NextState(state.StateUpdating, "Cluster is being updated.")
	case StateDeleting:
		return result.NextState(state.StateDeleting, "Cluster is being deleted.")
	}

	return result.NextState(nextState, "Cluster is ready.")
}

func (h *Handlerv20250312) getDependencies(ctx context.Context, cluster *akov2generated.Cluster) ([]client.Object, error) {
	var deps []client.Object

	if cluster.Spec.V20250312.GroupRef == nil {
		return deps, nil
	}

	groupRef := cluster.Spec.V20250312.GroupRef
	group := &akov2generated.Group{}
	err := h.kubeClient.Get(ctx, client.ObjectKey{Name: groupRef.Name, Namespace: cluster.GetNamespace()}, group)
	if err != nil {
		return deps, fmt.Errorf("failed to get Group %s/%s: %w", cluster.GetNamespace(), groupRef.Name, err)
	}

	deps = append(deps, group)

	return deps, nil
}
