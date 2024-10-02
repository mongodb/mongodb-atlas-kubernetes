package atlasdatabaseuser

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
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

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestReconcile(t *testing.T) {
	tests := map[string]struct {
		dbUser         *akov2.AtlasDatabaseUser
		interceptors   interceptor.Funcs
		expectedResult ctrl.Result
	}{
		"failed to retrieve user": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			interceptors: interceptor.Funcs{
				Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
					return errors.New("failed to get user")
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		},
		"user was not found": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user2",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			expectedResult: ctrl.Result{},
		},
		"skip reconciliation": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Annotations: map[string]string{
						"mongodb.com/atlas-reconciliation-policy": "skip",
					},
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			expectedResult: ctrl.Result{},
		},
		"handle user": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUser).
				WithInterceptorFuncs(tt.interceptors).
				Build()
			r := &AtlasDatabaseUserReconciler{
				Client: k8sClient,
				AtlasProvider: &atlasmock.TestProvider{
					IsSupportedFunc: func() bool {
						return true
					},
					IsCloudGovFunc: func() bool {
						return false
					},
					SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
						userAPI := mockadmin.NewDatabaseUsersApi(t)
						userAPI.EXPECT().GetDatabaseUser(context.Background(), "", "admin", "user1").
							Return(admin.GetDatabaseUserApiRequest{ApiService: userAPI})
						userAPI.EXPECT().GetDatabaseUserExecute(mock.AnythingOfType("admin.GetDatabaseUserApiRequest")).
							Return(nil, nil, nil)
						userAPI.EXPECT().CreateDatabaseUser(context.Background(), "", mock.AnythingOfType("*admin.CloudDatabaseUser")).
							Return(admin.CreateDatabaseUserApiRequest{ApiService: userAPI})
						userAPI.EXPECT().CreateDatabaseUserExecute(mock.AnythingOfType("admin.CreateDatabaseUserApiRequest")).
							Return(&admin.CloudDatabaseUser{}, nil, nil)

						clusterAPI := mockadmin.NewClustersApi(t)

						return &admin.APIClient{
							ClustersApi:      clusterAPI,
							DatabaseUsersApi: userAPI,
						}, "", nil
					},
				},
				EventRecorder: record.NewFakeRecorder(10),
				Log:           zaptest.NewLogger(t).Sugar(),
			}

			result, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "user1", Namespace: "default"}})
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestNotFound(t *testing.T) {
	t.Run("custom resource was not found", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		c := &AtlasDatabaseUserReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Equal(t, ctrl.Result{}, c.notFound(ctrl.Request{NamespacedName: types.NamespacedName{Name: "object", Namespace: "test"}}))
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.Level(0), logs.All()[0].Level)
		assert.Equal(t, "Object test/object doesn't exist, was it deleted after reconcile request?", logs.All()[0].Message)
	})
}

func TestFail(t *testing.T) {
	t.Run("failed to retrieve custom resource", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		c := &AtlasDatabaseUserReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Equal(
			t,
			ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			c.fail(ctrl.Request{NamespacedName: types.NamespacedName{Name: "object", Namespace: "test"}}, errors.New("failed to retrieve custom resource")),
		)
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.Level(2), logs.All()[0].Level)
		assert.Equal(t, "Failed to query object test/object: failed to retrieve custom resource", logs.All()[0].Message)
	})
}

func TestSkip(t *testing.T) {
	t.Run("skip reconciliation of custom resource", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		c := &AtlasDatabaseUserReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Equal(t, ctrl.Result{}, c.skip())
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.Level(0), logs.All()[0].Level)
		assert.Equal(t, "-> Skipping AtlasDatabaseUser reconciliation as annotation mongodb.com/atlas-reconciliation-policy=skip", logs.All()[0].Message)
	})
}

func TestTerminate(t *testing.T) {
	tests := map[string]struct {
		object         akov2.AtlasCustomResource
		condition      api.ConditionType
		reason         workflow.ConditionReason
		retry          bool
		err            error
		expectedResult ctrl.Result
		expectedLogs   []string
	}{
		"terminates reconciliation with retry": {
			object: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "ns-test",
				},
			},
			condition:      api.ProjectReadyType,
			reason:         workflow.Internal,
			retry:          true,
			err:            errors.New("failed to reconcile project"),
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedLogs: []string{
				"resource *v1.AtlasProject(ns-test/my-project) failed on condition ProjectReady: failed to reconcile project",
			},
		},
		"terminates reconciliation without retry": {
			object: &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "ns-test",
				},
			},
			condition:      api.StreamInstanceReadyType,
			reason:         workflow.StreamConnectionNotCreated,
			retry:          false,
			err:            errors.New("failed to reconcile stream instance"),
			expectedResult: ctrl.Result{},
			expectedLogs: []string{
				"resource *v1.AtlasStreamInstance(ns-test/my-project) failed on condition StreamInstanceReady: failed to reconcile stream instance",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			c := &AtlasDatabaseUserReconciler{
				Log: zap.New(core).Sugar(),
			}

			assert.Equal(t, tt.expectedResult, c.terminate(&workflow.Context{}, tt.object, tt.condition, tt.reason, tt.retry, tt.err))
			assert.Equal(t, len(tt.expectedLogs), logs.Len())
			for ix, log := range logs.All() {
				assert.Equal(t, zapcore.Level(2), log.Level)
				assert.Equal(t, tt.expectedLogs[ix], log.Message)
			}
		})
	}
}

func TestFindAtlasDatabaseUserForSecret(t *testing.T) {
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
				ObjectMeta: metav1.ObjectMeta{Name: "secret", Namespace: "ns"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: "ns"},
					Spec: akov2.AtlasDatabaseUserSpec{
						PasswordSecret: &common.ResourceRef{Name: "secret"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "user1", Namespace: "ns"}},
			},
		},
		{
			name: "different namespace",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "name", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDatabaseUser{
					ObjectMeta: metav1.ObjectMeta{Name: "user1", Namespace: "ns"},
					Spec: akov2.AtlasDatabaseUserSpec{
						PasswordSecret: &common.ResourceRef{Name: "secret"},
					},
				},
			},
			want: []reconcile.Request{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			secretsIndexer := indexer.NewAtlasDatabaseUserBySecretsIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(secretsIndexer.Object(), secretsIndexer.Name(), secretsIndexer.Keys).
				Build()
			reconciler := &AtlasDatabaseUserReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findAtlasDatabaseUserForSecret(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}
