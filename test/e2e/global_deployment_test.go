package e2e_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("UserLogin", Label("global-deployment"), func() {
	var testData *model.TestDataProvider

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Delete Resources", func() {
			actions.DeleteTestDataDeployments(testData)
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test *model.TestDataProvider, mapping []v1.CustomZoneMapping, ns []v1.ManagedNamespace) {
			testData = test
			actions.ProjectCreationFlow(test)
			globalClusterFlow(test, mapping, ns)
		},
		Entry("Test[gc-regular-deployment]: Deployment with global config", Label("gc-regular-deployment"),
			model.DataProvider(
				"gc-regular-deployment",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateRegularGeoshardedDeployment("gc-regular-deployment")),
			[]v1.CustomZoneMapping{
				{
					Zone:     "Zone 1",
					Location: "AO",
				},
			},
			[]v1.ManagedNamespace{
				{
					Collection:     "somecollection",
					Db:             "somedb",
					CustomShardKey: "somekey",
				},
			},
		),
		Entry("Test[gc-advanced-deployment]: Advanced", Label("gc-advanced-deployment"),
			model.DataProvider(
				"gc-advanced-deployment",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).WithInitialDeployments(data.CreateAdvancedGeoshardedDeployment("gc-advanced-deployment")),
			[]v1.CustomZoneMapping{
				{
					Zone:     "Zone 1",
					Location: "AO",
				},
				{
					Zone:     "Zone 2",
					Location: "CA",
				},
			},
			[]v1.ManagedNamespace{
				{
					Collection:             "somecollection",
					Db:                     "somedb",
					CustomShardKey:         "somekey",
					PresplitHashedZones:    toptr.MakePtr(true),
					IsCustomShardKeyHashed: toptr.MakePtr(true),
					NumInitialChunks:       4,
				},
			},
		),
	)
})

func globalClusterFlow(userData *model.TestDataProvider, mapping []v1.CustomZoneMapping, managedNamespace []v1.ManagedNamespace) {
	By("Apply deployment", func() {
		Expect(userData.InitialDeployments).ShouldNot(BeEmpty())
		userData.InitialDeployments[0].Namespace = userData.Resources.Namespace
		Expect(userData.K8SClient.Create(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			return userData.InitialDeployments[0].Status.StateName == status.StateIDLE
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})

	By("Applying global cluster config to Deployment", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.InitialDeployments[0].Name,
			Namespace: userData.InitialDeployments[0].Namespace,
		}, userData.InitialDeployments[0])).To(Succeed())
		if userData.InitialDeployments[0].Spec.DeploymentSpec != nil {
			userData.InitialDeployments[0].Spec.DeploymentSpec.ManagedNamespaces = managedNamespace
			userData.InitialDeployments[0].Spec.DeploymentSpec.CustomZoneMapping = mapping
		} else {
			userData.InitialDeployments[0].Spec.AdvancedDeploymentSpec.ManagedNamespaces = managedNamespace
			userData.InitialDeployments[0].Spec.AdvancedDeploymentSpec.CustomZoneMapping = mapping
		}

		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
	})

	By("Wait and check global zone mapping status", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == status.CustomZoneMappingReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Wait and check global managed namespaces status", func() {
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == status.ManagedNamespacesReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(10 * time.Minute).Should(BeTrue())
	})

	By("Delete global  cluster config and wait idle state of cluster", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
			Name:      userData.InitialDeployments[0].Name,
			Namespace: userData.InitialDeployments[0].Namespace,
		}, userData.InitialDeployments[0])).To(Succeed())
		if userData.InitialDeployments[0].Spec.DeploymentSpec != nil {
			userData.InitialDeployments[0].Spec.DeploymentSpec.ManagedNamespaces = nil
			userData.InitialDeployments[0].Spec.DeploymentSpec.CustomZoneMapping = nil
		} else {
			userData.InitialDeployments[0].Spec.AdvancedDeploymentSpec.ManagedNamespaces = nil
			userData.InitialDeployments[0].Spec.AdvancedDeploymentSpec.CustomZoneMapping = nil
		}
		Expect(userData.K8SClient.Update(userData.Context, userData.InitialDeployments[0])).To(Succeed())
		Eventually(func(g Gomega) bool {
			g.Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{
				Name:      userData.InitialDeployments[0].Name,
				Namespace: userData.InitialDeployments[0].Namespace,
			}, userData.InitialDeployments[0])).To(Succeed())
			for _, condition := range userData.InitialDeployments[0].Status.Conditions {
				if condition.Type == status.DeploymentReadyType {
					return condition.Status == corev1.ConditionTrue
				}
			}
			return false
		}).WithTimeout(30 * time.Minute).Should(BeTrue())
	})
}
