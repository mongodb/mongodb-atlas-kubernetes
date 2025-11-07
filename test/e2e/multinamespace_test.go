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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("Users can use clusterwide configuration with limitation to watch only particular namespaces", Label("multinamespaced"), func() {
	var listData []*model.TestDataProvider
	var watchedNamespace []string

	_ = AfterEach(func() {
		By("AfterEach. clean-up", func() {
			for i := range listData {
				testData := listData[i]
				actions.DeleteTestDataProject(testData)
				actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
				if testData.ManagerContext != nil {
					testData.ManagerContext.Done()
				}
			}
		})

	})

	// (Consider Shared Deployments when E2E tests could conflict with each other)
	It("Deploy deployment multi-namespaced operator and create resources in each of them", func(ctx SpecContext) {
		By("Set up test data configuration", func() {
			watched := model.DataProvider(ctx, "multinamespace-watched", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			watchedGlobal := model.DataProvider(ctx, "multinamespace-watched-global", model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			notWatched := model.DataProvider(ctx, "multinamespace-not-watched", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			notWatchedGlobal := model.DataProvider(ctx, "multinamespace-not-watched-global", model.NewEmptyAtlasKeyType().UseDefaultFullAccess().CreateAsGlobalLevelKey(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			listData = []*model.TestDataProvider{watched, watchedGlobal, notWatched, notWatchedGlobal}
			watchedNamespace = []string{watched.Resources.Namespace, watchedGlobal.Resources.Namespace}
		})
		firstData := listData[0]
		By("Run operator and deploy resources", func() {
			Expect(k8s.CreateNamespace(firstData.Context, firstData.K8SClient, config.DefaultOperatorNS)).Should(Succeed())
			watchedNamespace = append(watchedNamespace, config.DefaultOperatorNS)
			k8s.CreateDefaultSecret(firstData.Context, firstData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
			for i := range listData {
				deployMultiNSResources(listData[i])
			}
			deploy.MultiNamespaceOperator(firstData, watchedNamespace)
		})
		By("Check if operator working as expected: watched/not watched namespaces", func() {
			for i := range listData {
				multiNSFlow(listData[i], watchedNamespace)
			}
		})
	})
})

func deployMultiNSResources(testData *model.TestDataProvider) {
	By("User create namespace", func() {
		Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
	})
	By("Deploy resources", func() {
		if testData.Project.GetNamespace() == "" {
			testData.Project.Namespace = testData.Resources.Namespace
		}
		if !testData.Resources.AtlasKeyAccessType.GlobalLevelKey {
			testData.Project.Spec.ConnectionSecret = &common.ResourceRefNamespaced{Name: testData.Prefix, Namespace: testData.Project.Namespace}
		}
		By(fmt.Sprintf("Deploy Project %s", testData.Project.GetName()), func() {
			err := testData.K8SClient.Create(testData.Context, testData.Project)
			Expect(err).ShouldNot(HaveOccurred(), "Project %s was not created", testData.Project.GetName())
		})
	})
}

func multiNSFlow(data *model.TestDataProvider, watchedNamespaces []string) {
	isWatchedNamespace := stringutil.Contains(watchedNamespaces, data.Resources.Namespace)
	if isWatchedNamespace {
		watchedFlow(data)
	} else {
		notWatchedFlow(data)
	}
}

func watchedFlow(data *model.TestDataProvider) {
	By("Deploy secret", func() {
		actions.CreateConnectionAtlasKey(data)
	})
	By(fmt.Sprintf("Check if projects were deployed. Project name: %s, namespace: %s",
		data.Project.GetName(), data.Project.GetNamespace()), func() {
		Eventually(func(g Gomega) {
			condition, _ := k8s.GetProjectStatusCondition(data.Context, data.K8SClient, api.ReadyType, data.Resources.Namespace, data.Project.GetName())
			g.Expect(condition).Should(Equal("True"))
		}).Should(Succeed(), "kubernetes resource: Project status `Ready` should be True. Watched namespace")
	})
}

func notWatchedFlow(data *model.TestDataProvider) {
	By("Deploy secret", func() {
		actions.CreateConnectionAtlasKey(data)
	})
	By(fmt.Sprintf("Check if projects were deployed. Project name: %s, namespace: %s",
		data.Project.GetName(), data.Project.GetNamespace()), func() {
		Eventually(func(g Gomega) {
			condition, _ := k8s.GetProjectStatusCondition(data.Context, data.K8SClient, api.ReadyType, data.Resources.Namespace, data.Project.GetName())
			g.Expect(condition).Should(Equal(""))
		}).Should(Succeed(), "Kubernetes resource: Project status `Ready` should be empty. NOT Watched namespace")
	})
}
