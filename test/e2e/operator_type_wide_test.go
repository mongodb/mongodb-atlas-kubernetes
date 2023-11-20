package e2e_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper/e2e/model"
)

// DO NOT RUN THIS TEST IN PARALLEL WITH OTHER TESTS

var _ = Describe("Deployment wide operator can work with resources in different namespaces without conflict", Label("deployment-wide"), func() {
	var NortonData, NimnulData *model.TestDataProvider
	commonDeploymentName := "megadeployment"
	k8sClient, err := k8s.CreateNewClient()
	Expect(err).To(BeNil())
	ctx := context.Background()

	_ = AfterEach(func() {
		By("AfterEach. clean-up", func() {
			if CurrentSpecReport().Failed() {
				Expect(actions.SaveProjectsToFile(NortonData.Context, NortonData.K8SClient, NortonData.Resources.Namespace)).Should(Succeed())
				Expect(actions.SaveDeploymentsToFile(NortonData.Context, NortonData.K8SClient, NortonData.Resources.Namespace)).Should(Succeed())
				Expect(actions.SaveUsersToFile(NortonData.Context, NortonData.K8SClient, NortonData.Resources.Namespace)).Should(Succeed())

				Expect(actions.SaveProjectsToFile(NimnulData.Context, NimnulData.K8SClient, NimnulData.Resources.Namespace)).Should(Succeed())
				Expect(actions.SaveDeploymentsToFile(NimnulData.Context, NimnulData.K8SClient, NimnulData.Resources.Namespace)).Should(Succeed())
				Expect(actions.SaveUsersToFile(NimnulData.Context, NimnulData.K8SClient, NimnulData.Resources.Namespace)).Should(Succeed())
			}
			actions.DeleteTestDataUsers(NortonData)
			actions.DeleteTestDataUsers(NimnulData)
			actions.DeleteTestDataDeployments(NortonData)
			actions.DeleteTestDataProject(NortonData)
			actions.DeleteTestDataDeployments(NimnulData)
			actions.DeleteTestDataProject(NimnulData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*NortonData, *NimnulData})
		})
	})

	It("Deploy deployment wide operator and create resources in each of them", func() {
		By("Users can create deployments with the same name", func() {
			NortonData = model.DataProvider(
				"norton-wide",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).
				WithInitialDeployments(data.CreateDeploymentWithBackup(commonDeploymentName)).
				WithUsers(data.BasicUser("reader2", "reader2", data.WithSecretRef("dbuser-secret-u2"), data.WithReadWriteRole()))

			NimnulData = model.DataProvider(
				"nimnul-wide",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()).
				WithInitialDeployments(data.CreateBasicDeployment(commonDeploymentName)).
				WithUsers(data.BasicUser("reader2", "reader2", data.WithSecretRef("dbuser-secret-u2"), data.WithReadWriteRole()))
		})

		By("Initial preparation", func() {
			actions.CreateNamespaceAndSecrets(NortonData)
			actions.CreateNamespaceAndSecrets(NimnulData)
			k8s.CreateNamespace(ctx, k8sClient, config.DefaultOperatorNS)
			k8s.CreateDefaultSecret(ctx, k8sClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)

			mgr, err := k8s.BuildManager(&k8s.Config{
				GlobalAPISecret: client.ObjectKey{
					Namespace: config.DefaultOperatorNS,
					Name:      config.DefaultOperatorGlobalKey,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			go func(ctx context.Context) context.Context {
				err := mgr.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
				return ctx
			}(ctx)
		})

		By("Norton creates resources", func() {
			deploy.CreateProject(NortonData)
			deploy.CreateUsers(NortonData)

			deployment := NortonData.InitialDeployments[0]
			if deployment.Namespace == "" {
				deployment.Namespace = NortonData.Resources.Namespace
				deployment.Spec.Project.Namespace = NortonData.Resources.Namespace
			}
			err := k8sClient.Create(ctx, deployment)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Deployment was not created: %v", deployment))
		})

		By("Nimnul creates resources", func() {
			deploy.CreateProject(NimnulData)
			deploy.CreateUsers(NimnulData)
			deployment := NimnulData.InitialDeployments[0]
			if deployment.Namespace == "" {
				deployment.Namespace = NimnulData.Resources.Namespace
				deployment.Spec.Project.Namespace = NimnulData.Resources.Namespace
			}
			err := k8sClient.Create(ctx, deployment)
			Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("Deployment was not created: %v", deployment))
		})

		By("Check resources", func() {
			Eventually(kube.DeploymentReadyCondition(NortonData), time.Minute*60, time.Second*5).Should(Equal("True"), "Deployment was not created")
			Eventually(kube.DeploymentReadyCondition(NimnulData), time.Minute*60, time.Second*5).Should(Equal("True"), "Deployment was not created")
		})
	})
})
