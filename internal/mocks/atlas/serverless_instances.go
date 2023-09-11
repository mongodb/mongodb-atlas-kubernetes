package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ServerlessInstancesClientMock struct {
	ListFunc     func(projectID string) (*mongodbatlas.ClustersResponse, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(projectID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, instance *mongodbatlas.ServerlessCreateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.ServerlessCreateRequestParams

	UpdateFunc     func(projectID string, name string, instance *mongodbatlas.ServerlessUpdateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.ServerlessUpdateRequestParams

	DeleteFunc     func(projectID string, name string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *ServerlessInstancesClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) (*mongodbatlas.ClustersResponse, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.DeleteRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *ServerlessInstancesClientMock) Get(_ context.Context, projectID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, name)] = struct{}{}

	return c.GetFunc(projectID, name)
}

func (c *ServerlessInstancesClientMock) Create(_ context.Context, projectID string, instance *mongodbatlas.ServerlessCreateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.ServerlessCreateRequestParams{}
	}

	c.CreateRequests[projectID] = instance

	return c.CreateFunc(projectID, instance)
}

func (c *ServerlessInstancesClientMock) Update(_ context.Context, projectID string, name string, instance *mongodbatlas.ServerlessUpdateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.ServerlessUpdateRequestParams{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, name)] = instance

	return c.UpdateFunc(projectID, name, instance)
}

func (c *ServerlessInstancesClientMock) Delete(_ context.Context, projectID string, name string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, name)] = struct{}{}

	return c.DeleteFunc(projectID, name)
}
