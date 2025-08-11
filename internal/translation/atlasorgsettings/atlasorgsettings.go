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
