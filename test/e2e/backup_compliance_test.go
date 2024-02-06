package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Backup Compliance Configuration", Label("backup-compliance"), func() {
	var testData *model.TestDataProvider

	BeforeEach(func() {
		By("Setting up cloud environment", func() {
			actions.ProjectCreationFlow(testData)
		})
	})

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveUsersToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}

		By("Should clean up created resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)

			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Configures a backup compliance policy", func() {
		By("Creating a backup compliance policy in kubernetes", func() {
			compliancePolicy := v1.AtlasBackupCompliancePolicy{
				Spec: v1.AtlasBackupCompliancePolicySpec{
					AuthorizedEmail:         "test@example.com",
					AuthorizedUserFirstName: "John",
					AuthorizedUserLastName:  "Doe",
					CopyProtectionEnabled:   false,
					EncryptionAtRestEnabled: false,
					PITEnabled:              false,
					RestoreWindowDays:       42,
					ScheduledPolicyItems: []v1.AtlasBackupPolicyItem{
						{
							FrequencyType:     "monthly",
							FrequencyInterval: 4,
							RetentionUnit:     "days",
							RetentionValue:    14,
						},
					},
					OnDemandPolicy: v1.AtlasBackupPolicyItem{
						FrequencyType:     "weekly",
						FrequencyInterval: 1,
						RetentionUnit:     "weeks",
						RetentionValue:    3,
					},
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, &compliancePolicy)).Should(Succeed())
			actions.WaitForConditionsToBecomeTrue(testData, status.BackupComplianceReadyType, status.ReadyType)
		})
	})
})
