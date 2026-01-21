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

package state

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const FieldOwner = "mongodb-atlas-kubernetes-resource-handler"

type Patcher struct {
	patchedObj                   *unstructured.Unstructured
	obj                          client.Object
	statusChanged, objectChanged bool
	fieldOwner                   string
	err                          error
}

func NewPatcher(obj client.Object) *Patcher {
	patchedObj := &unstructured.Unstructured{}
	patchedObj.SetAPIVersion(obj.GetObjectKind().GroupVersionKind().GroupVersion().String())
	patchedObj.SetKind(obj.GetObjectKind().GroupVersionKind().Kind)
	patchedObj.SetName(obj.GetName())
	patchedObj.SetNamespace(obj.GetNamespace())
	patchedObj.SetGeneration(obj.GetGeneration())
	// Copy ResourceVersion to ensure Server-Side Apply works correctly with fake client
	// In controller-runtime 0.23.0+, SSA requires ResourceVersion to be present and correct
	if obj.GetResourceVersion() != "" {
		patchedObj.SetResourceVersion(obj.GetResourceVersion())
	}
	return &Patcher{patchedObj: patchedObj, obj: obj, fieldOwner: FieldOwner}
}

// UpdateStateTracker updates the state tracker annotation on the given object.
func (p *Patcher) UpdateStateTracker(dependencies ...client.Object) *Patcher {
	if p.err != nil {
		return p
	}

	stateTracker := ComputeStateTracker(p.obj, dependencies...)

	annotations := p.patchedObj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[AnnotationStateTracker] = stateTracker
	p.patchedObj.SetAnnotations(annotations)

	p.objectChanged = true
	return p
}

// UpdateStatus updates the status of the given object.
//
// Note: this method omits the "conditions" field from the status. To update conditions, use UpdateConditions().
func (p *Patcher) UpdateStatus() *Patcher {
	if p.err != nil {
		return p
	}

	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(p.obj)
	if err != nil {
		p.err = err
		return p
	}

	if err := unstructured.SetNestedField(p.patchedObj.Object, content["status"], "status"); err != nil {
		p.err = err
		return p
	}

	unstructured.RemoveNestedField(p.patchedObj.Object, "status", "conditions")

	p.statusChanged = true
	return p
}

// WithFieldOwner sets the field owner for the patch operation.
func (p *Patcher) WithFieldOwner(fieldOwner string) *Patcher {
	p.fieldOwner = fieldOwner
	return p
}

// UpdateConditions updates the status conditions of the given object.
//
// Note: this method only updates the "conditions" field in the status. To update remaining status fields, use UpdateStatus().
func (p *Patcher) UpdateConditions(conditions []metav1.Condition) *Patcher {
	if p.err != nil {
		return p
	}

	content := make([]interface{}, 0, len(conditions))
	for i := range conditions {
		c, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&conditions[i])
		if err != nil {
			p.err = err
			return p
		}
		content = append(content, c)
	}

	unstructured.RemoveNestedField(p.patchedObj.Object, "status")

	if err := unstructured.SetNestedSlice(p.patchedObj.Object, content, "status", "conditions"); err != nil {
		p.err = err
		return p
	}

	p.statusChanged = true
	return p
}

func (p *Patcher) patchObject(ctx context.Context, c client.Client) {
	if p.err != nil || !p.objectChanged {
		return
	}

	patchedCopy, err := p.copyPatchedObject(c)
	if err != nil {
		p.err = err
		return
	}

	applyConfig := client.ApplyConfigurationFromUnstructured(patchedCopy)
	err = c.Apply(ctx, applyConfig, client.FieldOwner(p.fieldOwner), client.ForceOwnership)
	p.err = err
}

func (p *Patcher) patchStatus(ctx context.Context, c client.Client) {
	if p.err != nil || !p.statusChanged {
		return
	}

	patchedCopy, err := p.copyPatchedObject(c)
	if err != nil {
		p.err = err
		return
	}

	// SSA Apply() method for sub-resources is not yet supported, so we use Patch here.
	// See the following issue for more details: https://github.com/kubernetes-sigs/controller-runtime/issues/3183
	err = c.Status().Patch(ctx, patchedCopy, client.Apply, client.FieldOwner(p.fieldOwner), client.ForceOwnership)
	p.err = err

	// After successful status patch, fetch the updated ResourceVersion from the client
	// This is necessary because the fake client in 0.23.0+ updates ResourceVersion after patching,
	// and we need it for subsequent object patch operations
	if err == nil {
		key := client.ObjectKeyFromObject(p.obj)
		currentObj := p.obj.DeepCopyObject().(client.Object)
		if fetchErr := c.Get(ctx, key, currentObj); fetchErr == nil {
			if rv := currentObj.GetResourceVersion(); rv != "" {
				p.patchedObj.SetResourceVersion(rv)
			}
		}
	}
}

func (p *Patcher) copyPatchedObject(c client.Client) (*unstructured.Unstructured, error) {
	patchedCopy := p.patchedObj.DeepCopy()
	if patchedCopy.GetObjectKind().GroupVersionKind().Empty() {
		gvk, err := apiutil.GVKForObject(p.obj, c.Scheme())
		if err != nil {
			return nil, err
		}
		patchedCopy.SetAPIVersion(gvk.GroupVersion().String())
		patchedCopy.SetKind(gvk.Kind)
	}
	return patchedCopy, nil
}

// Patch applies the patches to the given object and updates both status and the annotations if they were modified.
func (p *Patcher) Patch(ctx context.Context, c client.Client) error {
	if p.err != nil {
		return p.err
	}

	p.patchStatus(ctx, c)
	p.patchObject(ctx, c)

	return p.err
}
