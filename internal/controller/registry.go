package controller

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasbackupcompliancepolicy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlascustomrole"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasipaccesslist"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasprivateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlassearchindexconfig"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasstream"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
)

type Reconciler interface {
	reconcile.Reconciler
	For() (client.Object, builder.Predicates)
	SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error
}

type Registry struct {
	sharedPredicates      []predicate.Predicate
	deletionProtection    bool
	independentSyncPeriod time.Duration
	featureFlags          *featureflags.FeatureFlags

	logger      *zap.Logger
	reconcilers []Reconciler
}

func NewRegistry(predicates []predicate.Predicate, deletionProtection bool, logger *zap.Logger, independentSyncPeriod time.Duration, featureFlags *featureflags.FeatureFlags) *Registry {
	return &Registry{
		sharedPredicates:      predicates,
		deletionProtection:    deletionProtection,
		logger:                logger,
		independentSyncPeriod: independentSyncPeriod,
		featureFlags:          featureFlags,
	}
}

func (r *Registry) RegisterWithDryRunManager(mgr *dryrun.Manager, ap atlas.Provider) error {
	r.registerControllers(mgr, ap)

	for _, reconciler := range r.reconcilers {
		mgr.SetupReconciler(reconciler)
	}

	return nil
}

func (r *Registry) RegisterWithManager(mgr ctrl.Manager, skipNameValidation bool, ap atlas.Provider) error {
	r.registerControllers(mgr, ap)

	for _, reconciler := range r.reconcilers {
		if err := reconciler.SetupWithManager(mgr, skipNameValidation); err != nil {
			return fmt.Errorf("failed to set up with manager: %w", err)
		}
	}
	return nil
}

func (r *Registry) registerControllers(c cluster.Cluster, ap atlas.Provider) {
	if len(r.reconcilers) > 0 {
		return
	}

	var reconcilers []Reconciler
	reconcilers = append(reconcilers, atlasproject.NewAtlasProjectReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasdeployment.NewAtlasDeploymentReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasdatabaseuser.NewAtlasDatabaseUserReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.featureFlags, r.logger))
	reconcilers = append(reconcilers, atlasdatafederation.NewAtlasDataFederationReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasfederatedauth.NewAtlasFederatedAuthReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsInstanceReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsConnectionReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlassearchindexconfig.NewAtlasSearchIndexConfigReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasbackupcompliancepolicy.NewAtlasBackupCompliancePolicyReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlascustomrole.NewAtlasCustomRoleReconciler(c, r.deprecatedPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasprivateendpoint.NewAtlasPrivateEndpointReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasipaccesslist.NewAtlasIPAccessListReconciler(c, r.defaultPredicates(), ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	r.reconcilers = reconcilers
}

// deprecatedPredicates are to be phased out in favor of defaultPredicates
func (r *Registry) deprecatedPredicates() []predicate.Predicate {
	return append(r.sharedPredicates, watch.DeprecatedCommonPredicates())
}

// defaultPredicates minimize the reconciliations controllers actually do, avoiding
// spurious after delete handling and acting on finalizers setting or unsetting
func (r *Registry) defaultPredicates() []predicate.Predicate {
	return append(r.sharedPredicates, watch.DefaultPredicates[client.Object]())
}
