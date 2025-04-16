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

package k8s

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/operator"
)

var (
	setupSignalHandlerOnce sync.Once
	signalCancelledCtx     context.Context
)

func BuildCluster(initCfg *Config) (cluster.Cluster, error) {
	akoScheme := runtime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(akoScheme))
	utilruntime.Must(akov2.AddToScheme(akoScheme))

	config := mergeConfiguration(initCfg)
	logger := zaptest.NewLogger(
		GinkgoT(),
		zaptest.WrapOptions(
			zap.ErrorOutput(zapcore.Lock(zapcore.AddSync(GinkgoWriter))),
		),
	)
	ctrl.SetLogger(zapr.NewLogger(logger))
	setupLog := logger.Named("setup").Sugar()
	setupLog.Info("starting with configuration", zap.Any("config", *config))

	// Ensure all concurrent managers configured per test share a single exit signal handler
	setupSignalHandlerOnce.Do(func() {
		signalCancelledCtx = ctrl.SetupSignalHandler()
	})

	return operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), akoScheme, 5*time.Minute).
		WithConfig(ctrl.GetConfigOrDie()).
		WithNamespaces(collection.Keys(config.WatchedNamespaces)...).
		WithLogger(logger).
		WithMetricAddress(config.MetricsAddr).
		WithProbeAddress(config.ProbeAddr).
		WithLeaderElection(config.EnableLeaderElection).
		WithAtlasDomain(config.AtlasDomain).
		WithAPISecret(config.GlobalAPISecret).
		WithDeletionProtection(config.ObjectDeletionProtection).
		WithSkipNameValidation(true). // this is needed as this starts multiple controllers concurrently
		Build(signalCancelledCtx)
}

type Config struct {
	AtlasDomain                 string
	EnableLeaderElection        bool
	MetricsAddr                 string
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

	return config
}

type ManagerStart func(ctx context.Context) error
type ManagerConfig func(config *Config)

func managerDefaults() *Config {
	return &Config{
		AtlasDomain:                 "https://cloud-qa.mongodb.com/",
		EnableLeaderElection:        false,
		MetricsAddr:                 "0",
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

	c, err := BuildCluster(managerConfig)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		err = c.Start(ctx)
		if err != nil {
			return err
		}

		return nil
	}, nil
}
