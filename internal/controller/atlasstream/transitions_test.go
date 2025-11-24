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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestCreate(t *testing.T) {
	t.Run("should transition to in-progress state when creates a stream instance in Atlas", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "FRANKFURT_DEU",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
		}
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance, streamConnection).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamWorkspace(context.Background(), "my-project-id", mock.AnythingOfType("*admin.StreamsTenant")).
			Return(admin.CreateStreamWorkspaceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamWorkspaceExecute(mock.AnythingOfType("admin.CreateStreamWorkspaceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "FRANKFURT_DEU",
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					Connections: &[]admin.StreamsConnection{
						{
							Name: pointer.MakePtr("sample-connection"),
							Type: pointer.MakePtr("Sample"),
						},
					},
					GroupId: pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}

		result, err := reconciler.create(ctx, project, streamInstance)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceSetupInProgress), ctx.Conditions()[0].Reason)
		assert.Equal(t, "configuring stream instance in Atlas", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to a terminate state when fail to create a stream instance in Atlas", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "FRANKFURT_DEU",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
		}
		streamConnection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance, streamConnection).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamWorkspace(context.Background(), "my-project-id", mock.AnythingOfType("*admin.StreamsTenant")).
			Return(admin.CreateStreamWorkspaceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamWorkspaceExecute(mock.AnythingOfType("admin.CreateStreamWorkspaceApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to create instance"),
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}

		_, err := reconciler.create(ctx, project, streamInstance)
		assert.Error(t, err)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceNotCreated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to create instance", ctx.Conditions()[0].Message)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should handle deletion when deletion protection is disabled", func(t *testing.T) {
		t.Run("should successfully transition when deleting a stream instance in Atlas", func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "my-project-id",
				},
			}
			streamInstance := &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-stream-processing-instance",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: pointer.MakePtr(metav1.Now()),
				},
				Spec: akov2.AtlasStreamInstanceSpec{
					Name: "instance-0",
					Config: akov2.Config{
						Provider: "AWS",
						Region:   "FRANKFURT_DEU",
						Tier:     "SP30",
					},
					Project: common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ConnectionRegistry: []common.ResourceRefNamespaced{
						{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
				Status: status.AtlasStreamInstanceStatus{
					ID:        "instance-0-id",
					Hostnames: []string{"mdb://host1", "mdb://host2"},
					Connections: []status.StreamConnection{
						{
							Name: "sample-connection",
							ResourceRef: common.ResourceRefNamespaced{
								Name:      "my-sample-connection",
								Namespace: "default",
							},
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, streamInstance).
				Build()

			reconciler := &AtlasStreamsInstanceReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamWorkspace(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamWorkspaceExecute(mock.AnythingOfType("admin.DeleteStreamWorkspaceApiRequest")).
				Return(
					&http.Response{},
					nil,
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312009: &admin.APIClient{
						StreamsApi: streamsAPI,
					},
				},
			}

			result, err := reconciler.delete(ctx, project, streamInstance)
			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{}, result)
			assert.Empty(t, streamInstance.Finalizers)
		})

		t.Run("should successfully transition when keep resource annotation is set", func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "my-project-id",
				},
			}
			streamInstance := &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-stream-processing-instance",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: pointer.MakePtr(metav1.Now()),
					Annotations: map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
					},
				},
				Spec: akov2.AtlasStreamInstanceSpec{
					Name: "instance-0",
					Config: akov2.Config{
						Provider: "AWS",
						Region:   "FRANKFURT_DEU",
						Tier:     "SP30",
					},
					Project: common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ConnectionRegistry: []common.ResourceRefNamespaced{
						{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
				Status: status.AtlasStreamInstanceStatus{
					ID:        "instance-0-id",
					Hostnames: []string{"mdb://host1", "mdb://host2"},
					Connections: []status.StreamConnection{
						{
							Name: "sample-connection",
							ResourceRef: common.ResourceRefNamespaced{
								Name:      "my-sample-connection",
								Namespace: "default",
							},
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, streamInstance).
				Build()

			reconciler := &AtlasStreamsInstanceReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}
			ctx := &workflow.Context{
				Context: context.Background(),
			}

			result, err := reconciler.delete(ctx, project, streamInstance)
			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{}, result)
			assert.Empty(t, streamInstance.Finalizers)
		})

		t.Run("should fail transition when unable to delete a stream instance in Atlas", func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "my-project-id",
				},
			}
			streamInstance := &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-stream-processing-instance",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: pointer.MakePtr(metav1.Now()),
				},
				Spec: akov2.AtlasStreamInstanceSpec{
					Name: "instance-0",
					Config: akov2.Config{
						Provider: "AWS",
						Region:   "FRANKFURT_DEU",
						Tier:     "SP30",
					},
					Project: common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ConnectionRegistry: []common.ResourceRefNamespaced{
						{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
				Status: status.AtlasStreamInstanceStatus{
					ID:        "instance-0-id",
					Hostnames: []string{"mdb://host1", "mdb://host2"},
					Connections: []status.StreamConnection{
						{
							Name: "sample-connection",
							ResourceRef: common.ResourceRefNamespaced{
								Name:      "my-sample-connection",
								Namespace: "default",
							},
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, streamInstance).
				Build()

			reconciler := &AtlasStreamsInstanceReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamWorkspace(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamWorkspaceExecute(mock.AnythingOfType("admin.DeleteStreamWorkspaceApiRequest")).
				Return(
					&http.Response{},
					errors.New("failed to delete instance"),
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312009: &admin.APIClient{
						StreamsApi: streamsAPI,
					},
				},
			}

			_, err := reconciler.delete(ctx, project, streamInstance)
			assert.Error(t, err)
			assert.NotEmpty(t, streamInstance.Finalizers)
			assert.Len(t, ctx.Conditions(), 1)
			assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
			assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
			assert.Equal(t, string(workflow.StreamInstanceNotRemoved), ctx.Conditions()[0].Reason)
			assert.Equal(t, "failed to delete instance", ctx.Conditions()[0].Message)
		})
	})

	t.Run("should handle deletion when deletion protection is enabled", func(t *testing.T) {
		t.Run("should successfully transition when delete resource annotation is set", func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "my-project-id",
				},
			}
			streamInstance := &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-stream-processing-instance",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: pointer.MakePtr(metav1.Now()),
					Annotations: map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyDelete,
					},
				},
				Spec: akov2.AtlasStreamInstanceSpec{
					Name: "instance-0",
					Config: akov2.Config{
						Provider: "AWS",
						Region:   "FRANKFURT_DEU",
						Tier:     "SP30",
					},
					Project: common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ConnectionRegistry: []common.ResourceRefNamespaced{
						{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
				Status: status.AtlasStreamInstanceStatus{
					ID:        "instance-0-id",
					Hostnames: []string{"mdb://host1", "mdb://host2"},
					Connections: []status.StreamConnection{
						{
							Name: "sample-connection",
							ResourceRef: common.ResourceRefNamespaced{
								Name:      "my-sample-connection",
								Namespace: "default",
							},
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, streamInstance).
				Build()

			reconciler := &AtlasStreamsInstanceReconciler{
				Client:                   k8sClient,
				Log:                      zaptest.NewLogger(t).Sugar(),
				ObjectDeletionProtection: true,
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamWorkspace(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamWorkspaceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamWorkspaceExecute(mock.AnythingOfType("admin.DeleteStreamWorkspaceApiRequest")).
				Return(
					&http.Response{},
					nil,
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312009: &admin.APIClient{
						StreamsApi: streamsAPI,
					},
				},
			}

			result, err := reconciler.delete(ctx, project, streamInstance)
			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{}, result)
			assert.Empty(t, streamInstance.Finalizers)
		})

		t.Run("should successfully transition when deleting a stream instance in Atlas", func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
				Status: status.AtlasProjectStatus{
					ID: "my-project-id",
				},
			}
			streamInstance := &akov2.AtlasStreamInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-stream-processing-instance",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: pointer.MakePtr(metav1.Now()),
					Annotations: map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
					},
				},
				Spec: akov2.AtlasStreamInstanceSpec{
					Name: "instance-0",
					Config: akov2.Config{
						Provider: "AWS",
						Region:   "FRANKFURT_DEU",
						Tier:     "SP30",
					},
					Project: common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ConnectionRegistry: []common.ResourceRefNamespaced{
						{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
				Status: status.AtlasStreamInstanceStatus{
					ID:        "instance-0-id",
					Hostnames: []string{"mdb://host1", "mdb://host2"},
					Connections: []status.StreamConnection{
						{
							Name: "sample-connection",
							ResourceRef: common.ResourceRefNamespaced{
								Name:      "my-sample-connection",
								Namespace: "default",
							},
						},
					},
				},
			}
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, streamInstance).
				Build()

			reconciler := &AtlasStreamsInstanceReconciler{
				Client:                   k8sClient,
				Log:                      zaptest.NewLogger(t).Sugar(),
				ObjectDeletionProtection: true,
			}
			ctx := &workflow.Context{
				Context: context.Background(),
			}

			result, err := reconciler.delete(ctx, project, streamInstance)
			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{}, result)
			assert.Empty(t, streamInstance.Finalizers)
		})
	})

	t.Run("should fail transition when unable to clean finalizer", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-stream-processing-instance",
				Namespace:         "default",
				Finalizers:        []string{customresource.FinalizerLabel},
				DeletionTimestamp: pointer.MakePtr(metav1.Now()),
				Annotations: map[string]string{
					customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
				},
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "FRANKFURT_DEU",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
				Connections: []status.StreamConnection{
					{
						Name: "sample-connection",
						ResourceRef: common.ResourceRefNamespaced{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		_, err := reconciler.delete(ctx, project, streamInstance)
		assert.Error(t, err)
		assert.NotEmpty(t, streamInstance.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotRemoved), ctx.Conditions()[0].Reason)
		assert.Contains(t, ctx.Conditions()[0].Message, `atlasstreaminstances.atlas.mongodb.com "my-stream-processing-instance" not found`)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("should update a stream instance in Atlas", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
				Connections: []status.StreamConnection{
					{
						Name: "sample-connection",
						ResourceRef: common.ResourceRefNamespaced{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamWorkspace(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsDataProcessRegion")).
			Return(admin.UpdateStreamWorkspaceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamWorkspaceExecute(mock.AnythingOfType("admin.UpdateStreamWorkspaceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "DUBLIN_IRL",
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					Connections: &[]admin.StreamsConnection{
						{
							Name: pointer.MakePtr("sample-connection"),
							Type: pointer.MakePtr("Sample"),
						},
					},
					GroupId: pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}

		result, err := reconciler.update(ctx, project, streamInstance)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceSetupInProgress), ctx.Conditions()[0].Reason)
		assert.Equal(t, "configuring stream instance in Atlas", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to in-progress state when fail to update a stream instance in Atlas", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "FRANKFURT_DEU",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
				Connections: []status.StreamConnection{
					{
						Name: "sample-connection",
						ResourceRef: common.ResourceRefNamespaced{
							Name:      "my-sample-connection",
							Namespace: "default",
						},
					},
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamWorkspace(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsDataProcessRegion")).
			Return(admin.UpdateStreamWorkspaceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamWorkspaceExecute(mock.AnythingOfType("admin.UpdateStreamWorkspaceApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to update instance"),
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}

		_, err := reconciler.update(ctx, project, streamInstance)
		assert.Error(t, err)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceNotCreated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to update instance", ctx.Conditions()[0].Message)
	})
}

func TestCreateConnections(t *testing.T) {
	t.Run("should add an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamConnection(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.CreateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamConnectionExecute(mock.AnythingOfType("admin.CreateStreamConnectionApiRequest")).
			Return(
				&admin.StreamsConnection{
					Name: pointer.MakePtr("sample-connection"),
					Type: pointer.MakePtr("Sample"),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := createConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		})
		assert.NoError(t, err)
	})

	t.Run("should return error when fail adding an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamConnection(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.CreateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamConnectionExecute(mock.AnythingOfType("admin.CreateStreamConnectionApiRequest")).
			Return(nil, &http.Response{}, errors.New("failed to create connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := createConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		})
		assert.ErrorContains(t, err, "failed to create connection")
	})

	t.Run("should return error when fail mapping an instance connection", func(t *testing.T) { //nolint:dupl
		ctx := &workflow.Context{}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := createConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return nil, errors.New("failed to map connection")
		})
		assert.ErrorContains(t, err, "failed to map connection")
	})
}

func TestUpdateConnections(t *testing.T) {
	t.Run("should update an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.UpdateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamConnectionExecute(mock.AnythingOfType("admin.UpdateStreamConnectionApiRequest")).
			Return(
				&admin.StreamsConnection{
					Name: pointer.MakePtr("sample-connection"),
					Type: pointer.MakePtr("Sample"),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := updateConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		})
		assert.NoError(t, err)
	})

	t.Run("should return error when fail updating an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.UpdateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamConnectionExecute(mock.AnythingOfType("admin.UpdateStreamConnectionApiRequest")).
			Return(nil, &http.Response{}, errors.New("failed to update connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := updateConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		})
		assert.ErrorContains(t, err, "failed to update connection")
	})

	t.Run("should return error when fail mapping an instance connection", func(t *testing.T) { //nolint:dupl
		ctx := &workflow.Context{}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*akov2.AtlasStreamConnection{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-sample-connection",
					Namespace: "default",
				},
				Spec: akov2.AtlasStreamConnectionSpec{
					Name:           "sample-connection",
					ConnectionType: "Sample",
				},
			},
		}

		err := updateConnections(ctx, project, streamInstance, connections, func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return nil, errors.New("failed to map connection")
		})
		assert.ErrorContains(t, err, "failed to map connection")
	})
}

func TestDeleteConnections(t *testing.T) {
	t.Run("should delete an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			DeleteStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection").
			Return(admin.DeleteStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			DeleteStreamConnectionExecute(mock.AnythingOfType("admin.DeleteStreamConnectionApiRequest")).
			Return(
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*admin.StreamsConnection{
			{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			},
		}

		err := deleteConnections(ctx, project, streamInstance, connections)
		assert.NoError(t, err)
	})

	t.Run("should return error when fail deleting an instance connection", func(t *testing.T) {
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			DeleteStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection").
			Return(admin.DeleteStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			DeleteStreamConnectionExecute(mock.AnythingOfType("admin.DeleteStreamConnectionApiRequest")).
			Return(&http.Response{}, errors.New("failed to delete connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312009: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: akov2.AtlasProjectSpec{
				Name: "my-project",
			},
			Status: status.AtlasProjectStatus{
				ID: "my-project-id",
			},
		}
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				Config: akov2.Config{
					Provider: "AWS",
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection",
						Namespace: "default",
					},
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		connections := []*admin.StreamsConnection{
			{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			},
		}

		err := deleteConnections(ctx, project, streamInstance, connections)
		assert.ErrorContains(t, err, "failed to delete connection")
	})
}
