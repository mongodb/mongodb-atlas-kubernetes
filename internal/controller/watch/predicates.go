package watch

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// CommonPredicates returns the predicate which filter out the changes done to any field except for spec (e.g. status)
// Also we should reconcile if finalizers have changed (see https://blog.openshift.com/kubernetes-operators-best-practices/)
// This will be phased out gradually to be replaced by DefaultPredicates
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

// GlobalResyncAwareGenerationChangePredicate reconcile on unfrequent global
// resyncs or on spec generation changes, but ignore finalizer changes
func GlobalResyncAwareGenerationChangePredicate() predicate.Predicate {
	return predicate.Or[client.Object](
		predicate.Not[client.Object](predicate.ResourceVersionChangedPredicate{}), // for the global resync
		predicate.GenerationChangedPredicate{},
	)
}

// IgnoreDeletedPredicate ignore after deletion handling, use unless some after
// deletion cleanup is needed
func IgnoreDeletedPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}

// DefaultPredicates avoid spurious after deletion or finalizer changes handling
func DefaultPredicates() predicate.Predicate {
	return predicate.And(
		GlobalResyncAwareGenerationChangePredicate(),
		IgnoreDeletedPredicate(),
	)
}
