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

//nolint:dupl
package deployment

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestProductionAtlasDeployments_ListDeploymentConnections(t *testing.T) {
	t.Run("Shouldn't call the serverless or flex api if running in Gov", func(t *testing.T) {
		mockClustersAPI := mockadmin.NewClustersApi(t)
		mockClustersAPI.EXPECT().ListClusters(context.Background(), mock.Anything).Return(
			admin.ListClustersApiRequest{ApiService: mockClustersAPI})
		mockClustersAPI.EXPECT().ListClustersExecute(admin.ListClustersApiRequest{ApiService: mockClustersAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)

		mockFlexAPI := mockadmin.NewFlexClustersApi(t)
		mockFlexAPI.EXPECT().ListFlexClustersExecute(mock.Anything).Unset()
		ds := &ProductionAtlasDeployments{
			clustersAPI: mockClustersAPI,
			flexAPI:     mockFlexAPI,
			isGov:       true,
		}
		projectID := "testProjectID"
		_, err := ds.ListDeploymentConnections(context.Background(), projectID)
		assert.Nil(t, err)
	})

	t.Run("Should call the serverless and flex apis if not running in Gov", func(t *testing.T) {
		mockClustersAPI := mockadmin.NewClustersApi(t)
		mockClustersAPI.EXPECT().ListClusters(context.Background(), mock.Anything).Return(
			admin.ListClustersApiRequest{ApiService: mockClustersAPI})
		mockClustersAPI.EXPECT().ListClustersExecute(admin.ListClustersApiRequest{ApiService: mockClustersAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)

		mockFlexAPI := mockadmin.NewFlexClustersApi(t)
		mockFlexAPI.EXPECT().ListFlexClusters(context.Background(), mock.Anything).Return(
			admin.ListFlexClustersApiRequest{ApiService: mockFlexAPI})
		mockFlexAPI.EXPECT().ListFlexClustersExecute(
			admin.ListFlexClustersApiRequest{ApiService: mockFlexAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)

		ds := &ProductionAtlasDeployments{
			clustersAPI: mockClustersAPI,
			flexAPI:     mockFlexAPI,
			isGov:       false,
		}
		projectID := "testProjectID"
		_, err := ds.ListDeploymentConnections(context.Background(), projectID)
		assert.Nil(t, err)
	})

	t.Run("Should create connection for each cluster type", func(t *testing.T) {
		mockClustersAPI := mockadmin.NewClustersApi(t)
		mockClustersAPI.EXPECT().ListClusters(context.Background(), mock.Anything).Return(
			admin.ListClustersApiRequest{ApiService: mockClustersAPI})
		mockClustersAPI.EXPECT().ListClustersExecute(admin.ListClustersApiRequest{ApiService: mockClustersAPI}).Return(
			&admin.PaginatedClusterDescription20240805{
				Results: &[]admin.ClusterDescription20240805{
					{
						Name:              pointer.MakePtr("testCluster"),
						ConnectionStrings: &admin.ClusterConnectionStrings{StandardSrv: pointer.MakePtr("clusterSRV")},
					},
				},
			}, &http.Response{StatusCode: http.StatusOK}, nil)

		mockFlexAPI := mockadmin.NewFlexClustersApi(t)
		mockFlexAPI.EXPECT().ListFlexClusters(context.Background(), mock.Anything).Return(
			admin.ListFlexClustersApiRequest{ApiService: mockFlexAPI})
		mockFlexAPI.EXPECT().ListFlexClustersExecute(
			admin.ListFlexClustersApiRequest{ApiService: mockFlexAPI}).Return(
			&admin.PaginatedFlexClusters20241113{
				Results: &[]admin.FlexClusterDescription20241113{
					{
						Name:              pointer.MakePtr("testFlex"),
						ConnectionStrings: &admin.FlexConnectionStrings20241113{StandardSrv: pointer.MakePtr("flexSRV")},
					},
				},
			}, &http.Response{StatusCode: http.StatusOK}, nil)

		ds := &ProductionAtlasDeployments{
			clustersAPI: mockClustersAPI,
			flexAPI:     mockFlexAPI,
			isGov:       false,
		}
		projectID := "testProjectID"
		conns, err := ds.ListDeploymentConnections(context.Background(), projectID)
		assert.Nil(t, err)
		assert.Equal(t, len(conns), 2)
	})
}

func TestClusterExists(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		gov        bool
		result     bool
		err        error
	}{
		"should fail to assert a cluster exists in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster from atlas"))

				clusterAPI := mockadmin.NewClustersApi(t)
				return clusterAPI, flexAPI
			},
			err: errors.New("failed to get cluster from atlas"),
		},
		"should fail to assert a serverless instance exists in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get serverless instance from atlas"))

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to get serverless instance from atlas"),
		},
		"should return false when cluster doesn't exist": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				err := &admin.GenericOpenAPIError{}
				err.SetModel(admin.ApiError{ErrorCode: atlas.NonFlexInFlexAPI})

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, nil, err)

				return clusterAPI, flexAPI
			},
		},
		"should return false when serverless instance doesn't exist": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				return clusterAPI, flexAPI
			},
		},
		"should return a cluster exists": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				err := &admin.GenericOpenAPIError{}
				err.SetModel(admin.ApiError{ErrorCode: atlas.NonFlexInFlexAPI})

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(mock.Anything, "project-id", "cluster0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, &http.Response{}, err)

				return clusterAPI, flexAPI
			},
			result: true,
		},
		"should return a serverless instance exists": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(&admin.FlexClusterDescription20241113{}, nil, nil)

				return clusterAPI, flexAPI
			},
			result: true,
		},
		"should return false when asserting serverless instance exists in gov": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			gov:    true,
			result: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, tt.gov)

			result, err := service.ClusterExists(context.Background(), "project-id", tt.deployment.GetDeploymentName())
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to retrieve cluster from atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster from atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to get cluster from atlas"),
		},
		"should fail to retrieve serverless instance from atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.AnythingOfType("admin.GetFlexClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get serverless instance from atlas"))

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to get serverless instance from atlas"),
		},
		"should return nil when cluster doesn't exist": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(mock.Anything, "project-id", mock.Anything).
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.Anything).Return(nil, nil, atlasAPIError(atlas.NonFlexInFlexAPI))

				return clusterAPI, flexAPI
			},
		},
		"should return nil when serverless instance doesn't exist": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().GetFlexCluster(mock.Anything, "project-id", mock.Anything).
					Return(admin.GetFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().GetFlexClusterExecute(mock.Anything).Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				return clusterAPI, flexAPI
			},
		},
		"should return a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: expectedGeoShardedCluster(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			result, err := service.GetDeployment(context.Background(), "project-id", tt.deployment)
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestCreateDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to create cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().CreateCluster(context.Background(), "project-id", mock.AnythingOfType("*admin.ClusterDescription20240805")).
					Return(admin.CreateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().CreateClusterExecute(mock.AnythingOfType("admin.CreateClusterApiRequest")).
					Return(nil, nil, errors.New("failed to create cluster in atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				return clusterAPI, flexAPI
			},
			err: errors.New("failed to create cluster in atlas"),
		},
		"should fail to create flex cluster in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().CreateFlexCluster(context.Background(), "project-id", mock.AnythingOfType("*admin.FlexClusterDescriptionCreate20241113")).
					Return(admin.CreateFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().CreateFlexClusterExecute(mock.AnythingOfType("admin.CreateFlexClusterApiRequest")).
					Return(nil, nil, errors.New("failed to create flex cluster in atlas"))

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to create flex cluster in atlas"),
		},
		"should create a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().CreateCluster(context.Background(), "project-id", mock.AnythingOfType("*admin.ClusterDescription20240805")).
					Return(admin.CreateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().CreateClusterExecute(mock.AnythingOfType("admin.CreateClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: expectedGeoShardedCluster(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			result, err := service.CreateDeployment(context.Background(), NewDeployment("project-id", tt.deployment))
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestUpdateDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to update cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateCluster(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescription20240805")).
					Return(admin.UpdateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterExecute(mock.AnythingOfType("admin.UpdateClusterApiRequest")).
					Return(nil, nil, errors.New("failed to update cluster in atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)
				return clusterAPI, flexAPI
			},
			err: errors.New("failed to update cluster in atlas"),
		},
		"should fail to update flex cluster in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().UpdateFlexCluster(context.Background(), "project-id", "instance0", mock.AnythingOfType("*admin.FlexClusterDescriptionUpdate20241113")).
					Return(admin.UpdateFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().UpdateFlexClusterExecute(mock.AnythingOfType("admin.UpdateFlexClusterApiRequest")).
					Return(nil, nil, errors.New("failed to update flex cluster in atlas"))

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to update flex cluster in atlas"),
		},
		"should update a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateCluster(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescription20240805")).
					Return(admin.UpdateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterExecute(mock.AnythingOfType("admin.UpdateClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				return clusterAPI, flexAPI
			},
			result: expectedGeoShardedCluster(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			result, err := service.UpdateDeployment(context.Background(), NewDeployment("project-id", tt.deployment))
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestDeleteDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to delete cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().DeleteCluster(context.Background(), "project-id", "cluster0").
					Return(admin.DeleteClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().DeleteClusterExecute(mock.AnythingOfType("admin.DeleteClusterApiRequest")).
					Return(nil, errors.New("failed to delete cluster in atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to delete cluster in atlas"),
		},
		"should fail to delete flex cluster in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().DeleteFlexCluster(context.Background(), "project-id", "instance0").
					Return(admin.DeleteFlexClusterApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().DeleteFlexClusterExecute(mock.AnythingOfType("admin.DeleteFlexClusterApiRequest")).
					Return(nil, errors.New("failed to delete flex cluster in atlas"))

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to delete flex cluster in atlas"),
		},
		"should delete a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().DeleteCluster(context.Background(), "project-id", "cluster0").
					Return(admin.DeleteClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().DeleteClusterExecute(mock.AnythingOfType("admin.DeleteClusterApiRequest")).
					Return(nil, nil)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: expectedGeoShardedCluster(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			err := service.DeleteDeployment(context.Background(), NewDeployment("project-id", tt.deployment))
			require.Equal(t, tt.err, err)
		})
	}
}

func TestClusterWithProcessArgs(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to retrieve cluster process args from atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetProcessArgs(context.Background(), "project-id", "cluster0").
					Return(admin.GetProcessArgsApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetProcessArgsExecute(mock.AnythingOfType("admin.GetProcessArgsApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster process args from atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to get cluster process args from atlas"),
		},
		"should return process args with default settings": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetProcessArgs(context.Background(), "project-id", "cluster0").
					Return(admin.GetProcessArgsApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetProcessArgsExecute(mock.AnythingOfType("admin.GetProcessArgsApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs20240805{
							MinimumEnabledTlsProtocol: pointer.MakePtr("TLS1_2"),
							JavascriptEnabled:         pointer.MakePtr(true),
							NoTableScan:               pointer.MakePtr(false),
						},
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: &Cluster{
				ProcessArgs: &akov2.ProcessArgs{
					MinimumEnabledTLSProtocol: "TLS1_2",
					JavascriptEnabled:         pointer.MakePtr(true),
					NoTableScan:               pointer.MakePtr(false),
				},
			},
		},
		"should return process args": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetProcessArgs(context.Background(), "project-id", "cluster0").
					Return(admin.GetProcessArgsApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetProcessArgsExecute(mock.AnythingOfType("admin.GetProcessArgsApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs20240805{
							DefaultWriteConcern:              pointer.MakePtr("available"),
							JavascriptEnabled:                pointer.MakePtr(false),
							MinimumEnabledTlsProtocol:        pointer.MakePtr("TLS1_1"),
							NoTableScan:                      pointer.MakePtr(true),
							OplogMinRetentionHours:           pointer.MakePtr(12.0),
							OplogSizeMB:                      pointer.MakePtr(5),
							SampleRefreshIntervalBIConnector: pointer.MakePtr(10),
							SampleSizeBIConnector:            pointer.MakePtr(5),
						},
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: &Cluster{
				ProcessArgs: &akov2.ProcessArgs{
					DefaultWriteConcern:              "available",
					MinimumEnabledTLSProtocol:        "TLS1_1",
					JavascriptEnabled:                pointer.MakePtr(true),
					NoTableScan:                      pointer.MakePtr(false),
					OplogSizeMB:                      pointer.MakePtr(int64(5)),
					SampleSizeBIConnector:            pointer.MakePtr(int64(5)),
					SampleRefreshIntervalBIConnector: pointer.MakePtr(int64(10)),
					OplogMinRetentionHours:           "12.0",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			d := NewDeployment("project-id", tt.deployment)
			cluster := d.(*Cluster)
			err := service.ClusterWithProcessArgs(context.Background(), cluster)
			require.Equal(t, tt.err, err)

			expectedCluster := d.(*Cluster)
			assert.Equal(t, expectedCluster.ProcessArgs, cluster.ProcessArgs)
		})
	}
}

func TestUpdateProcessArgs(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.FlexClustersApi)
		result     Deployment
		err        error
	}{
		"should fail to construct cluster process args": {
			deployment: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "cluster0",
				},
				ProcessArgs: &akov2.ProcessArgs{
					OplogMinRetentionHours: "wrong",
				},
			},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			err: &strconv.NumError{Func: "ParseFloat", Num: "wrong", Err: errors.New("invalid syntax")},
		},
		"should fail to retrieve cluster process args from atlas": {
			deployment: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "cluster0",
				},
				ProcessArgs: &akov2.ProcessArgs{
					DefaultReadConcern:        "available",
					DefaultWriteConcern:       "available",
					MinimumEnabledTLSProtocol: "TLS1_1",
					FailIndexKeyTooLong:       pointer.MakePtr(true),
					JavascriptEnabled:         pointer.MakePtr(true),
					NoTableScan:               pointer.MakePtr(false),
					OplogMinRetentionHours:    "12.0",
				},
			},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateProcessArgs(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescriptionProcessArgs20240805")).
					Return(admin.UpdateProcessArgsApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateProcessArgsExecute(mock.AnythingOfType("admin.UpdateProcessArgsApiRequest")).
					Return(nil, nil, errors.New("failed to update cluster process args in atlas"))

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			err: errors.New("failed to update cluster process args in atlas"),
		},
		"should update process args": {
			deployment: &Cluster{
				ProjectID: "project-id",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name: "cluster0",
				},
				ProcessArgs: &akov2.ProcessArgs{
					DefaultReadConcern:        "available",
					DefaultWriteConcern:       "available",
					MinimumEnabledTLSProtocol: "TLS1_2",
					FailIndexKeyTooLong:       pointer.MakePtr(true),
					JavascriptEnabled:         pointer.MakePtr(true),
					NoTableScan:               pointer.MakePtr(false),
					OplogMinRetentionHours:    "12.0",
				},
			},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateProcessArgs(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescriptionProcessArgs20240805")).
					Return(admin.UpdateProcessArgsApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateProcessArgsExecute(mock.AnythingOfType("admin.UpdateProcessArgsApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs20240805{
							DefaultWriteConcern:              pointer.MakePtr("available"),
							JavascriptEnabled:                pointer.MakePtr(true),
							MinimumEnabledTlsProtocol:        pointer.MakePtr("TLS1_2"),
							NoTableScan:                      pointer.MakePtr(false),
							OplogMinRetentionHours:           pointer.MakePtr(12.0),
							OplogSizeMB:                      pointer.MakePtr(5),
							SampleRefreshIntervalBIConnector: pointer.MakePtr(10),
							SampleSizeBIConnector:            pointer.MakePtr(5),
						},
						nil,
						nil,
					)

				flexAPI := mockadmin.NewFlexClustersApi(t)

				return clusterAPI, flexAPI
			},
			result: &Cluster{
				ProcessArgs: &akov2.ProcessArgs{
					DefaultWriteConcern:              "available",
					MinimumEnabledTLSProtocol:        "TLS1_2",
					JavascriptEnabled:                pointer.MakePtr(true),
					NoTableScan:                      pointer.MakePtr(false),
					OplogSizeMB:                      pointer.MakePtr(int64(5)),
					SampleSizeBIConnector:            pointer.MakePtr(int64(5)),
					SampleRefreshIntervalBIConnector: pointer.MakePtr(int64(10)),
					OplogMinRetentionHours:           "12.0",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			cluster := tt.deployment.(*Cluster)
			err := service.UpdateProcessArgs(context.Background(), cluster)
			require.Equal(t, tt.err, err)

			expectedCluster := tt.deployment.(*Cluster)
			assert.Equal(t, expectedCluster.ProcessArgs, cluster.ProcessArgs)
		})
	}
}

func TestUpgradeCluster(t *testing.T) {
	tests := map[string]struct {
		currentDeployment Deployment
		targetDeployment  Deployment
		apiMocker         func() (admin.ClustersApi, admin.FlexClustersApi)
		result            Deployment
		err               error
	}{
		"should fail to upgrade shared cluster in atlas": {
			currentDeployment: &Cluster{},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				flexAPI := mockadmin.NewFlexClustersApi(t)
				return clusterAPI, flexAPI
			},
			err: errors.New("upgrade from shared to dedicated is not supported"),
		},
		"should fail to upgrade flex instance in atlas": {
			currentDeployment: &Flex{},
			targetDeployment: &Cluster{
				ProjectID: "project-id",
				customResource: &akov2.AtlasDeployment{
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{
							Name: "cluster0",
						},
					},
				},
			},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().TenantUpgrade(context.Background(), "project-id", mock.AnythingOfType("*admin.AtlasTenantClusterUpgradeRequest20240805")).
					Return(admin.TenantUpgradeApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().TenantUpgradeExecute(mock.AnythingOfType("admin.TenantUpgradeApiRequest")).
					Return(nil, &http.Response{}, errors.New("failed to upgrade flex cluster in atlas"))
				return clusterAPI, flexAPI
			},
			err: errors.New("failed to upgrade flex cluster in atlas"),
		},
		"should upgrade flex instance in atlas": {
			currentDeployment: &Flex{},
			targetDeployment: &Cluster{
				ProjectID: "project-id",
				customResource: &akov2.AtlasDeployment{
					Spec: akov2.AtlasDeploymentSpec{
						DeploymentSpec: &akov2.AdvancedDeploymentSpec{
							Name: "cluster0",
						},
					},
				},
			},
			apiMocker: func() (admin.ClustersApi, admin.FlexClustersApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				flexAPI := mockadmin.NewFlexClustersApi(t)
				flexAPI.EXPECT().TenantUpgrade(context.Background(), "project-id", mock.AnythingOfType("*admin.AtlasTenantClusterUpgradeRequest20240805")).
					Return(admin.TenantUpgradeApiRequest{ApiService: flexAPI})
				flexAPI.EXPECT().TenantUpgradeExecute(mock.AnythingOfType("admin.TenantUpgradeApiRequest")).
					Return(
						&admin.FlexClusterDescription20241113{GroupId: pointer.MakePtr("project-id")},
						&http.Response{},
						nil,
					)
				return clusterAPI, flexAPI
			},
			result: &Flex{
				FlexSpec: &akov2.FlexSpec{
					Tags:             []*akov2.TagSpec{},
					ProviderSettings: &akov2.FlexProviderSettings{},
				},
				ProjectID:  "project-id",
				Connection: &status.ConnectionStrings{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, flexAPI := tt.apiMocker()
			service := NewAtlasDeployments(clusterAPI, nil, flexAPI, false)

			result, err := service.UpgradeToDedicated(context.Background(), tt.currentDeployment, tt.targetDeployment)
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func atlasAPIError(code string) *admin.GenericOpenAPIError {
	err := admin.GenericOpenAPIError{}
	err.SetModel(admin.ApiError{ErrorCode: code})

	return &err
}

func geoShardedCluster() *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: v1.ObjectMeta{
			Name: "cluster0",
		},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:                         "cluster0",
				ClusterType:                  "GEOSHARDED",
				DiskSizeGB:                   pointer.MakePtr(40),
				BackupEnabled:                pointer.MakePtr(true),
				PitEnabled:                   pointer.MakePtr(true),
				Paused:                       pointer.MakePtr(false),
				TerminationProtectionEnabled: true,
				EncryptionAtRestProvider:     "AWS",
				RootCertType:                 "ISRGROOTX1",
				MongoDBMajorVersion:          "7.0",
				VersionReleaseSystem:         "LTS",
				BiConnector: &akov2.BiConnectorSpec{
					Enabled:        pointer.MakePtr(true),
					ReadPreference: "secondary",
				},
				Labels: []common.LabelSpec{
					{
						Key:   "B",
						Value: "B",
					},
					{
						Key:   "A",
						Value: "A",
					},
				},
				Tags: []*akov2.TagSpec{
					{
						Key:   "B",
						Value: "B",
					},
					{
						Key:   "A",
						Value: "A",
					},
				},
				MongoDBVersion: "7.3.3",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						ZoneName:  "Zone 1",
						NumShards: 1,
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "EU_WEST_1",
								Priority:     pointer.MakePtr(5),
								ElectableSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								ReadOnlySpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AnalyticsSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(1),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									DiskGB: &akov2.DiskGB{
										Enabled: pointer.MakePtr(true),
									},
									Compute: &akov2.ComputeSpec{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  "M30",
										MaxInstanceSize:  "M60",
									},
								},
							},
							{
								ProviderName: "AWS",
								RegionName:   "US_EAST_1",
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								ReadOnlySpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AnalyticsSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(1),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									DiskGB: &akov2.DiskGB{
										Enabled: pointer.MakePtr(true),
									},
									Compute: &akov2.ComputeSpec{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  "M30",
										MaxInstanceSize:  "M60",
									},
								},
							},
						},
					},
					{
						ZoneName:  "Zone 2",
						NumShards: 1,
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "EU_CENTRAL_1",
								Priority:     pointer.MakePtr(6),
								ElectableSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(2),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								ReadOnlySpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AnalyticsSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(1),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									DiskGB: &akov2.DiskGB{
										Enabled: pointer.MakePtr(true),
									},
									Compute: &akov2.ComputeSpec{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  "M30",
										MaxInstanceSize:  "M60",
									},
								},
							},
							{
								ProviderName: "AWS",
								RegionName:   "EU_WEST_1",
								Priority:     pointer.MakePtr(4),
								ElectableSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								ReadOnlySpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(3),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AnalyticsSpecs: &akov2.Specs{
									InstanceSize:  "M30",
									NodeCount:     pointer.MakePtr(1),
									EbsVolumeType: "STANDARD",
									DiskIOPS:      pointer.MakePtr(int64(3000)),
								},
								AutoScaling: &akov2.AdvancedAutoScalingSpec{
									DiskGB: &akov2.DiskGB{
										Enabled: pointer.MakePtr(true),
									},
									Compute: &akov2.ComputeSpec{
										Enabled:          pointer.MakePtr(true),
										ScaleDownEnabled: pointer.MakePtr(true),
										MinInstanceSize:  "M30",
										MaxInstanceSize:  "M60",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func expectedGeoShardedCluster() *Cluster {
	return &Cluster{
		ProjectID: "project-id",
		AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:                         "cluster0",
			ClusterType:                  "GEOSHARDED",
			DiskSizeGB:                   pointer.MakePtr(40),
			BackupEnabled:                pointer.MakePtr(true),
			PitEnabled:                   pointer.MakePtr(true),
			Paused:                       pointer.MakePtr(false),
			TerminationProtectionEnabled: true,
			EncryptionAtRestProvider:     "AWS",
			RootCertType:                 "ISRGROOTX1",
			MongoDBMajorVersion:          "7.0",
			VersionReleaseSystem:         "LTS",
			BiConnector: &akov2.BiConnectorSpec{
				Enabled:        pointer.MakePtr(true),
				ReadPreference: "secondary",
			},
			Labels: []common.LabelSpec{

				{
					Key:   "A",
					Value: "A",
				},
				{
					Key:   "B",
					Value: "B",
				},
			},
			Tags: []*akov2.TagSpec{
				{
					Key:   "A",
					Value: "A",
				},
				{
					Key:   "B",
					Value: "B",
				},
			},
			MongoDBVersion: "7.3.3",
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					ZoneName:  "Zone 1",
					NumShards: 1,
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "US_EAST_1",
							Priority:     pointer.MakePtr(7),
							ElectableSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							ReadOnlySpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AnalyticsSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(1),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M30",
									MaxInstanceSize:  "M60",
								},
							},
						},
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST_1",
							Priority:     pointer.MakePtr(5),
							ElectableSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							ReadOnlySpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AnalyticsSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(1),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M30",
									MaxInstanceSize:  "M60",
								},
							},
						},
					},
				},
				{
					ZoneName:  "Zone 2",
					NumShards: 1,
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ProviderName: "AWS",
							RegionName:   "EU_CENTRAL_1",
							Priority:     pointer.MakePtr(6),
							ElectableSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(2),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							ReadOnlySpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AnalyticsSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(1),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M30",
									MaxInstanceSize:  "M60",
								},
							},
						},
						{
							ProviderName: "AWS",
							RegionName:   "EU_WEST_1",
							Priority:     pointer.MakePtr(4),
							ElectableSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							ReadOnlySpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(3),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AnalyticsSpecs: &akov2.Specs{
								InstanceSize:  "M30",
								NodeCount:     pointer.MakePtr(1),
								EbsVolumeType: "STANDARD",
								DiskIOPS:      pointer.MakePtr(int64(3000)),
							},
							AutoScaling: &akov2.AdvancedAutoScalingSpec{
								DiskGB: &akov2.DiskGB{
									Enabled: pointer.MakePtr(true),
								},
								Compute: &akov2.ComputeSpec{
									Enabled:          pointer.MakePtr(true),
									ScaleDownEnabled: pointer.MakePtr(true),
									MinInstanceSize:  "M30",
									MaxInstanceSize:  "M60",
								},
							},
						},
					},
				},
			},
		},
		State:          "CREATING",
		MongoDBVersion: "7.3.3",
		Connection: &status.ConnectionStrings{
			Standard:    "standard-str",
			StandardSrv: "standard-srv-str",
			Private:     "private-str",
			PrivateSrv:  "private-srv-str",
			PrivateEndpoint: []status.PrivateEndpoint{
				{
					ConnectionString:                  "connection-str",
					SRVConnectionString:               "connection-srv-str",
					SRVShardOptimizedConnectionString: "connection-sharded-srv-str",
					Endpoints: []status.Endpoint{
						{
							ProviderName: "AWS",
							Region:       "US_EAST_1",
							EndpointID:   "arn-endpoint-id",
						},
					},
				},
			},
		},
		ReplicaSet: []status.ReplicaSet{
			{
				ID:       "replication-id-2",
				ZoneName: "Zone 2",
			},
			{
				ID:       "replication-id-1",
				ZoneName: "Zone 1",
			},
		},
		computeAutoscalingEnabled: true,
		instanceSizeOverride:      "M30",
	}
}

func atlasGeoShardedCluster() *admin.ClusterDescription20240805 {
	return &admin.ClusterDescription20240805{
		GroupId:                      pointer.MakePtr("project-id"),
		Name:                         pointer.MakePtr("cluster0"),
		ClusterType:                  pointer.MakePtr("GEOSHARDED"),
		BackupEnabled:                pointer.MakePtr(true),
		PitEnabled:                   pointer.MakePtr(true),
		Paused:                       pointer.MakePtr(false),
		TerminationProtectionEnabled: pointer.MakePtr(true),
		EncryptionAtRestProvider:     pointer.MakePtr("AWS"),
		RootCertType:                 pointer.MakePtr("ISRGROOTX1"),
		MongoDBMajorVersion:          pointer.MakePtr("7.0"),
		VersionReleaseSystem:         pointer.MakePtr("LTS"),
		BiConnector: &admin.BiConnector{
			Enabled:        pointer.MakePtr(true),
			ReadPreference: pointer.MakePtr("secondary"),
		},
		Labels: &[]admin.ComponentLabel{
			{
				Key:   pointer.MakePtr("B"),
				Value: pointer.MakePtr("B"),
			},
			{
				Key:   pointer.MakePtr("A"),
				Value: pointer.MakePtr("A"),
			},
		},
		Tags: &[]admin.ResourceTag{
			{
				Key:   "B",
				Value: "B",
			},
			{
				Key:   "A",
				Value: "A",
			},
		},
		ReplicationSpecs: &[]admin.ReplicationSpec20240805{
			{
				Id:       pointer.MakePtr("replication-id-2"),
				ZoneName: pointer.MakePtr("Zone 2"),
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("EU_CENTRAL_1"),
						Priority:     pointer.MakePtr(6),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(2),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							DiskGB: &admin.DiskGBAutoScaling{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  pointer.MakePtr("M30"),
								MaxInstanceSize:  pointer.MakePtr("M60"),
							},
						},
					},
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("EU_WEST_1"),
						Priority:     pointer.MakePtr(4),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							DiskGB: &admin.DiskGBAutoScaling{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  pointer.MakePtr("M30"),
								MaxInstanceSize:  pointer.MakePtr("M60"),
							},
						},
					},
				},
			},
			{
				Id:       pointer.MakePtr("replication-id-1"),
				ZoneName: pointer.MakePtr("Zone 1"),
				RegionConfigs: &[]admin.CloudRegionConfig20240805{
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("US_EAST_1"),
						Priority:     pointer.MakePtr(7),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							DiskGB: &admin.DiskGBAutoScaling{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  pointer.MakePtr("M30"),
								MaxInstanceSize:  pointer.MakePtr("M60"),
							},
						},
					},
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("EU_WEST_1"),
						Priority:     pointer.MakePtr(5),
						ElectableSpecs: &admin.HardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec20240805{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
							DiskSizeGB:    pointer.MakePtr(40.0),
						},
						AutoScaling: &admin.AdvancedAutoScalingSettings{
							DiskGB: &admin.DiskGBAutoScaling{
								Enabled: pointer.MakePtr(true),
							},
							Compute: &admin.AdvancedComputeAutoScaling{
								Enabled:          pointer.MakePtr(true),
								ScaleDownEnabled: pointer.MakePtr(true),
								MinInstanceSize:  pointer.MakePtr("M30"),
								MaxInstanceSize:  pointer.MakePtr("M60"),
							},
						},
					},
				},
			},
		},
		StateName:      pointer.MakePtr("CREATING"),
		MongoDBVersion: pointer.MakePtr("7.3.3"),
		ConnectionStrings: &admin.ClusterConnectionStrings{
			PrivateEndpoint: &[]admin.ClusterDescriptionConnectionStringsPrivateEndpoint{
				{
					ConnectionString:                  pointer.MakePtr("connection-str"),
					SrvConnectionString:               pointer.MakePtr("connection-srv-str"),
					SrvShardOptimizedConnectionString: pointer.MakePtr("connection-sharded-srv-str"),
					Endpoints: &[]admin.ClusterDescriptionConnectionStringsPrivateEndpointEndpoint{
						{
							ProviderName: pointer.MakePtr("AWS"),
							Region:       pointer.MakePtr("US_EAST_1"),
							EndpointId:   pointer.MakePtr("arn-endpoint-id"),
						},
					},
					Type: pointer.MakePtr("MONGOS"),
				},
			},
			Private:     pointer.MakePtr("private-str"),
			PrivateSrv:  pointer.MakePtr("private-srv-str"),
			Standard:    pointer.MakePtr("standard-str"),
			StandardSrv: pointer.MakePtr("standard-srv-str"),
		},
	}
}

func serverlessInstance() *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: v1.ObjectMeta{
			Name: "instance0",
		},
		Spec: akov2.AtlasDeploymentSpec{
			ServerlessSpec: &akov2.ServerlessSpec{
				Name: "instance0",
				ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
					ProviderName:        "SERVERLESS",
					BackingProviderName: "AWS",
					RegionName:          "US_EAST_1",
				},
				BackupOptions: akov2.ServerlessBackupOptions{
					ServerlessContinuousBackupEnabled: true,
				},
				TerminationProtectionEnabled: true,
				Tags: []*akov2.TagSpec{
					{
						Key:   "B",
						Value: "B",
					},
					{
						Key:   "A",
						Value: "A",
					},
				},
			},
		},
	}
}
