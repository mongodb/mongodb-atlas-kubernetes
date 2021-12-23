package e2e_test

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
)

const (
	EventuallyTimeout   = 100 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	PollingInterval     = 10 * time.Second
)

var (
	// default
	Platform   = "kind"
	K8sVersion = "v1.17.17"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
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
	mongocli.GetVersionOutput()
}

func checkUpAWSEnviroment() {
	// TODO
}

func checkUpAzureEnviroment() {
	Expect(os.Getenv("AZURE_CLIENT_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_CLIENT_ID environment variable for test with Azure")
	Expect(os.Getenv("AZURE_TENANT_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_TENANT_ID environment variable for test with Azure")
	Expect(os.Getenv("AZURE_CLIENT_SECRET")).ShouldNot(BeEmpty(), "Please, setup AZURE_CLIENT_SECRET environment variable for test with Azure")
	Expect(os.Getenv("AZURE_SUBSCRIPTION_ID")).ShouldNot(BeEmpty(), "Please, setup AZURE_SUBSCRIPTION_ID environment variable for test with Azure")
}
