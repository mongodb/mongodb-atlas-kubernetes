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
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdatafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasfederatedauth"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	atlasDomainDefault  = "https://cloud-qa.mongodb.com/"
	EventuallyTimeout   = 60 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

var (
	// This variable is "global" - as is visible only on the first ginkgo node
	testEnv *envtest.Environment

	// These variables are initialized once per each node
	k8sClient   client.Client
	atlasClient *admin.APIClient

	// These variables are per each test and are changed by each BeforeRun
	namespace         corev1.Namespace
	cfg               *rest.Config
	managerCancelFunc context.CancelFunc
	orgID             string
	publicKey         string
	privateKey        string
	atlasDomain       = atlasDomainDefault
)

func TestAPIs(t *testing.T) {
	if !control.Enabled("AKO_INT_TEST") {
		t.Skip("Skipping int tests, AKO_INT_TEST is not set")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Atlas Operator Namespaced Integration Test Suite")
}

// SynchronizedBeforeSuite uses the parallel "with singleton" pattern described by ginkgo
// http://onsi.github.io/ginkgo/#parallel-specs
// The first function starts the envtest (done only once by the 1st node). The second function is called on each of
// the ginkgo nodes and initializes all reconcilers and clients that will be used by the test.
var _ = SynchronizedBeforeSuite(func() []byte {
	if !control.Enabled("AKO_INT_TEST") {
		fmt.Println("Skipping int SynchronizedBeforeSuite, AKO_INT_TEST is not set")
		return nil
	}

	By("Validating configuration data is available", func() {
		Expect(os.Getenv("ATLAS_ORG_ID")).ToNot(BeEmpty())
		Expect(os.Getenv("ATLAS_PUBLIC_KEY")).ToNot(BeEmpty())
		Expect(os.Getenv("ATLAS_PRIVATE_KEY")).ToNot(BeEmpty())
	})

	var b bytes.Buffer

	By("Bootstrapping test environment", func() {
		useExistingCluster := os.Getenv("USE_EXISTING_CLUSTER") != ""
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:  []string{filepath.Join("..", "..", "config", "crd", "bases")},
			UseExistingCluster: &useExistingCluster,
		}

		cfg, err := testEnv.Start()
		Expect(err).ToNot(HaveOccurred())
		Expect(cfg).ToNot(BeNil())

		e := gob.NewEncoder(&b)
		err = e.Encode(*cfg)
		Expect(err).ToNot(HaveOccurred())

		GinkgoWriter.Printf("Api Server is listening on %s\n", cfg.Host)
	})

	return b.Bytes()
}, func(data []byte) {
	By("Setup test dependencies", func() {
		d := gob.NewDecoder(bytes.NewReader(data))
		err := d.Decode(&cfg)
		Expect(err).ToNot(HaveOccurred())

		orgID, publicKey, privateKey = os.Getenv("ATLAS_ORG_ID"), os.Getenv("ATLAS_PUBLIC_KEY"), os.Getenv("ATLAS_PRIVATE_KEY")

		if domain, found := os.LookupEnv("ATLAS_DOMAIN"); found && domain != "" {
			atlasDomain = domain
		}

		err = mdbv1.AddToScheme(scheme.Scheme)
		Expect(err).ToNot(HaveOccurred())

		// It's recommended to construct the client directly for tests
		// see https://github.com/kubernetes-sigs/controller-runtime/issues/343#issuecomment-469435686
		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).ToNot(HaveOccurred())
		Expect(k8sClient).ToNot(BeNil())

		atlasClient, err = atlas.NewClient(atlasDomain, publicKey, privateKey)
		Expect(err).ToNot(HaveOccurred())
		defaultTimeouts()
	})
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	By("Tearing down the test environment", func() {
		err := testEnv.Stop()
		Expect(err).ToNot(HaveOccurred())
	})
})

func defaultTimeouts() {
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultEventuallyPollingInterval(PollingInterval)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
}

// prepareControllers is a common function used by all the tests that creates the namespace and registers all the
// reconcilers there. Each of them listens only this specific namespace only, otherwise it's not possible to run in parallel
func prepareControllers(deletionProtection bool) (*corev1.Namespace, context.CancelFunc) {
	err := mdbv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	namespace = corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "test",
			GenerateName: "test",
		},
	}

	By("Creating the namespace " + namespace.GenerateName + "...")
	Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())
	Expect(namespace.Name).ToNot(BeEmpty())
	GinkgoWriter.Printf("Generated namespace %q\n", namespace.Name)

	// +kubebuilder:scaffold:scheme
	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))

	ctrl.SetLogger(zapr.NewLogger(logger))

	// Note on the syncPeriod - decreasing this to a smaller time allows to test its work for the long-running tests
	// (deployments, database users). The prod value is much higher
	syncPeriod := time.Minute * 30
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
		SyncPeriod:         &syncPeriod,
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

	featureFlags := featureflags.NewFeatureFlags(os.Environ)

	atlasProvider := atlas.NewProductionProvider(atlasDomain, kube.ObjectKey(namespace.Name, "atlas-operator-api-key"), k8sManager.GetClient())

	err = (&atlasproject.AtlasProjectReconciler{
		Client:                      k8sManager.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasProject").Sugar(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               k8sManager.GetEventRecorderFor("AtlasProject"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    deletionProtection,
		SubObjectDeletionProtection: deletionProtection,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasdeployment.AtlasDeploymentReconciler{
		Client:                      k8sManager.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDeployment").Sugar(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               k8sManager.GetEventRecorderFor("AtlasDeployment"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    deletionProtection,
		SubObjectDeletionProtection: deletionProtection,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		Client:                        k8sManager.GetClient(),
		Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		EventRecorder:                 k8sManager.GetEventRecorderFor("AtlasDatabaseUser"),
		ResourceWatcher:               watch.NewResourceWatcher(),
		AtlasProvider:                 atlasProvider,
		GlobalPredicates:              globalPredicates,
		ObjectDeletionProtection:      deletionProtection,
		SubObjectDeletionProtection:   deletionProtection,
		FeaturePreviewOIDCAuthEnabled: featureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasdatafederation.AtlasDataFederationReconciler{
		Client:                      k8sManager.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		EventRecorder:               k8sManager.GetEventRecorderFor("AtlasDatabaseUser"),
		ResourceWatcher:             watch.NewResourceWatcher(),
		AtlasProvider:               atlasProvider,
		GlobalPredicates:            globalPredicates,
		ObjectDeletionProtection:    deletionProtection,
		SubObjectDeletionProtection: deletionProtection,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasfederatedauth.AtlasFederatedAuthReconciler{
		Client:                      k8sManager.GetClient(),
		Log:                         logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
		ResourceWatcher:             watch.NewResourceWatcher(),
		GlobalPredicates:            globalPredicates,
		EventRecorder:               k8sManager.GetEventRecorderFor("AtlasFederatedAuth"),
		AtlasProvider:               atlasProvider,
		ObjectDeletionProtection:    deletionProtection,
		SubObjectDeletionProtection: deletionProtection,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	By("Starting controllers")

	var ctx context.Context
	ctx, managerCancelFunc = context.WithCancel(context.Background())

	go func() {
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	return &namespace, managerCancelFunc
}

func removeControllersAndNamespace() {
	// end the manager
	managerCancelFunc()

	By("Removing the namespace " + namespace.Name)
	err := k8sClient.Delete(context.Background(), &namespace)
	Expect(err).ToNot(HaveOccurred())
}

func secretData() map[string]string {
	return map[string]string{
		OrgID:         orgID,
		PublicAPIKey:  publicKey,
		PrivateAPIKey: privateKey,
	}
}
