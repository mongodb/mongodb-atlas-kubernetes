package deployment

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type Service struct {
	admin.ClustersApi
	admin.ServerlessInstancesApi
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*Service, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewFromAPIs(client.ClustersApi, client.ServerlessInstancesApi), nil
}

func NewFromAPIs(clusterService admin.ClustersApi, serverlessAPI admin.ServerlessInstancesApi) *Service {
	return &Service{ClustersApi: clusterService, ServerlessInstancesApi: serverlessAPI}
}

func (ds *Service) ListClusterDeploymentNames(ctx context.Context, projectID string) ([]string, error) {
	var deploymentNames []string
	clusters, _, err := ds.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, err
	}
	if clusters.Results == nil {
		return deploymentNames, nil
	}

	for _, d := range *clusters.Results {
		name := pointer.GetOrDefault(d.Name, "")
		if name != "" {
			deploymentNames = append(deploymentNames, name)
		}
	}
	return deploymentNames, nil
}

func (ds *Service) ListDeploymentConns(ctx context.Context, projectID string) ([]Conn, error) {
	clusters, _, err := ds.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, err
	}
	clusterConns := clustersToConns(clusters.GetResults())

	serverless, _, err := ds.ListServerlessInstances(ctx, projectID).Execute()
	if err != nil {
		return nil, err
	}
	serverlessConns := serverlessToConns(serverless.GetResults())

	return connSet(clusterConns, serverlessConns), nil
}

func (ds *Service) Exists(ctx context.Context, projectID, clusterName string) (bool, error) {
	_, _, err := ds.GetCluster(ctx, projectID, clusterName).Execute()
	if admin.IsErrorCode(err, atlas.ClusterNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ds *Service) IsReady(ctx context.Context, projectID, deploymentName string) (bool, error) {
	clusterStatus, _, err := ds.GetClusterStatus(ctx, projectID, deploymentName).Execute()
	if err != nil {
		return false, err
	}
	return clusterStatus.GetChangeStatus() == string(mongodbatlas.ChangeStatusApplied), nil
}
