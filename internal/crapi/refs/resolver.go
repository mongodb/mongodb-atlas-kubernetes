// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package refs

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/crapi/unstructured"
)

type encodeDecodeFunc func(any) (any, error)

type resolver struct {
	decoders           map[string]encodeDecodeFunc
	encoders           map[string]encodeDecodeFunc
	optionalExpansions map[string]struct{} // A set is better for lookups
	kubeObjectRegistry map[string]func(map[string]any) (client.Object, error)
}

// nilInstance use when some type should not be automatically instantiated by translation
func nilInstance(_ map[string]any) (client.Object, error) {
	return nil, nil
}

func newReferenceResolver() *resolver {
	return &resolver{

		decoders: map[string]encodeDecodeFunc{
			"v1/secrets": func(in any) (any, error) {
				s, ok := in.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string for secret decoding, but got %T", in)
				}
				return secretDecode(s)
			},
		},

		encoders: map[string]encodeDecodeFunc{
			"v1/secrets": func(in any) (any, error) {
				s, ok := in.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string for secret encoding, but got %T", in)
				}
				return secretEncode(s), nil
			},
		},

		optionalExpansions: map[string]struct{}{"groupRef": {}},

		kubeObjectRegistry: map[string]func(map[string]any) (client.Object, error){
			"v1/secrets":                            newKubeObjectFactory[corev1.Secret](),
			"atlas.generated.mongodb.com/v1/groups": nilInstance,
		},
	}
}

func newKubeObjectFactory[T any, P PtrClientObj[T]]() func(map[string]any) (client.Object, error) {
	return func(unstructured map[string]any) (client.Object, error) {
		obj := new(T)
		initializedObj, err := initObject(obj, unstructured)
		if err != nil {
			return nil, err
		}
		// Assert the concrete pointer type (*P) to the interface type.
		// This is guaranteed to be safe because of the PtrClientObj constraint
		return any(initializedObj).(client.Object), nil
	}
}

func initObject[T any](obj *T, unstructuredObj map[string]any) (*T, error) {
	if unstructuredObj != nil {
		if err := unstructured.FromUnstructured(obj, unstructuredObj); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func unstructuredKubeObjectFor(refSolver *resolver, gvr string) (map[string]any, error) {
	objCopy, err := kubeObjectFor(refSolver, gvr)
	if err != nil {
		return nil, fmt.Errorf("failed to get unstructured kube object for GVR %q: %w", gvr, err)
	}
	return unstructured.ToUnstructured(objCopy)
}

func kubeObjectFor(refSolver *resolver, gvr string) (client.Object, error) {
	return initializedKubeObjectFor(refSolver, gvr, nil)
}

func initializedKubeObjectFor(refSolver *resolver, gvr string, unstructuredData map[string]any) (client.Object, error) {
	objFn, ok := refSolver.kubeObjectRegistry[gvr]
	if !ok {
		return nil, fmt.Errorf("unsupported kube object for GVR %q", gvr)
	}
	return objFn(unstructuredData)
}
