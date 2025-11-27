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

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestEnsureServerlessPrivateEndpoints(t *testing.T) {
	t.Run("should fail when deployment is nil", func(t *testing.T) {
		result := ensureServerlessPrivateEndpoints(&workflow.Context{}, "project-id", nil)

		assert.Equal(
			t,
			workflow.Terminate(workflow.Internal, errors.New("serverless deployment spec is empty")),
			result,
		)
	})

	t.Run("should fail when serverless spec is nil", func(t *testing.T) {
		result := ensureServerlessPrivateEndpoints(&workflow.Context{}, "project-id", &akov2.AtlasDeployment{})

		assert.Equal(
			t,
			workflow.Terminate(workflow.Internal, errors.New("serverless deployment spec is empty")),
			result,
		)
	})

	t.Run("should fail when setting a GCP serverless instance with a private endpoint", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "GCP",
					},
					PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
						{
							Name: "spe-1",
						},
					},
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&workflow.Context{}, "project-id", &deployment)

		assert.Equal(
			t,
			workflow.Terminate(workflow.AtlasUnsupportedFeature, errors.New("serverless private endpoints are not supported for GCP")),
			result,
		)
	})

	t.Run("should succeed when setting a GCP serverless instance without private endpoints", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "GCP",
					},
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&workflow.Context{}, "project-id", &deployment)

		assert.Equal(t, workflow.OK(), result)
	})

	t.Run("should succeed when there are nothing to sync", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
					},
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.Anything).
			Return([]admin.ServerlessTenantEndpoint{}, &http.Response{}, nil)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&service, "project-id", &deployment)

		assert.Equal(t, workflow.OK(), result)
	})

	t.Run("should fail when error happens syncing private endpoints", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
					},
					PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
						{
							Name: "spe-1",
						},
					},
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.Anything).
			Return(nil, &http.Response{}, errors.New("connection failed"))
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&service, "project-id", &deployment)

		assert.Equal(t, result.GetError().Error(), "unable to retrieve list of serverless private endpoints from Atlas: connection failed")
		expected := workflow.Terminate(workflow.ServerlessPrivateEndpointFailed, fmt.Errorf("unable to retrieve list of serverless private endpoints from Atlas: connection failed"))
		assert.Equal(t, result.CloneWithoutError(), expected.CloneWithoutError())
	})

	t.Run("should succeed when syncing private endpoints still in progress", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
					},
					PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
						{
							Name: "spe-1",
						},
					},
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.Anything).
			Return([]admin.ServerlessTenantEndpoint{}, &http.Response{}, nil)
		speAPI.EXPECT().CreateServerlessPrivateEndpoint(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(admin.CreateServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().CreateServerlessPrivateEndpointExecute(mock.Anything).
			Return(
				&admin.ServerlessTenantEndpoint{
					Id:      pointer.MakePtr("spe-id"),
					Comment: pointer.MakePtr("spe-1"),
					Status:  pointer.MakePtr("RESERVATION_REQUESTED"),
				},
				&http.Response{},
				nil,
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&service, "project-id", &deployment)

		assert.Equal(
			t,
			workflow.InProgress(workflow.ServerlessPrivateEndpointInProgress, "Waiting serverless private endpoint to be configured"),
			result,
		)

		condition, ok := service.GetCondition(api.ServerlessPrivateEndpointReadyType)
		assert.True(t, ok)
		condition.LastTransitionTime = metav1.Time{} // not relevant for the unit test
		assert.Equal(
			t,
			api.Condition{
				Type:    api.ServerlessPrivateEndpointReadyType,
				Status:  corev1.ConditionFalse,
				Reason:  string(workflow.ServerlessPrivateEndpointInProgress),
				Message: "Waiting serverless private endpoint to be configured",
			},
			condition,
		)
	})

	t.Run("should succeed when finish syncing private endpoints", func(t *testing.T) {
		deployment := akov2.AtlasDeployment{
			Spec: akov2.AtlasDeploymentSpec{
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance-0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        "SERVERLESS",
						BackingProviderName: "AWS",
					},
					PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
						{
							Name:                    "spe-1",
							CloudProviderEndpointID: "aws-endpoint-id",
						},
					},
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.Anything).
			Return(
				[]admin.ServerlessTenantEndpoint{
					{
						Id:                      pointer.MakePtr("spe-id"),
						ProviderName:            pointer.MakePtr("AWS"),
						CloudProviderEndpointId: pointer.MakePtr("aws-endpoint-id"),
						Comment:                 pointer.MakePtr("spe-1"),
						Status:                  pointer.MakePtr(SPEStatusAvailable),
					},
				},
				&http.Response{},
				nil,
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}
		result := ensureServerlessPrivateEndpoints(&service, "project-id", &deployment)

		assert.Equal(t, workflow.OK(), result)

		condition, ok := service.GetCondition(api.ServerlessPrivateEndpointReadyType)
		assert.True(t, ok)
		condition.LastTransitionTime = metav1.Time{} // not relevant for the unit test
		assert.Equal(
			t,
			api.Condition{
				Type:   api.ServerlessPrivateEndpointReadyType,
				Status: corev1.ConditionTrue,
			},
			condition,
		)
	})
}

func TestSyncServerlessPrivateEndpoints(t *testing.T) {
	t.Run("should succeed adding, creating and deleting private endpoints", func(t *testing.T) {
		spec := akov2.ServerlessSpec{
			Name: "instance-0",
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				ProviderName:        "SERVERLESS",
				BackingProviderName: "AWS",
			},
			PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
				{
					Name: "spe-1",
				},
				{
					Name:                    "spe-2",
					CloudProviderEndpointID: "aws-endpoint-id",
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance-0").
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
			Return(
				[]admin.ServerlessTenantEndpoint{
					{
						Id:           pointer.MakePtr("spe-2-id"),
						ProviderName: pointer.MakePtr("AWS"),
						Comment:      pointer.MakePtr("spe-2"),
						Status:       pointer.MakePtr(SPEStatusReserved),
					},
					{
						Id:           pointer.MakePtr("spe-3-id"),
						ProviderName: pointer.MakePtr("AWS"),
						Comment:      pointer.MakePtr("spe-3"),
						Status:       pointer.MakePtr(SPEStatusAvailable),
					},
				},
				&http.Response{},
				nil,
			)
		speAPI.EXPECT().CreateServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", mock.AnythingOfType("*admin.ServerlessTenantCreateRequest")).
			Return(admin.CreateServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().CreateServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.CreateServerlessPrivateEndpointApiRequest")).
			Return(
				&admin.ServerlessTenantEndpoint{
					Id:           pointer.MakePtr("spe-1-id"),
					ProviderName: pointer.MakePtr("AWS"),
					Comment:      pointer.MakePtr("spe-1"),
					Status:       pointer.MakePtr("RESERVATION_REQUESTED"),
				},
				&http.Response{},
				nil,
			)
		speAPI.EXPECT().UpdateServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", "spe-2-id", mock.AnythingOfType("*admin.ServerlessTenantEndpointUpdate")).
			Return(admin.UpdateServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().UpdateServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.UpdateServerlessPrivateEndpointApiRequest")).
			Return(
				&admin.ServerlessTenantEndpoint{
					Id:                      pointer.MakePtr("spe-2-id"),
					ProviderName:            pointer.MakePtr("AWS"),
					CloudProviderEndpointId: pointer.MakePtr("aws-endpoint-id"),
					Comment:                 pointer.MakePtr("spe-2"),
					Status:                  pointer.MakePtr("INITIATING"),
				},
				&http.Response{},
				nil,
			)
		speAPI.EXPECT().DeleteServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", "spe-3-id").
			Return(admin.DeleteServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().DeleteServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.DeleteServerlessPrivateEndpointApiRequest")).
			Return(
				&http.Response{},
				nil,
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}

		finished, err := syncServerlessPrivateEndpoints(&service, "project-id", &spec)
		assert.NoError(t, err)
		assert.False(t, finished)
	})

	t.Run("should fail adding a private endpoint", func(t *testing.T) {
		spec := akov2.ServerlessSpec{
			Name: "instance-0",
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				ProviderName:        "SERVERLESS",
				BackingProviderName: "AWS",
			},
			PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
				{
					Name: "spe-1",
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance-0").
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
			Return([]admin.ServerlessTenantEndpoint{}, &http.Response{}, nil)
		speAPI.EXPECT().CreateServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", mock.AnythingOfType("*admin.ServerlessTenantCreateRequest")).
			Return(admin.CreateServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().CreateServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.CreateServerlessPrivateEndpointApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to create serverless private endpoint"),
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}

		finished, err := syncServerlessPrivateEndpoints(&service, "project-id", &spec)
		assert.ErrorContains(t, err, "failed to create serverless private endpoint")
		assert.False(t, finished)
	})

	t.Run("should fail updating a private endpoint", func(t *testing.T) {
		spec := akov2.ServerlessSpec{
			Name: "instance-0",
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				ProviderName:        "SERVERLESS",
				BackingProviderName: "AWS",
			},
			PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
				{
					Name:                    "spe-1",
					CloudProviderEndpointID: "endpoint-id",
				},
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance-0").
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
			Return(
				[]admin.ServerlessTenantEndpoint{
					{
						Id:           pointer.MakePtr("spe-1-id"),
						Comment:      pointer.MakePtr("spe-1"),
						Status:       pointer.MakePtr(SPEStatusReserved),
						ProviderName: pointer.MakePtr("AWS"),
					},
				},
				&http.Response{},
				nil,
			)
		speAPI.EXPECT().UpdateServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", "spe-1-id", mock.AnythingOfType("*admin.ServerlessTenantEndpointUpdate")).
			Return(admin.UpdateServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().UpdateServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.UpdateServerlessPrivateEndpointApiRequest")).
			Return(
				nil,
				&http.Response{},
				errors.New("failed to update serverless private endpoint"),
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}

		finished, err := syncServerlessPrivateEndpoints(&service, "project-id", &spec)
		assert.ErrorContains(t, err, "failed to update serverless private endpoint")
		assert.False(t, finished)
	})

	t.Run("should fail delete a private endpoint", func(t *testing.T) {
		spec := akov2.ServerlessSpec{
			Name: "instance-0",
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				ProviderName:        "SERVERLESS",
				BackingProviderName: "AWS",
			},
		}
		speAPI := mockadmin.NewServerlessPrivateEndpointsApi(t)
		speAPI.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance-0").
			Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speAPI})
		speAPI.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
			Return(
				[]admin.ServerlessTenantEndpoint{
					{
						Id:           pointer.MakePtr("spe-1-id"),
						Comment:      pointer.MakePtr("spe-1"),
						Status:       pointer.MakePtr(SPEStatusAvailable),
						ProviderName: pointer.MakePtr("AWS"),
					},
				},
				&http.Response{},
				nil,
			)
		speAPI.EXPECT().DeleteServerlessPrivateEndpoint(context.Background(), "project-id", "instance-0", "spe-1-id").
			Return(admin.DeleteServerlessPrivateEndpointApiRequest{ApiService: speAPI})
		speAPI.EXPECT().DeleteServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.DeleteServerlessPrivateEndpointApiRequest")).
			Return(
				&http.Response{},
				errors.New("failed to delete serverless private endpoint"),
			)
		service := workflow.Context{
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312006: &admin.APIClient{
					ServerlessPrivateEndpointsApi: speAPI,
				},
			},
		}

		finished, err := syncServerlessPrivateEndpoints(&service, "project-id", &spec)
		assert.ErrorContains(t, err, "failed to delete serverless private endpoint")
		assert.False(t, finished)
	})
}

func TestIsGCPWithPrivateEndpoints(t *testing.T) {
	t.Run("should return true when is GCP serverless instance containing private endpoint configuration", func(t *testing.T) {
		deployment := akov2.ServerlessSpec{
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				BackingProviderName: "GCP",
			},
			PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
				{
					Name: "spe-1",
				},
			},
		}

		assert.True(t, isGCPWithPrivateEndpoints(&deployment))
	})

	t.Run("should return false when is GCP serverless instance without private endpoint configuration", func(t *testing.T) {
		deployment := akov2.ServerlessSpec{
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				BackingProviderName: "GCP",
			},
		}

		assert.False(t, isGCPWithPrivateEndpoints(&deployment))
	})
}

func TestIsGCPWithoutPrivateEndpoints(t *testing.T) {
	t.Run("should return false when is GCP serverless instance containing private endpoint configuration", func(t *testing.T) {
		deployment := akov2.ServerlessSpec{
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				BackingProviderName: "GCP",
			},
			PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
				{
					Name: "spe-1",
				},
			},
		}

		assert.False(t, isGCPWithoutPrivateEndpoints(&deployment))
	})

	t.Run("should return true when is GCP serverless instance without private endpoint configuration", func(t *testing.T) {
		deployment := akov2.ServerlessSpec{
			ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
				BackingProviderName: "GCP",
			},
		}

		assert.True(t, isGCPWithoutPrivateEndpoints(&deployment))
	})
}

func TestSortTasks(t *testing.T) {
	t.Run("should sort one of each operation", func(t *testing.T) {
		spes := []akov2.ServerlessPrivateEndpoint{
			{
				Name: "spe-1",
			},
			{
				Name:                    "spe-2",
				CloudProviderEndpointID: "endpoint-id",
			},
		}
		atlas := []admin.ServerlessTenantEndpoint{
			{
				ProviderName: pointer.MakePtr("AWS"),
				Comment:      pointer.MakePtr("spe-2"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			{
				Comment: pointer.MakePtr("spe-3"),
			},
		}

		toCreate, toUpdate, toDelete := sortTasks(spes, atlas)

		assert.Equal(
			t,
			[]akov2.ServerlessPrivateEndpoint{
				{
					Name: "spe-1",
				},
			},
			toCreate,
		)
		assert.Equal(
			t,
			[]akov2.ServerlessPrivateEndpoint{
				{
					Name:                    "spe-2",
					CloudProviderEndpointID: "endpoint-id",
				},
			},
			toUpdate,
		)
		assert.Equal(
			t,
			[]admin.ServerlessTenantEndpoint{
				{
					Comment: pointer.MakePtr("spe-3"),
				},
			},
			toDelete,
		)
	})
}

func TestIsReadyToConnect(t *testing.T) {
	data := map[string]struct {
		spe      akov2.ServerlessPrivateEndpoint
		atlas    admin.ServerlessTenantEndpoint
		expected bool
	}{
		"should return false when private endpoint is not in RESERVED state": {
			spe: akov2.ServerlessPrivateEndpoint{},
			atlas: admin.ServerlessTenantEndpoint{
				Status: pointer.MakePtr("RESERVATION_REQUESTED"),
			},
			expected: false,
		},
		"should return false when a AWS private endpoint is in RESERVED state but miss endpoint ID": {
			spe: akov2.ServerlessPrivateEndpoint{},
			atlas: admin.ServerlessTenantEndpoint{
				ProviderName: pointer.MakePtr("AWS"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			expected: false,
		},
		"should return false when a Azure private endpoint is in RESERVED state but miss endpoint ID": {
			spe: akov2.ServerlessPrivateEndpoint{
				PrivateEndpointIPAddress: "some-ip-address",
			},
			atlas: admin.ServerlessTenantEndpoint{
				ProviderName: pointer.MakePtr("AZURE"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			expected: false,
		},
		"should return false when a Azure private endpoint is in RESERVED state but miss IP address": {
			spe: akov2.ServerlessPrivateEndpoint{
				CloudProviderEndpointID: "azure-endpoint-id",
			},
			atlas: admin.ServerlessTenantEndpoint{
				ProviderName: pointer.MakePtr("AZURE"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			expected: false,
		},
		"should return true when a Azure private endpoint is in RESERVED state and has connection data": {
			spe: akov2.ServerlessPrivateEndpoint{
				CloudProviderEndpointID:  "azure-endpoint-id",
				PrivateEndpointIPAddress: "some-ip-address",
			},
			atlas: admin.ServerlessTenantEndpoint{
				ProviderName: pointer.MakePtr("AZURE"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			expected: true,
		},
		"should return true when a AWS private endpoint is in RESERVED state and has connection data": {
			spe: akov2.ServerlessPrivateEndpoint{
				CloudProviderEndpointID: "aws-endpoint-id",
			},
			atlas: admin.ServerlessTenantEndpoint{
				ProviderName: pointer.MakePtr("AWS"),
				Status:       pointer.MakePtr(SPEStatusReserved),
			},
			expected: true,
		},
	}

	for desc, val := range data {
		t.Run(desc, func(t *testing.T) {
			spe := val.spe
			atlas := val.atlas
			assert.Equal(t, val.expected, isReadyToConnect(&spe, &atlas))
		})
	}
}

func TestCheckStatuses(t *testing.T) {
	data := map[string]struct {
		spes     []status.ServerlessPrivateEndpoint
		expected bool
	}{
		"should return true when nil": {
			spes:     nil,
			expected: true,
		},
		"should return true when empty": {
			spes:     []status.ServerlessPrivateEndpoint{},
			expected: true,
		},
		"should return true when all status are available": {
			spes: []status.ServerlessPrivateEndpoint{
				{
					Status: SPEStatusAvailable,
				},
				{
					Status: SPEStatusAvailable,
				},
			},
			expected: true,
		},
		"should return false when all status are not available": {
			spes: []status.ServerlessPrivateEndpoint{
				{
					Status: SPEStatusReserved,
				},
				{
					Status: SPEStatusDeleting,
				},
			},
			expected: false,
		},
		"should return false when at least one status is not available": {
			spes: []status.ServerlessPrivateEndpoint{
				{
					Status: SPEStatusReserved,
				},
				{
					Status: SPEStatusAvailable,
				},
				{
					Status: SPEStatusAvailable,
				},
			},
			expected: false,
		},
	}

	for desc, val := range data {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, val.expected, areSPEsAvailable(val.spes))
		})
	}
}
