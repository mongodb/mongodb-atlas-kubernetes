package e2e_test

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("UserLogin", Label("x509auth"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})
	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, certRef common.ResourceRefNamespaced) {
			testData = test
			actions.ProjectCreationFlow(test)
			x509Flow(test, &certRef)
		},
		Entry("Test[x509auth]: Can create project and add X.509 Auth to that project", Label("x509auth-basic"),
			model.DataProvider(
				"x509auth",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			common.ResourceRefNamespaced{
				Name: "x509cert",
			},
		),
	)
})

func x509Flow(testData *model.TestDataProvider, certRef *common.ResourceRefNamespaced) {
	By("Create X.509 cert secret", func() {
		Expect(certRef.Name).NotTo(BeEmpty(), "certRef.Name should not be empty")
		if certRef.Namespace == "" {
			certRef.Namespace = testData.Resources.Namespace
		}
		Expect(k8s.CreateCertificateX509(testData.Context, testData.K8SClient, certRef.Name, certRef.Namespace)).To(Succeed())
	})

	By("Add X.509 cert to the project", func() {
		Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name,
			Namespace: testData.Resources.Namespace}, testData.Project)).To(Succeed())
		testData.Project.Spec.X509CertRef = certRef
		Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
		actions.WaitForConditionsToBecomeTrue(testData, status.ReadyType)
	})

	By("Create User with X.509 cert", func() {
		userName := "CN=my-x509-authenticated-user,OU=organizationalunit,O=organization"
		x509User := data.BasicUser("x509user", "user1",
			data.WithReadWriteRole(),
			data.WithX509(userName),
		)
		testData.Users = append(testData.Users, x509User)
		deploy.CreateUsers(testData)
	})

	By("Deploy User", func() {
		By("Check database users Attributes", func() {
			Eventually(actions.CheckUserExistInAtlas(testData), "2m", "10s").Should(BeTrue())
			actions.CheckUsersAttributes(testData)
		})
	})
}
