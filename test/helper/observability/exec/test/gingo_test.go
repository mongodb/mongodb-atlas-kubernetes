package test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Top-Level", func() {
	Describe("Describe2", func() {
		Context("Context2", func() {
			It("It2.1", func() {
				GinkgoWriter.Println("It2.1 GinkgoParallelProcess()", GinkgoParallelProcess())
				//time.Sleep(1 * time.Second)
				GinkgoWriter.Println("It2.1 GinkgoParallelProcess()", GinkgoParallelProcess())
			})
			It("It2.2", func() {
				GinkgoWriter.Println("It2.2 GinkgoParallelProcess()", GinkgoParallelProcess())
				//time.Sleep(1 * time.Second)
				GinkgoWriter.Println("It2.2 GinkgoParallelProcess()", GinkgoParallelProcess())
			})
			It("It2.3", func() {
				GinkgoWriter.Println("It2.3 GinkgoParallelProcess()", GinkgoParallelProcess())
				//time.Sleep(1 * time.Second)
				GinkgoWriter.Println("It2.3 GinkgoParallelProcess()", GinkgoParallelProcess())
			})
		})
	})
})
