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

type EncryptionAtRestClientMock struct {
	CreateFunc     func(ear *mongodbatlas.EncryptionAtRest) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error)
	CreateRequests []*mongodbatlas.EncryptionAtRest

	GetFunc     func(projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error)
	GetRequests map[string]struct{}

	DeleteFunc     func(projectID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *EncryptionAtRestClientMock) Create(_ context.Context, ear *mongodbatlas.EncryptionAtRest) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = []*mongodbatlas.EncryptionAtRest{}
	}

	c.CreateRequests = append(c.CreateRequests, ear)

	return c.CreateFunc(ear)
}

func (c *EncryptionAtRestClientMock) Get(_ context.Context, projectID string) (*mongodbatlas.EncryptionAtRest, *mongodbatlas.Response, error) {
	if c.GetRequests == nil {
		c.GetRequests = map[string]struct{}{}
	}

	c.GetRequests[projectID] = struct{}{}

	return c.GetFunc(projectID)
}
func (c *EncryptionAtRestClientMock) Delete(_ context.Context, projectID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[projectID] = struct{}{}

	return c.DeleteFunc(projectID)
}
