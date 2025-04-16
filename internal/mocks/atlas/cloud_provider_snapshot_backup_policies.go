// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
