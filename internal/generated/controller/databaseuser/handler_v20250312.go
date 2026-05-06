// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package databaseuser

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312018/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	atlasapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

type Handlerv20250312 struct {
	kubeClient         client.Client
	atlasClient        *v20250312sdk.APIClient
	translator         crapi.Translator
	deletionProtection bool
}

func NewHandlerv20250312(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translator crapi.Translator, deletionProtection bool) *Handlerv20250312 {
	return &Handlerv20250312{
		atlasClient:        atlasClient,
		deletionProtection: deletionProtection,
		kubeClient:         kubeClient,
		translator:         translator,
	}
}

// HandleInitial handles the initial state for version v20250312
func (h *Handlerv20250312) HandleInitial(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to resolve DatabaseUser dependencies: %w", err))
	}

	atlasDBUser := &v20250312sdk.CloudDatabaseUser{}
	if err := h.translator.ToAPI(atlasDBUser, databaseuser, deps...); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate DatabaseUser to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.DatabaseUsersApi.CreateDatabaseUser(ctx, atlasDBUser.GroupId, atlasDBUser).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create DatabaseUser: %w", err))
	}

	databaseuserCopy := databaseuser.DeepCopy()
	if _, err := h.translator.FromAPI(databaseuserCopy, response); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate DatabaseUser from Atlas: %w", err))
	}

	if err := ctrlstate.NewPatcher(databaseuserCopy).UpdateStatus().UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to patch DatabaseUser status: %w", err))
	}

	return result.NextState(state.StateCreating, "DatabaseUser created. Waiting for clusters to apply changes.")
}

// HandleImportRequested handles the importrequested state for version v20250312.
// The annotation mongodb.com/external-id must be set to "groupId:databaseName:username".
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	externalID, ok := databaseuser.GetAnnotations()["mongodb.com/external-id"]
	if !ok {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-id"))
	}

	databaseName, username, err := parseExternalID(externalID)
	if err != nil {
		return result.Error(state.StateImportRequested, err)
	}

	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to resolve Cluster dependencies: %w", err))
	}

	params := &v20250312sdk.GetDatabaseUserApiParams{
		DatabaseName: databaseName,
		Username:     username,
	}
	err = h.translator.ToAPI(params, databaseuser, deps...)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate cluster API parameters to Atlas: %w", err))
	}

	response, _, err := h.atlasClient.DatabaseUsersApi.GetDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to get DatabaseUser %q: %w", externalID, err))
	}

	databaseuserCopy := databaseuser.DeepCopy()
	if _, err := h.translator.FromAPI(databaseuserCopy, response); err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to translate DatabaseUser from Atlas: %w", err))
	}

	if err := ctrlstate.NewPatcher(databaseuserCopy).UpdateStatus().UpdateStateTracker().Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to patch DatabaseUser status: %w", err))
	}

	return result.NextState(state.StateImported, "DatabaseUser imported.")
}

// HandleImported handles the imported state for version v20250312
func (h *Handlerv20250312) HandleImported(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	if expired, err := checkExpiry(databaseuser); err != nil {
		return result.Error(state.StateImported, fmt.Errorf("failed to check DatabaseUser expiry: %w", err))
	} else if expired {
		return result.NextState(state.StateDeletionRequested, "DatabaseUser has expired.")
	}

	return h.handleUpserted(ctx, state.StateImported, databaseuser)
}

// HandleCreating polls cluster readiness after a DatabaseUser has been created in Atlas.
func (h *Handlerv20250312) HandleCreating(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	groupID, _, _, err := h.resolveIdentity(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateCreating, err)
	}

	ready, err := h.deploymentsReady(ctx, groupID, clusterScopes(databaseuser))
	if err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to check cluster readiness: %w", err))
	}

	if !ready {
		return result.NextState(state.StateCreating, "Waiting for clusters to apply DatabaseUser changes.")
	}

	return result.NextState(state.StateCreated, "DatabaseUser is ready.")
}

// HandleCreated handles the created state for version v20250312
func (h *Handlerv20250312) HandleCreated(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	if expired, err := checkExpiry(databaseuser); err != nil {
		return result.Error(state.StateCreated, fmt.Errorf("failed to check DatabaseUser expiry: %w", err))
	} else if expired {
		return result.NextState(state.StateDeletionRequested, "DatabaseUser has expired.")
	}

	return h.handleUpserted(ctx, state.StateCreated, databaseuser)
}

// HandleUpdating polls cluster readiness after a DatabaseUser has been updated in Atlas.
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	groupID, _, _, err := h.resolveIdentity(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateUpdating, err)
	}

	ready, err := h.deploymentsReady(ctx, groupID, clusterScopes(databaseuser))
	if err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to check cluster readiness: %w", err))
	}

	if !ready {
		return result.NextState(state.StateUpdating, "Waiting for clusters to apply DatabaseUser changes.")
	}

	return result.NextState(state.StateUpdated, "DatabaseUser is up to date.")
}

// HandleUpdated handles the updated state for version v20250312
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	if expired, err := checkExpiry(databaseuser); err != nil {
		return result.Error(state.StateUpdated, fmt.Errorf("failed to check DatabaseUser expiry: %w", err))
	} else if expired {
		return result.NextState(state.StateDeletionRequested, "DatabaseUser has expired.")
	}

	return h.handleUpserted(ctx, state.StateUpdated, databaseuser)
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(databaseuser, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "DatabaseUser skipped deletion due to retention policy.")
	}

	groupID, databaseName, username, err := h.identityFromStatusOrDeps(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}

	_, err = h.atlasClient.DatabaseUsersApi.DeleteDatabaseUser(ctx, groupID, databaseName, username).Execute()
	if v20250312sdk.IsErrorCode(err, atlasapi.UserNotfound) || v20250312sdk.IsErrorCode(err, atlasapi.UsernameNotFound) {
		return result.NextState(state.StateDeleted, "DatabaseUser deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete DatabaseUser: %w", err))
	}

	return result.NextState(state.StateDeleting, "Deleting DatabaseUser.")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	groupID, databaseName, username, err := h.identityFromStatusOrDeps(ctx, databaseuser)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}

	_, _, err = h.atlasClient.DatabaseUsersApi.GetDatabaseUser(ctx, groupID, databaseName, username).Execute()
	switch {
	case v20250312sdk.IsErrorCode(err, atlasapi.UserNotfound) || v20250312sdk.IsErrorCode(err, atlasapi.UsernameNotFound):
		return result.NextState(state.StateDeleted, "DatabaseUser deleted.")
	case err != nil:
		return result.Error(state.StateDeleting, fmt.Errorf("failed to check DatabaseUser deletion: %w", err))
	}

	return result.NextState(state.StateDeleting, "Deleting DatabaseUser.")
}

// For returns the resource and predicates for the controller
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.DatabaseUser{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}

func (h *Handlerv20250312) handleUpserted(ctx context.Context, currentState state.ResourceState, databaseuser *akov2generated.DatabaseUser) (ctrlstate.Result, error) {
	// Fetch dependencies first so ShouldUpdate can detect Secret.ResourceVersion changes (password rotation).
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to resolve DatabaseUser dependencies: %w", err))
	}

	update, err := ctrlstate.ShouldUpdate(databaseuser, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "DatabaseUser is up to date. No update required.")
	}

	atlasDBUser := &v20250312sdk.CloudDatabaseUser{}
	if err := h.translator.ToAPI(atlasDBUser, databaseuser, deps...); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate DatabaseUser to Atlas: %w", err))
	}

	params := &v20250312sdk.UpdateDatabaseUserApiParams{
		GroupId:           atlasDBUser.GroupId,
		DatabaseName:      atlasDBUser.DatabaseName,
		Username:          atlasDBUser.Username,
		CloudDatabaseUser: atlasDBUser,
	}

	response, _, err := h.atlasClient.DatabaseUsersApi.UpdateDatabaseUserWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update DatabaseUser: %w", err))
	}

	databaseuserCopy := databaseuser.DeepCopy()
	if _, err := h.translator.FromAPI(databaseuserCopy, response); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to translate DatabaseUser from Atlas: %w", err))
	}

	// Pass deps to UpdateStateTracker so Secret.ResourceVersion is included in the hash.
	if err := ctrlstate.NewPatcher(databaseuserCopy).UpdateStateTracker(deps...).UpdateStatus().Patch(ctx, h.kubeClient); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to patch DatabaseUser: %w", err))
	}

	return result.NextState(state.StateUpdating, "DatabaseUser updated. Waiting for clusters to apply changes.")
}

// deploymentsReady returns true when all relevant clusters have applied the latest DatabaseUser changes.
// scopedClusters lists the cluster names to check; if empty, all clusters in the project are checked.
func (h *Handlerv20250312) deploymentsReady(ctx context.Context, groupID string, scopedClusters []string) (bool, error) {
	if len(scopedClusters) > 0 {
		for _, name := range scopedClusters {
			clusterStatus, _, err := h.atlasClient.ClustersApi.GetClusterStatus(ctx, groupID, name).Execute()
			if err != nil {
				return false, fmt.Errorf("failed to get cluster %q status: %w", name, err)
			}

			if clusterStatus.GetChangeStatus() != "APPLIED" {
				return false, nil
			}
		}

		return true, nil
	}

	// No scopes specified: the user has access to all clusters, so check them all.
	clusters, _, err := h.atlasClient.ClustersApi.ListClusters(ctx, groupID).Execute()
	if err != nil {
		return false, fmt.Errorf("failed to list clusters for group %q: %w", groupID, err)
	}

	for _, cluster := range clusters.GetResults() {
		name := cluster.GetName()
		if name == "" {
			continue
		}

		clusterStatus, _, err := h.atlasClient.ClustersApi.GetClusterStatus(ctx, groupID, name).Execute()
		if err != nil {
			return false, fmt.Errorf("failed to get cluster %q status: %w", name, err)
		}

		if clusterStatus.GetChangeStatus() != "APPLIED" {
			return false, nil
		}
	}

	return true, nil
}

// clusterScopes extracts the names of CLUSTER-type scopes from the DatabaseUser spec.
// Returns nil when scopes are not specified, which means the user has access to all clusters.
func clusterScopes(databaseuser *akov2generated.DatabaseUser) []string {
	if databaseuser.Spec.V20250312 == nil ||
		databaseuser.Spec.V20250312.Entry == nil ||
		databaseuser.Spec.V20250312.Entry.Scopes == nil {
		return nil
	}

	var names []string
	for _, scope := range *databaseuser.Spec.V20250312.Entry.Scopes {
		if scope.Type == "CLUSTER" {
			names = append(names, scope.Name)
		}
	}

	return names
}

// identityFromStatusOrDeps returns the DatabaseUser composite identity.
func (h *Handlerv20250312) identityFromStatusOrDeps(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (groupID, databaseName, username string, err error) {
	if groupID, databaseName, username, ok := readIdentityFromStatus(databaseuser); ok {
		return groupID, databaseName, username, nil
	}

	if groupID, databaseName, username, ok := readIdentityFromSpec(databaseuser); ok {
		return groupID, databaseName, username, nil
	}

	if groupID, databaseName, username, ok, lookupErr := h.readIdentityViaGroupRef(ctx, databaseuser); lookupErr != nil {
		return "", "", "", fmt.Errorf("identity not in status and group lookup failed: %w", lookupErr)
	} else if ok {
		return groupID, databaseName, username, nil
	}

	groupID, databaseName, username, err = h.resolveIdentity(ctx, databaseuser)
	if err != nil {
		return "", "", "", fmt.Errorf("identity not in status and dependency resolution failed: %w", err)
	}

	return groupID, databaseName, username, nil
}

// readIdentityFromSpec reads the composite key directly from the spec when groupId is set explicitly (not via groupRef).
func readIdentityFromSpec(databaseuser *akov2generated.DatabaseUser) (groupID, databaseName, username string, ok bool) {
	spec := databaseuser.Spec.V20250312
	if spec == nil || spec.Entry == nil {
		return "", "", "", false
	}
	if spec.GroupId == nil || *spec.GroupId == "" {
		return "", "", "", false
	}
	databaseName = spec.Entry.DatabaseName
	username = spec.Entry.Username
	if databaseName == "" || username == "" {
		return "", "", "", false
	}
	return *spec.GroupId, databaseName, username, true
}

// readIdentityViaGroupRef resolves identity when a groupRef is used instead of a direct groupId.
func (h *Handlerv20250312) readIdentityViaGroupRef(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (groupID, databaseName, username string, ok bool, err error) {
	spec := databaseuser.Spec.V20250312
	if spec == nil || spec.Entry == nil || spec.GroupRef == nil {
		return "", "", "", false, nil
	}
	databaseName = spec.Entry.DatabaseName
	username = spec.Entry.Username
	if databaseName == "" || username == "" {
		return "", "", "", false, nil
	}

	group := &akov2generated.Group{}
	if lookupErr := h.kubeClient.Get(ctx, client.ObjectKey{
		Name:      spec.GroupRef.Name,
		Namespace: databaseuser.GetNamespace(),
	}, group); lookupErr != nil {
		return "", "", "", false, fmt.Errorf("failed to get Group %s/%s: %w", databaseuser.GetNamespace(), spec.GroupRef.Name, lookupErr)
	}

	if group.Status.V20250312 == nil || group.Status.V20250312.Id == nil || *group.Status.V20250312.Id == "" {
		return "", "", "", false, fmt.Errorf("Group %s/%s has no ID in status", databaseuser.GetNamespace(), spec.GroupRef.Name)
	}

	return *group.Status.V20250312.Id, databaseName, username, true, nil
}

// readIdentityFromStatus reads the composite key (groupId, databaseName, username) from the
// DatabaseUser status. Returns ok=false if status has not been populated yet.
func readIdentityFromStatus(databaseuser *akov2generated.DatabaseUser) (groupID, databaseName, username string, ok bool) {
	s := databaseuser.Status.V20250312
	if s == nil || s.GroupId == "" || s.DatabaseName == "" || s.Username == "" {
		return "", "", "", false
	}

	return s.GroupId, s.DatabaseName, s.Username, true
}

// resolveIdentity calls getDependencies, runs ToAPI to collapse all references,
// and returns the three-part key (groupId, databaseName, username) needed for Atlas API calls.
func (h *Handlerv20250312) resolveIdentity(ctx context.Context, databaseuser *akov2generated.DatabaseUser) (groupID, databaseName, username string, err error) {
	deps, err := h.getDependencies(ctx, databaseuser)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to resolve DatabaseUser dependencies: %w", err)
	}

	atlasDBUser := &v20250312sdk.CloudDatabaseUser{}
	if err := h.translator.ToAPI(atlasDBUser, databaseuser, deps...); err != nil {
		return "", "", "", fmt.Errorf("failed to resolve DatabaseUser identity: %w", err)
	}

	return atlasDBUser.GroupId, atlasDBUser.DatabaseName, atlasDBUser.Username, nil
}

// checkExpiry returns true when the DatabaseUser's deleteAfterDate has passed. Works as the original
// version for the curated AtlasDatabaseUser CR
func checkExpiry(databaseuser *akov2generated.DatabaseUser) (bool, error) {
	if databaseuser.Spec.V20250312 == nil ||
		databaseuser.Spec.V20250312.Entry == nil ||
		databaseuser.Spec.V20250312.Entry.DeleteAfterDate == nil {
		return false, nil
	}

	deleteAfter, err := time.Parse(time.RFC3339, *databaseuser.Spec.V20250312.Entry.DeleteAfterDate)
	if err != nil {
		return false, fmt.Errorf("invalid deleteAfterDate %q: %w", *databaseuser.Spec.V20250312.Entry.DeleteAfterDate, err)
	}

	return deleteAfter.Before(time.Now()), nil
}

// parseExternalID parses the mongodb.com/external-id annotation value.
// Expected format: "databaseName:username".
func parseExternalID(externalID string) (databaseName, username string, err error) {
	parts := strings.SplitN(externalID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid mongodb.com/external-id %q: expected format \"databaseName:username\"", externalID)
	}
	return parts[0], parts[1], nil
}
