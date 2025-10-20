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

package groupalertsconfig

import (
	"context"
	"fmt"

	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set
func (h *GroupAlertsConfigHandler) getHandlerForResource(groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.StateHandler[v1.GroupAlertsConfig], error) {
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.GroupAlertsConfig]

	if groupalertsconfig.Spec.V20250312 != nil {
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
func (h *GroupAlertsConfigHandler) HandleInitial(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, groupalertsconfig)
}

// HandleImportRequested delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleImportRequested(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImportRequested(ctx, groupalertsconfig)
}

// HandleImported delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleImported(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImported(ctx, groupalertsconfig)
}

// HandleCreating delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleCreating(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreating(ctx, groupalertsconfig)
}

// HandleCreated delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleCreated(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreated(ctx, groupalertsconfig)
}

// HandleUpdating delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleUpdating(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdating(ctx, groupalertsconfig)
}

// HandleUpdated delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleUpdated(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdated(ctx, groupalertsconfig)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleDeletionRequested(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeletionRequested(ctx, groupalertsconfig)
}

// HandleDeleting delegates to the version-specific handler
func (h *GroupAlertsConfigHandler) HandleDeleting(ctx context.Context, groupalertsconfig *v1.GroupAlertsConfig) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(groupalertsconfig)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeleting(ctx, groupalertsconfig)
}

// For returns the resource and predicates for the controller
func (h *GroupAlertsConfigHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.GroupAlertsConfig{}
	// TODO: Add appropriate predicates
	return obj, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *GroupAlertsConfigHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("GroupAlertsConfig").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
