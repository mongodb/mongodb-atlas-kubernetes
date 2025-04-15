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

	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
)

type TeamDataContainer struct {
	ProjectTeam *akov2.Team
	Team        *akov2.AtlasTeam
	Context     *workflow.Context
}

func (r *AtlasProjectReconciler) ensureAssignedTeams(workflowCtx *workflow.Context, teamsService teams.TeamsService, project *akov2.AtlasProject) workflow.Result {
	teamsToAssign := map[string]*akov2.Team{}
	for _, entry := range project.Spec.Teams {
		assignedTeam := entry

		if assignedTeam.TeamRef.Name == "" {
			workflowCtx.Log.Warnf("missing team name. skipping assignment for entry %v", assignedTeam)
			continue
		}

		if assignedTeam.TeamRef.Namespace == "" {
			assignedTeam.TeamRef.Namespace = project.Namespace
		}

		team := &akov2.AtlasTeam{}
		teamReconciler := r.teamReconcile(team, workflowCtx, teamsService)
		_, err := teamReconciler(
			workflowCtx.Context,
			controllerruntime.Request{NamespacedName: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)
		if err != nil {
			workflowCtx.Log.Warnf("unable to reconcile team %s. skipping assignment. %s", assignedTeam.TeamRef.GetObject(""), err.Error())
			continue
		}

		teamsToAssign[team.Status.ID] = &assignedTeam
	}

	err := r.syncAssignedTeams(workflowCtx, teamsService, project.ID(), project, teamsToAssign)
	if err != nil {
		workflowCtx.SetConditionFalse(api.ProjectTeamsReadyType)
		return workflow.Terminate(workflow.ProjectTeamUnavailable, err)
	}

	workflowCtx.SetConditionTrue(api.ProjectTeamsReadyType)

	if len(project.Spec.Teams) == 0 {
		workflowCtx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(nil))
		workflowCtx.UnsetCondition(api.ProjectTeamsReadyType)
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAssignedTeams(ctx *workflow.Context, teamsService teams.TeamsService, projectID string, project *akov2.AtlasProject, teamsToAssign map[string]*akov2.Team) error {
	ctx.Log.Debug("fetching assigned teams from atlas")

	atlasAssignedTeams, err := teamsService.ListProjectTeams(ctx.Context, projectID)
	if err != nil {
		return err
	}

	projectTeamStatus := make([]status.ProjectTeamStatus, 0, len(teamsToAssign))
	currentProjectsStatus := map[string]status.ProjectTeamStatus{}
	for _, projectTeam := range project.Status.Teams {
		currentProjectsStatus[projectTeam.ID] = projectTeam
	}

	defer statushandler.Update(ctx, r.Client, r.EventRecorder, project)
	var teamErrors error

	toDelete := make([]*teams.AssignedTeam, 0, len(atlasAssignedTeams))

	for _, atlasAssignedTeam := range atlasAssignedTeams {
		if atlasAssignedTeam.TeamID == "" {
			continue
		}

		desiredTeam, ok := teamsToAssign[atlasAssignedTeam.TeamID]
		if !ok {
			toDelete = append(toDelete, &atlasAssignedTeam)
			continue
		}

		if !hasTeamRolesChanged(atlasAssignedTeam.Roles, desiredTeam.Roles) {
			currentProjectsStatus[atlasAssignedTeam.TeamID] = status.ProjectTeamStatus{
				ID:      atlasAssignedTeam.TeamID,
				TeamRef: desiredTeam.TeamRef,
			}
			delete(teamsToAssign, atlasAssignedTeam.TeamID)

			continue
		}

		ctx.Log.Debugf("removing team %s from project for later update", atlasAssignedTeam.TeamID)
		err = teamsService.Unassign(ctx.Context, projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			teamErrors = errors.Join(teamErrors, err)
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}
	}

	for _, atlasAssignedTeam := range toDelete {
		ctx.Log.Debugf("removing team %s from project", atlasAssignedTeam.TeamID)
		err = teamsService.Unassign(ctx.Context, projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			teamErrors = errors.Join(teamErrors, err)
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}

		teamRef := getTeamRefFromProjectStatus(project, atlasAssignedTeam.TeamID)
		if teamRef == nil {
			ctx.Log.Warnf("unable to find team %s status in the project", atlasAssignedTeam.TeamID)
		} else {
			if err = r.updateTeamState(ctx, project, teamRef, true); err != nil {
				teamErrors = errors.Join(teamErrors, err)
				ctx.Log.Warnf("failed to update team %s status with removed project: %s", atlasAssignedTeam.TeamID, err.Error())
			}
		}

		delete(currentProjectsStatus, atlasAssignedTeam.TeamID)
	}

	if len(teamsToAssign) > 0 {
		ctx.Log.Debug("assigning teams to project")
		projectTeams := make([]teams.AssignedTeam, 0, len(teamsToAssign))
		for teamID := range teamsToAssign {
			teamToAssign := teams.NewAssignedTeam(teamsToAssign[teamID], teamID)

			projectTeams = append(projectTeams, *teamToAssign)
			currentProjectsStatus[teamID] = status.ProjectTeamStatus{
				ID:      teamID,
				TeamRef: teamsToAssign[teamID].TeamRef,
			}

			if err = r.updateTeamState(ctx, project, &teamsToAssign[teamID].TeamRef, false); err != nil {
				teamErrors = errors.Join(teamErrors, err)
				ctx.Log.Warnf("failed to update team %s status with added project: %s", teamID, err.Error())
			}
		}

		err = teamsService.Assign(ctx.Context, &projectTeams, projectID)
		if err != nil {
			return errors.Join(teamErrors, err)
		}
	}

	for _, projectsStatus := range currentProjectsStatus {
		projectTeamStatus = append(projectTeamStatus, projectsStatus)
	}

	ctx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(&projectTeamStatus))

	return teamErrors
}

func (r *AtlasProjectReconciler) updateTeamState(ctx *workflow.Context, project *akov2.AtlasProject, teamRef *common.ResourceRefNamespaced, isRemoval bool) error {
	team := &akov2.AtlasTeam{}
	objKey := kube.ObjectKey(teamRef.Namespace, teamRef.Name)
	err := r.Client.Get(ctx.Context, objKey, team)
	if err != nil {
		return err
	}

	assignedProjects := make([]status.TeamProject, 0, len(team.Status.Projects)+1)

	if !isRemoval {
		assignedProjects = append(
			assignedProjects,
			status.TeamProject{
				ID:   project.Status.ID,
				Name: project.Spec.Name,
			},
		)
	}

	for _, projectStat := range team.Status.Projects {
		if projectStat.ID == project.Status.ID {
			continue
		}

		assignedProjects = append(assignedProjects, projectStat)
	}

	log := r.Log.With("atlasteam", teamRef)
	conditions := akov2.InitCondition(team, api.FalseCondition(api.ReadyType))
	teamCtx := workflow.NewContext(log, conditions, ctx.Context, team)

	if len(assignedProjects) == 0 {
		if err = customresource.ManageFinalizer(ctx.Context, r.Client, team, customresource.UnsetFinalizer); err != nil {
			return err
		}

		teamCtx.SetConditionTrueMsg(api.TeamUnmanaged, "This resource is only reconciled when associated to a project")
	} else {
		if err = customresource.ManageFinalizer(ctx.Context, r.Client, team, customresource.SetFinalizer); err != nil {
			return err
		}
	}

	teamCtx.SetConditionTrue(api.ReadyType)
	teamCtx.EnsureStatusOption(status.AtlasTeamSetProjects(assignedProjects))
	statushandler.Update(teamCtx, r.Client, r.EventRecorder, team)

	return nil
}

func getTeamRefFromProjectStatus(project *akov2.AtlasProject, teamID string) *common.ResourceRefNamespaced {
	for _, stat := range project.Status.Teams {
		if stat.ID == teamID {
			return &stat.TeamRef
		}
	}

	return nil
}

func hasTeamRolesChanged(current []string, desired []akov2.TeamRole) bool {
	desiredMap := map[string]struct{}{}
	for _, desiredRole := range desired {
		desiredMap[string(desiredRole)] = struct{}{}
	}

	for _, currentRole := range current {
		delete(desiredMap, currentRole)
	}

	return len(desiredMap) != 0
}
