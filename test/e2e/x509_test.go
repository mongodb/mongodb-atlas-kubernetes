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
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

var _ = Describe("UserLogin", Label("x509auth"), func() {
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
		func(ctx SpecContext, test func(context.Context) *model.TestDataProvider, certRef common.ResourceRefNamespaced) {
			testData = test(ctx)
			actions.ProjectCreationFlow(testData)
			x509Flow(testData, &certRef)
		},
		Entry("Test[x509auth]: Can create project and add X.509 Auth to that project", Label("focus-x509auth-basic"),
			func(ctx context.Context) *model.TestDataProvider {
				return model.DataProvider(ctx, "x509auth", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30000, []func(*model.TestDataProvider){}).WithProject(data.DefaultProject())
			},
			common.ResourceRefNamespaced{
				Name: "x509cert",
			},
		),
	)
})

func x509Flow(testData *model.TestDataProvider, certRef *common.ResourceRefNamespaced) {
	By("Create X.509 cert secret", func() {
		Expect(certRef.Name).NotTo(BeEmpty(), "certRef.Name should not be empty")
		if certRef.Namespace == "" {
			certRef.Namespace = testData.Resources.Namespace
		}
		Expect(k8s.CreateCertificateX509(testData.Context, testData.K8SClient, certRef.Name, certRef.Namespace)).To(Succeed())
	})

	By("Add X.509 cert to the project", func() {
		Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
			Name:      testData.Project.Name,
			Namespace: testData.Resources.Namespace,
		}, testData.Project)).To(Succeed())
		testData.Project.Spec.X509CertRef = certRef
		Expect(testData.K8SClient.Update(testData.Context, testData.Project)).To(Succeed())
		actions.WaitForConditionsToBecomeTrue(testData, api.ReadyType)
	})

	By("Create User with X.509 cert", func() {
		userName := "CN=my-x509-authenticated-user,OU=organizationalunit,O=organization"
		x509User := data.BasicUser("x509user", "user1",
			data.WithReadWriteRole(),
			data.WithX509(userName),
		)
		testData.Users = append(testData.Users, x509User)
		deploy.CreateUsers(testData)
	})

	By("Deploy User", func() {
		By("Check database users Attributes", func() {
			Eventually(actions.CheckUserExistInAtlas(testData), "2m", "10s").Should(BeTrue())
			actions.CheckUsersAttributes(testData)
		})
	})
}
