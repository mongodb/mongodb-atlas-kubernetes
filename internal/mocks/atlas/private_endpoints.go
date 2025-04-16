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

type PrivateEndpointsClientMock struct {
	CreateFunc     func(projectID string, endpoint *mongodbatlas.PrivateEndpointConnection) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.PrivateEndpointConnection

	GetFunc     func(projectID string, cloudProvider string, endpointServiceID string) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	ListFunc     func(projectID, providerName string) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	DeleteFunc     func(projectID string, cloudProvider string, endpointServiceID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}

	AddOnePrivateEndpointFunc     func(projectID string, cloudProvider string, endpointServiceID string, endpoint *mongodbatlas.InterfaceEndpointConnection) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error)
	AddOnePrivateEndpointRequests map[string]*mongodbatlas.InterfaceEndpointConnection

	GetOnePrivateEndpointFunc     func(projectID string, cloudProvider string, endpointServiceID string, privateEndpointID string) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error)
	GetOnePrivateEndpointRequests map[string]struct{}

	DeleteOnePrivateEndpointFunc     func(projectID string, cloudProvider string, endpointServiceID string, privateEndpointID string) (*mongodbatlas.Response, error)
	DeleteOnePrivateEndpointRequests map[string]struct{}

	UpdateRegionalizedPrivateEndpointSettingFunc     func(projectID string, enabled bool) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error)
	UpdateRegionalizedPrivateEndpointSettingRequests map[string]bool

	GetRegionalizedPrivateEndpointSettingFunc     func(projectID string) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error)
	GetRegionalizedPrivateEndpointSettingRequests map[string]struct{}
}

func (c *PrivateEndpointsClientMock) Create(_ context.Context, projectID string, endpoint *mongodbatlas.PrivateEndpointConnection) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.PrivateEndpointConnection{}
	}

	c.CreateRequests[projectID] = endpoint

	return c.CreateFunc(projectID, endpoint)
}

func (c *PrivateEndpointsClientMock) Get(_ context.Context, projectID string, cloudProvider string, endpointServiceID string) (*mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s.%s", projectID, cloudProvider, endpointServiceID)] = struct{}{}

	return c.GetFunc(projectID, cloudProvider, endpointServiceID)
}

func (c *PrivateEndpointsClientMock) List(_ context.Context, projectID string, cloudProvider string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.PrivateEndpointConnection, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[fmt.Sprintf("%s.%s", projectID, cloudProvider)] = struct{}{}

	return c.ListFunc(projectID, cloudProvider)
}

func (c *PrivateEndpointsClientMock) Delete(_ context.Context, projectID string, cloudProvider string, endpointServiceID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s.%s", projectID, cloudProvider, endpointServiceID)] = struct{}{}

	return c.DeleteFunc(projectID, cloudProvider, endpointServiceID)
}

func (c *PrivateEndpointsClientMock) AddOnePrivateEndpoint(_ context.Context, projectID string, cloudProvider string, endpointServiceID string, endpoint *mongodbatlas.InterfaceEndpointConnection) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error) {
	if c.AddOnePrivateEndpointRequests == nil {
		c.AddOnePrivateEndpointRequests = map[string]*mongodbatlas.InterfaceEndpointConnection{}
	}

	c.AddOnePrivateEndpointRequests[fmt.Sprintf("%s.%s.%s", projectID, cloudProvider, endpointServiceID)] = endpoint

	return c.AddOnePrivateEndpointFunc(projectID, cloudProvider, endpointServiceID, endpoint)
}

func (c *PrivateEndpointsClientMock) GetOnePrivateEndpoint(_ context.Context, projectID string, cloudProvider string, endpointServiceID string, privateEndpointID string) (*mongodbatlas.InterfaceEndpointConnection, *mongodbatlas.Response, error) {
	if c.GetOnePrivateEndpointRequests == nil {
		c.GetOnePrivateEndpointRequests = map[string]struct{}{}
	}

	c.GetOnePrivateEndpointRequests[fmt.Sprintf("%s.%s.%s.%s", projectID, cloudProvider, endpointServiceID, privateEndpointID)] = struct{}{}

	return c.GetOnePrivateEndpointFunc(projectID, cloudProvider, endpointServiceID, privateEndpointID)
}

func (c *PrivateEndpointsClientMock) DeleteOnePrivateEndpoint(_ context.Context, projectID string, cloudProvider string, endpointServiceID string, privateEndpointID string) (*mongodbatlas.Response, error) {
	if c.DeleteOnePrivateEndpointRequests == nil {
		c.DeleteOnePrivateEndpointRequests = map[string]struct{}{}
	}

	c.DeleteOnePrivateEndpointRequests[fmt.Sprintf("%s.%s.%s.%s", projectID, cloudProvider, endpointServiceID, privateEndpointID)] = struct{}{}

	return c.DeleteOnePrivateEndpointFunc(projectID, cloudProvider, endpointServiceID, privateEndpointID)
}

func (c *PrivateEndpointsClientMock) UpdateRegionalizedPrivateEndpointSetting(_ context.Context, projectID string, enabled bool) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error) {
	if c.UpdateRegionalizedPrivateEndpointSettingRequests == nil {
		c.UpdateRegionalizedPrivateEndpointSettingRequests = map[string]bool{}
	}

	c.UpdateRegionalizedPrivateEndpointSettingRequests[projectID] = enabled

	return c.UpdateRegionalizedPrivateEndpointSettingFunc(projectID, enabled)
}

func (c *PrivateEndpointsClientMock) GetRegionalizedPrivateEndpointSetting(_ context.Context, projectID string) (*mongodbatlas.RegionalizedPrivateEndpointSetting, *mongodbatlas.Response, error) {
	if c.GetRegionalizedPrivateEndpointSettingRequests == nil {
		c.GetRegionalizedPrivateEndpointSettingRequests = map[string]struct{}{}
	}

	c.GetRegionalizedPrivateEndpointSettingRequests[projectID] = struct{}{}

	return c.GetRegionalizedPrivateEndpointSettingFunc(projectID)
}
