package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type GlobalClustersClientMock struct {
	GetFunc     func(projectID string, clusterName string) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	AddManagedNamespaceFunc     func(projectID string, clusterName string, namespace *mongodbatlas.ManagedNamespace) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error)
	AddManagedNamespaceRequests map[string]*mongodbatlas.ManagedNamespace

	DeleteManagedNamespaceFunc     func(projectID string, clusterName string, namespace *mongodbatlas.ManagedNamespace) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error)
	DeleteManagedNamespaceRequests map[string]*mongodbatlas.ManagedNamespace

	AddCustomZoneMappingsFunc     func(projectID string, clusterName string, request *mongodbatlas.CustomZoneMappingsRequest) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error)
	AddCustomZoneMappingsRequests map[string]*mongodbatlas.CustomZoneMappingsRequest

	DeleteCustomZoneMappingsFunc     func(projectID string, clusterName string) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error)
	DeleteCustomZoneMappingsRequests map[string]struct{}
}

func (c *GlobalClustersClientMock) Get(
	_ context.Context,
	projectID string,
	clusterName string,
) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.GetFunc(projectID, clusterName)
}

func (c *GlobalClustersClientMock) AddManagedNamespace(
	_ context.Context,
	projectID string,
	clusterName string,
	namespace *mongodbatlas.ManagedNamespace,
) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
	if c.AddManagedNamespaceRequests == nil {
		c.AddManagedNamespaceRequests = map[string]*mongodbatlas.ManagedNamespace{}
	}

	c.AddManagedNamespaceRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = namespace

	return c.AddManagedNamespaceFunc(projectID, clusterName, namespace)
}

func (c *GlobalClustersClientMock) DeleteManagedNamespace(
	_ context.Context,
	projectID string,
	clusterName string,
	namespace *mongodbatlas.ManagedNamespace,
) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
	if c.DeleteManagedNamespaceRequests == nil {
		c.DeleteManagedNamespaceRequests = map[string]*mongodbatlas.ManagedNamespace{}
	}

	c.DeleteManagedNamespaceRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = namespace

	return c.DeleteManagedNamespaceFunc(projectID, clusterName, namespace)
}

func (c *GlobalClustersClientMock) AddCustomZoneMappings(
	_ context.Context,
	projectID string,
	clusterName string,
	request *mongodbatlas.CustomZoneMappingsRequest,
) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
	if c.AddCustomZoneMappingsRequests == nil {
		c.AddCustomZoneMappingsRequests = map[string]*mongodbatlas.CustomZoneMappingsRequest{}
	}

	c.AddCustomZoneMappingsRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = request

	return c.AddCustomZoneMappingsFunc(projectID, clusterName, request)
}

func (c *GlobalClustersClientMock) DeleteCustomZoneMappings(
	_ context.Context,
	projectID string,
	clusterName string,
) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
	if c.DeleteCustomZoneMappingsRequests == nil {
		c.DeleteCustomZoneMappingsRequests = map[string]struct{}{}
	}

	c.DeleteCustomZoneMappingsRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.DeleteCustomZoneMappingsFunc(projectID, clusterName)
}
