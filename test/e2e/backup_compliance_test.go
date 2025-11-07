// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

var _ = Describe("Backup Compliance Configuration", Label("backup-compliance"), func() {
	var testData *model.TestDataProvider
	var backupCompliancePolicy *akov2.AtlasBackupCompliancePolicy

	BeforeEach(func(ctx SpecContext) {
		By("Setting up cloud environment", func() {
			testData = model.DataProvider(ctx, "atlas-bcp", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30005, []func(*model.TestDataProvider){}).
				WithProject(data.DefaultProject()).
				WithInitialDeployments(data.CreateAdvancedDeployment("bcp-test-deployment").
					WithBackupScheduleRef(common.ResourceRefNamespaced{}))

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
		By("Creating a deployment with a non-compliant backup policy and backup schedule", func() {
			backupPolicy := &akov2.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-bkp-policy", testData.Project.Name),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasBackupPolicySpec{
					Items: []akov2.AtlasBackupPolicyItem{
						{
							FrequencyType:     "hourly",
							FrequencyInterval: 12,
							RetentionUnit:     "months",
							RetentionValue:    1,
						},
					},
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, backupPolicy)).Should(Succeed())

			backupSchedule := &akov2.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-bkp-schedule", testData.Project.Name),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasBackupScheduleSpec{
					PolicyRef: common.ResourceRefNamespaced{
						Name:      backupPolicy.Name,
						Namespace: testData.Resources.Namespace,
					},
					ReferenceHourOfDay:                19,
					ReferenceMinuteOfHour:             2,
					RestoreWindowDays:                 1,
					UseOrgAndGroupNamesInExportPrefix: true,
				},
			}
			Expect(testData.K8SClient.Create(testData.Context, backupSchedule)).Should(Succeed())

			testData.InitialDeployments[0].WithBackupScheduleRef(common.ResourceRefNamespaced{Name: backupSchedule.Name, Namespace: backupSchedule.Namespace})
			deploy.CreateInitialDeployments(testData)
		})

		By("Creating a backup compliance policy in kubernetes", func() {
			backupCompliancePolicy = &akov2.AtlasBackupCompliancePolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-bcp", testData.Project.Name),
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
							FrequencyType:     "daily",
							FrequencyInterval: 2,
							RetentionUnit:     "days",
							RetentionValue:    7,
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

		})

		By("Checking project for appropriate failure message", func() {
			proj := &akov2.AtlasProject{}
			Eventually(func() bool {
				Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), proj)).Should(Succeed())
				return checkForStatusReason(proj.Status.Conditions, api.BackupComplianceReadyType, workflow.ProjectBackupCompliancePolicyNotMet)
			}).WithTimeout(2 * time.Minute).Should(BeTrue())
		})

		By("Changing the BCP to override, and reapply", func() {
			_, err := akoretry.RetryUpdateOnConflict(
				ctx,
				testData.K8SClient,
				client.ObjectKeyFromObject(backupCompliancePolicy),
				func(bcp *akov2.AtlasBackupCompliancePolicy) {
					bcp.Spec.OverwriteBackupPolicies = true
				})
			Expect(err).To(BeNil())
			actions.WaitForConditionsToBecomeTrue(testData, api.BackupComplianceReadyType, api.ReadyType)
		})
	})
})

func checkForStatusReason(conditions []api.Condition, conditionType api.ConditionType, conditionReason workflow.ConditionReason) bool {
	for _, con := range conditions {
		if con.Type == conditionType {
			fmt.Println(con.Reason == string(conditionReason))
			return con.Reason == string(conditionReason)
		}
	}
	return false
}
