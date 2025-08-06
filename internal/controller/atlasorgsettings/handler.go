package atlasorgsettings

import (
	"context"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

func (h *AtlasOrgSettingsHandler) HandleInitial(ctx context.Context, orgSettings *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return ctrlstate.Result{}, nil
}

func (h *AtlasOrgSettingsHandler) HandleCreated(ctx context.Context, orgSettings *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return ctrlstate.Result{}, nil
}

func (h *AtlasOrgSettingsHandler) HandleUpdated(ctx context.Context, orgSettings *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return ctrlstate.Result{}, nil
}

func (h *AtlasOrgSettingsHandler) HandleDeletionRequested(ctx context.Context, orgSettings *akov2.AtlasOrgSettings) (ctrlstate.Result, error) {
	return ctrlstate.Result{}, nil
}
