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

const (
	vpcID = "vpc-097de88bae21b8f2e"
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

			// Eventually(IsStatusPEUpdated(&userData.Resources)).Should(BeTrue())
			// project := kube.GetProjectResource(userData.Resources.GetAtlasProjectFullKubeName(), userData.Resources.Namespace)
			// for _, pelist := range project.Status.PrivateEndpoints {
			// 	for _, peitem := range pelist.Endpoints {
			// 		Expect(peitem.EndpointID).ShouldNot(BeEmpty())
			// 	}
			// } // TODO if implemented, get status instead

		})
		By("Request Private Link by passing region and provider", func() {

		})
		// TODO get Service Name of PE
		serviceName := "com.amazonaws.vpce.eu-west-2.vpce-svc-0ce0c4e9a5d1f6472"
		// TODO create private link on AWS side
		By("Create VPC, subnet, Endpoint. Sample", func() {

			// TODO this is for sample. will be removed later
			session := aws.SessionAWS("eu-west-2")
			testID := "id-test-ksjs03jk"

			subnetID, err := session.CreateSubnet(vpcID, "10.0.0.0/28", testID)
			Expect(err).ShouldNot(HaveOccurred())
			getStatusSubnet := func(subnetID string) string {
				r, err := session.DescribeSubnetStatus(subnetID)
				if err != nil {
					return ""
				}
				return r
			}
			Eventually(getStatusSubnet(subnetID)).Should(Equal("available"))
			GinkgoWriter.Write([]byte("Subnet is available: " + subnetID))

			privateEndpointID, err := session.CreatePrivateEndpoint(vpcID, subnetID, serviceName, testID)
			Expect(err).ShouldNot(HaveOccurred())
			getStatusPE := func(privateEndpointID string) string {
				r, err := session.DescribePrivateEndpointStatus(privateEndpointID)
				if err != nil {
					return ""
				}
				return r
			}
			Eventually(getStatusPE(privateEndpointID)).Should(Equal("pendingAcceptance"))

			err = session.DeletePrivateLink(privateEndpointID)
			Expect(err).ShouldNot(HaveOccurred())
		})
		// TODO attach private link id to atlas
		// TODO check atlas status
		// TODO test app block
	})
})

func IsStatusPEUpdated(input *model.UserInputs) func() bool {
	return func() bool {
		result := kube.GetProjectResource(input.Namespace, input.GetAtlasProjectFullKubeName())
		return len(result.Status.PrivateEndpoints) == len(result.Spec.PrivateEndpoints)
	}
}
