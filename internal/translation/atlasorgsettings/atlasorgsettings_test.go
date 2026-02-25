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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestNewFromAKO(t *testing.T) {
	tests := []struct {
		name     string
		spec     akov2.AtlasOrgSettingsSpec
		expected *AtlasOrgSettings
	}{
		{
			name: "complete org settings spec",
			spec: akov2.AtlasOrgSettingsSpec{
				OrgID: "test-org-id",
				ConnectionSecretRef: &api.LocalObjectReference{
					Name: "test-secret",
				},
				ApiAccessListRequired:                  pointer.MakePtr(true),
				GenAIFeaturesEnabled:                   pointer.MakePtr(false),
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
				MultiFactorAuthRequired:                pointer.MakePtr(true),
				RestrictEmployeeAccess:                 pointer.MakePtr(false),
				SecurityContact:                        pointer.MakePtr("security@example.com"),
				StreamsCrossGroupEnabled:               pointer.MakePtr(true),
			},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID: "test-org-id",
					ConnectionSecretRef: &api.LocalObjectReference{
						Name: "test-secret",
					},
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
			},
		},
		{
			name: "minimal org settings spec",
			spec: akov2.AtlasOrgSettingsSpec{
				OrgID: "minimal-org-id",
			},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID: "minimal-org-id",
				},
			},
		},
		{
			name: "empty org settings spec",
			spec: akov2.AtlasOrgSettingsSpec{},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewFromAKO(tt.spec)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewFromAtlas(t *testing.T) {
	tests := []struct {
		name      string
		orgID     string
		atlasSpec *admin.OrganizationSettings
		expected  *AtlasOrgSettings
	}{
		{
			name:  "complete atlas organization settings",
			orgID: "test-org-id",
			atlasSpec: &admin.OrganizationSettings{
				ApiAccessListRequired:                  pointer.MakePtr(true),
				GenAIFeaturesEnabled:                   pointer.MakePtr(false),
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
				MultiFactorAuthRequired:                pointer.MakePtr(true),
				RestrictEmployeeAccess:                 pointer.MakePtr(false),
				SecurityContact:                        pointer.MakePtr("security@example.com"),
				StreamsCrossGroupEnabled:               pointer.MakePtr(true),
			},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  "test-org-id",
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
			},
		},
		{
			name:  "atlas settings with some nil values",
			orgID: "test-org-id-2",
			atlasSpec: &admin.OrganizationSettings{
				ApiAccessListRequired:                  pointer.MakePtr(false),
				GenAIFeaturesEnabled:                   nil,
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(0),
				MultiFactorAuthRequired:                nil,
				RestrictEmployeeAccess:                 pointer.MakePtr(true),
				SecurityContact:                        nil,
				StreamsCrossGroupEnabled:               pointer.MakePtr(false),
			},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  "test-org-id-2",
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(false),
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(0),
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 pointer.MakePtr(true),
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               pointer.MakePtr(false),
				},
			},
		},
		{
			name:      "nil atlas spec returns nil",
			orgID:     "test-org-id",
			atlasSpec: nil,
			expected:  nil,
		},
		{
			name:  "empty org id with valid atlas spec",
			orgID: "",
			atlasSpec: &admin.OrganizationSettings{
				ApiAccessListRequired: pointer.MakePtr(true),
			},
			expected: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  "",
					ConnectionSecretRef:                    nil,
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewFromAtlas(tt.orgID, tt.atlasSpec)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAtlas(t *testing.T) {
	tests := []struct {
		name        string
		orgSettings *AtlasOrgSettings
		expected    *admin.OrganizationSettings
	}{
		{
			name: "complete org settings to atlas",
			orgSettings: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID: "test-org-id",
					ConnectionSecretRef: &api.LocalObjectReference{
						Name: "test-secret",
					},
					ApiAccessListRequired:                  pointer.MakePtr(true),
					GenAIFeaturesEnabled:                   pointer.MakePtr(false),
					MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
					MultiFactorAuthRequired:                pointer.MakePtr(true),
					RestrictEmployeeAccess:                 pointer.MakePtr(false),
					SecurityContact:                        pointer.MakePtr("security@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(true),
				},
			},
			expected: &admin.OrganizationSettings{
				ApiAccessListRequired:                  pointer.MakePtr(true),
				GenAIFeaturesEnabled:                   pointer.MakePtr(false),
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
				MultiFactorAuthRequired:                pointer.MakePtr(true),
				RestrictEmployeeAccess:                 pointer.MakePtr(false),
				SecurityContact:                        pointer.MakePtr("security@example.com"),
			},
		},
		{
			name: "org settings with some nil values",
			orgSettings: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  "test-org-id",
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   pointer.MakePtr(true),
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                pointer.MakePtr(false),
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        pointer.MakePtr("admin@example.com"),
					StreamsCrossGroupEnabled:               pointer.MakePtr(false),
				},
			},
			expected: &admin.OrganizationSettings{
				ApiAccessListRequired:                  nil,
				GenAIFeaturesEnabled:                   pointer.MakePtr(true),
				MaxServiceAccountSecretValidityInHours: nil,
				MultiFactorAuthRequired:                pointer.MakePtr(false),
				RestrictEmployeeAccess:                 nil,
				SecurityContact:                        pointer.MakePtr("admin@example.com"),
			},
		},
		{
			name: "org settings with all nil values",
			orgSettings: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					OrgID:                                  "test-org-id",
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
			},
			expected: &admin.OrganizationSettings{
				ApiAccessListRequired:                  nil,
				GenAIFeaturesEnabled:                   nil,
				MaxServiceAccountSecretValidityInHours: nil,
				MultiFactorAuthRequired:                nil,
				RestrictEmployeeAccess:                 nil,
				SecurityContact:                        nil,
			},
		},
		{
			name:        "nil org settings returns nil",
			orgSettings: nil,
			expected:    nil,
		},
		{
			name: "empty org settings struct",
			orgSettings: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{},
			},
			expected: &admin.OrganizationSettings{
				ApiAccessListRequired:                  nil,
				GenAIFeaturesEnabled:                   nil,
				MaxServiceAccountSecretValidityInHours: nil,
				MultiFactorAuthRequired:                nil,
				RestrictEmployeeAccess:                 nil,
				SecurityContact:                        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAtlas(tt.orgSettings)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtlasOrgSettings_Equal(t *testing.T) {
	// Helper function to create a complete AtlasOrgSettings instance
	createCompleteSettings := func() *AtlasOrgSettings {
		return &AtlasOrgSettings{
			AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
				ApiAccessListRequired:                  pointer.MakePtr(true),
				GenAIFeaturesEnabled:                   pointer.MakePtr(false),
				MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
				MultiFactorAuthRequired:                pointer.MakePtr(true),
				RestrictEmployeeAccess:                 pointer.MakePtr(false),
				SecurityContact:                        pointer.MakePtr("security@example.com"),
				StreamsCrossGroupEnabled:               pointer.MakePtr(true),
			},
		}
	}

	tests := []struct {
		name     string
		a        *AtlasOrgSettings
		other    *AtlasOrgSettings
		expected bool
	}{
		{
			name:     "identical complete settings should be equal",
			a:        createCompleteSettings(),
			other:    createCompleteSettings(),
			expected: true,
		},
		{
			name:     "comparing with nil should return false",
			a:        createCompleteSettings(),
			other:    nil,
			expected: false,
		},
		{
			name: "different ApiAccessListRequired values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.ApiAccessListRequired = pointer.MakePtr(false)
				return s
			}(),
			expected: false,
		},
		{
			name: "different GenAIFeaturesEnabled values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.GenAIFeaturesEnabled = pointer.MakePtr(true)
				return s
			}(),
			expected: false,
		},
		{
			name: "different MaxServiceAccountSecretValidityInHours values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MaxServiceAccountSecretValidityInHours = pointer.MakePtr(48)
				return s
			}(),
			expected: false,
		},
		{
			name: "different MultiFactorAuthRequired values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MultiFactorAuthRequired = pointer.MakePtr(false)
				return s
			}(),
			expected: false,
		},
		{
			name: "different RestrictEmployeeAccess values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.RestrictEmployeeAccess = pointer.MakePtr(true)
				return s
			}(),
			expected: false,
		},
		{
			name: "different SecurityContact values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.SecurityContact = pointer.MakePtr("different@example.com")
				return s
			}(),
			expected: false,
		},
		{
			name: "different StreamsCrossGroupEnabled values should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.StreamsCrossGroupEnabled = pointer.MakePtr(false)
				return s
			}(),
			expected: false,
		},
		{
			name: "nil ApiAccessListRequired in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.ApiAccessListRequired = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil ApiAccessListRequired in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.ApiAccessListRequired = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil GenAIFeaturesEnabled in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.GenAIFeaturesEnabled = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil GenAIFeaturesEnabled in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.GenAIFeaturesEnabled = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil MaxServiceAccountSecretValidityInHours in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MaxServiceAccountSecretValidityInHours = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil MaxServiceAccountSecretValidityInHours in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MaxServiceAccountSecretValidityInHours = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil MultiFactorAuthRequired in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MultiFactorAuthRequired = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil MultiFactorAuthRequired in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.MultiFactorAuthRequired = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil RestrictEmployeeAccess in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.RestrictEmployeeAccess = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil RestrictEmployeeAccess in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.RestrictEmployeeAccess = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil SecurityContact in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.SecurityContact = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil SecurityContact in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.SecurityContact = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "nil StreamsCrossGroupEnabled in first object should not be equal",
			a: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.StreamsCrossGroupEnabled = nil
				return s
			}(),
			other:    createCompleteSettings(),
			expected: false,
		},
		{
			name: "nil StreamsCrossGroupEnabled in second object should not be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				s.StreamsCrossGroupEnabled = nil
				return s
			}(),
			expected: false,
		},
		{
			name: "both objects with all nil fields should not be equal",
			a: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
			},
			other: &AtlasOrgSettings{
				AtlasOrgSettingsSpec: akov2.AtlasOrgSettingsSpec{
					ApiAccessListRequired:                  nil,
					GenAIFeaturesEnabled:                   nil,
					MaxServiceAccountSecretValidityInHours: nil,
					MultiFactorAuthRequired:                nil,
					RestrictEmployeeAccess:                 nil,
					SecurityContact:                        nil,
					StreamsCrossGroupEnabled:               nil,
				},
			},
			expected: false,
		},
		{
			name: "self comparison should be equal",
			a:    createCompleteSettings(),
			other: func() *AtlasOrgSettings {
				s := createCompleteSettings()
				return s
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Equal(tt.other)
			assert.Equal(t, tt.expected, result, "Expected Equal() to return %v, but got %v", tt.expected, result)
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	t.Run("AKO -> Atlas -> AKO round trip", func(t *testing.T) {
		originalAKOSpec := akov2.AtlasOrgSettingsSpec{
			OrgID: "test-org-id",
			ConnectionSecretRef: &api.LocalObjectReference{
				Name: "test-secret",
			},
			ApiAccessListRequired:                  pointer.MakePtr(true),
			GenAIFeaturesEnabled:                   pointer.MakePtr(false),
			MaxServiceAccountSecretValidityInHours: pointer.MakePtr(24),
			MultiFactorAuthRequired:                pointer.MakePtr(true),
			RestrictEmployeeAccess:                 pointer.MakePtr(false),
			SecurityContact:                        pointer.MakePtr("security@example.com"),
			StreamsCrossGroupEnabled:               pointer.MakePtr(true),
		}

		// Convert AKO -> AtlasOrgSettings -> Atlas -> AtlasOrgSettings
		atlasOrgSettings := NewFromAKO(originalAKOSpec)
		atlasFormat := ToAtlas(atlasOrgSettings)
		roundTripSettings := NewFromAtlas(originalAKOSpec.OrgID, atlasFormat)

		// Verify core fields are preserved (note: ConnectionSecret is not in Atlas format)
		assert.Equal(t, originalAKOSpec.OrgID, roundTripSettings.OrgID)
		assert.Equal(t, originalAKOSpec.ApiAccessListRequired, roundTripSettings.ApiAccessListRequired)
		assert.Equal(t, originalAKOSpec.GenAIFeaturesEnabled, roundTripSettings.GenAIFeaturesEnabled)
		assert.Equal(t, originalAKOSpec.MaxServiceAccountSecretValidityInHours, roundTripSettings.MaxServiceAccountSecretValidityInHours)
		assert.Equal(t, originalAKOSpec.MultiFactorAuthRequired, roundTripSettings.MultiFactorAuthRequired)
		assert.Equal(t, originalAKOSpec.RestrictEmployeeAccess, roundTripSettings.RestrictEmployeeAccess)
		assert.Equal(t, originalAKOSpec.SecurityContact, roundTripSettings.SecurityContact)

		// Note: StreamsCrossGroupEnabled is not included in ToAtlas conversion, so it will be nil
		assert.Nil(t, roundTripSettings.StreamsCrossGroupEnabled)
	})

	t.Run("Atlas -> AKO -> Atlas round trip", func(t *testing.T) {
		originalAtlasSettings := &admin.OrganizationSettings{
			ApiAccessListRequired:                  pointer.MakePtr(false),
			GenAIFeaturesEnabled:                   pointer.MakePtr(true),
			MaxServiceAccountSecretValidityInHours: pointer.MakePtr(48),
			MultiFactorAuthRequired:                pointer.MakePtr(false),
			RestrictEmployeeAccess:                 pointer.MakePtr(true),
			SecurityContact:                        pointer.MakePtr("admin@example.com"),
			StreamsCrossGroupEnabled:               pointer.MakePtr(false),
		}

		// Convert Atlas -> AtlasOrgSettings -> Atlas
		atlasOrgSettings := NewFromAtlas("test-org", originalAtlasSettings)
		roundTripAtlas := ToAtlas(atlasOrgSettings)

		// Verify all fields are preserved except StreamsCrossGroupEnabled (not in ToAtlas)
		assert.Equal(t, originalAtlasSettings.ApiAccessListRequired, roundTripAtlas.ApiAccessListRequired)
		assert.Equal(t, originalAtlasSettings.GenAIFeaturesEnabled, roundTripAtlas.GenAIFeaturesEnabled)
		assert.Equal(t, originalAtlasSettings.MaxServiceAccountSecretValidityInHours, roundTripAtlas.MaxServiceAccountSecretValidityInHours)
		assert.Equal(t, originalAtlasSettings.MultiFactorAuthRequired, roundTripAtlas.MultiFactorAuthRequired)
		assert.Equal(t, originalAtlasSettings.RestrictEmployeeAccess, roundTripAtlas.RestrictEmployeeAccess)
		assert.Equal(t, originalAtlasSettings.SecurityContact, roundTripAtlas.SecurityContact)

		// Note: StreamsCrossGroupEnabled is not included in ToAtlas, so it's not in the round trip
		assert.Nil(t, roundTripAtlas.StreamsCrossGroupEnabled)
	})
}
