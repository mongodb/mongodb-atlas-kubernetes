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

package atlascustomrole

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
	"go.mongodb.org/atlas-sdk/v20250312011/mockadmin"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	mocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func Test_roleController_Reconcile(t *testing.T) {
	type fields struct {
		ctx                       *workflow.Context
		service                   func() customroles.CustomRoleService
		role                      *akov2.AtlasCustomRole
		deletionProtectionEnabled bool
		k8sObjects                []client.Object
	}
	tests := []struct {
		name   string
		fields fields
		want   workflow.DeprecatedResult
	}{
		{
			name: "Create custom role successfully",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{}, nil)
					s.EXPECT().Create(context.Background(), "testProjectID",
						mock.AnythingOfType("customroles.CustomRole")).Return(nil)
					return s
				},
				role: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testRole",
						Namespace: "testRoleNamespace",
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testRole",
							Namespace: "testRoleNamespace",
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name: "TestRoleName",
								InheritedRoles: []akov2.Role{
									{
										Name:     "read",
										Database: "main",
									},
								},
								Actions: []akov2.Action{
									{
										Name: "VIEW_ALL_HISTORY",
										Resources: []akov2.Resource{
											{
												Cluster:    pointer.MakePtr(true),
												Database:   pointer.MakePtr("main"),
												Collection: pointer.MakePtr("collection"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "Create custom role with error",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{}, nil)
					s.EXPECT().Create(context.Background(), "testProjectID",
						mock.AnythingOfType("customroles.CustomRole")).Return(fmt.Errorf("unable to create role"))
					return s
				},
				role: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotCreated, errors.New("unable to create role")),
		},
		{
			name: "Create custom role with error on Getting roles",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{}, fmt.Errorf("unable to Get roles"))
					return s
				},
				role: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.ProjectCustomRolesReady, errors.New("unable to Get roles")),
		},
		{
			name: "Update custom role successfully",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							// This has to be different from the one described below
							CustomRole: &akov2.CustomRole{
								Name:           "TestRoleName",
								InheritedRoles: nil,
								Actions:        nil,
							},
						}, nil)
					s.EXPECT().Update(context.Background(), "testProjectID", "TestRoleName",
						mock.AnythingOfType("customroles.CustomRole")).Return(nil)
					return s
				},
				role: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testRole",
						Namespace: "testRoleNamespace",
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testRole",
							Namespace: "testRoleNamespace",
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name: "TestRoleName",
								InheritedRoles: []akov2.Role{
									{
										Name:     "read",
										Database: "main",
									},
								},
								Actions: []akov2.Action{
									{
										Name: "VIEW_ALL_HISTORY",
										Resources: []akov2.Resource{
											{
												Cluster:    pointer.MakePtr(true),
												Database:   pointer.MakePtr("main"),
												Collection: pointer.MakePtr("collection"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "Update custom role with error",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							// This has to be different from the one described below
							CustomRole: &akov2.CustomRole{
								Name:           "TestRoleName",
								InheritedRoles: nil,
								Actions:        nil,
							},
						}, nil)
					s.EXPECT().Update(context.Background(), "testProjectID", "TestRoleName",
						mock.AnythingOfType("customroles.CustomRole")).Return(fmt.Errorf("unable to update custom role"))
					return s
				},
				role: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotUpdated, errors.New("unable to update custom role")),
		},
		{
			name: "Update custom role successfully no update",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							CustomRole: &akov2.CustomRole{
								Name: "TestRoleName",
								InheritedRoles: []akov2.Role{
									{
										Name:     "read",
										Database: "main",
									},
								},
								Actions: []akov2.Action{
									{
										Name: "VIEW_ALL_HISTORY",
										Resources: []akov2.Resource{
											{
												Cluster:    pointer.MakePtr(true),
												Database:   pointer.MakePtr("main"),
												Collection: pointer.MakePtr("collection"),
											},
										},
									},
								},
							},
						}, nil)
					return s
				},
				role: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "Delete custom role successfully",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							CustomRole: &akov2.CustomRole{
								Name:           "TestRoleName",
								InheritedRoles: nil,
								Actions:        nil,
							},
						}, nil)
					s.EXPECT().Delete(context.Background(), "testProjectID", "TestRoleName").Return(nil)
					return s
				},
				role: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "testRole",
						Namespace:         "testRoleNamespace",
						DeletionTimestamp: pointer.MakePtr(metav1.NewTime(time.Now())),
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testRole",
							Namespace: "testRoleNamespace",
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name: "TestRoleName",
								InheritedRoles: []akov2.Role{
									{
										Name:     "read",
										Database: "main",
									},
								},
								Actions: []akov2.Action{
									{
										Name: "VIEW_ALL_HISTORY",
										Resources: []akov2.Resource{
											{
												Cluster:    pointer.MakePtr(true),
												Database:   pointer.MakePtr("main"),
												Collection: pointer.MakePtr("collection"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: workflow.Deleted(),
		},
		{
			name: "DO NOT Delete custom role successfully with DeletionProtection enabled, just stop managing it",
			fields: fields{
				deletionProtectionEnabled: true,
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							CustomRole: &akov2.CustomRole{
								Name:           "TestRoleName",
								InheritedRoles: nil,
								Actions:        nil,
							},
						}, nil)
					return s
				},
				role: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "testRole",
						Namespace:         "testNamespace",
						DeletionTimestamp: pointer.MakePtr(metav1.NewTime(time.Now())),
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "testRole",
							Namespace:         "testNamespace",
							DeletionTimestamp: pointer.MakePtr(metav1.NewTime(time.Now())),
							Finalizers:        []string{customresource.FinalizerLabel},
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name: "TestRoleName",
								InheritedRoles: []akov2.Role{
									{
										Name:     "read",
										Database: "main",
									},
								},
								Actions: []akov2.Action{
									{
										Name: "VIEW_ALL_HISTORY",
										Resources: []akov2.Resource{
											{
												Cluster:    pointer.MakePtr(true),
												Database:   pointer.MakePtr("main"),
												Collection: pointer.MakePtr("collection"),
											},
										},
									},
								},
							},
							ProjectDualReference: akov2.ProjectDualReference{
								ExternalProjectRef: &akov2.ExternalProjectReference{
									ID: "testProjectID",
								},
							},
						},
						Status: status.AtlasCustomRoleStatus{},
					},
				},
			},
			want: workflow.Deleted(),
		},
		{
			name: "Delete custom role with error",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := mocks.NewCustomRoleServiceMock(t)
					s.EXPECT().Get(context.Background(), "testProjectID", "TestRoleName").
						Return(customroles.CustomRole{
							CustomRole: &akov2.CustomRole{
								Name:           "TestRoleName",
								InheritedRoles: nil,
								Actions:        nil,
							},
						}, nil)
					s.EXPECT().Delete(context.Background(), "testProjectID", "TestRoleName").
						Return(fmt.Errorf("unable to delete custom role"))
					return s
				},
				role: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						DeletionTimestamp: pointer.MakePtr(metav1.NewTime(time.Now())),
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name: "TestRoleName",
							InheritedRoles: []akov2.Role{
								{
									Name:     "read",
									Database: "main",
								},
							},
							Actions: []akov2.Action{
								{
									Name: "VIEW_ALL_HISTORY",
									Resources: []akov2.Resource{
										{
											Cluster:    pointer.MakePtr(true),
											Database:   pointer.MakePtr("main"),
											Collection: pointer.MakePtr("collection"),
										},
									},
								},
							},
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotDeleted, errors.New("unable to delete custom role")),
		},
	}
	for _, tt := range tests {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(tt.fields.k8sObjects...).
			Build()
		t.Run(tt.name, func(t *testing.T) {
			r := &roleController{
				ctx:       tt.fields.ctx,
				project:   &project.Project{ID: "testProjectID"},
				service:   tt.fields.service(),
				role:      tt.fields.role,
				dpEnabled: tt.fields.deletionProtectionEnabled,
				k8sClient: k8sClient,
			}
			if got := r.Reconcile(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args struct {
	ctx                       *workflow.Context
	akoCustomRole             *akov2.AtlasCustomRole
	deletionProtectionEnabled bool
	k8sObjects                []client.Object
}

func Test_handleCustomRole(t *testing.T) {
	tests := []struct {
		name       string
		args       args
		solveError error
		want       workflow.DeprecatedResult
	}{
		{
			name: "Create custom role successfully using external project ID",
			args: args{
				ctx: &workflow.Context{ //nolint:dupl
					Log:   zap.S(),
					OrgID: "",
					SdkClientSet: &atlas.ClientSet{
						SdkClient20250312011: &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								cdrAPI.EXPECT().GetCustomDbRole(context.Background(), "testProjectID", "testRole").
									Return(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().GetCustomDbRoleExecute(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI}).
									Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
								cdrAPI.EXPECT().CreateCustomDbRole(context.Background(), "testProjectID",
									mock.AnythingOfType("*admin.UserCustomDBRole")).
									Return(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().CreateCustomDbRoleExecute(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI}).
									Return(nil, nil, nil)
								return cdrAPI
							}(),
						},
					},
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testRole",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ExternalProjectRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testRole",
							Namespace: "testNamespace",
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name:           "testRole",
								InheritedRoles: nil,
								Actions:        nil,
							},
							ProjectDualReference: akov2.ProjectDualReference{
								ProjectRef: &common.ResourceRefNamespaced{
									Name:      "testProject",
									Namespace: "testNamespace",
								},
							},
						},
						Status: status.AtlasCustomRoleStatus{},
					},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "Create custom role successfully using external project reference",
			args: args{
				ctx: &workflow.Context{ //nolint:dupl
					Log:   zap.S(),
					OrgID: "",
					SdkClientSet: &atlas.ClientSet{
						SdkClient20250312011: &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								cdrAPI.EXPECT().GetCustomDbRole(context.Background(), "testProjectID", "testRole").
									Return(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().GetCustomDbRoleExecute(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI}).
									Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
								cdrAPI.EXPECT().CreateCustomDbRole(context.Background(), "testProjectID",
									mock.AnythingOfType("*admin.UserCustomDBRole")).
									Return(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().CreateCustomDbRoleExecute(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI}).
									Return(nil, nil, nil)
								return cdrAPI
							}(),
							ProjectsApi: func() admin.ProjectsApi {
								projectAPI := mockadmin.NewProjectsApi(t)
								projectAPI.EXPECT().GetGroupByName(mock.Anything, "testProject").
									Return(admin.GetGroupByNameApiRequest{ApiService: projectAPI})
								projectAPI.EXPECT().GetGroupByNameExecute(mock.Anything).
									Return(&admin.Group{Id: pointer.MakePtr("testProjectID")}, nil, nil)
								return projectAPI
							}(),
						},
					},
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testRole",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "testProject",
								Namespace: "testNamespace",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasProject{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testProject",
							Namespace: "testNamespace",
						},
						Spec: akov2.AtlasProjectSpec{
							Name: "testProject",
						},
						Status: status.AtlasProjectStatus{
							ID: "testProjectID",
						},
					},
					&akov2.AtlasCustomRole{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testRole",
							Namespace: "testNamespace",
						},
						Spec: akov2.AtlasCustomRoleSpec{
							Role: akov2.CustomRole{
								Name:           "testRole",
								InheritedRoles: nil,
								Actions:        nil,
							},
							ProjectDualReference: akov2.ProjectDualReference{
								ProjectRef: &common.ResourceRefNamespaced{
									Name:      "testProject",
									Namespace: "testNamespace",
								},
							},
						},
						Status: status.AtlasCustomRoleStatus{},
					},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "DO NOT create custom role if external project cannot be found",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClientSet: &atlas.ClientSet{
						SdkClient20250312011: &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								return cdrAPI
							}(),
							ProjectsApi: func() admin.ProjectsApi {
								notFound := &admin.GenericOpenAPIError{}
								notFound.SetModel(admin.ApiError{ErrorCode: "RESOURCE_NOT_FOUND"})

								projectAPI := mockadmin.NewProjectsApi(t)
								projectAPI.EXPECT().GetGroupByName(mock.Anything, "testProject").
									Return(admin.GetGroupByNameApiRequest{ApiService: projectAPI})
								projectAPI.EXPECT().GetGroupByNameExecute(mock.Anything).
									Return(nil, nil, notFound)
								return projectAPI
							}(),
						},
					},
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "testProject",
								Namespace: "testNamespace",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{
					&akov2.AtlasProject{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testProject",
							Namespace: "testNamespace",
						},
						Spec: akov2.AtlasProjectSpec{
							Name: "testProject",
						},
						Status: status.AtlasProjectStatus{
							ID: "",
						},
					},
				},
			},
			solveError: translation.ErrNotFound,
		},
		{
			name: "DO NOT create custom role if external project reference doesn't exist",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClientSet: &atlas.ClientSet{
						SdkClient20250312011: &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								return cdrAPI
							}(),
						},
					},
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectDualReference: akov2.ProjectDualReference{
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "testProject",
								Namespace: "testNamespace",
							},
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{},
			},
			solveError: fmt.Errorf("atlasprojects.atlas.mongodb.com \"testProject\" not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.args.k8sObjects...).
				Build()
			service := customroles.NewCustomRoles(tt.args.ctx.SdkClientSet.SdkClient20250312011.CustomDatabaseRolesApi)
			r := AtlasCustomRoleReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{Client: k8sClient},
			}
			prj, err := solveProjectID(t, &r, tt.args)
			if tt.solveError == nil {
				require.NoError(t, err)
				if got := handleCustomRole(tt.args.ctx, k8sClient, prj, service, tt.args.akoCustomRole, tt.args.deletionProtectionEnabled); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("handleCustomRole() = %v, want %v", got, tt.want)
				}
			} else {
				assert.ErrorContains(t, err, tt.solveError.Error())
			}
		})
	}
}

func solveProjectID(t *testing.T, r *AtlasCustomRoleReconciler, args args) (*project.Project, error) {
	t.Helper()
	if args.akoCustomRole.Spec.ProjectDualReference.ExternalProjectRef != nil {
		return &project.Project{ID: args.akoCustomRole.Spec.ExternalProjectRef.ID}, nil
	}
	return r.ResolveProject(args.ctx.Context, args.ctx.SdkClientSet.SdkClient20250312011, args.akoCustomRole)
}
