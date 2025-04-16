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

package atlassearchindexconfig

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
)

func TestAtlasSearchIndexConfigReconciler_Reconcile(t *testing.T) {
	t.Run("should skip silently if resource is not referenced", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "testSearchIndexConfig",
					Namespace: "mongodb-atlas-system",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("should skip reconciliation when annotation is set for AtlasSearchIndexConfig resource", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("should transition to invalid state when resource version is invalid", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
				Labels: map[string]string{
					customresource.ResourceVersion: "SomeRandomVersion",
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithStatusSubresource(searchIndexConfig).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		conditions := searchIndexConfig.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[1].Status)
		assert.Equal(t, string(workflow.AtlasResourceVersionIsInvalid), conditions[1].Reason)
		assert.Equal(t, "SomeRandomVersion is not a valid semver version for label mongodb.com/atlas-resource-version", conditions[1].Message)
	})

	t.Run("should transition to unsupported state when resource is not supported by platform", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithStatusSubresource(searchIndexConfig).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
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
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		conditions := searchIndexConfig.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Equal(t, string(workflow.AtlasGovUnsupported), conditions[0].Reason)
		assert.Equal(t, "the AtlasSearchIndexConfig is not supported by Atlas for government", conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to terminate state when failed to list instances", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		deploymentIndexer := indexer.NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithStatusSubresource(searchIndexConfig).
			WithIndex(
				deploymentIndexer.Object(),
				deploymentIndexer.Name(),
				deploymentIndexer.Keys,
			).
			WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list instances")
			}}).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
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
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		conditions := searchIndexConfig.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Equal(t, string(workflow.Internal), conditions[0].Reason)
		assert.Equal(t, "failed to list instances", conditions[0].Message)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to lock state when referred by an instance", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		atlasDeployment := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testAtlasDeployment",
				Namespace: "mongodb-atlas-system",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					SearchIndexes: []akov2.SearchIndex{
						{
							Name: "testSearchIndex",
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      searchIndexConfig.GetName(),
									Namespace: searchIndexConfig.GetNamespace(),
								},
							},
						},
					},
				},
			},
			Status: status.AtlasDeploymentStatus{},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		deploymentIndexer := indexer.NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig, atlasDeployment).
			WithStatusSubresource(searchIndexConfig).
			WithIndex(
				deploymentIndexer.Object(),
				deploymentIndexer.Name(),
				deploymentIndexer.Keys,
			).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
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
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.NotEmpty(t, searchIndexConfig.GetFinalizers())
		conditions := searchIndexConfig.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[0].Status)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to release state when not referred by any instance", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "testSearchIndexConfig",
				Namespace:  "mongodb-atlas-system",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		deploymentIndexer := indexer.NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithStatusSubresource(searchIndexConfig).
			WithIndex(
				deploymentIndexer.Object(),
				deploymentIndexer.Name(),
				deploymentIndexer.Keys,
			).
			Build()

		reconciler := &AtlasSearchIndexConfigReconciler{
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
					Name:      searchIndexConfig.GetName(),
					Namespace: searchIndexConfig.GetNamespace(),
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.Empty(t, searchIndexConfig.GetFinalizers())
		conditions := searchIndexConfig.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, api.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[0].Status)
		assert.Equal(t, api.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
	})

	t.Run("should transition to ready when finalizer is already set", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "testSearchIndexConfig",
				Namespace:  "mongodb-atlas-system",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		reconciler := &AtlasSearchIndexConfigReconciler{
			Log: zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.lock(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{}, result)
		assert.NotEmpty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})

	t.Run("should transition to terminate when failed to set finalizer", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithInterceptorFuncs(interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to set finalizer")
				},
			}).
			Build()
		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.lock(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.Empty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotSet), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to set finalizer", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to ready when setting finalizer", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			Build()
		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.lock(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{}, result)
		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.NotEmpty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})

	t.Run("should transition to ready when finalizer is already unset", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testSearchIndexConfig",
				Namespace: "mongodb-atlas-system",
			},
		}
		reconciler := &AtlasSearchIndexConfigReconciler{
			Log: zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.release(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{}, result)
		assert.Empty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})

	t.Run("should transition to terminate when failed to unset finalizer", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "testSearchIndexConfig",
				Namespace:  "mongodb-atlas-system",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			WithInterceptorFuncs(interceptor.Funcs{
				Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
					return errors.New("failed to set finalizer")
				},
			}).
			Build()
		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.release(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.NotEmpty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotRemoved), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to set finalizer", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to ready when unsetting finalizer", func(t *testing.T) {
		searchIndexConfig := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:       "testSearchIndexConfig",
				Namespace:  "mongodb-atlas-system",
				Finalizers: []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(searchIndexConfig).
			Build()
		reconciler := &AtlasSearchIndexConfigReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result := reconciler.release(ctx, searchIndexConfig)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(searchIndexConfig), searchIndexConfig))
		assert.Empty(t, searchIndexConfig.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
	})
}
