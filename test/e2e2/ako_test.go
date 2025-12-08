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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
)

const (
	AtlasProjectCRDName = "atlasprojects.atlas.mongodb.com"
)

var _ = Describe("Atlas Operator Start and Stop test", Ordered, Label("ako-start-stop"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace

	_ = BeforeAll(func() {
		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(GinkgoT())

		ctx = context.Background()
		client, err := kube.NewTestClient(false)
		Expect(err).To(Succeed())
		kubeClient = client
		Expect(kube.AssertCRDs(ctx, kubeClient, &apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: v1.ObjectMeta{Name: AtlasProjectCRDName},
		})).To(Succeed())
	})

	_ = AfterAll(func() {
		ako.Stop(GinkgoT())
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: utils.RandomName("ako-ns"),
		}}
		Expect(kubeClient.Create(ctx, testNamespace))
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
	})

	_ = AfterEach(func() {
		Expect(
			kubeClient.Delete(ctx, testNamespace),
		).To(Succeed())
		Eventually(func(g Gomega) bool {
			return kubeClient.Get(ctx, client.ObjectKeyFromObject(testNamespace), testNamespace) == nil
		}).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	It("AKO running", func() {
		testProject := akov2.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test-project",
				Namespace: testNamespace.Name,
			},
			Spec: akov2.AtlasProjectSpec{
				Name: utils.RandomName("test-project"),
			},
		}
		Expect(kubeClient.Create(ctx, &testProject)).To(Succeed())
		Eventually(func() bool {
			kubeProject := akov2.AtlasProject{}
			Expect(
				kubeClient.Get(ctx, client.ObjectKeyFromObject(&testProject), &kubeProject),
			).To(Succeed())
			for _, condition := range kubeProject.Status.Conditions {
				if condition.Type == "Ready" {
					return string(condition.Status) == string(metav1.ConditionTrue)
				}
			}
			return false
		}).WithTimeout(time.Minute).WithPolling(time.Second)
	})
})
