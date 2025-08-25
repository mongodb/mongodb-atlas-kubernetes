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

package atlasorgsettings

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/atlasorgsettings"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type reconcileRequest struct {
	svc atlasorgsettings.AtlasOrgSettingsService
	aos *akov2.AtlasOrgSettings
}

func (h *AtlasOrgSettingsHandler) newReconcileRequest(ctx context.Context, aos *akov2.AtlasOrgSettings) (*reconcileRequest, error) {
	var objKey *client.ObjectKey
	if aos.Spec.ConnectionSecretRef != nil && aos.Spec.ConnectionSecretRef.Name != "" {
		objKey = &client.ObjectKey{
			Namespace: aos.GetNamespace(),
			Name:      aos.Spec.ConnectionSecretRef.Name,
		}
	}

	cfg, err := reconciler.GetConnectionConfig(ctx, h.Client, objKey, &h.GlobalSecretRef)
	if err != nil {
		return nil, err
	}

	atlasSdk, err := h.AtlasProvider.SdkClientSet(ctx, cfg.Credentials, h.Log)
	if err != nil {
		return nil, err
	}
	return &reconcileRequest{
		svc: h.serviceBuilder(atlasSdk),
		aos: aos,
	}, nil
}

func (h *AtlasOrgSettingsHandler) upsert(ctx context.Context, currentState, nextState state.ResourceState,
	aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	reconcileCtx, err := h.newReconcileRequest(ctx, aos)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to create reconcile context: %w", err))
	}

	currentAtlasSettings, err := reconcileCtx.svc.Get(ctx, aos.Spec.OrgID)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to get current org settings from Atlas: %w", err))
	}

	desiredSettings := atlasorgsettings.NewFromAKO(aos.Spec)

	if !desiredSettings.Equal(currentAtlasSettings) {
		resp, apiErr := reconcileCtx.svc.Update(ctx, aos.Spec.OrgID, desiredSettings)
		if apiErr != nil {
			return result.Error(currentState, apiErr)
		}
		if resp == nil {
			return result.Error(currentState, fmt.Errorf("atlas returned OrgSettings which is nil after update"))
		}

		return result.NextState(nextState, "Updated")
	}

	return result.NextState(nextState, "Ready")
}

func (h *AtlasOrgSettingsHandler) unmanage(orgID string) (ctrlstate.Result, error) {
	return result.NextState(state.StateDeleted, fmt.Sprintf("unmanaged AtlasOrgSettings for orgID %s.", orgID))
}

func (h *AtlasOrgSettingsHandler) HandleInitial(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Updated AtlasOrgSettings.")
}

func (h *AtlasOrgSettingsHandler) HandleUpdated(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return h.upsert(ctx, state.StateUpdated, state.StateUpdated, aos)
}

func (h *AtlasOrgSettingsHandler) HandleDeletionRequested(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return h.unmanage(aos.Spec.OrgID)
}
