package atlasfederatedauth

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	mockadminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/mockadmin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestReconcile(t *testing.T) {
	t.Run("should reconcile successfully with existing configuration", func(t *testing.T) {
		orgID := "616ec36209c07e743422b7cc" //nolint:gosec
		projectID := "abc123"
		fedSettingsID := "651438d6cda56304464dd128" //nolint:gosec
		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "api-secret",
				Namespace: "default",
				Labels: map[string]string{
					"atlas.mongodb.com/type": "credentials",
				},
			},
			Data: map[string][]byte{
				"orgId":         []byte(orgID),
				"publicApiKey":  []byte("a1b2c3"),
				"privateApiKey": []byte("abcdef123456"),
			},
			Type: "Opaque",
		}
		project := akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "MyProject",
			},
			Status: status.AtlasProjectStatus{ID: projectID},
		}
		fedAuth := akov2.AtlasFederatedAuth{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-fed-auth",
				Namespace: "default",
			},
			Spec: akov2.AtlasFederatedAuthSpec{
				Enabled: true,
				ConnectionSecretRef: common.ResourceRefNamespaced{
					Name:      secret.Name,
					Namespace: secret.Namespace,
				},
				DomainAllowList:          []string{"qa-27092023.com", "cloud-qa.mongodb.com"},
				DomainRestrictionEnabled: pointer.MakePtr(true),
				SSODebugEnabled:          pointer.MakePtr(false),
				PostAuthRoleGrants:       []string{"ORG_OWNER"},
			},
		}
		sch := runtime.NewScheme()
		sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{})
		sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
		sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasProject{})
		sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasFederatedAuth{})
		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(&secret, &project, &fedAuth).
			WithStatusSubresource(&fedAuth).
			Build()

		logger := zaptest.NewLogger(t).Sugar()
		fedAuthAPI := mockadminv20241113001.NewFederatedAuthenticationApi(t)
		fedAuthAPI.EXPECT().GetFederationSettings(context.Background(), orgID).
			Return(adminv20241113001.GetFederationSettingsApiRequest{ApiService: fedAuthAPI})
		fedAuthAPI.EXPECT().GetFederationSettingsExecute(mock.Anything).
			Return(
				&adminv20241113001.OrgFederationSettings{
					Id:                     pointer.MakePtr(fedSettingsID),
					IdentityProviderId:     pointer.MakePtr("0oawce8e76SR9K7Tq357"),
					FederatedDomains:       &[]string{"qa-27092023.com", "cloud-qa.mongodb.com"},
					HasRoleMappings:        pointer.MakePtr(false),
					IdentityProviderStatus: pointer.MakePtr("ACTIVE"),
				},
				&http.Response{},
				nil,
			)
		fedAuthAPI.EXPECT().GetIdentityProvider(context.Background(), fedSettingsID, "0oawce8e76SR9K7Tq357").
			Return(adminv20241113001.GetIdentityProviderApiRequest{ApiService: fedAuthAPI})
		fedAuthAPI.EXPECT().GetIdentityProviderExecute(mock.Anything).
			Return(
				&adminv20241113001.FederationIdentityProvider{
					Id:        "65143bd1612f01218e885cf2",
					OktaIdpId: "0oawce8e76SR9K7Tq357",
				},
				&http.Response{},
				nil,
			)
		fedAuthAPI.EXPECT().GetConnectedOrgConfig(context.Background(), fedSettingsID, orgID).
			Return(adminv20241113001.GetConnectedOrgConfigApiRequest{ApiService: fedAuthAPI})
		fedAuthAPI.EXPECT().GetConnectedOrgConfigExecute(mock.Anything).
			Return(
				&adminv20241113001.ConnectedOrgConfig{
					OrgId:                    "616ec36209c07e743422b7cc",
					DomainAllowList:          &[]string{"qa-27092023.com", "cloud-qa.mongodb.com"},
					DomainRestrictionEnabled: true,
					IdentityProviderId:       pointer.MakePtr("0oawce8e76SR9K7Tq357"),
					PostAuthRoleGrants:       &[]string{"ORG_OWNER"},
				},
				&http.Response{},
				nil,
			)
		groupAPI := mockadmin.NewProjectsApi(t)
		groupAPI.EXPECT().ListProjects(context.Background()).
			Return(admin.ListProjectsApiRequest{ApiService: groupAPI})
		groupAPI.EXPECT().ListProjectsExecute(mock.Anything).
			Return(
				&admin.PaginatedAtlasGroup{
					Results: &[]admin.Group{
						{
							Id:   pointer.MakePtr(projectID),
							Name: "MyProject",
						},
					},
					TotalCount: pointer.MakePtr(1),
				},
				&http.Response{},
				nil,
			)
		atlasProvider := atlasmock.TestProvider{
			SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
				return &atlas.ClientSet{SdkClient20241113001: &adminv20241113001.APIClient{
					FederatedAuthenticationApi: fedAuthAPI,
				}}, orgID, nil
			},
			SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
				return &admin.APIClient{
					ProjectsApi: groupAPI,
				}, orgID, nil
			},
			IsCloudGovFunc: func() bool {
				return false
			},
			IsSupportedFunc: func() bool {
				return true
			},
		}

		reconciler := &AtlasFederatedAuthReconciler{
			Client:                      k8sClient,
			Log:                         logger,
			AtlasProvider:               &atlasProvider,
			EventRecorder:               record.NewFakeRecorder(10),
			ObjectDeletionProtection:    true,
			SubObjectDeletionProtection: true,
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: fedAuth.Namespace,
					Name:      fedAuth.Name,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: 0}, result)

		fedAuthResult := akov2.AtlasFederatedAuth{}
		err = k8sClient.Get(context.Background(), client.ObjectKeyFromObject(&fedAuth), &fedAuthResult)
		assert.NoError(t, err)
		assert.Condition(t, func() bool {
			expected := map[api.ConditionType]struct{}{
				api.ResourceVersionStatus:  {},
				api.FederatedAuthReadyType: {},
				api.ReadyType:              {},
			}

			for _, condition := range fedAuthResult.Status.Conditions {
				if _, ok := expected[condition.Type]; !ok {
					return false
				}

				if condition.Status != corev1.ConditionTrue {
					return false
				}

				delete(expected, condition.Type)
			}

			return len(expected) == 0
		})
	})
}

func TestFindAtlasFederatedAuthForSecret(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasProject{},
			want: nil,
		},
		{
			name: "same namespace",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns"},
			},
			initObjs: []client.Object{
				&akov2.AtlasFederatedAuth{
					ObjectMeta: metav1.ObjectMeta{Name: "auth1", Namespace: "ns"},
					Spec: akov2.AtlasFederatedAuthSpec{
						ConnectionSecretRef: common.ResourceRefNamespaced{Name: "name"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "auth1", Namespace: "ns"}},
			},
		},
		{
			name: "different namespace",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasFederatedAuth{
					ObjectMeta: metav1.ObjectMeta{Name: "auth1", Namespace: "ns"},
					Spec: akov2.AtlasFederatedAuthSpec{
						ConnectionSecretRef: common.ResourceRefNamespaced{Name: "name", Namespace: "ns2"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "auth1", Namespace: "ns"}},
			},
		},
		{
			name: "different reference",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasFederatedAuth{
					ObjectMeta: metav1.ObjectMeta{Name: "auth1", Namespace: "ns"},
					Spec: akov2.AtlasFederatedAuthSpec{
						ConnectionSecretRef: common.ResourceRefNamespaced{Name: "name"},
					},
				},
			},
			want: []reconcile.Request{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			indexer := indexer.NewAtlasFederatedAuthBySecretsIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(indexer.Object(), indexer.Name(), indexer.Keys).
				Build()
			reconciler := &AtlasFederatedAuthReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findAtlasFederatedAuthForSecret(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}
