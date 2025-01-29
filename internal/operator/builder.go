package operator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
)

const (
	DefaultAtlasDomain           = "https://cloud.mongodb.com/"
	DefaultSyncPeriod            = 3 * time.Hour
	DefaultIndependentSyncPeriod = 15 * time.Minute
	DefaultLeaderElectionID      = "06d035fb.mongodb.com"
)

type ManagerProvider interface {
	New(config *rest.Config, options manager.Options) (manager.Manager, error)
}

type ManagerProviderFunc func(config *rest.Config, options manager.Options) (manager.Manager, error)

func (f ManagerProviderFunc) New(config *rest.Config, options manager.Options) (manager.Manager, error) {
	return f(config, options)
}

type Builder struct {
	managerProvider              ManagerProvider
	scheme                       *runtime.Scheme
	minimumIndependentSyncPeriod time.Duration

	config                *rest.Config
	namespaces            []string
	logger                *zap.Logger
	syncPeriod            time.Duration
	independentSyncPeriod time.Duration
	metricAddress         string
	probeAddress          string
	leaderElection        bool
	leaderElectionID      string

	atlasDomain        string
	predicates         []predicate.Predicate
	apiSecret          client.ObjectKey
	atlasProvider      atlas.Provider
	featureFlags       *featureflags.FeatureFlags
	deletionProtection bool
	skipNameValidation bool
}

func (b *Builder) WithConfig(config *rest.Config) *Builder {
	b.config = config
	return b
}

func (b *Builder) WithNamespaces(namespaces ...string) *Builder {
	b.namespaces = namespaces
	return b
}

func (b *Builder) WithLogger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) WithSyncPeriod(period time.Duration) *Builder {
	b.syncPeriod = period
	return b
}

func (b *Builder) WithMetricAddress(address string) *Builder {
	b.metricAddress = address
	return b
}

func (b *Builder) WithProbeAddress(address string) *Builder {
	b.probeAddress = address
	return b
}

func (b *Builder) WithLeaderElection(enable bool) *Builder {
	b.leaderElection = enable
	return b
}

func (b *Builder) WithAtlasDomain(domain string) *Builder {
	b.atlasDomain = domain
	return b
}

func (b *Builder) WithPredicates(predicates []predicate.Predicate) *Builder {
	b.predicates = predicates
	return b
}

func (b *Builder) WithAPISecret(apiSecret client.ObjectKey) *Builder {
	b.apiSecret = apiSecret
	return b
}

func (b *Builder) WithAtlasProvider(provider atlas.Provider) *Builder {
	b.atlasProvider = provider
	return b
}

func (b *Builder) WithFeatureFlags(featureFlags *featureflags.FeatureFlags) *Builder {
	b.featureFlags = featureFlags
	return b
}

func (b *Builder) WithDeletionProtection(deletionProtection bool) *Builder {
	b.deletionProtection = deletionProtection
	return b
}

func (b *Builder) WithIndependentSyncPeriod(period time.Duration) *Builder {
	b.independentSyncPeriod = period
	return b
}

// WithSkipNameValidation skips name validation in controller-runtime
// to prevent duplicate controller names.
//
// Note: use this in tests only, setting this to true in a production setup will cause faulty behavior.
func (b *Builder) WithSkipNameValidation(skip bool) *Builder {
	b.skipNameValidation = skip
	return b
}

// Build builds the cluster object and configures operator controllers
func (b *Builder) Build(ctx context.Context) (cluster.Cluster, error) {
	mergeDefaults(b)

	if b.independentSyncPeriod < b.minimumIndependentSyncPeriod {
		return nil, errors.New("wrong value for independentSyncPeriod. Value should be greater or equal to 5")
	}

	cacheOpts := cache.Options{
		SyncPeriod: &b.syncPeriod,
	}

	if len(b.namespaces) == 0 {
		cacheOpts.ByObject = map[client.Object]cache.ByObject{
			&corev1.Secret{}: {
				Label: labels.SelectorFromSet(labels.Set{
					connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
				}),
			},
		}
	} else {
		cacheOpts.DefaultNamespaces = map[string]cache.Config{}
		for _, namespace := range b.namespaces {
			cacheOpts.DefaultNamespaces[namespace] = cache.Config{}
		}
	}

	controllerRegistry := controller.NewRegistry(b.predicates, b.deletionProtection, b.logger, b.independentSyncPeriod, b.featureFlags)

	mgr, err := b.managerProvider.New(
		b.config,
		ctrl.Options{
			Scheme:  b.scheme,
			Metrics: metricsserver.Options{BindAddress: b.metricAddress},
			WebhookServer: webhook.NewServer(webhook.Options{
				Port: 9443,
			}),
			Cache:                  cacheOpts,
			HealthProbeBindAddress: b.probeAddress,
			LeaderElection:         b.leaderElection,
			LeaderElectionID:       b.leaderElectionID,
		},
	)

	if err != nil {
		return nil, err
	}

	if err = mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		return nil, err
	}

	if err = mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		return nil, err
	}

	if b.atlasProvider == nil {
		b.atlasProvider = atlas.NewProductionProvider(b.atlasDomain, b.apiSecret, mgr.GetClient(), false)
	}

	if err := controllerRegistry.RegisterWithManager(mgr, b.skipNameValidation, b.atlasProvider); err != nil {
		return nil, err
	}

	if err := indexer.RegisterAll(ctx, mgr, b.logger); err != nil {
		return nil, fmt.Errorf("unable to create indexers: %w", err)
	}

	return mgr, nil
}

// NewBuilder return a new Builder to construct operator controllers
func NewBuilder(provider ManagerProvider, scheme *runtime.Scheme, minimumIndependentSyncPeriod time.Duration) *Builder {
	return &Builder{
		managerProvider:              provider,
		scheme:                       scheme,
		minimumIndependentSyncPeriod: minimumIndependentSyncPeriod,
	}
}

func mergeDefaults(b *Builder) {
	if b.config == nil {
		b.config = &rest.Config{}
	}

	if b.logger == nil {
		b.logger = ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.StacktraceLevel(zap.ErrorLevel))
	}

	if b.syncPeriod == 0 {
		b.syncPeriod = DefaultSyncPeriod
	}

	if b.independentSyncPeriod == 0 {
		b.independentSyncPeriod = DefaultIndependentSyncPeriod
	}

	if b.metricAddress == "" {
		b.metricAddress = "0"
	}

	if b.probeAddress == "" {
		b.probeAddress = "0"
	}

	if b.leaderElection {
		b.leaderElectionID = DefaultLeaderElectionID
	}

	if b.atlasDomain == "" {
		b.atlasDomain = DefaultAtlasDomain
	}

	if len(b.predicates) == 0 {
		b.predicates = []predicate.Predicate{
			watch.CommonPredicates(),
			watch.SelectNamespacesPredicate(b.namespaces),
		}
	}

	if b.featureFlags == nil {
		b.featureFlags = featureflags.NewFeatureFlags(os.Environ)
	}
}
