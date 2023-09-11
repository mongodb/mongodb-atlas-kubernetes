package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ProjectIPAccessListClientMock struct {
	ListFunc     func(projectID string) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(projectID string, entry string) (*mongodbatlas.ProjectIPAccessList, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error)
	CreateRequests map[string][]*mongodbatlas.ProjectIPAccessList

	DeleteFunc     func(projectID, entry string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *ProjectIPAccessListClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *ProjectIPAccessListClientMock) Get(_ context.Context, projectID string, entry string) (*mongodbatlas.ProjectIPAccessList, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, entry)] = struct{}{}

	return c.GetFunc(projectID, entry)
}

func (c *ProjectIPAccessListClientMock) Create(_ context.Context, projectID string, ipAccessLists []*mongodbatlas.ProjectIPAccessList) (*mongodbatlas.ProjectIPAccessLists, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string][]*mongodbatlas.ProjectIPAccessList{}
	}

	c.CreateRequests[projectID] = ipAccessLists

	return c.CreateFunc(projectID, ipAccessLists)
}

func (c *ProjectIPAccessListClientMock) Delete(_ context.Context, projectID, entry string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, entry)] = struct{}{}

	return c.DeleteFunc(projectID, entry)
}
