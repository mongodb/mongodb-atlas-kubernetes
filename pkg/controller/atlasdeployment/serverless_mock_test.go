package atlasdeployment

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
)

type serverlessClientMock struct {
	GetFn func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error)
}

func (sc *serverlessClientMock) List(_ context.Context, _ string, _ *mongodbatlas.ListOptions) (*mongodbatlas.ClustersResponse, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (sc *serverlessClientMock) Get(_ context.Context, groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	if sc.GetFn == nil {
		panic("GetFn not mocked for test")
	}
	return sc.GetFn(groupID, name)
}

func (sc *serverlessClientMock) Create(_ context.Context, _ string, _ *mongodbatlas.ServerlessCreateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (sc *serverlessClientMock) Update(_ context.Context, _ string, _ string, _ *mongodbatlas.ServerlessUpdateRequestParams) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (sc *serverlessClientMock) Delete(_ context.Context, _ string, _ string) (*mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}
