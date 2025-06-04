package test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"

	akoginkgo "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability/ginkgo"
)

func TestBooks(t *testing.T) {
	akoginkgo.RegisterCallbacks()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}
