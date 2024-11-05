package contract

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
)

func mustCreateK8sClient() client.Client {
	client, err := k8s.CreateNewClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create a Kubernetes client: %v", err))
	}
	return client
}

func defaultNamespace(name string) client.Object {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
}

func waitForReadyStatus(c client.Client, resources []client.Object, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for resources to be ready")
		case <-ticker.C:
			allReady := true
			for _, obj := range resources {
				if err := c.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
					return fmt.Errorf("failed to get object: %w", err)
				}
				if !isReady(obj) {
					allReady = false
					break
				}
			}
			if allReady {
				return nil
			}
		}
	}
}

func isReady(obj client.Object) bool {
	if cr, ok := obj.(api.AtlasCustomResource); ok {
		for _, condition := range cr.GetStatus().GetConditions() {
			if condition.Type == "Ready" {
				return condition.Status == corev1.ConditionTrue
			}
		}
		return false
	}
	switch o := obj.(type) {
	case *corev1.Secret:
		return true
	default:
		panic(fmt.Sprintf("failed to get ready state from unsupported object of type %T", o))
	}
}

func k8sRecreate(ctx context.Context, k8sClient client.Client, obj client.Object) error {
	err := k8sClient.Create(ctx, obj)
	if err == nil {
		return nil
	}
	if !k8sErrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create object: %w", err)
	}
	if err := k8sClient.Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update object: %w", err)
	}
	return nil
}
