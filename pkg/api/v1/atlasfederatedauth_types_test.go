package v1

import (
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	"testing"

	"github.com/go-test/deep"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/stretchr/testify/assert"
)

func Test_FederatedAuthSpec_ToAtlas(t *testing.T) {
	t.Run("Can convert valid spec to Atlas", func(t *testing.T) {
		orgID := "test-org"
		idpID := "test-idp"
		projectName := "test-project"
		projectId := "test-project-id"

		projectNameToID := map[string]string{
			projectName: projectId,
		}

		spec := &AtlasFederatedAuthSpec{
			Enabled:                  true,
			ConnectionSecretRef:      common.ResourceRefNamespaced{},
			DomainAllowList:          []string{"test.com"},
			DomainRestrictionEnabled: toptr.MakePtr(true),
			SSODebugEnabled:          toptr.MakePtr(true),
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

		assert.NoError(t, err, "ToAtlas() failed")
		assert.NotNil(t, result, "ToAtlas() result is nil")

		expected := &admin.ConnectedOrgConfig{
			DomainAllowList:          &spec.DomainAllowList,
			DomainRestrictionEnabled: *spec.DomainRestrictionEnabled,
			IdentityProviderId:       idpID,
			OrgId:                    orgID,
			PostAuthRoleGrants:       &spec.PostAuthRoleGrants,
			RoleMappings: &[]admin.AuthFederationRoleMapping{
				{
					ExternalGroupName: spec.RoleMappings[0].ExternalGroupName,
					Id:                &idpID,
					RoleAssignments: &[]admin.RoleAssignment{
						{
							GroupId: &projectId,
							Role:    &spec.RoleMappings[0].RoleAssignments[0].Role,
						},
					},
				},
			},
			UserConflicts: nil,
		}

		diff := deep.Equal(expected, result)
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
			DomainRestrictionEnabled: toptr.MakePtr(true),
			SSODebugEnabled:          toptr.MakePtr(true),
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
