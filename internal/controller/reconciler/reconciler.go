package reconciler

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

type AtlasReconciler struct {
	Client client.Client
	Log    *zap.SugaredLogger
}

func (r *AtlasReconciler) ResolveProject(ctx context.Context, sdkClient *admin.APIClient, pro project.ProjectReferrerObject) (*project.Project, error) {
	projectsService := project.NewProjectAPIService(sdkClient.ProjectsApi)
	pdr := pro.ProjectDualRef()
	if pdr.ProjectRef != nil {
		project, err := r.projectFromKube(ctx, pro, projectsService)
		if err != nil {
			return nil, fmt.Errorf("failed to query Kubernetes: %w", err)
		}
		return project, nil
	}
	prj, err := r.projectFromAtlas(ctx, projectsService, pdr)
	if err != nil {
		return nil, fmt.Errorf("failed to get Project from Atlas by ID: %w", err)
	}
	return prj, nil
}

func (r *AtlasReconciler) projectFromKube(ctx context.Context, pro project.ProjectReferrerObject, service project.ProjectService) (*project.Project, error) {
	kubeProject, err := r.fetchProject(ctx, pro)
	if err != nil {
		return nil, fmt.Errorf("failed to get Project from Kubernetes: %w", err)
	}
	return service.GetProjectByName(ctx, kubeProject.Spec.Name)
}

func (r *AtlasReconciler) projectFromAtlas(ctx context.Context, projectsService project.ProjectService, pdr *akov2.ProjectDualReference) (*project.Project, error) {
	return projectsService.GetProject(ctx, pdr.ExternalProjectRef.ID)
}

func (r *AtlasReconciler) ResolveCredentials(ctx context.Context, pro project.ProjectReferrerObject) (*client.ObjectKey, error) {
	creds := r.credentialsFor(pro)
	if creds != nil && creds.Name != "" {
		return creds, nil
	}
	project, err := r.fetchProject(ctx, pro)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, nil
	}
	return project.ConnectionSecretObjectKey(), nil
}

func (r *AtlasReconciler) credentialsFor(pro project.ProjectReferrerObject) *client.ObjectKey {
	key := client.ObjectKeyFromObject(pro)
	pdr := pro.ProjectDualRef()
	if pdr.ConnectionSecret == nil {
		return nil
	}
	key.Name = pdr.ConnectionSecret.Name
	return &key
}

func (r *AtlasReconciler) fetchProject(ctx context.Context, pro project.ProjectReferrerObject) (*akov2.AtlasProject, error) {
	pdr := pro.ProjectDualRef()
	if pdr.ProjectRef == nil {
		return nil, nil
	}
	project := akov2.AtlasProject{}
	ns := pro.GetNamespace()
	if pdr.ProjectRef.Namespace != "" {
		ns = pdr.ProjectRef.Namespace
	}
	key := client.ObjectKey{Name: pdr.ProjectRef.Name, Namespace: ns}
	err := r.Client.Get(ctx, key, &project)
	if err != nil {
		return nil, fmt.Errorf("can not fetch AtlasProject: %w", err)
	}
	return &project, nil
}
