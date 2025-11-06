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
	"os"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

var _ = Describe("Atlas for Government", Label("atlas-gov"), func() {
	var awsHelper *cloud.AwsAction
	var testData *model.TestDataProvider
	var managerStop context.CancelFunc
	projectName := fmt.Sprintf("atlas-gov-e2e-%s", uuid.New().String()[0:6])
	clusterName := fmt.Sprintf("%s-cluster", projectName)
	ctx := context.Background()

	BeforeEach(func(ctx SpecContext) {
		By("Setting up cloud environment", func() {
			checkUpAWSEnvironment()

			aws, err := cloud.NewAWSAction(ctx, GinkgoT())
			Expect(err).ToNot(HaveOccurred())
			awsHelper = aws
		})

		By("Setting up test environment", func() {
			testData = model.DataProvider(ctx, "atlas-gov", model.NewEmptyAtlasKeyType().CreateAsGlobalLevelKey(), 30005, []func(*model.TestDataProvider){})

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

	It("Manage all supported Atlas for Government features", Label("focus-atlas-gov-supported"), func(ctx SpecContext) {
		By("Preparing API Key for integrations", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pager-duty-service-key",
					Namespace: testData.Resources.Namespace,
					Labels: map[string]string{
						secretservice.TypeLabelKey: secretservice.CredLabelVal,
					},
				},
				StringData: map[string]string{"password": os.Getenv("PAGER_DUTY_SERVICE_KEY")},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())
		})

		By("Creating a project to be managed by the operator", func() {
			akoProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projectName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasProjectSpec{
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
					Auditing: &akov2.Auditing{
						AuditAuthorizationSuccess: false,
						AuditFilter:               `{"$or":[{"users":[]},{"$and":[{"users":{"$elemMatch":{"$or":[{"db":"admin"}]}}},{"atype":{"$in":["authenticate","dropDatabase","createUser","dropUser","dropAllUsersFromDatabase","dropAllRolesFromDatabase","shutdown"]}}]}]}`,
						Enabled:                   true,
					},
					Settings: &akov2.ProjectSettings{
						IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
						IsDataExplorerEnabled:                       pointer.MakePtr(false),
						IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
						IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
						IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
						IsSchemaAdvisorEnabled:                      pointer.MakePtr(true),
					},
					CustomRoles: []akov2.CustomRole{
						{
							Name:           "testRole",
							InheritedRoles: nil,
							Actions: []akov2.Action{
								{
									Name: "INSERT",
									Resources: []akov2.Resource{
										{
											Database:   pointer.MakePtr("testD"),
											Collection: pointer.MakePtr("testCollection"),
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
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.ValidationSucceeded),
					api.TrueCondition(api.ProjectReadyType),
					api.TrueCondition(api.IPAccessListReadyType),
					api.TrueCondition(api.IntegrationReadyType),
					api.TrueCondition(api.MaintenanceWindowReadyType),
					api.TrueCondition(api.AuditingReadyType),
					api.TrueCondition(api.ProjectSettingsReadyType),
					api.TrueCondition(api.ProjectCustomRolesReadyType),
					api.TrueCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring a Team", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())

			users, _, err := atlasClient.Client.MongoDBCloudUsersApi.
				ListProjectUsers(ctx, testData.Project.ID()).
				Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(users.GetResults()).ToNot(BeEmpty())

			usernames := make([]akov2.TeamUser, 0, users.GetTotalCount())
			for _, user := range users.GetResults() {
				usernames = append(usernames, akov2.TeamUser(user.GetUsername()))
			}

			akoTeam := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-team", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.TeamSpec{
					Name:      fmt.Sprintf("%s-team", projectName),
					Usernames: usernames,
				},
			}
			testData.Teams = []*akov2.AtlasTeam{akoTeam}
			Expect(testData.K8SClient.Create(ctx, testData.Teams[0]))

			testData.Project.Spec.Teams = []akov2.Team{
				{
					TeamRef: common.ResourceRefNamespaced{
						Name:      fmt.Sprintf("%s-team", projectName),
						Namespace: testData.Resources.Namespace,
					},
					Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ProjectTeamsReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Cloud Provider Access", func() {
			assumedRoleArn, err := cloudaccess.CreateAWSIAMRole(ctx, projectName)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.CloudProviderIntegrations = []akov2.CloudProviderIntegration{
				{
					ProviderName:      "AWS",
					IamAssumedRoleArn: assumedRoleArn,
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.CloudProviderIntegrations).ShouldNot(BeEmpty())
				g.Expect(testData.Project.Status.CloudProviderIntegrations[0].Status).Should(BeElementOf([2]string{status.CloudProviderIntegrationStatusCreated, status.CloudProviderIntegrationStatusFailedToAuthorize}))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())

			Expect(
				cloudaccess.AddAtlasStatementToAWSIAMRole(ctx, testData.Project.Status.CloudProviderIntegrations[0].AtlasAWSAccountArn, testData.Project.Status.CloudProviderIntegrations[0].AtlasAssumedRoleExternalID, projectName),
			).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.CloudProviderIntegrationReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Networking Peering", func() {
			awsAccountID, err := awsHelper.GetAccountID(ctx)
			Expect(err).ToNot(HaveOccurred())

			AwsVpcID, err := awsHelper.InitNetwork(ctx, projectName, "10.0.0.0/24", "us-west-1", map[string]string{"subnet-1": "10.0.0.0/24"}, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.NetworkPeers = []akov2.NetworkPeer{
				{
					ProviderName:        "AWS",
					AccepterRegionName:  "us-west-1",
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

			Expect(awsHelper.AcceptVpcPeeringConnection(ctx, testData.Project.Status.NetworkPeers[0].ConnectionID, "us-west-1")).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.NetworkPeerReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Encryption at Rest", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			atlasAccountARN := testData.Project.Status.CloudProviderIntegrations[0].AtlasAWSAccountArn
			awsRoleARN := testData.Project.Status.CloudProviderIntegrations[0].IamAssumedRoleArn
			atlasRoleID := testData.Project.Status.CloudProviderIntegrations[0].RoleID

			customerMasterKeyID, err := awsHelper.CreateKMS(ctx, fmt.Sprintf("%s-kms", projectName), "us-west-1", atlasAccountARN, awsRoleARN)
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
						secretservice.TypeLabelKey: secretservice.CredLabelVal,
					},
				},
				Data: map[string][]byte{
					"CustomerMasterKeyID": []byte(customerMasterKeyID),
					"RoleID":              []byte(atlasRoleID),
				},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())

			encryptionAtRest := &akov2.EncryptionAtRest{
				AwsKms: akov2.AwsKms{
					Enabled: pointer.MakePtr(true),
					Region:  "US_WEST_1",
					SecretRef: common.ResourceRefNamespaced{
						Name:      "aws-secret",
						Namespace: testData.Resources.Namespace,
					},
				},
			}

			_, err = akoretry.RetryUpdateOnConflict(
				ctx,
				testData.K8SClient,
				client.ObjectKeyFromObject(testData.Project),
				func(project *akov2.AtlasProject) {
					project.Spec.EncryptionAtRest = encryptionAtRest
				})
			Expect(err).To(BeNil())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.EncryptionAtRestReadyType))))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Configuring Private Endpoint", func() {
			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.PrivateEndpoints = []akov2.PrivateEndpoint{
				{
					Provider: "AWS",
					Region:   "us-west-1",
				},
			}
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.PrivateEndpoints).ShouldNot(BeEmpty())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.PrivateEndpointServiceReadyType))))
			}).WithTimeout(time.Minute * 15).WithPolling(time.Second * 20).Should(Succeed())

			peID, err := awsHelper.CreatePrivateEndpoint(
				ctx,
				testData.Project.Status.PrivateEndpoints[0].ServiceName,
				fmt.Sprintf("pe-%s-gov", testData.Project.Status.PrivateEndpoints[0].ID),
				"us-west-1",
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
			testData.Project.Spec.PrivateEndpoints[0].ID = peID
			Expect(testData.K8SClient.Update(ctx, testData.Project)).To(Succeed())

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElement(conditions.MatchCondition(api.TrueCondition(api.PrivateEndpointReadyType))))
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Project is in ready state", func() {
			expectedConditions := conditions.MatchConditions(
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ProjectReadyType),
				api.TrueCondition(api.IPAccessListReadyType),
				api.TrueCondition(api.IntegrationReadyType),
				api.TrueCondition(api.MaintenanceWindowReadyType),
				api.TrueCondition(api.AuditingReadyType),
				api.TrueCondition(api.ProjectSettingsReadyType),
				api.TrueCondition(api.ProjectCustomRolesReadyType),
				api.TrueCondition(api.ProjectTeamsReadyType),
				api.TrueCondition(api.CloudProviderIntegrationReadyType),
				api.TrueCondition(api.NetworkPeerReadyType),
				api.TrueCondition(api.EncryptionAtRestReadyType),
				api.TrueCondition(api.ReadyType),
			)

			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Cluster", func() {
			akoBackupPolicy := &akov2.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-policy", clusterName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasBackupPolicySpec{
					Items: []akov2.AtlasBackupPolicyItem{
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

			akoBackupSchedule := &akov2.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-schedule", clusterName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasBackupScheduleSpec{
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

			akoDeployment := &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      clusterName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      projectName,
							Namespace: testData.Resources.Namespace,
						},
					},
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:          clusterName,
						BackupEnabled: pointer.MakePtr(true),
						BiConnector: &akov2.BiConnectorSpec{
							Enabled:        pointer.MakePtr(true),
							ReadPreference: "secondary",
						},
						ClusterType:              "REPLICASET",
						DiskSizeGB:               pointer.MakePtr(40),
						EncryptionAtRestProvider: "AWS",
						Labels: []common.LabelSpec{
							{Key: "type", Value: "e2e-test"},
							{Key: "context", Value: "cloud-gov"},
						},
						MongoDBMajorVersion: "7.0",
						Paused:              pointer.MakePtr(false),
						PitEnabled:          pointer.MakePtr(true),
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								NumShards: 1,
								ZoneName:  "GOV1",
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ElectableSpecs: &akov2.Specs{
											DiskIOPS:     pointer.MakePtr(int64(3000)),
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
										AutoScaling: &akov2.AdvancedAutoScalingSpec{
											DiskGB: &akov2.DiskGB{
												Enabled: pointer.MakePtr(true),
											},
											Compute: &akov2.ComputeSpec{
												Enabled:          pointer.MakePtr(true),
												ScaleDownEnabled: pointer.MakePtr(true),
												MinInstanceSize:  "M20",
												MaxInstanceSize:  "M40",
											},
										},
										Priority:     pointer.MakePtr(7),
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
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
					ProcessArgs: &akov2.ProcessArgs{
						DefaultReadConcern:        "available",
						MinimumEnabledTLSProtocol: "TLS1_2",
						JavascriptEnabled:         pointer.MakePtr(true),
						NoTableScan:               pointer.MakePtr(false),
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDeployment))
		})

		By("Cluster is in ready state", func() {
			expectedConditions := conditions.MatchConditions(
				api.TrueCondition(api.DeploymentReadyType),
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)

			Eventually(func(g Gomega) {
				akoDeployment := &akov2.AtlasDeployment{}
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
						secretservice.TypeLabelKey: secretservice.CredLabelVal,
					},
				},
				StringData: map[string]string{"password": "myHardPass2MyDB"},
			}
			Expect(testData.K8SClient.Create(ctx, secret)).To(Succeed())
			akoDBUser := &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-dbuser", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      projectName,
							Namespace: testData.Resources.Namespace,
						},
					},
					DatabaseName: "admin",
					Labels: []common.LabelSpec{
						{Key: "type", Value: "e2e-test"},
						{Key: "context", Value: "cloud-gov"},
					},
					Roles: []akov2.RoleSpec{
						{
							RoleName:     "readAnyDatabase",
							DatabaseName: "admin",
						},
					},
					Scopes: []akov2.ScopeSpec{
						{
							Name: clusterName,
							Type: akov2.DeploymentScopeType,
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
			expectedConditions := conditions.MatchConditions(
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ValidationSucceeded),
				api.TrueCondition(api.ResourceVersionStatus),
			)

			Eventually(func(g Gomega) {
				akoDBUser := &akov2.AtlasDatabaseUser{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-dbuser", projectName)}, akoDBUser)).To(Succeed())
				g.Expect(akoDBUser.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	It("Fail to manage when there are non supported features for Atlas for Government", Label("focus-atlas-gov-unsupported"), func() {
		By("Creating a project to be managed by the operator", func() {
			akoProject := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      projectName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasProjectSpec{
					Name:                    projectName,
					RegionUsageRestrictions: "GOV_REGIONS_ONLY",
				},
			}
			testData.Project = akoProject

			Expect(testData.K8SClient.Create(ctx, testData.Project))
		})

		By("Project is ready", func() {
			Eventually(func(g Gomega) {
				expectedConditions := conditions.MatchConditions(
					api.TrueCondition(api.ValidationSucceeded),
					api.TrueCondition(api.ProjectReadyType),
					api.TrueCondition(api.ReadyType),
				)

				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(testData.Project), testData.Project)).To(Succeed())
				g.Expect(testData.Project.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Flex Cluster", func() {
			akoDeployment := &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      clusterName,
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name:      projectName,
							Namespace: testData.Resources.Namespace,
						},
					},
					FlexSpec: &akov2.FlexSpec{
						Name: clusterName,
						ProviderSettings: &akov2.FlexProviderSettings{
							BackingProviderName: "AWS",
							RegionName:          "US_GOV_WEST_1",
						},
						TerminationProtectionEnabled: false,
					},
				},
			}
			Expect(testData.K8SClient.Create(ctx, akoDeployment))
		})

		By("Flex is not supported in Atlas for government", func() {
			expectedConditions := conditions.MatchConditions(
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.AtlasGovUnsupported)).
					WithMessageRegexp("the AtlasDeployment is not supported by Atlas for government"),
				api.FalseCondition(api.ReadyType),
				api.TrueCondition(api.ResourceVersionStatus),
			)

			Eventually(func(g Gomega) {
				akoDeployment := &akov2.AtlasDeployment{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: clusterName}, akoDeployment)).To(Succeed())
				g.Expect(akoDeployment.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})

		By("Creating a Data Federation", func() {
			akoDataFederation := &akov2.AtlasDataFederation{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-data-federation", projectName),
					Namespace: testData.Resources.Namespace,
				},
				Spec: akov2.DataFederationSpec{
					Project: common.ResourceRefNamespaced{
						Name:      projectName,
						Namespace: testData.Resources.Namespace,
					},
					Name: fmt.Sprintf("%s-data-federation", projectName),
					Storage: &akov2.Storage{
						Databases: []akov2.Database{
							{
								Name: "test-db-1",
								Collections: []akov2.Collection{
									{
										Name: "test-collection-1",
										DataSources: []akov2.DataSource{
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
						Stores: []akov2.Store{
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
			expectedConditions := conditions.MatchConditions(
				api.FalseCondition(api.DataFederationReadyType),
				api.FalseCondition(api.ReadyType),
			)

			Eventually(func(g Gomega) {
				akoDataFederation := &akov2.AtlasDataFederation{}
				g.Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-data-federation", projectName)}, akoDataFederation)).To(Succeed())
				g.Expect(akoDataFederation.Status.Conditions).To(ContainElements(expectedConditions))
			}).WithTimeout(time.Minute * 5).WithPolling(time.Second * 20).Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Deleting DataFederation from the operator", func() {
			akoDataFederation := &akov2.AtlasDataFederation{}
			err := testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: fmt.Sprintf("%s-data-federation", projectName)}, akoDataFederation)
			if err == nil {
				Expect(testData.K8SClient.Delete(ctx, akoDataFederation)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(testData.K8SClient.Get(ctx, client.ObjectKeyFromObject(akoDataFederation), akoDataFederation)).ToNot(Succeed())
				}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 20).Should(Succeed())
			}
		})

		By("Deleting DatabaseUser from the operator", func() {
			akoDBUser := &akov2.AtlasDatabaseUser{}
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
			akoDeployment := &akov2.AtlasDeployment{}
			Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Namespace: testData.Resources.Namespace, Name: clusterName}, akoDeployment)).To(Succeed())
			err := testData.K8SClient.Delete(ctx, akoDeployment)
			Expect(err == nil || !k8serrors.IsNotFound(err)).To(BeTrue())

			Eventually(func(g Gomega) {
				_, _, err := atlasClient.Client.ClustersApi.GetCluster(ctx, testData.Project.ID(), clusterName).Execute()
				g.Expect(err).To(HaveOccurred())
			}).WithTimeout(time.Minute * 30).WithPolling(time.Second * 20).Should(Succeed())

			if akoDeployment.Spec.BackupScheduleRef.Name != "" {
				akoBackupSchedule := &akov2.AtlasBackupSchedule{}
				Expect(testData.K8SClient.Get(ctx, client.ObjectKey{Name: fmt.Sprintf("%s-schedule", clusterName), Namespace: testData.Resources.Namespace}, akoBackupSchedule)).To(Succeed())
				Expect(testData.K8SClient.Delete(ctx, akoBackupSchedule)).To(Succeed())

				akoBackupPolicy := &akov2.AtlasBackupPolicy{}
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
				g.Expect(testData.Project.Status.Conditions).ToNot(ContainElement(conditions.MatchCondition(api.TrueCondition(api.ProjectTeamsReadyType))))
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
