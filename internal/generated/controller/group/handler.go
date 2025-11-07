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
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set
func (h *GroupHandler) getHandlerForResource(ctx context.Context, group *v1.Group) (ctrlstate.StateHandler[v1.Group], error) {
	atlasClients, err := h.getSDKClientSet(ctx, group)
	if err != nil {
		return nil, err
	}

	translationReq, err := h.getTranslationRequest(ctx)
	if err != nil {
		return nil, err
	}

	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[v1.Group]

	if group.Spec.V20250312 != nil {
		versionCount++
		selectedHandler = h.handlerv20250312(h.Client, atlasClients.SdkClient20250312008, translationReq)
	}

	if versionCount == 0 {
		return nil, fmt.Errorf("no resource spec version specified - please set one of the available spec versions")
	}
	if versionCount > 1 {
		return nil, fmt.Errorf("multiple resource spec versions specified - please set only one spec version")
	}
	return selectedHandler, nil
}

func (h *GroupHandler) getSDKClientSet(ctx context.Context, group *v1.Group) (*atlas.ClientSet, error) {
	connectionConfig, err := reconciler.GetConnectionConfig(
		ctx,
		h.Client,
		&client.ObjectKey{
			Namespace: group.Namespace,
			Name:      group.Spec.ConnectionSecretRef.Name,
		},
		&h.GlobalSecretRef,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Atlas credentials: %w", err)
	}

	atlasClients, err := h.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, h.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Atlas SDK client: %w", err)
	}

	return atlasClients, nil
}

func (h *GroupHandler) getTranslationRequest(ctx context.Context) (*translate.Request, error) {
	groupCRD := &apiextensionsv1.CustomResourceDefinition{}
	err := h.Client.Get(ctx, client.ObjectKey{Name: "groups.atlas.generated.mongodb.com"}, groupCRD)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Group CRD: %w", err)
	}

	translator, err := translate.NewTranslator(groupCRD, "v1", "v20250312")
	if err != nil {
		return nil, fmt.Errorf("failed to setup translator: %w", err)
	}

	return &translate.Request{
		Translator:   translator,
		Dependencies: nil,
	}, nil
}

// HandleInitial delegates to the version-specific handler
func (h *GroupHandler) HandleInitial(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, group)
}

// HandleImportRequested delegates to the version-specific handler
func (h *GroupHandler) HandleImportRequested(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}
	return handler.HandleImportRequested(ctx, group)
}

// HandleImported delegates to the version-specific handler
func (h *GroupHandler) HandleImported(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateImported, err)
	}
	return handler.HandleImported(ctx, group)
}

// HandleCreating delegates to the version-specific handler
func (h *GroupHandler) HandleCreating(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateCreating, err)
	}
	return handler.HandleCreating(ctx, group)
}

// HandleCreated delegates to the version-specific handler
func (h *GroupHandler) HandleCreated(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateCreated, err)
	}
	return handler.HandleCreated(ctx, group)
}

// HandleUpdating delegates to the version-specific handler
func (h *GroupHandler) HandleUpdating(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateUpdating, err)
	}
	return handler.HandleUpdating(ctx, group)
}

// HandleUpdated delegates to the version-specific handler
func (h *GroupHandler) HandleUpdated(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateUpdated, err)
	}
	return handler.HandleUpdated(ctx, group)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *GroupHandler) HandleDeletionRequested(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}
	return handler.HandleDeletionRequested(ctx, group)
}

// HandleDeleting delegates to the version-specific handler
func (h *GroupHandler) HandleDeleting(ctx context.Context, group *v1.Group) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, group)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}
	return handler.HandleDeleting(ctx, group)
}

// For returns the resource and predicates for the controller
func (h *GroupHandler) For() (client.Object, builder.Predicates) {
	obj := &v1.Group{}
	return obj, builder.WithPredicates(h.predicates...)
}

// SetupWithManager sets up the controller with the Manager
func (h *GroupHandler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("Group").For(h.For()).WithOptions(defaultOptions).Complete(rec)
}
