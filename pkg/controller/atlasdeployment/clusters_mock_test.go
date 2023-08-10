package atlasdeployment

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
)

type advancedClustersClientMock struct {
	GetFn    func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error)
	DeleteFn func(groupID string, clusterName string, options *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error)
}

func (ac *advancedClustersClientMock) List(ctx context.Context, groupID string, options *mongodbatlas.ListOptions) (*mongodbatlas.AdvancedClustersResponse, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (ac *advancedClustersClientMock) Get(_ context.Context, groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	if ac.GetFn == nil {
		panic("GetFn not mocked for test")
	}
	return ac.GetFn(groupID, clusterName)
}

func (ac *advancedClustersClientMock) Create(ctx context.Context, groupID string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (ac *advancedClustersClientMock) Update(ctx context.Context, groupID string, clusterName string, cluster *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (ac *advancedClustersClientMock) Delete(ctx context.Context, groupID string, clusterName string, options *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error) {
	if ac.DeleteFn == nil {
		panic("GetFn not mocked for test")
	}
	return ac.DeleteFn(groupID, clusterName, options)
}

func (ac *advancedClustersClientMock) TestFailover(ctx context.Context, groupID string, clusterName string) (*mongodbatlas.Response, error) {
	panic("not implemented") // TODO: Implement
}
