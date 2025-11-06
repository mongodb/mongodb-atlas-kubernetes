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
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	akoretry "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/retry"
)

const (
	datadogEnvKey         = "DATADOG_KEY"
	pagerDutyEnvKey       = "PAGER_DUTY_SERVICE_KEY"
	integrationSecretName = "integration-secret"
)

var _ = Describe("Project Third-Party Integration", Label("integration-ns"), func() {
	var testData *model.TestDataProvider

	BeforeEach(func() {
		Expect(os.Getenv(datadogEnvKey)).ShouldNot(BeEmpty())
		Expect(os.Getenv(pagerDutyEnvKey)).ShouldNot(BeEmpty())
	})

	AfterEach(func() {
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

	DescribeTable("Integration can be configured in a project",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, integration project.Integration, envKeyName string, setSecret configSecret) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			integrationTest(testData, integration, os.Getenv(envKeyName), setSecret)
		},

		Entry("Users can integrate DATADOG on region US1", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "datatog-us1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "DATADOG",
				Region: "US",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region US3", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "datatog-us3", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "DATADOG",
				Region: "US3",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region US5", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "datatog-us5", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "DATADOG",
				Region: "US5",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region EU1", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "datatog-eu1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "DATADOG",
				Region: "EU",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate DATADOG on region AP1", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "datatog-ap1", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "DATADOG",
				Region: "AP1",
			},
			datadogEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.APIKeyRef = ref
			},
		),
		Entry("Users can integrate PagerDuty on region US", Label("focus-project-integration"),
			func(ctx SpecContext) *model.TestDataProvider {
				return model.DataProvider(ctx, "pager-duty-us", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			project.Integration{
				Type:   "PAGER_DUTY",
				Region: "US",
			},
			pagerDutyEnvKey,
			func(integration *project.Integration, ref common.ResourceRefNamespaced) {
				integration.ServiceKeyRef = ref
			},
		),
	)

	It("Project Integrations are not greedy", Label("focus-project-integration-not-greedy"), func(ctx SpecContext) {
		testData = model.DataProvider(ctx, "several-integrations", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30018, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject()).WithObjectDeletionProtection(false)
		actions.ProjectCreationFlow(testData)

		By("Create Secrets for integrations", func() {
			for _, secretName := range []string{"datadog-secret", "slack-secret", "webhook-secret"} {
				fakeSecret := os.Getenv(datadogEnvKey) // good for datadog and slack
				if secretName == "webhook-secret" {
					fakeSecret = utils.RandomName("fake-secret")
				}
				Expect(
					k8s.CreateUserSecret(ctx, testData.K8SClient, fakeSecret, secretName, testData.Resources.Namespace),
				).Should(Succeed())
			}
		})

		integrations := []project.Integration{
			{
				Type:   "DATADOG",
				Region: "US",
				APIKeyRef: common.ResourceRefNamespaced{
					Name:      "datadog-secret",
					Namespace: testData.Resources.Namespace,
				},
			},
			{
				Type:        "SLACK",
				ChannelName: "channel",
				TeamName:    "team",
				APITokenRef: common.ResourceRefNamespaced{
					Name:      "slack-secret",
					Namespace: testData.Resources.Namespace,
				},
			},
			{
				Type: "WEBHOOK",
				URL:  "https://www.example.com/path",
				SecretRef: common.ResourceRefNamespaced{
					Name:      "webhook-secret",
					Namespace: testData.Resources.Namespace,
				},
			},
		}

		projectKey := types.NamespacedName{
			Name:      testData.Project.Name,
			Namespace: testData.Resources.Namespace,
		}

		By("Add integrations", func() {
			_, err := akoretry.RetryUpdateOnConflict(ctx, testData.K8SClient, projectKey, func(project *akov2.AtlasProject) {
				project.Spec.Integrations = integrations
			})
			Expect(err).To(Succeed())
		})

		By("Integrations are ready", func() {
			actions.WaitForConditionsToBecomeTrue(testData, api.IntegrationReadyType, api.ReadyType)

			for _, integration := range integrations {
				_, err := atlasClient.GetIntegrationByType(testData.Project.ID(), integration.Type)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Skip reconciliation & remove one integration", func() {
			_, err := akoretry.RetryUpdateOnConflict(ctx, testData.K8SClient, projectKey, func(project *akov2.AtlasProject) {
				project.Annotations[customresource.ReconciliationPolicyAnnotation] = customresource.ReconciliationPolicySkip
				kubeIntegrations := project.Spec.Integrations
				project.Spec.Integrations = kubeIntegrations[:2]
			})
			Expect(err).To(Succeed())
		})

		By("Project reconciliation skipped", func() {
			// TODO: how to reliable wait for the skip reocniliation to have been evaluated?
			time.Sleep(5 * time.Second)
		})

		By("Resume reconciliation", func() {
			_, err := akoretry.RetryUpdateOnConflict(ctx, testData.K8SClient, projectKey, func(project *akov2.AtlasProject) {
				delete(project.Annotations, customresource.ReconciliationPolicyAnnotation)
			})
			Expect(err).To(Succeed())
		})

		var kubeProject *akov2.AtlasProject
		By("Change another integration", func() {
			var err error
			kubeProject, err = akoretry.RetryUpdateOnConflict(ctx, testData.K8SClient, projectKey, func(project *akov2.AtlasProject) {
				project.Spec.Integrations[1].ChannelName = "other-channel"
			})
			Expect(err).To(Succeed())
		})

		By("Integrations are ready again", func() {
			actions.WaitForConditionsToBecomeTrue(testData, api.IntegrationReadyType, api.ReadyType)

			for _, integration := range integrations {
				_, err := atlasClient.GetIntegrationByType(testData.Project.ID(), integration.Type)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Expect removed WEBHOOK integration to still exist in Atlas", func() {
			atlasClient, err := atlas.AClient()
			Expect(err).To(Succeed())
			_, _, err = atlasClient.Client.ThirdPartyIntegrationsApi.GetThirdPartyIntegration(
				ctx, kubeProject.Status.ID, "WEBHOOK",
			).Execute()
			Expect(err).To(Succeed())
		})
	})
})

func integrationTest(data *model.TestDataProvider, integration project.Integration, key string, setSecret configSecret) {
	By("Create Secret for integration", func() {
		Expect(k8s.CreateUserSecret(data.Context, data.K8SClient, key, integrationSecretName, data.Resources.Namespace)).Should(Succeed())

		setSecret(&integration, common.ResourceRefNamespaced{Name: integrationSecretName, Namespace: data.Resources.Namespace})
	})

	By("Add integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		data.Project.Spec.Integrations = append(data.Project.Spec.Integrations, integration)

		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
	})

	By("Integration is ready", func() {
		actions.WaitForConditionsToBecomeTrue(data, api.IntegrationReadyType, api.ReadyType)

		atlasIntegration, err := atlasClient.GetIntegrationByType(data.Project.ID(), integration.Type)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(strings.HasSuffix(key, strings.TrimLeft(atlasIntegration.GetApiKey(), "*"))).Should(BeTrue())
	})

	By("Delete integration", func() {
		Expect(data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Resources.Namespace}, data.Project)).Should(Succeed())
		data.Project.Spec.Integrations = []project.Integration{}

		Expect(data.K8SClient.Update(data.Context, data.Project)).Should(Succeed())
	})

	By("Delete integration check", func() {
		actions.CheckProjectConditionsNotSet(data, api.IntegrationReadyType)

		atlasIntegration, err := atlasClient.GetIntegrationByType(data.Project.ID(), integration.Type)
		Expect(err).Should(HaveOccurred())
		Expect(atlasIntegration).To(BeNil())
	})
}

type configSecret func(integration *project.Integration, ref common.ResourceRefNamespaced)
