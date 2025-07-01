// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

const (
	AnnotationContentHash = "mongodb.com/content-hash"
)

func (h *AtlasThirdPartyIntegrationHandler) HandleInitial(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	return h.upsert(ctx, state.StateInitial, state.StateCreated, integration)
}

func (h *AtlasThirdPartyIntegrationHandler) HandleCreated(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	return h.upsert(ctx, state.StateCreated, state.StateUpdated, integration)
}

func (h *AtlasThirdPartyIntegrationHandler) HandleUpdated(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	return h.upsert(ctx, state.StateUpdated, state.StateUpdated, integration)
}

func (h *AtlasThirdPartyIntegrationHandler) HandleDeletionRequested(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	req, err := h.newReconcileRequest(ctx, integration)
	if err != nil {
		// TODO is this good for all error cases?
		return h.unmanage(integration.Spec.Type)
	}

	if !h.deletionProtection {
		return h.delete(ctx, req, integration.Spec.Type)
	}
	return h.unmanage(integration.Spec.Type)
}

func (h *AtlasThirdPartyIntegrationHandler) upsert(ctx context.Context, currentState, nextState state.ResourceState, integration *akov2.AtlasThirdPartyIntegration) (ctrlstate.Result, error) {
	req, err := h.newReconcileRequest(ctx, integration)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to build reconcile request: %w", err))
	}

	integrationSpec, err := h.populateIntegration(ctx, integration)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to populate integration: %w", err))
	}
	atlasIntegration, err := req.Service.Get(ctx, req.Project.ID, integrationSpec.Type)
	if errors.Is(err, thirdpartyintegration.ErrNotFound) {
		return h.create(ctx, currentState, req, integrationSpec)
	}
	if err != nil {
		return result.Error(
			currentState,
			fmt.Errorf("Error getting %s Atlas Integration for project %s: %w",
				integrationSpec.Type, req.Project.ID, err),
		)
	}
	atlas := atlasIntegration.Comparable()
	spec := integrationSpec.Comparable()
	secretChanged, err := h.secretChanged(ctx, integration)
	if err != nil {
		return result.Error(
			currentState,
			fmt.Errorf("Error evaluating secret changes for %s Atlas Integration for project %s: %w",
				integrationSpec.Type, req.Project.ID, err),
		)
	}
	if secretChanged || !reflect.DeepEqual(atlas, spec) {
		return h.update(ctx, currentState, req, integrationSpec)
	}
	return result.NextState(
		nextState,
		fmt.Sprintf("Synced %s Atlas Third Party Integration for %s", integrationSpec.Type, req.Project.ID),
	)
}

func (h *AtlasThirdPartyIntegrationHandler) create(ctx context.Context, currentState state.ResourceState, req *reconcileRequest, integrationSpec *thirdpartyintegration.ThirdPartyIntegration) (ctrlstate.Result, error) {
	newIntegration, err := req.Service.Create(ctx, req.Project.ID, integrationSpec)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to create %s Atlas Third Party Integration for project %s: %w",
			integrationSpec.Type, req.Project.ID, err))
	}
	req.integration.Status.ID = newIntegration.ID
	if err := h.patchNonConditionStatus(ctx, req); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to record id for %s Atlas Third Party Integration for project %s: %w",
			integrationSpec.Type, req.Project.ID, err))
	}
	if err := h.ensureSecretHash(ctx, req.integration); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to ensure secret is hashed to detect further changes "+
			"for %s Atlas Third Party Integration for project %s: %w",
			integrationSpec.Type, req.Project.ID, err))
	}
	return result.NextState(
		state.StateCreated,
		fmt.Sprintf("Created Atlas Third Party Integration for %s", integrationSpec.Type),
	)
}

func (h *AtlasThirdPartyIntegrationHandler) update(ctx context.Context, currentState state.ResourceState, req *reconcileRequest, integrationSpec *thirdpartyintegration.ThirdPartyIntegration) (ctrlstate.Result, error) {
	updatedIntegration, err := req.Service.Update(ctx, req.Project.ID, integrationSpec)
	if req.integration.Status.ID == "" { // On imports, the ID might be unset
		req.integration.Status.ID = updatedIntegration.ID
		if err := h.patchNonConditionStatus(ctx, req); err != nil {
			return result.Error(currentState, fmt.Errorf("failed to record id for %s Atlas Third Party Integration for project %s: %w",
				integrationSpec.Type, req.Project.ID, err))
		}
	}
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update %s Atlas Third Party Integration for project %s: %w",
			integrationSpec.Type, req.Project.ID, err))
	}
	return result.NextState(
		state.StateUpdated,
		fmt.Sprintf("Updated Atlas Third Party Integration for %s", integrationSpec.Type),
	)
}

func (h *AtlasThirdPartyIntegrationHandler) delete(ctx context.Context, req *reconcileRequest, integrationType string) (ctrlstate.Result, error) {
	err := req.Service.Delete(ctx, req.Project.ID, integrationType)
	if errors.Is(err, thirdpartyintegration.ErrNotFound) {
		return h.unmanage(integrationType)
	}
	if err != nil {
		return result.Error(
			state.StateDeletionRequested,
			fmt.Errorf("Error deleting %s Atlas Integration for project %s: %w", integrationType, req.Project.ID, err),
		)
	}
	return h.unmanage(integrationType)
}

func (h *AtlasThirdPartyIntegrationHandler) unmanage(integrationType string) (ctrlstate.Result, error) {
	return result.NextState(
		state.StateDeleted,
		fmt.Sprintf("Deleted Atlas Third Party Integration for %s", integrationType),
	)
}

func (h *AtlasThirdPartyIntegrationHandler) patchNonConditionStatus(ctx context.Context, req *reconcileRequest) error {
	statusJSON, err := json.Marshal(req.integration)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}
	if err := h.Client.Status().Patch(ctx, req.integration, client.RawPatch(types.MergePatchType, statusJSON)); err != nil {
		return fmt.Errorf("failed to patch: %w", err)
	}
	return nil
}

func (h *AtlasThirdPartyIntegrationHandler) populateIntegration(ctx context.Context, integration *akov2.AtlasThirdPartyIntegration) (*thirdpartyintegration.ThirdPartyIntegration, error) {
	secrets, err := fetchIntegrationSecrets(ctx, h.Client, integration)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch integration secrets: %w", err)
	}
	internalIntegration, err := thirdpartyintegration.NewFromSpec(integration, secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to populate integration: %w", err)
	}
	return internalIntegration, nil
}
