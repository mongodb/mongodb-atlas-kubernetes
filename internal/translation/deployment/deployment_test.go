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
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

func TestProductionAtlasDeployments_ListDeploymentConnections(t *testing.T) {
	t.Run("Shouldn't call the serverless api if running in Gov", func(t *testing.T) {
		mockClustersAPI := mockadmin.NewClustersApi(t)
		mockClustersAPI.EXPECT().ListClusters(context.Background(), mock.Anything).Return(
			admin.ListClustersApiRequest{ApiService: mockClustersAPI})
		mockClustersAPI.EXPECT().ListClustersExecute(admin.ListClustersApiRequest{ApiService: mockClustersAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)

		mockServerlessAPI := mockadmin.NewServerlessInstancesApi(t)
		mockServerlessAPI.EXPECT().ListServerlessInstancesExecute(mock.Anything).Unset()
		ds := &ProductionAtlasDeployments{
			clustersAPI:   mockClustersAPI,
			serverlessAPI: mockServerlessAPI,
			isGov:         true,
		}
		projectID := "testProjectID"
		_, err := ds.ListDeploymentConnections(context.Background(), projectID)
		assert.Nil(t, err)
	})

	t.Run("Should call the serverless api if not running in Gov", func(t *testing.T) {
		mockClustersAPI := mockadmin.NewClustersApi(t)
		mockClustersAPI.EXPECT().ListClusters(context.Background(), mock.Anything).Return(
			admin.ListClustersApiRequest{ApiService: mockClustersAPI})
		mockClustersAPI.EXPECT().ListClustersExecute(admin.ListClustersApiRequest{ApiService: mockClustersAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)

		mockServerlessAPI := mockadmin.NewServerlessInstancesApi(t)
		mockServerlessAPI.EXPECT().ListServerlessInstances(context.Background(), mock.Anything).Return(
			admin.ListServerlessInstancesApiRequest{ApiService: mockServerlessAPI})
		mockServerlessAPI.EXPECT().ListServerlessInstancesExecute(
			admin.ListServerlessInstancesApiRequest{ApiService: mockServerlessAPI}).Return(
			nil, &http.Response{StatusCode: http.StatusOK}, nil)
		ds := &ProductionAtlasDeployments{
			clustersAPI:   mockClustersAPI,
			serverlessAPI: mockServerlessAPI,
			isGov:         false,
		}
		projectID := "testProjectID"
		_, err := ds.ListDeploymentConnections(context.Background(), projectID)
		assert.Nil(t, err)
	})
}

func TestClusterExists(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		gov        bool
		result     bool
		err        error
	}{
		"should fail to assert a cluster exists in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster from atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to get cluster from atlas"),
		},
		"should fail to assert a serverless instance exists in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, errors.New("failed to get serverless instance from atlas"))

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to get serverless instance from atlas"),
		},
		"should return false when cluster doesn't exist": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "cluster0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ProviderUnsupported))

				return clusterAPI, serverlessInstanceAPI
			},
		},
		"should return false when serverless instance doesn't exist": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceNotFound))

				return clusterAPI, serverlessInstanceAPI
			},
		},
		"should return a cluster exists": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: true,
		},
		"should return a serverless instance exists": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(
						atlasServerlessInstance(),
						nil,
						nil,
					)

				return clusterAPI, serverlessInstanceAPI
			},
			result: true,
		},
		"should return false when asserting serverless instance exists in gov": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ProviderUnsupported))

				return clusterAPI, serverlessInstanceAPI
			},
			gov:    true,
			result: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, tt.gov)

			result, err := service.ClusterExists(context.Background(), tt.deployment.GetProjectID(), tt.deployment.GetName())
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		result     Deployment
		err        error
	}{
		"should fail to retrieve cluster from atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster from atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to get cluster from atlas"),
		},
		"should fail to retrieve serverless instance from atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, errors.New("failed to get serverless instance from atlas"))

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to get serverless instance from atlas"),
		},
		"should return nil when cluster doesn't exist": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ClusterNotFound))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "cluster0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ProviderUnsupported))

				return clusterAPI, serverlessInstanceAPI
			},
		},
		"should return nil when serverless instance doesn't exist": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceNotFound))

				return clusterAPI, serverlessInstanceAPI
			},
		},
		"should return a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedGeoShardedCluster(),
		},
		"should return a serverless instance": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetCluster(context.Background(), "project-id", "instance0").
					Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
					Return(nil, nil, atlasAPIError(atlas.ServerlessInstanceFromClusterAPI))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().GetServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
					Return(
						atlasServerlessInstance(),
						nil,
						nil,
					)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedServerlessInstance(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			result, err := service.GetDeployment(context.Background(), tt.deployment.GetProjectID(), tt.deployment.GetName())
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestCreateDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		result     Deployment
		err        error
	}{
		"should fail to create cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().CreateCluster(context.Background(), "project-id", mock.AnythingOfType("*admin.AdvancedClusterDescription")).
					Return(admin.CreateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().CreateClusterExecute(mock.AnythingOfType("admin.CreateClusterApiRequest")).
					Return(nil, nil, errors.New("failed to create cluster in atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to create cluster in atlas"),
		},
		"should fail to create serverless instance in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().CreateServerlessInstance(context.Background(), "project-id", mock.AnythingOfType("*admin.ServerlessInstanceDescriptionCreate")).
					Return(admin.CreateServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().CreateServerlessInstanceExecute(mock.AnythingOfType("admin.CreateServerlessInstanceApiRequest")).
					Return(nil, nil, errors.New("failed to create serverless instance in atlas"))

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to create serverless instance in atlas"),
		},
		"should create a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().CreateCluster(context.Background(), "project-id", mock.AnythingOfType("*admin.AdvancedClusterDescription")).
					Return(admin.CreateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().CreateClusterExecute(mock.AnythingOfType("admin.CreateClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedGeoShardedCluster(),
		},
		"should create a serverless instance": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().CreateServerlessInstance(context.Background(), "project-id", mock.AnythingOfType("*admin.ServerlessInstanceDescriptionCreate")).
					Return(admin.CreateServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().CreateServerlessInstanceExecute(mock.AnythingOfType("admin.CreateServerlessInstanceApiRequest")).
					Return(
						atlasServerlessInstance(),
						nil,
						nil,
					)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedServerlessInstance(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			result, err := service.CreateDeployment(context.Background(), tt.deployment)
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestUpdateDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		result     Deployment
		err        error
	}{
		"should fail to update cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateCluster(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.AdvancedClusterDescription")).
					Return(admin.UpdateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterExecute(mock.AnythingOfType("admin.UpdateClusterApiRequest")).
					Return(nil, nil, errors.New("failed to update cluster in atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to update cluster in atlas"),
		},
		"should fail to update serverless instance in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().UpdateServerlessInstance(context.Background(), "project-id", "instance0", mock.AnythingOfType("*admin.ServerlessInstanceDescriptionUpdate")).
					Return(admin.UpdateServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().UpdateServerlessInstanceExecute(mock.AnythingOfType("admin.UpdateServerlessInstanceApiRequest")).
					Return(nil, nil, errors.New("failed to update serverless instance in atlas"))

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to update serverless instance in atlas"),
		},
		"should update a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateCluster(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.AdvancedClusterDescription")).
					Return(admin.UpdateClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterExecute(mock.AnythingOfType("admin.UpdateClusterApiRequest")).
					Return(
						atlasGeoShardedCluster(),
						nil,
						nil,
					)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedGeoShardedCluster(),
		},
		"should update a serverless instance": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().UpdateServerlessInstance(context.Background(), "project-id", "instance0", mock.AnythingOfType("*admin.ServerlessInstanceDescriptionUpdate")).
					Return(admin.UpdateServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().UpdateServerlessInstanceExecute(mock.AnythingOfType("admin.UpdateServerlessInstanceApiRequest")).
					Return(
						atlasServerlessInstance(),
						nil,
						nil,
					)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedServerlessInstance(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			result, err := service.UpdateDeployment(context.Background(), tt.deployment)
			require.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestDeleteDeployment(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		result     Deployment
		err        error
	}{
		"should fail to delete cluster in atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().DeleteCluster(context.Background(), "project-id", "cluster0").
					Return(admin.DeleteClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().DeleteClusterExecute(mock.AnythingOfType("admin.DeleteClusterApiRequest")).
					Return(nil, errors.New("failed to delete cluster in atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to delete cluster in atlas"),
		},
		"should fail to delete serverless instance in atlas": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().DeleteServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.DeleteServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().DeleteServerlessInstanceExecute(mock.AnythingOfType("admin.DeleteServerlessInstanceApiRequest")).
					Return(nil, nil, errors.New("failed to delete serverless instance in atlas"))

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to delete serverless instance in atlas"),
		},
		"should delete a cluster": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().DeleteCluster(context.Background(), "project-id", "cluster0").
					Return(admin.DeleteClusterApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().DeleteClusterExecute(mock.AnythingOfType("admin.DeleteClusterApiRequest")).
					Return(nil, nil)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedGeoShardedCluster(),
		},
		"should delete a serverless instance": {
			deployment: serverlessInstance(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)
				serverlessInstanceAPI.EXPECT().DeleteServerlessInstance(context.Background(), "project-id", "instance0").
					Return(admin.DeleteServerlessInstanceApiRequest{ApiService: serverlessInstanceAPI})
				serverlessInstanceAPI.EXPECT().DeleteServerlessInstanceExecute(mock.AnythingOfType("admin.DeleteServerlessInstanceApiRequest")).
					Return(nil, nil, nil)

				return clusterAPI, serverlessInstanceAPI
			},
			result: expectedServerlessInstance(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			err := service.DeleteDeployment(context.Background(), tt.deployment)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestClusterWithProcessArgs(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
		result     Deployment
		err        error
	}{
		"should fail to retrieve cluster process args from atlas": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetClusterAdvancedConfiguration(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.GetClusterAdvancedConfigurationApiRequest")).
					Return(nil, nil, errors.New("failed to get cluster process args from atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			err: errors.New("failed to get cluster process args from atlas"),
		},
		"should return process args with default settings": {
			deployment: geoShardedCluster(),
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetClusterAdvancedConfiguration(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.GetClusterAdvancedConfigurationApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs{
							MinimumEnabledTlsProtocol: pointer.MakePtr("TLS1_2"),
							JavascriptEnabled:         pointer.MakePtr(true),
							NoTableScan:               pointer.MakePtr(false),
						},
						nil,
						nil,
					)

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
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
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().GetClusterAdvancedConfiguration(context.Background(), "project-id", "cluster0").
					Return(admin.GetClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().GetClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.GetClusterAdvancedConfigurationApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs{
							DefaultReadConcern:               pointer.MakePtr("available"),
							DefaultWriteConcern:              pointer.MakePtr("available"),
							FailIndexKeyTooLong:              pointer.MakePtr(true),
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

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: &Cluster{
				ProcessArgs: &akov2.ProcessArgs{
					DefaultReadConcern:               "available",
					DefaultWriteConcern:              "available",
					MinimumEnabledTLSProtocol:        "TLS1_1",
					FailIndexKeyTooLong:              pointer.MakePtr(true),
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
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			cluster := tt.deployment.(*Cluster)
			err := service.ClusterWithProcessArgs(context.Background(), cluster)
			require.Equal(t, tt.err, err)

			expectedCluster := tt.deployment.(*Cluster)
			assert.Equal(t, expectedCluster.ProcessArgs, cluster.ProcessArgs)
		})
	}
}

func TestUpdateProcessArgs(t *testing.T) {
	tests := map[string]struct {
		deployment Deployment
		apiMocker  func() (admin.ClustersApi, admin.ServerlessInstancesApi)
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
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
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
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateClusterAdvancedConfiguration(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescriptionProcessArgs")).
					Return(admin.UpdateClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.UpdateClusterAdvancedConfigurationApiRequest")).
					Return(nil, nil, errors.New("failed to update cluster process args in atlas"))

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
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
			apiMocker: func() (admin.ClustersApi, admin.ServerlessInstancesApi) {
				clusterAPI := mockadmin.NewClustersApi(t)
				clusterAPI.EXPECT().UpdateClusterAdvancedConfiguration(context.Background(), "project-id", "cluster0", mock.AnythingOfType("*admin.ClusterDescriptionProcessArgs")).
					Return(admin.UpdateClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
				clusterAPI.EXPECT().UpdateClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.UpdateClusterAdvancedConfigurationApiRequest")).
					Return(
						&admin.ClusterDescriptionProcessArgs{
							DefaultReadConcern:               pointer.MakePtr("available"),
							DefaultWriteConcern:              pointer.MakePtr("available"),
							FailIndexKeyTooLong:              pointer.MakePtr(true),
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

				serverlessInstanceAPI := mockadmin.NewServerlessInstancesApi(t)

				return clusterAPI, serverlessInstanceAPI
			},
			result: &Cluster{
				ProcessArgs: &akov2.ProcessArgs{
					DefaultReadConcern:               "available",
					DefaultWriteConcern:              "available",
					MinimumEnabledTLSProtocol:        "TLS1_2",
					FailIndexKeyTooLong:              pointer.MakePtr(true),
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
			clusterAPI, serverlessInstanceAPI := tt.apiMocker()
			service := NewProductionAtlasDeployments(clusterAPI, serverlessInstanceAPI, false)

			cluster := tt.deployment.(*Cluster)
			err := service.UpdateProcessArgs(context.Background(), cluster)
			require.Equal(t, tt.err, err)

			expectedCluster := tt.deployment.(*Cluster)
			assert.Equal(t, expectedCluster.ProcessArgs, cluster.ProcessArgs)
		})
	}
}

func atlasAPIError(code string) *admin.GenericOpenAPIError {
	err := admin.GenericOpenAPIError{}
	err.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr(code)})

	return &err
}

func geoShardedCluster() *Cluster {
	return &Cluster{
		ProjectID: "project-id",
		//nolint:dupl
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
	}
}

func expectedGeoShardedCluster() *Cluster {
	return &Cluster{
		ProjectID: "project-id",
		//nolint:dupl
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

func atlasGeoShardedCluster() *admin.AdvancedClusterDescription {
	return &admin.AdvancedClusterDescription{
		GroupId:                      pointer.MakePtr("project-id"),
		Name:                         pointer.MakePtr("cluster0"),
		ClusterType:                  pointer.MakePtr("GEOSHARDED"),
		DiskSizeGB:                   pointer.MakePtr(40.0),
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
		ReplicationSpecs: &[]admin.ReplicationSpec{
			//nolint:dupl
			{
				Id:        pointer.MakePtr("replication-id-2"),
				ZoneName:  pointer.MakePtr("Zone 2"),
				NumShards: pointer.MakePtr(1),
				RegionConfigs: &[]admin.CloudRegionConfig{
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("EU_CENTRAL_1"),
						Priority:     pointer.MakePtr(6),
						ElectableSpecs: &admin.HardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(2),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
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
						ElectableSpecs: &admin.HardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
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
			//nolint:dupl
			{
				Id:        pointer.MakePtr("replication-id-1"),
				ZoneName:  pointer.MakePtr("Zone 1"),
				NumShards: pointer.MakePtr(1),
				RegionConfigs: &[]admin.CloudRegionConfig{
					{
						ProviderName: pointer.MakePtr("AWS"),
						RegionName:   pointer.MakePtr("US_EAST_1"),
						Priority:     pointer.MakePtr(7),
						ElectableSpecs: &admin.HardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
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
						ElectableSpecs: &admin.HardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						ReadOnlySpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(3),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
						},
						AnalyticsSpecs: &admin.DedicatedHardwareSpec{
							InstanceSize:  pointer.MakePtr("M30"),
							NodeCount:     pointer.MakePtr(1),
							EbsVolumeType: pointer.MakePtr("STANDARD"),
							DiskIOPS:      pointer.MakePtr(3000),
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

func serverlessInstance() *Serverless {
	return &Serverless{
		ProjectID: "project-id",
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
	}
}

func expectedServerlessInstance() *Serverless {
	return &Serverless{
		ProjectID: "project-id",
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
					Key:   "A",
					Value: "A",
				},
				{
					Key:   "B",
					Value: "B",
				},
			},
		},
		State:          "IDLE",
		MongoDBVersion: "7.3.3",
		Connection: &status.ConnectionStrings{
			StandardSrv: "standard-str",
			PrivateEndpoint: []status.PrivateEndpoint{
				{
					SRVConnectionString: "connection-srv-str",
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
	}
}

func atlasServerlessInstance() *admin.ServerlessInstanceDescription {
	return &admin.ServerlessInstanceDescription{
		GroupId: pointer.MakePtr("project-id"),
		Name:    pointer.MakePtr("instance0"),
		ProviderSettings: admin.ServerlessProviderSettings{
			ProviderName:        pointer.MakePtr("SERVERLESS"),
			BackingProviderName: "AWS",
			RegionName:          "US_EAST_1",
		},
		ServerlessBackupOptions: &admin.ClusterServerlessBackupOptions{
			ServerlessContinuousBackupEnabled: pointer.MakePtr(true),
		},
		TerminationProtectionEnabled: pointer.MakePtr(true),
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
		StateName:      pointer.MakePtr("IDLE"),
		MongoDBVersion: pointer.MakePtr("7.3.3"),
		ConnectionStrings: &admin.ServerlessInstanceDescriptionConnectionStrings{
			StandardSrv: pointer.MakePtr("standard-str"),
			PrivateEndpoint: &[]admin.ServerlessConnectionStringsPrivateEndpointList{
				{
					SrvConnectionString: pointer.MakePtr("connection-srv-str"),
					Endpoints: &[]admin.ServerlessConnectionStringsPrivateEndpointItem{
						{
							ProviderName: pointer.MakePtr("AWS"),
							Region:       pointer.MakePtr("US_EAST_1"),
							EndpointId:   pointer.MakePtr("arn-endpoint-id"),
						},
					},
					Type: pointer.MakePtr("MONGOS"),
				},
			},
		},
	}
}
