package atlasdatabaseuser

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
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
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
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
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
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
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    zaptest.NewLogger(t).Sugar(),
				},
				AtlasProvider: DefaultTestProvider(t),
				EventRecorder: record.NewFakeRecorder(10),
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
			AtlasReconciler: reconciler.AtlasReconciler{
				Log: zap.New(core).Sugar(),
			},
		}

		assert.Equal(t, ctrl.Result{}, c.notFound(ctrl.Request{NamespacedName: types.NamespacedName{Name: "object", Namespace: "test"}}))
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zapcore.Level(0), logs.All()[0].Level)
		assert.Equal(t, "object test/object doesn't exist, was it deleted after reconcile request?", logs.All()[0].Message)
	})
}

func TestFail(t *testing.T) {
	t.Run("failed to retrieve custom resource", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		c := &AtlasDatabaseUserReconciler{
			AtlasReconciler: reconciler.AtlasReconciler{
				Log: zap.New(core).Sugar(),
			},
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
			AtlasReconciler: reconciler.AtlasReconciler{
				Log: zap.New(core).Sugar(),
			},
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
				AtlasReconciler: reconciler.AtlasReconciler{
					Log: zap.New(core).Sugar(),
				},
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

func TestReady(t *testing.T) {
	tests := map[string]struct {
		dbUser          *akov2.AtlasDatabaseUser
		passwordVersion string
		interceptors    interceptor.Funcs

		expectedResult     ctrl.Result
		expectedConditions []api.Condition
	}{
		"fail to set finalizer": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			passwordVersion: "1",
			interceptors: interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to set finalizer")
				},
			},
			expectedResult: workflow.Terminate(workflow.AtlasFinalizerNotSet, errors.New("")).ReconcileResult(),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.AtlasFinalizerNotSet)).
					WithMessageRegexp("failed to set finalizer"),
			},
		},
		"fail to set last applied config": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			passwordVersion: "1",
			interceptors: interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					if patch.Type() == types.JSONPatchType {
						return nil
					}

					return errors.New("failed to set last applied config")
				},
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("")).ReconcileResult(),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DatabaseUserReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to set last applied config"),
			},
		},
		"don't requeue when it's a linked resource": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      "my-project",
							Namespace: "default",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			passwordVersion: "1",
			expectedResult:  workflow.OK().ReconcileResult(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.DatabaseUserReadyType),
			},
		},
		"requeue when it's a standalone resource": {
			dbUser: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: "project-id",
						},
						ConnectionSecret: &api.LocalObjectReference{
							Name: "user-creds",
						},
					},
					Username: "user1",
					PasswordSecret: &common.ResourceRef{
						Name: "user-pass",
					},
					DatabaseName: "admin",
				},
			},
			passwordVersion: "1",
			expectedResult:  workflow.Requeue(10 * time.Minute).ReconcileResult(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.DatabaseUserReadyType),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.dbUser).
				WithStatusSubresource(tt.dbUser).
				WithInterceptorFuncs(tt.interceptors).
				Build()

			logger := zaptest.NewLogger(t).Sugar()
			c := &AtlasDatabaseUserReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
				independentSyncPeriod: 10 * time.Minute,
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			assert.Equal(t, tt.expectedResult, c.ready(ctx, tt.dbUser, tt.passwordVersion))
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
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
				AtlasReconciler: reconciler.AtlasReconciler{
					Log:    zaptest.NewLogger(t).Sugar(),
					Client: k8sClient,
				},
			}
			got := reconciler.findAtlasDatabaseUserForSecret(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}
