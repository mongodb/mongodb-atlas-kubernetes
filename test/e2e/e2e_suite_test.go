package e2e_test

import (
	"fmt"
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
	ConfigAll     = "../../deploy/all-in-one.yaml" // basic configuration (release)
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

var _ = SynchronizedBeforeSuite(func() []byte {
	GinkgoWriter.Write([]byte("==============================Global FIRST Node Synchronized Before Each==============================\n"))
	GinkgoWriter.Write([]byte("SetUp Global Timeout\n"))
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
	checkUpMongoCLI()
	GinkgoWriter.Write([]byte("==============================End of Global FIRST Node Synchronized Before Each=======================\n"))
	return nil
}, func(_ []byte) {
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================Global Node %d Synchronized Before Each==============================\n", GinkgoParallelNode())))
	if GinkgoParallelNode() != 1 {
		Fail("Please Test suite cannot run in parallel") // TODO prepare configurations for parallel
	}
	GinkgoWriter.Write([]byte(fmt.Sprintf("==============================End of Global Node %d Synchronized Before Each========================\n", GinkgoParallelNode())))
})

var _ = BeforeEach(func() {
	GinkgoWriter.Write([]byte("==============================Global Before Each==============================\n"))
	GinkgoWriter.Write([]byte("========================End of Global Before Each==============================\n"))
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
