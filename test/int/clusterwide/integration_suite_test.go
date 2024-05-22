/*
Copyright 2020 MongoDB.

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

package int

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasstream"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	ctx     context.Context
	testEnv *envtest.Environment

	k8sClient   client.Client
	atlasClient *admin.APIClient

	managerCancelFunc context.CancelFunc
	atlasDomain       string
	orgID             string
	publicKey         string
	privateKey        string
)

const (
	atlasDomainDefault = "https://cloud-qa.mongodb.com/"
)

func TestAPIs(t *testing.T) {
	if !control.Enabled("AKO_INT_TEST") {
		t.Skip("Skipping int tests, AKO_INT_TEST is not set")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Atlas Operator Cluster-Wide Integration Test Suite")
}

var _ = BeforeSuite(func() {
	if !control.Enabled("AKO_INT_TEST") {
		fmt.Println("Skipping int BeforeSuite, AKO_INT_TEST is not set")

		return
	}

	ctx = ctrl.SetupSignalHandler()

	By("Validating configuration data is available", func() {
		orgID, publicKey, privateKey = os.Getenv("ATLAS_ORG_ID"), os.Getenv("ATLAS_PUBLIC_KEY"), os.Getenv("ATLAS_PRIVATE_KEY")
		Expect(orgID).ToNot(BeEmpty())
		Expect(publicKey).ToNot(BeEmpty())
		Expect(privateKey).ToNot(BeEmpty())

		if atlasDomain = os.Getenv("ATLAS_DOMAIN"); atlasDomain == "" {
			atlasDomain = atlasDomainDefault
		}
	})

	By("Bootstrapping test environment", func() {
		_, useExistingCluster := os.LookupEnv("USE_EXISTING_CLUSTER")
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:  []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
			UseExistingCluster: &useExistingCluster,
		}

		_, err := testEnv.Start()
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setup test dependencies", func() {
		err := akov2.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		// shallow copy global config
		ginkgoCfg := *testEnv.Config
		ginkgoCfg.UserAgent = "ginkgo"

		// It's recommended to construct the client directly for tests
		// see https://github.com/kubernetes-sigs/controller-runtime/issues/343#issuecomment-469435686
		k8sClient, err = client.New(&ginkgoCfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).ToNot(BeNil())

		atlasClient, err = atlas.NewClient(atlasDomain, publicKey, privateKey)
		Expect(err).ToNot(HaveOccurred())
	})

	var k8sManager manager.Manager
	var atlasProvider atlas.Provider
	var globalPredicates []predicate.Predicate
	var logger *zap.Logger
	var err error
	ctx, managerCancelFunc = context.WithCancel(ctx)

	By("Setting up operator manager", func() {
		logger = ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))
		ctrl.SetLogger(zapr.NewLogger(logger))

		syncPeriod := time.Minute * 30

		// shallow copy global config
		managerCfg := *testEnv.Config
		managerCfg.UserAgent = "AKO"
		k8sManager, err = ctrl.NewManager(&managerCfg, ctrl.Options{
			Scheme:  scheme.Scheme,
			Metrics: metricsserver.Options{BindAddress: "0"},
			Cache: cache.Options{
				SyncPeriod: &syncPeriod,
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up controllers dependencies", func() {
		globalPredicates = []predicate.Predicate{
			watch.CommonPredicates(), // ignore spurious changes. status changes etc.
			watch.SelectNamespacesPredicate(map[string]bool{ // select only desired namespaces
				"": true,
			}),
		}
		atlasProvider = atlas.NewProductionProvider(atlasDomain, kube.ObjectKey("default", "atlas-operator-api-key"), k8sManager.GetClient())

		err = indexer.RegisterAll(ctx, k8sManager, logger)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up project controller", func() {
		err = (&atlasproject.AtlasProjectReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
			DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
			GlobalPredicates:            globalPredicates,
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasProject"),
			AtlasProvider:               atlasProvider,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up deployment controller", func() {
		err = (&atlasdeployment.AtlasDeploymentReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
			GlobalPredicates:            globalPredicates,
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasDeployment"),
			AtlasProvider:               atlasProvider,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up database user controller", func() {
		featureFlags := featureflags.NewFeatureFlags(os.Environ)

		err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
			Client:                        k8sManager.GetClient(),
			Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
			EventRecorder:                 k8sManager.GetEventRecorderFor("AtlasDatabaseUser"),
			DeprecatedResourceWatcher:     watch.NewDeprecatedResourceWatcher(),
			AtlasProvider:                 atlasProvider,
			GlobalPredicates:              globalPredicates,
			SubObjectDeletionProtection:   false,
			FeaturePreviewOIDCAuthEnabled: featureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up data federation controller", func() {
		err = (&atlasdatafederation.AtlasDataFederationReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasDatabaseUser"),
			DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
			AtlasProvider:               atlasProvider,
			GlobalPredicates:            globalPredicates,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up federated authentication controller", func() {
		err = (&atlasfederatedauth.AtlasFederatedAuthReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
			DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
			GlobalPredicates:            globalPredicates,
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasFederatedAuth"),
			AtlasProvider:               atlasProvider,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up streams instance controller", func() {
		err = (&atlasstream.AtlasStreamsInstanceReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasStreamInstance").Sugar(),
			GlobalPredicates:            globalPredicates,
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasStreamInstance"),
			AtlasProvider:               atlasProvider,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(ctx, k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setting up streams connection controller", func() {
		err = (&atlasstream.AtlasStreamsConnectionReconciler{
			Client:                      k8sManager.GetClient(),
			Log:                         logger.Named("controllers").Named("AtlasStreamConnection").Sugar(),
			GlobalPredicates:            globalPredicates,
			EventRecorder:               k8sManager.GetEventRecorderFor("AtlasStreamConnection"),
			AtlasProvider:               atlasProvider,
			SubObjectDeletionProtection: false,
		}).SetupWithManager(ctx, k8sManager)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Starting controllers", func() {
		go func() {
			err = k8sManager.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		}()
	})
})

var _ = AfterSuite(func() {
	By("Tearing down the test environment", func() {
		if managerCancelFunc != nil {
			managerCancelFunc()
		}

		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	})
})
