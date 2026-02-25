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

package deployment

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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

	GetDeployment(ctx context.Context, projectID string, deployment *akov2.AtlasDeployment) (Deployment, error)
	CreateDeployment(ctx context.Context, deployment Deployment) (Deployment, error)
	UpdateDeployment(ctx context.Context, deployment Deployment) (Deployment, error)
	DeleteDeployment(ctx context.Context, deployment Deployment) error
	UpgradeToDedicated(ctx context.Context, currentDeployment, targetDeployment Deployment) (Deployment, error)
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
	flexAPI          admin.FlexClustersApi
	globalClusterAPI admin.GlobalClustersApi
	isGov            bool
}

func NewAtlasDeployments(clusterService admin.ClustersApi, globalClusterAPI admin.GlobalClustersApi, flexAPI admin.FlexClustersApi, isGov bool) *ProductionAtlasDeployments {
	return &ProductionAtlasDeployments{clustersAPI: clusterService, globalClusterAPI: globalClusterAPI, flexAPI: flexAPI, isGov: isGov}
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

	flex, _, err := ds.flexAPI.ListFlexClusters(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list flex clusters for project %s: %w", projectID, err)
	}
	for _, d := range flex.GetResults() {
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

	flex, _, flexErr := ds.flexAPI.ListFlexClusters(ctx, projectID).Execute()
	if flexErr != nil {
		return nil, fmt.Errorf("failed to list flex clusters for project %s: %w", projectID, err)
	}
	flexConns := flexToConnections(flex.GetResults())

	return connectionSet(clusterConns, flexConns), nil
}

func (ds *ProductionAtlasDeployments) ClusterExists(ctx context.Context, projectID, name string) (bool, error) {
	flex, err := ds.GetFlexCluster(ctx, projectID, name)
	if !admin.IsErrorCode(err, atlas.NonFlexInFlexAPI) && err != nil {
		return false, err
	}
	if flex != nil {
		return true, nil
	}

	cluster, err := ds.GetCluster(ctx, projectID, name)
	if !admin.IsErrorCode(err, atlas.ServerlessInstanceFromClusterAPI) && err != nil {
		return false, err
	}
	if cluster != nil {
		return true, nil
	}

	return false, nil
}

func (ds *ProductionAtlasDeployments) DeploymentIsReady(ctx context.Context, projectID, deploymentName string) (bool, error) {
	// although this is within the clusters API it seems to also reply for serverless deployments
	clusterStatus, _, err := ds.clustersAPI.GetClusterStatus(ctx, projectID, deploymentName).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to get cluster %q status %w", deploymentName, err)
	}
	return clusterStatus.GetChangeStatus() == "APPLIED", nil
}

func (ds *ProductionAtlasDeployments) GetFlexCluster(ctx context.Context, projectID, name string) (*Flex, error) {
	if ds.isGov {
		return nil, nil
	}

	flex, _, err := ds.flexAPI.GetFlexCluster(ctx, projectID, name).Execute()
	if err == nil {
		return flexFromAtlas(flex), nil
	}

	if !admin.IsErrorCode(err, atlas.ClusterNotFound) {
		return nil, err
	}

	return nil, nil
}

func (ds *ProductionAtlasDeployments) GetCluster(ctx context.Context, projectID, name string) (*Cluster, error) {
	cluster, _, err := ds.clustersAPI.GetCluster(ctx, projectID, name).Execute()
	if err == nil {
		return clusterFromAtlas(cluster), nil
	}

	if !admin.IsErrorCode(err, atlas.ClusterNotFound) {
		return nil, err
	}

	return nil, nil
}

func (ds *ProductionAtlasDeployments) GetDeployment(ctx context.Context, projectID string, deployment *akov2.AtlasDeployment) (Deployment, error) {
	if deployment == nil {
		return nil, errors.New("deployment is nil")
	}

	cluster, err := ds.GetCluster(ctx, projectID, deployment.GetDeploymentName())
	if !admin.IsErrorCode(err, atlas.ServerlessInstanceFromClusterAPI) && !admin.IsErrorCode(err, atlas.FlexFromClusterAPI) && err != nil {
		return nil, err
	}
	if cluster != nil {
		return cluster, nil
	}

	flex, err := ds.GetFlexCluster(ctx, projectID, deployment.GetDeploymentName())
	if !admin.IsErrorCode(err, atlas.NonFlexInFlexAPI) && err != nil {
		return nil, err
	}
	if flex != nil {
		return flex, nil
	}

	// not found
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
	case *Flex:
		flex, _, err := ds.flexAPI.CreateFlexCluster(ctx, deployment.GetProjectID(), flexCreateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}
		return flexFromAtlas(flex), nil
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
	case *Flex:
		flex, _, err := ds.flexAPI.UpdateFlexCluster(ctx, deployment.GetProjectID(), deployment.GetName(), flexUpdateToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return flexFromAtlas(flex), nil
	}

	return nil, errors.New("unable to create deployment: unknown type")
}

func (ds *ProductionAtlasDeployments) DeleteDeployment(ctx context.Context, deployment Deployment) error {
	var err error

	switch deployment.(type) {
	case *Cluster:
		_, err = ds.clustersAPI.DeleteCluster(ctx, deployment.GetProjectID(), deployment.GetName()).Execute()
	case *Flex:
		_, err = ds.flexAPI.DeleteFlexCluster(ctx, deployment.GetProjectID(), deployment.GetName()).Execute()
	}

	if err != nil {
		if admin.IsErrorCode(err, atlas.ClusterNotFound) {
			return nil
		}

		return err
	}

	return nil
}

func (ds *ProductionAtlasDeployments) UpgradeToDedicated(ctx context.Context, currentDeployment, targetDeployment Deployment) (Deployment, error) {
	switch currentDeployment.(type) {
	case *Cluster:
		return nil, errors.New("upgrade from shared to dedicated is not supported")
	case *Flex:
		d := targetDeployment.(*Cluster)
		flex, _, err := ds.flexAPI.TenantUpgrade(ctx, targetDeployment.GetProjectID(), flexUpgradeToAtlas(d)).Execute()
		if err != nil {
			return nil, err
		}

		return flexFromAtlas(flex), nil
	}

	return nil, errors.New("unable to upgrade deployment: unknown type")
}

func (ds *ProductionAtlasDeployments) ClusterWithProcessArgs(ctx context.Context, cluster *Cluster) error {
	config, _, err := ds.clustersAPI.GetProcessArgs(ctx, cluster.GetProjectID(), cluster.GetName()).Execute()
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

	config, _, err := ds.clustersAPI.UpdateProcessArgs(ctx, cluster.GetProjectID(), cluster.GetName(), processArgs).Execute()
	if err != nil {
		return err
	}

	cluster.ProcessArgs = processArgsFromAtlas(config)

	return nil
}

func (ds *ProductionAtlasDeployments) GetCustomZones(ctx context.Context, projectID, clusterName string) (map[string]string, error) {
	geosharding, _, err := ds.globalClusterAPI.GetClusterGlobalWrites(ctx, projectID, clusterName).Execute()
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
	_, _, err := ds.globalClusterAPI.DeleteCustomZoneMapping(ctx, projectID, clusterName).Execute()
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
	geosharding, _, err := ds.globalClusterAPI.GetClusterGlobalWrites(ctx, projectID, clusterName).Execute()
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
	params := &admin.DeleteManagedNamespacesApiParams{
		GroupId:     projectID,
		ClusterName: clusterName,
		Db:          &namespace.Db,
		Collection:  &namespace.Collection,
	}
	_, _, err := ds.globalClusterAPI.DeleteManagedNamespacesWithParams(ctx, params).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete managed namespace: %w", err)
	}
	return nil
}
