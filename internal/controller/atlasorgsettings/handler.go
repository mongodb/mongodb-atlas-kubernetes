package atlasorgsettings

import (
	"context"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func (h *AtlasOrgSettingsHandler) HandleInitial(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "Initialized")
}

func (h *AtlasOrgSettingsHandler) HandleImportRequested(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateImported, "Importing")
}

func (h *AtlasOrgSettingsHandler) HandleImported(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "Imported")
}

func (h *AtlasOrgSettingsHandler) HandleCreating(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "Creating")
}

func (h *AtlasOrgSettingsHandler) HandleCreated(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "Created")
}

func (h *AtlasOrgSettingsHandler) HandleUpdating(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Updating")
}

func (h *AtlasOrgSettingsHandler) HandleUpdated(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Updated")
}

func (h *AtlasOrgSettingsHandler) HandleDeletionRequested(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateDeleting, "DeletionRequested")
}

func (h *AtlasOrgSettingsHandler) HandleDeleting(ctx context.Context, aos *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return result.NextState(state.StateDeleted, "Deleted")
}
