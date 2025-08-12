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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/searchindex"
)

const FreeTier = "M0"

func (r *AtlasDeploymentReconciler) handleAdvancedDeployment(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, akoDeployment, atlasDeployment deployment.Deployment) (ctrl.Result, error) {
	if akoDeployment.GetCustomResource().Spec.UpgradeToDedicated && !atlasDeployment.IsDedicated() {
		if atlasDeployment.GetState() == status.StateUPDATING {
			return r.inProgress(ctx, akoDeployment.GetCustomResource(), atlasDeployment, workflow.DeploymentUpdating, "deployment is updating")
		}

		updatedDeployment, err := deploymentService.UpgradeToDedicated(ctx.Context, atlasDeployment, akoDeployment)

		if err != nil {
			return r.terminate(ctx, workflow.DedicatedMigrationFailed, fmt.Errorf("failed to upgrade cluster: %w", err))
		}

		return r.inProgress(ctx, akoDeployment.GetCustomResource(), updatedDeployment, workflow.DedicatedMigrationProgressing, "Cluster upgrade to dedicated instance initiated in Atlas. The process may take several minutes")
	}

	akoCluster, ok := akoDeployment.(*deployment.Cluster)
	if !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in AKO is not an advanced cluster"))
	}

	var atlasCluster *deployment.Cluster
	if atlasCluster, ok = atlasDeployment.(*deployment.Cluster); atlasDeployment != nil && !ok {
		return r.terminate(ctx, workflow.Internal, errors.New("deployment in Atlas is not an advanced cluster"))
	}

	if atlasCluster == nil {
		ctx.Log.Infof("Advanced Deployment %s doesn't exist in Atlas - creating", akoCluster.GetName())
		newDeployment, err := deploymentService.CreateDeployment(ctx.Context, akoCluster)
		if err != nil {
			return r.terminate(ctx, workflow.DeploymentNotCreatedInAtlas, err)
		}

		atlasCluster = newDeployment.(*deployment.Cluster)
	}

	switch atlasCluster.GetState() {
	case status.StateIDLE:
		if changes, occurred := deployment.ComputeChanges(akoCluster, atlasCluster); occurred {
			updatedDeployment, err := deploymentService.UpdateDeployment(ctx.Context, changes)
			if err != nil {
				return r.terminate(ctx, workflow.DeploymentNotUpdatedInAtlas, err)
			}

			return r.inProgress(ctx, akoCluster.GetCustomResource(), updatedDeployment, workflow.DeploymentUpdating, "deployment is updating")
		}

		transition := r.ensureBackupScheduleAndPolicy(ctx, deploymentService, akoCluster.GetProjectID(), akoCluster.GetCustomResource(), atlasCluster.ZoneID)
		if transition != nil {
			return transition(workflow.Internal)
		}

		transition = r.ensureAdvancedOptions(ctx, deploymentService, akoCluster, atlasCluster)
		if transition != nil {
			return transition(workflow.DeploymentAdvancedOptionsReady)
		}

		var results []workflow.DeprecatedResult
		if !r.AtlasProvider.IsCloudGov() {
			searchNodeResult := handleSearchNodes(ctx, akoCluster.GetCustomResource(), akoCluster.GetProjectID())
			results = append(results, searchNodeResult)
		}

		searchService := searchindex.NewSearchIndexes(ctx.SdkClientSet.SdkClient20250312002.AtlasSearchApi)
		result := handleSearchIndexes(ctx, r.Client, searchService, akoCluster.GetCustomResource(), akoCluster.GetProjectID())
		results = append(results, result)

		result = r.ensureCustomZoneMapping(
			ctx,
			deploymentService,
			akoCluster.GetProjectID(),
			akoCluster.GetCustomResource().Spec.DeploymentSpec.CustomZoneMapping,
			akoCluster.GetName(),
		)
		results = append(results, result)

		result = r.ensureManagedNamespaces(
			ctx,
			deploymentService,
			akoCluster.GetProjectID(),
			akoCluster.ClusterType,
			akoCluster.GetCustomResource().Spec.DeploymentSpec.ManagedNamespaces,
			akoCluster.GetName(),
		)
		results = append(results, result)

		for i := range results {
			if !results[i].IsOk() {
				return r.transitionFromResult(ctx, deploymentService, akoCluster.GetProjectID(), akoCluster.GetCustomResource(), results[i])(workflow.Internal)
			}
		}
		err := customresource.ApplyLastConfigApplied(ctx.Context, akoCluster.GetCustomResource(), r.Client)
		if err != nil {
			return r.terminate(ctx, workflow.Internal, err)
		}

		return r.ready(ctx, akoCluster, atlasCluster)
	case status.StateCREATING:
		return r.inProgress(ctx, akoCluster.GetCustomResource(), atlasCluster, workflow.DeploymentCreating, "deployment is provisioning")
	case status.StateUPDATING, status.StateREPAIRING:
		return r.inProgress(ctx, akoCluster.GetCustomResource(), atlasCluster, workflow.DeploymentUpdating, "deployment is updating")
	case status.StateDELETING, status.StateDELETED:
		return r.handleDeleted()
	default:
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("unknown deployment state: %s", atlasCluster.GetState()))
	}
}

func (r *AtlasDeploymentReconciler) ensureAdvancedOptions(ctx *workflow.Context, deploymentService deployment.AtlasDeploymentsService, deploymentInAKO, deploymentInAtlas *deployment.Cluster) transitionFn {
	if deploymentInAKO.IsTenant() {
		return nil
	}

	err := deploymentService.ClusterWithProcessArgs(ctx.Context, deploymentInAtlas)
	if err != nil {
		return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), err)
	}

	if deploymentInAKO.ProcessArgs != nil {
		if deploymentInAKO.ProcessArgs.DefaultReadConcern != "" {
			ctx.Log.Warn("Process Arg DefaultReadConcern is no longer available in Atlas. Setting this will have no effect.")
		}
		if deploymentInAKO.ProcessArgs.FailIndexKeyTooLong != nil {
			ctx.Log.Warn("Process Arg FailIndexKeyTooLong is no longer available in Atlas. Setting this will have no effect.")
		}
		if !reflect.DeepEqual(deploymentInAKO.ProcessArgs, deploymentInAtlas.ProcessArgs) {
			err = deploymentService.UpdateProcessArgs(ctx.Context, deploymentInAKO)
			if err != nil {
				return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), err)
			}

			return r.transitionFromLegacy(ctx, deploymentService, deploymentInAKO.GetProjectID(), deploymentInAKO.GetCustomResource(), nil)
		}
	}

	return nil
}
