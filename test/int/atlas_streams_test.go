package int

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/resources"
)

var _ = Describe("AtlasStreams", Label("int", "AtlasStreams"), func() {
	var testNamespace *corev1.Namespace
	var stopManager context.CancelFunc
	var connectionSecret corev1.Secret
	var projectName string
	var testProject *akov2.AtlasProject
	resourceName := "stream-instance-0"
	kafkaUserPassSecretName := "kafka-userpass"
	kafkaCertificateSecretName := "kafka-certificate" //nolint:gosec
	certificate := `-----BEGIN CERTIFICATE-----
MIIEITCCAwmgAwIBAgIUTLX+HHPxjMxw1pOXEu/+m+aXrgIwDQYJKoZIhvcNAQEL
BQAwgZ8xCzAJBgNVBAYTAkRFMQ8wDQYDVQQIDAZCZXJsaW4xDzANBgNVBAcMBkJl
cmxpbjEVMBMGA1UECgwMTW9uZ29EQiBHbWJoMRMwEQYDVQQLDApLdWJlcm5ldGVz
MRcwFQYDVQQDDA5BdGxhcyBPcGVyYXRvcjEpMCcGCSqGSIb3DQEJARYaaGVsZGVy
LnNhbnRhbmFAbW9uZ29kYi5jb20wHhcNMjQwNDIzMTE0NzI2WhcNMjcwMTE4MTE0
NzI2WjCBnzELMAkGA1UEBhMCREUxDzANBgNVBAgMBkJlcmxpbjEPMA0GA1UEBwwG
QmVybGluMRUwEwYDVQQKDAxNb25nb0RCIEdtYmgxEzARBgNVBAsMCkt1YmVybmV0
ZXMxFzAVBgNVBAMMDkF0bGFzIE9wZXJhdG9yMSkwJwYJKoZIhvcNAQkBFhpoZWxk
ZXIuc2FudGFuYUBtb25nb2RiLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAKoBtN0V9F8ZnbPJMKDZ0jHRw35Y/jtZpdN6z824nyRh4U4FeLaAOzex
EiHrxDt9IccxKcVc/9WAq7Pn1C42YJFy9dgLSD94TW4lJwLhAsGxI5bVy+ls6c3u
cpiPzaoUU1vx+Gg5ob+UefjAf7WxaRnuSiUpYPVVueZ218Hhc1W8yajfwLdshXiN
NaBox2Pu+ofsq5aM1T4MARsLODUJqzoQHR2275oFPNaz2BgBgRUDkICw+RPfjQ0X
lCkCtHy2QeBb5hGOi0lG89C9lbuEXb5YOzGG4Cc6snZGf21MGxXAXiL/KsBZrP5i
edABbwkXEgLk41OcwNgshuADM7iOd9sCAwEAAaNTMFEwHQYDVR0OBBYEFBiwIuyh
3sqgzfcgKb80FF1WByAIMB8GA1UdIwQYMBaAFBiwIuyh3sqgzfcgKb80FF1WByAI
MA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAB0iWV/hpK1WuxjS
h5HAfRxBCyWFIU14S7tQHTPuQANQAh3Zktkghpmc6hdNb3VjKzVUSTv9Ye6V22mh
Resf7PVWFvOdPoiJnmJjUQ5W3FUVZWOgx3rFlKO/5HOi5wRvBDyuZsTjIEJP5MOl
3lBs17FOVqM3iT785oabOEj/8LhkvdG9brobG8oAttUSPChiYbEtH83WqgeHnCWI
reLAKIvG8bFVaokdInEgoRt5uque70g0tqAje9MXqCodB96Lo1tk8yyvX4jWI2Pb
pe7aAzw79hIH3tyw+FHjZLgHAq77E14xBxMxvamSnsqGhvCkb7pRHD5+l4tg2k/N
YJZC5C0=
-----END CERTIFICATE-----
`

	BeforeEach(func() {
		By("Starting the operator", func() {
			testNamespace, stopManager = prepareControllers(false)
			Expect(testNamespace).ToNot(BeNil())
			Expect(stopManager).ToNot(BeNil())
		})

		By("Creating project connection secret", func() {
			connectionSecret = buildConnectionSecret(fmt.Sprintf("%s-atlas-key", testNamespace.Name))
			Expect(k8sClient.Create(context.Background(), &connectionSecret)).To(Succeed())
		})

		By("Creating a project", func() {
			testProject = &akov2.AtlasProject{}
			projectName = fmt.Sprintf("new-project-%s", testNamespace.Name)

			testProject = akov2.NewProject(testNamespace.Name, projectName, projectName).
				WithConnectionSecret(connectionSecret.Name)
			Expect(k8sClient.Create(context.Background(), testProject)).To(Succeed())

			Eventually(func() bool {
				return resources.CheckCondition(k8sClient, testProject, status.TrueCondition(status.ReadyType))
			}).WithTimeout(15 * time.Minute).WithPolling(PollingInterval).Should(BeTrue())
		})
	})

	Describe("When creating a stream instance with 2 connections", func() {
		It("Should successfully manage instance and connections", func() {
			var sampleConnection *akov2.AtlasStreamConnection
			var kafkaConnection *akov2.AtlasStreamConnection

			By("Adding a sample connection", func() {
				sampleConnection = &akov2.AtlasStreamConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "sample-connection",
						Namespace: testNamespace.Name,
					},
					Spec: akov2.AtlasStreamConnectionSpec{
						Name:           "sample_stream_solar",
						ConnectionType: "Sample",
					},
				}

				Expect(k8sClient.Create(context.Background(), sampleConnection)).To(Succeed())
				Expect(sampleConnection.GetFinalizers()).To(BeEmpty())
			})

			By("Adding a Kafka connection", func() {
				kafkaUserPassSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaUserPassSecretName,
						Namespace: testNamespace.Name,
						Labels: map[string]string{
							"atlas.mongodb.com/type": "credentials",
						},
					},
					StringData: map[string]string{
						"username": "kafka_user",
						"password": "kafka_pass",
					},
				}
				Expect(k8sClient.Create(context.Background(), kafkaUserPassSecret)).To(Succeed())

				kafkaCertificateSecret := &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      kafkaCertificateSecretName,
						Namespace: testNamespace.Name,
						Labels: map[string]string{
							"atlas.mongodb.com/type": "credentials",
						},
					},
					StringData: map[string]string{
						"certificate": certificate,
					},
				}
				Expect(k8sClient.Create(context.Background(), kafkaCertificateSecret)).To(Succeed())

				kafkaConnection = &akov2.AtlasStreamConnection{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kafka-connection",
						Namespace: testNamespace.Name,
					},
					Spec: akov2.AtlasStreamConnectionSpec{
						Name:           "kafka-config",
						ConnectionType: "Kafka",
						KafkaConfig: &akov2.StreamsKafkaConnection{
							Authentication: akov2.StreamsKafkaAuthentication{
								Mechanism: "SCRAM-512",
								Credentials: common.ResourceRefNamespaced{
									Name:      kafkaUserPassSecret.Name,
									Namespace: kafkaUserPassSecret.Namespace,
								},
							},
							BootstrapServers: "kafka.server1:9001,kafka.server2:9002",
							Security: akov2.StreamsKafkaSecurity{
								Protocol: "SSL",
								Certificate: common.ResourceRefNamespaced{
									Name:      kafkaCertificateSecret.Name,
									Namespace: kafkaCertificateSecret.Namespace,
								},
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), kafkaConnection)).To(Succeed())
				Expect(kafkaConnection.GetFinalizers()).To(BeEmpty())
			})

			By("Creating an instance", func() {
				streamInstance := &akov2.AtlasStreamInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: testNamespace.Name,
					},
					Spec: akov2.AtlasStreamInstanceSpec{
						Name: resourceName,
						Config: akov2.Config{
							Provider: "AWS",
							Region:   "VIRGINIA_USA",
							Tier:     "SP10",
						},
						Project: common.ResourceRefNamespaced{
							Name:      testProject.Name,
							Namespace: testProject.Namespace,
						},
						ConnectionRegistry: []common.ResourceRefNamespaced{
							{
								Name:      sampleConnection.Name,
								Namespace: sampleConnection.Namespace,
							},
							{
								Name:      kafkaConnection.Name,
								Namespace: kafkaConnection.Namespace,
							},
						},
					},
				}
				Expect(k8sClient.Create(context.Background(), streamInstance)).To(Succeed())

				checkInstanceIsReady(client.ObjectKeyFromObject(streamInstance))
			})

			By("Updating the instance", func() {
				Eventually(func(g Gomega) {
					streamInstance := &akov2.AtlasStreamInstance{}
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, streamInstance)).To(Succeed())

					streamInstance.Spec.Config.Region = "DUBLIN_IRL"
					g.Expect(k8sClient.Update(context.Background(), streamInstance)).To(Succeed())
				}).WithTimeout(time.Minute).WithPolling(PollingInterval)

				checkInstanceIsReady(client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name})
			})

			By("Updating a connection", func() {
				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(kafkaConnection), kafkaConnection)).To(Succeed())
					kafkaConnection.Spec.KafkaConfig.BootstrapServers = "kafka.server1:9001,kafka.server2:9002,kafka.server3:9003"
					g.Expect(k8sClient.Update(context.Background(), kafkaConnection)).To(Succeed())
				}).WithTimeout(time.Minute).WithPolling(PollingInterval)

				checkInstanceIsReady(client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name})
			})

			By("Updating a secret", func() {
				Eventually(func(g Gomega) {
					s := corev1.Secret{}
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: kafkaUserPassSecretName, Namespace: testNamespace.Name}, &s)).To(Succeed())
					s.Data["username"] = []byte("kafka_user_changed")
					g.Expect(k8sClient.Update(context.Background(), &s)).To(Succeed())
				}).WithTimeout(time.Minute).WithPolling(PollingInterval)

				checkInstanceIsReady(client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name})
			})

			By("Releasing a connection when removed from instance", func() {
				streamInstance := &akov2.AtlasStreamInstance{}
				Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, streamInstance)).To(Succeed())

				streamInstance.Spec.ConnectionRegistry = []common.ResourceRefNamespaced{
					{
						Name:      sampleConnection.Name,
						Namespace: sampleConnection.Namespace,
					},
				}
				Expect(k8sClient.Update(context.Background(), streamInstance)).To(Succeed())

				checkInstanceIsReady(client.ObjectKeyFromObject(streamInstance))

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(kafkaConnection), kafkaConnection)).To(Succeed())
					g.Expect(kafkaConnection.GetFinalizers()).To(BeEmpty())
				}).WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})

			By("Deleting instance and connections", func() {
				Expect(k8sClient.Delete(context.Background(), kafkaConnection)).To(Succeed())
				Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(kafkaConnection), kafkaConnection)).ToNot(Succeed())

				Expect(k8sClient.Delete(context.Background(), sampleConnection)).To(Succeed())
				Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(sampleConnection), sampleConnection)).To(Succeed())
				Expect(sampleConnection.DeletionTimestamp).ShouldNot(BeNil())

				streamInstance := &akov2.AtlasStreamInstance{}
				Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: resourceName, Namespace: testNamespace.Name}, streamInstance)).To(Succeed())
				Expect(k8sClient.Delete(context.Background(), streamInstance))

				Eventually(func(g Gomega) {
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(sampleConnection), sampleConnection)).ToNot(Succeed())
					g.Expect(k8sClient.Get(context.Background(), client.ObjectKeyFromObject(streamInstance), streamInstance)).ToNot(Succeed())
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			})
		})
	})

	AfterEach(func() {
		By("Deleting stream connection secrets", func() {
			Eventually(func(g Gomega) {
				secret := &corev1.Secret{}
				g.Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: kafkaUserPassSecretName, Namespace: testNamespace.Name}, secret)).To(Succeed())
				g.Expect(k8sClient.Delete(context.Background(), secret)).To(Succeed())
			}).WithTimeout(1 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

			Eventually(func(g Gomega) {
				secret := &corev1.Secret{}
				g.Expect(k8sClient.Get(context.Background(), client.ObjectKey{Name: kafkaCertificateSecretName, Namespace: testNamespace.Name}, secret)).To(Succeed())
				g.Expect(k8sClient.Delete(context.Background(), secret)).To(Succeed())
			}).WithTimeout(1 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
		})

		By("Deleting project", func() {
			if testProject != nil {
				projectID := testProject.ID()
				Expect(k8sClient.Delete(context.Background(), testProject)).To(Succeed())

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.ProjectsApi.GetProject(context.Background(), projectID).Execute()
					g.Expect(err).ToNot(BeNil())
					g.Expect(r).ToNot(BeNil())
					g.Expect(r.StatusCode).To(Equal(http.StatusNotFound))
				}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
			}
		})

		By("Deleting project connection secret", func() {
			Expect(k8sClient.Delete(context.Background(), &connectionSecret)).To(Succeed())
		})

		By("Stopping the operator", func() {
			stopManager()
			err := k8sClient.Delete(context.Background(), testNamespace)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})

func checkInstanceIsReady(instanceObjKey client.ObjectKey) {
	readyConditions := conditions.MatchConditions(
		status.TrueCondition(status.ReadyType),
		status.TrueCondition(status.ResourceVersionStatus),
		status.TrueCondition(status.StreamInstanceReadyType),
	)
	Eventually(func(g Gomega) {
		streamInstance := &akov2.AtlasStreamInstance{}
		g.Expect(k8sClient.Get(context.Background(), instanceObjKey, streamInstance)).To(Succeed())
		g.Expect(streamInstance.Status.Conditions).To(ContainElements(readyConditions))
	}).WithTimeout(5 * time.Minute).WithPolling(PollingInterval).Should(Succeed())

	Eventually(func(g Gomega) {
		streamInstance := &akov2.AtlasStreamInstance{}
		g.Expect(k8sClient.Get(context.Background(), instanceObjKey, streamInstance)).To(Succeed())

		for _, connectionRef := range streamInstance.Spec.ConnectionRegistry {
			connection := &akov2.AtlasStreamConnection{}
			g.Expect(k8sClient.Get(context.Background(), *connectionRef.GetObject(streamInstance.Namespace), connection)).To(Succeed())
			g.Expect(connection.GetFinalizers()).ToNot(BeEmpty())
		}
	}).WithTimeout(2 * time.Minute).WithPolling(PollingInterval).Should(Succeed())
}
