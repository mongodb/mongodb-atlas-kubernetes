// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
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
	"net/http"
	"time"

	k8s "github.com/crd2go/crd2go/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"

	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/resources"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/samples"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/specmatch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
)

var _ = Describe("NetworkContainer CRUD", Ordered, Label("networkcontainer"), func() {
	var ctx context.Context
	var kubeClient client.Client
	var ako operator.Operator
	var testNamespace *corev1.Namespace
	var orgID string
	var atlasClient *admin.APIClient

	_ = BeforeAll(func() {
		ctx = context.Background()

		deletionProtectionOff := false
		ako = runTestAKO(DefaultGlobalCredentials, control.MustEnvVar("OPERATOR_NAMESPACE"), deletionProtectionOff)
		ako.Start(ctx, GinkgoT())

		DeferCleanup(func() {
			if ako != nil {
				ako.Stop(GinkgoT())
			}
		})

		testClient, err := kube.NewTestClient()
		Expect(err).To(Succeed())
		kubeClient = testClient
		Expect(kube.AssertCRDNames(ctx, kubeClient,
			"groups.atlas.generated.mongodb.com",
			"networkcontainers.atlas.generated.mongodb.com",
		)).To(Succeed())

		atlasClient, orgID = newTestAtlasClient()
	})

	_ = BeforeEach(func() {
		testNamespace = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("nc-ns-%s", rand.String(6)),
		}}
		Expect(kubeClient.Create(ctx, testNamespace)).To(Succeed())
		Expect(ako.Running()).To(BeTrue(), "Operator must be running")
		Expect(resources.CopyCredentialsToNamespace(
			ctx,
			kubeClient,
			DefaultGlobalCredentials,
			control.MustEnvVar("OPERATOR_NAMESPACE"),
			testNamespace.Name, GinkGoFieldOwner),
		).To(Succeed())
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
		}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(time.Second).To(BeFalse())
	})

	// createGroup creates a Group in the test namespace and waits for it to be ready in Atlas.
	createGroup := func() *generatedv1.Group {
		groupName := fmt.Sprintf("test-group-%s", rand.String(6))
		params := testparams.New(orgID, control.MustEnvVar("OPERATOR_NAMESPACE"), DefaultGlobalCredentials).
			WithGroupName(groupName).
			WithNamespace(testNamespace.Name)
		objs := samples.MustLoadSampleObjects("atlas_generated_v1_group.yaml")
		Expect(len(objs)).To(Equal(1))
		group := objs[0].(*generatedv1.Group)
		applyTestParamsToGroup(group, params)
		Expect(kubeClient.Create(ctx, group)).To(Succeed())

		Eventually(func(g Gomega) {
			g.Expect(resources.CheckResourceReady(ctx, kubeClient, group)).To(Succeed())
		}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
		Expect(group.Status.V20250312).NotTo(BeNil())
		Expect(group.Status.V20250312.Id).NotTo(BeNil())
		return group
	}

	// deleteKubeObject deletes obj from Kubernetes and waits until it is gone.
	deleteKubeObject := func(obj client.Object) {
		Expect(kubeClient.Delete(ctx, obj)).To(Succeed())
		Eventually(func(g Gomega) {
			err := kubeClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
			g.Expect(apierrors.IsNotFound(err)).To(BeTrue())
		}).WithContext(ctx).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
	}

	// expectAtlasMatchesSpec GETs the container from Atlas and checks every field of the submitted entry came back as requested.
	// The expected value is the entry the test built, not the re-read spec, so values the operator writes back cannot mask a mismatch.
	expectAtlasMatchesSpec := func(groupID string, nc *generatedv1.NetworkContainer, submitted generatedv1.NetworkContainerSpecV20250312Entry) {
		Expect(nc.Spec.V20250312.Entry).To(Equal(&submitted), "operator must not rewrite the submitted spec")
		Expect(nc.Status.V20250312).NotTo(BeNil())
		Expect(nc.Status.V20250312.Id).NotTo(BeNil())
		atlasContainer, _, err := atlasClient.NetworkPeeringAPI.GetGroupContainer(ctx, groupID, *nc.Status.V20250312.Id).Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(specmatch.AtlasMatchesSpec(submitted, atlasContainer)).To(Succeed())
	}

	DescribeTable("Should create, match Atlas, update and delete a NetworkContainer",
		func(entry, updated generatedv1.NetworkContainerSpecV20250312Entry) {
			group := createGroup()
			groupID := *group.Status.V20250312.Id

			var nc *generatedv1.NetworkContainer
			By("Create NetworkContainer via groupRef and check Atlas matches the spec", func() {
				nc = newNetworkContainer(testNamespace.Name, group.GetName(), entry)
				Expect(kubeClient.Create(ctx, nc)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceReady(ctx, kubeClient, nc)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				expectAtlasMatchesSpec(groupID, nc, entry)
			})

			By("Update NetworkContainer and check Atlas matches the new spec", func() {
				Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(nc), nc)).To(Succeed())
				nc.Spec.V20250312.Entry = new(updated)
				Expect(kubeClient.Update(ctx, nc)).To(Succeed())

				Eventually(func(g Gomega) {
					g.Expect(resources.CheckResourceUpdated(ctx, kubeClient, nc)).To(Succeed())
				}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())
				expectAtlasMatchesSpec(groupID, nc, updated)
			})

			By("Delete NetworkContainer - should delete from Atlas", func() {
				containerID := *nc.Status.V20250312.Id
				deleteKubeObject(nc)

				Eventually(func(g Gomega) {
					_, r, err := atlasClient.NetworkPeeringAPI.GetGroupContainer(ctx, groupID, containerID).Execute()
					g.Expect(err).To(HaveOccurred())
					g.Expect(httputil.StatusCode(r)).To(Equal(http.StatusNotFound))
				}).WithContext(ctx).WithTimeout(time.Minute).WithPolling(5 * time.Second).Should(Succeed())
			})

			By("Delete prerequisite Group", func() {
				deleteKubeObject(group)
			})
		},
		Entry("AWS", Label("focus-networkcontainer-aws"),
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("AWS"), RegionName: new("US_EAST_1"), AtlasCidrBlock: new("10.128.0.0/21")},
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("AWS"), RegionName: new("US_EAST_1"), AtlasCidrBlock: new("10.129.0.0/21")},
		),
		Entry("Azure", Label("focus-networkcontainer-azure"),
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("AZURE"), Region: new("US_EAST_2"), AtlasCidrBlock: new("10.128.0.0/21")},
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("AZURE"), Region: new("US_EAST_2"), AtlasCidrBlock: new("10.129.0.0/21")},
		),
		Entry("GCP", Label("focus-networkcontainer-gcp"),
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("GCP"), AtlasCidrBlock: new("10.128.0.0/18")},
			generatedv1.NetworkContainerSpecV20250312Entry{ProviderName: new("GCP"), AtlasCidrBlock: new("10.129.0.0/18")},
		),
	)

	It("Should NOT delete from Atlas when ResourcePolicyKeep annotation is set", Label("focus-networkcontainer-kept"), func() {
		group := createGroup()
		groupID := *group.Status.V20250312.Id

		nc := newNetworkContainer(testNamespace.Name, group.GetName(), awsContainerEntry())
		nc.SetAnnotations(map[string]string{
			customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
		})
		Expect(kubeClient.Create(ctx, nc)).To(Succeed())
		Eventually(func(g Gomega) {
			g.Expect(resources.CheckResourceReady(ctx, kubeClient, nc)).To(Succeed())
		}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

		containerID := *nc.Status.V20250312.Id
		deleteKubeObject(nc)

		_, _, err := atlasClient.NetworkPeeringAPI.GetGroupContainer(ctx, groupID, containerID).Execute()
		Expect(err).ToNot(HaveOccurred(), "container must survive deletion of the Kubernetes object")

		_, err = atlasClient.NetworkPeeringAPI.DeleteGroupContainer(ctx, groupID, containerID).Execute()
		Expect(err).ToNot(HaveOccurred())
		deleteKubeObject(group)
	})

	It("Should import an existing Atlas container using the external-id annotation", Label("focus-networkcontainer-import"), func() {
		group := createGroup()
		groupID := *group.Status.V20250312.Id

		entry := awsContainerEntry()
		created, _, err := atlasClient.NetworkPeeringAPI.CreateGroupContainer(ctx, groupID, &admin.CloudProviderContainer{
			ProviderName:   entry.ProviderName,
			RegionName:     entry.RegionName,
			AtlasCidrBlock: entry.AtlasCidrBlock,
		}).Execute()
		Expect(err).ToNot(HaveOccurred())

		nc := newNetworkContainer(testNamespace.Name, group.GetName(), entry)
		nc.SetAnnotations(map[string]string{"mongodb.com/external-id": created.GetId()})
		Expect(kubeClient.Create(ctx, nc)).To(Succeed())
		Eventually(func(g Gomega) {
			g.Expect(resources.CheckResourceReady(ctx, kubeClient, nc)).To(Succeed())
		}).WithContext(ctx).WithTimeout(5 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

		Expect(nc.Status.V20250312.Id).To(Equal(created.Id))
		expectAtlasMatchesSpec(groupID, nc, entry)

		deleteKubeObject(nc)
		deleteKubeObject(group)
	})

	It("Should fail if groupRef points to non-existent Group", Label("focus-networkcontainer-fail"), func() {
		nc := newNetworkContainer(testNamespace.Name, "non-existent-group", awsContainerEntry())
		Expect(kubeClient.Create(ctx, nc)).To(Succeed())

		Eventually(func(g Gomega) {
			g.Expect(kubeClient.Get(ctx, client.ObjectKeyFromObject(nc), nc)).To(Succeed())
			g.Expect(meta.FindStatusCondition(nc.GetConditions(), "Ready")).NotTo(BeNil())
		}).WithContext(ctx).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())
		Expect(meta.FindStatusCondition(nc.GetConditions(), "Ready").Reason).To(Equal("Error"))

		nc.SetFinalizers([]string{})
		Expect(kubeClient.Update(ctx, nc)).To(Succeed())
		deleteKubeObject(nc)
	})
})

func awsContainerEntry() generatedv1.NetworkContainerSpecV20250312Entry {
	return generatedv1.NetworkContainerSpecV20250312Entry{
		ProviderName:   new("AWS"),
		RegionName:     new("US_EAST_1"),
		AtlasCidrBlock: new("10.128.0.0/21"),
	}
}

func newNetworkContainer(namespace, groupRefName string, entry generatedv1.NetworkContainerSpecV20250312Entry) *generatedv1.NetworkContainer {
	objs := samples.MustLoadSampleObjects("atlas_generated_v1_networkcontainer_with_groupref.yaml")
	Expect(len(objs)).To(Equal(1))
	nc := objs[0].(*generatedv1.NetworkContainer)
	nc.SetNamespace(namespace)
	nc.SetName(fmt.Sprintf("nc-%s", rand.String(6)))
	nc.Spec.V20250312.GroupRef = &k8s.LocalReference{Name: groupRefName}
	nc.Spec.V20250312.GroupId = nil
	nc.Spec.V20250312.Entry = new(entry)
	return nc
}
