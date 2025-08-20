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

package atlasorgsettings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/atlasorgsettings"
)

var testAtlasSecret = corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test-atlas-secret",
		Namespace: "default",
	},
	Data: map[string][]byte{
		"orgId":         []byte("test-org-id"),
		"publicApiKey":  []byte("test-public-key"),
		"privateApiKey": []byte("test-private-key"),
	},
}

func TestNewAtlasOrgSettingsReconciler(t *testing.T) {
	fakeCluster := &fakeCluster{}
	atlasProvider := &atlasmock.TestProvider{}
	logger := zaptest.NewLogger(t)
	globalSecretRef := types.NamespacedName{Name: "global-secret", Namespace: "default"}

	rec := NewAtlasOrgSettingsReconciler(
		fakeCluster,
		atlasProvider,
		logger,
		client.ObjectKey{Name: globalSecretRef.Name, Namespace: globalSecretRef.Namespace},
		true,
		false,
	)

	assert.NotNil(t, rec)
}

func TestAtlasOrgSettingsHandlerFor(t *testing.T) {
	handler := &AtlasOrgSettingsHandler{}
	obj, preds := handler.For()
	assert.IsType(t, &akov2.AtlasOrgSettings{}, obj)
	assert.NotNil(t, preds)
}

func TestSetupWithManager(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	mgr, err := manager.New(&rest.Config{}, manager.Options{Scheme: scheme})
	require.NoError(t, err)
	fakeMgr := &fakeManager{
		Manager: mgr,
		client:  fake.NewClientBuilder().WithScheme(scheme).Build(),
	}

	handler := &AtlasOrgSettingsHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			AtlasProvider:   nil,
			Client:          fakeMgr.GetClient(),
			Log:             &zap.SugaredLogger{},
			GlobalSecretRef: client.ObjectKey{},
		},
		deletionProtection: false,
		serviceBuilder: func(clientSet *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
			return nil
		},
	}

	require.NoError(t, handler.SetupWithManager(fakeMgr, &fakeReconciler{}, controller.Options{}))
}

func TestFindSecretsForOrgSettings(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	logger := zaptest.NewLogger(t)
	ctx := context.Background()

	testCases := []struct {
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
						Name:      "test-secret",
						Namespace: "default",
					},
				},
			},
			want: []reconcile.Request{},
		},
		{
			name: "hit on linked credential secret",
			objects: []client.Object{
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-credential",
						Namespace: "default",
					},
				},
				&akov2.AtlasOrgSettings{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "org-settings",
						Namespace: "default",
					},
					Spec: akov2.AtlasOrgSettingsSpec{
						OrgID: "test-org-id",
						ConnectionSecretRef: &api.LocalObjectReference{
							Name: "test-credential",
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "org-settings", Namespace: "default"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			credIndex := indexer.NewAtlasOrgSettingsByConnectionSecretIndexer(logger)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tc.objects...).
				WithIndex(&akov2.AtlasOrgSettings{}, indexer.AtlasOrgSettingsBySecretsIndex, credIndex.Keys).
				Build()

			handler := &AtlasOrgSettingsHandler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    logger.Sugar(),
				},
			}
			mapper := handler.findSecretsForOrgSettings()

			got := mapper(ctx, tc.objects[0])

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAtlasOrgSettingsHandler_NewReconcileContext(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))

	orgSettings := &akov2.AtlasOrgSettings{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-org-settings",
			Namespace: "default",
		},
		Spec: akov2.AtlasOrgSettingsSpec{
			OrgID: "test-org-id",
			ConnectionSecretRef: &api.LocalObjectReference{
				Name: "test-atlas-secret",
			},
			ApiAccessListRequired:   pointer.MakePtr(true),
			GenAIFeaturesEnabled:    pointer.MakePtr(false),
			MultiFactorAuthRequired: pointer.MakePtr(true),
		},
	}

	k8sClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(&testAtlasSecret, orgSettings).Build()

	handler := &AtlasOrgSettingsHandler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: k8sClient,
			AtlasProvider: &atlasmock.TestProvider{
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return &atlas.ClientSet{}, nil
				},
			},
		},
		deletionProtection: false,
		serviceBuilder: func(clientSet *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
			return nil
		},
	}

	ctx := context.Background()
	req, err := handler.newReconcileRequest(ctx, orgSettings)

	require.NoError(t, err)
	assert.NotNil(t, req)
}

func TestAtlasOrgSettingsHandler_ServiceBuilder(t *testing.T) {
	handler := &AtlasOrgSettingsHandler{
		serviceBuilder: func(clientSet *atlas.ClientSet) atlasorgsettings.AtlasOrgSettingsService {
			return nil
		},
	}

	clientSet := &atlas.ClientSet{}
	service := handler.serviceBuilder(clientSet)

	assert.Nil(t, service)
}

func TestAtlasOrgSettingsHandler_DeletionProtection(t *testing.T) {
	testCases := []struct {
		name               string
		deletionProtection bool
	}{
		{
			name:               "deletion protection enabled",
			deletionProtection: true,
		},
		{
			name:               "deletion protection disabled",
			deletionProtection: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := &AtlasOrgSettingsHandler{
				deletionProtection: tc.deletionProtection,
			}

			assert.Equal(t, tc.deletionProtection, handler.deletionProtection)
		})
	}
}

type fakeManager struct {
	controllerruntime.Manager
	client client.Client
}

func (f *fakeManager) GetClient() client.Client {
	return f.client
}

type fakeCluster struct {
	cluster.Cluster
}

func (m *fakeCluster) GetClient() client.Client {
	return &fakeClient{}
}

type fakeClient struct {
	client.Client
}

type fakeReconciler struct {
	reconcile.Reconciler
}
