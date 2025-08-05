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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
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
	strRequest := req.NamespacedName.String()
	r.Log.Infof("Reconcile started for ConnectionSecret request with %s", strRequest)

	ids, err := LoadRequestIdentifiers(ctx, r.Client, req.NamespacedName)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			r.Log.Debugf("ConnectionSecret not found; assuming it was deleted %s", strRequest)
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}

		r.Log.Errorf("Failed to parse ConnectionSecret request with %s: %v", strRequest, err)
		return workflow.Terminate("InvalidConnectionSecretName", err).ReconcileResult()
	}

	r.Log.Debugf("Identifiers loaded for ConnectionSecret request with %s", strRequest)

	// Loads the pair of AtlasDeployment and AtlasDatabaseUser via the indexers
	pair, err := LoadPairedResources(ctx, r.Client, ids, req.Namespace)
	if err != nil {
		switch {
		// This means there's no owner resources; the secret will be garbage collected
		case errors.Is(err, ErrNoPairedResourcesFound):
			r.Log.Debugf("No paired resources for ConnectionSecret request with %s", strRequest)
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()

		// This means an owner from the pair was deleted; the secret will be forcefully removed
		case errors.Is(err, ErrNoDeploymentFound), errors.Is(err, ErrNoUserFound):
			r.Log.Infof("Paired resource missing for ConnectionSecret request with %s — scheduling deletion", strRequest)
			return r.handleDelete(ctx, req, ids, pair)

		case errors.Is(err, ErrManyDeployments), errors.Is(err, ErrManyUsers):
			r.Log.Errorf("Ambiguous pairing (more than one) for ConnectionSecret request with %s", strRequest)
			return workflow.Terminate("AmbiguousConnectionResources", err).ReconcileResult()

		default:
			r.Log.Errorf("Failed to get paired resources ConnectionSecret request with %s: %v", strRequest, err)
			return workflow.Terminate("InvalidConnectionResources", err).ReconcileResult()
		}
	}

	r.Log.Debugf("Paired resource loaded for ConnectionSecret request with %s", strRequest)

	// If the user expired, delete connection secret
	expired, err := atlasdatabaseuser.IsExpired(pair.User)
	if err != nil {
		r.Log.Errorf("Failed to check expiration date for ConnectionSecret request with %s", strRequest)
		return workflow.Terminate("AmbiguousConnectionResources", err).ReconcileResult()
	}
	if expired {
		r.Log.Infof("Expired user for paired resource for ConnectionSecret request with %s — scheduling deletion", strRequest)
		return r.handleDelete(ctx, req, ids, pair)
	}

	// If the scope became invalid, delete connection secret
	if pair.InvalidScopes() {
		r.Log.Infof("Invalid scope for paired resource for ConnectionSecret request with %s — scheduling deletion", strRequest)
		return r.handleDelete(ctx, req, ids, pair)
	}

	// Checks that AtlasDeployment and AtlasDatabaseUser are ready before proceeding
	if ready, notReady := pair.IsReady(); !ready {
		r.Log.Debugf("Waiting till paired resources are ready for ConnectionSecret request with %s", strRequest)
		return workflow.InProgress("ConnectionSecretNotReady", fmt.Sprintf("Not ready: %s", strings.Join(notReady, ", "))).ReconcileResult()
	}

	// Create or update the k8s connection secret
	r.Log.Infof("Start create or update ConnectionSecret request with %s", strRequest)
	return r.handleUpdate(ctx, req, ids, pair)
}

func (r *ConnectionSecretReconciler) DeploymentWatcherPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc: func(e event.UpdateEvent) bool {
			newObj, ok := e.ObjectNew.(*akov2.AtlasDeployment)
			if !ok {
				return false
			}
			oldObj, ok := e.ObjectOld.(*akov2.AtlasDeployment)
			if !ok {
				return false
			}
			return !IsDeploymentReady(oldObj) && IsDeploymentReady(newObj)
		},
	}
}

func (r *ConnectionSecretReconciler) DatabaseUserWatcherPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc: func(e event.UpdateEvent) bool {
			newObj, ok := e.ObjectNew.(*akov2.AtlasDatabaseUser)
			if !ok {
				return false
			}
			oldObj, ok := e.ObjectOld.(*akov2.AtlasDatabaseUser)
			if !ok {
				return false
			}
			return !IsDatabaseUserReady(oldObj) && IsDatabaseUserReady(newObj)
		},
	}
}

func (r *ConnectionSecretReconciler) For() (client.Object, builder.Predicates) {
	// Filter out connection secrets based on the required labels
	labelPredicates := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		labels := obj.GetLabels()
		_, hasType := labels[TypeLabelKey]
		_, hasProject := labels[ProjectLabelKey]
		_, hasCluster := labels[ClusterLabelKey]
		return hasType && hasProject && hasCluster
	})

	predicates := append(r.GlobalPredicates, labelPredicates)
	return &corev1.Secret{}, builder.WithPredicates(predicates...)
}

func (r *ConnectionSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.newDeploymentMapFunc),
			builder.WithPredicates(predicate.Or(
				r.DeploymentWatcherPredicate(),
				predicate.GenerationChangedPredicate{},
			)),
		).
		Watches(
			&akov2.AtlasDatabaseUser{},
			handler.EnqueueRequestsFromMapFunc(r.newDatabaseUserMapFunc),
			builder.WithPredicates(predicate.Or(
				r.DatabaseUserWatcherPredicate(),
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

func (r *ConnectionSecretReconciler) newDeploymentMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	deployment, ok := obj.(*akov2.AtlasDeployment)
	if !ok {
		r.Log.Warnf("watching AtlasDeployment but got %T", obj)
		return nil
	}

	projectID, err := ResolveProjectIDFromDeployment(ctx, r.Client, deployment)
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

	projectID, err := ResolveProjectIDFromDatabaseUser(ctx, r.Client, user)
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
