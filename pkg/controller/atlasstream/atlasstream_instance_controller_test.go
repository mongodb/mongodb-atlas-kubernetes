package atlasstream

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/indexer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
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

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestReconcile(t *testing.T) {
	t.Run("should terminate silently when resource is not found", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})
}

func TestEnsureAtlasStreamsInstance(t *testing.T) {
	t.Run("should skip reconciliation when annotation is set", func(t *testing.T) {
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
				Annotations: map[string]string{
					customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
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
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("should transition to invalid state when resource version is invalid", func(t *testing.T) {
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
				Labels: map[string]string{
					customresource.ResourceVersion: "no-semver",
				},
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 2)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[1].Status)
		assert.Equal(t, string(workflow.AtlasResourceVersionIsInvalid), conditions[1].Reason)
		assert.Equal(t, "no-semver is not a valid semver version for label mongodb.com/atlas-resource-version", conditions[1].Message)
	})

	t.Run("should transition to unsupported state when resource is not supported by platform", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
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
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.AtlasGovUnsupported), conditions[2].Reason)
		assert.Equal(t, "the AtlasStreamInstance is not supported by Atlas for government", conditions[2].Message)
	})

	t.Run("should transition to terminate state when project resource is not found", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
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
					Name:      "my-stream-processing-instance",
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.Internal), conditions[2].Reason)
		assert.Equal(t, "atlasprojects.atlas.mongodb.com \"my-project\" not found", conditions[2].Message)
	})

	t.Run("should transition to terminate state when unable to configure sdk client", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to configure sdk client")
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.AtlasAPIAccessNotConfigured), conditions[2].Reason)
		assert.Equal(t, "failed to configure sdk client", conditions[2].Message)
	})

	t.Run("should transition to terminate state when Atlas API fails", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().GetStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.GetStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().GetStreamInstanceExecute(mock.AnythingOfType("admin.GetStreamInstanceApiRequest")).
			Return(nil, &http.Response{}, errors.New("failed to get instance"))

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{StreamsApi: streamsAPI}, "org-id", nil
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
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

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.Internal), conditions[2].Reason)
		assert.Equal(t, "failed to get instance", conditions[2].Message)
	})

	t.Run("should transition to ready state when everything is in sync", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().GetStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.GetStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().GetStreamInstanceExecute(mock.AnythingOfType("admin.GetStreamInstanceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "DUBLIN_IRL",
					},
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					GroupId:   pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().ListStreamConnections(context.Background(), "my-project-id", "instance-0").
			Return(admin.ListStreamConnectionsApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().ListStreamConnectionsExecute(mock.AnythingOfType("admin.ListStreamConnectionsApiRequest")).
			Return(
				&admin.PaginatedApiStreamsConnection{
					Results:    nil,
					TotalCount: pointer.MakePtr(0),
				},
				&http.Response{},
				nil,
			)

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{StreamsApi: streamsAPI}, "org-id", nil
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[2].Status)
	})

	t.Run("should transition to in-progress state when creating new instance", func(t *testing.T) {
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
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().GetStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.GetStreamInstanceApiRequest{ApiService: streamsAPI})
		notFound := admin.ApiError{}
		notFound.SetError(404)
		notFound.SetErrorCode(instanceNotFound)
		apiError := admin.GenericOpenAPIError{}
		apiError.SetModel(notFound)
		streamsAPI.EXPECT().GetStreamInstanceExecute(mock.AnythingOfType("admin.GetStreamInstanceApiRequest")).
			Return(
				nil,
				&http.Response{},
				&apiError,
			)
		streamsAPI.EXPECT().CreateStreamInstance(context.Background(), "my-project-id", mock.AnythingOfType("*admin.StreamsTenant")).
			Return(admin.CreateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().CreateStreamInstanceExecute(mock.AnythingOfType("admin.CreateStreamInstanceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "DUBLIN_IRL",
					},
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					GroupId:   pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{StreamsApi: streamsAPI}, "org-id", nil
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.StreamInstanceSetupInProgress), conditions[2].Reason)
		assert.Equal(t, "configuring stream instance in Atlas", conditions[2].Message)
	})

	t.Run("should transition succeed when deleting an instance", func(t *testing.T) {
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
					Region:   "DUBLIN_IRL",
					Tier:     "SP30",
				},
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().GetStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.GetStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().GetStreamInstanceExecute(mock.AnythingOfType("admin.GetStreamInstanceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "DUBLIN_IRL",
					},
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					GroupId:   pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().DeleteStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.DeleteStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().DeleteStreamInstanceExecute(mock.AnythingOfType("admin.DeleteStreamInstanceApiRequest")).
			Return(
				nil,
				&http.Response{},
				nil,
			)

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{StreamsApi: streamsAPI}, "org-id", nil
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)

		assert.Error(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
	})

	t.Run("should transition to ready state when updating an instance", func(t *testing.T) {
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
			},
			Status: status.AtlasStreamInstanceStatus{
				ID:        "instance-0-id",
				Hostnames: []string{"mdb://host1", "mdb://host2"},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance).
			WithStatusSubresource(streamInstance).
			Build()

		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().GetStreamInstance(context.Background(), "my-project-id", "instance-0").
			Return(admin.GetStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().GetStreamInstanceExecute(mock.AnythingOfType("admin.GetStreamInstanceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "FRANKFURT_DEU",
					},
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					GroupId:   pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().UpdateStreamInstance(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsDataProcessRegion")).
			Return(admin.UpdateStreamInstanceApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().UpdateStreamInstanceExecute(mock.AnythingOfType("admin.UpdateStreamInstanceApiRequest")).
			Return(
				&admin.StreamsTenant{
					Id:   pointer.MakePtr("instance-0-id"),
					Name: pointer.MakePtr("instance-0"),
					DataProcessRegion: &admin.StreamsDataProcessRegion{
						CloudProvider: "AWS",
						Region:        "DUBLIN_IRL",
					},
					StreamConfig: &admin.StreamConfig{
						Tier: pointer.MakePtr("SP30"),
					},
					Hostnames: pointer.MakePtr([]string{"mdb://host1", "mdb://host2"}),
					GroupId:   pointer.MakePtr("my-project-id"),
				},
				&http.Response{},
				nil,
			)

		reconciler := &AtlasStreamsInstanceReconciler{
			Client:        k8sClient,
			Log:           zaptest.NewLogger(t).Sugar(),
			EventRecorder: record.NewFakeRecorder(1),
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{StreamsApi: streamsAPI}, "org-id", nil
				},
			},
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      "my-stream-processing-instance",
					Namespace: "default",
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: workflow.DefaultRetry}, result)

		assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance))
		conditions := streamInstance.Status.GetConditions()
		assert.Len(t, conditions, 3)
		assert.Equal(t, status.ReadyType, conditions[0].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[0].Status)
		assert.Empty(t, conditions[0].Reason)
		assert.Empty(t, conditions[0].Message)
		assert.Equal(t, status.ResourceVersionStatus, conditions[1].Type)
		assert.Equal(t, corev1.ConditionTrue, conditions[1].Status)
		assert.Equal(t, status.StreamInstanceReadyType, conditions[2].Type)
		assert.Equal(t, corev1.ConditionFalse, conditions[2].Status)
		assert.Equal(t, string(workflow.StreamInstanceSetupInProgress), conditions[2].Reason)
		assert.Equal(t, "configuring stream instance in Atlas", conditions[2].Message)
	})
}

func TestFindStreamInstancesForStreamConnection(t *testing.T) {
	t.Run("should fail when watching wrong object", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasStreamsInstanceReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Nil(t, reconciler.findStreamInstancesForStreamConnection(context.Background(), &akov2.AtlasProject{}))
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zap.WarnLevel, logs.All()[0].Level)
		assert.Equal(t, "watching AtlasStreamConnection but got *v1.AtlasProject", logs.All()[0].Message)
	})

	t.Run("should return slice of requests for instances", func(t *testing.T) {
		instance1 := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance1",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance1",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "connection",
						Namespace: "default",
					},
				},
			},
		}
		instance2 := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance2",
				Namespace: "other-ns",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance2",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "connection",
						Namespace: "default",
					},
				},
			},
		}
		connection := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "connection",
				Namespace: "default",
			},
			Status: status.AtlasStreamConnectionStatus{
				Instances: []common.ResourceRefNamespaced{
					{
						Namespace: "ns1",
						Name:      "instance1",
					},
					{
						Namespace: "ns2",
						Name:      "instance2",
					},
				},
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection, instance1, instance2).
			WithIndex(
				&akov2.AtlasStreamInstance{},
				indexer.AtlasStreamInstancesByConnectionRegistry,
				indexer.AtlasStreamInstancesByConnectionRegistryIndices,
			).
			Build()
		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		requests := reconciler.findStreamInstancesForStreamConnection(context.Background(), connection)
		assert.Equal(
			t,
			[]ctrl.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "instance1",
					},
				},
				{
					NamespacedName: types.NamespacedName{
						Namespace: "other-ns",
						Name:      "instance2",
					},
				},
			},
			requests,
		)
	})
}
