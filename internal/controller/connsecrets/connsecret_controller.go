package connsecrets

import (
	"context"
	"fmt"

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
)

type ConnSecretReconciler struct {
	reconciler.AtlasReconciler
	Scheme             *runtime.Scheme
	EventRecorder      record.EventRecorder
	GlobalPredicates   []predicate.Predicate
	EndpointStrategies []AnyEndpointStrategy
}

func (r *ConnSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Parses the request name and fills up the identifiers: ProjectID, ClusterName, DatabaseUsername
	log := r.Log.With("ns", req.Namespace, "name", req.Name)
	log.Debugw("reconcile started")

	ids, err := r.LoadIdentifiers(ctx, req.NamespacedName)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			log.Debugw("connectionsecret not found; assuming deleted")
			return workflow.TerminateSilently(nil).WithoutRetry().ReconcileResult()
		}
		log.Errorw("failed to parse connectionsecret request", "reason", workflow.ConnSecretInvalidName, "error", err)
		return workflow.Terminate(workflow.ConnSecretInvalidName, err).ReconcileResult()
	}

	var (
		pair     *ConnSecretPair[any]
		strategy AnyEndpointStrategy
	)

	// We would need to know if we use a Deployment or Federation as Endpoint
	for _, s := range r.EndpointStrategies {
		p, err := s.LoadPair(ctx, r.Client, ids)
		if err == nil {
			pair, strategy = p, s
			break
		}
		if err == ErrNoEndpointFound || err == ErrNoPairedResourcesFound {
			continue
		}

		return ctrl.Result{}, err
	}

	expired, err := timeutil.IsExpired(pair.User.Spec.DeleteAfterDate)
	if err != nil {
		return workflow.Terminate(workflow.ConnSecretCheckExpirationFailed, err).ReconcileResult()
	}
	if expired {
		return r.handleDelete(ctx, req, ids, pair, strategy)
	}

	if !strategy.ValidScopes(pair) {
		log.Infow("invalid scope; scheduling deletion", "reason", workflow.ConnSecretInvalidScopes)
		return r.handleDelete(ctx, req, ids, pair, strategy)
	}

	// Checks that AtlasDeployment and AtlasDatabaseUser are ready before proceeding
	if ready := strategy.Ready(pair); !ready {
		return workflow.InProgress(workflow.ConnSecretNotReady, "not ready").ReconcileResult()
	}

	// Create or update the k8s connection secret
	log.Infow("creating/updating connection secret", "reason", workflow.ConnSecretUpsert)
	return r.handleUpsert(ctx, req, ids, pair, strategy)
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
			handler.EnqueueRequestsFromMapFunc(r.newDeploymentMapFunc),
			builder.WithPredicates(predicate.Or(
				watch.ReadyTransitionPredicate(func(d *akov2.AtlasDeployment) bool {
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

func (r *ConnSecretReconciler) generateConnectionSecretRequests(projectID string, deployments []akov2.AtlasDeployment, users []akov2.AtlasDatabaseUser) []reconcile.Request {
	var requests []reconcile.Request
	for _, d := range deployments {
		for _, u := range users {
			scopes := u.GetScopes(akov2.DeploymentScopeType)
			if len(scopes) != 0 && !stringutil.Contains(scopes, d.GetDeploymentName()) {
				continue
			}
			name := CreateInternalFormat(projectID, d.GetDeploymentName(), u.Spec.Username)
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: u.Namespace, Name: name},
			})
		}
	}
	return requests
}

func (r *ConnSecretReconciler) ResolveProjectId(ctx context.Context, ref akov2.ProjectDualReference, parentNamespace string) (string, error) {
	if ref.ExternalProjectRef != nil && ref.ExternalProjectRef.ID != "" {
		return ref.ExternalProjectRef.ID, nil
	}
	if ref.ProjectRef != nil && ref.ProjectRef.Name != "" {
		project := &akov2.AtlasProject{}
		if err := r.Client.Get(ctx, *ref.ProjectRef.GetObject(parentNamespace), project); err != nil {
			return "", fmt.Errorf("failed to resolve projectRef: %w", err)
		}
		return project.ID(), nil
	}
	return "", fmt.Errorf("missing both external and internal project references")
}

func (r *ConnSecretReconciler) newDeploymentMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	d, ok := obj.(*akov2.AtlasDeployment)
	if !ok {
		return nil
	}
	projectID, err := r.ResolveProjectId(ctx, d.Spec.ProjectDualReference, d.GetNamespace())
	if err != nil || projectID == "" {
		return nil
	}
	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserByProject, projectID),
	}); err != nil {
		return nil
	}
	return r.generateConnectionSecretRequests(projectID, []akov2.AtlasDeployment{*d}, users.Items)
}

func (r *ConnSecretReconciler) newDatabaseUserMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	u, ok := obj.(*akov2.AtlasDatabaseUser)
	if !ok {
		return nil
	}
	projectID, err := r.ResolveProjectId(ctx, u.Spec.ProjectDualReference, u.GetNamespace())
	if err != nil || projectID == "" {
		return nil
	}
	deps := &akov2.AtlasDeploymentList{}
	if err := r.Client.List(ctx, deps, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectID),
	}); err != nil {
		return nil
	}
	return r.generateConnectionSecretRequests(projectID, deps.Items, []akov2.AtlasDatabaseUser{*u})
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

	r.EndpointStrategies = []AnyEndpointStrategy{
		NewAnyEndpointStrategy(r.NewDeploymentEndpoint()),
		// NewAnyEndpointStrategy(df),
	}

	return r
}
