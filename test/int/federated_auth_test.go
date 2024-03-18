package int

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasFederatedAuth test", Label("AtlasFederatedAuth", "federated-auth"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var connectionSecret corev1.Secret

	var originalConnectedOrgConfig *admin.ConnectedOrgConfig
	var originalFederationSettings *admin.OrgFederationSettings
	var originalIdp *admin.FederationIdentityProvider

	resourceName := "fed-auth-test"
	ctx := context.Background()

	BeforeEach(func() {
		By("Checking if Federation Settings enabled for the org", func() {
			federationSettings, _, err := atlasClient.FederatedAuthenticationApi.GetFederationSettings(ctx, orgID).Execute()
			Expect(err).ToNot(HaveOccurred())
			Expect(federationSettings).ShouldNot(BeNil())

			originalFederationSettings = federationSettings
		})

		By("Getting original IDP", func() {
			identityProviders, _, err := atlasClient.FederatedAuthenticationApi.ListIdentityProviders(ctx, originalFederationSettings.GetId()).Execute()
			Expect(err).ToNot(HaveOccurred())

			for _, identityProvider := range identityProviders.GetResults() {
				idp := identityProvider
				if identityProvider.GetOktaIdpId() == originalFederationSettings.GetIdentityProviderId() {
					originalIdp = &idp
				}
			}

			Expect(originalIdp).ShouldNot(BeNil())
		})

		By("Getting existing org config", func() {
			connectedOrgConfig, _, err := atlasClient.FederatedAuthenticationApi.
				GetConnectedOrgConfig(ctx, originalFederationSettings.GetId(), orgID).
				Execute()
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
			atlasRoleMappings := originalConnectedOrgConfig.GetRoleMappings()
			roles := make([]akov2.RoleMapping, 0, len(atlasRoleMappings))
			for i := range atlasRoleMappings {
				atlasRole := atlasRoleMappings[i]
				newRole := akov2.RoleMapping{
					ExternalGroupName: atlasRole.ExternalGroupName,
					RoleAssignments:   []akov2.RoleAssignment{},
				}

				atlasRoleAssignments := atlasRole.GetRoleAssignments()
				for j := range atlasRoleAssignments {
					atlasRS := atlasRoleAssignments[j]
					project, _, err := atlasClient.ProjectsApi.GetProject(ctx, atlasRS.GetGroupId()).Execute()
					Expect(err).NotTo(HaveOccurred())
					Expect(project).NotTo(BeNil())

					newRS := akov2.RoleAssignment{
						ProjectName: project.Name,
						Role:        atlasRS.GetRole(),
					}
					newRole.RoleAssignments = append(newRole.RoleAssignments, newRS)
				}
				roles = append(roles, newRole)
			}
			roles = append(
				roles,
				akov2.RoleMapping{
					ExternalGroupName: "ako_team",
					RoleAssignments: []akov2.RoleAssignment{
						{Role: "ORG_OWNER"},
					},
				},
			)

			fedAuth := &akov2.AtlasFederatedAuth{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: testNamespace.Name,
				},
				Spec: akov2.AtlasFederatedAuthSpec{
					Enabled: true,
					ConnectionSecretRef: common.ResourceRefNamespaced{
						Name:      connectionSecret.Name,
						Namespace: connectionSecret.Namespace,
					},
					DomainAllowList:          append(originalConnectedOrgConfig.GetDomainAllowList(), "cloud-qa.mongodb.com", "mongodb.com"),
					DomainRestrictionEnabled: pointer.MakePtr(true),
					SSODebugEnabled:          pointer.MakePtr(false),
					PostAuthRoleGrants:       []string{"ORG_MEMBER"},
					RoleMappings:             roles,
				},
			}

			Expect(k8sClient.Create(ctx, fedAuth)).NotTo(HaveOccurred())
		})

		By("Federated Auth is ready", func() {
			Eventually(func(g Gomega) {
				fedAuth := &akov2.AtlasFederatedAuth{}
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())
				g.Expect(resources.CheckCondition(k8sClient, fedAuth, status.TrueCondition(status.ReadyType))).To(BeTrue())
			}).WithTimeout(10 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
		})

		By("Set initial config back", func() {
			fedAuth := &akov2.AtlasFederatedAuth{}
			Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())

			fedAuth.Spec.DomainAllowList = append(originalConnectedOrgConfig.GetDomainAllowList(), "mongodb.com")
			fedAuth.Spec.DomainRestrictionEnabled = &originalConnectedOrgConfig.DomainRestrictionEnabled
			fedAuth.Spec.SSODebugEnabled = originalIdp.SsoDebugEnabled
			fedAuth.Spec.PostAuthRoleGrants = originalConnectedOrgConfig.GetPostAuthRoleGrants()

			roleMappings := make([]akov2.RoleMapping, 0, len(fedAuth.Spec.RoleMappings))
			for _, roleMap := range fedAuth.Spec.RoleMappings {
				if roleMap.ExternalGroupName != "ako_team" {
					roleMappings = append(roleMappings, roleMap)
				}
			}
			fedAuth.Spec.RoleMappings = roleMappings

			Expect(k8sClient.Update(ctx, fedAuth)).NotTo(HaveOccurred())
		})

		By("Federated Auth is ready", func() {
			Eventually(func(g Gomega) {
				fedAuth := &akov2.AtlasFederatedAuth{}
				g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, fedAuth)).To(Succeed())
				g.Expect(resources.CheckCondition(k8sClient, fedAuth, status.TrueCondition(status.ReadyType))).To(BeTrue())
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
