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

package refs

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// kubeset holds to the main Kubernetes object being translated,
// and related existing & added Kubernetes dependencies
type kubeset struct {
	scheme *runtime.Scheme
	main   client.Object
	m      map[client.ObjectKey]client.Object
	added  []client.Object
}

func newKubeset(scheme *runtime.Scheme, main client.Object, deps []client.Object) *kubeset {
	m := map[client.ObjectKey]client.Object{}
	for _, obj := range deps {
		m[client.ObjectKeyFromObject(obj)] = obj
	}
	return &kubeset{scheme: scheme, main: main, m: m}
}

func (mc *kubeset) find(name string) client.Object {
	key := client.ObjectKey{Name: name, Namespace: mc.main.GetNamespace()}
	return mc.m[key]
}

func (mc *kubeset) has(name string) bool {
	return mc.find(name) != nil
}

func (mc *kubeset) add(obj client.Object) {
	mc.m[client.ObjectKeyFromObject(obj)] = obj
	mc.added = append(mc.added, obj)
}
