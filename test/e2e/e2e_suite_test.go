package e2e_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/utils"
)

const (
	EventuallyTimeout   = 100 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

var (
	// default
	Platform    = "kind"
	K8sVersion  = "v1.17.17"
	atlasClient *atlas.Atlas
)

func TestE2e(t *testing.T) {
	if !control.Enabled("AKO_E2E_TEST") {
		t.Skip("Skipping e2e tests, AKO_E2E_TEST is not set")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
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

// checkUpEnvironment initial check setup
func checkUpEnvironment() {
	Platform = os.Getenv("K8S_PLATFORM")
	K8sVersion = os.Getenv("K8S_VERSION")

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
	Expect(utils.SaveToFile(config.FileNameSAGCP, []byte(os.Getenv("GCP_SA_CRED")))).ShouldNot(HaveOccurred())
	Expect(os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.FileNameSAGCP)).ShouldNot(HaveOccurred())
	Expect(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")).ShouldNot(BeEmpty(), "Please, setup GOOGLE_APPLICATION_CREDENTIALS environment variable for test with GCP")
}
