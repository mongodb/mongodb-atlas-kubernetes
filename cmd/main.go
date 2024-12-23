/*
Copyright 2020 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
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

func main() {
	akoScheme := runtime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(akoScheme))
	utilruntime.Must(akov2.AddToScheme(akoScheme))

	ctx := ctrl.SetupSignalHandler()
	config := parseConfiguration()

	logger, err := initCustomZapLogger(config.LogLevel, config.LogEncoder)
	if err != nil {
		fmt.Printf("error instantiating logger: %v\r\n", err)
		os.Exit(1)
	}

	ctrl.SetLogger(zapr.NewLogger(logger))
	klog.SetLogger(zapr.NewLogger(logger))
	setupLog := logger.Named("setup").Sugar()
	setupLog.Info("starting with configuration", zap.Any("config", config), zap.Any("version", version.Version))

	mgr, err := operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), akoScheme, time.Duration(minimumIndependentSyncPeriod)*time.Minute).
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
		Build(ctx)
	if err != nil {
		setupLog.Error(err, "unable to start operator")
		os.Exit(1)
	}

	setupLog.Info(subobjectDeletionProtectionMessage)
	setupLog.Info("starting manager")
	if err = mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
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
}

// ParseConfiguration fills the 'OperatorConfig' from the flags passed to the program
func parseConfiguration() Config {
	var globalAPISecretName string
	config := Config{}
	flag.StringVar(&config.AtlasDomain, "atlas-domain", operator.DefaultAtlasDomain, "the Atlas URL domain name (with slash in the end).")
	flag.StringVar(&config.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&config.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&globalAPISecretName, "global-api-secret-name", "", "The name of the Secret that contains Atlas API keys. "+
		"It is used by the Operator if AtlasProject configuration doesn't contain API key reference. Defaults to <deployment_name>-api-key.")
	flag.BoolVar(&config.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&config.LogLevel, "log-level", "info", "Log level. Available values: debug | info | warn | error | dpanic | panic | fatal")
	flag.StringVar(&config.LogEncoder, "log-encoder", "json", "Log encoder. Available values: json | console")
	flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "Defines if the operator deletes Atlas resource "+
		"when a Custom Resource is deleted")
	flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "Defines if the operator overwrites "+
		"(and consequently delete) subresources that were not previously created by the operator. "+subobjectDeletionProtectionMessage)
	flag.IntVar(
		&config.IndependentSyncPeriod,
		"independent-sync-period",
		independentSyncPeriod,
		fmt.Sprintf("The default time, in minutes,  between reconciliations for independent custom resources. (default %d, minimum %d)", independentSyncPeriod, minimumIndependentSyncPeriod),
	)
	appVersion := flag.Bool("v", false, "prints application version")
	flag.Parse()

	if *appVersion {
		fmt.Println(version.Version)
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

	configureDeletionProtection(&config)

	config.FeatureFlags = featureflags.NewFeatureFlags(os.Environ)
	return config
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

func configureDeletionProtection(config *Config) {
	if config == nil {
		return
	}

	objectDeletionSet := false
	subObjectDeletionSet := false

	flag.Visit(func(f *flag.Flag) {
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
