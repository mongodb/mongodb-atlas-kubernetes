package e2e_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"

	"github.com/google/uuid"
	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("Project Deletion Protection", Label("project", "deletion-protection"), func() {
	var testData *model.TestDataProvider
	var managerStop context.CancelFunc
	var projectID, networkPeerID, awsRoleARN, awsAccountID, AwsVpcID string
	ctx := context.Background()

	BeforeEach(func() {
		checkUpAWSEnvironment()
		Expect(os.Getenv("PAGER_DUTY_SERVICE_KEY")).
			ShouldNot(BeEmpty(), "Please, setup PAGER_DUTY_SERVICE_KEY environment variable for test integration with Pager Duty")

		testData = model.DataProvider(
			"project-deletion-protection",
			model.NewEmptyAtlasKeyType().CreateAsGlobalLevelKey(),
			30005,
			[]func(*model.TestDataProvider){},
		)

		actions.CreateNamespaceAndSecrets(testData)

		managerStart, err := k8s.RunManager(
			k8s.WithGlobalKey(client.ObjectKey{Namespace: testData.Resources.Namespace, Name: config.DefaultOperatorGlobalKey}),
			k8s.WithNamespaces(testData.Resources.Namespace),
			k8s.WithObjectDeletionProtection(true),
			k8s.WithSubObjectDeletionProtection(true),
		)
		Expect(err).ToNot(HaveOccurred())

		cancelCtx, cancel := context.WithCancel(ctx)
		managerStop = cancel
		go func() {
			err := managerStart(cancelCtx)
			Expect(err).ToNot(HaveOccurred())
		}()
	})

	It("Reconcile Atlas Project when deletion protection is enabled", func() {
		projectName := fmt.Sprintf("project-deletion-protection-e2e-%s", uuid.New().String()[0:6])

		By("Creating a project outside the operator", func() {
			atlasProject, _, err := atlasClient.Client.Projects.Create(
				ctx, &mongodbatlas.Project{
					OrgID: os.Getenv("MCLI_ORG_ID"),
					Name:  projectName,
				},
				&mongodbatlas.CreateProjectOptions{},
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(atlasProject).ToNot(BeNil())

			projectID = atlasProject.ID
		})

		By("Adding IP Access List entry to the project", func() {
			_, _, err := atlasClient.Client.ProjectIPAccessList.Create(
				ctx,
				projectID,
				[]*mongodbatlas.ProjectIPAccessList{
					{
						CIDRBlock: "192.168.0.0/24",
						GroupID:   projectID,
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding Cloud Provider Access to the project", func() {
			assumedRoleArn, err := cloudaccess.CreateAWSIAMRole(projectName)
			Expect(err).ToNot(HaveOccurred())
			awsRoleARN = assumedRoleArn

			cloudProvider, _, err := atlasClient.Client.CloudProviderAccess.CreateRole(
				ctx,
				projectID,
				&mongodbatlas.CloudProviderAccessRoleRequest{
					ProviderName: "AWS",
				},
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(cloudaccess.AddAtlasStatementToAWSIAMRole(cloudProvider.AtlasAWSAccountARN, cloudProvider.AtlasAssumedRoleExternalID, projectName)).
				To(Succeed())

			Eventually(func(g Gomega) {
				_, _, err := atlasClient.Client.CloudProviderAccess.AuthorizeRole(
					ctx,
					projectID,
					cloudProvider.RoleID,
					&mongodbatlas.CloudProviderAccessRoleRequest{
						ProviderName:      "AWS",
						IAMAssumedRoleARN: toptr.MakePtr(assumedRoleArn),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
			}).WithTimeout(time.Minute).WithPolling(time.Second * 15).Should(Succeed())
		})

		By("Adding Network peering to the project", func() {
			aws, err := cloud.NewAWSAction(GinkgoT())
			Expect(err).ToNot(HaveOccurred())

			awsAccountID, err = aws.GetAccountID()
			Expect(err).ToNot(HaveOccurred())

			AwsVpcID, err = aws.InitNetwork(projectName, "10.0.0.0/24", "eu-west-2", map[string]string{"subnet-1": "10.0.0.0/24"}, true)
			Expect(err).ToNot(HaveOccurred())

			c, _, err := atlasClient.Client.Containers.Create(ctx, projectID, &mongodbatlas.Container{
				ProviderName:   "AWS",
				RegionName:     "EU_WEST_2",
				AtlasCIDRBlock: "192.168.224.0/21",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(c).ToNot(BeNil())

			p, _, err := atlasClient.Client.Peers.Create(ctx, projectID, &mongodbatlas.Peer{
				ProviderName:        "AWS",
				AccepterRegionName:  "eu-west-2",
				ContainerID:         c.ID,
				AWSAccountID:        awsAccountID,
				RouteTableCIDRBlock: "10.0.0.0/24",
				VpcID:               AwsVpcID,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(p).ToNot(BeNil())

			Eventually(func(g Gomega) {
				p, _, err = atlasClient.Client.Peers.Get(ctx, projectID, p.ID)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(p).ToNot(BeNil())
				g.Expect(p.StatusName).To(Equal("PENDING_ACCEPTANCE"))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			Expect(aws.AcceptVpcPeeringConnection(p.ConnectionID, "eu-west-2")).To(Succeed())

			Eventually(func(g Gomega) {
				pCheck, _, err := atlasClient.Client.Peers.Get(ctx, projectID, p.ID)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(pCheck).ToNot(BeNil())
				g.Expect(pCheck.StatusName).To(Equal("AVAILABLE"))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			networkPeerID = p.ID
		})

		By("Adding integration to the project", func() {
			_, _, err := atlasClient.Client.Integrations.Create(
				ctx,
				projectID,
				"PAGER_DUTY",
				&mongodbatlas.ThirdPartyIntegration{
					Type:       "PAGER_DUTY",
					Region:     "EU",
					ServiceKey: os.Getenv("PAGER_DUTY_SERVICE_KEY"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Creating a project to be managed by the operator", func() {
			akoProject := &mdbv1.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projectName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasProjectSpec{
					Name: projectName,
					CloudProviderAccessRoles: []mdbv1.CloudProviderAccessRole{
						{
							ProviderName:      "AWS",
							IamAssumedRoleArn: "an-aws-role-arn",
						},
					},
					ProjectIPAccessList: []project.IPAccessList{
						{
							CIDRBlock: "10.1.1.0/24",
						},
					},
					NetworkPeers: []mdbv1.NetworkPeer{
						{
							ProviderName:        "AWS",
							AccepterRegionName:  "eu-west-2",
							AtlasCIDRBlock:      "192.168.224.0/21",
							AWSAccountID:        awsAccountID,
							RouteTableCIDRBlock: "10.0.0.0/24",
							VpcID:               "wrong",
						},
					},
					Integrations: []project.Integration{
						{
							Type:   "DATADOG",
							Region: "EU",
							ServiceKeyRef: common.ResourceRefNamespaced{
								Name: "pager-duty-service-key",
							},
						},
					},
				},
			}
			testData.Project = akoProject

			Expect(testData.K8SClient.Create(ctx, testData.Project))
			time.Sleep(time.Second * 30)
		})

		By("Project is ready by all sub-resources are in conflict", func() {
			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.FalseCondition(status.IPAccessListReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile IP Access List due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.CloudProviderAccessReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.NetworkPeerReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Network Peering due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.IntegrationReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("IP Access List is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.ProjectIPAccessList[0].CIDRBlock = "192.168.0.0/24"
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.FalseCondition(status.CloudProviderAccessReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Cloud Provider Access due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.NetworkPeerReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Network Peering due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.IntegrationReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Cloud Provider Access is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.CloudProviderAccessRoles[0].IamAssumedRoleArn = awsRoleARN
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

				g.Expect(testData.Project.Status.CloudProviderAccessRoles).ToNot(HaveLen(0))
				g.Expect(testData.Project.Status.CloudProviderAccessRoles[0].Status).To(Equal("AUTHORIZED"))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			Expect(
				cloudaccess.AddAtlasStatementToAWSIAMRole(
					testData.Project.Status.CloudProviderAccessRoles[0].AtlasAWSAccountArn,
					testData.Project.Status.CloudProviderAccessRoles[0].AtlasAssumedRoleExternalID,
					projectName,
				),
			).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.FalseCondition(status.NetworkPeerReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Network Peering due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.IntegrationReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Network Peering is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.NetworkPeers[0].VpcID = AwsVpcID
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

				g.Expect(testData.Project.Status.NetworkPeers).ToNot(HaveLen(0))
				g.Expect(testData.Project.Status.NetworkPeers[0].StatusName).To(Equal("AVAILABLE"))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.FalseCondition(status.IntegrationReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Integrations due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Integration is ready after configured properly", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pager-duty-service-key",
					Namespace: testData.Resources.Namespace,
					Labels: map[string]string{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					},
				},
				StringData: map[string]string{"password": os.Getenv("PAGER_DUTY_SERVICE_KEY")},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())

			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Integrations[0].Type = "PAGER_DUTY"
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.TrueCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 1).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Deleting project from the operator", func() {
			Expect(testData.K8SClient.Delete(ctx, testData.Project)).To(Succeed())
		})

		By("Stopping the operator", func() {
			managerStop()
		})

		By("Deleting Network Peering", func() {
			if networkPeerID != "" {
				_, err := atlasClient.Client.Peers.Delete(ctx, projectID, networkPeerID)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func(g Gomega) {
					_, _, err := atlasClient.Client.Peers.Get(ctx, projectID, networkPeerID)
					g.Expect(err).To(HaveOccurred())
				}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
			}
		})

		By("Deleting Project", func() {
			Eventually(func(g Gomega) {
				_, err := atlasClient.Client.Projects.Delete(ctx, projectID)
				g.Expect(err).ToNot(HaveOccurred())
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Deleting AWS Role", func() {
			if awsRoleARN != "" {
				Expect(cloudaccess.DeleteAWSIAMRoleByArn(awsRoleARN)).To(Succeed())
			}
		})

		By("Clean up", func() {
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})
})
