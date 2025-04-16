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

package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	EventuallyTimeout   = 100 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

var (
	atlasClient *atlas.Atlas
)

func TestE2e(t *testing.T) {
	control.SkipTestUnless(t, "AKO_E2E_TEST")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Atlas Operator E2E Test Suite")
}

var _ = BeforeSuite(func() {
	if !control.Enabled("AKO_E2E_TEST") {
		fmt.Println("Skipping e2e BeforeSuite, AKO_E2E_TEST is not set")

		return
	}

	GinkgoWriter.Write([]byte("==============================Before==============================\n"))
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultEventuallyPollingInterval(PollingInterval)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
	checkUpEnvironment()
	GinkgoWriter.Write([]byte("========================End of Before==============================\n"))
})

var _ = ReportAfterSuite("Ensure test suite was not empty", func(r Report) {
	Expect(r.PreRunStats.SpecsThatWillRun > 0).To(BeTrue(), "Suite must run at least 1 test")
})

// checkUpEnvironment initial check setup
func checkUpEnvironment() {
	Expect(os.Getenv("MCLI_ORG_ID")).ShouldNot(BeEmpty(), "Please, setup MCLI_ORG_ID environment variable")
	Expect(os.Getenv("MCLI_PUBLIC_API_KEY")).ShouldNot(BeEmpty(), "Please, setup MCLI_PUBLIC_API_KEY environment variable")
	Expect(os.Getenv("MCLI_PRIVATE_API_KEY")).ShouldNot(BeEmpty(), "Please, setup MCLI_PRIVATE_API_KEY environment variable")
	Expect(os.Getenv("MCLI_OPS_MANAGER_URL")).ShouldNot(BeEmpty(), "Please, setup MCLI_OPS_MANAGER_URL environment variable")

	atlasClient = atlas.GetClientOrFail()
}

func checkUpAWSEnvironment() {
	Expect(os.Getenv("AWS_ACCESS_KEY_ID")).ShouldNot(BeEmpty(), "Please, setup AWS_ACCESS_KEY_ID environment variable for test with AWS")
	Expect(os.Getenv("AWS_SECRET_ACCESS_KEY")).ShouldNot(BeEmpty(), "Please, setup AWS_SECRET_ACCESS_KEY environment variable for test with AWS")
	Expect(os.Getenv("AWS_ACCOUNT_ARN_LIST")).ShouldNot(BeEmpty(), "Please, setup AWS_ACCOUNT_ARN_LIST environment variable for test with AWS")
}

func checkUpAzureEnvironment() {
	Expect(os.Getenv("AZURE_CLIENT_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_CLIENT_ID environment variable for test with Azure")
	Expect(os.Getenv("AZURE_TENANT_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_TENANT_ID environment variable for test with Azure")
	Expect(os.Getenv("AZURE_CLIENT_SECRET")).ShouldNot(BeEmpty(), "Please, setup AZURE_CLIENT_SECRET environment variable for test with Azure")
	Expect(os.Getenv("AZURE_SUBSCRIPTION_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_SUBSCRIPTION_ID environment variable for test with Azure")
}

func checkNSetUpGCPEnvironment() {
	Expect(os.Getenv("GCP_SA_CRED")).ShouldNot(BeEmpty(), "Please, setup GCP_SA_CRED environment variable for test with GCP (req. Service Account)")
	Expect(os.MkdirAll("../../output", 0750)).ShouldNot(HaveOccurred())
	Expect(utils.SaveToFile(config.FileNameSAGCP, []byte(os.Getenv("GCP_SA_CRED")))).ShouldNot(HaveOccurred())
	Expect(os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.FileNameSAGCP)).ShouldNot(HaveOccurred())
	Expect(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")).ShouldNot(BeEmpty(), "Please, setup GOOGLE_APPLICATION_CREDENTIALS environment variable for test with GCP")
}
