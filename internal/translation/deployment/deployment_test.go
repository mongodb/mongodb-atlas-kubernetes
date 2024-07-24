package deployment

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
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
