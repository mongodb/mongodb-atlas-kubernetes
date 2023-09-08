package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type AdvancedClustersClientMock struct {
	ListFunc     func(projectID string) (*mongodbatlas.AdvancedClustersResponse, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(projectID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.AdvancedCluster

	UpdateFunc     func(projectID string, clusterName string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.AdvancedCluster

	DeleteFunc     func(projectID string, clusterName string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}

	TestFailoverFunc     func(projectID string, clusterName string) (*mongodbatlas.Response, error)
	TestFailoverRequests map[string]struct{}
}

func (c *AdvancedClustersClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) (*mongodbatlas.AdvancedClustersResponse, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *AdvancedClustersClientMock) Get(_ context.Context, projectID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.GetFunc(projectID, clusterName)
}

func (c *AdvancedClustersClientMock) Create(_ context.Context, projectID string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.AdvancedCluster{}
	}

	c.CreateRequests[projectID] = cluster

	return c.CreateFunc(projectID, cluster)
}

func (c *AdvancedClustersClientMock) Update(_ context.Context, projectID string, clusterName string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.AdvancedCluster{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = cluster

	return c.UpdateFunc(projectID, clusterName, cluster)
}

func (c *AdvancedClustersClientMock) Delete(_ context.Context, projectID string, clusterName string, _ *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.DeleteFunc(projectID, clusterName)
}

func (c *AdvancedClustersClientMock) TestFailover(_ context.Context, projectID string, clusterName string) (*mongodbatlas.Response, error) {
	if c.TestFailoverRequests == nil {
		c.TestFailoverRequests = map[string]struct{}{}
	}

	c.TestFailoverRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.TestFailoverFunc(projectID, clusterName)
}
