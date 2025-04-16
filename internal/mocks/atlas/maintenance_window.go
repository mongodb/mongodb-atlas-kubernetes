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

	"go.mongodb.org/atlas/mongodbatlas"
)

type MaintenanceWindowClientMock struct {
	GetFunc     func(projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	UpdateFunc     func(projectID string, maintenance *mongodbatlas.MaintenanceWindow) (*mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.MaintenanceWindow

	DeferFunc     func(projectID string) (*mongodbatlas.Response, error)
	DeferRequests map[string]struct{}

	AutoDeferFunc     func(projectID string) (*mongodbatlas.Response, error)
	AutoDeferRequests map[string]struct{}

	ResetFunc     func(projectID string) (*mongodbatlas.Response, error)
	ResetRequests map[string]struct{}
}

func (c *MaintenanceWindowClientMock) Get(_ context.Context, projectID string) (*mongodbatlas.MaintenanceWindow, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[projectID] = struct{}{}

	return c.GetFunc(projectID)
}

func (c *MaintenanceWindowClientMock) Update(_ context.Context, projectID string, maintenance *mongodbatlas.MaintenanceWindow) (*mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.MaintenanceWindow{}
	}

	c.UpdateRequests[projectID] = maintenance

	return c.UpdateFunc(projectID, maintenance)
}

func (c *MaintenanceWindowClientMock) Defer(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.DeferRequests == nil {
		c.DeferRequests = map[string]struct{}{}
	}

	c.DeferRequests[projectID] = struct{}{}

	return c.DeferFunc(projectID)
}

func (c *MaintenanceWindowClientMock) AutoDefer(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.AutoDeferRequests == nil {
		c.AutoDeferRequests = map[string]struct{}{}
	}

	c.AutoDeferRequests[projectID] = struct{}{}

	return c.AutoDeferFunc(projectID)
}

func (c *MaintenanceWindowClientMock) Reset(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.ResetRequests == nil {
		c.ResetRequests = map[string]struct{}{}
	}

	c.ResetRequests[projectID] = struct{}{}

	return c.ResetFunc(projectID)
}
