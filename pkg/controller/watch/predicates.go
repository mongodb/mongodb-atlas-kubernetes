package watch

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// CommonPredicates returns the predicate which filter out the changes done to any field except for spec (e.g. status)
// Also we should reconcile if finalizers have changed (see https://blog.openshift.com/kubernetes-operators-best-practices/)
func CommonPredicates() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld.GetGeneration() == e.ObjectNew.GetGeneration() && reflect.DeepEqual(e.ObjectNew.GetFinalizers(), e.ObjectOld.GetFinalizers()) {
				return false
			}
			return true
		},
	}
}

// CommonPredicatesWithAnnotations returns the predicate which filter out the changes done to any field except for spec (e.g. status)
// Also we should reconcile if finalizers have changed (see https://blog.openshift.com/kubernetes-operators-best-practices/)
// and we should reconcile if the annotations have changed to allow other controllers to trigger the reconciliation using annotations.
func CommonPredicatesWithAnnotations() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld.GetGeneration() == e.ObjectNew.GetGeneration() && reflect.DeepEqual(e.ObjectNew.GetFinalizers(), e.ObjectOld.GetFinalizers()) && reflect.DeepEqual(e.ObjectNew.GetAnnotations(), e.ObjectOld.GetAnnotations()) {
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

func SelectNamespacesPredicate(namespaceMap map[string]bool) predicate.Funcs {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		if _, ok := namespaceMap[""]; ok {
			return true
		}

		if _, ok := namespaceMap[object.GetNamespace()]; ok {
			return true
		}

		return false
	})
}
