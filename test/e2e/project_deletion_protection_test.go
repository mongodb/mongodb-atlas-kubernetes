package e2e_test

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/cloud"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"github.com/google/uuid"
	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
)

var _ = Describe("Project Deletion Protection", Label("project", "deletion-protection"), func() {
	var testData *model.TestDataProvider
	var managerStop context.CancelFunc
	var projectID, networkPeerID, atlasAccountARN, atlasRoleID, teamID string
	var awsRoleARN, awsAccountID, AwsVpcID, customerMasterKeyID string
	var usernames []string
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

			atlasRoleID = cloudProvider.RoleID
			atlasAccountARN = cloudProvider.AtlasAWSAccountARN
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

		By("Adding Maintenance Window to the project", func() {
			_, err := atlasClient.Client.MaintenanceWindows.Update(ctx, projectID, &mongodbatlas.MaintenanceWindow{
				DayOfWeek: 7,
				HourOfDay: toptr.MakePtr(20),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding Auditing to the project", func() {
			_, _, err := atlasClient.Client.Auditing.Configure(ctx, projectID, &mongodbatlas.Auditing{
				AuditFilter: `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`,
				Enabled:     toptr.MakePtr(true),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding Settings to the project", func() {
			_, _, err := atlasClient.Client.Projects.UpdateProjectSettings(ctx, projectID, &mongodbatlas.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
				IsDataExplorerEnabled:                       toptr.MakePtr(true),
				IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
				IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
				IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding AWS Encryption At Rest to the project", func() {
			awsAction, err := cloud.NewAWSAction(GinkgoT())
			Expect(err).ToNot(HaveOccurred())
			customerMasterKeyID, err = awsAction.CreateKMS(fmt.Sprintf("%s-kms", projectName), "eu-west-2", atlasAccountARN, awsRoleARN)
			Expect(err).ToNot(HaveOccurred())

			_, _, err = atlasClient.Client.EncryptionsAtRest.Create(ctx, &mongodbatlas.EncryptionAtRest{
				GroupID: projectID,
				AwsKms: mongodbatlas.AwsKms{
					Enabled:             toptr.MakePtr(true),
					CustomerMasterKeyID: customerMasterKeyID,
					Region:              "EU_WEST_2",
					RoleID:              atlasRoleID,
				},
			})
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding Custom Roles to the project", func() {
			_, _, err := atlasClient.Client.CustomDBRoles.Create(
				ctx,
				projectID,
				&mongodbatlas.CustomDBRole{
					RoleName:       "testRole",
					InheritedRoles: nil,
					Actions: []mongodbatlas.Action{
						{
							Action: "INSERT",
							Resources: []mongodbatlas.Resource{
								{
									DB:         toptr.MakePtr("testDB"),
									Collection: toptr.MakePtr("testCollection"),
								},
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Adding Assign team to the project", func() {
			users, _, err := atlasClient.Client.AtlasUsers.List(ctx, projectID, &mongodbatlas.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(users).ToNot(BeEmpty())

			usernames = make([]string, 0, len(users))
			for _, user := range users {
				usernames = append(usernames, user.Username)
			}

			team := &mongodbatlas.Team{
				Name:      fmt.Sprintf("%s-team", projectName),
				Usernames: usernames,
			}

			team, _, err = atlasClient.Client.Teams.Create(ctx, os.Getenv("MCLI_ORG_ID"), team)
			Expect(err).ToNot(HaveOccurred())
			teamID = team.ID

			_, _, err = atlasClient.Client.Projects.AddTeamsToProject(
				ctx,
				projectID,
				[]*mongodbatlas.ProjectTeam{
					{
						TeamID:    team.ID,
						RoleNames: []string{"GROUP_OWNER"},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
		})

		By("Creating a project and team to be managed by the operator", func() {
			akoTeam := &mdbv1.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-team", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.TeamSpec{
					Name:      fmt.Sprintf("%s-team", projectName),
					Usernames: []mdbv1.TeamUser{"user1@mongodb.com"},
				},
			}
			testData.Teams = []*mdbv1.AtlasTeam{akoTeam}
			Expect(testData.K8SClient.Create(ctx, testData.Teams[0]))

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
					MaintenanceWindow: project.MaintenanceWindow{
						DayOfWeek: 1,
						HourOfDay: 20,
					},
					Auditing: &mdbv1.Auditing{
						AuditAuthorizationSuccess: false,
						AuditFilter:               `{"$or":[{"users":[]},{"$and":[{"users":{"$elemMatch":{"$or":[{"db":"admin"}]}}},{"atype":{"$in":["authenticate","dropDatabase","createUser","dropUser","dropAllUsersFromDatabase","dropAllRolesFromDatabase","shutdown"]}}]}]}`,
						Enabled:                   true,
					},
					Settings: &mdbv1.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
						IsDataExplorerEnabled:                       toptr.MakePtr(false),
						IsExtendedStorageSizesEnabled:               toptr.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
						IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
					},
					EncryptionAtRest: &mdbv1.EncryptionAtRest{
						AwsKms: mdbv1.AwsKms{
							Enabled: toptr.MakePtr(true),
							Region:  "EU_WEST_1",
							SecretRef: common.ResourceRefNamespaced{
								Name:      "aws-secret",
								Namespace: testData.Resources.Namespace,
							},
						},
					},
					CustomRoles: []mdbv1.CustomRole{
						{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions: []mdbv1.Action{
								{
									Name: "INSERT",
									Resources: []mdbv1.Resource{
										{
											Database:   toptr.MakePtr("testD"),
											Collection: toptr.MakePtr("testCollection"),
										},
									},
								},
							},
						},
					},
					Teams: []mdbv1.Team{
						{
							TeamRef: common.ResourceRefNamespaced{
								Name:      fmt.Sprintf("%s-team", projectName),
								Namespace: testData.Resources.Namespace,
							},
							Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
						},
					},
				},
			}
			testData.Project = akoProject

			Expect(testData.K8SClient.Create(ctx, &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-secret",
					Namespace: testData.Project.Namespace,
					Labels: map[string]string{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					},
				},
				Data: map[string][]byte{
					"CustomerMasterKeyID": []byte(customerMasterKeyID),
					"RoleID":              []byte(atlasRoleID),
				},
			})).To(Succeed())

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
					status.FalseCondition(status.MaintenanceWindowReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
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
					status.FalseCondition(status.MaintenanceWindowReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
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
					status.FalseCondition(status.MaintenanceWindowReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
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
					status.FalseCondition(status.MaintenanceWindowReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
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
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.FalseCondition(status.MaintenanceWindowReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Maintenance Window due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Maintenance Window is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.MaintenanceWindow.DayOfWeek = 7
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.FalseCondition(status.AuditingReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Auditing is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Auditing.AuditFilter = `{"atype":"authenticate","param":{"user":"auditReadOnly","db":"admin","mechanism":"SCRAM-SHA-1"}}`
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.FalseCondition(status.ProjectSettingsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Project Settings due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Maintenance Window is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Settings.IsDataExplorerEnabled = toptr.MakePtr(true)
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.TrueCondition(status.ProjectSettingsReadyType),
					status.FalseCondition(status.EncryptionAtRestReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Encryption At Rest due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Encryption At Rest is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.EncryptionAtRest.AwsKms.Region = "EU_WEST_2"
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.TrueCondition(status.ProjectSettingsReadyType),
					status.TrueCondition(status.EncryptionAtRestReadyType),
					status.FalseCondition(status.ProjectCustomRolesReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Custom Roles due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Custom Roles is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.CustomRoles[0].Actions[0].Resources[0].Database = toptr.MakePtr("testDB")
			Expect(testData.K8SClient.Update(context.TODO(), testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.FalseCondition(status.ReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.CloudProviderAccessReadyType),
					status.TrueCondition(status.NetworkPeerReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.TrueCondition(status.ProjectSettingsReadyType),
					status.TrueCondition(status.EncryptionAtRestReadyType),
					status.TrueCondition(status.ProjectCustomRolesReadyType),
					status.FalseCondition(status.ProjectTeamsReadyType).
						WithReason(string(workflow.AtlasDeletionProtection)).
						WithMessageRegexp("unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information"),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Team is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Teams[0]), testData.Teams[0])).To(Succeed())
			testData.Teams[0].Spec.Usernames = make([]mdbv1.TeamUser, 0, len(usernames))
			for _, username := range usernames {
				testData.Teams[0].Spec.Usernames = append(testData.Teams[0].Spec.Usernames, mdbv1.TeamUser(username))
			}
			Expect(testData.K8SClient.Update(context.TODO(), testData.Teams[0])).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Teams[0]), testData.Teams[0])).To(Succeed())
				g.Expect(testData.Teams[0].Status.Conditions).To(ContainElements(testutil.MatchCondition(status.TrueCondition(status.ReadyType))))
			}).WithTimeout(time.Minute * 1).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Assigned Teams is ready after configured properly", func() {
			Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Teams[0].Roles[0] = "GROUP_OWNER"
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
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.TrueCondition(status.ProjectSettingsReadyType),
					status.TrueCondition(status.EncryptionAtRestReadyType),
					status.TrueCondition(status.ProjectCustomRolesReadyType),
					status.TrueCondition(status.ProjectTeamsReadyType),
				)

				g.Expect(testData.K8SClient.Get(context.TODO(), client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Deleting project from the operator", func() {
			Expect(testData.K8SClient.Delete(ctx, testData.Teams[0])).To(Succeed())
			Expect(testData.K8SClient.Delete(ctx, testData.Project)).To(Succeed())
			time.Sleep(time.Second * 30)
		})

		By("Stopping the operator", func() {
			managerStop()
		})

		By("Deleting Team", func() {
			if teamID != "" {
				_, err := atlasClient.Client.Teams.RemoveTeamFromProject(ctx, projectID, teamID)
				Expect(err).ToNot(HaveOccurred())

				_, err = atlasClient.Client.Teams.RemoveTeamFromOrganization(ctx, os.Getenv("MCLI_ORG_ID"), teamID)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func(g Gomega) {
					_, _, err := atlasClient.Client.Teams.Get(ctx, os.Getenv("MCLI_ORG_ID"), teamID)
					g.Expect(err).To(HaveOccurred())
				}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
			}
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
