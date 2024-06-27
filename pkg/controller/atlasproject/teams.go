package atlasproject

import (
	"go.mongodb.org/atlas/mongodbatlas"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

type TeamDataContainer struct {
	ProjectTeam *akov2.Team
	Team        *akov2.AtlasTeam
	Context     *workflow.Context
}

func (r *AtlasProjectReconciler) ensureAssignedTeams(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	resourcesToWatch := make([]watch.WatchedObject, 0, len(project.Spec.Teams))

	defer func() {
		workflowCtx.AddResourcesToWatch(resourcesToWatch...)
		r.Log.Debugf("watching team resources: %v\r\n", r.DeprecatedResourceWatcher.WatchedResourcesSnapshot())
	}()

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
		teamReconciler := r.teamReconcile(team, project.ConnectionSecretObjectKey())
		_, err := teamReconciler(
			workflowCtx.Context,
			controllerruntime.Request{NamespacedName: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)
		if err != nil {
			workflowCtx.Log.Warnf("unable to reconcile team %s. skipping assignment. %s", assignedTeam.TeamRef.GetObject(""), err.Error())
			continue
		}

		resourcesToWatch = append(
			resourcesToWatch,
			watch.WatchedObject{ResourceKind: team.Kind, Resource: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)

		teamsToAssign[team.Status.ID] = &assignedTeam
	}

	err := r.syncAssignedTeams(workflowCtx, project.ID(), project, teamsToAssign)
	if err != nil {
		workflowCtx.SetConditionFalse(api.ProjectTeamsReadyType)
		return workflow.Terminate(workflow.ProjectTeamUnavailable, err.Error())
	}

	workflowCtx.SetConditionTrue(api.ProjectTeamsReadyType)

	if len(project.Spec.Teams) == 0 {
		workflowCtx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(nil))
		workflowCtx.UnsetCondition(api.ProjectTeamsReadyType)
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAssignedTeams(ctx *workflow.Context, projectID string, project *akov2.AtlasProject, teamsToAssign map[string]*akov2.Team) error {
	ctx.Log.Debug("fetching assigned teams from atlas")
	atlasAssignedTeams, _, err := ctx.Client.Projects.GetProjectTeamsAssigned(ctx.Context, projectID)
	if err != nil {
		return err
	}

	projectTeamStatus := make([]status.ProjectTeamStatus, 0, len(teamsToAssign))
	currentProjectsStatus := map[string]status.ProjectTeamStatus{}
	for _, projectTeam := range project.Status.Teams {
		currentProjectsStatus[projectTeam.ID] = projectTeam
	}

	defer statushandler.Update(ctx, r.Client, r.EventRecorder, project)

	toDelete := make([]*mongodbatlas.Result, 0, len(atlasAssignedTeams.Results))
	for _, atlasAssignedTeam := range atlasAssignedTeams.Results {
		desiredTeam, ok := teamsToAssign[atlasAssignedTeam.TeamID]
		if !ok {
			toDelete = append(toDelete, atlasAssignedTeam)

			continue
		}

		if !hasTeamRolesChanged(atlasAssignedTeam.RoleNames, desiredTeam.Roles) {
			currentProjectsStatus[atlasAssignedTeam.TeamID] = status.ProjectTeamStatus{
				ID:      atlasAssignedTeam.TeamID,
				TeamRef: desiredTeam.TeamRef,
			}
			delete(teamsToAssign, atlasAssignedTeam.TeamID)

			continue
		}

		ctx.Log.Debugf("removing team %s from project for later update", atlasAssignedTeam.TeamID)
		_, err = ctx.Client.Teams.RemoveTeamFromProject(ctx.Context, projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}
	}

	for _, atlasAssignedTeam := range toDelete {
		ctx.Log.Debugf("removing team %s from project", atlasAssignedTeam.TeamID)
		_, err = ctx.Client.Teams.RemoveTeamFromProject(ctx.Context, projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}

		teamRef := getTeamRefFromProjectStatus(project, atlasAssignedTeam.TeamID)
		if teamRef == nil {
			ctx.Log.Warnf("unable to find team %s status in the project", atlasAssignedTeam.TeamID)
		} else {
			if err = r.updateTeamState(ctx, project, teamRef, true); err != nil {
				ctx.Log.Warnf("failed to update team %s status with removed project: %s", atlasAssignedTeam.TeamID, err.Error())
			}
		}

		delete(currentProjectsStatus, atlasAssignedTeam.TeamID)
	}

	if len(teamsToAssign) > 0 {
		ctx.Log.Debug("assigning teams to project")
		projectTeams := make([]*mongodbatlas.ProjectTeam, 0, len(teamsToAssign))
		for teamID := range teamsToAssign {
			assignedTeam := teamsToAssign[teamID]
			projectTeams = append(projectTeams, assignedTeam.ToAtlas(teamID))
			currentProjectsStatus[teamID] = status.ProjectTeamStatus{
				ID:      teamID,
				TeamRef: assignedTeam.TeamRef,
			}

			if err = r.updateTeamState(ctx, project, &assignedTeam.TeamRef, false); err != nil {
				ctx.Log.Warnf("failed to update team %s status with added project: %s", teamID, err.Error())
			}
		}

		_, _, err = ctx.Client.Projects.AddTeamsToProject(ctx.Context, projectID, projectTeams)
		if err != nil {
			return err
		}
	}

	for _, projectsStatus := range currentProjectsStatus {
		projectTeamStatus = append(projectTeamStatus, projectsStatus)
	}

	ctx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(&projectTeamStatus))

	return nil
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
	teamCtx := workflow.NewContext(log, conditions, ctx.Context)

	atlasClient, orgID, err := r.AtlasProvider.Client(teamCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		return err
	}
	teamCtx.Client = atlasClient

	if len(assignedProjects) == 0 {
		log.Debugf("team %s has no project associated to it. removing from atlas.", team.Spec.Name)
		_, err = teamCtx.Client.Teams.RemoveTeamFromOrganization(ctx.Context, orgID, team.Status.ID)
		if err != nil {
			return err
		}

		teamCtx.EnsureStatusOption(status.AtlasTeamUnsetID())
	}

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
