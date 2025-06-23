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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
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
func NewK8sTest(ctx context.Context, crds ...*apiextensionsv1.CustomResourceDefinition) (client.Client, error) {
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
	testScheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(testScheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(testScheme))
	utilruntime.Must(akov2.AddToScheme(testScheme))
	utilruntime.Must(appsv1.AddToScheme(testScheme))
	return getKubeClient(testScheme)
}

func WithRenamedNamespace(obj client.Object, ns string) client.Object {
	renamed := obj.DeepCopyObject().(client.Object)
	renamed.SetNamespace(ns)
	return renamed
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

func assertCRD(ctx context.Context, kubeClient client.Client, targetCRD *apiextensionsv1.CustomResourceDefinition) error {
	crds := apiextensionsv1.CustomResourceDefinitionList{}
	if err := kubeClient.List(ctx, &crds, &client.ListOptions{}); err != nil {
		return fmt.Errorf("failed to list CRDs: %w", err)
	}
	for _, crd := range crds.Items {
		if crd.Name == targetCRD.Name {
			return nil
		}
	}
	return fmt.Errorf("%s not found", targetCRD)
}
