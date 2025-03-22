package atlasipaccesslist

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestReconcile(t *testing.T) {
	tests := map[string]struct {
		request        ctrl.Request
		provider       atlas.Provider
		expectedResult ctrl.Result
	}{
		"should fail to prepare source reconciliation": {
			request:        ctrl.Request{NamespacedName: types.NamespacedName{Name: "wrong", Namespace: "default"}},
			expectedResult: ctrl.Result{},
		},
		"should handle ip access list": {
			request: ctrl.Request{NamespacedName: types.NamespacedName{Name: "ip-access-list", Namespace: "default"}},
			//nolint:dupl
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					ialAPI := mockadmin.NewProjectIPAccessListApi(t)
					ialAPI.EXPECT().ListProjectIpAccessLists(mock.Anything, "123").
						Return(admin.ListProjectIpAccessListsApiRequest{ApiService: ialAPI})
					ialAPI.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
						Return(
							&admin.PaginatedNetworkAccess{
								Results: &[]admin.NetworkPermissionEntry{
									{
										CidrBlock: pointer.MakePtr("192.168.0.0/24"),
									},
								},
							},
							nil,
							nil,
						)
					ialAPI.EXPECT().GetProjectIpAccessListStatus(mock.Anything, "123", "192.168.0.0/24").
						Return(admin.GetProjectIpAccessListStatusApiRequest{ApiService: ialAPI})
					ialAPI.EXPECT().GetProjectIpAccessListStatusExecute(mock.AnythingOfType("admin.GetProjectIpAccessListStatusApiRequest")).
						Return(
							&admin.NetworkPermissionEntryStatus{STATUS: "ACTIVE"},
							nil,
							nil,
						)

					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProjectByName(mock.Anything, "my-project").
						Return(admin.GetProjectByNameApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectByNameExecute(mock.Anything).
						Return(&admin.Group{Id: pointer.MakePtr("123")}, nil, nil)

					return &atlas.ClientSet{
						SdkClient20231115008: &admin.APIClient{
							ProjectIPAccessListApi: ialAPI,
							ProjectsApi:            projectAPI,
						},
						SdkClient20241113001: &adminv20241113001.APIClient{},
					}, nil
				},
			},
			expectedResult: ctrl.Result{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			}
			ipAccessList := &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			require.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, ipAccessList, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("orgId"),
						"publicApiKey":  []byte("publicApiKey"),
						"privateApiKey": []byte("privateApiKey"),
					},
				}).
				WithStatusSubresource(ipAccessList).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasIPAccessListReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
				AtlasProvider: tt.provider,
				EventRecorder: record.NewFakeRecorder(10),
			}
			result, err := r.Reconcile(ctx.Context, tt.request)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
