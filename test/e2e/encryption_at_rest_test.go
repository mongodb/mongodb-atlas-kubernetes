package e2e_test

import (
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("UserLogin", Label("encryption-at-rest"), func() {
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
		By("Clean Cloud", func() {
			DeleteAllRoles(testData)
		})

		By("Delete Resources, Project with Encryption at rest", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, encAtRest v1.EncryptionAtRest, roles []cloudaccess.Role) {
			testData = test
			actions.ProjectCreationFlow(test)
			encryptionAtRestFlow(test, encAtRest, roles)
		},
		Entry("Test[encryption-at-rest-aws]: Can add Encryption at Rest to AWS project", Label("encryption-at-rest-aws"),
			model.DataProvider(
				"encryption-at-rest-aws",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			v1.EncryptionAtRest{
				AwsKms: v1.AwsKms{
					Enabled: toptr.MakePtr(true),
					// CustomerMasterKeyID: "",
					Region: "US_EAST_1",
					Valid:  toptr.MakePtr(true),
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
	By("Add cloud access role (AWS only)", func() {
		cloudAccessRolesFlow(userData, roles)
	})

	By("Create KMS", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

		Expect(len(userData.Project.Status.CloudProviderAccessRoles)).NotTo(Equal(0))
		aRole := userData.Project.Status.CloudProviderAccessRoles[0]

		fillKMSforAWS(&encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)
		fillVaultforAzure(&encAtRest)
		fillKMSforGCP(&encAtRest)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.EncryptionAtRest = &encAtRest
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.EncryptionAtRestReadyType, status.ReadyType)
	})

	By("Remove Encryption at Rest from the project", func() {
		removeAllEncryptionsSeparately(&encAtRest)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.EncryptionAtRest = &encAtRest
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Check if project returned back to the initial state", func() {
		actions.CheckProjectConditionsNotSet(userData, status.EncryptionAtRestReadyType)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

		Eventually(func(g Gomega) bool {
			areEmpty, err := checkIfEncryptionsAreDisabled(userData.Project.ID())
			g.Expect(err).ShouldNot(HaveOccurred())
			return areEmpty
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).
			Should(BeTrue(), "Encryption at Rest is not disabled")
	})
}

func fillKMSforAWS(encAtRest *v1.EncryptionAtRest, atlasAccountArn, assumedRoleArn string) {
	if (encAtRest.AwsKms == v1.AwsKms{}) {
		return
	}

	Expect(encAtRest.AwsKms.Region).NotTo(Equal(""))
	awsAction := cloud.NewAwsAction()
	CustomerMasterKeyID, err := awsAction.CreateKMS(config.AWSRegionUS, atlasAccountArn, assumedRoleArn)
	Expect(err).ToNot(HaveOccurred())
	Expect(CustomerMasterKeyID).NotTo(Equal(""))

	encAtRest.AwsKms.CustomerMasterKeyID = CustomerMasterKeyID
}

func fillVaultforAzure(encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.AzureKeyVault == v1.AzureKeyVault{}) {
		return
	}

	// todo: fill in
}

func fillKMSforGCP(encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.GoogleCloudKms == v1.GoogleCloudKms{}) {
		return
	}

	// todo: fill in
}

func removeAllEncryptionsSeparately(encAtRest *v1.EncryptionAtRest) {
	encAtRest.AwsKms = v1.AwsKms{}
	encAtRest.AzureKeyVault = v1.AzureKeyVault{}
	encAtRest.GoogleCloudKms = v1.GoogleCloudKms{}
}

func checkIfEncryptionsAreDisabled(projectID string) (areEmpty bool, err error) {
	atlasClient := atlas.GetClientOrFail()
	encryptionAtRest, err := atlasClient.GetEncryptioAtRest(projectID)
	if err != nil {
		return false, err
	}

	if encryptionAtRest == nil {
		return true, nil
	}

	awsEnabled := *encryptionAtRest.AwsKms.Enabled
	azureEnabled := *encryptionAtRest.AzureKeyVault.Enabled
	gcpEnabled := *encryptionAtRest.GoogleCloudKms.Enabled

	if awsEnabled || azureEnabled || gcpEnabled {
		return false, nil
	}

	return true, nil
}
