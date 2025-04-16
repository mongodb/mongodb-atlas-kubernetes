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

package actions

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

func DeleteDeploymentCRWithKeepAnnotation(testData *model.TestDataProvider) {
	By("Check that deployment CR have keep annotation", func() {
		deployment := &akov2.AtlasDeployment{}
		err := testData.K8SClient.Get(testData.Context, client.ObjectKey{Namespace: testData.Resources.Namespace,
			Name: testData.InitialDeployments[0].GetName()}, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(deployment.GetAnnotations()).To(HaveKeyWithValue(customresource.ResourcePolicyAnnotation, customresource.ResourcePolicyKeep), "Deployment CR should have keep annotation")
	})

	By(fmt.Sprintf("Try to delete CR %s", testData.InitialDeployments[0].GetName()), func() {
		err := kubecli.DeleteDeployment(testData.Context, testData.K8SClient, testData.Resources.Namespace, testData.InitialDeployments[0].GetName())
		Expect(err).NotTo(HaveOccurred())
		By("Checking Cluster still existed", func() {
			aClient := atlas.GetClientOrFail()
			deployment, err := aClient.GetDeployment(testData.Project.Status.ID, testData.InitialDeployments[0].AtlasName())
			Expect(err).NotTo(HaveOccurred())
			Expect(deployment.StateName).ShouldNot(Equal("DELETING"), "Deployment is being deleted despite the keep annotation")
		})
	})
}

func RedeployDeployment(testData *model.TestDataProvider) {
	By("Redeploying the cluster", func() {
		deployment := data.CreateDeploymentWithKeepPolicy(testData.InitialDeployments[0].GetName())
		deployment.Namespace = testData.Resources.Namespace
		Eventually(func() error {
			return testData.K8SClient.Create(testData.Context, deployment)
		}).WithTimeout(5 * time.Minute).WithPolling(20 * time.Second).Should(Succeed())
		deploymentForCheck := &akov2.AtlasDeployment{}
		Eventually(func() bool {
			err := testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: deployment.GetName(), Namespace: deployment.GetNamespace()}, deploymentForCheck)
			Expect(err).Should(BeNil(), fmt.Sprintf("Deployment not found: %v", deploymentForCheck))
			return deploymentForCheck.Status.StateName == status.StateIDLE
		}, time.Minute*10, time.Second*5).Should(BeTrue(), fmt.Sprintf("Deployment was not created: %v", deploymentForCheck))
		testData.InitialDeployments[0] = deploymentForCheck
	})
}

func RemoveKeepAnnotation(testData *model.TestDataProvider) {
	// Remove annotation so actions.AfterEachFinalCleanup can cleanup successfully
	By("Removing keep annotation", func() {
		deploymentUpdate := &akov2.AtlasDeployment{}
		err := testData.K8SClient.Get(testData.Context, client.ObjectKey{Namespace: testData.Resources.Namespace,
			Name: testData.InitialDeployments[0].GetName()}, deploymentUpdate)
		Expect(err).NotTo(HaveOccurred())
		annotations := deploymentUpdate.GetAnnotations()
		// remove keep annotations from map
		delete(annotations, customresource.ResourcePolicyAnnotation)
		deploymentUpdate.SetAnnotations(annotations)
		err = testData.K8SClient.Update(testData.Context, deploymentUpdate)
		Expect(err).NotTo(HaveOccurred())
		testData.InitialDeployments[0] = deploymentUpdate
	})

	By("Checking that keep annotation was removed", func() {
		deployment := &akov2.AtlasDeployment{}
		err := testData.K8SClient.Get(testData.Context, client.ObjectKey{Namespace: testData.Resources.Namespace,
			Name: testData.InitialDeployments[0].GetName()}, deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(deployment.GetAnnotations()).ToNot(HaveKey(customresource.ResourcePolicyAnnotation), "Deployment CR should not have keep annotation")
		testData.InitialDeployments[0] = deployment
	})
}
