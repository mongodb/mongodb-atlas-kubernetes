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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	helm "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/cli/helm"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

// HelmDefaultUpgradeResources helm update should change at least 1 field: databaseusers, project, deployment
func HelmDefaultUpgradeResources(data *model.TestDataProvider) {
	By("User use HELM upgrade command for changing atlas resources\n", func() {
		data.Resources.Project.Spec.ProjectIPAccessList[0].Comment = "updated"
		data.Resources.Users[0].DeleteAllRoles()
		data.Resources.Users[0].AddBuildInAdminRole()
		data.Resources.Users[0].Spec.ProjectRef.Name = data.Resources.GetAtlasProjectFullKubeName()
		generation, err := kubecli.GetDeploymentObservedGeneration(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
		Expect(err).NotTo(HaveOccurred())
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)

		By("Wait project creation", func() {
			WaitDeployment(data, generation)
			ExpectWithOffset(1, data.Resources.ProjectID).ShouldNot(BeEmpty())
		})
		aClient := atlas.GetClientOrFail()
		updatedDeployment, err := aClient.GetFlexInstance(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName())
		Expect(err).NotTo(HaveOccurred())
		CompareFlexSpec(data.Resources.Deployments[0].Spec, *updatedDeployment)
		Eventually(func() error {
			aClient := atlas.GetClientOrFail()
			user, err := aClient.GetDBUser("admin", data.Resources.Users[0].Spec.Username, data.Resources.ProjectID)
			if err != nil {
				return err
			}
			if user.GetRoles()[0].RoleName != model.RoleBuildInAdmin {
				return fmt.Errorf("user role %s not equal to %s", user.GetRoles()[0].RoleName, model.RoleBuildInAdmin)
			}
			return nil
		}).WithTimeout(7 * time.Minute).WithPolling(10 * time.Second).ShouldNot(HaveOccurred())
	})
}

// HelmUpgradeUsersRoleAddAdminUser helm update: add user+change user role
func HelmUpgradeUsersRoleAddAdminUser(data *model.TestDataProvider) {
	By("User change role for all users and add new database user\n", func() {
		for i := range data.Resources.Users {
			data.Resources.Users[i].WithProjectRef(data.Resources.Project.GetK8sMetaName())
			data.Resources.Users[i].AddCustomRole(model.RoleCustomReadWrite, "Ships", "")
		}
		newUser := *model.NewDBUser("only-one-admin").
			WithAuthDatabase("admin").
			WithProjectRef(data.Resources.Project.GetK8sMetaName()).
			WithSecretRef("new-user-secret").
			AddBuildInAdminRole()
		data.Resources.Users = append(data.Resources.Users, newUser)
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)
		CheckUsersAttributes(data)
	})
}

// HelmUpgradeDeleteFirstUser helm update: delete user
func HelmUpgradeDeleteFirstUser(data *model.TestDataProvider) {
	By("User delete database user from the Atlas\n", func() {
		data.Resources.Users = data.Resources.Users[1:]
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)
		CheckUsersAttributes(data)
	})
}

// HelmUpgradeChartVersions upgrade chart version of crd, operator, and
func HelmUpgradeChartVersions(data *model.TestDataProvider) {
	By("User update helm chart (used main-branch)", func() {
		generation, err := kubecli.GetDeploymentObservedGeneration(data.Context, data.K8SClient, data.Resources.Namespace, data.Resources.Deployments[0].ObjectMeta.GetName())
		Expect(err).NotTo(HaveOccurred())
		helm.UpgradeOperatorChart(data.Resources)
		helm.UpgradeAtlasDeploymentChartDev(data.Resources)

		By("Wait updating")
		WaitDeployment(data, generation+1)
		aClient := atlas.GetClientOrFail()
		updatedDeployment, err := aClient.GetFlexInstance(data.Resources.ProjectID, data.Resources.Deployments[0].Spec.GetDeploymentName())
		Expect(err).NotTo(HaveOccurred())
		CompareFlexSpec(data.Resources.Deployments[0].Spec, *updatedDeployment)
		CheckUsersAttributes(data)
	})
}
