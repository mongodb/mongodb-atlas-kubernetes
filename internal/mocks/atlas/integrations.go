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

type IntegrationsMock struct {
	CreateFunc  func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	ReplaceFunc func(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
	DeleteFunc  func(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.Response, error)
	GetFunc     func(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error)
	ListFunc    func(ctx context.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error)
}

func (im *IntegrationsMock) Create(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.CreateFunc(ctx, projectID, integrationType, integration)
}

func (im *IntegrationsMock) Replace(ctx context.Context, projectID string, integrationType string, integration *mongodbatlas.ThirdPartyIntegration) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.ReplaceFunc(ctx, projectID, integrationType, integration)
}

func (im *IntegrationsMock) Delete(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.Response, error) {
	return im.DeleteFunc(ctx, projectID, integrationType)
}

func (im *IntegrationsMock) Get(ctx context.Context, projectID string, integrationType string) (*mongodbatlas.ThirdPartyIntegration, *mongodbatlas.Response, error) {
	return im.GetFunc(ctx, projectID, integrationType)
}

func (im *IntegrationsMock) List(ctx context.Context, projectID string) (*mongodbatlas.ThirdPartyIntegrations, *mongodbatlas.Response, error) {
	return im.ListFunc(ctx, projectID)
}
