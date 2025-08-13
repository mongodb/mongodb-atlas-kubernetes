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
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

type ConnectionSecretReconciler struct {
	reconciler.AtlasReconciler
	Scheme           *runtime.Scheme
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
}

func (r *ConnectionSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Parses the request name and fills up the identifiers: ProjectID, ClusterName, DatabaseUsername
	log := r.Log.With("ns", req.Namespace, "name", req.Name)
	log.Debugw("reconcile started")

	ids, err := r.loadRequestIdentifiers(ctx, req.NamespacedName)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("connectionsecret not found; assuming deleted")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("failed to parse connectionsecret request", "reason", workflow.ConnSecretInvalidName, "error", err)
		return workflow.Terminate(workflow.ConnSecretInvalidName, err).ReconcileResult()
	}

	log.Debugw("identifiers loaded")

	// Loads the pair of AtlasDeployment and AtlasDatabaseUser via the indexers
	pair, err := r.loadPairedResources(ctx, ids)
	if err != nil {
		switch {
		// This means there's no owner resources; the secret will be garbage collected
		case errors.Is(err, ErrNoPairedResourcesFound):
			log.Debugw("no paired resources found")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()

		// This means an owner from the pair was deleted; the secret will be forcefully removed
		case errors.Is(err, ErrNoDeploymentFound), errors.Is(err, ErrNoUserFound):
			log.Infow("paired resource missing; scheduling deletion", "reason", workflow.ConnSecretOwnerMissing)
			return r.handleDelete(ctx, req, ids, pair)

		case errors.Is(err, ErrManyDeployments), errors.Is(err, ErrManyUsers):
			log.Errorw("ambiguous pairing; multiple matches", "reason", workflow.ConnSecretAmbiguousResources, "error", err)
			return workflow.Terminate(workflow.ConnSecretAmbiguousResources, err).ReconcileResult()

		default:
			log.Errorw("failed to load paired resources", "reason", workflow.ConnSecretInvalidResources, "error", err)
			return workflow.Terminate(workflow.ConnSecretInvalidResources, err).ReconcileResult()
		}
	}

	log.Debugw("paired resources loaded")

	// If the user expired, delete connection secret
	expired, err := timeutil.IsExpired(pair.User.Spec.DeleteAfterDate)
	if err != nil {
		log.Errorw("failed to check expiration date", "reason", workflow.ConnSecretCheckExpirationFailed, "error", err)
		return workflow.Terminate(workflow.ConnSecretCheckExpirationFailed, err).ReconcileResult()
	}
	if expired {
		log.Infow("user expired; scheduling deletion", "reason", workflow.ConnSecretUserExpired)
		return r.handleDelete(ctx, req, ids, pair)
	}

	// If the scope became invalid, delete connection secret
	if invalidScopes(pair) {
		log.Infow("invalid scope; scheduling deletion", "reason", workflow.ConnSecretInvalidScopes)
		return r.handleDelete(ctx, req, ids, pair)
	}

	// Checks that AtlasDeployment and AtlasDatabaseUser are ready before proceeding
	if ready, notReady := isReady(pair); !ready {
		log.Debugw("waiting for paired resources to become ready", "notReady", strings.Join(notReady, ","))
		return workflow.InProgress(workflow.ConnSecretNotReady, fmt.Sprintf("Not ready: %s", strings.Join(notReady, ", "))).ReconcileResult()
	}

	// Create or update the k8s connection secret
	log.Infow("creating/updating connection secret", "reason", workflow.ConnSecretUpsert)
	return r.handleUpsert(ctx, req, ids, pair)
}

func (r *ConnectionSecretReconciler) For() (client.Object, builder.Predicates) {
	preds := append(
		r.GlobalPredicates,
		watch.SecretLabelPredicate(TypeLabelKey, ProjectLabelKey, ClusterLabelKey),
	)
	return &corev1.Secret{}, builder.WithPredicates(preds...)
}

func (r *ConnectionSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.newDeploymentMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate((*akov2.AtlasDeployment).IsDeploymentReady),
				predicate.GenerationChangedPredicate{},
			)),
		).
		Watches(
			&akov2.AtlasDatabaseUser{},
			handler.EnqueueRequestsFromMapFunc(r.newDatabaseUserMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate((*akov2.AtlasDatabaseUser).IsDatabaseUserReady),
				predicate.GenerationChangedPredicate{},
			)),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:        ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation: pointer.MakePtr(skipNameValidation),
		}).
		Complete(r)
}

func (r *ConnectionSecretReconciler) generateConnectionSecretRequests(
	projectID string,
	deployments []akov2.AtlasDeployment,
	users []akov2.AtlasDatabaseUser,
) []reconcile.Request {
	var requests []reconcile.Request
	for _, d := range deployments {
		for _, u := range users {
			scopes := u.GetScopes(akov2.DeploymentScopeType)
			if len(scopes) != 0 && !stringutil.Contains(scopes, d.GetDeploymentName()) {
				continue
			}

			requestName := CreateInternalFormat(projectID, d.GetDeploymentName(), u.Spec.Username)
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: u.Namespace, // connection secrets always live in the namespace of the user
					Name:      requestName,
				},
			})
		}
	}
	return requests
}

func (r *ConnectionSecretReconciler) ResolveProjectId(ctx context.Context, ref akov2.ProjectDualReference, parentNamespace string) (string, error) {
	if ref.ExternalProjectRef != nil && ref.ExternalProjectRef.ID != "" {
		return ref.ExternalProjectRef.ID, nil
	}
	if ref.ProjectRef != nil && ref.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		if err := r.Client.Get(ctx, *ref.ProjectRef.GetObject(parentNamespace), project); err != nil {
			return "", fmt.Errorf("failed to resolve projectRef from deployment: %w", err)
		}
		return project.ID(), nil
	}
	return "", fmt.Errorf("missing both external and internal project references")
}

func (r *ConnectionSecretReconciler) newDeploymentMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	deployment, ok := obj.(*akov2.AtlasDeployment)
	if !ok {
		r.Log.Warnf("watching AtlasDeployment but got %T", obj)
		return nil
	}
	projectID, err := r.ResolveProjectId(ctx, deployment.Spec.ProjectDualReference, deployment.GetNamespace())
	if err != nil {
		r.Log.Errorw("Unable to resolve projectID for deployment", "error", err)
		return nil
	}
	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserByProject, projectID),
	}); err != nil {
		r.Log.Errorf("failed to list AtlasDatabaseUsers: %v", err)
		return nil
	}
	return r.generateConnectionSecretRequests(projectID, []akov2.AtlasDeployment{*deployment}, users.Items)
}
func (r *ConnectionSecretReconciler) newDatabaseUserMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	user, ok := obj.(*akov2.AtlasDatabaseUser)
	if !ok {
		r.Log.Warnf("watching AtlasDatabaseUser but got %T", obj)
		return nil
	}
	projectID, err := r.ResolveProjectId(ctx, user.Spec.ProjectDualReference, user.GetNamespace())
	if err != nil {
		r.Log.Errorw("Unable to resolve projectID for user", "error", err)
		return nil
	}
	deployments := &akov2.AtlasDeploymentList{}
	if err := r.Client.List(ctx, deployments, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectID),
	}); err != nil {
		r.Log.Errorf("failed to list AtlasDeployments: %v", err)
		return nil
	}
	return r.generateConnectionSecretRequests(projectID, deployments.Items, []akov2.AtlasDatabaseUser{*user})
}

func NewConnectionSecretReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef types.NamespacedName,
) *ConnectionSecretReconciler {
	return &ConnectionSecretReconciler{
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
}
