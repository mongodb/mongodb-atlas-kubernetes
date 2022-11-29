package atlasproject

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	controllerruntime "sigs.k8s.io/controller-runtime"

	"go.mongodb.org/atlas/mongodbatlas"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

type TeamDataContainer struct {
	ProjectTeam *v1.Team
	Team        *v1.AtlasTeam
	Context     *workflow.Context
}

func (r *AtlasProjectReconciler) ensureAssignedTeams(ctx *workflow.Context, projectID string, project *v1.AtlasProject) workflow.Result {
	resourcesToWatch := make([]watch.WatchedObject, 0, len(project.Spec.Teams))
	defer func() {
		r.EnsureMultiplesResourcesAreWatched(
			types.NamespacedName{Namespace: project.Namespace, Name: project.Name},
			r.Log, resourcesToWatch...,
		)
		r.Log.Debugf("watching team resources: %v\r\n", r.WatchedResources)
	}()

	teamsToAssign := map[string]*v1.Team{}
	for _, assignedTeam := range project.Spec.Teams {
		assignedTeam := assignedTeam

		if assignedTeam.TeamRef.Name == "" {
			ctx.Log.Warnf("missing team name. skiping assignement for entry %v", assignedTeam)

			continue
		}

		if assignedTeam.TeamRef.Namespace == "" {
			assignedTeam.TeamRef.Namespace = project.Namespace
		}

		team := &v1.AtlasTeam{}
		teamReconciler := r.teamReconcile(team, ctx.Connection)
		_, err := teamReconciler(
			context.Background(),
			controllerruntime.Request{NamespacedName: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)
		if err != nil {
			ctx.Log.Warnf("unable to reconcile team %s. skipping assignment. %s", assignedTeam.TeamRef.GetObject(""), err.Error())
			continue
		}

		resourcesToWatch = append(
			resourcesToWatch,
			watch.WatchedObject{ResourceKind: team.Kind, Resource: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)

		teamsToAssign[team.Status.ID] = &assignedTeam
	}

	err := r.syncAssignedTeams(ctx, projectID, project, teamsToAssign)
	if err != nil {
		ctx.SetConditionFalse(status.ProjectTeamsReadyType)
		return workflow.Terminate(workflow.ProjectTeamUnavailable, err.Error())
	}

	ctx.SetConditionTrue(status.ProjectTeamsReadyType)

	if len(project.Spec.Teams) == 0 {
		ctx.UnsetCondition(status.ProjectTeamsReadyType)
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAssignedTeams(ctx *workflow.Context, projectID string, project *v1.AtlasProject, teamsToAssign map[string]*v1.Team) error {
	ctx.Log.Debug("fetching assigned teams from atlas")
	atlasAssignedTeams, _, err := ctx.Client.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
	if err != nil {
		return err
	}

	projectTeamStatus := status.ProjectTeamStatus{
		Teams:  make([]status.ProjectTeamRef, 0, len(teamsToAssign)),
		Status: true,
	}
	currentProjectsStatus := map[string]status.ProjectTeamRef{}
	for _, projectTeam := range project.Status.Teams.Teams {
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
			currentProjectsStatus[atlasAssignedTeam.TeamID] = status.ProjectTeamRef{
				ID:      atlasAssignedTeam.TeamID,
				TeamRef: desiredTeam.TeamRef,
			}
			delete(teamsToAssign, atlasAssignedTeam.TeamID)

			continue
		}

		ctx.Log.Debugf("removing team %s from project for later update", atlasAssignedTeam.TeamID)
		_, err = ctx.Client.Teams.RemoveTeamFromProject(context.Background(), projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}
	}

	for _, atlasAssignedTeam := range toDelete {
		ctx.Log.Debugf("removing team %s from project", atlasAssignedTeam.TeamID)
		_, err = ctx.Client.Teams.RemoveTeamFromProject(context.Background(), projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			ctx.Log.Warnf("failed to remove team %s from project: %s", atlasAssignedTeam.TeamID, err.Error())
		}

		teamRef := getTeamRefFromProjectStatus(project, atlasAssignedTeam.TeamID)
		if teamRef == nil {
			ctx.Log.Warnf("unable to find team %s status in the project: %s", atlasAssignedTeam.TeamID, err.Error())
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
		for teamID, assignedTeam := range teamsToAssign {
			projectTeams = append(projectTeams, assignedTeam.ToAtlas(teamID))
			currentProjectsStatus[teamID] = status.ProjectTeamRef{
				ID:      teamID,
				TeamRef: assignedTeam.TeamRef,
			}

			if err = r.updateTeamState(ctx, project, &assignedTeam.TeamRef, false); err != nil {
				ctx.Log.Warnf("failed to update team %s status with added project: %s", teamID, err.Error())
			}
		}

		_, _, err = ctx.Client.Projects.AddTeamsToProject(context.Background(), projectID, projectTeams)
		if err != nil {
			projectTeamStatus.Status = false
			projectTeamStatus.Error = err.Error()

			return err
		}
	}

	for _, projectsStatus := range currentProjectsStatus {
		projectTeamStatus.Teams = append(projectTeamStatus.Teams, projectsStatus)
	}

	ctx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(&projectTeamStatus))

	return nil
}

func (r *AtlasProjectReconciler) updateTeamState(ctx *workflow.Context, project *v1.AtlasProject, teamRef *common.ResourceRefNamespaced, isRemoval bool) error {
	team := &v1.AtlasTeam{}
	objKey := kube.ObjectKey(teamRef.Namespace, teamRef.Name)
	err := r.Client.Get(context.Background(), objKey, team)
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
		if projectStat.ID == project.Status.ID && isRemoval {
			continue
		}

		assignedProjects = append(assignedProjects, projectStat)
	}

	log := r.Log.With("atlasteam", teamRef)
	teamCtx, err := createTeamContextFromParent(team, r.Client, ctx.Connection, r.AtlasDomain, log)
	if err != nil {
		return err
	}

	if len(assignedProjects) == 0 {
		log.Debugf("team %s has no project associated to it. removing from atlas.", team.Spec.Name)
		_, err = teamCtx.Client.Teams.RemoveTeamFromOrganization(context.Background(), teamCtx.Connection.OrgID, team.Status.ID)
		if err != nil {
			return err
		}

		teamCtx.EnsureStatusOption(status.AtlasTeamUnsetID())
	}

	teamCtx.EnsureStatusOption(status.AtlasTeamSetProjects(assignedProjects))
	statushandler.Update(teamCtx, r.Client, r.EventRecorder, team)

	return nil
}

func getTeamRefFromProjectStatus(project *v1.AtlasProject, teamID string) *common.ResourceRefNamespaced {
	for _, stat := range project.Status.Teams.Teams {
		if stat.ID == teamID {
			return &stat.TeamRef
		}
	}

	return nil
}

func hasTeamRolesChanged(current []string, desired []v1.TeamRole) bool {
	desiredMap := map[string]struct{}{}
	for _, desiredRole := range desired {
		desiredMap[string(desiredRole)] = struct{}{}
	}

	for _, currentRole := range current {
		delete(desiredMap, currentRole)
	}

	return len(desiredMap) != 0
}
