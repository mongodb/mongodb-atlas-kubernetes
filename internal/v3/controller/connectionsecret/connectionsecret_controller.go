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
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

type ConnectionSecretReconciler struct {
	reconciler.AtlasReconciler
	Scheme                *runtime.Scheme
	EventRecorder         record.EventRecorder
	GlobalPredicates      []predicate.Predicate
	ConnectionTargetKinds []ConnectionTarget
}

// Each connectionTarget would have to implement this interface (e.g. AtlasDeployment, AtlasDataFederation)
type ConnectionTarget interface {
	GetConnectionTargetType() string
	GetName() string
	IsReady() bool
	GetScopeType() akov2.ScopeType
	GetOwnerReferences() []metav1.OwnerReference
	GetProjectID(ctx context.Context) (string, error)
	SelectorByProjectID(projectID string) fields.Selector
	BuildConnectionData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnectionSecretData, error)
}

func (r *ConnectionSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("namespace", req.Namespace, "name", req.Name)

	// Fetch the AtlasDatabaseUser resource.
	user := &akov2.AtlasDatabaseUser{}
	err := r.Client.Get(ctx, req.NamespacedName, user)
	objectNotFound := err != nil && apiErrors.IsNotFound(err)
	failedToRetrieve := err != nil && !objectNotFound

	switch {
	case failedToRetrieve:
		return workflow.Terminate(workflow.ConnectionSecretInvalidUsername, err).ReconcileResult()
	case objectNotFound:
		log.Debugw("user not found; nothing to reconcile")
		return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
	}

	// Retrieve the project ID associated with the user.
	projectID, err := r.getUserProjectID(ctx, user)
	if err != nil {
		return workflow.Terminate(workflow.ConnectionSecretProjectIDNotLoaded, err).ReconcileResult()
	}

	// Load the connection targets for the project.
	connectionTargets, err := r.listConnectionTargetsByProject(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ConnectionSecretConnectionTargetsNotLoaded, err).ReconcileResult()
	}

	// Cleanup stale Secrets
	if err := r.cleanupStaleSecrets(ctx, req.Namespace, connectionTargets, projectID); err != nil {
		return workflow.Terminate(workflow.ConnectionSecretStaleSecretsNotCleaned, err).ReconcileResult()
	}

	// Verify if the AtlasDatabaseUser is ready.
	isUserReady := api.HasReadyCondition(user.Status.Conditions)
	if !isUserReady {
		log.Debugw("AtlasDatabaseUser not ready; nothing to reconcile")
		return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
	}

	// Delegate the batch upsert logic to handleBatchUpsert.
	return r.handleBatchUpsert(ctx, req, user, projectID, connectionTargets)
}

func (r *ConnectionSecretReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasDatabaseUser{}, builder.WithPredicates(
		predicate.Or(
			watch.ReadyTransitionPredicate(func(d *akov2.AtlasDatabaseUser) bool {
				return api.HasReadyCondition(d.Status.Conditions)
			}),
			predicate.GenerationChangedPredicate{},
		),
	)
}

func (r *ConnectionSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Owns(&corev1.Secret{}).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.newConnectionTargetMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(d *akov2.AtlasDeployment) bool {
					return api.HasReadyCondition(d.Status.Conditions)
				}),
				predicate.GenerationChangedPredicate{},
			)),
		).
		Watches(
			&akov2.AtlasDataFederation{},
			handler.EnqueueRequestsFromMapFunc(r.newConnectionTargetMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(d *akov2.AtlasDataFederation) bool {
					return api.HasReadyCondition(d.Status.Conditions)
				}),
				predicate.GenerationChangedPredicate{},
			)),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:        ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation: pointer.MakePtr(skipNameValidation),
		}).
		Complete(r)
}

func allowsByScopes(u *akov2.AtlasDatabaseUser, epName string, epType akov2.ScopeType) bool {
	scopes := u.Spec.Scopes
	filtered_scopes := u.GetScopes(epType)
	if len(scopes) == 0 || stringutil.Contains(filtered_scopes, epName) {
		return true
	}

	return false
}

func (r *ConnectionSecretReconciler) generateConnectionSecretRequests(users []akov2.AtlasDatabaseUser) []reconcile.Request {
	reqs := make([]reconcile.Request, 0, len(users))
	for _, u := range users {
		reqs = append(reqs, reconcile.Request{
			NamespacedName: types.NamespacedName{Namespace: u.Namespace, Name: u.Name},
		})
	}
	return reqs
}

// listConnectionTargetsByProject retrieves all of the connectionTargets that live under an AtlasProject
func (r *ConnectionSecretReconciler) listConnectionTargetsByProject(ctx context.Context, projectID string) ([]ConnectionTarget, error) {
	var out []ConnectionTarget

	for _, kind := range r.ConnectionTargetKinds {
		switch kind.(type) {
		case DataFederationConnectionTarget:
			list := &akov2.AtlasDataFederationList{}
			if err := r.Client.List(ctx, list, &client.ListOptions{
				FieldSelector: kind.SelectorByProjectID(projectID),
			}); err != nil {
				return nil, err
			}

			for i := range list.Items {
				out = append(out, DataFederationConnectionTarget{
					obj:             &list.Items[i],
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				})
			}

		case DeploymentConnectionTarget:
			list := &akov2.AtlasDeploymentList{}
			if err := r.Client.List(ctx, list, &client.ListOptions{
				FieldSelector: kind.SelectorByProjectID(projectID),
			}); err != nil {
				return nil, err
			}

			for i := range list.Items {
				out = append(out, DeploymentConnectionTarget{
					obj:             &list.Items[i],
					client:          r.Client,
					provider:        r.AtlasProvider,
					globalSecretRef: r.GlobalSecretRef,
					log:             r.Log,
				})
			}
		}
	}

	return out, nil
}

func (r *ConnectionSecretReconciler) cleanupStaleSecrets(ctx context.Context, namespace string, connectionTargets []ConnectionTarget, projectID string) error {
	log := r.Log.With("namespace", namespace)

	// Define the label selector to find relevant secrets.
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{Key: TypeLabelKey, Operator: metav1.LabelSelectorOpExists},
			{Key: ProjectLabelKey, Operator: metav1.LabelSelectorOpIn, Values: []string{projectID}},
			{Key: TargetLabelKey, Operator: metav1.LabelSelectorOpExists},
			{Key: DatabaseUserLabelKey, Operator: metav1.LabelSelectorOpExists},
		},
	}

	// Convert the label selector into a client-compatible format.
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return fmt.Errorf("failed to convert label selector: %w", err)
	}

	// Fetch all secrets in the specified namespace that match the label selector.
	secretList := &corev1.SecretList{}
	if err := r.Client.List(ctx, secretList, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	// Iterate through secrets and delete any that are stale.
	for _, secret := range secretList.Items {
		if err := r.checkAndDeleteStaleSecret(ctx, &secret, connectionTargets); err != nil {
			log.Errorw("Error cleaning up stale secret", "secretName", secret.Name, "error", err)
			return err
		}
	}

	return nil
}

// checkAndDeleteStaleSecret deletes a secret if its associated user or connected resource no longer exists.
func (r *ConnectionSecretReconciler) checkAndDeleteStaleSecret(ctx context.Context, secret *corev1.Secret, connectionTargets []ConnectionTarget) error {
	pendingDeletion := true
	for _, connectionTarget := range connectionTargets {
		if connectionTarget.GetName() == secret.Labels[TargetLabelKey] && connectionTarget.GetConnectionTargetType() == secret.Annotations[ConnectionTypelKey] {
			pendingDeletion = false
		}
	}

	if pendingDeletion {
		if err := r.Client.Delete(ctx, secret); err != nil {
			if apiErrors.IsNotFound(err) {
				r.Log.Debugw("Secret already deleted", "secretName", secret.Name)
				return nil
			}
			return fmt.Errorf("failed to delete secret: %w", err)
		}
	}

	return nil
}

// newConnectionTargetMapFunc maps a ConnectionTarget to requests by fetching all AtlasDatabaseUsers and creating a request for each
func (r *ConnectionSecretReconciler) newConnectionTargetMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	var ep ConnectionTarget

	// Case on the type of connectionTarget
	switch o := obj.(type) {
	case *akov2.AtlasDeployment:
		ep = DeploymentConnectionTarget{
			obj:             o,
			client:          r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		}
	case *akov2.AtlasDataFederation:
		ep = DataFederationConnectionTarget{
			obj:             o,
			client:          r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		}
	default:
		return nil
	}

	projectID, err := ep.GetProjectID(ctx)
	if err != nil || projectID == "" {
		return nil
	}

	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserByProject, projectID),
	}); err != nil {
		return nil
	}

	return r.generateConnectionSecretRequests(users.Items)
}

func NewConnectionSecretReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef types.NamespacedName,
) *ConnectionSecretReconciler {
	r := &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             logger.Named("controllers").Named("ConnectionSecret").Sugar(),
			GlobalSecretRef: globalSecretRef,
			AtlasProvider:   atlasProvider,
		},
		Scheme:           c.GetScheme(),
		EventRecorder:    c.GetEventRecorderFor("ConnectionSecret"),
		GlobalPredicates: predicates,
	}

	// Register all the connectionTarget types
	r.ConnectionTargetKinds = []ConnectionTarget{
		DeploymentConnectionTarget{
			client:          r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		},
		DataFederationConnectionTarget{
			client:          r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		},
	}

	return r
}

func (r *ConnectionSecretReconciler) getUserProjectID(ctx context.Context, user *akov2.AtlasDatabaseUser) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil user")
	}
	if user.Spec.ExternalProjectRef != nil && user.Spec.ExternalProjectRef.ID != "" {
		return user.Spec.ExternalProjectRef.ID, nil
	}
	return resolveProjectIDByKey(ctx, r.Client, user.AtlasProjectObjectKey())
}
