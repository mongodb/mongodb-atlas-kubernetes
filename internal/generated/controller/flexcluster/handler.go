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

package flexcluster

import (
	"context"
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	handler "sigs.k8s.io/controller-runtime/pkg/handler"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	atlas "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	reconciler "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	indexers "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers"
	translate "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate"
	indexer "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// getHandlerForResource selects the appropriate version-specific handler based on which resource spec version is set
func (h *Handler) getHandlerForResource(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.StateHandler[akov2generated.FlexCluster], error) {
	atlasClients, err := h.getSDKClientSet(ctx, flexcluster)
	if err != nil {
		return nil, err
	}
	// Check which resource spec version is set and validate that only one is specified
	var versionCount int
	var selectedHandler ctrlstate.StateHandler[akov2generated.FlexCluster]

	if flexcluster.Spec.V20250312 != nil {
		translationReq, err := getTranslationRequest(ctx, h.Client, "flexclusters.atlas.generated.mongodb.com", "v1", "v20250312")
		if err != nil {
			return nil, err
		}
		versionCount++
		selectedHandler = h.handlerv20250312(h.Client, atlasClients.SdkClient20250312006, translationReq, h.deletionProtection)
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
func (h *Handler) HandleInitial(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateInitial, err)
	}
	return handler.HandleInitial(ctx, flexcluster)
}

// HandleImportRequested delegates to the version-specific handler
func (h *Handler) HandleImportRequested(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}
	return handler.HandleImportRequested(ctx, flexcluster)
}

// HandleImported delegates to the version-specific handler
func (h *Handler) HandleImported(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateImported, err)
	}
	return handler.HandleImported(ctx, flexcluster)
}

// HandleCreating delegates to the version-specific handler
func (h *Handler) HandleCreating(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateCreating, err)
	}
	return handler.HandleCreating(ctx, flexcluster)
}

// HandleCreated delegates to the version-specific handler
func (h *Handler) HandleCreated(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateCreated, err)
	}
	return handler.HandleCreated(ctx, flexcluster)
}

// HandleUpdating delegates to the version-specific handler
func (h *Handler) HandleUpdating(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateUpdating, err)
	}
	return handler.HandleUpdating(ctx, flexcluster)
}

// HandleUpdated delegates to the version-specific handler
func (h *Handler) HandleUpdated(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateUpdated, err)
	}
	return handler.HandleUpdated(ctx, flexcluster)
}

// HandleDeletionRequested delegates to the version-specific handler
func (h *Handler) HandleDeletionRequested(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}
	return handler.HandleDeletionRequested(ctx, flexcluster)
}

// HandleDeleting delegates to the version-specific handler
func (h *Handler) HandleDeleting(ctx context.Context, flexcluster *akov2generated.FlexCluster) (ctrlstate.Result, error) {
	handler, err := h.getHandlerForResource(ctx, flexcluster)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}
	return handler.HandleDeleting(ctx, flexcluster)
}

// For returns the resource and predicates for the controller
func (h *Handler) For() (client.Object, builder.Predicates) {
	obj := &akov2generated.FlexCluster{}
	return obj, builder.WithPredicates(h.predicates...)
}
func (h *Handler) flexclusterForGroupMapFunc() handler.MapFunc {
	return indexer.ProjectsIndexMapperFunc(indexers.FlexClusterByGroupIndex, func() *akov2generated.FlexClusterList {
		return &akov2generated.FlexClusterList{}
	}, indexers.FlexClusterRequestsFromGroup, h.Client, h.Log)
}

// SetupWithManager sets up the controller with the Manager
func (h *Handler) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	h.Client = mgr.GetClient()
	return controllerruntime.NewControllerManagedBy(mgr).Named("FlexCluster").For(h.For()).Watches(&akov2generated.Group{}, handler.EnqueueRequestsFromMapFunc(h.flexclusterForGroupMapFunc()), builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).WithOptions(defaultOptions).Complete(rec)
}

// getSDKClientSet creates an Atlas SDK client set using credentials from the resource's connection secret
func (h *Handler) getSDKClientSet(ctx context.Context, flexcluster *akov2generated.FlexCluster) (*atlas.ClientSet, error) {
	var connectionSecretRef *client.ObjectKey
	if flexcluster.Spec.ConnectionSecretRef != nil {
		connectionSecretRef = &client.ObjectKey{
			Name:      flexcluster.Spec.ConnectionSecretRef.Name,
			Namespace: flexcluster.Namespace,
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

// getTranslationRequest creates a translation request for converting entities between API and AKO.
// This is a package-level function that can be called from any handler.
func getTranslationRequest(ctx context.Context, kubeClient client.Client, crdName string, storageVersion string, targetVersion string) (*translate.Request, error) {
	crd := &apiextensionsv1.CustomResourceDefinition{}
	err := kubeClient.Get(ctx, client.ObjectKey{Name: crdName}, crd)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve CRD %s: %w", crdName, err)
	}

	translator, err := translate.NewTranslator(crd, storageVersion, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to setup translator: %w", err)
	}

	return &translate.Request{
		Dependencies: nil,
		Translator:   translator,
	}, nil
}
