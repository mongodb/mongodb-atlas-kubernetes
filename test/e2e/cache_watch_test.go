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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/conditions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/data"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

var _ = Describe("Kubernetes cache watch test:", Label("cache-watch"), func() {
	DescribeTable("Cache gets", Label("watch-gets"),
		func(name string, namespaceIndexesToWatch []int, expectedNamespaceFound []bool) {
			testData := model.DataProvider(
				fmt.Sprintf("cache-%s", name),
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30008,
				[]func(*model.TestDataProvider){},
			)
			namespaces := setupNamespaces(testData, name, len(expectedNamespaceFound))
			defer clearNamespaces(testData, namespaces)
			setupSecrets(testData, namespaces)
			defer clearSecrets(testData, namespaces)

			mgrHandle := setupManager(testData, namespaces, namespaceIndexesToWatch)
			defer mgrHandle.stop()

			By("Using the manager cache to get all secrets in all namespaces and checking the expected results", func() {
				cache := mgrHandle.mgr.GetCache()
				for i, ns := range namespaces {
					secret := &corev1.Secret{}
					err := cache.Get(testData.Context, types.NamespacedName{Name: config.DefaultOperatorGlobalKey, Namespace: ns}, secret)
					if expectedNamespaceFound[i] {
						Expect(err).To(Succeed())
					} else {
						Expect(err).NotTo(Succeed())
						Expect(err.Error()).To(ContainSubstring("because of unknown namespace for the cache"))
					}
				}
			})
		},
		Entry("From all namespaces when no namespace config is set", Label("gets-all"),
			"gets-all", []int{}, []bool{true, true, true}),
		Entry("One namespace when only one namespace is configured", Label("gets-one"),
			"gets-one", []int{0}, []bool{true, false, false}),
		Entry("Two namespaces when only those two are configured", Label("gets-two"),
			"gets-two", []int{0, 1}, []bool{true, true, false}),
	)

	DescribeTable("Cache lists", Label("watch-lists"),
		func(name string, namespaceIndexesToWatch []int, expectedNamespaceFound []bool) {
			testData := model.DataProvider(
				fmt.Sprintf("cache-%s", name),
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30009,
				[]func(*model.TestDataProvider){},
			)
			namespaces := setupNamespaces(testData, name, len(expectedNamespaceFound))
			defer clearNamespaces(testData, namespaces)
			setupSecrets(testData, namespaces)
			defer clearSecrets(testData, namespaces)

			mgrHandle := setupManager(testData, namespaces, namespaceIndexesToWatch)
			defer mgrHandle.stop()

			By("Using the manager cache to list all secrets in all namespaces and checking the expected results", func() {
				cache := mgrHandle.mgr.GetCache()
				for i, ns := range namespaces {
					secrets := &corev1.SecretList{}
					err := cache.List(testData.Context, secrets, &client.ListOptions{Namespace: ns})
					if expectedNamespaceFound[i] {
						Expect(err).To(Succeed())
					} else {
						Expect(err).NotTo(Succeed())
						Expect(err.Error()).To(ContainSubstring("because of unknown namespace for the cache"))
					}
				}
			})
		},
		Entry("From all namespaces when no namespace config is set", Label("list-all"),
			"list-all", []int{}, []bool{true, true, true}),
		Entry("One namespace when only one namespace is configured", Label("list-one"),
			"list-one", []int{0}, []bool{true, false, false}),
		Entry("Two namespaces when only those two are configured", Label("list-two"),
			"list-two", []int{0, 1}, []bool{true, true, false}),
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
	DescribeTable("Reconciles",
		func(name string, namespaceIndexesToWatch []int, expectedNamespaceFound []bool) {
			testData := model.DataProvider(
				name,
				model.NewEmptyAtlasKeyType().UseDefaultFullAccess(),
				30010,
				[]func(*model.TestDataProvider){},
			)
			namespaces := setupNamespaces(testData, name, len(expectedNamespaceFound))
			defer clearNamespaces(testData, namespaces)
			setupSecrets(testData, namespaces)
			defer clearSecrets(testData, namespaces)

			mgrHandle := setupManager(testData, namespaces, namespaceIndexesToWatch)
			defer mgrHandle.stop()

			By("Launching an atlas project on each namespace, expect only listened namespaces to update status", func() {
				projectNames := make([]string, len(expectedNamespaceFound))
				for i, ns := range namespaces {
					project := data.DefaultProject()
					By("Create project", func() {
						project.Name = utils.RandomName(fmt.Sprintf("test-project-%s-ns%d", name, i))
						project.Namespace = ns
						project.Spec.Name = "" // Set empty name to force validation error
						projectNames[i] = project.Name

						Expect(testData.K8SClient.Create(context.Background(), project)).ToNot(HaveOccurred())
					})

					prj := &akov2.AtlasProject{}
					By("Wait project to be present in Kubernetes", func() {
						Eventually(func(g Gomega) bool {
							return g.Expect(
								testData.K8SClient.Get(testData.Context,
									types.NamespacedName{Name: project.Name, Namespace: ns}, prj),
							).To(Succeed())
						}).WithTimeout(time.Minute).Should(BeTrue())
					})

					if expectedNamespaceFound[i] {
						expectedCondition := api.Condition{
							Type:    "ProjectReady",
							Status:  "False",
							Reason:  "ProjectNotCreatedInAtlas",
							Message: "projectName is invalid because must be set",
						}
						By("Verify Kubernetes status got the expected error", func() {
							Eventually(func(g Gomega) bool {
								g.Expect(
									testData.K8SClient.Get(testData.Context,
										types.NamespacedName{Name: project.Name, Namespace: ns}, prj),
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
									testData.K8SClient.Get(testData.Context,
										types.NamespacedName{Name: project.Name, Namespace: ns}, prj),
								).To(Succeed())
								if prj.Status.Common.ObservedGeneration == expectedObservedGeneration {
									verifications += 1

								}
								return verifications > 15
							}).WithPolling(time.Second).WithTimeout(40 * time.Second).Should(BeTrue())
						})
					}
					By("Delete project", func() {
						Expect(testData.K8SClient.Delete(context.Background(), project)).ToNot(HaveOccurred())
					})
				}
			})
		},
		Entry("All namespaces when no namespace config is set", Label("reconcile-all"),
			"reconcile-all", []int{}, []bool{true, true, true}),
		Entry("One namespace when only one namespace is configured", Label("reconcile-one"),
			"reconcile-one", []int{0}, []bool{true, false, false}),
		Entry("Two namespaces when only those two are configured", Label("reconcile-two"),
			"reconcile-two", []int{0, 1}, []bool{true, true, false}),
	)
})

func setupNamespaces(testData *model.TestDataProvider, name string, size int) []string {
	namespaces := make([]string, size)
	By("Setting up test namespaces (and wait for them to be created)", func() {
		for i := 0; i < size; i++ {
			namespaces[i] = utils.RandomName(fmt.Sprintf("%s-ns%d", name, i))
			Expect(k8s.CreateNamespace(testData.Context, testData.K8SClient, namespaces[i])).To(Succeed())
		}
		for _, ns := range namespaces {
			Eventually(func(g Gomega) bool {
				namespace := &corev1.Namespace{}
				return g.Expect(
					testData.K8SClient.Get(testData.Context, types.NamespacedName{Name: ns}, namespace),
				).To(Succeed())
			}).WithTimeout(time.Minute).Should(BeTrue())
		}
	})
	return namespaces
}

func setupSecrets(testData *model.TestDataProvider, namespaces []string) {
	By("Setting up test secrets, one on each namespace (and wait for them)", func() {
		for _, ns := range namespaces {
			k8s.CreateDefaultSecret(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, ns)
		}
		for _, ns := range namespaces {
			Eventually(func(g Gomega) bool {
				secret := &corev1.Secret{}
				return g.Expect(
					testData.K8SClient.Get(testData.Context, types.NamespacedName{
						Name: config.DefaultOperatorGlobalKey, Namespace: ns}, secret),
				).To(Succeed())
			}).WithTimeout(time.Minute).Should(BeTrue())
		}
	})
}

func setupManager(testData *model.TestDataProvider, namespaces []string, namespaceIndexesToWatch []int) *managerHandle {
	mgrHandle := managerHandle{}
	mgrCtx, cancelFn := context.WithCancel(context.Background())
	mgrHandle.cancelFn = cancelFn
	By("Setting up the manager in the first namespace to watch on the given namespaces", func() {
		managerConfig := &k8s.Config{
			GlobalAPISecret: client.ObjectKey{
				Namespace: namespaces[0],
				Name:      config.DefaultOperatorGlobalKey,
			},
			FeatureFlags: featureflags.NewFeatureFlags(func() []string { return []string{} }),
		}
		watchedNamespaces := map[string]bool{}
		for _, idx := range namespaceIndexesToWatch {
			watchedNamespaces[namespaces[idx]] = true
		}
		if len(watchedNamespaces) > 0 {
			managerConfig.WatchedNamespaces = watchedNamespaces
		}
		var err error
		mgrHandle.mgr, err = k8s.BuildManager(managerConfig)
		Expect(err).NotTo(HaveOccurred())

		mgrHandle.wg.Add(1)
		go func(ctx context.Context) {
			defer mgrHandle.wg.Done()
			err := mgrHandle.mgr.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
		}(mgrCtx)
	})

	By("wait for cache to be started", func() {
		// the first namespace should always be cache-accessible
		cache := mgrHandle.mgr.GetCache()
		Eventually(func(g Gomega) bool {
			namespace := &corev1.Namespace{}
			return g.Expect(
				cache.Get(testData.Context, types.NamespacedName{Name: namespaces[0]}, namespace),
			).To(Succeed())
		}).WithTimeout(time.Minute).Should(BeTrue())
	})
	return &mgrHandle
}

type managerHandle struct {
	mgr      manager.Manager
	wg       sync.WaitGroup
	cancelFn context.CancelFunc
}

func (mgrHandle *managerHandle) stop() {
	mgrHandle.cancelFn()
	mgrHandle.wg.Wait()
}

func clearSecrets(testData *model.TestDataProvider, namespaces []string) {
	By("Clearing the secrets", func() {
		for _, ns := range namespaces {
			Expect(k8s.DeleteKey(testData.Context, testData.K8SClient, config.DefaultOperatorGlobalKey, ns)).To(Succeed())
		}
	})
}

func clearNamespaces(testData *model.TestDataProvider, namespaces []string) {
	By("Clearing namespaces", func() {
		for _, ns := range namespaces {
			Expect(k8s.DeleteNamespace(testData.Context, testData.K8SClient, ns)).To(Succeed())
		}
	})
}
