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

import "sigs.k8s.io/controller-runtime/pkg/client"

// context points to the main Kubernetes object being translated,
// and holds related existing & added Kubernetes dependencies
type context struct {
	main  client.Object
	m     map[client.ObjectKey]client.Object
	added []client.Object
}

func newMapContext(main client.Object, deps []client.Object) *context {
	m := map[client.ObjectKey]client.Object{}
	for _, obj := range deps {
		m[client.ObjectKeyFromObject(obj)] = obj
	}
	return &context{main: main, m: m}
}

func (mc *context) find(name string) client.Object {
	key := client.ObjectKey{Name: name, Namespace: mc.main.GetNamespace()}
	return mc.m[key]
}

func (mc *context) has(name string) bool {
	return mc.find(name) != nil
}

func (mc *context) add(obj client.Object) {
	mc.m[client.ObjectKeyFromObject(obj)] = obj
	mc.added = append(mc.added, obj)
}
