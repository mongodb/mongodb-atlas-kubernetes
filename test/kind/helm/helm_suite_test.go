package e2e_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
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
	control.SkipTestUnless(t, "AKO_KIND_HELM_TEST")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Atlas Operator Helm on Kind Test Suite")
}

var _ = BeforeSuite(func() {
	if !control.Enabled("AKO_KIND_HELM_TEST") {
		fmt.Println("Skipping helm on kind BeforeSuite, AKO_KIND_HELM_TEST is not set")

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
