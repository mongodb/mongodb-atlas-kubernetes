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
	"context"
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
	"k8s.io/client-go/kubernetes/scheme"
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
	control.SkipTestUnless(t, "AKO_INT_TEST")

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
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		}

		_, err := testEnv.Start()
		Expect(err).ToNot(HaveOccurred())
	})

	By("Setup test dependencies", func() {
		err := akov2.AddToScheme(scheme.Scheme)
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
		var ctx context.Context
		ctx, cancelManager = context.WithCancel(context.Background())

		logger := ctrzap.NewRaw(ctrzap.UseDevMode(true), ctrzap.WriteTo(GinkgoWriter), ctrzap.StacktraceLevel(zap.ErrorLevel))
		ctrl.SetLogger(zapr.NewLogger(logger))

		mgr, err := operator.NewBuilder(operator.ManagerProviderFunc(ctrl.NewManager), testEnv.Scheme, 5*time.Minute).
			WithConfig(testEnv.Config).
			WithLogger(logger).
			WithAtlasDomain(atlasDomain).
			WithSyncPeriod(30 * time.Minute).
			Build(ctx)
		Expect(err).ToNot(HaveOccurred())

		go func() {
			err = mgr.Start(ctx)
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

var _ = ReportAfterSuite("Ensure test suite was not empty", func(r Report) {
	Expect(r.PreRunStats.SpecsThatWillRun > 0).To(BeTrue(), "Suite must run at least 1 test")
})
