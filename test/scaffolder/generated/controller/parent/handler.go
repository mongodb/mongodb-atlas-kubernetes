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

package parent

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	handler "sigs.k8s.io/controller-runtime/pkg/handler"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	indexers "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/indexers"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/scaffolder/generated/types/v1"
)

// getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set
func (h *Handler) getHandlerForResource(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.StateHandler[akov2generated.Parent], error) {
	atlasClients, err := h.getSDKClientSet(ctx, parent)
	if err != nil {
		return nil, err
	}
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[akov2generated.Parent]

	if parent.Spec.Integrations != nil {
		translator, ok := h.translators["integrations"]
		if ok != true {
			return nil, errors.New("unsupported version integrations set in CR")
		}
		versionCount++
		selectedHandler = h.handlerintegrations(h.Client, atlasClients.SdkClient20250312013, translator, h.deletionProtection)
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
func (h *Handler) HandleInitial(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, parent)
}

// HandleImportRequested delegates to the version-specific handler
func (h *Handler) HandleImportRequested(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}
	return handler.HandleImportRequested(ctx, parent)
}

// HandleImported delegates to the version-specific handler
func (h *Handler) HandleImported(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateImported, err)
	}
	return handler.HandleImported(ctx, parent)
}

// HandleCreating delegates to the version-specific handler
func (h *Handler) HandleCreating(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateCreating, err)
	}
	return handler.HandleCreating(ctx, parent)
}

// HandleCreated delegates to the version-specific handler
func (h *Handler) HandleCreated(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateCreated, err)
	}
	return handler.HandleCreated(ctx, parent)
}

// HandleUpdating delegates to the version-specific handler
func (h *Handler) HandleUpdating(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateUpdating, err)
	}
	return handler.HandleUpdating(ctx, parent)
}

// HandleUpdated delegates to the version-specific handler
func (h *Handler) HandleUpdated(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateUpdated, err)
	}
	return handler.HandleUpdated(ctx, parent)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *Handler) HandleDeletionRequested(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}
	return handler.HandleDeletionRequested(ctx, parent)
}

// HandleDeleting delegates to the version-specific handler
func (h *Handler) HandleDeleting(ctx context.Context, parent *akov2generated.Parent) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, parent)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}
	return handler.HandleDeleting(ctx, parent)
}

// For returns the resource and predicates for the controller
func (h *Handler) For() (client.Object, builder.Predicates) {
	obj := &akov2generated.Parent{}
	return obj, builder.WithPredicates(h.predicates...)
}
func (h *Handler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("Parent").For(h.For()).Watches(&v1.Secret{}, handler.EnqueueRequestsFromMapFunc(indexers.NewParentBySecretMapFunc(h.Client)), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).WithOptions(defaultOptions).Complete(rec)
}

// getSDKClientSet creates an Atlas SDK client set using credentials from the resource's connection secret
func (h *Handler) getSDKClientSet(ctx context.Context, parent *akov2generated.Parent) (*atlas.ClientSet, error) {
	var connectionSecretRef *client.ObjectKey
	if parent.Spec.ConnectionSecretRef != nil {
		connectionSecretRef = &client.ObjectKey{
			Name:      parent.Spec.ConnectionSecretRef.Name,
			Namespace: parent.Namespace,
		}
	}

	connectionConfig, err := reconciler.GetConnectionConfig(ctx, h.Client, connectionSecretRef, &h.GlobalSecretRef)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Atlas credentials: %w", err)
	}

	clientSet, err := h.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, h.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Atlas SDK client: %w", err)
	}

	return clientSet, nil
}
