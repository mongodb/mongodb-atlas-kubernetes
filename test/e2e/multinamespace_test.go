package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

var _ = Describe("Users can use clusterwide configuration with limitation to watch only particular namespaces", Label("multinamespaced"), func() {
	var listData []model.TestDataProvider
	var watchedNamespace []string

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		By("AfterEach. clean-up", func() {
			if CurrentSpecReport().Failed() {
				GinkgoWriter.Write([]byte("Resources wasn't clean"))
				utils.SaveToFile(
					"output/operator-logs.txt",
					kubecli.GetManagerLogs(config.DefaultOperatorNS),
				)
				actions.SaveK8sResources(
					[]string{"deploy"},
					config.DefaultOperatorNS,
				)
				for _, data := range listData {
					actions.SaveK8sResources(
						[]string{"atlasprojects"},
						data.Resources.Namespace,
					)
				}
			}
		})
	})

	// (Consider Shared Clusters when E2E tests could conflict with each other)
	It("Deploy cluster multinamespaced operator and create resources in each of them", func() {
		By("Set up test data configuration", func() {
			watched1 := model.NewTestDataProvider(
				"multinamestace-watched1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{},
				30013,
				[]func(*model.TestDataProvider){},
			)
			watchedGlobal := model.NewTestDataProvider(
				"multinamestace-watched-global",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{},
				30013,
				[]func(*model.TestDataProvider){},
			)
			notWatched := model.NewTestDataProvider(
				"multinamestace-notwatched",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{},
				30013,
				[]func(*model.TestDataProvider){},
			)
			notWatchedGlobal := model.NewTestDataProvider(
				"multinamestace-notwatched",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess().CreateAsGlobalLevelKey(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{},
				30013,
				[]func(*model.TestDataProvider){},
			)
			listData = []model.TestDataProvider{watched1, watchedGlobal, notWatched, notWatchedGlobal}
			watchedNamespace = []string{watched1.Resources.Namespace, watchedGlobal.Resources.Namespace}
		})
		By("User Install CRD, cluster multinamespace Operator", func() {
			for i := range listData {
				actions.PrepareUsersConfigurations(&listData[i])
			}
			deploy.MultiNamespaceOperator(&listData[0], watchedNamespace)
			kubecli.CreateApiKeySecret(config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
		})
		By("Check if operator working as expected: watched/not watched namespaces", func() {
			crdsFile := listData[0].Resources.GetOperatorFolder() + "/crds.yaml"
			for i := range listData {
				testFlow(&listData[i], crdsFile, watchedNamespace)
			}
		})
	})
})

func testFlow(data *model.TestDataProvider, crdsFile string, watchedNamespece []string) {
	isWatchedNamespace := func(target string, list []string) bool {
		for _, item := range list {
			if item == target {
				return true
			}
		}
		return false
	}
	if isWatchedNamespace(data.Resources.Namespace, watchedNamespece) {
		watchedFlow(data, crdsFile)
	} else {
		notWatchedFlow(data, crdsFile)
	}
}

func watchedFlow(data *model.TestDataProvider, crdsFile string) {
	By("Deploy users resorces", func() {
		if !data.Resources.AtlasKeyAccessType.GlobalLevelKey {
			actions.CreateConnectionAtlasKey(data)
		}
		kubecli.Apply(crdsFile, "-n", data.Resources.Namespace)
		kubecli.Apply(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
	})
	By("Check if projects were deployed", func() {
		Eventually(
			kube.GetReadyProjectStatus(data),
		).Should(Equal("True"), "kubernetes resource: Project status `Ready` should be True. Watched namespace")
	})
	By("Get IDs for deletion", func() {
		resource, err := kube.GetProjectResource(data)
		Expect(err).Should(BeNil())
		data.Resources.ProjectID = resource.Status.ID
		Expect(data.Resources.ProjectID).ShouldNot(BeEmpty())
	})
	By("Delete Resources", func() {
		actions.DeleteUserResourcesProject(data)
	})
}

func notWatchedFlow(data *model.TestDataProvider, crdsFile string) {
	By("Deploy users resorces", func() {
		if !data.Resources.AtlasKeyAccessType.GlobalLevelKey {
			actions.CreateConnectionAtlasKey(data)
		}
		kubecli.Apply(crdsFile, "-n", data.Resources.Namespace)
		kubecli.Apply(data.Resources.ProjectPath, "-n", data.Resources.Namespace)
	})
	By("Check if projects were deployed", func() {
		Eventually(
			kube.GetReadyProjectStatus(data),
		).Should(BeEmpty(), "Kubernetes resource: Project status `Ready` should be False. NOT Watched namespace")
	})
}
