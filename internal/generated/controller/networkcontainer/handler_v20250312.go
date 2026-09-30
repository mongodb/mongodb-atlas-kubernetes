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

package networkcontainer

import (
	"context"
	"errors"
	"fmt"

	ctrlstate "github.com/crd2go/constate"
	state "github.com/crd2go/constate/state"
	crapi "github.com/crd2go/crapi"
	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312026/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
)

const (
	atlasContainerNotFound = "CLOUD_PROVIDER_CONTAINER_NOT_FOUND"
	atlasContainersInUse   = "CONTAINERS_IN_USE"
)

// Handlerv20250312 is the handler for the v20250312 version of the networkcontainer resource
type Handlerv20250312 struct {
	kubeClient         client.Client
	atlasClient        *v20250312sdk.APIClient
	translator         crapi.Translator
	deletionProtection bool
}

// NewHandlerv20250312 creates a new Handlerv20250312 handler
func NewHandlerv20250312(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translator crapi.Translator, deletionProtection bool) *Handlerv20250312 {
	return &Handlerv20250312{
		atlasClient:        atlasClient,
		deletionProtection: deletionProtection,
		kubeClient:         kubeClient,
		translator:         translator,
	}
}

// HandleInitial creates the network container in Atlas.
func (h *Handlerv20250312) HandleInitial(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, networkcontainer)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to resolve NetworkContainer dependencies: %w", err))
	}

	params := &v20250312sdk.CreateGroupContainerApiParams{CloudProviderContainer: &v20250312sdk.CloudProviderContainer{}}
	if err := h.translator.ToAPI(params, networkcontainer, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate NetworkContainer API parameters to Atlas: %w", err))
	}
	if err := h.translator.ToAPI(params.CloudProviderContainer, networkcontainer, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate NetworkContainer to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.NetworkPeeringAPI.CreateGroupContainerWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create NetworkContainer: %w", err))
	}

	if err := h.patchStatus(ctx, networkcontainer, response, deps...); err != nil {
		return result.Error(state.StateInitial, err)
	}

	return result.NextState(state.StateCreated, "NetworkContainer created.")
}

// HandleImportRequested handles the importrequested state for version v20250312.
// The annotation mongodb.com/external-id must be set to the Atlas container ID.
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	containerID, err := ctrlstate.GetExternalID(networkcontainer)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}

	deps, err := h.getDependencies(ctx, networkcontainer)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to resolve NetworkContainer dependencies: %w", err))
	}

	groupID, err := h.resolveGroupID(networkcontainer, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}

	response, _, err := h.atlasClient.NetworkPeeringAPI.GetGroupContainer(ctx, groupID, containerID).Execute()
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to get NetworkContainer %q in group %q: %w", containerID, groupID, err))
	}

	if err := h.patchStatus(ctx, networkcontainer, response); err != nil {
		return result.Error(state.StateImportRequested, err)
	}

	return result.NextState(state.StateImported, "NetworkContainer imported.")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateImported, networkcontainer)
}

// HandleCreating is not expected to be reached: Atlas creates containers synchronously.
func (h *Handlerv20250312) HandleCreating(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "NetworkContainer created.")
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateCreated, networkcontainer)
}

// HandleUpdating is not expected to be reached: Atlas updates containers synchronously.
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "NetworkContainer updated.")
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateUpdated, networkcontainer)
}

// HandleDeletionRequested deletes the network container from Atlas.
// Atlas refuses to delete a container that still has clusters or peering connections;
// that error is returned so the deletion is retried.
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(networkcontainer, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "NetworkContainer skipped deletion due to retention policy.")
	}

	containerID := statusContainerID(networkcontainer)
	if containerID == "" {
		return result.NextState(state.StateDeleted, "NetworkContainer deleted.")
	}

	deps, err := h.getDependencies(ctx, networkcontainer)
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to resolve NetworkContainer dependencies: %w", err))
	}

	groupID, err := h.resolveGroupID(networkcontainer, deps...)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}

	_, err = h.atlasClient.NetworkPeeringAPI.DeleteGroupContainer(ctx, groupID, containerID).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, atlasContainerNotFound):
		return result.NextState(state.StateDeleted, "NetworkContainer deleted.")
	case v20250312sdk.IsErrorCode(err, atlasContainersInUse):
		return result.Error(state.StateDeletionRequested, fmt.Errorf("NetworkContainer %q is still in use by clusters or peering connections: %w", containerID, err))
	case err != nil:
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete NetworkContainer %q: %w", containerID, err))
	}

	return result.NextState(state.StateDeleted, "NetworkContainer deleted.")
}

// HandleDeleting polls Atlas until the network container is gone.
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	containerID := statusContainerID(networkcontainer)
	if containerID == "" {
		return result.NextState(state.StateDeleted, "NetworkContainer deleted.")
	}

	deps, err := h.getDependencies(ctx, networkcontainer)
	if err != nil {
		return result.Error(state.StateDeleting, fmt.Errorf("failed to resolve NetworkContainer dependencies: %w", err))
	}

	groupID, err := h.resolveGroupID(networkcontainer, deps...)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}

	_, _, err = h.atlasClient.NetworkPeeringAPI.GetGroupContainer(ctx, groupID, containerID).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, atlasContainerNotFound):
		return result.NextState(state.StateDeleted, "NetworkContainer deleted.")
	case err != nil:
		return result.Error(state.StateDeleting, fmt.Errorf("failed to check NetworkContainer deletion: %w", err))
	}

	return result.NextState(state.StateDeleting, "Waiting for NetworkContainer deletion.")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.NetworkContainer{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}

// handleUpserted updates the network container in Atlas when the spec or its dependencies changed.
func (h *Handlerv20250312) handleUpserted(ctx context.Context, currentState state.ResourceState, networkcontainer *akov2generated.NetworkContainer) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, networkcontainer)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to resolve NetworkContainer dependencies: %w", err))
	}

	update, err := ctrlstate.ShouldUpdate(networkcontainer, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "NetworkContainer is up to date. No update required.")
	}

	containerID := statusContainerID(networkcontainer)
	if containerID == "" {
		return result.Error(currentState, errors.New("NetworkContainer status has no Atlas container ID"))
	}

	params := &v20250312sdk.UpdateGroupContainerApiParams{CloudProviderContainer: &v20250312sdk.CloudProviderContainer{}}
	if err := h.translator.ToAPI(params, networkcontainer, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate NetworkContainer API parameters to Atlas: %w", err))
	}
	if err := h.translator.ToAPI(params.CloudProviderContainer, networkcontainer, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate NetworkContainer to Atlas: %w", err))
	}
	params.ContainerId = containerID

	response, _, err := h.atlasClient.NetworkPeeringAPI.UpdateGroupContainerWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update NetworkContainer %q: %w", containerID, err))
	}

	if err := h.patchStatus(ctx, networkcontainer, response, deps...); err != nil {
		return result.Error(currentState, err)
	}

	return result.NextState(state.StateUpdated, "NetworkContainer updated.")
}

// resolveGroupID collapses groupId or groupRef into the Atlas project ID.
func (h *Handlerv20250312) resolveGroupID(networkcontainer *akov2generated.NetworkContainer, deps ...client.Object) (string, error) {
	params := &v20250312sdk.GetGroupContainerApiParams{}
	if err := h.translator.ToAPI(params, networkcontainer, deps...); err != nil {
		return "", fmt.Errorf("failed to resolve NetworkContainer groupId: %w", err)
	}
	if params.GroupId == "" {
		return "", errors.New("could not resolve groupId: neither groupId nor groupRef with a ready Group is available")
	}

	return params.GroupId, nil
}

// patchStatus copies the Atlas response into the status and refreshes the state tracker.
func (h *Handlerv20250312) patchStatus(ctx context.Context, networkcontainer *akov2generated.NetworkContainer, response *v20250312sdk.CloudProviderContainer, deps ...client.Object) error {
	networkcontainerCopy := networkcontainer.DeepCopy()
	if _, err := h.translator.FromAPI(networkcontainerCopy, response); err != nil {
		return fmt.Errorf("failed to translate NetworkContainer from Atlas: %w", err)
	}

	if err := ctrlstate.NewPatcher(networkcontainerCopy).UpdateStatus().UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return fmt.Errorf("failed to patch NetworkContainer status: %w", err)
	}

	return nil
}

// statusContainerID returns the Atlas container ID recorded in the status, or "" if not yet known.
func statusContainerID(networkcontainer *akov2generated.NetworkContainer) string {
	if networkcontainer.Status.V20250312 == nil || networkcontainer.Status.V20250312.Id == nil {
		return ""
	}

	return *networkcontainer.Status.V20250312.Id
}
