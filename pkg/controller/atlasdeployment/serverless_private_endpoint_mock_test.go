package atlasdeployment

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
)

type ServerlessPrivateEndpointClientMock struct {
	ListFn func(string, string, *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error)
}

func (spec ServerlessPrivateEndpointClientMock) List(_ context.Context, groupID string, instanceName string, opts *mongodbatlas.ListOptions) ([]mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
	if spec.ListFn == nil {
		panic("ListFn not mocked for test")
	}
	return spec.ListFn(groupID, instanceName, opts)
}

func (spec ServerlessPrivateEndpointClientMock) Create(ctx context.Context, groupID string, instanceName string, opts *mongodbatlas.ServerlessPrivateEndpointConnection) (*mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (spec ServerlessPrivateEndpointClientMock) Get(ctx context.Context, groupID string, instanceName string, opts string) (*mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (spec ServerlessPrivateEndpointClientMock) Delete(ctx context.Context, groupID string, instanceName string, opts string) (*mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (spec ServerlessPrivateEndpointClientMock) Update(_ context.Context, _ string, _ string, _ string, _ *mongodbatlas.ServerlessPrivateEndpointConnection) (*mongodbatlas.ServerlessPrivateEndpointConnection, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}
