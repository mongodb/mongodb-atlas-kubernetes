package deployment

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
)

type AtlasDeploymentsService interface {
	DeploymentService
	GlobalClusterService
}

type DeploymentService interface {
	ListDeploymentNames(ctx context.Context, projectID string) ([]string, error)
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

type GlobalClusterService interface {
	GetCustomZones(ctx context.Context, projectID, clusterName string) (map[string]string, error)
	CreateCustomZones(ctx context.Context, projectID, clusterName string, mappings []akov2.CustomZoneMapping) (map[string]string, error)
	DeleteCustomZones(ctx context.Context, projectID, clusterName string) error
	GetZoneMapping(ctx context.Context, projectID, deploymentName string) (map[string]string, error)
	GetManagedNamespaces(ctx context.Context, projectID, clusterName string) ([]akov2.ManagedNamespace, error)
	CreateManagedNamespace(ctx context.Context, projectID, clusterName string, ns *akov2.ManagedNamespace) error
	DeleteManagedNamespace(ctx context.Context, projectID, clusterName string, ns *akov2.ManagedNamespace) error
}

type ProductionAtlasDeployments struct {
	clustersAPI      admin.ClustersApi
	serverlessAPI    admin.ServerlessInstancesApi
	globalClusterAPI admin.GlobalClustersApi
	isGov            bool
}

func NewAtlasDeploymentsService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger, isGov bool) (*ProductionAtlasDeployments, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create versioned client: %w", err)
	}
	return NewAtlasDeployments(client.ClustersApi, client.ServerlessInstancesApi, client.GlobalClustersApi, isGov), nil
}

func NewAtlasDeployments(clusterService admin.ClustersApi, serverlessAPI admin.ServerlessInstancesApi, globalClusterAPI admin.GlobalClustersApi, isGov bool) *ProductionAtlasDeployments {
	return &ProductionAtlasDeployments{clustersAPI: clusterService, serverlessAPI: serverlessAPI, globalClusterAPI: globalClusterAPI, isGov: isGov}
}

func (ds *ProductionAtlasDeployments) ListDeploymentNames(ctx context.Context, projectID string) ([]string, error) {
	var deploymentNames []string
	clusters, _, err := ds.clustersAPI.ListClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters for project %s: %w", projectID, err)
	}
	for _, d := range clusters.GetResults() {
		name := pointer.GetOrDefault(d.Name, "")
		if name != "" {
			deploymentNames = append(deploymentNames, name)
		}
	}

	if ds.isGov {
		return deploymentNames, nil
	}

	serverless, _, err := ds.serverlessAPI.ListServerlessInstances(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list serverless deployments for project %s: %w", projectID, err)
	}
	for _, d := range serverless.GetResults() {
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

	if ds.isGov {
		return nil, nil
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

func (ds *ProductionAtlasDeployments) GetCustomZones(ctx context.Context, projectID, clusterName string) (map[string]string, error) {
	geosharding, _, err := ds.globalClusterAPI.GetManagedNamespace(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get global cluster: %w", err)
	}
	return geosharding.GetCustomZoneMapping(), nil
}

func (ds *ProductionAtlasDeployments) CreateCustomZones(ctx context.Context, projectID, clusterName string, mappings []akov2.CustomZoneMapping) (map[string]string, error) {
	geosharding, _, err := ds.globalClusterAPI.CreateCustomZoneMapping(ctx, projectID, clusterName, customZonesToAtlas(&mappings)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create custom zone: %w", err)
	}
	return geosharding.GetCustomZoneMapping(), nil
}

func (ds *ProductionAtlasDeployments) DeleteCustomZones(ctx context.Context, projectID, clusterName string) error {
	_, _, err := ds.globalClusterAPI.DeleteAllCustomZoneMappings(ctx, projectID, clusterName).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete custom zone: %w", err)
	}
	return nil
}

func (ds *ProductionAtlasDeployments) GetZoneMapping(ctx context.Context, projectID, deploymentName string) (map[string]string, error) {
	cluster, _, err := ds.clustersAPI.GetCluster(ctx, projectID, deploymentName).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}
	result := make(map[string]string, len(cluster.GetReplicationSpecs()))
	for _, rc := range cluster.GetReplicationSpecs() {
		result[rc.GetId()] = rc.GetZoneName()
	}
	return result, nil
}

func (ds *ProductionAtlasDeployments) GetManagedNamespaces(ctx context.Context, projectID, clusterName string) ([]akov2.ManagedNamespace, error) {
	geosharding, _, err := ds.globalClusterAPI.GetManagedNamespace(ctx, projectID, clusterName).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get global cluster: %w", err)
	}
	return managedNamespacesFromAtlas(geosharding), nil
}

func (ds *ProductionAtlasDeployments) CreateManagedNamespace(ctx context.Context, projectID, clusterName string, ns *akov2.ManagedNamespace) error {
	_, _, err := ds.globalClusterAPI.CreateManagedNamespace(ctx, projectID, clusterName, managedNamespaceToAtlas(ns)).Execute()
	if err != nil {
		return fmt.Errorf("failed to create managed namespace: %w", err)
	}
	return nil
}

func (ds *ProductionAtlasDeployments) DeleteManagedNamespace(ctx context.Context, projectID, clusterName string, namespace *akov2.ManagedNamespace) error {
	params := &admin.DeleteManagedNamespaceApiParams{
		GroupId:     projectID,
		ClusterName: clusterName,
		Db:          &namespace.Db,
		Collection:  &namespace.Collection,
	}
	_, _, err := ds.globalClusterAPI.DeleteManagedNamespaceWithParams(ctx, params).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete managed namespace: %w", err)
	}
	return nil
}
