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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

func TestAtlasCustomRoleReconciler_Reconcile(t *testing.T) {
	tests := map[string]struct {
		akoCustomRole  *akov2.AtlasCustomRole
		interceptors   interceptor.Funcs
		expected       ctrl.Result
		isSupported    bool
		sdkShouldError bool
		wantErr        bool
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
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
			interceptors: interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("failed to get custom role")
				},
			},
			wantErr: true,
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
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
					},
				},
				Status: status.AtlasCustomRoleStatus{},
			},
			wantErr: true,
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
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "test"},
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "test"},
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "testProjectID",
						},
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
				},
				Status: status.AtlasCustomRoleStatus{},
			},
			wantErr: true,
		},
	}
	version.Version = "1.0.0"
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.akoCustomRole, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("orgId"),
						"publicApiKey":  []byte("publicApiKey"),
						"privateApiKey": []byte("privateApiKey"),
					},
				}).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			r := &AtlasCustomRoleReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    zap.S(),
					AtlasProvider: &atlasmocks.TestProvider{
						SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
							if tt.sdkShouldError {
								return nil, fmt.Errorf("failed to create sdk")
							}
							cdrAPI := mockadmin.NewCustomDatabaseRolesApi(t)
							cdrAPI.EXPECT().GetCustomDbRole(mock.Anything, "testProjectID", "TestRoleName").
								Return(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI})
							cdrAPI.EXPECT().GetCustomDbRoleExecute(admin.GetCustomDbRoleApiRequest{ApiService: cdrAPI}).
								Return(&admin.UserCustomDBRole{}, &http.Response{StatusCode: http.StatusNotFound}, nil)
							cdrAPI.EXPECT().CreateCustomDbRole(mock.Anything, "testProjectID",
								mock.AnythingOfType("*admin.UserCustomDBRole")).
								Return(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI})
							cdrAPI.EXPECT().CreateCustomDbRoleExecute(admin.CreateCustomDbRoleApiRequest{ApiService: cdrAPI}).
								Return(nil, nil, nil)

							pAPI := mockadmin.NewProjectsApi(t)
							if tt.akoCustomRole.Spec.ExternalProjectRef != nil {
								grp := &admin.Group{
									Id:   &tt.akoCustomRole.Spec.ExternalProjectRef.ID,
									Name: tt.akoCustomRole.Spec.ExternalProjectRef.ID,
								}
								pAPI.EXPECT().GetGroup(context.Background(), tt.akoCustomRole.Spec.ExternalProjectRef.ID).
									Return(admin.GetGroupApiRequest{ApiService: pAPI})
								pAPI.EXPECT().GetGroupExecute(admin.GetGroupApiRequest{ApiService: pAPI}).
									Return(grp, nil, nil)
							}
							return &atlas.ClientSet{SdkClient20250312012: &admin.APIClient{
								CustomDatabaseRolesApi: cdrAPI,
								ProjectsApi:            pAPI,
							}}, nil
						},
						IsCloudGovFunc: func() bool {
							return false
						},
						IsSupportedFunc: func() bool {
							return tt.isSupported
						},
					},
				},
				Scheme:        testScheme,
				EventRecorder: record.NewFakeRecorder(10),
			}

			result, err := r.Reconcile(context.Background(), ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "role",
					Namespace: "default",
				},
			})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}
