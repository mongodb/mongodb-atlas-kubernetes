package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ThirdPartyIntegrationsClientMock struct {
	CreateFunc     func(projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.ThirdPartyIntegration

	ReplaceFunc     func(projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	ReplaceRequests map[string]*mongodbatlas.ThirdPartyIntegration

	DeleteFunc     func(projectID string, integrationType string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}

	GetFunc     func(projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	ListFunc     func(projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}
}

func (c *ThirdPartyIntegrationsClientMock) Create(_ context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.ThirdPartyIntegration{}
	}

	c.CreateRequests[fmt.Sprintf("%s.%s", projectID, integrationType)] = integration

	return c.CreateFunc(projectID, integrationType, integration)
}

func (c *ThirdPartyIntegrationsClientMock) Replace(_ context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	if c.ReplaceRequests == nil {
		c.ReplaceRequests = map[string]*mongodbatlas.ThirdPartyIntegration{}
	}

	c.ReplaceRequests[fmt.Sprintf("%s.%s", projectID, integrationType)] = integration

	return c.ReplaceFunc(projectID, integrationType, integration)
}

func (c *ThirdPartyIntegrationsClientMock) Delete(_ context.Context, projectID string, integrationType string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, integrationType)] = struct{}{}

	return c.DeleteFunc(projectID, integrationType)
}

func (c *ThirdPartyIntegrationsClientMock) Get(_ context.Context, projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, integrationType)] = struct{}{}

	return c.GetFunc(projectID, integrationType)
}

func (c *ThirdPartyIntegrationsClientMock) List(_ context.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}
