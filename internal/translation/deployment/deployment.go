package deployment

import (
	"context"
	"errors"
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

	GetDeployment(ctx context.Context, projectID, name string) (Deployment, error)
	CreateDeployment(ctx context.Context, deployment Deployment) (Deployment, error)
	UpdateDeployment(ctx context.Context, deployment Deployment) (Deployment, error)
	DeleteDeployment(ctx context.Context, deployment Deployment) error
	ClusterWithProcessArgs(ctx context.Context, cluster *Cluster) error
	UpdateProcessArgs(ctx context.Context, cluster *Cluster) error
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
	d, err := ds.GetDeployment(ctx, projectID, clusterName)
	if err != nil {
		return false, err
	}

	return d != nil, nil
}

func (ds *ProductionAtlasDeployments) DeploymentIsReady(ctx context.Context, projectID, deploymentName string) (bool, error) {
	// although this is within the clusters API it seems to also reply for serverless deployments
	clusterStatus, _, err := ds.clustersAPI.GetClusterStatus(ctx, projectID, deploymentName).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to get cluster %q status %w", deploymentName, err)
	}
	return clusterStatus.GetChangeStatus() == string(mongodbatlas.ChangeStatusApplied), nil
}

func (ds *ProductionAtlasDeployments) GetDeployment(ctx context.Context, projectID, name string) (Deployment, error) {
	cluster, _, err := ds.clustersAPI.GetCluster(ctx, projectID, name).Execute()
	if err == nil {
		return clusterFromAtlas(cluster), nil
	}

	if !admin.IsErrorCode(err, atlas.ClusterNotFound) && !admin.IsErrorCode(err, atlas.ServerlessInstanceFromClusterAPI) {
		return nil, err
	}

	serverless, _, err := ds.serverlessAPI.GetServerlessInstance(ctx, projectID, name).Execute()
	if err == nil {
		return serverlessFromAtlas(serverless), err
	}

	if !admin.IsErrorCode(err, atlas.ServerlessInstanceNotFound) && !admin.IsErrorCode(err, atlas.ProviderUnsupported) {
		return nil, err
	}

	return nil, nil
}

func (ds *ProductionAtlasDeployments) CreateDeployment(ctx context.Context, deployment Deployment) (Deployment, error) {
	switch d := deployment.(type) {
	case *Cluster:
		cluster, _, err := ds.clustersAPI.CreateCluster(ctx, deployment.GetProjectID(), clusterCreateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return clusterFromAtlas(cluster), nil
	case *Serverless:
		serverless, _, err := ds.serverlessAPI.CreateServerlessInstance(ctx, deployment.GetProjectID(), serverlessCreateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return serverlessFromAtlas(serverless), nil
	}

	return nil, errors.New("unable to create deployment: unknown type")
}

func (ds *ProductionAtlasDeployments) UpdateDeployment(ctx context.Context, deployment Deployment) (Deployment, error) {
	switch d := deployment.(type) {
	case *Cluster:
		cluster, _, err := ds.clustersAPI.UpdateCluster(ctx, deployment.GetProjectID(), deployment.GetName(), clusterUpdateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return clusterFromAtlas(cluster), nil
	case *Serverless:
		serverless, _, err := ds.serverlessAPI.UpdateServerlessInstance(ctx, deployment.GetProjectID(), deployment.GetName(), serverlessUpdateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return serverlessFromAtlas(serverless), nil
	}

	return nil, errors.New("unable to create deployment: unknown type")
}

func (ds *ProductionAtlasDeployments) DeleteDeployment(ctx context.Context, deployment Deployment) error {
	var err error

	switch deployment.(type) {
	case *Cluster:
		_, err = ds.clustersAPI.DeleteCluster(ctx, deployment.GetProjectID(), deployment.GetName()).Execute()
	case *Serverless:
		_, _, err = ds.serverlessAPI.DeleteServerlessInstance(ctx, deployment.GetProjectID(), deployment.GetName()).Execute()
	}

	if err != nil {
		if admin.IsErrorCode(err, atlas.ClusterNotFound) {
			return nil
		}

		return err
	}

	return nil
}

func (ds *ProductionAtlasDeployments) ClusterWithProcessArgs(ctx context.Context, cluster *Cluster) error {
	config, _, err := ds.clustersAPI.GetClusterAdvancedConfiguration(ctx, cluster.GetProjectID(), cluster.GetName()).Execute()
	if err != nil {
		return err
	}

	cluster.ProcessArgs = processArgsFromAtlas(config)

	return nil
}

func (ds *ProductionAtlasDeployments) UpdateProcessArgs(ctx context.Context, cluster *Cluster) error {
	processArgs, err := processArgsToAtlas(cluster.ProcessArgs)
	if err != nil {
		return err
	}

	config, _, err := ds.clustersAPI.UpdateClusterAdvancedConfiguration(ctx, cluster.GetProjectID(), cluster.GetName(), processArgs).Execute()
	if err != nil {
		return err
	}

	cluster.ProcessArgs = processArgsFromAtlas(config)

	return nil
}
