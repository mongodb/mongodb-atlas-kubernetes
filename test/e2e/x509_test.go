package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	common "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("UserLogin", Label("x509auth"), func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			By("Save logs to output directory ", func() {
				GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
					[]byte(kubecli.DescribeOperatorPod(data.Resources.Namespace)),
				)
				utils.SaveToFile(
					fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
					kubecli.GetManagerLogs(data.Resources.Namespace),
				)
				actions.SaveTestAppLogs(data.Resources)
				actions.SaveK8sResources(
					[]string{"deploy", "atlasprojects"},
					data.Resources.Namespace,
				)
			})

		}
		By("Delete Resources", func() {
			actions.DeleteUserResourcesProject(&data)
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, certRef common.ResourceRefNamespaced) {
			data = test
			x509Flow(&data, &certRef)
		},
		Entry("Test[x509auth]: Can create project and add X.509 Auth to that project", Label("x509auth-basic"),
			model.NewTestDataProvider(
				"x509auth",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_standard.yaml"},
				[]string{},
				[]model.DBUser{},
				30000,
				[]func(*model.TestDataProvider){},
			),
			common.ResourceRefNamespaced{
				Name:      "x509cert",
				Namespace: data.Resources.Namespace,
			},
		),
	)
})

func x509Flow(data *model.TestDataProvider, certRef *common.ResourceRefNamespaced) {
	By("Deploy Project with standart configuration", func() {
		actions.PrepareUsersConfigurations(data)
		deploy.NamespacedOperator(data)
		actions.DeployProjectAndWait(data, "1")
	})

	By("Create X.509 cert secret", func() {
		kubecli.CreateX509Secret(certRef.Name, certRef.Namespace)
	})

	By("Add X.509 cert to the project", func() {
		data.Resources.Project.WithX509(certRef)
		actions.DeployProject(data, "2")
	})

	By("Check if project statuses are updating, get project ID", func() {
		Eventually(kube.GetReadyProjectStatus(data)).Should(Equal("True"),
			"Atlasproject status.conditions are not True")

		actions.UpdateProjectID(data)
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})

	By("Create User with X.509 cert", func() {
		userName := "CN=my-x509-authenticated-user,OU=organizationalunit,O=organization"
		x509User := model.NewDBUser("my-x509-user").
			WithProjectRef(data.Resources.Project.GetK8sMetaName()).
			AddBuildInReadWriteRole().
			WithX509(userName)
		data.Resources.Users = append(data.Resources.Users, *x509User)
		actions.PrepareUsersConfigurations(data)
		actions.DeployUsers(data)
	})
}
