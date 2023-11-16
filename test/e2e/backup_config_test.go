package e2e_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
		func(test *model.TestDataProvider) {
			testData = test

			bucket := fmt.Sprintf("%s-%s", bucketName, testData.Resources.TestID)
			bucketPolicy := fmt.Sprintf("%s-%s", atlasBucketPolicyName, testData.Resources.TestID)
			role := fmt.Sprintf("%s-%s", atlasIAMRoleName, testData.Resources.TestID)

			actions.CreateProjectWithCloudProviderAccess(testData, role)
			setupAWSResource(testData.AWSResourcesGenerator, bucket, bucketPolicy, role)
			deploy.CreateInitialDeployments(testData)

			backupConfigFlow(test, bucket)
		},
		Entry(
			"Enable backup for a deployment",
			model.DataProvider(
				"deployment-backup-enabled",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30001,
				[]func(*model.TestDataProvider){},
			).
				WithProject(data.DefaultProject()).
				WithInitialDeployments(data.CreateAdvancedDeployment("backup-deployment")),
		),
	)
})

func backupConfigFlow(data *model.TestDataProvider, bucket string) {
	By("Enable backup for deployment", func() {
		Expect(data.K8SClient.Get(data.Context, client.ObjectKeyFromObject(data.InitialDeployments[0]), data.InitialDeployments[0])).To(Succeed())
		data.InitialDeployments[0].Spec.DeploymentSpec.BackupEnabled = toptr.MakePtr(true)
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
		bkpPolicy := &mdbv1.AtlasBackupPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-policy", data.Project.Name),
			},
			Spec: mdbv1.AtlasBackupPolicySpec{
				Items: []mdbv1.AtlasBackupPolicyItem{
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
				},
			},
		}
		Expect(data.K8SClient.Create(data.Context, bkpPolicy)).To(Succeed())

		bkpSchedule := &mdbv1.AtlasBackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-schedule", data.Project.Name),
			},
			Spec: mdbv1.AtlasBackupScheduleSpec{
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
			data.Project.Status.CloudProviderAccessRoles[0].RoleID,
		)
		Expect(err).Should(BeNil())
		Expect(exportBucket).ShouldNot(BeNil())

		backupSchedule := &mdbv1.AtlasBackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: data.Project.Namespace,
				Name:      fmt.Sprintf("%s-bkp-schedule", data.Project.Name),
			},
		}
		Expect(data.K8SClient.Get(data.Context, client.ObjectKeyFromObject(backupSchedule), backupSchedule)).To(Succeed())

		backupSchedule.Spec.AutoExportEnabled = true
		backupSchedule.Spec.Export = &mdbv1.AtlasBackupExportSpec{
			ExportBucketID: exportBucket.ID,
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

func setupAWSResource(gen *helper.AwsResourcesGenerator, bucket, bucketPolicy, role string) {
	Expect(gen.CreateBucket(bucket)).To(Succeed())
	gen.Cleanup(func() {
		Expect(gen.EmptyBucket(bucket)).To(Succeed())
		Expect(gen.DeleteBucket(bucket)).To(Succeed())
	})

	policyArn, err := gen.CreatePolicy(bucketPolicy, func() helper.IAMPolicy {
		return helper.BucketExportPolicy(bucket)
	})
	Expect(err).Should(BeNil())
	Expect(policyArn).ShouldNot(BeEmpty())
	gen.Cleanup(func() {
		Expect(gen.DeletePolicy(policyArn)).To(Succeed())
	})

	Expect(gen.AttachRolePolicy(role, policyArn)).To(Succeed())
	gen.Cleanup(func() {
		Expect(gen.DetachRolePolicy(role, policyArn)).To(Succeed())
	})
}
