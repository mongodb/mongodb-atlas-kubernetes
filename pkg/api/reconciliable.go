package api

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +k8s:deepcopy-gen=false

// Reconciliable is implemented by CRD objects used by indexes to trigger reconciliations
type Reconciliable interface {
	ReconciliableRequests() []reconcile.Request
}

// +k8s:deepcopy-gen=false

// ReconciliableList is a Reconciliable that is also a CRD list
type ReconciliableList interface {
	client.ObjectList
	Reconciliable
}

// ToRequest is a helper to turns CRD objects into reconcile requests.
// Most Reconciliable implementations may leverage it.
func ToRequest(obj client.Object) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		},
	}
}
