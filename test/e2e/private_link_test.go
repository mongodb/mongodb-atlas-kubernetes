package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// NOTES
// Feature unavailable in Free and Shared-Tier Clusters
// This feature is not available for M0 free clusters, M2, and M5 clusters.

// tag for test resources "atlas-operator-test" (config.Tag)

type privateEndpoint struct {
	provider string
	region   string
}

var _ = Describe("[privatelink-aws] UserLogin", func() {
	var data model.TestDataProvider

	_ = BeforeEach(func() {
		Eventually(kubecli.GetVersionOutput()).Should(Say(K8sVersion))
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + data.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentGinkgoTestDescription().Failed {
			GinkgoWriter.Write([]byte("Test has been failed. Trying to save logs...\n"))
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operatorDecribe.txt", data.Resources.Namespace),
				[]byte(kubecli.DescribeOperatorPod(data.Resources.Namespace)),
			)
			utils.SaveToFile(
				fmt.Sprintf("output/%s/operator-logs.txt", data.Resources.Namespace),
				kubecli.GetManagerLogs(data.Resources.Namespace),
			)
			actions.SaveTestAppLogs(data.Resources)
			actions.SaveK8sResources(
				[]string{"deploy", "atlasprojects"},
				data.Resources.Namespace,
			)
		} else {
			actions.AfterEachFinalCleanup([]model.TestDataProvider{data})
		}
	})

	DescribeTable("Namespaced operators working only with its own namespace with different configuration",
		func(test model.TestDataProvider, pe []privateEndpoint) {
			data = test
			privateFlow(test, pe)
		},
		Entry("Test: User has project which was updated with AWS PrivateEndpoint",
			model.NewTestDataProvider(
				"operator-plink-aws-1",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				40000,
				[]func(*model.TestDataProvider){},
			),
			[]privateEndpoint{
				{
					provider: "AWS",
					region:   "eu-west-2",
				},
			},
		),
		// Entry("Test: User has project which was updated with Azure PrivateEndpoint",
		// 	model.NewTestDataProvider(
		// 		"operator-plink-azure-1",
		// 		model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
		// 		[]string{"data/atlascluster_backup.yaml"},
		// 		[]string{},
		// 		[]model.DBUser{
		// 			*model.NewDBUser("user1").
		// 				WithSecretRef("dbuser-secret-u1").
		// 				AddBuildInAdminRole(),
		// 		},
		// 		40000,
		// 		[]func(*model.TestDataProvider){
		// 			actions.DeleteFirstUser,
		// 		},
		// 	),
		// ),
	)
})

func privateFlow(userData model.TestDataProvider, requstedPE []privateEndpoint) {
	By("Deploy Project with requested configuration", func() {
		actions.PrepareUsersConfigurations(&userData)
		deploy.NamespacedOperator(&userData)
		actions.DeployProject(&userData, "1")
	})
	By("Create Private Link and the rest users resources", func() {
		for _, pe := range requstedPE {
			userData.Resources.Project.WithPrivateLink(provider.ProviderName(pe.provider), pe.region)
		}
		actions.PrepareUsersConfigurations(&userData)
		actions.DeployProject(&userData, "2")
		// actions.DeployCluster(&userData, "1")
		// actions.DeployUsers(&userData)

		Eventually(kube.GetProjectPEndpointServiceStatus(&userData)).Should(Equal("True"), "Atlasproject status.conditions are not True")
		Expect(AllPEndpointUpdated(&userData)).Should(BeTrue(), "Error: Was created a different amount of endpoints")
	})
	By("Create Endpoint in requested Cloud Provider", func() {
		session := aws.SessionAWS("eu-west-2")
		vpcID, err := session.GetVPCID()
		Expect(err).ShouldNot(HaveOccurred())
		subnetID, err := session.GetSubnetID()
		Expect(err).ShouldNot(HaveOccurred())
		project, err := kube.GetProjectResource(&userData)
		Expect(err).ShouldNot(HaveOccurred())

		for i, peitem := range project.Status.PrivateEndpoints {
			serviceName := peitem.ServiceName
			Expect(serviceName).ShouldNot(BeEmpty())
			GinkgoWriter.Write([]byte("Subnet is available: " + subnetID))

			privateEndpointID, err := session.CreatePrivateEndpoint(vpcID, subnetID, serviceName, userData.Resources.Project.GetK8sMetaName())
			Expect(err).ShouldNot(HaveOccurred())
			getStatusPE := func(privateEndpointID string) func() string {
				return func() string {
					r, err := session.DescribePrivateEndpointStatus(privateEndpointID)
					if err != nil {
						return ""
					}
					return r
				}
			}
			Eventually(
				kube.GetProjectPEndpointServiceStatus(&userData),
			).Should(Equal("True"))
			Eventually(getStatusPE(privateEndpointID)).Should(Equal("pendingAcceptance"))

			By("Update PE ID from AWS", func() {
				userData.Resources.Project.UpdatePrivateLinkByOrder(i, privateEndpointID)
				actions.PrepareUsersConfigurations(&userData)
				actions.DeployProject(&userData, "3")
				Eventually(kube.GetReadyProjectStatus(&userData)).Should(Equal("True"), "Condition status is not 'True'")
				Eventually(
					kube.GetProjectPEndpointStatus(&userData),
				).Should(Equal("True"))
				Eventually(getStatusPE(privateEndpointID)).Should(Equal("available"))
			})

			By("Delete PE from AWS", func() {
				err = session.DeletePrivateLink(privateEndpointID)
				Expect(err).ShouldNot(HaveOccurred())

				status, err := session.DescribePrivateEndpointStatus(privateEndpointID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).Should(Equal("rejected"))
			})
		}
	})
	By("Delete Resources, Project with PEService", func() {
		// actions.DeleteUserResources(&userData)
		actions.DeleteUserResourcesProject(&userData)
	})
}

func AllPEndpointUpdated(data *model.TestDataProvider) bool {
	result, _ := kube.GetProjectResource(data)
	return len(result.Status.PrivateEndpoints) == len(result.Spec.PrivateEndpoints)
}
