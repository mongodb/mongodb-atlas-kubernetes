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

package watch

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// DeprecatedCommonPredicates returns the predicate which filter out the changes done to any field except for spec (e.g. status)
// Also we should reconcile if finalizers have changed (see https://blog.openshift.com/kubernetes-operators-best-practices/)
// This will be phased out gradually to be replaced by DefaultPredicates
func DeprecatedCommonPredicates() predicate.Funcs {
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
func GlobalResyncAwareGenerationChangePredicate[T metav1.Object]() predicate.TypedPredicate[T] {
	return predicate.Or[T](
		predicate.Not[T](predicate.TypedResourceVersionChangedPredicate[T]{}), // for the global resync
		predicate.TypedGenerationChangedPredicate[T]{},
	)
}

// IgnoreDeletedPredicate ignore after deletion handling, use unless some after
// deletion cleanup is needed
func IgnoreDeletedPredicate[T metav1.Object]() predicate.TypedPredicate[T] {
	return predicate.TypedFuncs[T]{
		DeleteFunc: func(e event.TypedDeleteEvent[T]) bool {
			return false
		},
	}
}

// DefaultPredicates avoid spurious after deletion or finalizer changes handling
func DefaultPredicates[T metav1.Object]() predicate.TypedPredicate[T] {
	return predicate.And[T](
		GlobalResyncAwareGenerationChangePredicate[T](),
		IgnoreDeletedPredicate[T](),
	)
}

type ReadyFunc[T any] func(obj T) bool

// ReadyTransitionPredicate filters out only those objects where the previous
// oldObject was not ready, but the new one is, or the object was deleted.
func ReadyTransitionPredicate[T any](ready ReadyFunc[T]) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectNew == nil || e.ObjectOld == nil {
				return false
			}

			newObj, ok := e.ObjectNew.(T)
			if !ok {
				return false
			}
			oldObj, ok := e.ObjectOld.(T)
			if !ok {
				return false
			}
			return !ready(oldObj) && ready(newObj)
		},
	}
}
