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

package customrole

import (
	"context"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312010/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	translate "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type Handlerv20250312 struct {
	kubeClient         client.Client
	atlasClient        *v20250312sdk.APIClient
	translationRequest *translate.Request
	deletionProtection bool
}

func NewHandlerv20250312(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translationRequest *translate.Request, deletionProtection bool) *Handlerv20250312 {
	return &Handlerv20250312{
		atlasClient:        atlasClient,
		deletionProtection: deletionProtection,
		kubeClient:         kubeClient,
		translationRequest: translationRequest,
	}
}

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement initial state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateUpdated, "Updated AtlasCustomRole.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement importrequested state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateImported, "Import completed")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement imported state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateUpdated, "Ready")
}

// HandleCreating handles the creating state for version v20250312
func (h *Handlerv20250312) HandleCreating(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement creating state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateCreated, "Resource created")
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement created state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateUpdated, "Ready")
}

// HandleUpdating handles the updating state for version v20250312
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement updating state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateUpdated, "Update completed")
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement updated state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateUpdated, "Ready")
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement deletionrequested state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateDeleting, "Deletion started")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, customrole *akov2generated.CustomRole) (ctrlstate.Result, error) {
	// TODO: Implement deleting state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateDeleted, "Deleted")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.CustomRole{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}
