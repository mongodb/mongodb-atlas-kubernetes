package e2e_test

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
)

var _ = Describe("Atlas for Government", Label("atlas-gov"), func() {
	var awsHelper *cloud.AwsAction
	var testData *model.TestDataProvider
	var managerStop context.CancelFunc
	projectName := fmt.Sprintf("atlas-gov-e2e-%s", uuid.New().String()[0:6])
	clusterName := fmt.Sprintf("%s-cluster", projectName)
	ctx := context.Background()

	BeforeEach(func() {
		By("Setting up cloud environment", func() {
			checkUpAWSEnvironment()

			aws, err := cloud.NewAWSAction(GinkgoT())
			Expect(err).ToNot(HaveOccurred())
			awsHelper = aws
		})

		By("Setting up test environment", func() {
			testData = model.DataProvider(
				"atlas-gov",
				model.NewEmptyAtlasKeyType().CreateAsGlobalLevelKey(),
				30005,
				[]func(*model.TestDataProvider){},
			)

			actions.CreateNamespaceAndSecrets(testData)
		})

		By("Setting up the operator", func() {
			managerStart, err := k8s.RunManager(
				k8s.WithAtlasDomain(os.Getenv("MCLI_OPS_MANAGER_URL")),
				k8s.WithGlobalKey(client.ObjectKey{Namespace: testData.Resources.Namespace, Name: config.DefaultOperatorGlobalKey}),
				k8s.WithNamespaces(testData.Resources.Namespace),
				k8s.WithObjectDeletionProtection(false),
				k8s.WithSubObjectDeletionProtection(false),
			)
			Expect(err).ToNot(HaveOccurred())

			cancelCtx, cancel := context.WithCancel(ctx)
			managerStop = cancel
			go func() {
				err := managerStart(cancelCtx)
				if err != nil {
					GinkgoWriter.Write([]byte(err.Error()))
				}
				Expect(err).ToNot(HaveOccurred())
			}()
		})
	})

	It("Manage all supported Atlas for Government features", func() {
		By("Preparing API Key for integrations", func() {
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
		})

		By("Creating a project to be managed by the operator", func() {
			akoProject := &mdbv1.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projectName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasProjectSpec{
					Name:                    projectName,
					RegionUsageRestrictions: "NONE",
					ProjectIPAccessList: []project.IPAccessList{
						{
							CIDRBlock: "10.0.0.0/24",
						},
					},
					Integrations: []project.Integration{
						{
							Type:   "PAGER_DUTY",
							Region: "US",
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
				},
			}
			testData.Project = akoProject

			Expect(testData.K8SClient.Create(ctx, testData.Project))
		})

		By("Project is ready", func() {
			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.TrueCondition(status.IPAccessListReadyType),
					status.TrueCondition(status.IntegrationReadyType),
					status.TrueCondition(status.MaintenanceWindowReadyType),
					status.TrueCondition(status.AuditingReadyType),
					status.TrueCondition(status.ProjectSettingsReadyType),
					status.TrueCondition(status.ProjectCustomRolesReadyType),
					status.TrueCondition(status.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring a Team", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

			users, _, err := atlasClient.Client.AtlasUsers.List(ctx, testData.Project.ID(), &mongodbatlas.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(users).ToNot(BeEmpty())

			usernames := make([]mdbv1.TeamUser, 0, len(users))
			for _, user := range users {
				usernames = append(usernames, mdbv1.TeamUser(user.Username))
			}

			akoTeam := &mdbv1.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-team", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.TeamSpec{
					Name:      fmt.Sprintf("%s-team", projectName),
					Usernames: usernames,
				},
			}
			testData.Teams = []*mdbv1.AtlasTeam{akoTeam}
			Expect(testData.K8SClient.Create(ctx, testData.Teams[0]))

			testData.Project.Spec.Teams = []mdbv1.Team{
				{
					TeamRef: common.ResourceRefNamespaced{
						Name:      fmt.Sprintf("%s-team", projectName),
						Namespace: testData.Resources.Namespace,
					},
					Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectTeamsReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Cloud Provider Access", func() {
			assumedRoleArn, err := cloudaccess.CreateAWSIAMRole(projectName)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.CloudProviderAccessRoles = []mdbv1.CloudProviderAccessRole{
				{
					ProviderName:      "AWS",
					IamAssumedRoleArn: assumedRoleArn,
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.CloudProviderAccessRoles).ShouldNot(BeEmpty())
				g.Expect(testData.Project.Status.CloudProviderAccessRoles[0].Status).Should(BeElementOf([2]string{status.CloudProviderAccessStatusCreated, status.CloudProviderAccessStatusFailedToAuthorize}))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			Expect(
				cloudaccess.AddAtlasStatementToAWSIAMRole(
					testData.Project.Status.CloudProviderAccessRoles[0].AtlasAWSAccountArn,
					testData.Project.Status.CloudProviderAccessRoles[0].AtlasAssumedRoleExternalID,
					projectName,
				),
			).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.CloudProviderAccessReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Networking Peering", func() {
			awsAccountID, err := awsHelper.GetAccountID()
			Expect(err).ToNot(HaveOccurred())

			AwsVpcID, err := awsHelper.InitNetwork(projectName, "10.0.0.0/24", "us-east-1", map[string]string{"subnet-1": "10.0.0.0/24"}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.NetworkPeers = []mdbv1.NetworkPeer{
				{
					ProviderName:        "AWS",
					AccepterRegionName:  "us-east-1",
					AtlasCIDRBlock:      "192.168.224.0/21",
					AWSAccountID:        awsAccountID,
					RouteTableCIDRBlock: "10.0.0.0/24",
					VpcID:               AwsVpcID,
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.NetworkPeers).ShouldNot(BeEmpty())
				g.Expect(testData.Project.Status.NetworkPeers[0].StatusName).Should(Equal("PENDING_ACCEPTANCE"))
			}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 20).Should(Succeed())

			Expect(awsHelper.AcceptVpcPeeringConnection(testData.Project.Status.NetworkPeers[0].ConnectionID, "us-east-1")).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.NetworkPeerReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Encryption at Rest", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			atlasAccountARN := testData.Project.Status.CloudProviderAccessRoles[0].AtlasAWSAccountArn
			awsRoleARN := testData.Project.Status.CloudProviderAccessRoles[0].IamAssumedRoleArn
			atlasRoleID := testData.Project.Status.CloudProviderAccessRoles[0].RoleID

			customerMasterKeyID, err := awsHelper.CreateKMS(fmt.Sprintf("%s-kms", projectName), "us-east-1", atlasAccountARN, awsRoleARN)
			Expect(err).ToNot(HaveOccurred())

			secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-secret",
					Namespace: testData.Resources.Namespace,
					Labels: map[string]string{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					},
				},
				Data: map[string][]byte{
					"CustomerMasterKeyID": []byte(customerMasterKeyID),
					"RoleID":              []byte(atlasRoleID),
				},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())

			testData.Project.Spec.EncryptionAtRest = &mdbv1.EncryptionAtRest{
				AwsKms: mdbv1.AwsKms{
					Enabled: toptr.MakePtr(true),
					Region:  "US_EAST_1",
					SecretRef: common.ResourceRefNamespaced{
						Name:      "aws-secret",
						Namespace: testData.Resources.Namespace,
					},
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.EncryptionAtRestReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Private Endpoint", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.PrivateEndpoints = []mdbv1.PrivateEndpoint{
				{
					Provider: "AWS",
					Region:   "us-east-1",
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.PrivateEndpoints).ShouldNot(BeEmpty())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.PrivateEndpointServiceReadyType))))
			}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 20).Should(Succeed())

			peID, err := awsHelper.CreatePrivateEndpoint(
				testData.Project.Status.PrivateEndpoints[0].ServiceName,
				fmt.Sprintf("pe-%s-gov", testData.Project.Status.PrivateEndpoints[0].ID),
				"us-east-1",
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.PrivateEndpoints[0].ID = peID
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(testutil.MatchCondition(status.TrueCondition(status.PrivateEndpointReadyType))))
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Project is in ready state", func() {
			expectedConditions := testutil.MatchConditions(
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ProjectReadyType),
				status.TrueCondition(status.IPAccessListReadyType),
				status.TrueCondition(status.IntegrationReadyType),
				status.TrueCondition(status.MaintenanceWindowReadyType),
				status.TrueCondition(status.AuditingReadyType),
				status.TrueCondition(status.ProjectSettingsReadyType),
				status.TrueCondition(status.ProjectCustomRolesReadyType),
				status.TrueCondition(status.ProjectTeamsReadyType),
				status.TrueCondition(status.CloudProviderAccessReadyType),
				status.TrueCondition(status.NetworkPeerReadyType),
				status.TrueCondition(status.EncryptionAtRestReadyType),
				status.TrueCondition(status.ReadyType),
			)

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Cluster", func() {
			akoBackupPolicy := &mdbv1.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-policy", clusterName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasBackupPolicySpec{
					Items: []mdbv1.AtlasBackupPolicyItem{
						{
							FrequencyType:     "hourly",
							FrequencyInterval: 12,
							RetentionUnit:     "days",
							RetentionValue:    7,
						},
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoBackupPolicy))

			akoBackupSchedule := &mdbv1.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-schedule", clusterName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasBackupScheduleSpec{
					PolicyRef: common.ResourceRefNamespaced{
						Name:      fmt.Sprintf("%s-policy", clusterName),
						Namespace: testData.Resources.Namespace,
					},
					AutoExportEnabled:                 false,
					Export:                            nil,
					ReferenceHourOfDay:                22,
					ReferenceMinuteOfHour:             30,
					RestoreWindowDays:                 7,
					UpdateSnapshots:                   true,
					UseOrgAndGroupNamesInExportPrefix: true,
					CopySettings:                      nil,
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoBackupSchedule))

			akoDeployment := &mdbv1.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      clusterName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasDeploymentSpec{
					Project: common.ResourceRefNamespaced{
						Name:      projectName,
						Namespace: testData.Resources.Namespace,
					},
					DeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
						Name:          clusterName,
						BackupEnabled: toptr.MakePtr(true),
						BiConnector: &mdbv1.BiConnectorSpec{
							Enabled:        toptr.MakePtr(true),
							ReadPreference: "secondary",
						},
						ClusterType:              "REPLICASET",
						DiskSizeGB:               toptr.MakePtr(40),
						EncryptionAtRestProvider: "AWS",
						Labels: []common.LabelSpec{
							{Key: "type", Value: "e2e-test"},
							{Key: "context", Value: "cloud-gov"},
						},
						MongoDBMajorVersion: "7.0",
						Paused:              toptr.MakePtr(false),
						PitEnabled:          toptr.MakePtr(true),
						ReplicationSpecs: []*mdbv1.AdvancedReplicationSpec{
							{
								NumShards: 1,
								ZoneName:  "GOV1",
								RegionConfigs: []*mdbv1.AdvancedRegionConfig{
									{
										ElectableSpecs: &mdbv1.Specs{
											DiskIOPS:     toptr.MakePtr(int64(3000)),
											InstanceSize: "M20",
											NodeCount:    toptr.MakePtr(3),
										},
										AutoScaling: &mdbv1.AdvancedAutoScalingSpec{
											DiskGB: &mdbv1.DiskGB{
												Enabled: toptr.MakePtr(true),
											},
											Compute: &mdbv1.ComputeSpec{
												Enabled:          toptr.MakePtr(true),
												ScaleDownEnabled: toptr.MakePtr(true),
												MinInstanceSize:  "M20",
												MaxInstanceSize:  "M40",
											},
										},
										Priority:     toptr.MakePtr(7),
										ProviderName: "AWS",
										RegionName:   "US_EAST_1",
									},
								},
							},
						},
						RootCertType:         "ISRGROOTX1",
						VersionReleaseSystem: "LTS",
					},
					BackupScheduleRef: common.ResourceRefNamespaced{
						Name:      fmt.Sprintf("%s-schedule", clusterName),
						Namespace: testData.Resources.Namespace,
					},
					ProcessArgs: &mdbv1.ProcessArgs{
						DefaultReadConcern:        "available",
						MinimumEnabledTLSProtocol: "TLS1_2",
						JavascriptEnabled:         toptr.MakePtr(true),
						NoTableScan:               toptr.MakePtr(false),
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDeployment))
		})

		By("Cluster is in ready state", func() {
			expectedConditions := testutil.MatchConditions(
				status.TrueCondition(status.DeploymentReadyType),
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)

			Eventually(func(g Gomega) {
				akoDeployment := &mdbv1.AtlasDeployment{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: clusterName}, akoDeployment)).To(Succeed())
				g.Expect(akoDeployment.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 45).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a DatabaseUser", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-dbuser-pass", projectName),
					Namespace: testData.Resources.Namespace,
					Labels: map[string]string{
						connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
					},
				},
				StringData: map[string]string{"password": "myHardPass2MyDB"},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())
			akoDBUser := &mdbv1.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-dbuser", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasDatabaseUserSpec{
					Project: common.ResourceRefNamespaced{
						Name:      projectName,
						Namespace: testData.Resources.Namespace,
					},
					DatabaseName: "admin",
					Labels: []common.LabelSpec{
						{Key: "type", Value: "e2e-test"},
						{Key: "context", Value: "cloud-gov"},
					},
					Roles: []mdbv1.RoleSpec{
						{
							RoleName:     "readAnyDatabase",
							DatabaseName: "admin",
						},
					},
					Scopes: []mdbv1.ScopeSpec{
						{
							Name: clusterName,
							Type: mdbv1.DeploymentScopeType,
						},
					},
					Username: fmt.Sprintf("%s-dbuser", projectName),
					PasswordSecret: &common.ResourceRef{
						Name: fmt.Sprintf("%s-dbuser-pass", projectName),
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDBUser))
		})

		By("DatabaseUser is in ready state", func() {
			expectedConditions := testutil.MatchConditions(
				status.TrueCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
				status.TrueCondition(status.ResourceVersionStatus),
			)

			Eventually(func(g Gomega) {
				akoDBUser := &mdbv1.AtlasDatabaseUser{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-dbuser", projectName)}, akoDBUser)).To(Succeed())
				g.Expect(akoDBUser.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	It("Fail to manage when there are non supported features for Atlas for Government", func() {
		By("Creating a project to be managed by the operator", func() {
			akoProject := &mdbv1.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projectName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasProjectSpec{
					Name:                    projectName,
					RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				},
			}
			testData.Project = akoProject

			Expect(testData.K8SClient.Create(ctx, testData.Project))
		})

		By("Project is ready", func() {
			Eventually(func(g Gomega) {
				expectedConditions := testutil.MatchConditions(
					status.TrueCondition(status.ValidationSucceeded),
					status.TrueCondition(status.ProjectReadyType),
					status.TrueCondition(status.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Serverless Cluster", func() {
			akoDeployment := &mdbv1.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      clusterName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.AtlasDeploymentSpec{
					Project: common.ResourceRefNamespaced{
						Name:      projectName,
						Namespace: testData.Resources.Namespace,
					},
					ServerlessSpec: &mdbv1.ServerlessSpec{
						Name: clusterName,
						ProviderSettings: &mdbv1.ProviderSettingsSpec{
							BackingProviderName: "AWS",
							ProviderName:        "SERVERLESS",
							RegionName:          "US_GOV_EAST_1",
						},
						TerminationProtectionEnabled: false,
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDeployment))
		})

		By("Serverless is not supported in Atlas for government", func() {
			expectedConditions := testutil.MatchConditions(
				status.FalseCondition(status.DeploymentReadyType),
				status.FalseCondition(status.ReadyType),
				status.TrueCondition(status.ValidationSucceeded),
			)

			Eventually(func(g Gomega) {
				akoDeployment := &mdbv1.AtlasDeployment{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: clusterName}, akoDeployment)).To(Succeed())
				g.Expect(akoDeployment.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Data Federation", func() {
			akoDataFederation := &mdbv1.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-data-federation", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: mdbv1.DataFederationSpec{
					Project: common.ResourceRefNamespaced{
						Name:      projectName,
						Namespace: testData.Resources.Namespace,
					},
					Name: fmt.Sprintf("%s-data-federation", projectName),
					Storage: &mdbv1.Storage{
						Databases: []mdbv1.Database{
							{
								Name: "test-db-1",
								Collections: []mdbv1.Collection{
									{
										Name: "test-collection-1",
										DataSources: []mdbv1.DataSource{
											{
												StoreName: "http-test",
												Urls: []string{
													"https://data.cityofnewyork.us/api/views/vfnx-vebw/rows.csv",
												},
											},
										},
									},
								},
							},
						},
						Stores: []mdbv1.Store{
							{
								Name:     "http-test",
								Provider: "http",
							},
						},
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDataFederation))
		})

		By("DataFederation is not supported in Atlas for government", func() {
			expectedConditions := testutil.MatchConditions(
				status.FalseCondition(status.DataFederationReadyType),
				status.FalseCondition(status.ReadyType),
			)

			Eventually(func(g Gomega) {
				akoDataFederation := &mdbv1.AtlasDataFederation{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-data-federation", projectName)}, akoDataFederation)).To(Succeed())
				g.Expect(akoDataFederation.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Deleting DataFederation from the operator", func() {
			akoDataFederation := &mdbv1.AtlasDataFederation{}
			err := testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-data-federation", projectName)}, akoDataFederation)
			if err == nil {
				Expect(testData.K8SClient.Delete(ctx, akoDataFederation)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(akoDataFederation), akoDataFederation)).ToNot(Succeed())
				}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
			}
		})

		By("Deleting DatabaseUser from the operator", func() {
			akoDBUser := &mdbv1.AtlasDatabaseUser{}
			err := testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-dbuser", projectName)}, akoDBUser)
			if err == nil {
				Expect(testData.K8SClient.Delete(ctx, akoDBUser)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(akoDBUser), akoDBUser)).ToNot(Succeed())
				}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
			}
		})

		By("Deleting cluster from the operator", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			akoDeployment := &mdbv1.AtlasDeployment{}
			Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: clusterName}, akoDeployment)).To(Succeed())
			Expect(testData.K8SClient.Delete(ctx, akoDeployment)).To(Succeed())

			Eventually(func(g Gomega) {
				_, _, err := atlasClient.Client.AdvancedClusters.Get(ctx, testData.Project.ID(), clusterName)
				g.Expect(err).To(HaveOccurred())
			}).WithTimeout(time.Minute * 30).WithPolling(time.Second * 20).Should(Succeed())

			if akoDeployment.Spec.BackupScheduleRef.Name != "" {
				akoBackupSchedule := &mdbv1.AtlasBackupSchedule{}
				Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-schedule", clusterName), Namespace: testData.Resources.Namespace}, akoBackupSchedule)).To(Succeed())
				Expect(testData.K8SClient.Delete(ctx, akoBackupSchedule)).To(Succeed())

				akoBackupPolicy := &mdbv1.AtlasBackupPolicy{}
				Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-policy", clusterName), Namespace: testData.Resources.Namespace}, akoBackupPolicy)).To(Succeed())
				Expect(testData.K8SClient.Delete(ctx, akoBackupPolicy)).To(Succeed())
			}
		})

		By("Deleting team from the operator", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.Teams = nil
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).ToNot(ContainElement(testutil.MatchCondition(status.TrueCondition(status.ProjectTeamsReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Deleting project from the operator", func() {
			Expect(testData.K8SClient.Delete(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).ToNot(Succeed())
			}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Stopping the operator", func() {
			managerStop()
		})

		By("Clean up", func() {
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})
})
