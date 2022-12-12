package e2e_test

import (
	"os"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = FDescribe("ARM 64 Sanity check", Label("arm64"), func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		if runtime.GOARCH != "arm64" {
			Skip("not an arm64 arch")
		}

		imageURL := os.Getenv("IMAGE_URL")
		Expect(imageURL).ShouldNot(BeEmpty(), "SetUP IMAGE_URL")
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		By("After each.", func() {
			GinkgoWriter.Write([]byte("\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
			GinkgoWriter.Write([]byte("===============================================\n"))
			if CurrentSpecReport().Failed() {
				SaveDump(&data)
				actions.DeleteGlobalKeyIfExist(data)
			}
		})
	})

	DescribeTable("Helm install and Atlas deployment in an ARM64 env",
		func(test model.TestDataProvider, deploymentType string) { // deploymentType - probably will be moved later ()
			data = test

			By("Install CRD using helm", func() {
				helm.InstallCRD(data.Resources)
			})
			By("Install operator using helm", func() {
				helm.InstallOperatorNamespacedSubmodule(data.Resources)
			})

			By("User deploy the deployment via helm", func() {
				helm.InstallDeploymentSubmodule(data.Resources)
			})
			waitDeploymentWithChecks(&data)
			By("Additional check for the current data set", func() {
				for _, check := range data.Actions {
					check(&data)
				}
			})
			deleteDeploymentAndOperator(&data)
		},
		Entry("Advanced deployment by helm chart",
			model.DataProviderWithResources(
				"helm-advanced",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_advanced_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30014,
				[]func(*model.TestDataProvider){},
			),
			"advanced",
		),
	)
})
