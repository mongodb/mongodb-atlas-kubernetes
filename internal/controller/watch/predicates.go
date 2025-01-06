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
			if e.ObjectOld.GetResourceVersion() == e.ObjectNew.GetResourceVersion() {
				// resource version didn't change, so this is a resync, allow reconciliation.
				return true
			}

			if e.ObjectOld.GetGeneration() == e.ObjectNew.GetGeneration() && reflect.DeepEqual(e.ObjectNew.GetFinalizers(), e.ObjectOld.GetFinalizers()) {
				return false
			}
			return true
		},
	}
}

func SelectNamespacesPredicate(namespaces []string) predicate.Funcs {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		if len(namespaces) == 0 {
			return true
		}

		for _, ns := range namespaces {
			if object.GetNamespace() == ns {
				return true
			}
		}

		return false
	})
}
