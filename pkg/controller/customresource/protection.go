package customresource

import (
	"context"
	"encoding/json"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	AnnotationLastAppliedConfiguration = "mongodb.com/last-applied-configuration"
	AnnotationLastSkippedConfiguration = "mongodb.com/last-skipped-configuration"
)

type OperatorChecker func(resource api.AtlasCustomResource) (bool, error)
type AtlasChecker func(resource api.AtlasCustomResource) (bool, error)

func IsOwner(resource api.AtlasCustomResource, protectionFlag bool, operatorChecker OperatorChecker, atlasChecker AtlasChecker) (bool, error) {
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

func ApplyLastConfigApplied(ctx context.Context, resource api.AtlasCustomResource, k8sClient client.Client) error {
	return applyLastSpec(ctx, resource, k8sClient, AnnotationLastAppliedConfiguration)
}

func ApplyLastConfigSkipped(ctx context.Context, resource api.AtlasCustomResource, k8sClient client.Client) error {
	return applyLastSpec(ctx, resource, k8sClient, AnnotationLastSkippedConfiguration)
}

func applyLastSpec(ctx context.Context, resource api.AtlasCustomResource, k8sClient client.Client, annotationKey string) error {
	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
	if err != nil {
		return err
	}

	js, err := json.Marshal(uObj["spec"])
	if err != nil {
		return err
	}

	// we just copied an api.AtlasCustomResource so it must be one
	resourceCopy := resource.DeepCopyObject().(api.AtlasCustomResource)

	annotations := resourceCopy.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[annotationKey] = string(js)
	resourceCopy.SetAnnotations(annotations)
	err = k8sClient.Patch(ctx, resourceCopy, client.MergeFrom(resource))
	if err != nil {
		return err
	}

	// retains current behavior
	resource.SetAnnotations(resourceCopy.GetAnnotations())
	return nil
}

func IsResourceManagedByOperator(resource api.AtlasCustomResource) (bool, error) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return false, nil
	}

	_, ok := annotations[AnnotationLastAppliedConfiguration]

	return ok, nil
}
