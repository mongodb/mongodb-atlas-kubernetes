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

package k8s

import (
	"bytes"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8scfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func GetPodLogsByDeployment(deploymentName, deploymentNS string, options corev1.PodLogOptions) ([]byte, error) {
	pods, err := GetAllDeploymentPods(deploymentName, deploymentNS)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}
	if len(pods) == 0 {
		return nil, fmt.Errorf("no pods found")
	}
	return GetPodLogs(options, deploymentNS, pods[0].Name)
}

func GetPodLogs(options corev1.PodLogOptions, ns string, podName string) ([]byte, error) {
	cfg, err := k8scfg.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get client set: %w", err)
	}

	podLogRequest := clientSet.CoreV1().
		Pods(ns).
		GetLogs(podName, &options)
	stream, err := podLogRequest.Stream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}
	defer stream.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stream: %w", err)
	}
	return buf.Bytes(), nil
}

func GetAllDeploymentPods(deploymentName, deploymentNS string) ([]corev1.Pod, error) {
	cfg, err := k8scfg.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get client set: %w", err)
	}

	deployment, err := clientSet.AppsV1().Deployments(deploymentNS).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	pods, err := clientSet.CoreV1().Pods(deploymentNS).List(context.Background(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	return pods.Items, nil
}

func GetDeployment(deploymentName, deploymentNS string) (*appsv1.Deployment, error) {
	cfg, err := k8scfg.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get client set: %w", err)
	}

	deployment, err := clientSet.AppsV1().Deployments(deploymentNS).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deployment, err
}
