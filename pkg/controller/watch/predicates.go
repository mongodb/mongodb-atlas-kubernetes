package watch

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// CommonPredicates returns the predicate which filter out the changes done to any field except for spec (e.g. status)
// Also we should reconcile if finalizers have changed (see https://blog.openshift.com/kubernetes-operators-best-practices/)
func CommonPredicates() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.MetaNew.GetGeneration() == e.MetaOld.GetGeneration() && reflect.DeepEqual(e.MetaNew.GetFinalizers(), e.MetaOld.GetFinalizers()) {
				return false
			}
			return true
		},
	}
}

// DeleteOnly returns a predicate that will filter out everything except the Delete event
func DeleteOnly() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(ce event.CreateEvent) bool {
			return false
		},
		UpdateFunc: func(ce event.UpdateEvent) bool {
			return false
		},
		GenericFunc: func(ce event.GenericEvent) bool {
			return false
		},
	}
}
