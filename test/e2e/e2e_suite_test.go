package e2e_test

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	EventuallyTimeout   = 60 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	// TODO data provider?
	ConfigAll     = "../../deploy/" // Released generated files
	ProjectSample = "data/atlasproject.yaml"
	ClusterSample = "data/atlascluster_basic.yaml"
)

var (
	// default
	Platform   = "kind"
	K8sVersion = "v1.17.17"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	GinkgoWriter.Write([]byte("==============================Before==============================\n"))
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
	checkUpMongoCLI()
	GinkgoWriter.Write([]byte("========================End of Before==============================\n"))
})

// setUpMongoCLI initial setup
func checkUpMongoCLI() {
	Platform = os.Getenv("K8s_PLATFORM")
	K8sVersion = os.Getenv("K8s_VERSION")
	// additional checks
	Expect(os.Getenv("MCLI_ORG_ID")).ShouldNot(BeEmpty())
	Expect(os.Getenv("MCLI_PUBLIC_API_KEY")).ShouldNot(BeEmpty())
	Expect(os.Getenv("MCLI_PRIVATE_API_KEY")).ShouldNot(BeEmpty())
	Expect(os.Getenv("MCLI_OPS_MANAGER_URL")).ShouldNot(BeEmpty())
}
