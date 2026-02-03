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

package atlasproject

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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

		projectAnnotations map[string]string
	}{
		{
			name: "No Roles in AKO or Atlas (no op)",
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
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
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{},
						&http.Response{},
						nil,
					)
				roleAPI.EXPECT().CreateCustomDbRole(context.Background(), "", mock.AnythingOfType("*admin.UserCustomDBRole")).
					Return(admin.CreateCustomDbRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().CreateCustomDbRoleExecute(mock.Anything).
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
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
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
				roleAPI.EXPECT().UpdateCustomDbRole(context.Background(), "", "test-role", mock.AnythingOfType("*admin.UpdateCustomDBRole")).
					Return(admin.UpdateCustomDbRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().UpdateCustomDbRoleExecute(mock.Anything).
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
			projectAnnotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: func() string {
					d, _ := json.Marshal(&akov2.AtlasProjectSpec{
						CustomRoles: []akov2.CustomRole{
							{
								Name: "test-role-1",
								InheritedRoles: []akov2.Role{
									{Name: "role3", Database: "db1"},
								},
								Actions: []akov2.Action{
									{
										Name: "action1",
										Resources: []akov2.Resource{
											{Database: pointer.MakePtr("db2")},
										},
									},
								},
							},
						},
					})
					return string(d)
				}(),
			},
			name: "Roles not in AKO but are in Atlas (Delete) if there were previous in AKO. Remove only those that were in AKO",
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
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
							{
								RoleName: "test-role-1",
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
				roleAPI.EXPECT().DeleteCustomDbRole(context.Background(), "", "test-role-1").
					Return(admin.DeleteCustomDbRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().DeleteCustomDbRoleExecute(mock.Anything).
					Return(
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
		{
			name: "Roles in AKO and in last applied config. Delete only those that were deleted from the spec",
			projectAnnotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: func() string {
					d, _ := json.Marshal(&akov2.AtlasProjectSpec{
						CustomRoles: []akov2.CustomRole{
							{
								Name: "test-role",
								InheritedRoles: []akov2.Role{
									{Name: "role", Database: "db"},
								},
								Actions: []akov2.Action{
									{
										Name: "action",
										Resources: []akov2.Resource{
											{
												Database:   pointer.MakePtr("db"),
												Cluster:    pointer.MakePtr(true),
												Collection: pointer.MakePtr("test-collection"),
											},
										},
									},
								},
							},
							{
								Name: "test-role-2",
								InheritedRoles: []akov2.Role{
									{Name: "role2", Database: "db2"},
								},
								Actions: []akov2.Action{
									{
										Name: "action2",
										Resources: []akov2.Resource{
											{Database: pointer.MakePtr("db2")},
										},
									},
								},
							},
						},
					})
					return string(d)
				}(),
			},
			roles: []akov2.CustomRole{
				{
					Name: "test-role",
					InheritedRoles: []akov2.Role{
						{Name: "role", Database: "db"},
					},
					Actions: []akov2.Action{
						{
							Name: "action",
							Resources: []akov2.Resource{
								{
									Database:   pointer.MakePtr("db"),
									Cluster:    pointer.MakePtr(true),
									Collection: pointer.MakePtr("test-collection"),
								},
							},
						},
					},
				},
			},
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
					Return(
						[]admin.UserCustomDBRole{
							{
								RoleName: "test-role",
								InheritedRoles: &[]admin.DatabaseInheritedRole{
									{Role: "role", Db: "db"},
								},
								Actions: &[]admin.DatabasePrivilegeAction{
									{
										Action: "action",
										Resources: &[]admin.DatabasePermittedNamespaceResource{
											{
												Db:         "db",
												Collection: "test-collection",
												Cluster:    true,
											},
										},
									},
								},
							},
							{
								RoleName: "test-role-1",
								InheritedRoles: &[]admin.DatabaseInheritedRole{
									{Role: "role1", Db: "db1"},
								},
								Actions: &[]admin.DatabasePrivilegeAction{
									{
										Action: "action1",
										Resources: &[]admin.DatabasePermittedNamespaceResource{
											{Db: "db1"},
										},
									},
								},
							},
							{
								RoleName: "test-role-2",
								InheritedRoles: &[]admin.DatabaseInheritedRole{
									{Role: "role2", Db: "db2"},
								},
								Actions: &[]admin.DatabasePrivilegeAction{
									{
										Action: "action2",
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
				roleAPI.EXPECT().DeleteCustomDbRole(context.Background(), "", "test-role-2").
					Return(admin.DeleteCustomDbRoleApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().DeleteCustomDbRoleExecute(mock.Anything).
					Return(
						&http.Response{},
						nil,
					)
				return roleAPI
			}(),
			isOK: true,
		},
		{
			name: "Roles not in AKO but are in Atlas (Do not Delete) and NO previous in AKO",
			roleAPI: func() *mockadmin.CustomDatabaseRolesApi {
				roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
				roleAPI.EXPECT().ListCustomDbRoles(context.Background(), "").
					Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI})
				roleAPI.EXPECT().ListCustomDbRolesExecute(mock.Anything).
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
				return roleAPI
			}(),
			isOK: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						CustomDatabaseRolesApi: tc.roleAPI,
					},
				},
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			project := akov2.DefaultProject("test-namespace", "test-connnection")
			project.Spec.CustomRoles = tc.roles
			project.Annotations = tc.projectAnnotations

			result := ensureCustomRoles(workflowCtx, project)

			assert.Equal(t, tc.isOK, result.IsOk())
		})
	}
}

func TestCustomRolesNonGreedyBehaviour(t *testing.T) {
	for _, tc := range []struct {
		title                  string
		lastAppliedCustomRoles []string
		specCustomRoles        []string
		atlasCustomRoles       []string
		wantRemoved            []string
	}{
		{
			title:                  "no last applied no removal in Atlas",
			lastAppliedCustomRoles: []string{},
			specCustomRoles:        []string{},
			atlasCustomRoles:       []string{"cr1", "cr2"},
			wantRemoved:            []string{},
		},
		{
			title:                  "removed from last applied removes from Atlas",
			lastAppliedCustomRoles: []string{"cr1", "cr2"},
			specCustomRoles:        []string{"cr1"},
			atlasCustomRoles:       []string{"cr1", "cr2"},
			wantRemoved:            []string{"cr2"},
		},
		{
			title:                  "removed all from last applied removes all from Atlas",
			lastAppliedCustomRoles: []string{"cr1", "cr2"},
			specCustomRoles:        []string{},
			atlasCustomRoles:       []string{"cr1", "cr2"},
			wantRemoved:            []string{"cr1", "cr2"},
		},
		{
			title:                  "not in last applied not removed from Atlas",
			lastAppliedCustomRoles: []string{"cr1"},
			specCustomRoles:        []string{"cr1"},
			atlasCustomRoles:       []string{"cr1", "cr2"},
			wantRemoved:            []string{},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			prj := newCustomRolesTestProject(tc.specCustomRoles)
			lastPrj := newCustomRolesTestProject(tc.lastAppliedCustomRoles)
			prj.Annotations[customresource.AnnotationLastAppliedConfiguration] = jsonize(t, lastPrj.Spec)

			roleAPI := mockadmin.NewCustomDatabaseRolesApi(t)
			roleAPI.EXPECT().ListCustomDbRoles(mock.Anything, mock.Anything).
				Return(admin.ListCustomDbRolesApiRequest{ApiService: roleAPI}).Once()
			roleAPI.EXPECT().ListCustomDbRolesExecute(
				mock.AnythingOfType("admin.ListCustomDbRolesApiRequest")).Return(
				synthesizeAtlasCustomRoles(tc.atlasCustomRoles), nil, nil,
			).Once()

			removals := len(tc.wantRemoved)
			if removals > 0 {
				roleAPI.EXPECT().DeleteCustomDbRole(
					mock.Anything, mock.Anything, mock.Anything,
				).Return(admin.DeleteCustomDbRoleApiRequest{ApiService: roleAPI}).Times(removals)
				roleAPI.EXPECT().DeleteCustomDbRoleExecute(
					mock.AnythingOfType("admin.DeleteCustomDbRoleApiRequest")).Return(
					nil, nil,
				).Times(removals)
			}

			workflowCtx := workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						CustomDatabaseRolesApi: roleAPI,
					},
				},
			}

			result := ensureCustomRoles(&workflowCtx, prj)
			require.Equal(t, workflow.OK(), result)
		})
	}
}

func newCustomRolesTestProject(customRoles []string) *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name:        "test-project",
			CustomRoles: synthesizeCustomRoles(customRoles),
		},
	}
}

func synthesizeCustomRoles(customRoles []string) []akov2.CustomRole {
	crs := make([]akov2.CustomRole, 0, len(customRoles))
	for _, name := range customRoles {
		crs = append(crs, akov2.CustomRole{
			Name:           name,
			InheritedRoles: []akov2.Role{},
			Actions:        []akov2.Action{},
		})
	}
	return crs
}

func synthesizeAtlasCustomRoles(customRoles []string) []admin.UserCustomDBRole {
	atlasRoles := make([]admin.UserCustomDBRole, 0, len(customRoles))
	for _, name := range customRoles {
		atlasRoles = append(atlasRoles, admin.UserCustomDBRole{
			RoleName:       name,
			Actions:        &[]admin.DatabasePrivilegeAction{},
			InheritedRoles: &[]admin.DatabaseInheritedRole{},
		})
	}
	return atlasRoles
}
