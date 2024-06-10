package e2e_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

var _ = Describe("Backup Compliance Configuration", Label("backup-compliance"), func() {
	var testData *model.TestDataProvider
	var backupCompliancePolicy *akov2.AtlasBackupCompliancePolicy

	BeforeEach(func() {
		By("Setting up cloud environment", func() {
			testData = model.DataProvider(
				"atlas-bcp",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30005,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject())
			actions.ProjectCreationFlow(testData)
		})
	})

	AfterEach(func() {
		By("Should clean up created resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)

			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Configures a backup compliance policy", func(ctx context.Context) {
		By("Creating a backup compliance policy in kubernetes", func() {
			backupCompliancePolicy = &akov2.AtlasBackupCompliancePolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-bcp",
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasBackupCompliancePolicySpec{
					AuthorizedEmail:         "test@example.com",
					AuthorizedUserFirstName: "John",
					AuthorizedUserLastName:  "Doe",
					CopyProtectionEnabled:   false,
					EncryptionAtRestEnabled: false,
					PITEnabled:              false,
					RestoreWindowDays:       42,
					ScheduledPolicyItems: []akov2.AtlasBackupPolicyItem{
						{
							FrequencyType:     "monthly",
							FrequencyInterval: 4,
							RetentionUnit:     "months",
							RetentionValue:    1,
						},
					},
					OnDemandPolicy: akov2.AtlasOnDemandPolicy{
						RetentionUnit:  "weeks",
						RetentionValue: 3,
					},
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, backupCompliancePolicy)).Should(Succeed())
		})
		By("Adding BCP to a Project", func() {
			_, err := akoretry.RetryUpdateOnConflict(ctx, testData.K8SClient, client.ObjectKeyFromObject(testData.Project), func(project *akov2.AtlasProject) {
				project.Spec.BackupCompliancePolicyRef = &common.ResourceRefNamespaced{
					Name:      backupCompliancePolicy.Name,
					Namespace: backupCompliancePolicy.Namespace,
				}
			})
			Expect(err).To(BeNil())
			actions.WaitForConditionsToBecomeTrue(testData, api.BackupComplianceReadyType, api.ReadyType)
		})
	})
})
