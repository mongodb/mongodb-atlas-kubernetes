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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("project-settings"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Project Settings Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, settings akov2.ProjectSettings) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			projectSettingsFlow(testData, &settings)
		},
		Entry("Test[project-settings]: User has project to which Project Settings was added", Label("focus-project-settings"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "project-settings", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			akov2.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
				IsDataExplorerEnabled:                       pointer.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 pointer.MakePtr(false),
				IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
				IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
			},
		),
	)
})

func projectSettingsFlow(userData *model.TestDataProvider, settings *akov2.ProjectSettings) {
	By("Add Project Settings to the project", func() {
		userData.Project.Spec.Settings = settings
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, api.ProjectSettingsReadyType, api.ReadyType)
	})

	By("Remove Project Settings from the project", func() {
		userData.Project.Spec.Settings = nil
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, api.ProjectSettingsReadyType)
	})
}
