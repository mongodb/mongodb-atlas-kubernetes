package customresource

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AnnotationLastAppliedConfiguration = "mongodb.com/last-applied-configuration"
)

type OperatorChecker func(resource akov2.AtlasCustomResource) (bool, error)
type AtlasChecker func(resource akov2.AtlasCustomResource) (bool, error)

func IsOwner(resource akov2.AtlasCustomResource, protectionFlag bool, operatorChecker OperatorChecker, atlasChecker AtlasChecker) (bool, error) {
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

	existInAtlas, err := atlasChecker(resource)
	if err != nil {
		return false, err
	}

	return !existInAtlas, nil
}

func ApplyLastConfigApplied(ctx context.Context, resource akov2.AtlasCustomResource, k8sClient client.Client) error {
	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
	if err != nil {
		return err
	}

	js, err := json.Marshal(uObj["spec"])
	if err != nil {
		return err
	}

	// we just copied an akov2.AtlasCustomResource so it must be one
	resourceCopy := resource.DeepCopyObject().(akov2.AtlasCustomResource)

	annotations := resourceCopy.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[AnnotationLastAppliedConfiguration] = string(js)
	resourceCopy.SetAnnotations(annotations)
	err = k8sClient.Patch(ctx, resourceCopy, client.MergeFrom(resource))
	if err != nil {
		return err
	}

	// retains current behavior
	resource.SetAnnotations(resourceCopy.GetAnnotations())
	return nil
}

func IsResourceManagedByOperator(resource akov2.AtlasCustomResource) (bool, error) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return false, nil
	}

	_, ok := annotations[AnnotationLastAppliedConfiguration]

	return ok, nil
}
