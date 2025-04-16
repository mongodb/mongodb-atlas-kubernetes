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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

// ListByDeploymentName returns all secrets in the specified namespace that have labels for 'projectID' and 'clusterName'
func ListByDeploymentName(ctx context.Context, k8sClient client.Client, namespace, projectID, clusterName string) ([]corev1.Secret, error) {
	return list(ctx, k8sClient, namespace, projectID, clusterName, "")
}

// ListByUserName returns all secrets in the specified namespace that have label for 'projectID' and data for 'userName'
func ListByUserName(ctx context.Context, k8sClient client.Client, namespace, projectID, userName string) ([]corev1.Secret, error) {
	return list(ctx, k8sClient, namespace, projectID, "", userName)
}

func list(ctx context.Context, k8sClient client.Client, namespace, projectID, clusterName, dbUserName string) ([]corev1.Secret, error) {
	secrets := corev1.SecretList{}
	var result []corev1.Secret
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			TypeLabelKey: CredLabelVal,
		}),
	}

	if namespace != "" {
		opts.Namespace = namespace
	}

	if err := k8sClient.List(ctx, &secrets, opts); err != nil {
		return nil, err
	}

	for _, s := range secrets.Items {
		if value, ok := s.Labels[ProjectLabelKey]; !ok || value != projectID {
			continue
		}
		if _, ok := s.Labels[ClusterLabelKey]; !ok {
			continue
		}
		if clusterName != "" && s.Labels[ClusterLabelKey] == kube.NormalizeLabelValue(clusterName) {
			result = append(result, s)
		}
		if dbUserName != "" {
			var userName []byte
			var ok bool
			if userName, ok = s.Data[userNameKey]; !ok {
				return nil, fmt.Errorf("secret %v is broken: missing the mandatory field %s", s.Name, userNameKey)
			}
			if string(userName) == dbUserName {
				result = append(result, s)
			}
		}
	}
	return result, nil
}
