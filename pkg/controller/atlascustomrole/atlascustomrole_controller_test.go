package atlascustomrole

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

func TestAtlasCustomRoleReconciler_Reconcile(t *testing.T) {
	tests := map[string]struct {
		akoCustomRole  *akov2.AtlasCustomRole
		interceptors   interceptor.Funcs
		expected       ctrl.Result
		isSupported    bool
		sdkShouldError bool
	}{
		"failed to retrieve custom role": {
			isSupported: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
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
			interceptors: interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("failed to get custom role")
				},
			},
			expected: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		},
		"custom role is not found": {
			isSupported: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role-not-found",
					Namespace: "default",
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
			expected: ctrl.Result{},
		},
		"custom role has invalid version": {
			isSupported: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "9.0.0",
					},
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
			expected: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		},
		"custom role resource unsupported": {
			isSupported: false,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
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
			expected: ctrl.Result{RequeueAfter: 0},
		},
		"custom role has skip annotation": {
			isSupported: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
					Annotations: map[string]string{
						"mongodb.com/atlas-reconciliation-policy": "skip",
					},
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
			expected: ctrl.Result{},
		},
		"custom role processing": {
			isSupported: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{Name: "test"},
					},
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
			expected: ctrl.Result{},
		},
		"custom role failed to create sdk client": {
			isSupported:    true,
			sdkShouldError: true,
			akoCustomRole: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "role",
					Namespace: "default",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{Name: "test"},
					},
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
			expected: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		},
	}
	version.Version = "1.0.0"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.akoCustomRole).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			r := &AtlasCustomRoleReconciler{
				Client:        k8sClient,
				Log:           zap.S(),
				Scheme:        testScheme,
				EventRecorder: record.NewFakeRecorder(10),
				AtlasProvider: &atlasmocks.TestProvider{
					SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
						if tt.sdkShouldError {
							return nil, "", fmt.Errorf("failed to create sdk")
						}
						cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
						cdrAPI.EXPECT().GetCustomDatabaseRole(context.Background(), "testProjectID", "TestRoleName").
							Return(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
						cdrAPI.EXPECT().GetCustomDatabaseRoleExecute(admin.GetCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
							Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
						cdrAPI.EXPECT().CreateCustomDatabaseRole(context.Background(), "testProjectID",
							mock.AnythingOfType("*admin.UserCustomDBRole")).
							Return(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI})
						cdrAPI.EXPECT().CreateCustomDatabaseRoleExecute(admin.CreateCustomDatabaseRoleApiRequest{ApiService: cdrAPI}).
							Return(nil, nil, nil)
						return &admin.APIClient{
							CustomDatabaseRolesApi: cdrAPI,
						}, "", nil
					},
					IsCloudGovFunc: func() bool {
						return false
					},
					IsSupportedFunc: func() bool {
						return tt.isSupported
					},
				},
			}

			result, err := r.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "role",
					Namespace: "default",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtlasCustomRoleReconciler_selectCredentials(t *testing.T) {
	tests := map[string]struct {
		k8sObjects []client.Object
		resultErr  bool
		result     *client.ObjectKey
		role       *akov2.AtlasCustomRole
	}{
		"Should select CustomRoleCredentials from CustomRole when externalProjectID is set": {
			k8sObjects: nil,
			resultErr:  false,
			result:     &client.ObjectKey{Name: "credentials", Namespace: "testNamespace"},
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ExternalProjectIDRef:  &akov2.ExternalProjectReference{ID: "testProjectID"},
					LocalCredentialHolder: api.LocalCredentialHolder{ConnectionSecret: &api.LocalObjectReference{Name: "credentials"}},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should NOT select CustomRoleCredentials from CustomRole when externalProjectID is set with empty local credentials": {
			k8sObjects: nil,
			resultErr:  true,
			result:     nil,
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ExternalProjectIDRef:  &akov2.ExternalProjectReference{ID: "testProjectID"},
					LocalCredentialHolder: api.LocalCredentialHolder{ConnectionSecret: nil},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should select CustomRoleCredentials when projectRef is set with LocalCredentials": {
			k8sObjects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "testCredentials",
							Namespace: "testNamespace",
						},
					},
					Status: status.AtlasProjectStatus{},
				},
			},
			resultErr: false,
			result:    &client.ObjectKey{Name: "credentials", Namespace: "testNamespace"},
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectRef:            &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
					LocalCredentialHolder: api.LocalCredentialHolder{ConnectionSecret: &api.LocalObjectReference{Name: "credentials"}},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should select AtlasProject credentials when projectRef is set but LocalCredentials are empty": {
			k8sObjects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "testCredentials",
							Namespace: "testNamespace",
						},
					},
					Status: status.AtlasProjectStatus{},
				},
				&v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testCredentials",
						Namespace: "testNamespace",
					},
				},
			},
			resultErr: false,
			result:    &client.ObjectKey{Name: "testCredentials", Namespace: "testNamespace"},
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should NOT select AtlasProject credentials when projectRef is set but AtlasProject credentials are not configured": {
			k8sObjects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec:   akov2.AtlasProjectSpec{},
					Status: status.AtlasProjectStatus{},
				},
				&v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testCredentials",
						Namespace: "testNamespace",
					},
				},
			},
			resultErr: true,
			result:    nil,
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should NOT select AtlasProject credentials when projectRef is set and AtlasProject doesn't exist": {
			k8sObjects: []client.Object{},
			resultErr:  true,
			result:     nil,
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
		"Should NOT select credentials when both projectRef and externalProjectId are empty": {
			k8sObjects: []client.Object{},
			resultErr:  true,
			result:     nil,
			role: &akov2.AtlasCustomRole{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec:   akov2.AtlasCustomRoleSpec{},
				Status: status.AtlasCustomRoleStatus{},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, v1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.k8sObjects...).Build()

			result, err := selectCredentials(context.Background(), k8sClient, tt.role)
			if err != nil && !tt.resultErr {
				t.Fatalf("Unexpected error: %v", err)
			}
			assert.Equal(t, result, tt.result)
		})
	}
}
