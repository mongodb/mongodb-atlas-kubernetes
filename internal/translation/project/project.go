package project

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

// ProjectReferrer is anything that holds a ProjectDualReference
type ProjectReferrer interface {
	ProjectDualRef() *akov2.ProjectDualReference
}

// ProjectReferrerObject is an project referrer that is also an Kubernetes Object
type ProjectReferrerObject interface {
	client.Object
	ProjectReferrer
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
	group, _, err := a.projectAPI.GetProjectByName(ctx, name).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "NOT_IN_GROUP") || admin.IsErrorCode(err, "RESOURCE_NOT_FOUND") {
			return nil, nil
		}

		return nil, err
	}

	return fromAtlas(group), err
}

func (a *ProjectAPI) GetProject(ctx context.Context, ID string) (*Project, error) {
	group, _, err := a.projectAPI.GetProject(ctx, ID).Execute()
	if err != nil {
		return nil, err
	}

	return fromAtlas(group), err
}

func (a *ProjectAPI) CreateProject(ctx context.Context, project *Project) error {
	group, _, err := a.projectAPI.CreateProject(ctx, toAtlas(project)).Execute()
	if err != nil {
		return err
	}

	project.OrgID = group.GetOrgId()
	project.ID = group.GetId()

	return nil
}

func (a *ProjectAPI) DeleteProject(ctx context.Context, project *Project) error {
	_, _, err := a.projectAPI.DeleteProject(ctx, project.ID).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "GROUP_NOT_FOUND") || admin.IsErrorCode(err, "RESOURCE_NOT_FOUND") {
			return nil
		}

		return err
	}

	return nil
}

func NewProjectAPIService(sdk admin.ProjectsApi) *ProjectAPI {
	return &ProjectAPI{
		projectAPI: sdk,
	}
}
