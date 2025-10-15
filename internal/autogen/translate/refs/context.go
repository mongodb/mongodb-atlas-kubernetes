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
