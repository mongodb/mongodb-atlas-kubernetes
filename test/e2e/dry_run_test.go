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

package e2e

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/dryrun"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/atlas"
	e2e_config "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/operator"
)

type (
	waitFunc      func() ([]*corev1.Event, bool)
	predicateFunc func([]*corev1.Event) bool
)

var _ = When("running in dry run mode", Label("dry-run"), Ordered, func() {
	var (
		testData               *model.TestDataProvider
		projectID, projectName string
	)

	BeforeAll(func(ctx context.Context) {
		atlasClient := atlas.GetClientOrFail()
		By("creating a project in Atlas")

		testID := uuid.New().String()[0:6]
		projectName = fmt.Sprintf("dry-run-%s", testID)
		group, _, err := atlasClient.Client.ProjectsApi.CreateGroup(ctx, &admin.Group{
			Name:  projectName,
			OrgId: atlasClient.OrgID,
		}).Execute()
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx context.Context) {
			_, err := atlasClient.Client.ProjectsApi.DeleteGroup(ctx, group.GetId()).Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		projectID = group.GetId()
	})

	AfterEach(func() {
		GinkgoWriter.Write([]byte("\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		GinkgoWriter.Write([]byte("Dry run test\n"))
		GinkgoWriter.Write([]byte("Operator namespace: " + testData.Resources.Namespace + "\n"))
		GinkgoWriter.Write([]byte("===============================================\n"))
		By("Delete Resources", func() {
			actions.AfterEachFinalCleanup([]model.TestDataProvider{*testData})
		})
	})

	BeforeEach(func(ctx SpecContext) {
		By("setting up secrets and a namespace", func() {
			testData = model.DataProvider(ctx, "dry-run", model.NewEmptyAtlasKeyType().UseDefaultFullAccess(), 40000, nil)
			actions.CreateNamespaceAndSecrets(testData)
		})
	})

	It("would create a project in Atlas if it does not exist", func(ctx context.Context) {
		By("creating an AtlasProject resource")
		prj := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "dry-run",
				Namespace:    testData.Resources.Namespace,
			},
			Spec: akov2.AtlasProjectSpec{
				Name: uuid.New().String()[0:6],
			},
		}
		err := testData.K8SClient.Create(ctx, prj)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx context.Context) {
			err := testData.K8SClient.Delete(ctx, prj)
			Expect(err).NotTo(HaveOccurred())
		})

		StartDryRunUntil(ctx, testData.K8SClient, testData.Resources.Namespace,
			and(
				messageEquals("Would create (POST) /api/atlas/v2/groups"),
				messageEquals("finished"),
			))
	})

	It("would change a project in Atlas", func(ctx context.Context) {
		By("creating an AtlasProject resource including changes")
		prj := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "dry-run",
				Namespace:    testData.Resources.Namespace,
			},
			Spec: akov2.AtlasProjectSpec{
				Name: projectName,
				MaintenanceWindow: project.MaintenanceWindow{
					DayOfWeek: 1,
					HourOfDay: 2,
					AutoDefer: true,
					StartASAP: true,
					Defer:     false,
				},
				Auditing: &akov2.Auditing{
					AuditFilter: "foo",
					Enabled:     true,
				},
			},
		}
		err := testData.K8SClient.Create(ctx, prj)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx context.Context) {
			err := testData.K8SClient.Delete(ctx, prj)
			Expect(err).NotTo(HaveOccurred())
		})

		StartDryRunUntil(ctx, testData.K8SClient, testData.Resources.Namespace,
			and(
				messageEquals(fmt.Sprintf("Would update (PATCH) /api/atlas/v2/groups/%s/maintenanceWindow", projectID)),
				messageEquals("finished"),
			))
	})

	It("would create a deployment in Atlas", func(ctx context.Context) {
		By("creating an AtlasDeployment resource")
		deployment := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "dry-run",
				Namespace:    testData.Resources.Namespace,
			},
			Spec: akov2.AtlasDeploymentSpec{
				ProjectDualReference: akov2.ProjectDualReference{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: projectID,
					},
					ConnectionSecret: &api.LocalObjectReference{
						Name: e2e_config.DefaultOperatorGlobalKey,
					},
				},
				FlexSpec: &akov2.FlexSpec{
					Name: "foo",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
				},
			},
		}
		err := testData.K8SClient.Create(ctx, deployment)
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func(ctx context.Context) {
			err := testData.K8SClient.Delete(ctx, deployment)
			Expect(err).NotTo(HaveOccurred())
		})

		StartDryRunUntil(ctx, testData.K8SClient, testData.Resources.Namespace,
			and(
				messageEquals("finished"),
				messageEquals(fmt.Sprintf("Would create (POST) /api/atlas/v2/groups/%s/flexClusters", projectID)),
			))
	})
})

func StartDryRunUntil(ctx context.Context, kubeClient client.Client, namespace string, predicate predicateFunc) {
	By("collecting events")
	waitForEvents := dryRunEventsFunc(ctx, kubeClient, 5*time.Minute, predicate)

	By("starting the operator in dry-run mode")
	o := operator.NewOperator(operator.DefaultOperatorEnv(namespace), GinkgoWriter, GinkgoWriter,
		"--log-level=debug",
		"--dry-run=true",
		"--global-api-secret-name=mongodb-atlas-operator-api-key",
		`--atlas-domain=https://cloud-qa.mongodb.com`,
	)
	t := GinkgoT()
	o.Start(t)
	DeferCleanup(func() {
		o.Wait(t)
	})

	By("waiting for events")
	_, result := waitForEvents()
	Expect(result).To(BeTrue())
}

func and(predicates ...predicateFunc) predicateFunc {
	return func(events []*corev1.Event) bool {
		for _, p := range predicates {
			if !p(events) {
				return false
			}
		}
		return true
	}
}

func messageEquals(v string) predicateFunc {
	return func(events []*corev1.Event) bool {
		for _, ev := range events {
			if ev.Message == v {
				return true
			}
		}
		return false
	}
}

// dryRunEventsFunc starts watching and collecting Kubernetes *corev1.Event having reason=DryRunReason in a separate goroutine.
// Before it starts watching events, it lists all existing events with reason=DryRunReason and ignores those.
// This allows to collect new events emitted after this function has been invoked.
//
// It returns a function that blocks until the given predicate returns true or the given timeout is reached.
// It returns all collected events, including those where the predicate returned true.
// If the timeout is reached it returns all collected events until the timeout was reached.
//
// Note for consumers: Do NOT assume ANY ordering (i.e. chronological) of the event slice both in the predicate and returned events.
func dryRunEventsFunc(ctx context.Context, kubeClient client.Client, timeout time.Duration, predicate predicateFunc) waitFunc {
	observedEvents := []*corev1.Event{}
	result := false

	listOpts := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("reason", dryrun.DryRunReason),
	}

	observedDryRuns := sets.New[string]()
	currentEvents := &corev1.EventList{}
	err := kubeClient.List(ctx, currentEvents, listOpts)
	Expect(err).NotTo(HaveOccurred())
	for _, event := range currentEvents.Items {
		if instance, ok := event.GetAnnotations()[dryrun.DryRunInstance]; ok {
			observedDryRuns.Insert(instance)
		}
	}

	cfg, err := config.GetConfig()
	Expect(err).NotTo(HaveOccurred())
	watchClient, err := client.NewWithWatch(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred())

	watch, err := watchClient.Watch(ctx, &corev1.EventList{}, listOpts)
	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(func() {
		watch.Stop()
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		eventCh := watch.ResultChan()
		for {
			select {
			case event := <-eventCh:
				ev, ok := event.Object.(*corev1.Event)
				if !ok {
					// TODO: change this!!!
					fmt.Printf("event %T is not of type *corev1.Event: %v\n", event.Object, event.Object)
					continue
				}

				instance, ok := ev.GetAnnotations()[dryrun.DryRunInstance]
				if !ok || observedDryRuns.Has(instance) {
					continue
				}

				observedEvents = append(observedEvents, ev)
				resultCopy := make([]*corev1.Event, 0, len(observedEvents))
				for _, ev := range observedEvents {
					resultCopy = append(resultCopy, ev.DeepCopy())
				}

				if predicate(resultCopy) {
					result = true
					return
				}

			case <-timeoutCtx.Done():
				return
			}
		}
	}()

	return func() ([]*corev1.Event, bool) {
		wg.Wait()
		return observedEvents, result
	}
}
