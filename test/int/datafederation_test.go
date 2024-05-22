package int

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasDataFederation", Label("AtlasDataFederation"), func() {
	const (
		interval                   = PollingInterval
		dataFederationInstanceName = "test-data-federation-aws"
	)

	var (
		connectionSecret      corev1.Secret
		createdProject        *akov2.AtlasProject
		createdDataFederation *akov2.AtlasDataFederation
		manualDeletion        bool
	)

	BeforeEach(func() {
		prepareControllers(false)
		Expect(testNamespace).ToNot(BeNil())

		manualDeletion = false

		connectionSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      ConnectionSecretName,
				Namespace: testNamespace.Name,
				Labels: map[string]string{
					connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
				},
			},
			StringData: secretData(),
		}
		By(fmt.Sprintf("Creating the Secret %s", kube.ObjectKeyFromObject(&connectionSecret)))
		Expect(k8sClient.Create(ctx, &connectionSecret)).To(Succeed())

		createdProject = akov2.DefaultProject(testNamespace.Name, connectionSecret.Name).WithIPAccessList(project.NewIPAccessList().WithCIDR("0.0.0.0/0"))
		if DeploymentDevMode {
			// While developing tests we need to reuse the same project
			createdProject.Spec.Name = "dev-test atlas-project"
		}
		By("Creating the project " + createdProject.Name)
		Expect(k8sClient.Create(ctx, createdProject)).To(Succeed())
		Eventually(func() bool {
			return resources.CheckCondition(k8sClient, createdProject, api.TrueCondition(api.ReadyType))
		}).WithTimeout(30 * time.Minute).WithPolling(interval).Should(BeTrue())
	})

	AfterEach(func() {
		if DeploymentDevMode {
			return
		}
		if manualDeletion && createdProject != nil {
			By("Deleting the deployment in Atlas manually", func() {
				// We need to remove the deployment in Atlas manually to let project get removed
				_, err := atlasClient.ClustersApi.
					DeleteCluster(ctx, createdProject.ID(), createdDataFederation.Name).
					Execute()
				Expect(err).NotTo(HaveOccurred())
				Eventually(checkAtlasDeploymentRemoved(createdProject.Status.ID, createdDataFederation.Name), 600, interval).Should(BeTrue())
				createdDataFederation = nil
			})
		}
		if createdProject != nil && createdProject.Status.ID != "" {
			if createdDataFederation != nil {
				By("Removing Atlas DataFederation " + createdDataFederation.Name)
				Expect(k8sClient.Delete(ctx, createdDataFederation)).To(Succeed())
				deploymentName := createdDataFederation.Name
				if customresource.IsResourcePolicyKeep(createdDataFederation) || customresource.ReconciliationShouldBeSkipped(createdDataFederation) {
					By("Removing Atlas DataFederation " + createdDataFederation.Name + " from Atlas manually")
					Expect(deleteAtlasDataFederation(createdProject.Status.ID, deploymentName)).To(Succeed())
				}
				Eventually(checkAtlasDataFederationRemoved(createdProject.Status.ID, deploymentName), 600, interval).Should(BeTrue())
			}

			By("Removing Atlas Project " + createdProject.Status.ID)
			Expect(k8sClient.Delete(ctx, createdProject)).To(Succeed())
			Eventually(checkAtlasProjectRemoved(createdProject.Status.ID), 60, interval).Should(BeTrue())
		}
	})

	Describe("DataFederation can be created with stores and databases", func() {
		It("Should Succeed", func() {
			By("Creating a DataFederation instance with DB and STORE", func() {
				createdDataFederation = akov2.NewDataFederationInstance(createdProject.Name, dataFederationInstanceName, testNamespace.Name)
				Expect(k8sClient.Create(ctx, createdDataFederation)).ShouldNot(HaveOccurred())

				Eventually(func(g Gomega) {
					df, _, err := atlasClient.DataFederationApi.
						GetFederatedDatabase(ctx, createdProject.ID(), createdDataFederation.Spec.Name).
						Execute()
					g.Expect(err).ShouldNot(HaveOccurred())
					g.Expect(df).NotTo(BeNil())
				}).WithTimeout(20 * time.Minute).WithPolling(15 * time.Second).ShouldNot(HaveOccurred())
			})

			By("Adding a new DB and STORE", func() {
				df := &akov2.AtlasDataFederation{}
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: testNamespace.Name,
					Name:      dataFederationInstanceName,
				}, df)).To(Succeed())

				dfu := df.WithStorage(&akov2.Storage{
					Databases: []akov2.Database{
						{
							Name: "test-db-1",
							Collections: []akov2.Collection{
								{
									Name: "test-collection-1",
									DataSources: []akov2.DataSource{
										{
											StoreName: "http-test",
											Urls: []string{
												"https://data.cityofnewyork.us/api/views/vfnx-vebw/rows.csv",
											},
										},
									},
								},
							},
						},
					},
					Stores: []akov2.Store{
						{
							Name:     "http-test",
							Provider: "http",
						},
					},
				})
				Expect(k8sClient.Update(ctx, dfu)).To(Succeed())
			})

			By("Checking the DataFederation is ready", func() {
				df := &akov2.AtlasDataFederation{}
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Namespace: testNamespace.Name,
					Name:      dataFederationInstanceName,
				}, df)).To(Succeed())
				Eventually(func() bool {
					return resources.CheckCondition(k8sClient, df, api.TrueCondition(api.ReadyType))
				}).WithTimeout(2 * time.Minute).WithPolling(20 * time.Second).Should(BeTrue())
			})

			By("Deleting the DataFederation instance", func() {
				Expect(k8sClient.Delete(ctx, createdDataFederation)).To(Succeed())
				createdDataFederation = nil
			})
		})
	})
})

func deleteAtlasDataFederation(projectID, dataFederationName string) error {
	_, _, err := atlasClient.DataFederationApi.
		DeleteFederatedDatabase(ctx, projectID, dataFederationName).
		Execute()

	return err
}

func checkAtlasDataFederationRemoved(projectID, dataFederation string) func() bool {
	return func() bool {
		_, r, err := atlasClient.DataFederationApi.
			GetFederatedDatabase(ctx, projectID, dataFederation).
			Execute()
		if err != nil {
			if r != nil && r.StatusCode == http.StatusNotFound {
				return true
			}
		}

		return false
	}
}
