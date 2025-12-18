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

package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/testparams"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

// CopySecretToNamespace copies a secret from one namespace to another.
// Returns the copied secret ready to be applied to the target namespace.
func CopySecretToNamespace(ctx context.Context, kubeClient client.Client, key client.ObjectKey, targetNamespace string) (*corev1.Secret, error) {
	secret := corev1.Secret{}
	if err := kubeClient.Get(ctx, key, &secret); err != nil {
		return nil, fmt.Errorf("failed to load original secret %v: %w", key, err)
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: targetNamespace,
			Labels:    secret.Labels,
		},
		Data: secret.Data,
	}, nil
}

// CopyCredentialsToNamespace copies the credentials secret from the operator namespace
// to the target namespace. The secret is applied with the specified field owner.
func CopyCredentialsToNamespace(ctx context.Context, kubeClient client.Client, credentialsName, operatorNamespace, targetNamespace string, fieldOwner client.FieldOwner) error {
	globalCredsKey := client.ObjectKey{
		Name:      credentialsName,
		Namespace: operatorNamespace,
	}
	credentialsSecret, err := CopySecretToNamespace(ctx, kubeClient, globalCredsKey, targetNamespace)
	if err != nil {
		return err
	}
	return kubeClient.Patch(ctx, credentialsSecret, client.Apply, client.ForceOwnership, fieldOwner)
}

// ApplyYAMLToNamespace applies YAML objects to a namespace after replacing placeholders.
// Returns the list of applied objects.
func ApplyYAMLToNamespace(ctx context.Context, kubeClient client.Client, yaml []byte, params *testparams.TestParams, namespace string, fieldOwner client.FieldOwner) ([]client.Object, error) {
	yamlStr := params.ReplaceYAML(string(yaml))
	objs := yml.MustParseObjects(strings.NewReader(yamlStr))
	for _, obj := range objs {
		obj.SetNamespace(namespace)
		if err := kubeClient.Patch(ctx, obj, client.Apply, client.ForceOwnership, fieldOwner); err != nil {
			return nil, fmt.Errorf("failed to apply object %s/%s: %w", obj.GetObjectKind().GroupVersionKind().Kind, obj.GetName(), err)
		}
	}
	return objs, nil
}

var (
	// ErrResourceNotReady indicates the resource is not in Ready state
	ErrResourceNotReady = errors.New("resource is not ready")
	// ErrResourceNotUpdated indicates the resource is not in Updated state
	ErrResourceNotUpdated = errors.New("resource is not updated")
	// ErrResourceNotDeleted indicates the resource still exists
	ErrResourceNotDeleted = errors.New("resource still exists")
)

// CheckResourceReady checks if a resource has Ready condition set to True.
// Returns nil if ready, ErrResourceNotReady if not ready, or an error if the resource cannot be fetched.
func CheckResourceReady(ctx context.Context, kubeClient client.Client, obj kube.ObjectWithStatus) error {
	key := client.ObjectKeyFromObject(obj)
	if err := kubeClient.Get(ctx, key, obj); err != nil {
		return fmt.Errorf("failed to get resource %s/%s: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	if condition := meta.FindStatusCondition(obj.GetConditions(), "Ready"); condition != nil {
		if condition.Status == metav1.ConditionTrue {
			return nil
		}
	}
	return ErrResourceNotReady
}

// CheckResourceUpdated checks if a resource is Ready and in Updated state.
// Returns nil if updated, ErrResourceNotUpdated if not updated, or an error if the resource cannot be fetched.
func CheckResourceUpdated(ctx context.Context, kubeClient client.Client, obj kube.ObjectWithStatus) error {
	key := client.ObjectKeyFromObject(obj)
	if err := kubeClient.Get(ctx, key, obj); err != nil {
		return fmt.Errorf("failed to get resource %s/%s: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	ready := false
	if condition := meta.FindStatusCondition(obj.GetConditions(), "Ready"); condition != nil {
		ready = (condition.Status == metav1.ConditionTrue)
	}
	if !ready {
		return ErrResourceNotUpdated
	}
	if condition := meta.FindStatusCondition(obj.GetConditions(), "State"); condition != nil {
		if state.ResourceState(condition.Reason) == state.StateUpdated {
			return nil
		}
	}
	return ErrResourceNotUpdated
}

// CheckResourceDeleted checks if a resource has been deleted from the cluster.
// Returns nil if deleted, ErrResourceNotDeleted if still exists, or an error if the check fails.
func CheckResourceDeleted(ctx context.Context, kubeClient client.Client, obj client.Object) error {
	key := client.ObjectKeyFromObject(obj)
	if err := kubeClient.Get(ctx, key, obj); err != nil {
		// Resource not found means it's deleted - success!
		if client.IgnoreNotFound(err) == nil {
			return nil
		}
		// Other errors are unexpected
		return fmt.Errorf("failed to check if resource %s/%s is deleted: %w", obj.GetNamespace(), obj.GetName(), err)
	}
	return ErrResourceNotDeleted
}
