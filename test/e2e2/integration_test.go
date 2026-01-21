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
	"embed"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

//go:embed integrations/*
var integrations embed.FS

const (
	AtlasThirdPartyIntegrationsCRDName = "atlasthirdpartyintegrations.atlas.mongodb.com"
)

// applyObject converts a client.Object to ApplyConfiguration and applies it using the new Apply() API
// This replaces the deprecated Patch() with client.Apply pattern
func applyObject(ctx context.Context, kubeClient client.Client, obj client.Object, fieldOwner client.FieldOwner) error {
	objUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return fmt.Errorf("failed to convert object to unstructured: %w", err)
	}
	objUnstructuredObj := &unstructured.Unstructured{Object: objUnstructured}
	applyConfig := client.ApplyConfigurationFromUnstructured(objUnstructuredObj)
	return kubeClient.Apply(ctx, applyConfig, fieldOwner, client.ForceOwnership)
}

var _ = Describe("Atlas Third-Party Integrations Controller", Ordered, Label("integrations-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace

	_ = BeforeAll(func() {
		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(GinkgoT())

		// Register cleanup - this should even when the process is interrupted with Ctrl+C
		// AfterAll is not reliable in such cases.
		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		ctx = context.Background()
		client, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = client
		Expect(kube.AssertCRDs(ctx, kubeClient, &apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: v1.ObjectMeta{Name: AtlasThirdPartyIntegrationsCRDName},
		})).To(Succeed())
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: utils.RandomName("integrations-ctlr-ns"),
		}}
		Expect(kubeClient.Create(ctx, testNamespace))
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

	DescribeTable("Integrations samples",
		func(objs []client.Object, updates []client.Object, wantReady string) {
			By("Prepare & apply test case objects", func() {
				for _, obj := range objs {
					objToApply := WithRandomAtlasProject(kube.WithRenamedNamespace(obj, testNamespace.Name))
					Expect(
						applyObject(ctx, kubeClient, objToApply, GinkGoFieldOwner),
					).To(Succeed())
				}
			})

			integration := akov2.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{Name: wantReady, Namespace: testNamespace.Name},
			}
			By("Wait integration to be Ready", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(
						kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration),
					).To(Succeed())
					if condition := meta.FindStatusCondition(integration.GetConditions(), "Ready"); condition != nil {
						return condition.Status == metav1.ConditionTrue
					}
					return false
				}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
			})

			By("Apply updates", func() {
				for _, objUpdate := range updates {
					objToPatch := WithRandomAtlasProject(kube.WithRenamedNamespace(objUpdate, testNamespace.Name))
					Expect(
						applyObject(ctx, kubeClient, objToPatch, GinkGoFieldOwner),
					).To(Succeed())
				}
			})

			By("Wait integration to be Ready & updated", func() {
				Eventually(func(g Gomega) bool {
					g.Expect(
						kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration),
					).To(Succeed())
					ready := false
					if condition := meta.FindStatusCondition(integration.GetConditions(), "Ready"); condition != nil {
						ready = (condition.Status == metav1.ConditionTrue)
					}
					if ready {
						if condition := meta.FindStatusCondition(integration.GetConditions(), "State"); condition != nil {
							return state.ResourceState(condition.Reason) == state.StateUpdated
						}
					}
					return false
				}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
			})

			By("Delete integration", func() {
				Expect(kubeClient.Delete(ctx, &integration)).To(Succeed())
			})

			By("Wait integration to be gone", func() {
				Eventually(func(g Gomega) error {
					err := kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration)
					return err
				}).WithTimeout(time.Minute).WithPolling(time.Second).NotTo(Succeed())
			})
		},
		Entry("Test[datadog]: Datadog integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/datadog.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/datadog.update.yml")),
			"atlas-datadog-integ",
		),
		Entry("Test[msteams]: Microsoft Teams integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/msteams.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/msteams.update.yml")),
			"atlas-msteams-integ",
		),
		Entry("Test[newrelic]: New Relic integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/newrelic.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/newrelic.update.yml")),
			"atlas-newrelic-integ",
		),
		Entry("Test[opsgenie]: Ops Genie integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/opsgenie.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/opsgenie.update.yml")),
			"atlas-opsgenie-integ",
		),
		Entry("Test[pagerduty]: PagerDuty integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/pagerduty.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/pagerduty.update.yml")),
			"atlas-pagerduty-integ",
		),
		Entry("Test[prometheus]: Prometheus integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/prometheus.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/prometheus.update.yml")),
			"atlas-prometheus-integ",
		),
		Entry("Test[slack]: Slack integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/slack.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/slack.update.yml")),
			"atlas-slack-integ",
		),
		Entry("Test[victorops]: Victor Ops integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/victorops.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/victorops.update.yml")),
			"atlas-victorops-integ",
		),
		Entry("Test[webhook]: Webhooks integration with a parent project",
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/webhook.sample.yml")),
			yml.MustParseObjects(yml.MustOpen(integrations, "integrations/webhook.update.yml")),
			"atlas-webhook-integ",
		),
	)

	It("Can handle isolated integrations", func() {
		project := akov2.AtlasProject{
			TypeMeta:   v1.TypeMeta{Kind: "AtlasProject", APIVersion: akov2.GroupVersion.String()},
			ObjectMeta: v1.ObjectMeta{Name: "atlas-project", Namespace: testNamespace.Name},
			Spec:       akov2.AtlasProjectSpec{Name: utils.RandomName("atlas-project")},
		}
		projectID := ""

		By("Create Atlas Project", func() {
			Expect(
				applyObject(ctx, kubeClient, &project, GinkGoFieldOwner),
			).To(Succeed())
		})

		By("Wait for Atlas Project ID", func() {
			Eventually(func(g Gomega) bool {
				projectInKube := akov2.AtlasProject{}
				g.Expect(
					kubeClient.Get(ctx, client.ObjectKeyFromObject(&project), &projectInKube),
				).To(Succeed())
				if projectInKube.Status.ID != "" {
					projectID = projectInKube.Status.ID
					return true
				}
				return false
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
		})

		integrationSecret := corev1.Secret{
			TypeMeta: v1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
			ObjectMeta: v1.ObjectMeta{
				Name:      "victor-ops-secret",
				Namespace: testNamespace.Name,
				Labels: map[string]string{
					secretservice.TypeLabelKey: secretservice.CredLabelVal,
				},
			},
			Data: map[string][]byte{
				"apiKey": ([]byte)("00000000-0000-0000-0000-000000000000"),
			},
		}
		integration := akov2.AtlasThirdPartyIntegration{
			TypeMeta:   v1.TypeMeta{Kind: "AtlasThirdPartyIntegration", APIVersion: akov2.GroupVersion.String()},
			ObjectMeta: v1.ObjectMeta{Name: "test-victor-ops-integration", Namespace: testNamespace.Name},
			Spec: akov2.AtlasThirdPartyIntegrationSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					ConnectionSecret: &api.LocalObjectReference{},
				},
				Type: "VICTOR_OPS",
				VictorOps: &akov2.VictorOpsIntegration{
					RoutingKey: "routing-key",
					APIKeySecretRef: api.LocalObjectReference{
						Name: integrationSecret.Name,
					},
				},
			},
		}

		By("Apply isolated integration", func() {
			globalCredsKey := client.ObjectKey{
				Name:      DefaultGlobalCredentials,
				Namespace: control.MustEnvVar("OPERATOR_NAMESPACE"),
			}
			credentialsSecret, err := copySecretToNamespace(ctx, kubeClient, globalCredsKey, testNamespace.Name)
			Expect(err).NotTo(HaveOccurred())
			integration.Spec.ConnectionSecret.Name = credentialsSecret.Name

			for _, obj := range []client.Object{credentialsSecret, &integrationSecret, &integration} {
				Expect(
					applyObject(ctx, kubeClient, obj, GinkGoFieldOwner),
				).To(Succeed())
			}
		})

		By("Wait integration to be Ready", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(
					kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration),
				).To(Succeed())
				if condition := meta.FindStatusCondition(integration.GetConditions(), "Ready"); condition != nil {
					return condition.Status == metav1.ConditionTrue
				}
				return false
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
		})

		By("Update integration", func() {
			updatedIntegration := akov2.AtlasThirdPartyIntegration{
				TypeMeta:   v1.TypeMeta{Kind: "AtlasThirdPartyIntegration", APIVersion: akov2.GroupVersion.String()},
				ObjectMeta: v1.ObjectMeta{Name: "test-victor-ops-integration", Namespace: testNamespace.Name},
				Spec:       integration.Spec,
			}
			updatedIntegration.Spec.VictorOps.RoutingKey = "another-routing-key"
			Expect(
				applyObject(ctx, kubeClient, &updatedIntegration, GinkGoFieldOwner),
			).To(Succeed())
		})

		By("Wait integration to be Ready & updated", func() {
			Eventually(func(g Gomega) bool {
				g.Expect(
					kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration),
				).To(Succeed())
				ready := false
				if condition := meta.FindStatusCondition(integration.GetConditions(), "Ready"); condition != nil {
					ready = (condition.Status == metav1.ConditionTrue)
				}
				if ready {
					if condition := meta.FindStatusCondition(integration.GetConditions(), "State"); condition != nil {
						return state.ResourceState(condition.Reason) == state.StateUpdated
					}
				}
				return false
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
		})

		By("Delete integration", func() {
			Expect(kubeClient.Delete(ctx, &integration)).To(Succeed())
		})

		By("Wait integration to be gone", func() {
			Eventually(func(g Gomega) error {
				err := kubeClient.Get(ctx, client.ObjectKeyFromObject(&integration), &integration)
				return err
			}).WithTimeout(time.Minute).WithPolling(time.Second).NotTo(Succeed())
		})
	})
})

func WithRandomAtlasProject(obj client.Object) client.Object {
	if project, ok := (obj).(*akov2.AtlasProject); ok {
		renamed := project.DeepCopy()
		renamed.Spec.Name = utils.RandomName(project.Spec.Name)
		return renamed
	}
	return obj
}

func copySecretToNamespace(ctx context.Context, kubeClient client.Client, key client.ObjectKey, ns string) (*corev1.Secret, error) {
	secret := corev1.Secret{}
	if err := kubeClient.Get(ctx, key, &secret); err != nil {
		return nil, fmt.Errorf("failed to load original secret %v: %w", key, err)
	}
	return &corev1.Secret{
		TypeMeta: v1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: v1.ObjectMeta{
			Name:      key.Name,
			Namespace: ns,
			Labels:    secret.Labels,
		},
		Data: secret.Data,
	}, nil
}
