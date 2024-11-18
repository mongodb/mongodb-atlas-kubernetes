package atlascustomrole

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
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
		want   workflow.Result
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotCreated, "unable to create role"),
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.ProjectCustomRolesReady, "unable to Get roles"),
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotUpdated, "unable to update custom role"),
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
							ExternalProjectIDRef: &akov2.ExternalProjectReference{
								ID: "testProjectID",
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
					s := translation.NewCustomRoleServiceMock(t)
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotDeleted, "unable to delete custom role"),
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

func Test_handleCustomRole(t *testing.T) {
	type args struct {
		ctx                       *workflow.Context
		akoCustomRole             *akov2.AtlasCustomRole
		deletionProtectionEnabled bool
		k8sObjects                []client.Object
	}
	tests := []struct {
		name string
		args args
		want workflow.Result
	}{
		{
			name: "Create custom role successfully using external project ID",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClient: func() *admin.APIClient {
						return &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								cdrAPI.EXPECT().GetCustomDatabaseRole(context.Background(), "testProjectID", "testRole").
									Return(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().GetCustomDatabaseRoleExecute(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
									Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
								cdrAPI.EXPECT().CreateCustomDatabaseRole(context.Background(), "testProjectID",
									mock.AnythingOfType("*admin.UserCustomDBRole")).
									Return(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().CreateCustomDatabaseRoleExecute(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
									Return(nil, nil, nil)
								return cdrAPI
							}(),
						}
					}(),
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
						ExternalProjectIDRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
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
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "testProject",
								Namespace: "testNamespace",
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
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClient: func() *admin.APIClient {
						return &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								cdrAPI.EXPECT().GetCustomDatabaseRole(context.Background(), "testProjectID", "testRole").
									Return(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().GetCustomDatabaseRoleExecute(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
									Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
								cdrAPI.EXPECT().CreateCustomDatabaseRole(context.Background(), "testProjectID",
									mock.AnythingOfType("*admin.UserCustomDBRole")).
									Return(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().CreateCustomDatabaseRoleExecute(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
									Return(nil, nil, nil)
								return cdrAPI
							}(),
						}
					}(),
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
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "testProject",
							Namespace: "testNamespace",
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
						Spec: akov2.AtlasProjectSpec{},
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
							ProjectRef: &common.ResourceRefNamespaced{
								Name:      "testProject",
								Namespace: "testNamespace",
							},
						},
						Status: status.AtlasCustomRoleStatus{},
					},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "DO NOT create custom role if external project reference is empty",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClient: func() *admin.APIClient {
						return &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								return cdrAPI
							}(),
						}
					}(),
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "testProject",
							Namespace: "testNamespace",
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
						Spec: akov2.AtlasProjectSpec{},
						Status: status.AtlasProjectStatus{
							ID: "",
						},
					},
				},
			},
			want: workflow.Terminate(workflow.ProjectCustomRolesReady, "the referenced AtlasProject resource 'testProject' doesn't have ID (status.ID is empty)"),
		},
		{
			name: "DO NOT create custom role if external project reference doesn't exist",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClient: func() *admin.APIClient {
						return &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								return cdrAPI
							}(),
						}
					}(),
					Context: context.Background(),
				},
				akoCustomRole: &akov2.AtlasCustomRole{
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "testProject",
							Namespace: "testNamespace",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
				k8sObjects: []client.Object{},
			},
			want: workflow.Terminate(workflow.ProjectCustomRolesReady, "atlasprojects.atlas.mongodb.com \"testProject\" not found"),
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
			if got := handleCustomRole(tt.args.ctx, k8sClient, tt.args.akoCustomRole, tt.args.deletionProtectionEnabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleCustomRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
