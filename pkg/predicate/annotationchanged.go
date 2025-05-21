package predicate

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func AnnotationChanged(key string) predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldAnn := e.ObjectOld.GetAnnotations()
			newAnn := e.ObjectNew.GetAnnotations()

			if oldAnn == nil && newAnn == nil {
				return false
			}

			oldValue, oldExists := oldAnn[key]
			newValue, newExists := newAnn[key]

			result := oldExists != newExists || oldValue != newValue
			return result
		},
	}
}
