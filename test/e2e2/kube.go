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

package e2e2

import (
	"context"
	"fmt"
	"os"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

// InitK8sTest initializes a test enviroment on Kubernetes.
// It requires:
// - A running kubernetes cluster with a local configuration accessing it.
// - The given target CRD being installed in that cluster
// - Certain environment variables, such as OPERATOR_NAMESPACE, being set
// - The operator running on a given Pod or as a process with a given PID
func InitK8sTest(ctx context.Context, targetCRD string) (client.Client, error) {
	if err := assertRequiredEnvVars("OPERATOR_NAMESPACE"); err != nil {
		return nil, fmt.Errorf("missing required test Kubernetes env vars: %w", err)
	}

	kubeClient, err := TestKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to setup Kubernetes test env client: %w", err)
	}

	if err := assertCRD(ctx, kubeClient, targetCRD); err != nil {
		return nil, fmt.Errorf("failed to asert for test-required CRD: %w", err)
	}

	if err := assertOperator(ctx, kubeClient); err != nil {
		return nil, fmt.Errorf("failed to asert for test operator running: %w", err)
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

func assertRequiredEnvVars(envVars ...string) error {
	missing := make([]string, 0, len(envVars))
	for _, envVar := range envVars {
		_, ok := os.LookupEnv(envVar)
		if !ok {
			missing = append(missing, envVar)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %v", missing)
	}
	return nil
}

func assertOperator(ctx context.Context, kubeClient client.Client) error {
	pid := os.Getenv("OPERATOR_PID")
	if pid != "" {
		return assertProcessIsRunning(pid)
	}
	pod := os.Getenv("OPERATOR_POD_NAME")
	if pod != "" {
		ns := os.Getenv("OPERATOR_NAMESPACE")
		return assertPod(ctx, kubeClient, pod, ns)
	}
	return fmt.Errorf("please set OPERATOR_PID or OPERATOR_POD_NAME to allow to check he operator is running")
}

func assertProcessIsRunning(pidString string) error {
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return fmt.Errorf("failed to convert %s to a numeric PID: %w", pidString, err)
	}
	if _, err := os.FindProcess(pid); err != nil {
		return fmt.Errorf("failed to find process for PID %d: %w", pid, err)
	}
	return nil
}

func assertPod(ctx context.Context, kubeClient client.Client, pod, ns string) error {
	podObj := corev1.Pod{}
	key := client.ObjectKey{Name: pod, Namespace: ns}
	if err := kubeClient.Get(ctx, key, &podObj, &client.GetOptions{}); err != nil {
		return fmt.Errorf("failed to get POD %s/%s: %w", ns, pod, err)
	}
	return nil
}
