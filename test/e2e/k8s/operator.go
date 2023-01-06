package k8s

import (
	"errors"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
)

func RunOperator(initCfg *Config) (manager.Manager, error) {
	scheme := runtime.NewScheme()
	setupLog := ctrl.Log.WithName("setup")

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(mdbv1.AddToScheme(scheme))

	ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.StacktraceLevel(zap.ErrorLevel))
	config := mergeConfiguration(initCfg)
	logger, err := initCustomZapLogger(config.LogLevel, config.LogEncoder, config.LogFileName)
	if err != nil {
		logger.Error("failed to initialize custom zap logger", zap.Error(err))
		return nil, err
	}

	logger.Info("starting with configuration", zap.Any("config", config))

	ctrl.SetLogger(zapr.NewLogger(logger))

	syncPeriod := time.Hour * 3

	var cacheFunc cache.NewCacheFunc
	if len(config.WatchedNamespaces) > 1 {
		var namespaces []string
		for ns := range config.WatchedNamespaces {
			namespaces = append(namespaces, ns)
		}
		cacheFunc = cache.MultiNamespacedCacheBuilder(namespaces)
	} else {
		cacheFunc = cache.BuilderWithOptions(cache.Options{
			SelectorsByObject: cache.SelectorsByObject{
				&corev1.Secret{}: {
					Label: labels.SelectorFromSet(labels.Set{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					}),
				},
			},
		})
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     config.MetricsAddr,
		Port:                   9443,
		Namespace:              config.Namespace,
		HealthProbeBindAddress: config.ProbeAddr,
		LeaderElection:         config.EnableLeaderElection,
		LeaderElectionID:       "06d035fb.mongodb.com",
		SyncPeriod:             &syncPeriod,
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

	if err = (&atlasdeployment.AtlasDeploymentReconciler{
		Client:           mgr.GetClient(),
		Log:              logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		Scheme:           mgr.GetScheme(),
		AtlasDomain:      config.AtlasDomain,
		GlobalAPISecret:  config.GlobalAPISecret,
		ResourceWatcher:  watch.NewResourceWatcher(),
		GlobalPredicates: globalPredicates,
		EventRecorder:    mgr.GetEventRecorderFor("AtlasDeployment"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDeployment")
		return nil, err
	}

	if err = (&atlasproject.AtlasProjectReconciler{
		Client:           mgr.GetClient(),
		Log:              logger.Named("controllers").Named("AtlasProject").Sugar(),
		Scheme:           mgr.GetScheme(),
		AtlasDomain:      config.AtlasDomain,
		ResourceWatcher:  watch.NewResourceWatcher(),
		GlobalAPISecret:  config.GlobalAPISecret,
		GlobalPredicates: globalPredicates,
		EventRecorder:    mgr.GetEventRecorderFor("AtlasProject"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasProject")
		return nil, err
	}

	if err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		Client:           mgr.GetClient(),
		Log:              logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		Scheme:           mgr.GetScheme(),
		AtlasDomain:      config.AtlasDomain,
		ResourceWatcher:  watch.NewResourceWatcher(),
		GlobalAPISecret:  config.GlobalAPISecret,
		GlobalPredicates: globalPredicates,
		EventRecorder:    mgr.GetEventRecorderFor("AtlasDatabaseUser"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AtlasDatabaseUser")
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
	AtlasDomain          string
	EnableLeaderElection bool
	MetricsAddr          string
	Namespace            string
	WatchedNamespaces    map[string]bool
	ProbeAddr            string
	GlobalAPISecret      client.ObjectKey
	LogLevel             string
	LogEncoder           string
	LogFileName          string
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
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}
	if config.LogEncoder == "" {
		config.LogEncoder = "json"
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

func initCustomZapLogger(level, encoding, logFileName string) (*zap.Logger, error) {
	lv := zap.AtomicLevel{}
	err := lv.UnmarshalText([]byte(strings.ToLower(level)))
	if err != nil {
		return nil, err
	}

	enc := strings.ToLower(encoding)
	if enc != "json" && enc != "console" {
		return nil, errors.New("'encoding' parameter can only by either 'json' or 'console'")
	}

	cfg := zap.Config{
		Level:             lv,
		OutputPaths:       []string{logFileName},
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          enc,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	return cfg.Build()
}
