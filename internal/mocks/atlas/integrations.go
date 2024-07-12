package atlas

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
)

type IntegrationsMock struct {
	CreateFunc  func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	ReplaceFunc func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	DeleteFunc  func(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.Response, error)
	GetFunc     func(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error)
	ListFunc    func(ctx context.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
}

func (im *IntegrationsMock) Create(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.CreateFunc(ctx, projectID, integrationType, integration)
}

func (im *IntegrationsMock) Replace(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.ReplaceFunc(ctx, projectID, integrationType, integration)
}

func (im *IntegrationsMock) Delete(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.Response, error) {
	return im.DeleteFunc(ctx, projectID, integrationType)
}

func (im *IntegrationsMock) Get(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error) {
	return im.GetFunc(ctx, projectID, integrationType)
}

func (im *IntegrationsMock) List(ctx context.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.ListFunc(ctx, projectID)
}
