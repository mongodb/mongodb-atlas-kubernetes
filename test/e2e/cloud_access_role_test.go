package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

const awsRoleNameBase = "atlas-operator-test-aws-role"

var _ = Describe("UserLogin", Label("cloud-access-role"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
		checkUpAWSEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Roles", func() {
			DeleteAllRoles(testData)
		})
		By("Delete Resources, Project with Cloud provider access roles", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, roles []cloudaccess.Role) {
			testData = test
			actions.ProjectCreationFlow(test)
			cloudAccessRolesFlow(test, roles)
		},
		Entry("Test[cloud-access-role-aws-1]: User has project which was updated with AWS custom role", Label("cloud-access-role-aws-1"),
			model.DataProvider(
				"cloud-access-role-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
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

func DeleteAllRoles(testData *model.TestDataProvider) {
	Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).Should(Succeed())
	errorList := cloudaccess.DeleteRoles(testData.Project.Spec.CloudProviderAccessRoles)
	Expect(len(errorList)).Should(Equal(0), errorList)
}

func cloudAccessRolesFlow(userData *model.TestDataProvider, roles []cloudaccess.Role) {
	By("Create AWS role", func() {
		err := cloudaccess.CreateRoles(roles)
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Create project with cloud access role", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		for _, role := range roles {
			userData.Project.Spec.CloudProviderAccessRoles = append(userData.Project.Spec.CloudProviderAccessRoles, role.AccessRole)
		}
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Establish connection between Atlas and cloud roles", func() {
		Eventually(func(g Gomega) bool {
			return EnsureAllRolesCreated(g, *userData, len(roles))
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Cloud access roles are not created")
		project := &v1.AtlasProject{}
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, project)).Should(Succeed())
		err := cloudaccess.AddAtlasStatementToRole(roles, project.Status.CloudProviderAccessRoles)
		Expect(err).ShouldNot(HaveOccurred())
		actions.WaitForConditionsToBecomeTrue(userData, status.CloudProviderAccessReadyType, status.ReadyType)
	})

	By("Check cloud access roles status state", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.CloudProviderAccessRoles).Should(HaveLen(len(roles)))
	})
}

func EnsureAllRolesCreated(g Gomega, testData model.TestDataProvider, rolesLen int) bool {
	project := &v1.AtlasProject{}
	g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, project)).Should(Succeed())

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
