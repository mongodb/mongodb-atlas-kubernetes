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

package team

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
func (h *TeamHandler) getHandlerForResource(team *v1.Team) (ctrlstate.StateHandler[v1.Team], error) {
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.Team]

	if team.Spec.V20250312 != nil {
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
func (h *TeamHandler) HandleInitial(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, team)
}

// HandleImportRequested delegates to the version-specific handler
func (h *TeamHandler) HandleImportRequested(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImportRequested(ctx, team)
}

// HandleImported delegates to the version-specific handler
func (h *TeamHandler) HandleImported(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImported(ctx, team)
}

// HandleCreating delegates to the version-specific handler
func (h *TeamHandler) HandleCreating(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreating(ctx, team)
}

// HandleCreated delegates to the version-specific handler
func (h *TeamHandler) HandleCreated(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreated(ctx, team)
}

// HandleUpdating delegates to the version-specific handler
func (h *TeamHandler) HandleUpdating(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdating(ctx, team)
}

// HandleUpdated delegates to the version-specific handler
func (h *TeamHandler) HandleUpdated(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdated(ctx, team)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *TeamHandler) HandleDeletionRequested(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeletionRequested(ctx, team)
}

// HandleDeleting delegates to the version-specific handler
func (h *TeamHandler) HandleDeleting(ctx context.Context, team *v1.Team) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(team)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeleting(ctx, team)
}

// For returns the resource and predicates for the controller
func (h *TeamHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.Team{}
	// TODO: Add appropriate predicates
	return obj, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *TeamHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("Team").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
