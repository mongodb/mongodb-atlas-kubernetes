package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/google/go-cmp/cmp"

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
	projectId string,
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
		// SELFNOTE: may need to include teamCtx.Client (ie old api client) if sdkAPI does not support all fields

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
				_, err := teamCtx.SdkClient.TeamsApi.RemoveProjectTeam(teamCtx.Context, projectId, team.Status.ID).Execute()
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
	var atlasTeam *admin.Team
	var err error

	if team.Status.ID != "" {
		atlasTeam, err = fetchTeamByID(workflowCtx, team.Status.ID)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}

		atlasTeam, err = renameTeam(workflowCtx, atlasTeam, team.Spec.Name)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
		}

		return *atlasTeam.Id, workflow.OK()
	}

	atlasTeam, err = fetchTeamByName(workflowCtx, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
	}

	if atlasTeam == nil {
		atlasTeam, err = team.ToAtlas()
		if err != nil {
			return "", workflow.Terminate(workflow.TeamInvalidSpec, err.Error())
		}

		atlasTeam, err = createTeam(workflowCtx, atlasTeam)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}
	}

	atlasTeam, err = renameTeam(workflowCtx, atlasTeam, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
	}

	return *atlasTeam.Id, workflow.OK()
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
	for _, atlasUser := range *atlasUsers.Results { // SELFNOTE: intuitively, this doesn't feel right... * to *admin.PaginatedApiAppUser child *Results
		atlasUsernamesMap[atlasUser.Username] = atlasUser
	}

	g, taskContext := errgroup.WithContext(workflowCtx.Context)

	for i := range *atlasUsers.Results {
		user := (*atlasUsers.Results)[i]
		if _, ok := usernamesMap[(*atlasUsers.Results)[i].Username]; !ok {
			g.Go(func() error {
				workflowCtx.Log.Debugf("removing user %s from team %s", user.Id, teamID)
				_, err := workflowCtx.SdkClient.TeamsApi.RemoveTeamUser(taskContext, workflowCtx.OrgID, teamID, *user.Id).Execute()
				return err
			})
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
				toAdd = append(toAdd, admin.AddUserToTeam{Id: *user.Id})
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

func fetchTeamByID(workflowCtx *workflow.Context, teamID string) (*admin.Team, error) {
	workflowCtx.Log.Debugf("fetching team %s from atlas", teamID)
	atlasTeamResponse, _, err := workflowCtx.SdkClient.TeamsApi.GetTeamById(workflowCtx.Context, workflowCtx.OrgID, teamID).Execute()
	if err != nil {
		return nil, err
	}
	atlasTeam := &admin.Team{Id: atlasTeamResponse.Id, Name: *atlasTeamResponse.Name, Links: atlasTeamResponse.Links}

	return atlasTeam, nil
}

func fetchTeamByName(workflowCtx *workflow.Context, teamName string) (*admin.Team, error) {
	workflowCtx.Log.Debugf("fetching team named %s from atlas", teamName)
	atlasTeamResponse, resp, err := workflowCtx.SdkClient.TeamsApi.GetTeamByName(workflowCtx.Context, workflowCtx.OrgID, teamName).Execute()
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return nil, err
	}
	atlasTeam := &admin.Team{Id: atlasTeamResponse.Id, Name: *atlasTeamResponse.Name, Links: atlasTeamResponse.Links}

	return atlasTeam, nil
}

func createTeam(workflowCtx *workflow.Context, atlasTeam *admin.Team) (*admin.Team, error) {
	workflowCtx.Log.Debugf("create team named %s in atlas", atlasTeam.Name)
	atlasTeam, _, err := workflowCtx.SdkClient.TeamsApi.CreateTeam(workflowCtx.Context, workflowCtx.OrgID, atlasTeam).Execute()
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func renameTeam(workflowCtx *workflow.Context, atlasTeam *admin.Team, newName string) (*admin.Team, error) {
	if atlasTeam.Name == newName {
		return atlasTeam, nil
	}
	teamUpdate := admin.TeamUpdate{Links: atlasTeam.Links, Name: newName}

	workflowCtx.Log.Debugf("updating name of team %s in atlas", atlasTeam.Id)
	atlasTeamResponse, _, err := workflowCtx.SdkClient.TeamsApi.RenameTeam(workflowCtx.Context, workflowCtx.OrgID, *atlasTeam.Id, &teamUpdate).Execute()
	if err != nil {
		return nil, err
	}
	atlasTeam.SetLinks(*atlasTeamResponse.Links)
	// atlasTeam.SetId(*atlasTeamResponse.Id)  //SELFNOTE: see if id/links can change or if only new name needs to be changed
	atlasTeam.SetName(*atlasTeamResponse.Name)
	return atlasTeam, nil
}

func teamsManagedByAtlas(workflowCtx *workflow.Context) customresource.AtlasChecker { // SELFNOTE: to check if get user work around is working properly (ie is the number of users accurately reflected in len(*atlasTeamUsers.Results) )
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
			var apiError *mongodbatlas.ErrorResponse // SELFNOTE: Feels a little wrong to mix and match like this but admin.ApiError returns in `second argument to errors.As must be a non-nil pointer to either a type that implements error, or to any interface type` for line below
			if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.NotInGroup || apiError.ErrorCode == atlas.ResourceNotFound) {
				return false, nil
			}

			return false, err
		}

		atlasTeamUsers, _, err := workflowCtx.SdkClient.TeamsApi.ListTeamUsers(workflowCtx.Context, workflowCtx.OrgID, *atlasTeam.Name).Execute()
		if err != nil {
			return false, nil
		}

		if team.Spec.Name != *atlasTeam.Name || len(*atlasTeamUsers.Results) == 0 {
			return false, err
		}

		usernames := make([]string, 0, len(team.Spec.Usernames))
		for _, username := range team.Spec.Usernames {
			usernames = append(usernames, string(username))
		}

		return cmp.Diff(usernames, *atlasTeamUsers.Results) != "", nil
	}
}
