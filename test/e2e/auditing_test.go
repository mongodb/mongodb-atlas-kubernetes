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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("auditing"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Auditing Test\n"))
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
		func(ctx SpecContext, test func(ctx2 context.Context) *model.TestDataProvider, auditing akov2.Auditing) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			auditingFlow(testData, &auditing)
		},
		Entry("Test[auditing]: User has project to which Auditing was added", Label("focus-auditing"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "auditing", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			akov2.Auditing{
				AuditAuthorizationSuccess: false,
				AuditFilter:               exampleFilter(),
				Enabled:                   true,
			},
		),
	)
})

func auditingFlow(userData *model.TestDataProvider, auditing *akov2.Auditing) {
	By("Add auditing to the project", func() {
		userData.Project.Spec.Auditing = auditing
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, api.AuditingReadyType, api.ReadyType)
	})

	By("Remove Auditing from the project", func() {
		userData.Project.Spec.Auditing = nil
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.CheckProjectConditionsNotSet(userData, api.AuditingReadyType)
	})
}

func exampleFilter() string {
	return `{"atype" : "authenticate", "param" : {"user" : "auditReadOnly", "db" : "admin", "mechanism" : "SCRAM-SHA-1"} }`
}
