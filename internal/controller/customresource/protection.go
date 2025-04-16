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

package customresource

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

const (
	AnnotationLastAppliedConfiguration = "mongodb.com/last-applied-configuration"
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

func PatchLastConfigApplied[S any](ctx context.Context, k8sClient client.Client, resource api.AtlasCustomResource, spec *S) error {
	if spec == nil {
		return fmt.Errorf("spec is nil")
	}

	js, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal %v: %w", spec, err)
	}

	resourceCopy := resource.DeepCopyObject().(api.AtlasCustomResource)
	annotations := resourceCopy.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[AnnotationLastAppliedConfiguration] = string(js)
	resourceCopy.SetAnnotations(annotations)

	err = k8sClient.Patch(ctx, resourceCopy, client.MergeFrom(resource))
	if err != nil {
		return fmt.Errorf("failed to patch resource: %w", err)
	}
	return nil
}

func ParseLastConfigApplied[S any](resource api.AtlasCustomResource) (*S, error) {
	var spec S
	lastAppliedJSON, ok := resource.GetAnnotations()[AnnotationLastAppliedConfiguration]
	if !ok {
		return nil, nil
	}

	err := json.Unmarshal([]byte(lastAppliedJSON), &spec)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON annotation value into a %T: %w", spec, err)
	}
	return &spec, nil
}

func IsResourceManagedByOperator(resource api.AtlasCustomResource) (bool, error) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		return false, nil
	}

	_, ok := annotations[AnnotationLastAppliedConfiguration]

	return ok, nil
}
