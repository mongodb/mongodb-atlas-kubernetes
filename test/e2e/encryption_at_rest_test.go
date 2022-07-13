package e2e_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("UserLogin", Label("encryption-at-rest"), func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnvironment()
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
		By("Clean Cloud", func() {
			DeleteAllRoles(&data)
		})

		By("Delete Resources, Project with PEService", func() {
			actions.DeleteUserResourcesProject(&data)
			actions.DeleteGlobalKeyIfExist(data)
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, encAtRest v1.EncryptionAtRest, roles []cloudaccess.Role) {
			data = test
			encryptionAtRestFlow(&data, encAtRest, roles)
		},
		Entry("Test[encryption-at-rest-aws]: Can create project and add Encryption at Rest to that project", Label("encryption-at-rest-aws"),
			model.NewTestDataProvider(
				"encryption-at-rest-aws",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_standard.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
			v1.EncryptionAtRest{
				AwsKms: v1.AwsKms{
					Enabled: toptr.Boolptr(true),
					// CustomerMasterKeyID: "",
					Region: "US_EAST_1",
					Valid:  toptr.Boolptr(true),
				},
			},
			[]cloudaccess.Role{
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: v1.CloudProviderAccessRole{
						ProviderName: "AWS",
					},
				},
			},
		),
	)
})

func encryptionAtRestFlow(userData *model.TestDataProvider, encAtRest v1.EncryptionAtRest, roles []cloudaccess.Role) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData)
		actions.DeployProjectAndWait(userData, "1")
	})

	By("Add cloud access role (AWS only)", func() {
		if len(roles) == 0 {
			return
		}

		err := cloudaccess.CreateRoles(roles)
		Expect(err).ShouldNot(HaveOccurred())

		for _, role := range roles {
			userData.Resources.Project.WithCloudAccessRole(role.AccessRole)
		}

		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData)

		Eventually(func() bool {
			return EnsureAllRolesCreated(*userData, len(roles))
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Cloud access roles are not created")

		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())

		err = cloudaccess.AddAtlasStatementToRole(roles, project.Status.CloudProviderAccessRoles)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(kube.GetProjectCloudAccessRolesStatus(userData), "2m", "20s").Should(Equal("True"), "Cloud Access Roles status should be True")
	})

	By("Create KMS", func() {
		Expect(encAtRest.AwsKms.Region).NotTo(Equal(""))
		awsAction := cloud.NewAwsAction()
		CustomerMasterKeyID, err := awsAction.CreateKMS(config.AWSRegionUS)
		Expect(err).ToNot(HaveOccurred())
		Expect(CustomerMasterKeyID).NotTo(Equal(""))

		encAtRest.AwsKms.CustomerMasterKeyID = CustomerMasterKeyID

		userData.Resources.Project.WithEncryptionAtRest(encAtRest)
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData)
	})

	By("Check Encryption at Rest status", func() {
		_, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Schedule key deletion", func() {
		keyID := encAtRest.AwsKms.CustomerMasterKeyID
		awsAction := cloud.NewAwsAction()
		err := awsAction.DeleteKMS(keyID, encAtRest.AwsKms.Region)
		Expect(err).ShouldNot(HaveOccurred())
	})
}
