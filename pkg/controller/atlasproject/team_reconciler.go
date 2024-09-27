package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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

		teamID, result := ensureTeamState(teamCtx, team)
		if !result.IsOk() {
			teamCtx.SetConditionFromResult(api.ReadyType, result)
			if result.IsWarning() {
				teamCtx.Log.Warnf("failed to ensure team state %v: %s", team.Spec, result.GetMessage())
			}

			return result.ReconcileResult(), nil
		}

		teamCtx.EnsureStatusOption(status.AtlasTeamSetID(teamID))

		result = ensureTeamUsersAreInSync(teamCtx, teamID, team)
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

func ensureTeamState(workflowCtx *workflow.Context, team *akov2.AtlasTeam) (string, workflow.Result) {
	var atlasTeamResponse *admin.TeamResponse
	var err error

	if team.Status.ID != "" {
		atlasTeamResponse, err = fetchTeamByID(workflowCtx, team.Status.ID)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}

		atlasTeamResponse, err = renameTeam(workflowCtx, atlasTeamResponse, team.Spec.Name)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
		}

		return atlasTeamResponse.GetId(), workflow.OK()
	}

	atlasTeamResponse, err = fetchTeamByName(workflowCtx, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
	}

	if atlasTeamResponse == nil {
		atlasTeam, err := team.ToAtlas()
		if err != nil {
			return "", workflow.Terminate(workflow.TeamInvalidSpec, err.Error())
		}

		atlasTeam, err = createTeam(workflowCtx, atlasTeam)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}

		return atlasTeam.GetId(), workflow.OK()
	}

	atlasTeamResponse, err = renameTeam(workflowCtx, atlasTeamResponse, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
	}

	return atlasTeamResponse.GetId(), workflow.OK()
}

func ensureTeamUsersAreInSync(workflowCtx *workflow.Context, teamID string, team *akov2.AtlasTeam) workflow.Result {
	atlasUsers, _, err := workflowCtx.SdkClient.TeamsApi.ListTeamUsers(workflowCtx.Context, workflowCtx.OrgID, teamID).Execute()
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	usernamesMap := map[string]struct{}{}
	for _, username := range team.Spec.Usernames {
		usernamesMap[string(username)] = struct{}{}
	}

	atlasUsernamesMap := map[string]admin.CloudAppUser{}
	for _, atlasUser := range atlasUsers.GetResults() {
		atlasUsernamesMap[atlasUser.Username] = atlasUser
	}

	g, taskContext := errgroup.WithContext(workflowCtx.Context)

	if atlasUsers.Results != nil {
		for _, user := range atlasUsers.GetResults() {
			if _, ok := usernamesMap[user.Username]; !ok {
				g.Go(func() error {
					workflowCtx.Log.Debugf("removing user %s from team %s", user.GetId(), teamID)
					_, err := workflowCtx.SdkClient.TeamsApi.RemoveTeamUser(taskContext, workflowCtx.OrgID, teamID, user.GetId()).Execute()
					return err
				})
			}
		}
	}

	if err = g.Wait(); err != nil {
		workflowCtx.Log.Warnf("failed to remove user(s) from team %s", teamID)

		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	g, taskContext = errgroup.WithContext(workflowCtx.Context)
	toAdd := make([]admin.AddUserToTeam, 0, len(team.Spec.Usernames))
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
				toAdd = append(toAdd, admin.AddUserToTeam{Id: user.GetId()})
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
	_, _, err = workflowCtx.SdkClient.TeamsApi.AddTeamUser(workflowCtx.Context, workflowCtx.OrgID, teamID, &toAdd).Execute()
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	return workflow.OK()
}

func fetchTeamByID(workflowCtx *workflow.Context, teamID string) (*admin.TeamResponse, error) {
	workflowCtx.Log.Debugf("fetching team %s from atlas", teamID)
	atlasTeamResponse, _, err := workflowCtx.SdkClient.TeamsApi.GetTeamById(workflowCtx.Context, workflowCtx.OrgID, teamID).Execute()
	if err != nil {
		return nil, err
	}

	return atlasTeamResponse, nil
}

func fetchTeamByName(workflowCtx *workflow.Context, teamName string) (*admin.TeamResponse, error) {
	workflowCtx.Log.Debugf("fetching team named %s from atlas", teamName)
	atlasTeamResponse, resp, err := workflowCtx.SdkClient.TeamsApi.GetTeamByName(workflowCtx.Context, workflowCtx.OrgID, teamName).Execute()
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return nil, err
	}

	return atlasTeamResponse, nil
}

func createTeam(workflowCtx *workflow.Context, atlasTeam *admin.Team) (*admin.Team, error) {
	workflowCtx.Log.Debugf("create team named %s in atlas", atlasTeam.Name)
	atlasTeam, _, err := workflowCtx.SdkClient.TeamsApi.CreateTeam(workflowCtx.Context, workflowCtx.OrgID, atlasTeam).Execute()
	if err != nil {
		return nil, err
	}
	return atlasTeam, nil
}

func renameTeam(workflowCtx *workflow.Context, atlasTeamResponse *admin.TeamResponse, newName string) (*admin.TeamResponse, error) {
	if atlasTeamResponse.GetName() == newName {
		return atlasTeamResponse, nil
	}

	workflowCtx.Log.Debugf("updating name of team %s in atlas", atlasTeamResponse.GetId())
	teamUpdate := admin.TeamUpdate{Name: newName}
	atlasTeamResponse, _, err := workflowCtx.SdkClient.TeamsApi.RenameTeam(workflowCtx.Context, workflowCtx.OrgID, atlasTeamResponse.GetId(), &teamUpdate).Execute()
	if err != nil {
		return nil, err
	}

	return atlasTeamResponse, nil
}

func teamsManagedByAtlas(workflowCtx *workflow.Context) customresource.AtlasChecker {
	return func(resource api.AtlasCustomResource) (bool, error) {
		team, ok := resource.(*akov2.AtlasTeam)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasTeams")
		}

		if team.Status.ID == "" {
			return false, nil
		}

		atlasTeam, _, err := workflowCtx.SdkClient.TeamsApi.GetTeamById(workflowCtx.Context, workflowCtx.OrgID, team.Status.ID).Execute()
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.NotInGroup || apiError.ErrorCode == atlas.ResourceNotFound) {
				return false, nil
			}

			return false, err
		}

		atlasTeamUsers, _, err := workflowCtx.SdkClient.TeamsApi.ListTeamUsers(workflowCtx.Context, workflowCtx.OrgID, atlasTeam.GetName()).Execute()
		if err != nil {
			return false, err
		}

		if team.Spec.Name != atlasTeam.GetName() || len(atlasTeamUsers.GetResults()) == 0 {
			return false, err
		}

		usernames := make([]string, 0, len(team.Spec.Usernames))
		for _, username := range team.Spec.Usernames {
			usernames = append(usernames, string(username))
		}

		atlasUsernames := make([]string, 0, len(atlasTeamUsers.GetResults()))
		for _, user := range atlasTeamUsers.GetResults() {
			atlasUsernames = append(atlasUsernames, user.Username)
		}

		return cmp.Diff(usernames, atlasUsernames) != "", nil
	}
}
