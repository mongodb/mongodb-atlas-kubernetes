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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

//go:embed configs/*
var configs embed.FS

const (
	AtlasThirdPartyIntegrationsCRDName = "atlasthirdpartyintegrations.atlas.nextapi.mongodb.com"
)

var _ = Describe("Atlas Third-Party Integrations Controller", Ordered, Label("integrations-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako *operator.Operator
	testNamespace := corev1.Namespace{}

	_ = BeforeAll(func() {
		// Launch one Operator instance for all tests
		deletionProtectionOff := false
		ako = runTestAKO(control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(GinkgoT())
	})

	_ = AfterAll(func() {
		ako.Stop(GinkgoT())
	})

	_ = BeforeEach(OncePerOrdered, func() {
		testNamespace = corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: utils.RandomName("integrations-ctlr-ns"),
		}}

		ctx = context.Background()
		client, err := kube.NewK8sTest(ctx, &apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: v1.ObjectMeta{Name: AtlasThirdPartyIntegrationsCRDName},
		})
		Expect(err).To(Succeed())
		kubeClient = client
		Expect(kubeClient.Create(ctx, &testNamespace))
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Atlas Third-Party Integrations Controller Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testNamespace.Name + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))

		By("Delete Kubernetes resources", func() {
			Expect(
				kubeClient.Delete(ctx, &testNamespace),
			).To(Succeed())
			Eventually(func(g Gomega) bool {
				return kubeClient.Get(ctx, client.ObjectKeyFromObject(&testNamespace), &testNamespace) == nil
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
		})
	})

	DescribeTable("Integrations samples",
		func(objs []client.Object, wantReady string) {
			By("Prepare and apply test case objects", func() {
				for _, obj := range objs {
					objToApply := WithRandomAtlasProject(kube.WithRenamedNamespace(obj, testNamespace.Name))
					Expect(
						kubeClient.Patch(ctx, objToApply, client.Apply, client.ForceOwnership, GinkGoFieldOwner),
					).To(Succeed())
				}
			})

			By("Wait main Object to be Ready", func() {
				integration := akov2next.AtlasThirdPartyIntegration{
					ObjectMeta: v1.ObjectMeta{Name: wantReady, Namespace: testNamespace.Name},
				}
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
		},
		Entry("Test[datadog]: Datadog integration with a parent project",
			Label("datadog"),
			yml.MustParseObjects(yml.MustOpen(configs, "configs/datadog.sample.yml")),
			"atlas-datadog-integ",
		),
	)
})

func WithRandomAtlasProject(obj client.Object) client.Object {
	if project, ok := (obj).(*akov2.AtlasProject); ok {
		renamed := project.DeepCopy()
		renamed.Spec.Name = utils.RandomName(project.Spec.Name)
		return renamed
	}
	return obj
}
