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
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	atlas_controllers "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestRenconcile(t *testing.T) {
	tests := map[string]struct {
		atlasSDKMocker func() *admin.APIClient
		interceptors   interceptor.Funcs
		project        *akov2.AtlasProject
		result         reconcile.Result
		conditions     []api.Condition
		finalizers     []string
	}{
		"should handle project": {
			atlasSDKMocker: func() *admin.APIClient {
				notFoundErr := &admin.GenericOpenAPIError{}
				notFoundErr.SetModel(admin.ApiError{ErrorCode: "NOT_IN_GROUP"})
				projectsAPI := mockadmin.NewProjectsApi(t)
				projectsAPI.EXPECT().GetGroupByName(mock.Anything, "my-project").
					Return(admin.GetGroupByNameApiRequest{ApiService: projectsAPI})
				projectsAPI.EXPECT().GetGroupByNameExecute(mock.AnythingOfType("admin.GetGroupByNameApiRequest")).
					Return(
						nil,
						&http.Response{},
						notFoundErr,
					)
				projectsAPI.EXPECT().CreateGroup(mock.Anything, mock.AnythingOfType("*admin.Group")).
					Return(admin.CreateGroupApiRequest{ApiService: projectsAPI})
				projectsAPI.EXPECT().CreateGroupExecute(mock.AnythingOfType("admin.CreateGroupApiRequest")).
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
				WithObjects(tt.project, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("orgId"),
						"publicApiKey":  []byte("publicApiKey"),
						"privateApiKey": []byte("privateApiKey"),
					},
				}).
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
					SdkClientSetFunc: func(ctx context.Context, creds *atlas_controllers.Credentials, log *zap.SugaredLogger) (*atlas_controllers.ClientSet, error) {
						return &atlas_controllers.ClientSet{
							SdkClient20250312013: tt.atlasSDKMocker(),
						}, nil
					},
				},
				GlobalSecretRef: client.ObjectKey{
					Namespace: "default",
					Name:      "secret",
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

func TestLastSpecFrom(t *testing.T) {
	tests := map[string]struct {
		annotations      map[string]string
		expectedLastSpec *akov2.AtlasProjectSpec
		expectedError    string
	}{

		"should return nil when there is no last spec": {},
		"should return error when last spec annotation is wrong": {
			annotations:   map[string]string{"mongodb.com/last-applied-configuration": "{wrong}"},
			expectedError: "invalid character 'w' looking for beginning of object key string",
		},
		"should return last spec": {
			annotations: map[string]string{"mongodb.com/last-applied-configuration": "{\"name\": \"my-project\"}"},
			expectedLastSpec: &akov2.AtlasProjectSpec{
				Name: "my-project",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &akov2.AtlasProject{}
			p.WithAnnotations(tt.annotations)
			lastSpec, err := lastAppliedSpecFrom(p)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedLastSpec, lastSpec)
		})
	}
}

func TestSkipClearsMigratedResourcesLastConfig(t *testing.T) {
	ctx := context.Background()
	prj := akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-project",
			Namespace:   "test-ns",
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name:                      "test-project",
			PrivateEndpoints:          []akov2.PrivateEndpoint{{}},
			CloudProviderAccessRoles:  []akov2.CloudProviderAccessRole{{}},
			CloudProviderIntegrations: []akov2.CloudProviderIntegration{{}},
			AlertConfigurations:       []akov2.AlertConfiguration{{}},
			NetworkPeers:              []akov2.NetworkPeer{{}},
			X509CertRef:               &common.ResourceRefNamespaced{},
			Integrations:              []project.Integration{{}},
			EncryptionAtRest:          &akov2.EncryptionAtRest{},
			Auditing:                  &akov2.Auditing{},
			Settings:                  &akov2.ProjectSettings{},
			CustomRoles:               []akov2.CustomRole{{}},
			Teams:                     []akov2.Team{{}},
			BackupCompliancePolicyRef: &common.ResourceRefNamespaced{},
			ConnectionSecret:          &common.ResourceRefNamespaced{},
			ProjectIPAccessList:       []project.IPAccessList{{}},
			MaintenanceWindow:         project.MaintenanceWindow{},
		},
	}
	prj.Annotations[customresource.AnnotationLastAppliedConfiguration] = jsonize(t, prj.Spec)
	prj.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(&prj).
		WithStatusSubresource(&prj).
		Build()

	req := ctrl.Request{NamespacedName: types.NamespacedName{
		Name:      prj.Name,
		Namespace: prj.Namespace,
	}}
	r := AtlasProjectReconciler{
		Client:        k8sClient,
		Log:           zaptest.NewLogger(t).Sugar(),
		EventRecorder: record.NewFakeRecorder(30),
		AtlasProvider: &atlas.TestProvider{
			IsCloudGovFunc: func() bool {
				return false
			},
			IsSupportedFunc: func() bool {
				return true
			},
		},
	}

	result, err := r.Reconcile(ctx, req)

	require.Equal(t, reconcile.Result{}, result)
	require.NoError(t, err)
	require.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(&prj), &prj))
	lastApplied, err := customresource.ParseLastConfigApplied[akov2.AtlasProjectSpec](&prj)
	require.NoError(t, err)
	wantLastApplied := &akov2.AtlasProjectSpec{
		Name:                      "test-project",
		PrivateEndpoints:          nil,
		CustomRoles:               nil,
		NetworkPeers:              nil,
		ProjectIPAccessList:       nil,
		CloudProviderAccessRoles:  []akov2.CloudProviderAccessRole{{}},
		CloudProviderIntegrations: []akov2.CloudProviderIntegration{{}},
		AlertConfigurations:       []akov2.AlertConfiguration{{}},
		X509CertRef:               &common.ResourceRefNamespaced{},
		EncryptionAtRest:          &akov2.EncryptionAtRest{},
		Auditing:                  &akov2.Auditing{},
		Settings:                  &akov2.ProjectSettings{},
		Teams:                     []akov2.Team{{}},
		BackupCompliancePolicyRef: &common.ResourceRefNamespaced{},
		ConnectionSecret:          &common.ResourceRefNamespaced{},
	}
	assert.Equal(t, wantLastApplied, lastApplied)
}

func TestSkipClearsMigratedResourcesLastConfigDoesNotPanic(t *testing.T) {
	ctx := context.Background()
	prj := akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-project",
			Namespace:   "test-ns",
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "test-project",
		},
	}
	prj.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))
	require.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(&prj).
		WithStatusSubresource(&prj).
		Build()

	req := ctrl.Request{NamespacedName: types.NamespacedName{
		Name:      prj.Name,
		Namespace: prj.Namespace,
	}}
	r := AtlasProjectReconciler{
		Client:        k8sClient,
		Log:           zaptest.NewLogger(t).Sugar(),
		EventRecorder: record.NewFakeRecorder(30),
		AtlasProvider: &atlas.TestProvider{
			IsCloudGovFunc: func() bool {
				return false
			},
			IsSupportedFunc: func() bool {
				return true
			},
		},
	}

	result, err := r.Reconcile(ctx, req)

	require.Equal(t, reconcile.Result{}, result)
	require.NoError(t, err)
	require.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(&prj), &prj))
	lastApplied, err := customresource.ParseLastConfigApplied[akov2.AtlasProjectSpec](&prj)
	require.NoError(t, err)
	wantLastApplied := (*akov2.AtlasProjectSpec)(nil)
	assert.Equal(t, wantLastApplied, lastApplied)
}

func jsonize(t *testing.T, obj any) string {
	t.Helper()

	js, err := json.Marshal(obj)
	require.NoError(t, err)
	return string(js)
}
