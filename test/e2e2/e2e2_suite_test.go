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

package e2e2_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.uber.org/zap/zaptest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	atlasctrl "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
)

const (
	EventuallyTimeout   = 100 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

const (
	// nolint:gosec
	DefaultGlobalCredentials = "mongodb-atlas-operator-api-key"
)

var GinkGoFieldOwner = client.FieldOwner("ginkgo")

var suiteCtx context.Context
var suiteCancel context.CancelFunc

func TestE2e(t *testing.T) {
	control.SkipTestUnless(t, "AKO_E2E2_TEST")

	initTestLogging(t)
	suiteCtx, suiteCancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Atlas Operator E2E2 Test Suite tests the operator binary")
}

var _ = BeforeSuite(func() {
	if !control.Enabled("AKO_E2E2_TEST") {
		fmt.Println("Skipping e2e2 BeforeSuite, AKO_E2E2_TEST is not set")
		return
	}
	GinkgoWriter.Write([]byte("==============================Before==============================\n"))
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultEventuallyPollingInterval(PollingInterval)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
	GinkgoWriter.Write([]byte("========================End of Before==============================\n"))
})

var _ = ReportAfterSuite("Ensure test suite was not empty", func(r Report) {
	Expect(r.PreRunStats.SpecsThatWillRun > 0).To(BeTrue(), "Suite must run at least 1 test")
})

func initTestLogging(t *testing.T) {
	logger := zaptest.NewLogger(t)
	logrLogger := zapr.NewLogger(logger)
	ctrllog.SetLogger(logrLogger.WithName("test"))
}

// nolint:unparam
func runTestAKO(globalCreds, ns string, deletionprotection bool) operator.Operator {
	args := []string{
		"--log-level=-9",
		fmt.Sprintf("--global-api-secret-name=%s", globalCreds),
		"--log-encoder=json",
		`--atlas-domain=https://cloud-qa.mongodb.com`,
	}
	args = append(args, fmt.Sprintf("--object-deletion-protection=%v", deletionprotection))
	return operator.NewOperator(operator.AllNamespacesOperatorEnv(ns), os.Stdout, os.Stderr, args...)
}

// newTestAtlasClient creates an Atlas API client for e2e2 tests.
// It reads credentials from environment variables and returns a configured client.
// Required environment variables:
//   - MCLI_ORG_ID: Atlas organization ID
//   - MCLI_PUBLIC_API_KEY: Atlas public API key
//   - MCLI_PRIVATE_API_KEY: Atlas private API key
//   - MCLI_OPS_MANAGER_URL: (optional) Atlas API base URL, defaults to "https://cloud.mongodb.com/"
func newTestAtlasClient() (*admin.APIClient, string) {
	orgID := os.Getenv("MCLI_ORG_ID")
	Expect(orgID).NotTo(BeEmpty(), "MCLI_ORG_ID environment variable must be set")

	publicKey := os.Getenv("MCLI_PUBLIC_API_KEY")
	Expect(publicKey).NotTo(BeEmpty(), "MCLI_PUBLIC_API_KEY environment variable must be set")

	privateKey := os.Getenv("MCLI_PRIVATE_API_KEY")
	Expect(privateKey).NotTo(BeEmpty(), "MCLI_PRIVATE_API_KEY environment variable must be set")

	atlasDomain := os.Getenv("MCLI_OPS_MANAGER_URL")
	if atlasDomain == "" {
		atlasDomain = "https://cloud-qa.mongodb.com/"
	}

	client, err := atlasctrl.NewClient(atlasDomain, publicKey, privateKey)
	Expect(err).ToNot(HaveOccurred())

	return client, orgID
}
