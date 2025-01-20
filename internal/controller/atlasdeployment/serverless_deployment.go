package atlasdeployment

import (
	"errors"
	"fmt"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func (r *AtlasDeploymentReconciler) handleServerlessInstance(ctx *workflow.Context, projectService project.ProjectService, deploymentService deployment.AtlasDeploymentsService, akoDeployment, atlasDeployment deployment.Deployment) (ctrl.Result, error) {
	akoServerless, ok := akoDeployment.(*deployment.Serverless)
	if !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in AKO is not a serverless cluster"))
	}
	atlasServerless, _ := atlasDeployment.(*deployment.Serverless)

	if atlasServerless == nil {
		ctx.Log.Infof("Serverless Instance %s doesn't exist in Atlas - creating", akoServerless.GetName())
		newServerlessDeployment, err := deploymentService.CreateDeployment(ctx.Context, akoServerless)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		atlasServerless = newServerlessDeployment.(*deployment.Serverless)
	}

	switch atlasServerless.GetState() {
	case status.StateIDLE:
		if !reflect.DeepEqual(akoServerless.ServerlessSpec, atlasServerless.ServerlessSpec) {
			_, err := deploymentService.UpdateDeployment(ctx.Context, akoServerless)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.DeploymentUpdating, "deployment is updating")
		}

		err := r.ensureConnectionSecrets(ctx, projectService, akoServerless, atlasServerless.GetConnection())
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentConnectionSecretsNotCreated, err)
		}

		// Note: Serverless Private endpoints keep theirs flows without translation layer (yet)
		result := ensureServerlessPrivateEndpoints(ctx, akoServerless.GetProjectID(), akoServerless.GetCustomResource())

		switch {
		case result.IsInProgress():
			return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.ServerlessPrivateEndpointInProgress, result.GetMessage())
		case !result.IsOk():
			return r.terminate(ctx, workflow.ServerlessPrivateEndpointFailed, errors.New(result.GetMessage()))
		}

		err = customresource.ApplyLastConfigApplied(ctx.Context, akoServerless.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, akoServerless.GetCustomResource(), atlasServerless)

	case status.StateCREATING:
		return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.DeploymentUpdating, "deployment is updating")
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", atlasServerless.GetState()))
	}
}
