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
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type AtlasOrgSettings struct {
	akov2.AtlasOrgSettingsSpec
}

func NewFromAKO(spec akov2.AtlasOrgSettingsSpec) *AtlasOrgSettings {
	return &AtlasOrgSettings{
		AtlasOrgSettingsSpec: spec,
	}
}

func NewFromAtlas(orgID string, atlasSpec *admin.OrganizationSettings) *AtlasOrgSettings {
	if atlasSpec == nil {
		return nil
	}

	return &AtlasOrgSettings{
		AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
			OrgID:                                  orgID,
			ConnectionSecret:                       nil,
			ApiAccessListRequired:                  atlasSpec.ApiAccessListRequired,
			GenAIFeaturesEnabled:                   atlasSpec.GenAIFeaturesEnabled,
			MaxServiceAccountSecretValidityInHours: atlasSpec.MaxServiceAccountSecretValidityInHours,
			MultiFactorAuthRequired:                atlasSpec.MultiFactorAuthRequired,
			RestrictEmployeeAccess:                 atlasSpec.RestrictEmployeeAccess,
			SecurityContact:                        atlasSpec.SecurityContact,
			StreamsCrossGroupEnabled:               atlasSpec.StreamsCrossGroupEnabled,
		},
	}
}

func ToAtlas(orgSettings *AtlasOrgSettings) *admin.OrganizationSettings {
	if orgSettings == nil {
		return nil
	}

	return &admin.OrganizationSettings{
		ApiAccessListRequired:                  orgSettings.ApiAccessListRequired,
		GenAIFeaturesEnabled:                   orgSettings.GenAIFeaturesEnabled,
		MaxServiceAccountSecretValidityInHours: orgSettings.MaxServiceAccountSecretValidityInHours,
		MultiFactorAuthRequired:                orgSettings.MultiFactorAuthRequired,
		RestrictEmployeeAccess:                 orgSettings.RestrictEmployeeAccess,
		SecurityContact:                        orgSettings.SecurityContact,
	}
}
