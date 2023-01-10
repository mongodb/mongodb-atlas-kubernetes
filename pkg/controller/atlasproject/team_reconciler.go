package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasProjectReconciler) teamReconcile(
	team *v1.AtlasTeam,
	connection atlas.Connection,
) reconcile.Func {
	return func(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
		log := r.Log.With("atlasteam", req.NamespacedName)

		result := customresource.PrepareResource(r.Client, req, team, log)
		if !result.IsOk() {
			return result.ReconcileResult(), nil
		}

		if shouldSkip := customresource.ReconciliationShouldBeSkipped(team); shouldSkip {
			log.Infow(fmt.Sprintf("-> Skipping AtlasTeam reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", team.Spec)
			return workflow.OK().ReconcileResult(), nil
		}

		teamCtx, err := createTeamContextFromParent(team, r.Client, connection, r.AtlasDomain, log)
		if err != nil {
			teamCtx.SetConditionFalse(status.ReadyType)
			return workflow.Terminate(workflow.Internal, err.Error()).ReconcileResult(), nil
		}

		defer statushandler.Update(teamCtx, r.Client, r.EventRecorder, team)

		resourceVersionIsValid := customresource.ValidateResourceVersion(teamCtx, team, r.Log)
		if !resourceVersionIsValid.IsOk() {
			r.Log.Debugf("team validation result: %v", resourceVersionIsValid)
			return resourceVersionIsValid.ReconcileResult(), nil
		}

		log.Infow("-> Starting AtlasTeam reconciliation", "spec", team.Spec)

		teamID, result := ensureTeamState(ctx, teamCtx, team)
		if !result.IsOk() {
			teamCtx.SetConditionFromResult(status.ReadyType, result)
			if result.IsWarning() {
				teamCtx.Log.Warnf("failed to ensure team state %v: %s", team.Spec, result.GetMessage())
			}

			return result.ReconcileResult(), nil
		}

		teamCtx.EnsureStatusOption(status.AtlasTeamSetID(teamID))

		result = ensureTeamUsersAreInSync(ctx, teamCtx, teamID, team)
		if !result.IsOk() {
			teamCtx.SetConditionFromResult(status.ReadyType, result)
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

			if err = r.Client.Update(ctx, team); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to update finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}

		if !team.GetDeletionTimestamp().IsZero() {
			if customresource.HaveFinalizer(team, customresource.FinalizerLabel) {
				log.Warnf("team %s is assigned to a project. Remove it from all projects before delete", team.Name)
			} else if customresource.ResourceShouldBeLeftInAtlas(team) {
				log.Infof("Not removing the Atlas Team from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
			} else {
				log.Infow("-> Starting AtlasTeam deletion", "spec", team.Spec)
				_, err := teamCtx.Client.Teams.RemoveTeamFromOrganization(ctx, teamCtx.Connection.OrgID, team.Status.ID)
				var apiError *mongodbatlas.ErrorResponse
				if errors.As(err, &apiError) && apiError.ErrorCode == atlas.NotInGroup {
					log.Infow("team does not exist", "projectID", team.Status.ID)
					return workflow.Terminate(workflow.TeamDoesNotExist, err.Error()).ReconcileResult(), nil
				}
			}
		}

		teamCtx.SetConditionTrue(status.ReadyType)
		return workflow.OK().ReconcileResult(), nil
	}
}

func createTeamContextFromParent(
	team *v1.AtlasTeam,
	kubeClient client.Client,
	atlasConnection atlas.Connection,
	atlasDomain string,
	logger *zap.SugaredLogger,
) (*workflow.Context, error) {
	teamCtx := customresource.MarkReconciliationStarted(kubeClient, team, logger)
	teamCtx.Connection = atlasConnection
	atlasClient, err := atlas.Client(atlasDomain, atlasConnection, logger)
	if err != nil {
		return nil, err
	}
	teamCtx.Client = atlasClient

	return teamCtx, nil
}

func ensureTeamState(ctx context.Context, workflowCtx *workflow.Context, team *v1.AtlasTeam) (string, workflow.Result) {
	var atlasTeam *mongodbatlas.Team
	var err error

	if team.Status.ID != "" {
		atlasTeam, err = fetchTeamByID(ctx, workflowCtx, team.Status.ID)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}

		atlasTeam, err = renameTeam(ctx, workflowCtx, atlasTeam, team.Spec.Name)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
		}

		return atlasTeam.ID, workflow.OK()
	}

	atlasTeam, err = fetchTeamByName(ctx, workflowCtx, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
	}

	if atlasTeam == nil {
		atlasTeam, err = team.ToAtlas()
		if err != nil {
			return "", workflow.Terminate(workflow.TeamInvalidSpec, err.Error())
		}

		atlasTeam, err = createTeam(ctx, workflowCtx, atlasTeam)
		if err != nil {
			return "", workflow.Terminate(workflow.TeamNotCreatedInAtlas, err.Error())
		}
	}

	atlasTeam, err = renameTeam(ctx, workflowCtx, atlasTeam, team.Spec.Name)
	if err != nil {
		return "", workflow.Terminate(workflow.TeamNotUpdatedInAtlas, err.Error())
	}

	return atlasTeam.ID, workflow.OK()
}

func ensureTeamUsersAreInSync(ctx context.Context, workflowCtx *workflow.Context, teamID string, team *v1.AtlasTeam) workflow.Result {
	atlasUsers, _, err := workflowCtx.Client.Teams.GetTeamUsersAssigned(ctx, workflowCtx.Connection.OrgID, teamID)
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	usernamesMap := map[string]struct{}{}
	for _, username := range team.Spec.Usernames {
		usernamesMap[string(username)] = struct{}{}
	}

	atlasUsernamesMap := map[string]mongodbatlas.AtlasUser{}
	for _, atlasUser := range atlasUsers {
		atlasUsernamesMap[atlasUser.Username] = atlasUser
	}

	g, taskContext := errgroup.WithContext(ctx)

	for _, user := range atlasUsers {
		if _, ok := usernamesMap[user.Username]; !ok {
			g.Go(func() error {
				workflowCtx.Log.Debugf("removing user %s from team %s", user.ID, teamID)
				_, err := workflowCtx.Client.Teams.RemoveUserToTeam(taskContext, workflowCtx.Connection.OrgID, teamID, user.ID)

				return err
			})
		}
	}

	if err = g.Wait(); err != nil {
		workflowCtx.Log.Warnf("failed to remove user(s) from team %s", teamID)

		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	g, taskContext = errgroup.WithContext(ctx)
	toAdd := make([]string, 0, len(team.Spec.Usernames))
	lock := sync.Mutex{}
	for _, username := range team.Spec.Usernames {
		if _, ok := atlasUsernamesMap[string(username)]; !ok {
			g.Go(func() error {
				user, _, err := workflowCtx.Client.AtlasUsers.GetByName(taskContext, string(username))

				if err != nil {
					return err
				}

				lock.Lock()
				toAdd = append(toAdd, user.ID)
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
	_, _, err = workflowCtx.Client.Teams.AddUsersToTeam(ctx, workflowCtx.Connection.OrgID, teamID, toAdd)
	if err != nil {
		return workflow.Terminate(workflow.TeamUsersNotReady, err.Error())
	}

	return workflow.OK()
}

func fetchTeamByID(ctx context.Context, workflowCtx *workflow.Context, teamID string) (*mongodbatlas.Team, error) {
	workflowCtx.Log.Debugf("fetching team %s from atlas", teamID)
	atlasTeam, _, err := workflowCtx.Client.Teams.Get(ctx, workflowCtx.Connection.OrgID, teamID)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func fetchTeamByName(ctx context.Context, workflowCtx *workflow.Context, teamName string) (*mongodbatlas.Team, error) {
	workflowCtx.Log.Debugf("fetching team named %s from atlas", teamName)
	atlasTeam, resp, err := workflowCtx.Client.Teams.GetOneTeamByName(ctx, workflowCtx.Connection.OrgID, teamName)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return nil, err
	}

	return atlasTeam, nil
}

func createTeam(ctx context.Context, workflowCtx *workflow.Context, atlasTeam *mongodbatlas.Team) (*mongodbatlas.Team, error) {
	workflowCtx.Log.Debugf("create team named %s in atlas", atlasTeam.Name)
	atlasTeam, _, err := workflowCtx.Client.Teams.Create(ctx, workflowCtx.Connection.OrgID, atlasTeam)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}

func renameTeam(ctx context.Context, workflowCtx *workflow.Context, atlasTeam *mongodbatlas.Team, newName string) (*mongodbatlas.Team, error) {
	if atlasTeam.Name == newName {
		return atlasTeam, nil
	}

	workflowCtx.Log.Debugf("updating name of team %s in atlas", atlasTeam.ID)
	atlasTeam, _, err := workflowCtx.Client.Teams.Rename(ctx, workflowCtx.Connection.OrgID, atlasTeam.ID, newName)
	if err != nil {
		return nil, err
	}

	return atlasTeam, nil
}
