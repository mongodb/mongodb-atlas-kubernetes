package reconciler

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type Reconciler struct {
	client.Client
	AtlasProvider atlas.Provider
	Log           *zap.SugaredLogger
}

func (r *Reconciler) SelectCredentials(ctx context.Context, projectRefs *api.ProjectReferences, resource api.AtlasCustomResource) (*client.ObjectKey, error) {
	if projectRefs.ExternalProjectRef != nil {
		if projectRefs.ConnectionSecret == nil {
			return nil, errors.New("externalProjectIDRef is set but connectionSecret is missing")
		}
		return &client.ObjectKey{Name: resource.GetName(), Namespace: resource.GetNamespace()}, nil
	}
	if projectRefs.Project != nil {
		if projectRefs.ConnectionSecret != nil {
			return &client.ObjectKey{Name: projectRefs.ConnectionSecret.Name, Namespace: resource.GetNamespace()}, nil
		}
		project := &akov2.AtlasProject{}
		err := r.Client.Get(ctx,
			client.ObjectKey{Name: projectRefs.Project.Name, Namespace: projectRefs.Project.Namespace}, project)
		if err != nil {
			return nil, fmt.Errorf("can not read AtlasProject %q from Kubernetes: %w", projectRefs.Project.Name, err)
		}
		return project.ConnectionSecretObjectKey(), nil
	}
	return nil, errors.New("either 'externalProjectIDRef' or 'projectRef' must be set for the AtlasCustomRole resource")
}

func (r *Reconciler) GetProjectID(ctx context.Context, projectRefs *api.ProjectReferences, atlasClient *admin.APIClient, ns string) (string, error) {
	if projectRefs.ExternalProjectRef != nil {
		projectsService := project.NewProjectAPIService(atlasClient.ProjectsApi)
		if _, err := projectsService.GetProject(ctx, projectRefs.ExternalProjectRef.ID); err != nil {
			return "", fmt.Errorf("failed to verify Atlas Project from external reference: %w", err)
		}
		return projectRefs.ExternalProjectRef.ID, nil
	} else {
		atlasProject := &akov2.AtlasProject{}
		if err := r.Client.Get(ctx, objectKey(projectRefs, ns), atlasProject); err != nil {
			return "", fmt.Errorf("failed to get Project from Kubernetes: %w", err)
		}
		return atlasProject.ID(), nil
	}
}

func (r *Reconciler) Skip(ctx context.Context, typeName string, resource api.AtlasCustomResource, spec any) (ctrl.Result, error) {
	msg := fmt.Sprintf("-> Skipping %s reconciliation as annotation %s=%s",
		typeName, customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip)
	r.Log.Infow(msg, "spec", spec)
	if !resource.GetDeletionTimestamp().IsZero() {
		if err := customresource.ManageFinalizer(ctx, r.Client, resource, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.Internal, err.Error())
			r.Log.Errorw("Failed to remove finalizer", "terminate", err)

			return result.ReconcileResult(), nil
		}
	}

	return workflow.OK().ReconcileResult(), nil
}

func (r *Reconciler) Invalidate(invalid workflow.Result) (ctrl.Result, error) {
	// note: ValidateResourceVersion already set the state so we don't have to do it here.
	r.Log.Debugf("AtlasNetworkPeering is invalid: %v", invalid)
	return invalid.ReconcileResult(), nil
}

func (r *Reconciler) Unsupport(ctx *workflow.Context, typeName string) (ctrl.Result, error) {
	unsupported := workflow.Terminate(
		workflow.AtlasGovUnsupported,
		fmt.Sprintf("the %s is not supported by Atlas for government", typeName),
	).WithoutRetry()
	ctx.SetConditionFromResult(api.ReadyType, unsupported)
	return unsupported.ReconcileResult(), nil
}

func objectKey(projectRefs *api.ProjectReferences, fallbackNamespace string) client.ObjectKey {
	ns := fallbackNamespace
	if projectRefs.Project.Namespace != "" {
		ns = projectRefs.Project.Namespace
	}
	return kube.ObjectKey(ns, projectRefs.Project.Name)
}
