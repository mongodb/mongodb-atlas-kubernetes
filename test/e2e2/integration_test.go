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
	AtlasThirdPartyIntegrationsCRD = "atlasthirdpartyintegrations.atlas.nextapi.mongodb.com"
)

var _ = Describe("Atlas Third-Party Integrations Controller", Ordered, Label("integrations-ctlr"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako *operator.Operator
	testNamespace := ""

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
		testNamespace = utils.RandomName("integrations-ctlr-ns")

		ctx = context.Background()
		client, err := kube.NewK8sTest(ctx, AtlasThirdPartyIntegrationsCRD)
		Expect(err).To(Succeed())
		kubeClient = client
		Expect(kube.CreateNamespace(ctx, kubeClient, testNamespace))
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
	})

	_ = AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Atlas Third-Party Integrations Controller Test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testNamespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))

		By("Delete Kubernetes resources", func() {
			Expect(kube.WipeNamespace(ctx, kubeClient, testNamespace)).To(Succeed())
			Eventually(func(g Gomega) bool {
				return kube.HasNamespace(ctx, kubeClient, testNamespace)
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
		})
	})

	DescribeTable("Integrations samples",
		func(objs []client.Object, wantReady string) {
			Expect(kube.Apply(ctx, kubeClient, fix(objs, testNamespace)...)).To(Succeed())

			integration := akov2next.AtlasThirdPartyIntegration{
				ObjectMeta: v1.ObjectMeta{
					Name:      wantReady,
					Namespace: testNamespace,
				},
			}

			key := client.ObjectKeyFromObject(&integration)
			Eventually(func(g Gomega) bool {
				ok, err := kube.AssertObjReady(ctx, kubeClient, key, &integration)
				Expect(err).To(Succeed())
				return ok
			}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeTrue())
		},
		Entry("Test[datadog]: Datadog integration with a parent project",
			Label("datadog"),
			yml.MustParseCRs(yml.MustOpen(configs, "configs/datadog.sample.yml")),
			"atlas-datadog-integ",
		),
	)
})

func fix(objs []client.Object, namespace string) []client.Object {
	allFixed := make([]client.Object, 0, len(objs))
	for _, obj := range objs {
		fixed := kube.SetNamespace(obj, namespace)
		if project, ok := (fixed).(*akov2.AtlasProject); ok {
			randomizeProjectName(project)
		}
		allFixed = append(allFixed, fixed)
	}
	return allFixed
}

func randomizeProjectName(project *akov2.AtlasProject) {
	project.Spec.Name = utils.RandomName(project.Spec.Name)
}
