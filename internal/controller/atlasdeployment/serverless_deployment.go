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
)

func (r *AtlasDeploymentReconciler) handleServerlessInstance(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService,
	akoDeployment, atlasDeployment deployment.Deployment) (ctrl.Result, error) {
	akoServerless, ok := akoDeployment.(*deployment.Serverless)
	if !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in AKO is not a serverless cluster"))
	}

	var atlasServerless *deployment.Serverless
	if atlasServerless, ok = atlasDeployment.(*deployment.Serverless); atlasDeployment != nil && !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in Atlas is not a serverless cluster"))
	}

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

		// Note: Serverless Private endpoints keep theirs flows without translation layer (yet)
		result := ensureServerlessPrivateEndpoints(ctx, akoServerless.GetProjectID(), akoServerless.GetCustomResource())

		switch {
		case result.IsInProgress():
			return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.ServerlessPrivateEndpointInProgress, result.GetMessage())
		case !result.IsOk():
			return r.terminate(ctx, workflow.ServerlessPrivateEndpointFailed, errors.New(result.GetMessage()))
		}

		err := customresource.ApplyLastConfigApplied(ctx.Context, akoServerless.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, akoServerless, atlasServerless)
	case status.StateCREATING:
		return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, akoServerless.GetCustomResource(), atlasServerless, workflow.DeploymentUpdating, "deployment is updating")
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", atlasServerless.GetState()))
	}
}
