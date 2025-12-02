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

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
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
	resp, httpResp, err := a.orgSettingsAPI.GetOrgSettings(ctx, orgID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get AtlasOrgSettings: %w", err)
	}
	statusCode := httputil.StatusCode(httpResp)
	if statusCode != 200 {
		return nil, fmt.Errorf("failed to get AtlasOrgSettings: expected status code 200. Got: %d", statusCode)
	}

	return NewFromAtlas(orgID, resp), nil
}

func (a *AtlasOrgSettingsServiceImpl) Update(ctx context.Context, orgID string, aos *AtlasOrgSettings) (*AtlasOrgSettings, error) {
	atlasOrgSettings := ToAtlas(aos)
	if atlasOrgSettings == nil {
		return nil, nil
	}

	resp, httpResp, err := a.orgSettingsAPI.UpdateOrgSettings(ctx, orgID, atlasOrgSettings).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update AtlasOrgSettings: %w", err)
	}
	statusCode := httputil.StatusCode(httpResp)
	if statusCode != 201 && statusCode != 200 {
		return nil, fmt.Errorf("failed to update AtlasOrgSettings: expected status code 200 or 201. Got: %d", statusCode)
	}

	return NewFromAtlas(orgID, resp), nil
}
