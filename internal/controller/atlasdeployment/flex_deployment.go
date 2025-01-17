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

func (r *AtlasDeploymentReconciler) handleFlexInstance(ctx *workflow.Context, projectService project.ProjectService, deploymentService deployment.AtlasDeploymentsService, akoDeployment, atlasDeployment deployment.Deployment) (ctrl.Result, error) {
	akoFlex, ok := akoDeployment.(*deployment.Flex)
	if !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in AKO is not a serverless cluster"))
	}
	atlasFlex, _ := atlasDeployment.(*deployment.Flex)

	if atlasFlex == nil {
		ctx.Log.Infof("Flex Instance %s doesn't exist in Atlas - creating", akoFlex.GetName())
		newFlexDeployment, err := deploymentService.CreateDeployment(ctx.Context, akoFlex)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		atlasFlex = newFlexDeployment.(*deployment.Flex)
	}

	switch atlasFlex.GetState() {
	case status.StateIDLE:
		if !reflect.DeepEqual(akoFlex.FlexSpec, atlasFlex.FlexSpec) {
			_, err := deploymentService.UpdateDeployment(ctx.Context, akoFlex)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, akoFlex.GetCustomResource(), atlasFlex, workflow.DeploymentUpdating, "deployment is updating")
		}

		err := r.ensureConnectionSecrets(ctx, projectService, akoFlex, atlasFlex.GetConnection())
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentConnectionSecretsNotCreated, err)
		}

		err = customresource.ApplyLastConfigApplied(ctx.Context, akoFlex.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, akoFlex.GetCustomResource(), atlasFlex)

	case status.StateCREATING:
		return r.inProgress(ctx, akoFlex.GetCustomResource(), atlasFlex, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, akoFlex.GetCustomResource(), atlasFlex, workflow.DeploymentUpdating, "deployment is updating")
	case status.StateDELETING, status.StateDELETED:
		return r.handleDeleted()
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", atlasFlex.GetState()))
	}
}
