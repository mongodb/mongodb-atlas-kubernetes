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
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlascluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdatabaseuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	EventuallyTimeout   = 60 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

var (
	// This variable is "global" - as is visible only on the first ginkgo node
	testEnv *envtest.Environment

	// These variables are initialized once per each node
	k8sClient   client.Client
	atlasClient *mongodbatlas.Client
	connection  atlas.Connection

	// These variables are per each test and are changed by each BeforeRun
	namespace         corev1.Namespace
	cfg               *rest.Config
	managerCancelFunc context.CancelFunc
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Project Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

// SynchronizedBeforeSuite uses the parallel "with singleton" pattern described by ginkgo
// http://onsi.github.io/ginkgo/#parallel-specs
// The first function starts the envtest (done only once by the 1st node). The second function is called on each of
// the ginkgo nodes and initializes all reconcilers and clients that will be used by the test.
var _ = SynchronizedBeforeSuite(func() []byte {
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
	}

	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	fmt.Printf("Api Server is listening on %s\n", cfg.Host)
	return []byte(cfg.Host)
}, func(data []byte) {
	if os.Getenv("USE_EXISTING_CLUSTER") != "" {
		var err error
		// For the existing cluster we read the kubeconfig
		cfg, err = config.GetConfig()
		if err != nil {
			panic("Failed to read the config for existing cluster")
		}
	} else {
		// This is the host that was serialized on the 1st node by the function above.
		host := string(data)
		// copied from Environment.Start()
		cfg = &rest.Config{
			Host: host,
			// gotta go fast during tests -- we don't really care about overwhelming our test API server
			QPS:   1000.0,
			Burst: 2000.0,
		}
	}

	err := mdbv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// It's recommended to construct the client directly for tests
	// see https://github.com/kubernetes-sigs/controller-runtime/issues/343#issuecomment-469435686
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	atlasClient, connection = prepareAtlasClient()
	defaultTimeouts()
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func defaultTimeouts() {
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultEventuallyPollingInterval(PollingInterval)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
}

func prepareAtlasClient() (*mongodbatlas.Client, atlas.Connection) {
	orgID, publicKey, privateKey := os.Getenv("ATLAS_ORG_ID"), os.Getenv("ATLAS_PUBLIC_KEY"), os.Getenv("ATLAS_PRIVATE_KEY")
	if orgID == "" || publicKey == "" || privateKey == "" {
		Fail(`All of the "ATLAS_ORG_ID", "ATLAS_PUBLIC_KEY", and "ATLAS_PRIVATE_KEY" environment variables must be set!`)
	}
	withDigest := httputil.Digest(publicKey, privateKey)
	httpClient, err := httputil.DecorateClient(&http.Client{Transport: http.DefaultTransport}, withDigest)
	Expect(err).ToNot(HaveOccurred())
	aClient, err := mongodbatlas.New(httpClient, mongodbatlas.SetBaseURL("https://cloud-qa.mongodb.com/api/atlas/v1.0/"))
	Expect(err).ToNot(HaveOccurred())

	return aClient, atlas.Connection{
		OrgID:      orgID,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

// prepareControllers is a common function used by all the tests that creates the namespace and registers all the
// reconcilers there. Each of them listens only this specific namespace only, otherwise it's not possible to run in parallel
func prepareControllers() {
	err := mdbv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	namespace = corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "test",
			GenerateName: "test",
		},
	}

	By("Creating the namespace " + namespace.Name)
	Expect(k8sClient.Create(context.Background(), &namespace)).ToNot(HaveOccurred())

	// +kubebuilder:scaffold:scheme
	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))

	ctrl.SetLogger(zapr.NewLogger(logger))

	// Note on the syncPeriod - decreasing this to a smaller time allows to test its work for the long-running tests
	// (clusters, database users). The prod value is much higher
	syncPeriod := time.Minute * 30
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		Namespace:          namespace.Name,
		MetricsBindAddress: "0",
		SyncPeriod:         &syncPeriod,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasproject.AtlasProjectReconciler{
		Client:          k8sManager.GetClient(),
		Log:             logger.Named("controllers").Named("AtlasProject").Sugar(),
		AtlasDomain:     "https://cloud-qa.mongodb.com",
		ResourceWatcher: watch.NewResourceWatcher(),
		GlobalAPISecret: kube.ObjectKey(namespace.Name, "atlas-operator-api-key"),
		EventRecorder:   k8sManager.GetEventRecorderFor("AtlasProject"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlascluster.AtlasClusterReconciler{
		Client:          k8sManager.GetClient(),
		Log:             logger.Named("controllers").Named("AtlasCluster").Sugar(),
		AtlasDomain:     "https://cloud-qa.mongodb.com",
		GlobalAPISecret: kube.ObjectKey(namespace.Name, "atlas-operator-api-key"),
		EventRecorder:   k8sManager.GetEventRecorderFor("AtlasCluster"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasdatabaseuser.AtlasDatabaseUserReconciler{
		Client:          k8sManager.GetClient(),
		Log:             logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		AtlasDomain:     "https://cloud-qa.mongodb.com",
		EventRecorder:   k8sManager.GetEventRecorderFor("AtlasDatabaseUser"),
		ResourceWatcher: watch.NewResourceWatcher(),
		GlobalAPISecret: kube.ObjectKey(namespace.Name, "atlas-operator-api-key"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	By("Starting controllers")

	var ctx context.Context
	ctx, managerCancelFunc = context.WithCancel(context.Background())

	go func() {
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()
}

func removeControllersAndNamespace() {
	// end the manager
	managerCancelFunc()

	By("Removing the namespace " + namespace.Name)
	err := k8sClient.Delete(context.Background(), &namespace)
	Expect(err).ToNot(HaveOccurred())
}
