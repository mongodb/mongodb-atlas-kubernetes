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

package kube

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

const (
	Pause = time.Second
)

type ObjectWithStatus interface {
	client.Object
	GetConditions() []metav1.Condition
}

// NewK8sTest initializes a test environment on Kubernetes.
// It requires:
// - A running Kubernetes cluster with a local configuration bound to it.
// - The given set CRDs installed in that cluster
func NewK8sTest(ctx context.Context, crds ...string) (client.Client, error) {
	kubeClient, err := TestKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to setup Kubernetes test env client: %w", err)
	}

	for _, targetCRD := range crds {
		if err := assertCRD(ctx, kubeClient, targetCRD); err != nil {
			return nil, fmt.Errorf("failed to asert for test-required CRD: %w", err)
		}
	}
	return kubeClient, nil
}

// TestKubeClient returns a Kubernetes client for tests.
// It requires a running Kubernetes cluster and a local configuration to it.
// It supports core Kubernetes types, production and experimental CRDs.
func TestKubeClient() (client.Client, error) {
	testScheme, err := getTestScheme(
		corev1.AddToScheme,
		apiextensionsv1.AddToScheme,
		akov2.AddToScheme,
		akov2next.AddToScheme)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Kubernetes test env scheme: %w", err)
	}
	return getKubeClient(testScheme)
}

func Apply(ctx context.Context, kubeClient client.Client, objs ...client.Object) error {
	for i, obj := range objs {
		if err := apply(ctx, kubeClient, obj); err != nil {
			return fmt.Errorf("failed to apply object %d: %w", (i + 1), err)
		}
	}
	return nil
}

func apply(ctx context.Context, kubeClient client.Client, obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	old := obj.DeepCopyObject().(client.Object)
	err := kubeClient.Get(ctx, key, old)
	switch {
	case err == nil:
		obj = obj.DeepCopyObject().(client.Object)
		obj.SetResourceVersion(old.GetResourceVersion())
		if err := kubeClient.Update(ctx, obj); err != nil {
			return fmt.Errorf("failed to update %s: %w", key, err)
		}
	case apierrors.IsNotFound(err):
		if err := kubeClient.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create %s: %w", key, err)
		}
	default:
		return fmt.Errorf("failed to apply %s: %w", key, err)
	}
	return nil
}

func CreateNamespace(ctx context.Context, kubeClient client.Client, namespace string) error {
	return kubeClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
}

func WipeNamespace(ctx context.Context, kubeClient client.Client, namespace string) error {
	return kubeClient.Delete(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
}

func HasNamespace(ctx context.Context, kubeClient client.Client, namespace string) bool {
	return kubeClient.Get(ctx, client.ObjectKey{Name: namespace}, &corev1.Namespace{}) == nil
}

func AssertObjReady(ctx context.Context, kubeClient client.Client, key client.ObjectKey, obj ObjectWithStatus) (bool, error) {
	err := kubeClient.Get(ctx, key, obj)
	if err != nil {
		return false, fmt.Errorf("failed to get object %v: %w", key, err)
	}
	for _, condition := range obj.GetConditions() {
		if condition.Type == "Ready" && condition.Status == metav1.ConditionTrue {
			return true, nil
		}
	}
	return false, nil
}

func AssertObjExists(ctx context.Context, kubeClient client.Client, key client.ObjectKey, obj client.Object) (bool, error) {
	err := kubeClient.Get(ctx, key, obj)
	if err != nil {
		return false, fmt.Errorf("failed to get object %v: %w", key, err)
	}
	return true, nil
}

func SetNamespace(obj client.Object, ns string) client.Object {
	renamed := obj.DeepCopyObject().(client.Object)
	renamed.SetNamespace(ns)
	return renamed
}

func getTestScheme(addToSchemeFunctions ...func(*runtime.Scheme) error) (*runtime.Scheme, error) {
	testScheme := runtime.NewScheme()
	for _, addToSchemeFn := range addToSchemeFunctions {
		if err := addToSchemeFn(testScheme); err != nil {
			return nil, fmt.Errorf("failed to add to testScheme: %w", err)
		}
	}
	return testScheme, nil
}

func getKubeClient(scheme *runtime.Scheme) (client.Client, error) {
	restCfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config (is cluster configured?): %w", err)
	}
	kubeClient, err := client.New(restCfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes client (is cluster up?): %w", err)
	}
	return kubeClient, nil
}

func assertCRD(ctx context.Context, kubeClient client.Client, targetCRD string) error {
	crds := apiextensionsv1.CustomResourceDefinitionList{}
	if err := kubeClient.List(ctx, &crds, &client.ListOptions{}); err != nil {
		return fmt.Errorf("failed to list CRDs: %w", err)
	}
	for _, crd := range crds.Items {
		if crd.Name == targetCRD {
			return nil
		}
	}
	return fmt.Errorf("%s not found", targetCRD)
}
