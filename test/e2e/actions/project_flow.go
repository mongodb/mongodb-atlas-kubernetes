package actions

import (
	"fmt"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func ProjectCreationFlow(userData *model.TestDataProvider) {
	By("Deploy Project with requested configuration", func() {
		PrepareUsersConfigurations(userData)
		deploy.NamespacedOperator(userData) // TODO: how to deploy operator by code?
		By(fmt.Sprintf("Create namespace %s", userData.Resources.Namespace))
		Expect(kubecli.CreateNamespace(userData.Context, userData.K8SClient, userData.Resources.Namespace)).Should(Succeed())
		kubecli.CreateDefaultSecret(userData.Context, userData.K8SClient, config.DefaultOperatorGlobalKey, userData.Resources.Namespace)
		if !userData.Resources.AtlasKeyAccessType.GlobalLevelKey {
			CreateConnectionAtlasKey(userData)
		}
		deploy.CreateProject(userData)
	})
}
