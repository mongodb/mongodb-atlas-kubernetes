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

package project

import (
	"context"
	"errors"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
)

// ProjectReferrerObject is a Kube client object that includes references to Atlas projects.
type ProjectReferrerObject interface {
	client.Object
	ProjectDualRef() *akov2.ProjectDualReference
}

type ProjectService interface {
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	GetProject(ctx context.Context, ID string) (*Project, error)
	CreateProject(ctx context.Context, project *Project) error
	DeleteProject(ctx context.Context, project *Project) error
}

type ProjectAPI struct {
	projectAPI admin.ProjectsApi
}

func (a *ProjectAPI) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	group, _, err := a.projectAPI.GetGroupByName(ctx, name).Execute()
	if err != nil {
		return nil, translateError(err)
	}

	return fromAtlas(group), nil
}

func (a *ProjectAPI) GetProject(ctx context.Context, ID string) (*Project, error) {
	group, _, err := a.projectAPI.GetGroup(ctx, ID).Execute()
	if err != nil {
		return nil, translateError(err)
	}

	return fromAtlas(group), nil
}

func (a *ProjectAPI) CreateProject(ctx context.Context, project *Project) error {
	group, _, err := a.projectAPI.CreateGroup(ctx, toAtlas(project)).Execute()
	if err != nil {
		return err
	}

	project.OrgID = group.GetOrgId()
	project.ID = group.GetId()

	return nil
}

func (a *ProjectAPI) DeleteProject(ctx context.Context, project *Project) error {
	_, err := a.projectAPI.DeleteGroup(ctx, project.ID).Execute()
	err = translateError(err)
	if err != nil && !errors.Is(err, translation.ErrNotFound) {
		return err
	}

	return nil
}

func NewProjectAPIService(sdk admin.ProjectsApi) *ProjectAPI {
	return &ProjectAPI{
		projectAPI: sdk,
	}
}

func translateError(err error) error {
	switch {
	case admin.IsErrorCode(err, "RESOURCE_NOT_FOUND"):
	case admin.IsErrorCode(err, "GROUP_NOT_FOUND"):
	case admin.IsErrorCode(err, "NOT_IN_GROUP"):
	default:
		return err
	}

	return errors.Join(translation.ErrNotFound, err)
}
