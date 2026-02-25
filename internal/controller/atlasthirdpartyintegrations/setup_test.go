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

package integrations

import (
	"context"
	"testing"

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
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/thirdpartyintegration"
)

var fakeAtlasSecret = corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name: "fake-atlas-secret",
	},
	Data: map[string][]byte{
		"orgId":         ([]byte)("fake-org"),
		"publicApiKey":  ([]byte)("pubkey"),
		"privateApiKey": ([]byte)("-"),
	},
}

var fakeProject = akov2.AtlasProject{
	ObjectMeta: metav1.ObjectMeta{Name: "fake-project"},
	Spec: v1.AtlasProjectSpec{
		Name: "fake-project",
	},
}

var referenceFakeProject = v1.ProjectDualReference{
	ProjectRef: &common.ResourceRefNamespaced{
		Name: "fake-project",
	},
	ConnectionSecret: &api.LocalObjectReference{
		Name: "fake-atlas-secret",
	},
}

func TestNewAtlasThirdPartyIntegrationsReconciler(t *testing.T) {
	fakeCluster := &fakeCluster{}
	atlasProvider := &atlasmock.TestProvider{}
	logger := zaptest.NewLogger(t)
	globalSecretRef := types.NamespacedName{Name: "global-secret", Namespace: "default"}

	rec := NewAtlasThirdPartyIntegrationsReconciler(
		fakeCluster, atlasProvider, true, logger, globalSecretRef, false,
	)
	assert.NotNil(t, rec)
}

func TestAtlasThirdPartyIntegrationHandlerFor(t *testing.T) {
	handler := &AtlasThirdPartyIntegrationHandler{}
	obj, preds := handler.For()
	assert.IsType(t, &akov2.AtlasThirdPartyIntegration{}, obj)
	assert.NotNil(t, preds)
}

func TestSetupWithManager(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	mgr, err := manager.New(&rest.Config{}, manager.Options{Scheme: scheme})
	require.NoError(t, err)
	fakeMgr := &fakeManager{
		Manager: mgr,
		client:  fake.NewClientBuilder().WithScheme(scheme).Build(),
	}

	handler := &AtlasThirdPartyIntegrationHandler{
		StateHandler: nil,
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   nil,
			Client:          fakeMgr.GetClient(),
			Log:             &zap.SugaredLogger{},
			GlobalSecretRef: client.ObjectKey{},
		},
		deletionProtection: false,
	}

	require.NoError(t, handler.SetupWithManager(fakeMgr, &fakeReconciler{}, controller.Options{}))
}

func TestNewReconcileRequest(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	integration := akov2.AtlasThirdPartyIntegration{
		Spec: akov2.AtlasThirdPartyIntegrationSpec{
			ProjectDualReference: referenceFakeProject,
			Type:                 "WEBHOOK",
			Webhook: &akov2.WebhookIntegration{
				URLSecretRef: api.LocalObjectReference{
					Name: "webhook-secret",
				},
			},
		},
	}
	k8sClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&fakeAtlasSecret, &fakeProject, &integration).Build()

	h := AtlasThirdPartyIntegrationHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: k8sClient,
			AtlasProvider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					projectAPI := mockFindFakeParentProject(t)
					integrationsAPI := mockadmin.NewThirdPartyIntegrationsApi(t)
					return &atlas.ClientSet{
						SdkClient20250312013: &admin.APIClient{
							ProjectsApi:               projectAPI,
							ThirdPartyIntegrationsApi: integrationsAPI,
						},
					}, nil
				},
			},
		},
		deletionProtection: false,
		serviceBuilder:     thirdpartyintegration.NewThirdPartyIntegrationServiceFromClientSet,
	}
	ctx := context.Background()
	req, err := h.newReconcileRequest(ctx, &integration)

	require.NoError(t, err)
	assert.NotEmpty(t, req)
}

func TestIntegrationForSecretMapFunc(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	logger := zaptest.NewLogger(t)
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		objects []client.Object
		want    []reconcile.Request
	}{
		{
			name:    "nil on non-secret object",
			objects: []client.Object{&corev1.ConfigMap{}},
			want:    nil,
		},
		{
			name: "empty on non-linked secret",
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-secret",
					},
				},
			},
			want: []reconcile.Request{},
		},
		{
			name: "hit on linked credential",
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-credential",
					},
				},
				&akov2.AtlasThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name: "integration",
					},
					Spec: akov2.AtlasThirdPartyIntegrationSpec{
						ProjectDualReference: v1.ProjectDualReference{
							ConnectionSecret: &api.LocalObjectReference{
								Name: "test-credential",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "integration"}},
			},
		},
		{
			name: "hit on linked credential",
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-prometheus-secret",
					},
				},
				&akov2.AtlasThirdPartyIntegration{
					ObjectMeta: metav1.ObjectMeta{
						Name: "integration-2",
					},
					Spec: akov2.AtlasThirdPartyIntegrationSpec{
						Type: "PROMETHEUS",
						Prometheus: &akov2.PrometheusIntegration{
							PrometheusCredentialsSecretRef: api.LocalObjectReference{
								Name: "test-prometheus-secret",
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "integration-2"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			secretIndex := indexer.NewAtlasThirdPartyIntegrationBySecretsIndexer(logger)
			credIndex := indexer.NewAtlasThirdPartyIntegrationByCredentialIndexer(logger)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).
				WithIndex(&akov2.AtlasThirdPartyIntegration{}, indexer.AtlasThirdPartyIntegrationBySecretsIndex, secretIndex.Keys).
				WithIndex(&akov2.AtlasThirdPartyIntegration{}, indexer.AtlasThirdPartyIntegrationCredentialsIndex, credIndex.Keys).
				Build()

			handler := &AtlasThirdPartyIntegrationHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    logger.Sugar(),
				},
			}
			mapper := handler.integrationForSecretMapFunc()

			got := mapper(ctx, tc.objects[0])

			assert.Equal(t, tc.want, got)
		})
	}
}

func mockFindFakeParentProject(t *testing.T) *mockadmin.ProjectsApi {
	projectAPI := mockadmin.NewProjectsApi(t)
	projectAPI.EXPECT().GetGroupByName(mock.Anything, "fake-project").
		Return(admin.GetGroupByNameApiRequest{ApiService: projectAPI})
	projectAPI.EXPECT().GetGroupByNameExecute(mock.Anything).
		Return(&admin.Group{Id: pointer.MakePtr("testProjectID")}, nil, nil)
	return projectAPI
}

type fakeManager struct {
	controllerruntime.Manager
	client client.Client
}

func (f *fakeManager) GetClient() client.Client { return f.client }

type fakeClient struct {
	client.Client
}

type fakeCluster struct {
	cluster.Cluster
}

func (m *fakeCluster) GetClient() client.Client {
	return fakeClient{}
}

type fakeReconciler struct {
	reconcile.Reconciler
}
