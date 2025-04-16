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

type CloudProviderAccessClientMock struct {
	ListRolesFunc     func(projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error)
	ListRolesRequests map[string]struct{}

	GetRoleFunc     func(projectID string, roleID string) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error)
	GetRoleRequests map[string]struct{}

	CreateRoleFunc     func(projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error)
	CreateRoleRequests map[string]*mongodbatlas.CloudProviderAccessRoleRequest

	AuthorizeRoleFunc     func(projectID, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error)
	AuthorizeRoleRequests map[string]*mongodbatlas.CloudProviderAccessRoleRequest

	DeauthorizeRoleFunc     func(cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error)
	DeauthorizeRoleRequests []*mongodbatlas.CloudProviderDeauthorizationRequest
}

func (c *CloudProviderAccessClientMock) ListRoles(_ context.Context, projectID string) (*mongodbatlas.CloudProviderAccessRoles, *mongodbatlas.Response, error) {
	if c.ListRolesRequests == nil {
		c.ListRolesRequests = map[string]struct{}{}
	}

	c.ListRolesRequests[projectID] = struct{}{}

	return c.ListRolesFunc(projectID)
}

func (c *CloudProviderAccessClientMock) GetRole(_ context.Context, projectID string, roleID string) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	if c.GetRoleRequests == nil {
		c.GetRoleRequests = map[string]struct{}{}
	}

	c.GetRoleRequests[fmt.Sprintf("%s.%s", projectID, roleID)] = struct{}{}

	return c.GetRoleFunc(projectID, roleID)
}

func (c *CloudProviderAccessClientMock) CreateRole(_ context.Context, projectID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	if c.CreateRoleRequests == nil {
		c.CreateRoleRequests = map[string]*mongodbatlas.CloudProviderAccessRoleRequest{}
	}

	c.CreateRoleRequests[projectID] = cpa

	return c.CreateRoleFunc(projectID, cpa)
}

func (c *CloudProviderAccessClientMock) AuthorizeRole(_ context.Context, projectID string, roleID string, cpa *mongodbatlas.CloudProviderAccessRoleRequest) (*mongodbatlas.CloudProviderAccessRole, *mongodbatlas.Response, error) {
	if c.AuthorizeRoleRequests == nil {
		c.AuthorizeRoleRequests = map[string]*mongodbatlas.CloudProviderAccessRoleRequest{}
	}

	c.AuthorizeRoleRequests[fmt.Sprintf("%s.%s", projectID, roleID)] = cpa

	return c.AuthorizeRoleFunc(projectID, roleID, cpa)
}
func (c *CloudProviderAccessClientMock) DeauthorizeRole(_ context.Context, cpa *mongodbatlas.CloudProviderDeauthorizationRequest) (*mongodbatlas.Response, error) {
	if c.DeauthorizeRoleRequests == nil {
		c.DeauthorizeRoleRequests = []*mongodbatlas.CloudProviderDeauthorizationRequest{}
	}

	c.DeauthorizeRoleRequests = append(c.DeauthorizeRoleRequests, cpa)

	return c.DeauthorizeRoleFunc(cpa)
}
