package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type CloudProviderSnapshotBackupPoliciesClientMock struct {
	GetFunc     func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	UpdateFunc     func(projectID string, clusterName string, backup *mongodbatlas.CloudProviderSnapshotBackupPolicy) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.CloudProviderSnapshotBackupPolicy

	DeleteFunc     func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *CloudProviderSnapshotBackupPoliciesClientMock) Get(_ context.Context, projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.GetFunc(projectID, clusterName)
}

func (c *CloudProviderSnapshotBackupPoliciesClientMock) Update(_ context.Context, projectID string, clusterName string, backup *mongodbatlas.CloudProviderSnapshotBackupPolicy) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.CloudProviderSnapshotBackupPolicy{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = backup

	return c.UpdateFunc(projectID, clusterName, backup)
}

func (c *CloudProviderSnapshotBackupPoliciesClientMock) Delete(_ context.Context, projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, clusterName)] = struct{}{}

	return c.DeleteFunc(projectID, clusterName)
}
