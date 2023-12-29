package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	controllerruntime "sigs.k8s.io/controller-runtime"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/kube"
)

type TeamDataContainer struct {
	ProjectTeam *v1.Team
	Team        *v1.AtlasTeam
	Context     *workflow.Context
}

func (r *AtlasProjectReconciler) ensureAssignedTeams(workflowCtx *workflow.Context, project *v1.AtlasProject, protected bool) workflow.Result {
	resourcesToWatch := make([]watch.WatchedObject, 0, len(project.Spec.Teams))
	defer func() {
		workflowCtx.AddResourcesToWatch(resourcesToWatch...)
		r.Log.Debugf("watching team resources: %v\r\n", r.WatchedResources)
	}()

	teamsToAssign := map[string]*v1.Team{}
	for _, entry := range project.Spec.Teams {
		assignedTeam := entry

		if assignedTeam.TeamRef.Name == "" {
			workflowCtx.Log.Warnf("missing team name. skipping assignment for entry %v", assignedTeam)

			continue
		}

		if assignedTeam.TeamRef.Namespace == "" {
			assignedTeam.TeamRef.Namespace = project.Namespace
		}

		team := &v1.AtlasTeam{}
		teamReconciler := r.teamReconcile(team, project.ConnectionSecretObjectKey())
		_, err := teamReconciler(
			context.Background(),
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

	canReconcile, err := canAssignedTeamsReconcile(workflowCtx, r.Client, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.ProjectTeamsReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.ProjectTeamsReadyType, result)

		return result
	}

	err = r.syncAssignedTeams(workflowCtx, project.ID(), project, teamsToAssign)
	if err != nil {
		workflowCtx.SetConditionFalse(status.ProjectTeamsReadyType)
		return workflow.Terminate(workflow.ProjectTeamUnavailable, err.Error())
	}

	workflowCtx.SetConditionTrue(status.ProjectTeamsReadyType)

	if len(project.Spec.Teams) == 0 {
		workflowCtx.EnsureStatusOption(status.AtlasProjectSetTeamsOption(nil))
		workflowCtx.UnsetCondition(status.ProjectTeamsReadyType)
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) syncAssignedTeams(ctx *workflow.Context, projectID string, project *v1.AtlasProject, teamsToAssign map[string]*v1.Team) error {
	ctx.Log.Debug("fetching assigned teams from atlas")
	atlasAssignedTeams, _, err := ctx.Client.Projects.GetProjectTeamsAssigned(context.Background(), projectID)
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

		_, _, err = ctx.Client.Projects.AddTeamsToProject(context.Background(), projectID, projectTeams)
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
		if projectStat.ID == project.Status.ID {
			continue
		}

		assignedProjects = append(assignedProjects, projectStat)
	}

	log := r.Log.With("atlasteam", teamRef)
	teamCtx := customresource.MarkReconciliationStarted(r.Client, team, log, ctx.Context)

	atlasClient, orgID, err := r.AtlasProvider.Client(teamCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		return err
	}
	teamCtx.Client = atlasClient

	if len(assignedProjects) == 0 {
		log.Debugf("team %s has no project associated to it. removing from atlas.", team.Spec.Name)
		_, err = teamCtx.Client.Teams.RemoveTeamFromOrganization(context.Background(), orgID, team.Status.ID)
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
	for _, stat := range project.Status.Teams {
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

type assignedTeamInfo struct {
	ID    string
	Roles []string
}

func canAssignedTeamsReconcile(workflowCtx *workflow.Context, k8sClient client.Client, protected bool, akoProject *v1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &v1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	atlasAssignedTeams, _, err := workflowCtx.Client.Projects.GetProjectTeamsAssigned(workflowCtx.Context, akoProject.ID())
	if err != nil {
		return false, err
	}

	if atlasAssignedTeams == nil || atlasAssignedTeams.TotalCount == 0 {
		return true, nil
	}

	atlasAssignedTeamsInfo := make([]assignedTeamInfo, 0, atlasAssignedTeams.TotalCount)
	for _, atlasAssignedTeam := range atlasAssignedTeams.Results {
		if atlasAssignedTeam != nil {
			atlasAssignedTeamsInfo = append(
				atlasAssignedTeamsInfo,
				assignedTeamInfo{
					ID:    atlasAssignedTeam.TeamID,
					Roles: atlasAssignedTeam.RoleNames,
				},
			)
		}
	}

	lastAssignedTeamsInfo, err := collectTeams(workflowCtx.Context, k8sClient, latestConfig, akoProject.Namespace)
	if err != nil {
		return false, err
	}

	if cmp.Diff(atlasAssignedTeamsInfo, lastAssignedTeamsInfo, cmpopts.EquateEmpty()) == "" {
		return true, nil
	}

	currentAssignedTeamsInfo, err := collectTeams(workflowCtx.Context, k8sClient, &akoProject.Spec, akoProject.Namespace)
	if err != nil {
		return false, err
	}

	return cmp.Diff(atlasAssignedTeamsInfo, currentAssignedTeamsInfo, cmpopts.EquateEmpty()) == "", nil
}

func collectTeams(ctx context.Context, k8sClient client.Client, projectSpec *v1.AtlasProjectSpec, projectNamespace string) ([]assignedTeamInfo, error) {
	teams := make([]assignedTeamInfo, 0, len(projectSpec.Teams))

	for _, assignedTeam := range projectSpec.Teams {
		team := &v1.AtlasTeam{}
		err := k8sClient.Get(ctx, *assignedTeam.TeamRef.GetObject(projectNamespace), team)
		if err != nil {
			if !apiErrors.IsNotFound(err) {
				return nil, err
			}
		}

		info := assignedTeamInfo{
			ID: team.Status.ID,
		}
		for _, role := range assignedTeam.Roles {
			info.Roles = append(info.Roles, string(role))
		}

		teams = append(teams, info)
	}

	return teams, nil
}
