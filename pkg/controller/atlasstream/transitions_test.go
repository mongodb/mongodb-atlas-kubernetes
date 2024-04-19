package atlasstream

import (
	"context"
	"errors"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("should transition to a ready state when creates a stream instance in Atlas", func(t *testing.T) {
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamInstance(context.Background(), "my-project-id", mock.AnythingOfType("*admin.StreamsTenant")).
			Return(admin.CreateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamInstanceExecute(mock.AnythingOfType("admin.CreateStreamInstanceApiRequest")).
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
			SdkClient: &admin.APIClient{
				StreamsApi: streamsAPI,
			},
		}

		connectionMapper := func(conn *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		}

		result, err := reconciler.create(ctx, project, streamInstance, connectionMapper)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.Len(t, ctx.Conditions(), 2)
		assert.Equal(t, status.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
		assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[1].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[1].Status)
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamInstance(context.Background(), "my-project-id", mock.AnythingOfType("*admin.StreamsTenant")).
			Return(admin.CreateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamInstanceExecute(mock.AnythingOfType("admin.CreateStreamInstanceApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to create instance"),
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClient: &admin.APIClient{
				StreamsApi: streamsAPI,
			},
		}

		connectionMapper := func(conn *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return &admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			}, nil
		}

		result, err := reconciler.create(ctx, project, streamInstance, connectionMapper)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceNotCreated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to create instance", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to a terminate state when fail to map connection resource", func(t *testing.T) {
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		connectionMapper := func(conn *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return nil, errors.New("failed to map sample-connection")
		}

		result, err := reconciler.create(ctx, project, streamInstance, connectionMapper)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceNotCreated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to map sample-connection", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to a terminate state when fail to get connection resource", func(t *testing.T) {
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
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			Build()

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		connectionMapper := func(conn *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
			return nil, nil
		}

		result, err := reconciler.create(ctx, project, streamInstance, connectionMapper)
		assert.NoError(t, err)
		assert.Equal(
			t,
			ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			result,
		)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamInstanceNotCreated), ctx.Conditions()[0].Reason)
		assert.Contains(t, ctx.Conditions()[0].Message, "failed to retrieve connection {my-sample-connection default}:")
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

			reconciler := &InstanceReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamInstance(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamInstanceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamInstanceExecute(mock.AnythingOfType("admin.DeleteStreamInstanceApiRequest")).
				Return(
					nil,
					&http.Response{},
					nil,
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClient: &admin.APIClient{
					StreamsApi: streamsAPI,
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

			reconciler := &InstanceReconciler{
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

			reconciler := &InstanceReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamInstance(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamInstanceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamInstanceExecute(mock.AnythingOfType("admin.DeleteStreamInstanceApiRequest")).
				Return(
					nil,
					&http.Response{},
					errors.New("failed to delete instance"),
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClient: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			}

			result, err := reconciler.delete(ctx, project, streamInstance)
			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			}, result)
			assert.NotEmpty(t, streamInstance.Finalizers)
			assert.Len(t, ctx.Conditions(), 1)
			assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[0].Type)
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

			reconciler := &InstanceReconciler{
				Client:                   k8sClient,
				Log:                      zaptest.NewLogger(t).Sugar(),
				ObjectDeletionProtection: true,
			}
			streamsAPI := mockadmin.NewStreamsApi(t)
			streamsAPI.EXPECT().
				DeleteStreamInstance(context.Background(), "my-project-id", "instance-0").
				Return(admin.DeleteStreamInstanceApiRequest{ApiService: streamsAPI})
			streamsAPI.EXPECT().
				DeleteStreamInstanceExecute(mock.AnythingOfType("admin.DeleteStreamInstanceApiRequest")).
				Return(
					nil,
					&http.Response{},
					nil,
				)
			ctx := &workflow.Context{
				Context: context.Background(),
				SdkClient: &admin.APIClient{
					StreamsApi: streamsAPI,
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

			reconciler := &InstanceReconciler{
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		result, err := reconciler.delete(ctx, project, streamInstance)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{
			RequeueAfter: workflow.DefaultRetry,
		}, result)
		assert.NotEmpty(t, streamInstance.Finalizers)
		assert.Len(t, ctx.Conditions(), 1)
		assert.Equal(t, status.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.AtlasFinalizerNotRemoved), ctx.Conditions()[0].Reason)
		assert.Contains(t, ctx.Conditions()[0].Message, "failed to get &{{%!t(string=) %!t(string=)} {%!t(string=my-stream-processing-instance)")
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamInstance(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsDataProcessRegion")).
			Return(admin.UpdateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamInstanceExecute(mock.AnythingOfType("admin.UpdateStreamInstanceApiRequest")).
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
			SdkClient: &admin.APIClient{
				StreamsApi: streamsAPI,
			},
		}

		err := reconciler.update(ctx, project, streamInstance)
		assert.NoError(t, err)
	})

	t.Run("should return error when fail to update a stream instance in Atlas", func(t *testing.T) {
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

		reconciler := &InstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			UpdateStreamInstance(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsDataProcessRegion")).
			Return(admin.UpdateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamInstanceExecute(mock.AnythingOfType("admin.UpdateStreamInstanceApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to update instance"),
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClient: &admin.APIClient{
				StreamsApi: streamsAPI,
			},
		}

		err := reconciler.update(ctx, project, streamInstance)
		assert.ErrorContains(t, err, "failed to update instance")
	})
}
