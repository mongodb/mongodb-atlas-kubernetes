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

var _ = BeforeSuite(func(done Done) {
	logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))

	ctrl.SetLogger(zapr.NewLogger(logger))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "config", "crd", "bases")},
	}

	atlasClient, connection = prepareAtlasClient()

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = mdbv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&atlasproject.AtlasProjectReconciler{
		Client: k8sManager.GetClient(),
		Log:    *logger.Named("controllers").Named("AtlasProject").Sugar(),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
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
