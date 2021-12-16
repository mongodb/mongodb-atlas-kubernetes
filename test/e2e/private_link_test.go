package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/deploy"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/aws"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

// NOTES
// Feature unavailable in Free and Shared-Tier Clusters
// This feature is not available for M0 free clusters, M2, and M5 clusters.

// tag for test resources "atlas-operator-test" (config.Tag)

var _ = Describe("[privatelink-aws] UserLogin", func() { // TODO probably it will be test-data-table later(?)
	It("User can deploy his resource", func() {
		var userData model.TestDataProvider
		By("Create user resources with No Private Link-------------------", func() {
			userData = model.NewTestDataProvider(
				"operator-private-link",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_backup.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30000,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			)
			actions.PrepareUsersConfigurations(&userData)
			deploy.NamespacedOperator(&userData)
			actions.DeployProject(&userData, "1")
		})
		By("Create Private Link and the rest users resources---------------------", func() {
			userData.Resources.Project.WithPrivateLink("AWS", "eu-west-2")
			actions.PrepareUsersConfigurations(&userData)
			actions.DeployProject(&userData, "2")
			// actions.DeployCluster(&userData, "1")
			// actions.DeployUsers(&userData)

			// TODO if implemented, check conditions statuses too + change Eventually >> Expect
			Eventually(AllPEndpointUpdated(&userData.Resources)).Should(BeTrue())
		})
		By("Create VPC, subnet, Endpoint. Sample", func() {
			session := aws.SessionAWS("eu-west-2")
			vpcID, err := session.GetVPCID()
			Expect(err).ShouldNot(HaveOccurred())
			subnetID, err := session.GetSubnetID() // subnetID := "subnet-01db67bf8e8c7a87b"
			Expect(err).ShouldNot(HaveOccurred())
			project := kube.GetProjectResource(userData.Resources.Namespace, userData.Resources.GetAtlasProjectFullKubeName())
			for i, peitem := range project.Status.PrivateEndpoints {
				serviceName := peitem.ServiceName
				Expect(serviceName).ShouldNot(BeEmpty())
				GinkgoWriter.Write([]byte("Subnet is available: " + subnetID))

				privateEndpointID, err := session.CreatePrivateEndpoint(vpcID, subnetID, serviceName, userData.Resources.Project.GetK8sMetaName())
				Expect(err).ShouldNot(HaveOccurred())
				getStatusPE := func(privateEndpointID string) string {
					r, err := session.DescribePrivateEndpointStatus(privateEndpointID)
					if err != nil {
						return ""
					}
					return r
				}
				Eventually(getStatusPE(privateEndpointID)).Should(Equal("pendingAcceptance"))

				By("Update PE ID from AWS", func() {
					userData.Resources.Project.UpdatePrivateLinkByOrder(i, privateEndpointID)
					actions.PrepareUsersConfigurations(&userData)
					actions.DeployProject(&userData, "3")
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
		// TODO test app block
	})
})

func AllPEndpointUpdated(input *model.UserInputs) func() bool {
	return func() bool {
		result := kube.GetProjectResource(input.Namespace, input.GetAtlasProjectFullKubeName())
		return len(result.Status.PrivateEndpoints) == len(result.Spec.PrivateEndpoints)
	}
}
