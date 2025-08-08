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

package connectionsecret

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
)

const ConnectionSecretsEnsuredEvent = "ConnectionSecretsEnsured"

func ReapOrphanConnectionSecrets(ctx context.Context, k8sClient client.Client, projectID, namespace string, projectDeploymentNames []string) ([]string, error) {
	secretList := &corev1.SecretList{}
	labelSelector := labels.SelectorFromSet(labels.Set{TypeLabelKey: CredLabelVal, ProjectLabelKey: projectID})
	err := k8sClient.List(context.Background(), secretList, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed listing possible orphan secrets: %w", err)
	}

	removedOrphanSecrets := []string{}
	for _, secret := range secretList.Items {
		clusterName, ok := secret.Labels[ClusterLabelKey]
		if !ok {
			continue
		}
		if clusterExists := stringutil.Contains(projectDeploymentNames, clusterName); clusterExists {
			continue
		}
		if err := k8sClient.Delete(ctx, &secret); err != nil {
			return nil, fmt.Errorf("failed to remove orphan connection Secret: %w", err)
		} else {
			removedOrphanSecrets = append(removedOrphanSecrets, fmt.Sprintf("%s/%s", namespace, secret.Name))
		}
	}
	return removedOrphanSecrets, nil
}
