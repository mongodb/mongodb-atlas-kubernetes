package atlasorgsettings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

func TestNewFromAKO(t *testing.T) {
	spec := akov2.AtlasOrgSettingsSpec{
		OrgID: "test-org",
	}
	result := NewFromAKO(spec)
	assert.NotNil(t, result)
	assert.Equal(t, "test-org", result.OrgID)
}

func TestNewFromAtlas_NilSpec(t *testing.T) {
	result := NewFromAtlas("org-id", nil)
	assert.Nil(t, result)
}

func TestNewFromAtlas_ValidSpec(t *testing.T) {
	atlasSpec := &admin.OrganizationSettings{
		ApiAccessListRequired:                  admin.PtrBool(true),
		GenAIFeaturesEnabled:                   admin.PtrBool(false),
		MaxServiceAccountSecretValidityInHours: admin.PtrInt(24),
		MultiFactorAuthRequired:                admin.PtrBool(true),
		RestrictEmployeeAccess:                 admin.PtrBool(false),
		SecurityContact:                        admin.PtrString("security@example.com"),
		StreamsCrossGroupEnabled:               admin.PtrBool(true),
	}
	result := NewFromAtlas("org-id", atlasSpec)
	assert.NotNil(t, result)
	assert.Equal(t, "org-id", result.OrgID)
	assert.Equal(t, atlasSpec.ApiAccessListRequired, result.ApiAccessListRequired)
	assert.Equal(t, atlasSpec.GenAIFeaturesEnabled, result.GenAIFeaturesEnabled)
	assert.Equal(t, atlasSpec.MaxServiceAccountSecretValidityInHours, result.MaxServiceAccountSecretValidityInHours)
	assert.Equal(t, atlasSpec.MultiFactorAuthRequired, result.MultiFactorAuthRequired)
	assert.Equal(t, atlasSpec.RestrictEmployeeAccess, result.RestrictEmployeeAccess)
	assert.Equal(t, atlasSpec.SecurityContact, result.SecurityContact)
	assert.Equal(t, atlasSpec.StreamsCrossGroupEnabled, result.StreamsCrossGroupEnabled)
}

func TestToAtlas_NilInput(t *testing.T) {
	result := ToAtlas(nil)
	assert.Nil(t, result)
}

func TestToAtlas_ValidInput(t *testing.T) {
	orgSettings := &AtlasOrgSettings{
		AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
			ApiAccessListRequired:                  admin.PtrBool(true),
			GenAIFeaturesEnabled:                   admin.PtrBool(false),
			MaxServiceAccountSecretValidityInHours: admin.PtrInt(12),
			MultiFactorAuthRequired:                admin.PtrBool(true),
			RestrictEmployeeAccess:                 admin.PtrBool(false),
			SecurityContact:                        admin.PtrString("contact@org.com"),
		},
	}
	result := ToAtlas(orgSettings)
	assert.NotNil(t, result)
	assert.Equal(t, orgSettings.ApiAccessListRequired, result.ApiAccessListRequired)
	assert.Equal(t, orgSettings.GenAIFeaturesEnabled, result.GenAIFeaturesEnabled)
	assert.Equal(t, orgSettings.MaxServiceAccountSecretValidityInHours, result.MaxServiceAccountSecretValidityInHours)
	assert.Equal(t, orgSettings.MultiFactorAuthRequired, result.MultiFactorAuthRequired)
	assert.Equal(t, orgSettings.RestrictEmployeeAccess, result.RestrictEmployeeAccess)
	assert.Equal(t, orgSettings.SecurityContact, result.SecurityContact)
}
