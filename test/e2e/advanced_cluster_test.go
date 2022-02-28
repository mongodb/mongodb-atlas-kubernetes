package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/helm"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

var _ = Describe("HELM advanced cluster.", Label("helm-advanced-cluster"), func() {
	var data model.TestDataProvider

	It("User can deploy operator namespaces by using HELM", func() {
		By("User creates configuration for a new Project and Advanced Cluster", func() {
			data = model.NewTestDataProvider(
				"helm-wide",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_advanced_helm.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("reader2").
						WithSecretRef("dbuser-secret-u2").
						AddCustomRole(model.RoleCustomReadWrite, "Ships", "").
						WithAuthDatabase("admin"),
				},
				30014,
				[]func(*model.TestDataProvider){},
			)
			// helm template has equal ObjectMeta.Name and Spec.Name
			data.Resources.Clusters[0].ObjectMeta.Name = "advanced-cluster-helm"
			data.Resources.Clusters[0].Spec.AdvancedClusterSpec.Name = "advanced-cluster-helm"
		})
		By("User use helm for deploying operator", func() {
			helm.InstallOperatorWideSubmodule(data.Resources)
		})
		By("User deploy cluster by helm", func() {
			helm.InstallClusterSubmodule(data.Resources)
		})
		waitClusterWithChecks(&data)
		deleteClusterAndOperator(&data)
	})
})
