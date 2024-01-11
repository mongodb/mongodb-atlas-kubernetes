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
	"go.mongodb.org/atlas-sdk/v20231115003/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient     client.Client
	testEnv       *envtest.Environment
	cancelManager context.CancelFunc
	atlasClient   *admin.APIClient
	namespace     corev1.Namespace
	atlasDomain   string
	orgID         string
	publicKey     string
	privateKey    string
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
		useExistingCluster := os.Getenv("USE_EXISTING_CLUSTER") != ""
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:  []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
			UseExistingCluster: &useExistingCluster,
		}

		_, err := testEnv.Start()
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setup test dependencies", func() {
		err := mdbv1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		// It's recommended to construct the client directly for tests
		// see https://github.com/kubernetes-sigs/controller-runtime/issues/343#issuecomment-469435686
		k8sClient, err = client.New(testEnv.Config, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).ToNot(BeNil())

		atlasClient, err = atlas.NewClient(atlasDomain, publicKey, privateKey)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Start the operator", func() {
		logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))
		ctrl.SetLogger(zapr.NewLogger(logger))
		syncPeriod := time.Minute * 30
		// The manager watches ALL namespaces
		k8sManager, err := ctrl.NewManager(testEnv.Config, ctrl.Options{
			Scheme:     scheme.Scheme,
			SyncPeriod: &syncPeriod,
		})
		Expect(err).ToNot(HaveOccurred())

		// globalPredicates should be used for general controller Predicates
		// that should be applied to all controllers in order to limit the
		// resources they receive events for.
		globalPredicates := []predicate.Predicate{
			watch.CommonPredicates(), // ignore spurious changes. status changes etc.
			watch.SelectNamespacesPredicate(map[string]bool{ // select only desired namespaces
				namespace.Name: true,
			}),
		}

		atlasProvider := atlas.NewProductionProvider(atlasDomain, kube.ObjectKey(namespace.Name, "atlas-operator-api-key"), k8sManager.GetClient())

		err = (&atlasproject.AtlasProjectReconciler{
			Client:           k8sManager.GetClient(),
			Log:              logger.Named("controllers").Named("AtlasProject").Sugar(),
			ResourceWatcher:  watch.NewResourceWatcher(),
			EventRecorder:    k8sManager.GetEventRecorderFor("AtlasProject"),
			AtlasProvider:    atlasProvider,
			GlobalPredicates: globalPredicates,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = (&atlasdeployment.AtlasDeploymentReconciler{
			Client:           k8sManager.GetClient(),
			Log:              logger.Named("controllers").Named("AtlasDeployment").Sugar(),
			ResourceWatcher:  watch.NewResourceWatcher(),
			EventRecorder:    k8sManager.GetEventRecorderFor("AtlasDeployment"),
			AtlasProvider:    atlasProvider,
			GlobalPredicates: globalPredicates,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
			Client:           k8sManager.GetClient(),
			Log:              logger.Named("controllers").Named("AtlasDeployment").Sugar(),
			ResourceWatcher:  watch.NewResourceWatcher(),
			EventRecorder:    k8sManager.GetEventRecorderFor("AtlasDeployment"),
			AtlasProvider:    atlasProvider,
			GlobalPredicates: globalPredicates,
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())

		var ctx context.Context
		ctx, cancelManager = context.WithCancel(context.Background())

		go func() {
			err = k8sManager.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		}()
	})
})

var _ = AfterSuite(func() {
	By("Tearing down the test environment", func() {
		if cancelManager != nil {
			cancelManager()
		}
		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	})
})
