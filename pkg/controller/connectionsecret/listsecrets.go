package connectionsecret

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// ListByClusterName returns all secrets in the specified namespace that have labels for 'projectID' and 'clusterName'
func ListByClusterName(k8sClient client.Client, namespace, projectID, clusterName string) ([]corev1.Secret, error) {
	return list(k8sClient, namespace, projectID, clusterName, "")
}

// ListByUserName returns all secrets in the specified namespace that have label for 'projectID' and data for 'userName'
func ListByUserName(k8sClient client.Client, namespace, projectID, userName string) ([]corev1.Secret, error) {
	return list(k8sClient, namespace, projectID, "", userName)
}

func list(k8sClient client.Client, namespace, projectID, clusterName, dbUserName string) ([]corev1.Secret, error) {
	secrets := corev1.SecretList{}
	var result []corev1.Secret
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			CredLabelKey: CredLabelVal,
		}),
	}
	if err := k8sClient.List(context.Background(), &secrets, client.InNamespace(namespace), opts); err != nil {
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
