// different ways to deploy operator
package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	corev1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kustomize"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// CopyKustomizeNamespaceOperator create copy of `/deploy/namespaced` folder with kustomization file for overriding namespace
func prepareMultiNamespaceOperatorResources(input model.UserInputs, watchedNamespaces []string) {
	fullPath := input.GetOperatorFolder()
	err := os.Mkdir(fullPath, os.ModePerm)
	Expect(err).ShouldNot(HaveOccurred())
	utils.CopyFile(config.DefaultClusterWideCRDConfig, filepath.Join(fullPath, "crds.yaml"))
	utils.CopyFile(config.DefaultClusterWideOperatorConfig, filepath.Join(fullPath, "multinamespace-config.yaml"))
	namespaces := strings.Join(watchedNamespaces, ",")
	patchWatch := []byte(
		"apiVersion: apps/v1\n" +
			"kind: Deployment\n" +
			"metadata:\n" +
			"  name: mongodb-atlas-operator\n" +
			"spec:\n" +
			"  template:\n" +
			"    spec:\n" +
			"      containers:\n" +
			"      - name: manager\n" +
			"        env:\n" +
			"        - name: WATCH_NAMESPACE\n" +
			"          value: \"" + namespaces + "\"",
	)
	err = utils.SaveToFile(filepath.Join(fullPath, "patch.yaml"), patchWatch)
	Expect(err).ShouldNot(HaveOccurred())
	kustomization := []byte(
		"resources:\n" +
			"- multinamespace-config.yaml\n" +
			"patches:\n" +
			"- path: patch.yaml\n" +
			"  target:\n" +
			"    group: apps\n" +
			"    version: v1\n" +
			"    kind: Deployment\n" +
			"    name: mongodb-atlas-operator",
	)
	err = utils.SaveToFile(filepath.Join(fullPath, "kustomization.yaml"), kustomization)
	Expect(err).ShouldNot(HaveOccurred())
}

func CheckOperatorRunning(data *model.TestDataProvider, namespace string) {
	By("Check Operator is running", func() {
		Eventually(
			func(g Gomega) string {
				status, err := k8s.GetPodStatus(data.Context, data.K8SClient, namespace)
				g.Expect(err).ShouldNot(HaveOccurred())
				return status
			},
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})
}

func MultiNamespaceOperator(data *model.TestDataProvider, watchNamespace []string) {
	prepareMultiNamespaceOperatorResources(data.Resources, watchNamespace)
	By("Deploy multinamespaced Operator \n", func() {
		kustomOperatorPath := data.Resources.GetOperatorFolder() + "/final.yaml"
		utils.SaveToFile(kustomOperatorPath, kustomize.Build(data.Resources.GetOperatorFolder()))
		kubecli.Apply(kustomOperatorPath)
		CheckOperatorRunning(data, config.DefaultOperatorNS)
	})
}

func DeleteProject(testData *model.TestDataProvider) {
	By("Delete Project", func() {
		projectId := testData.Project.Status.ID
		Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: testData.Project.Name, Namespace: testData.Project.Namespace}, testData.Project)).Should(Succeed(), "Get project failed")
		Expect(testData.K8SClient.Delete(testData.Context, testData.Project)).Should(Succeed(), "Delete project failed")
		aClient := atlas.GetClientOrFail()
		Eventually(func(g Gomega) bool {
			return aClient.IsProjectExists(g, projectId)
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(BeTrue(), "Project was not deleted in Atlas")
	})
}

func DeleteUsers(testData *model.TestDataProvider) {
	By("Delete Users", func() {
		for _, user := range testData.Users {
			Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: user.Name, Namespace: user.Namespace}, user)).Should(Succeed(), "Get user failed")
			Expect(testData.K8SClient.Delete(testData.Context, user)).Should(Succeed(), "Delete user failed")
		}
	})
}

func DeleteInitialDeployments(testData *model.TestDataProvider) {
	By("Delete initial deployments", func() {
		for _, deployment := range testData.InitialDeployments {
			projectId := testData.Project.Status.ID
			deploymentName := deployment.Spec.DeploymentSpec.Name
			Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: deployment.Name,
				Namespace: testData.Resources.Namespace}, deployment)).Should(Succeed(), "Get deployment failed")
			Expect(testData.K8SClient.Delete(testData.Context, deployment)).Should(Succeed(), "Deployment %s was not deleted", deployment.Name)
			aClient := atlas.GetClientOrFail()
			Eventually(func() bool {
				return aClient.IsDeploymentExist(projectId, deploymentName)
			}).WithTimeout(15*time.Minute).WithPolling(20*time.Second).Should(BeFalse(), "Deployment should be deleted in Atlas")
		}
	})
}

func CreateProject(testData *model.TestDataProvider) {
	if testData.Project.GetNamespace() == "" {
		testData.Project.Namespace = testData.Resources.Namespace
	}
	By(fmt.Sprintf("Deploy Project %s", testData.Project.GetName()), func() {
		err := testData.K8SClient.Create(testData.Context, testData.Project)
		Expect(err).ShouldNot(HaveOccurred(), "Project %s was not created", testData.Project.GetName())
		Eventually(kube.ProjectReadyCondition(testData)).WithTimeout(5*time.Minute).WithPolling(20*time.Second).
			Should(Not(Equal("False")), "Project %s should be ready", testData.Project.GetName())
	})
	By(fmt.Sprintf("Wait for Project %s", testData.Project.GetName()), func() {
		Eventually(func() bool {
			statuses := kube.GetProjectStatus(testData)
			return statuses.ID != ""
		}, 5*time.Minute, 5*time.Second).Should(BeTrue(), "Project %s is not ready", kube.GetProjectStatus(testData))
	})
}

func CreateInitialDeployments(testData *model.TestDataProvider) {
	By("Deploy Initial Deployments", func() {
		for _, deployment := range testData.InitialDeployments {
			if deployment.Namespace == "" {
				deployment.Namespace = testData.Resources.Namespace
				deployment.Spec.Project.Namespace = testData.Resources.Namespace
			}
			err := testData.K8SClient.Create(testData.Context, deployment)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Deployment was not created: %v", deployment))
			Eventually(kube.DeploymentReadyCondition(testData), time.Minute*60, time.Second*5).Should(Equal("True"), "Deployment was not created")
		}
	})
}

func CreateUsers(testData *model.TestDataProvider) {
	By("Deploy Users", func() {
		for _, user := range testData.Users {
			if user.Namespace == "" {
				user.Namespace = testData.Resources.Namespace
				user.Spec.Project.Namespace = testData.Resources.Namespace
			}
			if user.Spec.PasswordSecret != nil {
				secret := utils.UserSecretPassword()
				Expect(k8s.CreateUserSecret(testData.Context, testData.K8SClient, secret,
					user.Spec.PasswordSecret.Name, testData.Resources.Namespace)).Should(Succeed(),
					"Create user secret failed")
			}
			err := testData.K8SClient.Create(testData.Context, user)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("User was not created: %v", user))
			Eventually(func(g Gomega) {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: user.GetName(), Namespace: user.GetNamespace()}, user))
				for _, condition := range user.Status.Conditions {
					if condition.Type == status.ReadyType {
						g.Expect(condition.Status).Should(Equal(corev1.ConditionTrue), "User should be ready")
					}
				}
			}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).Should(Succeed(), "User was not created")
		}
	})
}
