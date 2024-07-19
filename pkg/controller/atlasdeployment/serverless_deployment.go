package atlasdeployment

import (
	"errors"
	"fmt"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDeploymentReconciler) handleServerlessInstance(ctx *workflow.Context, deploymentInAKO, deploymentInAtlas *deployment.Serverless) (ctrl.Result, error) {
	if deploymentInAtlas == nil {
		ctx.Log.Infof("Serverless Instance %s doesn't exist in Atlas - creating", deploymentInAKO.GetName())
		newServerlessDeployment, err := r.deploymentService.CreateDeployment(ctx.Context, deploymentInAKO)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		deploymentInAtlas = newServerlessDeployment.(*deployment.Serverless)
	}

	switch deploymentInAtlas.GetState() {
	case status.StateIDLE:
		if !reflect.DeepEqual(deploymentInAKO.ServerlessSpec, deploymentInAtlas.ServerlessSpec) {
			_, err := r.deploymentService.UpdateDeployment(ctx.Context, deploymentInAKO)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
		}

		err := r.ensureConnectionSecrets(ctx, deploymentInAKO, deploymentInAtlas.GetConnection())
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentConnectionSecretsNotCreated, err)
		}

		// Note: Serverless Private endpoints keep theirs flows without translation layer (yet)
		result := ensureServerlessPrivateEndpoints(ctx, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource())

		switch {
		case result.IsInProgress():
			return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.ServerlessPrivateEndpointInProgress, result.GetMessage())
		case !result.IsOk():
			return r.terminate(ctx, workflow.ServerlessPrivateEndpointFailed, errors.New(result.GetMessage()))
		}

		err = customresource.ApplyLastConfigApplied(ctx.Context, deploymentInAKO.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas)

	case status.StateCREATING:
		return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", deploymentInAtlas.GetState()))
	}
}
