package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type NetworkPeeringClientMock struct {
	ListFunc     func(projectID string) ([]mongodbatlas.Peer, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	GetFunc     func(projectID string, peerID string) (*mongodbatlas.Peer, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, peer *mongodbatlas.Peer) (*mongodbatlas.Peer, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.Peer

	UpdateFunc     func(projectID string, peerID string, peer *mongodbatlas.Peer) (*mongodbatlas.Peer, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.Peer

	DeleteFunc     func(projectID string, peerID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *NetworkPeeringClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ContainersListOptions) ([]mongodbatlas.Peer, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *NetworkPeeringClientMock) Get(_ context.Context, projectID string, peerID string) (*mongodbatlas.Peer, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, peerID)] = struct{}{}

	return c.GetFunc(projectID, peerID)
}

func (c *NetworkPeeringClientMock) Create(_ context.Context, projectID string, peer *mongodbatlas.Peer) (*mongodbatlas.Peer, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.Peer{}
	}

	c.CreateRequests[projectID] = peer

	return c.CreateFunc(projectID, peer)
}

func (c *NetworkPeeringClientMock) Update(_ context.Context, projectID string, peerID string, peer *mongodbatlas.Peer) (*mongodbatlas.Peer, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.Peer{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, peerID)] = peer

	return c.UpdateFunc(projectID, peerID, peer)
}

func (c *NetworkPeeringClientMock) Delete(_ context.Context, projectID string, peerID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, peerID)] = struct{}{}

	return c.DeleteFunc(projectID, peerID)
}
