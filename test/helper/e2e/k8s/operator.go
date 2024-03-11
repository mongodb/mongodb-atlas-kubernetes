package k8s

import (
	"context"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
)

func BuildManager(initCfg *Config) (manager.Manager, error) {
	scheme := runtime.NewScheme()
	setupLog := ctrl.Log.WithName("setup")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))

	ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.StacktraceLevel(zap.ErrorLevel))
	config := mergeConfiguration(initCfg)

	logger := zaptest.NewLogger(
		GinkgoT(),
		zaptest.WrapOptions(
			zap.ErrorOutput(zapcore.Lock(zapcore.AddSync(GinkgoWriter))),
		),
	)

	logger.Info("starting with configuration", zap.Any("config", config))

	ctrl.SetLogger(zapr.NewLogger(logger))

	syncPeriod := time.Hour * 3

	logger.Info("starting manager", zap.Any("config", config))

	var cacheFunc cache.NewCacheFunc
	if len(config.WatchedNamespaces) > 1 {
		var namespaces []string
		for ns := range config.WatchedNamespaces {
			namespaces = append(namespaces, ns)
		}
		cacheFunc = controller.MultiNamespacedCacheBuilder(namespaces)
	} else {
		cacheFunc = controller.CustomLabelSelectorCacheBuilder(
			&corev1.Secret{},
			labels.SelectorFromSet(labels.Set{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			}),
		)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:  scheme,
		Metrics: metricsserver.Options{BindAddress: config.MetricsAddr},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port: 9443,
		}),
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{config.Namespace: {}},
			SyncPeriod:        &syncPeriod,
		},
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "06d035fb.mongodb.com",
		NewCache:               cacheFunc,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}

	// globalPredicates should be used for general controller Predicates
	// that should be applied to all controllers in order to limit the
	// resources they receive events for.
	globalPredicates := []predicate.Predicate{
		watch.CommonPredicates(),                                  // ignore spurious changes. status changes etc.
		watch.SelectNamespacesPredicate(config.WatchedNamespaces), // select only desired namespaces
	}

	atlasProvider := atlas.NewProductionProvider(config.AtlasDomain, config.GlobalAPISecret, mgr.GetClient())

	if err = (&atlasdeployment.AtlasDeploymentReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDeployment"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDeployment")
		return nil, err
	}

	if err = (&atlasproject.AtlasProjectReconciler{
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:                      mgr.GetScheme(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               mgr.GetEventRecorderFor("AtlasProject"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasProject")
		return nil, err
	}

	if err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		ResourceWatcher:               watch.NewResourceWatcher(),
		Client:                        mgr.GetClient(),
		Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:                        mgr.GetScheme(),
		EventRecorder:                 mgr.GetEventRecorderFor("AtlasDatabaseUser"),
		AtlasProvider:                 atlasProvider,
		GlobalPredicates:              globalPredicates,
		ObjectDeletionProtection:      config.ObjectDeletionProtection,
		SubObjectDeletionProtection:   config.SubObjectDeletionProtection,
		FeaturePreviewOIDCAuthEnabled: config.FeatureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDatabaseUser")
		return nil, err
	}

	if err = (&atlasdatafederation.AtlasDataFederationReconciler{
		ResourceWatcher:             watch.NewResourceWatcher(),
		Client:                      mgr.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		Scheme:                      mgr.GetScheme(),
		EventRecorder:               mgr.GetEventRecorderFor("AtlasDataFederation"),
		AtlasProvider:               atlasProvider,
		GlobalPredicates:            globalPredicates,
		ObjectDeletionProtection:    config.ObjectDeletionProtection,
		SubObjectDeletionProtection: config.SubObjectDeletionProtection,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDataFederation")
		return nil, err
	}

	if err = mgr.AddHealthzCheck("health", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return nil, err
	}
	if err = mgr.AddReadyzCheck("check", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return nil, err
	}
	return mgr, nil
}

type Config struct {
	AtlasDomain                 string
	EnableLeaderElection        bool
	MetricsAddr                 string
	Namespace                   string
	WatchedNamespaces           map[string]bool
	ProbeAddr                   string
	GlobalAPISecret             client.ObjectKey
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	FeatureFlags                *featureflags.FeatureFlags
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func mergeConfiguration(initCfg *Config) *Config {
	config := initCfg
	if config.AtlasDomain == "" {
		config.AtlasDomain = "https://cloud-qa.mongodb.com/"
	}
	if config.MetricsAddr == "" {
		// random port
		config.MetricsAddr = ":0"
	}
	if config.ProbeAddr == "" {
		// random port
		config.ProbeAddr = ":0"
	}

	watchedNamespace := ""
	if config.WatchedNamespaces == nil {
		config.WatchedNamespaces = make(map[string]bool)
		for _, namespace := range strings.Split(watchedNamespace, ",") {
			namespace = strings.TrimSpace(namespace)
			config.WatchedNamespaces[namespace] = true
		}
	}

	if len(config.WatchedNamespaces) == 1 && config.Namespace == "" {
		config.Namespace = watchedNamespace
	}

	return config
}

type ManagerStart func(ctx context.Context) error
type ManagerConfig func(config *Config)

func managerDefaults() *Config {
	return &Config{
		AtlasDomain:                 "https://cloud-qa.mongodb.com/",
		EnableLeaderElection:        false,
		MetricsAddr:                 "0",
		Namespace:                   "mongodb-atlas-system",
		WatchedNamespaces:           map[string]bool{},
		ProbeAddr:                   "0",
		GlobalAPISecret:             client.ObjectKey{},
		ObjectDeletionProtection:    false,
		SubObjectDeletionProtection: false,
		FeatureFlags:                featureflags.NewFeatureFlags(os.Environ),
	}
}

func WithAtlasDomain(domain string) ManagerConfig {
	return func(config *Config) {
		config.AtlasDomain = domain
	}
}

func WithNamespaces(namespaces ...string) ManagerConfig {
	return func(config *Config) {
		for _, namespace := range namespaces {
			config.WatchedNamespaces[namespace] = true
		}

		if len(namespaces) == 1 {
			config.Namespace = namespaces[0]
		}
	}
}

func WithObjectDeletionProtection(flag bool) ManagerConfig {
	return func(config *Config) {
		config.ObjectDeletionProtection = flag
	}
}

func WithSubObjectDeletionProtection(flag bool) ManagerConfig {
	return func(config *Config) {
		config.SubObjectDeletionProtection = flag
	}
}

func WithGlobalKey(key client.ObjectKey) ManagerConfig {
	return func(config *Config) {
		config.GlobalAPISecret = key
	}
}

func RunManager(withConfigs ...ManagerConfig) (ManagerStart, error) {
	managerConfig := managerDefaults()

	for _, withConfig := range withConfigs {
		withConfig(managerConfig)
	}

	mgr, err := BuildManager(managerConfig)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		err = mgr.Start(ctx)
		if err != nil {
			return err
		}

		return nil
	}, nil
}
