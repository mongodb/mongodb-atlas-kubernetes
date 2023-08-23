package int

import (
	"context"
	"fmt"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// nolint:dupl
var _ = Describe("AtlasBackupSchedule Deletion Protected",
	Ordered,
	Label("AtlasDeployment", "AtlasBackupSchedule", "deletion-protection", "deletion-protection-backup"), func() {
		var testNamespace *corev1.Namespace
		var stopManager context.CancelFunc
		var connectionSecret corev1.Secret
		var testProject *mdbv1.AtlasProject
		var testDeployment *mdbv1.AtlasDeployment

		BeforeAll(func() {
			By("Starting the operator with protection ON", func() {
				testNamespace, stopManager = prepareControllers(true)
				Expect(testNamespace).ToNot(BeNil())
				Expect(stopManager).ToNot(BeNil())
			})

			By("Creating project connection secret", func() {
				connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
				Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Creating a project with deletion annotation", func() {
				testProject = mdbv1.DefaultProject(testNamespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
				customresource.SetAnnotation( // this test project must be deleted
					testProject,
					customresource.ResourcePolicyAnnotation,
					customresource.ResourcePolicyDelete,
				)
				Expect(k8sClient.Create(context.TODO(), testProject, &client.CreateOptions{})).To(Succeed())

				Eventually(func() bool {
					return testutil.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
				}).WithTimeout(3 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})

		AfterAll(func() {
			By("Deleting deployment from Atlas", func() {
				if testDeployment == nil {
					return
				}

				Expect(deleteAtlasDeployment(testProject.Status.ID, testDeployment.Spec.DeploymentSpec.Name)).ToNot(HaveOccurred())
			})
			By("Deleting project from k8s and atlas", func() {
				Expect(k8sClient.Delete(context.TODO(), testProject, &client.DeleteOptions{})).To(Succeed())
				Eventually(
					checkAtlasProjectRemoved(testProject.Status.ID),
				).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Deleting project connection secret", func() {
				Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
			})

			By("Stopping the operator", func() {
				stopManager()
				err := k8sClient.Delete(context.Background(), testNamespace)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("Should not process BackupSchedule with deletion protection ON", func() {
			var bsPolicy *mdbv1.AtlasBackupPolicy
			var bsSchedule *mdbv1.AtlasBackupSchedule
			By("Creating AtlasBackupPolicy resource", func() {
				bsPolicy = &mdbv1.AtlasBackupPolicy{
					TypeMeta: metav1.TypeMeta{
						Kind:       "atlas.mongodb.com/v1",
						APIVersion: "AtlasBackupPolicy",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:        "bs-policy",
						Namespace:   testNamespace.Name,
						Labels:      map[string]string{},
						Annotations: map[string]string{},
					},
					Spec: mdbv1.AtlasBackupPolicySpec{
						Items: []mdbv1.AtlasBackupPolicyItem{
							{
								FrequencyType:     "daily",
								FrequencyInterval: 5,
								RetentionUnit:     "days",
								RetentionValue:    20,
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), bsPolicy)).To(Succeed())
			})

			By("Creating AtlasBackupSchedule resource", func() {
				bsSchedule = &mdbv1.AtlasBackupSchedule{
					TypeMeta: metav1.TypeMeta{
						Kind:       "atlas.mongodb.com/v1",
						APIVersion: "AtlasBackupSchedule",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bs-schedule",
						Namespace: testNamespace.Name,
					},
					Spec: mdbv1.AtlasBackupScheduleSpec{
						AutoExportEnabled: false,
						PolicyRef: common.ResourceRefNamespaced{
							Name:      bsPolicy.Name,
							Namespace: bsPolicy.Namespace,
						},
						ReferenceHourOfDay:                10,
						ReferenceMinuteOfHour:             10,
						RestoreWindowDays:                 10,
						UpdateSnapshots:                   false,
						UseOrgAndGroupNamesInExportPrefix: false,
						CopySettings:                      []mdbv1.CopySetting{},
					},
				}
				Expect(k8sClient.Create(context.Background(), bsSchedule)).To(Succeed())
			})

			By("Creating a deployment with backups enabled (default)", func() {
				testDeployment = mdbv1.DefaultAWSDeployment(testNamespace.Name, testProject.Name)
				testDeployment.Spec.DeploymentSpec.ProviderBackupEnabled = toptr.MakePtr(true)
				Expect(k8sClient.Create(context.Background(), testDeployment)).To(Succeed())
			})

			By("Deployment should be Ready", func() {
				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(k8sClient, testDeployment, status.TrueCondition(status.ReadyType), validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})

			By("Add custom BackupSchedule for the Deployment", func() {
				Eventually(func(g Gomega) {
					deployment := &mdbv1.AtlasDeployment{}
					g.Expect(
						k8sClient.Get(context.Background(),
							kube.ObjectKeyFromObject(testDeployment),
							deployment),
					).To(Succeed())

					deployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
						Name:      bsSchedule.Name,
						Namespace: bsSchedule.Namespace,
					}

					g.Expect(k8sClient.Update(context.Background(), deployment)).To(Succeed())
				}).WithTimeout(2 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
			})

			By("Deployment should NOT be Ready", func() {
				Eventually(func(g Gomega) bool {
					return testutil.CheckCondition(
						k8sClient,
						testDeployment,
						status.FalseCondition(status.DeploymentReadyType).
							WithReason(string(workflow.Internal)).
							WithMessageRegexp(atlasdeployment.BackupProtected),
						validateDeploymentUpdatingFunc(g))
				}).WithTimeout(30 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
			})
		})
	})
