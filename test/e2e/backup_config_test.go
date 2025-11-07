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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	helper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

const (
	atlasIAMRoleName      = "atlas-role"
	atlasBucketPolicyName = "atlas-bucket-export-policy"
	bucketName            = "cloud-backup-snapshot"
)

var _ = Describe("Deployment Backup Configuration", Label("backup-config"), func() {
	var testData *model.TestDataProvider

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

	DescribeTable("Configure backup for a deployment",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider) {
			testData = test(ctx)

			bucket := fmt.Sprintf("%s-%s", bucketName, testData.Resources.TestID)
			bucketPolicy := fmt.Sprintf("%s-%s", atlasBucketPolicyName, testData.Resources.TestID)
			role := fmt.Sprintf("%s-%s", atlasIAMRoleName, testData.Resources.TestID)

			actions.CreateProjectWithCloudProviderAccess(ctx, testData, role)
			setupAWSResource(ctx, testData.AWSResourcesGenerator, bucket, bucketPolicy, role)
			deploy.CreateInitialDeployments(testData)

			backupConfigFlow(testData, bucket)
		},
		Entry(
			"Enable backup for a deployment",
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "deployment-backup-enabled", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30001, []func(*model.TestDataProvider){}).
					WithProject(data.DefaultProject()).
					WithInitialDeployments(data.CreateAdvancedDeployment("backup-deployment"))
			},
		),
	)
})

func backupConfigFlow(data *model.TestDataProvider, bucket string) {
	By("Enable backup for deployment", func() {
		Expect(data.K8SClient.Get(data.Context, client.ObjectKeyFromObject(data.InitialDeployments[0]), data.InitialDeployments[0])).To(Succeed())
		data.InitialDeployments[0].Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
		Expect(data.K8SClient.Update(data.Context, data.InitialDeployments[0])).To(Succeed())

		Eventually(func(g Gomega) bool {
			objectKey := types.NamespacedName{
				Name:      data.InitialDeployments[0].Name,
				Namespace: data.InitialDeployments[0].Namespace,
			}
			g.Expect(data.K8SClient.Get(data.Context, objectKey, data.InitialDeployments[0])).To(Succeed())
			return data.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})

	By("Configure backup schedule and policy for the deployment", func() {
		bkpPolicy := &akov2.AtlasBackupPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-policy", data.Project.Name),
			},
			Spec: akov2.AtlasBackupPolicySpec{
				Items: []akov2.AtlasBackupPolicyItem{
					{
						FrequencyInterval: 6,
						FrequencyType:     "hourly",
						RetentionValue:    2,
						RetentionUnit:     "days",
					},
					{
						FrequencyInterval: 1,
						FrequencyType:     "daily",
						RetentionValue:    7,
						RetentionUnit:     "days",
					},
					{
						FrequencyInterval: 1,
						FrequencyType:     "weekly",
						RetentionValue:    4,
						RetentionUnit:     "weeks",
					},
					{
						FrequencyInterval: 1,
						FrequencyType:     "monthly",
						RetentionValue:    12,
						RetentionUnit:     "months",
					},
					{
						FrequencyInterval: 1,
						FrequencyType:     "yearly",
						RetentionValue:    1,
						RetentionUnit:     "years",
					},
				},
			},
		}
		Expect(data.K8SClient.Create(data.Context, bkpPolicy)).To(Succeed())

		bkpSchedule := &akov2.AtlasBackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-schedule", data.Project.Name),
			},
			Spec: akov2.AtlasBackupScheduleSpec{
				PolicyRef: common.ResourceRefNamespaced{
					Namespace: data.Project.Namespace,
					Name:      fmt.Sprintf("%s-bkp-policy", data.Project.Name),
				},
				ReferenceHourOfDay:                19,
				ReferenceMinuteOfHour:             2,
				RestoreWindowDays:                 1,
				UseOrgAndGroupNamesInExportPrefix: true,
			},
		}
		Expect(data.K8SClient.Create(data.Context, bkpSchedule)).To(Succeed())

		Expect(data.K8SClient.Get(data.Context, client.ObjectKeyFromObject(data.InitialDeployments[0]), data.InitialDeployments[0])).To(Succeed())
		data.InitialDeployments[0].Spec.BackupScheduleRef = common.ResourceRefNamespaced{
			Namespace: data.Project.Namespace,
			Name:      fmt.Sprintf("%s-bkp-schedule", data.Project.Name),
		}
		Expect(data.K8SClient.Update(data.Context, data.InitialDeployments[0])).To(Succeed())

		Eventually(func(g Gomega) bool {
			g.Expect(data.K8SClient.Get(data.Context, types.NamespacedName{
				Name:      data.InitialDeployments[0].Name,
				Namespace: data.InitialDeployments[0].Namespace,
			}, data.InitialDeployments[0])).To(Succeed())

			return data.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})

	By("Configure auto export to AWS bucket", func() {
		aClient := atlas.GetClientOrFail()
		exportBucket, err := aClient.CreateExportBucket(
			data.Project.ID(),
			bucket,
			data.Project.Status.CloudProviderIntegrations[0].RoleID,
		)
		Expect(err).Should(BeNil())
		Expect(exportBucket).ShouldNot(BeNil())

		backupSchedule := &akov2.AtlasBackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-schedule", data.Project.Name),
			},
		}
		Expect(data.K8SClient.Get(data.Context, client.ObjectKeyFromObject(backupSchedule), backupSchedule)).To(Succeed())

		backupSchedule.Spec.AutoExportEnabled = true
		backupSchedule.Spec.Export = &akov2.AtlasBackupExportSpec{
			ExportBucketID: exportBucket.GetId(),
			FrequencyType:  "monthly",
		}
		Expect(data.K8SClient.Update(data.Context, backupSchedule)).To(Succeed())

		Eventually(func(g Gomega) bool {
			g.Expect(data.K8SClient.Get(data.Context, types.NamespacedName{
				Name:      data.InitialDeployments[0].Name,
				Namespace: data.InitialDeployments[0].Namespace,
			}, data.InitialDeployments[0])).To(Succeed())

			return data.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})
}

func setupAWSResource(ctx context.Context, gen *helper.AwsResourcesGenerator, bucket, bucketPolicy, role string) {
	Expect(gen.CreateBucket(ctx, bucket)).To(Succeed())
	DeferCleanup(func(ctx SpecContext) {
		Expect(gen.EmptyBucket(ctx, bucket)).To(Succeed())
		Expect(gen.DeleteBucket(ctx, bucket)).To(Succeed())
	})

	policyArn, err := gen.CreatePolicy(ctx, bucketPolicy, func() helper.IAMPolicy {
		return helper.BucketExportPolicy(bucket)
	})
	Expect(err).Should(BeNil())
	Expect(policyArn).ShouldNot(BeEmpty())
	DeferCleanup(func(ctx SpecContext) {
		Expect(gen.DeletePolicy(ctx, policyArn)).To(Succeed())
	})

	Expect(gen.AttachRolePolicy(ctx, role, policyArn)).To(Succeed())
	DeferCleanup(func(ctx SpecContext) {
		Expect(gen.DetachRolePolicy(ctx, role, policyArn)).To(Succeed())
	})
}
