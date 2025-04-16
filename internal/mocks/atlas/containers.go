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

type ContainerClientMock struct {
	ListFunc     func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	ListAllFunc     func(projectID string) ([]mongodbatlas.Container, *mongodbatlas.Response, error)
	ListAllRequests map[string]struct{}

	GetFunc     func(projectID string, containerID string) (*mongodbatlas.Container, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	CreateFunc     func(projectID string, container *mongodbatlas.Container) (*mongodbatlas.Container, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.Container

	UpdateFunc     func(projectID string, containerID string, container *mongodbatlas.Container) (*mongodbatlas.Container, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.Container

	DeleteFunc     func(projectID string, containerID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *ContainerClientMock) List(_ context.Context, projectID string, _ *mongodbatlas.ContainersListOptions) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}

func (c *ContainerClientMock) ListAll(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.Container, *mongodbatlas.Response, error) {
	if c.ListAllRequests == nil {
		c.ListAllRequests = map[string]struct{}{}
	}

	c.ListAllRequests[projectID] = struct{}{}

	return c.ListAllFunc(projectID)
}

func (c *ContainerClientMock) Get(_ context.Context, projectID string, containerID string) (*mongodbatlas.Container, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[fmt.Sprintf("%s.%s", projectID, containerID)] = struct{}{}

	return c.GetFunc(projectID, containerID)
}

func (c *ContainerClientMock) Create(_ context.Context, projectID string, container *mongodbatlas.Container) (*mongodbatlas.Container, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.Container{}
	}

	c.CreateRequests[projectID] = container

	return c.CreateFunc(projectID, container)
}

func (c *ContainerClientMock) Update(_ context.Context, projectID string, containerID string, container *mongodbatlas.Container) (*mongodbatlas.Container, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.Container{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, containerID)] = container

	return c.UpdateFunc(projectID, containerID, container)
}

func (c *ContainerClientMock) Delete(_ context.Context, projectID string, containerID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, containerID)] = struct{}{}

	return c.DeleteFunc(projectID, containerID)
}
