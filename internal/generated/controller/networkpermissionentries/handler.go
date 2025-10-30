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

package networkpermissionentries

import (
	"context"
	"fmt"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set
func (h *NetworkPermissionEntriesHandler) getHandlerForResource(networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.StateHandler[v1.NetworkPermissionEntries], error) {
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.NetworkPermissionEntries]

	if networkpermissionentries.Spec.V20250312 != nil {
		versionCount++
		selectedHandler = h.handlerv20250312
	}

	if versionCount == 0 {
		return nil, fmt.Errorf("no resource spec version specified - please set one of the available spec versions")
	}
	if versionCount > 1 {
		return nil, fmt.Errorf("multiple resource spec versions specified - please set only one spec version")
	}
	return selectedHandler, nil
}

// HandleInitial delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleInitial(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, networkpermissionentries)
}

// HandleImportRequested delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleImportRequested(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImportRequested(ctx, networkpermissionentries)
}

// HandleImported delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleImported(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImported(ctx, networkpermissionentries)
}

// HandleCreating delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleCreating(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreating(ctx, networkpermissionentries)
}

// HandleCreated delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleCreated(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreated(ctx, networkpermissionentries)
}

// HandleUpdating delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleUpdating(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdating(ctx, networkpermissionentries)
}

// HandleUpdated delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleUpdated(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdated(ctx, networkpermissionentries)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleDeletionRequested(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeletionRequested(ctx, networkpermissionentries)
}

// HandleDeleting delegates to the version-specific handler
func (h *NetworkPermissionEntriesHandler) HandleDeleting(ctx context.Context, networkpermissionentries *v1.NetworkPermissionEntries) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(networkpermissionentries)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeleting(ctx, networkpermissionentries)
}

// For returns the resource and predicates for the controller
func (h *NetworkPermissionEntriesHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.NetworkPermissionEntries{}
	// TODO: Add appropriate predicates
	return obj, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *NetworkPermissionEntriesHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("NetworkPermissionEntries").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
