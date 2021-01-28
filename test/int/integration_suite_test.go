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
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlascluster"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasproject"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/httputil"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	ctrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg         *rest.Config
	k8sClient   client.Client
	testEnv     *envtest.Environment
	k8sManager  ctrl.Manager
	atlasClient *mongodbatlas.Client
	connection  atlas.Connection
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

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	// We cannot serialize the 'cfg' as it contains functions. Serializing only the host url is enough.
	return []byte(cfg.Host)
}, func(data []byte) {
	By("Preparing controllers and atlas client")

	// This is the host that was serialized on the 1st node by the function above.
	host := string(data)
	// copied from Environment.Start()
	cfg = &rest.Config{
		Host: host,
		// gotta go fast during tests -- we don't really care about overwhelming our test API server
		QPS:   1000.0,
		Burst: 2000.0,
	}
	err := mdbv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme
	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))

	ctrl.SetLogger(zapr.NewLogger(logger))

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme:             scheme.Scheme,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasproject.AtlasProjectReconciler{
		Client:      k8sManager.GetClient(),
		Log:         logger.Named("controllers").Named("AtlasProject").Sugar(),
		AtlasDomain: "https://cloud-qa.mongodb.com",
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&atlascluster.AtlasClusterReconciler{
		Client:      k8sManager.GetClient(),
		Log:         logger.Named("controllers").Named("AtlasCluster").Sugar(),
		AtlasDomain: "https://cloud-qa.mongodb.com",
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	atlasClient, connection = prepareAtlasClient()
})

var _ = SynchronizedAfterSuite(func() {

}, func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

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
