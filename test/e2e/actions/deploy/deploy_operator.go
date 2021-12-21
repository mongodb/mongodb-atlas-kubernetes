// different ways to deploy operator
package deploy

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kustomize"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// prepareNamespaceOperatorResources create copy of `/deploy/namespaced` folder with kustomization file for overriding namespace
func prepareNamespaceOperatorResources(input model.UserInputs) {
	fullPath := input.GetOperatorFolder()
	os.Mkdir(fullPath, os.ModePerm)
	utils.CopyFile(config.DefaultNamespacedCRDConfig, filepath.Join(fullPath, "crds.yaml"))
	utils.CopyFile(config.DefaultNamespacedOperatorConfig, filepath.Join(fullPath, "namespaced-config.yaml"))
	data := []byte(
		"namespace: " + input.Namespace + "\n" +
			"resources:" + "\n" +
			"- crds.yaml" + "\n" +
			"- namespaced-config.yaml",
	)
	utils.SaveToFile(filepath.Join(fullPath, "kustomization.yaml"), data)
}

// CopyKustomizeNamespaceOperator create copy of `/deploy/namespaced` folder with kustomization file for overriding namespace
func prepareWideOperatorResources(input model.UserInputs) {
	fullPath := input.GetOperatorFolder()
	os.Mkdir(fullPath, os.ModePerm)
	utils.CopyFile(config.DefaultClusterWideCRDConfig, filepath.Join(fullPath, "crds.yaml"))
	utils.CopyFile(config.DefaultClusterWideOperatorConfig, filepath.Join(fullPath, "clusterwide-config.yaml"))
}

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

func NamespacedOperator(data *model.TestDataProvider) {
	prepareNamespaceOperatorResources(data.Resources)
	By("Deploy namespaced Operator\n", func() {
		kubecli.Apply("-k", data.Resources.GetOperatorFolder())
		Eventually(
			kubecli.GetPodStatus(data.Resources.Namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})
}

func ClusterWideOperator(data *model.TestDataProvider) {
	prepareWideOperatorResources(data.Resources)
	By("Deploy clusterwide Operator \n", func() {
		kubecli.Apply("-k", data.Resources.GetOperatorFolder())
		Eventually(
			kubecli.GetPodStatus(config.DefaultOperatorNS),
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
		Eventually(
			kubecli.GetPodStatus(config.DefaultOperatorNS),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})
}
