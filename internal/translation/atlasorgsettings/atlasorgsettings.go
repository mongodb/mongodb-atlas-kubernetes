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
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

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

func (a *AtlasOrgSettings) Equal(other *AtlasOrgSettings) bool {
	if other == nil {
		return false
	}

	return a.ApiAccessListRequired != nil && other.ApiAccessListRequired != nil &&
		*a.ApiAccessListRequired == *other.ApiAccessListRequired &&
		a.GenAIFeaturesEnabled != nil && other.GenAIFeaturesEnabled != nil &&
		*a.GenAIFeaturesEnabled == *other.GenAIFeaturesEnabled &&
		a.MaxServiceAccountSecretValidityInHours != nil && other.MaxServiceAccountSecretValidityInHours != nil &&
		*a.MaxServiceAccountSecretValidityInHours == *other.MaxServiceAccountSecretValidityInHours &&
		a.MultiFactorAuthRequired != nil && other.MultiFactorAuthRequired != nil &&
		*a.MultiFactorAuthRequired == *other.MultiFactorAuthRequired &&
		a.RestrictEmployeeAccess != nil && other.RestrictEmployeeAccess != nil &&
		*a.RestrictEmployeeAccess == *other.RestrictEmployeeAccess &&
		a.SecurityContact != nil && other.SecurityContact != nil &&
		*a.SecurityContact == *other.SecurityContact &&
		a.StreamsCrossGroupEnabled != nil && other.StreamsCrossGroupEnabled != nil &&
		*a.StreamsCrossGroupEnabled == *other.StreamsCrossGroupEnabled
}

func NewFromAtlas(orgID string, atlasSpec *admin.OrganizationSettings) *AtlasOrgSettings {
	if atlasSpec == nil {
		return nil
	}

	return &AtlasOrgSettings{
		AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
			OrgID:                                  orgID,
			ConnectionSecretRef:                    nil,
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
