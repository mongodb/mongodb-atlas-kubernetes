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

package group

import (
	"context"
	"errors"
	"fmt"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312014/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
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

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	atlasGroup := &v20250312sdk.Group{}
	err := h.translator.ToAPI(atlasGroup, group)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate group to Atlas: %w", err))
	}

	params := &v20250312sdk.CreateGroupApiParams{Group: atlasGroup, ProjectOwnerId: group.Spec.V20250312.ProjectOwnerId}
	response, _, err := h.atlasClient.ProjectsApi.CreateGroupWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create group: %w", err))
	}

	groupCopy := group.DeepCopy()
	_, err = h.translator.FromAPI(groupCopy, response)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate group from Atlas: %w", err))
	}

	if err := ctrlstate.NewPatcher(groupCopy).UpdateStatus().UpdateStateTracker().Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to patch group status: %w", err))
	}

	return result.NextState(state.StateCreated, "Group created.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	id, ok := group.GetAnnotations()["mongodb.com/external-id"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-id"))
	}

	response, _, err := h.atlasClient.ProjectsApi.GetGroup(ctx, id).Execute()
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to get Group with id %s: %w", id, err))
	}

	groupCopy := group.DeepCopy()
	_, err = h.translator.FromAPI(groupCopy, response)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate Group from Atlas: %w", err))
	}

	if err := ctrlstate.NewPatcher(groupCopy).UpdateStatus().UpdateStateTracker().Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to patch Group status: %w", err))
	}

	return result.NextState(state.StateImported, "Group imported.")
}

func (h *Handlerv20250312) HandleImported(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateImported, group)
}

func (h *Handlerv20250312) HandleCreating(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	panic("unsupported state")
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateCreated, group)
}

// HandleUpdating handles the updating state for version v20250312
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	panic("unsupported state")
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return h.handleUpserted(ctx, state.StateUpdated, group)
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	dependents := h.getDependents(ctx, group)
	if len(dependents) > 0 {
		return result.NextState(state.StateDeletionRequested, fmt.Sprintf("failed to delete group because %v resources depend on it.", len(dependents)))
	}

	if customresource.IsResourcePolicyKeepOrDefault(group, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "Group deleted.")
	}

	if group.Status.V20250312 == nil || group.Status.V20250312.Id == nil {
		return result.NextState(state.StateDeleted, "Group deleted.")
	}

	_, err := h.atlasClient.ProjectsApi.DeleteGroup(ctx, *group.Status.V20250312.Id).Execute()
	if v20250312sdk.IsErrorCode(err, "GROUP_NOT_FOUND") {
		return result.NextState(state.StateDeleted, "Group deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete group: %w", err))
	}

	return result.NextState(state.StateDeleting, "Deleting group.")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	_, _, err := h.atlasClient.ProjectsApi.GetGroup(ctx, *group.Status.V20250312.Id).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, "GROUP_NOT_FOUND"):
		return result.NextState(state.StateDeleted, "Deleted")
	case err != nil:
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete group: %w", err))
	}
	return result.NextState(state.StateDeleting, "Deleting Group.")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.Group{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}

func (h *Handlerv20250312) handleUpserted(ctx context.Context, currentState state.ResourceState, group *akov2generated.Group) (ctrlstate.Result, error) {
	update, err := ctrlstate.ShouldUpdate(group)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "Group is up to date. No update required.")
	}

	atlasGroupUpdate := &v20250312sdk.GroupUpdate{}
	err = h.translator.ToAPI(atlasGroupUpdate, group)
	if err != nil {
		return result.Error(currentState, err)
	}

	params := &v20250312sdk.UpdateGroupApiParams{
		GroupId:     *group.Status.V20250312.Id,
		GroupUpdate: atlasGroupUpdate,
	}

	response, _, err := h.atlasClient.ProjectsApi.UpdateGroupWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(currentState, err)
	}

	groupCopy := group.DeepCopy()
	_, err = h.translator.FromAPI(groupCopy, response)
	if err != nil {
		return result.Error(currentState, err)
	}

	if err := ctrlstate.NewPatcher(groupCopy).UpdateStateTracker().UpdateStatus().Patch(ctx, h.kubeClient); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to patch group: %w", err))
	}

	return result.NextState(state.StateUpdated, "Group is updated.")
}
