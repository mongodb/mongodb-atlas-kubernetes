package atlasorgsettings

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

type AtlasOrgSettingsService interface {
	Get(ctx context.Context, orgID string) (*AtlasOrgSettings, error)
	Update(ctx context.Context, orgID string, aos *AtlasOrgSettings) (*AtlasOrgSettings, error)
}

type AtlasOrgSettingsServiceImpl struct {
	orgSettingsAPI admin.OrganizationsApi
}

func NewAtlasOrgSettingsService(api admin.OrganizationsApi) AtlasOrgSettingsService {
	return &AtlasOrgSettingsServiceImpl{
		orgSettingsAPI: api,
	}
}

func (a *AtlasOrgSettingsServiceImpl) Get(ctx context.Context, orgID string) (*AtlasOrgSettings, error) {
	resp, httpResp, err := a.orgSettingsAPI.GetOrganizationSettings(ctx, orgID).Execute()
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, err
	}

	return NewFromAtlas(orgID, resp), nil
}

func (a *AtlasOrgSettingsServiceImpl) Update(ctx context.Context, orgID string, aos *AtlasOrgSettings) (*AtlasOrgSettings, error) {
	atlasOrgSettings := ToAtlas(aos)
	if atlasOrgSettings == nil {
		return nil, nil
	}

	resp, httpResp, err := a.orgSettingsAPI.UpdateOrganizationSettings(ctx, orgID, atlasOrgSettings).Execute()
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 201 && httpResp.StatusCode != 200 {
		return nil, err
	}

	return NewFromAtlas(orgID, resp), nil
}
