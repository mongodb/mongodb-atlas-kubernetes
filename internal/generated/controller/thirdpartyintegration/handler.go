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

package thirdpartyintegration

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
func (h *ThirdPartyIntegrationHandler) getHandlerForResource(thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.StateHandler[v1.ThirdPartyIntegration], error) {
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.ThirdPartyIntegration]

	if thirdpartyintegration.Spec.V20250312 != nil {
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
func (h *ThirdPartyIntegrationHandler) HandleInitial(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, thirdpartyintegration)
}

// HandleImportRequested delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleImportRequested(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImportRequested(ctx, thirdpartyintegration)
}

// HandleImported delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleImported(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImported(ctx, thirdpartyintegration)
}

// HandleCreating delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleCreating(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreating(ctx, thirdpartyintegration)
}

// HandleCreated delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleCreated(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreated(ctx, thirdpartyintegration)
}

// HandleUpdating delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleUpdating(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdating(ctx, thirdpartyintegration)
}

// HandleUpdated delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleUpdated(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdated(ctx, thirdpartyintegration)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleDeletionRequested(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeletionRequested(ctx, thirdpartyintegration)
}

// HandleDeleting delegates to the version-specific handler
func (h *ThirdPartyIntegrationHandler) HandleDeleting(ctx context.Context, thirdpartyintegration *v1.ThirdPartyIntegration) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(thirdpartyintegration)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeleting(ctx, thirdpartyintegration)
}

// For returns the resource and predicates for the controller
func (h *ThirdPartyIntegrationHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.ThirdPartyIntegration{}
	// TODO: Add appropriate predicates
	return obj, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *ThirdPartyIntegrationHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("ThirdPartyIntegration").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
