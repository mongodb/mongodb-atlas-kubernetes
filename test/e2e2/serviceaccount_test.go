// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e2_test

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb-forks/digest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v20250312018 "go.mongodb.org/atlas-sdk/v20250312018/admin"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/accesstoken"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
)

type serviceAccountCreds struct {
	clientID     string
	clientSecret string
}

func createAtlasServiceAccount(ctx context.Context, orgID string) (*serviceAccountCreds, func()) {
	publicKey := control.MustEnvVar("ATLAS_PUBLIC_KEY")
	privateKey := control.MustEnvVar("ATLAS_PRIVATE_KEY")

	transport := digest.NewTransport(publicKey, privateKey)
	httpClient, err := transport.Client()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	atlasClient, err := v20250312018.NewClient(
		v20250312018.UseBaseURL("https://cloud-qa.mongodb.com"),
		v20250312018.UseHTTPClient(httpClient),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	saName := utils.RandomName("ako-e2e-sa")
	sa, _, err := atlasClient.ServiceAccountsApi.
		CreateOrgServiceAccount(ctx, orgID, &v20250312018.OrgServiceAccountRequest{
			Name:                    saName,
			Description:             fmt.Sprintf("AKO e2e test service account %s", saName),
			Roles:                   []string{"ORG_OWNER"},
			SecretExpiresAfterHours: 8,
		}).Execute()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "failed to create Atlas service account")
	ExpectWithOffset(1, sa.ClientId).NotTo(BeNil())
	ExpectWithOffset(1, sa.Secrets).NotTo(BeNil())
	ExpectWithOffset(1, len(*sa.Secrets)).To(BeNumerically(">", 0))

	secret := (*sa.Secrets)[0]
	ExpectWithOffset(1, secret.Secret).NotTo(BeNil())

	clientID := *sa.ClientId
	cleanup := func() {
		_, delErr := atlasClient.ServiceAccountsApi.
			DeleteOrgServiceAccount(ctx, clientID, orgID).Execute()
		if delErr != nil {
			GinkgoWriter.Printf("WARNING: failed to delete service account %s: %v\n", clientID, delErr)
		}
	}

	return &serviceAccountCreds{
		clientID:     clientID,
		clientSecret: *secret.Secret,
	}, cleanup
}

func createServiceAccountCredentialSecret(
	ctx context.Context,
	kubeClient client.Client,
	name, namespace string,
	creds *serviceAccountCreds,
	orgID string,
) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		Data: map[string][]byte{
			"orgId":        []byte(orgID),
			"clientId":     []byte(creds.clientID),
			"clientSecret": []byte(creds.clientSecret),
		},
	}
	ExpectWithOffset(1, kubeClient.Create(ctx, secret)).To(Succeed())
	return secret
}

func waitForAccessTokenSecret(ctx context.Context, kubeClient client.Client, secret *corev1.Secret) *corev1.Secret {
	expectedName, _ := accesstoken.DeriveSecretName(secret.Namespace, secret.Name)
	tokenSecret := &corev1.Secret{}
	EventuallyWithOffset(1, func(g Gomega) bool {
		err := kubeClient.Get(ctx, client.ObjectKey{
			Name:      expectedName,
			Namespace: secret.Namespace,
		}, tokenSecret)
		return err == nil
	}).WithTimeout(2*time.Minute).WithPolling(2*time.Second).Should(BeTrue(),
		"Expected access token secret to be created by the operator")
	return tokenSecret
}

var _ = Describe("Service Account Controller", Ordered, Label("service-account"), func() {
	ctx := context.Background()
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string

	_ = BeforeAll(func() {
		orgID = control.MustEnvVar("ATLAS_ORG_ID")
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), false)
		ako.Start(ctx, GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		client, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = client
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: utils.RandomName("sa-ctlr-ns"),
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
	})

	_ = AfterEach(func() {
		if kubeClient == nil {
			return
		}
		Expect(
			kubeClient.Delete(ctx, testNamespace),
		).To(Succeed())
		Eventually(func(g Gomega) bool {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace) == nil
		}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	It("creates an access token secret for a service account credential secret", func() {
		By("Creating an Atlas Service Account via the API")
		saCreds, cleanupSA := createAtlasServiceAccount(ctx, orgID)
		DeferCleanup(cleanupSA)

		credentialSecret := createServiceAccountCredentialSecret(ctx, kubeClient, "sa-test-credentials", testNamespace.Name, saCreds, orgID)

		By("Waiting for the access token secret to be created by the operator")
		tokenSecret := waitForAccessTokenSecret(ctx, kubeClient, credentialSecret)

		By("Verifying the access token secret has correct fields")

		Expect(tokenSecret.Data).To(HaveKey("accessToken"))
		Expect(tokenSecret.Data).To(HaveKey("expiry"))
		Expect(string(tokenSecret.Data["accessToken"])).NotTo(BeEmpty())

		Expect(tokenSecret.Labels).To(HaveKeyWithValue(
			secretservice.TypeLabelKey, secretservice.CredLabelVal,
		))

		Expect(tokenSecret.OwnerReferences).To(HaveLen(1))
		Expect(tokenSecret.OwnerReferences[0].Name).To(Equal("sa-test-credentials"))
	})

	It("does not create a token for an API key credential secret", func() {
		apiKeySecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "api-key-credentials",
				Namespace: testNamespace.Name,
				Labels: map[string]string{
					secretservice.TypeLabelKey: secretservice.CredLabelVal,
				},
			},
			Data: map[string][]byte{
				"orgId":         []byte("test-org-id"),
				"publicApiKey":  []byte("test-public-key"),
				"privateApiKey": []byte("test-private-key"),
			},
		}

		By("Creating the API key credential secret")
		Expect(kubeClient.Create(ctx, apiKeySecret)).To(Succeed())

		By("Verifying no access token secret is created")
		expectedName, _ := accesstoken.DeriveSecretName(apiKeySecret.Namespace, apiKeySecret.Name)
		Consistently(func(g Gomega) bool {
			tokenSecret := &corev1.Secret{}
			err := kubeClient.Get(ctx, client.ObjectKey{
				Name:      expectedName,
				Namespace: apiKeySecret.Namespace,
			}, tokenSecret)
			return apierrors.IsNotFound(err)
		}).WithTimeout(15*time.Second).WithPolling(2*time.Second).Should(BeTrue(),
			"API key secret should not get an access token secret")
	})

	It("refreshes the access token when credentials are rotated in place", func() {
		By("Creating two Atlas Service Accounts representing the pre- and post-rotation credentials")
		firstCreds, cleanupFirst := createAtlasServiceAccount(ctx, orgID)
		DeferCleanup(cleanupFirst)
		secondCreds, cleanupSecond := createAtlasServiceAccount(ctx, orgID)
		DeferCleanup(cleanupSecond)

		credentialSecret := createServiceAccountCredentialSecret(ctx, kubeClient, "sa-rotate-credentials", testNamespace.Name, firstCreds, orgID)

		By("Waiting for the initial access token secret to be minted from the first credentials")
		initialToken := waitForAccessTokenSecret(ctx, kubeClient, credentialSecret)
		initialAccessToken := string(initialToken.Data["accessToken"])
		Expect(initialAccessToken).NotTo(BeEmpty())
		initialHash, err := accesstoken.CredentialsHash(firstCreds.clientID, firstCreds.clientSecret)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(initialToken.Data["credentialsHash"])).To(Equal(initialHash),
			"initial token secret must record the hash of the first credentials")

		By("Rotating clientId/clientSecret on the credential Secret in place")
		Eventually(func(g Gomega) {
			fresh := &corev1.Secret{}
			g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(credentialSecret), fresh)).To(Succeed())
			fresh.Data["clientId"] = []byte(secondCreds.clientID)
			fresh.Data["clientSecret"] = []byte(secondCreds.clientSecret)
			g.Expect(kubeClient.Update(ctx, fresh)).To(Succeed())
		}).WithTimeout(30 * time.Second).WithPolling(2 * time.Second).Should(Succeed())

		By("Waiting for the access token secret to be refreshed with the second credentials")
		expectedHash, err := accesstoken.CredentialsHash(secondCreds.clientID, secondCreds.clientSecret)
		Expect(err).NotTo(HaveOccurred())
		Eventually(func(g Gomega) {
			tokenName, _ := accesstoken.DeriveSecretName(credentialSecret.Namespace, credentialSecret.Name)
			refreshed := &corev1.Secret{}
			g.Expect(kubeClient.Get(ctx, client.ObjectKey{Name: tokenName, Namespace: credentialSecret.Namespace}, refreshed)).To(Succeed())
			g.Expect(string(refreshed.Data["credentialsHash"])).To(Equal(expectedHash),
				"credentials hash must reflect the rotated clientId/clientSecret")
			g.Expect(string(refreshed.Data["accessToken"])).NotTo(Equal(initialAccessToken),
				"bearer token must be replaced after rotation")
			g.Expect(string(refreshed.Data["accessToken"])).NotTo(BeEmpty())
		}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
	})

	It("creates an AtlasProject using service account auth", func() {
		By("Creating an Atlas Service Account via the API")
		saCreds, cleanupSA := createAtlasServiceAccount(ctx, orgID)
		DeferCleanup(cleanupSA)

		resourcePrefix := utils.RandomName("sa-project")
		credSecretName := resourcePrefix + "-creds"

		By("Creating the service account credential secret")
		credentialSecret := createServiceAccountCredentialSecret(ctx, kubeClient, credSecretName, testNamespace.Name, saCreds, orgID)

		By("Waiting for the access token to be ready")
		waitForAccessTokenSecret(ctx, kubeClient, credentialSecret)

		By("Creating the AtlasProject referencing the SA credential secret")
		atlasProject := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-project",
				Namespace: testNamespace.Name,
			},
			Spec: akov2.AtlasProjectSpec{
				Name:                      resourcePrefix,
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				ConnectionSecret: &common.ResourceRefNamespaced{
					Name: credSecretName,
				},
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "0.0.0.0/1",
						Comment:   "Everyone has access. For the test purpose only.",
					},
					{
						CIDRBlock: "128.0.0.0/1",
						Comment:   "Everyone has access. For the test purpose only.",
					},
				},
			},
		}
		Expect(kubeClient.Create(ctx, atlasProject)).To(Succeed())

		By("Waiting for the AtlasProject to become Ready")
		Eventually(func(g Gomega) {
			condition, err := k8s.GetProjectStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, "test-project")
			g.Expect(err).To(BeNil())
			g.Expect(condition).To(Equal("True"))
		}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

		By("Verifying the project has an Atlas ID in status")
		projectInKube := &akov2.AtlasProject{}
		Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(atlasProject), projectInKube)).To(Succeed())
		Expect(projectInKube.Status.ID).NotTo(BeEmpty())

		By("Deleting the AtlasProject")
		Expect(kubeClient.Delete(ctx, atlasProject)).To(Succeed())
		Eventually(func(g Gomega) error {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(atlasProject), &akov2.AtlasProject{})
		}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).ShouldNot(Succeed())
	})

	It("creates an AtlasProject and AtlasDeployment (flex) using service account auth", func() {
		By("Creating an Atlas Service Account via the API")
		saCreds, cleanupSA := createAtlasServiceAccount(ctx, orgID)
		DeferCleanup(cleanupSA)

		resourcePrefix := utils.RandomName("sa-flex")
		credSecretName := resourcePrefix + "-creds"

		By("Creating the service account credential secret")
		credentialSecret := createServiceAccountCredentialSecret(ctx, kubeClient, credSecretName, testNamespace.Name, saCreds, orgID)

		By("Waiting for the access token to be ready")
		waitForAccessTokenSecret(ctx, kubeClient, credentialSecret)

		By("Creating the AtlasProject")
		atlasProject := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-project",
				Namespace: testNamespace.Name,
			},
			Spec: akov2.AtlasProjectSpec{
				Name:                      resourcePrefix,
				RegionUsageRestrictions:   "NONE",
				WithDefaultAlertsSettings: true,
				ConnectionSecret: &common.ResourceRefNamespaced{
					Name: credSecretName,
				},
				ProjectIPAccessList: []project.IPAccessList{
					{
						CIDRBlock: "0.0.0.0/1",
						Comment:   "Everyone has access. For the test purpose only.",
					},
					{
						CIDRBlock: "128.0.0.0/1",
						Comment:   "Everyone has access. For the test purpose only.",
					},
				},
			},
		}
		Expect(kubeClient.Create(ctx, atlasProject)).To(Succeed())

		By("Waiting for the AtlasProject to become Ready")
		Eventually(func(g Gomega) {
			condition, err := k8s.GetProjectStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, "test-project")
			g.Expect(err).To(BeNil())
			g.Expect(condition).To(Equal("True"))
		}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

		By("Creating the AtlasDeployment (flex) referencing the project")
		flexDeployment := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "flex",
				Namespace: testNamespace.Name,
			},
			Spec: akov2.AtlasDeploymentSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{
						Name: "test-project",
					},
				},
				FlexSpec: &akov2.FlexSpec{
					Name: resourcePrefix + "-flex",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
				},
			},
		}
		Expect(kubeClient.Create(ctx, flexDeployment)).To(Succeed())

		By("Waiting for the AtlasDeployment to become Ready")
		Eventually(func(g Gomega) {
			condition, err := k8s.GetDeploymentStatusCondition(ctx, kubeClient, api.ReadyType, testNamespace.Name, "flex")
			g.Expect(err).To(BeNil())
			g.Expect(condition).To(Equal("True"))
		}).WithTimeout(10 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())

		By("Deleting the AtlasDeployment")
		Expect(kubeClient.Delete(ctx, flexDeployment)).To(Succeed())
		Eventually(func(g Gomega) error {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(flexDeployment), &akov2.AtlasDeployment{})
		}).WithTimeout(10 * time.Minute).WithPolling(10 * time.Second).ShouldNot(Succeed())

		By("Deleting the AtlasProject")
		Expect(kubeClient.Delete(ctx, atlasProject)).To(Succeed())
		Eventually(func(g Gomega) error {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(atlasProject), &akov2.AtlasProject{})
		}).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).ShouldNot(Succeed())
	})
})
