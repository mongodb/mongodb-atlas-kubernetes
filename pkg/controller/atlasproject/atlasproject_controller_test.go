package atlasproject

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestReconcile(t *testing.T) {
	tests := map[string]struct {
		atlasClientMocker func() *mongodbatlas.Client
		atlasSDKMocker    func() *admin.APIClient
		interceptors      interceptor.Funcs
		project           *akov2.AtlasProject
		result            reconcile.Result
		conditions        []api.Condition
		finalizers        []string
	}{
		"should handle project": {
			atlasClientMocker: func() *mongodbatlas.Client {
				return nil
			},
			atlasSDKMocker: func() *admin.APIClient {
				notFoundErr := &admin.GenericOpenAPIError{}
				notFoundErr.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr("NOT_IN_GROUP")})
				projectsAPI := mockadmin.NewProjectsApi(t)
				projectsAPI.EXPECT().GetProjectByName(context.Background(), "my-project").
					Return(admin.GetProjectByNameApiRequest{ApiService: projectsAPI})
				projectsAPI.EXPECT().GetProjectByNameExecute(mock.AnythingOfType("admin.GetProjectByNameApiRequest")).
					Return(
						nil,
						&http.Response{},
						notFoundErr,
					)
				projectsAPI.EXPECT().CreateProject(context.Background(), mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateProjectApiRequest{ApiService: projectsAPI})
				projectsAPI.EXPECT().CreateProjectExecute(mock.AnythingOfType("admin.CreateProjectApiRequest")).
					Return(
						&admin.Group{
							OrgId:                     "my-org-id",
							Id:                        pointer.MakePtr("my-project-id"),
							Name:                      "my-project",
							ClusterCount:              0,
							RegionUsageRestrictions:   pointer.MakePtr("NONE"),
							WithDefaultAlertsSettings: pointer.MakePtr(true),
							Tags: &[]admin.ResourceTag{
								{
									Key:   "test",
									Value: "AKO",
								},
							},
						},
						&http.Response{},
						nil,
					)

				return &admin.APIClient{
					ProjectsApi: projectsAPI,
				}
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			},
			result: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			conditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.ValidationSucceeded),
				api.FalseCondition(api.ProjectReadyType).
					WithReason(string(workflow.ProjectBeingConfiguredInAtlas)).
					WithMessageRegexp("configuring project in Atlas"),
			},
			finalizers: []string{customresource.FinalizerLabel},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			require.NoError(t, corev1.AddToScheme(testScheme))
			instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.project).
				WithStatusSubresource(tt.project).
				WithIndex(
					instancesIndexer.Object(),
					instancesIndexer.Name(),
					instancesIndexer.Keys,
				).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			reconciler := &AtlasProjectReconciler{
				Client:        k8sClient,
				Log:           logger,
				EventRecorder: record.NewFakeRecorder(30),
				AtlasProvider: &atlas.TestProvider{
					IsCloudGovFunc: func() bool {
						return false
					},
					IsSupportedFunc: func() bool {
						return true
					},
					ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
						return tt.atlasClientMocker(), "", nil
					},
					SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
						return tt.atlasSDKMocker(), "", nil
					},
				},
			}

			result, err := reconciler.Reconcile(context.Background(), reconcile.Request{NamespacedName: types.NamespacedName{Name: "my-project", Namespace: "default"}})
			atlasProject := akov2.AtlasProject{}
			require.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(tt.project), &atlasProject))
			require.NoError(t, err)
			assert.Equal(t, tt.result, result)
			assert.True(
				t,
				cmp.Equal(
					tt.conditions,
					atlasProject.Status.GetConditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
			assert.Equal(t, tt.finalizers, atlasProject.Finalizers)
		})
	}
}

func TestFindProjectsForBCP(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &akov2.AtlasBackupCompliancePolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-bcp", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "test-bcp", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns247"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "test-bcp", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "different-test-bcp", Namespace: "ns1"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns2"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns247"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByBackupCompliancePolicyIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[akov2.AtlasBackupCompliancePolicy](indexer.AtlasProjectByBackupCompliancePolicyIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}

func TestFindProjectsForConnectionSecret(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret"},
						AlertConfigurations: []akov2.AlertConfiguration{
							{
								Notifications: []akov2.Notification{
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}}, // double entry
									{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{OpsGenieAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{ServiceKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "test-secret"}},
								},
							},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns3"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-4", Namespace: "ns4"},
					Spec: akov2.AtlasProjectSpec{
						AlertConfigurations: []akov2.AlertConfiguration{
							{
								Notifications: []akov2.Notification{
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}},
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}}, // double entry
									{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}},
									{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "test-secret2", Namespace: "ns1"}},
								},
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns1"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns2"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-4", Namespace: "ns4"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByConnectionSecretIndexer(zaptest.NewLogger(t))

			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[corev1.Secret](indexer.AtlasProjectBySecretsIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", tc.want, got)
			}
		})
	}
}

func TestFindProjectsForTeams(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{Name: "test-team", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team"}},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team", Namespace: "ns1"}},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns3"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team"}},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns1"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns2"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByTeamIndexer(zaptest.NewLogger(t))

			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[akov2.AtlasTeam](indexer.AtlasProjectByTeamIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", tc.want, got)
			}
		})
	}
}
