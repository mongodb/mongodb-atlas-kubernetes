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

package connsecretsgeneric

import (
	"context"
	"errors"
	"fmt"

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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

type ConnSecretReconciler struct {
	reconciler.AtlasReconciler
	Scheme           *runtime.Scheme
	EventRecorder    record.EventRecorder
	GlobalPredicates []predicate.Predicate
	EndpointKinds    []Endpoint // Register all kinds of endpoints
}

type Endpoint interface {
	GetName() string
	IsReady() bool
	GetProjectRef(ctx context.Context) string
	GetProjectID(ctx context.Context) (string, error)
	GetProjectName(ctx context.Context) (string, error)

	ListObj() client.ObjectList
	ExtractList(client.ObjectList) ([]Endpoint, error)
	SelectorByProject(projectRef string) fields.Selector
	SelectorByProjectAndName(ids *ConnSecretIdentifiers) fields.Selector

	BuildConnData(ctx context.Context, user *akov2.AtlasDatabaseUser) (ConnSecretData, error)
}

type ConnSecretPair struct {
	ProjectID string
	User      *akov2.AtlasDatabaseUser
	Endpoint  Endpoint
}

func (r *ConnSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("ns", req.Namespace, "name", req.Name)

	ids, err := r.LoadIdentifiers(ctx, req.NamespacedName)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("connectionsecret not found; assuming deleted")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("failed to parse connectionsecret request", "reason", workflow.ConnSecretInvalidName, "error", err)
		return workflow.Terminate("", err).ReconcileResult()
	}

	pair, err := r.LoadPair(ctx, ids)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoPairedResourcesFound):
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		case errors.Is(err, ErrNoEndpointFound), errors.Is(err, ErrNoUserFound):
			return r.handleDelete(ctx, req, ids, pair)
		case errors.Is(err, ErrManyEndpoints), errors.Is(err, ErrManyUsers):
			return workflow.Terminate("", err).ReconcileResult()
		default:
			return workflow.Terminate("", err).ReconcileResult()
		}
	}

	expired, err := timeutil.IsExpired(pair.User.Spec.DeleteAfterDate)
	if err != nil {
		return workflow.Terminate(workflow.ConnSecretCheckExpirationFailed, err).ReconcileResult()
	}
	if expired {
		if pair.Endpoint == nil {
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		return r.handleDelete(ctx, req, ids, pair)
	}

	if pair.Endpoint == nil {
		return r.handleDelete(ctx, req, ids, pair)
	}

	if !allowsByScopes(pair.User, pair.Endpoint.GetName()) {
		r.Log.Infow("invalid scope; scheduling deletion", "reason", workflow.ConnSecretInvalidScopes)
		return r.handleDelete(ctx, req, ids, pair)
	}

	if !(pair.User.IsDatabaseUserReady() && pair.Endpoint.IsReady()) {
		return workflow.InProgress(workflow.ConnSecretNotReady, "not ready").ReconcileResult()
	}

	r.Log.Infow("creating/updating connection secret", "reason", workflow.ConnSecretUpsert)
	return r.handleUpsert(ctx, req, ids, pair)
}

func (r *ConnSecretReconciler) For() (client.Object, builder.Predicates) {
	preds := append(r.GlobalPredicates, watch.SecretLabelPredicate(TypeLabelKey, ProjectLabelKey, ClusterLabelKey))
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

func (r *ConnSecretReconciler) generateConnectionSecretRequests(projectID string, endpoints []Endpoint, users []akov2.AtlasDatabaseUser) []reconcile.Request {
	var reqs []reconcile.Request
	for _, ep := range endpoints {
		for _, u := range users {
			if !allowsByScopes(&u, ep.GetName()) {
				continue
			}
			name := CreateInternalFormat(projectID, ep.GetName(), u.Spec.Username)
			reqs = append(reqs, reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: u.Namespace, Name: name},
			})
		}
	}
	return reqs
}

func (r *ConnSecretReconciler) ResolveProjectId(ctx context.Context, ref akov2.ProjectDualReference, parentNamespace string) (string, string, error) {
	if ref.ExternalProjectRef != nil && ref.ExternalProjectRef.ID != "" {
		return "", ref.ExternalProjectRef.ID, nil
	}
	if ref.ProjectRef != nil && ref.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		if err := r.Client.Get(ctx, *ref.ProjectRef.GetObject(parentNamespace), project); err != nil {
			return "", "", fmt.Errorf("failed to resolve projectRef: %w", err)
		}
		return ref.ProjectRef.GetObject(parentNamespace).String(), project.ID(), nil
	}
	return "", "", fmt.Errorf("missing both external and internal project references")
}

func (r *ConnSecretReconciler) listEndpointsByProject(ctx context.Context, projectRef string, projectID string) ([]Endpoint, error) {
	var out []Endpoint
	for _, kind := range r.EndpointKinds {
		ref := kind.GetProjectRef(ctx)
		if ref == "PROJECTID" {
			ref = projectID // ProjectID used by deployment
		} else {
			ref = projectRef // ProjectRef used by federation
		}

		list := kind.ListObj()
		if err := r.Client.List(ctx, list, &client.ListOptions{
			FieldSelector: kind.SelectorByProject(ref),
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

func (r *ConnSecretReconciler) newEndpointMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	var ep Endpoint
	switch o := obj.(type) {
	case *akov2.AtlasDeployment:
		ep = DeploymentEndpoint{obj: o, r: r}
	case *akov2.AtlasDataFederation:
		ep = FederationEndpoint{obj: o, r: r}
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

func (r *ConnSecretReconciler) newDatabaseUserMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	u, ok := obj.(*akov2.AtlasDatabaseUser)
	if !ok {
		return nil
	}
	projectRef, projectID, err := r.ResolveProjectId(ctx, u.Spec.ProjectDualReference, u.GetNamespace())
	if err != nil {
		return nil
	}

	// The user should connect to all endpoint types
	endpoints, err := r.listEndpointsByProject(ctx, projectRef, projectID)
	if err != nil {
		return nil
	}

	return r.generateConnectionSecretRequests(projectID, endpoints, []akov2.AtlasDatabaseUser{*u})
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

	// Register kinds to try (order matters)
	r.EndpointKinds = []Endpoint{
		DeploymentEndpoint{r: r},
		FederationEndpoint{r: r},
	}

	return r
}
