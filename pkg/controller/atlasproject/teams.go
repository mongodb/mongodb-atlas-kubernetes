package atlasproject

import (
	"context"

	controllerruntime "sigs.k8s.io/controller-runtime"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.mongodb.org/atlas/mongodbatlas"
	"k8s.io/apimachinery/pkg/types"
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

	teamsToAssign := map[string]TeamDataContainer{}
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
		teamReconciler := teamReconcile(team, r.Client, r.EventRecorder, ctx.Connection, r.AtlasDomain, r.Log)
		_, err := teamReconciler(
			context.Background(),
			controllerruntime.Request{NamespacedName: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)
		if err != nil {
			ctx.Log.Warnf("unable to reconcile team %s. skipping assignment. %s", assignedTeam.TeamRef.GetObject(project.Namespace), err.Error())
			continue
		}

		resourcesToWatch = append(
			resourcesToWatch,
			watch.WatchedObject{ResourceKind: team.Kind, Resource: types.NamespacedName{Name: assignedTeam.TeamRef.Name, Namespace: assignedTeam.TeamRef.Namespace}},
		)

		teamsToAssign[team.Status.ID] = TeamDataContainer{
			ProjectTeam: &assignedTeam,
			Team:        team,
			Context:     nil,
		}
	}

	err := r.syncAssignedTeams(ctx, projectID, teamsToAssign)
	if err != nil {
		ctx.SetConditionFalse(status.ProjectTeamsReadyType)
		return workflow.Terminate(workflow.ProjectTeamUnavailable, err.Error())
	}

	ctx.SetConditionTrue(status.ProjectTeamsReadyType)
	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAssignedTeams(ctx *workflow.Context, projectID string, teamsToAssign map[string]TeamDataContainer) error {
	ctx.Log.Debug("fetching assigned teams from atlas")
	atlasAssignedTeams, _, err := ctx.Client.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
	if err != nil {
		return err
	}

	for _, atlasAssignedTeam := range atlasAssignedTeams.Results {
		_, ok := teamsToAssign[atlasAssignedTeam.TeamID]
		if ok {
			delete(teamsToAssign, atlasAssignedTeam.TeamID)

			continue
		}

		ctx.Log.Debugf("removing team %s from project", atlasAssignedTeam.TeamID)
		_, err = ctx.Client.Teams.RemoveTeamFromProject(context.Background(), projectID, atlasAssignedTeam.TeamID)
		if err != nil {
			ctx.Log.Warnf("failed to remove team %s from project", atlasAssignedTeam.TeamID)
		}

		//teamData.Context.EnsureStatusOption(status.AtlasTeamRemoveProject(projectID, ""))
		//statushandler.Update(teamData.Context, r.Client, r.EventRecorder, teamData.Team)
	}

	if len(teamsToAssign) == 0 {
		return nil
	}

	ctx.Log.Debug("assigning teams to project")
	projectTeams := make([]*mongodbatlas.ProjectTeam, 0, len(teamsToAssign))
	for teamID, assignedTeam := range teamsToAssign {
		projectTeams = append(projectTeams, assignedTeam.ProjectTeam.ToAtlas(teamID))

		//assignedTeam.Context.EnsureStatusOption(status.AtlasTeamAddProject(projectID, ""))
		//statushandler.Update(assignedTeam.Context, r.Client, r.EventRecorder, assignedTeam.Team)
	}

	_, _, err = ctx.Client.Projects.AddTeamsToProject(context.Background(), projectID, projectTeams)
	if err != nil {
		return err
	}

	return nil
}
