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

package connectionsecret

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/controller/connectionsecret/target"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
)

func (r *ConnectionSecretReconciler) handleUpsert(ctx context.Context, ids *ConnectionSecretIdentifiers, user *generatedv1.DatabaseUser, connectionTarget target.ConnectionTargetInstance) (reconcile.Result, error) {
	if user == nil || user.Spec.V20250312 == nil || user.Spec.V20250312.Entry == nil || user.Spec.V20250312.Entry.PasswordSecretRef == nil {
		return reconcile.Result{}, nil // nothing to do
	}

	secret := &corev1.Secret{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      user.Spec.V20250312.Entry.PasswordSecretRef.Name,
		Namespace: user.GetNamespace(),
	}, secret)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get password secret: %w", err)
	}

	password, exists := secret.Data[*user.Spec.V20250312.Entry.PasswordSecretRef.Key] // key is defaulted so cannot be nil
	if !exists {
		return reconcile.Result{}, fmt.Errorf("secret does not contain key %q", *user.Spec.V20250312.Entry.PasswordSecretRef.Key)
	}

	// create the connection data that will populate secret.stringData
	data := connectionTarget.BuildConnectionData(ctx)
	if data == nil {
		return ctrl.Result{}, nil // nothing to do
	}

	data.DBUserName = user.Spec.V20250312.Entry.Username
	data.Password = string(password)

	if err := r.ensureSecret(ctx, ids, user, connectionTarget, data); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ConnectionSecretReconciler) handleBatchUpsert(
	ctx context.Context,
	req ctrl.Request,
	user *generatedv1.DatabaseUser,
	projectID string,
	connectionTargetInstances []target.ConnectionTargetInstance,
) (ctrl.Result, error) {
	if user.Spec.V20250312 == nil || user.Spec.V20250312.Entry == nil {
		return ctrl.Result{}, errors.New("user spec has no entry")
	}

	databaseUsername := user.Spec.V20250312.Entry.Username
	deleteAfterDate := ""
	if user.Spec.V20250312.Entry.DeleteAfterDate != nil {
		deleteAfterDate = *user.Spec.V20250312.Entry.DeleteAfterDate
	}

	groupName := ""
	if len(connectionTargetInstances) > 0 {
		atlasClient, err := r.getSDKClientSet(ctx, user)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("unable to create atlas client: %w", err)
		}

		group, _, err := atlasClient.SdkClient20250312012.ProjectsApi.GetGroup(ctx, projectID).Execute()
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("unable to get group: %w", err)
		}
		groupName = group.GetName()
	}

	activeTargets := 0
	for _, connectionTarget := range connectionTargetInstances {
		// Construct connection secret identifier for the current connection target.
		connectionSecretIdentifier := ConnectionSecretIdentifiers{
			ProjectID:        projectID,
			ProjectName:      groupName,
			TargetName:       connectionTarget.GetName(),
			DatabaseUsername: databaseUsername,
			ConnectionType:   connectionTarget.GetConnectionTargetType(),
		}

		// Check if the user is expired.
		expired, err := timeutil.IsExpired(deleteAfterDate)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to check expiration date on user: %w", err)
		}

		//Delete secret if the user is expired.
		if expired {
			result, err := r.handleDelete(ctx, req, &connectionSecretIdentifier)
			if err != nil {
				return result, err
			}
			continue
		}

		// Check that scopes are still valid.
		if !allowsByScopes(user, connectionTarget.GetName(), connectionTarget.GetScopeType()) {
			result, err := r.handleDelete(ctx, req, &connectionSecretIdentifier)
			if err != nil {
				return result, err
			}
			continue
		}

		// Ensure connectionTarget readiness
		if !(connectionTarget.IsReady()) {
			continue
		}

		// Handle the upsert of the connection secret.
		result, err := r.handleUpsert(ctx, &connectionSecretIdentifier, user, connectionTarget)
		if err != nil {
			return result, err
		}
		activeTargets++
	}

	var connectionSecretCondition metav1.Condition
	if activeTargets > 0 {
		connectionSecretCondition = metav1.Condition{
			Type:               ConnectionSecretReady,
			Status:             metav1.ConditionTrue,
			Reason:             "Settled",
			Message:            fmt.Sprintf("%v connection secrets are ready", activeTargets),
			LastTransitionTime: metav1.Now(),
		}
	}

	changed := meta.SetStatusCondition(user.Status.Conditions, connectionSecretCondition)
	if !changed {
		return ctrl.Result{}, nil
	}

	patcher := state.NewPatcher(user).
		WithFieldOwner(FieldOwner).
		UpdateConditions([]metav1.Condition{connectionSecretCondition})

	if err := patcher.Patch(ctx, r.Client); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update status condition on error: %w", err)
	}

	return ctrl.Result{}, nil
}

func getScopes(u *generatedv1.DatabaseUser, scopeType string) []string {
	if u.Spec.V20250312 == nil || u.Spec.V20250312.Entry.Scopes == nil {
		return nil
	}

	var scopeClusters []string
	for _, scope := range *u.Spec.V20250312.Entry.Scopes {
		if scope.Type == scopeType {
			scopeClusters = append(scopeClusters, scope.Name)
		}
	}
	return scopeClusters
}

func (r *ConnectionSecretReconciler) getUserGroupId(ctx context.Context, user *generatedv1.DatabaseUser) (string, error) {
	if user == nil || (user.Spec.V20250312 == nil) || (user.Spec.V20250312.GroupRef == nil && user.Spec.V20250312.GroupId == nil) {
		return "", fmt.Errorf("cannot get project ID")
	}

	if user.Spec.V20250312.GroupId != nil {
		return *user.Spec.V20250312.GroupId, nil
	}

	group := &generatedv1.Group{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Name:      user.Spec.V20250312.GroupRef.Name,
		Namespace: user.GetNamespace(),
	}, group)
	if err != nil {
		return "", fmt.Errorf("failed to get Group: %w", err)
	}

	if group.Status.V20250312 == nil || group.Status.V20250312.Id == nil {
		return "", fmt.Errorf("group does not have a valid project ID")
	}

	return *group.Status.V20250312.Id, nil
}

func allowsByScopes(u *generatedv1.DatabaseUser, epName string, epType string) bool {
	var scopes []generatedv1.Scopes
	if u.Spec.V20250312 != nil && u.Spec.V20250312.Entry != nil && u.Spec.V20250312.Entry.Scopes != nil {
		scopes = *u.Spec.V20250312.Entry.Scopes
	}

	filtered_scopes := getScopes(u, epType)
	if len(scopes) == 0 || stringutil.Contains(filtered_scopes, epName) {
		return true
	}

	return false
}
