package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type CustomRolesClientMock struct {
	ListFunc     func(projectID string) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(projectID string, customRoleID string) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, customRole *mongodbatlas.CustomDBRole) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.CustomDBRole

	UpdateFunc     func(projectID string, customRoleID string, customRole *mongodbatlas.CustomDBRole) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.CustomDBRole

	DeleteFunc     func(projectID string, customRoleID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *CustomRolesClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) (*[]mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *CustomRolesClientMock) Get(_ context.Context, projectID string, customRoleID string) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, customRoleID)] = struct{}{}

	return c.GetFunc(projectID, customRoleID)
}

func (c *CustomRolesClientMock) Create(_ context.Context, projectID string, customRole *mongodbatlas.CustomDBRole) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.CustomDBRole{}
	}

	c.CreateRequests[projectID] = customRole

	return c.CreateFunc(projectID, customRole)
}

func (c *CustomRolesClientMock) Update(_ context.Context, projectID string, customRoleID string, customRole *mongodbatlas.CustomDBRole) (*mongodbatlas.CustomDBRole, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.CustomDBRole{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, customRoleID)] = customRole

	return c.UpdateFunc(projectID, customRoleID, customRole)
}

func (c *CustomRolesClientMock) Delete(_ context.Context, projectID string, customRoleID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, customRoleID)] = struct{}{}

	return c.DeleteFunc(projectID, customRoleID)
}
