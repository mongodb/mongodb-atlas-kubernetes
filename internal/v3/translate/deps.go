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
//

package translate

import "sigs.k8s.io/controller-runtime/pkg/client"

const SetFallbackNamespace = "."

type DependencyRepo interface {
	MainObject() client.Object
	Find(name, namespace string) client.Object
	Add(obj client.Object)
	Added() []client.Object
}

type Dependencies struct {
	mainObj client.Object
	deps    map[string]client.Object
	added   []client.Object
}

// NewDependencies creates a set of Kubernetes client.Objects
func NewDependencies(mainObj client.Object, objs ...client.Object) *Dependencies {
	deps := map[string]client.Object{}
	for _, obj := range objs {
		deps[client.ObjectKeyFromObject(obj).String()] = obj
	}
	return &Dependencies{
		mainObj: mainObj,
		deps:    deps,
		added:   []client.Object{},
	}
}

// MainObject retried the main object for this dependecny repository
func (d *Dependencies) MainObject() client.Object {
	return d.mainObj
}

// Find looks for an object within the dependencies by name and namespace
func (d *Dependencies) Find(name, namespace string) client.Object {
	ns := namespace
	if ns == SetFallbackNamespace {
		ns = d.mainObj.GetNamespace()
	}
	return d.deps[client.ObjectKey{Name: name, Namespace: ns}.String()]
}

// Add appends an object to the added list and records it in the general set
func (d *Dependencies) Add(obj client.Object) {
	if obj.GetNamespace() == SetFallbackNamespace {
		obj.SetNamespace(d.mainObj.GetNamespace())
	}
	d.deps[client.ObjectKeyFromObject(obj).String()] = obj
	for i := range d.added {
		if d.added[i].GetName() == obj.GetName() && d.added[i].GetNamespace() == obj.GetNamespace() {
			d.added[i] = obj
			return
		}
	}
	d.added = append(d.added, obj)
}

// Added dumps an array of all dependencies added to the set after creation
func (d *Dependencies) Added() []client.Object {
	return d.added
}
