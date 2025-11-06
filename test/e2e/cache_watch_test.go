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

package e2e_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

type tc struct {
	namespaces      []string
	watchNamespaces []string
	wantToFind      []string
}

var _ = Describe("Kubernetes cache watch test:", Label("cache-watch"), func() {
	DescribeTable("Cache gets", Label("focus-watch-gets"),
		func(ctx context.Context, testCase *tc) {
			testData := model.DataProvider(ctx, fmt.Sprintf("cache-%s", CurrentSpecReport().LeafNodeLabels[0]), model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30008, []func(*model.TestDataProvider){})
			namespaces := setupNamespaces(ctx, testData, testCase.namespaces...)
			defer clearNamespaces(ctx, testData, namespaces)
			setupSecrets(ctx, testData, namespaces)
			defer clearSecrets(ctx, testData, namespaces)

			c, stop := setupCluster(ctx, namespaces, testCase.watchNamespaces)
			defer stop()

			wantToFindSet := sets.NewString(testCase.wantToFind...)
			By("Using the manager cache to get all secrets in all namespaces and checking the expected results", func() {
				cache := c.GetCache()
				for _, ns := range namespaces {
					err := cache.Get(ctx, types.NamespacedName{Name: config.DefaultOperatorGlobalKey, Namespace: ns.GetName()}, &corev1.Secret{})

					if wantToFindSet.Has(ns.GetGenerateName()) {
						Expect(err).To(Succeed())
					} else {
						Expect(err).NotTo(Succeed())
						Expect(err.Error()).To(ContainSubstring("because of unknown namespace for the cache"))
					}
				}
			})
		},

		Entry("From all namespaces when no namespace config is set", Label("focus-gets-all"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: nil,
			wantToFind:      []string{"ns1", "ns2", "ns3"},
		}),

		Entry("One namespace when only one namespace is configured", Label("focus-gets-one"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1"},
			wantToFind:      []string{"ns1"},
		}),

		Entry("Two namespaces when only those two are configured", Label("focus-gets-two"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1", "ns2"},
			wantToFind:      []string{"ns1", "ns2"},
		}),
	)

	DescribeTable("Cache lists", Label("focus-watch-lists"),
		func(ctx context.Context, testCase *tc) {
			testData := model.DataProvider(ctx, fmt.Sprintf("cache-%s", CurrentSpecReport().LeafNodeLabels[0]), model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30009, []func(*model.TestDataProvider){})

			namespaces := setupNamespaces(ctx, testData, testCase.namespaces...)
			defer clearNamespaces(ctx, testData, namespaces)
			setupSecrets(ctx, testData, namespaces)
			defer clearSecrets(ctx, testData, namespaces)

			c, stop := setupCluster(ctx, namespaces, testCase.watchNamespaces)
			defer stop()

			wantToFindSet := sets.NewString(testCase.wantToFind...)
			By("Using the manager cache to list all secrets in all namespaces and checking the expected results", func() {
				cache := c.GetCache()
				for _, ns := range namespaces {
					err := cache.List(ctx, &corev1.SecretList{}, &client.ListOptions{Namespace: ns.GetName()})

					if wantToFindSet.Has(ns.GetGenerateName()) {
						Expect(err).To(Succeed())
					} else {
						Expect(err).NotTo(Succeed())
						Expect(err.Error()).To(ContainSubstring("because of unknown namespace for the cache"))
					}
				}
			})
		},

		Entry("From all namespaces when no namespace config is set", Label("focus-list-all"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: nil,
			wantToFind:      []string{"ns1", "ns2", "ns3"},
		}),

		Entry("One namespace when only one namespace is configured", Label("focus-list-one"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1"},
			wantToFind:      []string{"ns1"},
		}),

		Entry("Two namespaces when only those two are configured", Label("focus-list-two"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1", "ns2"},
			wantToFind:      []string{"ns1", "ns2"},
		}),
	)
})

// Reconcile tests cannot be run all at once, there are races between:
//
// - A reader goroutine from the Kubernetes event watcher.
//  See https://github.com/kubernetes/client-go/blob/52e5651101edcb2ecd8463ffdc281053cf6e63d4/tools/record/event.go#L395
//
// - A writer AddToScheme on test client helper code.
//   See `test/helper/e2e/k8s.CreateNewClient()`

var _ = Describe("Reconciles test:", func() {
	DescribeTable("Reconciles", Ordered,
		func(ctx context.Context, testCase *tc) {
			testData := model.DataProvider(ctx, fmt.Sprintf("reconcile-%s", CurrentSpecReport().LeafNodeLabels[0]), model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 30010, []func(*model.TestDataProvider){})
			namespaces := setupNamespaces(ctx, testData, testCase.namespaces...)
			defer clearNamespaces(ctx, testData, namespaces)
			setupSecrets(ctx, testData, namespaces)
			defer clearSecrets(ctx, testData, namespaces)

			_, stop := setupCluster(ctx, namespaces, testCase.watchNamespaces)
			defer stop()

			By("Launching an atlas project on each namespace, expect only listened namespaces to update status", func() {
				projectNames := make([]string, len(namespaces))
				wantToFindSet := sets.NewString(testCase.wantToFind...)

				for i, ns := range namespaces {
					project := data.DefaultProject()
					By("Create project", func() {
						project.Name = ""
						project.GenerateName = "test-project-" + CurrentSpecReport().LeafNodeLabels[0]
						project.Namespace = ns.GetName()
						project.Spec.Name = "" // Set empty name to force validation error
						Expect(testData.K8SClient.Create(context.Background(), project)).ToNot(HaveOccurred())
						projectNames[i] = project.Name
					})

					prj := &akov2.AtlasProject{}
					By("Wait project to be present in Kubernetes", func() {
						Eventually(func(g Gomega) bool {
							return g.Expect(
								testData.K8SClient.Get(ctx,
									types.NamespacedName{Name: project.Name, Namespace: ns.GetName()}, prj),
							).To(Succeed())
						}).WithTimeout(time.Minute).Should(BeTrue())
					})

					if wantToFindSet.Has(ns.GetGenerateName()) {
						expectedCondition := api.Condition{
							Type:    "ProjectReady",
							Status:  "False",
							Reason:  "ProjectNotCreatedInAtlas",
							Message: "groupName is empty and must be specified",
						}
						By("Verify Kubernetes status got the expected error", func() {
							Eventually(func(g Gomega) bool {
								g.Expect(
									testData.K8SClient.Get(ctx,
										types.NamespacedName{Name: project.Name, Namespace: ns.GetName()}, prj),
								).To(Succeed())
								match, err := ContainElement(conditions.MatchCondition(expectedCondition)).Match(prj.GetStatus().GetConditions())
								g.Expect(err).To(Succeed())
								return match
							}).WithPolling(time.Second).WithTimeout(time.Minute).Should(BeTrue())
						})
					} else {
						expectedObservedGeneration := prj.Status.ObservedGeneration
						verifications := 0
						By("Verify Kubernetes status remains untouched", func() {
							Eventually(func(g Gomega) bool {
								g.Expect(
									testData.K8SClient.Get(ctx,
										types.NamespacedName{Name: project.Name, Namespace: ns.GetName()}, prj),
								).To(Succeed())
								if prj.Status.Common.ObservedGeneration == expectedObservedGeneration {
									verifications += 1
								}
								return verifications > 15
							}).WithPolling(time.Second).WithTimeout(40 * time.Second).Should(BeTrue())
						})
					}
					By("Delete project", func() {
						Expect(testData.K8SClient.Delete(ctx, project)).ToNot(HaveOccurred())
					})
				}
			})
		},

		Entry("All namespaces when no namespace config is set", Label("reconcile-all"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: nil,
			wantToFind:      []string{"ns1", "ns2", "ns3"},
		}),

		Entry("One namespace when only one namespace is configured", Label("reconcile-one"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1"},
			wantToFind:      []string{"ns1"},
		}),

		Entry("Two namespaces when only those two are configured", Label("reconcile-two"), &tc{
			namespaces:      []string{"ns1", "ns2", "ns3"},
			watchNamespaces: []string{"ns1", "ns2"},
			wantToFind:      []string{"ns1", "ns2"},
		}),
	)
})

func setupNamespaces(ctx context.Context, testData *model.TestDataProvider, generateNames ...string) []*corev1.Namespace {
	result := make([]*corev1.Namespace, 0, len(generateNames))
	By("Setting up test namespaces (and wait for them to be created)", func() {
		for i := 0; i < len(generateNames); i++ {
			namespace, err := k8s.CreateRandomNamespace(ctx, testData.K8SClient, generateNames[i])
			Expect(err).To(BeNil())
			result = append(result, namespace)
		}
		for _, ns := range result {
			Eventually(func(g Gomega) bool {
				namespace := &corev1.Namespace{}
				return g.Expect(
					testData.K8SClient.Get(ctx, types.NamespacedName{Name: ns.GetName()}, namespace),
				).To(Succeed())
			}).WithTimeout(time.Minute).Should(BeTrue())
		}
	})
	return result
}

func setupSecrets(ctx context.Context, testData *model.TestDataProvider, namespaces []*corev1.Namespace) {
	By("Setting up test secrets, one on each namespace (and wait for them)", func() {
		for _, ns := range namespaces {
			k8s.CreateDefaultSecret(ctx, testData.K8SClient, config.DefaultOperatorGlobalKey, ns.GetName())
		}
		for _, ns := range namespaces {
			Eventually(func(g Gomega) bool {
				secret := &corev1.Secret{}
				return g.Expect(
					testData.K8SClient.Get(ctx, types.NamespacedName{
						Name: config.DefaultOperatorGlobalKey, Namespace: ns.GetName(),
					}, secret),
				).To(Succeed())
			}).WithTimeout(time.Minute).Should(BeTrue())
		}
	})
}

type stopper func()

func setupCluster(ctx context.Context, namespaces []*corev1.Namespace, wantToWatch []string) (cluster.Cluster, stopper) {
	var (
		wg sync.WaitGroup
		c  cluster.Cluster
	)

	wantToWatchSet := sets.NewString(wantToWatch...)
	watchedNamespaces := make(map[string]bool)
	for _, ns := range namespaces {
		if wantToWatchSet.Has(ns.GetGenerateName()) {
			watchedNamespaces[ns.GetName()] = true
		}
	}

	mgrCtx, cancelFn := context.WithCancel(ctx)
	By("Setting up the manager in the first namespace to watch on the given namespaces", func() {
		managerConfig := &k8s.Config{
			GlobalAPISecret: client.ObjectKey{
				Namespace: namespaces[0].GetName(),
				Name:      config.DefaultOperatorGlobalKey,
			},
			FeatureFlags: featureflags.NewFeatureFlags(func() []string { return []string{} }),
		}

		if len(watchedNamespaces) > 0 {
			managerConfig.WatchedNamespaces = watchedNamespaces
		}

		var err error
		c, err = k8s.BuildCluster(managerConfig)
		Expect(err).NotTo(HaveOccurred())

		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			err := c.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
		}(mgrCtx)
	})

	By("wait for cache to be started", func() {
		// the first namespace should always be cache-accessible
		Eventually(func(g Gomega) bool {
			return g.Expect(
				c.GetCache().Get(ctx, types.NamespacedName{Name: namespaces[0].GetName()}, &corev1.Namespace{}),
			).To(Succeed())
		}).WithTimeout(time.Minute).Should(BeTrue())
	})

	stopper := func() {
		cancelFn()
		wg.Wait()
	}

	return c, stopper
}

func clearSecrets(ctx context.Context, testData *model.TestDataProvider, namespaces []*corev1.Namespace) {
	By("Clearing the secrets", func() {
		for _, ns := range namespaces {
			Expect(k8s.DeleteKey(ctx, testData.K8SClient, config.DefaultOperatorGlobalKey, ns.GetName())).To(Succeed())
		}
	})
}

func clearNamespaces(ctx context.Context, testData *model.TestDataProvider, namespaces []*corev1.Namespace) {
	By("Clearing namespaces", func() {
		for _, ns := range namespaces {
			Expect(k8s.DeleteNamespace(ctx, testData.K8SClient, ns.GetName())).To(Succeed())
		}
	})
}
