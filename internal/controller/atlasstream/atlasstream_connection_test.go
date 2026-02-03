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
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestHandleConnectionRegistry(t *testing.T) {
	t.Run("should handle connection registry transitions", func(t *testing.T) {
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
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection1",
						Namespace: "default",
					},
					{
						Name:      "my-sample-connection2",
						Namespace: "default",
					},
					{
						Name:      "my-sample-connection4",
						Namespace: "default",
					},
				},
			},
		}
		streamConnection1 := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection1",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection1",
				ConnectionType: "Sample",
			},
		}
		streamConnection2 := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection2",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection2",
				ConnectionType: "Sample",
			},
		}
		streamConnection4 := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection4",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection4",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamInstance, streamConnection1, streamConnection2, streamConnection4).
			Build()
		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().
			CreateStreamConnection(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.CreateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamConnectionExecute(mock.AnythingOfType("admin.CreateStreamConnectionApiRequest")).
			Return(
				&admin.StreamsConnection{
					Name: pointer.MakePtr("sample-connection1"),
					Type: pointer.MakePtr("Sample"),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().
			UpdateStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection2", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.UpdateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamConnectionExecute(mock.AnythingOfType("admin.UpdateStreamConnectionApiRequest")).
			Return(
				&admin.StreamsConnection{
					Name: pointer.MakePtr("sample-connection2"),
					Type: pointer.MakePtr("Sample"),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().
			DeleteStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection3").
			Return(admin.DeleteStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			DeleteStreamConnectionExecute(mock.AnythingOfType("admin.DeleteStreamConnectionApiRequest")).
			Return(
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().ListStreamConnections(context.Background(), "my-project-id", "instance-0").
			Return(admin.ListStreamConnectionsApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().ListStreamConnectionsExecute(mock.AnythingOfType("admin.ListStreamConnectionsApiRequest")).
			Return(
				&admin.PaginatedApiStreamsConnection{
					Results: &[]admin.StreamsConnection{
						{
							Name:        pointer.MakePtr("sample-connection2"),
							Type:        pointer.MakePtr("Cluster"),
							ClusterName: pointer.MakePtr("my-cluster"),
							DbRoleToExecute: &admin.DBRoleToExecute{
								Role: pointer.MakePtr("readWrite"),
								Type: pointer.MakePtr("BUILT_IN"),
							},
						},
						{
							Name: pointer.MakePtr("sample-connection3"),
							Type: pointer.MakePtr("Sample"),
						},
						{
							Name: pointer.MakePtr("sample-connection4"),
							Type: pointer.MakePtr("Sample"),
						},
					},
					TotalCount: pointer.MakePtr(2),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		atlasInstance := &admin.StreamsTenant{
			Name: pointer.MakePtr("instance-0"),
		}

		result, err := reconciler.handleConnectionRegistry(ctx, project, streamInstance, atlasInstance)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.Equal(t, api.ReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[0].Status)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[1].Type)
		assert.Equal(t, corev1.ConditionTrue, ctx.Conditions()[1].Status)
	})

	t.Run("should transition to terminate state when failing to sort operations", func(t *testing.T) {
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
		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}
		streamsAPI := mockadmin.NewStreamsApi(t)
		streamsAPI.EXPECT().ListStreamConnections(context.Background(), "my-project-id", "instance-0").
			Return(admin.ListStreamConnectionsApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().ListStreamConnectionsExecute(mock.AnythingOfType("admin.ListStreamConnectionsApiRequest")).
			Return(
				&admin.PaginatedApiStreamsConnection{
					Results: &[]admin.StreamsConnection{
						{
							Name:        pointer.MakePtr("sample-connection2"),
							Type:        pointer.MakePtr("Cluster"),
							ClusterName: pointer.MakePtr("my-cluster"),
							DbRoleToExecute: &admin.DBRoleToExecute{
								Role: pointer.MakePtr("readWrite"),
								Type: pointer.MakePtr("BUILT_IN"),
							},
						},
						{
							Name: pointer.MakePtr("sample-connection3"),
							Type: pointer.MakePtr("Sample"),
						},
					},
					TotalCount: pointer.MakePtr(2),
				},
				&http.Response{},
				nil,
			)
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		atlasInstance := &admin.StreamsTenant{
			Name: pointer.MakePtr("instance-0"),
		}

		_, err := reconciler.handleConnectionRegistry(ctx, project, streamInstance, atlasInstance)
		assert.Error(t, err)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamConnectionNotConfigured), ctx.Conditions()[0].Reason)
		assert.Contains(t, ctx.Conditions()[0].Message, "failed to retrieve connection {my-sample-connection default}")
	})

	t.Run("should transition to terminate state when failing to create connections", func(t *testing.T) {
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
		streamsAPI.EXPECT().
			CreateStreamConnection(context.Background(), "my-project-id", "instance-0", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.CreateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			CreateStreamConnectionExecute(mock.AnythingOfType("admin.CreateStreamConnectionApiRequest")).
			Return(nil, &http.Response{}, errors.New("failed to create connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		atlasInstance := &admin.StreamsTenant{
			Name:        pointer.MakePtr("instance-0"),
			Connections: &[]admin.StreamsConnection{},
		}

		_, err := reconciler.handleConnectionRegistry(ctx, project, streamInstance, atlasInstance)
		assert.Error(t, err)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamConnectionNotCreated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to create connection", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to terminate state when failing to update connections", func(t *testing.T) {
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
		streamsAPI.EXPECT().ListStreamConnections(context.Background(), "my-project-id", "instance-0").
			Return(admin.ListStreamConnectionsApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().ListStreamConnectionsExecute(mock.AnythingOfType("admin.ListStreamConnectionsApiRequest")).
			Return(
				&admin.PaginatedApiStreamsConnection{
					Results: &[]admin.StreamsConnection{
						{
							Name:        pointer.MakePtr("sample-connection"),
							Type:        pointer.MakePtr("Cluster"),
							ClusterName: pointer.MakePtr("my-cluster"),
							DbRoleToExecute: &admin.DBRoleToExecute{
								Role: pointer.MakePtr("readWrite"),
								Type: pointer.MakePtr("BUILT_IN"),
							},
						},
					},
					TotalCount: pointer.MakePtr(1),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().
			UpdateStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection", mock.AnythingOfType("*admin.StreamsConnection")).
			Return(admin.UpdateStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			UpdateStreamConnectionExecute(mock.AnythingOfType("admin.UpdateStreamConnectionApiRequest")).
			Return(nil, &http.Response{}, errors.New("failed to update connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		atlasInstance := &admin.StreamsTenant{
			Name: pointer.MakePtr("instance-0"),
		}

		_, err := reconciler.handleConnectionRegistry(ctx, project, streamInstance, atlasInstance)
		assert.Error(t, err)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamConnectionNotUpdated), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to update connection", ctx.Conditions()[0].Message)
	})

	t.Run("should transition to terminate state when failing to delete connections", func(t *testing.T) {
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
		streamsAPI.EXPECT().ListStreamConnections(context.Background(), "my-project-id", "instance-0").
			Return(admin.ListStreamConnectionsApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().ListStreamConnectionsExecute(mock.AnythingOfType("admin.ListStreamConnectionsApiRequest")).
			Return(
				&admin.PaginatedApiStreamsConnection{
					Results: &[]admin.StreamsConnection{
						{
							Name:        pointer.MakePtr("sample-connection"),
							Type:        pointer.MakePtr("Cluster"),
							ClusterName: pointer.MakePtr("my-cluster"),
							DbRoleToExecute: &admin.DBRoleToExecute{
								Role: pointer.MakePtr("readWrite"),
								Type: pointer.MakePtr("BUILT_IN"),
							},
						},
					},
					TotalCount: pointer.MakePtr(1),
				},
				&http.Response{},
				nil,
			)
		streamsAPI.EXPECT().
			DeleteStreamConnection(context.Background(), "my-project-id", "instance-0", "sample-connection").
			Return(admin.DeleteStreamConnectionApiRequest{ApiService: streamsAPI})
		streamsAPI.EXPECT().
			DeleteStreamConnectionExecute(mock.AnythingOfType("admin.DeleteStreamConnectionApiRequest")).
			Return(&http.Response{}, errors.New("failed to delete connection"))
		ctx := &workflow.Context{
			Context: context.Background(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312013: &admin.APIClient{
					StreamsApi: streamsAPI,
				},
			},
		}
		atlasInstance := &admin.StreamsTenant{
			Name: pointer.MakePtr("instance-0"),
		}

		_, err := reconciler.handleConnectionRegistry(ctx, project, streamInstance, atlasInstance)
		assert.Error(t, err)
		assert.Equal(t, api.StreamInstanceReadyType, ctx.Conditions()[0].Type)
		assert.Equal(t, corev1.ConditionFalse, ctx.Conditions()[0].Status)
		assert.Equal(t, string(workflow.StreamConnectionNotRemoved), ctx.Conditions()[0].Reason)
		assert.Equal(t, "failed to delete connection", ctx.Conditions()[0].Message)
	})
}

func TestSortConnectionRegistryTasks(t *testing.T) {
	t.Run("should return error when unable to retrieve resource from kube", func(t *testing.T) {
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
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()
		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{}

		ops, err := reconciler.sortConnectionRegistryTasks(ctx, streamInstance, []admin.StreamsConnection{})
		assert.ErrorContains(t, err, "failed to retrieve connection {my-sample-connection default}")
		assert.Nil(t, ops)
	})

	t.Run("should sort operation to transition connections", func(t *testing.T) {
		streamInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-stream-processing-instance",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Name: "instance-0",
				ConnectionRegistry: []common.ResourceRefNamespaced{
					{
						Name:      "my-sample-connection1",
						Namespace: "default",
					},
					{
						Name:      "my-sample-connection2",
						Namespace: "default",
					},
				},
			},
		}
		streamConnection1 := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection1",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection1",
				ConnectionType: "Sample",
			},
		}
		streamConnection2 := &akov2.AtlasStreamConnection{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-sample-connection2",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection2",
				ConnectionType: "Sample",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(streamInstance, streamConnection1, streamConnection2).
			Build()
		reconciler := &AtlasStreamsInstanceReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{}
		atlasInstanceConnections := []admin.StreamsConnection{
			{
				Name:        pointer.MakePtr("sample-connection2"),
				Type:        pointer.MakePtr("Cluster"),
				ClusterName: pointer.MakePtr("my-cluster"),
				DbRoleToExecute: &admin.DBRoleToExecute{
					Role: pointer.MakePtr("readWrite"),
					Type: pointer.MakePtr("BUILT_IN"),
				},
			},
			{
				Name: pointer.MakePtr("sample-connection3"),
				Type: pointer.MakePtr("Sample"),
			},
		}

		ops, err := reconciler.sortConnectionRegistryTasks(ctx, streamInstance, atlasInstanceConnections)
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]*akov2.AtlasStreamConnection{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "my-sample-connection1",
						Namespace:       "default",
						ResourceVersion: "999",
					},
					Spec: akov2.AtlasStreamConnectionSpec{
						Name:           "sample-connection1",
						ConnectionType: "Sample",
					},
				},
			},
			ops.Create,
		)
		assert.Equal(
			t,
			[]*akov2.AtlasStreamConnection{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:            "my-sample-connection2",
						Namespace:       "default",
						ResourceVersion: "999",
					},
					Spec: akov2.AtlasStreamConnectionSpec{
						Name:           "sample-connection2",
						ConnectionType: "Sample",
					},
				},
			},
			ops.Update,
		)
		assert.Equal(
			t,
			[]*admin.StreamsConnection{
				{
					Name: pointer.MakePtr("sample-connection3"),
					Type: pointer.MakePtr("Sample"),
				},
			},
			ops.Delete,
		)
	})
}

func TestStreamConnectionToAtlas(t *testing.T) {
	t.Run("should map a sample connection configuration", func(t *testing.T) {
		k8sClient := fake.NewClientBuilder().
			Build()
		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection",
				ConnectionType: "Sample",
			},
		}

		mapFunc := streamConnectionToAtlas(context.Background(), k8sClient)
		conn, err := mapFunc(&akoConnection)
		assert.NoError(t, err)
		assert.Equal(
			t,
			&admin.StreamsConnection{
				Name: pointer.MakePtr("sample-connection"),
				Type: pointer.MakePtr("Sample"),
			},
			conn,
		)
	})

	t.Run("should map a cluster connection configuration", func(t *testing.T) {
		k8sClient := fake.NewClientBuilder().
			Build()
		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "cluster-connection",
				ConnectionType: "Cluster",
				ClusterConfig: &akov2.ClusterConnectionConfig{
					Name: "my-cluster",
					Role: akov2.StreamsClusterDBRole{
						Name:     "read",
						RoleType: "BUILT_IN",
					},
				},
			},
		}

		mapFunc := streamConnectionToAtlas(context.Background(), k8sClient)
		conn, err := mapFunc(&akoConnection)
		assert.NoError(t, err)
		assert.Equal(
			t,
			&admin.StreamsConnection{
				Name:        pointer.MakePtr("cluster-connection"),
				Type:        pointer.MakePtr("Cluster"),
				ClusterName: pointer.MakePtr("my-cluster"),
				DbRoleToExecute: &admin.DBRoleToExecute{
					Role: pointer.MakePtr("read"),
					Type: pointer.MakePtr("BUILT_IN"),
				},
			},
			conn,
		)
	})

	t.Run("should map a kafka connection configuration", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		secretCreds := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-creds",
			},
			Data: map[string][]byte{
				"username": []byte("my-user"),
				"password": []byte("my-pass"),
			},
			Type: "opaque",
		}
		secretCert := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-cert",
			},
			Data: map[string][]byte{
				"certificate": []byte("hash"),
			},
			Type: "opaque",
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(&secretCreds, &secretCert).
			Build()

		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "kafka-connection",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Mechanism: "SCRAM-512",
						Credentials: common.ResourceRefNamespaced{
							Name:      "my-kafka-creds",
							Namespace: "default",
						},
					},
					BootstrapServers: "kafka:server1,kafka:server2",
					Security: akov2.StreamsKafkaSecurity{
						Protocol: "SSL",
						Certificate: common.ResourceRefNamespaced{
							Name:      "my-kafka-cert",
							Namespace: "default",
						},
					},
					Config: map[string]string{
						"option1": "value1",
					},
				},
			},
		}

		mapFunc := streamConnectionToAtlas(context.Background(), k8sClient)
		conn, err := mapFunc(&akoConnection)
		assert.NoError(t, err)
		assert.Equal(
			t,
			&admin.StreamsConnection{
				Name: pointer.MakePtr("kafka-connection"),
				Type: pointer.MakePtr("Kafka"),
				Authentication: &admin.StreamsKafkaAuthentication{
					Mechanism: pointer.MakePtr("SCRAM-512"),
					Username:  pointer.MakePtr("my-user"),
					Password:  pointer.MakePtr("my-pass"),
				},
				BootstrapServers: pointer.MakePtr("kafka:server1,kafka:server2"),
				Security: &admin.StreamsKafkaSecurity{
					BrokerPublicCertificate: pointer.MakePtr("hash"),
					Protocol:                pointer.MakePtr("SSL"),
				},
				Config: &map[string]string{
					"option1": "value1",
				},
			},
			conn,
		)
	})

	t.Run("should return error to map a kafka configuration when fail to get authentication data", func(t *testing.T) { //nolint:dupl
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		secretCreds := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-creds",
			},
			Data: map[string][]byte{
				"username":    []byte("my-user"),
				"misspelling": []byte("my-pass"),
			},
			Type: "opaque",
		}
		secretCert := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-cert",
			},
			Data: map[string][]byte{
				"certificate": []byte("hash"),
			},
			Type: "opaque",
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(&secretCreds, &secretCert).
			Build()

		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "kafka-connection",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Mechanism: "SCRAM-512",
						Credentials: common.ResourceRefNamespaced{
							Name:      "my-kafka-creds",
							Namespace: "default",
						},
					},
					BootstrapServers: "kafka:server1,kafka:server2",
					Security: akov2.StreamsKafkaSecurity{
						Protocol: "SSL",
						Certificate: common.ResourceRefNamespaced{
							Name:      "my-kafka-cert",
							Namespace: "default",
						},
					},
					Config: map[string]string{
						"option1": "value1",
					},
				},
			},
		}

		mapFunc := streamConnectionToAtlas(context.Background(), k8sClient)
		conn, err := mapFunc(&akoConnection)
		assert.ErrorContains(t, err, "key password is not present in the secret default/my-kafka-creds")
		assert.Nil(t, conn)
	})

	t.Run("should return error to map a kafka configuration when fail to get certificate data", func(t *testing.T) { //nolint:dupl
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		secretCreds := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-creds",
			},
			Data: map[string][]byte{
				"username": []byte("my-user"),
				"password": []byte("my-pass"),
			},
			Type: "opaque",
		}
		secretCert := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-kafka-cert",
			},
			Data: map[string][]byte{
				"misspelling": []byte("hash"),
			},
			Type: "opaque",
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(&secretCreds, &secretCert).
			Build()

		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "kafka-connection",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Mechanism: "SCRAM-512",
						Credentials: common.ResourceRefNamespaced{
							Name:      "my-kafka-creds",
							Namespace: "default",
						},
					},
					BootstrapServers: "kafka:server1,kafka:server2",
					Security: akov2.StreamsKafkaSecurity{
						Protocol: "SSL",
						Certificate: common.ResourceRefNamespaced{
							Name:      "my-kafka-cert",
							Namespace: "default",
						},
					},
					Config: map[string]string{
						"option1": "value1",
					},
				},
			},
		}

		mapFunc := streamConnectionToAtlas(context.Background(), k8sClient)
		conn, err := mapFunc(&akoConnection)
		assert.ErrorContains(t, err, "key certificate is not present in the secret default/my-kafka-cert")
		assert.Nil(t, conn)
	})
}

func TestGetSecretData(t *testing.T) {
	t.Run("should return error when secret doesn't exist", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		data, err := getSecretData(context.Background(), k8sClient, client.ObjectKey{Namespace: "default", Name: "my-secret"}, "certificate")
		assert.ErrorContains(t, err, "failed to retrieve secret default/my-secret: secrets \"my-secret\" not found")
		assert.Nil(t, data)
	})

	t.Run("should return error when keys are not found in the secret", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-secret",
			},
			Data: map[string][]byte{
				"wrong-key": {'a'},
			},
			Type: "opaque",
		}

		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(&secret).
			Build()

		data, err := getSecretData(context.Background(), k8sClient, client.ObjectKey{Namespace: "default", Name: "my-secret"}, "certificate")
		assert.ErrorContains(t, err, "key certificate is not present in the secret default/my-secret")
		assert.Nil(t, data)
	})

	t.Run("should return keys from secret", func(t *testing.T) {
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))

		secret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "my-secret",
			},
			Data: map[string][]byte{
				"certificate": []byte("hash"),
			},
			Type: "opaque",
		}

		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(&secret).
			Build()

		data, err := getSecretData(context.Background(), k8sClient, client.ObjectKey{Namespace: "default", Name: "my-secret"}, "certificate")
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"certificate": "hash"}, data)
	})
}

func TestHasStreamConnectionChanged(t *testing.T) {
	t.Run("should return an error when unable to map connection data", func(t *testing.T) {
		atlasConnection := admin.StreamsConnection{}
		akoConnection := akov2.AtlasStreamConnection{}

		ok, err := hasStreamConnectionChanged(
			&akoConnection,
			&atlasConnection,
			func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
				return nil, errors.New("failed to map connection")
			},
		)

		assert.ErrorContains(t, err, "failed to map connection")
		assert.False(t, ok)
	})

	t.Run("should return false when data haven't changed", func(t *testing.T) {
		atlasConnection := admin.StreamsConnection{
			Name: pointer.MakePtr("sample-connection"),
			Type: pointer.MakePtr("Sample"),
		}
		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "sample-connection",
				ConnectionType: "Sample",
			},
		}

		ok, err := hasStreamConnectionChanged(
			&akoConnection,
			&atlasConnection,
			func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
				return &admin.StreamsConnection{
					Name: pointer.MakePtr("sample-connection"),
					Type: pointer.MakePtr("Sample"),
				}, nil
			},
		)

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should return true when data have changed", func(t *testing.T) {
		atlasConnection := admin.StreamsConnection{
			Name: pointer.MakePtr("sample-connection"),
			Type: pointer.MakePtr("Sample"),
		}
		akoConnection := akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "cluster-connection",
				ConnectionType: "Cluster",
				ClusterConfig: &akov2.ClusterConnectionConfig{
					Name: "my-cluster",
					Role: akov2.StreamsClusterDBRole{
						Name:     "read",
						RoleType: "BUILT_IN",
					},
				},
			},
		}

		ok, err := hasStreamConnectionChanged(
			&akoConnection,
			&atlasConnection,
			func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
				return &admin.StreamsConnection{
					Name:        pointer.MakePtr("cluster-connection"),
					Type:        pointer.MakePtr("Cluster"),
					ClusterName: pointer.MakePtr("my-cluster"),
					DbRoleToExecute: &admin.DBRoleToExecute{
						Role: pointer.MakePtr("read"),
						Type: pointer.MakePtr("BUILT_IN"),
					},
				}, nil
			},
		)

		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
