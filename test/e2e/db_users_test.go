package e2e

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Operator watch all namespace should create connection secrets for database users in any namaspace", Label("deployment-ns"), func() {
	var testData *model.TestDataProvider
	secondNamespace := "second-namespace"

	_ = AfterEach(func() {
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveDeploymentsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
			Expect(actions.SaveUsersToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		actions.DeleteTestDataUsers(testData)
		actions.DeleteTestDataDeployments(testData)
		actions.DeleteTestDataProject(testData)
		actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		Expect(k8s.DeleteNamespace(testData.Context, testData.K8SClient, secondNamespace)).Should(Succeed())
	})

	It("Operator run on global namespace", func() {
		By("Setting up test data", func() {
			project := data.DefaultProject()
			project.Namespace = config.DefaultOperatorNS

			deployment := data.CreateBasicDeployment("dbusers-operator-global")
			deployment.Namespace = config.DefaultOperatorNS

			testData = model.DataProvider(
				"dbusers-operator-global",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			).WithProject(project).
				WithInitialDeployments(deployment).
				WithUsers(
					data.BasicUser(
						"reader1",
						"reader1",
						data.WithSecretRef("dbuser-secret-u1"),
						data.WithReadWriteRole(),
						data.WithNamespace(config.DefaultOperatorNS),
					),
					data.BasicUser(
						"reader2",
						"reader2",
						data.WithProject(project),
						data.WithSecretRef("dbuser-secret-u2"),
						data.WithReadWriteRole(),
						data.WithNamespace(secondNamespace),
					),
				)
			testData.Resources.Namespace = config.DefaultOperatorNS
		})
		By("Running operator watching global namespace", func() {
			k8s.CreateNamespace(testData.Context, testData.K8SClient, config.DefaultOperatorNS)
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, config.DefaultOperatorNS)
			k8s.CreateNamespace(testData.Context, testData.K8SClient, secondNamespace)
			logPath := path.Join("output", fmt.Sprintf("dbusers-operator-global-%s", testData.Resources.Namespace))

			mgr, err := k8s.RunOperator(&k8s.Config{
				GlobalAPISecret: client.ObjectKey{
					Namespace: config.DefaultOperatorNS,
					Name:      config.DefaultOperatorGlobalKey,
				},
				LogDir: logPath,
			})
			Expect(err).NotTo(HaveOccurred())

			go func(ctx context.Context) context.Context {
				err := mgr.Start(ctx)
				Expect(err).NotTo(HaveOccurred())
				return ctx
			}(testData.Context)
		})
		By("Creating project and database users resources", func() {
			deploy.CreateProject(testData)
			deploy.CreateUsers(testData)

			Eventually(func(g Gomega) bool {
				for i := range testData.Users {
					dbUser := testData.Users[i]

					g.Expect(testData.K8SClient.Get(testData.Context, client.ObjectKeyFromObject(dbUser), dbUser)).To(Succeed())
					g.Expect(testutil.CheckCondition(testData.K8SClient, dbUser, status.TrueCondition(status.ReadyType))).To(BeTrue())
				}

				return true
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

			Expect(countConnectionSecrets(testData.K8SClient, testData.Project.Spec.Name)).To(Equal(0))
		})
		By("Should create deployment and connection secrets for all users", func() {
			Expect(testData.K8SClient.Create(testData.Context, testData.InitialDeployments[0])).To(Succeed())

			Eventually(func(g Gomega) bool {
				g.Expect(testData.K8SClient.Get(testData.Context, types.NamespacedName{
					Name:      testData.InitialDeployments[0].Name,
					Namespace: testData.InitialDeployments[0].Namespace,
				}, testData.InitialDeployments[0])).To(Succeed())

				return testData.InitialDeployments[0].Status.StateName == status.StateIDLE
			}).WithTimeout(30 * time.Minute).Should(BeTrue())

			Eventually(func(g Gomega) bool {
				return g.Expect(countConnectionSecrets(testData.K8SClient, testData.Project.Spec.Name)).To(Equal(2))
			}).WithTimeout(5 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())
		})
	})
})

func countConnectionSecrets(k8sClient client.Client, projectName string) int {
	secretList := corev1.SecretList{}
	Expect(k8sClient.List(context.Background(), &secretList)).To(Succeed())

	names := make([]string, 0)
	for _, item := range secretList.Items {
		fmt.Println(item.Name, projectName, kube.NormalizeIdentifier(projectName))
		if strings.HasPrefix(item.Name, kube.NormalizeIdentifier(projectName)) {
			names = append(names, item.Name)
		}
	}

	return len(names)
}
