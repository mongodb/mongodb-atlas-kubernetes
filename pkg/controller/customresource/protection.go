package customresource

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AnnotationLastAppliedConfiguration = "mongodb.com/last-applied-configuration"
)

type OperatorChecker func(resource mdbv1.AtlasCustomResource) (bool, error)
type AtlasChecker func(resource mdbv1.AtlasCustomResource) (bool, error)

func IsOwner(resource mdbv1.AtlasCustomResource, protectionFlag bool, operatorChecker OperatorChecker, atlasChecker AtlasChecker) (bool, error) {
	if !protectionFlag {
		return true, nil
	}

	wasManaged, err := operatorChecker(resource)
	if err != nil {
		return false, err
	}

	if wasManaged {
		return true, nil
	}

	managedByAtlas, err := atlasChecker(resource)
	if err != nil {
		return false, err
	}

	// Operator owns the resoruce if it is not managed by the Atlas
	return !managedByAtlas, nil
}

func IsResourceProtected(resource mdbv1.AtlasCustomResource, protectionFlag bool) bool {
	if policy, ok := resource.GetAnnotations()[ResourcePolicyAnnotation]; ok {
		return policy == ResourcePolicyKeep
	}

	return protectionFlag
}

func ApplyLastConfigApplied(ctx context.Context, resource mdbv1.AtlasCustomResource, k8sClient client.Client) error {
	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
	if err != nil {
		return err
	}

	js, err := json.Marshal(uObj["spec"])
	if err != nil {
		return err
	}

	annotations := resource.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[AnnotationLastAppliedConfiguration] = string(js)
	resource.SetAnnotations(annotations)

	return k8sClient.Update(ctx, resource, &client.UpdateOptions{})
}

func IsResourceManagedByOperator(resource mdbv1.AtlasCustomResource) (bool, error) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return false, nil
	}

	_, ok := annotations[AnnotationLastAppliedConfiguration]

	return ok, nil
}
