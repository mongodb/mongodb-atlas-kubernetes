// Copyright 2024 MongoDB Inc
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
	"net/http"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

// convention regarding the state machine:
// handle... prefixed methods are state handling methods,
// any other method is a state transition method.
type searchNodeController struct {
	ctx        *workflow.Context
	deployment *akov2.AtlasDeployment
	projectID  string
}

func handleSearchNodes(ctx *workflow.Context, deployment *akov2.AtlasDeployment, projectID string) workflow.DeprecatedResult {
	ctx.Log.Debug("starting search node processing")
	defer ctx.Log.Debug("finished search node processing")

	s := searchNodeController{
		ctx:        ctx,
		deployment: deployment,
		projectID:  projectID,
	}

	// determine the current controller state and invoke state handling methods
	c, ok := ctx.GetCondition(api.SearchNodesReadyType)
	if ok {
		switch reason := workflow.ConditionReason(c.Reason); reason {
		case workflow.SearchNodesCreating, workflow.SearchNodesUpdating:
			return s.handleUpserting(reason)
		case workflow.SearchNodesDeleting:
			return s.handleDeleting()
		}
	}
	// anything else we assume we are in pending state
	return s.handlePending()
}

// handlePending handles the initial pending state. The following states are handled here:
//
// 1. Pending Unmanaged: search nodes are not created but deployment is ready:
//
//	Conditions:
//	  Status:                True
//	  Type:                  DeploymentReady
//
// 2. Pending Idle: search nodes are created and idle:
//
//	Conditions:
//		Status:                True
//		Type:                  SearchNodesReady
//
// 3. Pending Terminated: some error occurred in other states:
//
//	Conditions:
//		Status:                False
//		Type:                  SearchNodesReady
//		Message:               ErrorSearchNodes... | InternalError
//
// The following transitions can happen here:
// - create: transitions to "creating" state when search nodes are in AKO but not in Atlas.
// - update: transitions to "updating" state when search nodes are in AKO and in Atlas.
// - delete: transitions to "deleting" state when no search nodes are in AKO but in Atlas.
// - unmanage: transitions stays in "pending unmanaged" state unsetting "idle" status when search nodes are neither in AKO nor in Atlas.
func (s *searchNodeController) handlePending() workflow.DeprecatedResult {
	atlasNodes, found, err := s.getAtlasSearchDeployment()
	if err != nil {
		// transition back to pending, something went wrong.
		return s.terminate(workflow.Internal, err)
	}

	akoEmpty := len(s.deployment.Spec.DeploymentSpec.SearchNodes) == 0
	atlasEmpty := !found

	switch {
	case !akoEmpty && atlasEmpty:
		// If no nodes configured in atlas, but some in the operator - create them.
		return s.create()
	case !akoEmpty && !atlasEmpty:
		// If nodes already configured in atlas and in the operator - update them.
		return s.update(atlasNodes)
	case akoEmpty && !atlasEmpty:
		// If no nodes configured in the operator, but some in atlas - delete them.
		return s.delete()
	default:
		// akoEmpty && atlasEmpty
		// If no search nodes are in AKO and no nodes in Atlas - unmanage.
		return s.unmanage()
	}
}

// handleUpserting handles creating or updating search nodes. The following state is handled here:
//
//	Conditions:
//	  Status:                False
//	  Type:                  SearchNodesReady
//	  Reason:                SearchNodesCreating | SearchNodesUpdating
//
// The following transitions can happen here:
// - terminate: when an error occurred or when the operation has been aborted.
// - idle: when search nodes are marked as IDLE in Atlas.
func (s *searchNodeController) handleUpserting(state workflow.ConditionReason) workflow.DeprecatedResult {
	if len(s.deployment.Spec.DeploymentSpec.SearchNodes) == 0 {
		return s.terminate(workflow.ErrorSearchNodesOperationAborted, errors.New("aborting update/create: no search nodes specified"))
	}

	atlasNodes, found, err := s.getAtlasSearchDeployment()
	if err != nil {
		return s.terminate(workflow.ErrorSearchNodesNotUpsertedInAtlas, err)
	}
	if !found {
		return s.terminate(workflow.ErrorSearchNodesNotUpsertedInAtlas, errors.New("no search nodes found in Atlas"))
	}

	hasChanged := !reflect.DeepEqual(s.deployment.Spec.DeploymentSpec.SearchNodesToAtlas(), atlasNodes.GetSpecs())
	switch {
	case hasChanged:
		return s.terminate(workflow.ErrorSearchNodesOperationAborted, errors.New("aborting update/create: spec has changed"))
	case atlasNodes.GetStateName() != "IDLE":
		return s.progress(
			state,
			fmt.Sprintf("search nodes are not ready yet, Atlas state: %q", atlasNodes.GetStateName()),
			"waiting for search nodes to become ready",
		)
	default:
		return s.idle()
	}
}

// handleDeleting handles deleting search nodes. The following state is handled here:
//
//	Conditions:
//	  Status:                False
//	  Type:                  SearchNodesReady
//	  Reason:                SearchNodesDeleting
//
// The following transitions can happen here:
// - terminate: when an error occurred or when the operation has been aborted.
// - unmanage: when there are no search nodes anymore in Atlas.
func (s *searchNodeController) handleDeleting() workflow.DeprecatedResult {
	if len(s.deployment.Spec.DeploymentSpec.SearchNodes) > 0 {
		return s.terminate(workflow.ErrorSearchNodesOperationAborted, errors.New("aborting deletion: search nodes are specified"))
	}

	atlasNodes, found, err := s.getAtlasSearchDeployment()
	switch {
	case err != nil:
		return s.terminate(workflow.ErrorSearchNodesNotUpsertedInAtlas, err)
	case found:
		return s.progress(
			workflow.SearchNodesDeleting,
			fmt.Sprintf("search nodes are being deleted, Atlas state: %q", atlasNodes.GetStateName()),
			"deleting search nodes",
		)
	default:
		return s.unmanage()
	}
}

// create executes the actual creation of search nodes in Atlas and transitions to the following states:
// - creating: after search nodes have been created in Atlas
// - terminated: when an error occurred.
func (s *searchNodeController) create() workflow.DeprecatedResult {
	s.ctx.Log.Debugf("creating search nodes %v", s.deployment.Spec.DeploymentSpec.SearchNodes)
	resp, _, err := s.ctx.SdkClientSet.SdkClient20250312011.AtlasSearchApi.CreateClusterSearchDeployment(s.ctx.Context, s.projectID, s.deployment.GetDeploymentName(), &admin.ApiSearchDeploymentRequest{
		Specs: s.deployment.Spec.DeploymentSpec.SearchNodesToAtlas(),
	}).Execute()
	if err != nil {
		return s.terminate(workflow.ErrorSearchNodesNotUpsertedInAtlas, err)
	}

	return s.progress(
		workflow.SearchNodesCreating,
		fmt.Sprintf("search nodes are not ready yet: Atlas state is %q", resp.GetStateName()),
		"creating search nodes",
	)
}

// update updates search nodes in Atlas, if necessary. It transitions to the following states:
// - updating: after search nodes have been updated in Atlas.
// - idle: if no changes are necessary.
// - terminated: when an error occurred.
func (s *searchNodeController) update(atlasNodes *admin.ApiSearchDeploymentResponse) workflow.DeprecatedResult {
	s.ctx.Log.Debugf("updating search nodes %v", s.deployment.Spec.DeploymentSpec.SearchNodes)
	currentAkoNodesAsAtlas := s.deployment.Spec.DeploymentSpec.SearchNodesToAtlas()
	// We can deepequal without normalization here because there is only ever 1 spec in the array
	if !reflect.DeepEqual(currentAkoNodesAsAtlas, atlasNodes.GetSpecs()) {
		updateResponse, _, err := s.ctx.SdkClientSet.SdkClient20250312011.AtlasSearchApi.UpdateClusterSearchDeployment(
			s.ctx.Context, s.projectID, s.deployment.GetDeploymentName(), &admin.ApiSearchDeploymentRequest{
				Specs: s.deployment.Spec.DeploymentSpec.SearchNodesToAtlas(),
			}).Execute()
		if err != nil {
			return s.terminate(workflow.ErrorSearchNodesNotUpsertedInAtlas, err)
		}

		return s.progress(
			workflow.SearchNodesUpdating,
			fmt.Sprintf("search nodes are not ready yet: Atlas state is %q", updateResponse.GetStateName()),
			"updating search nodes",
		)
	}
	s.ctx.Log.Debug("search nodes in AKO and Atlas are equal")

	// even if both Atlas and AKO are equal, check if existing atlas nodes are potentially are not IDLE (i.e. UPDATING).
	// In this case we should continue staying in "updating" state.
	if atlasNodes.GetStateName() != "IDLE" {
		return s.progress(
			workflow.SearchNodesUpdating,
			fmt.Sprintf("search nodes are not ready yet: Atlas state is %q", atlasNodes.GetStateName()),
			"updating search nodes",
		)
	}

	// no changes needed, continue transitioning to pending state, marking search node status as idle.
	return s.idle()
}

// delete deletes search nodes in Atlas. It transitions to the following states:
// - deleting: after search nodes have been deleted in Atlas.
// - terminated: when an error occurred.
func (s *searchNodeController) delete() workflow.DeprecatedResult {
	s.ctx.Log.Debug("deleting search nodes")
	_, err := s.ctx.SdkClientSet.SdkClient20250312011.AtlasSearchApi.DeleteClusterSearchDeployment(s.ctx.Context, s.projectID, s.deployment.GetDeploymentName()).Execute()
	if err != nil {
		return s.terminate(workflow.ErrorSearchNodesNotDeletedInAtlas, err)
	}

	return s.progress(
		workflow.SearchNodesDeleting,
		"deleting search nodes",
		"deleting search nodes",
	)
}

// progress transitions to the given state in the "SearchNodesReady" status and sets the given fineMsg as its message.
// further it returns a coarse grained progress state for bubbling the chain.
func (s *searchNodeController) progress(state workflow.ConditionReason, fineMsg, coarseMsg string) workflow.DeprecatedResult {
	var (
		fineProgress   = workflow.InProgress(state, fineMsg)
		coarseProgress = workflow.InProgress(state, coarseMsg)
	)

	s.ctx.SetConditionFromResult(api.SearchNodesReadyType, fineProgress)
	return coarseProgress
}

// terminate transitions to pending state if an error occurred.
func (s *searchNodeController) terminate(reason workflow.ConditionReason, err error) workflow.DeprecatedResult {
	s.ctx.Log.Error(err)
	result := workflow.Terminate(reason, err)
	s.ctx.SetConditionFromResult(api.SearchNodesReadyType, result)
	return result
}

// unmanage transitions to pending state if no search nodes are managed.
func (s *searchNodeController) unmanage() workflow.DeprecatedResult {
	s.ctx.UnsetCondition(api.SearchNodesReadyType)
	return workflow.OK()
}

// idle transitions to idle state search nodes that are ready and idle.
func (s *searchNodeController) idle() workflow.DeprecatedResult {
	s.ctx.SetConditionTrue(api.SearchNodesReadyType)
	return workflow.OK()
}

func (s *searchNodeController) getAtlasSearchDeployment() (*admin.ApiSearchDeploymentResponse, bool, error) {
	atlasNodes, _, err := s.ctx.SdkClientSet.SdkClient20250312011.AtlasSearchApi.GetClusterSearchDeployment(s.ctx.Context, s.projectID, s.deployment.GetDeploymentName()).Execute()
	if err != nil {
		apiError, ok := admin.AsError(err)
		// TODO: Currently 400, should be be 404: CLOUDP-239015
		if ok && (apiError.GetError() == http.StatusBadRequest || apiError.GetError() == http.StatusNotFound) {
			s.ctx.Log.Debug("no search nodes in atlas found")
			return nil, false, nil
		} else {
			return nil, false, err
		}
	}

	if atlasNodes == nil || len(atlasNodes.GetSpecs()) == 0 {
		return nil, false, nil
	}

	return atlasNodes, true, nil
}
