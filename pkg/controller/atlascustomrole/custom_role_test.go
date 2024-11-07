package atlascustomrole

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"

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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{}, nil)
					s.EXPECT().Create(context.Background(), "testProjectID",
						mock.AnythingOfType("customroles.CustomRole")).Return(nil)
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{}, nil)
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.AtlasCustomRoleNotCreated, "unable to create role"),
		},
		{
			name: "Create custom role with error on listing roles",
			fields: fields{
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := translation.NewCustomRoleServiceMock(t)
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{}, fmt.Errorf("unable to list roles"))
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.Terminate(workflow.ProjectCustomRolesReady, "unable to list roles"),
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
								// This has to be different from the one described below
								CustomRole: &akov2.CustomRole{
									Name:           "TestRoleName",
									InheritedRoles: nil,
									Actions:        nil,
								},
							},
						}, nil)
					s.EXPECT().Update(context.Background(), "testProjectID", "TestRoleName",
						mock.AnythingOfType("customroles.CustomRole")).Return(nil)
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
								// This has to be different from the one described below
								CustomRole: &akov2.CustomRole{
									Name:           "TestRoleName",
									InheritedRoles: nil,
									Actions:        nil,
								},
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
						ProjectIDRef: akov2.ExternalProjectReference{
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
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
						ProjectIDRef: akov2.ExternalProjectReference{
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
								CustomRole: &akov2.CustomRole{
									Name:           "TestRoleName",
									InheritedRoles: nil,
									Actions:        nil,
								},
							},
						}, nil)
					s.EXPECT().Delete(context.Background(), "testProjectID", "TestRoleName").Return(nil)
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.OK(),
		},
		{
			name: "DO NOT Delete custom role successfully with DeletionProtection enabled",
			fields: fields{
				deletionProtectionEnabled: true,
				ctx: &workflow.Context{
					Log:     zap.S(),
					OrgID:   "",
					Context: context.Background(),
				},
				service: func() customroles.CustomRoleService {
					s := translation.NewCustomRoleServiceMock(t)
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
								CustomRole: &akov2.CustomRole{
									Name:           "TestRoleName",
									InheritedRoles: nil,
									Actions:        nil,
								},
							},
						}, nil)
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
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.OK(),
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
					s.EXPECT().List(context.Background(), "testProjectID").
						Return([]customroles.CustomRole{
							{
								CustomRole: &akov2.CustomRole{
									Name:           "TestRoleName",
									InheritedRoles: nil,
									Actions:        nil,
								},
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
						ProjectIDRef: akov2.ExternalProjectReference{
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
		t.Run(tt.name, func(t *testing.T) {
			r := &roleController{
				ctx:       tt.fields.ctx,
				service:   tt.fields.service(),
				role:      tt.fields.role,
				dpEnabled: tt.fields.deletionProtectionEnabled,
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
	}
	tests := []struct {
		name string
		args args
		want workflow.Result
	}{
		{
			name: "Create custom role successfully",
			args: args{
				ctx: &workflow.Context{
					Log:   zap.S(),
					OrgID: "",
					SdkClient: func() *admin.APIClient {
						return &admin.APIClient{
							CustomDatabaseRolesApi: func() admin.CustomDatabaseRolesApi {
								cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
								cdrAPI.EXPECT().ListCustomDatabaseRoles(context.Background(), "testProjectID").
									Return(admin.ListCustomDatabaseRolesApiRequest{ApiService: cdrAPI})
								cdrAPI.EXPECT().ListCustomDatabaseRolesExecute(admin.ListCustomDatabaseRolesApiRequest{ApiService: cdrAPI}).
									Return([]admin.UserCustomDBRole{}, nil, nil)
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
					Spec: akov2.AtlasCustomRoleSpec{
						Role: akov2.CustomRole{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions:        nil,
						},
						ProjectIDRef: akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
					Status: status.AtlasCustomRoleStatus{},
				},
			},
			want: workflow.OK(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleCustomRole(tt.args.ctx, tt.args.akoCustomRole, tt.args.deletionProtectionEnabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleCustomRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
