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
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compare"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Alert configuration tests", Label("alert-config", "alert-configs-table"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(ctx context.Context) *model.TestDataProvider, alertConfigurations []akov2.AlertConfiguration) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			alertConfigFlow(testData, alertConfigurations)
		},
		Entry("Test[alert-configs-1]: Project with 2 identical alert configs", Label("focus-alert-configs-1"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "alert-configs-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.AlertConfiguration{
				{
					EventTypeName:    "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:          true,
					SeverityOverride: "CRITICAL",
					Threshold: &akov2.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Notifications: []akov2.Notification{
						{
							IntervalMin:  5,
							DelayMin:     pointer.MakePtr(5),
							EmailEnabled: pointer.MakePtr(true),
							SMSEnabled:   pointer.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &akov2.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "2", // make it a different alert config
						Units:     "HOURS",
					},
					Notifications: []akov2.Notification{
						{
							IntervalMin:  5,
							DelayMin:     pointer.MakePtr(5),
							EmailEnabled: pointer.MakePtr(true),
							SMSEnabled:   pointer.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
			},
		),
		Entry("Test[alert-configs-2]: Project with 2 different alert configs", Label("focus-alert-configs-2"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "alert-configs-2", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.AlertConfiguration{
				{
					EventTypeName: "JOINED_GROUP",
					Enabled:       true,
					Notifications: []akov2.Notification{
						{
							IntervalMin:  60,
							DelayMin:     pointer.MakePtr(0),
							EmailEnabled: pointer.MakePtr(true),
							SMSEnabled:   pointer.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &akov2.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Notifications: []akov2.Notification{
						{
							IntervalMin:  5,
							DelayMin:     pointer.MakePtr(5),
							EmailEnabled: pointer.MakePtr(true),
							SMSEnabled:   pointer.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
			},
		),
		Entry("Test[alert-configs-3]: Project with an alert config containing a matcher", Label("focus-alert-configs-3"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "alert-configs-3", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			[]akov2.AlertConfiguration{
				{
					EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
					Enabled:       true,
					Threshold: &akov2.Threshold{
						Operator:  "LESS_THAN",
						Threshold: "1",
						Units:     "HOURS",
					},
					Matchers: []akov2.Matcher{
						{
							FieldName: "CLUSTER_NAME",
							Operator:  "STARTS_WITH",
							Value:     "ako_e2e_test_",
						},
					},
					Notifications: []akov2.Notification{
						{
							IntervalMin:  5,
							DelayMin:     pointer.MakePtr(5),
							EmailEnabled: pointer.MakePtr(true),
							SMSEnabled:   pointer.MakePtr(false),
							Roles: []string{
								"GROUP_OWNER",
							},
							TypeName: "GROUP",
						},
					},
				},
			},
		),
	)
})

func alertConfigFlow(userData *model.TestDataProvider, alertConfigs []akov2.AlertConfiguration) {
	By("Enable Alert Config Sync on Atlas Project", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.Project.Name,
			Namespace: userData.Project.Namespace,
		}, userData.Project)).Should(Succeed())
		userData.Project.Spec.AlertConfigurationSyncEnabled = true
		userData.Project.Spec.AlertConfigurations = append(userData.Project.Spec.AlertConfigurations, alertConfigs...)
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Wait for Alert Configurations to activate", func() {
		actions.WaitForConditionsToBecomeTrue(userData, api.AlertConfigurationReadyType, api.ReadyType)
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name, Namespace: userData.Project.Namespace}, userData.Project)).Should(Succeed())
		Expect(userData.Project.Status.AlertConfigurations).Should(HaveLen(len(alertConfigs)))
	})

	By("Check alert configurations have no errors and match configured configs", func() {
		var err error
		alertConfigurations, _, err := atlasClient.Client.AlertConfigurationsApi.
			ListAlertConfigs(userData.Context, userData.Project.ID()).
			Execute()
		By("No errors listing alert configs", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Check config counts match configured count", func() {
			Expect(alertConfigurations.GetTotalCount()).Should(Equal(len(alertConfigs)), "Atlas alert configurations", alertConfigurations)
		})

		By("ID sets in Atlas matches the status IDs", func() {
			atlasIDList := make([]string, 0, alertConfigurations.GetTotalCount())
			for _, alertConfig := range alertConfigurations.GetResults() {
				atlasIDList = append(atlasIDList, alertConfig.GetId())
			}
			statusIDList := make([]string, 0, len(userData.Project.Status.AlertConfigurations))
			for _, alertConfig := range userData.Project.Status.AlertConfigurations {
				statusIDList = append(statusIDList, alertConfig.ID)
			}
			Expect(compare.IsEqualWithoutOrder(statusIDList, atlasIDList)).Should(BeTrue())
		})

		By("Each Atlas alert config matches its Kubernetes config", func() {
			atlasConvertedSpecs := []*admin.GroupAlertsConfig{}
			for i := range alertConfigs {
				akoConfig, err := alertConfigs[i].ToAtlas()
				Expect(err).ToNot(HaveOccurred())
				atlasConvertedSpecs = append(atlasConvertedSpecs, akoConfig)
			}
			atlasConfigs := alertConfigurations.GetResults()
			for _, atlasConfig := range atlasConfigs {
				normalizedAtlasConfig := normalizeAtlasAlertConfig(atlasConfig)
				Expect(atlasConvertedSpecs).To(ContainElement(&normalizedAtlasConfig))
			}
		})
	})
}

func normalizeAtlasAlertConfig(atlasConfig admin.GroupAlertsConfig) admin.GroupAlertsConfig {
	atlasConfig.Id = nil
	atlasConfig.GroupId = nil
	atlasConfig.Created = nil
	atlasConfig.Updated = nil
	atlasConfig.Links = nil

	notifications := atlasConfig.GetNotifications()
	for j := range notifications {
		notifications[j].NotifierId = nil
		notifications[j].DatadogApiKey = pointer.MakePtr("")
		notifications[j].OpsGenieApiKey = pointer.MakePtr("")
		notifications[j].ServiceKey = pointer.MakePtr("")
		notifications[j].ApiToken = pointer.MakePtr("")
		notifications[j].VictorOpsApiKey = pointer.MakePtr("")
		notifications[j].VictorOpsRoutingKey = pointer.MakePtr("")
	}
	atlasConfig.SetNotifications(notifications)

	return atlasConfig
}

var _ = Describe("Alert configuration with secrets test", Label("alert-config", "alert-config-datadog"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		Expect(os.Getenv("DATADOG_KEY")).ShouldNot(BeEmpty(), "Please setup DATADOG_KEY environment variable")
	})

	alertConfigs := []akov2.AlertConfiguration{
		{
			EventTypeName: "REPLICATION_OPLOG_WINDOW_RUNNING_OUT",
			Enabled:       true,
			Threshold: &akov2.Threshold{
				Operator:  "LESS_THAN",
				Threshold: "1",
				Units:     "HOURS",
			},
			Notifications: []akov2.Notification{
				{
					IntervalMin:  5,
					DelayMin:     pointer.MakePtr(5),
					EmailEnabled: pointer.MakePtr(true),
					SMSEnabled:   pointer.MakePtr(false),
					Roles: []string{
						"GROUP_OWNER",
					},
					TypeName: "DATADOG",
				},
			},
		},
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "datadog-creds",
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"DatadogAPIKey": []byte(os.Getenv("DATADOG_KEY")),
		},
	}

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
	})

	It("Should be able to create AtlasProject with Alert Config and secrets", func(ctx SpecContext) {
		testData = model.DataProvider(ctx, "alert-configs-1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())

		By("Creating an AtlasProject", func() {
			actions.ProjectCreationFlow(testData)
		})

		By("Creating Datadog credentials secret", func() {
			secret.Namespace = testData.Project.Namespace
			Expect(testData.K8SClient.Create(testData.Context, secret)).To(Succeed())
		})

		By("Configuring the Datadog alert using secret ref", func() {
			alertConfigs[0].Notifications[0].DatadogAPIKeyRef = common.ResourceRefNamespaced{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			}

			Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).Should(Succeed())
			testData.Project.Spec.AlertConfigurationSyncEnabled = true
			testData.Project.Spec.AlertConfigurations = alertConfigs
			Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
		})

		By("Verifying the Datadog config in Atlas", func() {
			atlasClient := atlas.GetClientOrFail()
			Eventually(func(g Gomega) {
				atlasAlertConfigs, _, err := atlasClient.Client.AlertConfigurationsApi.
					ListAlertConfigs(testData.Context, testData.Project.ID()).
					Execute()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(atlasAlertConfigs.GetTotalCount()).Should(Equal(len(alertConfigs)))
				g.Expect(atlasAlertConfigs.GetResults()[0].GetNotifications()[0].GetDatadogApiKey()).ShouldNot(BeEmpty())
			}).WithPolling(10 * time.Second).WithTimeout(5 * time.Minute).Should(Succeed())
		})
	})
})
