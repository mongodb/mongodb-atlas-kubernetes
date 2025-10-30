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

package organization

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
func (h *OrganizationHandler) getHandlerForResource(organization *v1.Organization) (ctrlstate.StateHandler[v1.Organization], error) {
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.Organization]

	if organization.Spec.V20250312 != nil {
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
func (h *OrganizationHandler) HandleInitial(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, organization)
}

// HandleImportRequested delegates to the version-specific handler
func (h *OrganizationHandler) HandleImportRequested(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImportRequested(ctx, organization)
}

// HandleImported delegates to the version-specific handler
func (h *OrganizationHandler) HandleImported(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleImported(ctx, organization)
}

// HandleCreating delegates to the version-specific handler
func (h *OrganizationHandler) HandleCreating(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreating(ctx, organization)
}

// HandleCreated delegates to the version-specific handler
func (h *OrganizationHandler) HandleCreated(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleCreated(ctx, organization)
}

// HandleUpdating delegates to the version-specific handler
func (h *OrganizationHandler) HandleUpdating(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdating(ctx, organization)
}

// HandleUpdated delegates to the version-specific handler
func (h *OrganizationHandler) HandleUpdated(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleUpdated(ctx, organization)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *OrganizationHandler) HandleDeletionRequested(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeletionRequested(ctx, organization)
}

// HandleDeleting delegates to the version-specific handler
func (h *OrganizationHandler) HandleDeleting(ctx context.Context, organization *v1.Organization) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(organization)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleDeleting(ctx, organization)
}

// For returns the resource and predicates for the controller
func (h *OrganizationHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.Organization{}
	// TODO: Add appropriate predicates
	return obj, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *OrganizationHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("Organization").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
