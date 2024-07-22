package atlasproject

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestEnsureCustomRoles(t *testing.T) {
	testRole := []akov2.CustomRole{
		{
			Name: "test-role",
			InheritedRoles: []akov2.Role{
				{Name: "role1", Database: "db1"},
				{Name: "role2", Database: "db2"},
			},
			Actions: []akov2.Action{
				{
					Name: "action1",
					Resources: []akov2.Resource{
						{Cluster: pointer.MakePtr(true)},
						{Database: pointer.MakePtr("db1")},
					},
				},
				{
					Name: "action2",
					Resources: []akov2.Resource{
						{
							Database:   pointer.MakePtr("db2"),
							Collection: pointer.MakePtr("test-collection"),
						},
					},
				},
			},
		},
	}

	for _, tc := range []struct {
		name string

		roles []akov2.CustomRole

		roleAPI *mockadmin.CustomDatabaseRolesApi

		isOK bool
	}{
		{
			name: "No Roles in AKO or Atlas (no op)",
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDatabaseRoles(context.Background(), "").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDatabaseRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{},
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
		{
			name:  "Roles in AKO, but not Atlas (Create)",
			roles: testRole,
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDatabaseRoles(context.Background(), "").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDatabaseRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{},
						&http.Response{},
						nil,
					)
				roleAPI.EXPECT().CreateCustomDatabaseRole(context.Background(), "", mock.AnythingOfType("*admin.UserCustomDBRole")).
					Return(admin.CreateCustomDatabaseRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().CreateCustomDatabaseRoleExecute(mock.Anything).
					Return(
						&admin.UserCustomDBRole{},
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
		{
			name:  "Roles in AKO and in Atlas (Update)",
			roles: testRole,
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDatabaseRoles(context.Background(), "").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDatabaseRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{
							{
								RoleName: "test-role",
								InheritedRoles: &[]admin.DatabaseInheritedRole{
									{Role: "role3", Db: "db1"},
								},
								Actions: &[]admin.DatabasePrivilegeAction{
									{
										Action: "action1",
										Resources: &[]admin.DatabasePermittedNamespaceResource{
											{Db: "db2"},
										},
									},
								},
							},
						},
						&http.Response{},
						nil,
					)
				roleAPI.EXPECT().UpdateCustomDatabaseRole(context.Background(), "", "test-role", mock.AnythingOfType("*admin.UpdateCustomDBRole")).
					Return(admin.UpdateCustomDatabaseRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().UpdateCustomDatabaseRoleExecute(mock.Anything).
					Return(
						&admin.UserCustomDBRole{},
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
		{
			name: "Roles not in AKO but are in Atlas (Delete)",
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDatabaseRoles(context.Background(), "").
					Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDatabaseRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{
							{
								RoleName: "test-role",
								InheritedRoles: &[]admin.DatabaseInheritedRole{
									{Role: "role3", Db: "db1"},
								},
								Actions: &[]admin.DatabasePrivilegeAction{
									{
										Action: "action1",
										Resources: &[]admin.DatabasePermittedNamespaceResource{
											{Db: "db2"},
										},
									},
								},
							},
						},
						&http.Response{},
						nil,
					)
				roleAPI.EXPECT().DeleteCustomDatabaseRole(context.Background(), "", "test-role").
					Return(admin.DeleteCustomDatabaseRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().DeleteCustomDatabaseRoleExecute(mock.Anything).
					Return(
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				SdkClient: &admin.APIClient{
					CustomDatabaseRolesApi: tc.roleAPI,
				},
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			project := akov2.DefaultProject("test-namespace", "test-connnection")
			project.Spec.CustomRoles = tc.roles

			result := ensureCustomRoles(workflowCtx, project)

			assert.Equal(t, tc.isOK, result.IsOk())
		})
	}
}
