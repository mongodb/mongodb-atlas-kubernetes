package e2e_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas/mongodbatlas"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/testutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloud"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/cloudaccess"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/api/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
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
				encAtRest.AwsKms.RoleID = test.Project.Status.CloudProviderAccessRoles[0].IamAssumedRoleArn
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
					AccessRole: v1.CloudProviderAccessRole{
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
					KeyVaultName:      KeyVaultName,
					ResourceGroupName: cloud.ResourceGroupName,
					Secret:            os.Getenv(AzureClientSecret),
					TenantID:          os.Getenv(DirectoryID),
					SubscriptionID:    os.Getenv(SubscriptionID),
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
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: os.Getenv("GCP_SA_CRED"),
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

		var aRole status.CloudProviderAccessRole
		if len(userData.Project.Status.CloudProviderAccessRoles) > 0 {
			aRole = userData.Project.Status.CloudProviderAccessRoles[0]
		}

		fillKMSforAWS(&encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)
		fillVaultforAzure(&encAtRest)
		fillKMSforGCP(&encAtRest)

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

func fillKMSforAWS(encAtRest *v1.EncryptionAtRest, atlasAccountArn, assumedRoleArn string) {
	if (encAtRest.AwsKms == v1.AwsKms{}) {
		return
	}

	Expect(encAtRest.AwsKms.Region).NotTo(Equal(""))
	awsAction, err := cloud.NewAWSAction(GinkgoT())
	Expect(err).ToNot(HaveOccurred())
	CustomerMasterKeyID, err := awsAction.CreateKMS(config.AWSRegionUS, atlasAccountArn, assumedRoleArn)
	Expect(err).ToNot(HaveOccurred())
	Expect(CustomerMasterKeyID).NotTo(Equal(""))

	encAtRest.AwsKms.CustomerMasterKeyID = CustomerMasterKeyID
}

func fillVaultforAzure(encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.AzureKeyVault == v1.AzureKeyVault{}) {
		return
	}

	azAction, err := cloud.NewAzureAction(GinkgoT(), os.Getenv(SubscriptionID), cloud.ResourceGroupName)
	Expect(err).ToNot(HaveOccurred())

	keyID, err := azAction.CreateKeyVault(KeyName)
	Expect(err).ToNot(HaveOccurred())

	encAtRest.AzureKeyVault.KeyIdentifier = keyID
}

func fillKMSforGCP(encAtRest *v1.EncryptionAtRest) {
	if (encAtRest.GoogleCloudKms == v1.GoogleCloudKms{}) {
		return
	}

	gcpAction, err := cloud.NewGCPAction(GinkgoT(), cloud.GoogleProjectID)
	Expect(err).ToNot(HaveOccurred())

	keyID, err := gcpAction.CreateKMS()
	Expect(err).ToNot(HaveOccurred())

	encAtRest.GoogleCloudKms.KeyVersionResourceID = keyID
}

func removeAllEncryptionsSeparately(encAtRest *v1.EncryptionAtRest) {
	encAtRest.AwsKms = v1.AwsKms{}
	encAtRest.AzureKeyVault = v1.AzureKeyVault{}
	encAtRest.GoogleCloudKms = v1.GoogleCloudKms{}
}

func checkIfEncryptionsAreDisabled(projectID string) (areEmpty bool, err error) {
	atlasClient := atlas.GetClientOrFail()
	encryptionAtRest, err := atlasClient.GetEncryptioAtRest(projectID)
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

var _ = Describe("Encryption at rest AWS", Label("encryption-at-rest"), Ordered, func() {
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
				AccessRole: v1.CloudProviderAccessRole{
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
				projectData, _, err := atlasClient.Client.Projects.GetOneProjectByName(userData.Context, userData.Project.Spec.Name)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(projectData).NotTo(BeNil())
				ginkgo.GinkgoLogr.Info("Project ID", projectData.ID)
				projectID = projectData.ID
				return nil
			}).WithTimeout(2 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		var atlasRoles *mongodbatlas.CloudProviderAccessRoles
		By("Add cloud access role (AWS only)", func() {
			cloudAccessRolesFlow(userData, roles)
		})

		By("Fetching project CPAs", func() {
			var err error
			atlasRoles, _, err = atlasClient.Client.CloudProviderAccess.ListRoles(userData.Context, projectID)
			Expect(err).NotTo(HaveOccurred())
			Expect(atlasRoles).NotTo(BeNil())
			Expect(len(atlasRoles.AWSIAMRoles)).NotTo(BeZero())
		})

		By("Create KMS with AWS RoleID", func() {
			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

			Expect(len(userData.Project.Status.CloudProviderAccessRoles)).NotTo(Equal(0))
			aRole := userData.Project.Status.CloudProviderAccessRoles[0]

			fillKMSforAWS(&encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)
			fillVaultforAzure(&encAtRest)
			fillKMSforGCP(&encAtRest)

			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
			userData.Project.Spec.EncryptionAtRest = &encAtRest

			var roleARNToSet string
			for _, r := range atlasRoles.AWSIAMRoles {
				if r.IAMAssumedRoleARN == aRole.IamAssumedRoleArn {
					GinkgoWriter.Println("FOUND ROLE ID >>>> ", r.IAMAssumedRoleARN, r.RoleID)
					roleARNToSet = r.IAMAssumedRoleARN
					break
				}
			}
			GinkgoWriter.Println(" NO ROLE ID FOUND >>>> ")

			Expect(roleARNToSet).NotTo(BeEmpty())
			userData.Project.Spec.EncryptionAtRest.AwsKms.RoleID = roleARNToSet
			Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
			actions.WaitForConditionsToBecomeTrue(userData, status.EncryptionAtRestReadyType, status.ReadyType)
		})
	})

	It("Should be able to create Encryption at REST on AWS with data from the Secret", func() {

		testData = model.DataProvider(
			"encryption-at-rest-aws",
			model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
			40000,
			[]func(*model.TestDataProvider){},
		).WithProject(data.DefaultProject())

		roles := []cloudaccess.Role{
			{
				Name: utils.RandomName(awsRoleNameBase),
				AccessRole: v1.CloudProviderAccessRole{
					ProviderName: "AWS",
				},
			},
		}

		userData := testData
		encAtRest := v1.EncryptionAtRest{
			AwsKms: v1.AwsKms{
				Enabled: toptr.MakePtr(true),
			},
		}

		By("Creating a project", func() {
			actions.ProjectCreationFlow(testData)
		})

		var projectID string
		By("Getting a project ID by name from Atlas", func() {
			Eventually(func(g Gomega) error {
				projectData, _, err := atlasClient.Client.Projects.GetOneProjectByName(userData.Context, userData.Project.Spec.Name)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(projectData).NotTo(BeNil())
				ginkgo.GinkgoLogr.Info("Project ID", projectData.ID)
				projectID = projectData.ID
				return nil
			}).WithTimeout(2 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		var atlasRoles *mongodbatlas.CloudProviderAccessRoles
		By("Add cloud access role (AWS only)", func() {
			cloudAccessRolesFlow(userData, roles)
		})

		By("Fetching project CPAs", func() {
			var err error
			atlasRoles, _, err = atlasClient.Client.CloudProviderAccess.ListRoles(userData.Context, projectID)
			Expect(err).NotTo(HaveOccurred())
			Expect(atlasRoles).NotTo(BeNil())
			Expect(len(atlasRoles.AWSIAMRoles)).NotTo(BeZero())
		})

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
				"Region": []byte("US_EAST_1"),
			},
		}

		By("Create KMS with AWS RoleID", func() {
			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())

			Expect(len(userData.Project.Status.CloudProviderAccessRoles)).NotTo(Equal(0))
			aRole := userData.Project.Status.CloudProviderAccessRoles[0]

			encAtRest.AwsKms.Region = string(secret.Data["Region"])

			fillKMSforAWS(&encAtRest, aRole.AtlasAWSAccountArn, aRole.IamAssumedRoleArn)

			Expect(userData.K8SClient.Get(userData.Context, types.NamespacedName{Name: userData.Project.Name,
				Namespace: userData.Resources.Namespace}, userData.Project)).Should(Succeed())
			userData.Project.Spec.EncryptionAtRest = &encAtRest

			var roleARNToSet string
			for _, r := range atlasRoles.AWSIAMRoles {
				if r.IAMAssumedRoleARN == aRole.IamAssumedRoleArn {
					GinkgoWriter.Println("FOUND ROLE ID >>>> ", r.IAMAssumedRoleARN, r.RoleID)
					roleARNToSet = r.IAMAssumedRoleARN
					break
				}
			}
			GinkgoWriter.Println(" NO ROLE ID FOUND >>>> ")

			Expect(roleARNToSet).NotTo(BeEmpty())

			secret.Data["RoleID"] = []byte(roleARNToSet)
			secret.Data["CustomerMasterKeyID"] = []byte(encAtRest.AwsKms.CustomerMasterKeyID)
			userData.Project.Spec.EncryptionAtRest.AwsKms.CustomerMasterKeyID = ""
			userData.Project.Spec.EncryptionAtRest.AwsKms.RoleID = roleARNToSet
			userData.Project.Spec.EncryptionAtRest.AwsKms.SecretRef = common.ResourceRefNamespaced{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			}

			By("Creating a secret for AWS KMS", func() {
				Expect(userData.K8SClient.Create(userData.Context, secret)).To(Succeed())
			})

			Expect(userData.K8SClient.Update(userData.Context, userData.Project)).Should(Succeed())
			actions.WaitForConditionsToBecomeTrue(userData, status.EncryptionAtRestReadyType, status.ReadyType)
		})
	})
})

func configureManager(testData *model.TestDataProvider) {
	mgr := actions.PrepareOperatorConfigurations(testData)
	ctx := context.Background()
	go func(ctx context.Context) context.Context {
		err := mgr.Start(ctx)
		Expect(err).NotTo(HaveOccurred())
		return ctx
	}(ctx)
	testData.ManagerContext = ctx
}

func createProjectWithValidationError(testData *model.TestDataProvider, errMsg string) {
	if testData.Project.GetNamespace() == "" {
		testData.Project.Namespace = testData.Resources.Namespace
	}
	By(fmt.Sprintf("Deploy Broken Project %s", testData.Project.GetName()), func() {
		err := testData.K8SClient.Create(testData.Context, testData.Project)
		Expect(err).ShouldNot(HaveOccurred(), "Project %s was not created", testData.Project.GetName())
		expectedCondition :=
			status.FalseCondition(status.ValidationSucceeded).WithReason(string(workflow.Internal)).WithMessageRegexp(errMsg)
		Eventually(func() bool {
			return testutil.CheckCondition(testData.K8SClient, testData.Project, expectedCondition)
		}).WithPolling(3 * time.Second).WithTimeout(40 * time.Second).Should(BeTrue())
	})
}

func withProperUrls(properties string) string {
	urls := `"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/619108922856-compute%40developer.gserviceaccount.com"`
	return fmt.Sprintf(`{%s, %s}`, urls, properties)
}

func repeat(unit string, times int) string {
	var buf strings.Builder
	for i := 0; i < times; i++ {
		buf.WriteString(unit)
	}
	return buf.String()
}

func yamlMultiline(indentation int, s string) string {
	indentPrefix := repeat(" ", indentation)
	var buf strings.Builder
	fmt.Fprintf(&buf, "|\n")
	scanner := bufio.NewScanner(bytes.NewBufferString(s))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(&buf, "%s%s\n", indentPrefix, line)
	}
	if err := scanner.Err(); err != nil {
		return err.Error()
	}
	return buf.String()
}

const projectWithGceEncryptionFmt = `apiVersion: atlas.mongodb.com/v1
kind: AtlasProject
metadata:
  name: my-project
spec:
  name: Test Atlas Operator Project
  encryptionAtRest:
    googleCloudKms:
      enabled: true
      keyVersionResourceID: %s
      serviceAccountKey: %s`

// composeGoogleEncryptionAtRestProjectYAML produces something like this YAML:
//
//	apiVersion: atlas.mongodb.com/v1
//	kind: AtlasProject
//	metadata:
//	  name: my-project
//	spec:
//	  name: Test Atlas Operator Project
//	  encryptionAtRest:
//	    googleCloudKms:
//	      enabled: true
//	      keyVersionResourceID: projects/...
//	      serviceAccountKey: |
//	        {
//	          "type": "service_account",
//	          "project_id": "...",
//	          ...
//	        }
func composeGoogleEncryptionAtRestProjectYAML(serviceAccountKey, keyVersionResourceID string) string {
	return fmt.Sprintf(projectWithGceEncryptionFmt, keyVersionResourceID, yamlMultiline(8, serviceAccountKey))
}

var _ = Describe("Encryption at rest GCP key validation", Label("encryption-at-rest"), func() {
	var testData *model.TestDataProvider
	_ = BeforeEach(func() {
		checkUpEnvironment()
		checkNSetUpGCPEnvironment()
		testData = model.DataProvider(
			"ear-gcp-key-validation",
			model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
			40000,
			[]func(*model.TestDataProvider){},
		).WithProject(data.DefaultProject())
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		if CurrentSpecReport().Failed() {
			Expect(actions.SaveProjectsToFile(testData.Context, testData.K8SClient, testData.Resources.Namespace)).Should(Succeed())
		}

		By("Delete Resources, Project with Encryption at rest", func() {
			actions.DeleteTestDataProject(testData)
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	DescribeTable("fails if the service account key",
		func(encryption *v1.EncryptionAtRest, errMsg string) {
			testData.Project.Spec.EncryptionAtRest = encryption
			configureManager(testData)
			createProjectWithValidationError(testData, errMsg)
		},
		Entry(
			"is missing",
			&v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: "",
				},
			},
			"missing Google Service Account Key but GCP KMS is enabled",
		),
		Entry(
			"is an empty JSON object",
			&v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: "{}",
				},
			},
			"invalid empty service account key",
		),
		Entry(
			"is an empty JSON array",
			&v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: "[]",
				},
			},
			"cannot unmarshal array into Go value",
		),
		Entry(
			"has a bad PEM string",
			&v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: withProperUrls(`"private_key":"-----BEGIN PRIVATE KEY-----\nMIIEvQblah\n-----END PRIVATE KEY-----\n"`),
				},
			},
			"failed to decode PEM block",
		),
		Entry(
			"contains a bad URL",
			&v1.EncryptionAtRest{
				GoogleCloudKms: v1.GoogleCloudKms{
					Enabled:           toptr.MakePtr(true),
					ServiceAccountKey: withProperUrls(`"token_uri": "http//badurl.example"`),
				},
			},
			"invalid URL address",
		),
	)

	It("correct project works", func() {
		projectYAML := composeGoogleEncryptionAtRestProjectYAML(
			os.Getenv("GCP_SA_CRED"),
			os.Getenv("GOOGLE_KEY_VERSION_RESOURCE_ID"),
		)
		Expect(yaml.Unmarshal(([]byte)(projectYAML), testData.Project)).ToNot(HaveOccurred())

		actions.ProjectCreationFlow(testData)
	})
})
