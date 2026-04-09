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

package dryrun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	dryRunComponent   = "DryRun Manager"
	DryRunInstance    = "mongodb.com/dry-run-instance"
	DryRunFinishedMsg = "finished"
	DryRunReason      = "DryRun"
)

type reconciler interface {
	reconcile.Reconciler
	For() (client.Object, builder.Predicates)
}

type terminationAwareReconciler struct {
	reconciler
}

func (t *terminationAwareReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	clearTerminationErrors()
	result, err := t.reconciler.Reconcile(ctx, req)
	if err != nil {
		return result, err
	}
	return result, terminationError()
}

// Manager is a controller-runtime runnable
// that acts similar to controller-runtime's Manager
// but executing dry-run functionality.
type Manager struct {
	cluster.Cluster
	reconcilers  []reconciler
	logger       *zap.Logger
	instanceUID  string
	eventsClient corev1client.EventsGetter
	namespaces   []string
}

func NewManager(c cluster.Cluster, eventsClient corev1client.EventsGetter, logger *zap.Logger, namespaces []string) (*Manager, error) {
	mgr := &Manager{
		Cluster:      c,
		logger:       logger.Named("dry-run-manager"),
		instanceUID:  uuid.New().String(),
		eventsClient: eventsClient,
		namespaces:   []string{metav1.NamespaceAll},
	}

	if len(namespaces) > 0 {
		mgr.namespaces = namespaces
	}

	return mgr, nil
}

func (m *Manager) SetupReconciler(r reconciler) {
	m.reconcilers = append(m.reconcilers, &terminationAwareReconciler{reconciler: r})
}

//nolint:unparam
func (m *Manager) eventf(ctx context.Context, object runtime.Object, eventType, reason, messageFmt string, args ...any) error {
	ref, err := reference.GetReference(m.Cluster.GetScheme(), object)
	if err != nil {
		return fmt.Errorf("unable to get reference from object: %w", err)
	}

	t := metav1.Time{Time: time.Now()}
	namespace := ref.Namespace
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}

	ev := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.%x", ref.Name, t.UnixNano()),
			Namespace: namespace,
			// Note: labels are not encouraged on events as filtering events by labels can overload API-server,
			// see https://github.com/kubernetes/kubernetes/pull/115058.
			Annotations: map[string]string{
				DryRunInstance: m.instanceUID,
			},
		},
		InvolvedObject:      *ref,
		Reason:              reason,
		Message:             fmt.Sprintf(messageFmt, args...),
		FirstTimestamp:      t,
		LastTimestamp:       t,
		Count:               1,
		Type:                eventType,
		ReportingController: dryRunComponent,
		Source: corev1.EventSource{
			Component: dryRunComponent,
		},
	}

	_, err = m.eventsClient.Events(ev.GetNamespace()).Create(ctx, ev, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to create event: %w", err)
	}

	return nil
}

func (m *Manager) executeDryRun(ctx context.Context) error {
	if err := m.dryRunReconcilers(ctx); err != nil {
		return err
	}

	if err := m.eventf(ctx, m.object(), corev1.EventTypeNormal, DryRunReason, DryRunFinishedMsg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) dryRunReconcilers(ctx context.Context) error {
	enableErrors()

	if !m.Cluster.GetCache().WaitForCacheSync(ctx) {
		return errors.New("cluster cache sync failed")
	}

	for _, reconciler := range m.reconcilers {
		originalResource, _ := reconciler.For()
		resource := originalResource.DeepCopyObject() // don't mutate the prototype

		// build GVK
		if resource.GetObjectKind().GroupVersionKind().Empty() {
			if err := buildGVK(m.Cluster, resource); err != nil {
				return err
			}
		}

		gvk := resource.GetObjectKind().GroupVersionKind()
		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(schema.GroupVersionKind{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind + "List"})

		for _, namespace := range m.namespaces {
			if err := m.Cluster.GetClient().List(ctx, list, client.InNamespace(namespace)); err != nil {
				return fmt.Errorf("unable to list resources: %w", err)
			}

			for _, item := range list.Items {
				req := reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&item)}
				_, err := reconciler.Reconcile(ctx, req)
				if err != nil {
					if err := m.reportError(ctx, &item, err); err != nil {
						return err
					}
				}
				if err := m.eventf(ctx, &item, corev1.EventTypeNormal, DryRunReason, "done"); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func buildGVK(c cluster.Cluster, resource runtime.Object) error {
	gvks, _, err := c.GetScheme().ObjectKinds(resource)
	if err != nil {
		return fmt.Errorf("unable to determine GVK for resource %T: %w", resource, err)
	}
	if len(gvks) == 0 {
		return fmt.Errorf("no GVKs present for resource %T", resource)
	}
	objectKind, ok := resource.(schema.ObjectKind)
	if !ok {
		return fmt.Errorf("unable to set GVK for resource %T: %w", resource, err)
	}
	objectKind.SetGroupVersionKind(gvks[len(gvks)-1]) // set the latest version, it's what our local specs follow
	return nil
}

type unwrapError interface {
	Unwrap() error
}

type multiUnwrapError interface {
	Unwrap() []error
}

// reportDryRunErrors emits events for every instance of a DryRunError that are wrapped in the given error.
// Wrapped errors in Go are a tree that can be traversed recursively.
func (m *Manager) reportDryRunErrors(ctx context.Context, obj runtime.Object, err error) error {
	if err == nil {
		return nil
	}

	if unwrapped, ok := err.(unwrapError); ok {
		return m.reportDryRunErrors(ctx, obj, unwrapped.Unwrap())
	}

	if e, ok := err.(multiUnwrapError); ok {
		for _, unwrapped := range e.Unwrap() {
			if err := m.reportDryRunErrors(ctx, obj, unwrapped); err != nil {
				return err
			}
		}
		return nil
	}

	dryRunErr := &DryRunError{}
	if ok := errors.As(err, &dryRunErr); ok {
		return m.eventf(ctx, obj, corev1.EventTypeNormal, DryRunReason, "%s", dryRunErr.Msg)
	}

	// last resort: if some error in the error tree is not wrapped using either errors.Join() and/or fmt.Errorf("%w"...),
	// then detect a potential dry run error by the dry run error prefix.
	// this prevents false positives.
	if strings.Contains(err.Error(), dryRunErrorPrefix) {
		return m.eventf(ctx, obj, corev1.EventTypeNormal, DryRunReason, "%s", err.Error())
	}

	return nil
}

func (m *Manager) reportError(ctx context.Context, obj runtime.Object, err error) error {
	if err == nil {
		return nil
	}

	if containsDryRunErrors(err) {
		return m.reportDryRunErrors(ctx, obj, err)
	}

	m.logger.Error(err.Error())
	return m.eventf(ctx, obj, corev1.EventTypeWarning, DryRunReason, "%s", err.Error())
}

func (m *Manager) object() runtime.Object {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      os.Getenv("JOB_NAME"),
			Namespace: os.Getenv("JOB_NAMESPACE"),
		},
	}
}

// Start executes the dry-run and returns immediately.
// In contrast to controller-runtime's Manager it doesn't periodically reconcile
// but exits after the dry-run pass.
//
// This method blocks until the dry-run is complete.
func (m *Manager) Start(ctx context.Context) error {
	var (
		wg         sync.WaitGroup
		clusterErr error
	)

	cancelCtx, stopCluster := context.WithCancel(ctx)
	wg.Go(func() {
		// this blocks until it errors out or the context is canceled
		// where we instruct the Cluster to stop.
		if err := m.Cluster.Start(cancelCtx); err != nil {
			clusterErr = fmt.Errorf("cluster start failed: %w", err)
		}
	})

	err := m.executeDryRun(cancelCtx)
	if err != nil {
		stopCluster() // opportunistically stop the Cluster object if an error happened
		return err
	}

	stopCluster()
	wg.Wait()
	return clusterErr
}
