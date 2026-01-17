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

package atlasproject

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

// handleProject creates the project if it doesn't exist yet. Returns the project ID
func (r *AtlasProjectReconciler) handleProject(ctx *workflow.Context, orgID string, atlasProject *akov2.AtlasProject, services *AtlasProjectServices) (ctrl.Result, error) {
	projectInAtlas, err := services.projectService.GetProjectByName(ctx.Context, atlasProject.Spec.Name)
	if err != nil && !errors.Is(err, translation.ErrNotFound) {
		return r.terminate(ctx, workflow.ProjectNotCreatedInAtlas, err)
	}

	wasDeleted := !atlasProject.GetDeletionTimestamp().IsZero()
	existInAtlas := projectInAtlas != nil

	switch {
	case !existInAtlas && !wasDeleted:
		return r.create(ctx, orgID, atlasProject, services.projectService)
	case existInAtlas && wasDeleted:
		return r.delete(ctx, services, orgID, atlasProject)
	case !existInAtlas && wasDeleted:
		return r.release(ctx, atlasProject)
	case existInAtlas && !wasDeleted && atlasProject.Status.ID == "":
		// short circuit the "manage" state,
		// there is no need to wait another reconcile cycle to continue.
		_, _ = r.manage(ctx, atlasProject, projectInAtlas.ID)
		atlasProject.Status.ID = projectInAtlas.ID
	}

	ctx.SetConditionTrue(api.ProjectReadyType)
	r.EventRecorder.Event(atlasProject, "Normal", string(api.ProjectReadyType), "")

	results := r.ensureProjectResources(ctx, atlasProject, services)
	for i := range results {
		if !results[i].IsOk() {
			logIfWarning(ctx, results[i])
			return results[i].ReconcileResult()
		}
	}

	err = customresource.ApplyLastConfigApplied(ctx.Context, atlasProject, r.Client)
	if err != nil {
		return r.terminate(ctx, workflow.Internal, err)
	}

	return r.ready(ctx, projectInAtlas.ID)
}

func (r *AtlasProjectReconciler) create(ctx *workflow.Context, orgID string, atlasProject *akov2.AtlasProject, projectService project.ProjectService) (ctrl.Result, error) {
	projectInAKO := project.NewProject(atlasProject, orgID)
	err := projectService.CreateProject(ctx.Context, projectInAKO)
	if err != nil {
		return r.terminate(ctx, workflow.ProjectNotCreatedInAtlas, err)
	}

	err = customresource.ApplyLastConfigApplied(ctx.Context, atlasProject, r.Client)
	if err != nil {
		return r.terminate(ctx, workflow.Internal, err)
	}

	return r.manage(ctx, atlasProject, projectInAKO.ID)
}

func (r *AtlasProjectReconciler) terminate(ctx *workflow.Context, errorCondition workflow.ConditionReason, err error) (ctrl.Result, error) {
	r.Log.Error(err)
	terminated := workflow.Terminate(errorCondition, err)
	ctx.SetConditionFromResult(api.ProjectReadyType, terminated)

	return terminated.ReconcileResult()
}

func (r *AtlasProjectReconciler) delete(ctx *workflow.Context, services *AtlasProjectServices, orgID string, atlasProject *akov2.AtlasProject) (ctrl.Result, error) {
	hasDeps, err := r.hasDependencies(ctx, atlasProject)
	if err != nil {
		return r.terminate(ctx, workflow.Internal, fmt.Errorf("failed to determine if project has dependencies: %w", err))
	}

	if hasDeps {
		return r.terminate(ctx, workflow.Internal, errors.New("the project cannot be deleted until dependencies were removed"))
	}

	if customresource.HaveFinalizer(atlasProject, customresource.FinalizerLabel) {
		if customresource.IsResourcePolicyKeepOrDefault(atlasProject, r.ObjectDeletionProtection) {
			r.Log.Info("Not removing Project from Atlas as per configuration")
		} else {
			if result := DeleteAllPrivateEndpoints(ctx, atlasProject); !result.IsOk() {
				return r.terminate(ctx, workflow.ServerlessPrivateEndpointReady, errors.New(result.GetMessage()))
			}
			if result := DeleteOwnedNetworkPeers(ctx.Context, atlasProject, ctx.SdkClientSet.SdkClient20250312012.NetworkPeeringApi, ctx.Log); !result.IsOk() {
				return r.terminate(ctx, workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New(result.GetMessage()))
			}

			err = r.syncAssignedTeams(ctx, services.teamsService, atlasProject.ID(), atlasProject, nil)
			if err != nil {
				ctx.SetConditionFalse(api.ProjectTeamsReadyType)
				return r.terminate(ctx, workflow.TeamNotCleaned, err)
			}

			if err = services.projectService.DeleteProject(ctx.Context, project.NewProject(atlasProject, orgID)); err != nil {
				return r.terminate(ctx, workflow.Internal, err)
			}
		}

		if err = customresource.ManageFinalizer(ctx.Context, r.Client, atlasProject, customresource.UnsetFinalizer); err != nil {
			return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
		}
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasProjectReconciler) ready(ctx *workflow.Context, projectID string) (ctrl.Result, error) {
	ctx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))
	result := workflow.OK()
	ctx.SetConditionFromResult(api.ProjectReadyType, result)
	ctx.SetConditionFromResult(api.ReadyType, result)

	return result.ReconcileResult()
}

func (r *AtlasProjectReconciler) release(ctx *workflow.Context, atlasProject *akov2.AtlasProject) (ctrl.Result, error) {
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasProject, customresource.UnsetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotRemoved, err)
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasProjectReconciler) manage(ctx *workflow.Context, atlasProject *akov2.AtlasProject, projectID string) (ctrl.Result, error) {
	r.Log.Debugw("Add deletion finalizer", "name", customresource.FinalizerLabel)
	if err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasProject, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	ctx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))
	result := workflow.InProgress(workflow.ProjectBeingConfiguredInAtlas, "configuring project in Atlas")
	ctx.SetConditionFromResult(api.ProjectReadyType, result)

	return result.ReconcileResult()
}

func (r *AtlasProjectReconciler) hasDependencies(ctx *workflow.Context, project *akov2.AtlasProject) (bool, error) {
	streamInstances := &akov2.AtlasStreamInstanceList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasStreamInstanceByProjectIndex,
			client.ObjectKeyFromObject(project).String(),
		),
	}
	err := r.Client.List(ctx.Context, streamInstances, listOps)
	if err != nil {
		return false, err
	}

	if len(streamInstances.Items) > 0 {
		return true, nil
	}

	customRoles := &akov2.AtlasCustomRoleList{}
	listOps = &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasCustomRoleByProject,
			client.ObjectKeyFromObject(project).String()),
	}
	err = r.Client.List(ctx.Context, customRoles, listOps)
	if err != nil {
		return false, err
	}

	if len(customRoles.Items) > 0 {
		return true, nil
	}

	privateEndpoints := &akov2.AtlasPrivateEndpointList{}
	listOps = &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasPrivateEndpointByProjectIndex,
			client.ObjectKeyFromObject(project).String(),
		),
	}
	err = r.Client.List(ctx.Context, privateEndpoints, listOps)
	if err != nil {
		return false, err
	}

	return len(privateEndpoints.Items) > 0, nil
}
