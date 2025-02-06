package reconciler

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

type AtlasReconciler struct {
	Client client.Client
	Log    *zap.SugaredLogger
}

func (r *AtlasReconciler) ResolveProject(ctx context.Context, sdkClient *admin.APIClient, pro project.ProjectReferrerObject, orgID string) (*project.Project, error) {
	pdr := pro.ProjectDualRef()
	if pdr.ProjectRef != nil {
		project, err := r.projectFromKube(ctx, pro, orgID)
		if err != nil {
			return nil, fmt.Errorf("failed to query Kubernetes: %w", err)
		}
		return project, nil
	}
	projectsService := project.NewProjectAPIService(sdkClient.ProjectsApi)
	prj, err := r.projectFromAtlas(ctx, projectsService, pdr)
	if err != nil {
		return nil, fmt.Errorf("failed to get Project from Atlas by ID: %w", err)
	}
	return prj, nil
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

func (r *AtlasReconciler) Skip(ctx context.Context, typeName string, resource api.AtlasCustomResource, spec any) (ctrl.Result, error) {
	msg := fmt.Sprintf("-> Skipping %s reconciliation as annotation %s=%s",
		typeName, customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip)
	r.Log.Infow(msg, "spec", spec)
	if !resource.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, resource, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err)
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult(), nil
		}
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasReconciler) Invalidate(typeName string, invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("%T is invalid: %v", typeName, invalid)
	return invalid.ReconcileResult(), nil
}

func (r *AtlasReconciler) Unsupport(ctx *workflow.Context, typeName string) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported,
		fmt.Errorf("the %s is not supported by Atlas for government", typeName),
	).WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

func (r *AtlasReconciler) projectFromKube(ctx context.Context, pro project.ProjectReferrerObject, orgID string) (*project.Project, error) {
	kubeProject, err := r.fetchProject(ctx, pro)
	if err != nil {
		return nil, fmt.Errorf("failed to get Project from Kubernetes: %w", err)
	}
	return project.NewProject(kubeProject, orgID), nil
}

func (r *AtlasReconciler) projectFromAtlas(ctx context.Context, projectsService project.ProjectService, pdr *akov2.ProjectDualReference) (*project.Project, error) {
	return projectsService.GetProject(ctx, pdr.ExternalProjectRef.ID)
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
