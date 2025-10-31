// Copyright 2020 The Kubernetes authors.
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

package run

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"
	"k8s.io/utils/env"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

const (
	objectDeletionProtectionFlag       = "object-deletion-protection"
	subobjectDeletionProtectionFlag    = "subobject-deletion-protection"
	objectDeletionProtectionEnvVar     = "OBJECT_DELETION_PROTECTION"
	subobjectDeletionProtectionEnvVar  = "SUBOBJECT_DELETION_PROTECTION"
	objectDeletionProtectionDefault    = true
	subobjectDeletionProtectionDefault = false
	subobjectDeletionProtectionMessage = "Note: sub-object deletion protection is IGNORED because it does not work deterministically."
	independentSyncPeriod              = 15 // time in minutes
	minimumIndependentSyncPeriod       = 5  // time in minutes
)

func Run(ctx context.Context, fs *flag.FlagSet, args []string) error {
	akoScheme := apiruntime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(akoScheme))
	utilruntime.Must(akov2.AddToScheme(akoScheme))

	config, err := parseConfiguration(fs, args)
	if err != nil {
		return fmt.Errorf("error parsing configuration: %w", err)
	}

	logger, err := initCustomZapLogger(config.LogLevel, config.LogEncoder)
	if err != nil {
		return fmt.Errorf("error instantiating logger: %w", err)
	}

	logrLogger := zapr.NewLogger(logger)
	ctrl.SetLogger(logrLogger.WithName("ctrl"))
	klog.SetLogger(logrLogger.WithName("klog"))
	setupLog := logger.Named("setup").Sugar()
	if version.IsExperimental() {
		setupLog.Warn("Experimental features enabled!")
		utilruntime.Must(akov2next.AddToScheme(akoScheme))
	}
	setupLog.Info("starting with configuration", zap.Any("config", config), zap.Any("version", version.Version))

	runnable, err := operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), akoScheme, time.Duration(minimumIndependentSyncPeriod)*time.Minute).
		WithConfig(ctrl.GetConfigOrDie()).
		WithNamespaces(collection.Keys(config.WatchedNamespaces)...).
		WithLogger(logger).
		WithMetricAddress(config.MetricsAddr).
		WithProbeAddress(config.ProbeAddr).
		WithLeaderElection(config.EnableLeaderElection).
		WithAtlasDomain(config.AtlasDomain).
		WithAPISecret(config.GlobalAPISecret).
		WithDeletionProtection(config.ObjectDeletionProtection).
		WithIndependentSyncPeriod(time.Duration(config.IndependentSyncPeriod) * time.Minute).
		WithDryRun(config.DryRun).
		WithMaxConcurrentReconciles(config.MaxConcurrentReconciles).
		Build(ctx)
	if err != nil {
		setupLog.Error(err, "unable to start operator")
		return fmt.Errorf("unable to start operator: %w", err)
	}

	setupLog.Info(subobjectDeletionProtectionMessage)
	setupLog.Info("starting manager")
	if err = runnable.Start(ctx); err != nil {
		setupLog.Errorf("error running manager: %v", err)
		return fmt.Errorf("error running manager: %w", err)
	}
	return nil
}

type Config struct {
	AtlasDomain                 string
	EnableLeaderElection        bool
	MetricsAddr                 string
	WatchedNamespaces           map[string]bool
	ProbeAddr                   string
	GlobalAPISecret             client.ObjectKey
	LogLevel                    string
	LogEncoder                  string
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	IndependentSyncPeriod       int
	FeatureFlags                *featureflags.FeatureFlags
	DryRun                      bool
	MaxConcurrentReconciles     int
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func parseConfiguration(fs *flag.FlagSet, args []string) (Config, error) {
	var globalAPISecretName string
	config := Config{}
	fs.StringVar(&config.AtlasDomain, "atlas-domain", operator.DefaultAtlasDomain, "the Atlas URL domain name (with slash in the end).")
	fs.StringVar(&config.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&config.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.StringVar(&globalAPISecretName, "global-api-secret-name", "", "The name of the Secret that contains Atlas API keys. "+
		"It is used by the Operator if AtlasProject configuration doesn't contain API key reference. Defaults to <deployment_name>-api-key.")
	fs.BoolVar(&config.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	fs.StringVar(&config.LogLevel, "log-level", "info", "Log level. Available values: debug | info | warn | error | dpanic | panic | fatal or a numeric value from -9 to 5, where -9 is the most verbose and 5 is the least verbose.")
	fs.StringVar(&config.LogEncoder, "log-encoder", "json", "Log encoder. Available values: json | console")
	fs.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "Defines if the operator deletes Atlas resource "+
		"when a Custom Resource is deleted")
	fs.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "Defines if the operator overwrites "+
		"(and consequently delete) subresources that were not previously created by the operator. "+subobjectDeletionProtectionMessage)
	fs.IntVar(
		&config.IndependentSyncPeriod,
		"independent-sync-period",
		independentSyncPeriod,
		fmt.Sprintf("The default time, in minutes,  between reconciliations for independent custom resources. (default %d, minimum %d)", independentSyncPeriod, minimumIndependentSyncPeriod),
	)
	fs.BoolVar(&config.DryRun, "dry-run", false, "If set, the operator will not perform any changes to the Atlas resources, run all reconcilers only Once and emit events for all planned changes")
	config.MaxConcurrentReconciles, _ = env.GetInt("MDB_MAX_CONCURRENT_RECONCILES", 5) // errors yield default value

	appVersion := fs.Bool("v", false, "prints application version")
	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("failed to parse arguments: %w", err)
	}

	if *appVersion {
		runVersion()
		os.Exit(0)
	}

	config.GlobalAPISecret = operatorGlobalKeySecretOrDefault(globalAPISecretName)

	// dev note: we pass the watched namespace as the env variable to use the Kubernetes Downward API. Unfortunately
	// there is no way to use it for container arguments
	watchedNamespace := strings.TrimSpace(os.Getenv("WATCH_NAMESPACE"))
	if watchedNamespace != "" {
		config.WatchedNamespaces = make(map[string]bool)
		for _, namespace := range strings.Split(watchedNamespace, ",") {
			namespace = strings.TrimSpace(namespace)
			config.WatchedNamespaces[namespace] = true
		}
	}

	configureDeletionProtection(fs, &config)

	config.FeatureFlags = featureflags.NewFeatureFlags(os.Environ)
	return config, nil
}

func operatorGlobalKeySecretOrDefault(secretNameOverride string) client.ObjectKey {
	secretName := secretNameOverride
	if secretName == "" {
		operatorPodName := os.Getenv("OPERATOR_POD_NAME")
		if operatorPodName == "" {
			log.Fatal(`"OPERATOR_POD_NAME" environment variable must be set!`)
		}
		deploymentName, err := kube.ParseDeploymentNameFromPodName(operatorPodName)
		if err != nil {
			log.Fatalf(`Failed to get Operator Deployment name from "OPERATOR_POD_NAME" environment variable: %s`, err.Error())
		}
		secretName = deploymentName + "-api-key"
	}
	operatorNamespace := os.Getenv("OPERATOR_NAMESPACE")
	if operatorNamespace == "" {
		log.Fatal(`"OPERATOR_NAMESPACE" environment variable must be set!`)
	}

	return client.ObjectKey{Namespace: operatorNamespace, Name: secretName}
}

func initCustomZapLogger(level, encoding string) (*zap.Logger, error) {
	lv := zap.AtomicLevel{}

	i64, err := strconv.ParseInt(level, 10, 8)
	numericLevel := int8(i64)
	if err != nil {
		// not a numeric level, try to unmarshal it as a zapcore.Level ("debug", "info", "warn", "error", "dpanic", "panic", or "fatal")
		err := lv.UnmarshalText([]byte(strings.ToLower(level)))
		if err != nil {
			return nil, err
		}
	} else {
		// numeric level:
		// 1. configure klog if the numeric log level is negative and the absolute value of the negative numeric value represents the klog level.
		// 2. configure the atomic zap level based on the numeric value (5..-9).

		var klogLevel int8 = 0
		if numericLevel < 0 {
			klogLevel = -numericLevel
		}

		klogFlagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		klog.InitFlags(klogFlagSet)
		if err := klogFlagSet.Set("v", strconv.Itoa(int(klogLevel))); err != nil {
			return nil, err
		}

		lv = zap.NewAtomicLevelAt(zapcore.Level(numericLevel))
	}

	enc := strings.ToLower(encoding)
	if enc != "json" && enc != "console" {
		return nil, errors.New("'encoding' parameter can only by either 'json' or 'console'")
	}

	cfg := zap.Config{
		Level:             lv,
		OutputPaths:       []string{"stdout"},
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

func configureDeletionProtection(fs *flag.FlagSet, config *Config) {
	if config == nil {
		return
	}

	objectDeletionSet := false
	subObjectDeletionSet := false

	fs.Visit(func(f *flag.Flag) {
		if f.Name == objectDeletionProtectionFlag {
			objectDeletionSet = true
		}

		if f.Name == subobjectDeletionProtectionFlag {
			subObjectDeletionSet = true
		}
	})

	if !objectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(objectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			config.ObjectDeletionProtection = true
		case "false":
			config.ObjectDeletionProtection = false
		default:
			config.ObjectDeletionProtection = objectDeletionProtectionDefault
		}
	}

	if !subObjectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(subobjectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			config.SubObjectDeletionProtection = true
		case "false":
			config.SubObjectDeletionProtection = false
		default:
			config.SubObjectDeletionProtection = subobjectDeletionProtectionDefault
		}
	}
}

func runVersion() {
	v := map[string]string{
		"Version":      version.Version,
		"GitCommit":    version.GitCommit,
		"GoVersion":    runtime.Version(),
		"Platform":     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"BuildTime":    version.BuildTime,
		"Experimental": version.Experimental,
	}
	orderedKeys := []string{
		"Version",
		"GitCommit",
		"GoVersion",
		"Platform",
		"BuildTime",
		"Experimental",
	}

	maxKeyLength := 0
	for k := range v {
		if len(k) > maxKeyLength {
			maxKeyLength = len(k)
		}
	}

	// Create the format string dynamically
	// e.g., "%-12s: %s\n"
	format := fmt.Sprintf("%%-%ds: %%s\n", maxKeyLength)

	for _, key := range orderedKeys {
		value := v[key]
		fmt.Printf(format, key, value)
	}
}
