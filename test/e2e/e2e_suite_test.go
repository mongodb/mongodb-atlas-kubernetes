package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"time"
)

const (
	EventuallyTimeout   = 60 * time.Second
	ConsistentlyTimeout = 1 * time.Second
	//TODO data provider?
	ConfigAll     = "../../deploy/all-in-one.yaml" // basic configuration (release)
	ProjectSample = "data/atlasproject.yaml"
	ClusterSample = "data/atlascluster_basic.yaml"
)

var (
	//default
	Platform   = "kind"
	K8sVersion = "v1.17.17"
)

func TestE2e(t *testing.T) {
	setUpMongoCLI()
	SetDefaultEventuallyTimeout(EventuallyTimeout)
	SetDefaultConsistentlyDuration(ConsistentlyTimeout)
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

//setUpMongoCLI initial setup
func setUpMongoCLI() {
	Platform = os.Getenv("K8s_PLATFORM")
	K8sVersion = os.Getenv("K8s_VERSION")
}
