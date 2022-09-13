package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

const awsRoleNameBase = "atlas-operator-test-aws-role"

var _ = Describe("UserLogin", Label("cloud-access-role"), func() {
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
			SaveDump(&data)
		}
		By("Clean Roles", func() {
			DeleteAllRoles(&data)
		})
		By("Delete Resources, Project with Cloud provider access roles", func() {
			actions.DeleteUserResourcesProject(&data)
			actions.DeleteGlobalKeyIfExist(data)
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, roles []cloudaccess.Role) {
			data = test
			cloudAccessRolesFlow(&data, roles)
		},
		Entry("Test[cloud-access-role-aws-1]: User has project which was updated with AWS custom role", Label("cloud-access-role-aws-1"),
			model.NewTestDataProvider(
				"cloud-access-role-aws-1",
				model.AProject{},
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlasdeployment_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
			[]cloudaccess.Role{
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: v1.CloudProviderAccessRole{
						ProviderName: "AWS",
						// IamAssumedRoleArn will be filled after role creation
					},
				},
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: v1.CloudProviderAccessRole{
						ProviderName: "AWS",
						// IamAssumedRoleArn will be filled after role creation
					},
				},
			},
		),
	)
})

func DeleteAllRoles(data *model.TestDataProvider) {
	project, err := kube.GetProjectResource(data)
	Expect(err).ShouldNot(HaveOccurred())
	errorList := cloudaccess.DeleteRoles(project.Spec.CloudProviderAccessRoles)
	Expect(len(errorList)).Should(Equal(0), errorList)
}

func cloudAccessRolesFlow(userData *model.TestDataProvider, roles []cloudaccess.Role) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData)
		actions.DeployProjectAndWait(userData, 1)
	})

	By("Create AWS role", func() {
		err := cloudaccess.CreateRoles(roles)
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Create project with cloud access role", func() {
		for _, role := range roles {
			userData.Resources.Project.WithCloudAccessRole(role.AccessRole)
		}
		actions.PrepareUsersConfigurations(userData)
		actions.DeployProject(userData)
	})

	By("Establish connection between Atlas and cloud roles", func() {
		Eventually(func() bool {
			return EnsureAllRolesCreated(*userData, len(roles))
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Cloud access roles are not created")
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		err = cloudaccess.AddAtlasStatementToRole(roles, project.Status.CloudProviderAccessRoles)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(kube.GetProjectCloudAccessRolesStatus(userData), "2m", "20s").Should(Equal("True"), "Cloud Access Roles status should be True")
	})

	By("Check cloud access roles status state", func() {
		Eventually(kube.GetReadyProjectStatus(userData)).Should(Equal("True"), "Condition status 'Ready' is not 'True'")
		project, err := kube.GetProjectResource(userData)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(project.Status.CloudProviderAccessRoles).Should(HaveLen(len(roles)))
	})
}

func EnsureAllRolesCreated(data model.TestDataProvider, rolesLen int) bool {
	project, err := kube.GetProjectResource(&data)
	if err != nil {
		return false
	}
	if len(project.Status.CloudProviderAccessRoles) != rolesLen {
		return false
	}
	for _, role := range project.Status.CloudProviderAccessRoles {
		if role.Status != status.StatusCreated {
			return false
		}
	}
	return true
}
