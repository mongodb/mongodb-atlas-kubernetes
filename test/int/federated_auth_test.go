package int

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.mongodb.org/atlas/mongodbatlas"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/testutil"
)

var _ = Describe("AtlasFederatedAuth test", Label("AtlasFederatedAuth", "federated-auth"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var connectionSecret corev1.Secret

	var originalConnectedOrgConfig *mongodbatlas.FederatedSettingsConnectedOrganization
	var originalFederationSettings *mongodbatlas.FederatedSettings
	var originalIdp *mongodbatlas.FederatedSettingsIdentityProvider

	resourceName := "fed-auth-test"
	ctx := context.Background()

	BeforeEach(func() {
		By("Checking if Federation Settings enabled for the org", func() {
			federationSettings, _, err := atlasClient.FederatedSettings.Get(ctx, connection.OrgID)
			Expect(err).ToNot(HaveOccurred())
			Expect(federationSettings).ShouldNot(BeNil())

			originalFederationSettings = federationSettings
		})

		By("Getting original IDP", func() {
			idp, _, err := atlasClient.FederatedSettings.GetIdentityProvider(ctx, originalFederationSettings.ID, originalFederationSettings.IdentityProviderID)
			Expect(err).ToNot(HaveOccurred())
			Expect(idp).ShouldNot(BeNil())

			originalIdp = idp
		})

		By("Getting existing org config", func() {
			connectedOrgConfig, _, err := atlasClient.FederatedSettings.GetConnectedOrg(ctx, originalFederationSettings.ID, connection.OrgID)
			Expect(err).ToNot(HaveOccurred())
			Expect(connectedOrgConfig).ShouldNot(BeNil())

			originalConnectedOrgConfig = connectedOrgConfig
		})

		By("Starting the operator with protection OFF", func() {
			testNamespace, stopManager = prepareControllers(false)
			Expect(testNamespace).ShouldNot(BeNil())
			Expect(stopManager).ShouldNot(BeNil())
		})

		By("Creating project connection secret", func() {
			connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
			Expect(k8sClient.Create(ctx, &connectionSecret)).To(Succeed())
		})
	})

	It("Should be able to update existing Organization's federations settings", func() {
		By("Creating a FederatedAuthConfig resource", func() {
			roles := []mdbv1.RoleMapping{}

			for i := range originalConnectedOrgConfig.RoleMappings {
				atlasRole := *(originalConnectedOrgConfig.RoleMappings[i])
				newRole := mdbv1.RoleMapping{
					ExternalGroupName: atlasRole.ExternalGroupName,
					RoleAssignments:   []mdbv1.RoleAssignment{},
				}

				for j := range atlasRole.RoleAssignments {
					atlasRS := atlasRole.RoleAssignments[j]
					project, _, err := atlasClient.Projects.GetOneProject(ctx, atlasRS.GroupID)
					Expect(err).NotTo(HaveOccurred())
					Expect(project).NotTo(BeNil())

					newRS := mdbv1.RoleAssignment{
						ProjectName: project.Name,
						Role:        atlasRS.Role,
					}
					newRole.RoleAssignments = append(newRole.RoleAssignments, newRS)
				}
				roles = append(roles, newRole)
			}

			fedAuth := &mdbv1.AtlasFederatedAuth{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: testNamespace.Name,
				},
				Spec: mdbv1.AtlasFederatedAuthSpec{
					Enabled: true,
					ConnectionSecretRef: common.ResourceRefNamespaced{
						Name:      connectionSecret.Name,
						Namespace: connectionSecret.Namespace,
					},
					DomainAllowList:          append(originalConnectedOrgConfig.DomainAllowList, "cloud-qa.mongodb.com"),
					DomainRestrictionEnabled: toptr.MakePtr(true),
					SSODebugEnabled:          toptr.MakePtr(false),
					PostAuthRoleGrants:       []string{"ORG_MEMBER"},
					RoleMappings:             roles,
				},
			}

			Expect(k8sClient.Create(ctx, fedAuth)).NotTo(HaveOccurred())
		})

		By("Federated Auth is ready", func() {
			Eventually(func(g Gomega) {
				fedAuth := &mdbv1.AtlasFederatedAuth{}
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())
				g.Expect(testutil.CheckCondition(k8sClient, fedAuth, status.TrueCondition(status.ReadyType))).To(BeTrue())
			}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
		})

		By("Set initial config back", func() {
			fedAuth := &mdbv1.AtlasFederatedAuth{}
			Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())

			fedAuth.Spec.DomainAllowList = originalConnectedOrgConfig.DomainAllowList
			fedAuth.Spec.DomainRestrictionEnabled = originalConnectedOrgConfig.DomainRestrictionEnabled
			fedAuth.Spec.SSODebugEnabled = originalIdp.SsoDebugEnabled
			fedAuth.Spec.PostAuthRoleGrants = originalConnectedOrgConfig.PostAuthRoleGrants

			Expect(k8sClient.Update(ctx, fedAuth)).NotTo(HaveOccurred())
		})

		By("Federated Auth is ready", func() {
			Eventually(func(g Gomega) {
				fedAuth := &mdbv1.AtlasFederatedAuth{}
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())
				g.Expect(testutil.CheckCondition(k8sClient, fedAuth, status.TrueCondition(status.ReadyType))).To(BeTrue())
			}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Should delete connection secret", func() {
			Expect(k8sClient.Delete(ctx, &connectionSecret)).To(Succeed())
		})

		By("Should stop the operator", func() {
			stopManager()
			Expect(k8sClient.Delete(ctx, testNamespace)).ToNot(HaveOccurred())
		})
	})
})
