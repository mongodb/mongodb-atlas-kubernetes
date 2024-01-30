package e2e_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	AzureClientID     = "AZURE_CLIENT_ID"
	KeyVaultName      = "ako-kms-test"
	AzureClientSecret = "AZURE_CLIENT_SECRET" //#nosec G101 -- False positive; this is the env var, not the secret itself
	AzureEnvironment  = "AZURE"
	KeyName           = "encryption-at-rest-test-key"
)

var _ = Describe("Encryption at REST test", Label("encryption-at-rest"), func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		checkUpAWSEnvironment()
		checkUpAzureEnvironment()
		checkNSetUpGCPEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Cloud", func() {
			DeleteAllRoles(testData)
		})

		By("Delete Resources, Project with Encryption at rest", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("Encryption at rest for AWS, GCP, Azure",
		func(test *model.TestDataProvider, encAtRest v1.EncryptionAtRest, roles []cloudaccess.Role) {
			testData = test
			actions.ProjectCreationFlow(test)

			if roles != nil {
				cloudAccessRolesFlow(test, roles)
			}

			encryptionAtRestFlow(test, encAtRest)
		},
		Entry("Test[encryption-at-rest-aws]: Can add Encryption at Rest to AWS project", Label("encryption-at-rest-aws"),
			model.DataProvider(
				"encryption-at-rest-aws",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			v1.EncryptionAtRest{
				AwsKms: v1.AwsKms{
					Enabled: toptr.MakePtr(true),
					Region:  "US_EAST_1",
				},
			},
			[]cloudaccess.Role{
				{
					Name: utils.RandomName(awsRoleNameBase),
					AccessRole: v1.CloudProviderIntegration{
						ProviderName: "AWS",
					},
				},
			},
		),
		Entry("Test[encryption-at-rest-azure]: Can add Encryption at Rest to Azure project", Label("encryption-at-rest-azure"),
			model.DataProvider(
				"encryption-at-rest-azure",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			v1.EncryptionAtRest{
				AzureKeyVault: v1.AzureKeyVault{
					AzureEnvironment:  AzureEnvironment,
					ClientID:          os.Getenv(AzureClientID),
					Enabled:           toptr.MakePtr(true),
					ResourceGroupName: cloud.ResourceGroupName,
					TenantID:          os.Getenv(DirectoryID),
				},
			},
			nil,
		),
		Entry("Test[encryption-at-rest-gcp]: Can add Encryption at Rest to GCP project", Label("encryption-at-rest-gcp"),
			model.DataProvider(
				"encryption-at-rest-gcp",
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				40000,
				[]func(*model.TestDataProvider){},
			).WithProject(data.DefaultProject()),
			v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled: toptr.MakePtr(true),
				},
			},
			nil,
		),
	)
})

func encryptionAtRestFlow(userData *model.TestDataProvider, encAtRest v1.EncryptionAtRest) {
	By("Create KMS", func() {
		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

		var aRole status.CloudProviderIntegration
		if len(userData.Project.Status.CloudProviderIntegrations) > 0 {
			aRole = userData.Project.Status.CloudProviderIntegrations[0]
		}

		fillKMSforAWS(userData, &encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)
		fillVaultforAzure(userData, &encAtRest)
		fillKMSforGCP(userData, &encAtRest)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.EncryptionAtRest = &encAtRest
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
		actions.WaitForConditionsToBecomeTrue(userData, status.EncryptionAtRestReadyType, status.ReadyType)
	})

	By("Remove Encryption at Rest from the project", func() {
		removeAllEncryptionsSeparately(&encAtRest)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
		userData.Project.Spec.EncryptionAtRest = &encAtRest
		Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
	})

	By("Check if project returned back to the initial state", func() {
		actions.CheckProjectConditionsNotSet(userData, status.EncryptionAtRestReadyType)

		Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
			Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

		Eventually(func(g Gomega) bool {
			areEmpty, err := checkIfEncryptionsAreDisabled(userData.Project.ID())
			g.Expect(err).ShouldNot(HaveOccurred())
			return areEmpty
		}).WithTimeout(5*time.Minute).WithPolling(20*time.Second).
			Should(BeTrue(), "Encryption at Rest is not disabled")
	})
}

func fillKMSforAWS(userData *model.TestDataProvider, encAtRest *v1.EncryptionAtRest, atlasAccountArn, assumedRoleArn string) {
	if (encAtRest.AwsKms == v1.AwsKms{}) {
		return
	}

	alias := fmt.Sprintf("%s-kms", userData.Project.Spec.Name)

	Expect(encAtRest.AwsKms.Region).NotTo(Equal(""))
	awsAction, err := cloud.NewAWSAction(GinkgoT())
	Expect(err).ToNot(HaveOccurred())
	CustomerMasterKeyID, err := awsAction.CreateKMS(alias, config.AWSRegionUS, atlasAccountArn, assumedRoleArn)
	Expect(err).ToNot(HaveOccurred())
	Expect(CustomerMasterKeyID).NotTo(Equal(""))

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-secret",
			Namespace: userData.Resources.Namespace,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"CustomerMasterKeyID": []byte(CustomerMasterKeyID),
			"RoleID":              []byte(assumedRoleArn),
		},
	}

	Expect(userData.K8SClient.Create(context.Background(), secret)).To(Succeed())

	encAtRest.AwsKms.SecretRef = common.ResourceRefNamespaced{
		Name:      "aws-secret",
		Namespace: userData.Resources.Namespace,
	}
}

func fillVaultforAzure(userData *model.TestDataProvider, encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.AzureKeyVault == v1.AzureKeyVault{}) {
		return
	}

	azAction, err := cloud.NewAzureAction(GinkgoT(), os.Getenv(SubscriptionID), cloud.ResourceGroupName)
	Expect(err).ToNot(HaveOccurred())

	keyID, err := azAction.CreateKeyVault(KeyName)
	Expect(err).ToNot(HaveOccurred())

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "az-secret",
			Namespace: userData.Resources.Namespace,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"KeyIdentifier":  []byte(keyID),
			"KeyVaultName":   []byte(KeyVaultName),
			"Secret":         []byte(os.Getenv(AzureClientSecret)),
			"SubscriptionID": []byte(os.Getenv(SubscriptionID)),
		},
	}
	Expect(userData.K8SClient.Create(context.Background(), secret)).To(Succeed())

	encAtRest.AzureKeyVault.SecretRef = common.ResourceRefNamespaced{
		Name:      "az-secret",
		Namespace: userData.Resources.Namespace,
	}
}

func fillKMSforGCP(userData *model.TestDataProvider, encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.GoogleCloudKms == v1.GoogleCloudKms{}) {
		return
	}

	gcpAction, err := cloud.NewGCPAction(GinkgoT(), cloud.GoogleProjectID)
	Expect(err).ToNot(HaveOccurred())

	keyID, err := gcpAction.CreateKMS()
	Expect(err).ToNot(HaveOccurred())

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gcp-secret",
			Namespace: userData.Resources.Namespace,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"ServiceAccountKey":    []byte(os.Getenv("GCP_SA_CRED")),
			"KeyVersionResourceID": []byte(keyID),
		},
	}
	Expect(userData.K8SClient.Create(context.Background(), secret)).To(Succeed())

	encAtRest.GoogleCloudKms.SecretRef = common.ResourceRefNamespaced{
		Name:      "gcp-secret",
		Namespace: userData.Resources.Namespace,
	}
}

func removeAllEncryptionsSeparately(encAtRest *v1.EncryptionAtRest) {
	encAtRest.AwsKms = v1.AwsKms{}
	encAtRest.AzureKeyVault = v1.AzureKeyVault{}
	encAtRest.GoogleCloudKms = v1.GoogleCloudKms{}
}

func checkIfEncryptionsAreDisabled(projectID string) (areEmpty bool, err error) {
	atlasClient := atlas.GetClientOrFail()
	encryptionAtRest, err := atlasClient.GetEncryptionAtRest(projectID)
	if err != nil {
		return false, err
	}

	if encryptionAtRest == nil {
		return true, nil
	}

	awsEnabled := *encryptionAtRest.AwsKms.Enabled
	azureEnabled := *encryptionAtRest.AzureKeyVault.Enabled
	gcpEnabled := *encryptionAtRest.GoogleCloudKms.Enabled

	if awsEnabled || azureEnabled || gcpEnabled {
		return false, nil
	}

	return true, nil
}

var _ = Describe("Encryption at rest AWS", Label("encryption-at-rest", "encryption-at-rest-aws"), Ordered, func() {
	var testData *model.TestDataProvider

	_ = BeforeEach(func() {
		checkUpEnvironment()
		checkUpAWSEnvironment()
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}
		By("Clean Roles", func() {
			DeleteAllRoles(testData)
		})
		By("Delete Resources, Project with Cloud provider access roles", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	It("Should be able to create Encryption at REST on AWS with RoleID equal to AWS ARN", func() {

		testData = model.DataProvider(
			"encryption-at-rest-aws",
			model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
			40000,
			[]func(*model.TestDataProvider){},
		).WithProject(data.DefaultProject())

		roles := []cloudaccess.Role{
			{
				Name: utils.RandomName(awsRoleNameBase),
				AccessRole: v1.CloudProviderIntegration{
					ProviderName: "AWS",
				},
			},
		}
		userData := testData
		encAtRest := v1.EncryptionAtRest{
			AwsKms: v1.AwsKms{
				Enabled: toptr.MakePtr(true),
				Region:  "US_EAST_1",
			},
		}

		By("Creating a project", func() {
			actions.ProjectCreationFlow(testData)
		})

		var projectID string
		By("Getting a project ID by name from Atlas", func() {
			Eventually(func(g Gomega) error {
				projectData, _, err := atlasClient.Client.ProjectsApi.
					GetProjectByName(userData.Context, userData.Project.Spec.Name).
					Execute()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(projectData).NotTo(BeNil())
				ginkgo.GinkgoLogr.Info("Project ID", projectData.GetId())
				projectID = projectData.GetId()
				return nil
			}).WithTimeout(2 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		var atlasRoles *admin.CloudProviderAccessRoles
		By("Add cloud access role (AWS only)", func() {
			cloudAccessRolesFlow(userData, roles)
		})

		By("Fetching project CPAs", func() {
			var err error
			atlasRoles, _, err = atlasClient.Client.CloudProviderAccessApi.
				ListCloudProviderAccessRoles(userData.Context, projectID).
				Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(atlasRoles).NotTo(BeNil())
			Expect(len(atlasRoles.GetAwsIamRoles())).NotTo(BeZero())
		})

		By("Create KMS with AWS RoleID", func() {
			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

			Expect(len(userData.Project.Status.CloudProviderIntegrations)).NotTo(Equal(0))
			aRole := userData.Project.Status.CloudProviderIntegrations[0]

			fillKMSforAWS(userData, &encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)

			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
			userData.Project.Spec.EncryptionAtRest = &encAtRest

			Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
			actions.WaitForConditionsToBecomeTrue(userData, status.EncryptionAtRestReadyType, status.ReadyType)
		})
	})
})
