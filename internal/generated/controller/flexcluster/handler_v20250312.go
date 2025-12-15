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

package flexcluster

import (
	"context"
	"errors"
	"fmt"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312009/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
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

func (h *Handlerv20250312) getDependencies(ctx context.Context, flexcluster *akov2generated.FlexCluster) ([]client.Object, error) {
	var result []client.Object

	if flexcluster.Spec.V20250312.GroupRef != nil {
		group := &akov2generated.Group{}

		err := h.kubeClient.Get(ctx, client.ObjectKey{
			Name:      flexcluster.Spec.V20250312.GroupRef.Name,
			Namespace: flexcluster.GetNamespace(),
		}, group)

		if err != nil {
			return nil, fmt.Errorf("failed to get group  %w", err)
		}

		result = append(result, group)
	}

	return result, nil
}

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to get dependencies: %w", err))
	}

	body := &v20250312sdk.FlexClusterDescriptionCreate20241113{}
	params := &v20250312sdk.CreateFlexClusterApiParams{
		FlexClusterDescriptionCreate20241113: body,
	}

	if err := h.translator.ToAPI(params, flexcluster, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate flex api params: %w", err))
	}

	if err := h.translator.ToAPI(body, flexcluster, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate flex create description: %w", err))
	}

	atlasFlexCluster, _, err := h.atlasClient.FlexClustersApi.CreateFlexClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create flex cluster: %w", err))
	}
	newFlexCluster := flexcluster.DeepCopy()
	if _, err := h.translator.FromAPI(newFlexCluster, atlasFlexCluster); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate flex create response: %w", err))
	}

	err = h.kubeClient.Status().Patch(ctx, newFlexCluster, client.MergeFrom(flexcluster))
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to patch flex cluster status: %w", err))
	}
	return result.NextState(state.StateCreating, "Creating Flex Cluster.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	externalName, ok := flexcluster.GetAnnotations()["mongodb.com/external-name"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing mongodb.com/external-name"))
	}

	externalGroupID, ok := flexcluster.GetAnnotations()["mongodb.com/external-group-id"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing mongodb.com/external-group-id"))
	}
	flexClusterCopy := flexcluster.DeepCopy()
	flexClusterCopy.Spec.V20250312.Entry.Name = externalName
	flexClusterCopy.Spec.V20250312.GroupId = &externalGroupID
	_, err := h.patchStatus(ctx, flexClusterCopy)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}
	return result.NextState(state.StateImported, "Imported Flex Cluster.")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, flexcluster, state.StateCreated, state.StateUpdating)
}

// HandleCreating handles the creating state for version v20250312
func (h *Handlerv20250312) HandleCreating(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	return h.handleUpserting(ctx, flexcluster, state.StateCreating, state.StateCreated)
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, flexcluster, state.StateCreated, state.StateUpdating)
}

// HandleUpdating handles the updating state for version v20250312
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	return h.handleUpserting(ctx, flexcluster, state.StateUpdating, state.StateUpdated)
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, flexcluster, state.StateUpdated, state.StateUpdating)
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(flexcluster, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "Flex Cluster deleted.")
	}

	if flexcluster.Status.V20250312 == nil {
		return result.NextState(state.StateDeleted, "Flex Cluster is unamanged.")
	}

	deps, err := h.getDependencies(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to get dependencies: %w", err))
	}

	params := &v20250312sdk.DeleteFlexClusterApiParams{}
	if err := h.translator.ToAPI(params, flexcluster, deps...); err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to translate flex api params: %w", err))
	}

	_, err = h.atlasClient.FlexClustersApi.DeleteFlexClusterWithParams(ctx, params).Execute()

	switch {
	case v20250312sdk.IsErrorCode(err, "CLUSTER_NOT_FOUND"):
		return result.NextState(state.StateDeleted, "Flex Cluster was deleted in Atlas.")
	case err != nil:
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete flex cluster: %w", err))
	}

	return result.NextState(state.StateDeleting, "Deleting Flex Cluster.")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateDeleting, fmt.Errorf("failed to get dependencies: %w", err))
	}

	params := &v20250312sdk.GetFlexClusterApiParams{}
	if err := h.translator.ToAPI(params, flexcluster, deps...); err != nil {
		return result.Error(state.StateDeleting, fmt.Errorf("failed to translate flex api params: %w", err))
	}

	_, _, err = h.atlasClient.FlexClustersApi.GetFlexClusterWithParams(ctx, params).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, "CLUSTER_NOT_FOUND"):
		return result.NextState(state.StateDeleted, "Deleted")
	case err != nil:
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete flexcluster: %w", err))
	}
	return result.NextState(state.StateDeleting, "Deleting Flex Cluster.")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.FlexCluster{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	panic("do not setup Handlerv20250312")
}

// HandleUpserting handles the creating and updating state for flex version v20250312
func (h *Handlerv20250312) handleUpserting(ctx context.Context, flexcluster *akov2generated.FlexCluster, currentState, finalState state.ResourceState) (ctrlstate.Result, error) {
	atlasFlexCluster, err := h.patchStatus(ctx, flexcluster)
	if err != nil {
		return result.Error(currentState, err)
	}
	if atlasFlexCluster.GetStateName() == "CREATING" || atlasFlexCluster.GetStateName() == "UPDATING" {
		return result.NextState(currentState, "Upserting Flex Cluster.")
	}
	return result.NextState(finalState, "Upserted Flex Cluster.")
}

// HandleIdle handles the creating and updating state for flex version v20250312
func (h *Handlerv20250312) handleIdle(ctx context.Context, flexcluster *akov2generated.FlexCluster, currentState, finalState state.ResourceState) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, flexcluster)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to get dependencies: %w", err))
	}

	update, err := ctrlstate.ShouldUpdate(flexcluster, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "Flex cluster up to date. No update required.")
	}

	body := &v20250312sdk.FlexClusterDescriptionUpdate20241113{}
	params := &v20250312sdk.UpdateFlexClusterApiParams{
		FlexClusterDescriptionUpdate20241113: body,
	}

	// translate parameters
	if err := h.translator.ToAPI(params, flexcluster, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate update flex cluster parameters: %w", err))
	}

	// translate body
	if err := h.translator.ToAPI(body, flexcluster, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate update flex cluster description: %w", err))
	}

	atlasFlexCluster, _, err := h.atlasClient.FlexClustersApi.UpdateFlexClusterWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to get update cluster: %w", err))
	}

	flexclusterCopy := flexcluster.DeepCopy()
	if _, err := h.translator.FromAPI(flexclusterCopy, atlasFlexCluster); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate update cluster response: %w", err))
	}

	if err := ctrlstate.
		NewPatcher(flexclusterCopy).
		UpdateStateTracker(deps...).
		UpdateStatus().
		Patch(ctx, h.kubeClient); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to patch cluster: %w", err))
	}

	return result.NextState(finalState, "Updating Flex Cluster.")
}

func (h *Handlerv20250312) patchStatus(ctx context.Context, flexcluster *akov2generated.FlexCluster) (*v20250312sdk.FlexClusterDescription20241113, error) {
	deps, err := h.getDependencies(ctx, flexcluster)
	if err != nil {
		return nil, err
	}
	params := &v20250312sdk.GetFlexClusterApiParams{}
	if err := h.translator.ToAPI(params, flexcluster, deps...); err != nil {
		return nil, fmt.Errorf("failed to translate update flex cluster parameters: %w", err)
	}

	atlasFlexCluster, _, err := h.atlasClient.FlexClustersApi.GetFlexClusterWithParams(ctx, params).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}

	flexclusterCopy := flexcluster.DeepCopy()
	if _, err := h.translator.FromAPI(flexclusterCopy, atlasFlexCluster, deps...); err != nil {
		return nil, fmt.Errorf("failed to translate get cluster response: %w", err)
	}

	if err := ctrlstate.
		NewPatcher(flexclusterCopy).
		UpdateStateTracker(deps...).
		UpdateStatus().
		Patch(ctx, h.kubeClient); err != nil {
		return nil, fmt.Errorf("failed to patch cluster: %w", err)
	}

	return atlasFlexCluster, nil
}
