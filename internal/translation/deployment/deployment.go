package deployment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type AtlasDeploymentsService interface {
	ListClusterNames(ctx context.Context, projectID string) ([]string, error)
	ListDeploymentConnections(ctx context.Context, projectID string) ([]Connection, error)
	ClusterExists(ctx context.Context, projectID, clusterName string) (bool, error)
	DeploymentIsReady(ctx context.Context, projectID, deploymentName string) (bool, error)
}

type ProductionAtlasDeployments struct {
	clustersAPI   admin.ClustersApi
	serverlessAPI admin.ServerlessInstancesApi
	isGov         bool
}

func NewAtlasDeploymentsService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger, isGov bool) (*ProductionAtlasDeployments, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned client: %w", err)
	}
	return NewProductionAtlasDeployments(client.ClustersApi, client.ServerlessInstancesApi, isGov), nil
}

func NewProductionAtlasDeployments(clusterService admin.ClustersApi, serverlessAPI admin.ServerlessInstancesApi, isGov bool) *ProductionAtlasDeployments {
	return &ProductionAtlasDeployments{clustersAPI: clusterService, serverlessAPI: serverlessAPI, isGov: isGov}
}

func (ds *ProductionAtlasDeployments) ListClusterNames(ctx context.Context, projectID string) ([]string, error) {
	var deploymentNames []string
	clusters, _, err := ds.clustersAPI.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list cluster names for project %s: %w", projectID, err)
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

func (ds *ProductionAtlasDeployments) ListDeploymentConnections(ctx context.Context, projectID string) ([]Connection, error) {
	clusters, _, err := ds.clustersAPI.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters for project %s: %w", projectID, err)
	}
	clusterConns := clustersToConnections(clusters.GetResults())

	if ds.isGov {
		return clusterConns, nil
	}

	serverless, _, serverlessErr := ds.serverlessAPI.ListServerlessInstances(ctx, projectID).Execute()
	if serverlessErr != nil {
		return nil, fmt.Errorf("failed to list serverless deployments for project %s: %w", projectID, err)
	}
	serverlessConns := serverlessToConnections(serverless.GetResults())

	return connectionSet(clusterConns, serverlessConns), nil
}

func (ds *ProductionAtlasDeployments) ClusterExists(ctx context.Context, projectID, clusterName string) (bool, error) {
	_, _, err := ds.clustersAPI.GetCluster(ctx, projectID, clusterName).Execute()
	if admin.IsErrorCode(err, atlas.ClusterNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get cluster %q: %w", clusterName, err)
	}
	return true, nil
}

func (ds *ProductionAtlasDeployments) DeploymentIsReady(ctx context.Context, projectID, deploymentName string) (bool, error) {
	// although this is within the clusters API it seems to also reply for serverless deployments
	clusterStatus, _, err := ds.clustersAPI.GetClusterStatus(ctx, projectID, deploymentName).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to get cluster %q status %w", deploymentName, err)
	}
	return clusterStatus.GetChangeStatus() == string(mongodbatlas.ChangeStatusApplied), nil
}
