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

package atlasstream

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
)

func TestConnectionReconcile(t *testing.T) {
	t.Run("should terminate silently when resource is not found", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})
}

func TestEnsureAtlasStreamConnection(t *testing.T) {
	t.Run("should skip reconciliation when annotation is set", func(t *testing.T) {
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-connection",
				Namespace: "default",
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
				},
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("should unset finalizer and delete if reconciliation is skipped", func(t *testing.T) {
		now := metav1.Now()
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-stream-processing-connection",
				Namespace:         "default",
				DeletionTimestamp: &now,
				Finalizers:        []string{"mongodbatlas/finalizer"},
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
				},
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		err = k8sClient.Delete(context.Background(), streamConnection)
		assert.True(t, apierrors.IsNotFound(err))
	})

	t.Run("should terminate upon failed patching when unsetting finalizer", func(t *testing.T) {
		now := metav1.Now()
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-stream-processing-connection",
				Namespace:         "default",
				DeletionTimestamp: &now,
				Finalizers:        []string{"mongodbatlas/finalizer"},
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
				},
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			WithInterceptorFuncs(interceptor.Funcs{Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				return errors.New("failed to patch")
			}}).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)
	})

	t.Run("should transition to invalid state when resource version is invalid", func(t *testing.T) {
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-connection",
				Namespace: "default",
				Labels: map[string]string{
					customresource.ResourceVersion: "no-semver",
				},
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			WithStatusSubresource(streamConnection).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamConnection), streamConnection))
		conditions := streamConnection.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[1].Status)
		assert.Equal(t, string(workflow.AtlasResourceVersionIsInvalid), conditions[1].Reason)
		assert.Equal(t, "no-semver is not a valid semver version for label mongodb.com/atlas-resource-version", conditions[1].Message)
	})

	t.Run("should transition to unsupported state when resource is not supported by platform", func(t *testing.T) {
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-connection",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			WithStatusSubresource(streamConnection).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamConnection), streamConnection))
		conditions := streamConnection.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Equal(t, string(workflow.AtlasGovUnsupported), conditions[0].Reason)
		assert.Equal(t, "the AtlasStreamConnection is not supported by Atlas for government", conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to terminate state when failed to list instances", func(t *testing.T) {
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-connection",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		streamInstanceIndexer := indexer.NewAtlasStreamInstanceByConnectionIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			WithStatusSubresource(streamConnection).
			WithIndex(
				streamInstanceIndexer.Object(),
				streamInstanceIndexer.Name(),
				streamInstanceIndexer.Keys,
			).
			WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list instances")
			}}).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamConnection), streamConnection))
		conditions := streamConnection.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Equal(t, string(workflow.Internal), conditions[0].Reason)
		assert.Equal(t, "failed to list instances", conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to lock state when referred by an instance", func(t *testing.T) {
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-stream-processing-connection",
						Namespace: "default",
					},
				},
			},
		}
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-connection",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		streamInstanceIndexer := indexer.NewAtlasStreamInstanceByConnectionIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance, streamConnection).
			WithStatusSubresource(streamConnection).
			WithIndex(
				streamInstanceIndexer.Object(),
				streamInstanceIndexer.Name(),
				streamInstanceIndexer.Keys,
			).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamConnection), streamConnection))
		assert.NotEmpty(t, streamConnection.GetFinalizers())
		conditions := streamConnection.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[0].Status)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to release state when not referred by any instance", func(t *testing.T) {
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "my-stream-processing-connection",
				Namespace:  "default",
				Finalizers: []string{customresource.FinalizerLabel},
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-conn",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		streamInstanceIndexer := indexer.NewAtlasStreamInstanceByConnectionIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamConnection).
			WithStatusSubresource(streamConnection).
			WithIndex(
				streamInstanceIndexer.Object(),
				streamInstanceIndexer.Name(),
				streamInstanceIndexer.Keys,
			).
			Build()

		reconciler := &AtlasStreamsConnectionReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-connection",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamConnection), streamConnection))
		assert.Empty(t, streamConnection.GetFinalizers())
		conditions := streamConnection.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[0].Status)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})
}

func TestFindStreamConnectionsForStreamInstances(t *testing.T) {
	t.Run("should fail when watching wrong object", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasStreamsConnectionReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Nil(t, reconciler.findStreamConnectionsForStreamInstances(context.Background(), &akov2.AtlasProject{}))
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zap.WarnLevel, logs.All()[0].Level)
		assert.Equal(t, "watching AtlasStreamInstance but got *v1.AtlasProject", logs.All()[0].Message)
	})

	t.Run("should return slice of requests for instances", func(t *testing.T) {
		instance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance1",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance1",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "connection1",
						Namespace: "default",
					},
					{
						Name:      "connection2",
						Namespace: "other-ns",
					},
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(instance).
			Build()
		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		requests := reconciler.findStreamConnectionsForStreamInstances(context.Background(), instance)
		assert.Equal(
			t,
			[]ctrl.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      "connection1",
						Namespace: "default",
					},
				},
				{
					NamespacedName: types.NamespacedName{
						Name:      "connection2",
						Namespace: "other-ns",
					},
				},
			},
			requests,
		)
	})
}

func TestLock(t *testing.T) {
	t.Run("should transition to ready when finalizer is already set", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		reconciler := &AtlasStreamsConnectionReconciler{
			Log: zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.lock(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.NotEmpty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})

	t.Run("should transition to terminate when failed to set finalizer", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "conn1",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection).
			WithInterceptorFuncs(interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to set finalizer")
				},
			}).
			Build()
		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.lock(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(connection), connection))
		assert.Empty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotSet), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to set finalizer", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to ready when setting finalizer", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "conn1",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection).
			Build()
		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.lock(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(connection), connection))
		assert.NotEmpty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})
}

func TestRelease(t *testing.T) {
	t.Run("should transition to ready when finalizer is already unset", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{},
		}
		reconciler := &AtlasStreamsConnectionReconciler{
			Log: zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.release(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.Empty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})

	t.Run("should transition to terminate when failed to unset finalizer", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "conn1",
				Namespace:  "default",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection).
			WithInterceptorFuncs(interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to set finalizer")
				},
			}).
			Build()
		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.release(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(connection), connection))
		assert.NotEmpty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotRemoved), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to set finalizer", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to ready when unsetting finalizer", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "conn1",
				Namespace:  "default",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection).
			Build()
		reconciler := &AtlasStreamsConnectionReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.release(ctx, connection)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(connection), connection))
		assert.Empty(t, connection.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})
}
