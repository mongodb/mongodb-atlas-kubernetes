package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Annotations base test.", Label("test-label"), func() {

	DescribeTable("something",
		func(i int) {
			Expect(i).Should(Equal(10))
		},
		Entry("aaaa", Label("label2"),
			10,
		),
	)
})
