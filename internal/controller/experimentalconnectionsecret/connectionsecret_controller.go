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

package experimentalconnectionsecret

import (
	"context"
	"errors"

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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
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

type ConnSecretReconciler struct {
	reconciler.AtlasReconciler
	Scheme           *runtime.Scheme
	EventRecorder    record.EventRecorder
	GlobalPredicates []predicate.Predicate
	EndpointKinds    []Endpoint // Endpoints are generic
}

// Each endpoint would have to implement this interface (e.g. AtlasDeployment, AtlasDataFederation)
type Endpoint interface {
	GetName() string
	IsReady() bool
	GetScopeType() akov2.ScopeType
	GetProjectID(ctx context.Context) (string, error)
	GetConnectionType() string

	ListObj() client.ObjectList
	ExtractList(client.ObjectList) ([]Endpoint, error)
	SelectorByProject(projectID string) fields.Selector
	SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector

	BuildConnData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnSecretData, error)
}

// Each connection secret needs a paired resource: User and Endpoint
type ConnSecretPair struct {
	ProjectID string
	User      *akov2.AtlasDatabaseUser
	Endpoint  Endpoint
}

func (r *ConnSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("namespace", req.Namespace, "name", req.Name)
	log.Info("reconciliation started")

	// Parse the request and load up the identifiers
	ids, err := r.loadIdentifiers(ctx, req.NamespacedName)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("Connection secret not found; assuming deleted")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("failed to parse connection secret request", "error", err)
		return workflow.Terminate(workflow.ConnSecretInvalidName, err).ReconcileResult()
	}

	// Load the paired resource
	pair, err := r.loadPair(ctx, ids)
	if err != nil {
		switch {
		case errors.Is(err, ErrMissingPairing):
			log.Debugw("paired resource is missing; scheduling deletion of connection secrets")
			return r.handleDelete(ctx, req, ids)
		case errors.Is(err, ErrAmbiguousPairing):
			log.Errorw("failed to load paired resources; ambigous parent resources", "error", err)
			return workflow.Terminate(workflow.ConnSecretPairNotLoaded, err).ReconcileResult()
		default:
			log.Errorw("failed to load paired resource", "error", err)
			return workflow.Terminate(workflow.ConnSecretPairNotLoaded, err).ReconcileResult()
		}
	}

	// Check if user is expired
	expired, err := timeutil.IsExpired(pair.User.Spec.DeleteAfterDate)
	if err != nil {
		log.Errorw("failed to check expiration date on user", "error", err)
		return workflow.Terminate(workflow.ConnSecretUserExpired, err).ReconcileResult()
	}
	if expired {
		log.Debugw("user is expired; scheduling deletion of connection secrets")
		return r.handleDelete(ctx, req, ids)
	}

	// Check that scopes are still valid
	if !allowsByScopes(pair.User, pair.Endpoint.GetName(), pair.Endpoint.GetScopeType()) {
		log.Infow("invalid scope; scheduling deletion of connection secrets")
		return r.handleDelete(ctx, req, ids)
	}

	// Paired resource must be ready
	if !(pair.User.IsDatabaseUserReady() && pair.Endpoint.IsReady()) {
		log.Debugw("waiting on paired resource to be ready")
		return workflow.InProgress(workflow.ConnSecretNotReady, "resources not ready").ReconcileResult()
	}

	return r.handleUpsert(ctx, req, ids, pair)
}

func (r *ConnSecretReconciler) For() (client.Object, builder.Predicates) {
	preds := append(r.GlobalPredicates, watch.SecretLabelPredicate(TypeLabelKey, ProjectLabelKey, ClusterLabelKey, DatabaseUserLabelKey))
	return &corev1.Secret{}, builder.WithPredicates(preds...)
}

func (r *ConnSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.newEndpointMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(d *akov2.AtlasDeployment) bool {
					return api.HasReadyCondition(d.Status.Conditions)
				}),
				predicate.GenerationChangedPredicate{},
			)),
		).
		Watches(
			&akov2.AtlasDataFederation{},
			handler.EnqueueRequestsFromMapFunc(r.newEndpointMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(d *akov2.AtlasDataFederation) bool {
					return api.HasReadyCondition(d.Status.Conditions)
				}),
				predicate.GenerationChangedPredicate{},
			)),
		).
		Watches(
			&akov2.AtlasDatabaseUser{},
			handler.EnqueueRequestsFromMapFunc(r.newDatabaseUserMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(u *akov2.AtlasDatabaseUser) bool {
					return api.HasReadyCondition(u.Status.Conditions)
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

func (r *ConnSecretReconciler) generateConnectionSecretRequests(projectID string, endpoints []Endpoint, users []akov2.AtlasDatabaseUser) []reconcile.Request {
	var reqs []reconcile.Request
	for _, ep := range endpoints {
		for _, u := range users {
			if !allowsByScopes(&u, ep.GetName(), ep.GetScopeType()) {
				continue
			}

			name := CreateInternalFormat(projectID, ep.GetName(), u.Spec.Username, ep.GetConnectionType())
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: u.Namespace, Name: name},
			})
		}
	}
	return reqs
}

// TODO: create indexers for DataFederation by projectID

// listEndpointsByProject retrives all of the Endpoints that live under an AtlasProject
func (r *ConnSecretReconciler) listEndpointsByProject(ctx context.Context, projectID string) ([]Endpoint, error) {
	var out []Endpoint
	for _, kind := range r.EndpointKinds {
		list := kind.ListObj()
		if err := r.Client.List(ctx, list, &client.ListOptions{
			FieldSelector: kind.SelectorByProject(projectID),
		}); err != nil {
			return nil, err
		}

		eps, err := kind.ExtractList(list)
		if err != nil {
			return nil, err
		}

		out = append(out, eps...)
	}

	return out, nil
}

// newEndpointMapFunc maps an Endpoint to requests by fetching all AtlasDatabaseUsers and creating a request for each
func (r *ConnSecretReconciler) newEndpointMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	var ep Endpoint

	// Case on the type of endpoint
	switch o := obj.(type) {
	case *akov2.AtlasDeployment:
		ep = DeploymentEndpoint{
			obj: o, k8s: r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		}
	case *akov2.AtlasDataFederation:
		ep = FederationEndpoint{
			obj:             o,
			k8s:             r.Client,
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

	return r.generateConnectionSecretRequests(projectID, []Endpoint{ep}, users.Items)
}

// newDatabaseUserMapFunc maps an AtlasDatabaseUser to requests by fetching all endpoints and creating a request for each
func (r *ConnSecretReconciler) newDatabaseUserMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	user, ok := obj.(*akov2.AtlasDatabaseUser)
	if !ok {
		return nil
	}
	projectID, err := r.getUserProjectID(ctx, user)
	if err != nil {
		return nil
	}

	endpoints, err := r.listEndpointsByProject(ctx, projectID)
	if err != nil {
		return nil
	}

	return r.generateConnectionSecretRequests(projectID, endpoints, []akov2.AtlasDatabaseUser{*user})
}

func NewConnectionSecretReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef types.NamespacedName,
) *ConnSecretReconciler {
	r := &ConnSecretReconciler{
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

	// Register all the endpoint types
	r.EndpointKinds = []Endpoint{
		DeploymentEndpoint{
			k8s:             r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		},
		FederationEndpoint{
			k8s:             r.Client,
			provider:        r.AtlasProvider,
			globalSecretRef: r.GlobalSecretRef,
			log:             r.Log,
		},
	}

	return r
}
