package atlasdeployment

import (
	"fmt"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func (r *AtlasDeploymentReconciler) handleFlexInstance(ctx *workflow.Context, projectService project.ProjectService, deploymentService deployment.AtlasDeploymentsService, deploymentInAKO, deploymentInAtlas *deployment.Flex) (ctrl.Result, error) {
	if deploymentInAtlas == nil {
		ctx.Log.Infof("Flex Instance %s doesn't exist in Atlas - creating", deploymentInAKO.GetName())
		newFlexDeployment, err := deploymentService.CreateDeployment(ctx.Context, deploymentInAKO)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		deploymentInAtlas = newFlexDeployment.(*deployment.Flex)
	}

	switch deploymentInAtlas.GetState() {
	case status.StateIDLE:
		if !reflect.DeepEqual(deploymentInAKO.FlexSpec, deploymentInAtlas.FlexSpec) {
			_, err := deploymentService.UpdateDeployment(ctx.Context, deploymentInAKO)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, deploymentInAKO.GetCustomResource(), deploymentInAtlas, workflow.DeploymentUpdating, "deployment is updating")
		}

		err := r.ensureConnectionSecrets(ctx, projectService, deploymentInAKO, deploymentInAtlas.GetConnection())
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentConnectionSecretsNotCreated, err)
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
	case status.StateDELETING, status.StateDELETED:
		return r.deleted()
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", deploymentInAtlas.GetState()))
	}
}
