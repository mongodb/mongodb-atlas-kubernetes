package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/teams"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) teamReconcile(
	team *akov2.AtlasTeam,
	connectionSecretKey *client.ObjectKey,
) reconcile.Func {
	return func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		log := r.Log.With("atlasteam", req.NamespacedName)

		result := customresource.PrepareResource(ctx, r.Client, req, team, log)
		if !result.IsOk() {
			return result.ReconcileResult(), nil
		}

		if customresource.ReconciliationShouldBeSkipped(team) {
			log.Infow(fmt.Sprintf("-> Skipping AtlasTeam reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", team.Spec)
			return workflow.OK().ReconcileResult(), nil
		}

		conditions := akov2.InitCondition(team, api.FalseCondition(api.ReadyType))
		teamCtx := workflow.NewContext(log, conditions, ctx)
		log.Infow("-> Starting AtlasTeam reconciliation", "spec", team.Spec)
		defer statushandler.Update(teamCtx, r.Client, r.EventRecorder, team)

		resourceVersionIsValid := customresource.ValidateResourceVersion(teamCtx, team, r.Log)
		if !resourceVersionIsValid.IsOk() {
			r.Log.Debugf("team validation result: %v", resourceVersionIsValid)
			return resourceVersionIsValid.ReconcileResult(), nil
		}

		if !r.AtlasProvider.IsResourceSupported(team) {
			result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasTeam is not supported by Atlas for government").
				WithoutRetry()
			setCondition(teamCtx, api.ReadyType, result)
			return result.ReconcileResult(), nil
		}

		atlasClient, orgID, err := r.AtlasProvider.SdkClient(teamCtx.Context, connectionSecretKey, log)
		if err != nil {
			result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
			setCondition(teamCtx, api.ReadyType, result)
			return result.ReconcileResult(), nil
		}
		teamCtx.OrgID = orgID
		teamCtx.SdkClient = atlasClient

		teamID, result := r.ensureTeamState(teamCtx, team)
		if !result.IsOk() {
			teamCtx.SetConditionFromResult(api.ReadyType, result)
			if result.IsWarning() {
				teamCtx.Log.Warnf("failed to ensure team state %v: %s", team.Spec, result.GetMessage())
			}

			return result.ReconcileResult(), nil
		}

		teamCtx.EnsureStatusOption(status.AtlasTeamSetID(teamID))

		result = r.ensureTeamUsersAreInSync(teamCtx, teamID, team)
		if !result.IsOk() {
			teamCtx.SetConditionFromResult(api.ReadyType, result)
			return result.ReconcileResult(), nil
		}

		if team.GetDeletionTimestamp().IsZero() {
			if len(team.Status.Projects) > 0 {
				log.Debugw("Adding deletion finalizer", "name", customresource.FinalizerLabel)
				customresource.SetFinalizer(team, customresource.FinalizerLabel)
			} else {
				log.Debugw("Removing deletion finalizer", "name", customresource.FinalizerLabel)
				customresource.UnsetFinalizer(team, customresource.FinalizerLabel)
			}

			if err = r.Client.Update(teamCtx.Context, team); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to update finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}

		if !team.GetDeletionTimestamp().IsZero() {
			if customresource.HaveFinalizer(team, customresource.FinalizerLabel) {
				log.Warnf("team %s is assigned to a project. Remove it from all projects before delete", team.Name)
			} else if customresource.IsResourcePolicyKeepOrDefault(team, r.ObjectDeletionProtection) {
				log.Info("Not removing Team from Atlas as per configuration")
				return workflow.OK().ReconcileResult(), nil
			} else {
				log.Infow("-> Starting AtlasTeam deletion", "spec", team.Spec)
				_, _, err := teamCtx.SdkClient.TeamsApi.DeleteTeam(teamCtx.Context, orgID, team.Status.ID).Execute()
				var apiError *mongodbatlas.ErrorResponse
				if errors.As(err, &apiError) && apiError.ErrorCode == atlas.NotInGroup {
					log.Infow("team does not exist", "projectID", team.Status.ID)
					return workflow.Terminate(workflow.TeamDoesNotExist, err.Error()).ReconcileResult(), nil
				}
			}
		}

		err = customresource.ApplyLastConfigApplied(teamCtx.Context, team, r.Client)
		if err != nil {
			result = workflow.Terminate(workflow.Internal, err.Error())
			teamCtx.SetConditionFromResult(api.ReadyType, result)
			log.Error(result.GetMessage())

			return result.ReconcileResult(), nil
		}

		teamCtx.SetConditionTrue(api.ReadyType)
		return workflow.OK().ReconcileResult(), nil
	}
}

func (r *AtlasProjectReconciler) ensureTeamState(workflowCtx *workflow.Context, team *akov2.AtlasTeam) (string, workflow.Result) {
	var atlasAssignedTeam *teams.AssignedTeam
	var err error

	if team.Status.ID != "" {
		atlasAssignedTeam, err = r.fetchTeamByID(workflowCtx, team.Status.ID)
	} else {
		atlasAssignedTeam, err = r.fetchTeamByName(workflowCtx, team.Spec.Name)
	}
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
	}

	if atlasAssignedTeam == nil {
		desiredAtlasTeam := teams.NewTeam(&team.Spec, team.Status.ID)
		if desiredAtlasTeam == nil {
			return "", workflow.Terminate(workflow.TeamInvalidSpec, "teamspec is invalid")
		}

		atlasTeam, err := r.createTeam(workflowCtx, desiredAtlasTeam)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}
		return atlasTeam.TeamID, workflow.OK()
	}

	atlasAssignedTeam, err = r.renameTeam(workflowCtx, atlasAssignedTeam, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
	}

	return atlasAssignedTeam.TeamID, workflow.OK()
}

func (r *AtlasProjectReconciler) ensureTeamUsersAreInSync(workflowCtx *workflow.Context, teamID string, team *akov2.AtlasTeam) workflow.Result {
	atlasUsers, err := r.teamsService.GetTeamUsers(workflowCtx.Context, workflowCtx.OrgID, teamID)
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	usernamesMap := map[string]struct{}{}
	for _, username := range team.Spec.Usernames {
		usernamesMap[string(username)] = struct{}{}
	}

	atlasUsernamesMap := map[string]teams.TeamUser{}
	for _, atlasUser := range atlasUsers {
		atlasUsernamesMap[atlasUser.Username] = atlasUser
	}

	g, taskContext := errgroup.WithContext(workflowCtx.Context)

	for _, user := range atlasUsers {
		if _, ok := usernamesMap[user.Username]; !ok {
			g.Go(func() error {
				workflowCtx.Log.Debugf("removing user %s from team %s", user.UserID, teamID)
				err := r.teamsService.RemoveUser(workflowCtx.Context, workflowCtx.OrgID, teamID, user.UserID)
				return err
			})
		}
	}

	if err = g.Wait(); err != nil {
		workflowCtx.Log.Warnf("failed to remove user(s) from team %s", teamID)

		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	g, taskContext = errgroup.WithContext(workflowCtx.Context)
	toAdd := make([]teams.TeamUser, 0, len(team.Spec.Usernames))
	lock := sync.Mutex{}
	for i := range team.Spec.Usernames {
		username := team.Spec.Usernames[i]
		if _, ok := atlasUsernamesMap[string(username)]; !ok {
			g.Go(func() error {
				user, _, err := workflowCtx.SdkClient.MongoDBCloudUsersApi.GetUserByUsername(taskContext, string(username)).Execute()

				if err != nil {
					return err
				}

				lock.Lock()
				toAdd = append(toAdd, teams.TeamUser{UserID: user.GetId()})
				lock.Unlock()

				return nil
			})
		}
	}

	if err = g.Wait(); err != nil {
		workflowCtx.Log.Warnf("failed to retrieve users to add to the team %s", teamID)

		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	if len(toAdd) == 0 {
		return workflow.OK()
	}

	workflowCtx.Log.Debugf("Adding users to team %s", teamID)
	err = r.teamsService.AddUsers(workflowCtx.Context, &toAdd, workflowCtx.OrgID, teamID)
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	return workflow.OK()
}

func (r *AtlasProjectReconciler) fetchTeamByID(workflowCtx *workflow.Context, teamID string) (*teams.AssignedTeam, error) {
	workflowCtx.Log.Debugf("fetching team %s from atlas", teamID)
	atlasTeam, err := r.teamsService.GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, teamID)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func (r *AtlasProjectReconciler) fetchTeamByName(workflowCtx *workflow.Context, teamName string) (*teams.AssignedTeam, error) {
	workflowCtx.Log.Debugf("fetching team named %s from atlas", teamName)
	atlasTeam, err := r.teamsService.GetTeamByName(workflowCtx.Context, workflowCtx.OrgID, teamName)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func (r *AtlasProjectReconciler) createTeam(workflowCtx *workflow.Context, desiredAtlasTeam *teams.Team) (*teams.Team, error) {
	workflowCtx.Log.Debugf("create team named %s in atlas", desiredAtlasTeam.TeamName)
	atlasTeam, err := r.teamsService.Create(workflowCtx.Context, desiredAtlasTeam, workflowCtx.OrgID)
	if err != nil {
		return nil, err
	}
	return atlasTeam, nil
}

func (r *AtlasProjectReconciler) renameTeam(workflowCtx *workflow.Context, at *teams.AssignedTeam, newName string) (*teams.AssignedTeam, error) {
	if at.TeamName == newName {
		return at, nil
	}
	workflowCtx.Log.Debugf("updating name of team %s in atlas", at.TeamID)
	atlasTeam, err := r.teamsService.RenameTeam(workflowCtx.Context, at, workflowCtx.OrgID, newName)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func (r *AtlasProjectReconciler) teamsManagedByAtlas(workflowCtx *workflow.Context) customresource.AtlasChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		team, ok := resource.(*akov2.AtlasTeam)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasTeams")
		}

		if team.Status.ID == "" {
			return false, nil
		}

		atlasTeam, err := r.teamsService.GetTeamByID(workflowCtx.Context, workflowCtx.OrgID, team.Status.ID)
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.NotInGroup || apiError.ErrorCode == atlas.ResourceNotFound) {
				return false, nil
			}

			return false, err
		}

		atlasTeamUsers, err := r.teamsService.GetTeamUsers(workflowCtx.Context, workflowCtx.OrgID, team.Status.ID)
		if err != nil {
			return false, err
		}

		if len(atlasTeamUsers) == 0 || team.Spec.Name != atlasTeam.TeamName {
			return false, err
		}

		usernames := make([]string, 0, len(team.Spec.Usernames))
		for _, username := range team.Spec.Usernames {
			usernames = append(usernames, string(username))
		}

		atlasUsernames := make([]string, 0, len(atlasTeamUsers))
		for _, user := range atlasTeamUsers {
			atlasUsernames = append(atlasUsernames, user.Username)
		}

		return cmp.Diff(usernames, atlasUsernames) != "", nil
	}
}
