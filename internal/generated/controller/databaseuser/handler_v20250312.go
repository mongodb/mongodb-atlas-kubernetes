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

package databaseuser

import (
	"context"
	"errors"
	"fmt"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312011/admin"
	v1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

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

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("Failed to get dependencies: %w", err))
	}

	body := &v20250312sdk.CloudDatabaseUser{}
	params := &v20250312sdk.CreateDatabaseUserApiParams{
		CloudDatabaseUser: body,
	}

	if err := h.translator.ToAPI(params, databaseuser, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate params: %w", err))
	}

	if err := h.translator.ToAPI(body, databaseuser, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate body: %w", err))
	}

	_, _, err = h.atlasClient.DatabaseUsersApi.CreateDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create datebaseuser: %w", err))
	}

	if err := ctrlstate.NewPatcher(databaseuser).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateCreated, fmt.Errorf("failed to update state tracker: %w", err))
	}

	return result.NextState(state.StateCreated, "Created Database User.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("Failed to get dependencies: %w", err))
	}

	groupId, ok := databaseuser.GetAnnotations()["mongodb.com/external-group-id"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-id"))
	}

	databaseName, ok := databaseuser.GetAnnotations()["mongodb.com/external-database-name"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-database-name"))
	}

	username, ok := databaseuser.GetAnnotations()["mongodb.com/external-username"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-username"))
	}

	params := &v20250312sdk.GetDatabaseUserApiParams{
		GroupId:      groupId,
		DatabaseName: databaseName,
		Username:     username,
	}
	_, _, err = h.atlasClient.DatabaseUsersApi.GetDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create datebaseuser: %w", err))
	}

	if err := ctrlstate.NewPatcher(databaseuser).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateImported, fmt.Errorf("failed to update state tracker: %w", err))
	}

	return result.NextState(state.StateImported, "Import completed")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, state.StateImported, databaseuser)
}

// HandleCreating handles the creating state for version v20250312
func (h *Handlerv20250312) HandleCreating(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	panic("unsupported state")
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, state.StateCreated, databaseuser)
}

// HandleUpdating handles the updating state for version v20250312
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	panic("unsupported state")
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	return h.handleIdle(ctx, state.StateUpdated, databaseuser)
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("Failed to get dependencies: %w", err))
	}

	params := &v20250312sdk.DeleteDatabaseUserApiParams{}
	if err := h.translator.ToAPI(params, databaseuser, deps...); err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to translate params: %w", err))
	}
	_, err = h.atlasClient.DatabaseUsersApi.DeleteDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete datebaseuser: %w", err))
	}

	return result.NextState(state.StateDeleted, "User deleted.")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	panic("unsupported state")
}

func (h *Handlerv20250312) handleIdle(ctx context.Context, currentState state.ResourceState, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("Failed to get dependencies: %w", err))
	}

	update, err := ctrlstate.ShouldUpdate(databaseuser, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "Database user up to date. No update required.")
	}

	body := &v20250312sdk.CloudDatabaseUser{}
	params := &v20250312sdk.UpdateDatabaseUserApiParams{
		CloudDatabaseUser: body,
	}

	if err := h.translator.ToAPI(params, databaseuser, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate params: %w", err))
	}

	if err := h.translator.ToAPI(body, databaseuser, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate body: %w", err))
	}

	_, _, err = h.atlasClient.DatabaseUsersApi.UpdateDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update datebaseuser: %w", err))
	}

	if err := ctrlstate.NewPatcher(databaseuser).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update state tracker: %w", err))
	}

	return result.NextState(state.StateUpdated, "Updated Database User.")
}

func (h *Handlerv20250312) getDependencies(ctx context.Context, databaseuser *akov2generated.DatabaseUser) ([]client.Object, error) {
	var deps []client.Object

	// Check if passwordSecretRef is present
	if databaseuser.Spec.V20250312 != nil && databaseuser.Spec.V20250312.Entry != nil && databaseuser.Spec.V20250312.Entry.PasswordSecretRef != nil {
		secret := &v1.Secret{}
		err := h.kubeClient.Get(ctx, client.ObjectKey{
			Name:      databaseuser.Spec.V20250312.Entry.PasswordSecretRef.Name,
			Namespace: databaseuser.GetNamespace(),
		}, secret)
		if err != nil {
			return deps, fmt.Errorf("failed to get Secret %s/%s: %w", databaseuser.GetNamespace(), databaseuser.Spec.V20250312.Entry.PasswordSecretRef.Name, err)
		}

		deps = append(deps, secret)
	}

	// Check if groupRef is present
	if databaseuser.Spec.V20250312 != nil && databaseuser.Spec.V20250312.GroupRef != nil {
		group := &akov2generated.Group{}
		err := h.kubeClient.Get(ctx, client.ObjectKey{
			Name:      databaseuser.Spec.V20250312.GroupRef.Name,
			Namespace: databaseuser.GetNamespace(),
		}, group)
		if err != nil {
			return deps, fmt.Errorf("failed to get Group %s/%s: %w", databaseuser.GetNamespace(), databaseuser.Spec.V20250312.GroupRef.Name, err)
		}

		deps = append(deps, group)
	}

	return deps, nil
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.DatabaseUser{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}
