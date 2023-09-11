package atlas

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
)

type AuditingClientMock struct {
	GetFunc     func(projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	ConfigureFunc     func(projectID string, auditing *mongodbatlas.Auditing) (*mongodbatlas.Auditing, *mongodbatlas.Response, error)
	ConfigureRequests map[string]*mongodbatlas.Auditing
}

func (c *AuditingClientMock) Get(_ context.Context, projectID string) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[projectID] = struct{}{}

	return c.GetFunc(projectID)
}
func (c *AuditingClientMock) Configure(_ context.Context, projectID string, auditing *mongodbatlas.Auditing) (*mongodbatlas.Auditing, *mongodbatlas.Response, error) {
	if c.ConfigureRequests == nil {
		c.ConfigureRequests = map[string]*mongodbatlas.Auditing{}
	}

	c.ConfigureRequests[projectID] = auditing

	return c.ConfigureFunc(projectID, auditing)
}
