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

package v1

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func Test_FederatedAuthSpec_ToAtlas(t *testing.T) {
	t.Run("Can convert valid spec to Atlas", func(t *testing.T) {
		orgID := "test-org"
		idpID := "test-idp"
		projectName := "test-project"
		projectID := "test-project-id"

		projectNameToID := map[string]string{
			projectName: projectID,
		}

		spec := &AtlasFederatedAuthSpec{
			Enabled:                     true,
			ConnectionSecretRef:         common.ResourceRefNamespaced{},
			DomainAllowList:             []string{"test.com"},
			DomainRestrictionEnabled:    pointer.MakePtr(true),
			DataAccessIdentityProviders: &[]string{"test-123", "test-456"},
			SSODebugEnabled:             pointer.MakePtr(true),
			PostAuthRoleGrants:          []string{"role-3", "role-4"},
			RoleMappings: []RoleMapping{
				{
					ExternalGroupName: "test-group",
					RoleAssignments: []RoleAssignment{
						{
							ProjectName: projectName,
							Role:        "test-role",
						},
					},
				},
			},
		}

		result, err := spec.ToAtlas(orgID, idpID, projectNameToID)

		assert.NoError(t, err, "ToAtlas() failed")
		assert.NotNil(t, result, "ToAtlas() result is nil")

		expected := &admin.ConnectedOrgConfig{
			DomainAllowList:               &spec.DomainAllowList,
			DomainRestrictionEnabled:      *spec.DomainRestrictionEnabled,
			DataAccessIdentityProviderIds: spec.DataAccessIdentityProviders,
			IdentityProviderId:            &idpID,
			OrgId:                         orgID,
			PostAuthRoleGrants:            &spec.PostAuthRoleGrants,
			RoleMappings: &[]admin.AuthFederationRoleMapping{
				{
					ExternalGroupName: spec.RoleMappings[0].ExternalGroupName,
					RoleAssignments: &[]admin.ConnectedOrgConfigRoleAssignment{
						{
							GroupId: &projectID,
							Role:    &spec.RoleMappings[0].RoleAssignments[0].Role,
						},
					},
				},
			},
			UserConflicts: nil,
		}

		diff := deep.Equal(expected, result)
		if diff != nil {
			t.Log(cmp.Diff(expected, result))
		}
		assert.Nil(t, diff, diff)
	})

	t.Run("Should return an error when project is not available", func(t *testing.T) {
		orgID := "test-org"
		idpID := "test-idp"
		projectName := "test-project"

		projectNameToID := map[string]string{}

		spec := &AtlasFederatedAuthSpec{
			Enabled:                  true,
			ConnectionSecretRef:      common.ResourceRefNamespaced{},
			DomainAllowList:          []string{"test.com"},
			DomainRestrictionEnabled: pointer.MakePtr(true),
			SSODebugEnabled:          pointer.MakePtr(true),
			PostAuthRoleGrants:       []string{"role-3", "role-4"},
			RoleMappings: []RoleMapping{
				{
					ExternalGroupName: "test-group",
					RoleAssignments: []RoleAssignment{
						{
							ProjectName: projectName,
							Role:        "test-role",
						},
					},
				},
			},
		}

		result, err := spec.ToAtlas(orgID, idpID, projectNameToID)

		assert.Error(t, err, "ToAtlas() should fail")
		assert.NotNil(t, result, "ToAtlas() result should not be nil")
	})
}
