// diffrent ways to deploy operator
package actions

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO make it private
// CopyKustomizeNamespaceOperator create copy of `/deploy/namespaced` folder with kustomization file for overriding namespace
func CopyKustomizeNamespaceOperator(input model.UserInputs) {
	fullPath := input.GetOperatorFolder()
	os.Mkdir(fullPath, os.ModePerm)
	utils.CopyFile("../../deploy/namespaced/crds.yaml", filepath.Join(fullPath, "crds.yaml"))
	utils.CopyFile("../../deploy/namespaced/namespaced-config.yaml", filepath.Join(fullPath, "namespaced-config.yaml"))
	data := []byte(
		"namespace: " + input.Namespace + "\n" +
			"resources:" + "\n" +
			"- crds.yaml" + "\n" +
			"- namespaced-config.yaml",
	)
	utils.SaveToFile(filepath.Join(fullPath, "kustomization.yaml"), data)
}


func DeployNamespacedOperatorKuber(data *model.TestDataProvider) {
	By("Create namespaced Operator\n", func() {
		CopyKustomizeNamespaceOperator(data.Resources)
		kube.Apply("-k", data.Resources.GetOperatorFolder())
		Eventually(
			kube.GetPodStatus(data.Resources.Namespace),
			"5m", "3s",
		).Should(Equal("Running"), "The operator should successfully run")
	})
}
