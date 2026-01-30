// Copyright 2020 MongoDB Inc
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
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/operator"
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
	control.SkipTestUnless(t, "AKO_INT_TEST")

	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	utilruntime.Must(akov2.AddToScheme(scheme.Scheme))

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
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
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

		// shallow copy global config
		ginkgoCfg := *cfg
		ginkgoCfg.UserAgent = "ginkgo"

		// It's recommended to construct the client directly for tests
		// see https://github.com/kubernetes-sigs/controller-runtime/issues/343#issuecomment-469435686
		k8sClient, err = client.New(&ginkgoCfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).ToNot(HaveOccurred())
		Expect(k8sClient).ToNot(BeNil())

		atlasClient, err = atlas.NewClient(atlasDomain, publicKey, privateKey)
		Expect(err).ToNot(HaveOccurred())

		atlasClient, err = admin.NewClient(
			admin.UseBaseURL(atlasDomain),
			admin.UseDigestAuth(publicKey, privateKey),
		)
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

var _ = ReportAfterSuite("Ensure test suite was not empty", func(r Report) {
	Expect(r.PreRunStats.SpecsThatWillRun > 0).To(BeTrue(), "Suite must run at least 1 test")
})

func defaultTimeouts() {
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultEventuallyPollingInterval(PollingInterval)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
}

// prepareControllers is a common function used by all the tests that creates the namespace and registers all the
// reconcilers there. Each of them listens only this specific namespace only, otherwise it's not possible to run in parallel
func prepareControllers(deletionProtection bool) (*corev1.Namespace, context.CancelFunc) {
	var ctx context.Context
	ctx, managerCancelFunc = context.WithCancel(context.Background())
	namespace = corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "test",
			GenerateName: "test",
		},
	}

	By("Creating the namespace " + namespace.GenerateName + "...")
	Expect(k8sClient.Create(ctx, &namespace)).ToNot(HaveOccurred())
	Expect(namespace.Name).ToNot(BeEmpty())
	GinkgoWriter.Printf("Generated namespace %q\n", namespace.Name)

	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))
	ctrl.SetLogger(zapr.NewLogger(logger))

	// shallow copy global config
	managerCfg := *cfg
	managerCfg.UserAgent = "AKO"
	mgr, err := operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), scheme.Scheme, 5*time.Minute).
		WithConfig(&managerCfg).
		WithNamespaces(namespace.Name).
		WithLogger(logger).
		WithAtlasDomain(atlasDomain).
		WithSyncPeriod(30 * time.Minute).
		WithAPISecret(client.ObjectKey{Name: "atlas-operator-api-key", Namespace: namespace.Name}).
		WithDeletionProtection(deletionProtection).
		WithSkipNameValidation(true). // this is needed as this starts multiple controllers concurrently
		Build(ctx)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = mgr.Start(ctx)
		Expect(err).ToNot(HaveOccurred())
	}()

	return &namespace, managerCancelFunc
}

func prepareControllersWithSyncPeriod(deletionProtection bool, syncPeriod time.Duration) (*corev1.Namespace, context.CancelFunc) {
	var ctx context.Context
	ctx, managerCancelFunc = context.WithCancel(context.Background())
	namespace = corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "test",
			GenerateName: "test",
		},
	}

	By("Creating the namespace " + namespace.GenerateName + "...")
	Expect(k8sClient.Create(ctx, &namespace)).ToNot(HaveOccurred())
	Expect(namespace.Name).ToNot(BeEmpty())
	GinkgoWriter.Printf("Generated namespace %q\n", namespace.Name)

	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))
	ctrl.SetLogger(zapr.NewLogger(logger))

	// shallow copy global config
	managerCfg := *cfg
	managerCfg.UserAgent = "AKO"
	mgr, err := operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), scheme.Scheme, 5*time.Minute).
		WithConfig(&managerCfg).
		WithNamespaces(namespace.Name).
		WithLogger(logger).
		WithAtlasDomain(atlasDomain).
		WithSyncPeriod(syncPeriod).
		WithAPISecret(client.ObjectKey{Name: "atlas-operator-api-key", Namespace: namespace.Name}).
		WithDeletionProtection(deletionProtection).
		Build(ctx)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = mgr.Start(ctx)
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
