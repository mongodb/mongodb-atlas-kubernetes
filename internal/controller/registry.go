package controller

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasnetworkpeering"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasprivateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlassearchindexconfig"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlasstream"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
)

type ManagerAware interface {
	SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error
}

type AkoReconciler interface {
	reconcile.Reconciler
	ManagerAware
}

type Registry struct {
	predicates            []predicate.Predicate
	deletionProtection    bool
	independentSyncPeriod time.Duration
	featureFlags          *featureflags.FeatureFlags

	logger      *zap.Logger
	reconcilers []AkoReconciler
}

func NewRegistry(predicates []predicate.Predicate, deletionProtection bool, logger *zap.Logger, independentSyncPeriod time.Duration, featureFlags *featureflags.FeatureFlags) *Registry {
	return &Registry{
		predicates:            predicates,
		deletionProtection:    deletionProtection,
		logger:                logger,
		independentSyncPeriod: independentSyncPeriod,
		featureFlags:          featureFlags,
	}
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
	var reconcilers []AkoReconciler
	reconcilers = append(reconcilers, atlasproject.NewAtlasProjectReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasdeployment.NewAtlasDeploymentReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasdatabaseuser.NewAtlasDatabaseUserReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.featureFlags, r.logger))
	reconcilers = append(reconcilers, atlasdatafederation.NewAtlasDataFederationReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasfederatedauth.NewAtlasFederatedAuthReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsInstanceReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasstream.NewAtlasStreamsConnectionReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlassearchindexconfig.NewAtlasSearchIndexConfigReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlasbackupcompliancepolicy.NewAtlasBackupCompliancePolicyReconciler(c, r.predicates, ap, r.deletionProtection, r.logger))
	reconcilers = append(reconcilers, atlascustomrole.NewAtlasCustomRoleReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasprivateendpoint.NewAtlasPrivateEndpointReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasipaccesslist.NewAtlasIPAccessListReconciler(c, r.predicates, ap, r.deletionProtection, r.independentSyncPeriod, r.logger))
	reconcilers = append(reconcilers, atlasnetworkpeering.NewAtlasNetworkPeeringsReconciler(c, r.predicates, ap, r.deletionProtection, r.logger, r.independentSyncPeriod))
	r.reconcilers = reconcilers
}
